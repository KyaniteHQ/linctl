//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addReleaseCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ReleaseList, client.ReleaseSummary]{
			Use:           "release",
			Short:         "Read Linear releases",
			ListShort:     "List visible Linear releases",
			LimitHelp:     "maximum releases to return",
			GetUse:        "get RELEASE_ID",
			GetShort:      "Get one release by id",
			LoadList:      loadReleaseList,
			PageWithItems: releasePageWithItems,
			LoadGet:       loadRelease,
			WriteItem:     writeRelease,
		},
	)
	addReleaseSearchCommand(ctx, command, options)
	addReleaseHistoryCommand(ctx, command, options)
	addReleaseLinksCommand(ctx, command, options)
}

func addReleaseSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	var limit int
	command := &cobra.Command{
		Use:   "search TERM",
		Short: "Search Linear releases",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			releases, err := client.SearchReleases(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			items := releases.Releases
			items, err = sortByJSONField(items, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(items)); err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, releasePageWithItems(releases, items))
			}
			for _, release := range items {
				if err := writeRelease(command, options, release); err != nil {
					return err
				}
			}
			return nil
		},
	}
	command.Flags().IntVar(&limit, "limit", 20, "maximum releases to return")
	root.AddCommand(command)
}

func addReleaseHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "history RELEASE_ID",
		Short: "List history records associated with one Linear release",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadReleaseHistory,
				releaseHistoryPageWithItems,
				writeReleaseHistory,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum history records to return")
	root.AddCommand(command)
}

func addReleaseLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "links RELEASE_ID",
		Short: "List external links associated with one Linear release",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadReleaseLinks,
				releaseLinksPageWithItems,
				writeEntityExternalLink,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum links to return")
	root.AddCommand(command)
}

func addReleaseNoteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ReleaseNoteList, client.ReleaseNoteSummary]{
			Use:           "release-note",
			Short:         "Read Linear release notes",
			ListShort:     "List visible Linear release notes",
			LimitHelp:     "maximum release notes to return",
			GetUse:        "get RELEASE_NOTE_ID",
			GetShort:      "Get one release note by id",
			LoadList:      loadReleaseNoteList,
			PageWithItems: releaseNotePageWithItems,
			LoadGet:       loadReleaseNote,
			WriteItem:     writeReleaseNote,
		},
	)
}

func writeRelease(command *cobra.Command, options *rootOptions, release client.ReleaseSummary) error {
	if wrote, err := writeIDOnly(command, options, release.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, release)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s [%s] pipeline %s stage %s issues %d",
		release.ID,
		release.Name,
		emptyDash(release.Version),
		release.PipelineName,
		release.StageName,
		release.IssueCount,
	)
}

func writeReleaseHistory(command *cobra.Command, options *rootOptions, history client.ReleaseHistorySummary) error {
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
		"%s release %s entries %d",
		history.ID,
		history.ReleaseID,
		history.EntryCount,
	)
}

func writeEntityExternalLink(
	command *cobra.Command,
	options *rootOptions,
	link client.EntityExternalLinkSummary,
) error {
	if wrote, err := writeIDOnly(command, options, link.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, link)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s order %g",
		link.ID,
		link.Label,
		link.URL,
		link.SortOrder,
	)
}

func writeReleaseNote(command *cobra.Command, options *rootOptions, note client.ReleaseNoteSummary) error {
	if wrote, err := writeIDOnly(command, options, note.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, note)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s pipeline %s releases %d",
		note.ID,
		emptyDash(note.Title),
		note.PipelineName,
		note.ReleaseCount,
	)
}

func loadReleaseList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ReleaseList, []client.ReleaseSummary, error) {
	releases, err := client.ListReleases(ctx, runtime.graphqlClient, limit)
	return releases, releases.Releases, err
}

func loadRelease(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ReleaseSummary, error) {
	return client.GetReleaseByID(ctx, runtime.graphqlClient, id)
}

func releasePageWithItems(page client.ReleaseList, releases []client.ReleaseSummary) client.ReleaseList {
	page.Releases = releases
	return page
}

func loadReleaseHistory(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ReleaseHistoryList, []client.ReleaseHistorySummary, error) {
	history, err := client.ListReleaseHistory(ctx, runtime.graphqlClient, args[0], limit)
	return history, history.History, err
}

func releaseHistoryPageWithItems(
	page client.ReleaseHistoryList,
	history []client.ReleaseHistorySummary,
) client.ReleaseHistoryList {
	page.History = history
	return page
}

func loadReleaseLinks(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.EntityExternalLinkList, []client.EntityExternalLinkSummary, error) {
	links, err := client.ListReleaseLinks(ctx, runtime.graphqlClient, args[0], limit)
	return links, links.Links, err
}

func releaseLinksPageWithItems(
	page client.EntityExternalLinkList,
	links []client.EntityExternalLinkSummary,
) client.EntityExternalLinkList {
	page.Links = links
	return page
}

func loadReleaseNoteList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ReleaseNoteList, []client.ReleaseNoteSummary, error) {
	notes, err := client.ListReleaseNotes(ctx, runtime.graphqlClient, limit)
	return notes, notes.ReleaseNotes, err
}

func loadReleaseNote(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.ReleaseNoteSummary, error) {
	return client.GetReleaseNoteByID(ctx, runtime.graphqlClient, id)
}

func releaseNotePageWithItems(
	page client.ReleaseNoteList,
	notes []client.ReleaseNoteSummary,
) client.ReleaseNoteList {
	page.ReleaseNotes = notes
	return page
}
