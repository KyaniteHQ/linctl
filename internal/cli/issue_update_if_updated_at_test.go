package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

// issueReadClient answers the pre-write issue read with a fixed body and counts
// every request, so a test can prove that a refused update sent nothing.
type issueReadClient struct {
	delegate graphql.Client
	issue    string
	requests int
}

func (client *issueReadClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	client.requests++
	if request.OpName == "issue" {
		return json.Unmarshal([]byte(`{"data":{"issue":`+client.issue+`}}`), response)
	}

	return client.delegate.MakeRequest(ctx, request, response)
}

func issueBodyWithUpdatedAt(updatedAt string, teamID string, teamKey string) string {
	return strings.Replace(
		commandIssueJSONWithTeam("LIT-1", "Existing", "todo-state", "Todo", "unstarted", teamID, teamKey),
		`"assignee":null,`, `"assignee":null,"updatedAt":"`+updatedAt+`",`, 1,
	)
}

func Test_IssueUpdate_if_updated_at_guards_the_write(t *testing.T) {
	const current = "2026-09-16T14:26:12.000Z"
	tests := []struct {
		name         string
		issue        string
		ifUpdatedAt  string
		code         string
		wantMutation int
	}{
		{
			name:        "stale timestamp is a conflict",
			issue:       issueBodyWithUpdatedAt("2026-09-16T14:30:00.000Z", "team-id", "LIT"),
			ifUpdatedAt: "2026-09-16T14:26:12Z",
			code:        "CONFLICT",
		},
		{
			name:         "same instant at another precision writes",
			issue:        issueBodyWithUpdatedAt(current, "team-id", "LIT"),
			ifUpdatedAt:  "2026-09-16T14:26:12Z",
			wantMutation: 1,
		},
		{
			name:        "malformed timestamp is invalid before any request",
			issue:       issueBodyWithUpdatedAt(current, "team-id", "LIT"),
			ifUpdatedAt: "yesterday",
			code:        "INVALID_WRITE",
		},
		{
			name:        "target mismatch wins over a stale timestamp",
			issue:       issueBodyWithUpdatedAt("2026-09-16T14:30:00.000Z", "ops-team-id", "OPS"),
			ifUpdatedAt: "2026-09-16T14:26:12Z",
			code:        "TARGET_MISMATCH",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := &issueReadClient{delegate: commandFlowFakeClient{}, issue: test.issue}
			fake := &directWriteCaptureClient{operation: "IssueUpdate", delegate: reader}
			restore := useCommandRuntime(t, fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&bytes.Buffer{})
			command.SetArgs([]string{
				"issue", "update", "LIT-1", "--title", "New title", "--if-updated-at", test.ifUpdatedAt,
			})

			err := command.ExecuteContext(context.Background())

			require.Equal(t, test.wantMutation, fake.calls)
			if test.code == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Equal(t, test.code, errorCode(err))
			if test.code == "CONFLICT" {
				require.Contains(t, err.Error(), "2026-09-16T14:30:00.000Z")
				require.Contains(t, err.Error(), "2026-09-16T14:26:12Z")
			}
			if test.code == "INVALID_WRITE" {
				require.Zero(t, reader.requests)
			}
		})
	}
}
