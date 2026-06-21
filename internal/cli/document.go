package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addDocumentCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	documentCommand := &cobra.Command{
		Use:   "document",
		Short: "Read Linear documents",
	}
	addDocumentListCommand(ctx, documentCommand, options)
	addDocumentGetCommand(ctx, documentCommand, options)
	addDocumentCommentsCommand(ctx, documentCommand, options)
	root.AddCommand(documentCommand)
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

func addDocumentListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "list",
		Short: "List visible documents",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadDocumentList,
				documentPageWithItems,
				writeDocument,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum documents to return")
	root.AddCommand(command)
}

func addDocumentGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get DOCUMENT_ID",
		Short: "Get one Document by id or slug",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			document, err := client.GetDocumentByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeDocument(command, options, document)
		},
	})
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
