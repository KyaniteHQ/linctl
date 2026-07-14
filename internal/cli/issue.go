package cli

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addIssueCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	issueCommand := newGroupCommand("issue", "Read and write Linear issues")
	addIssueListCommand(ctx, issueCommand, options)
	addIssueSearchCommand(ctx, issueCommand, options)
	addIssueFigmaFileKeySearchCommand(ctx, issueCommand, options)
	addIssuePriorityValuesCommand(ctx, issueCommand, options)
	addIssueFilterSuggestionCommand(ctx, issueCommand, options)
	addIssueTitleSuggestionCommand(ctx, issueCommand, options)
	addIssueVCSBranchSearchCommand(ctx, issueCommand, options)
	addIssueGetCommand(ctx, issueCommand, options)
	addIssueDepsCommand(ctx, issueCommand, options)
	addIssueChildCommands(ctx, issueCommand, options, issueChildCommandBundleForIssue())
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
	addIssueAddLabelCommand(ctx, issueCommand, options)
	addIssueRemoveLabelCommand(ctx, issueCommand, options)
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

			return runIssueList(ctx, command, options, issueListClientFor(runtime), limit, flags)
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

			return runIssueSearch(ctx, command, options, issueSearchClientFor(runtime), args[0], limit)
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
				writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(preflightReadListCommand(command, loadIssueFigmaFileKeySearch))
}

func loadIssueFigmaFileKeySearch(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.SearchIssuesByFigmaFileKey(ctx, runtime.graphqlClient, args[0], limit)

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
			values, err := client.ListIssuePriorityValues(ctx, runtime.graphqlClient)
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
			suggestion, err := client.GetIssueFilterSuggestion(
				ctx,
				runtime.graphqlClient,
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
			suggestion, err := client.GetIssueTitleSuggestionFromCustomerRequest(
				ctx, runtime.graphqlClient, args[0],
			)
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
			issue, err := client.GetIssueByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeIssue(command, options, issue)
		},
	})
}

func issueChildCommandBundleForIssue() issueChildCommandBundle {
	return issueChildCommandBundle{
		Argument:          "ISSUE_ID",
		Text:              directIssueChildCommandText(),
		Attachments:       client.ListIssueAttachments,
		BotActor:          client.GetIssueBotActor,
		Children:          client.ListIssueChildren,
		Documents:         client.ListIssueDocuments,
		FormerAttachments: client.ListIssueFormerAttachments,
		FormerNeeds:       client.ListIssueFormerNeeds,
		History:           client.ListIssueHistory,
		InverseRelations:  client.ListIssueInverseRelations,
		Labels:            client.ListIssueLabels,
		Needs:             client.ListIssueNeeds,
		Relations:         client.ListIssueRelationsForIssue,
		Releases:          client.ListIssueReleases,
		SharedAccess:      client.GetIssueSharedAccess,
		StateHistory:      client.ListIssueStateHistory,
		Subscribers:       client.ListIssueSubscribers,
	}
}

func addIssueCommentMetadataListCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch issueChildListFetcher[client.IssueCommentMetadataList],
) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueCommentMetadataList, client.CommentMetadataSummary]{
		Use:       use,
		Short:     short,
		LimitHelp: limitHelp,
		Args:      cobra.ExactArgs(1),
		Load: func(
			_ context.Context, runtime commandRuntime, args []string, limit int,
		) (client.IssueCommentMetadataList, []client.CommentMetadataSummary, error) {
			list, err := fetch(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Comments, err
		},
		WriteItem: writeCommentMetadata,
	})
}
