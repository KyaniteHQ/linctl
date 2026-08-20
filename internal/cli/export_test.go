package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func runExportFlow(
	t *testing.T,
	fake commandFlowFakeClient,
	args []string,
) (stdout string, stderr string, err error) {
	t.Helper()
	restore := useCommandRuntime(t, fake)
	defer restore()
	outBuf := bytes.Buffer{}
	errBuf := bytes.Buffer{}
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&outBuf)
	command.SetErr(&errBuf)
	command.SetArgs(args)

	err = command.ExecuteContext(context.Background())

	return outBuf.String(), errBuf.String(), err
}

func Test_CommandFlows_issue_export_writes_document(t *testing.T) {
	dir := t.TempDir()

	stdout, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"issue", "export", "LIT-1", dir})
	require.NoError(t, err)

	path := filepath.Join(dir, "LIT-1.md")
	require.Contains(t, stdout, path)
	require.Contains(t, stdout, "(1 comments, 1 attachments)")

	data, err := os.ReadFile(path) //nolint:gosec // G304: path is the test-controlled temp dir.
	require.NoError(t, err)
	document := string(data)
	require.Contains(t, document, "# LIT-1 — Detail issue")
	require.Contains(t, document, "- URL: https://linear.app/kyanite/issue/LIT-1")
	require.Contains(t, document, "## Description\n\nExisting description")
	require.Contains(t, document, "## Comments (1)")
	require.Contains(t, document, "### Omer — 2026-06-19T12:00:00Z\n\nFirst comment")
	require.Contains(t, document, "## Attachments (1)")
	require.Contains(t, document, "- [Linked PR](<https://github.com/kyanite/linctl/pull/1>)")
}

func Test_CommandFlows_issue_export_honors_output_flags(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "LIT-1.md")

	idOnly, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"--id-only", "issue", "export", "LIT-1", dir})
	require.NoError(t, err)
	require.Equal(t, path+"\n", idOnly)

	jsonOut, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"--json", "issue", "export", "LIT-1", dir})
	require.NoError(t, err)
	require.Contains(t, jsonOut, `"identifier": "LIT-1"`)
	require.Contains(t, jsonOut, `"comments": 1`)

	quiet, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"--quiet", "issue", "export", "LIT-1", dir})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_CommandFlows_issue_export_notes_truncation(t *testing.T) {
	dir := t.TempDir()

	stdout, stderr, err := runExportFlow(
		t,
		commandFlowFakeClient{truncatedExport: true},
		[]string{"--json", "issue", "export", "LIT-1", dir},
	)
	require.NoError(t, err)
	require.Contains(t, stderr, "export capped at")
	require.Contains(t, stdout, `"truncated": true`)
}

func Test_CommandFlows_issue_export_surfaces_read_errors(t *testing.T) {
	for _, operation := range []string{"issue", "issue_comments", "issue_attachments"} {
		t.Run(operation, func(t *testing.T) {
			dir := t.TempDir()

			_, _, err := runExportFlow(
				t,
				commandFlowFakeClient{failOperation: operation},
				[]string{"issue", "export", "LIT-1", dir},
			)

			require.Error(t, err)
		})
	}
}

func Test_runIssueExport_surfaces_directory_errors(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-dir")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o600))

	_, _, err := runExportFlow(
		t,
		commandFlowFakeClient{},
		[]string{"issue", "export", "LIT-1", filepath.Join(file, "sub")},
	)

	require.Error(t, err)
}

func Test_runIssueExport_surfaces_note_writer_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{truncatedExport: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetErr(commandFailingWriter{})
	command.SetArgs([]string{"issue", "export", "LIT-1", t.TempDir()})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
}

func Test_writeExportDocument_refuses_identifier_outside_the_linear_grammar(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "export")

	_, err := writeExportDocument(dir, "../escape", "body")

	require.ErrorContains(t, err, "not a valid export filename")
	require.NoDirExists(t, dir)
}

func Test_issueExportLeaf_refuses_identifier_outside_the_linear_grammar(t *testing.T) {
	_, err := issueExportLeaf("../escape")

	require.ErrorContains(t, err, "not a valid Linear issue identifier")
}

func Test_writeExportDocument_surfaces_write_errors(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "LIT-1.md"), 0o750))

	_, err := writeExportDocument(dir, "LIT-1", "body")

	require.Error(t, err)
}

func Test_writeExportDocument_leaves_no_destination_on_write_failure(t *testing.T) {
	dir := t.TempDir()
	stubAtomicTempWriteError(t, errors.New("write boom"))

	_, err := writeExportDocument(dir, "LIT-1", "body")

	require.ErrorContains(t, err, "write boom")
	require.NoFileExists(t, filepath.Join(dir, "LIT-1.md"))
	entries, readErr := os.ReadDir(dir)
	require.NoError(t, readErr)
	require.Empty(t, entries)
}

func Test_renderExportDescription_falls_back_when_empty(t *testing.T) {
	require.Contains(t, renderExportDescription("   "), "_No description._")
}

func Test_renderExportComments_falls_back_when_empty(t *testing.T) {
	require.Contains(t, renderExportComments(nil), "_No comments._")
}

