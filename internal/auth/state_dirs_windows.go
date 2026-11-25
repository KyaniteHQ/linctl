//go:build windows

package auth

func userConfigDir(env Env) (string, error) {
	return requiredEnv(activeEnv(env), "APPDATA")
}

func userStateDir(env Env) (string, error) {
	activeEnv := activeEnv(env)
	if localAppData, ok := activeEnv.Lookup("LOCALAPPDATA"); ok && localAppData != "" {
		return localAppData, nil
	}

	return requiredEnv(activeEnv, "APPDATA")
}
