package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

// issueDraft is the local preview rendered by issue create --dry-run.
type issueDraft struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StateType   string `json:"state_type,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

// applyIssueTemplate fills the issue title and description from a Linear template
// when --template is set. Explicit --title/--description values take precedence;
// the template only supplies defaults for fields left empty.
func applyIssueTemplate(
	ctx context.Context,
	templates issueTemplateReader,
	request *client.IssueCreateRequest,
	templateID string,
) error {
	if templateID == "" {
		return nil
	}
	content, err := templates.GetIssueTemplateContent(ctx, templateID)
	if err != nil {
		return err
	}
	if request.Title == "" {
		request.Title = content.Title
	}
	if request.Description == "" {
		request.Description = content.Description
	}

	return nil
}

// applyIssueSections fills each NAME=VALUE markdown section into the description.
func applyIssueSections(request *client.IssueCreateRequest, sections []string) error {
	for _, section := range sections {
		name, value, found := strings.Cut(section, "=")
		if !found || strings.TrimSpace(name) == "" {
			return fmt.Errorf("invalid --section %q: expected NAME=VALUE", section)
		}
		request.Description = applyMarkdownSection(request.Description, name, value)
	}

	return nil
}

// applyMarkdownSection replaces the body under the first heading matching name
// (case-insensitive, any heading level) with value. When no such heading exists,
// a new "## name" section is appended.
func applyMarkdownSection(description string, name string, value string) string {
	lines := strings.Split(description, "\n")
	start := markdownHeadingIndex(lines, name)
	if start < 0 {
		return appendMarkdownSection(description, name, value)
	}
	end := nextHeadingIndex(lines, start+1)
	rebuilt := make([]string, 0, len(lines)+1)
	rebuilt = append(rebuilt, lines[:start+1]...)
	rebuilt = append(rebuilt, value)
	rebuilt = append(rebuilt, lines[end:]...)

	return strings.Join(rebuilt, "\n")
}

func markdownHeadingIndex(lines []string, name string) int {
	target := strings.ToLower(strings.TrimSpace(name))
	for index, line := range lines {
		if title, ok := markdownHeadingTitle(line); ok && strings.ToLower(title) == target {
			return index
		}
	}

	return -1
}

func nextHeadingIndex(lines []string, from int) int {
	for index := from; index < len(lines); index++ {
		if _, ok := markdownHeadingTitle(lines[index]); ok {
			return index
		}
	}

	return len(lines)
}

// markdownHeadingTitle returns the title of an ATX heading line ("## Title") and
// whether the line is a heading at all.
func markdownHeadingTitle(line string) (string, bool) {
	trimmed := strings.TrimLeft(line, " ")
	hashes := 0
	for hashes < len(trimmed) && trimmed[hashes] == '#' {
		hashes++
	}
	if hashes == 0 || hashes >= len(trimmed) || trimmed[hashes] != ' ' {
		return "", false
	}

	return strings.TrimSpace(trimmed[hashes:]), true
}

func appendMarkdownSection(description string, name string, value string) string {
	section := "## " + strings.TrimSpace(name) + "\n\n" + value
	if strings.TrimSpace(description) == "" {
		return section
	}

	return strings.TrimRight(description, "\n") + "\n\n" + section
}

// writeIssueDraft renders the assembled issue for --dry-run without creating it.
func writeIssueDraft(command *cobra.Command, options *rootOptions, request client.IssueCreateRequest) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, issueDraft{
			Title:       request.Title,
			Description: request.Description,
			StateType:   request.StateType,
			Priority:    request.Priority,
		})
	}

	return render.WriteLine(command.OutOrStdout(), "%s\n\n%s", emptyDash(request.Title), request.Description)
}
