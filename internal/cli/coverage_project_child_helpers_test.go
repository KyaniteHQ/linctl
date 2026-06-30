package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_ProjectChildListCommand_covers_helper_branches(t *testing.T) {
	type item struct {
		ID string `json:"id"`
	}
	type list struct {
		Items []item `json:"items"`
	}

	tests := []struct {
		name        string
		options     rootOptions
		fetch       func(commandRuntime, string, int) (list, error)
		sortList    func(list) (list, error)
		writeItem   readListItemWriter[item]
		requirement func(*testing.T, string, error)
	}{
		{
			name: "runtime error",
			fetch: func(commandRuntime, string, int) (list, error) {
				return list{}, nil
			},
			sortList: func(value list) (list, error) {
				return value, nil
			},
			writeItem: func(*cobra.Command, *rootOptions, item) error {
				return nil
			},
			requirement: func(t *testing.T, _ string, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "runtime failed")
			},
		},
		{
			name: "fetch error",
			fetch: func(commandRuntime, string, int) (list, error) {
				return list{}, errors.New("fetch failed")
			},
			sortList: func(value list) (list, error) {
				return value, nil
			},
			writeItem: func(*cobra.Command, *rootOptions, item) error {
				return nil
			},
			requirement: func(t *testing.T, _ string, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "fetch failed")
			},
		},
		{
			name:    "empty error",
			options: rootOptions{failOnEmpty: true},
			fetch: func(commandRuntime, string, int) (list, error) {
				return list{}, nil
			},
			sortList: func(value list) (list, error) {
				return value, nil
			},
			writeItem: func(*cobra.Command, *rootOptions, item) error {
				return nil
			},
			requirement: func(t *testing.T, _ string, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "empty result")
			},
		},
		{
			name: "sort error",
			fetch: func(commandRuntime, string, int) (list, error) {
				return list{Items: []item{{ID: "item-id"}}}, nil
			},
			sortList: func(list) (list, error) {
				return list{}, errors.New("sort failed")
			},
			writeItem: func(*cobra.Command, *rootOptions, item) error {
				return nil
			},
			requirement: func(t *testing.T, _ string, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "sort failed")
			},
		},
		{
			name:    "json output",
			options: rootOptions{json: true},
			fetch: func(commandRuntime, string, int) (list, error) {
				return list{Items: []item{{ID: "item-id"}}}, nil
			},
			sortList: func(value list) (list, error) {
				return value, nil
			},
			writeItem: func(*cobra.Command, *rootOptions, item) error {
				return nil
			},
			requirement: func(t *testing.T, output string, err error) {
				require.NoError(t, err)
				require.JSONEq(t, `{"items":[{"id":"item-id"}]}`, output)
			},
		},
		{
			name: "write error",
			fetch: func(commandRuntime, string, int) (list, error) {
				return list{Items: []item{{ID: "item-id"}}}, nil
			},
			sortList: func(value list) (list, error) {
				return value, nil
			},
			writeItem: func(*cobra.Command, *rootOptions, item) error {
				return errors.New("write failed")
			},
			requirement: func(t *testing.T, _ string, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "write failed")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			root := &cobra.Command{Use: "root"}
			options := test.options
			if options.sortOrder == "" {
				options.sortOrder = "asc"
			}
			originalBuildCommandRuntime := buildCommandRuntime
			if test.name == "runtime error" {
				buildCommandRuntime = func(context.Context, *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, errors.New("runtime failed")
				}
			} else {
				buildCommandRuntime = func(context.Context, *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, nil
				}
			}
			t.Cleanup(func() {
				buildCommandRuntime = originalBuildCommandRuntime
			})
			addChildListCommand(
				context.Background(),
				root,
				&options,
				"children PROJECT_ID",
				"List children",
				"children",
				test.fetch,
				func(value list) int {
					return len(value.Items)
				},
				test.sortList,
				test.writeItem,
				func(value list) []item {
					return value.Items
				},
			)
			root.SetOut(&output)
			root.SetArgs([]string{"children", "project-id"})

			err := root.ExecuteContext(context.Background())

			test.requirement(t, output.String(), err)
		})
	}
}

