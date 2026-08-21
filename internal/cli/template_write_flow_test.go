package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

const cliTemplateID = "7c9e6679-7425-40de-944b-e07fc1f90ae7"

func cliIssueTemplateJSON(name string, data string) string {
	return `{
		"id":"` + cliTemplateID + `",
		"name":"` + name + `",
		"type":"issue",
		"templateData":` + data + `,
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"pipeline":null
	}`
}

func templateWriteFake() graphqlPayloadOverride {
	entity := cliIssueTemplateJSON("Bug report", `{"title":"Bug"}`)

	return graphqlPayloadOverride{
		inner: commandFlowFakeClient{},
		payloads: map[string]string{
			"TemplateCreate":  `{"templateCreate":{"success":true,"template":` + entity + `}}`,
			"TemplateUpdate":  `{"templateUpdate":{"success":true,"template":` + entity + `}}`,
			"templateContent": `{"template":` + entity + `}`,
		},
	}
}

func writeTemplateDataFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "template.json")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

	return path
}

func Test_TemplateWriteCommands_have_no_bypass_or_scope_flags(t *testing.T) {
	root := NewRootCommand(context.Background(), BuildInfo{})
	for _, path := range []string{"template create", "template update"} {
		command, _, err := root.Find(strings.Fields(path))
		require.NoError(t, err)
		require.Nil(t, command.Flags().Lookup("force"))
		require.Nil(t, command.Flags().Lookup("confirm"))
		require.Nil(t, command.Flags().Lookup("org-wide"))
		require.Nil(t, command.Flags().Lookup("pipeline-id"))
		require.Nil(t, command.Flags().Lookup("team-id"))
	}
	update, _, err := root.Find([]string{"template", "update"})
	require.NoError(t, err)
	require.Nil(t, update.Flags().Lookup("type"))
	require.Nil(t, update.Flags().Lookup("id"))
}

func Test_TemplateWriteCommandFlows_cover_output_modes(t *testing.T) {
	dataFile := writeTemplateDataFile(t, `{"title":"Bug"}`)
	createArgs := []string{
		"template", "create",
		"--id", cliTemplateID,
		"--name", "Bug report",
		"--type", "issue",
		"--data-file", dataFile,
	}
	updateArgs := []string{"template", "update", cliTemplateID, "--name", "Bug report"}
	contentArgs := []string{"template", "content", cliTemplateID}
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "create human", args: createArgs, want: cliTemplateID + " Bug report [issue] team LIT\n"},
		{name: "create json", args: append([]string{"--json", "--compact"}, createArgs...), want: `"data":{"title":"Bug"}`},
		{name: "create id-only", args: append([]string{"--id-only"}, createArgs...), want: cliTemplateID + "\n"},
		{name: "create quiet", args: append([]string{"--quiet"}, createArgs...), want: ""},
		{name: "update human", args: updateArgs, want: cliTemplateID + " Bug report [issue] team LIT\n"},
		{
			name: "update data-file",
			args: []string{"template", "update", cliTemplateID, "--data-file", dataFile},
			want: cliTemplateID + " Bug report [issue] team LIT\n",
		},
		{name: "content json", args: append([]string{"--json", "--fields", "id,name,type,data"}, contentArgs...), want: `"data"`},
		{name: "content human", args: contentArgs, want: cliTemplateID + " Bug report [issue] team LIT\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, templateWriteFake())
			defer restore()
			var stdout bytes.Buffer
			err := execute(context.Background(), BuildInfo{}, strings.NewReader(""), &stdout, &bytes.Buffer{}, test.args)

			require.NoError(t, err)
			require.Contains(t, stdout.String(), test.want)
		})
	}
}

func Test_TemplateWriteCommandFlows_report_runtime_and_writer_errors(t *testing.T) {
	dataFile := writeTemplateDataFile(t, `{"title":"Bug"}`)
	commands := [][]string{
		{"template", "create", "--id", cliTemplateID, "--name", "Bug report", "--type", "issue", "--data-file", dataFile},
		{"template", "update", cliTemplateID, "--name", "Bug report"},
		{"template", "content", cliTemplateID},
	}
	for _, args := range commands {
		t.Run("runtime "+strings.Join(args[:2], " "), func(t *testing.T) {
			original := buildCommandRuntime
			buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			}
			defer func() { buildCommandRuntime = original }()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "runtime failed")
		})
		t.Run("writer "+strings.Join(args[:2], " "), func(t *testing.T) {
			restore := useCommandRuntime(t, templateWriteFake())
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(commandFailingWriter{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "write failed")
		})
	}
}

func Test_TemplateUpdate_rejects_data_file_stdin_before_mutation(t *testing.T) {
	fake := &mutationCallCounter{inner: templateWriteFake()}
	restore := useCommandRuntime(t, fake)
	defer restore()
	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(""),
		&bytes.Buffer{},
		&bytes.Buffer{},
		[]string{"template", "update", cliTemplateID, "--data-file", "-"},
	)

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Zero(t, fake.count("TemplateUpdate"))
}

func Test_TemplateCreate_rejects_data_file_failures_before_mutation(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{name: "stdin", path: "-"},
		{name: "missing", path: filepath.Join(t.TempDir(), "missing.json")},
		{name: "array", path: writeTemplateDataFile(t, `[]`)},
		{name: "scalar", path: writeTemplateDataFile(t, `123`)},
		{name: "null", path: writeTemplateDataFile(t, `null`)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fake := &mutationCallCounter{inner: templateWriteFake()}
			restore := useCommandRuntime(t, fake)
			defer restore()
			err := execute(
				context.Background(),
				BuildInfo{},
				strings.NewReader(""),
				&bytes.Buffer{},
				&bytes.Buffer{},
				[]string{
					"template", "create",
					"--id", cliTemplateID,
					"--name", "Bug report",
					"--type", "issue",
					"--data-file", test.path,
				},
			)

			require.ErrorIs(t, err, client.ErrWriteInvalid)
			require.Zero(t, fake.count("TemplateCreate"))
		})
	}
}

func Test_TemplateGet_stays_data_free(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	var stdout bytes.Buffer
	err := execute(
		context.Background(),
		BuildInfo{},
		strings.NewReader(""),
		&stdout,
		&bytes.Buffer{},
		[]string{"--json", "--compact", "template", "get", "template-id"},
	)

	require.NoError(t, err)
	require.NotContains(t, stdout.String(), `"data"`)
	require.NotContains(t, stdout.String(), "templateData")
}

type mutationCallCounter struct {
	inner graphqlPayloadOverride
	ops   []string
}

func (client *mutationCallCounter) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	client.ops = append(client.ops, request.OpName)

	return client.inner.MakeRequest(ctx, request, response)
}

func (client *mutationCallCounter) count(opName string) int {
	total := 0
	for _, operation := range client.ops {
		if operation == opName {
			total++
		}
	}

	return total
}
