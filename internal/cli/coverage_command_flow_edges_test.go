package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
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

			require.ErrorContains(t, err, "sort field")
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
		{name: "issue list updated-after", args: []string{"--fail-on-empty", "issue", "list", "--updated-after", "2026-07-01"}, fake: commandFlowFakeClient{emptyIssueUpdatedAfter: true}},
		{name: "issue list updated-before", args: []string{"--fail-on-empty", "issue", "list", "--updated-before", "2026-07-30"}, fake: commandFlowFakeClient{emptyIssueUpdatedBefore: true}},
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

			require.ErrorContains(t, err, "empty result")
		})
	}
}

func Test_CommandFlows_cover_issue_list_filter_validation(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{
			[]string{"issue", "list", "--state", "started", "--all-teams"},
			"--all-teams cannot combine with team-scoped filters",
		},
		{
			[]string{"issue", "list", "--mine", "--assignee", "assignee-id"},
			"use only one of --mine or --assignee",
		},
		{
			[]string{"issue", "list", "--created-after", "2026-06-01", "--created-since", "2026-06-01"},
			"use only one of --created-after or --created-since",
		},
		{
			[]string{"issue", "list", "--blocks", "--blocked-by", "LIT-1"},
			"--blocked-by resolves issue relations and cannot combine with other filters",
		},
	}
	for _, testCase := range tests {
		t.Run(strings.Join(testCase.args, " "), func(t *testing.T) {
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(testCase.args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, testCase.want)
		})
	}
}

