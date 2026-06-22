package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

// exportPageLimit caps the comments and attachments pulled into one export.
const exportPageLimit = 250

// issueExportResult is the structured confirmation of a written export.
type issueExportResult struct {
	Path        string `json:"path"`
	Identifier  string `json:"identifier"`
	Comments    int    `json:"comments"`
	Attachments int    `json:"attachments"`
	Truncated   bool   `json:"truncated,omitempty"`
}

func addIssueExportCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "export ISSUE_ID DIR",
		Short: "Export an issue's description, comments, and attachment URLs to a directory",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			return runIssueExport(ctx, command, options, args[0], args[1])
		},
	})
}

func runIssueExport(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	id string,
	dir string,
) error {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return err
	}
	detail, err := client.GetIssueDetail(ctx, runtime.graphqlClient, id)
	if err != nil {
		return err
	}
	comments, err := client.ListIssueComments(ctx, runtime.graphqlClient, id, exportPageLimit)
	if err != nil {
		return err
	}
	attachments, err := client.ListIssueAttachments(ctx, runtime.graphqlClient, id, exportPageLimit)
	if err != nil {
		return err
	}
	document := renderIssueExport(detail, comments.Comments, attachments.Attachments)
	path, err := writeExportDocument(dir, detail.Summary.Identifier, document)
	if err != nil {
		return err
	}
	truncated := comments.HasNextPage || attachments.HasNextPage
	if truncated {
		const note = "export capped at %d comments/attachments; more pages exist"
		if noteErr := writeNote(command, note, exportPageLimit); noteErr != nil {
			return noteErr
		}
	}

	return writeIssueExport(command, options, issueExportResult{
		Path:        path,
		Identifier:  detail.Summary.Identifier,
		Comments:    len(comments.Comments),
		Attachments: len(attachments.Attachments),
		Truncated:   truncated,
	})
}

// writeExportDocument creates dir if needed and writes the assembled export.
func writeExportDocument(dir string, identifier string, document string) (string, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", fmt.Errorf("create %s: %w", dir, err)
	}
	path := filepath.Join(dir, identifier+".md")
	if err := os.WriteFile(path, []byte(document), 0o600); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}

	return path, nil
}

// renderIssueExport assembles the metadata header, description, comments, and
// attachment URLs of one issue into a single markdown document.
func renderIssueExport(
	detail client.IssueDetail,
	comments []client.IssueCommentSummary,
	attachments []client.AttachmentSummary,
) string {
	sections := []string{
		renderExportHeader(detail.Summary),
		renderExportDescription(detail.Description),
		renderExportComments(comments),
		renderExportAttachments(attachments),
	}

	return strings.Join(sections, "\n") + "\n"
}

func renderExportHeader(summary client.IssueSummary) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "# %s — %s\n\n", summary.Identifier, summary.Title)
	for _, field := range exportHeaderFields(summary) {
		fmt.Fprintf(&builder, "- %s: %s\n", field.label, field.value)
	}

	return builder.String()
}

type exportField struct {
	label string
	value string
}

func exportHeaderFields(summary client.IssueSummary) []exportField {
	candidates := []exportField{
		{"URL", summary.URL},
		{"State", summary.State},
		{"Priority", summary.PriorityLabel},
		{"Assignee", summary.Assignee},
		{"Team", summary.Team},
		{"Project", summary.Project},
		{"Created", summary.CreatedAt},
	}
	fields := make([]exportField, 0, len(candidates))
	for _, field := range candidates {
		if field.value != "" {
			fields = append(fields, field)
		}
	}

	return fields
}

func renderExportDescription(description string) string {
	body := strings.TrimSpace(description)
	if body == "" {
		body = "_No description._"
	}

	return "## Description\n\n" + body + "\n"
}

func renderExportComments(comments []client.IssueCommentSummary) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "## Comments (%d)\n\n", len(comments))
	if len(comments) == 0 {
		fmt.Fprint(&builder, "_No comments._\n")

		return builder.String()
	}
	for _, comment := range comments {
		fmt.Fprint(&builder, renderExportComment(comment))
	}

	return builder.String()
}

func renderExportComment(comment client.IssueCommentSummary) string {
	author := comment.DisplayName
	if author == "" {
		author = comment.UserName
	}
	if author == "" {
		author = "Unknown"
	}
	body := strings.TrimSpace(comment.Body)
	if body == "" {
		body = "_(empty)_"
	}

	return fmt.Sprintf("### %s — %s\n\n%s\n\n", author, comment.CreatedAt, body)
}

func renderExportAttachments(attachments []client.AttachmentSummary) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "## Attachments (%d)\n\n", len(attachments))
	if len(attachments) == 0 {
		fmt.Fprint(&builder, "_No attachments._\n")

		return builder.String()
	}
	for _, attachment := range attachments {
		title := attachment.Title
		if title == "" {
			title = attachment.URL
		}
		// Wrap the destination in <> so a URL containing ')' does not close the
		// Markdown link early (CommonMark angle-bracket link destination).
		fmt.Fprintf(&builder, "- [%s](<%s>)\n", title, attachment.URL)
	}

	return builder.String()
}

func writeIssueExport(command *cobra.Command, options *rootOptions, result issueExportResult) error {
	if wrote, err := writeIDOnly(command, options, result.Path); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, result)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s (%d comments, %d attachments)",
		result.Path, result.Comments, result.Attachments,
	)
}
