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
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadSemanticSearch,
				semanticSearchPageWithItems,
				writeSemanticSearchResult,
			)
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
	return writeItem(command, options, result, result.ID,
		func(command *cobra.Command, _ *rootOptions, result client.SemanticSearchResultSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s %s",
				result.Type,
				result.ID,
				emptyDash(result.Key),
				result.Title,
			)
		})
}

func loadSemanticSearch(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.SemanticSearchList, []client.SemanticSearchResultSummary, error) {
	page, err := client.SearchSemantic(ctx, runtime.graphqlClient, args[0], limit)
	return page, page.Results, err
}

func semanticSearchPageWithItems(
	page client.SemanticSearchList,
	results []client.SemanticSearchResultSummary,
) client.SemanticSearchList {
	page.Results = results
	return page
}
