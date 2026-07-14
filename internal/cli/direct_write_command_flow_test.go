package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

type expectedWriteVariable struct {
	path  []string
	value string
}

type directWriteCaptureClient struct {
	operation string
	variables []expectedWriteVariable
	delegate  graphql.Client
	calls     int
}

func (client *directWriteCaptureClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName == client.operation {
		client.calls++
		for _, variable := range client.variables {
			if err := requireRequestVariable(request, variable.path, variable.value, strings.Join(variable.path, ".")); err != nil {
				return err
			}
		}
	}

	return client.delegate.MakeRequest(ctx, request, response)
}

func Test_DirectWriteCommandFlows_forward_request_variables(t *testing.T) {
	documentContentFile := writeTempTextFile(t, "document content from file")
	projectContentFile := writeTempTextFile(t, "project content from file")
	projectUpdateBodyFile := writeTempTextFile(t, "project update from file")
	tests := []struct {
		name      string
		args      []string
		stdin     string
		operation string
		variables []expectedWriteVariable
		cycle     bool
	}{
		{
			name:      "comment update body",
			args:      []string{"comment", "update", "comment-id", "--body", "updated body"},
			operation: "CommentUpdate",
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "comment-id"}, {path: []string{"input", "body"}, value: "updated body"}},
		},
		{
			name:      "comment delete id",
			args:      []string{"comment", "delete", "comment-id"},
			operation: "CommentDelete",
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "comment-id"}},
		},
		{
			name:      "comment resolve id",
			args:      []string{"comment", "resolve", "comment-id"},
			operation: "CommentResolve",
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "comment-id"}},
		},
		{
			name:      "comment unresolve id",
			args:      []string{"comment", "unresolve", "comment-id"},
			operation: "CommentUnresolve",
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "comment-id"}},
		},
		{
			name: "Cycle create fields",
			args: []string{
				"cycle", "create", "--name", "July", "--description", "delivery",
				"--starts-at", "2026-07-01T00:00:00Z", "--ends-at", "2026-07-15T00:00:00Z",
				"--completed-at", "2026-07-14T00:00:00Z",
			},
			operation: "CycleCreate",
			cycle:     true,
			variables: []expectedWriteVariable{
				{path: []string{"input", "name"}, value: "July"},
				{path: []string{"input", "description"}, value: "delivery"},
				{path: []string{"input", "startsAt"}, value: "2026-07-01T00:00:00Z"},
				{path: []string{"input", "endsAt"}, value: "2026-07-15T00:00:00Z"},
				{path: []string{"input", "completedAt"}, value: "2026-07-14T00:00:00Z"},
			},
		},
		{
			name: "Cycle update id and fields",
			args: []string{
				"cycle", "update", "cycle-id", "--name", "Late July", "--description", "shipped",
				"--starts-at", "2026-07-02T00:00:00Z", "--ends-at", "2026-07-16T00:00:00Z",
				"--completed-at", "2026-07-15T00:00:00Z",
			},
			operation: "CycleUpdate",
			cycle:     true,
			variables: []expectedWriteVariable{
				{path: []string{"id"}, value: "cycle-id"},
				{path: []string{"input", "name"}, value: "Late July"},
				{path: []string{"input", "description"}, value: "shipped"},
				{path: []string{"input", "startsAt"}, value: "2026-07-02T00:00:00Z"},
				{path: []string{"input", "endsAt"}, value: "2026-07-16T00:00:00Z"},
				{path: []string{"input", "completedAt"}, value: "2026-07-15T00:00:00Z"},
			},
		},
		{
			name:      "Cycle archive id",
			args:      []string{"cycle", "archive", "cycle-id"},
			operation: "CycleArchive",
			cycle:     true,
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "cycle-id"}},
		},
		{
			name:      "document create content from stdin",
			args:      []string{"document", "create", "--title", "Runbook", "--content", "-"},
			stdin:     "document content from stdin",
			operation: "DocumentCreate",
			variables: []expectedWriteVariable{{path: []string{"input", "title"}, value: "Runbook"}, {path: []string{"input", "content"}, value: "document content from stdin"}},
		},
		{
			name:      "document update id and content from file",
			args:      []string{"document", "update", "document-id", "--title", "Updated runbook", "--content-file", documentContentFile},
			operation: "DocumentUpdate",
			variables: []expectedWriteVariable{
				{path: []string{"id"}, value: "document-id"},
				{path: []string{"input", "title"}, value: "Updated runbook"},
				{path: []string{"input", "content"}, value: "document content from file"},
			},
		},
		{
			name:      "project create name and content from file",
			args:      []string{"project", "create", "--name", "Launch", "--description", "release", "--content-file", projectContentFile},
			operation: "ProjectCreate",
			variables: []expectedWriteVariable{
				{path: []string{"input", "name"}, value: "Launch"},
				{path: []string{"input", "description"}, value: "release"},
				{path: []string{"input", "content"}, value: "project content from file"},
			},
		},
		{
			name:      "project update id name and content from file",
			args:      []string{"project", "update", "project-id", "--name", "Launch two", "--content-file", projectContentFile},
			operation: "ProjectUpdate",
			variables: []expectedWriteVariable{
				{path: []string{"id"}, value: "project-id"},
				{path: []string{"input", "name"}, value: "Launch two"},
				{path: []string{"input", "content"}, value: "project content from file"},
			},
		},
		{
			name:      "project archive id",
			args:      []string{"project", "archive", "project-id"},
			operation: "ProjectArchive",
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "project-id"}},
		},
		{
			name:      "project add-label ids",
			args:      []string{"project", "add-label", "project-id", "project-label-id"},
			operation: "ProjectAddLabel",
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "project-id"}, {path: []string{"labelId"}, value: "project-label-id"}},
		},
		{
			name:      "project remove-label ids",
			args:      []string{"project", "remove-label", "project-id", "project-label-id"},
			operation: "ProjectRemoveLabel",
			variables: []expectedWriteVariable{{path: []string{"id"}, value: "project-id"}, {path: []string{"labelId"}, value: "project-label-id"}},
		},
		{
			name:      "project update body from stdin and canonical health",
			args:      []string{"project-update", "create", "project-id", "--health", "at-risk", "--body", "-"},
			stdin:     "project update from stdin",
			operation: "ProjectUpdateCreate",
			variables: []expectedWriteVariable{
				{path: []string{"input", "projectId"}, value: "project-id"},
				{path: []string{"input", "body"}, value: "project update from stdin"},
				{path: []string{"input", "health"}, value: "atRisk"},
			},
		},
		{
			name:      "project update body from file",
			args:      []string{"project-update", "create", "project-id", "--body-file", projectUpdateBodyFile},
			operation: "ProjectUpdateCreate",
			variables: []expectedWriteVariable{
				{path: []string{"input", "projectId"}, value: "project-id"},
				{path: []string{"input", "body"}, value: "project update from file"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			delegate := graphql.Client(commandFlowFakeClient{})
			if test.cycle {
				delegate = cycleCommandFlowFakeClient{}
			}
			fake := &directWriteCaptureClient{
				operation: test.operation,
				variables: test.variables,
				delegate:  delegate,
			}
			restore := useCommandRuntime(t, fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&bytes.Buffer{})
			if test.stdin != "" {
				command.SetIn(strings.NewReader(test.stdin))
			}
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, 1, fake.calls)
		})
	}
}

func Test_DirectProjectLabelCommandFlows_propagate_client_errors(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		operation string
		message   string
	}{
		{
			name:      "add",
			args:      []string{"project", "add-label", "project-id", "project-label-id"},
			operation: "ProjectAddLabel",
			message:   "add label to project project-id",
		},
		{
			name:      "remove",
			args:      []string{"project", "remove-label", "project-id", "project-label-id"},
			operation: "ProjectRemoveLabel",
			message:   "remove label from project project-id",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: test.operation})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, test.message)
		})
	}
}
