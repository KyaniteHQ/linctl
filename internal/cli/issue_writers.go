package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func writeIssueBotActor(command *cobra.Command, options *rootOptions, actor client.IssueBotActor) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, actor)
	}
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
}

func writeIssueStateSpan(command *cobra.Command, options *rootOptions, span client.IssueStateSpanSummary) error {
	if wrote, err := writeIDOnly(command, options, span.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, span)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s %s -> %s",
		span.ID,
		emptyDash(span.StateName),
		emptyDash(span.StateType),
		span.StartedAt,
		emptyDash(span.EndedAt),
	)
}

func writeIssueHistory(command *cobra.Command, options *rootOptions, history client.IssueHistorySummary) error {
	if wrote, err := writeIDOnly(command, options, history.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, history)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s issue %s updated_description %t",
		history.ID,
		history.IssueID,
		history.UpdatedDescription,
	)
}

func writeCustomerNeedMetadata(
	command *cobra.Command,
	options *rootOptions,
	need client.CustomerNeedMetadataSummary,
) error {
	if wrote, err := writeIDOnly(command, options, need.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, need)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s priority %.0f",
		need.ID,
		emptyDash(need.CustomerName),
		emptyDash(need.Issue),
		need.Priority,
	)
}

func writeIssueSharedAccess(
	command *cobra.Command,
	options *rootOptions,
	access client.IssueSharedAccessSummary,
) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, access)
	}

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
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, values)
	}
	for _, value := range values {
		if err := render.WriteLine(command.OutOrStdout(), "%d %s", value.Priority, value.Label); err != nil {
			return err
		}
	}

	return nil
}

func writeIssueFilterSuggestion(
	command *cobra.Command,
	options *rootOptions,
	suggestion client.IssueFilterSuggestion,
) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, suggestion)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"log_id=%s filter=%s",
		emptyDash(suggestion.LogID),
		emptyDash(string(suggestion.Filter)),
	)
}

func writeIssueTitleSuggestion(
	command *cobra.Command,
	options *rootOptions,
	suggestion client.IssueTitleSuggestion,
) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, suggestion)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"log_id=%s title=%s",
		emptyDash(suggestion.LogID),
		emptyDash(suggestion.Title),
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
