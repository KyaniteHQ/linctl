// Package auth manages local OAuth application and token state.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"
)

const (
	appConfigFileName = "auth-app.json"
	tokenFileName     = "auth-token.json" //nolint:gosec // File name, not a hardcoded credential.
)

var (
	runtimeGOOS    = runtime.GOOS
	chmodFile      = os.Chmod
	writeTempFile  = (*os.File).Write
	syncTempFile   = (*os.File).Sync
	closeTempFile  = (*os.File).Close
	lockAuthFile   = lockStateFile
	unlockAuthFile = unlockStateFile
)

var authFileLocks = struct {
	sync.Mutex
	byPath map[string]*contextLock
}{byPath: map[string]*contextLock{}}

type contextLock struct {
	available chan struct{}
}

func newContextLock() *contextLock {
	lock := &contextLock{available: make(chan struct{}, 1)}
	lock.available <- struct{}{}
	return lock
}

func (lock *contextLock) acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-lock.available:
		return nil
	}
}

func (lock *contextLock) release() {
	lock.available <- struct{}{}
}

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
//
// Secret-bearing fields (ClientSecret) must appear in String/GoString/MarshalJSON
// redaction below AND in persistedAppConfig for on-disk persistence.
type AppConfig struct {
	ClientID     string   `json:"client_id,omitempty"`
	ClientSecret string   `json:"client_secret,omitempty"`
	RedirectURI  string   `json:"redirect_uri,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

// String redacts secret material for fmt's %v, %+v, and %s verbs.
func (app AppConfig) String() string {
	return fmt.Sprintf(
		"auth.AppConfig{ClientID:%s, ClientSecret:%s, RedirectURI:%s, Scopes:%v}",
		presenceValue(app.ClientID), presenceValue(app.ClientSecret), presenceValue(app.RedirectURI), app.Scopes,
	)
}

// GoString redacts secret material for fmt's %#v verb.
func (app AppConfig) GoString() string {
	return app.String()
}

// MarshalJSON redacts secret material so json.Marshal never emits raw values.
func (app AppConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURI  string   `json:"redirect_uri"`
		Scopes       []string `json:"scopes"`
	}{
		ClientID:     presenceValue(app.ClientID),
		ClientSecret: presenceValue(app.ClientSecret),
		RedirectURI:  presenceValue(app.RedirectURI),
		Scopes:       app.Scopes,
	})
}

// persistedAppConfig carries the same fields as AppConfig but drops its redaction
// methods, so marshaling it writes real values to disk.
type persistedAppConfig AppConfig

// TokenState is the saved OAuth token material.
//
// Secret-bearing fields (AccessToken, RefreshToken) must appear in
// String/GoString/MarshalJSON redaction below AND in persistedTokenState for
// on-disk persistence.
type TokenState struct {
	AccessToken  string     `json:"access_token,omitempty"`
	RefreshToken string     `json:"refresh_token,omitempty"`
	TokenType    string     `json:"token_type,omitempty"`
	Scopes       []string   `json:"scopes,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
	Actor        string     `json:"actor,omitempty"`
	GrantType    string     `json:"grant_type,omitempty"`
}

// Equal reports whether all persisted token fields match.
func (token TokenState) Equal(other TokenState) bool {
	return reflect.DeepEqual(token, other)
}

// String redacts secret material for fmt's %v, %+v, and %s verbs.
func (token TokenState) String() string {
	return fmt.Sprintf(
		"auth.TokenState{AccessToken:%s, RefreshToken:%s, TokenType:%s, Scopes:%v, "+
			"ExpiresAt:%v, Actor:%s, GrantType:%s}",
		presenceValue(token.AccessToken), presenceValue(token.RefreshToken), token.TokenType, token.Scopes,
		token.ExpiresAt, presenceValue(token.Actor), token.GrantType,
	)
}

// GoString redacts secret material for fmt's %#v verb.
func (token TokenState) GoString() string {
	return token.String()
}

// MarshalJSON redacts secret material so json.Marshal never emits raw values.
func (token TokenState) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		AccessToken  string     `json:"access_token"`
		RefreshToken string     `json:"refresh_token"`
		TokenType    string     `json:"token_type"`
		Scopes       []string   `json:"scopes"`
		ExpiresAt    *time.Time `json:"expires_at"`
		Actor        string     `json:"actor"`
		GrantType    string     `json:"grant_type"`
	}{
		AccessToken:  presenceValue(token.AccessToken),
		RefreshToken: presenceValue(token.RefreshToken),
		TokenType:    token.TokenType,
		Scopes:       token.Scopes,
		ExpiresAt:    token.ExpiresAt,
		Actor:        presenceValue(token.Actor),
		GrantType:    token.GrantType,
	})
}

// persistedTokenState carries the same fields as TokenState but drops its
// redaction methods, so marshaling it writes real values to disk.
type persistedTokenState TokenState

