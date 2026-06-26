package cli

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

type commandRuntime struct {
	config        config.Resolved
	fileClient    httpDoer
	graphqlClient graphql.Client
	logger        *slog.Logger
}

var buildCommandRuntime = newCommandRuntime

func newCommandRuntime(ctx context.Context, options *rootOptions) (commandRuntime, error) {
	logger := newDiagnosticLogger(options.debug, os.Getenv("LINCTL_DEBUG_JSON") == "1", os.Stderr)
	override := targetOverride(options)
	resolvedConfig, err := config.Load(ctx, config.LoadRequest{
		GlobalPath:      defaultGlobalConfigPath(),
		RepoPath:        ".linctl.toml",
		ProfileOverride: options.profile,
		TargetOverride:  override,
	})
	if err != nil {
		return commandRuntime{}, err
	}
	applyTargetOverrideFlagSemantics(&resolvedConfig, options)
	if resolvedConfig.Token == "" {
		return commandRuntime{}, errors.New("missing Linear token: set LINCTL_TOKEN or LINEAR_API_KEY")
	}

	logger.Debug(
		"runtime ready",
		"profile", resolvedConfig.Profile,
		"org", resolvedConfig.Target.OrgID,
		"team_key", resolvedConfig.Target.TeamKey,
		"team_id", resolvedConfig.Target.TeamID,
		"project", resolvedConfig.Target.ProjectID,
		"timeout", options.timeout.String(),
	)

	return commandRuntime{
		config:     resolvedConfig,
		fileClient: &http.Client{Timeout: options.timeout},
		logger:     logger,
		graphqlClient: client.NewTransport(client.TransportConfig{
			Token:            client.PersonalAPIToken(resolvedConfig.Token),
			Timeout:          options.timeout,
			DiagnosticWriter: newTransportDiagnosticWriter(logger, options.debug),
		}),
	}, nil
}

func (runtime commandRuntime) resolveTarget(ctx context.Context) (client.ResolvedTarget, error) {
	target, err := client.ResolveTarget(ctx, runtime.graphqlClient, runtime.config.Target)
	logTargetResolution(runtime.log(), target, err)

	return target, err
}

func defaultGlobalConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".config", "linctl", "config.toml")
}

func targetOverride(options *rootOptions) config.Target {
	return config.Target{
		OrgID:     options.orgID,
		TeamKey:   options.team,
		TeamID:    options.teamID,
		ProjectID: options.project,
	}
}

func applyTargetOverrideFlagSemantics(resolved *config.Resolved, options *rootOptions) {
	if options.orgID == "" && options.team == "" && options.teamID == "" {
		return
	}

	resolved.Target.TeamID = options.teamID
}
