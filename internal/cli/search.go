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
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			page, err := client.SearchDocuments(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(page.Documents)); err != nil {
				return err
			}
			documents, err := sortByJSONField(page.Documents, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			page.Documents = documents
			if options.json {
				return writeJSONValue(command, options, page)
			}
			for _, document := range documents {
				if err := writeSearchDocument(command, options, document); err != nil {
					return err
				}
			}

			return nil
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
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			page, err := client.SearchIssues(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(page.Issues)); err != nil {
				return err
			}
			issues, err := sortByJSONField(page.Issues, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			page.Issues = issues
			if options.json {
				return writeJSONValue(command, options, page)
			}
			for _, issue := range issues {
				if err := writeSearchIssue(command, options, issue); err != nil {
					return err
				}
			}

			return nil
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
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			page, err := client.SearchProjects(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(page.Projects)); err != nil {
				return err
			}
			projects, err := sortByJSONField(page.Projects, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			page.Projects = projects
			if options.json {
				return writeJSONValue(command, options, page)
			}
			for _, project := range projects {
				if err := writeSearchProject(command, options, project); err != nil {
					return err
				}
			}

			return nil
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
