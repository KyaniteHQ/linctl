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
linctl organization labels --limit 20
linctl organization project-labels --limit 20
linctl organization teams --limit 20
linctl organization templates --limit 20
linctl organization users --limit 20
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
linctl search documents "agent search" --limit 20
linctl search issues "agent search" --limit 20
linctl search projects "agent search" --limit 20
linctl release-pipeline list --limit 20
linctl release-pipeline get release-pipeline-id
linctl release-pipeline releases release-pipeline-id --limit 20
linctl release-pipeline stages release-pipeline-id --limit 20
linctl release-pipeline teams release-pipeline-id --limit 20
linctl release-stage list --limit 20
linctl release-stage get release-stage-id
linctl release-stage releases release-stage-id --limit 20
linctl release list --limit 20
linctl release search "mobile" --limit 20
linctl release get release-id
linctl release history release-id --limit 20
linctl release documents release-id --limit 20
linctl release issues release-id --limit 20
linctl release links release-id --limit 20
linctl external-link get external-link-id
linctl release-note list --limit 20
linctl release-note get release-note-id
linctl issue-to-release list --limit 20
linctl issue-to-release get issue-to-release-id
linctl external-user list --limit 20
linctl external-user get external-user-id
linctl current --json
linctl next --dry-run
linctl done
linctl issue id
linctl issue title
linctl issue url
linctl issue branch LIT-123
linctl issue deps LIT-123 --limit 20
linctl issue attachments LIT-123 --limit 20
linctl issue children LIT-123 --limit 20
linctl issue documents LIT-123 --limit 20
linctl issue former-attachments LIT-123 --limit 20
linctl issue former-needs LIT-123 --limit 20
linctl issue history LIT-123 --limit 20
linctl issue inverse-relations LIT-123 --limit 20
linctl issue labels LIT-123 --limit 20
linctl issue needs LIT-123 --limit 20
linctl issue relations LIT-123 --limit 20
linctl issue releases LIT-123 --limit 20
linctl issue shared-access LIT-123
linctl issue bot-actor LIT-123
linctl issue state-history LIT-123 --limit 20
linctl issue subscribers LIT-123 --limit 20
linctl issue vcs-branch-search get omer/example-branch
linctl issue vcs-branch-search attachments omer/example-branch --limit 20
linctl issue vcs-branch-search bot-actor omer/example-branch
linctl issue vcs-branch-search children omer/example-branch --limit 20
linctl issue vcs-branch-search documents omer/example-branch --limit 20
linctl issue vcs-branch-search former-attachments omer/example-branch --limit 20
linctl issue vcs-branch-search comments omer/example-branch --limit 20
linctl issue vcs-branch-search former-needs omer/example-branch --limit 20
linctl issue vcs-branch-search history omer/example-branch --limit 20
linctl issue vcs-branch-search inverse-relations omer/example-branch --limit 20
linctl issue vcs-branch-search labels omer/example-branch --limit 20
linctl issue vcs-branch-search needs omer/example-branch --limit 20
linctl issue vcs-branch-search relations omer/example-branch --limit 20
linctl issue vcs-branch-search releases omer/example-branch --limit 20
linctl issue vcs-branch-search shared-access omer/example-branch
linctl issue vcs-branch-search state-history omer/example-branch --limit 20
linctl issue vcs-branch-search subscribers omer/example-branch --limit 20
linctl issue-relation list --limit 20
linctl issue-relation get issue-relation-id
linctl issue pr LIT-123
linctl issue comments LIT-123 --limit 20
linctl comment list --limit 20
linctl comment get comment-id
linctl comment bot-actor comment-id
linctl comment children comment-id --limit 20
linctl comment created-issues comment-id --limit 20
linctl issue start LIT-123
linctl issue reply LIT-123 comment-id --body "thread reply"
linctl issue usage
linctl cycle list --limit 20
linctl cycle get cycle-id
linctl cycle issues cycle-id --limit 20
linctl cycle uncompleted-issues cycle-id --limit 20
linctl cycle create --starts-at 2026-07-01T00:00:00Z --ends-at 2026-07-15T00:00:00Z --name "Planning"
linctl cycle update cycle-id --name "Updated planning"
linctl cycle archive cycle-id
linctl sprint current
linctl sprint report cycle-id --limit 20
linctl project attachments project-id --limit 20
linctl project documents project-id --limit 20
linctl project external-links project-id --limit 20
linctl project history project-id --limit 20
linctl project initiative-links project-id --limit 20
linctl project initiatives project-id --limit 20
linctl project inverse-relations project-id --limit 20
linctl project issues project-id --limit 20
linctl project comments project-id --limit 20
linctl project labels project-id --limit 20
linctl project members project-id --limit 20
linctl project needs project-id --limit 20
linctl project relations project-id --limit 20
linctl project teams project-id --limit 20
linctl project updates project-id --limit 20
linctl project filter-suggestion "started projects"
linctl project-update list --limit 20
linctl project-update get project-update-id
linctl project-update comments project-update-id --limit 20
linctl project-milestone all --limit 20
linctl project-milestone list project-id --limit 20
linctl project-milestone get project-milestone-id
linctl project-milestone issues project-milestone-id --limit 20
linctl project-milestone create project-id --name "Launch milestone"
linctl project-milestone update project-milestone-id --target-date 2026-06-30
linctl project-status list --limit 20
linctl project-status get project-status-id
linctl project-label list --limit 20
linctl project-label get project-label-id
linctl project-label children project-label-id --limit 20
linctl project-label projects project-label-id --limit 20
linctl project-relation list --limit 20
linctl project-relation get project-relation-id
linctl document list --limit 20
linctl document get document-id
linctl document comments document-id --limit 20
linctl label list --limit 20
linctl label get label-id
linctl label children label-id --limit 20
linctl label issues label-id --limit 20
linctl team list --limit 20
linctl team get team-id
linctl team cycles team-id --limit 20
linctl team issues team-id --limit 20
linctl team labels team-id --limit 20
linctl team members team-id --limit 20
linctl team memberships team-id --limit 20
linctl team projects team-id --limit 20
linctl team release-pipelines team-id --limit 20
linctl team states team-id --limit 20
linctl team git-automation-states team-id --limit 20
linctl team templates team-id --limit 20
linctl team-membership list --limit 20
linctl team-membership get team-membership-id
linctl user list --limit 20
linctl user get user-id
linctl user me
linctl user drafts --limit 20
linctl user settings get
linctl user settings notification-categories
linctl user settings notification-category assignments
linctl user settings notification-channels
linctl user settings notification-delivery
linctl user settings mobile-delivery
linctl user settings mobile-schedule
linctl user settings mobile-schedule-day monday
linctl user settings theme --device-type desktop --mode light
linctl user settings custom-theme --device-type desktop --mode light
linctl user settings custom-sidebar-theme --device-type desktop --mode light
linctl user assigned-issues user-id --limit 20
linctl user created-issues user-id --limit 20
linctl user delegated-issues user-id --limit 20
linctl user team-memberships user-id --limit 20
linctl user teams user-id --limit 20
linctl user my-assigned-issues --limit 20
linctl user my-created-issues --limit 20
linctl user my-delegated-issues --limit 20
linctl user my-team-memberships --limit 20
linctl user my-teams --limit 20
linctl workflow-state list --limit 20
linctl workflow-state get workflow-state-id
linctl workflow-state issues workflow-state-id --limit 20
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
linctl initiative documents initiative-id --limit 20
linctl initiative projects initiative-id --limit 20
linctl initiative-relation list --limit 20
linctl initiative-relation get initiative-relation-id
linctl initiative-to-project list --limit 20
linctl initiative-to-project get initiative-to-project-id
linctl initiative-update list --limit 20
linctl initiative-update get initiative-update-id
linctl initiative-update comments initiative-update-id --limit 20
linctl roadmap list --limit 20
linctl roadmap get roadmap-id
linctl roadmap projects roadmap-id --limit 20
linctl roadmap-to-project list --limit 20
linctl roadmap-to-project get roadmap-to-project-id
linctl custom-view list --limit 20
linctl custom-view subscribers custom-view-id
linctl custom-view get custom-view-id
linctl custom-view initiatives custom-view-id --limit 20
linctl custom-view issues custom-view-id --limit 20
linctl custom-view organization-preferences custom-view-id
linctl custom-view organization-preferences values custom-view-id
linctl custom-view projects custom-view-id --limit 20
linctl custom-view user-preferences custom-view-id
linctl custom-view user-preferences values custom-view-id
linctl custom-view preference-values custom-view-id
linctl customer list --limit 20
linctl customer get customer-id
linctl customer-need list --limit 20
linctl customer-need get customer-need-id
linctl customer-need project-attachment customer-need-id
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
linctl attachment issue get attachment-id
linctl attachment issue attachments attachment-id --limit 20
linctl attachment issue bot-actor attachment-id
linctl attachment issue children attachment-id --limit 20
linctl attachment issue comments attachment-id --limit 20
linctl attachment issue documents attachment-id --limit 20
linctl attachment issue former-attachments attachment-id --limit 20
linctl attachment issue former-needs attachment-id --limit 20
linctl attachment issue history attachment-id --limit 20
linctl attachment issue inverse-relations attachment-id --limit 20
linctl attachment issue labels attachment-id --limit 20
linctl attachment issue needs attachment-id --limit 20
linctl attachment issue relations attachment-id --limit 20
linctl attachment issue releases attachment-id --limit 20
linctl attachment issue shared-access attachment-id
linctl attachment issue state-history attachment-id --limit 20
linctl attachment issue subscribers attachment-id --limit 20
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
linctl issue attachments LIT-123 --limit 20
linctl issue children LIT-123 --limit 20
linctl issue documents LIT-123 --limit 20
linctl issue former-attachments LIT-123 --limit 20
linctl issue history LIT-123 --limit 20
linctl issue inverse-relations LIT-123 --limit 20
linctl issue labels LIT-123 --limit 20
linctl issue relations LIT-123 --limit 20
linctl issue releases LIT-123 --limit 20
linctl issue bot-actor LIT-123
linctl issue state-history LIT-123 --limit 20
linctl issue subscribers LIT-123 --limit 20
linctl issue pr LIT-123
linctl next --dry-run
linctl application info app-client-id
linctl agent-activity list --limit 20
linctl agent-activity get agent-activity-id
linctl agent-skill list --limit 20
linctl agent-skill get agent-skill-id
linctl audit-entry types
linctl organization exists kyanite
linctl organization labels --limit 20
linctl organization project-labels --limit 20
linctl organization teams --limit 20
linctl organization templates --limit 20
linctl organization users --limit 20
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
linctl search documents "agent search" --limit 20
linctl search issues "agent search" --limit 20
linctl search projects "agent search" --limit 20
linctl release-pipeline list --limit 20
linctl release-pipeline get release-pipeline-id
linctl release-pipeline releases release-pipeline-id --limit 20
linctl release-pipeline stages release-pipeline-id --limit 20
linctl release-pipeline teams release-pipeline-id --limit 20
linctl release-stage list --limit 20
linctl release-stage get release-stage-id
linctl release-stage releases release-stage-id --limit 20
linctl release list --limit 20
linctl release search "mobile" --limit 20
linctl release get release-id
linctl release history release-id --limit 20
linctl release documents release-id --limit 20
linctl release issues release-id --limit 20
linctl release links release-id --limit 20
linctl external-link get external-link-id
linctl release-note list --limit 20
linctl release-note get release-note-id
linctl issue-to-release list --limit 20
linctl issue-to-release get issue-to-release-id
linctl external-user list --limit 20
linctl external-user get external-user-id
linctl cycle list --limit 20
linctl cycle get cycle-id
linctl cycle issues cycle-id --limit 20
linctl cycle uncompleted-issues cycle-id --limit 20
linctl cycle create --starts-at 2026-07-01T00:00:00Z --ends-at 2026-07-15T00:00:00Z --name "Planning"
linctl cycle update cycle-id --name "Updated planning"
linctl cycle archive cycle-id
linctl sprint current
linctl sprint report cycle-id --limit 20
linctl project attachments project-id --limit 20
linctl project documents project-id --limit 20
linctl project external-links project-id --limit 20
linctl project history project-id --limit 20
linctl project initiative-links project-id --limit 20
linctl project initiatives project-id --limit 20
linctl project inverse-relations project-id --limit 20
linctl project issues project-id --limit 20
linctl project labels project-id --limit 20
linctl project members project-id --limit 20
linctl project needs project-id --limit 20
linctl project relations project-id --limit 20
linctl project teams project-id --limit 20
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
linctl issue-relation list --limit 20
linctl issue-relation get issue-relation-id
linctl comment list --limit 20
linctl comment get comment-id
linctl comment bot-actor comment-id
linctl comment children comment-id --limit 20
linctl comment created-issues comment-id --limit 20
linctl project-milestone create project-id --name "Launch milestone"
linctl project-milestone update project-milestone-id --name "Renamed milestone"
linctl project-update list --limit 20
linctl project-update get project-update-id
linctl project-status list --limit 20
linctl project-status get project-status-id
linctl project-label list --limit 20
linctl project-label get project-label-id
linctl project-label children project-label-id --limit 20
linctl project-label projects project-label-id --limit 20
linctl project-relation list --limit 20
linctl project-relation get project-relation-id
linctl document list --limit 20
linctl label list --limit 20
linctl label children label-id --limit 20
linctl label issues label-id --limit 20
linctl team cycles team-id --limit 20
linctl team issues team-id --limit 20
linctl team labels team-id --limit 20
linctl team members team-id --limit 20
linctl team memberships team-id --limit 20
linctl team projects team-id --limit 20
linctl team release-pipelines team-id --limit 20
linctl team states team-id --limit 20
linctl team git-automation-states team-id --limit 20
linctl team templates team-id --limit 20
linctl team-membership list --limit 20
linctl team-membership get team-membership-id
linctl user me
linctl user drafts --limit 20
linctl user settings get
linctl user settings notification-categories
linctl user settings notification-category assignments
linctl user settings notification-channels
linctl user settings notification-delivery
linctl user settings mobile-delivery
linctl user settings mobile-schedule
linctl user settings mobile-schedule-day monday
linctl user settings theme --device-type desktop --mode light
linctl user settings custom-theme --device-type desktop --mode light
linctl user settings custom-sidebar-theme --device-type desktop --mode light
linctl user assigned-issues user-id --limit 20
linctl user created-issues user-id --limit 20
linctl user delegated-issues user-id --limit 20
linctl user team-memberships user-id --limit 20
linctl user teams user-id --limit 20
linctl user my-assigned-issues --limit 20
linctl user my-created-issues --limit 20
linctl user my-delegated-issues --limit 20
linctl user my-team-memberships --limit 20
linctl user my-teams --limit 20
linctl workflow-state list --limit 20
linctl workflow-state get workflow-state-id
linctl workflow-state issues workflow-state-id --limit 20
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
linctl initiative documents initiative-id --limit 20
linctl initiative projects initiative-id --limit 20
linctl initiative-relation list --limit 20
linctl initiative-relation get initiative-relation-id
linctl initiative-to-project list --limit 20
linctl initiative-to-project get initiative-to-project-id
linctl initiative-update list --limit 20
linctl initiative-update get initiative-update-id
linctl roadmap list --limit 20
linctl roadmap get roadmap-id
linctl roadmap projects roadmap-id --limit 20
linctl roadmap-to-project list --limit 20
linctl roadmap-to-project get roadmap-to-project-id
linctl custom-view list --limit 20
linctl custom-view subscribers custom-view-id
linctl custom-view get custom-view-id
linctl custom-view initiatives custom-view-id --limit 20
linctl custom-view issues custom-view-id --limit 20
linctl custom-view organization-preferences custom-view-id
linctl custom-view organization-preferences values custom-view-id
linctl custom-view projects custom-view-id --limit 20
linctl custom-view user-preferences custom-view-id
linctl custom-view user-preferences values custom-view-id
linctl custom-view preference-values custom-view-id
linctl customer list --limit 20
linctl customer get customer-id
linctl customer-need list --limit 20
linctl customer-need get customer-need-id
linctl customer-need project-attachment customer-need-id
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
linctl attachment issue get attachment-id
linctl attachment issue attachments attachment-id --limit 20
linctl attachment issue bot-actor attachment-id
linctl attachment issue children attachment-id --limit 20
linctl attachment issue comments attachment-id --limit 20
linctl attachment issue documents attachment-id --limit 20
linctl attachment issue former-attachments attachment-id --limit 20
linctl attachment issue former-needs attachment-id --limit 20
linctl attachment issue history attachment-id --limit 20
linctl attachment issue inverse-relations attachment-id --limit 20
linctl attachment issue labels attachment-id --limit 20
linctl attachment issue needs attachment-id --limit 20
linctl attachment issue relations attachment-id --limit 20
linctl attachment issue releases attachment-id --limit 20
linctl attachment issue shared-access attachment-id
linctl attachment issue state-history attachment-id --limit 20
linctl attachment issue subscribers attachment-id --limit 20
linctl --fail-on-empty --sort title --order asc issue list
linctl --format minimal issue get LIT-123
```

Issue, project, Cycle, and ProjectMilestone writes require a pinned target. Team-scoped creates compare
org/team; resource-scoped updates and archives resolve the resource first and compare the pinned project
when configured. Application, AgentActivity, AgentSkill, ExternalUser, AuditEntry, Organization, rate-limit, notification, release-pipeline, release-stage, release, release-note, external-link, issue-to-release, semantic-search, typed search, comment, IssueRelation, ProjectUpdate, ProjectRelation, document, label, team, TeamMembership, user, workflow-state, time-schedule, TriageResponsibility, SLA configuration, template, initiative, initiative-relation, initiative-to-project, initiative-update, roadmap, roadmap-to-project, custom-view, customer, customer-need, customer-status, customer-tier, favorite, emoji, and attachment commands are read-only in the current CLI.

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