func presenceValue(value string) string {
	if value == "" {
		return "missing"
	}

	return "set"
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
	App      persistedAppConfig            `json:"app,omitempty"`
	Profiles map[string]persistedAppConfig `json:"profiles,omitempty"`
}

type tokenFile struct {
	Token    persistedTokenState            `json:"token,omitempty"`
	Profiles map[string]persistedTokenState `json:"profiles,omitempty"`
}

type osEnv struct{}

// Store reads and writes local auth state.
type Store struct {
	paths Paths
}

// TokenStateTransaction updates one profile's token while its state-file lock is held.
type TokenStateTransaction func(current TokenState) (TokenState, error)

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

	if err := withAuthFileLock(ctx, store.paths.AppConfigPath, "auth app config", "save auth state", func() error {
		return writeJSON(store.paths.AppConfigPath, appConfigFileFromState(state), "auth app config")
	}); err != nil {
		return err
	}
	return withAuthFileLock(ctx, store.paths.TokenPath, "auth token state", "save auth state", func() error {
		return writeJSON(store.paths.TokenPath, tokenFileFromState(state), "auth token state")
	})
}

// mutateAuthFile reloads one auth state file under its cross-process lock,
// applies mutate, and persists the result. It is the shared shape behind every
// profile-scoped save and clear. The mutate callback must not call Store
// persistence methods, because the lock is not reentrant.
func mutateAuthFile[File any](
	ctx context.Context,
	path string,
	label string,
	operation string,
	mutate func(*File),
) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%s context: %w", operation, err)
	}

	return withAuthFileLock(ctx, path, label, operation, func() error {
		state, err := readJSON[File](path, label)
		if err != nil {
			return err
		}
		mutate(&state)

		return writeJSON(path, state, label)
	})
}

// SaveAppConfig writes OAuth app configuration without touching token state.
func (store Store) SaveAppConfig(ctx context.Context, profile string, app AppConfig) error {
	return mutateAuthFile(
		ctx, store.paths.AppConfigPath, "auth app config", "save auth app config",
		func(appState *appConfigFile) {
			if profile == "" {
				appState.App = persistedAppConfig(app)

				return
			}
			if appState.Profiles == nil {
				appState.Profiles = map[string]persistedAppConfig{}
			}
			appState.Profiles[profile] = persistedAppConfig(app)
		},
	)
}

// SaveTokenState writes OAuth token state without touching app configuration.
func (store Store) SaveTokenState(ctx context.Context, profile string, token TokenState) error {
	return mutateAuthFile(
		ctx, store.paths.TokenPath, "auth token state", "save auth token state",
		func(tokenState *tokenFile) {
			if profile == "" {
				tokenState.Token = persistedTokenState(token)

				return
			}
			if tokenState.Profiles == nil {
				tokenState.Profiles = map[string]persistedTokenState{}
			}
			tokenState.Profiles[profile] = persistedTokenState(token)
		},
	)
}

// TransactTokenState reloads, updates, and persists one profile's token under one cross-process lock.
// The transaction callback must not call Store persistence methods because the lock is not reentrant.
func (store Store) TransactTokenState(
	ctx context.Context,
	profile string,
	transaction TokenStateTransaction,
) (TokenState, error) {
	if err := ctx.Err(); err != nil {
		return TokenState{}, fmt.Errorf("transact auth token state context: %w", err)
	}

	var result TokenState
	err := withAuthFileLock(ctx, store.paths.TokenPath, "auth token state", "transact auth token state", func() error {
		tokenState, err := readJSON[tokenFile](store.paths.TokenPath, "auth token state")
		if err != nil {
			return err
		}
		current := tokenStateForProfile(tokenState, profile)
		updated, err := transaction(current)
		if err != nil {
			return err
		}
		result = updated
		if current.Equal(updated) {
			return nil
		}
		setTokenStateForProfile(&tokenState, profile, updated)

		return writeJSON(store.paths.TokenPath, tokenState, "auth token state")
	})

	return result, err
}

// ClearTokenState removes saved OAuth token material while preserving app config.
func (store Store) ClearTokenState(ctx context.Context, profile string) error {
	return mutateAuthFile(
		ctx, store.paths.TokenPath, "auth token state", "clear auth token state",
		func(tokenState *tokenFile) {
			if profile == "" {
				tokenState.Token = persistedTokenState{}

				return
			}
			delete(tokenState.Profiles, profile)
		},
	)
}

// ClearAppConfig removes saved OAuth app configuration while preserving token state.
func (store Store) ClearAppConfig(ctx context.Context, profile string) error {
	return mutateAuthFile(
		ctx, store.paths.AppConfigPath, "auth app config", "clear auth app config",
		func(appState *appConfigFile) {
			if profile == "" {
				appState.App = persistedAppConfig{}

				return
			}
			delete(appState.Profiles, profile)
		},
	)
}

