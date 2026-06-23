package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

// errorEnvelope is the machine-readable failure shape emitted to stderr so
// agents can branch on a stable error_code instead of parsing prose.
type errorEnvelope struct {
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// errorCode maps a command error to a stable machine error_code, preferring the
// client sentinels (matched through wrapping) over the not-found heuristic.
func errorCode(err error) string {
	switch {
	case errors.Is(err, client.ErrTargetMismatch):
		return "TARGET_MISMATCH"
	case errors.Is(err, client.ErrRateLimited):
		return "RATE_LIMITED"
	case errors.Is(err, client.ErrMutationFailed):
		return "MUTATION_FAILED"
	case errors.Is(err, client.ErrWriteInvalid):
		return "INVALID_WRITE"
	case errors.Is(err, client.ErrGraphQL):
		return "GRAPHQL_ERROR"
	case errors.Is(err, client.ErrNotFound):
		return "NOT_FOUND"
	default:
		return "INTERNAL"
	}
}

// writeErrorEnvelope emits the structured error envelope (one JSON line) to the
// given writer. json.Encoder cannot fail to marshal these two strings, so the
// only error it returns is a write failure.
func writeErrorEnvelope(writer io.Writer, err error) error {
	return json.NewEncoder(writer).Encode(errorEnvelope{
		ErrorCode: errorCode(err),
		Message:   err.Error(),
	})
}

func writeJSONValue(command *cobra.Command, options *rootOptions, value any) error {
	if options.quiet {
		return nil
	}
	projected, err := projectJSONFields(value, options.fields)
	if err != nil {
		return err
	}

	return render.WriteJSON(command.OutOrStdout(), projected, options.compact)
}

func writeIDOnly(command *cobra.Command, options *rootOptions, id string) (bool, error) {
	if !options.idOnly {
		return false, nil
	}
	if options.quiet {
		return true, nil
	}
	if id == "" {
		return true, errors.New("id-only output: id is empty")
	}

	return true, render.WriteLine(command.OutOrStdout(), "%s", id)
}

// deletionResult is the structured confirmation returned by guarded delete commands.
type deletionResult struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// writeDeletion renders the confirmation for a guarded delete across all output
// modes: id-only, quiet, JSON envelope, then a plain "<id> deleted" line.
func writeDeletion(command *cobra.Command, options *rootOptions, id string) error {
	if wrote, err := writeIDOnly(command, options, id); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, deletionResult{ID: id, Status: "deleted"})
	}

	return render.WriteLine(command.OutOrStdout(), "%s deleted", id)
}

func ensureNonEmpty(options *rootOptions, count int) error {
	if options.failOnEmpty && count == 0 {
		return errors.New("empty result")
	}

	return nil
}

func sortByJSONField[T any](items []T, field string, order string) ([]T, error) {
	if field == "" {
		return items, nil
	}
	if order != "asc" && order != "desc" {
		return nil, fmt.Errorf("invalid sort order %q: use asc or desc", order)
	}

	sortedItems := slices.Clone(items)
	sort.SliceStable(sortedItems, func(leftIndex int, rightIndex int) bool {
		leftValue, leftErr := jsonFieldValue(sortedItems[leftIndex], field)
		rightValue, rightErr := jsonFieldValue(sortedItems[rightIndex], field)
		if leftErr != nil || rightErr != nil {
			return false
		}
		if order == "desc" {
			return rightValue < leftValue
		}

		return leftValue < rightValue
	})

	for _, item := range sortedItems {
		if _, err := jsonFieldValue(item, field); err != nil {
			return nil, err
		}
	}

	return sortedItems, nil
}

func jsonFieldValue(value any, field string) (string, error) {
	raw, err := jsonRoundTrip(value)
	if err != nil {
		return "", err
	}

	current := any(raw)
	for _, part := range strings.Split(field, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return "", fmt.Errorf("sort field %q is not an object path", field)
		}
		next, ok := object[part]
		if !ok {
			return "", fmt.Errorf("sort field %q is not present", field)
		}
		current = next
	}

	return fmt.Sprint(current), nil
}

func normalizedHumanFormat(options *rootOptions) (string, error) {
	switch options.format {
	case "":
		return "compact", nil
	case "minimal", "compact", "full":
		return options.format, nil
	default:
		return "", fmt.Errorf("invalid format %q: use minimal, compact, or full", options.format)
	}
}

