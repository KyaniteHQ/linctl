package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_CommandFlows_execute_read_and_write_commands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{name: "target", args: []string{"target"}, contains: "org org-id team LIT/team-id project project-id confirmed true"},
		{name: "whoami", args: []string{"whoami"}, contains: "Omer <omer@example.com>"},
		{name: "issue list", args: []string{"issue", "list", "--limit", "1"}, contains: "LIT-1 Listed issue [Todo]"},
		{name: "issue get", args: []string{"issue", "get", "LIT-1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "issue create", args: []string{"issue", "create", "--title", "Created issue"}, contains: "LIT-2 Created issue [Todo]"},
		{name: "issue update", args: []string{"issue", "update", "LIT-1", "--title", "Updated issue"}, contains: "LIT-1 Updated issue [Todo]"},
		{name: "issue comment", args: []string{"issue", "comment", "LIT-1", "--body", "Looks good"}, contains: "comment comment-id on LIT-1"},
		{name: "issue close", args: []string{"issue", "close", "LIT-1"}, contains: "LIT-1 Closed issue [Done]"},
		{name: "project list", args: []string{"project", "list", "--limit", "1"}, contains: "project-id Listed project [Backlog]"},
		{name: "project get", args: []string{"project", "get", "project-id"}, contains: "project-id Detail project [Backlog]"},
		{name: "project members", args: []string{"project", "members", "project-id", "--limit", "1"}, contains: "user-id Omer"},
		{name: "project create", args: []string{"project", "create", "--name", "Created project"}, contains: "project-id Created project [Backlog]"},
		{name: "project update", args: []string{"project", "update", "project-id", "--name", "Updated project"}, contains: "project-id Updated project [Started]"},
		{name: "project archive", args: []string{"project", "archive", "project-id"}, contains: "project-id Archived project [Canceled]"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), test.contains)
		})
	}
}

func Test_CommandFlows_resolve_current_issue_from_branch(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"current"})

	require.NoError(t, err)
	require.Contains(t, output, "LIT-1 Detail issue [Todo]")
}

func Test_CommandFlows_print_current_issue_as_json(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"--json", "current"})

	require.NoError(t, err)
	require.Contains(t, output, `"identifier": "LIT-1"`)
}

func Test_CommandFlows_report_current_issue_errors(t *testing.T) {
	t.Run("missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"current"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("runtime failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"current"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("issue lookup failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"current"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return testCommandRuntime(commandFlowFakeClient{failOperation: "IssueByID"}), nil
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})
}

func runCurrentCommandInGitBranch(t *testing.T, args []string) (string, error) {
	t.Helper()

	return runCurrentCommandInGitBranchWithRuntime(t, args, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return testCommandRuntime(commandFlowFakeClient{}), nil
	})
}

func runCurrentCommandInGitBranchWithRuntime(
	t *testing.T,
	args []string,
	runtimeBuilder func(context.Context, *rootOptions) (commandRuntime, error),
) (string, error) {
	t.Helper()

	dir := t.TempDir()
	runGitCommand(t, dir, "init")
	runGitCommand(t, dir, "checkout", "-b", "feature/LIT-1-current")
	t.Chdir(dir)

	output := bytes.Buffer{}
	original := buildCommandRuntime
	buildCommandRuntime = runtimeBuilder
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs(args)

	err := command.ExecuteContext(context.Background())
	return output.String(), err
}

