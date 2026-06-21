//nolint:dupl // Search command glue is intentionally uniform across search surfaces.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addSemanticSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 20
	command := &cobra.Command{
		Use:   "semantic-search QUERY",
		Short: "Search Linear issues, projects, initiatives, and documents semantically",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			results, err := client.SearchSemantic(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(results.Results)); err != nil {
				return err
			}
			items, err := sortByJSONField(results.Results, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			results.Results = items
			if options.json {
				return writeJSONValue(command, options, results)
			}
			for _, result := range items {
				if err := writeSemanticSearchResult(command, options, result); err != nil {
					return err
				}
			}

			return nil
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum semantic search results to return")
	root.AddCommand(command)
}

func writeSemanticSearchResult(
	command *cobra.Command,
	options *rootOptions,
	result client.SemanticSearchResultSummary,
) error {
	if wrote, err := writeIDOnly(command, options, result.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, result)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s %s",
		result.Type,
		result.ID,
		emptyDash(result.Key),
		result.Title,
	)
}
