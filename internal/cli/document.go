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
	addDocumentCreateCommand(ctx, documentCommand, options)
	addDocumentUpdateCommand(ctx, documentCommand, options)
}

func addDocumentCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.DocumentCreateRequest{}
	contentFile := ""
	command := &cobra.Command{
		Use:   "create",
		Short: "Create a document in the pinned target",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			if err := resolveDocumentContent(command, &request.Content, contentFile); err != nil {
				return err
			}
			document, err := client.CreateDocument(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeDocument(command, options, document)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "document title")
	command.Flags().StringVar(&request.Content, "content", "", "document content as markdown; use - to read stdin")
	command.Flags().StringVar(&contentFile, "content-file", "", "read document content from file")
	root.AddCommand(command)
}

func addDocumentUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.DocumentUpdateRequest{}
	contentFile := ""
	command := &cobra.Command{
		Use:   "update DOCUMENT_ID",
		Short: "Update a document after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			request.ID = args[0]
			if err := resolveDocumentContent(command, &request.Content, contentFile); err != nil {
				return err
			}
			document, err := client.UpdateDocument(ctx, runtime.graphqlClient, runtime.config.Target, request)
			if err != nil {
				return err
			}

			return writeDocument(command, options, document)
		},
	}
	command.Flags().StringVar(&request.Title, "title", "", "new document title")
	command.Flags().StringVar(&request.Content, "content", "", "new document content as markdown; use - to read stdin")
	command.Flags().StringVar(&contentFile, "content-file", "", "read new document content from file")
	root.AddCommand(command)
}

// resolveDocumentContent resolves the document content from --content (with "-"
// reading stdin) and the mutually exclusive --content-file.
func resolveDocumentContent(command *cobra.Command, content *string, contentFile string) error {
	if err := resolveBodyFlag(command, content); err != nil {
		return err
	}

	return resolveFileFlag(content, contentFile, "content")
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
	return writeItem(command, options, document, document.ID,
		func(command *cobra.Command, _ *rootOptions, document client.DocumentSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s [%s]",
				document.ID,
				document.Title,
				emptyDash(document.ParentType),
			)
		})
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
