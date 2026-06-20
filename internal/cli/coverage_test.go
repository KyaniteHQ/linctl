package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CliRenderHelpers_write_text_and_json_output(t *testing.T) {
	issue := client.IssueSummary{
		Identifier: "LIT-1",
		Title:      "Ship coverage",
		State:      "Todo",
		URL:        "https://linear.app/issue/LIT-1",
	}
	project := client.ProjectSummary{
		ID:   "project-id",
		Name: "Coverage",
		URL:  "https://linear.app/project/project-id",
		Status: client.ProjectStatus{
			Name: "Backlog",
		},
	}
	projectUpdate := client.ProjectUpdateSummary{
		ID:          "project-update-id",
		Body:        "First update",
		Health:      "onTrack",
		DisplayName: "Omer",
	}
	cycle := client.CycleSummary{
		ID:       "cycle-id",
		Name:     "Planning cycle",
		Status:   "active",
		StartsAt: "2026-07-01T00:00:00Z",
		EndsAt:   "2026-07-15T00:00:00Z",
		Progress: 0.5,
	}
	milestone := client.ProjectMilestoneSummary{
		ID:         "project-milestone-id",
		Name:       "Launch milestone",
		TargetDate: "2026-06-30",
		Status:     "next",
		Progress:   0.5,
	}
	document := client.DocumentSummary{
		ID:         "document-id",
		Title:      "Spec",
		ParentType: "project",
	}
	label := client.LabelSummary{
		ID:    "label-id",
		Name:  "Bug",
		Color: "#ff0000",
	}
	team := client.TeamSummary{
		ID:   "team-id",
		Key:  "LIT",
		Name: "linctl",
	}
	user := client.UserSummary{
		ID:          "user-id",
		DisplayName: "Omer",
		Email:       "omer@example.com",
	}
	comment := client.CommentSummary{
		ID:          "comment-id",
		Body:        "First comment",
		DisplayName: "Omer",
	}
	workflowState := client.WorkflowStateSummary{
		ID:   "workflow-state-id",
		Name: "Started",
		Type: "started",
	}

	textOut := bytes.Buffer{}
	textCommand := &cobra.Command{}
	textCommand.SetOut(&textOut)
	textOptions := rootOptions{}

	require.NoError(t, writeIssue(textCommand, &textOptions, issue))
	require.NoError(t, writeCycle(textCommand, &textOptions, cycle))
	require.NoError(t, writeProject(textCommand, &textOptions, project))
	require.NoError(t, writeProjectUpdate(textCommand, &textOptions, projectUpdate))
	require.NoError(t, writeProjectMilestone(textCommand, &textOptions, milestone))
	require.NoError(t, writeDocument(textCommand, &textOptions, document))
	require.NoError(t, writeLabel(textCommand, &textOptions, label))
	require.NoError(t, writeTeam(textCommand, &textOptions, team))
	require.NoError(t, writeUser(textCommand, &textOptions, user))
	require.NoError(t, writeComment(textCommand, &textOptions, comment))
	require.NoError(t, writeWorkflowState(textCommand, &textOptions, workflowState))
	require.Equal(
		t,
		"LIT-1 Ship coverage [Todo]\ncycle-id Planning cycle [active]\n"+
			"project-id Coverage [Backlog]\nproject-update-id onTrack Omer First update\n"+
			"project-milestone-id Launch milestone [next]\n"+
			"document-id Spec [project]\nlabel-id Bug #ff0000\nteam-id LIT linctl\n"+
			"user-id Omer <omer@example.com>\ncomment-id Omer First comment\nworkflow-state-id Started [started]\n",
		textOut.String(),
	)

	jsonOut := bytes.Buffer{}
	jsonCommand := &cobra.Command{}
	jsonCommand.SetOut(&jsonOut)
	jsonOptions := rootOptions{json: true}

	require.NoError(t, writeIssue(jsonCommand, &jsonOptions, issue))
	require.NoError(t, writeCycle(jsonCommand, &jsonOptions, cycle))
	require.NoError(t, writeProject(jsonCommand, &jsonOptions, project))
	require.NoError(t, writeProjectUpdate(jsonCommand, &jsonOptions, projectUpdate))
	require.NoError(t, writeProjectMilestone(jsonCommand, &jsonOptions, milestone))
	require.NoError(t, writeDocument(jsonCommand, &jsonOptions, document))
	require.NoError(t, writeLabel(jsonCommand, &jsonOptions, label))
	require.NoError(t, writeTeam(jsonCommand, &jsonOptions, team))
	require.NoError(t, writeUser(jsonCommand, &jsonOptions, user))
	require.NoError(t, writeComment(jsonCommand, &jsonOptions, comment))
	require.NoError(t, writeWorkflowState(jsonCommand, &jsonOptions, workflowState))
	require.Contains(t, jsonOut.String(), `"identifier": "LIT-1"`)
	require.Contains(t, jsonOut.String(), `"name": "Planning cycle"`)
	require.Contains(t, jsonOut.String(), `"name": "Coverage"`)
	require.Contains(t, jsonOut.String(), `"body": "First update"`)
	require.Contains(t, jsonOut.String(), `"name": "Launch milestone"`)
	require.Contains(t, jsonOut.String(), `"title": "Spec"`)
	require.Contains(t, jsonOut.String(), `"color": "#ff0000"`)
	require.Contains(t, jsonOut.String(), `"key": "LIT"`)
	require.Contains(t, jsonOut.String(), `"email": "omer@example.com"`)
	require.Contains(t, jsonOut.String(), `"body": "First comment"`)
	require.Contains(t, jsonOut.String(), `"type": "started"`)
}

