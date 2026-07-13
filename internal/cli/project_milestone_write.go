package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectMilestoneDeleteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "delete PROJECT_MILESTONE_ID",
		Short: "Hard delete a ProjectMilestone after pinned-target comparison; cannot be undone via linctl",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return runProjectMilestoneDelete(ctx, command, options, commandAdapterFor(runtime), args[0])
		},
	})
}

func runProjectMilestoneDelete(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	deleter projectMilestoneDeleter,
	projectMilestoneID string,
) error {
	deletedID, err := deleter.DeleteProjectMilestone(ctx, projectMilestoneID)
	if err != nil {
		return err
	}

	return writeProjectMilestoneDeletion(command, options, deletedID)
}

// writeProjectMilestoneDeletion renders the ProjectMilestone delete confirmation.
// It reuses the standard deletionResult JSON shape but overrides the human line
// with an explicit irreversibility warning: ProjectMilestone delete is linctl's
// one approved hard delete, and there is no restore path via linctl.
func writeProjectMilestoneDeletion(command *cobra.Command, options *rootOptions, id string) error {
	if wrote, err := writeIDOnly(command, options, id); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, deletionResult{ID: id, Status: "deleted"})
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"hard deleted ProjectMilestone %s: cannot be undone via linctl",
		id,
	)
}

func addProjectMilestoneCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectMilestoneWriteCommand(ctx, root, options, projectMilestoneWriteSpec{
		Use:             "create PROJECT_ID",
		Short:           "Create a ProjectMilestone in a pinned project",
		NameHelp:        "ProjectMilestone name",
		DescriptionHelp: "ProjectMilestone description",
		TargetDateHelp:  "ProjectMilestone target date",
		Run: func(
			ctx context.Context,
			command *cobra.Command,
			options *rootOptions,
			adapter commandClientAdapter,
			id string,
			flags projectMilestoneWriteFlags,
		) error {
			request := client.ProjectMilestoneCreateRequest{
				ProjectID:   id,
				Name:        flags.Name,
				Description: flags.Description,
				TargetDate:  flags.TargetDate,
			}

			return runProjectMilestoneCreate(ctx, command, options, adapter, request)
		},
	})
}

func addProjectMilestoneUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectMilestoneWriteCommand(ctx, root, options, projectMilestoneWriteSpec{
		Use:             "update PROJECT_MILESTONE_ID",
		Short:           "Update a ProjectMilestone after pinned-target comparison",
		NameHelp:        "new ProjectMilestone name",
		DescriptionHelp: "new ProjectMilestone description",
		TargetDateHelp:  "new ProjectMilestone target date",
		Run: func(
			ctx context.Context,
			command *cobra.Command,
			options *rootOptions,
			adapter commandClientAdapter,
			id string,
			flags projectMilestoneWriteFlags,
		) error {
			request := client.ProjectMilestoneUpdateRequest{
				ID:          id,
				Name:        flags.Name,
				Description: flags.Description,
				TargetDate:  flags.TargetDate,
			}

			return runProjectMilestoneUpdate(ctx, command, options, adapter, request)
		},
	})
}

type projectMilestoneWriteFlags struct {
	Name        string
	Description string
	TargetDate  string
}

type projectMilestoneWriteSpec struct {
	Use             string
	Short           string
	NameHelp        string
	DescriptionHelp string
	TargetDateHelp  string
	Run             func(
		context.Context,
		*cobra.Command,
		*rootOptions,
		commandClientAdapter,
		string,
		projectMilestoneWriteFlags,
	) error
}

func addProjectMilestoneWriteCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	spec projectMilestoneWriteSpec,
) {
	flags := projectMilestoneWriteFlags{}
	command := &cobra.Command{
		Use:   spec.Use,
		Short: spec.Short,
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}

			return spec.Run(ctx, command, options, commandAdapterFor(runtime), args[0], flags)
		},
	}
	command.Flags().StringVar(&flags.Name, "name", "", spec.NameHelp)
	command.Flags().StringVar(&flags.Description, "description", "", spec.DescriptionHelp)
	command.Flags().StringVar(&flags.TargetDate, "target-date", "", spec.TargetDateHelp)
	root.AddCommand(command)
}

func runProjectMilestoneCreate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	creator projectMilestoneCreator,
	request client.ProjectMilestoneCreateRequest,
) error {
	milestone, err := creator.CreateProjectMilestone(ctx, request)
	if err != nil {
		return err
	}

	return writeProjectMilestone(command, options, milestone)
}

func runProjectMilestoneUpdate(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	updater projectMilestoneUpdater,
	request client.ProjectMilestoneUpdateRequest,
) error {
	milestone, err := updater.UpdateProjectMilestone(ctx, request)
	if err != nil {
		return err
	}

	return writeProjectMilestone(command, options, milestone)
}
