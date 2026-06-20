package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_CommandFlows_execute_read_and_write_commands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
		fake     commandFlowFakeClient
	}{
		{name: "target", args: []string{"target"}, contains: "org org-id team LIT/team-id project project-id confirmed true"},
		{name: "doctor", args: []string{"doctor"}, contains: "config ok\n token set\n target confirmed LIT/team-id project project-id"},
		{name: "whoami", args: []string{"whoami"}, contains: "Omer <omer@example.com>"},
		{name: "next dry run", args: []string{"next", "--dry-run"}, contains: "LIT-27 Next issue [Todo]"},
		{name: "issue list", args: []string{"issue", "list", "--limit", "1"}, contains: "LIT-1 Listed issue [Todo]"},
		{name: "issue list state filter", args: []string{"issue", "list", "--state", "started", "--limit", "1"}, contains: "LIT-2 Started issue [Started]", fake: commandFlowFakeClient{expectedStateType: "started"}},
		{name: "issue list project filter", args: []string{"issue", "list", "--project", "project-id", "--limit", "1"}, contains: "LIT-4 Project issue [Todo]", fake: commandFlowFakeClient{expectedProjectID: "project-id"}},
		{name: "issue list mine filter", args: []string{"issue", "list", "--mine", "--limit", "1"}, contains: "LIT-5 Mine issue [Todo]", fake: commandFlowFakeClient{expectedAssigneeID: "user-id"}},
		{name: "issue list assignee filter", args: []string{"issue", "list", "--assignee", "assignee-id", "--limit", "1"}, contains: "LIT-6 Assigned issue [Todo]", fake: commandFlowFakeClient{expectedAssigneeID: "assignee-id"}},
		{name: "issue list label filter", args: []string{"issue", "list", "--label", "label-id", "--limit", "1"}, contains: "LIT-7 Labeled issue [Todo]", fake: commandFlowFakeClient{expectedLabelID: "label-id"}},
		{name: "issue list cycle filter", args: []string{"issue", "list", "--cycle", "cycle-id", "--limit", "1"}, contains: "LIT-8 Cycle issue [Todo]", fake: commandFlowFakeClient{expectedCycleID: "cycle-id"}},
		{name: "issue list created-after filter", args: []string{"issue", "list", "--created-after", "2026-06-01", "--limit", "1"}, contains: "LIT-9 Recent issue [Todo]", fake: commandFlowFakeClient{expectedCreatedAfter: "2026-06-01"}},
		{name: "issue list created-since filter", args: []string{"issue", "list", "--created-since", "2026-06-01", "--limit", "1"}, contains: "LIT-9 Recent issue [Todo]", fake: commandFlowFakeClient{expectedCreatedAfter: "2026-06-01"}},
		{name: "issue list created-before filter", args: []string{"issue", "list", "--created-before", "2026-06-30", "--limit", "1"}, contains: "LIT-19 Older issue [Todo]", fake: commandFlowFakeClient{expectedCreatedBefore: "2026-06-30"}},
		{name: "issue list has blockers filter", args: []string{"issue", "list", "--has-blockers", "--limit", "1"}, contains: "LIT-21 Blocked issue [Todo]"},
		{name: "issue list blocks filter", args: []string{"issue", "list", "--blocks", "--limit", "1"}, contains: "LIT-22 Blocking issue [Todo]"},
		{name: "issue list blocked by filter", args: []string{"issue", "list", "--blocked-by", "LIT-1", "--limit", "1"}, contains: "LIT-23 Blocked by issue [Todo]", fake: commandFlowFakeClient{expectedBlockedBy: "LIT-1"}},
		{name: "issue list all teams", args: []string{"issue", "list", "--all-teams", "--limit", "1"}, contains: "LIT-20 All-team issue [Todo]"},
		{name: "issue search", args: []string{"issue", "search", "needle", "--limit", "1"}, contains: "LIT-3 Search result [Todo]", fake: commandFlowFakeClient{expectedSearchQuery: "needle"}},
		{name: "issue get", args: []string{"issue", "get", "LIT-1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "issue deps", args: []string{"issue", "deps", "LIT-1", "--limit", "2"}, contains: "blocked_by:\nLIT-24 Blocker issue [Todo]", fake: commandFlowFakeClient{expectedIssueDeps: "LIT-1"}},
		{name: "issue pr", args: []string{"issue", "pr", "LIT-1"}, contains: `gh pr create --title "LIT-1 Detail issue" --body "https://linear.app/kyanite/issue/LIT-1"`},
		{name: "issue create", args: []string{"issue", "create", "--title", "Created issue"}, contains: "LIT-2 Created issue [Todo]"},
		{name: "issue update", args: []string{"issue", "update", "LIT-1", "--title", "Updated issue"}, contains: "LIT-1 Updated issue [Todo]"},
		{name: "issue update append", args: []string{"issue", "update", "LIT-1", "--append", "Progress note"}, contains: "LIT-1 Updated issue [Todo]", fake: commandFlowFakeClient{expectedUpdateDescription: "Existing description\n\nProgress note"}},
		{name: "issue start", args: []string{"issue", "start", "LIT-1"}, contains: "LIT-1 Started issue [Started]", fake: commandFlowFakeClient{expectedStartAssigneeID: "user-id", expectedStartStateID: "started-state"}},
		{name: "issue comment", args: []string{"issue", "comment", "LIT-1", "--body", "Looks good"}, contains: "comment comment-id on LIT-1"},
		{name: "issue reply", args: []string{"issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"}, contains: "comment comment-id on LIT-1", fake: commandFlowFakeClient{expectedCommentBody: "Reply body", expectedCommentParentID: "comment-id"}},
		{name: "issue comments", args: []string{"issue", "comments", "LIT-1", "--limit", "1"}, contains: "comment-id Omer First comment"},
		{name: "comment list", args: []string{"comment", "list", "--limit", "1"}, contains: "comment-id Omer First comment"},
		{name: "comment get", args: []string{"comment", "get", "comment-id"}, contains: "comment-id Omer First comment"},
		{name: "issue close", args: []string{"issue", "close", "LIT-1"}, contains: "LIT-1 Closed issue [Done]"},
		{name: "project list", args: []string{"project", "list", "--limit", "1"}, contains: "project-id Listed project [Backlog]"},
		{name: "project get", args: []string{"project", "get", "project-id"}, contains: "project-id Detail project [Backlog]"},
		{name: "project members", args: []string{"project", "members", "project-id", "--limit", "1"}, contains: "user-id Omer"},
		{name: "project updates", args: []string{"project", "updates", "project-id", "--limit", "1"}, contains: "project-update-id onTrack Omer First update"},
		{name: "project update list", args: []string{"project-update", "list", "--limit", "1"}, contains: "project-update-id onTrack Omer First update"},
		{name: "project update get", args: []string{"project-update", "get", "project-update-id"}, contains: "project-update-id onTrack Omer First update"},
		{name: "project milestone list", args: []string{"project-milestone", "list", "project-id", "--limit", "1"}, contains: "project-milestone-id Launch milestone [next]"},
		{name: "project milestone create", args: []string{"project-milestone", "create", "project-id", "--name", "Created milestone"}, contains: "project-milestone-id Created milestone [next]"},
		{name: "project milestone update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"}, contains: "project-milestone-id Updated milestone [done]"},
		{name: "project create", args: []string{"project", "create", "--name", "Created project"}, contains: "project-id Created project [Backlog]"},
		{name: "project update", args: []string{"project", "update", "project-id", "--name", "Updated project"}, contains: "project-id Updated project [Started]"},
		{name: "project archive", args: []string{"project", "archive", "project-id"}, contains: "project-id Archived project [Canceled]"},
		{name: "document list", args: []string{"document", "list", "--limit", "1"}, contains: "document-id Spec [project]"},
		{name: "document get", args: []string{"document", "get", "document-id"}, contains: "document-id Team note [team]"},
		{name: "label list", args: []string{"label", "list", "--limit", "1"}, contains: "label-id Bug #ff0000"},
		{name: "label get", args: []string{"label", "get", "label-id"}, contains: "label-id Bug #ff0000"},
		{name: "team list", args: []string{"team", "list", "--limit", "1"}, contains: "team-id LIT linctl"},
		{name: "team get", args: []string{"team", "get", "team-id"}, contains: "team-id LIT linctl"},
		{name: "team members", args: []string{"team", "members", "team-id", "--limit", "1"}, contains: "user-id Omer <omer@example.com>"},
		{name: "user list", args: []string{"user", "list", "--limit", "1"}, contains: "user-id Omer <omer@example.com>"},
		{name: "user get", args: []string{"user", "get", "user-id"}, contains: "user-id Omer <omer@example.com>"},
		{name: "user me", args: []string{"user", "me"}, contains: "user-id Omer <omer@example.com>"},
		{name: "workflow state list", args: []string{"workflow-state", "list", "--limit", "1"}, contains: "workflow-state-id Started [started]"},
		{name: "workflow state get", args: []string{"workflow-state", "get", "workflow-state-id"}, contains: "workflow-state-id Started [started]"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), test.contains)
		})
	}
}

