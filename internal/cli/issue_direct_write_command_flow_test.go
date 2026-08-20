package cli

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

type expectedIssueWriteNumber struct {
	path  []string
	value float64
}

type issueWriteCaptureClient struct {
	directWriteCaptureClient
	numbers   []expectedIssueWriteNumber
	stateType string
}

func (client *issueWriteCaptureClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName == "WorkflowStatesByType" && client.stateType != "" {
		if err := requireRequestVariable(request, []string{"stateType"}, client.stateType, "state type"); err != nil {
			return err
		}
	}
	if request.OpName == client.operation {
		for _, number := range client.numbers {
			actual, err := requestVariable[float64](request, number.path...)
			if err != nil {
				return err
			}
			if actual != number.value {
				return fmt.Errorf("%v = %v", number.path, actual)
			}
		}
	}

	return client.directWriteCaptureClient.MakeRequest(ctx, request, response)
}

func Test_IssueDirectWriteCommandFlows_forward_mutation_variables(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		operation string
		variables []expectedWriteVariable
		numbers   []expectedIssueWriteNumber
		stateType string
	}{
		{
			name:      "create",
			args:      []string{"issue", "create", "--title", "Direct create"},
			operation: "IssueCreate",
			variables: []expectedWriteVariable{
				{path: []string{"input", "title"}, value: "Direct create"},
				{path: []string{"input", "teamId"}, value: "team-id"},
			},
		},
		{
			name: "update",
			args: []string{
				"issue", "update", "LIT-1", "--title", "Direct update",
				"--state", "done", "--priority", "urgent", "--estimate", "8",
			},
			operation: "IssueUpdate",
			variables: []expectedWriteVariable{
				{path: []string{"id"}, value: "LIT-1"},
				{path: []string{"input", "title"}, value: "Direct update"},
				{path: []string{"input", "stateId"}, value: "done-state"},
			},
			numbers: []expectedIssueWriteNumber{
				{path: []string{"input", "priority"}, value: 1},
				{path: []string{"input", "estimate"}, value: 8},
			},
		},
		{
			name:      "start",
			args:      []string{"issue", "start", "LIT-1"},
			operation: "IssueUpdate",
			variables: []expectedWriteVariable{
				{path: []string{"id"}, value: "LIT-1"},
				{path: []string{"input", "assigneeId"}, value: "user-id"},
				{path: []string{"input", "stateId"}, value: "type-state-id"},
			},
		},
		{
			name:      "comment",
			args:      []string{"issue", "comment", "LIT-1", "--body", "Direct comment"},
			operation: "IssueCommentCreate",
			variables: []expectedWriteVariable{
				{path: []string{"input", "issueId"}, value: "LIT-1"},
				{path: []string{"input", "body"}, value: "Direct comment"},
			},
		},
		{
			name:      "reply",
			args:      []string{"issue", "reply", "LIT-1", "comment-id", "--body", "Direct reply"},
			operation: "IssueCommentCreate",
			variables: []expectedWriteVariable{
				{path: []string{"input", "issueId"}, value: "LIT-1"},
				{path: []string{"input", "parentId"}, value: "comment-id"},
			},
		},
		{
			name:      "close",
			args:      []string{"issue", "close", "LIT-1"},
			operation: "IssueClose",
			variables: []expectedWriteVariable{
				{path: []string{"id"}, value: "LIT-1"},
				{path: []string{"input", "stateId"}, value: "done-state"},
			},
		},
		{
			name: "link",
			args: []string{
				"issue", "link", "https://example.com/direct", "LIT-1",
				"--title", "Direct link", "--subtitle", "review",
			},
			operation: "AttachmentLinkURL",
			variables: []expectedWriteVariable{
				{path: []string{"input", "issueId"}, value: "issue-id"},
				{path: []string{"input", "url"}, value: "https://example.com/direct"},
				{path: []string{"input", "title"}, value: "Direct link"},
				{path: []string{"input", "subtitle"}, value: "review"},
			},
		},
		{
			name:      "relate",
			args:      []string{"issue", "relate", "LIT-1", "LIT-2", "--type", "related"},
			operation: "IssueRelationCreate",
			variables: []expectedWriteVariable{
				{path: []string{"input", "issueId"}, value: "issue-id"},
				{path: []string{"input", "relatedIssueId"}, value: "related-issue-id"},
				{path: []string{"input", "type"}, value: "related"},
			},
		},
		{
			name:      "unrelate",
			args:      []string{"issue", "unrelate", "issue-relation-id"},
			operation: "IssueRelationDelete",
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "issue-relation-id"}},
		},
		{
			name:      "add label",
			args:      []string{"issue", "add-label", "LIT-1", "label-id"},
			operation: "IssueAddLabel",
			variables: []expectedWriteVariable{
				{path: []string{"id"}, value: "issue-id"},
				{path: []string{"labelId"}, value: "label-id"},
			},
		},
		{
			name:      "remove label",
			args:      []string{"issue", "remove-label", "LIT-1", "label-id"},
			operation: "IssueRemoveLabel",
			variables: []expectedWriteVariable{
				{path: []string{"id"}, value: "issue-id"},
				{path: []string{"labelId"}, value: "label-id"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			flow := commandFlowFakeClient{}
			after := ""
			relation := ""
			flow.afterIssue = &after
			flow.afterRelation = &relation
			fake := &issueWriteCaptureClient{
				directWriteCaptureClient: directWriteCaptureClient{
					operation: test.operation,
					variables: test.variables,
					delegate:  flow,
				},
				numbers:   test.numbers,
				stateType: test.stateType,
			}
			restore := useCommandRuntime(t, fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&bytes.Buffer{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, 1, fake.calls)
		})
	}
}

func Test_IssueDirectLabelWrites_propagate_mutation_errors(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		operation string
	}{
		{
			name:      "add label",
			args:      []string{"issue", "add-label", "LIT-1", "label-id"},
			operation: "IssueAddLabel",
		},
		{
			name:      "remove label",
			args:      []string{"issue", "remove-label", "LIT-1", "label-id"},
			operation: "IssueRemoveLabel",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: test.operation})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "operation failed")
		})
	}
}
