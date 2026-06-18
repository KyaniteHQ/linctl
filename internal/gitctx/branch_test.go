package gitctx

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseIssueIdentifier_returns_linear_identifier_from_branch_name(t *testing.T) {
	// Given
	branch := "omer/feature/LIT-123-current-command"

	// When
	identifier, ok := ParseIssueIdentifier(branch)

	// Then
	require.True(t, ok)
	require.Equal(t, "LIT-123", identifier)
}

func Test_ParseIssueIdentifier_returns_linear_identifier_from_jj_trailer(t *testing.T) {
	// Given
	description := "Implement current command\n\nLinear-issue: LIT-456\n"

	// When
	identifier, ok := ParseIssueIdentifier(description)

	// Then
	require.True(t, ok)
	require.Equal(t, "LIT-456", identifier)
}

func Test_CurrentIssueIdentifier_reads_git_branch_when_issue_named(t *testing.T) {
	// Given
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "checkout", "-b", "feature/LIT-789-current")

	// When
	identifier, err := CurrentIssueIdentifier(context.Background(), repo)

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-789", identifier)
}

func Test_CurrentIssueIdentifier_returns_error_when_no_issue_reference_exists(t *testing.T) {
	// Given
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "checkout", "-b", "feature/no-issue")

	// When
	_, err := CurrentIssueIdentifier(context.Background(), repo)

	// Then
	require.ErrorIs(t, err, ErrIssueReferenceMissing)
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	//nolint:gosec // Test helper runs fixed git commands with test-controlled arguments.
	command := exec.Command("git", args...)
	command.Dir = filepath.Clean(dir)
	output, err := command.CombinedOutput()
	require.NoError(t, err, string(output))
}
