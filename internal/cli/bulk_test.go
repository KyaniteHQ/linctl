package cli

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func runBulkFlow(t *testing.T, fake graphql.Client, args []string) (stdout string, err error) {
	t.Helper()
	restore := useCommandRuntime(t, fake)
	defer restore()
	outBuf := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&outBuf)
	command.SetArgs(args)

	err = command.ExecuteContext(context.Background())

	return outBuf.String(), err
}

func writeImportFile(t *testing.T, name string, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}

func writeImportDryRunConfig(t *testing.T, content string) {
	t.Helper()
	dir := t.TempDir()
	t.Chdir(dir)
	configureTestAuthEnvironment(t)
	require.NoError(t, os.WriteFile(".linctl.toml", []byte(content), 0o600))
}

func writePinnedImportDryRunConfig(t *testing.T) {
	t.Helper()

	writeImportDryRunConfig(t, `
[target]
org_id = "org-id"
team_key = "LIT"
team_id = "team-id"
`)
}

func Test_CommandFlows_issue_import_accepts_an_uppercase_extension(t *testing.T) {
	path := writeImportFile(t, "rows.JSON", `[{"title":"First"}]`)

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path})

	require.NoError(t, err)
	require.Equal(t, "LIT-2 Created issue [Todo]\n", stdout)
}

type importTitleCaptureClient struct {
	mu     sync.Mutex
	titles []string
}

func (capture *importTitleCaptureClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName == "IssueCreate" {
		title, err := requestVariable[string](request, "input", "title")
		if err != nil {
			return err
		}
		capture.mu.Lock()
		capture.titles = append(capture.titles, title)
		capture.mu.Unlock()
	}

	return commandFlowFakeClient{}.MakeRequest(ctx, request, response)
}

func Test_CommandFlows_issue_import_creates_from_json(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First"},{"team":"LIT","title":"Second"}]`)
	capture := &importTitleCaptureClient{}

	stdout, err := runBulkFlow(t, capture, []string{"issue", "import", path})

	require.NoError(t, err)
	require.Equal(t, "LIT-2 Created issue [Todo]\nLIT-2 Created issue [Todo]\n", stdout)
	require.ElementsMatch(t, []string{"First", "Second"}, capture.titles)
}

func Test_CommandFlows_issue_import_creates_from_csv(t *testing.T) {
	path := writeImportFile(t, "rows.csv", "title,team\nFirst,LIT\n")

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path})

	require.NoError(t, err)
	require.Contains(t, stdout, "LIT-2 Created issue [Todo]")
}

func Test_CommandFlows_issue_import_strips_utf8_bom(t *testing.T) {
	// A spreadsheet-exported CSV carries a leading UTF-8 BOM; without stripping it
	// the first header cell would otherwise be read as BOM+"title" and every row fails validation.
	path := writeImportFile(t, "rows.csv", "\ufeff"+"title,team\nFirst,LIT\n")

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path})

	require.NoError(t, err)
	require.Contains(t, stdout, "LIT-2 Created issue [Todo]")
}

func Test_CommandFlows_issue_import_json_result(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--json", "issue", "import", path})

	require.NoError(t, err)
	require.Contains(t, stdout, `"count": 1`)
	require.Contains(t, stdout, `"identifier": "LIT-2"`)
}

func Test_CommandFlows_issue_import_json_fields_projects_issues(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)

	stdout, err := runBulkFlow(
		t, commandFlowFakeClient{}, []string{"--json", "--fields", "identifier", "issue", "import", path},
	)

	require.NoError(t, err)
	require.Contains(t, stdout, `"identifier": "LIT-2"`)
	require.NotContains(t, stdout, `"count"`)
}

func Test_CommandFlows_issue_import_quiet_result(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--quiet", "issue", "import", path})

	require.NoError(t, err)
	require.Empty(t, stdout)
}

func Test_CommandFlows_issue_import_dry_run_previews(t *testing.T) {
	writePinnedImportDryRunConfig(t)
	path := writeImportFile(t, "rows.json", `[{"title":"First","priority":"high","state":"  In Review  "}]`)

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path, "--dry-run"})
	require.NoError(t, err)
	require.Contains(t, stdout, `would create "First" state=In Review priority=2`)

	jsonOut, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--json", "issue", "import", path, "--dry-run"})
	require.NoError(t, err)
	require.Contains(t, jsonOut, `"dry_run": true`)
	require.Contains(t, jsonOut, `"state": "In Review"`)
	require.NotContains(t, jsonOut, `"state_type"`)

	quiet, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--quiet", "issue", "import", path, "--dry-run"})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_CommandFlows_issue_import_dry_run_does_not_build_auth_runtime(t *testing.T) {
	writePinnedImportDryRunConfig(t)
	path := writeImportFile(t, "rows.csv", "title,team\nFirst,LIT\n")
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime should not be built")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	output := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"issue", "import", path, "--dry-run"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `would create "First"`)
}

