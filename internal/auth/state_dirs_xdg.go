//go:build !windows && !darwin

package auth

import (
	"errors"
	"path/filepath"
)

func userConfigDir(env Env) (string, error) {
	activeEnv := activeEnv(env)
	if configHome, ok := activeEnv.Lookup("XDG_CONFIG_HOME"); ok && configHome != "" {
		if !filepath.IsAbs(configHome) {
			return "", errors.New("XDG_CONFIG_HOME must be absolute")
		}

		return configHome, nil
	}
	home, err := requiredEnv(activeEnv, "HOME")
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config"), nil
}

func userStateDir(env Env) (string, error) {
	activeEnv := activeEnv(env)
	if stateHome, ok := activeEnv.Lookup("XDG_STATE_HOME"); ok && stateHome != "" {
		if !filepath.IsAbs(stateHome) {
			return "", errors.New("XDG_STATE_HOME must be absolute")
		}

		return stateHome, nil
	}
	home, err := requiredEnv(activeEnv, "HOME")
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".local", "state"), nil
}
