// Package config loads linctl configuration from files, profiles, and environment variables.
package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/KyaniteHQ/linctl/internal/auth"
	"github.com/pelletier/go-toml/v2"
)

// Env resolves environment variables by key.
type Env interface {
	Lookup(key string) (string, bool)
}

// ErrProfileNotFound marks an explicitly requested profile that does not exist.
var ErrProfileNotFound = errors.New("profile not found")

const (
	oauthAccessTokenEnv  = "LINCTL_OAUTH_ACCESS_TOKEN" //nolint:gosec // Environment variable name, not a secret.
	oauthClientIDEnv     = "LINCTL_OAUTH_CLIENT_ID"
	oauthClientSecretEnv = "LINCTL_OAUTH_CLIENT_SECRET" //nolint:gosec // Environment variable name, not a secret.
	oauthRedirectURIEnv  = "LINCTL_OAUTH_REDIRECT_URI"
	oauthScopesEnv       = "LINCTL_OAUTH_SCOPES"
)

// OAuthTokenSource is the resolved OAuth token material used by runtime code.
type OAuthTokenSource struct {
	AccessToken string
	App         auth.AppConfig
}

// Resolve returns the current OAuth access token from this source boundary.
func (source OAuthTokenSource) Resolve(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("resolve oauth token source: %w", err)
	}

	return source.AccessToken, nil
}

// Target is the pinned Linear write target.
type Target struct {
	OrgID     string `toml:"org_id"`
	TeamKey   string `toml:"team_key"`
	TeamID    string `toml:"team_id"`
	ProjectID string `toml:"project_id"`
}

// LoadRequest describes the config sources to load.
type LoadRequest struct {
	Env             Env
	GlobalPath      string
	RepoPath        string
	AuthStatePaths  auth.Paths
	ProfileOverride string
	TargetOverride  Target
}

// Resolved is the effective linctl configuration.
type Resolved struct {
	Profile string
	Auth    OAuthTokenSource
	Token   string
	Target  Target
}

type osEnv struct{}

type fileConfig struct {
	Profile  string                   `toml:"profile"`
	Token    string                   `toml:"token"`
	Target   Target                   `toml:"target"`
	Profiles map[string]profileConfig `toml:"profiles"`
}

type profileConfig struct {
	Token  string `toml:"token"`
	Target Target `toml:"target"`
}

// Load resolves config with repo config overriding global config, then flags and env.
func Load(ctx context.Context, request LoadRequest) (Resolved, error) {
	if err := ctx.Err(); err != nil {
		return Resolved{}, fmt.Errorf("load config context: %w", err)
	}

	globalConfig, err := readConfigFile(request.GlobalPath)
	if err != nil {
		return Resolved{}, err
	}
	repoConfig, err := readConfigFile(request.RepoPath)
	if err != nil {
		return Resolved{}, err
	}

	mergedConfig := mergeConfig(globalConfig, repoConfig)
	profileName := firstNonEmpty(request.ProfileOverride, mergedConfig.Profile)
	profile, err := resolveProfile(mergedConfig, profileName)
	if err != nil {
		return Resolved{}, err
	}
	target := mergeTarget(mergeTarget(mergedConfig.Target, profile.Target), request.TargetOverride)

	authSource, err := resolveOAuthTokenSource(ctx, request.Env, request.AuthStatePaths, profileName)
	if err != nil {
		return Resolved{}, err
	}

	return Resolved{
		Profile: profileName,
		Auth:    authSource,
		Token:   authSource.AccessToken,
		Target:  target,
	}, nil
}

func resolveProfile(config fileConfig, profileName string) (profileConfig, error) {
	if profileName == "" {
		return profileConfig{}, nil
	}
	profile, ok := config.Profiles[profileName]
	if !ok {
		return profileConfig{}, fmt.Errorf("%w: %s", ErrProfileNotFound, profileName)
	}

	return profile, nil
}

