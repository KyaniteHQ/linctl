package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addSLAConfigurationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "sla-configuration",
		Short: "Read Linear SLA configurations",
	}

	listCommand := &cobra.Command{
		Use:   "list TEAM_ID_OR_KEY",
		Short: "List active SLA configurations for a team",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			configurations, err := client.ListSLAConfigurations(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(configurations.SLAConfigurations)); err != nil {
				return err
			}
			items, err := sortByJSONField(configurations.SLAConfigurations, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			configurations.SLAConfigurations = items
			if options.json {
				return writeJSONValue(command, options, configurations)
			}
			for _, configuration := range items {
				if err := writeSLAConfiguration(command, options, configuration); err != nil {
					return err
				}
			}

			return nil
		},
	}

	command.AddCommand(listCommand)
	root.AddCommand(command)
}

func writeSLAConfiguration(
	command *cobra.Command,
	options *rootOptions,
	configuration client.SLAConfigurationSummary,
) error {
	return writeItem(command, options, configuration, configuration.ID,
		func(command *cobra.Command, _ *rootOptions, configuration client.SLAConfigurationSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s sla %s type %s removes %t",
				configuration.ID,
				configuration.Name,
				slaValue(configuration.SLA),
				emptyDash(configuration.SLAType),
				configuration.RemovesSLA,
			)
		})
}

func slaValue(value float64) string {
	if value == 0 {
		return "-"
	}

	return strconv.FormatFloat(value, 'f', -1, 64)
}
