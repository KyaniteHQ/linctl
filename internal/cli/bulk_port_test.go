package cli

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// fakeBulkCreator is an in-memory bulkIssueCreator: it records each request and
// can fail at a chosen row, so the import create loop's accumulation and
// per-row error wrapping are tested without canned GraphQL JSON.
type fakeBulkCreator struct {
	results  []client.IssueSummary
	failAt   int // 1-based row to fail at; 0 never fails
	failErr  error
	calls    int
	requests []client.IssueCreateRequest
}

func (creator *fakeBulkCreator) CreateIssue(
	_ context.Context,
	request client.IssueCreateRequest,
) (client.IssueSummary, error) {
	creator.calls++
	creator.requests = append(creator.requests, request)
	if creator.failAt == creator.calls {
		return client.IssueSummary{}, creator.failErr
	}
	if index := creator.calls - 1; index < len(creator.results) {
		return creator.results[index], nil
	}

	return client.IssueSummary{}, nil
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