func Test_CliOutputHelpers_cover_machine_output_edges(t *testing.T) {
	command := &cobra.Command{}
	output := bytes.Buffer{}
	command.SetOut(&output)
	issue := client.IssueSummary{
		ID:         "issue-id",
		Identifier: "LIT-1",
		Title:      "Ship coverage",
		State:      "Todo",
		Project:    "Pinned project",
		URL:        "https://linear.app/issue/LIT-1",
	}
	project := client.ProjectSummary{
		ID:   "project-id",
		Name: "Coverage",
		URL:  "https://linear.app/project/project-id",
		Status: client.ProjectStatus{
			Name: "Backlog",
		},
	}
	projectUpdate := client.ProjectUpdateSummary{
		ID:          "project-update-id",
		Body:        "First update",
		Health:      "onTrack",
		DisplayName: "Omer",
	}
	cycle := client.CycleSummary{
		ID:       "cycle-id",
		Name:     "Planning cycle",
		Status:   "active",
		StartsAt: "2026-07-01T00:00:00Z",
		EndsAt:   "2026-07-15T00:00:00Z",
		Progress: 0.5,
	}
	milestone := client.ProjectMilestoneSummary{
		ID:         "project-milestone-id",
		Name:       "Launch milestone",
		Status:     "next",
		TargetDate: "2026-06-30",
		Progress:   0.5,
	}
	document := client.DocumentSummary{
		ID:         "document-id",
		Title:      "Spec",
		ParentType: "project",
	}
	label := client.LabelSummary{
		ID:    "label-id",
		Name:  "Bug",
		Color: "#ff0000",
	}
	team := client.TeamSummary{
		ID:   "team-id",
		Key:  "LIT",
		Name: "linctl",
	}
	user := client.UserSummary{
		ID:          "user-id",
		DisplayName: "Omer",
		Email:       "omer@example.com",
	}
	comment := client.CommentSummary{
		ID:          "comment-id",
		Body:        "First comment",
		DisplayName: "Omer",
	}
	workflowState := client.WorkflowStateSummary{
		ID:   "workflow-state-id",
		Name: "Started",
		Type: "started",
	}

	require.NoError(t, writeIssue(command, &rootOptions{format: "full"}, issue))
	require.NoError(t, writeIssue(command, &rootOptions{idOnly: true}, issue))
	require.NoError(t, writeCycle(command, &rootOptions{format: "minimal"}, cycle))
	require.NoError(t, writeCycle(command, &rootOptions{format: "full"}, cycle))
	require.NoError(t, writeCycle(command, &rootOptions{idOnly: true}, cycle))
	require.NoError(t, writeProject(command, &rootOptions{format: "minimal"}, project))
	require.NoError(t, writeProject(command, &rootOptions{format: "full"}, project))
	require.NoError(t, writeProject(command, &rootOptions{idOnly: true}, project))
	require.NoError(t, writeProjectUpdate(command, &rootOptions{idOnly: true}, projectUpdate))
	require.NoError(t, writeProjectMilestone(command, &rootOptions{format: "minimal"}, milestone))
	require.NoError(t, writeProjectMilestone(command, &rootOptions{format: "full"}, milestone))
	require.NoError(t, writeProjectMilestone(command, &rootOptions{idOnly: true}, milestone))
	require.NoError(t, writeDocument(command, &rootOptions{idOnly: true}, document))
	require.NoError(t, writeLabel(command, &rootOptions{idOnly: true}, label))
	require.NoError(t, writeTeam(command, &rootOptions{idOnly: true}, team))
	require.NoError(t, writeUser(command, &rootOptions{idOnly: true}, user))
	require.NoError(t, writeComment(command, &rootOptions{idOnly: true}, comment))
	require.NoError(t, writeWorkflowState(command, &rootOptions{idOnly: true}, workflowState))
	require.Contains(t, output.String(), "project=Pinned project")
	require.Contains(t, output.String(), "issue-id")
	require.Contains(t, output.String(), "starts_at=2026-07-01T00:00:00Z")
	require.Contains(t, output.String(), "cycle-id")
	require.Contains(t, output.String(), "project-id")
	require.Contains(t, output.String(), "project-update-id")
	require.Contains(t, output.String(), "target_date=2026-06-30")
	require.Contains(t, output.String(), "project-milestone-id")
	require.Contains(t, output.String(), "document-id")
	require.Contains(t, output.String(), "label-id")
	require.Contains(t, output.String(), "team-id")
	require.Contains(t, output.String(), "user-id")
	require.Contains(t, output.String(), "comment-id")
	require.Contains(t, output.String(), "workflow-state-id")
	require.Equal(t, "-", emptyDash(""))

	quietOutput := bytes.Buffer{}
	quietCommand := &cobra.Command{}
	quietCommand.SetOut(&quietOutput)
	require.NoError(t, writeJSONValue(quietCommand, &rootOptions{quiet: true}, issue))
	require.NoError(t, writeCycle(quietCommand, &rootOptions{quiet: true}, cycle))
	require.NoError(t, writeProject(quietCommand, &rootOptions{quiet: true}, project))
	require.NoError(t, writeProjectUpdate(quietCommand, &rootOptions{quiet: true}, projectUpdate))
	require.NoError(t, writeProjectMilestone(quietCommand, &rootOptions{quiet: true}, milestone))
	require.NoError(t, writeDocument(quietCommand, &rootOptions{quiet: true}, document))
	require.NoError(t, writeLabel(quietCommand, &rootOptions{quiet: true}, label))
	require.NoError(t, writeTeam(quietCommand, &rootOptions{quiet: true}, team))
	require.NoError(t, writeUser(quietCommand, &rootOptions{quiet: true}, user))
	require.NoError(t, writeComment(quietCommand, &rootOptions{quiet: true}, comment))
	require.NoError(t, writeWorkflowState(quietCommand, &rootOptions{quiet: true}, workflowState))
	require.NoError(t, writeScalar(quietCommand, &rootOptions{quiet: true}, "title", "quiet"))
	wrote, err := writeIDOnly(quietCommand, &rootOptions{idOnly: true, quiet: true}, "issue-id")
	require.NoError(t, err)
	require.True(t, wrote)
	require.Empty(t, quietOutput.String())

	scalarJSONOutput := bytes.Buffer{}
	scalarJSONCommand := &cobra.Command{}
	scalarJSONCommand.SetOut(&scalarJSONOutput)
	require.NoError(t, writeScalar(scalarJSONCommand, &rootOptions{json: true}, "title", "Ship coverage"))
	require.Contains(t, scalarJSONOutput.String(), `"title": "Ship coverage"`)

	wrote, err = writeIDOnly(command, &rootOptions{idOnly: true}, "")
	require.Error(t, err)
	require.True(t, wrote)
	require.Contains(t, err.Error(), "id is empty")

	require.NoError(t, ensureNonEmpty(&rootOptions{}, 0))
	require.Error(t, writeIssue(command, &rootOptions{format: "wide"}, issue))
	require.Error(t, writeCycle(command, &rootOptions{format: "wide"}, cycle))
	require.Error(t, writeProject(command, &rootOptions{format: "wide"}, project))
	require.Error(t, writeProjectMilestone(command, &rootOptions{format: "wide"}, milestone))
	_, err = normalizedHumanFormat(&rootOptions{format: "wide"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid format")

	err = writeJSONValue(command, &rootOptions{json: true, fields: "missing"}, issue)
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"missing\" is not present")
}

func Test_CliOutputHelpers_cover_json_projection_and_sort_edges(t *testing.T) {
	projected, err := projectJSONFields(
		map[string]any{"issues": []any{map[string]any{"identifier": "LIT-1", "state": map[string]any{"name": "Todo"}}}},
		"identifier,state.name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"issues": []any{map[string]any{"identifier": "LIT-1", "state": map[string]any{"name": "Todo"}}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"projects": []any{map[string]any{"id": "project-id", "status": map[string]any{"name": "Backlog"}}}},
		"id,status.name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"projects": []any{map[string]any{"id": "project-id", "status": map[string]any{"name": "Backlog"}}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"members": []any{map[string]any{"id": "user-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"members": []any{map[string]any{"id": "user-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"comments": []any{map[string]any{"id": "comment-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"comments": []any{map[string]any{"id": "comment-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"documents": []any{map[string]any{"id": "document-id", "title": "Spec"}}},
		"id,title",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"documents": []any{map[string]any{"id": "document-id", "title": "Spec"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"labels": []any{map[string]any{"id": "label-id", "color": "#ff0000"}}},
		"id,color",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"labels": []any{map[string]any{"id": "label-id", "color": "#ff0000"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"teams": []any{map[string]any{"id": "team-id", "key": "LIT"}}},
		"id,key",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"teams": []any{map[string]any{"id": "team-id", "key": "LIT"}},
	}, projected)

	projected, err = projectJSONFields(
		map[string]any{"users": []any{map[string]any{"id": "user-id", "display_name": "Omer"}}},
		"id,display_name",
	)
	require.NoError(t, err)
	require.Equal(t, map[string]any{
		"users": []any{map[string]any{"id": "user-id", "display_name": "Omer"}},
	}, projected)

	projected, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "identifier")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	projected, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "identifier,, ")
	require.NoError(t, err)
	require.Equal(t, map[string]any{"identifier": "LIT-1"}, projected)

	_, err = projectJSONFields(map[string]any{"bad": func() {}}, "bad")
	require.Error(t, err)
	require.Contains(t, err.Error(), "marshal output")

	_, err = projectJSONFields([]string{"not-an-object"}, "id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode output")

	_, err = projectJSONFields(map[string]any{"issues": []any{"bad-item"}}, "identifier")
	require.Error(t, err)
	require.Contains(t, err.Error(), "item is not an object")

	_, err = projectJSONFields(map[string]any{"issues": []any{map[string]any{"title": "Missing id"}}}, "identifier")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"identifier\" is not present")

	_, err = projectJSONFields(map[string]any{"identifier": "LIT-1"}, "missing")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"missing\" is not present")

	_, err = projectJSONFields(map[string]any{"state": "Todo"}, "state.name")
	require.Error(t, err)
	require.Contains(t, err.Error(), "field \"state\" is not an object")

	items := []client.IssueSummary{
		{Identifier: "LIT-2", Title: "Zebra"},
		{Identifier: "LIT-1", Title: "Alpha"},
	}
	sortedItems, err := sortByJSONField(items, "", "asc")
	require.NoError(t, err)
	require.Equal(t, items, sortedItems)

	sortedItems, err = sortByJSONField(items, "title", "asc")
	require.NoError(t, err)
	require.Equal(t, "Alpha", sortedItems[0].Title)

	_, err = sortByJSONField(items, "title", "sideways")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid sort order")

	_, err = sortByJSONField(items, "missing", "asc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "sort field \"missing\" is not present")

	_, err = sortByJSONField([]map[string]any{{"state": "Todo"}}, "state.name", "asc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not an object path")

	_, err = jsonFieldValue(map[string]any{"bad": func() {}}, "bad")
	require.Error(t, err)
	require.Contains(t, err.Error(), "marshal output")

	destination := map[string]any{}
	require.NoError(t, copyJSONPath(map[string]any{"id": "issue-id"}, destination, nil))
	require.Empty(t, destination)
}

func Test_CommandFlows_cover_output_error_and_quiet_branches(t *testing.T) {
	quietCommands := [][]string{
		{"--quiet", "target"},
		{"--quiet", "whoami"},
		{"--quiet", "issue", "deps", "LIT-1"},
		{"--quiet", "issue", "pr", "LIT-1"},
		{"--quiet", "usage"},
	}
	for _, args := range quietCommands {
		t.Run("quiet "+args[len(args)-1], func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Empty(t, output.String())
		})
	}

	errorCommands := [][]string{
		{"--sort", "missing", "issue", "list"},
		{"--sort", "missing", "issue", "list", "--project", "project-id"},
		{"--sort", "missing", "issue", "list", "--mine"},
		{"--sort", "missing", "issue", "list", "--assignee", "assignee-id"},
		{"--sort", "missing", "issue", "list", "--label", "label-id"},
		{"--sort", "missing", "issue", "list", "--cycle", "cycle-id"},
		{"--sort", "missing", "issue", "list", "--created-after", "2026-06-01"},
		{"--sort", "missing", "issue", "list", "--created-since", "2026-06-01"},
		{"--sort", "missing", "issue", "list", "--created-before", "2026-06-30"},
		{"--sort", "missing", "issue", "list", "--has-blockers"},
		{"--sort", "missing", "issue", "list", "--blocks"},
		{"--sort", "missing", "issue", "list", "--blocked-by", "LIT-1"},
		{"--sort", "missing", "issue", "list", "--all-teams"},
		{"--sort", "missing", "issue", "comments", "LIT-1"},
		{"--sort", "missing", "issue", "search", "needle"},
		{"--sort", "missing", "project", "list"},
		{"--sort", "missing", "project", "members", "project-id"},
	}
	for _, args := range errorCommands {
		t.Run("sort error "+args[len(args)-1], func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "sort field")
		})
	}

	emptyCommands := []struct {
		name string
		args []string
		fake commandFlowFakeClient
	}{
		{name: "issue list project", args: []string{"--fail-on-empty", "issue", "list", "--project", "project-id"}, fake: commandFlowFakeClient{emptyIssueProject: true}},
		{name: "issue list mine", args: []string{"--fail-on-empty", "issue", "list", "--mine"}, fake: commandFlowFakeClient{emptyIssueMine: true}},
		{name: "issue list assignee", args: []string{"--fail-on-empty", "issue", "list", "--assignee", "assignee-id"}, fake: commandFlowFakeClient{emptyIssueMine: true}},
		{name: "issue list label", args: []string{"--fail-on-empty", "issue", "list", "--label", "label-id"}, fake: commandFlowFakeClient{emptyIssueLabel: true}},
		{name: "issue list cycle", args: []string{"--fail-on-empty", "issue", "list", "--cycle", "cycle-id"}, fake: commandFlowFakeClient{emptyIssueCycle: true}},
		{name: "issue list created-after", args: []string{"--fail-on-empty", "issue", "list", "--created-after", "2026-06-01"}, fake: commandFlowFakeClient{emptyIssueCreatedAfter: true}},
		{name: "issue list created-since", args: []string{"--fail-on-empty", "issue", "list", "--created-since", "2026-06-01"}, fake: commandFlowFakeClient{emptyIssueCreatedAfter: true}},
		{name: "issue list created-before", args: []string{"--fail-on-empty", "issue", "list", "--created-before", "2026-06-30"}, fake: commandFlowFakeClient{emptyIssueCreatedBefore: true}},
		{name: "issue list has blockers", args: []string{"--fail-on-empty", "issue", "list", "--has-blockers"}, fake: commandFlowFakeClient{emptyIssueHasBlockers: true}},
		{name: "issue list blocks", args: []string{"--fail-on-empty", "issue", "list", "--blocks"}, fake: commandFlowFakeClient{emptyIssueBlocks: true}},
		{name: "issue list blocked by", args: []string{"--fail-on-empty", "issue", "list", "--blocked-by", "LIT-1"}, fake: commandFlowFakeClient{emptyIssueBlockedBy: true}},
		{name: "issue list all teams", args: []string{"--fail-on-empty", "issue", "list", "--all-teams"}, fake: commandFlowFakeClient{emptyIssueAllTeams: true}},
		{name: "issue comments", args: []string{"--fail-on-empty", "issue", "comments", "LIT-1"}, fake: commandFlowFakeClient{emptyIssueComments: true}},
		{name: "issue search", args: []string{"--fail-on-empty", "issue", "search", "needle"}, fake: commandFlowFakeClient{emptyIssueSearch: true}},
		{name: "project list", args: []string{"--fail-on-empty", "project", "list"}, fake: commandFlowFakeClient{emptyProjectList: true}},
		{name: "project members", args: []string{"--fail-on-empty", "project", "members", "project-id"}, fake: commandFlowFakeClient{emptyProjectMembers: true}},
	}
	for _, test := range emptyCommands {
		t.Run("empty "+test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "empty result")
		})
	}
}

