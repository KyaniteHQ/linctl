package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addSprintCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	sprintCommand := &cobra.Command{
		Use:   "sprint",
		Short: "Read Linear Cycle sprint reports",
	}
	addSprintCurrentCommand(ctx, sprintCommand, options)
	addSprintReportCommand(ctx, sprintCommand, options)
	root.AddCommand(sprintCommand)
}

func addSprintCurrentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "current",
		Short: "Show the active Cycle for the resolved team",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			target, err := runtime.resolveTarget(ctx)
			if err != nil {
				return err
			}
			cycle, err := client.CurrentCycleByTeam(ctx, runtime.graphqlClient, target.Team.ID)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, cycle)
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s [%s]",
				cycle.ID,
				cycle.Name,
				cycle.Status,
			)
		},
	})
}

func addSprintReportCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "report CYCLE_ID",
		Short: "Show one Cycle with assigned issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			report, err := client.GetSprintReport(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(report.Issues)); err != nil {
				return err
			}
			report.Issues, err = sortByJSONField(report.Issues, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, report)
			}
			if err := render.WriteLine(
				command.OutOrStdout(),
				"%s %s [%s]",
				report.Cycle.ID,
				report.Cycle.Name,
				report.Cycle.Status,
			); err != nil {
				return err
			}
			for _, issue := range report.Issues {
				if err := render.WriteLine(
					command.OutOrStdout(),
					"%s %s [%s]",
					issue.Identifier,
					issue.Title,
					issue.State,
				); err != nil {
					return err
				}
			}

			return nil
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to include")
	root.AddCommand(command)
}
