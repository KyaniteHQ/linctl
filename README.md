# linctl

`linctl` is a schema-aligned Go CLI for Linear.

It is built for agent-safe daily coordination: reads are lightweight, writes re-resolve the active Linear
token and fail closed unless the resolved org/team/project matches the pinned target.

## Install

### Clean Linux Machine

These commands start from a fresh Ubuntu 24.04 environment with no project tools installed:

```bash
apt-get update
apt-get install -y build-essential ca-certificates curl git tar

curl -fsSL https://go.dev/dl/go1.26.4.linux-amd64.tar.gz -o /tmp/go.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf /tmp/go.tar.gz
export PATH="/usr/local/go/bin:$PATH"

git clone https://github.com/KyaniteHQ/linctl.git
cd linctl

go run ./cmd/linctl usage
go run ./cmd/linctl --version
```

From source:

```bash
go install github.com/KyaniteHQ/linctl/cmd/linctl@latest
```

After the first tagged release:

```bash
brew install --cask KyaniteHQ/linctl/linctl
```

## Configure

Create `.linctl.toml` in a repo:

```toml
[target]
org_id = "linear-org-id"
team_key = "LIT"
team_id = "linear-team-id"
project_id = "optional-linear-project-id"
```

Inject credentials with `LINCTL_TOKEN` or `LINEAR_API_KEY`; do not commit tokens.

## Usage

```bash
linctl usage
linctl target --json
linctl doctor
linctl application info app-client-id
linctl agent-activity list --limit 20
linctl agent-activity get agent-activity-id
linctl agent-skill list --limit 20
linctl agent-skill get agent-skill-id
linctl audit-entry types
linctl organization exists kyanite
linctl organization templates --limit 20
linctl rate-limit status
linctl notification list --limit 20
linctl notification get notification-id
linctl notification subscription list --limit 20
linctl notification subscription get notification-subscription-id
linctl triage-responsibility list --limit 20
linctl triage-responsibility get triage-responsibility-id
linctl triage-responsibility manual-selection triage-responsibility-id
linctl sla-configuration list team-id
linctl semantic-search "agent search" --limit 20
linctl release-pipeline list --limit 20
linctl release-pipeline get release-pipeline-id
linctl release-pipeline releases release-pipeline-id --limit 20
linctl release-pipeline stages release-pipeline-id --limit 20
linctl release-stage list --limit 20
linctl release-stage get release-stage-id
linctl release-stage releases release-stage-id --limit 20
linctl release list --limit 20
linctl release search "mobile" --limit 20
linctl release get release-id
linctl release history release-id --limit 20
linctl release links release-id --limit 20
linctl external-link get external-link-id
linctl release-note list --limit 20
linctl release-note get release-note-id
linctl current --json
linctl next --dry-run
linctl done
linctl issue id
linctl issue title
linctl issue url
linctl issue branch LIT-123
linctl issue deps LIT-123 --limit 20
linctl issue pr LIT-123
linctl issue comments LIT-123 --limit 20
linctl comment list --limit 20
linctl comment get comment-id
linctl issue start LIT-123
linctl issue reply LIT-123 comment-id --body "thread reply"
linctl issue usage
linctl cycle list --limit 20
linctl cycle get cycle-id
linctl cycle create --starts-at 2026-07-01T00:00:00Z --ends-at 2026-07-15T00:00:00Z --name "Planning"
linctl cycle update cycle-id --name "Updated planning"
linctl cycle archive cycle-id
linctl sprint current
linctl sprint report cycle-id --limit 20
linctl project updates project-id --limit 20
linctl project-update list --limit 20
linctl project-update get project-update-id
linctl project-milestone list project-id --limit 20
linctl project-milestone get project-milestone-id
linctl project-milestone create project-id --name "Launch milestone"
linctl project-milestone update project-milestone-id --target-date 2026-06-30
linctl project-status list --limit 20
linctl project-status get project-status-id
linctl document list --limit 20
linctl document get document-id
linctl label list --limit 20
linctl label get label-id
linctl team list --limit 20
linctl team get team-id
linctl team members team-id --limit 20
linctl user list --limit 20
linctl user get user-id
linctl user me
linctl user drafts --limit 20
linctl workflow-state list --limit 20
linctl workflow-state get workflow-state-id
linctl time-schedule list --limit 20
linctl time-schedule get time-schedule-id
linctl template list --limit 20
linctl template get template-id
linctl initiative list --limit 20
linctl initiative get initiative-id
linctl initiative history initiative-id --limit 20
linctl initiative links initiative-id --limit 20
linctl initiative sub-initiatives initiative-id --limit 20
linctl initiative updates initiative-id --limit 20
linctl initiative-relation list --limit 20
linctl initiative-relation get initiative-relation-id
linctl initiative-to-project list --limit 20
linctl initiative-to-project get initiative-to-project-id
linctl initiative-update list --limit 20
linctl initiative-update get initiative-update-id
linctl roadmap list --limit 20
linctl roadmap get roadmap-id
linctl custom-view list --limit 20
linctl custom-view subscribers custom-view-id
linctl custom-view get custom-view-id
linctl custom-view initiatives custom-view-id --limit 20
linctl custom-view organization-preferences custom-view-id
linctl custom-view organization-preferences values custom-view-id
linctl custom-view preference-values custom-view-id
linctl customer list --limit 20
linctl customer get customer-id
linctl customer-need list --limit 20
linctl customer-need get customer-need-id
linctl customer-status list --limit 20
linctl customer-status get customer-status-id
linctl customer-tier list --limit 20
linctl customer-tier get customer-tier-id
linctl favorite list --limit 20
linctl favorite children favorite-folder-id --limit 20
linctl favorite get favorite-id
linctl emoji list --limit 20
linctl emoji get emoji-id
linctl attachment list --limit 20
linctl attachment url https://example.com/spec --limit 20
linctl attachment get attachment-id
linctl project usage
```

