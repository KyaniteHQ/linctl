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

// ParseIssueIdentifier extracts the first Linear issue identifier from text.
func ParseIssueIdentifier(text string) (string, bool) {
	identifier := issueIdentifierPattern.FindString(text)
	if identifier == "" {
		return "", false
	}

	return identifier, true
}

// CurrentIssueIdentifier derives the active Linear issue from git or jj checkout context.
func CurrentIssueIdentifier(ctx context.Context, dir string) (string, error) {
	branch, branchErr := currentGitBranch(ctx, dir)
	if branchErr == nil {
		identifier, ok := ParseIssueIdentifier(branch)
		if ok {
			return identifier, nil
		}
		branchErr = fmt.Errorf("%w: git branch %q", ErrIssueReferenceMissing, branch)
	}

	description, descriptionErr := currentJJDescription(ctx, dir)
	if descriptionErr == nil {
		identifier, ok := ParseIssueIdentifier(description)
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