func Test_CommandFlows_read_issue_comment_body_from_stdin(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{expectedCommentBody: "stdin body\nsecond line"})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(strings.NewReader("stdin body\nsecond line"))
	command.SetOut(&output)
	command.SetArgs([]string{"issue", "comment", "LIT-1", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "comment comment-id on LIT-1")
}

func Test_CommandFlows_read_issue_text_from_files(t *testing.T) {
	descriptionFile := writeTempTextFile(t, "description from file")
	appendFile := writeTempTextFile(t, "append from file")
	commentFile := writeTempTextFile(t, "comment from file")
	replyFile := writeTempTextFile(t, "reply from file")

	tests := []struct {
		name string
		args []string
		fake commandFlowFakeClient
	}{
		{
			name: "create description",
			args: []string{"issue", "create", "--title", "Created issue", "--description-file", descriptionFile},
			fake: commandFlowFakeClient{expectedCreateDescription: "description from file"},
		},
		{
			name: "update append",
			args: []string{"issue", "update", "LIT-1", "--append-file", appendFile},
			fake: commandFlowFakeClient{expectedUpdateDescription: "Existing description\n\nappend from file"},
		},
		{
			name: "comment body",
			args: []string{"issue", "comment", "LIT-1", "--body-file", commentFile},
			fake: commandFlowFakeClient{expectedCommentBody: "comment from file"},
		},
		{
			name: "reply body",
			args: []string{"issue", "reply", "LIT-1", "comment-id", "--body-file", replyFile},
			fake: commandFlowFakeClient{expectedCommentBody: "reply from file", expectedCommentParentID: "comment-id"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.NotEmpty(t, output.String())
		})
	}
}

func Test_CommandFlows_report_issue_text_file_errors(t *testing.T) {
	textFile := writeTempTextFile(t, "from file")
	missingFile := filepath.Join(t.TempDir(), "missing.md")
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "create description conflict",
			args:     []string{"issue", "create", "--title", "Created issue", "--description", "inline", "--description-file", textFile},
			contains: "description and description-file are mutually exclusive",
		},
		{
			name:     "update description conflict",
			args:     []string{"issue", "update", "LIT-1", "--description", "inline", "--description-file", textFile},
			contains: "description and description-file are mutually exclusive",
		},
		{
			name:     "update append conflict",
			args:     []string{"issue", "update", "LIT-1", "--append", "inline", "--append-file", textFile},
			contains: "append and append-file are mutually exclusive",
		},
		{
			name:     "comment body conflict",
			args:     []string{"issue", "comment", "LIT-1", "--body", "inline", "--body-file", textFile},
			contains: "body and body-file are mutually exclusive",
		},
		{
			name:     "reply body conflict",
			args:     []string{"issue", "reply", "LIT-1", "comment-id", "--body", "inline", "--body-file", textFile},
			contains: "body and body-file are mutually exclusive",
		},
		{
			name:     "missing file",
			args:     []string{"issue", "comment", "LIT-1", "--body-file", missingFile},
			contains: "read body from file",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.contains)
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

func Test_CommandFlows_report_next_errors(t *testing.T) {
	t.Run("requires dry run", func(t *testing.T) {
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"next"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "next requires --dry-run")
	})

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
			return testCommandRuntime(commandFlowFakeClient{failOperation: "IssueByID"}), nil
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

