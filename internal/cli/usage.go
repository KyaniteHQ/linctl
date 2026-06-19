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
			"Core commands: target, whoami, current, issue, project. " +
			"Use --json for structured output, --profile for named configs, --org/--team/--project for explicit " +
			"target overrides, and --timeout for request bounds. " +
			"Write flow: resolve the active token, compare it to the pinned target, perform the mutation only " +
			"on match, " +
			"then return the created or updated entity. " +
			"Start every unfamiliar repo with linctl target --json so the active token, org, team, and project " +
			"are visible " +
			"before work starts. " +
			"Use linctl current when the branch carries an issue key. " +
			"Use domain guidance before writes: linctl issue usage or linctl project usage. " +
			"For test runs, create namespaced throwaway resources and archive them after the observable check.",
	},
	"issue": {
		Topic: "issue",
		Text: "issue commands cover the safe Linear issue loop. " +
			"Use linctl issue list --limit 50 to inspect the resolved team, and linctl issue get LIT-123 to read one " +
			"issue by identifier or id. " +
			"Writes require a pinned org/team target: linctl issue create --title \"...\" --description \"...\"; " +
			"linctl issue update LIT-123 --title \"...\" --description \"...\"; " +
			"linctl issue comment LIT-123 --body \"...\"; linctl issue close LIT-123. " +
			"If .linctl.toml also pins project_id, writes to existing issues compare the issue's resolved project " +
			"before " +
			"mutating, so same-team wrong-project writes are refused. " +
			"Use --json for automation and parse the returned id, identifier, state, url, team, and project fields. " +
			"For branch-driven work, linctl current derives LIT-123 from the git branch or a jj Linear-issue " +
			"trailer and " +
			"then uses the same issue get path. " +
			"Recommended agent flow: run linctl target --json, run linctl issue list --json --limit 20 to " +
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
			"Project create is team-scoped: linctl project create --name \"linctl-it-<runid>\" " +
			"--description \"...\"; " +
			"it compares only org/team because the project does not exist yet. " +
			"Project update and archive are resource-scoped: linctl project update PROJECT_ID --name \"...\" " +
			"--description \"...\" and linctl project archive PROJECT_ID both resolve the project first and " +
			"refuse if " +
			"the pinned project_id differs. " +
			"Prefer namespaced throwaway projects for tests, archive them after verification, and use --json " +
			"when another " +
			"agent or script will consume the result. " +
			"Recommended agent flow: run linctl target --json, run linctl project list --json --limit 20, create the " +
			"namespaced project, list again and match the returned id, then archive with --project set to that " +
			"new id if the " +
			"repo target pins a different fixture project. " +
			"That explicit override still goes through target comparison; it is not a bypass. " +
			"Use project members for read-only membership inspection. " +
			"Do not hard-delete projects in v1; cleanup means archive, and a failed cleanup should be reported " +
			"with the " +
			"project id so it can be retried safely.",
	},
}

func addUsageCommand(root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "usage [overview|issue|project]",
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
	if options.json {
		return render.WriteJSON(command.OutOrStdout(), payload)
	}

	return render.WriteLine(command.OutOrStdout(), "%s", payload.Text)
}