func Test_renderExportComment_resolves_author_and_body(t *testing.T) {
	fromUserName := renderExportComment(client.IssueCommentSummary{UserName: "omer", Body: "hi"})
	require.Contains(t, fromUserName, "### omer — ")

	unknown := renderExportComment(client.IssueCommentSummary{Body: "hi"})
	require.Contains(t, unknown, "### Unknown — ")

	emptyBody := renderExportComment(client.IssueCommentSummary{DisplayName: "Omer"})
	require.Contains(t, emptyBody, "_(empty)_")
}

func Test_renderExportAttachments_falls_back_to_url_and_empty(t *testing.T) {
	require.Contains(t, renderExportAttachments(nil), "_No attachments._")

	titled := renderExportAttachments([]client.AttachmentSummary{{URL: "https://x/y"}})
	require.Contains(t, titled, "- [https://x/y](<https://x/y>)")
}

func Test_CommandFlows_project_export_writes_document(t *testing.T) {
	dir := t.TempDir()

	stdout, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"project", "export", "project-id", dir})
	require.NoError(t, err)

	path := filepath.Join(dir, "project-id.md")
	require.Contains(t, stdout, path)
	require.Contains(t, stdout, "(1 attachments)")

	data, err := os.ReadFile(path) //nolint:gosec // G304: path is the test-controlled temp dir.
	require.NoError(t, err)
	document := string(data)
	require.Contains(t, document, "# Detail project")
	require.Contains(t, document, "- ID: project-id")
	require.Contains(t, document, "- URL: https://linear.app/kyanite/project/project-id")
	require.Contains(t, document, "## Content\n\nExisting project content")
	require.Contains(t, document, "## Attachments (1)")
	require.Contains(t, document, "- [Linked PR](<https://github.com/kyanite/linctl/pull/1>)")
}

func Test_CommandFlows_project_export_honors_output_flags(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "project-id.md")

	idOnly, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"--id-only", "project", "export", "project-id", dir})
	require.NoError(t, err)
	require.Equal(t, path+"\n", idOnly)

	jsonOut, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"--json", "project", "export", "project-id", dir})
	require.NoError(t, err)
	require.Contains(t, jsonOut, `"slug_id": "Detail project"`)
	require.Contains(t, jsonOut, `"attachments": 1`)

	quiet, _, err := runExportFlow(t, commandFlowFakeClient{}, []string{"--quiet", "project", "export", "project-id", dir})
	require.NoError(t, err)
	require.Empty(t, quiet)
}

func Test_CommandFlows_project_export_notes_truncation(t *testing.T) {
	dir := t.TempDir()

	stdout, stderr, err := runExportFlow(
		t,
		commandFlowFakeClient{truncatedExport: true},
		[]string{"--json", "project", "export", "project-id", dir},
	)
	require.NoError(t, err)
	require.Contains(t, stderr, "export capped at")
	require.Contains(t, stdout, `"truncated": true`)
}

func Test_CommandFlows_project_export_surfaces_read_errors(t *testing.T) {
	for _, operation := range []string{"project", "project_attachments"} {
		t.Run(operation, func(t *testing.T) {
			dir := t.TempDir()

			_, _, err := runExportFlow(
				t,
				commandFlowFakeClient{failOperation: operation},
				[]string{"project", "export", "project-id", dir},
			)

			require.Error(t, err)
		})
	}
}

func Test_runProjectExport_surfaces_directory_errors(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-dir")
	require.NoError(t, os.WriteFile(file, []byte("x"), 0o600))

	_, _, err := runExportFlow(
		t,
		commandFlowFakeClient{},
		[]string{"project", "export", "project-id", filepath.Join(file, "sub")},
	)

	require.Error(t, err)
}

func Test_runProjectExport_surfaces_note_writer_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{truncatedExport: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetErr(commandFailingWriter{})
	command.SetArgs([]string{"project", "export", "project-id", t.TempDir()})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
}

func Test_renderExportContent_falls_back_when_empty(t *testing.T) {
	require.Contains(t, renderExportContent("   "), "_No content._")
}

func Test_projectExportLeaf_prefers_slug_then_id(t *testing.T) {
	slug, err := projectExportLeaf(client.ProjectSummary{ID: "project-id", SlugID: "hive-spec"})
	require.NoError(t, err)
	require.Equal(t, "hive-spec", slug)

	id, err := projectExportLeaf(client.ProjectSummary{ID: "project-id", SlugID: "Detail project"})
	require.NoError(t, err)
	require.Equal(t, "project-id", id)

	_, err = projectExportLeaf(client.ProjectSummary{ID: "../escape", SlugID: "bad name"})
	require.ErrorContains(t, err, "not a valid export filename")
}

func Test_runIssueExport_surfaces_identifier_errors(t *testing.T) {
	_, _, err := runExportFlow(
		t,
		commandFlowFakeClient{invalidExportLeaf: true},
		[]string{"issue", "export", "LIT-1", t.TempDir()},
	)

	require.ErrorContains(t, err, "not a valid Linear issue identifier")
}

func Test_runProjectExport_surfaces_filename_errors(t *testing.T) {
	_, _, err := runExportFlow(
		t,
		commandFlowFakeClient{invalidExportLeaf: true},
		[]string{"project", "export", "project-id", t.TempDir()},
	)

	require.ErrorContains(t, err, "not a valid export filename")
}
