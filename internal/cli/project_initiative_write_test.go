package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

type initiativeLinkCaptureClient struct {
	delegate graphql.Client
	create   *graphql.Request
}

func (client *initiativeLinkCaptureClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName == "InitiativeToProjectCreate" {
		client.create = request
	}

	return client.delegate.MakeRequest(ctx, request, response)
}

func Test_ProjectInitiativeLinkCommand_sends_sort_order_only_when_given(t *testing.T) {
	tests := []struct {
		name      string
		extraArgs []string
		input     string
	}{
		{
			name: "given", extraArgs: []string{"--sort-order", "3.5"},
			input: `{"initiativeId":"roadmap-initiative-id","projectId":"project-id","sortOrder":3.5}`,
		},
		{name: "omitted", input: `{"initiativeId":"roadmap-initiative-id","projectId":"project-id"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			capture := &initiativeLinkCaptureClient{delegate: commandFlowFakeClient{}}
			restore := useCommandRuntime(t, capture)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&bytes.Buffer{})
			command.SetArgs(append(
				[]string{"project", "initiative-link", "project-id", "roadmap-initiative-id"}, test.extraArgs...,
			))

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.NotNil(t, capture.create)
			payload, marshalErr := json.Marshal(capture.create.Variables)
			require.NoError(t, marshalErr)
			var sent struct {
				Input json.RawMessage `json:"input"`
			}
			require.NoError(t, json.Unmarshal(payload, &sent))
			require.JSONEq(t, test.input, string(sent.Input))
		})
	}
}
