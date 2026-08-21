package cli

import (
	"cmp"
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
	addIssueMoveTeamCommand(ctx, issueCommand, options)
	addIssueMoveProjectCommand(ctx, issueCommand, options)
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
	status := ""
	flags := issueListFlagValues{}
	command := &cobra.Command{
		Use:   "list",
		Short: "List issues for the resolved team",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			normalizedState, err := normalizeAndNote(
				command, "state", cmp.Or(flags.stateType, status), normalizedStateType,
			)
			if err != nil {
				return err
			}
			resolved := flags
			resolved.stateType = normalizedState
			if err := validateIssueListFilters(resolved); err != nil {
				return err
			}
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runIssueList(ctx, command, options, issueListClientFor(runtime), limit, resolved)
		},
	}
	bindIssueListFlags(command, &limit, &flags)
	command.Flags().StringVar(&status, "status", "", "alias for --state")
	annotateReadCollectionCommand(command, collectionKeyForPage[client.IssueList]())
	registerStateCompletion(ctx, command, options, workflowStateTypeCandidates)
	registerFlagCompletion(command, "project", flagCompletion(ctx, options, projectIDCandidates))
	root.AddCommand(command)
}

func bindIssueListFlags(command *cobra.Command, limit *int, flags *issueListFlagValues) {
	command.Flags().IntVar(limit, "limit", *limit, "maximum issues to return")
	command.Flags().StringVar(&flags.stateType, "state", flags.stateType, "filter by workflow state type")
	command.Flags().StringVar(&flags.projectID, "project", flags.projectID, "filter by Linear project id")
	command.Flags().StringVar(&flags.assigneeID, "assignee", flags.assigneeID, "filter by Linear assignee user id")
	command.Flags().StringVar(&flags.labelID, "label", flags.labelID, "filter by Linear issue label id")
	command.Flags().StringVar(&flags.cycleID, "cycle", flags.cycleID, "filter by Linear cycle id")
	command.Flags().StringVar(
		&flags.createdAfter, "created-after", flags.createdAfter, "filter by the earliest created-at date",
	)
	command.Flags().StringVar(&flags.createdSince, "created-since", flags.createdSince, "alias for --created-after")
	command.Flags().StringVar(
		&flags.createdBefore, "created-before", flags.createdBefore, "filter by the latest created-at date",
	)
	command.Flags().StringVar(
		&flags.updatedAfter, "updated-after", flags.updatedAfter, "filter by the earliest updated-at date",
	)
	command.Flags().StringVar(
		&flags.updatedBefore, "updated-before", flags.updatedBefore, "filter by the latest updated-at date",
	)
	command.Flags().BoolVar(
		&flags.hasBlockers, "has-blockers", flags.hasBlockers, "filter to issues blocked by another issue",
	)
	command.Flags().BoolVar(&flags.blocks, "blocks", flags.blocks, "filter to issues blocking another issue")
	command.Flags().StringVar(
		&flags.blockedBy,
		"blocked-by",
		flags.blockedBy,
		"filter to issues blocked by an issue id or identifier",
	)
	command.Flags().BoolVar(
		&flags.allTeams, "all-teams", flags.allTeams, "list issues across every visible Linear team",
	)
	command.Flags().BoolVar(&flags.mine, "mine", flags.mine, "filter to issues assigned to the authenticated user")
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
	updatedAfter  string
	updatedBefore string
	hasBlockers   bool
	blocks        bool
	blockedBy     string
	allTeams      bool
	mine          bool
}

// validateIssueListFilters rejects only genuinely conflicting filter
// combinations; every remaining combination composes into one IssueFilter.
func validateIssueListFilters(flags issueListFlagValues) error {
	teamScopedCount := 0
	for _, active := range []bool{
		flags.stateType != "",
		flags.projectID != "",
		flags.assigneeID != "",
		flags.labelID != "",
		flags.cycleID != "",
		flags.createdAfter != "",
		flags.createdSince != "",
		flags.createdBefore != "",
		flags.updatedAfter != "",
		flags.updatedBefore != "",
		flags.hasBlockers,
		flags.blocks,
		flags.blockedBy != "",
		flags.mine,
	} {
		if active {
			teamScopedCount++
		}
	}
	if flags.allTeams && teamScopedCount > 0 {
		return errors.New("issue list filters: --all-teams cannot combine with team-scoped filters")
	}
	if flags.mine && flags.assigneeID != "" {
		return errors.New("issue list filters: use only one of --mine or --assignee")
	}
	if flags.createdAfter != "" && flags.createdSince != "" {
		return errors.New("issue list filters: use only one of --created-after or --created-since")
	}
	if flags.blockedBy != "" && teamScopedCount > 1 {
		return errors.New(
			"issue list filters: --blocked-by resolves issue relations and cannot combine with other filters",
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
			CreatedAfter:  cmp.Or(flags.createdAfter, flags.createdSince),
			CreatedBefore: flags.createdBefore,
			UpdatedAfter:  flags.updatedAfter,
			UpdatedBefore: flags.updatedBefore,
			HasBlockers:   flags.hasBlockers,
			Blocks:        flags.blocks,
			BlockedBy:     flags.blockedBy,
		})
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
	annotateReadCollectionCommand(command, collectionKeyForPage[client.IssueList]())
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
		return writeJSONValue(command, options, issues)
	}

	return writeIssues(command, options, issues.Issues)
}

func addIssueFigmaFileKeySearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.IssueList, client.IssueSummary]{
		Use:       "figma-file-key-search FILE_KEY",
		Short:     "Search issues linked to a Figma file key",
		LimitHelp: "issues",
		Args:      cobra.ExactArgs(1),
		Load:      loadIssueFigmaFileKeySearch,
		WriteItem: writeIssue,
	})
}

func loadIssueFigmaFileKeySearch(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, error) {
	issues, err := client.SearchIssuesByFigmaFileKey(ctx, runtime.graphqlClient, args[0], limit)

	return issues, err
}

func addIssuePriorityValuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
		Use:   "priority-values",
		Short: "List the issue priority values",
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
	addCommandWithSafety(root, CommandSafetyRead, command)
}

func addIssueTitleSuggestionCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
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
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
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
