package cli

import (
	"context"
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
			"Set a pinned target in .linctl.toml with org_id, team_key, team_id, and an optional project_id. " +
			"Use --config PATH to select another repo target config explicitly. " +
			"Reads are free. Writes fail closed against the pinned target. " +
			"Core commands: target, doctor, whoami, current, next, done, issue, comment, cycle, sprint, project, " +
			"project-milestone, organization, rate-limit, document, label, team, user, workflow-state. " +
			"Use initiative and initiative-to-project for new strategic planning reads. " +
			"Use initiative-update list and get for InitiativeUpdate reads. " +
			"Use initiative-update create for a guarded InitiativeUpdate write. " +
			"The roadmap and roadmap-to-project commands read deprecated Linear Roadmap data for compatibility. " +
			"Use --json for structured output. Use --profile to select a named auth config. " +
			"Use --org, --team, and --project to set the target explicitly. Use --timeout to bound one command. " +
			"For scripts, add --compact or --fields to --json. Use --id-only to chain commands. " +
			"Use --quiet for a write that succeeds. Use --fail-on-empty for a monitor. " +
			"Use --sort and --order for a deterministic list. " +
			"Use --format minimal, compact, or full for human output. " +
			"A write has four steps. linctl resolves the active OAuth auth state. " +
			"linctl compares that state with the pinned target. linctl makes the change only on a match. " +
			"linctl then returns the created or updated entity. " +
			"Start in an unfamiliar repo with linctl target --json. " +
			"That command shows the active auth target, the organization, the team, and the project. " +
			"Use linctl current when the branch carries an issue key. " +
			"Use linctl doctor to check config load, token presence, and pinned-target confirmation. " +
			"Use linctl next --dry-run to see the top-ranked unblocked issue without a branch or a worktree. " +
			"Read the domain guidance before a write with linctl issue usage or linctl project usage. " +
			"For a test run, create namespaced throwaway resources. Archive them after the check.",
	},
	"issue": {
		Topic: "issue",
		Text: "The issue commands cover the safe Linear issue loop. " +
			"Use linctl issue list --limit 50 to inspect the resolved team. " +
			"Add one filter to that command to narrow the queue. " +
			"--state started selects one workflow state type on list. " +
			"On create and update, --state selects an exact workflow-state name. " +
			"--project PROJECT_ID selects one project. " +
			"--mine selects the issues assigned to the authenticated user. " +
			"--assignee USER_ID selects the issues assigned to one Linear user id. " +
			"--label LABEL_ID selects the issues with one Linear label id. " +
			"--cycle CYCLE_ID selects the issues attached to one Linear Cycle id. " +
			"--created-after DATE selects the issues created on or after a date. " +
			"--created-since DATE is an alias for --created-after. " +
			"--created-before DATE selects the issues created on or before a date. " +
			"--updated-after DATE and --updated-before DATE do the same for the update date. " +
			"--has-blockers selects the issues blocked by another issue. " +
			"--blocks selects the issues that block another issue. " +
			"--blocked-by ISSUE selects the issues blocked by that issue. " +
			"--all-teams reads across every visible team. " +
			"Use linctl issue search \"text\" to search the text of the resolved team. " +
			"Use linctl issue deps ISSUE to inspect the parent, the children, and the blocking relations. " +
			"Use linctl issue pr ISSUE to show a gh pr create plan. " +
			"Use linctl issue get LIT-123 to read one issue by identifier or by id. " +
			"A write needs a pinned organization and team. " +
			"Create an issue with linctl issue create --title \"...\" --description \"...\". " +
			"Use --description-file FILE to read the description from a file. " +
			"Update an issue with linctl issue update LIT-123 --title \"...\" --description \"...\". " +
			"Select one started state by exact name: linctl issue update LIT-123 --state \"In Review\". " +
			"Relate two issues with linctl issue relate A B --type related. " +
			"Name every allowed project with --allowed-project when the issues sit in different projects. " +
			"Append to the description with linctl issue update LIT-123 --append \"note\" or --append-file FILE. " +
			"Use linctl issue start LIT-123 to assign the issue and move it to the started state. " +
			"That command assigns a human viewer, or delegates the issue to an assignable app. " +
			"Comment with linctl issue comment LIT-123 --body \"...\" or --body-file FILE. " +
			"Reply with linctl issue reply LIT-123 COMMENT_ID --body \"...\" or --body-file FILE. " +
			"Close an issue with linctl issue close LIT-123. " +
			"If .linctl.toml also pins project_id, a write to an existing issue compares the project of that issue. " +
			"linctl thus refuses a write to the correct team but the wrong project. " +
			"Use --json for automation. The result has the id, identifier, state, url, team, and project fields. " +
			"For branch work, linctl current reads LIT-123 from the git branch or from a jj Linear-issue trailer. " +
			"linctl current then uses the same path as issue get. " +
			"linctl done closes that Current Issue through the guarded close path. " +
			"Use --fields identifier,title,state with --json for a compact agent queue. " +
			"Use --id-only to chain commands. " +
			"Use --fail-on-empty --sort title --order asc for a monitor list. " +
			"The recommended agent flow has three steps. First, run linctl target --json. " +
			"Second, run linctl issue list --json --limit 20 with the filter that matches the task. " +
			"Third, run exactly one write command with a concrete title, body, or status change. " +
			"If a write fails with a target mismatch, do not retry with different auth. " +
			"Read the expected and resolved ids, then correct the local target configuration. " +
			"For a temporary QA issue, use a linctl-it-<runid> title prefix. " +
			"Check the issue with issue get or issue list, then close or archive it through the cleanup path. " +
			"Keep every comment short. Never write a secret, a private log, or unredacted user data.",
	},
	"project": {
		Topic: "project",
		Text: "The project commands cover the safe Linear project loop. " +
			"Use linctl project list --limit 50 to list the projects attached to the resolved team. " +
			"Use linctl project get PROJECT_ID to inspect one project. " +
			"Use linctl project export PROJECT_ID DIR to write the project content markdown to a local file. " +
			"Use linctl project members PROJECT_ID to list the current members. " +
			"Use linctl project updates PROJECT_ID --limit 20 to read the project status history. " +
			"Use linctl project-milestone list PROJECT_ID --limit 20 to read the ProjectMilestone context. " +
			"Project create is team-scoped. " +
			"Run linctl project create --name \"linctl-it-<runid>\" --description \"...\". " +
			"That command compares only the organization and the team, because the project does not exist yet. " +
			"Project update and archive are resource-scoped. " +
			"Run linctl project update PROJECT_ID --name \"...\" --description \"...\", " +
			"or linctl project archive PROJECT_ID. " +
			"Both commands resolve the project first, then refuse the write when the pinned project_id differs. " +
			"ProjectMilestone create and update are resource-scoped project writes. " +
			"Run linctl project-milestone create PROJECT_ID --name \"...\". " +
			"Run linctl project-milestone update PROJECT_MILESTONE_ID --name \"...\" --target-date YYYY-MM-DD. " +
			"Both commands compare the resolved project before the write. " +
			"For a test, create a namespaced throwaway project, then archive it after the check. " +
			"Use --json when another agent or script reads the result. " +
			"The recommended agent flow has five steps. First, run linctl target --json. " +
			"Second, run linctl project list --json --limit 20. " +
			"Third, read linctl project updates PROJECT_ID --limit 20 when the status matters, " +
			"and linctl project-milestone list PROJECT_ID --limit 20 when the milestone matters. " +
			"Fourth, create the namespaced project, list again, and match the returned id. " +
			"Fifth, archive that project with --project set to the new id " +
			"when the repo target pins a different fixture project. " +
			"That explicit override still goes through the target comparison. It is not a bypass. " +
			"Use project members to inspect the membership without a write. " +
			"Never hard-delete a project. Cleanup means archive. " +
			"Report a failed cleanup with the project id, so a person can retry it safely.",
	},
	"cycle": {
		Topic: "cycle",
		Text: "The cycle commands cover the Linear Cycles of the resolved team. " +
			"Use linctl cycle list --limit 20 to list the Cycles with their derived status. " +
			"Use linctl cycle get CYCLE_ID to inspect one Cycle by id or by slug. " +
			"Cycle writes are team-scoped. " +
			"Run linctl cycle create --starts-at START --ends-at END --name \"...\". " +
			"Run linctl cycle update CYCLE_ID --name \"...\". Run linctl cycle archive CYCLE_ID. " +
			"All three commands compare the pinned team before the write. " +
			"Use linctl sprint current to read the active Cycle. " +
			"Use linctl sprint report CYCLE_ID --limit 20 to read the issue status of one Cycle. " +
			"Sprint is a read-only report alias over Cycle. Never add a Sprint mutation.",
	},
}

func addUsageCommand(_ context.Context, root *cobra.Command, options *rootOptions) {
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
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
	addCommandWithSafety(root, CommandSafetyRead, &cobra.Command{
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
