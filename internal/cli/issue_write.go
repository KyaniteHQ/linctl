package cli

import (
	"cmp"
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

// issueCreateFlags collects the non-request inputs of the issue create command.
type issueCreateFlags struct {
	descriptionFile string
	templateID      string
	sections        []string
	state           string
	status          string
	priority        string
	dryRun          bool
	estimate        int
}

func addIssueCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCreateRequest{}
	flags := issueCreateFlags{}
	command := &cobra.Command{
		Use:   "create",
		Short: "Create an issue in the pinned target",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			var runtime commandRuntime
			var templateReader issueTemplateReader
			if issueCreateRequiresRuntime(flags) {
				var err error
				runtime, err = buildCommandRuntime(ctx, options)
				if err != nil {
					return err
				}
				templateReader = issueTemplateClient{graphqlClient: runtime.graphqlClient}
			}
			var estimate *int
			if command.Flags().Changed("estimate") {
				estimate = &flags.estimate
			}
			assembled, err := assembleIssueCreate(ctx, command, templateReader, request, flags, estimate)
			if err != nil {
				return err
			}
			if flags.dryRun {
				return writeIssueDraft(command, options, assembled)
			}
			issue, err := client.CreateIssue(ctx, runtime.graphqlClient, runtime.config.Target, assembled)
			if err != nil {
				return err
			}
			if err := noteCreatedProject(command, assembled.ParentID, issue); err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "issue title")
	command.Flags().StringVar(&request.Description, "description", "", "issue description")
	command.Flags().StringVar(&flags.descriptionFile, "description-file", "", "read issue description from file")
	command.Flags().StringVar(
		&flags.templateID, "template", "",
		"apply a Linear template by id to set the default title and description",
	)
	command.Flags().StringArrayVar(
		&flags.sections, "section", nil,
		"fill a markdown section with NAME=VALUE, and repeat the flag for more sections",
	)
	command.Flags().BoolVar(&flags.dryRun, "dry-run", false, "show the assembled issue, and do not create it")
	command.Flags().StringVar(
		&flags.state, "state", "",
		"set the workflow state type, for example started or completed",
	)
	command.Flags().StringVar(&flags.status, "status", "", "alias for --state")
	command.Flags().StringVar(
		&flags.priority, "priority", "",
		"set the priority to urgent, high, medium, low, none, or a number from 0 to 4",
	)
	command.Flags().StringVar(&request.AssigneeID, "assignee", "", "assign the issue to a user id")
	command.Flags().StringArrayVar(
		&request.LabelIDs, "label", nil,
		"attach a label by id, and repeat the flag for more labels",
	)
	command.Flags().StringVar(&request.DueDate, "due-date", "", "set the due date in YYYY-MM-DD format")
	command.Flags().IntVar(
		&flags.estimate, "estimate", 0,
		"set the estimate, which linctl validates against the team configuration",
	)
	command.Flags().StringVar(
		&request.ParentID, "parent", "",
		"create the issue as a sub-issue of this parent issue id",
	)
	command.Flags().StringVar(
		&request.ProjectMilestoneID, "milestone", "",
		"assign to a project milestone id, which needs a pinned project",
	)
	registerStateCompletion(ctx, command, options)
	addWriteCommand(root, WriteEffectGuarded, command)
}

// noteCreatedProject reports on stderr which project an unparented create landed
// in. A create with a parent already proves its project through the parent guard,
// so it needs no note. The note is informational: the pinned project is the
// intended destination, and only the caller knows whether the pin was the one
// they meant.
func noteCreatedProject(command *cobra.Command, parentID string, issue client.IssueSummary) error {
	if parentID != "" {
		return nil
	}
	if issue.Project == "" {
		return writeNote(command, "created in team %s with no project", issue.Team)
	}

	return writeNote(command, "created in project %q", issue.Project)
}

func issueCreateRequiresRuntime(flags issueCreateFlags) bool {
	return !flags.dryRun || flags.templateID != ""
}

func assembleIssueCreate(
	ctx context.Context,
	command *cobra.Command,
	templateReader issueTemplateReader,
	request client.IssueCreateRequest,
	flags issueCreateFlags,
	estimate *int,
) (client.IssueCreateRequest, error) {
	if err := resolveFileFlag(command, &request.Description, flags.descriptionFile, "description"); err != nil {
		return client.IssueCreateRequest{}, err
	}
	if err := applyIssueTemplate(ctx, templateReader, &request, flags.templateID); err != nil {
		return client.IssueCreateRequest{}, err
	}
	if err := applyIssueSections(&request, flags.sections); err != nil {
		return client.IssueCreateRequest{}, err
	}
	stateType, normalizedPriority, normErr := applyIssueWriteNormalization(
		command, flags.state, flags.status, flags.priority,
	)
	if normErr != nil {
		return client.IssueCreateRequest{}, normErr
	}
	request.StateType = stateType
	request.Priority = normalizedPriority
	request.Estimate = estimate
	return request, nil
}

