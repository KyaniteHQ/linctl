package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addCycleCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := cycleWriteFlags{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.CycleSummary]{
		Use:       "create",
		Short:     "Create a Cycle in the pinned team",
		Args:      cobra.NoArgs,
		Configure: bindCycleWriteFlags(&flags, ""),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, _ []string,
		) (client.CycleSummary, error) {
			request := client.CycleCreateRequest{
				Name:        flags.Name,
				Description: flags.Description,
				StartsAt:    flags.StartsAt,
				EndsAt:      flags.EndsAt,
				CompletedAt: flags.CompletedAt,
			}

			return client.CreateCycle(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeCycle,
	})
}

func addCycleUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	flags := cycleWriteFlags{}
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.CycleSummary]{
		Use:       "update CYCLE_ID",
		Short:     "Update a Cycle after pinned-target comparison",
		Args:      cobra.ExactArgs(1),
		Configure: bindCycleWriteFlags(&flags, "new "),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.CycleSummary, error) {
			request := client.CycleUpdateRequest{
				ID:          args[0],
				Name:        flags.Name,
				Description: flags.Description,
				StartsAt:    flags.StartsAt,
				EndsAt:      flags.EndsAt,
				CompletedAt: flags.CompletedAt,
			}

			return client.UpdateCycle(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeCycle,
	})
}

func addCycleArchiveCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.CycleSummary]{
		Use:   "archive CYCLE_ID",
		Short: "Archive a Cycle after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.CycleSummary, error) {
			return client.ArchiveCycle(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeCycle,
	})
}

type cycleWriteFlags struct {
	Name        string
	Description string
	StartsAt    string
	EndsAt      string
	CompletedAt string
}

func bindCycleWriteFlags(flags *cycleWriteFlags, helpPrefix string) func(*cobra.Command) {
	return func(command *cobra.Command) {
		command.Flags().StringVar(&flags.Name, "name", "", helpPrefix+"Cycle name")
		command.Flags().StringVar(&flags.Description, "description", "", helpPrefix+"Cycle description")
		command.Flags().StringVar(&flags.StartsAt, "starts-at", "", helpPrefix+"Cycle start time")
		command.Flags().StringVar(&flags.EndsAt, "ends-at", "", helpPrefix+"Cycle end time")
		command.Flags().StringVar(&flags.CompletedAt, "completed-at", "", helpPrefix+"Cycle completion time")
	}
}
