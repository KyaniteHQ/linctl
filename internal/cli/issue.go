//nolint:dupl // Issue child read commands intentionally share the same list-command shape.
package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	issueCommand := &cobra.Command{
		Use:   "issue",
		Short: "Read and write Linear issues",
	}
	addIssueListCommand(ctx, issueCommand, options)
	addIssueSearchCommand(ctx, issueCommand, options)
	addIssueVCSBranchSearchCommand(ctx, issueCommand, options)
	addIssueGetCommand(ctx, issueCommand, options)
	addIssueDepsCommand(ctx, issueCommand, options)
	addIssueAttachmentsCommand(ctx, issueCommand, options)
	addIssueBotActorCommand(ctx, issueCommand, options)
	addIssueChildrenCommand(ctx, issueCommand, options)
	addIssueDocumentsCommand(ctx, issueCommand, options)
	addIssueFormerAttachmentsCommand(ctx, issueCommand, options)
	addIssueHistoryCommand(ctx, issueCommand, options)
	addIssueInverseRelationsCommand(ctx, issueCommand, options)
	addIssueLabelsCommand(ctx, issueCommand, options)
	addIssueRelationsCommand(ctx, issueCommand, options)
	addIssueReleasesCommand(ctx, issueCommand, options)
	addIssueStateHistoryCommand(ctx, issueCommand, options)
	addIssueSubscribersCommand(ctx, issueCommand, options)
	addIssuePRCommand(ctx, issueCommand, options)
	addIssueCreateCommand(ctx, issueCommand, options)
	addIssueUpdateCommand(ctx, issueCommand, options)
	addIssueStartCommand(ctx, issueCommand, options)
	addIssueCommentCommand(ctx, issueCommand, options)
	addIssueReplyCommand(ctx, issueCommand, options)
	addIssueCommentsCommand(ctx, issueCommand, options)
	addIssueCloseCommand(ctx, issueCommand, options)
	addIssueCurrentCommands(ctx, issueCommand, options)
	addDomainUsageCommand(issueCommand, options, "issue")
	root.AddCommand(issueCommand)
}

func addIssueListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	stateType := ""
	projectID := ""
	assigneeID := ""
	labelID := ""
	cycleID := ""
	createdAfter := ""
	createdSince := ""
	createdBefore := ""
	hasBlockers := false
	blocks := false
	blockedBy := ""
	allTeams := false
	mine := false
	command := &cobra.Command{
		Use:   "list",
		Short: "List issues for the resolved team",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := validateIssueListFilters(
				stateType,
				projectID,
				assigneeID,
				labelID,
				cycleID,
				createdAfter,
				createdSince,
				createdBefore,
				hasBlockers,
				blocks,
				blockedBy,
				allTeams,
				mine,
			); err != nil {
				return err
			}
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issues, err := issueList(ctx, runtime, limit, issueListFlagValues{
				stateType:     stateType,
				projectID:     projectID,
				assigneeID:    assigneeID,
				labelID:       labelID,
				cycleID:       cycleID,
				createdAfter:  createdAfter,
				createdSince:  createdSince,
				createdBefore: createdBefore,
				hasBlockers:   hasBlockers,
				blocks:        blocks,
				blockedBy:     blockedBy,
				allTeams:      allTeams,
				mine:          mine,
			})
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(issues.Issues)); err != nil {
				return err
			}
			issues.Issues, err = sortByJSONField(issues.Issues, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, issues)
			}

			return writeIssues(command, options, issues.Issues)
		},
	}
	bindIssueListFlags(
		command,
		&limit,
		&stateType,
		&projectID,
		&assigneeID,
		&labelID,
		&cycleID,
		&createdAfter,
		&createdSince,
		&createdBefore,
		&hasBlockers,
		&blocks,
		&blockedBy,
		&allTeams,
		&mine,
	)
	root.AddCommand(command)
}

