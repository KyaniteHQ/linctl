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
	identifier, err := CurrentIssueIdentifierForTeam(context.Background(), repo, "")

	// Then
	require.NoError(t, err)
	require.Equal(t, "LIT-789", identifier)
}

func Test_ParseIssueIdentifier_returns_uppercased_identifier_from_lowercase_branch_name(t *testing.T) {
	// Given
	branch := "omer/lit-123-current-command"

	// When
	identifier, ok := ParseIssueIdentifier(branch)

	// Then
	require.True(t, ok)
	require.Equal(t, "LIT-123", identifier)
}

func Test_ParseIssueIdentifier_still_matches_uppercase_identifier_in_mixed_case_branch(t *testing.T) {
	// Given
	branch := "feature/LIT-42-add-thing"

	// When
	identifier, ok := ParseIssueIdentifier(branch)

	// Then
	require.True(t, ok)
	require.Equal(t, "LIT-42", identifier)
}

func Test_ParseIssueIdentifier_prefers_uppercase_match_over_earlier_lowercase_false_positive(t *testing.T) {
	// Given
	text := "fix-1 branch also references LIT-2"

	// When
	identifier, ok := ParseIssueIdentifier(text)

	// Then
	require.True(t, ok)
	require.Equal(t, "LIT-2", identifier)
}

func Test_ParseIssueIdentifierForTeam_rejects_match_outside_pinned_team(t *testing.T) {
	// Given
	branch := "omer/fix-123-thing"

	// When
	identifier, ok := ParseIssueIdentifierForTeam(branch, "LIT")

	// Then
	require.False(t, ok)
	require.Empty(t, identifier)
}

func Test_ParseIssueIdentifierForTeam_accepts_match_within_pinned_team(t *testing.T) {
	// Given
	branch := "omer/lit-123-thing"

	// When
	identifier, ok := ParseIssueIdentifierForTeam(branch, "LIT")

	// Then
	require.True(t, ok)
	require.Equal(t, "LIT-123", identifier)
}

func Test_ParseIssueIdentifierForTeam_keeps_unfiltered_behavior_when_team_key_empty(t *testing.T) {
	// Given
	branch := "omer/lit-123-thing"

	// When
	identifier, ok := ParseIssueIdentifierForTeam(branch, "")

	// Then
	require.True(t, ok)
	require.Equal(t, "LIT-123", identifier)
}

func Test_CurrentIssueIdentifier_returns_error_when_no_issue_reference_exists(t *testing.T) {
	// Given
	repo := t.TempDir()
	runGit(t, repo, "init")
	runGit(t, repo, "checkout", "-b", "feature/no-issue")

	// When
	_, err := CurrentIssueIdentifierForTeam(context.Background(), repo, "")

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
