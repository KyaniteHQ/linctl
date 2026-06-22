package cli

import (
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

func Test_completionValues_returns_nothing_when_runtime_fails(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(context.Context, *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() { buildCommandRuntime = original }()

	values, directive := completionValues(context.Background(), &rootOptions{}, teamKeyCandidates)

	require.Nil(t, values)
	require.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func Test_completionValues_returns_nothing_when_load_fails(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "Teams"})
	defer restore()

	values, directive := completionValues(context.Background(), &rootOptions{}, teamKeyCandidates)

	require.Nil(t, values)
	require.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func Test_completionValues_returns_team_keys_on_success(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	values, directive := completionValues(context.Background(), &rootOptions{}, teamKeyCandidates)

	require.NotEmpty(t, values)
	require.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func Test_projectIDCandidates_returns_id_and_name(t *testing.T) {
	candidates, err := projectIDCandidates(context.Background(), testCommandRuntime(commandFlowFakeClient{}))

	require.NoError(t, err)
	require.NotEmpty(t, candidates)
	require.Contains(t, candidates[0], "\t")
}

func Test_teamKeyCandidates_returns_error_when_list_fails(t *testing.T) {
	_, err := teamKeyCandidates(
		context.Background(), testCommandRuntime(commandFlowFakeClient{failOperation: "Teams"}),
	)

	require.Error(t, err)
}

func Test_projectIDCandidates_returns_error_when_list_fails(t *testing.T) {
	_, err := projectIDCandidates(
		context.Background(), testCommandRuntime(commandFlowFakeClient{failOperation: "projects"}),
	)

	require.Error(t, err)
}

func Test_workflowStateTypeCandidates_returns_error_when_list_fails(t *testing.T) {
	_, err := workflowStateTypeCandidates(
		context.Background(), testCommandRuntime(commandFlowFakeClient{failOperation: "workflowStates"}),
	)

	require.Error(t, err)
}

func Test_workflowStateTypeCandidates_dedupes_and_skips_empty(t *testing.T) {
	runtime := commandRuntime{graphqlClient: payloadGraphQLClient{payload: `{"workflowStates":{"nodes":[
		{"id":"1","name":"A","type":"started","position":1},
		{"id":"2","name":"B","type":"started","position":2},
		{"id":"3","name":"C","type":"","position":3},
		{"id":"4","name":"D","type":"completed","position":4}
	],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`}}

	types, err := workflowStateTypeCandidates(context.Background(), runtime)

	require.NoError(t, err)
	require.Equal(t, []string{"started", "completed"}, types)
}

func Test_flagCompletion_invokes_loader(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	completion := flagCompletion(context.Background(), &rootOptions{}, projectIDCandidates)
	values, directive := completion(nil, nil, "")

	require.NotEmpty(t, values)
	require.Contains(t, values[0], "\t")
	require.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func Test_firstArgCompletion_completes_first_argument(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	completion := firstArgCompletion(context.Background(), &rootOptions{}, teamKeyCandidates)
	values, directive := completion(nil, nil, "")

	require.NotEmpty(t, values)
	require.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}

func Test_firstArgCompletion_skips_after_first_argument(t *testing.T) {
	completion := firstArgCompletion(context.Background(), &rootOptions{}, teamKeyCandidates)
	values, directive := completion(nil, []string{"already-set"}, "")

	require.Nil(t, values)
	require.Equal(t, cobra.ShellCompDirectiveNoFileComp, directive)
}