func projectJSONFields(value any, fields string) (any, error) {
	paths := fieldPaths(fields)
	if len(paths) == 0 {
		return value, nil
	}

	raw, err := jsonRoundTrip(value)
	if err != nil {
		return nil, err
	}

	if projected, ok, err := projectCollection(raw, paths); ok || err != nil {
		return projected, err
	}

	projected := map[string]any{}
	for _, path := range paths {
		if err := copyJSONPath(raw, projected, path); err != nil {
			return nil, err
		}
	}

	return projected, nil
}

func fieldPaths(fields string) [][]string {
	if strings.TrimSpace(fields) == "" {
		return nil
	}

	parts := strings.Split(fields, ",")
	paths := make([][]string, 0, len(parts))
	for _, part := range parts {
		field := strings.TrimSpace(part)
		if field == "" {
			continue
		}
		paths = append(paths, strings.Split(field, "."))
	}

	return paths
}

func jsonRoundTrip(value any) (map[string]any, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("project json fields: marshal output: %w", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, fmt.Errorf("project json fields: decode output: %w", err)
	}

	return raw, nil
}

// collectionKeys is the explicit allowlist of collection field names that list
// pages emit (each list page carries exactly one such array plus scalar
// pagination/context fields). It is deliberately NOT generic top-level []any
// detection: some detail responses embed an incidental array that is not a
// collection (for example a TimeScheduleSummary's "entries", or a
// ProjectSummary's "teams"), so treating "the single top-level array" as the
// collection would wrongly project per-element instead of over the object.
// Equally, list pages are not all paginated (AuditEntryTypeList,
// SemanticSearchList, SLAConfigurationList, TemplateList carry no
// has_next_page), so a pagination marker cannot stand in for the allowlist
// either. Multi-array responses (IssueDependencyGraph) and detail objects fall
// through to whole-object projection. When adding a new list command, add its
// collection key here.
var collectionKeys = []string{
	"issues",
	"associations",
	"cycles",
	"projects",
	"members",
	"comments",
	"updates",
	"milestones",
	"documents",
	"labels",
	"teams",
	"users",
	"memberships",
	"drafts",
	"initiatives",
	"notifications",
	"notification_subscriptions",
	"release_pipelines",
	"release_stages",
	"releases",
	"history",
	"links",
	"release_notes",
	"customers",
	"customer_needs",
	"customer_statuses",
	"customer_tiers",
	"relations",
	"roadmaps",
	"time_schedules",
	"triage_responsibilities",
	"sla_configurations",
	"results",
	"templates",
	"workflow_states",
	"agent_activities",
	"agent_skills",
	"external_users",
	"audit_entry_types",
	"favorites",
	"emojis",
	"attachments",
	"custom_views",
	"project_labels",
	"project_statuses",
	"spans",
	"git_automation_states",
}

// projectCollection projects --fields over the items of a list-page envelope.
func projectCollection(raw map[string]any, paths [][]string) (map[string]any, bool, error) {
	for _, key := range collectionKeys {
		items, ok := raw[key].([]any)
		if !ok {
			continue
		}

		projectedItems := make([]any, 0, len(items))
		for _, item := range items {
			itemMap, ok := item.(map[string]any)
			if !ok {
				return nil, true, fmt.Errorf("project json fields: %s item is not an object", key)
			}
			projectedItem := map[string]any{}
			for _, path := range paths {
				if err := copyJSONPath(itemMap, projectedItem, path); err != nil {
					return nil, true, err
				}
			}
			projectedItems = append(projectedItems, projectedItem)
		}

		return map[string]any{key: projectedItems}, true, nil
	}

	return nil, false, nil
}

func copyJSONPath(source map[string]any, destination map[string]any, path []string) error {
	if len(path) == 0 {
		return nil
	}

	value, ok := source[path[0]]
	if !ok {
		return fmt.Errorf("project json fields: field %q is not present", strings.Join(path, "."))
	}
	if len(path) == 1 {
		destination[path[0]] = value
		return nil
	}

	sourceChild, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("project json fields: field %q is not an object", path[0])
	}
	destinationChild, ok := destination[path[0]].(map[string]any)
	if !ok {
		destinationChild = map[string]any{}
		destination[path[0]] = destinationChild
	}

	return copyJSONPath(sourceChild, destinationChild, path[1:])
}
