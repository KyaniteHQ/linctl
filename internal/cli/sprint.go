package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addSprintCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	sprintCommand := newGroupCommand("sprint", "Read Linear Cycle sprint reports")
	addSprintCurrentCommand(ctx, sprintCommand, options)
	addSprintReportCommand(ctx, sprintCommand, options)
	root.AddCommand(sprintCommand)
}

func addSprintCurrentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
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

			return writeCycle(command, options, cycle)
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
			// Combined JSON keeps Cycle and issues in one object so --fields
			// still projects the issues collection.
			if options.json && !options.idOnly && !options.quiet {
				return writeJSONValue(command, options, report)
			}
			if err := writeCycle(command, options, report.Cycle); err != nil {
				return err
			}

			return writeIssues(command, options, report.Issues)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to include")
	annotateReadCollectionCommand(command, collectionKeyForPage[client.SprintReport]())
	root.AddCommand(command)
}
