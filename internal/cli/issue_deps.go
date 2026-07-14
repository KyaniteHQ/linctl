package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueDepsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "deps ISSUE_ID",
		Short: "Show issue dependencies",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			dependencies, err := client.GetIssueDependencies(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, dependencies)
			}

			return writeIssueDependencies(command, options, dependencies)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum related issues per section")
	root.AddCommand(command)
}

func writeIssueDependencies(
	command *cobra.Command,
	options *rootOptions,
	dependencies client.IssueDependencyGraph,
) error {
	if options.quiet {
		return nil
	}
	if err := render.WriteLine(command.OutOrStdout(), "issue %s", dependencies.Identifier); err != nil {
		return err
	}
	if dependencies.Parent != nil {
		parent := []client.IssueSummary{*dependencies.Parent}
		if err := writeIssueDependencySection(command, options, "parent", parent); err != nil {
			return err
		}
	}
	if err := writeIssueDependencySection(command, options, "children", dependencies.Children); err != nil {
		return err
	}
	if err := writeIssueDependencySection(command, options, "blocks", dependencies.Blocks); err != nil {
		return err
	}

	return writeIssueDependencySection(command, options, "blocked_by", dependencies.BlockedBy)
}

func writeIssueDependencySection(
	command *cobra.Command,
	options *rootOptions,
	name string,
	issues []client.IssueSummary,
) error {
	if err := render.WriteLine(command.OutOrStdout(), "%s:", name); err != nil {
		return err
	}

	return writeIssues(command, options, issues)
}
