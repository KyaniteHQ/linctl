package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := newGroupCommand("search", "Search Linear issues, projects, and documents")
	addSearchDocumentsCommand(ctx, command, options)
	addSearchIssuesCommand(ctx, command, options)
	addSearchProjectsCommand(ctx, command, options)
	root.AddCommand(command)
}

func addSearchDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 20
	command := &cobra.Command{
		Use:   "documents QUERY",
		Short: "Search Linear documents by text",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadSearchDocuments,
				writeSearchDocument,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum document search results to return")
	root.AddCommand(preflightReadListCommand(command, loadSearchDocuments))
}

func addSearchIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 20
	command := &cobra.Command{
		Use:   "issues QUERY",
		Short: "Search Linear issues by text",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadSearchIssues,
				writeSearchIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issue search results to return")
	root.AddCommand(preflightReadListCommand(command, loadSearchIssues))
}

func addSearchProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 20
	command := &cobra.Command{
		Use:   "projects QUERY",
		Short: "Search Linear projects by text",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadSearchProjects,
				writeSearchProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum project search results to return")
	root.AddCommand(preflightReadListCommand(command, loadSearchProjects))
}

func writeSearchDocument(command *cobra.Command, options *rootOptions, document client.SearchDocumentSummary) error {
	return writeItem(command, options, document, document.ID,
		func(command *cobra.Command, _ *rootOptions, document client.SearchDocumentSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s [%s]",
				document.ID,
				document.Title,
				emptyDash(document.ParentType),
			)
		})
}

func writeSearchIssue(command *cobra.Command, options *rootOptions, issue client.SearchIssueSummary) error {
	return writeItem(command, options, issue, issue.ID,
		func(command *cobra.Command, _ *rootOptions, issue client.SearchIssueSummary) error {
			return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", issue.Identifier, issue.Title, issue.StateName)
		})
}

func writeSearchProject(command *cobra.Command, options *rootOptions, project client.SearchProjectSummary) error {
	return writeItem(command, options, project, project.ID,
		func(command *cobra.Command, _ *rootOptions, project client.SearchProjectSummary) error {
			return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", project.ID, project.Name, project.Status.Name)
		})
}

func loadSearchDocuments(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.SearchDocumentList, []client.SearchDocumentSummary, error) {
	page, err := client.SearchDocuments(ctx, runtime.graphqlClient, args[0], limit)
	return page, page.Documents, err
}

func loadSearchIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.SearchIssueList, []client.SearchIssueSummary, error) {
	page, err := client.SearchIssues(ctx, runtime.graphqlClient, args[0], limit)
	return page, page.Issues, err
}

func loadSearchProjects(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.SearchProjectList, []client.SearchProjectSummary, error) {
	page, err := client.SearchProjects(ctx, runtime.graphqlClient, args[0], limit)
	return page, page.Projects, err
}
