package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addCycleCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := cycleWriteFlags{}
	command := &cobra.Command{
		Use:   "create",
		Short: "Create a Cycle in the pinned team",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			request := client.CycleCreateRequest{
				Name:        flags.Name,
				Description: flags.Description,
				StartsAt:    flags.StartsAt,
				EndsAt:      flags.EndsAt,
				CompletedAt: flags.CompletedAt,
			}

			return runCycleWriteCommand(ctx, command, options, func(runtime commandRuntime) (
				client.CycleSummary,
				error,
			) {
				return client.CreateCycle(ctx, runtime.graphqlClient, runtime.config.Target, request)
			})
		},
	}
	bindCycleWriteFlags(command, &flags, "")
	root.AddCommand(command)
}

func addCycleUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := cycleWriteFlags{}
	command := &cobra.Command{
		Use:   "update CYCLE_ID",
		Short: "Update a Cycle after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			request := client.CycleUpdateRequest{
				ID:          args[0],
				Name:        flags.Name,
				Description: flags.Description,
				StartsAt:    flags.StartsAt,
				EndsAt:      flags.EndsAt,
				CompletedAt: flags.CompletedAt,
			}

			return runCycleWriteCommand(ctx, command, options, func(runtime commandRuntime) (
				client.CycleSummary,
				error,
			) {
				return client.UpdateCycle(ctx, runtime.graphqlClient, runtime.config.Target, request)
			})
		},
	}
	bindCycleWriteFlags(command, &flags, "new ")
	root.AddCommand(command)
}

func addCycleArchiveCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "archive CYCLE_ID",
		Short: "Archive a Cycle after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runCycleWriteCommand(ctx, command, options, func(runtime commandRuntime) (
				client.CycleSummary,
				error,
			) {
				return client.ArchiveCycle(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
			})
		},
	})
}

type cycleWriteFlags struct {
	Name        string
	Description string
	StartsAt    string
	EndsAt      string
	CompletedAt string
}

func bindCycleWriteFlags(command *cobra.Command, flags *cycleWriteFlags, helpPrefix string) {
	command.Flags().StringVar(&flags.Name, "name", "", helpPrefix+"Cycle name")
	command.Flags().StringVar(&flags.Description, "description", "", helpPrefix+"Cycle description")
	command.Flags().StringVar(&flags.StartsAt, "starts-at", "", helpPrefix+"Cycle start time")
	command.Flags().StringVar(&flags.EndsAt, "ends-at", "", helpPrefix+"Cycle end time")
	command.Flags().StringVar(&flags.CompletedAt, "completed-at", "", helpPrefix+"Cycle completion time")
}

func runCycleWriteCommand(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	write func(commandRuntime) (client.CycleSummary, error),
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	cycle, err := write(runtime)
	if err != nil {
		return err
	}

	return writeCycle(command, options, cycle)
}