Script-friendly output controls are global:

```bash
linctl --json --compact issue get LIT-123
linctl --json --fields identifier,title,state issue list --limit 20
linctl issue list --state started --limit 20
linctl issue list --project project-id --limit 20
linctl issue list --mine --limit 20
linctl issue list --assignee user-id --limit 20
linctl issue list --label label-id --limit 20
linctl issue list --cycle cycle-id --limit 20
linctl issue list --created-after 2026-06-01 --limit 20
linctl issue list --created-since 2026-06-01 --limit 20
linctl issue list --created-before 2026-06-30 --limit 20
linctl issue list --has-blockers --limit 20
linctl issue list --blocks --limit 20
linctl issue list --blocked-by LIT-123 --limit 20
linctl issue list --all-teams --limit 20
linctl issue search "needle" --limit 20
linctl issue deps LIT-123 --limit 20
linctl issue pr LIT-123
linctl next --dry-run
linctl application info app-client-id
linctl agent-activity list --limit 20
linctl agent-activity get agent-activity-id
linctl agent-skill list --limit 20
linctl agent-skill get agent-skill-id
linctl audit-entry types
linctl organization exists kyanite
linctl organization templates --limit 20
linctl rate-limit status
linctl notification list --limit 20
linctl notification get notification-id
linctl notification subscription list --limit 20
linctl notification subscription get notification-subscription-id
linctl triage-responsibility list --limit 20
linctl triage-responsibility get triage-responsibility-id
linctl triage-responsibility manual-selection triage-responsibility-id
linctl sla-configuration list team-id
linctl semantic-search "agent search" --limit 20
linctl release-pipeline list --limit 20
linctl release-pipeline get release-pipeline-id
linctl release-pipeline releases release-pipeline-id --limit 20
linctl release-pipeline stages release-pipeline-id --limit 20
linctl release-stage list --limit 20
linctl release-stage get release-stage-id
linctl release-stage releases release-stage-id --limit 20
linctl release list --limit 20
linctl release search "mobile" --limit 20
linctl release get release-id
linctl release history release-id --limit 20
linctl release links release-id --limit 20
linctl external-link get external-link-id
linctl release-note list --limit 20
linctl release-note get release-note-id
linctl cycle list --limit 20
linctl cycle get cycle-id
linctl cycle create --starts-at 2026-07-01T00:00:00Z --ends-at 2026-07-15T00:00:00Z --name "Planning"
linctl cycle update cycle-id --name "Updated planning"
linctl cycle archive cycle-id
linctl sprint current
linctl sprint report cycle-id --limit 20
linctl issue start LIT-123
linctl done
linctl --id-only issue create --title "small task"
linctl issue create --title "small task" --description-file ./issue.md
linctl --quiet issue update LIT-123 --title "retitled"
linctl issue update LIT-123 --append "progress note"
linctl issue update LIT-123 --append-file ./progress.md
printf 'progress note\n' | linctl issue comment LIT-123 --body -
linctl issue comment LIT-123 --body-file ./comment.md
linctl issue reply LIT-123 comment-id --body "thread reply"
linctl issue reply LIT-123 comment-id --body-file ./reply.md
linctl comment list --limit 20
linctl comment get comment-id
linctl project-milestone create project-id --name "Launch milestone"
linctl project-milestone update project-milestone-id --name "Renamed milestone"
linctl project-update list --limit 20
linctl project-update get project-update-id
linctl project-status list --limit 20
linctl project-status get project-status-id
linctl document list --limit 20
linctl label list --limit 20
linctl team members team-id --limit 20
linctl user me
linctl user drafts --limit 20
linctl workflow-state list --limit 20
linctl workflow-state get workflow-state-id
linctl time-schedule list --limit 20
linctl time-schedule get time-schedule-id
linctl template list --limit 20
linctl template get template-id
linctl initiative list --limit 20
linctl initiative get initiative-id
linctl initiative history initiative-id --limit 20
linctl initiative links initiative-id --limit 20
linctl initiative sub-initiatives initiative-id --limit 20
linctl initiative updates initiative-id --limit 20
linctl initiative-relation list --limit 20
linctl initiative-relation get initiative-relation-id
linctl initiative-to-project list --limit 20
linctl initiative-to-project get initiative-to-project-id
linctl initiative-update list --limit 20
linctl initiative-update get initiative-update-id
linctl roadmap list --limit 20
linctl roadmap get roadmap-id
linctl custom-view list --limit 20
linctl custom-view subscribers custom-view-id
linctl custom-view get custom-view-id
linctl custom-view initiatives custom-view-id --limit 20
linctl custom-view organization-preferences custom-view-id
linctl custom-view organization-preferences values custom-view-id
linctl custom-view preference-values custom-view-id
linctl customer list --limit 20
linctl customer get customer-id
linctl customer-need list --limit 20
linctl customer-need get customer-need-id
linctl customer-status list --limit 20
linctl customer-status get customer-status-id
linctl customer-tier list --limit 20
linctl customer-tier get customer-tier-id
linctl favorite list --limit 20
linctl favorite children favorite-folder-id --limit 20
linctl favorite get favorite-id
linctl emoji list --limit 20
linctl emoji get emoji-id
linctl attachment list --limit 20
linctl attachment url https://example.com/spec --limit 20
linctl attachment get attachment-id
linctl --fail-on-empty --sort title --order asc issue list
linctl --format minimal issue get LIT-123
```

