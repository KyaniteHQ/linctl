package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type fakeIssueWorkflow struct {
	template      client.IssueTemplateContent
	templateErr   error
	resolved      client.ResolvedTarget
	resolveErr    error
	listAll       client.IssueList
	listAllCalls  int
	listTeam      client.IssueList
	listTeamID    string
	listFilters   client.IssueListFilters
	listTeamCalls int
	searchList    client.IssueList
	searchTeamID  string
	searchQuery   string
	searchLimit   int
	searchCalls   int
	searchErr     error
	nextList      client.IssueList
	nextTeamID    string
	nextLimit     int
	nextCalls     int
	started       client.IssueSummary
	startID       string
	startErr      error
}

func (workflow *fakeIssueWorkflow) GetIssueTemplateContent(
	_ context.Context,
	_ string,
) (client.IssueTemplateContent, error) {
	return workflow.template, workflow.templateErr
}

func (workflow *fakeIssueWorkflow) ResolveTarget(_ context.Context) (client.ResolvedTarget, error) {
	return workflow.resolved, workflow.resolveErr
}

func (workflow *fakeIssueWorkflow) ListIssues(_ context.Context, _ int) (client.IssueList, error) {
	workflow.listAllCalls++

	return workflow.listAll, nil
}

func (workflow *fakeIssueWorkflow) ListIssuesByTeam(
	_ context.Context,
	teamID string,
	_ int,
	filters client.IssueListFilters,
) (client.IssueList, error) {
	workflow.listTeamCalls++
	workflow.listTeamID = teamID
	workflow.listFilters = filters

	return workflow.listTeam, nil
}

func (workflow *fakeIssueWorkflow) SearchIssuesByTeam(
	_ context.Context,
	teamID string,
	query string,
	limit int,
) (client.IssueList, error) {
	workflow.searchCalls++
	workflow.searchTeamID = teamID
	workflow.searchQuery = query
	workflow.searchLimit = limit

	return workflow.searchList, workflow.searchErr
}

func (workflow *fakeIssueWorkflow) ListNextIssuesByTeam(
	_ context.Context,
	teamID string,
	limit int,
) (client.IssueList, error) {
	workflow.nextCalls++
	workflow.nextTeamID = teamID
	workflow.nextLimit = limit

	return workflow.nextList, nil
}

func (workflow *fakeIssueWorkflow) StartIssue(
	_ context.Context,
	issueID string,
) (client.IssueSummary, error) {
	workflow.startID = issueID

	return workflow.started, workflow.startErr
}

func bufferedCommand() (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	command := &cobra.Command{}
	var stdout, stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)

	return command, &stdout, &stderr
}

func Test_assembleIssueCreate_applies_template_and_normalization(t *testing.T) {
	command, _, stderr := bufferedCommand()
	workflow := &fakeIssueWorkflow{
		template: client.IssueTemplateContent{Title: "Template title", Description: "Template body"},
	}
	estimate := 5

	request, err := assembleIssueCreate(
		context.Background(),
		command,
		workflow,
		client.IssueCreateRequest{},
		issueCreateFlags{templateID: "tmpl-1", state: "in progress", priority: "high"},
		&estimate,
	)

	require.NoError(t, err)
	require.Equal(t, "Template title", request.Title)
	require.Equal(t, "Template body", request.Description)
	require.Equal(t, "started", request.StateType)
	require.Equal(t, "2", request.Priority)
	require.NotNil(t, request.Estimate)
	require.Equal(t, 5, *request.Estimate)
	require.Contains(t, stderr.String(), "normalized")
}

func Test_assembleIssueCreate_stops_on_template_error(t *testing.T) {
	command, _, _ := bufferedCommand()
	workflow := &fakeIssueWorkflow{templateErr: errors.New("no such template")}

	_, err := assembleIssueCreate(
		context.Background(), command, workflow,
		client.IssueCreateRequest{Title: "X"}, issueCreateFlags{templateID: "missing"}, nil,
	)

	require.ErrorContains(t, err, "no such template")
}

func Test_issueList_lists_across_teams_without_resolving_target(t *testing.T) {
	workflow := &fakeIssueWorkflow{
		listAll:    client.IssueList{Issues: []client.IssueSummary{{Identifier: "LIT-1"}}},
		resolveErr: errors.New("must not resolve"),
	}

	list, err := issueList(context.Background(), workflow, 50, issueListFlagValues{allTeams: true})

	require.NoError(t, err)
	require.Equal(t, 1, workflow.listAllCalls)
	require.Equal(t, 0, workflow.listTeamCalls)
	require.Len(t, list.Issues, 1)
}

func Test_issueList_resolves_team_and_assembles_filters(t *testing.T) {
	workflow := &fakeIssueWorkflow{
		resolved: client.ResolvedTarget{
			Team:   client.TargetTeam{ID: "team-id"},
			Viewer: client.TargetViewer{ID: "viewer-id"},
		},
		listTeam: client.IssueList{Issues: []client.IssueSummary{{Identifier: "LIT-2"}}},
	}

	_, err := issueList(context.Background(), workflow, 50, issueListFlagValues{mine: true, stateType: "started"})

	require.NoError(t, err)
	require.Equal(t, 1, workflow.listTeamCalls)
	require.Equal(t, "team-id", workflow.listTeamID)
	require.Equal(t, "viewer-id", workflow.listFilters.AssigneeID)
	require.Equal(t, "started", workflow.listFilters.StateType)
}

func Test_runIssueSearch_uses_the_resolved_team(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	workflow := &fakeIssueWorkflow{
		resolved: client.ResolvedTarget{Team: client.TargetTeam{ID: "team-id"}},
		searchList: client.IssueList{Issues: []client.IssueSummary{
			clientIssue("LIT-13", "Search from workflow"),
		}},
	}

	err := runIssueSearch(context.Background(), command, &rootOptions{}, workflow, "needle", 7)

	require.NoError(t, err)
	require.Equal(t, 1, workflow.searchCalls)
	require.Equal(t, "team-id", workflow.searchTeamID)
	require.Equal(t, "needle", workflow.searchQuery)
	require.Equal(t, 7, workflow.searchLimit)
	require.Contains(t, stdout.String(), "LIT-13 Search from workflow")
}
