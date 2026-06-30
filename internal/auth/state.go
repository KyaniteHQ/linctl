// Package auth manages local OAuth application and token state.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	appConfigFileName = "auth-app.json"
	tokenFileName     = "auth-token.json" //nolint:gosec // File name, not a hardcoded credential.
)

var (
	runtimeGOOS = runtime.GOOS
	chmodFile   = os.Chmod
)

// Env resolves environment variables by key.
type Env interface {
	Lookup(key string) (string, bool)
}

// Paths identifies the local auth files.
type Paths struct {
	AppConfigPath string
	TokenPath     string
}

// AppConfig is the saved OAuth application material.
type AppConfig struct {
	ClientID     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
	RedirectURI  string   `json:"redirect_uri,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

// TokenState is the saved OAuth token material.
type TokenState struct {
	AccessToken  string     `json:"access_token,omitempty"`
	RefreshToken string     `json:"refresh_token,omitempty"`
	TokenType    string     `json:"token_type,omitempty"`
	Scopes       []string   `json:"scopes,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	Actor        string     `json:"actor,omitempty"`
	GrantType    string     `json:"grant_type,omitempty"`
}

// ProfileState is auth state scoped to one auth profile.
type ProfileState struct {
	App   AppConfig  `json:"app,omitempty"`
	Token TokenState `json:"token,omitempty"`
}

// State groups local auth state while keeping app config and tokens separate.
type State struct {
	App      AppConfig               `json:"app,omitempty"`
	Token    TokenState              `json:"token,omitempty"`
	Profiles map[string]ProfileState `json:"profiles,omitempty"`
}

type appConfigFile struct {
	App      AppConfig            `json:"app,omitempty"`
	Profiles map[string]AppConfig `json:"profiles,omitempty"`
}

type tokenFile struct {
	Token    TokenState            `json:"token,omitempty"`
	Profiles map[string]TokenState `json:"profiles,omitempty"`
}

type osEnv struct{}

// Store reads and writes local auth state.
type Store struct {
	paths Paths
}

// NewStore returns a local auth state store.
func NewStore(paths Paths) Store {
	return Store{paths: paths}
}

// DefaultPaths returns OS-native user paths for linctl auth state.
func DefaultPaths(env Env) (Paths, error) {
	configDir, err := userConfigDir(env)
	if err != nil {
		return Paths{}, fmt.Errorf("resolve auth app config path: %w", err)
	}
	stateDir, err := userStateDir(env)
	if err != nil {
		return Paths{}, fmt.Errorf("resolve auth token state path: %w", err)
	}

	return Paths{
		AppConfigPath: filepath.Join(configDir, "linctl", appConfigFileName),
		TokenPath:     filepath.Join(stateDir, "linctl", tokenFileName),
	}, nil
}

// Load reads local auth state. Missing files resolve as empty state.
func (store Store) Load(ctx context.Context) (State, error) {
	if err := ctx.Err(); err != nil {
		return State{}, fmt.Errorf("load auth state context: %w", err)
	}

	appState, err := readJSON[appConfigFile](store.paths.AppConfigPath, "auth app config")
	if err != nil {
		return State{}, err
	}
	tokenState, err := readJSON[tokenFile](store.paths.TokenPath, "auth token state")
	if err != nil {
		return State{}, err
	}

	return mergeFiles(appState, tokenState), nil
}

// Save writes local auth state.
func (store Store) Save(ctx context.Context, state State) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("save auth state context: %w", err)
	}

	if err := writeJSON(store.paths.AppConfigPath, appConfigFileFromState(state), "auth app config"); err != nil {
		return err
	}
	return writeJSON(store.paths.TokenPath, tokenFileFromState(state), "auth token state")
}

// SaveAppConfig writes OAuth app configuration without touching token state.
func (store Store) SaveAppConfig(ctx context.Context, profile string, app AppConfig) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("save auth app config context: %w", err)
	}

	appState, err := readJSON[appConfigFile](store.paths.AppConfigPath, "auth app config")
	if err != nil {
		return err
	}
	if profile == "" {
		appState.App = app
	} else {
		if appState.Profiles == nil {
			appState.Profiles = map[string]AppConfig{}
		}
		appState.Profiles[profile] = app
	}

	return writeJSON(store.paths.AppConfigPath, appState, "auth app config")
}

// SaveTokenState writes OAuth token state without touching app configuration.
func (store Store) SaveTokenState(ctx context.Context, profile string, token TokenState) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("save auth token state context: %w", err)
	}

	tokenState, err := readJSON[tokenFile](store.paths.TokenPath, "auth token state")
	if err != nil {
		return err
	}
	if profile == "" {
		tokenState.Token = token
	} else {
		if tokenState.Profiles == nil {
			tokenState.Profiles = map[string]TokenState{}
		}
		tokenState.Profiles[profile] = token
	}

	return writeJSON(store.paths.TokenPath, tokenState, "auth token state")
}

// ClearTokenState removes saved OAuth token material while preserving app config.
func (store Store) ClearTokenState(ctx context.Context, profile string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("clear auth token state context: %w", err)
	}

	tokenState, err := readJSON[tokenFile](store.paths.TokenPath, "auth token state")
	if err != nil {
		return err
	}
	if profile == "" {
		tokenState.Token = TokenState{}
	} else if tokenState.Profiles != nil {
		delete(tokenState.Profiles, profile)
	}
	return writeJSON(store.paths.TokenPath, tokenState, "auth token state")
}

