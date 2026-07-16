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
			Use:       "document",
			Short:     "Read Linear documents",
			ListShort: "List visible documents",
			LimitHelp: "maximum documents to return",
			GetUse:    "get DOCUMENT_ID",
			GetShort:  "Get one Document by id or slug",
			LoadList:  loadDocumentList,
			LoadGet:   loadDocument,
			WriteItem: writeDocument,
		},
	)
	addDocumentCommentsCommand(ctx, documentCommand, options)
	addDocumentCreateCommand(ctx, documentCommand, options)
	addDocumentUpdateCommand(ctx, documentCommand, options)
}

func addDocumentCreateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.DocumentCreateRequest{}
	contentFile := ""
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.DocumentSummary]{
		Use:   "create",
		Short: "Create a document in the pinned target",
		Args:  cobra.NoArgs,
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Title, "title", "", "document title")
			command.Flags().StringVar(
				&request.Content, "content", "", "document content as markdown; use - to read stdin",
			)
			command.Flags().StringVar(&contentFile, "content-file", "", "read document content from file")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, _ []string,
		) (client.DocumentSummary, error) {
			if err := resolveDocumentContent(command, &request.Content, contentFile); err != nil {
				return client.DocumentSummary{}, err
			}

			return client.CreateDocument(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeDocument,
	})
}

func addDocumentUpdateCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := client.DocumentUpdateRequest{}
	contentFile := ""
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.DocumentSummary]{
		Use:   "update DOCUMENT_ID",
		Short: "Update a document after pinned-target comparison",
		Args:  cobra.ExactArgs(1),
		Configure: func(command *cobra.Command) {
			command.Flags().StringVar(&request.Title, "title", "", "new document title")
			command.Flags().StringVar(
				&request.Content, "content", "", "new document content as markdown; use - to read stdin",
			)
			command.Flags().StringVar(&contentFile, "content-file", "", "read new document content from file")
		},
		Run: func(
			ctx context.Context, command *cobra.Command, runtime commandRuntime, args []string,
		) (client.DocumentSummary, error) {
			request.ID = args[0]

			if err := resolveDocumentContent(command, &request.Content, contentFile); err != nil {
				return client.DocumentSummary{}, err
			}

			return client.UpdateDocument(ctx, runtime.graphqlClient, runtime.config.Target, request)
		},
		Write: writeDocument,
	})
}

// resolveDocumentContent resolves the document content from --content (with "-"
// reading stdin) and the mutually exclusive --content-file.
func resolveDocumentContent(command *cobra.Command, content *string, contentFile string) error {
	if err := resolveBodyFlag(command, content); err != nil {
		return err
	}

	return resolveFileFlag(command, content, contentFile, "content")
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
				writeCommentMetadata,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum comments to return")
	root.AddCommand(preflightReadListCommand(command, loadDocumentCommentList))
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