// issueUpdateFlags collects the non-request inputs of the issue update command.
type issueUpdateFlags struct {
	descriptionFile string
	appendFile      string
	state           string
	status          string
	priority        string
}

func addIssueUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueUpdateRequest{}
	flags := issueUpdateFlags{}
	estimate := 0
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueSummary]{
		Use:   "update ISSUE_ID",
		Short: "Update an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Title, "title", "", "new issue title")
			command.Flags().StringVar(&request.Description, "description", "", "new issue description")
			command.Flags().StringVar(
				&flags.descriptionFile, "description-file", "", "read new issue description from file",
			)
			command.Flags().StringVar(&request.Append, "append", "", "text to append to the issue description")
			command.Flags().StringVar(&flags.appendFile, "append-file", "", "read text to append from file")
			command.Flags().StringVar(
				&flags.state, "state", "",
				"set the workflow state type, for example started or completed",
			)
			command.Flags().StringVar(&flags.status, "status", "", "alias for --state")
			command.Flags().StringVar(
				&flags.priority, "priority", "",
				"set the priority to urgent, high, medium, low, none, or a number from 0 to 4",
			)
			command.Flags().StringVar(&request.AssigneeID, "assignee", "", "reassign the issue to a user id")
			command.Flags().StringArrayVar(
				&request.LabelIDs, "label", nil,
				"replace the labels with these ids, and repeat the flag for more labels",
			)
			command.Flags().StringVar(&request.DueDate, "due-date", "", "set the due date in YYYY-MM-DD format")
			command.Flags().BoolVar(&request.ClearDueDate, "clear-due-date", false, "clear the due date")
			command.Flags().IntVar(
				&estimate, "estimate", 0,
				"set the estimate, which linctl validates against the team configuration",
			)
			command.Flags().BoolVar(&request.ClearEstimate, "clear-estimate", false, "clear the estimate")
			command.Flags().StringVar(
				&request.ProjectMilestoneID, "milestone", "",
				"assign to a project milestone id, which needs a pinned project",
			)
			command.Flags().BoolVar(&request.ClearMilestone, "clear-milestone", false, "clear the milestone")
			registerStateCompletion(ctx, command, options)
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueSummary, error) {
			request.ID = args[0]
			var resolvedEstimate *int
			if command.Flags().Changed("estimate") {
				resolvedEstimate = &estimate
			}

			assembled, err := assembleIssueUpdate(command, request, flags, resolvedEstimate)
			if err != nil {
				return client.IssueSummary{}, err
			}

			return client.UpdateIssue(ctx, runtime.graphqlClient, runtime.config.Target, assembled)
		},
		Write: writeIssue,
	})
}

func addIssueMoveTeamCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueMoveTeamRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueSummary]{
		Use:   "move-team ISSUE_ID",
		Short: "Move an issue from the pinned team to another team in the organization",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.TeamKey, "to-team", "", "destination team key")
			command.Flags().StringVar(&request.TeamID, "to-team-id", "", "destination team id")
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueSummary, error) {
			request.IssueID = args[0]

			return client.MoveIssueTeam(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeIssue,
	})
}

func addIssueMoveProjectCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueMoveProjectRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueSummary]{
		Use:   "move-project ISSUE_ID",
		Short: "Move an issue from the pinned project to another project on the team",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.ProjectID, "to-project-id", "", "destination project id")
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueSummary, error) {
			request.IssueID = args[0]

			return client.MoveIssueProject(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeIssue,
	})
}

func assembleIssueUpdate(
	command *cobra.Command,
	request client.IssueUpdateRequest,
	flags issueUpdateFlags,
	estimate *int,
) (client.IssueUpdateRequest, error) {
	if err := resolveFileFlag(command, &request.Description, flags.descriptionFile, "description"); err != nil {
		return client.IssueUpdateRequest{}, err
	}
	if err := resolveFileFlag(command, &request.Append, flags.appendFile, "append"); err != nil {
		return client.IssueUpdateRequest{}, err
	}
	stateType, normalizedPriority, normErr := applyIssueWriteNormalization(
		command, flags.state, flags.status, flags.priority,
	)
	if normErr != nil {
		return client.IssueUpdateRequest{}, normErr
	}
	request.StateType = stateType
	request.Priority = normalizedPriority
	request.Estimate = estimate
	return request, nil
}

