package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CommandFlows_cover_output_error_and_quiet_branches(t *testing.T) {
	quietCommands := [][]string{
		{"--quiet", "target"},
		{"--quiet", "whoami"},
		{"--quiet", "issue", "deps", "LIT-1"},
		{"--quiet", "issue", "pr", "LIT-1"},
		{"--quiet", "usage"},
	}
	for _, args := range quietCommands {
		t.Run("quiet "+args[len(args)-1], func(t *testing.T) {
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

	errorCommands := [][]string{
		{"--sort", "missing", "issue", "list"},
		{"--sort", "missing", "issue", "list", "--project", "project-id"},
		{"--sort", "missing", "issue", "list", "--mine"},
		{"--sort", "missing", "issue", "list", "--assignee", "assignee-id"},
		{"--sort", "missing", "issue", "list", "--label", "label-id"},
		{"--sort", "missing", "issue", "list", "--cycle", "cycle-id"},
		{"--sort", "missing", "issue", "list", "--created-after", "2026-06-01"},
		{"--sort", "missing", "issue", "list", "--created-since", "2026-06-01"},
		{"--sort", "missing", "issue", "list", "--created-before", "2026-06-30"},
		{"--sort", "missing", "issue", "list", "--has-blockers"},
		{"--sort", "missing", "issue", "list", "--blocks"},
		{"--sort", "missing", "issue", "list", "--blocked-by", "LIT-1"},
		{"--sort", "missing", "issue", "list", "--all-teams"},
		{"--sort", "missing", "issue", "comments", "LIT-1"},
		{"--sort", "missing", "issue", "search", "needle"},
		{"--sort", "missing", "project", "list"},
		{"--sort", "missing", "project", "members", "project-id"},
	}
	for _, args := range errorCommands {
		t.Run("sort error "+args[len(args)-1], func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "sort field")
		})
	}

	emptyCommands := []struct {
		name string
		args []string
		fake commandFlowFakeClient
	}{
		{name: "issue list project", args: []string{"--fail-on-empty", "issue", "list", "--project", "project-id"}, fake: commandFlowFakeClient{emptyIssueProject: true}},
		{name: "issue list mine", args: []string{"--fail-on-empty", "issue", "list", "--mine"}, fake: commandFlowFakeClient{emptyIssueMine: true}},
		{name: "issue list assignee", args: []string{"--fail-on-empty", "issue", "list", "--assignee", "assignee-id"}, fake: commandFlowFakeClient{emptyIssueMine: true}},
		{name: "issue list label", args: []string{"--fail-on-empty", "issue", "list", "--label", "label-id"}, fake: commandFlowFakeClient{emptyIssueLabel: true}},
		{name: "issue list cycle", args: []string{"--fail-on-empty", "issue", "list", "--cycle", "cycle-id"}, fake: commandFlowFakeClient{emptyIssueCycle: true}},
		{name: "issue list created-after", args: []string{"--fail-on-empty", "issue", "list", "--created-after", "2026-06-01"}, fake: commandFlowFakeClient{emptyIssueCreatedAfter: true}},
		{name: "issue list created-since", args: []string{"--fail-on-empty", "issue", "list", "--created-since", "2026-06-01"}, fake: commandFlowFakeClient{emptyIssueCreatedAfter: true}},
		{name: "issue list created-before", args: []string{"--fail-on-empty", "issue", "list", "--created-before", "2026-06-30"}, fake: commandFlowFakeClient{emptyIssueCreatedBefore: true}},
		{name: "issue list has blockers", args: []string{"--fail-on-empty", "issue", "list", "--has-blockers"}, fake: commandFlowFakeClient{emptyIssueHasBlockers: true}},
		{name: "issue list blocks", args: []string{"--fail-on-empty", "issue", "list", "--blocks"}, fake: commandFlowFakeClient{emptyIssueBlocks: true}},
		{name: "issue list blocked by", args: []string{"--fail-on-empty", "issue", "list", "--blocked-by", "LIT-1"}, fake: commandFlowFakeClient{emptyIssueBlockedBy: true}},
		{name: "issue list all teams", args: []string{"--fail-on-empty", "issue", "list", "--all-teams"}, fake: commandFlowFakeClient{emptyIssueAllTeams: true}},
		{name: "issue comments", args: []string{"--fail-on-empty", "issue", "comments", "LIT-1"}, fake: commandFlowFakeClient{emptyIssueComments: true}},
		{name: "issue search", args: []string{"--fail-on-empty", "issue", "search", "needle"}, fake: commandFlowFakeClient{emptyIssueSearch: true}},
		{name: "issue figma file key search", args: []string{"--fail-on-empty", "issue", "figma-file-key-search", "figma-key"}, fake: commandFlowFakeClient{emptyIssueFigmaSearch: true}},
		{name: "project list", args: []string{"--fail-on-empty", "project", "list"}, fake: commandFlowFakeClient{emptyProjectList: true}},
		{name: "project all", args: []string{"--fail-on-empty", "project", "all"}, fake: commandFlowFakeClient{emptyProjectList: true}},
		{name: "project members", args: []string{"--fail-on-empty", "project", "members", "project-id"}, fake: commandFlowFakeClient{emptyProjectMembers: true}},
	}
	for _, test := range emptyCommands {
		t.Run("empty "+test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "empty result")
		})
	}
}