func Test_CommandFlows_report_runtime_and_writer_errors(t *testing.T) {
	t.Run("runtime error returns from command", func(t *testing.T) {
		commands := [][]string{
			{"target"},
			{"doctor"},
			{"whoami"},
			{"next", "--dry-run"},
			{"issue", "list"},
			{"issue", "search", "needle"},
			{"issue", "get", "LIT-1"},
			{"issue", "deps", "LIT-1"},
			{"issue", "pr", "LIT-1"},
			{"issue", "create", "--title", "Created issue"},
			{"issue", "update", "LIT-1", "--title", "Updated issue"},
			{"issue", "start", "LIT-1"},
			{"issue", "comment", "LIT-1", "--body", "Looks good"},
			{"issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"},
			{"issue", "comments", "LIT-1"},
			{"issue", "close", "LIT-1"},
			{"project", "list"},
			{"project", "get", "project-id"},
			{"project", "members", "project-id"},
			{"project", "updates", "project-id"},
			{"project-milestone", "list", "project-id"},
			{"project-milestone", "get", "project-milestone-id"},
			{"project-milestone", "create", "project-id", "--name", "Created milestone"},
			{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"},
			{"project", "create", "--name", "Created project"},
			{"project", "update", "project-id", "--name", "Updated project"},
			{"project", "archive", "project-id"},
			{"document", "list"},
			{"document", "get", "document-id"},
			{"label", "list"},
			{"label", "get", "label-id"},
			{"team", "list"},
			{"team", "get", "team-id"},
			{"team", "members", "team-id"},
			{"user", "list"},
			{"user", "get", "user-id"},
			{"user", "me"},
		}
		for _, args := range commands {
			t.Run(strings.Join(args, " "), func(t *testing.T) {
				original := buildCommandRuntime
				buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, errors.New("runtime failed")
				}
				defer func() {
					buildCommandRuntime = original
				}()
				command := NewRootCommand(context.Background(), BuildInfo{})
				command.SetArgs(args)

				err := command.ExecuteContext(context.Background())

				require.Error(t, err)
				require.Contains(t, err.Error(), "runtime failed")
			})
		}
	})

	t.Run("writeIssues returns writer errors", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssues(command, &rootOptions{}, []client.IssueSummary{{Identifier: "LIT-1", Title: "Broken", State: "Todo"}})

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("doctor returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"doctor"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "list"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project members returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "members", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project updates returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "updates", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project milestone list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project-milestone", "list", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("document list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"document", "list"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("usage returns writer errors", func(t *testing.T) {
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"usage"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_print_json_for_read_and_comment_commands(t *testing.T) {
	tests := [][]string{
		{"--json", "target"},
		{"--json", "doctor"},
		{"--json", "whoami"},
		{"--json", "next", "--dry-run"},
		{"--json", "issue", "list", "--limit", "1"},
		{"--json", "issue", "search", "needle", "--limit", "1"},
		{"--json", "issue", "deps", "LIT-1", "--limit", "2"},
		{"--json", "issue", "pr", "LIT-1"},
		{"--json", "issue", "start", "LIT-1"},
		{"--json", "issue", "comment", "LIT-1", "--body", "Looks good"},
		{"--json", "issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"},
		{"--json", "--fields", "id,display_name", "issue", "comments", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,display_name", "comment", "list", "--limit", "1"},
		{"--json", "comment", "get", "comment-id"},
		{"--json", "project", "list", "--limit", "1"},
		{"--json", "project", "members", "project-id", "--limit", "1"},
		{"--json", "--fields", "id,health,display_name", "project", "updates", "project-id", "--limit", "1"},
		{"--json", "--fields", "id,health,project_id", "project-update", "list", "--limit", "1"},
		{"--json", "project-update", "get", "project-update-id"},
		{"--json", "--fields", "id,name,status", "project-milestone", "list", "project-id", "--limit", "1"},
		{"--json", "--fields", "id,title,parent_type", "document", "list", "--limit", "1"},
		{"--json", "--fields", "id,name,color", "label", "list", "--limit", "1"},
		{"--json", "--fields", "id,key,name", "team", "list", "--limit", "1"},
		{"--json", "--fields", "id,display_name,email", "team", "members", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,display_name,email", "user", "list", "--limit", "1"},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), "{")
		})
	}
}

func Test_CommandFlows_print_compact_json_when_compact_flag_is_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "--compact", "issue", "get", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `{"id":"issue-id"`)
	require.NotContains(t, output.String(), "\n  ")
}

func Test_CommandFlows_project_json_fields_when_fields_flag_is_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "--fields", "identifier,title,state", "issue", "get", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"identifier": "LIT-1"`)
	require.Contains(t, output.String(), `"title": "Detail issue"`)
	require.Contains(t, output.String(), `"state": "Todo"`)
	require.NotContains(t, output.String(), `"url"`)
	require.NotContains(t, output.String(), `"project_id"`)
}

func Test_CommandFlows_print_only_id_when_id_only_flag_is_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--id-only", "issue", "get", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Equal(t, "issue-id\n", output.String())
}

func Test_CommandFlows_suppress_success_output_when_quiet_flag_is_set(t *testing.T) {
	tests := [][]string{
		{"--quiet", "doctor"},
		{"--quiet", "issue", "get", "LIT-1"},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Empty(t, output.String())
		})
	}
}

func Test_CommandFlows_fail_on_empty_list_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyIssueList: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "issue", "list"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_fail_on_empty_project_updates_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyProjectUpdates: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "project", "updates", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_allow_empty_project_updates_without_fail_on_empty(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyProjectUpdates: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"project", "updates", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Empty(t, output.String())
}

func Test_CommandFlows_report_project_updates_sort_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--sort", "missing", "project", "updates", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), `sort field "missing" is not present`)
}

func Test_CommandFlows_fail_on_empty_project_milestones_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyProjectMilestones: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "project-milestone", "list", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_allow_empty_project_milestones_without_fail_on_empty(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyProjectMilestones: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"project-milestone", "list", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Empty(t, output.String())
}

func Test_CommandFlows_report_project_milestone_sort_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--sort", "missing", "project-milestone", "list", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), `sort field "missing" is not present`)
}

func Test_CommandFlows_get_project_milestone(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"project-milestone", "get", "project-milestone-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "project-milestone-id Launch milestone [next]")
}

func Test_CommandFlows_get_project_milestone_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "project-milestone", "get", "project-milestone-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"id": "project-milestone-id"`)
	require.Contains(t, output.String(), `"status": "next"`)
}

func Test_CommandFlows_report_project_milestone_get_runtime_error(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"project-milestone", "get", "project-milestone-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime failed")
}

