//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addInitiativeCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.InitiativeList, client.InitiativeSummary]{
			Use:           "initiative",
			Short:         "Read Linear initiatives",
			ListShort:     "List visible initiatives",
			LimitHelp:     "maximum initiatives to return",
			GetUse:        "get INITIATIVE_ID",
			GetShort:      "Get one initiative by id or slug",
			LoadList:      loadInitiativeList,
			PageWithItems: initiativePageWithItems,
			LoadGet:       loadInitiative,
			WriteItem:     writeInitiative,
		},
	)
	addInitiativeHistoryCommand(ctx, command, options)
	addInitiativeLinksCommand(ctx, command, options)
	addSubInitiativesCommand(ctx, command, options)
	addInitiativeScopedUpdatesCommand(ctx, command, options)
	addInitiativeDocumentsCommand(ctx, command, options)
	addInitiativeProjectsCommand(ctx, command, options)
}

func addInitiativeHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "history INITIATIVE_ID",
		Short: "List history records associated with one Linear initiative",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadInitiativeHistory,
				initiativeHistoryPageWithItems,
				writeInitiativeHistory,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum history records to return")
	root.AddCommand(command)
}

func addInitiativeLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "links INITIATIVE_ID",
		Short: "List external links associated with one Linear initiative",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadInitiativeLinks,
				releaseLinksPageWithItems,
				writeEntityExternalLink,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum links to return")
	root.AddCommand(command)
}

func addSubInitiativesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "sub-initiatives INITIATIVE_ID",
		Short: "List sub-initiatives associated with one Linear initiative",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadSubInitiatives,
				initiativePageWithItems,
				writeInitiative,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum sub-initiatives to return")
	root.AddCommand(command)
}

func addInitiativeScopedUpdatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "updates INITIATIVE_ID",
		Short: "List status updates associated with one Linear initiative",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadInitiativeScopedUpdates,
				initiativeUpdatePageWithItems,
				writeInitiativeUpdate,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum initiative updates to return")
	root.AddCommand(command)
}

func addInitiativeDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "documents INITIATIVE_ID",
		Short: "List documents associated with one Linear initiative",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadInitiativeDocuments,
				documentPageWithItems,
				writeDocument,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum documents to return")
	root.AddCommand(command)
}

func addInitiativeProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "projects INITIATIVE_ID",
		Short: "List projects directly associated with one Linear initiative",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadInitiativeProjects,
				projectPageWithItems,
				writeProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum projects to return")
	root.AddCommand(command)
}

func writeInitiative(
	command *cobra.Command,
	options *rootOptions,
	initiative client.InitiativeSummary,
) error {
	if wrote, err := writeIDOnly(command, options, initiative.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, initiative)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", initiative.ID, initiative.Name, initiative.Status)
}

func writeInitiativeHistory(
	command *cobra.Command,
	options *rootOptions,
	history client.InitiativeHistorySummary,
) error {
	if wrote, err := writeIDOnly(command, options, history.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, history)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s initiative %s entries %d",
		history.ID,
		history.InitiativeID,
		history.EntryCount,
	)
}

func loadInitiativeList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.InitiativeList, []client.InitiativeSummary, error) {
	initiatives, err := client.ListInitiatives(ctx, runtime.graphqlClient, limit)
	return initiatives, initiatives.Initiatives, err
}

func loadInitiative(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.InitiativeSummary, error) {
	return client.GetInitiativeByID(ctx, runtime.graphqlClient, id)
}

func initiativePageWithItems(
	page client.InitiativeList,
	initiatives []client.InitiativeSummary,
) client.InitiativeList {
	page.Initiatives = initiatives
	return page
}

func loadInitiativeHistory(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.InitiativeHistoryList, []client.InitiativeHistorySummary, error) {
	history, err := client.ListInitiativeHistory(ctx, runtime.graphqlClient, args[0], limit)
	return history, history.History, err
}

func initiativeHistoryPageWithItems(
	page client.InitiativeHistoryList,
	history []client.InitiativeHistorySummary,
) client.InitiativeHistoryList {
	page.History = history
	return page
}

func loadInitiativeLinks(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.EntityExternalLinkList, []client.EntityExternalLinkSummary, error) {
	links, err := client.ListInitiativeLinks(ctx, runtime.graphqlClient, args[0], limit)
	return links, links.Links, err
}

func loadSubInitiatives(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.InitiativeList, []client.InitiativeSummary, error) {
	initiatives, err := client.ListSubInitiatives(ctx, runtime.graphqlClient, args[0], limit)
	return initiatives, initiatives.Initiatives, err
}

func loadInitiativeScopedUpdates(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.InitiativeUpdateList, []client.InitiativeUpdateSummary, error) {
	updates, err := client.ListInitiativeUpdatesForInitiative(ctx, runtime.graphqlClient, args[0], limit)
	return updates, updates.Updates, err
}

func loadInitiativeDocuments(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.DocumentList, []client.DocumentSummary, error) {
	documents, err := client.ListInitiativeDocuments(ctx, runtime.graphqlClient, args[0], limit)
	return documents, documents.Documents, err
}

func loadInitiativeProjects(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ProjectList, []client.ProjectSummary, error) {
	projects, err := client.ListInitiativeProjects(ctx, runtime.graphqlClient, args[0], limit)
	return projects, projects.Projects, err
}