func Test_CommandFlows_report_runtime_and_writer_errors(t *testing.T) {
	t.Run("runtime error returns from command", func(t *testing.T) {
		commands := [][]string{
			{"target"},
			{"whoami"},
			{"issue", "list"},
			{"issue", "get", "LIT-1"},
			{"issue", "create", "--title", "Created issue"},
			{"issue", "update", "LIT-1", "--title", "Updated issue"},
			{"issue", "comment", "LIT-1", "--body", "Looks good"},
			{"issue", "close", "LIT-1"},
			{"project", "list"},
			{"project", "get", "project-id"},
			{"project", "members", "project-id"},
			{"project", "create", "--name", "Created project"},
			{"project", "update", "project-id", "--name", "Updated project"},
			{"project", "archive", "project-id"},
		}
		for _, args := range commands {
			t.Run(strings.Join(args, " "), func(t *testing.T) {
				original := buildCommandRuntime
				buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, errors.New("runtime failed")
				}
				defer func() {
					buildCommandRuntime = original
				}()
				command := NewRootCommand(context.Background(), BuildInfo{})
				command.SetArgs(args)

				err := command.ExecuteContext(context.Background())

				require.Error(t, err)
				require.Contains(t, err.Error(), "runtime failed")
			})
		}
	})

	t.Run("writeIssues returns writer errors", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssues(command, &rootOptions{}, []client.IssueSummary{{Identifier: "LIT-1", Title: "Broken", State: "Todo"}})

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "list"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project members returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "members", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("usage returns writer errors", func(t *testing.T) {
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"usage"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_print_json_for_read_and_comment_commands(t *testing.T) {
	tests := [][]string{
		{"--json", "target"},
		{"--json", "whoami"},
		{"--json", "issue", "list", "--limit", "1"},
		{"--json", "issue", "comment", "LIT-1", "--body", "Looks good"},
		{"--json", "project", "list", "--limit", "1"},
		{"--json", "project", "members", "project-id", "--limit", "1"},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), "{")
		})
	}
}

func Test_CommandFlows_report_operation_errors(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		operation string
		contains  string
	}{
		{name: "target resolve", args: []string{"target"}, operation: "Teams", contains: "resolve teams"},
		{name: "whoami resolve", args: []string{"whoami"}, operation: "Viewer", contains: "resolve viewer"},
		{name: "issue list target resolve", args: []string{"issue", "list"}, operation: "Teams", contains: "resolve teams"},
		{name: "issue list", args: []string{"issue", "list"}, operation: "IssuesByTeam", contains: "list issues"},
		{name: "issue get", args: []string{"issue", "get", "LIT-1"}, operation: "IssueByID", contains: "get issue LIT-1"},
		{name: "issue create", args: []string{"issue", "create", "--title", "Created issue"}, operation: "IssueCreate", contains: "create issue"},
		{name: "issue update", args: []string{"issue", "update", "LIT-1", "--title", "Updated issue"}, operation: "IssueUpdate", contains: "update issue LIT-1"},
		{name: "issue comment", args: []string{"issue", "comment", "LIT-1", "--body", "Looks good"}, operation: "IssueCommentCreate", contains: "comment on issue LIT-1"},
		{name: "issue close", args: []string{"issue", "close", "LIT-1"}, operation: "IssueClose", contains: "close issue LIT-1"},
		{name: "project list target resolve", args: []string{"project", "list"}, operation: "Teams", contains: "resolve teams"},
		{name: "project list", args: []string{"project", "list"}, operation: "Projects", contains: "list projects"},
		{name: "project get", args: []string{"project", "get", "project-id"}, operation: "ProjectByID", contains: "get project project-id"},
		{name: "project members", args: []string{"project", "members", "project-id"}, operation: "ProjectMembers", contains: "list project members project-id"},
		{name: "project create", args: []string{"project", "create", "--name", "Created project"}, operation: "ProjectCreate", contains: "create project"},
		{name: "project update", args: []string{"project", "update", "project-id", "--name", "Updated project"}, operation: "ProjectUpdate", contains: "update project project-id"},
		{name: "project archive", args: []string{"project", "archive", "project-id"}, operation: "ProjectArchive", contains: "archive project project-id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: test.operation})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.contains)
		})
	}
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()

	command := exec.Command("git", args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	require.NoError(t, err, string(output))
}

type commandFailingWriter struct{}

func (writer commandFailingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failed")
}

func useCommandRuntime(t *testing.T, graphqlClient graphql.Client) func() {
	t.Helper()

	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return testCommandRuntime(graphqlClient), nil
	}

	return func() {
		buildCommandRuntime = original
	}
}