func Test_IssueList_accepts_combined_team_scoped_filters(t *testing.T) {
	restore := useCommandRuntime(t, graphqlPayloadOverride{
		inner: commandFlowFakeClient{},
		payloads: map[string]string{
			"IssuesByTeamFiltered": `{"issues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		},
	})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&bytes.Buffer{})
	command.SetArgs([]string{
		"issue", "list",
		"--state", "started",
		"--project", "project-id",
		"--assignee", "user-id",
		"--label", "label-id",
		"--cycle", "cycle-id",
		"--created-after", "2026-06-01",
		"--created-before", "2026-06-30",
		"--updated-after", "2026-07-01",
		"--updated-before", "2026-07-30",
		"--has-blockers",
		"--blocks",
	})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
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

		require.ErrorContains(t, err, "linear issue reference missing")
	})

	t.Run("title missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "title"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "linear issue reference missing")
	})

	t.Run("title runtime error", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(
			t,
			[]string{"issue", "title"},
			func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			},
		)

		require.ErrorContains(t, err, "runtime failed")
	})

	t.Run("url lookup error", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(
			t,
			[]string{"issue", "url"},
			func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return testCommandRuntime(commandFlowFakeClient{failOperation: "issue"}), nil
			},
		)

		require.ErrorContains(t, err, "get issue LIT-1")
	})

	t.Run("branch argument lookup error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "branch", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get issue LIT-1")
	})
}

func Test_CommandFlows_cover_issue_comment_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"issue", "comment", "LIT-1", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "read body from stdin")
}

func Test_CommandFlows_cover_issue_reply_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"issue", "reply", "LIT-1", "comment-id", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "read body from stdin")
}

func Test_CommandFlows_cover_document_create_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"document", "create", "--title", "x", "--content", "-"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "read body from stdin")
}

func Test_CommandFlows_cover_document_update_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"document", "update", "document-id", "--content", "-"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "read body from stdin")
}

func Test_CommandFlows_cover_comment_update_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"comment", "update", "comment-id", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "read body from stdin")
}

func Test_CommandFlows_cover_project_update_create_stdin_read_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(commandFailingReader{})
	command.SetArgs([]string{"project-update", "create", "project-id", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "read body from stdin")
}

func Test_CommandFlows_cover_issue_comments_operation_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue_comments"})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "comments", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.ErrorContains(t, err, "list issue comments LIT-1")
}

func Test_CommandFlows_cover_comment_child_error_branches(t *testing.T) {
	t.Run("issue vcs branch search operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issueVcsBranchSearch"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "get", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get issue by vcs branch omer/branch")
	})

	t.Run("issue vcs branch bot actor operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issueVcsBranchSearch_botActor"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "bot-actor", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get issue vcs branch bot actor omer/branch")
	})

	t.Run("issue vcs branch shared access operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issueVcsBranchSearch_sharedAccess"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "vcs-branch-search", "shared-access", "omer/branch"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get issue vcs branch shared access omer/branch")
	})

	t.Run("attachment issue operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "attachmentIssue"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "get", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get attachment issue attachment-id")
	})

	t.Run("attachment issue bot actor operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "attachmentIssue_botActor"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "bot-actor", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get attachment issue bot actor attachment-id")
	})

	t.Run("attachment issue shared access operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "attachmentIssue_sharedAccess"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"attachment", "issue", "shared-access", "attachment-id"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get attachment issue shared access attachment-id")
	})

	t.Run("issue bot actor operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue_botActor"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "bot-actor", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get issue bot actor LIT-1")
	})

	t.Run("issue shared access operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "issue_sharedAccess"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"issue", "shared-access", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get issue shared access LIT-1")
	})

	t.Run("bot actor operation error", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "comment_botActor"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"comment", "bot-actor", "comment-id"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "get comment bot actor comment-id")
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

			require.ErrorContains(t, err, "runtime failed")
		})
	}

	operationErrorCommands := []struct {
		args      []string
		operation string
		contains  string
	}{
		{args: []string{"user", "settings", "get"}, operation: "userSettings", contains: "get user settings"},
		{args: []string{"user", "settings", "notification-categories"}, operation: "userSettings_notificationCategoryPreferences", contains: "get user settings notification categories"},
		{args: []string{"user", "settings", "notification-category", "assignments"}, operation: "userSettings_notificationCategoryPreferences", contains: "get user settings category assignments"},
		{args: []string{"user", "settings", "notification-channels"}, operation: "userSettings_notificationChannelPreferences", contains: "get user settings notification channels"},
		{args: []string{"user", "settings", "notification-delivery"}, operation: "userSettings_notificationDeliveryPreferences", contains: "get user settings notification delivery"},
		{args: []string{"user", "settings", "mobile-delivery"}, operation: "userSettings_notificationDeliveryPreferences_mobile", contains: "get user settings mobile delivery"},
		{args: []string{"user", "settings", "mobile-schedule"}, operation: "userSettings_notificationDeliveryPreferences_mobile_schedule", contains: "get user settings mobile schedule"},
		{args: []string{"user", "settings", "mobile-schedule-day", "monday"}, operation: "userSettings_notificationDeliveryPreferences_mobile_schedule", contains: "get user settings mobile schedule monday"},
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

			require.ErrorContains(t, err, test.contains)
		})
	}
}

func Test_UserSettings_honors_output_modes(t *testing.T) {
	nullMobile := payloadGraphQLClient{
		payload: `{"userSettings":{"notificationDeliveryPreferences":{"mobile":null}}}`,
	}
	tests := []struct {
		name string
		args []string
		fake graphql.Client
		want string
	}{
		{
			name: "settings id-only",
			args: []string{"--id-only", "user", "settings", "get"},
			fake: commandFlowFakeClient{},
			want: "settings-id\n",
		},
		{
			name: "settings quiet",
			args: []string{"--quiet", "user", "settings", "get"},
			fake: commandFlowFakeClient{},
			want: "",
		},
		{
			name: "schedule day quiet",
			args: []string{"--quiet", "user", "settings", "mobile-schedule-day", "monday"},
			fake: commandFlowFakeClient{},
			want: "",
		},
		{
			name: "schedule day json",
			args: []string{"--json", "--compact", "user", "settings", "mobile-schedule-day", "monday"},
			fake: commandFlowFakeClient{},
			want: `{"start":"09:00","end":"18:00"}` + "\n",
		},
		{
			name: "schedule day without a schedule",
			args: []string{"user", "settings", "mobile-schedule-day", "monday"},
			fake: nullMobile,
			want: "monday start=- end=-\n",
		},
		{
			name: "mobile delivery quiet",
			args: []string{"--quiet", "user", "settings", "mobile-delivery"},
			fake: commandFlowFakeClient{},
			want: "",
		},
		{
			name: "null mobile schedule json",
			args: []string{"--json", "user", "settings", "mobile-schedule"},
			fake: nullMobile,
			want: "null\n",
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
			require.Equal(t, test.want, output.String())
		})
	}
}

func Test_UserSettings_theme_reports_a_preset_linctl_does_not_know_as_sent(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, payloadGraphQLClient{
		payload: `{"userSettings":{"theme":{"preset":"aurora","custom":null}}}`,
	})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "--compact", "user", "settings", "theme"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Equal(t, `{"preset":"aurora"}`+"\n", output.String())
}

func Test_UserSettings_prints_none_when_mobile_delivery_is_null(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, payloadGraphQLClient{
		payload: `{"userSettings":{"notificationDeliveryPreferences":{"mobile":null}}}`,
	})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"user", "settings", "mobile-delivery"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Equal(t, "mobile delivery none\n", output.String())
}

func Test_UserSettings_projects_json_fields_for_nullable_mobile_delivery(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{
		"--json", "--fields", "notifications_disabled", "user", "settings", "mobile-delivery",
	})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"notifications_disabled"`)
	require.NotContains(t, output.String(), `"schedule"`)
}

type commandFailingReader struct{}

func (reader commandFailingReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}