func (env osEnv) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func readConfigFile(path string) (fileConfig, error) {
	if path == "" {
		return fileConfig{Profiles: map[string]profileConfig{}}, nil
	}

	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return fileConfig{Profiles: map[string]profileConfig{}}, nil
	}
	if err != nil {
		return fileConfig{}, fmt.Errorf("read config %s: %w", path, err)
	}
	//nolint:gosec // Config paths are explicit user/repo inputs; loading them is the feature.
	data, err := os.ReadFile(path)
	if err != nil {
		return fileConfig{}, fmt.Errorf("read config %s: %w", path, err)
	}

	var config fileConfig
	if err := toml.Unmarshal(data, &config); err != nil {
		return fileConfig{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	if config.Profiles == nil {
		config.Profiles = map[string]profileConfig{}
	}
	return config, nil
}

func mergeConfig(base fileConfig, overlay fileConfig) fileConfig {
	merged := fileConfig{
		Profile:  firstNonEmpty(overlay.Profile, base.Profile),
		Token:    firstNonEmpty(overlay.Token, base.Token),
		Target:   mergeTarget(base.Target, overlay.Target),
		Profiles: map[string]profileConfig{},
	}
	for name, profile := range base.Profiles {
		merged.Profiles[name] = profile
	}
	for name, profile := range overlay.Profiles {
		baseProfile := merged.Profiles[name]
		merged.Profiles[name] = profileConfig{
			Token:  firstNonEmpty(profile.Token, baseProfile.Token),
			Target: mergeTarget(baseProfile.Target, profile.Target),
		}
	}

	return merged
}

func mergeTarget(base Target, overlay Target) Target {
	return Target{
		OrgID:     firstNonEmpty(overlay.OrgID, base.OrgID),
		TeamKey:   firstNonEmpty(overlay.TeamKey, base.TeamKey),
		TeamID:    firstNonEmpty(overlay.TeamID, base.TeamID),
		ProjectID: firstNonEmpty(overlay.ProjectID, base.ProjectID),
	}
}

func resolveOAuthTokenSource(
	ctx context.Context,
	env Env,
	authStatePaths auth.Paths,
	profileName string,
) (OAuthTokenSource, error) {
	activeEnv := env
	if activeEnv == nil {
		activeEnv = osEnv{}
	}
	if token, ok := activeEnv.Lookup(oauthAccessTokenEnv); ok && token != "" {
		if isPersonalAPIKeyShape(token) {
			return OAuthTokenSource{}, fmt.Errorf("%s contains personal API key-shaped material", oauthAccessTokenEnv)
		}

		return OAuthTokenSource{AccessToken: token, App: resolveOAuthAppConfig(activeEnv)}, nil
	}

	envApp := resolveOAuthAppConfig(activeEnv)
	if authStatePaths == (auth.Paths{}) {
		return OAuthTokenSource{App: envApp}, nil
	}
	state, err := auth.NewStore(authStatePaths).Load(ctx)
	if err != nil {
		return OAuthTokenSource{}, err
	}
	profileState := state.Profile(profileName)
	if !oauthAppConfigEmpty(envApp) {
		profileState.App = envApp
	}
	if isPersonalAPIKeyShape(profileState.Token.AccessToken) {
		return OAuthTokenSource{}, errors.New("local auth state contains personal API key-shaped material")
	}

	return OAuthTokenSource{
		AccessToken: profileState.Token.AccessToken,
		App:         profileState.App,
	}, nil
}

func resolveOAuthAppConfig(env Env) auth.AppConfig {
	clientID, _ := env.Lookup(oauthClientIDEnv)
	clientSecret, _ := env.Lookup(oauthClientSecretEnv)
	redirectURI, _ := env.Lookup(oauthRedirectURIEnv)
	scopeText, _ := env.Lookup(oauthScopesEnv)

	return auth.AppConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Scopes:       splitOAuthScopes(scopeText),
	}
}

func firstNonEmpty(primary string, fallback string) string {
	if primary != "" {
		return primary
	}

	return fallback
}

func isPersonalAPIKeyShape(value string) bool {
	return strings.HasPrefix(value, "lin_api_")
}

func splitOAuthScopes(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n'
	})
}

func oauthAppConfigEmpty(app auth.AppConfig) bool {
	return app.ClientID == "" &&
		app.ClientSecret == "" &&
		app.RedirectURI == "" &&
		len(app.Scopes) == 0
}