func Test_ProjectHistoryWriter_covers_output_modes(t *testing.T) {
	history := client.ProjectHistorySummary{
		ID:         "project-history-id",
		ProjectID:  "project-id",
		EntryCount: 1,
	}
	tests := []struct {
		name     string
		options  rootOptions
		expected string
	}{
		{
			name:     "id only",
			options:  rootOptions{idOnly: true},
			expected: "project-history-id\n",
		},
		{
			name:     "quiet",
			options:  rootOptions{quiet: true},
			expected: "",
		},
		{
			name:     "json",
			options:  rootOptions{json: true},
			expected: `{"id":"project-history-id","project_id":"project-id","entry_count":1,"entries":null,"created_at":"","updated_at":""}` + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			command := &cobra.Command{}
			command.SetOut(&output)

			err := writeProjectHistory(command, &test.options, history)

			require.NoError(t, err)
			if test.options.json {
				require.JSONEq(t, test.expected, output.String())
				return
			}
			require.Equal(t, test.expected, output.String())
		})
	}
}

func Test_CommandFlows_project_child_lists_route_output_policy(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expected   string
		assertJSON bool
	}{
		{
			name:       "members json emits page object",
			args:       []string{"--json", "project", "members", "project-id", "--limit", "1"},
			expected:   `{"project_id":"project-id","project_name":"Detail project","members":[{"id":"user-id","name":"omer","display_name":"Omer","email":"omer@example.com"}],"has_next_page":false}`,
			assertJSON: true,
		},
		{
			name:     "members quiet emits nothing",
			args:     []string{"--quiet", "project", "members", "project-id", "--limit", "1"},
			expected: "",
		},
		{
			name:     "members id-only emits ids",
			args:     []string{"--id-only", "project", "members", "project-id", "--limit", "1"},
			expected: "user-id\n",
		},
		{
			name:     "members human output",
			args:     []string{"project", "members", "project-id", "--limit", "1"},
			expected: "user-id Omer\n",
		},
		{
			name:       "updates json emits page object",
			args:       []string{"--json", "project", "updates", "project-id", "--limit", "1"},
			expected:   `{"project_id":"project-id","project_name":"Detail project","updates":[{"id":"project-update-id","health":"onTrack","created_at":"2026-06-19T12:00:00Z","updated_at":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","user_id":"user-id","name":"omer","display_name":"Omer"}],"has_next_page":false}`,
			assertJSON: true,
		},
		{
			name:     "updates quiet emits nothing",
			args:     []string{"--quiet", "project", "updates", "project-id", "--limit", "1"},
			expected: "",
		},
		{
			name:     "updates id-only emits ids",
			args:     []string{"--id-only", "project", "updates", "project-id", "--limit", "1"},
			expected: "project-update-id\n",
		},
		{
			name:     "updates human output",
			args:     []string{"project", "updates", "project-id", "--limit", "1"},
			expected: "project-update-id onTrack Omer\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			if test.assertJSON {
				require.JSONEq(t, test.expected, output.String())
				return
			}
			require.Equal(t, test.expected, output.String())
		})
	}
}

