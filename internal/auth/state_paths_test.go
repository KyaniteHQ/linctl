package auth

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DefaultPaths_use_os_config_and_state_locations(t *testing.T) {
	root := t.TempDir()
	env := staticEnv{}
	var wantAppConfigPath string
	var wantTokenPath string
	switch runtime.GOOS {
	case "darwin":
		env["HOME"] = root
		base := filepath.Join(root, "Library", "Application Support", "linctl")
		wantAppConfigPath = filepath.Join(base, appConfigFileName)
		wantTokenPath = filepath.Join(base, tokenFileName)
	case "windows":
		env["APPDATA"] = root
		env["LOCALAPPDATA"] = filepath.Join(root, "local")
		wantAppConfigPath = filepath.Join(root, "linctl", appConfigFileName)
		wantTokenPath = filepath.Join(root, "local", "linctl", tokenFileName)
	default:
		env["XDG_CONFIG_HOME"] = filepath.Join(root, "config")
		env["XDG_STATE_HOME"] = filepath.Join(root, "state")
		wantAppConfigPath = filepath.Join(root, "config", "linctl", appConfigFileName)
		wantTokenPath = filepath.Join(root, "state", "linctl", tokenFileName)
	}

	got, err := DefaultPaths(env)

	require.NoError(t, err)
	require.Equal(t, wantAppConfigPath, got.AppConfigPath)
	require.Equal(t, wantTokenPath, got.TokenPath)
	require.NotContains(t, got.AppConfigPath, ".linctl.toml")
	require.NotContains(t, got.TokenPath, ".linctl.toml")
}

func Test_DefaultPaths_fall_back_when_preferred_location_is_unset(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name              string
		goos              string
		env               staticEnv
		wantAppConfigPath string
		wantTokenPath     string
	}{
		{
			name: "windows app data fallback",
			goos: "windows",
			env: staticEnv{
				"APPDATA": filepath.Join(root, "appdata"),
			},
			wantAppConfigPath: filepath.Join(root, "appdata", "linctl", appConfigFileName),
			wantTokenPath:     filepath.Join(root, "appdata", "linctl", tokenFileName),
		},
		{
			name:              "linux home fallback",
			goos:              "linux",
			env:               staticEnv{"HOME": root},
			wantAppConfigPath: filepath.Join(root, ".config", "linctl", appConfigFileName),
			wantTokenPath:     filepath.Join(root, ".local", "state", "linctl", tokenFileName),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.goos != runtime.GOOS {
				t.Skip("location fallback belongs to " + tt.goos)
			}

			got, err := DefaultPaths(tt.env)

			require.NoError(t, err)
			require.Equal(t, tt.wantAppConfigPath, got.AppConfigPath)
			require.Equal(t, tt.wantTokenPath, got.TokenPath)
		})
	}
}