func bindIssueListFlags(
	command *cobra.Command,
	limit *int,
	stateType *string,
	projectID *string,
	assigneeID *string,
	labelID *string,
	cycleID *string,
	createdAfter *string,
	createdSince *string,
	createdBefore *string,
	hasBlockers *bool,
	blocks *bool,
	blockedBy *string,
	allTeams *bool,
	mine *bool,
) {
	command.Flags().IntVar(limit, "limit", *limit, "maximum issues to return")
	command.Flags().StringVar(stateType, "state", *stateType, "filter by workflow state type")
	command.Flags().StringVar(projectID, "project", *projectID, "filter by Linear project id")
	command.Flags().StringVar(assigneeID, "assignee", *assigneeID, "filter by Linear assignee user id")
	command.Flags().StringVar(labelID, "label", *labelID, "filter by Linear issue label id")
	command.Flags().StringVar(cycleID, "cycle", *cycleID, "filter by Linear cycle id")
	command.Flags().StringVar(createdAfter, "created-after", *createdAfter, "filter by created-at date lower bound")
	command.Flags().StringVar(createdSince, "created-since", *createdSince, "alias for --created-after")
	command.Flags().StringVar(createdBefore, "created-before", *createdBefore, "filter by created-at date upper bound")
	command.Flags().BoolVar(hasBlockers, "has-blockers", *hasBlockers, "filter to issues blocked by another issue")
	command.Flags().BoolVar(blocks, "blocks", *blocks, "filter to issues blocking another issue")
	command.Flags().StringVar(
		blockedBy,
		"blocked-by",
		*blockedBy,
		"filter to issues blocked by an issue id or identifier",
	)
	command.Flags().BoolVar(allTeams, "all-teams", *allTeams, "list issues across every visible Linear team")
	command.Flags().BoolVar(mine, "mine", *mine, "filter to issues assigned to the authenticated user")
}

type issueListFlagValues struct {
	stateType     string
	projectID     string
	assigneeID    string
	labelID       string
	cycleID       string
	createdAfter  string
	createdSince  string
	createdBefore string
	hasBlockers   bool
	blocks        bool
	blockedBy     string
	allTeams      bool
	mine          bool
}

func validateIssueListFilters(
	stateType string,
	projectID string,
	assigneeID string,
	labelID string,
	cycleID string,
	createdAfter string,
	createdSince string,
	createdBefore string,
	hasBlockers bool,
	blocks bool,
	blockedBy string,
	allTeams bool,
	mine bool,
) error {
	filterCount := 0
	if stateType != "" {
		filterCount++
	}
	if projectID != "" {
		filterCount++
	}
	if assigneeID != "" {
		filterCount++
	}
	if labelID != "" {
		filterCount++
	}
	if cycleID != "" {
		filterCount++
	}
	if createdAfter != "" {
		filterCount++
	}
	if createdSince != "" {
		filterCount++
	}
	if createdBefore != "" {
		filterCount++
	}
	if hasBlockers {
		filterCount++
	}
	if blocks {
		filterCount++
	}
	if blockedBy != "" {
		filterCount++
	}
	if allTeams {
		filterCount++
	}
	if mine {
		filterCount++
	}
	if filterCount > 1 {
		return errors.New(
			"issue list filters: use only one of --state, --project, --assignee, " +
				"--label, --cycle, --created-after, --created-since, --created-before, " +
				"--has-blockers, --blocks, --blocked-by, --all-teams, or --mine",
		)
	}

	return nil
}

