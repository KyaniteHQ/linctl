//go:build darwin

package auth

import "path/filepath"

func userConfigDir(env Env) (string, error) {
	home, err := requiredEnv(activeEnv(env), "HOME")
	if err != nil {
		return "", err
	}

	return filepath.Join(home, "Library", "Application Support"), nil
}

func userStateDir(env Env) (string, error) {
	return userConfigDir(env)
}
