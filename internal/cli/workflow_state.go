package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addWorkflowStateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	workflowStateCommand := &cobra.Command{
		Use:   "workflow-state",
		Short: "Read Linear workflow states",
	}

	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List visible workflow states",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadWorkflowStateList,
				workflowStatePageWithItems,
				writeWorkflowState,
			)
		},
	}
	listCommand.Flags().IntVar(&limit, "limit", limit, "maximum workflow states to return")
	getCommand := &cobra.Command{
		Use:   "get WORKFLOW_STATE_ID",
		Short: "Get one workflow state by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			state, err := client.GetWorkflowStateByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeWorkflowState(command, options, state)
		},
	}
	workflowStateCommand.AddCommand(listCommand, getCommand)
	root.AddCommand(workflowStateCommand)
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

func workflowStatePageWithItems(
	page client.WorkflowStateList,
	states []client.WorkflowStateSummary,
) client.WorkflowStateList {
	page.WorkflowStates = states
	return page
}
