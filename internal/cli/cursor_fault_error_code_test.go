package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type graphqlPayloadOverride struct {
	inner    graphql.Client
	payloads map[string]string
}

func (client graphqlPayloadOverride) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	payload, ok := client.payloads[request.OpName]
	if !ok {
		return client.inner.MakeRequest(ctx, request, response)
	}

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

func Test_CursorFaultSites_report_stable_error_codes(t *testing.T) {
	nullCursorPageInfo := `"pageInfo":{"hasNextPage":true,"endCursor":null}`
	tests := []struct {
		name    string
		args    []string
		payload map[string]string
		code    string
		sent    error
	}{
		{
			name: "collectNodePages via project list",
			args: []string{"project", "list"},
			payload: map[string]string{
				"Projects": `{"team":{"projects":{"nodes":[` +
					commandProjectJSON("Listed project", "Backlog", "backlog") +
					`],` + nullCursorPageInfo + `}}}`,
			},
			code: "GRAPHQL_ERROR",
			sent: client.ErrGraphQL,
		},
		{
			name: "findTeamByKey via project add-team",
			args: []string{"project", "add-team", "project-id", "--to-team", "OPS"},
			payload: map[string]string{
				"teams_list": `{"teams":{"nodes":[{
					"id":"team-id","key":"LIT","name":"linctl",
					"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}
				}],` + nullCursorPageInfo + `}}`,
			},
			code: "GRAPHQL_ERROR",
			sent: client.ErrGraphQL,
		},
		{
			name: "listAllProjectTeamIDs via project add-team",
			args: []string{"project", "add-team", "project-id", "--to-team-id", "ops-team-id"},
			payload: map[string]string{
				"project": `{"project":` + truncatedProjectDetailJSON() + `}`,
				"project_teams": `{"project":{"id":"project-id","name":"Detail project","teams":{"nodes":[{
					"id":"team-id","key":"LIT","name":"linctl","description":"",
					"archivedAt":null,
					"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}
				}],` + nullCursorPageInfo + `}}}`,
			},
			code: "TARGET_MISMATCH",
			sent: client.ErrTargetMismatch,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restore := useCommandRuntime(t, graphqlPayloadOverride{
				inner:    commandFlowFakeClient{},
				payloads: tt.payload,
			})
			defer restore()
			var stdout, stderr bytes.Buffer
			err := execute(
				context.Background(),
				BuildInfo{},
				strings.NewReader(""),
				&stdout,
				&stderr,
				tt.args,
			)

			require.ErrorIs(t, err, tt.sent)
			require.Equal(t, tt.code, errorCode(err))
			require.Contains(t, stderr.String(), `"error_code":"`+tt.code+`"`)
		})
	}
}

func truncatedProjectDetailJSON() string {
	return `{
		"id":"project-id",
		"name":"Detail project",
		"description":"description",
		"content":"Existing project content",
		"archivedAt":null,
		"slugId":"Detail project",
		"url":"https://linear.app/kyanite/project/project-id",
		"priority":0,
		"status":{"id":"status-id","name":"Backlog","type":"backlog"},
		"lead":null,
		"teams":{
			"nodes":[{"id":"team-id","key":"LIT","name":"linctl"}],
			"pageInfo":{"hasNextPage":true}
		}
	}`
}