// ClearAppConfig removes saved OAuth app configuration while preserving token state.
func (store Store) ClearAppConfig(ctx context.Context, profile string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("clear auth app config context: %w", err)
	}

	appState, err := readJSON[appConfigFile](store.paths.AppConfigPath, "auth app config")
	if err != nil {
		return err
	}
	if profile == "" {
		appState.App = AppConfig{}
	} else if appState.Profiles != nil {
		delete(appState.Profiles, profile)
	}
	return writeJSON(store.paths.AppConfigPath, appState, "auth app config")
}

// Profile returns the auth state selected by a profile name.
func (state State) Profile(profile string) ProfileState {
	if profile == "" {
		return ProfileState{App: state.App, Token: state.Token}
	}
	profileState, ok := state.Profiles[profile]
	if !ok {
		return ProfileState{}
	}

	return profileState
}

func (env osEnv) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func userConfigDir(env Env) (string, error) {
	activeEnv := activeEnv(env)
	if runtimeGOOS == "windows" {
		return requiredEnv(activeEnv, "APPDATA")
	}
	if runtimeGOOS == "darwin" {
		home, err := requiredEnv(activeEnv, "HOME")
		if err != nil {
			return "", err
		}

		return filepath.Join(home, "Library", "Application Support"), nil
	}
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
	if runtimeGOOS == "windows" {
		if localAppData, ok := activeEnv.Lookup("LOCALAPPDATA"); ok && localAppData != "" {
			return localAppData, nil
		}

		return requiredEnv(activeEnv, "APPDATA")
	}
	if runtimeGOOS == "darwin" {
		return userConfigDir(env)
	}
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

func activeEnv(env Env) Env {
	if env != nil {
		return env
	}

	return osEnv{}
}

func requiredEnv(env Env, key string) (string, error) {
	value, ok := env.Lookup(key)
	if !ok || value == "" {
		return "", fmt.Errorf("%s is not set", key)
	}

	return value, nil
}

func readJSON[T any](path string, label string) (T, error) {
	var zero T
	if path == "" {
		return zero, nil
	}

	//nolint:gosec // Auth paths are resolved from user-specific config/state directories or explicit tests.
	data, err := os.ReadFile(path)
	if errorsIsNotExist(err) {
		return zero, nil
	}
	if err != nil {
		return zero, fmt.Errorf("read %s %s: %w", label, path, err)
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return zero, fmt.Errorf("parse %s %s: %w", label, path, err)
	}

	return value, nil
}

func writeJSON(path string, value any, label string) error {
	if path == "" {
		return nil
	}

	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("encode %s: %w", label, err)
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create %s directory: %w", label, err)
	}
	if err := chmodIfSupported(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("secure %s directory: %w", label, err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write %s %s: %w", label, path, err)
	}
	if err := chmodIfSupported(path, 0o600); err != nil {
		return fmt.Errorf("secure %s %s: %w", label, path, err)
	}

	return nil
}

func errorsIsNotExist(err error) bool {
	return err != nil && os.IsNotExist(err)
}

func appConfigFileFromState(state State) appConfigFile {
	profiles := map[string]AppConfig{}
	for name, profile := range state.Profiles {
		if !appConfigEmpty(profile.App) {
			profiles[name] = profile.App
		}
	}
	if len(profiles) == 0 {
		profiles = nil
	}

	return appConfigFile{App: state.App, Profiles: profiles}
}

func tokenFileFromState(state State) tokenFile {
	profiles := map[string]TokenState{}
	for name, profile := range state.Profiles {
		if !tokenStateEmpty(profile.Token) {
			profiles[name] = profile.Token
		}
	}
	if len(profiles) == 0 {
		profiles = nil
	}

	return tokenFile{Token: state.Token, Profiles: profiles}
}

func appConfigEmpty(app AppConfig) bool {
	return app.ClientID == "" &&
		app.ClientSecret == "" &&
		app.RedirectURI == "" &&
		len(app.Scopes) == 0
}

func tokenStateEmpty(token TokenState) bool {
	return token.AccessToken == "" &&
		token.RefreshToken == "" &&
		token.TokenType == "" &&
		len(token.Scopes) == 0 &&
		token.ExpiresAt == nil &&
		token.Actor == "" &&
		token.GrantType == ""
}

func mergeFiles(appState appConfigFile, tokenState tokenFile) State {
	state := State{
		App:   appState.App,
		Token: tokenState.Token,
	}
	for name, app := range appState.Profiles {
		if state.Profiles == nil {
			state.Profiles = map[string]ProfileState{}
		}
		profile := state.Profiles[name]
		profile.App = app
		state.Profiles[name] = profile
	}
	for name, token := range tokenState.Profiles {
		if state.Profiles == nil {
			state.Profiles = map[string]ProfileState{}
		}
		profile := state.Profiles[name]
		profile.Token = token
		state.Profiles[name] = profile
	}

	return state
}

func chmodIfSupported(path string, mode os.FileMode) error {
	if runtimeGOOS == "windows" {
		return nil
	}

	return chmodFile(path, mode)
}