func issueList(
	ctx context.Context,
	runtime commandRuntime,
	limit int,
	flags issueListFlagValues,
) (client.IssueList, error) {
	if flags.allTeams {
		return client.ListIssues(ctx, runtime.graphqlClient, limit)
	}

	target, err := runtime.resolveTarget(ctx)
	if err != nil {
		return client.IssueList{}, err
	}

	return client.ListIssuesByTeam(ctx, runtime.graphqlClient, target.Team.ID, limit,
		client.IssueListFilters{
			StateType:     flags.stateType,
			ProjectID:     flags.projectID,
			AssigneeID:    issueListAssigneeID(target, flags.assigneeID, flags.mine),
			LabelID:       flags.labelID,
			CycleID:       flags.cycleID,
			CreatedAfter:  issueListCreatedAfter(flags.createdAfter, flags.createdSince),
			CreatedBefore: flags.createdBefore,
			HasBlockers:   flags.hasBlockers,
			Blocks:        flags.blocks,
			BlockedBy:     flags.blockedBy,
		})
}

func issueListCreatedAfter(createdAfter string, createdSince string) string {
	if createdAfter != "" {
		return createdAfter
	}

	return createdSince
}

func issueListAssigneeID(target client.ResolvedTarget, assigneeID string, mine bool) string {
	if assigneeID != "" {
		return assigneeID
	}
	if mine {
		return target.Viewer.ID
	}

	return ""
}

func addIssueSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "search QUERY",
		Short: "Search issues for the resolved team",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			target, err := runtime.resolveTarget(ctx)
			if err != nil {
				return err
			}
			issues, err := client.SearchIssuesByTeam(ctx, runtime.graphqlClient, target.Team.ID, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(issues.Issues)); err != nil {
				return err
			}
			issues.Issues, err = sortByJSONField(issues.Issues, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, issues)
			}

			return writeIssues(command, options, issues.Issues)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addIssueVCSBranchSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	branchCommand := &cobra.Command{
		Use:   "vcs-branch-search",
		Short: "Read the issue matched by a VCS branch",
	}
	root.AddCommand(branchCommand)

	branchCommand.AddCommand(&cobra.Command{
		Use:   "get BRANCH_NAME",
		Short: "Get the issue matched by a VCS branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.GetIssueByVCSBranch(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
	addIssueVCSBranchAttachmentsCommand(ctx, branchCommand, options)
	addIssueVCSBranchBotActorCommand(ctx, branchCommand, options)
	addIssueVCSBranchChildrenCommand(ctx, branchCommand, options)
	addIssueVCSBranchDocumentsCommand(ctx, branchCommand, options)
	addIssueVCSBranchFormerAttachmentsCommand(ctx, branchCommand, options)
	addIssueVCSBranchHistoryCommand(ctx, branchCommand, options)
	addIssueVCSBranchInverseRelationsCommand(ctx, branchCommand, options)
	addIssueVCSBranchLabelsCommand(ctx, branchCommand, options)
	addIssueVCSBranchRelationsCommand(ctx, branchCommand, options)
	addIssueVCSBranchReleasesCommand(ctx, branchCommand, options)
	addIssueVCSBranchStateHistoryCommand(ctx, branchCommand, options)
	addIssueVCSBranchSubscribersCommand(ctx, branchCommand, options)
}

func addIssueVCSBranchAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"attachments BRANCH_NAME",
		"List attachments for the issue matched by a VCS branch",
		"attachments",
		func(runtime commandRuntime, branchName string, limit int) (client.AttachmentList, error) {
			return client.ListIssueVCSBranchAttachments(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		func(command *cobra.Command, item client.AttachmentSummary) error {
			return writeAttachment(command, options, item)
		},
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addIssueVCSBranchBotActorCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "bot-actor BRANCH_NAME",
		Short: "Show bot actor metadata for the issue matched by a VCS branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			actor, err := client.GetIssueVCSBranchBotActor(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssueBotActor(command, options, actor)
		},
	})
}

func addIssueVCSBranchChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"children BRANCH_NAME",
		"List child issues for the issue matched by a VCS branch",
		"child issues",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueList, error) {
			return client.ListIssueVCSBranchChildren(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.IssueList) int {
			return len(list.Issues)
		},
		func(list client.IssueList) (client.IssueList, error) {
			items, err := sortByJSONField(list.Issues, options.sortField, options.sortOrder)
			list.Issues = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueSummary) error {
			return writeIssue(command, options, item)
		},
		func(list client.IssueList) []client.IssueSummary {
			return list.Issues
		},
	)
}

func addIssueVCSBranchDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"documents BRANCH_NAME",
		"List documents for the issue matched by a VCS branch",
		"documents",
		func(runtime commandRuntime, branchName string, limit int) (client.DocumentList, error) {
			return client.ListIssueVCSBranchDocuments(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.DocumentList) int {
			return len(list.Documents)
		},
		func(list client.DocumentList) (client.DocumentList, error) {
			items, err := sortByJSONField(list.Documents, options.sortField, options.sortOrder)
			list.Documents = items
			return list, err
		},
		func(command *cobra.Command, item client.DocumentSummary) error {
			return writeDocument(command, options, item)
		},
		func(list client.DocumentList) []client.DocumentSummary {
			return list.Documents
		},
	)
}

func addIssueVCSBranchFormerAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"former-attachments BRANCH_NAME",
		"List former attachments for the issue matched by a VCS branch",
		"former attachments",
		func(runtime commandRuntime, branchName string, limit int) (client.AttachmentList, error) {
			return client.ListIssueVCSBranchFormerAttachments(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		func(command *cobra.Command, item client.AttachmentSummary) error {
			return writeAttachment(command, options, item)
		},
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addIssueVCSBranchHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"history BRANCH_NAME",
		"List history metadata for the issue matched by a VCS branch",
		"history entries",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueHistoryList, error) {
			return client.ListIssueVCSBranchHistory(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.IssueHistoryList) int {
			return len(list.History)
		},
		func(list client.IssueHistoryList) (client.IssueHistoryList, error) {
			items, err := sortByJSONField(list.History, options.sortField, options.sortOrder)
			list.History = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueHistorySummary) error {
			return writeIssueHistory(command, options, item)
		},
		func(list client.IssueHistoryList) []client.IssueHistorySummary {
			return list.History
		},
	)
}

func addIssueVCSBranchInverseRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueRelationChildListCommand(
		ctx,
		root,
		options,
		"inverse-relations BRANCH_NAME",
		"List inverse relations for the issue matched by a VCS branch",
		"inverse relations",
		func(
			ctx context.Context,
			runtime commandRuntime,
			branchName string,
			limit int,
		) (client.IssueRelationList, error) {
			return client.ListIssueVCSBranchInverseRelations(ctx, runtime.graphqlClient, branchName, limit)
		},
	)
}

func addIssueVCSBranchLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"labels BRANCH_NAME",
		"List labels for the issue matched by a VCS branch",
		"labels",
		func(runtime commandRuntime, branchName string, limit int) (client.LabelList, error) {
			return client.ListIssueVCSBranchLabels(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.LabelList) int {
			return len(list.Labels)
		},
		func(list client.LabelList) (client.LabelList, error) {
			items, err := sortByJSONField(list.Labels, options.sortField, options.sortOrder)
			list.Labels = items
			return list, err
		},
		func(command *cobra.Command, item client.LabelSummary) error {
			return writeLabel(command, options, item)
		},
		func(list client.LabelList) []client.LabelSummary {
			return list.Labels
		},
	)
}

func addIssueVCSBranchRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueRelationChildListCommand(
		ctx,
		root,
		options,
		"relations BRANCH_NAME",
		"List relations for the issue matched by a VCS branch",
		"relations",
		func(
			ctx context.Context,
			runtime commandRuntime,
			branchName string,
			limit int,
		) (client.IssueRelationList, error) {
			return client.ListIssueVCSBranchRelations(ctx, runtime.graphqlClient, branchName, limit)
		},
	)
}

func addIssueVCSBranchReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"releases BRANCH_NAME",
		"List releases for the issue matched by a VCS branch",
		"releases",
		func(runtime commandRuntime, branchName string, limit int) (client.ReleaseList, error) {
			return client.ListIssueVCSBranchReleases(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.ReleaseList) int {
			return len(list.Releases)
		},
		func(list client.ReleaseList) (client.ReleaseList, error) {
			items, err := sortByJSONField(list.Releases, options.sortField, options.sortOrder)
			list.Releases = items
			return list, err
		},
		func(command *cobra.Command, item client.ReleaseSummary) error {
			return writeRelease(command, options, item)
		},
		func(list client.ReleaseList) []client.ReleaseSummary {
			return list.Releases
		},
	)
}

func addIssueVCSBranchStateHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"state-history BRANCH_NAME",
		"List workflow state history for the issue matched by a VCS branch",
		"state spans",
		func(runtime commandRuntime, branchName string, limit int) (client.IssueStateHistoryList, error) {
			return client.ListIssueVCSBranchStateHistory(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.IssueStateHistoryList) int {
			return len(list.Spans)
		},
		func(list client.IssueStateHistoryList) (client.IssueStateHistoryList, error) {
			items, err := sortByJSONField(list.Spans, options.sortField, options.sortOrder)
			list.Spans = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueStateSpanSummary) error {
			return writeIssueStateSpan(command, options, item)
		},
		func(list client.IssueStateHistoryList) []client.IssueStateSpanSummary {
			return list.Spans
		},
	)
}

func addIssueVCSBranchSubscribersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"subscribers BRANCH_NAME",
		"List subscribers for the issue matched by a VCS branch",
		"subscribers",
		func(runtime commandRuntime, branchName string, limit int) (client.UserList, error) {
			return client.ListIssueVCSBranchSubscribers(ctx, runtime.graphqlClient, branchName, limit)
		},
		func(list client.UserList) int {
			return len(list.Users)
		},
		func(list client.UserList) (client.UserList, error) {
			items, err := sortByJSONField(list.Users, options.sortField, options.sortOrder)
			list.Users = items
			return list, err
		},
		func(command *cobra.Command, item client.UserSummary) error {
			return writeUser(command, options, item)
		},
		func(list client.UserList) []client.UserSummary {
			return list.Users
		},
	)
}

func addIssueGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get ISSUE_ID",
		Short: "Get one issue by id or identifier",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			issue, err := client.GetIssueByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}

func addIssueAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"attachments ISSUE_ID",
		"List issue attachments",
		"attachments",
		func(runtime commandRuntime, issueID string, limit int) (client.AttachmentList, error) {
			return client.ListIssueAttachments(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		func(command *cobra.Command, item client.AttachmentSummary) error {
			return writeAttachment(command, options, item)
		},
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addIssueBotActorCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "bot-actor ISSUE_ID",
		Short: "Show issue bot actor metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			actor, err := client.GetIssueBotActor(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssueBotActor(command, options, actor)
		},
	})
}

func addIssueChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"children ISSUE_ID",
		"List issue children",
		"child issues",
		func(runtime commandRuntime, issueID string, limit int) (client.IssueList, error) {
			return client.ListIssueChildren(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.IssueList) int {
			return len(list.Issues)
		},
		func(list client.IssueList) (client.IssueList, error) {
			items, err := sortByJSONField(list.Issues, options.sortField, options.sortOrder)
			list.Issues = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueSummary) error {
			return writeIssue(command, options, item)
		},
		func(list client.IssueList) []client.IssueSummary {
			return list.Issues
		},
	)
}

func addIssueDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"documents ISSUE_ID",
		"List issue documents",
		"documents",
		func(runtime commandRuntime, issueID string, limit int) (client.DocumentList, error) {
			return client.ListIssueDocuments(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.DocumentList) int {
			return len(list.Documents)
		},
		func(list client.DocumentList) (client.DocumentList, error) {
			items, err := sortByJSONField(list.Documents, options.sortField, options.sortOrder)
			list.Documents = items
			return list, err
		},
		func(command *cobra.Command, item client.DocumentSummary) error {
			return writeDocument(command, options, item)
		},
		func(list client.DocumentList) []client.DocumentSummary {
			return list.Documents
		},
	)
}

func addIssueFormerAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"former-attachments ISSUE_ID",
		"List former issue attachments",
		"former attachments",
		func(runtime commandRuntime, issueID string, limit int) (client.AttachmentList, error) {
			return client.ListIssueFormerAttachments(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		func(command *cobra.Command, item client.AttachmentSummary) error {
			return writeAttachment(command, options, item)
		},
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addIssueHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"history ISSUE_ID",
		"List issue history metadata",
		"history entries",
		func(runtime commandRuntime, issueID string, limit int) (client.IssueHistoryList, error) {
			return client.ListIssueHistory(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.IssueHistoryList) int {
			return len(list.History)
		},
		func(list client.IssueHistoryList) (client.IssueHistoryList, error) {
			items, err := sortByJSONField(list.History, options.sortField, options.sortOrder)
			list.History = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueHistorySummary) error {
			return writeIssueHistory(command, options, item)
		},
		func(list client.IssueHistoryList) []client.IssueHistorySummary {
			return list.History
		},
	)
}

func addIssueInverseRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueRelationChildListCommand(
		ctx,
		root,
		options,
		"inverse-relations ISSUE_ID",
		"List issue inverse relations",
		"inverse relations",
		func(ctx context.Context, runtime commandRuntime, issueID string, limit int) (client.IssueRelationList, error) {
			return client.ListIssueInverseRelations(ctx, runtime.graphqlClient, issueID, limit)
		},
	)
}

func addIssueLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"labels ISSUE_ID",
		"List issue labels",
		"labels",
		func(runtime commandRuntime, issueID string, limit int) (client.LabelList, error) {
			return client.ListIssueLabels(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.LabelList) int {
			return len(list.Labels)
		},
		func(list client.LabelList) (client.LabelList, error) {
			items, err := sortByJSONField(list.Labels, options.sortField, options.sortOrder)
			list.Labels = items
			return list, err
		},
		func(command *cobra.Command, item client.LabelSummary) error {
			return writeLabel(command, options, item)
		},
		func(list client.LabelList) []client.LabelSummary {
			return list.Labels
		},
	)
}

func addIssueRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueRelationChildListCommand(
		ctx,
		root,
		options,
		"relations ISSUE_ID",
		"List issue relations",
		"relations",
		func(ctx context.Context, runtime commandRuntime, issueID string, limit int) (client.IssueRelationList, error) {
			return client.ListIssueRelationsForIssue(ctx, runtime.graphqlClient, issueID, limit)
		},
	)
}

func addIssueReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"releases ISSUE_ID",
		"List issue releases",
		"releases",
		func(runtime commandRuntime, issueID string, limit int) (client.ReleaseList, error) {
			return client.ListIssueReleases(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.ReleaseList) int {
			return len(list.Releases)
		},
		func(list client.ReleaseList) (client.ReleaseList, error) {
			items, err := sortByJSONField(list.Releases, options.sortField, options.sortOrder)
			list.Releases = items
			return list, err
		},
		func(command *cobra.Command, item client.ReleaseSummary) error {
			return writeRelease(command, options, item)
		},
		func(list client.ReleaseList) []client.ReleaseSummary {
			return list.Releases
		},
	)
}

func addIssueStateHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"state-history ISSUE_ID",
		"List issue workflow state history",
		"state spans",
		func(runtime commandRuntime, issueID string, limit int) (client.IssueStateHistoryList, error) {
			return client.ListIssueStateHistory(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.IssueStateHistoryList) int {
			return len(list.Spans)
		},
		func(list client.IssueStateHistoryList) (client.IssueStateHistoryList, error) {
			items, err := sortByJSONField(list.Spans, options.sortField, options.sortOrder)
			list.Spans = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueStateSpanSummary) error {
			return writeIssueStateSpan(command, options, item)
		},
		func(list client.IssueStateHistoryList) []client.IssueStateSpanSummary {
			return list.Spans
		},
	)
}

func addIssueSubscribersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		"subscribers ISSUE_ID",
		"List issue subscribers",
		"subscribers",
		func(runtime commandRuntime, issueID string, limit int) (client.UserList, error) {
			return client.ListIssueSubscribers(ctx, runtime.graphqlClient, issueID, limit)
		},
		func(list client.UserList) int {
			return len(list.Users)
		},
		func(list client.UserList) (client.UserList, error) {
			items, err := sortByJSONField(list.Users, options.sortField, options.sortOrder)
			list.Users = items
			return list, err
		},
		func(command *cobra.Command, item client.UserSummary) error {
			return writeUser(command, options, item)
		},
		func(list client.UserList) []client.UserSummary {
			return list.Users
		},
	)
}

func addIssueRelationChildListCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch func(context.Context, commandRuntime, string, int) (client.IssueRelationList, error),
) {
	addIssueChildListCommand(
		ctx,
		root,
		options,
		use,
		short,
		limitHelp,
		func(runtime commandRuntime, issueID string, limit int) (client.IssueRelationList, error) {
			return fetch(ctx, runtime, issueID, limit)
		},
		func(list client.IssueRelationList) int {
			return len(list.Relations)
		},
		func(list client.IssueRelationList) (client.IssueRelationList, error) {
			items, err := sortByJSONField(list.Relations, options.sortField, options.sortOrder)
			list.Relations = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueRelationSummary) error {
			return writeIssueRelation(command, options, item)
		},
		func(list client.IssueRelationList) []client.IssueRelationSummary {
			return list.Relations
		},
	)
}

func addIssueChildListCommand[List any, Item any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch func(commandRuntime, string, int) (List, error),
	count func(List) int,
	sortList func(List) (List, error),
	writeItem func(*cobra.Command, Item) error,
	items func(List) []Item,
) {
	limit := 50
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			list, err := fetch(runtime, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, count(list)); err != nil {
				return err
			}
			list, err = sortList(list)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, list)
			}
			for _, item := range items(list) {
				if err := writeItem(command, item); err != nil {
					return err
				}
			}

			return nil
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum "+limitHelp+" to return")
	root.AddCommand(command)
}

func writeIssueBotActor(command *cobra.Command, options *rootOptions, actor client.IssueBotActor) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, actor)
	}
	if actor.Bot == nil {
		return render.WriteLine(command.OutOrStdout(), "%s bot -", actor.IssueID)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s bot %s %s [%s]",
		actor.IssueID,
		emptyDash(actor.Bot.ID),
		emptyDash(actor.Bot.Name),
		actor.Bot.Type,
	)
}

func writeIssueStateSpan(command *cobra.Command, options *rootOptions, span client.IssueStateSpanSummary) error {
	if wrote, err := writeIDOnly(command, options, span.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, span)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s %s -> %s",
		span.ID,
		emptyDash(span.StateName),
		emptyDash(span.StateType),
		span.StartedAt,
		emptyDash(span.EndedAt),
	)
}

func writeIssueHistory(command *cobra.Command, options *rootOptions, history client.IssueHistorySummary) error {
	if wrote, err := writeIDOnly(command, options, history.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, history)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s issue %s updated_description %t",
		history.ID,
		history.IssueID,
		history.UpdatedDescription,
	)
}

func writeIssues(command *cobra.Command, options *rootOptions, issues []client.IssueSummary) error {
	for _, issue := range issues {
		if err := writeIssue(command, options, issue); err != nil {
			return err
		}
	}

	return nil
}
