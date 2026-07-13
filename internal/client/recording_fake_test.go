package client

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

// recordingGraphQLClient wraps any of the package's map-based test fakes and
// records the variables sent for each outbound operation, so tests can assert
// on WHAT linctl sends to Linear rather than only the canned response it gets
// back.
type recordingGraphQLClient struct {
	inner    graphql.Client
	requests []recordedRequest
}

type recordedRequest struct {
	OpName    string
	Variables json.RawMessage
}

func (client *recordingGraphQLClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	variables, err := json.Marshal(request.Variables)
	if err != nil {
		variables = json.RawMessage("null")
	}
	client.requests = append(client.requests, recordedRequest{
		OpName:    request.OpName,
		Variables: variables,
	})

	return client.inner.MakeRequest(ctx, request, response)
}

// variablesFor returns the variables of the LAST recorded request for opName,
// failing the test if the operation was never sent.
func (client *recordingGraphQLClient) variablesFor(t *testing.T, opName string) json.RawMessage {
	t.Helper()

	var found json.RawMessage
	for _, request := range client.requests {
		if request.OpName == opName {
			found = request.Variables
		}
	}
	require.NotNil(t, found, "operation %s was never sent", opName)

	return found
}
