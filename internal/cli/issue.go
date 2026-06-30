//nolint:dupl // Issue child read commands intentionally share the same list-command shape.
package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	issueCommand := &cobra.Command{
		Use:   "issue",
		Short: "Read and write Linear issues",
	}
	addIssueListCommand(ctx, issueCommand, options)
	addIssueSearchCommand(ctx, issueCommand, options)
	addIssueFigmaFileKeySearchCommand(ctx, issueCommand, options)
	addIssuePriorityValuesCommand(ctx, issueCommand, options)
	addIssueFilterSuggestionCommand(ctx, issueCommand, options)
	addIssueTitleSuggestionCommand(ctx, issueCommand, options)
	addIssueVCSBranchSearchCommand(ctx, issueCommand, options)
	addIssueGetCommand(ctx, issueCommand, options)
	addIssueDepsCommand(ctx, issueCommand, options)
	addIssueAttachmentsCommand(ctx, issueCommand, options)
	addIssueBotActorCommand(ctx, issueCommand, options)
	addIssueChildrenCommand(ctx, issueCommand, options)
	addIssueDocumentsCommand(ctx, issueCommand, options)
	addIssueFormerAttachmentsCommand(ctx, issueCommand, options)
	addIssueFormerNeedsCommand(ctx, issueCommand, options)
	addIssueHistoryCommand(ctx, issueCommand, options)
	addIssueInverseRelationsCommand(ctx, issueCommand, options)
	addIssueLabelsCommand(ctx, issueCommand, options)
	addIssueNeedsCommand(ctx, issueCommand, options)
	addIssueRelationsCommand(ctx, issueCommand, options)
	addIssueReleasesCommand(ctx, issueCommand, options)
	addIssueSharedAccessCommand(ctx, issueCommand, options)
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
	addIssueRelateCommand(ctx, issueCommand, options)
	addIssueUnrelateCommand(ctx, issueCommand, options)
	addIssueLinkCommand(ctx, issueCommand, options)
	addIssueOpenCommand(ctx, issueCommand, options)
	addIssueExportCommand(ctx, issueCommand, options)
	addIssueImportCommand(ctx, issueCommand, options)
	addIssueBulkExportCommand(ctx, issueCommand, options)
	addIssueCurrentCommands(ctx, issueCommand, options)
	addDomainUsageCommand(issueCommand, options, "issue")
	root.AddCommand(issueCommand)
}

func addIssueListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	stateType := ""
	status := ""
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
			normalizedState, err := normalizeAndNote(
				command, "state", mergedStateFlag(stateType, status), normalizedStateType,
			)
			if err != nil {
				return err
			}
			flags := issueListFlagValues{
				stateType:     normalizedState,
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
			}
			if err := validateIssueListFilters(flags); err != nil {
				return err
			}
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runIssueList(ctx, command, options, issueAdapterFor(runtime), limit, flags)
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
	command.Flags().StringVar(&status, "status", "", "alias for --state")
	registerStateCompletion(ctx, command, options)
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

func validateIssueListFilters(flags issueListFlagValues) error {
	filterCount := 0
	for _, active := range []bool{
		flags.stateType != "",
		flags.projectID != "",
		flags.assigneeID != "",
		flags.labelID != "",
		flags.cycleID != "",
		flags.createdAfter != "",
		flags.createdSince != "",
		flags.createdBefore != "",
		flags.hasBlockers,
		flags.blocks,
		flags.blockedBy != "",
		flags.allTeams,
		flags.mine,
	} {
		if active {
			filterCount++
		}
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

func runIssueList(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	reader issueReader,
	limit int,
	flags issueListFlagValues,
) error {
	issues, err := issueList(ctx, reader, limit, flags)
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
}

func issueList(
	ctx context.Context,
	reader issueReader,
	limit int,
	flags issueListFlagValues,
) (client.IssueList, error) {
	if flags.allTeams {
		return reader.ListIssues(ctx, limit)
	}

	target, err := reader.ResolveTarget(ctx)
	if err != nil {
		return client.IssueList{}, err
	}

	return reader.ListIssuesByTeam(ctx, target.Team.ID, limit,
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

			return runIssueSearch(ctx, command, options, issueAdapterFor(runtime), args[0], limit)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func runIssueSearch(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	searcher issueSearcher,
	query string,
	limit int,
) error {
	target, err := searcher.ResolveTarget(ctx)
	if err != nil {
		return err
	}
	issues, err := searcher.SearchIssuesByTeam(ctx, target.Team.ID, query, limit)
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
		annotateReadCollectionCommand(command, collectionKeyForPage[client.IssueList]())
		return writeJSONValue(command, options, issues)
	}

	return writeIssues(command, options, issues.Issues)
}

func addIssueFigmaFileKeySearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "figma-file-key-search FILE_KEY",
		Short: "Search issues linked to a Figma file key",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadIssueFigmaFileKeySearch,
				issuePageWithItems,
				writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func loadIssueFigmaFileKeySearch(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := issueAdapterFor(runtime).SearchIssuesByFigmaFileKey(ctx, args[0], limit)

	return issues, issues.Issues, err
}

func addIssuePriorityValuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "priority-values",
		Short: "List issue priority values",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			values, err := issueAdapterFor(runtime).ListIssuePriorityValues(ctx)
			if err != nil {
				return err
			}

			return writeIssuePriorityValues(command, options, values)
		},
	})
}

func addIssueFilterSuggestionCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	teamID := ""
	projectID := ""
	command := &cobra.Command{
		Use:   "filter-suggestion PROMPT",
		Short: "Suggest an issue filter from a text prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			if teamID != "" && projectID != "" {
				return errors.New("issue filter suggestion: use only one of --team-id or --project-id")
			}
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			suggestion, err := issueAdapterFor(runtime).GetIssueFilterSuggestion(
				ctx,
				args[0],
				teamID,
				projectID,
			)
			if err != nil {
				return err
			}

			return writeIssueFilterSuggestion(command, options, suggestion)
		},
	}
	command.Flags().StringVar(&teamID, "team-id", teamID, "optional team id for team-scoped issue views")
	command.Flags().StringVar(&projectID, "project-id", projectID, "optional project id for project-scoped issue views")
	root.AddCommand(command)
}

func addIssueTitleSuggestionCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "title-suggestion REQUEST",
		Short: "Suggest an issue title from customer request text",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			suggestion, err := issueAdapterFor(runtime).GetIssueTitleSuggestionFromCustomerRequest(ctx, args[0])
			if err != nil {
				return err
			}

			return writeIssueTitleSuggestion(command, options, suggestion)
		},
	})
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
			issue, err := issueAdapterFor(runtime).GetIssueByID(ctx, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}

func addIssueAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"attachments ISSUE_ID",
		"List issue attachments",
		"attachments",
		func(runtime commandRuntime, issueID string, limit int) (client.AttachmentList, error) {
			return issueAdapterFor(runtime).ListIssueAttachments(ctx, issueID, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		writeAttachment,
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
			actor, err := issueAdapterFor(runtime).GetIssueBotActor(ctx, args[0])
			if err != nil {
				return err
			}

			return writeIssueBotActor(command, options, actor)
		},
	})
}

func addIssueChildrenCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"children ISSUE_ID",
		"List issue children",
		"child issues",
		func(runtime commandRuntime, issueID string, limit int) (client.IssueList, error) {
			return issueAdapterFor(runtime).ListIssueChildren(ctx, issueID, limit)
		},
		func(list client.IssueList) int {
			return len(list.Issues)
		},
		func(list client.IssueList) (client.IssueList, error) {
			items, err := sortByJSONField(list.Issues, options.sortField, options.sortOrder)
			list.Issues = items
			return list, err
		},
		writeIssue,
		func(list client.IssueList) []client.IssueSummary {
			return list.Issues
		},
	)
}

func addIssueDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"documents ISSUE_ID",
		"List issue documents",
		"documents",
		func(runtime commandRuntime, issueID string, limit int) (client.DocumentList, error) {
			return issueAdapterFor(runtime).ListIssueDocuments(ctx, issueID, limit)
		},
		func(list client.DocumentList) int {
			return len(list.Documents)
		},
		func(list client.DocumentList) (client.DocumentList, error) {
			items, err := sortByJSONField(list.Documents, options.sortField, options.sortOrder)
			list.Documents = items
			return list, err
		},
		writeDocument,
		func(list client.DocumentList) []client.DocumentSummary {
			return list.Documents
		},
	)
}

func addIssueFormerAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"former-attachments ISSUE_ID",
		"List former issue attachments",
		"former attachments",
		func(runtime commandRuntime, issueID string, limit int) (client.AttachmentList, error) {
			return issueAdapterFor(runtime).ListIssueFormerAttachments(ctx, issueID, limit)
		},
		func(list client.AttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.AttachmentList) (client.AttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		writeAttachment,
		func(list client.AttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addIssueFormerNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueCustomerNeedMetadataListCommand(
		ctx,
		root,
		options,
		"former-needs ISSUE_ID",
		"List body-free former issue customer needs",
		"former customer needs",
		func(runtime commandRuntime, issueID string, limit int) (client.IssueCustomerNeedMetadataList, error) {
			return issueAdapterFor(runtime).ListIssueFormerNeeds(ctx, issueID, limit)
		},
	)
}

func addIssueHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"history ISSUE_ID",
		"List issue history metadata",
		"history entries",
		func(runtime commandRuntime, issueID string, limit int) (client.IssueHistoryList, error) {
			return issueAdapterFor(runtime).ListIssueHistory(ctx, issueID, limit)
		},
		func(list client.IssueHistoryList) int {
			return len(list.History)
		},
		func(list client.IssueHistoryList) (client.IssueHistoryList, error) {
			items, err := sortByJSONField(list.History, options.sortField, options.sortOrder)
			list.History = items
			return list, err
		},
		writeIssueHistory,
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
			return issueAdapterFor(runtime).ListIssueInverseRelations(ctx, issueID, limit)
		},
	)
}

func addIssueLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"labels ISSUE_ID",
		"List issue labels",
		"labels",
		func(runtime commandRuntime, issueID string, limit int) (client.LabelList, error) {
			return issueAdapterFor(runtime).ListIssueLabels(ctx, issueID, limit)
		},
		func(list client.LabelList) int {
			return len(list.Labels)
		},
		func(list client.LabelList) (client.LabelList, error) {
			items, err := sortByJSONField(list.Labels, options.sortField, options.sortOrder)
			list.Labels = items
			return list, err
		},
		writeLabel,
		func(list client.LabelList) []client.LabelSummary {
			return list.Labels
		},
	)
}

func addIssueNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addIssueCustomerNeedMetadataListCommand(
		ctx,
		root,
		options,
		"needs ISSUE_ID",
		"List body-free issue customer needs",
		"customer needs",
		func(runtime commandRuntime, issueID string, limit int) (client.IssueCustomerNeedMetadataList, error) {
			return issueAdapterFor(runtime).ListIssueNeeds(ctx, issueID, limit)
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
			return issueAdapterFor(runtime).ListIssueRelationsForIssue(ctx, issueID, limit)
		},
	)
}

func addIssueReleasesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"releases ISSUE_ID",
		"List issue releases",
		"releases",
		func(runtime commandRuntime, issueID string, limit int) (client.ReleaseList, error) {
			return issueAdapterFor(runtime).ListIssueReleases(ctx, issueID, limit)
		},
		func(list client.ReleaseList) int {
			return len(list.Releases)
		},
		func(list client.ReleaseList) (client.ReleaseList, error) {
			items, err := sortByJSONField(list.Releases, options.sortField, options.sortOrder)
			list.Releases = items
			return list, err
		},
		writeRelease,
		func(list client.ReleaseList) []client.ReleaseSummary {
			return list.Releases
		},
	)
}

func addIssueSharedAccessCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "shared-access ISSUE_ID",
		Short: "Show issue shared-access metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			access, err := issueAdapterFor(runtime).GetIssueSharedAccess(ctx, args[0])
			if err != nil {
				return err
			}

			return writeIssueSharedAccess(command, options, access)
		},
	})
}

func addIssueStateHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"state-history ISSUE_ID",
		"List issue workflow state history",
		"state spans",
		func(runtime commandRuntime, issueID string, limit int) (client.IssueStateHistoryList, error) {
			return issueAdapterFor(runtime).ListIssueStateHistory(ctx, issueID, limit)
		},
		func(list client.IssueStateHistoryList) int {
			return len(list.Spans)
		},
		func(list client.IssueStateHistoryList) (client.IssueStateHistoryList, error) {
			items, err := sortByJSONField(list.Spans, options.sortField, options.sortOrder)
			list.Spans = items
			return list, err
		},
		writeIssueStateSpan,
		func(list client.IssueStateHistoryList) []client.IssueStateSpanSummary {
			return list.Spans
		},
	)
}

func addIssueSubscribersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"subscribers ISSUE_ID",
		"List issue subscribers",
		"subscribers",
		func(runtime commandRuntime, issueID string, limit int) (client.UserList, error) {
			return issueAdapterFor(runtime).ListIssueSubscribers(ctx, issueID, limit)
		},
		func(list client.UserList) int {
			return len(list.Users)
		},
		func(list client.UserList) (client.UserList, error) {
			items, err := sortByJSONField(list.Users, options.sortField, options.sortOrder)
			list.Users = items
			return list, err
		},
		writeUser,
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
	addChildListCommand(
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
		writeIssueRelation,
		func(list client.IssueRelationList) []client.IssueRelationSummary {
			return list.Relations
		},
	)
}

func addIssueCommentMetadataListCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch func(commandRuntime, string, int) (client.IssueCommentMetadataList, error),
) {
	addChildListCommand(
		ctx,
		root,
		options,
		use,
		short,
		limitHelp,
		fetch,
		func(list client.IssueCommentMetadataList) int {
			return len(list.Comments)
		},
		func(list client.IssueCommentMetadataList) (client.IssueCommentMetadataList, error) {
			items, err := sortByJSONField(list.Comments, options.sortField, options.sortOrder)
			list.Comments = items
			return list, err
		},
		writeCommentMetadata,
		func(list client.IssueCommentMetadataList) []client.CommentMetadataSummary {
			return list.Comments
		},
	)
}

func addIssueCustomerNeedMetadataListCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch func(commandRuntime, string, int) (client.IssueCustomerNeedMetadataList, error),
) {
	addChildListCommand(
		ctx,
		root,
		options,
		use,
		short,
		limitHelp,
		fetch,
		func(list client.IssueCustomerNeedMetadataList) int {
			return len(list.Needs)
		},
		func(list client.IssueCustomerNeedMetadataList) (client.IssueCustomerNeedMetadataList, error) {
			items, err := sortByJSONField(list.Needs, options.sortField, options.sortOrder)
			list.Needs = items
			return list, err
		},
		writeCustomerNeedMetadata,
		func(list client.IssueCustomerNeedMetadataList) []client.CustomerNeedMetadataSummary {
			return list.Needs
		},
	)
}