func Test_CommandFlows_issue_import_dry_run_rejects_team_mismatch_locally(t *testing.T) {
	writePinnedImportDryRunConfig(t)
	path := writeImportFile(t, "rows.json", `[{"team":"OTH","title":"First"}]`)
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime should not be built")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "import", path, "--dry-run"})

	err := command.ExecuteContext(context.Background())

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Equal(t, "INVALID_WRITE", errorCode(err))
	require.ErrorContains(t, err, "does not match pinned target team")
}

func Test_CommandFlows_issue_import_dry_run_requires_local_team_key(t *testing.T) {
	writeImportDryRunConfig(t, `
[target]
org_id = "org-id"
team_id = "team-id"
`)
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "import", path, "--dry-run"})

	err := command.ExecuteContext(context.Background())

	require.ErrorIs(t, err, client.ErrTargetNotConfigured)
	require.ErrorContains(t, err, "set team_key")
}

func Test_CommandFlows_issue_import_dry_run_surfaces_local_format_and_read_errors(t *testing.T) {
	writePinnedImportDryRunConfig(t)

	_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", "rows.txt", "--dry-run"})
	require.ErrorContains(t, err, "unsupported data format")

	_, err = runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", "missing.json", "--dry-run"})
	require.ErrorContains(t, err, "read missing.json")
}

func Test_CommandFlows_issue_import_dry_run_surfaces_config_load_error(t *testing.T) {
	writeImportDryRunConfig(t, "[")
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)

	_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path, "--dry-run"})

	require.ErrorContains(t, err, "parse config")
}

func Test_CommandFlows_issue_import_non_dry_run_still_uses_runtime(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "import", path})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "runtime failed")
}

func Test_CommandFlows_issue_import_rejects_team_mismatch(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"team":"OTH","title":"First"}]`)

	_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path})

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.Equal(t, "INVALID_WRITE", errorCode(err))
	require.ErrorContains(t, err, "does not match pinned target team")
}

func Test_CommandFlows_issue_import_surfaces_input_errors(t *testing.T) {
	cases := map[string]string{
		"missing title":    `[{"title":"  "}]`,
		"invalid state":    `[{"title":"First","state":"nope"}]`,
		"invalid priority": `[{"title":"First","priority":"highest"}]`,
		"invalid json":     `{not json`,
	}
	for name, content := range cases {
		t.Run(name, func(t *testing.T) {
			path := writeImportFile(t, "rows.json", content)

			_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path})

			require.Error(t, err)
			if name == "missing title" {
				require.ErrorIs(t, err, client.ErrWriteInvalid)
				require.Equal(t, "INVALID_WRITE", errorCode(err))
			}
		})
	}
}

func Test_CommandFlows_issue_import_refuses_a_state_type_as_a_state_name(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First","state":"completed"}]`)

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path})

	require.ErrorContains(t, err, `unknown workflow state "completed"`)
	require.Equal(t, "INVALID_WRITE", errorCode(err))
	require.Empty(t, stdout)
}

func Test_CommandFlows_issue_import_selects_exact_state_name(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First","state":"In Review"}]`)
	fake := &issueWriteCaptureClient{
		directWriteCaptureClient: directWriteCaptureClient{
			operation: "IssueCreate",
			variables: []expectedWriteVariable{
				{path: []string{"input", "title"}, value: "First"},
				{path: []string{"input", "stateId"}, value: "in-review-state"},
			},
			delegate: commandFlowFakeClient{},
		},
	}

	stdout, err := runBulkFlow(t, fake, []string{"issue", "import", path})

	require.NoError(t, err)
	require.Equal(t, 1, fake.calls)
	require.Contains(t, stdout, "LIT-2")
}

func Test_CommandFlows_issue_import_surfaces_csv_errors(t *testing.T) {
	ragged := writeImportFile(t, "rows.csv", "title,team\nFirst\n")
	_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", ragged})
	require.Error(t, err)

	empty := writeImportFile(t, "rows.csv", "")
	_, err = runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", empty})
	require.Error(t, err)
}

func Test_CommandFlows_issue_import_surfaces_read_and_format_errors(t *testing.T) {
	_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", "missing.json"})
	require.Error(t, err)

	_, err = runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", "rows.txt"})
	require.Error(t, err)
}

func Test_CommandFlows_issue_import_surfaces_create_errors(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)

	_, err := runBulkFlow(
		t,
		commandFlowFakeClient{failOperation: "IssueCreate"},
		[]string{"issue", "import", path},
	)

	require.Error(t, err)
}

