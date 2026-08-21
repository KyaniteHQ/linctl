//nolint:dupl // Declarative registration only; addReadListGetCommand and addChildListCommand own the behavior.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addWorkflowStateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	workflowStateCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.WorkflowStateList, client.WorkflowStateSummary]{
			Use:       "workflow-state",
			Short:     "Read and write Linear workflow states",
			ListShort: "List visible workflow states",
			LimitHelp: "maximum workflow states to return",
			GetUse:    "get WORKFLOW_STATE_ID",
			GetShort:  "Get one workflow state by id",
			LoadList:  clientList(client.ListWorkflowStates),
			LoadGet:   clientGet(client.GetWorkflowStateByID),
			WriteItem: writeWorkflowState,
		},
	)
	addWorkflowStateIssuesCommand(ctx, workflowStateCommand, options)
	addWorkflowStateCreateCommand(ctx, workflowStateCommand, options)
	addWorkflowStateUpdateCommand(ctx, workflowStateCommand, options)
}

func addWorkflowStateIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"issues WORKFLOW_STATE_ID",
		"List issues currently in one workflow state",
		"issues",
		client.ListWorkflowStateIssues,
		writeIssue,
	)
}

func writeWorkflowState(command *cobra.Command, options *rootOptions, state client.WorkflowStateSummary) error {
	return writeItemLine(command, options, state, state.ID, "%s %s [%s]", state.ID, state.Name, state.Type)
}
