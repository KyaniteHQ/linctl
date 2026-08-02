package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

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
	Attachments       childListFetcher[client.AttachmentList]
	BotActor          issueChildValueFetcher[client.IssueBotActor]
	Children          childListFetcher[client.IssueList]
	Documents         childListFetcher[client.DocumentList]
	FormerAttachments childListFetcher[client.AttachmentList]
	FormerNeeds       childListFetcher[client.IssueCustomerNeedMetadataList]
	History           childListFetcher[client.IssueHistoryList]
	InverseRelations  childListFetcher[client.IssueRelationList]
	Labels            childListFetcher[client.LabelList]
	Needs             childListFetcher[client.IssueCustomerNeedMetadataList]
	Relations         childListFetcher[client.IssueRelationList]
	Releases          childListFetcher[client.ReleaseList]
	SharedAccess      issueChildValueFetcher[client.IssueSharedAccessSummary]
	StateHistory      childListFetcher[client.IssueStateHistoryList]
	Subscribers       childListFetcher[client.UserList]
}

func addIssueChildCommands(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	bundle issueChildCommandBundle,
) {
	addIssueChildListCommand(
		ctx, root, options, "attachments", bundle.Argument, bundle.Text.Attachments, "attachments",
		bundle.Attachments, writeAttachment,
	)
	addIssueChildValueCommand(
		ctx, root, options, "bot-actor", bundle.Argument, bundle.Text.BotActor,
		bundle.BotActor, writeIssueBotActor,
	)
	addIssueChildListCommand(
		ctx, root, options, "children", bundle.Argument, bundle.Text.Children, "child issues",
		bundle.Children, writeIssue,
	)
	addIssueChildListCommand(
		ctx, root, options, "documents", bundle.Argument, bundle.Text.Documents, "documents",
		bundle.Documents, writeDocument,
	)
	addIssueChildListCommand(
		ctx, root, options, "former-attachments", bundle.Argument, bundle.Text.FormerAttachments, "former attachments",
		bundle.FormerAttachments, writeAttachment,
	)
	addIssueChildListCommand(
		ctx, root, options, "former-needs", bundle.Argument, bundle.Text.FormerNeeds, "former customer needs",
		bundle.FormerNeeds, writeCustomerNeedMetadata,
	)
	addIssueChildListCommand(
		ctx, root, options, "history", bundle.Argument, bundle.Text.History, "history entries",
		bundle.History, writeIssueHistory,
	)
	addIssueChildListCommand(
		ctx, root, options, "inverse-relations", bundle.Argument, bundle.Text.InverseRelations, "inverse relations",
		bundle.InverseRelations, writeIssueRelation,
	)
	addIssueChildListCommand(
		ctx, root, options, "labels", bundle.Argument, bundle.Text.Labels, "labels",
		bundle.Labels, writeLabel,
	)
	addIssueChildListCommand(
		ctx, root, options, "needs", bundle.Argument, bundle.Text.Needs, "customer needs",
		bundle.Needs, writeCustomerNeedMetadata,
	)
	addIssueChildListCommand(
		ctx, root, options, "relations", bundle.Argument, bundle.Text.Relations, "relations",
		bundle.Relations, writeIssueRelation,
	)
	addIssueChildListCommand(
		ctx, root, options, "releases", bundle.Argument, bundle.Text.Releases, "releases",
		bundle.Releases, writeRelease,
	)
	addIssueChildValueCommand(
		ctx, root, options, "shared-access", bundle.Argument, bundle.Text.SharedAccess,
		bundle.SharedAccess, writeIssueSharedAccess,
	)
	addIssueChildListCommand(
		ctx, root, options, "state-history", bundle.Argument, bundle.Text.StateHistory, "state spans",
		bundle.StateHistory, writeIssueStateSpan,
	)
	addIssueChildListCommand(
		ctx, root, options, "subscribers", bundle.Argument, bundle.Text.Subscribers, "subscribers",
		bundle.Subscribers, writeUser,
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
	fetch childListFetcher[Page],
	writeItem readListItemWriter[Item],
) {
	addChildListCommand(ctx, root, options, name+" "+argument, short, limitHelp, fetch, writeItem)
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
	addReadGetCommand(ctx, root, options, readGetSpec[Value]{
		Use:   name + " " + argument,
		Short: short,
		Load: func(ctx context.Context, runtime commandRuntime, id string) (Value, error) {
			return load(ctx, runtime.graphqlClient, id)
		},
		Write: writeValue,
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
