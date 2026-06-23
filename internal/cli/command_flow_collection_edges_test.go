package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CommandFlows_customer_need_project_attachment_handles_missing_attachment(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{missingCustomerNeedAttachment: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"customer-need", "project-attachment", "customer-need-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Equal(t, "customer-need-id project_attachment -\n", output.String())
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

func Test_CommandFlows_fail_on_empty_issue_child_list_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyIssueChildren: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "issue", "children", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_issue_child_list_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "children", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime failed")
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

func Test_CommandFlows_project_comment_children_omit_body_from_json(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "project comments", args: []string{"project", "comments", "project-id", "--json"}},
		{name: "project update comments", args: []string{"project-update", "comments", "project-update-id", "--json"}},
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
			require.Contains(t, output.String(), `"comments"`)
			require.NotContains(t, output.String(), `"body"`)
		})
	}
}

func Test_CommandFlows_project_child_reads_cover_json_and_sort_branches(t *testing.T) {
	t.Run("project milestone issues json", func(t *testing.T) {
		output := bytes.Buffer{}
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(&output)
		command.SetArgs([]string{"project-milestone", "issues", "project-milestone-id", "--json"})

		err := command.ExecuteContext(context.Background())

		require.NoError(t, err)
		require.Contains(t, output.String(), `"project_milestone_id"`)
		require.Contains(t, output.String(), `"issues"`)
	})

	t.Run("project comments sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "project", "comments", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("project comments text output", func(t *testing.T) {
		output := bytes.Buffer{}
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(&output)
		command.SetArgs([]string{"project", "comments", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.NoError(t, err)
		require.Contains(t, output.String(), "comment-id Omer 2026-06-19T12:00:00Z")
	})

	t.Run("release search json", func(t *testing.T) {
		output := bytes.Buffer{}
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(&output)
		command.SetArgs([]string{"release", "search", "mobile", "--json"})

		err := command.ExecuteContext(context.Background())

		require.NoError(t, err)
		require.Contains(t, output.String(), `"releases"`)
	})

	t.Run("release search fail on empty", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{emptyReleaseSearch: true})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--fail-on-empty", "release", "search", "mobile"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "empty result")
	})
}

func Test_CommandFlows_label_child_reads_cover_json_pages(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "label children json",
			args:     []string{"label", "children", "label-id", "--json"},
			contains: `"label_name": "Bug"`,
		},
		{
			name:     "label issues json",
			args:     []string{"label", "issues", "label-id", "--json"},
			contains: `"identifier": "LIT-1"`,
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
			require.Contains(t, output.String(), test.contains)
		})
	}
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

func Test_CommandFlows_fail_on_empty_sla_configurations_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptySLAConfigurations: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "sla-configuration", "list", "team-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_fail_on_empty_semantic_search_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptySemanticSearch: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "semantic-search", "agent search"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_fail_on_empty_typed_search_when_fail_on_empty_flag_is_set(t *testing.T) {
	tests := []struct {
		name string
		args []string
		fake commandFlowFakeClient
	}{
		{
			name: "documents",
			args: []string{"--fail-on-empty", "search", "documents", "agent search"},
			fake: commandFlowFakeClient{emptySearchDocuments: true},
		},
		{
			name: "issues",
			args: []string{"--fail-on-empty", "search", "issues", "agent search"},
			fake: commandFlowFakeClient{emptySearchIssues: true},
		},
		{
			name: "projects",
			args: []string{"--fail-on-empty", "search", "projects", "agent search"},
			fake: commandFlowFakeClient{emptySearchProjects: true},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

func Test_CommandFlows_semantic_search_honors_id_only_and_quiet(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		output string
	}{
		{name: "id only", args: []string{"--id-only", "semantic-search", "agent search"}, output: "issue-id\n"},
		{name: "quiet", args: []string{"--quiet", "semantic-search", "agent search"}, output: ""},
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
			require.Equal(t, test.output, output.String())
		})
	}
}

func Test_CommandFlows_typed_search_writers_emit_json(t *testing.T) {
	tests := []struct {
		name  string
		write func(*cobra.Command, *rootOptions) error
		want  string
	}{
		{
			name: "document",
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSearchDocument(command, options, client.SearchDocumentSummary{
					ID:    "search-document-id",
					Title: "Search spec",
				})
			},
			want: `{"id":"search-document-id","title":"Search spec","slug_id":"","url":""}`,
		},
		{
			name: "issue",
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSearchIssue(command, options, client.SearchIssueSummary{
					ID:         "search-issue-id",
					Identifier: "LIT-30",
					Title:      "Search issue",
				})
			},
			want: `{"id":"search-issue-id","identifier":"LIT-30","title":"Search issue","url":"","team_id":"","team_key":"","team_name":"","state_id":"","state_name":"","state_type":""}`,
		},
		{
			name: "project",
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSearchProject(command, options, client.SearchProjectSummary{
					ID:   "search-project-id",
					Name: "Search project",
				})
			},
			want: `{"id":"search-project-id","name":"Search project","slug_id":"","url":"","status":{"id":"","name":"","type":""},"teams":null}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			command := &cobra.Command{}
			command.SetOut(&output)

			err := test.write(command, &rootOptions{json: true})

			require.NoError(t, err)
			require.JSONEq(t, test.want, output.String())
		})
	}
}

func Test_CommandFlows_typed_search_honors_id_only_and_quiet(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		output string
	}{
		{name: "documents id only", args: []string{"--id-only", "search", "documents", "agent search"}, output: "search-document-id\n"},
		{name: "documents quiet", args: []string{"--quiet", "search", "documents", "agent search"}, output: ""},
		{name: "issues id only", args: []string{"--id-only", "search", "issues", "agent search"}, output: "search-issue-id\n"},
		{name: "issues quiet", args: []string{"--quiet", "search", "issues", "agent search"}, output: ""},
		{name: "projects id only", args: []string{"--id-only", "search", "projects", "agent search"}, output: "search-project-id\n"},
		{name: "projects quiet", args: []string{"--quiet", "search", "projects", "agent search"}, output: ""},
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
			require.Equal(t, test.output, output.String())
		})
	}
}

func Test_CommandFlows_user_drafts_honor_list_controls(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		fake   commandFlowFakeClient
		output string
	}{
		{name: "id only", args: []string{"--id-only", "user", "drafts"}, output: "draft-id\n"},
		{name: "quiet", args: []string{"--quiet", "user", "drafts"}, output: ""},
		{
			name:   "sort",
			args:   []string{"--sort", "parent_key", "--order", "desc", "user", "drafts"},
			output: "draft-id issue LIT-3 Draft issue\n",
		},
		{
			name: "empty",
			args: []string{"--fail-on-empty", "user", "drafts"},
			fake: commandFlowFakeClient{emptyViewerDrafts: true},
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

			if test.name == "empty" {
				require.Error(t, err)
				require.Contains(t, err.Error(), "empty result")
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.output, output.String())
		})
	}
}

func Test_CommandFlows_user_drafts_json_uses_projected_page(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "--sort", "parent_key", "user", "drafts"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"drafts"`)
	require.Contains(t, output.String(), `"parent_key": "LIT-3"`)
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
