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
	addIssueGetCommand(ctx, issueCommand, options)
	addIssueDepsCommand(ctx, issueCommand, options)
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

func writeIssues(command *cobra.Command, options *rootOptions, issues []client.IssueSummary) error {
	for _, issue := range issues {
		if err := writeIssue(command, options, issue); err != nil {
			return err
		}
	}

	return nil
}
