package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

// payloadGraphQLClient returns a fixed GraphQL data payload, mirroring how
// commandFlowFakeClient unmarshals canned responses.
type payloadGraphQLClient struct {
	payload string
}

func (client payloadGraphQLClient) MakeRequest(
	_ context.Context,
	_ *graphql.Request,
	response *graphql.Response,
) error {
	return json.Unmarshal([]byte(`{"data":`+client.payload+`}`), response)
}

func Test_Completion_prints_no_candidates_when_linear_is_unreachable(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		failOperation string
		failRuntime   bool
	}{
		{name: "runtime fails", args: []string{"issue", "list", "--project", ""}, failRuntime: true},
		{name: "team list fails", args: []string{"init", "--team", ""}, failOperation: "teams_list"},
		{name: "project list fails", args: []string{"issue", "list", "--project", ""}, failOperation: "projects"},
		{name: "state type list fails", args: []string{"issue", "list", "--state", ""}, failOperation: "workflowStates"},
		{name: "state name list fails", args: []string{"issue", "create", "--state", ""}, failOperation: "workflowStates"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: test.failOperation})
			defer restore()
			if test.failRuntime {
				buildCommandRuntime = func(context.Context, *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, errors.New("runtime failed")
				}
			}

			output := completeCommand(t, test.args)

			require.Equal(t, ":4\n", output)
		})
	}
}

func Test_IssueList_completes_distinct_state_types(t *testing.T) {
	restore := useCommandRuntime(t, payloadGraphQLClient{payload: `{"workflowStates":{"nodes":[
		{"id":"1","name":"A","type":"started","position":1},
		{"id":"2","name":"B","type":"started","position":2},
		{"id":"3","name":"C","type":"","position":3},
		{"id":"4","name":"D","type":"completed","position":4}
	],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`})
	defer restore()

	output := completeCommand(t, []string{"issue", "list", "--state", ""})

	require.Equal(t, "started\ncompleted\n:4\n", output)
}

func Test_IssueCreate_completes_distinct_state_names(t *testing.T) {
	restore := useCommandRuntime(t, payloadGraphQLClient{payload: `{"workflowStates":{"nodes":[
		{"id":"1","name":"In Review","type":"started","position":1},
		{"id":"2","name":"in review","type":"started","position":2},
		{"id":"3","name":"","type":"started","position":3},
		{"id":"4","name":"Todo","type":"unstarted","position":4}
	],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`})
	defer restore()

	output := completeCommand(t, []string{"issue", "create", "--state", ""})

	require.Equal(t, "In Review\tstarted\nTodo\tunstarted\n:4\n", output)
}

func Test_IssueUpdate_completes_state_names(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	output := completeCommand(t, []string{"issue", "update", "--state", ""})

	require.Contains(t, output, "Started")
}

func Test_Init_completes_local_project_and_team_flags(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	projectOutput := completeCommand(t, []string{"init", "--project", ""})
	require.Contains(t, projectOutput, "project-id")
	require.Contains(t, projectOutput, ":4")

	teamOutput := completeCommand(t, []string{"init", "--team", ""})
	require.Contains(t, teamOutput, "LIT")
	require.Contains(t, teamOutput, ":4")
}

func Test_IssueList_completes_local_project_flag(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	output := completeCommand(t, []string{"issue", "list", "--project", ""})

	require.Equal(t, "project-id\tListed project\n:4\n", output)
}

func Test_ProjectGet_completes_only_the_first_argument(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	require.Equal(t, "project-id\tListed project\n:4\n", completeCommand(t, []string{"project", "get", ""}))
	require.Equal(t, ":4\n", completeCommand(t, []string{"project", "get", "project-id", ""}))
}

func completeCommand(t *testing.T, args []string) string {
	t.Helper()
	stdout := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&stdout)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs(append([]string{cobra.ShellCompRequestCmd}, args...))
	require.NoError(t, command.ExecuteContext(context.Background()))

	return stdout.String()
}