// applyIssueWriteNormalization merges the --state/--status alias pair and
// normalizes both the state type and the priority string. It emits a note to
// stderr when an alias was expanded to its canonical form.
func applyIssueWriteNormalization(
	command *cobra.Command,
	state string,
	status string,
	priority string,
) (stateType string, normalizedPriority string, err error) {
	stateType, err = normalizeAndNote(command, "state", cmp.Or(state, status), normalizedStateType)
	if err != nil {
		return "", "", err
	}
	normalizedPriority, err = normalizeAndNote(command, "priority", priority, normalizedPriorityValue)
	if err != nil {
		return "", "", err
	}

	return stateType, normalizedPriority, nil
}

func addIssueStartCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueSummary]{
		Use:   "start ISSUE_ID",
		Short: "Assign and start an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueSummary, error) {
			return client.StartIssue(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeIssue,
	})
}

func addIssueCommentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCommentRequest{}
	bodyFile := ""
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueCommentResult]{
		Use:   "comment ISSUE_ID",
		Short: "Comment on an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Body, "body", "", "comment body")
			command.Flags().StringVar(&bodyFile, "body-file", "", "read comment body from file")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueCommentResult, error) {
			request.ID = args[0]

			return runIssueBodyWrite(ctx, command, runtime, request, bodyFile)
		},
		Write: writeIssueComment,
	})
}

func runIssueBodyWrite(
	ctx context.Context,
	command *cobra.Command,
	runtime commandRuntime,
	request client.IssueCommentRequest,
	bodyFile string,
) (client.IssueCommentResult, error) {
	if err := resolveBodyOrFileFlag(command, &request.Body, bodyFile, "body"); err != nil {
		return client.IssueCommentResult{}, err
	}

	return client.CommentOnIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
}

func writeIssueComment(command *cobra.Command, options *rootOptions, comment client.IssueCommentResult) error {
	return writeItemLine(
		command, options, comment, comment.ID,
		"comment %s on %s", comment.ID, comment.Issue.Identifier,
	)
}

func addIssueReplyCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCommentRequest{}
	bodyFile := ""
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueCommentResult]{
		Use:   "reply ISSUE_ID COMMENT_ID",
		Short: "Reply to an issue comment after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Body, "body", "", "reply body")
			command.Flags().StringVar(&bodyFile, "body-file", "", "read reply body from file")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueCommentResult, error) {
			request.ID = args[0]
			request.ParentID = args[1]

			return runIssueBodyWrite(ctx, command, runtime, request, bodyFile)
		},
		Write: writeIssueComment,
	})
}

func addIssueCloseCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.IssueSummary]{
		Use:   "close ISSUE_ID",
		Short: "Move an issue to the completed workflow state",
		Args:  cobra.ExactArgs(1),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.IssueSummary, error) {
			return client.CloseIssue(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeIssue,
	})
}

func writeIssue(command *cobra.Command, options *rootOptions, issue client.IssueSummary) error {
	return writeItem(command, options, issue, issue.ID, issueHumanLine)
}

func issueHumanLine(command *cobra.Command, options *rootOptions, issue client.IssueSummary) error {
	format, err := normalizedHumanFormat(options)
	if err != nil {
		return err
	}
	if format == "minimal" {
		return render.WriteLine(command.OutOrStdout(), "%s", issue.Identifier)
	}
	if format == "full" {
		return render.WriteLine(
			command.OutOrStdout(),
			"%s %s [%s] project=%s url=%s",
			issue.Identifier,
			issue.Title,
			issue.State,
			emptyDash(issue.Project),
			issue.URL,
		)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", issue.Identifier, issue.Title, issue.State)
}

func emptyDash(value string) string {
	if value == "" {
		return "-"
	}

	return value
}

func addIssueLinkCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.AttachmentLinkRequest{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.AttachmentSummary]{
		Use:   "link URL ISSUE_ID",
		Short: "Attach a URL to an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Title, "title", "", "attachment title")
			command.Flags().StringVar(&request.Subtitle, "subtitle", "", "attachment subtitle")
		},
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.AttachmentSummary, error) {
			request.URL = args[0]
			request.IssueID = args[1]

			return client.LinkIssueAttachment(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeAttachmentLink,
	})
}

func writeAttachmentLink(command *cobra.Command, options *rootOptions, attachment client.AttachmentSummary) error {
	return writeItem(command, options, attachment, attachment.ID,
		func(command *cobra.Command, _ *rootOptions, attachment client.AttachmentSummary) error {
			return render.WriteLine(command.OutOrStdout(), "%s %s", attachment.ID, attachment.URL)
		})
}
