package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func writeIssueBotActor(command *cobra.Command, options *rootOptions, actor client.IssueBotActor) error {
	return writeBotActor(command, options, actor, actor.IssueID, actor.Bot)
}

// writeBotActor renders one entity's optional bot actor: the full payload goes
// to JSON output while the human line carries the owning entity id plus the
// bot identity or a dash when no bot is attached.
func writeBotActor[T any](
	command *cobra.Command,
	options *rootOptions,
	actor T,
	entityID string,
	bot *client.ActorBotSummary,
) error {
	return writeItem(command, options, actor, entityID,
		func(command *cobra.Command, _ *rootOptions, _ T) error {
			if bot == nil {
				return render.WriteLine(command.OutOrStdout(), "%s bot -", entityID)
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"%s bot %s %s [%s]",
				entityID,
				emptyDash(bot.ID),
				emptyDash(bot.Name),
				bot.Type,
			)
		})
}

func writeIssueStateSpan(command *cobra.Command, options *rootOptions, span client.IssueStateSpanSummary) error {
	return writeItem(command, options, span, span.ID,
		func(command *cobra.Command, _ *rootOptions, span client.IssueStateSpanSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s %s -> %s",
				span.ID,
				emptyDash(span.StateName),
				emptyDash(span.StateType),
				span.StartedAt,
				emptyDash(span.EndedAt),
			)
		})
}

func writeIssueHistory(command *cobra.Command, options *rootOptions, history client.IssueHistorySummary) error {
	return writeItemLine(
		command, options, history, history.ID,
		"%s issue %s updated_description %t", history.ID, history.IssueID, history.UpdatedDescription,
	)
}

func writeCustomerNeedMetadata(
	command *cobra.Command,
	options *rootOptions,
	need client.CustomerNeedMetadataSummary,
) error {
	return writeItemLine(
		command, options, need, need.ID,
		"%s %s %s priority %.0f", need.ID, emptyDash(need.CustomerName), emptyDash(need.Issue), need.Priority,
	)
}

func writeIssueSharedAccess(
	command *cobra.Command,
	options *rootOptions,
	access client.IssueSharedAccessSummary,
) error {
	return writeItemLine(
		command, options, access, access.IssueID,
		"%s %s shared=%t shared_with=%d viewer_shared_only=%t disallowed=%s",
		access.IssueID,
		access.Identifier,
		access.IsShared,
		access.SharedWithCount,
		access.ViewerHasOnlySharedAccess,
		issueSharedAccessFieldsText(access.DisallowedIssueFields),
	)
}

func issueSharedAccessFieldsText(fields []string) string {
	if len(fields) == 0 {
		return "-"
	}

	return strings.Join(fields, ",")
}

func writeIssuePriorityValues(
	command *cobra.Command,
	options *rootOptions,
	values []client.IssuePriorityValue,
) error {
	return writeItemNoID(command, options, values,
		func(command *cobra.Command, _ *rootOptions, values []client.IssuePriorityValue) error {
			for _, value := range values {
				if err := render.WriteLine(command.OutOrStdout(), "%d %s", value.Priority, value.Label); err != nil {
					return err
				}
			}

			return nil
		})
}

func writeIssueFilterSuggestion(
	command *cobra.Command,
	options *rootOptions,
	suggestion client.IssueFilterSuggestion,
) error {
	return writeItemLine(
		command, options, suggestion, suggestion.LogID,
		"log_id=%s filter=%s", emptyDash(suggestion.LogID), emptyDash(string(suggestion.Filter)),
	)
}

func writeIssueTitleSuggestion(
	command *cobra.Command,
	options *rootOptions,
	suggestion client.IssueTitleSuggestion,
) error {
	return writeItemLine(
		command, options, suggestion, suggestion.LogID,
		"log_id=%s title=%s", emptyDash(suggestion.LogID), emptyDash(suggestion.Title),
	)
}

func writeIssues(command *cobra.Command, options *rootOptions, issues []client.IssueSummary) error {
	for _, issue := range issues {
		if err := writeIssue(command, options, issue); err != nil {
			return err
		}
	}

	return nil
}