func Test_CommandFlows_cover_issue_list_filter_validation(t *testing.T) {
	tests := [][]string{
		{"issue", "list", "--state", "started", "--project", "project-id"},
		{"issue", "list", "--state", "started", "--mine"},
		{"issue", "list", "--state", "started", "--assignee", "assignee-id"},
		{"issue", "list", "--state", "started", "--label", "label-id"},
		{"issue", "list", "--state", "started", "--cycle", "cycle-id"},
		{"issue", "list", "--state", "started", "--created-after", "2026-06-01"},
		{"issue", "list", "--state", "started", "--created-since", "2026-06-01"},
		{"issue", "list", "--created-after", "2026-06-01", "--created-since", "2026-06-01"},
		{"issue", "list", "--state", "started", "--created-before", "2026-06-30"},
		{"issue", "list", "--state", "started", "--has-blockers"},
		{"issue", "list", "--has-blockers", "--blocks"},
		{"issue", "list", "--blocks", "--blocked-by", "LIT-1"},
		{"issue", "list", "--state", "started", "--all-teams"},
	}
	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "use only one")
		})
	}
}

func Test_CommandFlows_cover_issue_current_error_branches(t *testing.T) {
	t.Run("id missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("title missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "title"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("title runtime error", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(
			t,
			[]string{"issue", "title"},
			func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			},
		)

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("url lookup error", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(
			t,
			[]string{"issue", "url"},
			func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return testCommandRuntime(commandFlowFakeClient{failOperation: "issue"}), nil
			},
		)

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})

	t.Run("branch argument lookup error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "branch", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})
}

func Test_CommandFlows_cover_issue_comment_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"issue", "comment", "LIT-1", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_issue_reply_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"issue", "reply", "LIT-1", "comment-id", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_document_create_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"document", "create", "--title", "x", "--content", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_document_update_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"document", "update", "document-id", "--content", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_comment_update_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"comment", "update", "comment-id", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_project_update_create_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"project-update", "create", "project-id", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "read body from stdin")
}

func Test_CommandFlows_cover_issue_comments_error_branches(t *testing.T) {
	t.Run("operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue_comments"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "comments", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "list issue comments LIT-1")
	})

	t.Run("writer error", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssueComments(command, []client.IssueCommentSummary{{ID: "comment-id", DisplayName: "Omer", Body: "body"}})

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_cover_comment_child_error_and_projection_branches(t *testing.T) {
	t.Run("project filter suggestion runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"project", "filter-suggestion", "started projects"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("issue vcs branch search runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "get", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("issue vcs branch search operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issueVcsBranchSearch"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "get", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue by vcs branch omer/branch")
	})

	t.Run("issue vcs branch bot actor runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "bot-actor", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("issue vcs branch bot actor operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issueVcsBranchSearch_botActor"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "bot-actor", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue vcs branch bot actor omer/branch")
	})

	t.Run("issue vcs branch shared access operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issueVcsBranchSearch_sharedAccess"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "shared-access", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue vcs branch shared access omer/branch")
	})

	t.Run("issue vcs branch shared access runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "shared-access", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("attachment issue runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "get", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("attachment issue operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "attachmentIssue"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "get", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get attachment issue attachment-id")
	})

	t.Run("attachment issue bot actor runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "bot-actor", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("attachment issue bot actor operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "attachmentIssue_botActor"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "bot-actor", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get attachment issue bot actor attachment-id")
	})

	t.Run("attachment issue shared access operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "attachmentIssue_sharedAccess"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "shared-access", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get attachment issue shared access attachment-id")
	})

	t.Run("attachment issue shared access runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "shared-access", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("issue bot actor runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "bot-actor", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("issue bot actor operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue_botActor"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "bot-actor", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue bot actor LIT-1")
	})

	t.Run("issue shared access operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue_sharedAccess"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "shared-access", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue shared access LIT-1")
	})

	t.Run("issue shared access runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "shared-access", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("comment bot actor runtime error", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		}
		defer func() {
			buildCommandRuntime = original
		}()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"comment", "bot-actor", "comment-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("bot actor operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "comment_botActor"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"comment", "bot-actor", "comment-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "get comment bot actor comment-id")
	})

	t.Run("child page projection", func(t *testing.T) {
		page := commentChildPageWithItems(
			client.CommentChildList{CommentID: "comment-id"},
			[]client.CommentMetadataSummary{{ID: "child-comment-id"}},
		)

		require.Equal(t, "child-comment-id", page.Comments[0].ID)
	})
}

