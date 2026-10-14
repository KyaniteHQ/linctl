// Package config loads linctl configuration from files and profiles.
package config

import (
	"cmp"
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
	// Transitions is the credential's workflow transition allowlist. It is read
	// from the top-level [transitions] table, never from [target], so it carries
	// no toml tag of its own.
	Transitions Transitions `toml:"-" json:"-"`
	// States names workflow states that guarded commands move an issue to. It is
	// read from the top-level [states] table, never from [target].
	States States `toml:"-" json:"-"`
}

// States holds optional state names that replace a command's default pick.
// Close names the completed state `issue close` and `done` move an issue to;
// empty means the team's lowest-position completed state.
type States struct {
	Close string `toml:"close"`
}

// Transitions is a workflow transition allowlist keyed by the current state
// name. Each entry lists the state names an issue may move to from that state.
// Names are compared case-insensitively. An empty allowlist permits every
// transition; a non-empty one refuses any state change it does not list.
type Transitions map[string][]string

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
	Profile     string                   `toml:"profile"`
	Target      Target                   `toml:"target"`
	Transitions Transitions              `toml:"transitions"`
	States      States                   `toml:"states"`
	Profiles    map[string]profileConfig `toml:"profiles"`
}

type profileConfig struct {
	Target      Target      `toml:"target"`
	Transitions Transitions `toml:"transitions"`
	States      States      `toml:"states"`
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
	profileName := cmp.Or(request.ProfileOverride, mergedConfig.Profile)
	profile, err := resolveProfile(mergedConfig, profileName)
	if err != nil {
		return Resolved{}, err
	}
	target := mergeTarget(mergeTarget(mergedConfig.Target, profile.Target), request.TargetOverride)
	target.Transitions = mergeTransitions(mergedConfig.Transitions, profile.Transitions)
	target.States = mergeStates(mergedConfig.States, profile.States)
	override := request.TargetOverride
	if override.OrgID != "" || override.TeamKey != "" || override.TeamID != "" {
		// An explicit org or team override invalidates a pinned team id: the id
		// only survives when the override itself carries one.
		target.TeamID = override.TeamID
	}

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
		Profile:     cmp.Or(overlay.Profile, base.Profile),
		Target:      mergeTarget(base.Target, overlay.Target),
		Transitions: mergeTransitions(base.Transitions, overlay.Transitions),
		States:      mergeStates(base.States, overlay.States),
		Profiles:    map[string]profileConfig{},
	}
	for name, profile := range base.Profiles {
		merged.Profiles[name] = profile
	}
	for name, profile := range overlay.Profiles {
		baseProfile := merged.Profiles[name]
		merged.Profiles[name] = profileConfig{
			Target:      mergeTarget(baseProfile.Target, profile.Target),
			Transitions: mergeTransitions(baseProfile.Transitions, profile.Transitions),
			States:      mergeStates(baseProfile.States, profile.States),
		}
	}

	return merged
}

func mergeTarget(base Target, overlay Target) Target {
	return Target{
		OrgID:       cmp.Or(overlay.OrgID, base.OrgID),
		TeamKey:     cmp.Or(overlay.TeamKey, base.TeamKey),
		TeamID:      cmp.Or(overlay.TeamID, base.TeamID),
		ProjectID:   cmp.Or(overlay.ProjectID, base.ProjectID),
		Transitions: mergeTransitions(base.Transitions, overlay.Transitions),
		States:      mergeStates(base.States, overlay.States),
	}
}

// mergeTransitions replaces the whole allowlist when the overlay sets one: a
// narrower repo or profile allowlist must not be widened by entries from a
// broader global one.
func mergeTransitions(base Transitions, overlay Transitions) Transitions {
	if len(overlay) > 0 {
		return overlay
	}

	return base
}

// mergeStates lets an overlay that names a state replace the base one; an
// overlay that names none keeps the base.
func mergeStates(base States, overlay States) States {
	return States{Close: cmp.Or(overlay.Close, base.Close)}
}
