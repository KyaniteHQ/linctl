package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Load_resolves_repo_profile_and_env_token_when_present(t *testing.T) {
	// Given
	root := t.TempDir()
	globalPath := filepath.Join(root, "global.toml")
	repoPath := filepath.Join(root, "repo.toml")
	require.NoError(t, os.WriteFile(globalPath, []byte(`
profile = "global"

[profiles.global]
token = "global-token"

[profiles.global.target]
org_id = "global-org"
team_key = "GLOBAL"
team_id = "global-team"
`), 0o600))
	require.NoError(t, os.WriteFile(repoPath, []byte(`
profile = "repo"

[profiles.repo]
token = "repo-token"

[profiles.repo.target]
org_id = "repo-org"
team_key = "REPO"
team_id = "repo-team"
project_id = "repo-project"
`), 0o600))

	// When
	resolved, err := Load(context.Background(), LoadRequest{
		GlobalPath: globalPath,
		RepoPath:   repoPath,
		Env: StaticEnv{
			"LINCTL_TOKEN": "env-token",
		},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "repo", resolved.Profile)
	require.Equal(t, "env-token", resolved.Token)
	require.Equal(t, Target{
		OrgID:     "repo-org",
		TeamKey:   "REPO",
		TeamID:    "repo-team",
		ProjectID: "repo-project",
	}, resolved.Target)
}

func Test_Load_applies_explicit_profile_and_target_overrides_when_present(t *testing.T) {
	// Given
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
profile = "default"

[profiles.default]
token = "default-token"

[profiles.default.target]
org_id = "default-org"
team_key = "DEF"
team_id = "default-team"

[profiles.other]
token = "other-token"

[profiles.other.target]
org_id = "other-org"
team_key = "OTH"
team_id = "other-team"
`), 0o600))

	// When
	resolved, err := Load(context.Background(), LoadRequest{
		GlobalPath:      configPath,
		ProfileOverride: "other",
		TargetOverride: Target{
			ProjectID: "override-project",
		},
		Env: StaticEnv{},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "other", resolved.Profile)
	require.Equal(t, "other-token", resolved.Token)
	require.Equal(t, Target{
		OrgID:     "other-org",
		TeamKey:   "OTH",
		TeamID:    "other-team",
		ProjectID: "override-project",
	}, resolved.Target)
}

func Test_Load_keeps_profile_targets_separate_when_multiple_workspaces_exist(t *testing.T) {
	// Given
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[profiles.personal]
token = "personal-token"

[profiles.personal.target]
org_id = "personal-org"
team_key = "PER"
team_id = "personal-team"

[profiles.work]
token = "work-token"

[profiles.work.target]
org_id = "work-org"
team_key = "WRK"
team_id = "work-team"
project_id = "work-project"
`), 0o600))

	// When
	resolved, err := Load(context.Background(), LoadRequest{
		GlobalPath:      configPath,
		ProfileOverride: "work",
		Env: StaticEnv{
			"LINCTL_TOKEN": "env-token",
		},
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "work", resolved.Profile)
	require.Equal(t, "env-token", resolved.Token)
	require.Equal(t, Target{
		OrgID:     "work-org",
		TeamKey:   "WRK",
		TeamID:    "work-team",
		ProjectID: "work-project",
	}, resolved.Target)
}

func Test_Load_refuses_unknown_explicit_profile(t *testing.T) {
	// Given
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[profiles.default]
token = "default-token"
`), 0o600))

	// When
	_, err := Load(context.Background(), LoadRequest{
		GlobalPath:      configPath,
		ProfileOverride: "missing",
		Env:             StaticEnv{},
	})

	// Then
	require.ErrorIs(t, err, ErrProfileNotFound)
}
