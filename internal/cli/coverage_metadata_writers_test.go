package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_CliRenderHelpers_write_issue_child_metadata_output(t *testing.T) {
	customerNeed := client.CustomerNeedMetadataSummary{
		ID:           "customer-need-id",
		CustomerName: "Acme",
		Issue:        "LIT-1",
		Priority:     2,
	}
	sharedAccess := client.IssueSharedAccessSummary{
		IssueID:                   "issue-id",
		Identifier:                "LIT-1",
		IsShared:                  true,
		ViewerHasOnlySharedAccess: false,
		SharedWithCount:           2,
		DisallowedIssueFields:     []string{"description", "priority"},
	}

	textOutput := bytes.Buffer{}
	textCommand := &cobra.Command{}
	textCommand.SetOut(&textOutput)
	require.NoError(t, writeCustomerNeedMetadata(textCommand, &rootOptions{}, customerNeed))
	require.NoError(t, writeIssueSharedAccess(textCommand, &rootOptions{}, sharedAccess))
	require.Contains(t, textOutput.String(), "customer-need-id Acme LIT-1 priority 2")
	require.Contains(t, textOutput.String(), "issue-id LIT-1 shared=true")
	require.Contains(t, textOutput.String(), "disallowed=description,priority")

	jsonOutput := bytes.Buffer{}
	jsonCommand := &cobra.Command{}
	jsonCommand.SetOut(&jsonOutput)
	require.NoError(t, writeCustomerNeedMetadata(jsonCommand, &rootOptions{json: true}, customerNeed))
	require.NoError(t, writeIssueSharedAccess(jsonCommand, &rootOptions{json: true}, sharedAccess))
	require.Contains(t, jsonOutput.String(), `"customer_name": "Acme"`)
	require.Contains(t, jsonOutput.String(), `"shared_with_count": 2`)

	quietOutput := bytes.Buffer{}
	quietCommand := &cobra.Command{}
	quietCommand.SetOut(&quietOutput)
	require.NoError(t, writeCustomerNeedMetadata(quietCommand, &rootOptions{quiet: true}, customerNeed))
	require.NoError(t, writeIssueSharedAccess(quietCommand, &rootOptions{quiet: true}, sharedAccess))
	require.Empty(t, quietOutput.String())

	idOnlyOutput := bytes.Buffer{}
	idOnlyCommand := &cobra.Command{}
	idOnlyCommand.SetOut(&idOnlyOutput)
	require.NoError(t, writeCustomerNeedMetadata(idOnlyCommand, &rootOptions{idOnly: true}, customerNeed))
	require.Equal(t, "customer-need-id\n", idOnlyOutput.String())

	emptyFieldsOutput := bytes.Buffer{}
	emptyFieldsCommand := &cobra.Command{}
	emptyFieldsCommand.SetOut(&emptyFieldsOutput)
	sharedAccess.DisallowedIssueFields = nil
	require.NoError(t, writeIssueSharedAccess(emptyFieldsCommand, &rootOptions{}, sharedAccess))
	require.Contains(t, emptyFieldsOutput.String(), "disallowed=-")
}

func Test_CliRenderHelpers_write_project_filter_suggestion_output(t *testing.T) {
	suggestion := client.ProjectFilterSuggestion{
		Filter: json.RawMessage(`{"status":{"type":{"eq":"started"}}}`),
		LogID:  "filter-log-id",
	}

	textOutput := bytes.Buffer{}
	textCommand := &cobra.Command{}
	textCommand.SetOut(&textOutput)
	require.NoError(t, writeProjectFilterSuggestion(textCommand, &rootOptions{}, suggestion))
	require.Contains(t, textOutput.String(), `log_id=filter-log-id`)
	require.Contains(t, textOutput.String(), `filter={"status":{"type":{"eq":"started"}}}`)

	jsonOutput := bytes.Buffer{}
	jsonCommand := &cobra.Command{}
	jsonCommand.SetOut(&jsonOutput)
	require.NoError(t, writeProjectFilterSuggestion(jsonCommand, &rootOptions{json: true}, suggestion))
	require.Contains(t, jsonOutput.String(), `"log_id": "filter-log-id"`)

	quietOutput := bytes.Buffer{}
	quietCommand := &cobra.Command{}
	quietCommand.SetOut(&quietOutput)
	require.NoError(t, writeProjectFilterSuggestion(quietCommand, &rootOptions{quiet: true}, suggestion))
	require.Empty(t, quietOutput.String())
}

func Test_CliRenderHelpers_write_issue_utility_output(t *testing.T) {
	priorityValues := []client.IssuePriorityValue{{Priority: 1, Label: "Urgent"}}
	filterSuggestion := client.IssueFilterSuggestion{
		Filter: json.RawMessage(`{"state":{"type":{"eq":"started"}}}`),
		LogID:  "issue-filter-log-id",
	}
	titleSuggestion := client.IssueTitleSuggestion{Title: "Improve exports", LogID: "title-log-id"}

	textOutput := bytes.Buffer{}
	textCommand := &cobra.Command{}
	textCommand.SetOut(&textOutput)
	require.NoError(t, writeIssuePriorityValues(textCommand, &rootOptions{}, priorityValues))
	require.NoError(t, writeIssueFilterSuggestion(textCommand, &rootOptions{}, filterSuggestion))
	require.NoError(t, writeIssueTitleSuggestion(textCommand, &rootOptions{}, titleSuggestion))
	require.Contains(t, textOutput.String(), "1 Urgent")
	require.Contains(t, textOutput.String(), `log_id=issue-filter-log-id`)
	require.Contains(t, textOutput.String(), `log_id=title-log-id title=Improve exports`)

	jsonOutput := bytes.Buffer{}
	jsonCommand := &cobra.Command{}
	jsonCommand.SetOut(&jsonOutput)
	require.NoError(t, writeIssuePriorityValues(jsonCommand, &rootOptions{json: true}, priorityValues))
	require.NoError(t, writeIssueFilterSuggestion(jsonCommand, &rootOptions{json: true}, filterSuggestion))
	require.NoError(t, writeIssueTitleSuggestion(jsonCommand, &rootOptions{json: true}, titleSuggestion))
	require.Contains(t, jsonOutput.String(), `"label": "Urgent"`)
	require.Contains(t, jsonOutput.String(), `"log_id": "issue-filter-log-id"`)
	require.Contains(t, jsonOutput.String(), `"title": "Improve exports"`)

	errorCommand := &cobra.Command{}
	errorCommand.SetOut(commandFailingWriter{})
	require.Error(t, writeIssuePriorityValues(errorCommand, &rootOptions{}, priorityValues))
	require.Error(t, writeIssueFilterSuggestion(errorCommand, &rootOptions{}, filterSuggestion))
	require.Error(t, writeIssueTitleSuggestion(errorCommand, &rootOptions{}, titleSuggestion))
	require.Error(t, writeProjectStatusProjectCount(
		errorCommand,
		&rootOptions{},
		client.ProjectStatusProjectCount{ProjectStatusID: "project-status-id", Count: 12},
	))
}
