package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addProjectMilestoneCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectMilestoneWriteCommand(ctx, root, options, projectMilestoneWriteSpec{
		Use:             "create PROJECT_ID",
		Short:           "Create a ProjectMilestone in a pinned project",
		NameHelp:        "ProjectMilestone name",
		DescriptionHelp: "ProjectMilestone description",
		TargetDateHelp:  "ProjectMilestone target date",
		Run: func(ctx context.Context, runtime commandRuntime, id string, flags projectMilestoneWriteFlags) (
			client.ProjectMilestoneSummary,
			error,
		) {
			request := client.ProjectMilestoneCreateRequest{
				ProjectID:   id,
				Name:        flags.Name,
				Description: flags.Description,
				TargetDate:  flags.TargetDate,
			}

			return client.CreateProjectMilestone(ctx, runtime.graphqlClient, runtime.config.Target, request)
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
		Run: func(ctx context.Context, runtime commandRuntime, id string, flags projectMilestoneWriteFlags) (
			client.ProjectMilestoneSummary,
			error,
		) {
			request := client.ProjectMilestoneUpdateRequest{
				ID:          id,
				Name:        flags.Name,
				Description: flags.Description,
				TargetDate:  flags.TargetDate,
			}

			return client.UpdateProjectMilestone(ctx, runtime.graphqlClient, runtime.config.Target, request)
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
		commandRuntime,
		string,
		projectMilestoneWriteFlags,
	) (client.ProjectMilestoneSummary, error)
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
			return runProjectMilestoneWriteCommand(ctx, command, options, func(runtime commandRuntime) (
				client.ProjectMilestoneSummary,
				error,
			) {
				return spec.Run(ctx, runtime, args[0], flags)
			})
		},
	}
	command.Flags().StringVar(&flags.Name, "name", "", spec.NameHelp)
	command.Flags().StringVar(&flags.Description, "description", "", spec.DescriptionHelp)
	command.Flags().StringVar(&flags.TargetDate, "target-date", "", spec.TargetDateHelp)
	root.AddCommand(command)
}

func runProjectMilestoneWriteCommand(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	write func(commandRuntime) (client.ProjectMilestoneSummary, error),
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	milestone, err := write(runtime)
	if err != nil {
		return err
	}

	return writeProjectMilestone(command, options, milestone)
}
