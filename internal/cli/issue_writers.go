package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func writeIssueBotActor(command *cobra.Command, options *rootOptions, actor client.IssueBotActor) error {
	return writeItem(command, options, actor, actor.IssueID,
		func(command *cobra.Command, _ *rootOptions, actor client.IssueBotActor) error {
			if actor.Bot == nil {
				return render.WriteLine(command.OutOrStdout(), "%s bot -", actor.IssueID)
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"%s bot %s %s [%s]",
				actor.IssueID,
				emptyDash(actor.Bot.ID),
				emptyDash(actor.Bot.Name),
				actor.Bot.Type,
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
	return writeItem(command, options, history, history.ID,
		func(command *cobra.Command, _ *rootOptions, history client.IssueHistorySummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s issue %s updated_description %t",
				history.ID,
				history.IssueID,
				history.UpdatedDescription,
			)
		})
}

func writeCustomerNeedMetadata(
	command *cobra.Command,
	options *rootOptions,
	need client.CustomerNeedMetadataSummary,
) error {
	return writeItem(command, options, need, need.ID,
		func(command *cobra.Command, _ *rootOptions, need client.CustomerNeedMetadataSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s priority %.0f",
				need.ID,
				emptyDash(need.CustomerName),
				emptyDash(need.Issue),
				need.Priority,
			)
		})
}

func writeIssueSharedAccess(
	command *cobra.Command,
	options *rootOptions,
	access client.IssueSharedAccessSummary,
) error {
	return writeItem(command, options, access, access.IssueID,
		func(command *cobra.Command, _ *rootOptions, access client.IssueSharedAccessSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s shared=%t shared_with=%d viewer_shared_only=%t disallowed=%s",
				access.IssueID,
				access.Identifier,
				access.IsShared,
				access.SharedWithCount,
				access.ViewerHasOnlySharedAccess,
				issueSharedAccessFieldsText(access.DisallowedIssueFields),
			)
		})
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
	return writeItem(command, options, suggestion, suggestion.LogID,
		func(command *cobra.Command, _ *rootOptions, suggestion client.IssueFilterSuggestion) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"log_id=%s filter=%s",
				emptyDash(suggestion.LogID),
				emptyDash(string(suggestion.Filter)),
			)
		})
}

func writeIssueTitleSuggestion(
	command *cobra.Command,
	options *rootOptions,
	suggestion client.IssueTitleSuggestion,
) error {
	return writeItem(command, options, suggestion, suggestion.LogID,
		func(command *cobra.Command, _ *rootOptions, suggestion client.IssueTitleSuggestion) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"log_id=%s title=%s",
				emptyDash(suggestion.LogID),
				emptyDash(suggestion.Title),
			)
		})
}

func writeIssues(command *cobra.Command, options *rootOptions, issues []client.IssueSummary) error {
	for _, issue := range issues {
		if err := writeIssue(command, options, issue); err != nil {
			return err
		}
	}

	return nil
}
