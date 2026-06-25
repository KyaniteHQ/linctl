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
	profile     string
	orgID       string
	team        string
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
	}
	command.SetVersionTemplate("linctl {{.Version}}\n")

	flags := command.PersistentFlags()
	flags.BoolVar(&options.json, "json", false, "emit JSON output")
	flags.BoolVar(&options.compact, "compact", false, "emit compact JSON when --json is set")
	flags.StringVar(&options.fields, "fields", "", "comma-separated JSON fields to emit")
	flags.BoolVar(&options.idOnly, "id-only", false, "emit only Linear ids")
	flags.BoolVar(&options.quiet, "quiet", false, "suppress successful output")
	flags.BoolVar(&options.failOnEmpty, "fail-on-empty", false, "exit non-zero when a list result is empty")
	flags.StringVar(&options.sortField, "sort", "", "JSON field to sort list output by")
	flags.StringVar(&options.sortOrder, "order", "asc", "sort order: asc or desc")
	flags.StringVar(&options.format, "format", "compact", "human output format: minimal, compact, or full")
	flags.StringVar(&options.profile, "profile", "", "config profile to load")
	flags.StringVar(&options.orgID, "org", "", "pinned Linear organization id")
	flags.StringVar(&options.team, "team", "", "pinned Linear team key or id")
	flags.StringVar(&options.project, "project", "", "pinned Linear project id")
	flags.DurationVar(&options.timeout, "timeout", options.timeout, "request timeout")
	flags.BoolVar(&options.debug, "debug", false, "emit debug diagnostics to stderr")

	addCommands(ctx, command, &options)
	registerGlobalCompletions(ctx, command, &options)
	command.SetContext(ctx)

	return command
}

func addCommands(ctx context.Context, command *cobra.Command, options *rootOptions) {
	addUsageCommand(command, options)
	addTargetCommand(ctx, command, options)
	addDoctorCommand(ctx, command, options)
	addWhoamiCommand(ctx, command, options)
	addApplicationCommand(ctx, command, options)
	addAgentActivityCommand(ctx, command, options)
	addAgentSkillCommand(ctx, command, options)
	addAgentSessionCommand(ctx, command, options)
	addExternalUserCommand(ctx, command, options)
	addAuditEntryCommand(ctx, command, options)
	addOrganizationCommand(ctx, command, options)
	addRateLimitCommand(ctx, command, options)
	addNotificationCommand(ctx, command, options)
	addReleasePipelineCommand(ctx, command, options)
	addReleaseStageCommand(ctx, command, options)
	addReleaseCommand(ctx, command, options)
	addExternalLinkCommand(ctx, command, options)
	addReleaseNoteCommand(ctx, command, options)
	addIssueToReleaseCommand(ctx, command, options)
	addIssueCommand(ctx, command, options)
	addIssueRelationCommand(ctx, command, options)
	addNextCommand(ctx, command, options)
	addFilesCommand(ctx, command, options)
	addCurrentCommand(ctx, command, options)
	addDoneCommand(ctx, command, options)
	addCommentCommand(ctx, command, options)
	addProjectCommand(ctx, command, options)
	addProjectUpdateReadCommand(ctx, command, options)
	addProjectMilestoneCommand(ctx, command, options)
	addProjectStatusCommand(ctx, command, options)
	addProjectLabelCommand(ctx, command, options)
	addProjectRelationCommand(ctx, command, options)
	addCycleCommand(ctx, command, options)
	addSprintCommand(ctx, command, options)
	addDocumentCommand(ctx, command, options)
	addLabelCommand(ctx, command, options)
	addTeamCommand(ctx, command, options)
	addTeamMembershipCommand(ctx, command, options)
	addUserCommand(ctx, command, options)
	addWorkflowStateCommand(ctx, command, options)
	addTimeScheduleCommand(ctx, command, options)
	addTriageResponsibilityCommand(ctx, command, options)
	addSLAConfigurationCommand(ctx, command, options)
	addSearchCommand(ctx, command, options)
	addSemanticSearchCommand(ctx, command, options)
	addTemplateCommand(ctx, command, options)
	addInitiativeCommand(ctx, command, options)
	addInitiativeRelationCommand(ctx, command, options)
	addInitiativeToProjectCommand(ctx, command, options)
	addInitiativeUpdateCommand(ctx, command, options)
	addRoadmapCommand(ctx, command, options)
	addRoadmapToProjectCommand(ctx, command, options)
	addCustomViewCommand(ctx, command, options)
	addCustomerCommand(ctx, command, options)
	addCustomerNeedCommand(ctx, command, options)
	addCustomerStatusCommand(ctx, command, options)
	addCustomerTierCommand(ctx, command, options)
	addFavoriteCommand(ctx, command, options)
	addEmojiCommand(ctx, command, options)
	addAttachmentCommand(ctx, command, options)
}

// Execute runs linctl with process stdio.
func Execute(ctx context.Context, build BuildInfo) error {
	return execute(ctx, build, os.Stdin, os.Stdout, os.Stderr, nil)
}

// execute runs linctl with injectable streams and args so the failure path
// (the structured error envelope) is testable. On any error it emits a single
// JSON error envelope to stderr for machine consumers and still returns the
// error so main can print the human-readable line and set the exit code.
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
		//nolint:errcheck // best-effort structured error on stderr; the error is still returned
		_ = writeErrorEnvelope(command.ErrOrStderr(), err)
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
