package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CommandFlows_report_normalization_note_write_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "create state note write failure",
			args: []string{"issue", "create", "--title", "T", "--state", "todo"},
		},
		{
			name: "create priority note write failure",
			args: []string{"issue", "create", "--title", "T", "--priority", "high"},
		},
		{
			name: "update state note write failure",
			args: []string{"issue", "update", "LIT-1", "--state", "todo"},
		},
		{
			name: "update priority note write failure",
			args: []string{"issue", "update", "LIT-1", "--priority", "high"},
		},
		{
			name: "list status note write failure",
			args: []string{"issue", "list", "--status", "todo"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetErr(commandFailingWriter{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "write failed")
		})
	}
}

func Test_CommandFlows_resolve_current_issue_from_branch(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"current"})

	require.NoError(t, err)
	require.Contains(t, output, "LIT-1 Detail issue [Todo]")
}

func Test_CommandFlows_print_current_issue_as_json(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"--json", "current"})

	require.NoError(t, err)
	require.Contains(t, output, `"identifier": "LIT-1"`)
}

func Test_CommandFlows_print_current_issue_identifier_from_issue_id(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"issue", "id"})

	require.NoError(t, err)
	require.Equal(t, "LIT-1\n", output)
}

func Test_CommandFlows_print_current_issue_title_from_issue_title(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"issue", "title"})

	require.NoError(t, err)
	require.Equal(t, "Detail issue\n", output)
}

func Test_CommandFlows_print_current_issue_url_from_issue_url(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"issue", "url"})

	require.NoError(t, err)
	require.Equal(t, "https://linear.app/kyanite/issue/LIT-1\n", output)
}

func Test_CommandFlows_print_issue_branch_from_issue_branch(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"issue", "branch", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Equal(t, "lit-1-detail-issue\n", output.String())
}

func Test_CommandFlows_print_issue_pr_from_current_branch(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"issue", "pr"})

	require.NoError(t, err)
	require.Contains(t, output, `gh pr create --title "LIT-1 Detail issue"`)
}

func Test_CommandFlows_close_current_issue_from_done(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"done"})

	require.NoError(t, err)
	require.Contains(t, output, "LIT-1 Closed issue [Done]")
}

func clientIssue(identifier string, title string) client.IssueSummary {
	return client.IssueSummary{
		ID:         identifier + "-id",
		Identifier: identifier,
		Title:      title,
		State:      "Todo",
	}
}

func Test_runCurrentIssueRead_reads_current_issue_through_the_port(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeIssuePort{
		gotIssue: clientIssue("LIT-9", "Current from port"),
	}

	err := runCurrentIssueRead(context.Background(), command, &rootOptions{}, port, "LIT-9")

	require.NoError(t, err)
	require.Equal(t, "LIT-9", port.getIssueID)
	require.Contains(t, stdout.String(), "LIT-9 Current from port")
}

func Test_resolveIssueArgumentWithReader_reads_issue_argument_through_the_port(t *testing.T) {
	port := &fakeIssuePort{
		gotIssue: clientIssue("LIT-10", "Issue current from port"),
	}

	issue, err := resolveIssueArgumentWithReader(context.Background(), port, "LIT-10")

	require.NoError(t, err)
	require.Equal(t, "LIT-10", port.getIssueID)
	require.Equal(t, "Issue current from port", issue.Title)
}

func Test_CommandFlows_report_next_errors(t *testing.T) {
	t.Run("empty candidate list", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{emptyNextIssues: true})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"next", "--dry-run"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "next issue not found")
	})

	t.Run("empty candidate list with fail on empty", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{emptyNextIssues: true})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--fail-on-empty", "next", "--dry-run"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "empty result")
	})
}

func Test_CommandFlows_rank_next_issue_candidates(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{rankedNextIssues: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"next", "--dry-run"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "LIT-30 Unblocks checkout [Todo]")
}

func swapCheckoutBranch(fn func(context.Context, string) error) func() {
	original := checkoutBranch
	checkoutBranch = fn

	return func() { checkoutBranch = original }
}

func Test_CommandFlows_next_starts_picked_issue(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"next"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "LIT-1")
}

