package config

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Load_resolves_repo_profile_when_present(t *testing.T) {
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
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "repo", resolved.Profile)
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
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "other", resolved.Profile)
	require.Equal(t, Target{
		OrgID:     "other-org",
		TeamKey:   "OTH",
		TeamID:    "other-team",
		ProjectID: "override-project",
	}, resolved.Target)
}

func Test_Load_keeps_profile_targets_separate_when_multiple_targets_exist(t *testing.T) {
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
	})

	// Then
	require.NoError(t, err)
	require.Equal(t, "work", resolved.Profile)
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
	})

	// Then
	require.ErrorIs(t, err, ErrProfileNotFound)
}

func Test_Load_ignores_legacy_config_token_without_permission_gate(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`token = "file-token"`), 0o644))

	resolved, err := Load(context.Background(), LoadRequest{
		RepoPath: configPath,
	})

	require.NoError(t, err)
	require.Empty(t, resolved.Profile)
}

func Test_Load_ignores_legacy_config_tokens_for_product_auth(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
token = "file-token"

[profiles.default]
token = "profile-token"
`), 0o600))

	_, err := Load(context.Background(), LoadRequest{
		RepoPath:        configPath,
		ProfileOverride: "default",
	})

	require.NoError(t, err)
}

func Test_Load_allows_broad_config_permissions_without_tokens(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows file modes do not expose unix group/world permission bits")
	}
	root := t.TempDir()
	configPath := filepath.Join(root, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
[target]
team_key = "ENG"
`), 0o644))

	resolved, err := Load(context.Background(), LoadRequest{
		RepoPath: configPath,
	})

	require.NoError(t, err)
	require.Equal(t, "ENG", resolved.Target.TeamKey)
}

func Test_Load_reports_config_read_error(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config-dir")
	require.NoError(t, os.Mkdir(configPath, 0o700))

	_, err := Load(context.Background(), LoadRequest{
		RepoPath: configPath,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read config")
}

func Test_Load_reports_bad_config_path_error(t *testing.T) {
	_, err := Load(context.Background(), LoadRequest{
		RepoPath: "bad\x00path",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "read config")
}
