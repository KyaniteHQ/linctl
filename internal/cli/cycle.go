package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCycleCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	cycleCommand := newGroupCommand("cycle", "Read and write Linear Cycles")
	addCycleListCommand(ctx, cycleCommand, options)
	addCycleGetCommand(ctx, cycleCommand, options)
	addCycleIssuesCommand(ctx, cycleCommand, options)
	addCycleUncompletedIssuesCommand(ctx, cycleCommand, options)
	addCycleCreateCommand(ctx, cycleCommand, options)
	addCycleUpdateCommand(ctx, cycleCommand, options)
	addCycleArchiveCommand(ctx, cycleCommand, options)
	addDomainUsageCommand(cycleCommand, options, "cycle")
	root.AddCommand(cycleCommand)
}

func addCycleListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.CycleList, client.CycleSummary]{
		Use:       "list",
		Short:     "List Cycles for the resolved team",
		LimitHelp: "Cycles",
		Args:      cobra.NoArgs,
		Load:      loadCyclesByTeam,
		WriteItem: writeCycle,
	})
}

func loadCyclesByTeam(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.CycleList, []client.CycleSummary, error) {
	target, err := runtime.resolveTarget(ctx)
	if err != nil {
		return client.CycleList{}, nil, err
	}
	cycles, err := client.ListCyclesByTeam(ctx, runtime.graphqlClient, target.Team.ID, limit)

	return cycles, cycles.Cycles, err
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

func addCycleIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"issues CYCLE_ID",
		"List Issues assigned to one Cycle",
		"Issues",
		client.ListCycleIssues,
		cycleIssueListItems,
		writeIssue,
	)
}

func addCycleUncompletedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"uncompleted-issues CYCLE_ID",
		"List Issues left open when one Cycle closed",
		"Issues",
		client.ListCycleUncompletedIssuesUponClose,
		cycleIssueListItems,
		writeIssue,
	)
}

func cycleIssueListItems(list client.CycleIssueList) []client.IssueSummary {
	return list.Issues
}

func writeCycle(command *cobra.Command, options *rootOptions, cycle client.CycleSummary) error {
	return writeItem(command, options, cycle, cycle.ID, cycleHumanLine)
}

func cycleHumanLine(command *cobra.Command, options *rootOptions, cycle client.CycleSummary) error {
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
