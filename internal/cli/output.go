package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/render"
)

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

func projectCollection(raw map[string]any, paths [][]string) (map[string]any, bool, error) {
	for _, key := range []string{
		"issues",
		"projects",
		"members",
		"comments",
		"updates",
		"milestones",
		"documents",
		"labels",
		"teams",
		"users",
		"notifications",
		"notification_subscriptions",
		"customers",
		"customer_needs",
		"customer_statuses",
		"customer_tiers",
		"roadmaps",
		"time_schedules",
	} {
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