func Test_CommandFlows_cover_user_settings_error_and_writer_branches(t *testing.T) {
	runtimeErrorCommands := [][]string{
		{"user", "settings", "get"},
		{"user", "settings", "notification-categories"},
		{"user", "settings", "notification-category", "assignments"},
		{"user", "settings", "notification-channels"},
		{"user", "settings", "notification-delivery"},
		{"user", "settings", "mobile-delivery"},
		{"user", "settings", "mobile-schedule"},
		{"user", "settings", "mobile-schedule-day", "monday"},
		{"user", "settings", "theme"},
		{"user", "settings", "custom-theme"},
		{"user", "settings", "custom-sidebar-theme"},
	}
	for _, args := range runtimeErrorCommands {
		t.Run("runtime "+strings.Join(args, " "), func(t *testing.T) {
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

	operationErrorCommands := []struct {
		args      []string
		operation string
		contains  string
	}{
		{args: []string{"user", "settings", "get"}, operation: "userSettings", contains: "get user settings"},
		{args: []string{"user", "settings", "notification-categories"}, operation: "userSettings_notificationCategoryPreferences", contains: "get user settings notification categories"},
		{args: []string{"user", "settings", "notification-category", "assignments"}, operation: "userSettings_notificationCategoryPreferences_assignments", contains: "get user settings category assignments"},
		{args: []string{"user", "settings", "notification-channels"}, operation: "userSettings_notificationChannelPreferences", contains: "get user settings notification channels"},
		{args: []string{"user", "settings", "notification-delivery"}, operation: "userSettings_notificationDeliveryPreferences", contains: "get user settings notification delivery"},
		{args: []string{"user", "settings", "mobile-delivery"}, operation: "userSettings_notificationDeliveryPreferences_mobile", contains: "get user settings mobile delivery"},
		{args: []string{"user", "settings", "mobile-schedule"}, operation: "userSettings_notificationDeliveryPreferences_mobile_schedule", contains: "get user settings mobile schedule"},
		{args: []string{"user", "settings", "mobile-schedule-day", "monday"}, operation: "userSettings_notificationDeliveryPreferences_mobile_schedule_monday", contains: "get user settings mobile schedule monday"},
		{args: []string{"user", "settings", "theme"}, operation: "userSettings_theme", contains: "get user settings theme"},
		{args: []string{"user", "settings", "custom-theme"}, operation: "userSettings_theme_custom", contains: "get user settings custom theme"},
		{args: []string{"user", "settings", "custom-sidebar-theme"}, operation: "userSettings_theme_custom_sidebar", contains: "get user settings custom sidebar theme"},
	}
	for _, test := range operationErrorCommands {
		t.Run("operation "+test.operation, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: test.operation})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.contains)
		})
	}

	command := &cobra.Command{}
	output := bytes.Buffer{}
	command.SetOut(&output)
	settings := client.UserSettingsSummary{ID: "settings-id", UserID: "user-id"}
	require.NoError(t, writeUserSettings(command, &rootOptions{idOnly: true}, settings))
	require.NoError(t, writeUserSettings(command, &rootOptions{quiet: true}, settings))
	require.NoError(t, writeUserSettingsValue(command, &rootOptions{quiet: true}, settings, "settings"))
	require.NoError(t, writeUserSettingsValue(command, &rootOptions{json: true}, settings, "settings"))
	require.NoError(t, writeUserSettingsNullableValue(command, &rootOptions{quiet: true}, nil, "nullable"))
	require.NoError(t, writeUserSettingsNullableValue(command, &rootOptions{json: true}, nil, "nullable"))
	require.NoError(t, writeUserSettingsNullableValue(command, &rootOptions{}, nil, "nullable"))
	require.Empty(t, pointerString(nil))
	require.Equal(t, "value", pointerString(stringPointerForUserSettingsTest("value")))

	err := runUserSettingsThemeCommand(
		context.Background(),
		command,
		&rootOptions{},
		commandRuntime{},
		"unknown",
		"desktop",
		"light",
	)
	require.Error(t, err)
}

func stringPointerForUserSettingsTest(value string) *string {
	return &value
}

type commandFailingReader struct{}

func (reader commandFailingReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}
