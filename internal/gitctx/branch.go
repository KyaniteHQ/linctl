// Package gitctx derives Linear context from the current VCS checkout.
package gitctx

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ErrIssueReferenceMissing marks a checkout without a Linear issue reference.
var ErrIssueReferenceMissing = errors.New("linear issue reference missing")

var issueIdentifierPattern = regexp.MustCompile(`\b[A-Z][A-Z0-9]+-[0-9]+\b`)

var issueIdentifierPatternAnyCase = regexp.MustCompile(`(?i)\b[a-z][a-z0-9]+-[0-9]+\b`)

// ParseIssueIdentifier extracts the first Linear issue identifier from text.
// An explicit uppercase identifier anywhere in the text always wins; failing
// that, the first case-insensitive match is normalized to uppercase. This
// tolerates Linear's own generated branch names (the `branchName` field and
// the "copy branch name" UI action), which are lowercase.
func ParseIssueIdentifier(text string) (string, bool) {
	identifier := issueIdentifierPattern.FindString(text)
	if identifier != "" {
		return identifier, true
	}

	identifier = issueIdentifierPatternAnyCase.FindString(text)
	if identifier == "" {
		return "", false
	}

	return strings.ToUpper(identifier), true
}

// ParseIssueIdentifierForTeam extracts the first Linear issue identifier from
// text, restricted to the given team key. An empty team key keeps the
// unfiltered behavior of ParseIssueIdentifier, so reads without a Pinned
// Target are unaffected. A non-empty team key rejects any match whose
// team-key portion does not equal teamKey, bounding the false-positive risk
// that case-insensitive matching introduces (e.g. "fix-123", "bug-42").
func ParseIssueIdentifierForTeam(text string, teamKey string) (string, bool) {
	if teamKey == "" {
		return ParseIssueIdentifier(text)
	}

	wantTeamKey := strings.ToUpper(teamKey)
	for _, match := range issueIdentifierPatternAnyCase.FindAllString(text, -1) {
		identifier := strings.ToUpper(match)
		matchTeamKey, _, found := strings.Cut(identifier, "-")
		if found && matchTeamKey == wantTeamKey {
			return identifier, true
		}
	}

	return "", false
}

// CurrentIssueIdentifierForTeam derives the active Linear issue from git or jj
// checkout context, restricted to the given team key. An empty team key keeps
// the unfiltered behavior, so callers without a Pinned Target are unaffected.
func CurrentIssueIdentifierForTeam(ctx context.Context, dir string, teamKey string) (string, error) {
	branch, branchErr := currentGitBranch(ctx, dir)
	if branchErr == nil {
		identifier, ok := ParseIssueIdentifierForTeam(branch, teamKey)
		if ok {
			return identifier, nil
		}
		branchErr = fmt.Errorf("%w: git branch %q", ErrIssueReferenceMissing, branch)
	}

	description, descriptionErr := currentJJDescription(ctx, dir)
	if descriptionErr == nil {
		identifier, ok := ParseIssueIdentifierForTeam(description, teamKey)
		if ok {
			return identifier, nil
		}
		descriptionErr = fmt.Errorf("%w: jj description has no identifier", ErrIssueReferenceMissing)
	}

	return "", fmt.Errorf(
		"%w: git branch: %w; jj description: %w",
		ErrIssueReferenceMissing,
		branchErr,
		descriptionErr,
	)
}

func currentGitBranch(ctx context.Context, dir string) (string, error) {
	command := exec.CommandContext(ctx, "git", "branch", "--show-current")
	command.Dir = filepath.Clean(dir)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git branch --show-current: %w: %s", err, strings.TrimSpace(string(output)))
	}
	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "", fmt.Errorf("%w: git branch empty", ErrIssueReferenceMissing)
	}

	return branch, nil
}

func currentJJDescription(ctx context.Context, dir string) (string, error) {
	command := exec.CommandContext(ctx, "jj", "log", "-r", "@", "--no-graph", "-T", "description")
	command.Dir = filepath.Clean(dir)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("jj log -r @ --no-graph -T description: %w: %s", err, strings.TrimSpace(string(output)))
	}
	description := strings.TrimSpace(string(output))
	if description == "" {
		return "", fmt.Errorf("%w: jj description empty", ErrIssueReferenceMissing)
	}

	return description, nil
}
