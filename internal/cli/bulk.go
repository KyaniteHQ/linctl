package cli

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

const (
	formatJSON             = "json"
	formatCSV              = "csv"
	defaultBulkExportLimit = 250
)

// issueImportRow is one issue parsed from a CSV or JSON import file.
type issueImportRow struct {
	Team        string `json:"team,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	State       string `json:"state,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

// issueImportPlan is the normalized create plan shown by an import dry run.
type issueImportPlan struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	StateType   string `json:"state_type,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

// issueImportPreview is the structured dry-run output of an import.
type issueImportPreview struct {
	Count  int               `json:"count"`
	DryRun bool              `json:"dry_run"`
	Issues []issueImportPlan `json:"issues"`
}

// issueImportResult is the structured confirmation of a completed import.
type issueImportResult struct {
	Count  int                   `json:"count"`
	Issues []client.IssueSummary `json:"issues"`
}

// bulkExportResult is the structured confirmation of a completed bulk export.
type bulkExportResult struct {
	Path  string `json:"path"`
	Count int    `json:"count"`
}

// dataFormat resolves the import/export encoding from a file extension.
func dataFormat(path string) (string, error) {
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".json":
		return formatJSON, nil
	case ".csv":
		return formatCSV, nil
	default:
		return "", fmt.Errorf("unsupported data format %q: use .json or .csv", ext)
	}
}

func addIssueImportCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	dryRun := false
	command := &cobra.Command{
		Use:   "import FILE",
		Short: "Create issues from a CSV or JSON file in the pinned target",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runIssueImport(ctx, command, options, args[0], dryRun)
		},
	}
	command.Flags().BoolVar(&dryRun, "dry-run", false, "render the rows that would be created without writing")
	root.AddCommand(command)
}

func runIssueImport(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	path string,
	dryRun bool,
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	format, err := dataFormat(path)
	if err != nil {
		return err
	}
	rows, err := parseImportFile(format, path)
	if err != nil {
		return err
	}
	requests, err := buildImportRequests(rows, runtime.config.Target.TeamKey)
	if err != nil {
		return err
	}
	if dryRun {
		return writeImportPreview(command, options, requests)
	}

	return createImportedIssues(ctx, command, options, issueAdapterFor(runtime), requests)
}

// bulkIssueCreator is the narrow Command Port the import create loop depends on.
// Defined by its consumer (it needs only CreateIssue, not the wider issue port)
// and satisfied in production by issueClientAdapter, so the per-row error
// wrapping is tested against an in-memory fake rather than canned GraphQL JSON.
type bulkIssueCreator interface {
	CreateIssue(ctx context.Context, request client.IssueCreateRequest) (client.IssueSummary, error)
}

// The shared production adapter satisfies the narrow bulk port; this assertion
// fails the build if CreateIssue's shape drifts, keeping the write-guard
// forwarding intact rather than letting an adapter quietly stop satisfying it.
var _ bulkIssueCreator = issueClientAdapter{}

func createImportedIssues(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	creator bulkIssueCreator,
	requests []client.IssueCreateRequest,
) error {
	// --quiet renders nothing, so skip accumulating the created summaries; the
	// default, --json, and --id-only modes still consume them downstream.
	var created []client.IssueSummary
	if !options.quiet {
		created = make([]client.IssueSummary, 0, len(requests))
	}
	for index, request := range requests {
		issue, err := creator.CreateIssue(ctx, request)
		if err != nil {
			return fmt.Errorf("import row %d %q: %w", index+1, request.Title, err)
		}
		if !options.quiet {
			created = append(created, issue)
		}
	}

	return writeImportResult(command, options, created)
}

func parseImportFile(format string, path string) ([]issueImportRow, error) {
	//nolint:gosec // G304: the import command's purpose is to read the user-named file.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	// Strip a leading UTF-8 BOM so spreadsheet-exported files parse: a BOM would
	// otherwise glue onto the first CSV header cell or break json.Unmarshal.
	data = bytes.TrimPrefix(data, []byte("\ufeff"))
	if format == formatCSV {
		return parseCSVRows(data)
	}

	return parseJSONRows(data)
}

func parseJSONRows(data []byte) ([]issueImportRow, error) {
	rows := []issueImportRow{}
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, fmt.Errorf("parse json import: %w", err)
	}

	return rows, nil
}

func parseCSVRows(data []byte) ([]issueImportRow, error) {
	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse csv import: %w", err)
	}
	if len(records) == 0 {
		return nil, errors.New("parse csv import: file is empty")
	}
	header := csvHeaderIndex(records[0])
	rows := make([]issueImportRow, 0, len(records)-1)
	for _, record := range records[1:] {
		rows = append(rows, issueImportRow{
			Team:        csvField(record, header, "team"),
			Title:       csvField(record, header, "title"),
			Description: csvField(record, header, "description"),
			State:       csvField(record, header, "state"),
			Priority:    csvField(record, header, "priority"),
		})
	}

	return rows, nil
}

func csvHeaderIndex(header []string) map[string]int {
	index := make(map[string]int, len(header))
	for position, name := range header {
		index[strings.ToLower(strings.TrimSpace(name))] = position
	}

	return index
}

// csvField returns the trimmed value for a named column. csv.Reader guarantees
// every record has the same field count as the header, so a resolved header
// index is always in range.
func csvField(record []string, header map[string]int, name string) string {
	if position, ok := header[name]; ok {
		return strings.TrimSpace(record[position])
	}

	return ""
}

func buildImportRequests(rows []issueImportRow, pinnedTeamKey string) ([]client.IssueCreateRequest, error) {
	requests := make([]client.IssueCreateRequest, 0, len(rows))
	for index, row := range rows {
		request, err := importRowToRequest(row, pinnedTeamKey)
		if err != nil {
			return nil, fmt.Errorf("import row %d: %w", index+1, err)
		}
		requests = append(requests, request)
	}

	return requests, nil
}

// importRowToRequest validates one row against the pinned target and normalizes
// its state and priority before it becomes a guarded create request.
func importRowToRequest(row issueImportRow, pinnedTeamKey string) (client.IssueCreateRequest, error) {
	if strings.TrimSpace(row.Title) == "" {
		return client.IssueCreateRequest{}, errors.New("title is required")
	}
	if team := strings.TrimSpace(row.Team); team != "" && team != strings.TrimSpace(pinnedTeamKey) {
		return client.IssueCreateRequest{}, fmt.Errorf(
			"team %q does not match pinned target team %q", team, pinnedTeamKey,
		)
	}
	stateType, err := normalizeOptional(row.State, normalizedStateType)
	if err != nil {
		return client.IssueCreateRequest{}, err
	}
	priority, err := normalizeOptional(row.Priority, normalizedPriorityValue)
	if err != nil {
		return client.IssueCreateRequest{}, err
	}

	return client.IssueCreateRequest{
		Title:       row.Title,
		Description: row.Description,
		StateType:   stateType,
		Priority:    priority,
	}, nil
}

// normalizeOptional normalizes raw with normalize, leaving an empty value
// untouched. Unlike normalizeAndNote it emits no stderr note, so a bulk import
// stays quiet on each normalized row.
func normalizeOptional(raw string, normalize func(string) (string, bool, error)) (string, error) {
	if raw == "" {
		return "", nil
	}
	value, _, err := normalize(raw)

	return value, err
}

func writeImportPreview(command *cobra.Command, options *rootOptions, requests []client.IssueCreateRequest) error {
	if options.quiet {
		return nil
	}
	plans := importPlans(requests)
	if options.json {
		return writeJSONValue(command, options, issueImportPreview{
			Count:  len(plans),
			DryRun: true,
			Issues: plans,
		})
	}
	for _, plan := range plans {
		if err := render.WriteLine(
			command.OutOrStdout(),
			"would create %q state=%s priority=%s",
			plan.Title, emptyDash(plan.StateType), emptyDash(plan.Priority),
		); err != nil {
			return err
		}
	}

	return nil
}

func importPlans(requests []client.IssueCreateRequest) []issueImportPlan {
	plans := make([]issueImportPlan, 0, len(requests))
	for _, request := range requests {
		plans = append(plans, issueImportPlan{
			Title:       request.Title,
			Description: request.Description,
			StateType:   request.StateType,
			Priority:    request.Priority,
		})
	}

	return plans
}

func writeImportResult(command *cobra.Command, options *rootOptions, created []client.IssueSummary) error {
	if options.json {
		return writeJSONValue(command, options, issueImportResult{Count: len(created), Issues: created})
	}

	return writeIssues(command, options, created)
}

func addIssueBulkExportCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := defaultBulkExportLimit
	command := &cobra.Command{
		Use:   "bulk-export FILE",
		Short: "Write the resolved team's issues to a CSV or JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runIssueBulkExport(ctx, command, options, args[0], limit)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to export")
	root.AddCommand(command)
}

func runIssueBulkExport(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	path string,
	limit int,
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	format, err := dataFormat(path)
	if err != nil {
		return err
	}
	target, err := runtime.resolveTarget(ctx)
	if err != nil {
		return err
	}
	issues, err := client.ListIssuesByTeam(ctx, runtime.graphqlClient, target.Team.ID, limit, client.IssueListFilters{})
	if err != nil {
		return err
	}
	if err := writeIssueFile(path, format, issues.Issues); err != nil {
		return err
	}

	return writeBulkExportResult(command, options, path, len(issues.Issues))
}

func writeIssueFile(path string, format string, issues []client.IssueSummary) error {
	//nolint:gosec // G304: the export command's purpose is to write the user-named file.
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %s: %w", path, err)
	}
	defer closeQuietly(file)

	return encodeIssues(file, format, issues)
}

func encodeIssues(writer io.Writer, format string, issues []client.IssueSummary) error {
	if format == formatCSV {
		return encodeIssuesCSV(writer, issues)
	}

	return render.WriteJSON(writer, issues, false)
}

func encodeIssuesCSV(writer io.Writer, issues []client.IssueSummary) error {
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.WriteAll(issuesToCSVRecords(issues)); err != nil {
		return fmt.Errorf("encode csv export: %w", err)
	}

	return nil
}

func issuesToCSVRecords(issues []client.IssueSummary) [][]string {
	records := make([][]string, 0, 1+len(issues))
	records = append(records, []string{"identifier", "title", "state", "priority", "assignee", "project", "url"})
	for _, issue := range issues {
		records = append(records, []string{
			issue.Identifier,
			issue.Title,
			issue.State,
			issue.PriorityLabel,
			issue.Assignee,
			issue.Project,
			issue.URL,
		})
	}

	return records
}

func writeBulkExportResult(command *cobra.Command, options *rootOptions, path string, count int) error {
	return writeItem(command, options, bulkExportResult{Path: path, Count: count}, path,
		func(command *cobra.Command, _ *rootOptions, result bulkExportResult) error {
			return render.WriteLine(command.OutOrStdout(), "%s (%d issues)", result.Path, result.Count)
		})
}
