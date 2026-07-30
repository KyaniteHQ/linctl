package config

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

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
	require.Contains(t, string(body), "[target]")
	require.NotContains(t, string(body), "token")
	require.NotContains(t, string(body), "secret")
	require.NotContains(t, string(body), "profile")
	require.NotContains(t, string(body), "client")

	keys, err := SerializedPinKeys(body)
	require.NoError(t, err)
	require.Equal(t, []string{"org_id", "team_key", "team_id"}, keys)
}

func Test_WritePin_target_struct_tags_match_pin_keys(t *testing.T) {
	targetType := reflect.TypeOf(Target{})
	tags := make([]string, 0, targetType.NumField())
	for index := range targetType.NumField() {
		tag := targetType.Field(index).Tag.Get("toml")
		require.NotEmpty(t, tag)
		tags = append(tags, tag)
	}
	require.Equal(t, PinFileKeys(), tags)
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

func Test_WritePin_encode_write_close_seams(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		original := marshalPin
		marshalPin = func(any) ([]byte, error) { return nil, errors.New("encode boom") }
		t.Cleanup(func() { marshalPin = original })
		err := WritePin(filepath.Join(t.TempDir(), "pin.toml"), Target{
			OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id",
		})
		require.ErrorContains(t, err, "encode")
	})
	t.Run("write", func(t *testing.T) {
		original := openPinFile
		openPinFile = func(string, int, os.FileMode) (pinWriter, error) {
			return failingPinWriter{writeErr: errors.New("write boom")}, nil
		}
		t.Cleanup(func() { openPinFile = original })
		err := WritePin(filepath.Join(t.TempDir(), "pin.toml"), Target{
			OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id",
		})
		require.ErrorContains(t, err, "write boom")
	})
	t.Run("close", func(t *testing.T) {
		original := openPinFile
		openPinFile = func(string, int, os.FileMode) (pinWriter, error) {
			return failingPinWriter{closeErr: errors.New("close boom")}, nil
		}
		t.Cleanup(func() { openPinFile = original })
		err := WritePin(filepath.Join(t.TempDir(), "pin.toml"), Target{
			OrgID: "org-id", TeamKey: "LIT", TeamID: "team-id",
		})
		require.ErrorContains(t, err, "close boom")
	})
}

func Test_SerializedPinKeys_edges(t *testing.T) {
	keys, err := SerializedPinKeys([]byte(`
[target]
org_id = "org-id"
team_key = "LIT"
team_id = "team-id"
project_id = "project-id"
`))
	require.NoError(t, err)
	require.Equal(t, []string{"org_id", "team_key", "team_id", "project_id"}, keys)

	_, err = SerializedPinKeys([]byte(`not = [toml`))
	require.Error(t, err)

	_, err = SerializedPinKeys([]byte(`
[target]
org_id = "token-looking-org"
team_key = "LIT"
team_id = "team-id"
`))
	// "token" substring in org value trips the auth-material guard
	require.Error(t, err)
	require.Contains(t, err.Error(), "auth material")
}

type failingPinWriter struct {
	writeErr error
	closeErr error
}

func (writer failingPinWriter) Write(payload []byte) (int, error) {
	if writer.writeErr != nil {
		return 0, writer.writeErr
	}

	return len(payload), nil
}

func (writer failingPinWriter) Close() error {
	return writer.closeErr
}