func Test_CommandFlows_cover_issue_deps_writer_error(t *testing.T) {
	t.Run("issue header", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("section header", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssueDependencySection(command, &rootOptions{}, "children", nil)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("parent issue", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 3})
		parent := client.IssueSummary{Identifier: "LIT-2", Title: "Parent", State: "Todo"}
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1", Parent: &parent}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("children section", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 2})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("blocks section", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 3})
		dependencies := client.IssueDependencyGraph{Identifier: "LIT-1"}

		err := writeIssueDependencies(command, &rootOptions{}, dependencies)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_cover_rate_limit_writer_errors(t *testing.T) {
	status := client.RateLimitStatus{
		Identifier: "api-key",
		Kind:       "api",
		Limits: []client.RateLimit{
			{Type: "complexity", AllowedAmount: 1000, RemainingAmount: 900, Reset: 1720000000000},
		},
	}

	t.Run("header", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeRateLimitStatus(command, &rootOptions{}, status)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("limit", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(&countFailingWriter{failAt: 2})

		err := writeRateLimitStatus(command, &rootOptions{}, status)

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_cover_audit_entry_type_writer_errors(t *testing.T) {
	command := &cobra.Command{}
	command.SetOut(commandFailingWriter{})
	types := client.AuditEntryTypeList{
		AuditEntryTypes: []client.AuditEntryTypeSummary{
			{Type: "user_login", Description: "User logged in"},
		},
	}

	err := writeAuditEntryTypes(command, &rootOptions{}, types)

	require.Error(t, err)
	require.Contains(t, err.Error(), "write line")
}

type countFailingWriter struct {
	failAt int
	writes int
}

func (writer *countFailingWriter) Write(content []byte) (int, error) {
	writer.writes++
	if writer.writes == writer.failAt {
		return 0, errors.New("write failed")
	}

	return len(content), nil
}

func Test_CliHelpers_resolve_target_overrides_and_project_ids(t *testing.T) {
	options := rootOptions{
		orgID:   "org-id",
		team:    "LIT",
		teamID:  "team-id",
		project: "project-id",
	}

	target := targetOverride(&options)

	require.Equal(t, "org-id", target.OrgID)
	require.Equal(t, "LIT", target.TeamKey)
	require.Equal(t, "team-id", target.TeamID)
	require.Equal(t, "project-id", target.ProjectID)
	require.Empty(t, projectID(nil))
	require.Equal(t, "project-id", projectID(&client.ResolvedProject{ID: "project-id"}))
	require.NotEmpty(t, defaultGlobalConfigPath())
}

func Test_CliHelpers_clear_config_team_id_when_org_or_team_override_changes_target(t *testing.T) {
	resolved := config.Resolved{Target: config.Target{
		OrgID:     "configured-org",
		TeamKey:   "CFG",
		TeamID:    "configured-team-id",
		ProjectID: "project-id",
	}}

	applyTargetOverrideFlagSemantics(&resolved, &rootOptions{team: "LIT"})

	require.Equal(t, config.Target{
		OrgID:     "configured-org",
		TeamKey:   "CFG",
		ProjectID: "project-id",
	}, resolved.Target)
}

func Test_CommandRuntime_loads_config_and_requires_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	configureTestAuthEnvironment(t)
	require.NoError(t, os.WriteFile(".linctl.toml", []byte(`
[target]
org_id = "org-id"
team_key = "LIT"
team_id = "team-id"
project_id = "project-id"
`), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{})
	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeNotConfigured), errorCode(err))
	require.Contains(t, err.Error(), "run linctl auth configure")

	t.Setenv("LINCTL_TOKEN", "test-token")
	_, err = newCommandRuntime(context.Background(), &rootOptions{})
	require.Error(t, err)
	require.Equal(t, string(auth.ErrorCodeNotConfigured), errorCode(err))
	require.Contains(t, err.Error(), "run linctl auth configure")

	t.Setenv("LINCTL_OAUTH_ACCESS_TOKEN", "test-token")
	runtime, err := newCommandRuntime(context.Background(), &rootOptions{})
	require.NoError(t, err)
	require.Equal(t, "project-id", runtime.config.Target.ProjectID)
	require.NotNil(t, runtime.graphqlClient)
}

func Test_CommandRuntime_reports_config_load_errors(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())
	require.NoError(t, os.WriteFile(".linctl.toml", []byte("[target\n"), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "parse config")
}

func Test_DefaultGlobalConfigPath_returns_empty_when_home_is_unset(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "") // Windows resolves the home dir from USERPROFILE, not HOME.

	require.Empty(t, defaultGlobalConfigPath())
}

func Test_WriteUsage_reports_unknown_topics(t *testing.T) {
	command := &cobra.Command{}

	err := writeUsage(command, &rootOptions{}, "missing")

	require.Error(t, err)
	require.Contains(t, err.Error(), `unknown usage topic "missing"`)
}
