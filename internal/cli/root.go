// Package cli owns the linctl command-line surface.
package cli

import (
	"context"
	"fmt"
	"io"
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
	timeout     time.Duration
	json        bool
	compact     bool
	fields      string
	idOnly      bool
	quiet       bool
	failOnEmpty bool
	sortField   string
	sortOrder   string
	format      string
	configPath  string
	profile     string
	orgID       string
	team        string
	teamID      string
	project     string
	debug       bool
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
		PersistentPreRunE: func(command *cobra.Command, _ []string) error {
			return validateCommandFlags(command)
		},
	}
	command.SetVersionTemplate("linctl {{.Version}}\n")

	flags := command.PersistentFlags()
	flags.BoolVar(&options.json, "json", false, "emit JSON output")
	flags.BoolVar(&options.compact, "compact", false, "emit compact JSON when --json is set")
	flags.StringVar(&options.fields, "fields", "", "comma-separated JSON fields to emit")
	flags.BoolVar(&options.idOnly, "id-only", false, "emit only Linear ids")
	flags.BoolVar(&options.quiet, "quiet", false, "do not show output when a command succeeds")
	flags.BoolVar(&options.failOnEmpty, "fail-on-empty", false, "exit non-zero when a list result is empty")
	flags.StringVar(&options.sortField, "sort", "", "sort the list output by this JSON field")
	flags.StringVar(&options.sortOrder, "order", "asc", "sort order: asc or desc")
	flags.StringVar(&options.format, "format", "compact", "human output format: minimal, compact, or full")
	flags.StringVar(&options.configPath, "config", "", "path to the repo target config (default .linctl.toml)")
	flags.StringVar(&options.profile, "profile", "", "config profile to load")
	flags.StringVar(&options.orgID, "org", "", "pinned Linear organization id")
	flags.StringVar(&options.team, "team", "", "pinned Linear team key")
	flags.StringVar(&options.teamID, "team-id", "", "pinned Linear team id")
	flags.StringVar(&options.project, "project", "", "pinned Linear project id")
	flags.DurationVar(&options.timeout, "timeout", options.timeout, "total per-command deadline across retries")
	flags.BoolVar(&options.debug, "debug", false, "emit debug diagnostics to stderr")

	addCommands(ctx, command, &options)
	registerGlobalCompletions(ctx, command, &options)
	command.SetContext(ctx)

	return command
}

func validateCommandFlags(command *cobra.Command) error {
	if command.Flags().Lookup("limit") == nil {
		return nil
	}
	limit, err := command.Flags().GetInt("limit")
	if err != nil {
		return fmt.Errorf("read --limit: %w", err)
	}
	if limit <= 0 {
		return fmt.Errorf("invalid --limit %d: must be greater than 0", limit)
	}

	return nil
}

// commandRegistrar registers one top-level command group on the root command.
type commandRegistrar func(context.Context, *cobra.Command, *rootOptions)

// commandRegistrars is the whole top-level command surface, in registration
// order. It is a list rather than a sequence of calls so adding a group never
// grows a function body.
var commandRegistrars = []commandRegistrar{
	addUsageCommand,
	addAuthCommand,
	addTargetCommand,
	addInitCommand,
	addDoctorCommand,
	addWhoamiCommand,
	addApplicationCommand,
	addAgentActivityCommand,
	addAgentSkillCommand,
	addAgentSessionCommand,
	addExternalUserCommand,
	addAuditEntryCommand,
	addOrganizationCommand,
	addRateLimitCommand,
	addNotificationCommand,
	addReleasePipelineCommand,
	addReleaseStageCommand,
	addReleaseCommand,
	addExternalLinkCommand,
	addReleaseNoteCommand,
	addIssueToReleaseCommand,
	addIssueCommand,
	addIssueRelationCommand,
	addNextCommand,
	addFilesCommand,
	addCurrentCommand,
	addDoneCommand,
	addCommentCommand,
	addProjectCommand,
	addProjectUpdateReadCommand,
	addProjectMilestoneCommand,
	addProjectStatusCommand,
	addProjectLabelCommand,
	addProjectRelationCommand,
	addCycleCommand,
	addSprintCommand,
	addDocumentCommand,
	addLabelCommand,
	addTeamCommand,
	addTeamMembershipCommand,
	addUserCommand,
	addWorkflowStateCommand,
	addTimeScheduleCommand,
	addTriageResponsibilityCommand,
	addSLAConfigurationCommand,
	addSearchCommand,
	addSemanticSearchCommand,
	addTemplateCommand,
	addInitiativeCommand,
	addInitiativeLabelCommand,
	addInitiativeRelationCommand,
	addInitiativeToProjectCommand,
	addInitiativeUpdateCommand,
	addRoadmapCommand,
	addRoadmapToProjectCommand,
	addCustomViewCommand,
	addCustomerCommand,
	addCustomerNeedCommand,
	addCustomerStatusCommand,
	addCustomerTierCommand,
	addFavoriteCommand,
	addEmojiCommand,
	addAttachmentCommand,
}

func addCommands(ctx context.Context, command *cobra.Command, options *rootOptions) {
	for _, register := range commandRegistrars {
		register(ctx, command, options)
	}
}

// Execute runs linctl with process stdio.
func Execute(ctx context.Context, build BuildInfo) error {
	return execute(ctx, build, os.Stdin, os.Stdout, os.Stderr, nil)
}

// execute runs linctl with injectable streams and args so the failure path
// (the structured error envelope) is testable. On any error it emits exactly
// one error representation to stderr for machine consumers: the JSON error
// envelope, or a plain-text fallback line if the envelope itself fails to
// write. main only sets the exit code from the returned error.
func execute(
	ctx context.Context,
	build BuildInfo,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	args []string,
) error {
	command := NewRootCommand(ctx, build)
	command.SetIn(stdin)
	command.SetOut(stdout)
	command.SetErr(stderr)
	if args != nil {
		command.SetArgs(args)
	}

	err := command.ExecuteContext(ctx)
	if err != nil {
		if envelopeErr := writeErrorEnvelope(command.ErrOrStderr(), err); envelopeErr != nil {
			//nolint:errcheck // fallback path; the original error is still returned
			fmt.Fprintln(command.ErrOrStderr(), err)
		}
	}

	return err
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
