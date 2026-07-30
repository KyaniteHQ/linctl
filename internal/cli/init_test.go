package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_Init_writes_pin_for_single_visible_team(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"init"})
	require.NoError(t, command.ExecuteContext(context.Background()))

	resolved, err := config.Load(context.Background(), config.LoadRequest{RepoPath: ".linctl.toml"})
	require.NoError(t, err)
	require.Equal(t, config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	}, resolved.Target)

	body, err := os.ReadFile(".linctl.toml")
	require.NoError(t, err)
	require.Contains(t, string(body), "[target]")
	require.NotContains(t, string(body), "token")
	require.NotContains(t, string(body), "secret")
}

func Test_Init_refuses_existing_pin(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, config.WritePin(".linctl.toml", config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	}))
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"init"})
	err := command.ExecuteContext(context.Background())

	require.ErrorIs(t, err, config.ErrPinExists)
	require.Equal(t, "INVALID_WRITE", errorCode(err))
}

func Test_Init_requires_team_flag_when_multiple_teams(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	restore := useCommandRuntime(t, multiTeamInitFakeClient{})
	defer restore()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"init"})
	err := command.ExecuteContext(context.Background())

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "linctl team list")
	_, statErr := os.Stat(".linctl.toml")
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func Test_Init_selects_team_by_key_when_multiple(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	restore := useCommandRuntime(t, multiTeamInitFakeClient{})
	defer restore()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"init", "--team", "OPS"})
	require.NoError(t, command.ExecuteContext(context.Background()))

	resolved, err := config.Load(context.Background(), config.LoadRequest{RepoPath: ".linctl.toml"})
	require.NoError(t, err)
	require.Equal(t, "OPS", resolved.Target.TeamKey)
	require.Equal(t, "ops-team-id", resolved.Target.TeamID)
}

func Test_Init_pins_verified_project(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"init", "--project", "project-id"})
	require.NoError(t, command.ExecuteContext(context.Background()))

	resolved, err := config.Load(context.Background(), config.LoadRequest{RepoPath: ".linctl.toml"})
	require.NoError(t, err)
	require.Equal(t, "project-id", resolved.Target.ProjectID)
}

func Test_Init_rejects_project_not_on_team(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	restore := useCommandRuntime(t, foreignProjectInitFakeClient{})
	defer restore()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"init", "--project", "foreign-project-id"})
	err := command.ExecuteContext(context.Background())

	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "not attached")
	_, statErr := os.Stat(".linctl.toml")
	require.ErrorIs(t, statErr, os.ErrNotExist)
}

func Test_Init_json_and_quiet(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		var stdout bytes.Buffer
		command.SetOut(&stdout)
		command.SetArgs([]string{"--json", "init"})
		require.NoError(t, command.ExecuteContext(context.Background()))
		require.Contains(t, stdout.String(), `"org_id"`)
		require.Contains(t, stdout.String(), `"team_key"`)
	})
	t.Run("quiet", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		var stdout bytes.Buffer
		command.SetOut(&stdout)
		command.SetArgs([]string{"--quiet", "init"})
		require.NoError(t, command.ExecuteContext(context.Background()))
		require.Empty(t, stdout.String())
		require.FileExists(t, filepath.Join(dir, ".linctl.toml"))
	})
}

func Test_selectInitTeam_table(t *testing.T) {
	single := client.TeamList{Teams: []client.TeamSummary{{
		ID: "team-id", Key: "LIT", Name: "linctl", OrgID: "org-id",
	}}}
	multi := client.TeamList{Teams: []client.TeamSummary{
		{ID: "team-id", Key: "LIT", Name: "linctl", OrgID: "org-id"},
		{ID: "ops-team-id", Key: "OPS", Name: "OPS", OrgID: "org-id"},
	}}

	selected, err := selectInitTeam(single, "", "")
	require.NoError(t, err)
	require.Equal(t, "team-id", selected.ID)

	_, err = selectInitTeam(client.TeamList{}, "", "")
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "no teams")

	_, err = selectInitTeam(multi, "", "")
	require.ErrorIs(t, err, client.ErrWriteInvalid)

	_, err = selectInitTeam(multi, "MISSING", "")
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "team key")

	_, err = selectInitTeam(multi, "", "missing-id")
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "team id")

	selected, err = selectInitTeam(multi, "OPS", "")
	require.NoError(t, err)
	require.Equal(t, "ops-team-id", selected.ID)

	selected, err = selectInitTeam(multi, "", "ops-team-id")
	require.NoError(t, err)
	require.Equal(t, "OPS", selected.Key)

	_, err = selectInitTeam(client.TeamList{
		Teams:       []client.TeamSummary{{ID: "team-id", Key: "LIT", OrgID: "org-id"}},
		HasNextPage: true,
	}, "", "")
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "multiple teams")
}

func Test_Init_rejects_both_team_flags(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"init", "--team", "LIT", "--team-id", "team-id"})
	err := command.ExecuteContext(context.Background())
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "use only one")
}

func Test_Init_runtime_and_list_errors(t *testing.T) {
	t.Run("runtime", func(t *testing.T) {
		original := buildCommandRuntime
		buildCommandRuntime = func(context.Context, *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime boom")
		}
		t.Cleanup(func() { buildCommandRuntime = original })
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"init"})
		require.ErrorContains(t, command.ExecuteContext(context.Background()), "runtime boom")
	})
	t.Run("list teams", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "teams_list"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"init"})
		require.Error(t, command.ExecuteContext(context.Background()))
	})
	t.Run("project get", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: "project"})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"init", "--project", "project-id"})
		require.Error(t, command.ExecuteContext(context.Background()))
	})
}

func Test_verifyInitProject_truncated_teams(t *testing.T) {
	err := verifyInitProject(
		context.Background(),
		truncatedTeamsProjectFake{},
		client.TeamSummary{ID: "team-id", Key: "LIT"},
		"project-id",
	)
	require.ErrorIs(t, err, client.ErrWriteInvalid)
	require.ErrorContains(t, err, "could not be fully verified")
}

type truncatedTeamsProjectFake struct{}

func (truncatedTeamsProjectFake) MakeRequest(
	_ context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName != "project" {
		return errors.New("unexpected operation " + request.OpName)
	}
	payload := `{
		"project":{
			"id":"project-id","name":"Truncated","description":"","slugId":"t",
			"url":"https://linear.app/p","priority":0,
			"status":{"id":"s","name":"Backlog","type":"backlog"},
			"lead":null,
			"teams":{"nodes":[{"id":"other","key":"OTH","name":"Other"}],
			"pageInfo":{"hasNextPage":true,"endCursor":"c"}}
		}
	}`

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

// multiTeamInitFakeClient returns two teams on teams_list.
type multiTeamInitFakeClient struct {
	commandFlowFakeClient
}

func (client multiTeamInitFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName == "teams_list" {
		payload := `{"teams":{"nodes":[` +
			`{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}},` +
			`{"id":"ops-team-id","key":"OPS","name":"OPS","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}` +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`

		return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
	}

	return client.commandFlowFakeClient.MakeRequest(ctx, request, response)
}

// foreignProjectInitFakeClient serves a project not attached to the selected team.
type foreignProjectInitFakeClient struct {
	commandFlowFakeClient
}

func (client foreignProjectInitFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName == "project" {
		payload := `{"project":` + commandProjectJSONWithTeams(
			"Foreign",
			"Backlog",
			"backlog",
			[]commandProjectTeam{{ID: "other-team", Key: "OTH", Name: "Other"}},
		) + `}`

		return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
	}

	return client.commandFlowFakeClient.MakeRequest(ctx, request, response)
}