func Test_CommandFlows_report_project_milestone_write_runtime_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "create", args: []string{"project-milestone", "create", "project-id", "--name", "Created milestone"}},
		{name: "update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := buildCommandRuntime
			buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			}
			defer func() {
				buildCommandRuntime = original
			}()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "runtime failed")
		})
	}
}

func Test_CommandFlows_report_project_milestone_get_writer_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"project-milestone", "get", "project-milestone-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func Test_CommandFlows_report_project_milestone_write_writer_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "create", args: []string{"project-milestone", "create", "project-id", "--name", "Created milestone"}},
		{name: "update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(commandFailingWriter{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "write failed")
		})
	}
}

func Test_CommandFlows_sort_issue_list_when_sort_flags_are_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{multiIssueList: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--sort", "title", "--order", "desc", "issue", "list"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Less(t, strings.Index(output.String(), "Zebra issue"), strings.Index(output.String(), "Alpha issue"))
}

func Test_CommandFlows_print_minimal_human_output_when_format_flag_is_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--format", "minimal", "issue", "get", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Equal(t, "LIT-1\n", output.String())
}

func Test_CommandFlows_print_workflow_state_list_as_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "workflow-state", "list", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"workflow_states": [`)
	require.Contains(t, output.String(), `"team_key": "LIT"`)
}

func Test_CommandFlows_report_operation_errors(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		operation string
		contains  string
	}{
		{name: "target resolve", args: []string{"target"}, operation: "Teams", contains: "resolve teams"},
		{name: "doctor target resolve", args: []string{"doctor"}, operation: "Teams", contains: "resolve teams"},
		{name: "whoami resolve", args: []string{"whoami"}, operation: "Viewer", contains: "resolve viewer"},
		{name: "next target resolve", args: []string{"next", "--dry-run"}, operation: "Teams", contains: "resolve teams"},
		{name: "next issues", args: []string{"next", "--dry-run"}, operation: "NextIssuesByTeam", contains: "list next issues"},
		{name: "issue list target resolve", args: []string{"issue", "list"}, operation: "Teams", contains: "resolve teams"},
		{name: "issue list", args: []string{"issue", "list"}, operation: "IssuesByTeam", contains: "list issues"},
		{name: "issue list project filter", args: []string{"issue", "list", "--project", "project-id"}, operation: "IssuesByTeamProject", contains: "list issues"},
		{name: "issue list mine filter", args: []string{"issue", "list", "--mine"}, operation: "IssuesByTeamAssignee", contains: "list issues"},
		{name: "issue list assignee filter", args: []string{"issue", "list", "--assignee", "assignee-id"}, operation: "IssuesByTeamAssignee", contains: "list issues"},
		{name: "issue list label filter", args: []string{"issue", "list", "--label", "label-id"}, operation: "IssuesByTeamLabel", contains: "list issues"},
		{name: "issue list cycle filter", args: []string{"issue", "list", "--cycle", "cycle-id"}, operation: "IssuesByTeamCycle", contains: "list issues"},
		{name: "issue list created-after filter", args: []string{"issue", "list", "--created-after", "2026-06-01"}, operation: "IssuesByTeamCreatedAfter", contains: "list issues"},
		{name: "issue list created-since filter", args: []string{"issue", "list", "--created-since", "2026-06-01"}, operation: "IssuesByTeamCreatedAfter", contains: "list issues"},
		{name: "issue list created-before filter", args: []string{"issue", "list", "--created-before", "2026-06-30"}, operation: "IssuesByTeamCreatedBefore", contains: "list issues"},
		{name: "issue list has blockers filter", args: []string{"issue", "list", "--has-blockers"}, operation: "IssuesByTeamHasBlockers", contains: "list issues"},
		{name: "issue list blocks filter", args: []string{"issue", "list", "--blocks"}, operation: "IssuesByTeamBlocks", contains: "list issues"},
		{name: "issue list blocked by filter", args: []string{"issue", "list", "--blocked-by", "LIT-1"}, operation: "IssueBlockedIssues", contains: "list issues"},
		{name: "issue list all teams", args: []string{"issue", "list", "--all-teams"}, operation: "AllTeamIssues", contains: "list issues"},
		{name: "issue search target resolve", args: []string{"issue", "search", "needle"}, operation: "Teams", contains: "resolve teams"},
		{name: "issue search", args: []string{"issue", "search", "needle"}, operation: "IssueSearch", contains: "search issues"},
		{name: "issue get", args: []string{"issue", "get", "LIT-1"}, operation: "IssueByID", contains: "get issue LIT-1"},
		{name: "issue deps", args: []string{"issue", "deps", "LIT-1"}, operation: "IssueDependencies", contains: "get issue dependencies LIT-1"},
		{name: "issue pr", args: []string{"issue", "pr", "LIT-1"}, operation: "IssueByID", contains: "get issue LIT-1"},
		{name: "issue create", args: []string{"issue", "create", "--title", "Created issue"}, operation: "IssueCreate", contains: "create issue"},
		{name: "issue update", args: []string{"issue", "update", "LIT-1", "--title", "Updated issue"}, operation: "IssueUpdate", contains: "update issue LIT-1"},
		{name: "issue start state", args: []string{"issue", "start", "LIT-1"}, operation: "StartedWorkflowStates", contains: "list started workflow states"},
		{name: "issue start update", args: []string{"issue", "start", "LIT-1"}, operation: "IssueUpdate", contains: "start issue LIT-1"},
		{name: "issue comment", args: []string{"issue", "comment", "LIT-1", "--body", "Looks good"}, operation: "IssueCommentCreate", contains: "comment on issue LIT-1"},
		{name: "issue reply", args: []string{"issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"}, operation: "IssueCommentCreate", contains: "comment on issue LIT-1"},
		{name: "comment list", args: []string{"comment", "list"}, operation: "comments", contains: "list comments"},
		{name: "comment get", args: []string{"comment", "get", "comment-id"}, operation: "comment", contains: "get comment comment-id"},
		{name: "issue close", args: []string{"issue", "close", "LIT-1"}, operation: "IssueClose", contains: "close issue LIT-1"},
		{name: "project list target resolve", args: []string{"project", "list"}, operation: "Teams", contains: "resolve teams"},
		{name: "project list", args: []string{"project", "list"}, operation: "Projects", contains: "list projects"},
		{name: "project get", args: []string{"project", "get", "project-id"}, operation: "ProjectByID", contains: "get project project-id"},
		{name: "project members", args: []string{"project", "members", "project-id"}, operation: "ProjectMembers", contains: "list project members project-id"},
		{name: "project updates", args: []string{"project", "updates", "project-id"}, operation: "ProjectUpdates", contains: "list project updates project-id"},
		{name: "project update list", args: []string{"project-update", "list"}, operation: "projectUpdates", contains: "list project updates"},
		{name: "project update get", args: []string{"project-update", "get", "project-update-id"}, operation: "projectUpdate", contains: "get project update project-update-id"},
		{name: "project milestone list", args: []string{"project-milestone", "list", "project-id"}, operation: "ProjectMilestones", contains: "list project milestones project-id"},
		{name: "project milestone get", args: []string{"project-milestone", "get", "project-milestone-id"}, operation: "ProjectMilestoneByID", contains: "get project milestone project-milestone-id"},
		{name: "project milestone create", args: []string{"project-milestone", "create", "project-id", "--name", "Created milestone"}, operation: "ProjectMilestoneCreate", contains: "create project milestone"},
		{name: "project milestone update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"}, operation: "ProjectMilestoneUpdate", contains: "update project milestone project-milestone-id"},
		{name: "project create", args: []string{"project", "create", "--name", "Created project"}, operation: "ProjectCreate", contains: "create project"},
		{name: "project update", args: []string{"project", "update", "project-id", "--name", "Updated project"}, operation: "ProjectUpdate", contains: "update project project-id"},
		{name: "project archive", args: []string{"project", "archive", "project-id"}, operation: "ProjectArchive", contains: "archive project project-id"},
		{name: "document list", args: []string{"document", "list"}, operation: "Documents", contains: "list documents"},
		{name: "document get", args: []string{"document", "get", "document-id"}, operation: "document", contains: "get document document-id"},
		{name: "label list", args: []string{"label", "list"}, operation: "IssueLabels", contains: "list labels"},
		{name: "label get", args: []string{"label", "get", "label-id"}, operation: "IssueLabelByID", contains: "get label label-id"},
		{name: "team list", args: []string{"team", "list"}, operation: "Teams", contains: "list teams"},
		{name: "team get", args: []string{"team", "get", "team-id"}, operation: "TeamByID", contains: "get team team-id"},
		{name: "team members", args: []string{"team", "members", "team-id"}, operation: "TeamMembers", contains: "list team members team-id"},
		{name: "user list", args: []string{"user", "list"}, operation: "Users", contains: "list users"},
		{name: "user get", args: []string{"user", "get", "user-id"}, operation: "UserByID", contains: "get user user-id"},
		{name: "user me", args: []string{"user", "me"}, operation: "ViewerUser", contains: "get viewer user"},
		{name: "workflow state list", args: []string{"workflow-state", "list"}, operation: "workflowStates", contains: "list workflow states"},
		{name: "workflow state get", args: []string{"workflow-state", "get", "workflow-state-id"}, operation: "workflowState", contains: "get workflow state workflow-state-id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: test.operation})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.contains)
		})
	}
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()

	command := exec.Command("git", args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	require.NoError(t, err, string(output))
}

func writeTempTextFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "body.md")
	err := os.WriteFile(path, []byte(content), 0o600)
	require.NoError(t, err)

	return path
}

type commandFailingWriter struct{}

func (writer commandFailingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failed")
}

func useCommandRuntime(t *testing.T, graphqlClient graphql.Client) func() {
	t.Helper()

	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return testCommandRuntime(graphqlClient), nil
	}

	return func() {
		buildCommandRuntime = original
	}
}

func testCommandRuntime(graphqlClient graphql.Client) commandRuntime {
	return commandRuntime{
		config: config.Resolved{
			Token: "test-token",
			Target: config.Target{
				OrgID:     "org-id",
				TeamKey:   "LIT",
				TeamID:    "team-id",
				ProjectID: "project-id",
			},
		},
		graphqlClient: graphqlClient,
	}
}

type commandFlowFakeClient struct {
	emptyIssueList            bool
	emptyIssueComments        bool
	emptyIssueProject         bool
	emptyIssueMine            bool
	emptyIssueLabel           bool
	emptyIssueCycle           bool
	emptyIssueCreatedAfter    bool
	emptyIssueCreatedBefore   bool
	emptyIssueHasBlockers     bool
	emptyIssueBlocks          bool
	emptyIssueBlockedBy       bool
	emptyIssueAllTeams        bool
	emptyIssueSearch          bool
	emptyNextIssues           bool
	rankedNextIssues          bool
	expectedStateType         string
	expectedProjectID         string
	expectedAssigneeID        string
	expectedLabelID           string
	expectedCycleID           string
	expectedCreatedAfter      string
	expectedCreatedBefore     string
	expectedBlockedBy         string
	expectedIssueDeps         string
	expectedSearchQuery       string
	emptyProjectList          bool
	emptyProjectMembers       bool
	emptyProjectUpdates       bool
	emptyProjectMilestones    bool
	expectedCommentBody       string
	expectedCommentParentID   string
	expectedCreateDescription string
	expectedUpdateDescription string
	expectedStartAssigneeID   string
	expectedStartStateID      string
	failOperation             string
	multiIssueList            bool
}

func (client commandFlowFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if request.OpName == client.failOperation {
		return errors.New("operation failed")
	}
	if err := client.requireExpectedVariables(request); err != nil {
		return err
	}

	payload, err := commandFlowPayload(request.OpName, client)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

func (client commandFlowFakeClient) requireExpectedVariables(request *graphql.Request) error {
	if client.expectedCreateDescription != "" && request.OpName == "IssueCreate" {
		return requireRequestVariable(
			request,
			[]string{"input", "description"},
			client.expectedCreateDescription,
			"create description",
		)
	}
	if client.expectedCommentBody != "" && request.OpName == "IssueCommentCreate" {
		return requireRequestVariable(request, []string{"input", "body"}, client.expectedCommentBody, "comment body")
	}
	if client.expectedCommentParentID != "" && request.OpName == "IssueCommentCreate" {
		return requireRequestVariable(
			request,
			[]string{"input", "parentId"},
			client.expectedCommentParentID,
			"comment parent id",
		)
	}
	if client.expectedUpdateDescription != "" && request.OpName == "IssueUpdate" {
		return requireRequestVariable(
			request,
			[]string{"input", "description"},
			client.expectedUpdateDescription,
			"update description",
		)
	}
	if err := client.requireExpectedIssueListVariables(request); err != nil {
		return err
	}
	if client.expectedSearchQuery != "" && request.OpName == "IssueSearch" {
		return requireRequestVariable(request, []string{"query"}, client.expectedSearchQuery, "search query")
	}
	if client.expectedIssueDeps != "" && request.OpName == "IssueDependencies" {
		return requireRequestVariable(request, []string{"id"}, client.expectedIssueDeps, "issue deps id")
	}
	return client.requireExpectedIssueStartVariables(request)
}

func (client commandFlowFakeClient) requireExpectedIssueStartVariables(request *graphql.Request) error {
	if client.expectedStartAssigneeID != "" && request.OpName == "IssueUpdate" {
		if err := requireRequestVariable(
			request,
			[]string{"input", "assigneeId"},
			client.expectedStartAssigneeID,
			"start assignee id",
		); err != nil {
			return err
		}
	}
	if client.expectedStartStateID != "" && request.OpName == "IssueUpdate" {
		return requireRequestVariable(request, []string{"input", "stateId"}, client.expectedStartStateID, "start state id")
	}

	return nil
}

func (client commandFlowFakeClient) requireExpectedIssueListVariables(request *graphql.Request) error {
	if client.expectedStateType != "" && request.OpName == "IssuesByTeamState" {
		return requireRequestVariable(request, []string{"stateType"}, client.expectedStateType, "state type")
	}
	if client.expectedProjectID != "" && request.OpName == "IssuesByTeamProject" {
		return requireRequestVariable(request, []string{"projectID"}, client.expectedProjectID, "project id")
	}
	if client.expectedAssigneeID != "" && request.OpName == "IssuesByTeamAssignee" {
		return requireRequestVariable(request, []string{"assigneeID"}, client.expectedAssigneeID, "assignee id")
	}
	if client.expectedLabelID != "" && request.OpName == "IssuesByTeamLabel" {
		return requireRequestVariable(request, []string{"labelID"}, client.expectedLabelID, "label id")
	}
	if client.expectedCycleID != "" && request.OpName == "IssuesByTeamCycle" {
		return requireRequestVariable(request, []string{"cycleID"}, client.expectedCycleID, "cycle id")
	}

	return client.requireExpectedDependencyIssueListVariables(request)
}

func (client commandFlowFakeClient) requireExpectedDependencyIssueListVariables(request *graphql.Request) error {
	if client.expectedCreatedAfter != "" && request.OpName == "IssuesByTeamCreatedAfter" {
		return requireRequestVariable(request, []string{"createdAfter"}, client.expectedCreatedAfter, "created after")
	}
	if client.expectedCreatedBefore != "" && request.OpName == "IssuesByTeamCreatedBefore" {
		return requireRequestVariable(request, []string{"createdBefore"}, client.expectedCreatedBefore, "created before")
	}
	if client.expectedBlockedBy != "" && request.OpName == "IssueBlockedIssues" {
		return requireRequestVariable(request, []string{"id"}, client.expectedBlockedBy, "blocked by issue")
	}

	return nil
}

func requireRequestVariable(request *graphql.Request, keys []string, expected string, label string) error {
	actual, err := requestVariable[string](request, keys...)
	if err != nil {
		return err
	}
	if actual != expected {
		return errors.New(label + " = " + actual)
	}

	return nil
}

func commandFlowPayload(operation string, fake commandFlowFakeClient) (string, error) {
	switch operation {
	case "Viewer":
		return `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`, nil
	case "Teams":
		return `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, nil
	case "TargetProject":
		return `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`, nil
	}
	if payload, ok := commandFlowIssuePayload(operation, fake); ok {
		return payload, nil
	}
	if payload, ok := commandFlowProjectPayload(operation, fake); ok {
		return payload, nil
	}
	if payload, ok := commandFlowPeopleAndReferencePayload(operation); ok {
		return payload, nil
	}

	return "", errors.New("missing fake response for " + operation)
}

