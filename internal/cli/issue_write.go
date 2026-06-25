package cli

import (
	"context"
	"fmt"
	"io"
	"os"

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
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			var estimate *int
			if command.Flags().Changed("estimate") {
				estimate = &flags.estimate
			}

			return runIssueCreate(ctx, command, options, issueAdapterFor(runtime), request, flags, estimate)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "issue title")
	command.Flags().StringVar(&request.Description, "description", "", "issue description")
	command.Flags().StringVar(&flags.descriptionFile, "description-file", "", "read issue description from file")
	command.Flags().StringVar(
		&flags.templateID, "template", "",
		"apply a Linear template by id for title/description defaults",
	)
	command.Flags().StringArrayVar(&flags.sections, "section", nil, "fill a markdown section: NAME=VALUE (repeatable)")
	command.Flags().BoolVar(&flags.dryRun, "dry-run", false, "render the assembled issue without creating it")
	command.Flags().StringVar(&flags.state, "state", "", "set workflow state type (e.g. started, completed)")
	command.Flags().StringVar(&flags.status, "status", "", "alias for --state")
	command.Flags().StringVar(&flags.priority, "priority", "", "set priority (urgent/high/medium/low/none or 0-4)")
	command.Flags().StringVar(&request.AssigneeID, "assignee", "", "assign the issue to a user id")
	command.Flags().StringArrayVar(&request.LabelIDs, "label", nil, "attach a label by id (repeatable)")
	command.Flags().StringVar(&request.DueDate, "due-date", "", "set the due date (YYYY-MM-DD)")
	command.Flags().IntVar(&flags.estimate, "estimate", 0, "set the estimate (validated against team config)")
	command.Flags().StringVar(&request.ParentID, "parent", "", "create as a sub-issue of a parent issue id")
	registerStateCompletion(ctx, command, options)
	root.AddCommand(command)
}

func runIssueCreate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	creator issueCreator,
	request client.IssueCreateRequest,
	flags issueCreateFlags,
	estimate *int,
) error {
	if err := resolveFileFlag(&request.Description, flags.descriptionFile, "description"); err != nil {
		return err
	}
	if err := applyIssueTemplate(ctx, creator, &request, flags.templateID); err != nil {
		return err
	}
	if err := applyIssueSections(&request, flags.sections); err != nil {
		return err
	}
	stateType, normalizedPriority, normErr := applyIssueWriteNormalization(
		command, flags.state, flags.status, flags.priority,
	)
	if normErr != nil {
		return normErr
	}
	request.StateType = stateType
	request.Priority = normalizedPriority
	request.Estimate = estimate
	if flags.dryRun {
		return writeIssueDraft(command, options, request)
	}
	issue, err := creator.CreateIssue(ctx, request)
	if err != nil {
		return err
	}

	return writeIssue(command, options, issue)
}

func addIssueUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueUpdateRequest{}
	descriptionFile := ""
	appendFile := ""
	state := ""
	status := ""
	priority := ""
	estimate := 0
	command := &cobra.Command{
		Use:   "update ISSUE_ID",
		Short: "Update an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]
			if err := resolveFileFlag(&request.Description, descriptionFile, "description"); err != nil {
				return err
			}
			if err := resolveFileFlag(&request.Append, appendFile, "append"); err != nil {
				return err
			}
			stateType, normalizedPriority, normErr := applyIssueWriteNormalization(command, state, status, priority)
			if normErr != nil {
				return normErr
			}
			request.StateType = stateType
			request.Priority = normalizedPriority
			if command.Flags().Changed("estimate") {
				request.Estimate = &estimate
			}
			issue, err := client.UpdateIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "new issue title")
	command.Flags().StringVar(&request.Description, "description", "", "new issue description")
	command.Flags().StringVar(&descriptionFile, "description-file", "", "read new issue description from file")
	command.Flags().StringVar(&request.Append, "append", "", "text to append to the issue description")
	command.Flags().StringVar(&appendFile, "append-file", "", "read text to append from file")
	command.Flags().StringVar(&state, "state", "", "set workflow state type (e.g. started, completed)")
	command.Flags().StringVar(&status, "status", "", "alias for --state")
	command.Flags().StringVar(&priority, "priority", "", "set priority (urgent/high/medium/low/none or 0-4)")
	command.Flags().StringVar(&request.AssigneeID, "assignee", "", "reassign the issue to a user id")
	command.Flags().StringArrayVar(&request.LabelIDs, "label", nil, "set labels by id (repeatable, replaces existing)")
	command.Flags().StringVar(&request.DueDate, "due-date", "", "set the due date (YYYY-MM-DD)")
	command.Flags().BoolVar(&request.ClearDueDate, "clear-due-date", false, "clear the due date")
	command.Flags().IntVar(&estimate, "estimate", 0, "set the estimate (validated against team config)")
	command.Flags().BoolVar(&request.ClearEstimate, "clear-estimate", false, "clear the estimate")
	registerStateCompletion(ctx, command, options)
	root.AddCommand(command)
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
	stateType, err = normalizeAndNote(command, "state", mergedStateFlag(state, status), normalizedStateType)
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
	root.AddCommand(&cobra.Command{
		Use:   "start ISSUE_ID",
		Short: "Assign and start an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.StartIssue(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}

func addIssueCommentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCommentRequest{}
	bodyFile := ""
	command := &cobra.Command{
		Use:   "comment ISSUE_ID",
		Short: "Comment on an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			request.ID = args[0]
			return runIssueBodyWriteCommand(ctx, command, options, request, bodyFile)
		},
	}
	command.Flags().StringVar(&request.Body, "body", "", "comment body")
	command.Flags().StringVar(&bodyFile, "body-file", "", "read comment body from file")
	root.AddCommand(command)
}

