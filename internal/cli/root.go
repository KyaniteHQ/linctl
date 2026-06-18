// Package cli owns the linctl command-line surface.
package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// BuildInfo contains version metadata injected by release builds.
type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

type rootOptions struct {
	timeout time.Duration
	json    bool
	profile string
	orgID   string
	team    string
	project string
}

// NewRootCommand builds the linctl root command.
func NewRootCommand(ctx context.Context, build BuildInfo) *cobra.Command {
	options := rootOptions{
		timeout: 30 * time.Second,
	}

	command := &cobra.Command{
		Use:           "linctl",
		Short:         "Schema-aligned Linear CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       build.versionText(),
	}
	command.SetVersionTemplate("linctl {{.Version}}\n")

	flags := command.PersistentFlags()
	flags.BoolVar(&options.json, "json", false, "emit JSON output")
	flags.StringVar(&options.profile, "profile", "", "config profile to load")
	flags.StringVar(&options.orgID, "org", "", "pinned Linear organization id")
	flags.StringVar(&options.team, "team", "", "pinned Linear team key or id")
	flags.StringVar(&options.project, "project", "", "pinned Linear project id")
	flags.DurationVar(&options.timeout, "timeout", options.timeout, "request timeout")

	addUsageCommand(command, &options)
	addTargetCommand(ctx, command, &options)
	addWhoamiCommand(ctx, command, &options)
	addIssueCommand(ctx, command, &options)
	addCurrentCommand(ctx, command, &options)
	addProjectCommand(ctx, command, &options)
	command.SetContext(ctx)

	return command
}

// Execute runs linctl with process stdio.
func Execute(ctx context.Context, build BuildInfo) error {
	command := NewRootCommand(ctx, build)
	command.SetIn(os.Stdin)
	command.SetOut(os.Stdout)
	command.SetErr(os.Stderr)

	return command.ExecuteContext(ctx)
}

func (build BuildInfo) versionText() string {
	version := defaultString(build.Version, "dev")
	commit := defaultString(build.Commit, "unknown")
	date := defaultString(build.Date, "unknown")

	return fmt.Sprintf("%s (commit %s, built %s)", version, commit, date)
}

func defaultString(value string, fallback string) string {
	if value != "" {
		return value
	}

	return fallback
}
