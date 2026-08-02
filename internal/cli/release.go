package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addReleaseCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ReleaseList, client.ReleaseSummary]{
			Use:       "release",
			Short:     "Read Linear releases",
			ListShort: "List visible Linear releases",
			LimitHelp: "maximum releases to return",
			GetUse:    "get RELEASE_ID",
			GetShort:  "Get one release by id",
			LoadList:  clientList(client.ListReleases),
			LoadGet:   clientGet(client.GetReleaseByID),
			WriteItem: writeRelease,
		},
	)
	addReleaseSearchCommand(ctx, command, options)
	addReleaseHistoryCommand(ctx, command, options)
	addReleaseDocumentsCommand(ctx, command, options)
	addReleaseIssuesCommand(ctx, command, options)
	addReleaseLinksCommand(ctx, command, options)
}

func addReleaseSearchCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ReleaseList, client.ReleaseSummary]{
		Use:       "search TERM",
		Short:     "Search Linear releases",
		LimitHelp: "releases",
		Limit:     20,
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ReleaseList, error) {
			return client.SearchReleases(ctx, runtime.graphqlClient, args[0], limit)
		},
		WriteItem: writeRelease,
	})
}

func addReleaseHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"history RELEASE_ID",
		"List history records associated with one Linear release",
		"history records",
		client.ListReleaseHistory,
		writeReleaseHistory,
	)
}

func addReleaseDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"documents RELEASE_ID",
		"List documents associated with one Linear release",
		"documents",
		client.ListReleaseDocuments,
		writeDocument,
	)
}

func addReleaseIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"issues RELEASE_ID",
		"List issues associated with one Linear release",
		"issues",
		client.ListReleaseIssues,
		writeIssue,
	)
}

func addReleaseLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"links RELEASE_ID",
		"List external links associated with one Linear release",
		"links",
		client.ListReleaseLinks,
		writeEntityExternalLink,
	)
}

func addReleaseNoteCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.ReleaseNoteList, client.ReleaseNoteSummary]{
			Use:       "release-note",
			Short:     "Read Linear release notes",
			ListShort: "List visible Linear release notes",
			LimitHelp: "maximum release notes to return",
			GetUse:    "get RELEASE_NOTE_ID",
			GetShort:  "Get one release note by id",
			LoadList:  clientList(client.ListReleaseNotes),
			LoadGet:   clientGet(client.GetReleaseNoteByID),
			WriteItem: writeReleaseNote,
		},
	)
}

func writeRelease(command *cobra.Command, options *rootOptions, release client.ReleaseSummary) error {
	return writeItemLine(
		command, options, release, release.ID,
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
	return writeItemLine(
		command, options, history, history.ID,
		"%s release %s entries %d", history.ID, history.ReleaseID, history.EntryCount,
	)
}

func writeEntityExternalLink(
	command *cobra.Command,
	options *rootOptions,
	link client.EntityExternalLinkSummary,
) error {
	return writeItemLine(
		command, options, link, link.ID,
		"%s %s %s order %g", link.ID, link.Label, link.URL, link.SortOrder,
	)
}

func writeReleaseNote(command *cobra.Command, options *rootOptions, note client.ReleaseNoteSummary) error {
	return writeItemLine(
		command, options, note, note.ID,
		"%s %s pipeline %s releases %d", note.ID, emptyDash(note.Title), note.PipelineName, note.ReleaseCount,
	)
}
