package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// fakeBulkCreator is an in-memory bulkIssueCreator: it records each request and
// can fail at a chosen row, so the import create loop's accumulation and
// per-row error wrapping are tested without canned GraphQL JSON or a real
// batch target resolution.
type fakeBulkCreator struct {
	results        []client.IssueSummary
	failAt         int // 1-based row to fail at; 0 never fails
	failErr        error
	calls          int
	requests       []client.IssueCreateRequest
	sawConcurrency int
}

func (creator *fakeBulkCreator) CreateIssues(
	_ context.Context,
	requests []client.IssueCreateRequest,
	concurrency int,
) ([]client.IssueCreateOutcome, error) {
	creator.sawConcurrency = concurrency
	outcomes := make([]client.IssueCreateOutcome, len(requests))
	for index, request := range requests {
		creator.calls++
		creator.requests = append(creator.requests, request)
		outcome := client.IssueCreateOutcome{Index: index}
		if creator.failAt == creator.calls {
			outcome.Err = creator.failErr
		} else if index < len(creator.results) {
			outcome.Issue = creator.results[index]
		}
		outcomes[index] = outcome
	}

	return outcomes, nil
}

func Test_createImportedIssues_creates_each_row_through_the_port(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	creator := &fakeBulkCreator{
		results: []client.IssueSummary{
			{Identifier: "LIT-1", Title: "First", State: "Todo"},
			{Identifier: "LIT-2", Title: "Second", State: "Todo"},
		},
	}
	requests := []client.IssueCreateRequest{{Title: "First"}, {Title: "Second"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{}, creator, requests)

	require.NoError(t, err)
	require.Equal(t, 2, creator.calls)
	require.Equal(t, "First", creator.requests[0].Title)
	require.Equal(t, "Second", creator.requests[1].Title)
	require.Contains(t, stdout.String(), "LIT-1")
	require.Contains(t, stdout.String(), "LIT-2")
}

func Test_createImportedIssues_wraps_the_failing_row(t *testing.T) {
	command, _, _ := bufferedCommand()
	creator := &fakeBulkCreator{
		results: []client.IssueSummary{{Identifier: "LIT-1"}},
		failAt:  2,
		failErr: errors.New("boom"),
	}
	requests := []client.IssueCreateRequest{{Title: "First"}, {Title: "Second"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{}, creator, requests)

	require.ErrorContains(t, err, "import row 2")
	require.ErrorContains(t, err, "Second")
	require.ErrorContains(t, err, "boom")
	require.Equal(t, 2, creator.calls)
}

func Test_createImportedIssues_reports_partial_success_in_json_mode(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	creator := &fakeBulkCreator{
		results: []client.IssueSummary{{Identifier: "LIT-1", Title: "First", State: "Todo"}},
		failAt:  2,
		failErr: errors.New("boom"),
	}
	requests := []client.IssueCreateRequest{{Title: "First"}, {Title: "Second"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{json: true}, creator, requests)

	require.Error(t, err)
	var result issueImportResult
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &result))
	require.Equal(t, 1, result.Count)
	require.Len(t, result.Issues, 1)
	require.Equal(t, "LIT-1", result.Issues[0].Identifier)
	require.Len(t, result.Failures, 1)
	require.Equal(t, 2, result.Failures[0].Row)
	require.Equal(t, "Second", result.Failures[0].Title)
	require.Contains(t, result.Failures[0].Error, "boom")
}

func Test_createImportedIssues_json_success_output_has_no_failures_field(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	creator := &fakeBulkCreator{
		results: []client.IssueSummary{{Identifier: "LIT-1", Title: "First", State: "Todo"}},
	}
	requests := []client.IssueCreateRequest{{Title: "First"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{json: true}, creator, requests)

	require.NoError(t, err)
	require.NotContains(t, stdout.String(), "failures")
}

func Test_createImportedIssues_quiet_mode_still_reports_failures_on_stderr(t *testing.T) {
	command, stdout, stderr := bufferedCommand()
	creator := &fakeBulkCreator{
		results: []client.IssueSummary{{Identifier: "LIT-1", Title: "First", State: "Todo"}},
		failAt:  2,
		failErr: errors.New("boom"),
	}
	requests := []client.IssueCreateRequest{{Title: "First"}, {Title: "Second"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{quiet: true}, creator, requests)

	require.Error(t, err)
	require.Empty(t, stdout.String())
	require.Contains(t, stderr.String(), "row 2")
	require.Contains(t, stderr.String(), "Second")
	require.Contains(t, stderr.String(), "boom")
}

func Test_createImportedIssues_surfaces_stdout_write_errors(t *testing.T) {
	command := &cobra.Command{}
	command.SetOut(failingWriter{err: errors.New("stdout closed")})
	command.SetErr(&bytes.Buffer{})
	creator := &fakeBulkCreator{
		results: []client.IssueSummary{{Identifier: "LIT-1", Title: "First", State: "Todo"}},
	}
	requests := []client.IssueCreateRequest{{Title: "First"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{}, creator, requests)

	require.ErrorContains(t, err, "stdout closed")
}

func Test_createImportedIssues_surfaces_stderr_write_errors_for_failures(t *testing.T) {
	command := &cobra.Command{}
	command.SetOut(&bytes.Buffer{})
	command.SetErr(failingWriter{err: errors.New("stderr closed")})
	creator := &fakeBulkCreator{
		failAt:  1,
		failErr: errors.New("boom"),
	}
	requests := []client.IssueCreateRequest{{Title: "First"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{}, creator, requests)

	require.ErrorContains(t, err, "stderr closed")
}

func Test_createImportedIssues_partial_report_with_batch_creator(t *testing.T) {
	command, stdout, stderr := bufferedCommand()
	creator := &fakeBulkCreator{
		results: []client.IssueSummary{
			{Identifier: "LIT-1", Title: "First", State: "Todo"},
			{},
			{Identifier: "LIT-3", Title: "Third", State: "Todo"},
		},
		failAt:  2,
		failErr: errors.New("boom"),
	}
	requests := []client.IssueCreateRequest{{Title: "First"}, {Title: "Second"}, {Title: "Third"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{}, creator, requests)

	require.ErrorContains(t, err, "import completed with 1 failed rows")
	require.Equal(t, 3, creator.calls, "all rows should go through one batch call, not one call per row")
	require.Equal(t, importBatchConcurrency, creator.sawConcurrency)
	require.Contains(t, stdout.String(), "LIT-1")
	require.Contains(t, stdout.String(), "LIT-3")
	require.Contains(t, stderr.String(), "row 2")
}

func Test_createImportedIssues_surfaces_batch_resolution_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	creator := &failingBatchCreator{err: client.ErrTargetMismatch}
	requests := []client.IssueCreateRequest{{Title: "First"}}

	err := createImportedIssues(context.Background(), command, &rootOptions{}, creator, requests)

	require.ErrorIs(t, err, client.ErrTargetMismatch)
}

// failingBatchCreator is a bulkIssueCreator whose CreateIssues call itself
// fails (e.g. target resolution), rather than an individual row.
type failingBatchCreator struct {
	err error
}

func (creator *failingBatchCreator) CreateIssues(
	_ context.Context,
	_ []client.IssueCreateRequest,
	_ int,
) ([]client.IssueCreateOutcome, error) {
	return nil, creator.err
}