func runIssueBodyWriteCommand(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	request client.IssueCommentRequest,
	bodyFile string,
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	if err := resolveBodyFlag(command, &request.Body); err != nil {
		return err
	}
	if err := resolveFileFlag(&request.Body, bodyFile, "body"); err != nil {
		return err
	}
	comment, err := client.CommentOnIssue(ctx, runtime.graphqlClient, runtime.config.Target, request)
	if err != nil {
		return err
	}
	if options.json {
		return writeJSONValue(command, options, comment)
	}

	return render.WriteLine(command.OutOrStdout(), "comment %s on %s", comment.ID, comment.Issue.Identifier)
}

func addIssueReplyCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.IssueCommentRequest{}
	bodyFile := ""
	command := &cobra.Command{
		Use:   "reply ISSUE_ID COMMENT_ID",
		Short: "Reply to an issue comment after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			request.ID = args[0]
			request.ParentID = args[1]
			return runIssueBodyWriteCommand(ctx, command, options, request, bodyFile)
		},
	}
	command.Flags().StringVar(&request.Body, "body", "", "reply body")
	command.Flags().StringVar(&bodyFile, "body-file", "", "read reply body from file")
	root.AddCommand(command)
}

func resolveBodyFlag(command *cobra.Command, body *string) error {
	if *body != "-" {
		return nil
	}
	data, err := io.ReadAll(command.InOrStdin())
	if err != nil {
		return fmt.Errorf("read body from stdin: %w", err)
	}
	*body = string(data)

	return nil
}

func resolveFileFlag(value *string, path string, label string) error {
	if path == "" {
		return nil
	}
	if *value != "" {
		return fmt.Errorf("%s and %s-file are mutually exclusive", label, label)
	}

	//nolint:gosec // The path is an explicit CLI input for reading issue text.
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s from file %s: %w", label, path, err)
	}
	*value = string(data)

	return nil
}

func addIssueCloseCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "close ISSUE_ID",
		Short: "Move an issue to the completed workflow state",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runIssueClose(ctx, command, options, issueAdapterFor(runtime), args[0])
		},
	})
}

func runIssueClose(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	closer issueCloser,
	issueID string,
) error {
	issue, err := closer.CloseIssue(ctx, issueID)
	if err != nil {
		return err
	}

	return writeIssue(command, options, issue)
}

func writeIssue(command *cobra.Command, options *rootOptions, issue client.IssueSummary) error {
	if wrote, err := writeIDOnly(command, options, issue.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, issue)
	}

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
	command := &cobra.Command{
		Use:   "link URL ISSUE_ID",
		Short: "Attach a URL to an issue after pinned-target comparison",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.URL = args[0]
			request.IssueID = args[1]
			attachment, err := client.LinkIssueAttachment(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeAttachmentLink(command, options, attachment)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "attachment title")
	command.Flags().StringVar(&request.Subtitle, "subtitle", "", "attachment subtitle")
	root.AddCommand(command)
}

func writeAttachmentLink(command *cobra.Command, options *rootOptions, attachment client.AttachmentSummary) error {
	if wrote, err := writeIDOnly(command, options, attachment.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, attachment)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s", attachment.ID, attachment.URL)
}
