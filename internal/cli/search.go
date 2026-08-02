package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := newGroupCommand("search", "Search Linear issues, projects, and documents")
	addSearchDocumentsCommand(ctx, command, options)
	addSearchIssuesCommand(ctx, command, options)
	addSearchProjectsCommand(ctx, command, options)
	root.AddCommand(command)
}

func addSearchDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.SearchDocumentList, client.SearchDocumentSummary]{
		Use:       "documents QUERY",
		Short:     "Search Linear documents by text",
		LimitHelp: "document search results",
		Limit:     20,
		Args:      cobra.ExactArgs(1),
		Load:      loadSearchDocuments,
		WriteItem: writeSearchDocument,
	})
}

func addSearchIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.SearchIssueList, client.SearchIssueSummary]{
		Use:       "issues QUERY",
		Short:     "Search Linear issues by text",
		LimitHelp: "issue search results",
		Limit:     20,
		Args:      cobra.ExactArgs(1),
		Load:      loadSearchIssues,
		WriteItem: writeSearchIssue,
	})
}

func addSearchProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.SearchProjectList, client.SearchProjectSummary]{
		Use:       "projects QUERY",
		Short:     "Search Linear projects by text",
		LimitHelp: "project search results",
		Limit:     20,
		Args:      cobra.ExactArgs(1),
		Load:      loadSearchProjects,
		WriteItem: writeSearchProject,
	})
}

func writeSearchDocument(command *cobra.Command, options *rootOptions, document client.SearchDocumentSummary) error {
	return writeItemLine(
		command, options, document, document.ID,
		"%s %s [%s]", document.ID, document.Title, emptyDash(document.ParentType),
	)
}

func writeSearchIssue(command *cobra.Command, options *rootOptions, issue client.SearchIssueSummary) error {
	return writeItemLine(
		command, options, issue, issue.ID,
		"%s %s [%s]", issue.Identifier, issue.Title, issue.StateName,
	)
}

func writeSearchProject(command *cobra.Command, options *rootOptions, project client.SearchProjectSummary) error {
	return writeItemLine(
		command, options, project, project.ID,
		"%s %s [%s]", project.ID, project.Name, project.Status.Name,
	)
}

func loadSearchDocuments(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.SearchDocumentList, error) {
	page, err := client.SearchDocuments(ctx, runtime.graphqlClient, args[0], limit)
	return page, err
}

func loadSearchIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.SearchIssueList, error) {
	page, err := client.SearchIssues(ctx, runtime.graphqlClient, args[0], limit)
	return page, err
}

func loadSearchProjects(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.SearchProjectList, error) {
	page, err := client.SearchProjects(ctx, runtime.graphqlClient, args[0], limit)
	return page, err
}