func Test_CommandFlows_cover_issue_list_filter_validation(t *testing.T) {
	tests := [][]string{
		{"issue", "list", "--state", "started", "--project", "project-id"},
		{"issue", "list", "--state", "started", "--mine"},
		{"issue", "list", "--state", "started", "--assignee", "assignee-id"},
		{"issue", "list", "--state", "started", "--label", "label-id"},
		{"issue", "list", "--state", "started", "--cycle", "cycle-id"},
		{"issue", "list", "--state", "started", "--created-after", "2026-06-01"},
		{"issue", "list", "--state", "started", "--created-since", "2026-06-01"},
		{"issue", "list", "--created-after", "2026-06-01", "--created-since", "2026-06-01"},
		{"issue", "list", "--state", "started", "--created-before", "2026-06-30"},
		{"issue", "list", "--state", "started", "--has-blockers"},
		{"issue", "list", "--has-blockers", "--blocks"},
		{"issue", "list", "--blocks", "--blocked-by", "LIT-1"},
		{"issue", "list", "--state", "started", "--all-teams"},
	}
	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "use only one")
		})
	}
}

func Test_CommandFlows_cover_issue_current_error_branches(t *testing.T) {
	t.Run("id missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("title missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "title"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("title runtime error", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(
			t,
			[]string{"issue", "title"},
			func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			},
		)

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("url lookup error", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(
			t,
			[]string{"issue", "url"},
			func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return testCommandRuntime(commandFlowFakeClient{failOperation: "issue"}), nil
			},
		)

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})

	t.Run("branch argument lookup error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "branch", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})
}