func Test_CommandFlows_next_checkout_creates_branch_then_starts(t *testing.T) {
	called := false
	restoreCheckout := swapCheckoutBranch(func(_ context.Context, _ string) error {
		called = true

		return nil
	})
	defer restoreCheckout()
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"next", "--checkout"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.True(t, called)
	require.Contains(t, output.String(), "LIT-1")
}

func Test_CommandFlows_next_checkout_failure_aborts(t *testing.T) {
	restoreCheckout := swapCheckoutBranch(func(_ context.Context, _ string) error {
		return errors.New("checkout boom")
	})
	defer restoreCheckout()
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"next", "--checkout"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "checkout boom")
}

func Test_runNextWithPicker_reads_and_starts_through_the_port(t *testing.T) {
	command, stdout, _ := bufferedCommand()
	port := &fakeIssuePort{
		resolved: client.ResolvedTarget{Team: client.TargetTeam{ID: "team-id"}},
		nextList: client.IssueList{Issues: []client.IssueSummary{
			{
				Identifier: "LIT-12",
				Title:      "Picked from port",
				BranchName: "lit-12-picked-from-port",
				State:      "Todo",
			},
		}},
		started: clientIssue("LIT-12", "Started from port"),
	}

	err := runNextWithPicker(
		context.Background(),
		command,
		&rootOptions{},
		port,
		nextFlags{limit: 7},
	)

	require.NoError(t, err)
	require.Equal(t, 1, port.nextCalls)
	require.Equal(t, "team-id", port.nextTeamID)
	require.Equal(t, 7, port.nextLimit)
	require.Equal(t, "LIT-12", port.startID)
	require.Contains(t, stdout.String(), "LIT-12 Started from port")
}

func Test_CommandFlows_next_surfaces_start_failure(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "IssueUpdate"})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"next"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
}

func Test_runGitCheckoutBranch_creates_and_fails(t *testing.T) {
	t.Run("creates a branch in a repo", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		t.Chdir(dir)

		err := runGitCheckoutBranch(context.Background(), "linctl-it-next")

		require.NoError(t, err)
	})

	t.Run("fails outside a repo", func(t *testing.T) {
		t.Chdir(t.TempDir())

		err := runGitCheckoutBranch(context.Background(), "linctl-it-next")

		require.Error(t, err)
		require.Contains(t, err.Error(), "git checkout -b")
	})
}

func Test_CommandFlows_report_current_issue_errors(t *testing.T) {
	t.Run("missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"current"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("done missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"done"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("runtime failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"current"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("done runtime failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"done"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("issue lookup failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"current"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return testCommandRuntime(commandFlowFakeClient{failOperation: "issue"}), nil
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})

	t.Run("done close failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"done"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return testCommandRuntime(commandFlowFakeClient{failOperation: "IssueClose"}), nil
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "close issue LIT-1")
	})
}

func runCurrentCommandInGitBranch(t *testing.T, args []string) (string, error) {
	t.Helper()

	return runCurrentCommandInGitBranchWithRuntime(t, args, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return testCommandRuntime(commandFlowFakeClient{}), nil
	})
}

// runCurrentCommandInGitBranchWithRuntime swaps the package-level
// buildCommandRuntime and changes the process working directory via t.Chdir.
// Both are process-wide side effects, so tests using this helper (and the
// others in this file that swap buildCommandRuntime) must NOT call t.Parallel()
// — concurrent execution would race on the shared builder and cwd. The
// structural fix is to thread the runtime builder through as an argument; until
// then this sequential constraint is load-bearing.
func runCurrentCommandInGitBranchWithRuntime(
	t *testing.T,
	args []string,
	runtimeBuilder func(context.Context, *rootOptions) (commandRuntime, error),
) (string, error) {
	t.Helper()

	dir := t.TempDir()
	runGitCommand(t, dir, "init")
	runGitCommand(t, dir, "checkout", "-b", "feature/LIT-1-current")
	t.Chdir(dir)

	output := bytes.Buffer{}
	original := buildCommandRuntime
	buildCommandRuntime = runtimeBuilder
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs(args)

	err := command.ExecuteContext(context.Background())
	return output.String(), err
}
