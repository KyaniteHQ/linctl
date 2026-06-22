package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addDocumentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	documentCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.DocumentList, client.DocumentSummary]{
			Use:           "document",
			Short:         "Read Linear documents",
			ListShort:     "List visible documents",
			LimitHelp:     "maximum documents to return",
			GetUse:        "get DOCUMENT_ID",
			GetShort:      "Get one Document by id or slug",
			LoadList:      loadDocumentList,
			PageWithItems: documentPageWithItems,
			LoadGet:       loadDocument,
			WriteItem:     writeDocument,
		},
	)
	addDocumentCommentsCommand(ctx, documentCommand, options)
}

func addDocumentCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "comments DOCUMENT_ID",
		Short: "List document comments without body content",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadDocumentCommentList,
				documentCommentPageWithItems,
				writeCommentMetadata,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum comments to return")
	root.AddCommand(command)
}

func writeDocument(command *cobra.Command, options *rootOptions, document client.DocumentSummary) error {
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

func loadDocumentList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.DocumentList, []client.DocumentSummary, error) {
	documents, err := client.ListDocuments(ctx, runtime.graphqlClient, limit)
	return documents, documents.Documents, err
}

func documentPageWithItems(page client.DocumentList, documents []client.DocumentSummary) client.DocumentList {
	page.Documents = documents
	return page
}

func loadDocument(ctx context.Context, runtime commandRuntime, id string) (client.DocumentSummary, error) {
	return client.GetDocumentByID(ctx, runtime.graphqlClient, id)
}

func loadDocumentCommentList(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.DocumentCommentList, []client.CommentMetadataSummary, error) {
	comments, err := client.ListDocumentComments(ctx, runtime.graphqlClient, args[0], limit)
	return comments, comments.Comments, err
}

func documentCommentPageWithItems(
	page client.DocumentCommentList,
	comments []client.CommentMetadataSummary,
) client.DocumentCommentList {
	page.Comments = comments
	return page
}
