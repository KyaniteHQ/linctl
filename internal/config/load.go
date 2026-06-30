// Package config loads linctl configuration from files and profiles.
package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// ErrProfileNotFound marks an explicitly requested profile that does not exist.
var ErrProfileNotFound = errors.New("profile not found")

// Target is the pinned Linear write target.
type Target struct {
	OrgID     string `toml:"org_id"`
	TeamKey   string `toml:"team_key"`
	TeamID    string `toml:"team_id"`
	ProjectID string `toml:"project_id"`
}

// LoadRequest describes the config sources to load.
type LoadRequest struct {
	GlobalPath      string
	RepoPath        string
	ProfileOverride string
	TargetOverride  Target
}

// Resolved is the effective linctl configuration.
type Resolved struct {
	Profile string
	Target  Target
}

type fileConfig struct {
	Profile  string                   `toml:"profile"`
	Target   Target                   `toml:"target"`
	Profiles map[string]profileConfig `toml:"profiles"`
}

type profileConfig struct {
	Target Target `toml:"target"`
}

// Load resolves config with repo config overriding global config, then explicit overrides.
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

	return Resolved{
		Profile: profileName,
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

func readConfigFile(path string) (fileConfig, error) {
	if path == "" {
		return fileConfig{Profiles: map[string]profileConfig{}}, nil
	}

	//nolint:gosec // Config paths are explicit user/repo inputs; loading them is the feature.
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return fileConfig{Profiles: map[string]profileConfig{}}, nil
	}
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
		Target:   mergeTarget(base.Target, overlay.Target),
		Profiles: map[string]profileConfig{},
	}
	for name, profile := range base.Profiles {
		merged.Profiles[name] = profile
	}
	for name, profile := range overlay.Profiles {
		baseProfile := merged.Profiles[name]
		merged.Profiles[name] = profileConfig{
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

func firstNonEmpty(primary string, fallback string) string {
	if primary != "" {
		return primary
	}

	return fallback
}
