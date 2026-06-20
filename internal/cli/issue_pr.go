package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

type pullRequestPlan struct {
	Title   string   `json:"title"`
	Body    string   `json:"body"`
	Command []string `json:"command"`
}

func addIssuePRCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "pr [ISSUE_ID]",
		Short: "Print a gh pr create command for an issue",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			issue, err := resolveIssueArgument(ctx, options, args)
			if err != nil {
				return err
			}

			return writePullRequestPlan(command, options, pullRequestPlanFromIssue(issue))
		},
	})
}

func pullRequestPlanFromIssue(issue client.IssueSummary) pullRequestPlan {
	title := issue.Identifier + " " + issue.Title
	body := issue.URL

	return pullRequestPlan{
		Title: title,
		Body:  body,
		Command: []string{
			"gh",
			"pr",
			"create",
			"--title",
			title,
			"--body",
			body,
		},
	}
}

func writePullRequestPlan(command *cobra.Command, options *rootOptions, plan pullRequestPlan) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, plan)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"gh pr create --title %s --body %s",
		strconv.Quote(plan.Title),
		strconv.Quote(plan.Body),
	)
}