func Test_CommandFlows_cover_issue_comment_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"issue", "comment", "LIT-1", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_issue_reply_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"issue", "reply", "LIT-1", "comment-id", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_issue_comments_error_branches(t *testing.T) {
	t.Run("operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "IssueComments"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "comments", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue comments LIT-1")
	})

	t.Run("writer error", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssueComments(command, []client.IssueCommentSummary{{ID: "comment-id", DisplayName: "Omer", Body: "body"}})

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_cover_issue_deps_writer_error(t *testing.T) {
	t.Run("issue header", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("section header", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssueDependencySection(command, &rootOptions{}, "children", nil)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("parent issue", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 3})
		parent := client.IssueSummary{Identifier: "LIT-2", Title: "Parent", State: "Todo"}
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1", Parent: &parent}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("children section", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 2})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("blocks section", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 3})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

type countFailingWriter struct {
	failAt int
	writes int
}

func (writer *countFailingWriter) Write(content []byte) (int, error) {
	writer.writes++
	if writer.writes == writer.failAt {
		return 0, errors.New("write failed")
	}

	return len(content), nil
}

func Test_CliHelpers_resolve_target_overrides_and_project_ids(t *testing.T) {
	options := rootOptions{
		orgID:   "org-id",
		team:    "LIT",
		project: "project-id",
	}

	target := targetOverride(&options)

	require.Equal(t, "org-id", target.OrgID)
	require.Equal(t, "LIT", target.TeamKey)
	require.Equal(t, "project-id", target.ProjectID)
	require.Empty(t, projectID(nil))
	require.Equal(t, "project-id", projectID(&client.ResolvedProject{ID: "project-id"}))
	require.NotEmpty(t, defaultGlobalConfigPath())
}

type commandFailingReader struct{}

func (reader commandFailingReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}

func Test_CommandRuntime_loads_config_and_requires_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("LINCTL_TOKEN", "")
	t.Setenv("LINEAR_API_KEY", "")
	require.NoError(t, os.WriteFile(".linctl.toml", []byte(`
[target]
org_id = "org-id"
team_key = "LIT"
team_id = "team-id"
project_id = "project-id"
`), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing Linear token")

	t.Setenv("LINCTL_TOKEN", "test-token")
	runtime, err := newCommandRuntime(context.Background(), &rootOptions{})
	require.NoError(t, err)
	require.Equal(t, "project-id", runtime.config.Target.ProjectID)
	require.NotNil(t, runtime.graphqlClient)
}

func Test_CommandRuntime_reports_config_load_errors(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())
	require.NoError(t, os.WriteFile(".linctl.toml", []byte("[target\n"), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "parse config")
}

func Test_DefaultGlobalConfigPath_returns_empty_when_home_is_unset(t *testing.T) {
	t.Setenv("HOME", "")

	require.Empty(t, defaultGlobalConfigPath())
}

func Test_WriteUsage_reports_unknown_topics(t *testing.T) {
	command := &cobra.Command{}

	err := writeUsage(command, &rootOptions{}, "missing")

	require.Error(t, err)
	require.Contains(t, err.Error(), `unknown usage topic "missing"`)
}
