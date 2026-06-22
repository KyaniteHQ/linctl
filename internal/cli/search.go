//nolint:dupl // Typed search command glue is intentionally uniform across result kinds.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "search",
		Short: "Search Linear issues, projects, and documents",
	}
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
				searchDocumentPageWithItems,
				writeSearchDocument,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum document search results to return")
	root.AddCommand(command)
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
				searchIssuePageWithItems,
				writeSearchIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issue search results to return")
	root.AddCommand(command)
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
				searchProjectPageWithItems,
				writeSearchProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum project search results to return")
	root.AddCommand(command)
}

func writeSearchDocument(
	command *cobra.Command,
	options *rootOptions,
	document client.SearchDocumentSummary,
) error {
	if wrote, err := writeIDOnly(command, options, document.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, document)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s [%s]",
		document.ID,
		document.Title,
		emptyDash(document.ParentType),
	)
}

func writeSearchIssue(
	command *cobra.Command,
	options *rootOptions,
	issue client.SearchIssueSummary,
) error {
	if wrote, err := writeIDOnly(command, options, issue.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, issue)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", issue.Identifier, issue.Title, issue.StateName)
}

func writeSearchProject(
	command *cobra.Command,
	options *rootOptions,
	project client.SearchProjectSummary,
) error {
	if wrote, err := writeIDOnly(command, options, project.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, project)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", project.ID, project.Name, project.Status.Name)
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

func searchDocumentPageWithItems(
	page client.SearchDocumentList,
	documents []client.SearchDocumentSummary,
) client.SearchDocumentList {
	page.Documents = documents
	return page
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

func searchIssuePageWithItems(
	page client.SearchIssueList,
	issues []client.SearchIssueSummary,
) client.SearchIssueList {
	page.Issues = issues
	return page
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

func searchProjectPageWithItems(
	page client.SearchProjectList,
	projects []client.SearchProjectSummary,
) client.SearchProjectList {
	page.Projects = projects
	return page
}