func testCommandRuntime(graphqlClient graphql.Client) commandRuntime {
	return commandRuntime{
		config: config.Resolved{
			Token: "test-token",
			Target: config.Target{
				OrgID:     "org-id",
				TeamKey:   "LIT",
				TeamID:    "team-id",
				ProjectID: "project-id",
			},
		},
		graphqlClient: graphqlClient,
	}
}

type commandFlowFakeClient struct {
	failOperation string
}

func (client commandFlowFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if request.OpName == client.failOperation {
		return errors.New("operation failed")
	}

	payload, err := commandFlowPayload(request.OpName)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

func commandFlowPayload(operation string) (string, error) {
	switch operation {
	case "Viewer":
		return `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`, nil
	case "Teams":
		return `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, nil
	case "TargetProject":
		return `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`, nil
	}
	if payload, ok := commandFlowIssuePayload(operation); ok {
		return payload, nil
	}
	if payload, ok := commandFlowProjectPayload(operation); ok {
		return payload, nil
	}

	return "", errors.New("missing fake response for " + operation)
}

func commandFlowIssuePayload(operation string) (string, bool) {
	switch operation {
	case "IssuesByTeam":
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-1", "Listed issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssueByID":
		return `{"issue":` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `}`, true
	case "IssueCreate":
		return `{"issueCreate":{"success":true,"issue":` + commandIssueJSON("LIT-2", "Created issue", "todo-state", "Todo", "unstarted") + `}}`, true
	case "IssueUpdate":
		return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSON("LIT-1", "Updated issue", "todo-state", "Todo", "unstarted") + `}}`, true
	case "IssueCommentCreate":
		return `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"Looks good","url":"https://linear.app/comment/comment-id","issue":` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `}}}`, true
	case "CompletedWorkflowStates":
		return `{"workflowStates":{"nodes":[{"id":"done-state","name":"Done","type":"completed","position":1}]}}`, true
	case "IssueClose":
		return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSON("LIT-1", "Closed issue", "done-state", "Done", "completed") + `}}`, true
	default:
		return "", false
	}
}

func commandFlowProjectPayload(operation string) (string, bool) {
	switch operation {
	case "Projects":
		return `{"team":{"projects":{"nodes":[` + commandProjectJSON("Listed project", "Backlog", "backlog") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "ProjectByID":
		return `{"project":` + commandProjectJSON("Detail project", "Backlog", "backlog") + `}`, true
	case "ProjectMembers":
		return `{"project":{"id":"project-id","name":"Detail project","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com"}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "ProjectCreate":
		return `{"projectCreate":{"success":true,"project":` + commandProjectJSON("Created project", "Backlog", "backlog") + `}}`, true
	case "ProjectUpdate":
		return `{"projectUpdate":{"success":true,"project":` + commandProjectJSON("Updated project", "Started", "started") + `}}`, true
	case "ProjectArchive":
		return `{"projectArchive":{"success":true,"entity":` + commandProjectJSON("Archived project", "Canceled", "canceled") + `}}`, true
	default:
		return "", false
	}
}

func commandIssueJSON(identifier string, title string, stateID string, state string, stateType string) string {
	return `{
		"id":"issue-id",
		"identifier":"` + identifier + `",
		"title":"` + title + `",
		"url":"https://linear.app/kyanite/issue/` + identifier + `",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"state":{"id":"` + stateID + `","name":"` + state + `","type":"` + stateType + `"},
		"assignee":null,
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

func commandProjectJSON(name string, status string, statusType string) string {
	return `{
		"id":"project-id",
		"name":"` + name + `",
		"description":"description",
		"slugId":"` + name + `",
		"url":"https://linear.app/kyanite/project/project-id",
		"priority":0,
		"status":{"id":"status-id","name":"` + status + `","type":"` + statusType + `"},
		"lead":null,
		"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl"}]}
	}`
}

var _ graphql.Client = commandFlowFakeClient{}
