package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCycleCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	cycleCommand := &cobra.Command{
		Use:   "cycle",
		Short: "Read and write Linear Cycles",
	}
	addCycleListCommand(ctx, cycleCommand, options)
	addCycleGetCommand(ctx, cycleCommand, options)
	addCycleCreateCommand(ctx, cycleCommand, options)
	addCycleUpdateCommand(ctx, cycleCommand, options)
	addCycleArchiveCommand(ctx, cycleCommand, options)
	addDomainUsageCommand(cycleCommand, options, "cycle")
	root.AddCommand(cycleCommand)
}

func addCycleListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "list",
		Short: "List Cycles for the resolved team",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runCycleListCommand(ctx, command, options, limit)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum Cycles to return")
	root.AddCommand(command)
}

func runCycleListCommand(ctx context.Context, command *cobra.Command, options *rootOptions, limit int) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	target, err := runtime.resolveTarget(ctx)
	if err != nil {
		return err
	}
	cycles, err := client.ListCyclesByTeam(ctx, runtime.graphqlClient, target.Team.ID, limit)
	if err != nil {
		return err
	}
	if err := ensureNonEmpty(options, len(cycles.Cycles)); err != nil {
		return err
	}
	cycles.Cycles, err = sortByJSONField(cycles.Cycles, options.sortField, options.sortOrder)
	if err != nil {
		return err
	}
	if options.json {
		return writeJSONValue(command, options, cycles)
	}
	for _, cycle := range cycles.Cycles {
		if err := writeCycle(command, options, cycle); err != nil {
			return err
		}
	}

	return nil
}

func addCycleGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get CYCLE_ID",
		Short: "Get one Cycle by id or slug",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			cycle, err := client.GetCycleByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCycle(command, options, cycle)
		},
	})
}

func writeCycle(command *cobra.Command, options *rootOptions, cycle client.CycleSummary) error {
	if wrote, err := writeIDOnly(command, options, cycle.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, cycle)
	}

	format, err := normalizedHumanFormat(options)
	if err != nil {
		return err
	}
	if format == "minimal" {
		return render.WriteLine(command.OutOrStdout(), "%s", cycle.ID)
	}
	if format == "full" {
		return render.WriteLine(
			command.OutOrStdout(),
			"%s %s [%s] starts_at=%s ends_at=%s progress=%0.2f",
			cycle.ID,
			cycle.Name,
			cycle.Status,
			cycle.StartsAt,
			cycle.EndsAt,
			cycle.Progress,
		)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", cycle.ID, cycle.Name, cycle.Status)
}