func withAuthFileLock(
	ctx context.Context,
	path string,
	label string,
	operation string,
	action func() error,
) (err error) {
	if path == "" {
		return action()
	}

	processLock := authFileLock(path)
	if err := processLock.acquire(ctx); err != nil {
		return fmt.Errorf("%s context: %w", operation, err)
	}
	defer processLock.release()

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create %s lock directory: %w", label, err)
	}
	if err := chmodIfSupported(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("secure %s lock directory: %w", label, err)
	}
	lockPath := path + ".lock"
	//nolint:gosec // Auth paths are resolved from user-specific config/state directories or explicit tests.
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("open %s lock %s: %w", label, lockPath, err)
	}
	defer func() {
		err = errors.Join(err, lockFile.Close())
	}()
	if err := chmodIfSupported(lockPath, 0o600); err != nil {
		return fmt.Errorf("secure %s lock %s: %w", label, lockPath, err)
	}
	if err := lockAuthFile(ctx, lockFile); err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil && errors.Is(err, ctxErr) {
			return fmt.Errorf("%s context: %w", operation, ctxErr)
		}
		return fmt.Errorf("acquire %s lock %s: %w", label, lockPath, err)
	}
	defer func() {
		err = errors.Join(err, unlockAuthFile(lockFile))
	}()
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%s context: %w", operation, err)
	}

	return action()
}

func tokenStateForProfile(state tokenFile, profile string) TokenState {
	if profile == "" {
		return TokenState(state.Token)
	}

	return TokenState(state.Profiles[profile])
}

func setTokenStateForProfile(state *tokenFile, profile string, token TokenState) {
	if profile == "" {
		state.Token = persistedTokenState(token)
		return
	}
	if state.Profiles == nil {
		state.Profiles = map[string]persistedTokenState{}
	}
	if tokenStateEmpty(token) {
		delete(state.Profiles, profile)
		return
	}

	state.Profiles[profile] = persistedTokenState(token)
}

func authFileLock(path string) *contextLock {
	cleanPath := filepath.Clean(path)
	authFileLocks.Lock()
	defer authFileLocks.Unlock()

	lock := authFileLocks.byPath[cleanPath]
	if lock == nil {
		lock = newContextLock()
		authFileLocks.byPath[cleanPath] = lock
	}

	return lock
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

	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("write %s %s: %w", label, path, err)
	}
	tmpPath := tmp.Name()
	defer func() {
		if err != nil {
			_ = os.Remove(tmpPath) //nolint:errcheck // temp cleanup is best effort after a failed write.
		}
	}()

	if err = chmodIfSupported(tmpPath, 0o600); err != nil {
		_ = tmp.Close() //nolint:errcheck // the original error is still returned.
		return fmt.Errorf("secure %s %s: %w", label, tmpPath, err)
	}
	if _, err = writeTempFile(tmp, data); err != nil {
		_ = tmp.Close() //nolint:errcheck // the original error is still returned.
		return fmt.Errorf("write %s %s: %w", label, path, err)
	}
	if err = syncTempFile(tmp); err != nil {
		_ = tmp.Close() //nolint:errcheck // the original error is still returned.
		return fmt.Errorf("write %s %s: %w", label, path, err)
	}
	if err = closeTempFile(tmp); err != nil {
		return fmt.Errorf("write %s %s: %w", label, path, err)
	}
	if err = os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("write %s %s: %w", label, path, err)
	}
	if err = chmodIfSupported(path, 0o600); err != nil {
		return fmt.Errorf("secure %s %s: %w", label, path, err)
	}

	return nil
}

func errorsIsNotExist(err error) bool {
	return err != nil && os.IsNotExist(err)
}

func appConfigFileFromState(state State) appConfigFile {
	profiles := map[string]persistedAppConfig{}
	for name, profile := range state.Profiles {
		if !appConfigEmpty(profile.App) {
			profiles[name] = persistedAppConfig(profile.App)
		}
	}
	if len(profiles) == 0 {
		profiles = nil
	}

	return appConfigFile{App: persistedAppConfig(state.App), Profiles: profiles}
}

func tokenFileFromState(state State) tokenFile {
	profiles := map[string]persistedTokenState{}
	for name, profile := range state.Profiles {
		if !tokenStateEmpty(profile.Token) {
			profiles[name] = persistedTokenState(profile.Token)
		}
	}
	if len(profiles) == 0 {
		profiles = nil
	}

	return tokenFile{Token: persistedTokenState(state.Token), Profiles: profiles}
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
		App:   AppConfig(appState.App),
		Token: TokenState(tokenState.Token),
	}
	for name, app := range appState.Profiles {
		if state.Profiles == nil {
			state.Profiles = map[string]ProfileState{}
		}
		profile := state.Profiles[name]
		profile.App = AppConfig(app)
		state.Profiles[name] = profile
	}
	for name, token := range tokenState.Profiles {
		if state.Profiles == nil {
			state.Profiles = map[string]ProfileState{}
		}
		profile := state.Profiles[name]
		profile.Token = TokenState(token)
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
