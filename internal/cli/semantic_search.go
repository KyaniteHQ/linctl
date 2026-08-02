package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addSemanticSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.SemanticSearchList, client.SemanticSearchResultSummary]{
		Use:       "semantic-search QUERY",
		Short:     "Search Linear issues, projects, initiatives, and documents by meaning",
		LimitHelp: "semantic search results",
		Limit:     20,
		Args:      cobra.ExactArgs(1),
		Load:      loadSemanticSearch,
		WriteItem: writeSemanticSearchResult,
	})
}

func writeSemanticSearchResult(
	command *cobra.Command,
	options *rootOptions,
	result client.SemanticSearchResultSummary,
) error {
	return writeItemLine(
		command, options, result, result.ID,
		"%s %s %s %s", result.Type, result.ID, emptyDash(result.Key), result.Title,
	)
}

func loadSemanticSearch(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.SemanticSearchList, error) {
	page, err := client.SearchSemantic(ctx, runtime.graphqlClient, args[0], limit)
	return page, err
}
