package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addWorkflowStateCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := workflowStateWriteFlags{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.WorkflowStateSummary]{
		Use:   "create",
		Short: "Create a WorkflowState in the pinned team",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			bindWorkflowStateCreateFlags(command, &flags)
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, _ []string,
		) (client.WorkflowStateSummary, error) {
			return client.CreateWorkflowState(
				ctx, runtime.graphqlClient, runtime.config.Target, workflowStateCreateRequest(command, flags),
			)
		},
		Write: writeWorkflowState,
	})
}

func addWorkflowStateUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := workflowStateWriteFlags{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.WorkflowStateSummary]{
		Use:   "update WORKFLOW_STATE_ID",
		Short: "Update a WorkflowState after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			bindWorkflowStateUpdateFlags(command, &flags)
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.WorkflowStateSummary, error) {
			return client.UpdateWorkflowState(
				ctx, runtime.graphqlClient, runtime.config.Target, workflowStateUpdateRequest(command, flags, args[0]),
			)
		},
		Write: writeWorkflowState,
	})
}

type workflowStateWriteFlags struct {
	ID          string
	Name        string
	Type        string
	Color       string
	Description string
	Position    float64
}

func bindWorkflowStateCreateFlags(command *cobra.Command, flags *workflowStateWriteFlags) {
	command.Flags().StringVar(&flags.ID, "id", "", "caller-supplied WorkflowState UUID v4")
	command.Flags().StringVar(&flags.Name, "name", "", "WorkflowState name")
	command.Flags().StringVar(&flags.Type, "type", "", "WorkflowState type")
	command.Flags().StringVar(&flags.Color, "color", "", "WorkflowState color")
	bindWorkflowStateOptionalFlags(command, flags, "")
}

func bindWorkflowStateUpdateFlags(command *cobra.Command, flags *workflowStateWriteFlags) {
	command.Flags().StringVar(&flags.Name, "name", "", "new WorkflowState name")
	command.Flags().StringVar(&flags.Color, "color", "", "new WorkflowState color")
	bindWorkflowStateOptionalFlags(command, flags, "new ")
}

func bindWorkflowStateOptionalFlags(command *cobra.Command, flags *workflowStateWriteFlags, helpPrefix string) {
	command.Flags().StringVar(&flags.Description, "description", "", helpPrefix+"WorkflowState description")
	command.Flags().Float64Var(&flags.Position, "position", 0, helpPrefix+"WorkflowState position")
}

func workflowStateCreateRequest(
	command *cobra.Command,
	flags workflowStateWriteFlags,
) client.WorkflowStateCreateRequest {
	request := client.WorkflowStateCreateRequest{
		ID:    flags.ID,
		Name:  flags.Name,
		Type:  flags.Type,
		Color: flags.Color,
	}
	if command.Flags().Changed("description") {
		request.Description = &flags.Description
	}
	if command.Flags().Changed("position") {
		request.Position = &flags.Position
	}

	return request
}

func workflowStateUpdateRequest(
	command *cobra.Command,
	flags workflowStateWriteFlags,
	id string,
) client.WorkflowStateUpdateRequest {
	request := client.WorkflowStateUpdateRequest{ID: id}
	if command.Flags().Changed("name") {
		request.Name = &flags.Name
	}
	if command.Flags().Changed("color") {
		request.Color = &flags.Color
	}
	if command.Flags().Changed("description") {
		request.Description = &flags.Description
	}
	if command.Flags().Changed("position") {
		request.Position = &flags.Position
	}

	return request
}
