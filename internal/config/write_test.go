package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
)

func Test_WritePin_round_trips_through_Load(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".linctl.toml")
	want := Target{
		OrgID:     "org-id",
		TeamKey:   "LIT",
		TeamID:    "team-id",
		ProjectID: "project-id",
	}

	require.NoError(t, WritePin(path, want))

	resolved, err := Load(context.Background(), LoadRequest{RepoPath: path})
	require.NoError(t, err)
	require.Equal(t, want, resolved.Target)
}

func Test_WritePin_refuses_overwrite(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".linctl.toml")
	require.NoError(t, WritePin(path, Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	}))

	err := WritePin(path, Target{
		OrgID:   "other-org",
		TeamKey: "OTH",
		TeamID:  "other-team",
	})

	require.ErrorIs(t, err, ErrPinExists)
	body, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	require.Contains(t, string(body), "org-id")
	require.NotContains(t, string(body), "other-org")
}

func Test_WritePin_serializes_only_target_table(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".linctl.toml")
	require.NoError(t, WritePin(path, Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	}))

	body, err := os.ReadFile(path)
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, toml.Unmarshal(body, &decoded))
	require.Equal(t, map[string]any{
		"target": map[string]any{
			"org_id":     "org-id",
			"team_key":   "LIT",
			"team_id":    "team-id",
			"project_id": "",
		},
	}, decoded)
}

func Test_WritePin_requires_core_fields(t *testing.T) {
	err := WritePin(filepath.Join(t.TempDir(), ".linctl.toml"), Target{OrgID: "org-id"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "org_id, team_key, and team_id are required")
}

func Test_WritePin_requires_path(t *testing.T) {
	err := WritePin("", Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "path is required")
}

func Test_WritePin_open_failure_non_exist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-dir", ".linctl.toml")
	err := WritePin(path, Target{OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "write pin")
	require.NotErrorIs(t, err, ErrPinExists)
}
