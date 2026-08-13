package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/KyaniteHQ/linctl/internal/client"
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
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
		Use:   "export ISSUE_ID DIR",
		Short: "Export the description, the comments, and the attachment URLs of an issue to a directory",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			return runIssueExport(ctx, command, options, args[0], args[1])
		},
	})
}

func addProjectExportCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "export PROJECT_ID DIR",
		Short: "Export the content and the attachment URLs of a project to a directory",
		Args:  cobra.ExactArgs(2),
		RunE: func(command *cobra.Command, args []string) error {
			return runProjectExport(ctx, command, options, args[0], args[1])
		},
	}
	command.ValidArgsFunction = firstArgCompletion(ctx, options, projectIDCandidates)
	addCommandWithSafety(root, CommandSafetyRead, command)
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
	detail := client.IssueDetail{}
	comments := client.IssueCommentList{}
	attachments := client.AttachmentList{}
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		var detailErr error
		detail, detailErr = client.GetIssueDetail(groupCtx, runtime.graphqlClient, id)
		return detailErr
	})
	group.Go(func() error {
		var commentsErr error
		comments, commentsErr = client.ListIssueComments(groupCtx, runtime.graphqlClient, id, exportPageLimit)
		return commentsErr
	})
	group.Go(func() error {
		var attachmentsErr error
		attachments, attachmentsErr = client.ListIssueAttachments(groupCtx, runtime.graphqlClient, id, exportPageLimit)
		return attachmentsErr
	})
	if err := group.Wait(); err != nil {
		return err
	}
	document := renderIssueExport(detail, comments.Comments, attachments.Attachments)
	leaf, err := issueExportLeaf(detail.Summary.Identifier)
	if err != nil {
		return err
	}
	path, err := writeExportDocument(dir, leaf, document)
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

// exportIdentifierPattern pins the Linear issue identifier grammar the export
// filename is built from; a response value outside it must not reach the
// filesystem, so the identifier cannot smuggle path separators into the join.
var exportIdentifierPattern = regexp.MustCompile(`^[A-Z][A-Z0-9]+-[0-9]+$`)

// exportSafeLeafPattern is the shared filesystem rule for an export filename
// stem: one path segment, no separators, no parent-directory tokens.
var exportSafeLeafPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// writeExportDocument creates dir if needed and writes the assembled export.
func writeExportDocument(dir string, leaf string, document string) (string, error) {
	if !exportSafeLeafPattern.MatchString(leaf) {
		return "", fmt.Errorf("refusing export: %q is not a valid export filename", leaf)
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return "", fmt.Errorf("create %s: %w", dir, err)
	}
	path := filepath.Join(dir, leaf+".md")
	if err := os.WriteFile(path, []byte(document), 0o600); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}

	return path, nil
}

func issueExportLeaf(identifier string) (string, error) {
	if !exportIdentifierPattern.MatchString(identifier) {
		return "", fmt.Errorf("refusing export: %q is not a valid Linear issue identifier", identifier)
	}

	return identifier, nil
}

func projectExportLeaf(summary client.ProjectSummary) (string, error) {
	if exportSafeLeafPattern.MatchString(summary.SlugID) {
		return summary.SlugID, nil
	}
	if exportSafeLeafPattern.MatchString(summary.ID) {
		return summary.ID, nil
	}

	return "", fmt.Errorf("refusing export: %q is not a valid export filename", summary.SlugID)
}

// projectExportResult is the structured confirmation of a written project export.
type projectExportResult struct {
	Path        string `json:"path"`
	SlugID      string `json:"slug_id"`
	Attachments int    `json:"attachments"`
	Truncated   bool   `json:"truncated,omitempty"`
}

func runProjectExport(
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
	detail := client.ProjectDetail{}
	attachments := client.ProjectAttachmentList{}
	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		var detailErr error
		detail, detailErr = client.GetProjectDetail(groupCtx, runtime.graphqlClient, id)
		return detailErr
	})
	group.Go(func() error {
		var attachmentsErr error
		attachments, attachmentsErr = client.ListProjectAttachments(
			groupCtx, runtime.graphqlClient, id, exportPageLimit,
		)
		return attachmentsErr
	})
	if err := group.Wait(); err != nil {
		return err
	}
	leaf, err := projectExportLeaf(detail.Summary)
	if err != nil {
		return err
	}
	document := renderProjectExport(detail, attachments.Attachments)
	path, err := writeExportDocument(dir, leaf, document)
	if err != nil {
		return err
	}
	truncated := attachments.HasNextPage
	if truncated {
		const note = "export capped at %d attachments; more pages exist"
		if noteErr := writeNote(command, note, exportPageLimit); noteErr != nil {
			return noteErr
		}
	}

	return writeProjectExport(command, options, projectExportResult{
		Path:        path,
		SlugID:      detail.Summary.SlugID,
		Attachments: len(attachments.Attachments),
		Truncated:   truncated,
	})
}

func renderProjectExport(detail client.ProjectDetail, attachments []client.AttachmentSummary) string {
	sections := []string{
		renderProjectExportHeader(detail.Summary),
		renderExportContent(detail.Content),
		renderExportAttachments(attachments),
	}

	return strings.Join(sections, "\n") + "\n"
}

func renderProjectExportHeader(summary client.ProjectSummary) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "# %s\n\n", summary.Name)
	for _, field := range projectExportHeaderFields(summary) {
		fmt.Fprintf(&builder, "- %s: %s\n", field.label, field.value)
	}

	return builder.String()
}

func projectExportHeaderFields(summary client.ProjectSummary) []exportField {
	candidates := []exportField{
		{"ID", summary.ID},
		{"Slug", summary.SlugID},
		{"URL", summary.URL},
		{"Status", summary.Status.Name},
		{"Lead", summary.Lead},
		{"Description", summary.Description},
	}
	fields := make([]exportField, 0, len(candidates))
	for _, field := range candidates {
		if field.value != "" {
			fields = append(fields, field)
		}
	}

	return fields
}

func renderExportContent(content string) string {
	body := strings.TrimSpace(content)
	if body == "" {
		body = "_No content._"
	}

	return "## Content\n\n" + body + "\n"
}

func writeProjectExport(command *cobra.Command, options *rootOptions, result projectExportResult) error {
	return writeItemLine(
		command, options, result, result.Path,
		"%s (%d attachments)", result.Path, result.Attachments,
	)
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
	return writeItemLine(
		command, options, result, result.Path,
		"%s (%d comments, %d attachments)", result.Path, result.Comments, result.Attachments,
	)
}
