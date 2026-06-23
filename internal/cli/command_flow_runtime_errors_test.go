package cli

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CommandFlows_report_runtime_and_writer_errors(t *testing.T) {
	t.Run("runtime error returns from command", func(t *testing.T) {
		commands := [][]string{
			{"target"},
			{"doctor"},
			{"whoami"},
			{"application", "info", "app-client-id"},
			{"agent-activity", "list"},
			{"agent-activity", "get", "agent-activity-id"},
			{"agent-skill", "list"},
			{"agent-skill", "get", "agent-skill-id"},
			{"audit-entry", "types"},
			{"triage-responsibility", "list"},
			{"triage-responsibility", "get", "triage-responsibility-id"},
			{"triage-responsibility", "manual-selection", "triage-responsibility-id"},
			{"organization", "exists", "kyanite"},
			{"semantic-search", "agent search"},
			{"search", "documents", "agent search"},
			{"search", "issues", "agent search"},
			{"search", "projects", "agent search"},
			{"rate-limit", "status"},
			{"release", "list"},
			{"release", "search", "mobile"},
			{"release", "get", "release-id"},
			{"external-link", "get", "release-link-id"},
			{"release-note", "list"},
			{"release-note", "get", "release-note-id"},
			{"next", "--dry-run"},
			{"files", "upload", "asset.txt"},
			{"issue", "list"},
			{"issue", "search", "needle"},
			{"issue", "figma-file-key-search", "figma-key"},
			{"issue", "priority-values"},
			{"issue", "filter-suggestion", "started issues"},
			{"issue", "title-suggestion", "Customer asks for faster exports"},
			{"issue", "get", "LIT-1"},
			{"issue", "deps", "LIT-1"},
			{"issue", "pr", "LIT-1"},
			{"issue", "create", "--title", "Created issue"},
			{"issue", "create", "--title", "Created issue", "--state", "todo"},
			{"issue", "create", "--title", "Created issue", "--priority", "high"},
			{"issue", "update", "LIT-1", "--title", "Updated issue"},
			{"issue", "update", "LIT-1", "--state", "done"},
			{"issue", "update", "LIT-1", "--priority", "2"},
			{"issue", "start", "LIT-1"},
			{"issue", "comment", "LIT-1", "--body", "Looks good"},
			{"issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"},
			{"issue", "comments", "LIT-1"},
			{"issue", "close", "LIT-1"},
			{"issue", "relate", "LIT-1", "LIT-2", "--type", "related"},
			{"issue", "unrelate", "issue-relation-id"},
			{"issue", "open", "LIT-1"},
			{"issue", "export", "LIT-1", "."},
			{"issue", "import", "rows.json"},
			{"issue", "bulk-export", "out.json"},
			{"project", "list"},
			{"project", "all"},
			{"project", "get", "project-id"},
			{"project", "members", "project-id"},
			{"project", "updates", "project-id"},
			{"project-milestone", "list", "project-id"},
			{"project-milestone", "get", "project-milestone-id"},
			{"project-milestone", "create", "project-id", "--name", "Created milestone"},
			{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"},
			{"project-update", "create", "project-id", "--health", "on-track", "--body", "Posted update"},
			{"project-status", "project-count", "project-status-id"},
			{"project", "create", "--name", "Created project"},
			{"project", "update", "project-id", "--name", "Updated project"},
			{"project", "archive", "project-id"},
			{"project", "open", "project-id"},
			{"document", "list"},
			{"document", "get", "document-id"},
			{"document", "create", "--title", "Created doc"},
			{"document", "update", "document-id", "--title", "Updated doc"},
			{"comment", "update", "comment-id", "--body", "New body"},
			{"comment", "delete", "comment-id"},
			{"label", "list"},
			{"label", "get", "label-id"},
			{"team", "list"},
			{"team", "get", "team-id"},
			{"team", "members", "team-id"},
			{"user", "list"},
			{"user", "get", "user-id"},
			{"user", "me"},
			{"user", "assigned-issues", "user-id"},
			{"user", "created-issues", "user-id"},
			{"user", "delegated-issues", "user-id"},
			{"user", "team-memberships", "user-id"},
			{"user", "teams", "user-id"},
			{"user", "my-assigned-issues"},
			{"user", "my-created-issues"},
			{"user", "my-delegated-issues"},
			{"user", "my-team-memberships"},
			{"user", "my-teams"},
			{"custom-view", "subscribers", "custom-view-id"},
			{"custom-view", "initiatives", "custom-view-id"},
			{"custom-view", "issues", "custom-view-id"},
			{"custom-view", "organization-preferences", "custom-view-id"},
			{"custom-view", "organization-preferences", "values", "custom-view-id"},
			{"custom-view", "projects", "custom-view-id"},
			{"custom-view", "user-preferences", "custom-view-id"},
			{"custom-view", "user-preferences", "values", "custom-view-id"},
			{"custom-view", "preference-values", "custom-view-id"},
			{"customer-need", "project-attachment", "customer-need-id"},
			{"sla-configuration", "list", "team-id"},
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

	t.Run("project all returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "all"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project all reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "project", "all"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("issue figma file key search reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{expectedIssueFigmaFileKey: "figma-key"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "issue", "figma-file-key-search", "figma-key"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("issue filter suggestion rejects conflicting scope flags", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{
			"issue",
			"filter-suggestion",
			"started issues",
			"--team-id",
			"team-id",
			"--project-id",
			"project-id",
		})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "use only one of --team-id or --project-id")
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

	t.Run("release search returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"release", "search", "mobile"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("release search reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "release", "search", "mobile"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("SLA configuration list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"sla-configuration", "list", "team-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("SLA configuration list reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "sla-configuration", "list", "team-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("semantic search returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"semantic-search", "agent search"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("semantic search reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "semantic-search", "agent search"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("typed search returns writer errors", func(t *testing.T) {
		tests := [][]string{
			{"search", "documents", "agent search"},
			{"search", "issues", "agent search"},
			{"search", "projects", "agent search"},
		}
		for _, args := range tests {
			t.Run(strings.Join(args[:2], " "), func(t *testing.T) {
				restore := useCommandRuntime(t, commandFlowFakeClient{})
				defer restore()
				command := NewRootCommand(context.Background(), BuildInfo{})
				command.SetOut(commandFailingWriter{})
				command.SetArgs(args)

				err := command.ExecuteContext(context.Background())

				require.Error(t, err)
				require.Contains(t, err.Error(), "write line")
			})
		}
	})

	t.Run("typed search reports sort errors", func(t *testing.T) {
		tests := [][]string{
			{"--sort", "missing", "search", "documents", "agent search"},
			{"--sort", "missing", "search", "issues", "agent search"},
			{"--sort", "missing", "search", "projects", "agent search"},
		}
		for _, args := range tests {
			t.Run(strings.Join(args[2:4], " "), func(t *testing.T) {
				restore := useCommandRuntime(t, commandFlowFakeClient{})
				defer restore()
				command := NewRootCommand(context.Background(), BuildInfo{})
				command.SetArgs(args)

				err := command.ExecuteContext(context.Background())

				require.Error(t, err)
				require.Contains(t, err.Error(), `sort field "missing" is not present`)
			})
		}
	})

	t.Run("issue child list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"issue", "history", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("issue child list reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "issue", "children", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
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
