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
		sortOrder *float64
	}{
		{name: "given", extraArgs: []string{"--sort-order", "3.5"}, sortOrder: new(3.5)},
		{name: "omitted", sortOrder: nil},
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
				Input struct {
					SortOrder *float64 `json:"sortOrder"`
				} `json:"input"`
			}
			require.NoError(t, json.Unmarshal(payload, &sent))
			require.Equal(t, test.sortOrder, sent.Input.SortOrder)
		})
	}
}
