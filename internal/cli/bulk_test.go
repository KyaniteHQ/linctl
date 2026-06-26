package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func runBulkFlow(t *testing.T, fake commandFlowFakeClient, args []string) (stdout string, err error) {
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

func Test_dataFormat_resolves_and_rejects(t *testing.T) {
	format, err := dataFormat("rows.JSON")
	require.NoError(t, err)
	//nolint:testifylint // formatJSON is the literal "json", not JSON content for JSONEq.
	require.Equal(t, formatJSON, format)

	format, err = dataFormat("rows.csv")
	require.NoError(t, err)
	require.Equal(t, formatCSV, format)

	_, err = dataFormat("rows.txt")
	require.Error(t, err)
}

func Test_CommandFlows_issue_import_creates_from_json(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First"},{"team":"LIT","title":"Second"}]`)

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path})

	require.NoError(t, err)
	require.Equal(t, "LIT-2 Created issue [Todo]\nLIT-2 Created issue [Todo]\n", stdout)
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

func Test_CommandFlows_issue_import_quiet_result(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First"}]`)

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--quiet", "issue", "import", path})

	require.NoError(t, err)
	require.Empty(t, stdout)
}

func Test_CommandFlows_issue_import_dry_run_previews(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"title":"First","priority":"high","state":"started"}]`)

	stdout, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path, "--dry-run"})
	require.NoError(t, err)
	require.Contains(t, stdout, `would create "First" state=started priority=2`)

	jsonOut, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--json", "issue", "import", path, "--dry-run"})
	require.NoError(t, err)
	require.Contains(t, jsonOut, `"dry_run": true`)
	require.Contains(t, jsonOut, `"state_type": "started"`)

	quiet, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"--quiet", "issue", "import", path, "--dry-run"})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_CommandFlows_issue_import_rejects_team_mismatch(t *testing.T) {
	path := writeImportFile(t, "rows.json", `[{"team":"OTH","title":"First"}]`)

	_, err := runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "import", path})

	require.Error(t, err)
	require.Contains(t, err.Error(), "does not match pinned target team")
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
		})
	}
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
	info, err := os.Stat(jsonPath)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), info.Mode().Perm())

	csvPath := filepath.Join(dir, "issues.csv")
	_, err = runBulkFlow(t, commandFlowFakeClient{}, []string{"issue", "bulk-export", csvPath})
	require.NoError(t, err)

	csvData, err := os.ReadFile(csvPath) //nolint:gosec // G304: test-controlled temp path.
	require.NoError(t, err)
	require.Contains(t, string(csvData), "identifier,title,state,priority,assignee,project,url")
	require.Contains(t, string(csvData), "LIT-1,Listed issue,Todo")
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
		commandFlowFakeClient{failOperation: "IssuesByTeam"},
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

func Test_encodeIssues_surfaces_writer_errors(t *testing.T) {
	issues := []client.IssueSummary{{Identifier: "LIT-1", Title: "First"}}

	require.Error(t, encodeIssues(commandFailingWriter{}, formatJSON, issues))
	require.Error(t, encodeIssues(commandFailingWriter{}, formatCSV, issues))
}