func commandFlowPeopleAndReferencePayload(operation string) (string, bool) {
	switch operation {
	case "Documents":
		return `{"documents":{"nodes":[` + commandDocumentJSON(
			"Spec",
			`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
		) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "document":
		return `{"document":` + commandDocumentJSON(
			"Team note",
			`"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null`,
		) + `}`, true
	case "IssueLabels":
		return `{"issueLabels":{"nodes":[` + commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssueLabelByID":
		return `{"issueLabel":` + commandLabelJSON("") + `}`, true
	case "TeamByID":
		return `{"team":` + commandTeamJSON(true) + `}`, true
	case "TeamMembers":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","members":{"nodes":[` + commandUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "Users":
		return `{"users":{"nodes":[` + commandUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "UserByID":
		return `{"user":` + commandUserJSON() + `}`, true
	case "ViewerUser":
		return `{"viewer":` + commandUserJSON() + `}`, true
	case "workflowStates":
		return `{"workflowStates":{"nodes":[` + commandWorkflowStateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "workflowState":
		return `{"workflowState":` + commandWorkflowStateJSON() + `}`, true
	case "comments":
		return `{"comments":{"nodes":[` + commandTopLevelCommentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "comment":
		return `{"comment":` + commandTopLevelCommentJSON() + `}`, true
	default:
		return "", false
	}
}

func commandFlowIssuePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowIssueReadPayload(operation, fake); ok {
		return payload, true
	}

	return commandFlowIssueWritePayload(operation, fake)
}

func commandFlowIssueReadPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowIssueListPayload(operation, fake); ok {
		return payload, true
	}

	switch operation {
	case "IssueSearch":
		if fake.emptyIssueSearch {
			return `{"issues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-3", "Search result", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "NextIssuesByTeam":
		if fake.emptyNextIssues {
			return emptyCommandIssuesPayload(), true
		}
		if fake.rankedNextIssues {
			return `{"issues":{"nodes":[` +
				commandIssueWithNextRankJSON("LIT-28", "Low priority standalone", 4, "Low", "2026-05-01T12:00:00Z", 0) + `,` +
				commandIssueWithNextRankJSON("LIT-29", "Urgent standalone", 1, "Urgent", "2026-06-01T12:00:00Z", 0) + `,` +
				commandIssueWithNextRankJSON("LIT-30", "Unblocks checkout", 2, "High", "2026-06-10T12:00:00Z", 2) +
				`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"issues":{"nodes":[` + commandIssueWithNextRankJSON("LIT-27", "Next issue", 0, "No priority", "2026-06-01T12:00:00Z", 0) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssueByID":
		return `{"issue":` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `}`, true
	case "IssueDependencies":
		return commandFlowIssueDependenciesPayload(), true
	case "IssueComments":
		if fake.emptyIssueComments {
			return `{"issue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[{"id":"comment-id","body":"First comment","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","parentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueDependenciesPayload() string {
	return `{"issue":{
		"id":"issue-id",
		"identifier":"LIT-1",
		"parent":` + commandIssueJSON("LIT-25", "Parent issue", "todo-state", "Todo", "unstarted") + `,
		"children":{
			"nodes":[` + commandIssueJSON("LIT-26", "Child issue", "todo-state", "Todo", "unstarted") + `],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		},
		"relations":{
			"nodes":[{"id":"blocks-relation","type":"blocks","relatedIssue":` + commandIssueJSON("LIT-23", "Blocked issue", "todo-state", "Todo", "unstarted") + `}],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		},
		"inverseRelations":{
			"nodes":[{"id":"blocked-by-relation","type":"blocks","issue":` + commandIssueJSON("LIT-24", "Blocker issue", "todo-state", "Todo", "unstarted") + `}],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		}
	}}`
}

func commandFlowIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowBroadIssueListPayload(operation, fake); ok {
		return payload, true
	}
	if payload, ok := commandFlowDependencyIssueListPayload(operation, fake); ok {
		return payload, true
	}

	switch operation {
	case "IssuesByTeamState":
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-2", "Started issue", "started-state", "Started", fake.expectedStateType) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamProject":
		if fake.emptyIssueProject {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-4", "Project issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamAssignee":
		return commandFlowAssigneeIssueListPayload(fake), true
	case "IssuesByTeamLabel":
		if fake.emptyIssueLabel {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-7", "Labeled issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCycle":
		if fake.emptyIssueCycle {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-8", "Cycle issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCreatedAfter":
		if fake.emptyIssueCreatedAfter {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-9", "Recent issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCreatedBefore":
		if fake.emptyIssueCreatedBefore {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-19", "Older issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	default:
		return "", false
	}
}

func commandFlowDependencyIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "IssuesByTeamHasBlockers":
		if fake.emptyIssueHasBlockers {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-21", "Blocked issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamBlocks":
		if fake.emptyIssueBlocks {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-22", "Blocking issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssueBlockedIssues":
		if fake.emptyIssueBlockedBy {
			return `{"issue":{"id":"issue-id","identifier":"LIT-1","relations":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","relations":{"nodes":[{"id":"relation-id","type":"blocks","relatedIssue":` + commandIssueJSON("LIT-23", "Blocked by issue", "todo-state", "Todo", "unstarted") + `}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowBroadIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "AllTeamIssues":
		if fake.emptyIssueAllTeams {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-20", "All-team issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeam":
		return commandFlowUnfilteredIssueListPayload(fake), true
	default:
		return "", false
	}
}

func commandFlowUnfilteredIssueListPayload(fake commandFlowFakeClient) string {
	if fake.emptyIssueList {
		return emptyCommandIssuesPayload()
	}
	if fake.multiIssueList {
		return `{"issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Alpha issue", "todo-state", "Todo", "unstarted") + `,` +
			commandIssueJSON("LIT-2", "Zebra issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
	}

	return `{"issues":{"nodes":[` + commandIssueJSON("LIT-1", "Listed issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func commandFlowAssigneeIssueListPayload(fake commandFlowFakeClient) string {
	if fake.emptyIssueMine {
		return emptyCommandIssuesPayload()
	}
	if fake.expectedAssigneeID == "assignee-id" {
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-6", "Assigned issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
	}

	return `{"issues":{"nodes":[` + commandIssueJSON("LIT-5", "Mine issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func emptyCommandIssuesPayload() string {
	return `{"issues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func commandFlowIssueWritePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "IssueCreate":
		return `{"issueCreate":{"success":true,"issue":` + commandIssueJSON("LIT-2", "Created issue", "todo-state", "Todo", "unstarted") + `}}`, true
	case "IssueUpdate":
		if fake.expectedStartStateID != "" {
			return `{"issueUpdate":{"success":true,"issue":` +
				commandIssueJSON("LIT-1", "Started issue", "started-state", "Started", "started") + `}}`, true
		}
		return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSON("LIT-1", "Updated issue", "todo-state", "Todo", "unstarted") + `}}`, true
	case "IssueCommentCreate":
		return `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"Looks good","url":"https://linear.app/comment/comment-id","issue":` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `}}}`, true
	case "CompletedWorkflowStates":
		return `{"workflowStates":{"nodes":[{"id":"done-state","name":"Done","type":"completed","position":1}]}}`, true
	case "StartedWorkflowStates":
		return `{"workflowStates":{"nodes":[{"id":"started-state","name":"Started","type":"started","position":1}]}}`, true
	case "IssueClose":
		return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSON("LIT-1", "Closed issue", "done-state", "Done", "completed") + `}}`, true
	default:
		return "", false
	}
}

func commandFlowProjectPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowProjectReadPayload(operation, fake); ok {
		return payload, true
	}

	return commandFlowProjectWritePayload(operation)
}

func commandFlowProjectReadPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "Projects":
		if fake.emptyProjectList {
			return `{"team":{"projects":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"team":{"projects":{"nodes":[` + commandProjectJSON("Listed project", "Backlog", "backlog") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "ProjectByID":
		return `{"project":` + commandProjectJSON("Detail project", "Backlog", "backlog") + `}`, true
	case "ProjectMembers":
		if fake.emptyProjectMembers {
			return `{"project":{"id":"project-id","name":"Detail project","members":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com"}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "ProjectUpdates":
		if fake.emptyProjectUpdates {
			return `{"project":{"id":"project-id","name":"Detail project","projectUpdates":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","projectUpdates":{"nodes":[{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projectUpdates":
		if fake.emptyProjectUpdates {
			return `{"projectUpdates":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"projectUpdates":{"nodes":[` + commandProjectUpdateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectUpdate":
		return `{"projectUpdate":` + commandProjectUpdateJSON() + `}`, true
	case "ProjectMilestones":
		if fake.emptyProjectMilestones {
			return `{"project":{"id":"project-id","name":"Detail project","projectMilestones":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "ProjectMilestoneByID":
		return `{"projectMilestone":` + commandProjectMilestoneJSON("Launch milestone", "next") + `}`, true
	default:
		return "", false
	}
}

func commandFlowProjectWritePayload(operation string) (string, bool) {
	switch operation {
	case "ProjectMilestoneCreate":
		return `{"projectMilestoneCreate":{"success":true,"projectMilestone":` + commandProjectMilestoneJSON("Created milestone", "next") + `}}`, true
	case "ProjectMilestoneUpdate":
		return `{"projectMilestoneUpdate":{"success":true,"projectMilestone":` + commandProjectMilestoneJSON("Updated milestone", "done") + `}}`, true
	case "ProjectCreate":
		return `{"projectCreate":{"success":true,"project":` + commandProjectJSON("Created project", "Backlog", "backlog") + `}}`, true
	case "ProjectUpdate":
		return `{"projectUpdate":{"success":true,"project":` + commandProjectJSON("Updated project", "Started", "started") + `}}`, true
	case "ProjectArchive":
		return `{"projectArchive":{"success":true,"entity":` + commandProjectJSON("Archived project", "Canceled", "canceled") + `}}`, true
	default:
		return "", false
	}
}

func commandIssueJSON(identifier string, title string, stateID string, state string, stateType string) string {
	return `{
		"id":"issue-id",
		"description":"Existing description",
		"identifier":"` + identifier + `",
		"title":"` + title + `",
		"branchName":"` + strings.ToLower(identifier) + `-` + strings.ToLower(strings.ReplaceAll(title, " ", "-")) + `",
		"url":"https://linear.app/kyanite/issue/` + identifier + `",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"state":{"id":"` + stateID + `","name":"` + state + `","type":"` + stateType + `"},
		"assignee":null,
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

func commandIssueWithNextRankJSON(
	identifier string,
	title string,
	priority int,
	priorityLabel string,
	createdAt string,
	unblocksCount int,
) string {
	return strings.TrimSuffix(commandIssueJSON(identifier, title, "todo-state", "Todo", "unstarted"), "\n\t}") +
		`,
		"priority":` + strconv.Itoa(priority) + `,
		"priorityLabel":"` + priorityLabel + `",
		"createdAt":"` + createdAt + `",
		"relations":{"nodes":[` + commandBlockingRelationsJSON(unblocksCount) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}
	}`
}

func commandBlockingRelationsJSON(count int) string {
	relations := make([]string, 0, count)
	for i := range count {
		relations = append(relations, fmt.Sprintf(`{"type":"blocks","relatedIssue":{"id":"blocked-%d","state":{"type":"unstarted"}}}`, i))
	}

	return strings.Join(relations, ",")
}

func commandProjectJSON(name string, status string, statusType string) string {
	return `{
		"id":"project-id",
		"name":"` + name + `",
		"description":"description",
		"slugId":"` + name + `",
		"url":"https://linear.app/kyanite/project/project-id",
		"priority":0,
		"status":{"id":"status-id","name":"` + status + `","type":"` + statusType + `"},
		"lead":null,
		"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl"}]}
	}`
}

func commandProjectUpdateJSON() string {
	return `{
		"id":"project-update-id",
		"body":"First update",
		"health":"onTrack",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"url":"https://linear.app/project-update/project-update-id",
		"project":{"id":"project-id","name":"Pinned project"},
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandProjectMilestoneJSON(name string, status string) string {
	return `{
		"id":"project-milestone-id",
		"name":"` + name + `",
		"description":"milestone body",
		"targetDate":"2026-06-30",
		"status":"` + status + `",
		"progress":0.5,
		"sortOrder":1,
		"project":` + commandProjectJSON("Pinned project", "Backlog", "backlog") + `
	}`
}

func commandDocumentJSON(title string, parents string) string {
	return `{
		"id":"document-id",
		"title":"` + title + `",
		"slugId":"document-slug",
		"archivedAt":null,
		` + parents + `
	}`
}

func commandLabelJSON(description string) string {
	descriptionPayload := "null"
	if description != "" {
		descriptionPayload = `"` + description + `"`
	}

	return `{
		"id":"label-id",
		"name":"Bug",
		"description":` + descriptionPayload + `,
		"color":"#ff0000",
		"isGroup":false,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandTeamJSON(includeDescription bool) string {
	descriptionPayload := "null"
	if includeDescription {
		descriptionPayload = `"team body"`
	}

	return `{
		"id":"team-id",
		"key":"LIT",
		"name":"linctl",
		"description":` + descriptionPayload + `,
		"archivedAt":null,
		"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}
	}`
}

func commandUserJSON() string {
	return `{
		"id":"user-id",
		"name":"omer",
		"displayName":"Omer",
		"email":"omer@example.com",
		"active":true,
		"guest":false,
		"admin":true
	}`
}

func commandWorkflowStateJSON() string {
	return `{
		"id":"workflow-state-id",
		"name":"Started",
		"type":"started",
		"color":"#f2c94c",
		"position":2,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandTopLevelCommentJSON() string {
	return `{
		"id":"comment-id",
		"body":"First comment",
		"url":"https://linear.app/comment/comment-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"editedAt":null,
		"resolvedAt":null,
		"parentId":null,
		"issueId":"issue-id",
		"projectId":null,
		"projectUpdateId":null,
		"initiativeId":null,
		"initiativeUpdateId":null,
		"documentContentId":null,
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

var _ graphql.Client = commandFlowFakeClient{}

func requestVariable[T comparable](request *graphql.Request, keys ...string) (T, error) {
	var zero T
	payload, err := json.Marshal(request.Variables)
	if err != nil {
		return zero, err
	}
	var variables map[string]any
	if err := json.Unmarshal(payload, &variables); err != nil {
		return zero, err
	}
	current := any(variables)
	for _, key := range keys {
		object, ok := current.(map[string]any)
		if !ok {
			return zero, errors.New("request variable is not an object")
		}
		value, ok := object[key]
		if !ok {
			return zero, errors.New("request variable missing " + key)
		}
		current = value
	}
	value, ok := current.(T)
	if !ok {
		return zero, errors.New("request variable has unexpected type")
	}

	return value, nil
}