Issue, project, Cycle, and ProjectMilestone writes require a pinned target. Team-scoped creates compare
org/team; resource-scoped updates and archives resolve the resource first and compare the pinned project
when configured. Application, AgentActivity, AgentSkill, AuditEntry, Organization, rate-limit, notification, release-pipeline, release-stage, release, release-note, external-link, semantic-search, comment, ProjectUpdate, document, label, team, user, workflow-state, time-schedule, TriageResponsibility, SLA configuration, template, initiative, initiative-relation, initiative-to-project, initiative-update, roadmap, custom-view, customer, customer-need, customer-status, customer-tier, favorite, emoji, and attachment commands are read-only in the current CLI.

## Development

After following the clean-machine setup above:

```bash
go generate ./...
git diff --exit-code -- internal/client/generated.go
go build ./...
go vet ./...
go test -race -shuffle=on -count=1 ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --timeout 5m ./...
```

The temporary integration fixture is configured in `test/integration-config.json`. Inject
`LINCTL_TEST_TOKEN` from secret storage only when running live integration tests:

```bash
LINCTL_TEST_TOKEN=<token> go test -count=1 -tags=integration ./internal/client
```

Or run the complete live smoke harness:

```bash
task live-smoke
```

The smoke harness builds a temporary CLI binary, runs read-only CLI checks, then runs the integration-tagged
client round trips. Write checks must use disposable `linctl-it-<runid>` resources and archive them during cleanup.
