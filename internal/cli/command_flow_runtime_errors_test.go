package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_CommandFlows_emit_guarded_write_target_error_codes(t *testing.T) {
	tests := []struct {
		name        string
		target      config.Target
		wantErr     error
		wantCode    string
		wantMessage string
	}{
		{
			name:        "missing pinned target",
			target:      config.Target{},
			wantErr:     client.ErrTargetNotConfigured,
			wantCode:    `"error_code":"TARGET_NOT_CONFIGURED"`,
			wantMessage: "set org_id and team_key in .linctl.toml",
		},
		{
			name: "resolved target mismatch",
			target: config.Target{
				OrgID:     "other-org",
				TeamKey:   "LIT",
				TeamID:    "team-id",
				ProjectID: "project-id",
			},
			wantErr:     client.ErrTargetMismatch,
			wantCode:    `"error_code":"TARGET_MISMATCH"`,
			wantMessage: "cannot reach team_key=LIT team_id=team-id",
		},
		{
			name: "team key override without a team id still fails closed",
			target: config.Target{
				OrgID:   "other-org",
				TeamKey: "STG",
			},
			wantErr:     client.ErrTargetMismatch,
			wantCode:    `"error_code":"TARGET_MISMATCH"`,
			wantMessage: "cannot reach team_key=STG",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := buildCommandRuntime
			buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				runtime := testCommandRuntime(commandFlowFakeClient{})
				runtime.config.Target = test.target
				return runtime, nil
			}
			defer func() {
				buildCommandRuntime = original
			}()

			var stdout bytes.Buffer
			var stderr bytes.Buffer

			err := execute(
				context.Background(),
				BuildInfo{},
				strings.NewReader(""),
				&stdout,
				&stderr,
				[]string{"issue", "create", "--title", "Created issue"},
			)

			require.ErrorIs(t, err, test.wantErr)
			require.Empty(t, stdout.String())
			require.Contains(t, stderr.String(), test.wantCode)
			require.Contains(t, stderr.String(), test.wantMessage)
		})
	}
}

func Test_CommandFlows_report_runtime_and_writer_errors(t *testing.T) {
	t.Run("runtime error returns from command", func(t *testing.T) {
		commands := [][]string{
			{"target"},
			{"doctor"},
			{"whoami"},
			{"application", "info", "app-client-id"},
			{"agent-activity", "list"},
			{"audit-entry", "types"},
			{"rate-limit", "status"},
			{"next", "--dry-run"},
			{"files", "upload", "asset.txt"},
			{"issue", "list"},
			{"issue", "search", "needle"},
			{"issue", "priority-values"},
			{"issue", "filter-suggestion", "started issues"},
			{"issue", "title-suggestion", "Customer asks for faster exports"},
			{"issue", "get", "LIT-1"},
			{"issue", "deps", "LIT-1"},
			{"issue", "pr", "LIT-1"},
			{"issue", "create", "--title", "Created issue"},
			{"issue", "update", "LIT-1", "--title", "Updated issue"},
			{"issue", "open", "LIT-1"},
			{"issue", "export", "LIT-1", "."},
			{"issue", "import", "rows.json"},
			{"issue", "bulk-export", "out.json"},
			{"project", "export", "project-id", "."},
			{"project", "open", "project-id"},
			{"user", "me"},
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

				require.ErrorContains(t, err, "runtime failed")
			})
		}
	})

	t.Run("issue list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"issue", "list"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "write line")
	})

	t.Run("doctor returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"doctor"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "write line")
	})

	t.Run("project list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "list"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "write line")
	})

	t.Run("project all reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "project", "all"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, `sort field "missing" is not present`)
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

		require.ErrorContains(t, err, "use only one of --team-id or --project-id")
	})

	t.Run("SLA configuration list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"sla-configuration", "list", "team-id"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "write line")
	})

	t.Run("SLA configuration list reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "sla-configuration", "list", "team-id"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, `sort field "missing" is not present`)
	})

	t.Run("usage returns writer errors", func(t *testing.T) {
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"usage"})

		err := command.ExecuteContext(context.Background())

		require.ErrorContains(t, err, "write line")
	})
}
