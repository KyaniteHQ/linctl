package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type issueChildListFetcher[Page any] func(
	context.Context,
	graphql.Client,
	string,
	int,
) (Page, error)

type issueChildValueFetcher[Value any] func(context.Context, graphql.Client, string) (Value, error)

type issueChildCommandText struct {
	Attachments       string
	BotActor          string
	Children          string
	Documents         string
	FormerAttachments string
	FormerNeeds       string
	History           string
	InverseRelations  string
	Labels            string
	Needs             string
	Relations         string
	Releases          string
	SharedAccess      string
	StateHistory      string
	Subscribers       string
}

type issueChildCommandBundle struct {
	Argument          string
	Text              issueChildCommandText
	Attachments       issueChildListFetcher[client.AttachmentList]
	BotActor          issueChildValueFetcher[client.IssueBotActor]
	Children          issueChildListFetcher[client.IssueList]
	Documents         issueChildListFetcher[client.DocumentList]
	FormerAttachments issueChildListFetcher[client.AttachmentList]
	FormerNeeds       issueChildListFetcher[client.IssueCustomerNeedMetadataList]
	History           issueChildListFetcher[client.IssueHistoryList]
	InverseRelations  issueChildListFetcher[client.IssueRelationList]
	Labels            issueChildListFetcher[client.LabelList]
	Needs             issueChildListFetcher[client.IssueCustomerNeedMetadataList]
	Relations         issueChildListFetcher[client.IssueRelationList]
	Releases          issueChildListFetcher[client.ReleaseList]
	SharedAccess      issueChildValueFetcher[client.IssueSharedAccessSummary]
	StateHistory      issueChildListFetcher[client.IssueStateHistoryList]
	Subscribers       issueChildListFetcher[client.UserList]
}

func addIssueChildCommands(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	bundle issueChildCommandBundle,
) {
	addIssueChildListCommand(
		ctx, root, options, "attachments", bundle.Argument, bundle.Text.Attachments, "attachments",
		bundle.Attachments, attachmentListItems, writeAttachment,
	)
	addIssueChildValueCommand(
		ctx, root, options, "bot-actor", bundle.Argument, bundle.Text.BotActor,
		bundle.BotActor, writeIssueBotActor,
	)
	addIssueChildListCommand(
		ctx, root, options, "children", bundle.Argument, bundle.Text.Children, "child issues",
		bundle.Children, issueListItems, writeIssue,
	)
	addIssueChildListCommand(
		ctx, root, options, "documents", bundle.Argument, bundle.Text.Documents, "documents",
		bundle.Documents, documentListItems, writeDocument,
	)
	addIssueChildListCommand(
		ctx, root, options, "former-attachments", bundle.Argument, bundle.Text.FormerAttachments, "former attachments",
		bundle.FormerAttachments, attachmentListItems, writeAttachment,
	)
	addIssueChildListCommand(
		ctx, root, options, "former-needs", bundle.Argument, bundle.Text.FormerNeeds, "former customer needs",
		bundle.FormerNeeds, customerNeedMetadataListItems, writeCustomerNeedMetadata,
	)
	addIssueChildListCommand(
		ctx, root, options, "history", bundle.Argument, bundle.Text.History, "history entries",
		bundle.History, issueHistoryListItems, writeIssueHistory,
	)
	addIssueChildListCommand(
		ctx, root, options, "inverse-relations", bundle.Argument, bundle.Text.InverseRelations, "inverse relations",
		bundle.InverseRelations, issueRelationListItems, writeIssueRelation,
	)
	addIssueChildListCommand(
		ctx, root, options, "labels", bundle.Argument, bundle.Text.Labels, "labels",
		bundle.Labels, labelListItems, writeLabel,
	)
	addIssueChildListCommand(
		ctx, root, options, "needs", bundle.Argument, bundle.Text.Needs, "customer needs",
		bundle.Needs, customerNeedMetadataListItems, writeCustomerNeedMetadata,
	)
	addIssueChildListCommand(
		ctx, root, options, "relations", bundle.Argument, bundle.Text.Relations, "relations",
		bundle.Relations, issueRelationListItems, writeIssueRelation,
	)
	addIssueChildListCommand(
		ctx, root, options, "releases", bundle.Argument, bundle.Text.Releases, "releases",
		bundle.Releases, releaseListItems, writeRelease,
	)
	addIssueChildValueCommand(
		ctx, root, options, "shared-access", bundle.Argument, bundle.Text.SharedAccess,
		bundle.SharedAccess, writeIssueSharedAccess,
	)
	addIssueChildListCommand(
		ctx, root, options, "state-history", bundle.Argument, bundle.Text.StateHistory, "state spans",
		bundle.StateHistory, issueStateHistoryListItems, writeIssueStateSpan,
	)
	addIssueChildListCommand(
		ctx, root, options, "subscribers", bundle.Argument, bundle.Text.Subscribers, "subscribers",
		bundle.Subscribers, userListItems, writeUser,
	)
}

