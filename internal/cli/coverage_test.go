package cli

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CliRenderHelpers_write_text_and_json_output(t *testing.T) {
	issue := client.IssueSummary{
		Identifier: "LIT-1",
		Title:      "Ship coverage",
		State:      "Todo",
	}
	project := client.ProjectSummary{
		ID:   "project-id",
		Name: "Coverage",
		Status: client.ProjectStatus{
			Name: "Backlog",
		},
	}

	textOut := bytes.Buffer{}
	textCommand := &cobra.Command{}
	textCommand.SetOut(&textOut)
	textOptions := rootOptions{}

	require.NoError(t, writeIssue(textCommand, &textOptions, issue))
	require.NoError(t, writeProject(textCommand, &textOptions, project))
	require.Equal(t, "LIT-1 Ship coverage [Todo]\nproject-id Coverage [Backlog]\n", textOut.String())

	jsonOut := bytes.Buffer{}
	jsonCommand := &cobra.Command{}
	jsonCommand.SetOut(&jsonOut)
	jsonOptions := rootOptions{json: true}

	require.NoError(t, writeIssue(jsonCommand, &jsonOptions, issue))
	require.NoError(t, writeProject(jsonCommand, &jsonOptions, project))
	require.Contains(t, jsonOut.String(), `"identifier": "LIT-1"`)
	require.Contains(t, jsonOut.String(), `"name": "Coverage"`)
}

func Test_CliHelpers_resolve_target_overrides_and_project_ids(t *testing.T) {
	options := rootOptions{
		orgID:   "org-id",
		team:    "LIT",
		project: "project-id",
	}

	target := targetOverride(&options)

	require.Equal(t, "org-id", target.OrgID)
	require.Equal(t, "LIT", target.TeamKey)
	require.Equal(t, "project-id", target.ProjectID)
	require.Empty(t, projectID(nil))
	require.Equal(t, "project-id", projectID(&client.ResolvedProject{ID: "project-id"}))
	require.NotEmpty(t, defaultGlobalConfigPath())
}

func Test_CommandRuntime_loads_config_and_requires_token(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())
	t.Setenv("LINCTL_TOKEN", "")
	t.Setenv("LINEAR_API_KEY", "")
	require.NoError(t, os.WriteFile(".linctl.toml", []byte(`
[target]
org_id = "org-id"
team_key = "LIT"
team_id = "team-id"
project_id = "project-id"
`), 0o600))

	_, err := newCommandRuntime(context.Background(), &rootOptions{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing Linear token")

	t.Setenv("LINCTL_TOKEN", "test-token")
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

	require.Empty(t, defaultGlobalConfigPath())
}

func Test_WriteUsage_reports_unknown_topics(t *testing.T) {
	command := &cobra.Command{}

	err := writeUsage(command, &rootOptions{}, "missing")

	require.Error(t, err)
	require.Contains(t, err.Error(), `unknown usage topic "missing"`)
}
