//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addWorkflowStateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	workflowStateCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.WorkflowStateList, client.WorkflowStateSummary]{
			Use:           "workflow-state",
			Short:         "Read Linear workflow states",
			ListShort:     "List visible workflow states",
			LimitHelp:     "maximum workflow states to return",
			GetUse:        "get WORKFLOW_STATE_ID",
			GetShort:      "Get one workflow state by id",
			LoadList:      loadWorkflowStateList,
			PageWithItems: workflowStatePageWithItems,
			LoadGet:       loadWorkflowState,
			WriteItem:     writeWorkflowState,
		},
	)
	addWorkflowStateIssuesCommand(ctx, workflowStateCommand, options)
}

func addWorkflowStateIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "issues WORKFLOW_STATE_ID",
		Short: "List issues currently in one workflow state",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadWorkflowStateIssues,
				workflowStateIssuePageWithItems,
				writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func writeWorkflowState(
	command *cobra.Command,
	options *rootOptions,
	state client.WorkflowStateSummary,
) error {
	if wrote, err := writeIDOnly(command, options, state.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, state)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", state.ID, state.Name, state.Type)
}

func loadWorkflowStateList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.WorkflowStateList, []client.WorkflowStateSummary, error) {
	states, err := client.ListWorkflowStates(ctx, runtime.graphqlClient, limit)
	return states, states.WorkflowStates, err
}

func loadWorkflowState(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.WorkflowStateSummary, error) {
	return client.GetWorkflowStateByID(ctx, runtime.graphqlClient, id)
}

func workflowStatePageWithItems(
	page client.WorkflowStateList,
	states []client.WorkflowStateSummary,
) client.WorkflowStateList {
	page.WorkflowStates = states
	return page
}

func loadWorkflowStateIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.WorkflowStateIssueList, []client.IssueSummary, error) {
	issues, err := client.ListWorkflowStateIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, issues.Issues, err
}

func workflowStateIssuePageWithItems(
	page client.WorkflowStateIssueList,
	issues []client.IssueSummary,
) client.WorkflowStateIssueList {
	page.Issues = issues
	return page
}