func Test_runIssueImport_surfaces_preview_writer_errors(t *testing.T) {
	writePinnedImportDryRunConfig(t)
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"issue", "import", path, "--dry-run"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
}

func Test_CommandFlows_issue_bulk_export_writes_file(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "issues.json")

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "bulk-export", jsonPath})
	require.NoError(t, err)
	require.Contains(t, stdout, jsonPath+" (1 issues)")

	data, err := os.ReadFile(jsonPath) //nolint:gosec // G304: test-controlled temp path.
	require.NoError(t, err)
	require.Contains(t, string(data), `"identifier": "LIT-1"`)
	if runtime.GOOS != "windows" {
		info, err := os.Stat(jsonPath)
		require.NoError(t, err)
		require.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}

	csvPath := filepath.Join(dir, "issues.csv")
	_, err = runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "bulk-export", csvPath})
	require.NoError(t, err)

	csvData, err := os.ReadFile(csvPath) //nolint:gosec // G304: test-controlled temp path.
	require.NoError(t, err)
	records, err := csv.NewReader(bytes.NewReader(csvData)).ReadAll()
	require.NoError(t, err)
	require.Equal(t, [][]string{
		{"identifier", "title", "state", "priority", "assignee", "project", "url"},
		{"LIT-1", "Listed issue", "Todo", "No priority", "", "Pinned project", "https://linear.app/kyanite/issue/LIT-1"},
	}, records)
}

func Test_CommandFlows_issue_bulk_export_honors_output_flags(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "issues.json")

	idOnly, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--id-only", "issue", "bulk-export", jsonPath})
	require.NoError(t, err)
	require.Equal(t, jsonPath+"\n", idOnly)

	out, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--json", "issue", "bulk-export", jsonPath})
	require.NoError(t, err)
	require.Contains(t, out, `"count": 1`)

	quiet, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--quiet", "issue", "bulk-export", jsonPath})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_CommandFlows_issue_bulk_export_surfaces_errors(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "issues.json")

	_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "bulk-export", "out.txt"})
	require.Error(t, err)

	_, err = runBulkFlow(
		t,
		commandFlowFakeClient{failOperation: "Viewer"},
		[]string{"issue", "bulk-export", jsonPath},
	)
	require.Error(t, err)

	_, err = runBulkFlow(
		t,
		commandFlowFakeClient{failOperation: "IssuesByTeamFiltered"},
		[]string{"issue", "bulk-export", jsonPath},
	)
	require.Error(t, err)

	_, err = runBulkFlow(
		t,
		commandFlowFakeClient{},
		[]string{"issue", "bulk-export", filepath.Join(dir, "missing-dir", "issues.json")},
	)
	require.Error(t, err)
}

type formulaIssueListClient struct {
	delegate graphql.Client
}

func (fake formulaIssueListClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName != "IssuesByTeamFiltered" {
		return fake.delegate.MakeRequest(ctx, request, response)
	}
	payload, err := json.Marshal(map[string]any{"issues": map[string]any{
		"nodes": []map[string]any{{
			"id":            "issue-id",
			"identifier":    "=CMD",
			"title":         `=HYPERLINK("http://evil/","click")`,
			"branchName":    "lit-1-formula",
			"url":           "\rCR",
			"priority":      0,
			"priorityLabel": "-1+1",
			"team":          map[string]any{"id": "team-id", "key": "LIT", "name": "LIT"},
			"state":         map[string]any{"id": "todo-state", "name": "+cmd", "type": "unstarted"},
			"assignee":      map[string]any{"id": "user-id", "name": "omer", "displayName": "@SUM(1)"},
			"project":       map[string]any{"id": "project-id", "name": "\tTAB"},
		}},
		"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
	}})
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(`{"data":`+string(payload)+`}`), response)
}

func Test_CommandFlows_issue_bulk_export_neutralizes_csv_formula_cells(t *testing.T) {
	path := filepath.Join(t.TempDir(), "issues.csv")

	_, err := runBulkFlow(t, formulaIssueListClient{delegate: commandFlowFakeClient{}}, []string{"issue", "bulk-export", path})

	require.NoError(t, err)
	data, err := os.ReadFile(path) //nolint:gosec // G304: test-controlled temp path.
	require.NoError(t, err)
	records, err := csv.NewReader(strings.NewReader(string(data))).ReadAll()
	require.NoError(t, err)
	require.Equal(t, [][]string{
		{"identifier", "title", "state", "priority", "assignee", "project", "url"},
		{"'=CMD", `'=HYPERLINK("http://evil/","click")`, "'+cmd", "'-1+1", "'@SUM(1)", "'\tTAB", "'\rCR"},
	}, records)
}
