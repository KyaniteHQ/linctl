//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addIssueToReleaseCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand[
		client.IssueToReleaseList,
		client.IssueToReleaseSummary,
	](ctx, root, options, readListGetSpec[client.IssueToReleaseList, client.IssueToReleaseSummary]{
		Use:           "issue-to-release",
		Short:         "Read Linear Issue-to-Release associations",
		ListShort:     "List visible Issue-to-Release associations",
		LimitHelp:     "maximum Issue-to-Release associations to return",
		GetUse:        "get ISSUE_TO_RELEASE_ID",
		GetShort:      "Get one Issue-to-Release association by id",
		LoadList:      loadIssueToReleaseList,
		PageWithItems: issueToReleasePageWithItems,
		LoadGet:       loadIssueToRelease,
		WriteItem:     writeIssueToRelease,
	})
}

func writeIssueToRelease(
	command *cobra.Command,
	options *rootOptions,
	association client.IssueToReleaseSummary,
) error {
	if wrote, err := writeIDOnly(command, options, association.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, association)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s issue %s -> release %s",
		association.ID,
		association.IssueID,
		association.ReleaseID,
	)
}

func loadIssueToReleaseList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.IssueToReleaseList, []client.IssueToReleaseSummary, error) {
	associations, err := client.ListIssueToReleases(ctx, runtime.graphqlClient, limit)
	return associations, associations.Associations, err
}

func loadIssueToRelease(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.IssueToReleaseSummary, error) {
	return client.GetIssueToReleaseByID(ctx, runtime.graphqlClient, id)
}

func issueToReleasePageWithItems(
	page client.IssueToReleaseList,
	associations []client.IssueToReleaseSummary,
) client.IssueToReleaseList {
	page.Associations = associations
	return page
}