func addIssueChildListCommand[Page any, Item any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	name string,
	argument string,
	short string,
	limitHelp string,
	fetch issueChildListFetcher[Page],
	items func(Page) []Item,
	writeItem readListItemWriter[Item],
) {
	addListCommand(ctx, root, options, listCommandSpec[Page, Item]{
		Use:       name + " " + argument,
		Short:     short,
		LimitHelp: limitHelp,
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (Page, []Item, error) {
			page, err := fetch(ctx, runtime.graphqlClient, args[0], limit)
			return page, items(page), err
		},
		WriteItem: writeItem,
	})
}

func addIssueChildValueCommand[Value any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	name string,
	argument string,
	short string,
	load issueChildValueFetcher[Value],
	writeValue readListItemWriter[Value],
) {
	root.AddCommand(&cobra.Command{
		Use:   name + " " + argument,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			value, err := load(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeValue(command, options, value)
		},
	})
}

func directIssueChildCommandText() issueChildCommandText {
	return issueChildCommandText{
		Attachments:       "List issue attachments",
		BotActor:          "Show issue bot actor metadata",
		Children:          "List issue children",
		Documents:         "List issue documents",
		FormerAttachments: "List former issue attachments",
		FormerNeeds:       "List body-free former issue customer needs",
		History:           "List issue history metadata",
		InverseRelations:  "List issue inverse relations",
		Labels:            "List issue labels",
		Needs:             "List body-free issue customer needs",
		Relations:         "List issue relations",
		Releases:          "List issue releases",
		SharedAccess:      "Show issue shared-access metadata",
		StateHistory:      "List issue workflow state history",
		Subscribers:       "List issue subscribers",
	}
}

func scopedIssueChildCommandText(scope string) issueChildCommandText {
	return issueChildCommandText{
		Attachments:       "List attachments for " + scope,
		BotActor:          "Show bot actor metadata for " + scope,
		Children:          "List child issues for " + scope,
		Documents:         "List documents for " + scope,
		FormerAttachments: "List former attachments for " + scope,
		FormerNeeds:       "List body-free former customer needs for " + scope,
		History:           "List history metadata for " + scope,
		InverseRelations:  "List inverse relations for " + scope,
		Labels:            "List labels for " + scope,
		Needs:             "List body-free customer needs for " + scope,
		Relations:         "List relations for " + scope,
		Releases:          "List releases for " + scope,
		SharedAccess:      "Show shared-access metadata for " + scope,
		StateHistory:      "List workflow state history for " + scope,
		Subscribers:       "List subscribers for " + scope,
	}
}

func attachmentListItems(list client.AttachmentList) []client.AttachmentSummary {
	return list.Attachments
}

func issueListItems(list client.IssueList) []client.IssueSummary {
	return list.Issues
}

func documentListItems(list client.DocumentList) []client.DocumentSummary {
	return list.Documents
}

func customerNeedMetadataListItems(
	list client.IssueCustomerNeedMetadataList,
) []client.CustomerNeedMetadataSummary {
	return list.Needs
}

func issueHistoryListItems(list client.IssueHistoryList) []client.IssueHistorySummary {
	return list.History
}

func issueRelationListItems(list client.IssueRelationList) []client.IssueRelationSummary {
	return list.Relations
}

func labelListItems(list client.LabelList) []client.LabelSummary {
	return list.Labels
}

func releaseListItems(list client.ReleaseList) []client.ReleaseSummary {
	return list.Releases
}

func issueStateHistoryListItems(list client.IssueStateHistoryList) []client.IssueStateSpanSummary {
	return list.Spans
}

func userListItems(list client.UserList) []client.UserSummary {
	return list.Users
}
