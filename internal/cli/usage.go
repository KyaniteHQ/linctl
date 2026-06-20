package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/render"
)

type usagePayload struct {
	Topic string `json:"topic"`
	Text  string `json:"text"`
}

var usageTopics = map[string]usagePayload{
	"overview": {
		Topic: "overview",
		Text: "linctl is a schema-aligned Linear CLI for safe daily coordination. " +
			"Configure a pinned target with org_id, team_key, team_id, and optional project_id in .linctl.toml, " +
			"then use reads freely and writes fail-closed against that target. " +
			"Core commands: target, doctor, whoami, current, next, done, issue, comment, cycle, sprint, project, " +
			"project-milestone, organization, rate-limit, document, label, team, user, workflow-state. " +
			"Use --json for structured output, --profile for named configs, --org/--team/--project for explicit " +
			"target overrides, and --timeout for request bounds. " +
			"For scripts, combine --json with --compact or --fields, use --id-only for chaining, --quiet for " +
			"successful writes, --fail-on-empty for monitors, --sort/--order for deterministic lists, and " +
			"--format minimal|compact|full for human output. " +
			"Write flow: resolve the active token, compare it to the pinned target, perform the mutation only " +
			"on match, " +
			"then return the created or updated entity. " +
			"Start every unfamiliar repo with linctl target --json so the active token, org, team, and project " +
			"are visible " +
			"before work starts. " +
			"Use linctl current when the branch carries an issue key, linctl doctor to check config/token/target " +
			"health, and linctl next --dry-run to preview the top-ranked unblocked issue without creating a " +
			"branch or worktree. " +
			"Use domain guidance before writes: linctl issue usage or linctl project usage. " +
			"For test runs, create namespaced throwaway resources and archive them after the observable check.",
	},
	"issue": {
		Topic: "issue",
		Text: "issue commands cover the safe Linear issue loop. " +
			"Use linctl issue list --limit 50 to inspect the resolved team, linctl issue list --state started " +
			"for a workflow state type queue, linctl issue list --project PROJECT_ID for a project queue, " +
			"linctl issue list --mine for issues assigned to the authenticated user, " +
			"linctl issue list --assignee USER_ID for issues assigned to a Linear user id, " +
			"linctl issue list --label LABEL_ID for issues with a Linear label id, " +
			"linctl issue list --cycle CYCLE_ID for issues attached to a Linear Cycle id, " +
			"linctl issue list --created-after DATE for issues created on or after a date, " +
			"linctl issue list --created-since DATE as an alias for created-after, " +
			"linctl issue list --created-before DATE for issues created on or before a date, " +
			"linctl issue list --has-blockers for issues blocked by another issue, " +
			"linctl issue list --blocks for issues blocking another issue, " +
			"linctl issue list --blocked-by ISSUE for issues blocked by that issue, " +
			"linctl issue list --all-teams for broad read-only issue inspection, " +
			"linctl issue search \"text\" for resolved-team text search, " +
			"linctl issue deps ISSUE to inspect parent, child, and blocking relationships, " +
			"linctl issue pr ISSUE to print a gh pr create title/body plan, " +
			"and linctl issue get LIT-123 to read one " +
			"issue by identifier or id. " +
			"Writes require a pinned org/team target: linctl issue create --title \"...\" --description \"...\" " +
			"or --description-file FILE; " +
			"linctl issue update LIT-123 --title \"...\" --description \"...\"; " +
			"linctl issue update LIT-123 --append \"progress note\" or --append-file FILE; " +
			"linctl issue start LIT-123 to assign the issue to you and move it to started; " +
			"linctl issue comment LIT-123 --body \"...\" or --body-file FILE; " +
			"linctl issue reply LIT-123 COMMENT_ID --body \"...\" or --body-file FILE; " +
			"linctl issue close LIT-123. " +
			"If .linctl.toml also pins project_id, writes to existing issues compare the issue's resolved project " +
			"before " +
			"mutating, so same-team wrong-project writes are refused. " +
			"Use --json for automation and parse the returned id, identifier, state, url, team, and project fields. " +
			"For branch-driven work, linctl current derives LIT-123 from the git branch or a jj Linear-issue " +
			"trailer and " +
			"then uses the same issue get path; linctl done closes that current issue through the guarded " +
			"close path. " +
			"Use --fields identifier,title,state with --json for compact agent queues, --id-only for chaining, " +
			"and --fail-on-empty --sort title --order asc for monitor-style lists. " +
			"Recommended agent flow: run linctl target --json, run linctl issue list --json --limit 20, " +
			"linctl issue list --state started --limit 20, linctl issue list --project PROJECT_ID --limit 20, " +
			"linctl issue list --mine --limit 20, " +
			"linctl issue list --assignee USER_ID --limit 20, " +
			"linctl issue list --label LABEL_ID --limit 20, " +
			"linctl issue list --cycle CYCLE_ID --limit 20, " +
			"linctl issue list --created-after 2026-06-01 --limit 20, " +
			"linctl issue list --created-since 2026-06-01 --limit 20, " +
			"linctl issue list --created-before 2026-06-30 --limit 20, " +
			"linctl issue list --has-blockers --limit 20, " +
			"linctl issue list --blocks --limit 20, " +
			"linctl issue list --blocked-by LIT-123 --limit 20, " +
			"linctl issue list --all-teams --limit 20, " +
			"linctl issue deps LIT-123 --limit 20, " +
			"linctl issue pr LIT-123, " +
			"linctl next --dry-run, " +
			"or linctl issue search \"text\" --limit 20 to " +
			"confirm the " +
			"visible queue, then perform exactly one write command with a concrete title, body, or status change. " +
			"If a write fails with target mismatch, do not retry with a different token blindly; inspect the " +
			"expected and " +
			"resolved ids and fix the local target configuration first. " +
			"For temporary QA issues, use a linctl-it-<runid> title prefix, verify via issue get or issue list, " +
			"then close " +
			"or archive through the cleanup path used by the test harness. " +
			"Keep comments concise and avoid pasting secrets, private logs, or unredacted user data.",
	},
	"project": {
		Topic: "project",
		Text: "project commands cover the safe Linear project loop. " +
			"Use linctl project list --limit 50 to list projects attached to the resolved team, " +
			"linctl project get PROJECT_ID to inspect one project, and linctl project members PROJECT_ID to " +
			"list current " +
			"members. " +
			"Use linctl project updates PROJECT_ID --limit 20 for read-only project status history, and " +
			"linctl project-milestone list PROJECT_ID --limit 20 for ProjectMilestone context. " +
			"Project create is team-scoped: linctl project create --name \"linctl-it-<runid>\" " +
			"--description \"...\"; " +
			"it compares only org/team because the project does not exist yet. " +
			"Project update and archive are resource-scoped: linctl project update PROJECT_ID --name \"...\" " +
			"--description \"...\" and linctl project archive PROJECT_ID both resolve the project first and " +
			"refuse if " +
			"the pinned project_id differs. " +
			"ProjectMilestone create and update are resource-scoped project writes: use " +
			"linctl project-milestone create PROJECT_ID --name \"...\" and " +
			"linctl project-milestone update PROJECT_MILESTONE_ID --name \"...\" --target-date YYYY-MM-DD; " +
			"both compare the resolved project before writing. " +
			"Prefer namespaced throwaway projects for tests, archive them after verification, and use --json " +
			"when another " +
			"agent or script will consume the result. " +
			"Recommended agent flow: run linctl target --json, run linctl project list --json --limit 20, " +
			"inspect linctl project updates PROJECT_ID --limit 20 when status context matters, inspect " +
			"linctl project-milestone list PROJECT_ID --limit 20 when milestone context matters, create the " +
			"namespaced project, list again and match the returned id, then archive with --project set to that " +
			"new id if the " +
			"repo target pins a different fixture project. " +
			"That explicit override still goes through target comparison; it is not a bypass. " +
			"Use project members for read-only membership inspection. " +
			"Do not hard-delete projects in v1; cleanup means archive, and a failed cleanup should be reported " +
			"with the " +
			"project id so it can be retried safely.",
	},
	"cycle": {
		Topic: "cycle",
		Text: "cycle commands cover Linear Cycles for the resolved team. " +
			"Use linctl cycle list --limit 20 to list Cycles with derived status, " +
			"and linctl cycle get CYCLE_ID to inspect one Cycle by id or slug. " +
			"Cycle writes are team-scoped: linctl cycle create --starts-at START --ends-at END " +
			"--name \"...\", linctl cycle update CYCLE_ID --name \"...\", and " +
			"linctl cycle archive CYCLE_ID all compare the pinned team before writing. " +
			"Use linctl sprint current for the active Cycle report alias, " +
			"and linctl sprint report CYCLE_ID --limit 20 for Cycle issue status. " +
			"Sprint is a read-only report alias over Cycle; do not create Sprint mutations.",
	},
}

func addUsageCommand(root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "usage [overview|issue|project|cycle]",
		Short: "Show compact linctl usage guidance",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			topic := "overview"
			if len(args) == 1 {
				topic = args[0]
			}

			return writeUsage(command, options, topic)
		},
	})
}

func addDomainUsageCommand(root *cobra.Command, options *rootOptions, topic string) {
	root.AddCommand(&cobra.Command{
		Use:   "usage",
		Short: "Show compact usage guidance for this domain",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return writeUsage(command, options, topic)
		},
	})
}

func writeUsage(command *cobra.Command, options *rootOptions, topic string) error {
	payload, ok := usageTopics[topic]
	if !ok {
		return fmt.Errorf("unknown usage topic %q", topic)
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, payload)
	}

	return render.WriteLine(command.OutOrStdout(), "%s", payload.Text)
}
