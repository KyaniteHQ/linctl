package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

type commandRuntime struct {
	config        config.Resolved
	graphqlClient graphql.Client
}

var buildCommandRuntime = newCommandRuntime

func newCommandRuntime(ctx context.Context, options *rootOptions) (commandRuntime, error) {
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
	if resolvedConfig.Token == "" {
		return commandRuntime{}, errors.New("missing Linear token: set LINCTL_TOKEN or LINEAR_API_KEY")
	}

	return commandRuntime{
		config: resolvedConfig,
		graphqlClient: client.NewTransport(client.TransportConfig{
			Token:   client.PersonalAPIToken(resolvedConfig.Token),
			Timeout: options.timeout,
		}),
	}, nil
}

func (runtime commandRuntime) resolveTarget(ctx context.Context) (client.ResolvedTarget, error) {
	return client.ResolveTarget(ctx, runtime.graphqlClient, runtime.config.Target)
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
		ProjectID: options.project,
	}
}
