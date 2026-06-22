package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func runIssueCreateCommand(t *testing.T, args []string) (string, error) {
	t.Helper()
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs(args)

	err := command.ExecuteContext(context.Background())

	return output.String(), err
}

func Test_IssueCreate_dry_run_renders_title_and_description(t *testing.T) {
	output, err := runIssueCreateCommand(t, []string{
		"issue", "create", "--dry-run", "--title", "Draft issue", "--description", "Body text",
	})

	require.NoError(t, err)
	require.Contains(t, output, "Draft issue")
	require.Contains(t, output, "Body text")
}

func Test_IssueCreate_dry_run_uses_template_defaults(t *testing.T) {
	output, err := runIssueCreateCommand(t, []string{
		"issue", "create", "--dry-run", "--template", "template-id",
	})

	require.NoError(t, err)
	require.Contains(t, output, "Template title")
	require.Contains(t, output, "Reproduce here")
}

func Test_IssueCreate_dry_run_json_emits_draft(t *testing.T) {
	output, err := runIssueCreateCommand(t, []string{
		"--json", "issue", "create", "--dry-run", "--title", "Draft issue", "--description", "Body",
	})

	require.NoError(t, err)
	require.Contains(t, output, `"title"`)
	require.Contains(t, output, "Draft issue")
}

func Test_IssueCreate_dry_run_quiet_is_silent(t *testing.T) {
	output, err := runIssueCreateCommand(t, []string{
		"--quiet", "issue", "create", "--dry-run", "--title", "Draft issue",
	})

	require.NoError(t, err)
	require.Empty(t, output)
}

func Test_IssueCreate_section_replaces_matching_heading(t *testing.T) {
	output, err := runIssueCreateCommand(t, []string{
		"issue", "create", "--dry-run",
		"--description", "## Steps\n\nold steps\n\n## Notes\n\nkeep me",
		"--section", "Steps=fresh steps",
	})

	require.NoError(t, err)
	require.Contains(t, output, "fresh steps")
	require.NotContains(t, output, "old steps")
	require.Contains(t, output, "keep me")
}

func Test_IssueCreate_section_replaces_trailing_heading(t *testing.T) {
	output, err := runIssueCreateCommand(t, []string{
		"issue", "create", "--dry-run",
		"--description", "## Intro\n\nhi\n\n## Steps\n\nold",
		"--section", "Steps=last value",
	})

	require.NoError(t, err)
	require.Contains(t, output, "last value")
	require.NotContains(t, output, "old")
}

func Test_IssueCreate_section_appends_when_heading_absent(t *testing.T) {
	output, err := runIssueCreateCommand(t, []string{
		"issue", "create", "--dry-run", "--description", "Intro paragraph", "--section", "Extra=more detail",
	})

	require.NoError(t, err)
	require.Contains(t, output, "## Extra")
	require.Contains(t, output, "more detail")
}

func Test_IssueCreate_section_appends_to_empty_description(t *testing.T) {
	output, err := runIssueCreateCommand(t, []string{
		"issue", "create", "--dry-run", "--section", "Context=only section",
	})

	require.NoError(t, err)
	require.Contains(t, output, "## Context")
	require.Contains(t, output, "only section")
}

func Test_IssueCreate_section_rejects_missing_equals(t *testing.T) {
	_, err := runIssueCreateCommand(t, []string{
		"issue", "create", "--dry-run", "--section", "noequals",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid --section")
}

func Test_markdownHeadingTitle_classifies_lines(t *testing.T) {
	title, ok := markdownHeadingTitle("## Steps")
	require.True(t, ok)
	require.Equal(t, "Steps", title)

	title, ok = markdownHeadingTitle("  #   Spaced heading  ")
	require.True(t, ok)
	require.Equal(t, "Spaced heading", title)

	_, ok = markdownHeadingTitle("plain text")
	require.False(t, ok)

	_, ok = markdownHeadingTitle("###")
	require.False(t, ok)

	_, ok = markdownHeadingTitle("##nospace")
	require.False(t, ok)
}
