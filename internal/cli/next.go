package cli

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// checkoutBranch creates and switches to a git branch. It is a package var so
// tests can substitute the git invocation.
var checkoutBranch = runGitCheckoutBranch

func runGitCheckoutBranch(ctx context.Context, branchName string) error {
	// branchName is a Linear-provided issue branch name passed as a discrete argv
	// argument (no shell), so there is no command-injection surface here.
	//nolint:gosec // G204: trusted Linear branch name passed as an explicit argv arg, not a shell string.
	command := exec.CommandContext(ctx, "git", "checkout", "-b", branchName)
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git checkout -b %s: %w: %s", branchName, err, strings.TrimSpace(string(output)))
	}

	return nil
}

func addNextCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	dryRun := false
	checkout := false
	limit := 20
	command := &cobra.Command{
		Use:   "next",
		Short: "Pick the next unblocked issue and start it",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runNext(ctx, command, options, nextFlags{dryRun: dryRun, checkout: checkout, limit: limit})
		},
	}
	command.Flags().BoolVar(&dryRun, "dry-run", dryRun, "preview the pick without starting it or creating a branch")
	command.Flags().BoolVar(&checkout, "checkout", checkout, "git checkout -b the issue branch before starting it")
	command.Flags().IntVar(&limit, "limit", limit, "maximum candidate issues to inspect")
	root.AddCommand(command)
}

// nextFlags collects the inputs of the next command.
type nextFlags struct {
	dryRun   bool
	checkout bool
	limit    int
}

func runNext(ctx context.Context, command *cobra.Command, options *rootOptions, flags nextFlags) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	return runNextWithPicker(ctx, command, options, nextIssueAdapterFor(runtime), flags)
}

type nextIssueClientAdapter struct {
	issueClientAdapter
	runtime commandRuntime
}

func nextIssueAdapterFor(runtime commandRuntime) nextIssueClientAdapter {
	return nextIssueClientAdapter{issueClientAdapter: issueAdapterFor(runtime), runtime: runtime}
}

func (adapter nextIssueClientAdapter) ResolveTarget(ctx context.Context) (client.ResolvedTarget, error) {
	return adapter.runtime.resolveTarget(ctx)
}

func runNextWithPicker(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	picker nextIssuePicker,
	flags nextFlags,
) error {
	target, err := picker.ResolveTarget(ctx)
	if err != nil {
		return err
	}
	issues, err := picker.ListNextIssuesByTeam(ctx, target.Team.ID, flags.limit)
	if err != nil {
		return err
	}
	if err := ensureNonEmpty(options, len(issues.Issues)); err != nil {
		return err
	}
	if len(issues.Issues) == 0 {
		return errors.New("next issue not found")
	}
	picked := issues.Issues[0]
	if flags.dryRun {
		return writeIssue(command, options, picked)
	}
	if flags.checkout {
		if err := checkoutBranch(ctx, picked.BranchName); err != nil {
			return err
		}
	}
	started, err := picker.StartIssue(ctx, picked.Identifier)
	if err != nil {
		return err
	}

	return writeIssue(command, options, started)
}
