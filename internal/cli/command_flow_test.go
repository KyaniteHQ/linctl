package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_CommandFlows_execute_read_and_write_commands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
		fake     commandFlowFakeClient
	}{
		{name: "target", args: []string{"target"}, contains: "org org-id team LIT/team-id project project-id confirmed true"},
		{name: "doctor", args: []string{"doctor"}, contains: "config ok\n token set\n target confirmed LIT/team-id project project-id"},
		{name: "whoami", args: []string{"whoami"}, contains: "Omer <omer@example.com>"},
		{name: "application info", args: []string{"application", "info", "app-client-id"}, contains: "app-id Demo App by Kyanite", fake: commandFlowFakeClient{expectedApplicationClientID: "app-client-id"}},
		{name: "agent activity list", args: []string{"agent-activity", "list", "--limit", "1"}, contains: "agent-activity-id session agent-session-id [action] signal continue"},
		{name: "agent activity get", args: []string{"agent-activity", "get", "agent-activity-id"}, contains: "agent-activity-id session agent-session-id [action] signal continue"},
		{name: "agent skill list", args: []string{"agent-skill", "list", "--limit", "1"}, contains: "agent-skill-id Triage Helper shared true recent 3"},
		{name: "agent skill get", args: []string{"agent-skill", "get", "agent-skill-id"}, contains: "agent-skill-id Triage Helper shared true recent 3"},
		{name: "external user list", args: []string{"external-user", "list", "--limit", "1"}, contains: "external-user-id External User @external last_seen 2026-06-19T12:00:00Z"},
		{name: "external user get", args: []string{"external-user", "get", "external-user-id"}, contains: "external-user-id External User @external last_seen 2026-06-19T12:00:00Z"},
		{name: "audit entry types", args: []string{"audit-entry", "types"}, contains: "user_login User logged in"},
		{name: "organization exists", args: []string{"organization", "exists", "kyanite"}, contains: "kyanite exists true success true", fake: commandFlowFakeClient{expectedOrganizationURLKey: "kyanite"}},
		{name: "organization labels", args: []string{"organization", "labels", "--limit", "1"}, contains: "label-id Bug #ff0000"},
		{name: "organization project labels", args: []string{"organization", "project-labels", "--limit", "1"}, contains: "project-label-id Roadmap #f2c94c"},
		{name: "organization teams", args: []string{"organization", "teams", "--limit", "1"}, contains: "team-id LIT linctl"},
		{name: "organization templates", args: []string{"organization", "templates", "--limit", "1"}, contains: "template-id Bug report [issue] team LIT"},
		{name: "organization users", args: []string{"organization", "users", "--limit", "1"}, contains: "user-id Omer <omer@example.com>"},
		{name: "rate limit status", args: []string{"rate-limit", "status"}, contains: "api api-key\ncomplexity remaining 900/1000 reset 1720000000000"},
		{name: "notification list", args: []string{"notification", "list", "--limit", "1"}, contains: "notification-id issueMention [mentions] Mentioned you"},
		{name: "notification get", args: []string{"notification", "get", "notification-id"}, contains: "notification-id issueMention [mentions] Mentioned you"},
		{name: "notification subscription list", args: []string{"notification", "subscription", "list", "--limit", "1"}, contains: "notification-subscription-id project Roadmap active true"},
		{name: "notification subscription get", args: []string{"notification", "subscription", "get", "notification-subscription-id"}, contains: "notification-subscription-id project Roadmap active true"},
		{name: "triage responsibility list", args: []string{"triage-responsibility", "list", "--limit", "1"}, contains: "triage-responsibility-id team LIT action notify current Omer"},
		{name: "triage responsibility get", args: []string{"triage-responsibility", "get", "triage-responsibility-id"}, contains: "triage-responsibility-id team LIT action notify current Omer"},
		{name: "triage responsibility manual selection", args: []string{"triage-responsibility", "manual-selection", "triage-responsibility-id"}, contains: "triage-responsibility-id manual users user-id,other-user-id"},
		{name: "SLA configuration list", args: []string{"sla-configuration", "list", "team-id"}, contains: "sla-configuration-id First response sla 3600000 type all removes false"},
		{name: "semantic search", args: []string{"semantic-search", "agent search", "--limit", "2"}, contains: "issue issue-id LIT-3 Search result", fake: commandFlowFakeClient{expectedSemanticSearchQuery: "agent search"}},
		{name: "search documents", args: []string{"search", "documents", "agent search", "--limit", "1"}, contains: "search-document-id Search spec [team]", fake: commandFlowFakeClient{expectedTypedSearchTerm: "agent search"}},
		{name: "search issues", args: []string{"search", "issues", "agent search", "--limit", "1"}, contains: "LIT-30 Search issue [Todo]", fake: commandFlowFakeClient{expectedTypedSearchTerm: "agent search"}},
		{name: "search projects", args: []string{"search", "projects", "agent search", "--limit", "1"}, contains: "search-project-id Search project [Backlog]", fake: commandFlowFakeClient{expectedTypedSearchTerm: "agent search"}},
		{name: "release pipeline list", args: []string{"release-pipeline", "list", "--limit", "1"}, contains: "release-pipeline-id Production production releases 4"},
		{name: "release pipeline get", args: []string{"release-pipeline", "get", "release-pipeline-id"}, contains: "release-pipeline-id Production production releases 4"},
		{name: "release pipeline releases", args: []string{"release-pipeline", "releases", "release-pipeline-id", "--limit", "1"}, contains: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3"},
		{name: "release pipeline stages", args: []string{"release-pipeline", "stages", "release-pipeline-id", "--limit", "1"}, contains: "release-stage-id Started [started] pipeline Production"},
		{name: "release pipeline teams", args: []string{"release-pipeline", "teams", "release-pipeline-id", "--limit", "1"}, contains: "team-id LIT linctl"},
		{name: "release stage list", args: []string{"release-stage", "list", "--limit", "1"}, contains: "release-stage-id Started [started] pipeline Production"},
		{name: "release stage get", args: []string{"release-stage", "get", "release-stage-id"}, contains: "release-stage-id Started [started] pipeline Production"},
		{name: "release stage releases", args: []string{"release-stage", "releases", "release-stage-id", "--limit", "1"}, contains: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3"},
		{name: "release list", args: []string{"release", "list", "--limit", "1"}, contains: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3"},
		{name: "release search", args: []string{"release", "search", "mobile", "--limit", "1"}, contains: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3", fake: commandFlowFakeClient{expectedReleaseSearchTerm: "mobile"}},
		{name: "release get", args: []string{"release", "get", "release-id"}, contains: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3"},
		{name: "release history", args: []string{"release", "history", "release-id", "--limit", "1"}, contains: "release-history-id release release-id entries 1"},
		{name: "release documents", args: []string{"release", "documents", "release-id", "--limit", "1"}, contains: "document-id Spec [project]"},
		{name: "release issues", args: []string{"release", "issues", "release-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "release links", args: []string{"release", "links", "release-id", "--limit", "1"}, contains: "release-link-id Runbook https://example.com/runbook order 1.5"},
		{name: "external link get", args: []string{"external-link", "get", "release-link-id"}, contains: "release-link-id Runbook https://example.com/runbook order 1.5"},
		{name: "release note list", args: []string{"release-note", "list", "--limit", "1"}, contains: "release-note-id Launch notes pipeline Production releases 2"},
		{name: "release note get", args: []string{"release-note", "get", "release-note-id"}, contains: "release-note-id Launch notes pipeline Production releases 2"},
		{name: "issue to release list", args: []string{"issue-to-release", "list", "--limit", "1"}, contains: "issue-to-release-id issue issue-id -> release release-id"},
		{name: "issue to release get", args: []string{"issue-to-release", "get", "issue-to-release-id"}, contains: "issue-to-release-id issue issue-id -> release release-id"},
		{name: "next dry run", args: []string{"next", "--dry-run"}, contains: "LIT-27 Next issue [Todo]"},
		{name: "issue list", args: []string{"issue", "list", "--limit", "1"}, contains: "LIT-1 Listed issue [Todo]"},
		{name: "issue list state filter", args: []string{"issue", "list", "--state", "started", "--limit", "1"}, contains: "LIT-2 Started issue [Started]", fake: commandFlowFakeClient{expectedStateType: "started"}},
		{name: "issue list project filter", args: []string{"issue", "list", "--project", "project-id", "--limit", "1"}, contains: "LIT-4 Project issue [Todo]", fake: commandFlowFakeClient{expectedProjectID: "project-id"}},
		{name: "issue list mine filter", args: []string{"issue", "list", "--mine", "--limit", "1"}, contains: "LIT-5 Mine issue [Todo]", fake: commandFlowFakeClient{expectedAssigneeID: "user-id"}},
		{name: "issue list assignee filter", args: []string{"issue", "list", "--assignee", "assignee-id", "--limit", "1"}, contains: "LIT-6 Assigned issue [Todo]", fake: commandFlowFakeClient{expectedAssigneeID: "assignee-id"}},
		{name: "issue list label filter", args: []string{"issue", "list", "--label", "label-id", "--limit", "1"}, contains: "LIT-7 Labeled issue [Todo]", fake: commandFlowFakeClient{expectedLabelID: "label-id"}},
		{name: "issue list cycle filter", args: []string{"issue", "list", "--cycle", "cycle-id", "--limit", "1"}, contains: "LIT-8 Cycle issue [Todo]", fake: commandFlowFakeClient{expectedCycleID: "cycle-id"}},
		{name: "issue list created-after filter", args: []string{"issue", "list", "--created-after", "2026-06-01", "--limit", "1"}, contains: "LIT-9 Recent issue [Todo]", fake: commandFlowFakeClient{expectedCreatedAfter: "2026-06-01"}},
		{name: "issue list created-since filter", args: []string{"issue", "list", "--created-since", "2026-06-01", "--limit", "1"}, contains: "LIT-9 Recent issue [Todo]", fake: commandFlowFakeClient{expectedCreatedAfter: "2026-06-01"}},
		{name: "issue list created-before filter", args: []string{"issue", "list", "--created-before", "2026-06-30", "--limit", "1"}, contains: "LIT-19 Older issue [Todo]", fake: commandFlowFakeClient{expectedCreatedBefore: "2026-06-30"}},
		{name: "issue list has blockers filter", args: []string{"issue", "list", "--has-blockers", "--limit", "1"}, contains: "LIT-21 Blocked issue [Todo]"},
		{name: "issue list blocks filter", args: []string{"issue", "list", "--blocks", "--limit", "1"}, contains: "LIT-22 Blocking issue [Todo]"},
		{name: "issue list blocked by filter", args: []string{"issue", "list", "--blocked-by", "LIT-1", "--limit", "1"}, contains: "LIT-23 Blocked by issue [Todo]", fake: commandFlowFakeClient{expectedBlockedBy: "LIT-1"}},
		{name: "issue list all teams", args: []string{"issue", "list", "--all-teams", "--limit", "1"}, contains: "LIT-20 All-team issue [Todo]"},
		{name: "issue search", args: []string{"issue", "search", "needle", "--limit", "1"}, contains: "LIT-3 Search result [Todo]", fake: commandFlowFakeClient{expectedSearchQuery: "needle"}},
		{name: "issue vcs branch search get", args: []string{"issue", "vcs-branch-search", "get", "omer/branch", "--json"}, contains: `"identifier": "LIT-40"`},
		{name: "issue vcs branch attachments", args: []string{"issue", "vcs-branch-search", "attachments", "omer/branch", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "issue vcs branch bot actor", args: []string{"issue", "vcs-branch-search", "bot-actor", "omer/branch"}, contains: "issue-id bot bot-actor-id GitHub [github]"},
		{name: "issue vcs branch children", args: []string{"issue", "vcs-branch-search", "children", "omer/branch", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "issue vcs branch documents", args: []string{"issue", "vcs-branch-search", "documents", "omer/branch", "--limit", "1"}, contains: "document-id Spec [issue]"},
		{name: "issue vcs branch former attachments", args: []string{"issue", "vcs-branch-search", "former-attachments", "omer/branch", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "issue vcs branch history", args: []string{"issue", "vcs-branch-search", "history", "omer/branch", "--limit", "1"}, contains: "issue-history-id issue issue-id updated_description true"},
		{name: "issue vcs branch inverse relations", args: []string{"issue", "vcs-branch-search", "inverse-relations", "omer/branch", "--limit", "1"}, contains: "issue-relation-id blocks LIT-1 -> LIT-2"},
		{name: "issue vcs branch labels", args: []string{"issue", "vcs-branch-search", "labels", "omer/branch", "--limit", "1"}, contains: "label-id Bug #ff0000"},
		{name: "issue vcs branch relations", args: []string{"issue", "vcs-branch-search", "relations", "omer/branch", "--limit", "1"}, contains: "issue-relation-id blocks LIT-1 -> LIT-2"},
		{name: "issue vcs branch releases", args: []string{"issue", "vcs-branch-search", "releases", "omer/branch", "--limit", "1"}, contains: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3"},
		{name: "issue vcs branch state history", args: []string{"issue", "vcs-branch-search", "state-history", "omer/branch", "--limit", "1"}, contains: "issue-state-span-id Started started 2026-06-19T12:00:00Z -> -"},
		{name: "issue vcs branch subscribers", args: []string{"issue", "vcs-branch-search", "subscribers", "omer/branch", "--limit", "1"}, contains: "user-id Omer <omer@example.com>"},
		{name: "issue get", args: []string{"issue", "get", "LIT-1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "issue deps", args: []string{"issue", "deps", "LIT-1", "--limit", "2"}, contains: "blocked_by:\nLIT-24 Blocker issue [Todo]", fake: commandFlowFakeClient{expectedIssueDeps: "LIT-1"}},
		{name: "issue attachments", args: []string{"issue", "attachments", "LIT-1", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "issue bot actor", args: []string{"issue", "bot-actor", "LIT-1"}, contains: "issue-id bot bot-actor-id GitHub [github]"},
		{name: "issue children", args: []string{"issue", "children", "LIT-1", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "issue documents", args: []string{"issue", "documents", "LIT-1", "--limit", "1"}, contains: "document-id Spec [issue]"},
		{name: "issue former attachments", args: []string{"issue", "former-attachments", "LIT-1", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "issue history", args: []string{"issue", "history", "LIT-1", "--limit", "1"}, contains: "issue-history-id issue issue-id updated_description true"},
		{name: "issue inverse relations", args: []string{"issue", "inverse-relations", "LIT-1", "--limit", "1"}, contains: "issue-relation-id blocks LIT-1 -> LIT-2"},
		{name: "issue labels", args: []string{"issue", "labels", "LIT-1", "--limit", "1"}, contains: "label-id Bug #ff0000"},
		{name: "issue relations", args: []string{"issue", "relations", "LIT-1", "--limit", "1"}, contains: "issue-relation-id blocks LIT-1 -> LIT-2"},
		{name: "issue releases", args: []string{"issue", "releases", "LIT-1", "--limit", "1"}, contains: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3"},
		{name: "issue state history", args: []string{"issue", "state-history", "LIT-1", "--limit", "1"}, contains: "issue-state-span-id Started started 2026-06-19T12:00:00Z -> -"},
		{name: "issue subscribers", args: []string{"issue", "subscribers", "LIT-1", "--limit", "1"}, contains: "user-id Omer <omer@example.com>"},
		{name: "issue relation list", args: []string{"issue-relation", "list", "--limit", "1"}, contains: "issue-relation-id blocks LIT-1 -> LIT-2"},
		{name: "issue relation get", args: []string{"issue-relation", "get", "issue-relation-id"}, contains: "issue-relation-id blocks LIT-1 -> LIT-2"},
		{name: "issue pr", args: []string{"issue", "pr", "LIT-1"}, contains: `gh pr create --title "LIT-1 Detail issue" --body "https://linear.app/kyanite/issue/LIT-1"`},
		{name: "issue create", args: []string{"issue", "create", "--title", "Created issue"}, contains: "LIT-2 Created issue [Todo]"},
		{name: "issue update", args: []string{"issue", "update", "LIT-1", "--title", "Updated issue"}, contains: "LIT-1 Updated issue [Todo]"},
		{name: "issue update append", args: []string{"issue", "update", "LIT-1", "--append", "Progress note"}, contains: "LIT-1 Updated issue [Todo]", fake: commandFlowFakeClient{expectedUpdateDescription: "Existing description\n\nProgress note"}},
		{name: "issue start", args: []string{"issue", "start", "LIT-1"}, contains: "LIT-1 Started issue [Started]", fake: commandFlowFakeClient{expectedStartAssigneeID: "user-id", expectedStartStateID: "started-state"}},
		{name: "issue comment", args: []string{"issue", "comment", "LIT-1", "--body", "Looks good"}, contains: "comment comment-id on LIT-1"},
		{name: "issue reply", args: []string{"issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"}, contains: "comment comment-id on LIT-1", fake: commandFlowFakeClient{expectedCommentBody: "Reply body", expectedCommentParentID: "comment-id"}},
		{name: "issue comments", args: []string{"issue", "comments", "LIT-1", "--limit", "1"}, contains: "comment-id Omer First comment"},
		{name: "comment list", args: []string{"comment", "list", "--limit", "1"}, contains: "comment-id Omer First comment"},
		{name: "comment get", args: []string{"comment", "get", "comment-id"}, contains: "comment-id Omer First comment"},
		{name: "comment bot actor", args: []string{"comment", "bot-actor", "comment-id"}, contains: "comment-id bot bot-actor-id GitHub [github]"},
		{name: "comment children", args: []string{"comment", "children", "comment-id", "--limit", "1"}, contains: "child-comment-id Omer 2026-06-19T12:00:00Z"},
		{name: "comment created issues", args: []string{"comment", "created-issues", "comment-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "issue close", args: []string{"issue", "close", "LIT-1"}, contains: "LIT-1 Closed issue [Done]"},
		{name: "project list", args: []string{"project", "list", "--limit", "1"}, contains: "project-id Listed project [Backlog]"},
		{name: "project get", args: []string{"project", "get", "project-id"}, contains: "project-id Detail project [Backlog]"},
		{name: "project attachments", args: []string{"project", "attachments", "project-id", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "project documents", args: []string{"project", "documents", "project-id", "--limit", "1"}, contains: "document-id Spec [project]"},
		{name: "project external links", args: []string{"project", "external-links", "project-id", "--limit", "1"}, contains: "release-link-id Runbook https://example.com/runbook order 1.5"},
		{name: "project history", args: []string{"project", "history", "project-id", "--limit", "1"}, contains: "project-history-id project project-id entries 1"},
		{name: "project initiative links", args: []string{"project", "initiative-links", "project-id", "--limit", "1"}, contains: "initiative-to-project-id Platform -> Pinned project order 1"},
		{name: "project initiatives", args: []string{"project", "initiatives", "project-id", "--limit", "1"}, contains: "initiative-id Platform [Active]"},
		{name: "project inverse relations", args: []string{"project", "inverse-relations", "project-id", "--limit", "1"}, contains: "project-relation-id blocks Pinned project -> Related project"},
		{name: "project issues", args: []string{"project", "issues", "project-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "project labels", args: []string{"project", "labels", "project-id", "--limit", "1"}, contains: "project-label-id Roadmap #f2c94c"},
		{name: "project members", args: []string{"project", "members", "project-id", "--limit", "1"}, contains: "user-id Omer"},
		{name: "project needs", args: []string{"project", "needs", "project-id", "--limit", "1"}, contains: "customer-need-id Acme LIT-1 priority 1"},
		{name: "project relations", args: []string{"project", "relations", "project-id", "--limit", "1"}, contains: "project-relation-id blocks Pinned project -> Related project"},
		{name: "project teams", args: []string{"project", "teams", "project-id", "--limit", "1"}, contains: "team-id LIT linctl"},
		{name: "project updates", args: []string{"project", "updates", "project-id", "--limit", "1"}, contains: "project-update-id onTrack Omer First update"},
		{name: "project update list", args: []string{"project-update", "list", "--limit", "1"}, contains: "project-update-id onTrack Omer First update"},
		{name: "project update get", args: []string{"project-update", "get", "project-update-id"}, contains: "project-update-id onTrack Omer First update"},
		{name: "project milestone list", args: []string{"project-milestone", "list", "project-id", "--limit", "1"}, contains: "project-milestone-id Launch milestone [next]"},
		{name: "project milestone create", args: []string{"project-milestone", "create", "project-id", "--name", "Created milestone"}, contains: "project-milestone-id Created milestone [next]"},
		{name: "project milestone update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"}, contains: "project-milestone-id Updated milestone [done]"},
		{name: "project status list", args: []string{"project-status", "list", "--limit", "1"}, contains: "project-status-id Backlog [backlog] #bec2c8"},
		{name: "project status get", args: []string{"project-status", "get", "project-status-id"}, contains: "project-status-id Backlog [backlog] #bec2c8"},
		{name: "project label list", args: []string{"project-label", "list", "--limit", "1"}, contains: "project-label-id Roadmap #f2c94c"},
		{name: "project label get", args: []string{"project-label", "get", "project-label-id"}, contains: "project-label-id Roadmap #f2c94c"},
		{name: "project label children", args: []string{"project-label", "children", "project-label-id", "--limit", "1"}, contains: "child-project-label-id Mobile #56ccf2"},
		{name: "project label projects", args: []string{"project-label", "projects", "project-label-id", "--limit", "1"}, contains: "project-id Listed project [Backlog]"},
		{name: "project relation list", args: []string{"project-relation", "list", "--limit", "1"}, contains: "project-relation-id blocks Pinned project -> Related project"},
		{name: "project relation get", args: []string{"project-relation", "get", "project-relation-id"}, contains: "project-relation-id blocks Pinned project -> Related project"},
		{name: "project create", args: []string{"project", "create", "--name", "Created project"}, contains: "project-id Created project [Backlog]"},
		{name: "project update", args: []string{"project", "update", "project-id", "--name", "Updated project"}, contains: "project-id Updated project [Started]"},
		{name: "project archive", args: []string{"project", "archive", "project-id"}, contains: "project-id Archived project [Canceled]"},
		{name: "document list", args: []string{"document", "list", "--limit", "1"}, contains: "document-id Spec [project]"},
		{name: "document get", args: []string{"document", "get", "document-id"}, contains: "document-id Team note [team]"},
		{name: "label list", args: []string{"label", "list", "--limit", "1"}, contains: "label-id Bug #ff0000"},
		{name: "label get", args: []string{"label", "get", "label-id"}, contains: "label-id Bug #ff0000"},
		{name: "label children", args: []string{"label", "children", "label-id", "--limit", "1"}, contains: "child-label-id Mobile #56ccf2"},
		{name: "label issues", args: []string{"label", "issues", "label-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "team list", args: []string{"team", "list", "--limit", "1"}, contains: "team-id LIT linctl"},
		{name: "team get", args: []string{"team", "get", "team-id"}, contains: "team-id LIT linctl"},
		{name: "team cycles", args: []string{"team", "cycles", "team-id", "--limit", "1"}, contains: "cycle-id Planning cycle [active]"},
		{name: "team issues", args: []string{"team", "issues", "team-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "team labels", args: []string{"team", "labels", "team-id", "--limit", "1"}, contains: "label-id Bug #ff0000"},
		{name: "team members", args: []string{"team", "members", "team-id", "--limit", "1"}, contains: "user-id Omer <omer@example.com>"},
		{name: "team memberships", args: []string{"team", "memberships", "team-id", "--limit", "1"}, contains: "team-membership-id LIT Omer owner true order 1.50"},
		{name: "team projects", args: []string{"team", "projects", "team-id", "--limit", "1"}, contains: "project-id Listed project [Backlog]"},
		{name: "team release pipelines", args: []string{"team", "release-pipelines", "team-id", "--limit", "1"}, contains: "release-pipeline-id Production production releases 4"},
		{name: "team states", args: []string{"team", "states", "team-id", "--limit", "1"}, contains: "workflow-state-id Started [started]"},
		{name: "team templates", args: []string{"team", "templates", "team-id", "--limit", "1"}, contains: "template-id Bug report [issue] team LIT"},
		{name: "team membership list", args: []string{"team-membership", "list", "--limit", "1"}, contains: "team-membership-id LIT Omer owner true order 1.50"},
		{name: "team membership get", args: []string{"team-membership", "get", "team-membership-id"}, contains: "team-membership-id LIT Omer owner true order 1.50"},
		{name: "user list", args: []string{"user", "list", "--limit", "1"}, contains: "user-id Omer <omer@example.com>"},
		{name: "user get", args: []string{"user", "get", "user-id"}, contains: "user-id Omer <omer@example.com>"},
		{name: "user me", args: []string{"user", "me"}, contains: "user-id Omer <omer@example.com>"},
		{name: "user drafts", args: []string{"user", "drafts", "--limit", "1"}, contains: "draft-id issue LIT-3 Draft issue"},
		{name: "user assigned issues", args: []string{"user", "assigned-issues", "user-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "user created issues", args: []string{"user", "created-issues", "user-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "user delegated issues", args: []string{"user", "delegated-issues", "user-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "user team memberships", args: []string{"user", "team-memberships", "user-id", "--limit", "1"}, contains: "team-membership-id LIT Omer owner true order 1.50"},
		{name: "user teams", args: []string{"user", "teams", "user-id", "--limit", "1"}, contains: "team-id LIT linctl"},
		{name: "user my assigned issues", args: []string{"user", "my-assigned-issues", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "user my created issues", args: []string{"user", "my-created-issues", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "user my delegated issues", args: []string{"user", "my-delegated-issues", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "user my team memberships", args: []string{"user", "my-team-memberships", "--limit", "1"}, contains: "team-membership-id LIT Omer owner true order 1.50"},
		{name: "user my teams", args: []string{"user", "my-teams", "--limit", "1"}, contains: "team-id LIT linctl"},
		{name: "workflow state list", args: []string{"workflow-state", "list", "--limit", "1"}, contains: "workflow-state-id Started [started]"},
		{name: "workflow state get", args: []string{"workflow-state", "get", "workflow-state-id"}, contains: "workflow-state-id Started [started]"},
		{name: "time schedule list", args: []string{"time-schedule", "list", "--limit", "1"}, contains: "time-schedule-id Primary on-call entries 1"},
		{name: "time schedule get", args: []string{"time-schedule", "get", "time-schedule-id"}, contains: "time-schedule-id Primary on-call entries 1"},
		{name: "template list", args: []string{"template", "list", "--limit", "1"}, contains: "template-id Bug report [issue] team LIT"},
		{name: "template get", args: []string{"template", "get", "template-id"}, contains: "template-id Bug report [issue] team LIT"},
		{name: "initiative list", args: []string{"initiative", "list", "--limit", "1"}, contains: "initiative-id Platform [Active]"},
		{name: "initiative get", args: []string{"initiative", "get", "initiative-id"}, contains: "initiative-id Platform [Active]"},
		{name: "initiative history", args: []string{"initiative", "history", "initiative-id", "--limit", "1"}, contains: "initiative-history-id initiative initiative-id entries 1"},
		{name: "initiative links", args: []string{"initiative", "links", "initiative-id", "--limit", "1"}, contains: "release-link-id Runbook https://example.com/runbook order 1.5"},
		{name: "initiative sub-initiatives", args: []string{"initiative", "sub-initiatives", "initiative-id", "--limit", "1"}, contains: "child-initiative-id Child platform [Planned]"},
		{name: "initiative updates", args: []string{"initiative", "updates", "initiative-id", "--limit", "1"}, contains: "initiative-update-id onTrack Omer First initiative update"},
		{name: "initiative documents", args: []string{"initiative", "documents", "initiative-id", "--limit", "1"}, contains: "document-id Spec [project]"},
		{name: "initiative projects", args: []string{"initiative", "projects", "initiative-id", "--limit", "1"}, contains: "project-id Listed project [Backlog]"},
		{name: "initiative relation list", args: []string{"initiative-relation", "list", "--limit", "1"}, contains: "initiative-relation-id Platform -> Child initiative order 1.50"},
		{name: "initiative relation get", args: []string{"initiative-relation", "get", "initiative-relation-id"}, contains: "initiative-relation-id Platform -> Child initiative order 1.50"},
		{name: "initiative to project list", args: []string{"initiative-to-project", "list", "--limit", "1"}, contains: "initiative-to-project-id Platform -> Pinned project order 1"},
		{name: "initiative to project get", args: []string{"initiative-to-project", "get", "initiative-to-project-id"}, contains: "initiative-to-project-id Platform -> Pinned project order 1"},
		{name: "initiative update list", args: []string{"initiative-update", "list", "--limit", "1"}, contains: "initiative-update-id onTrack Omer First initiative update"},
		{name: "initiative update get", args: []string{"initiative-update", "get", "initiative-update-id"}, contains: "initiative-update-id onTrack Omer First initiative update"},
		{name: "roadmap list", args: []string{"roadmap", "list", "--limit", "1"}, contains: "roadmap-id Platform roadmap platform-roadmap"},
		{name: "roadmap get", args: []string{"roadmap", "get", "roadmap-id"}, contains: "roadmap-id Platform roadmap platform-roadmap"},
		{name: "roadmap to project list", args: []string{"roadmap-to-project", "list", "--limit", "1"}, contains: "roadmap-to-project-id Platform roadmap -> Pinned project order 1"},
		{name: "roadmap to project get", args: []string{"roadmap-to-project", "get", "roadmap-to-project-id"}, contains: "roadmap-to-project-id Platform roadmap -> Pinned project order 1"},
		{name: "custom view list", args: []string{"custom-view", "list", "--limit", "1"}, contains: "custom-view-id My issues [Issue]"},
		{name: "custom view subscribers", args: []string{"custom-view", "subscribers", "custom-view-id"}, contains: "custom-view-id has_subscribers true"},
		{name: "custom view get", args: []string{"custom-view", "get", "custom-view-id"}, contains: "custom-view-id My issues [Issue]"},
		{name: "custom view initiatives", args: []string{"custom-view", "initiatives", "custom-view-id", "--limit", "1"}, contains: "initiative-id Platform [Active]"},
		{name: "custom view issues", args: []string{"custom-view", "issues", "custom-view-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "custom view organization preferences", args: []string{"custom-view", "organization-preferences", "custom-view-id"}, contains: "custom-view-id organization preferences organization customView layout list"},
		{name: "custom view organization preference values", args: []string{"custom-view", "organization-preferences", "values", "custom-view-id"}, contains: "custom-view-id preference values layout list ordering priority"},
		{name: "custom view projects", args: []string{"custom-view", "projects", "custom-view-id", "--limit", "1"}, contains: "project-id Pinned project [Backlog]"},
		{name: "custom view user preferences", args: []string{"custom-view", "user-preferences", "custom-view-id"}, contains: "custom-view-id user preferences user customView layout board"},
		{name: "custom view user preference values", args: []string{"custom-view", "user-preferences", "values", "custom-view-id"}, contains: "custom-view-id preference values layout board ordering updatedAt"},
		{name: "custom view preference values", args: []string{"custom-view", "preference-values", "custom-view-id"}, contains: "custom-view-id preference values layout board ordering updatedAt"},
		{name: "customer list", args: []string{"customer", "list", "--limit", "1"}, contains: "customer-id Acme [Active] needs 3"},
		{name: "customer get", args: []string{"customer", "get", "customer-id"}, contains: "customer-id Acme [Active] needs 3"},
		{name: "customer need list", args: []string{"customer-need", "list", "--limit", "1"}, contains: "customer-need-id Acme LIT-1 priority 1"},
		{name: "customer need get", args: []string{"customer-need", "get", "customer-need-id"}, contains: "customer-need-id Acme LIT-1 priority 1"},
		{name: "customer status list", args: []string{"customer-status", "list", "--limit", "1"}, contains: "customer-status-id Active #00ff00 1"},
		{name: "customer status get", args: []string{"customer-status", "get", "customer-status-id"}, contains: "customer-status-id Active #00ff00 1"},
		{name: "customer tier list", args: []string{"customer-tier", "list", "--limit", "1"}, contains: "customer-tier-id Enterprise #0000ff 2"},
		{name: "customer tier get", args: []string{"customer-tier", "get", "customer-tier-id"}, contains: "customer-tier-id Enterprise #0000ff 2"},
		{name: "favorite list", args: []string{"favorite", "list", "--limit", "1"}, contains: "favorite-id [issue]"},
		{name: "favorite children", args: []string{"favorite", "children", "favorite-folder-id", "--limit", "1"}, contains: "favorite-child-id [project]"},
		{name: "favorite get", args: []string{"favorite", "get", "favorite-id"}, contains: "favorite-id [issue]"},
		{name: "emoji list", args: []string{"emoji", "list", "--limit", "1"}, contains: "emoji-id party [custom]"},
		{name: "emoji get", args: []string{"emoji", "get", "emoji-id"}, contains: "emoji-id party [custom]"},
		{name: "attachment list", args: []string{"attachment", "list", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "attachment url", args: []string{"attachment", "url", "https://github.com/kyanite/linctl/pull/1", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "attachment get", args: []string{"attachment", "get", "attachment-id"}, contains: "attachment-id Linked PR [github]"},
		{name: "attachment issue get", args: []string{"attachment", "issue", "get", "attachment-id"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "attachment issue attachments", args: []string{"attachment", "issue", "attachments", "attachment-id", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "attachment issue bot actor", args: []string{"attachment", "issue", "bot-actor", "attachment-id"}, contains: "issue-id bot bot-actor-id GitHub [github]"},
		{name: "attachment issue children", args: []string{"attachment", "issue", "children", "attachment-id", "--limit", "1"}, contains: "LIT-1 Detail issue [Todo]"},
		{name: "attachment issue documents", args: []string{"attachment", "issue", "documents", "attachment-id", "--limit", "1"}, contains: "document-id Spec [issue]"},
		{name: "attachment issue former attachments", args: []string{"attachment", "issue", "former-attachments", "attachment-id", "--limit", "1"}, contains: "attachment-id Linked PR [github]"},
		{name: "attachment issue history", args: []string{"attachment", "issue", "history", "attachment-id", "--limit", "1"}, contains: "issue-history-id issue issue-id updated_description true"},
		{name: "attachment issue inverse relations", args: []string{"attachment", "issue", "inverse-relations", "attachment-id", "--limit", "1"}, contains: "issue-relation-id blocks LIT-1 -> LIT-2"},
		{name: "attachment issue labels", args: []string{"attachment", "issue", "labels", "attachment-id", "--limit", "1"}, contains: "label-id Bug #ff0000"},
		{name: "attachment issue relations", args: []string{"attachment", "issue", "relations", "attachment-id", "--limit", "1"}, contains: "issue-relation-id blocks LIT-1 -> LIT-2"},
		{name: "attachment issue releases", args: []string{"attachment", "issue", "releases", "attachment-id", "--limit", "1"}, contains: "release-id Mobile 1.2.3 [v1.2.3] pipeline Production stage Started issues 3"},
		{name: "attachment issue state history", args: []string{"attachment", "issue", "state-history", "attachment-id", "--limit", "1"}, contains: "issue-state-span-id Started started 2026-06-19T12:00:00Z -> -"},
		{name: "attachment issue subscribers", args: []string{"attachment", "issue", "subscribers", "attachment-id", "--limit", "1"}, contains: "user-id Omer <omer@example.com>"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), test.contains)
		})
	}
}

func Test_CommandFlows_read_issue_comment_body_from_stdin(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{expectedCommentBody: "stdin body\nsecond line"})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetIn(strings.NewReader("stdin body\nsecond line"))
	command.SetOut(&output)
	command.SetArgs([]string{"issue", "comment", "LIT-1", "--body", "-"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "comment comment-id on LIT-1")
}

func Test_CommandFlows_read_issue_text_from_files(t *testing.T) {
	descriptionFile := writeTempTextFile(t, "description from file")
	appendFile := writeTempTextFile(t, "append from file")
	commentFile := writeTempTextFile(t, "comment from file")
	replyFile := writeTempTextFile(t, "reply from file")

	tests := []struct {
		name string
		args []string
		fake commandFlowFakeClient
	}{
		{
			name: "create description",
			args: []string{"issue", "create", "--title", "Created issue", "--description-file", descriptionFile},
			fake: commandFlowFakeClient{expectedCreateDescription: "description from file"},
		},
		{
			name: "update append",
			args: []string{"issue", "update", "LIT-1", "--append-file", appendFile},
			fake: commandFlowFakeClient{expectedUpdateDescription: "Existing description\n\nappend from file"},
		},
		{
			name: "comment body",
			args: []string{"issue", "comment", "LIT-1", "--body-file", commentFile},
			fake: commandFlowFakeClient{expectedCommentBody: "comment from file"},
		},
		{
			name: "reply body",
			args: []string{"issue", "reply", "LIT-1", "comment-id", "--body-file", replyFile},
			fake: commandFlowFakeClient{expectedCommentBody: "reply from file", expectedCommentParentID: "comment-id"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.NotEmpty(t, output.String())
		})
	}
}

func Test_CommandFlows_report_issue_text_file_errors(t *testing.T) {
	textFile := writeTempTextFile(t, "from file")
	missingFile := filepath.Join(t.TempDir(), "missing.md")
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "create description conflict",
			args:     []string{"issue", "create", "--title", "Created issue", "--description", "inline", "--description-file", textFile},
			contains: "description and description-file are mutually exclusive",
		},
		{
			name:     "update description conflict",
			args:     []string{"issue", "update", "LIT-1", "--description", "inline", "--description-file", textFile},
			contains: "description and description-file are mutually exclusive",
		},
		{
			name:     "update append conflict",
			args:     []string{"issue", "update", "LIT-1", "--append", "inline", "--append-file", textFile},
			contains: "append and append-file are mutually exclusive",
		},
		{
			name:     "comment body conflict",
			args:     []string{"issue", "comment", "LIT-1", "--body", "inline", "--body-file", textFile},
			contains: "body and body-file are mutually exclusive",
		},
		{
			name:     "reply body conflict",
			args:     []string{"issue", "reply", "LIT-1", "comment-id", "--body", "inline", "--body-file", textFile},
			contains: "body and body-file are mutually exclusive",
		},
		{
			name:     "missing file",
			args:     []string{"issue", "comment", "LIT-1", "--body-file", missingFile},
			contains: "read body from file",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.contains)
		})
	}
}

func Test_CommandFlows_resolve_current_issue_from_branch(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"current"})

	require.NoError(t, err)
	require.Contains(t, output, "LIT-1 Detail issue [Todo]")
}

func Test_CommandFlows_print_current_issue_as_json(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"--json", "current"})

	require.NoError(t, err)
	require.Contains(t, output, `"identifier": "LIT-1"`)
}

func Test_CommandFlows_print_current_issue_identifier_from_issue_id(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"issue", "id"})

	require.NoError(t, err)
	require.Equal(t, "LIT-1\n", output)
}

func Test_CommandFlows_print_current_issue_title_from_issue_title(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"issue", "title"})

	require.NoError(t, err)
	require.Equal(t, "Detail issue\n", output)
}

func Test_CommandFlows_print_current_issue_url_from_issue_url(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"issue", "url"})

	require.NoError(t, err)
	require.Equal(t, "https://linear.app/kyanite/issue/LIT-1\n", output)
}

func Test_CommandFlows_print_issue_branch_from_issue_branch(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"issue", "branch", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Equal(t, "lit-1-detail-issue\n", output.String())
}

func Test_CommandFlows_print_issue_pr_from_current_branch(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"issue", "pr"})

	require.NoError(t, err)
	require.Contains(t, output, `gh pr create --title "LIT-1 Detail issue"`)
}

func Test_CommandFlows_close_current_issue_from_done(t *testing.T) {
	output, err := runCurrentCommandInGitBranch(t, []string{"done"})

	require.NoError(t, err)
	require.Contains(t, output, "LIT-1 Closed issue [Done]")
}

func Test_CommandFlows_report_next_errors(t *testing.T) {
	t.Run("requires dry run", func(t *testing.T) {
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"next"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "next requires --dry-run")
	})

	t.Run("empty candidate list", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{emptyNextIssues: true})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"next", "--dry-run"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "next issue not found")
	})

	t.Run("empty candidate list with fail on empty", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{emptyNextIssues: true})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--fail-on-empty", "next", "--dry-run"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "empty result")
	})
}

func Test_CommandFlows_rank_next_issue_candidates(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{rankedNextIssues: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"next", "--dry-run"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "LIT-30 Unblocks checkout [Todo]")
}

func Test_CommandFlows_report_current_issue_errors(t *testing.T) {
	t.Run("missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"current"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("done missing issue reference", func(t *testing.T) {
		dir := t.TempDir()
		runGitCommand(t, dir, "init")
		runGitCommand(t, dir, "checkout", "-b", "main")
		t.Chdir(dir)
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"done"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "linear issue reference missing")
	})

	t.Run("runtime failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"current"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("done runtime failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"done"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return commandRuntime{}, errors.New("runtime failed")
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "runtime failed")
	})

	t.Run("issue lookup failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"current"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return testCommandRuntime(commandFlowFakeClient{failOperation: "issue"}), nil
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")
	})

	t.Run("done close failure", func(t *testing.T) {
		_, err := runCurrentCommandInGitBranchWithRuntime(t, []string{"done"}, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
			return testCommandRuntime(commandFlowFakeClient{failOperation: "IssueClose"}), nil
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "close issue LIT-1")
	})
}

func runCurrentCommandInGitBranch(t *testing.T, args []string) (string, error) {
	t.Helper()

	return runCurrentCommandInGitBranchWithRuntime(t, args, func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return testCommandRuntime(commandFlowFakeClient{}), nil
	})
}

func runCurrentCommandInGitBranchWithRuntime(
	t *testing.T,
	args []string,
	runtimeBuilder func(context.Context, *rootOptions) (commandRuntime, error),
) (string, error) {
	t.Helper()

	dir := t.TempDir()
	runGitCommand(t, dir, "init")
	runGitCommand(t, dir, "checkout", "-b", "feature/LIT-1-current")
	t.Chdir(dir)

	output := bytes.Buffer{}
	original := buildCommandRuntime
	buildCommandRuntime = runtimeBuilder
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs(args)

	err := command.ExecuteContext(context.Background())
	return output.String(), err
}

func Test_CommandFlows_report_runtime_and_writer_errors(t *testing.T) {
	t.Run("runtime error returns from command", func(t *testing.T) {
		commands := [][]string{
			{"target"},
			{"doctor"},
			{"whoami"},
			{"application", "info", "app-client-id"},
			{"agent-activity", "list"},
			{"agent-activity", "get", "agent-activity-id"},
			{"agent-skill", "list"},
			{"agent-skill", "get", "agent-skill-id"},
			{"audit-entry", "types"},
			{"triage-responsibility", "list"},
			{"triage-responsibility", "get", "triage-responsibility-id"},
			{"triage-responsibility", "manual-selection", "triage-responsibility-id"},
			{"organization", "exists", "kyanite"},
			{"semantic-search", "agent search"},
			{"search", "documents", "agent search"},
			{"search", "issues", "agent search"},
			{"search", "projects", "agent search"},
			{"rate-limit", "status"},
			{"release", "list"},
			{"release", "search", "mobile"},
			{"release", "get", "release-id"},
			{"external-link", "get", "release-link-id"},
			{"release-note", "list"},
			{"release-note", "get", "release-note-id"},
			{"next", "--dry-run"},
			{"issue", "list"},
			{"issue", "search", "needle"},
			{"issue", "get", "LIT-1"},
			{"issue", "deps", "LIT-1"},
			{"issue", "pr", "LIT-1"},
			{"issue", "create", "--title", "Created issue"},
			{"issue", "update", "LIT-1", "--title", "Updated issue"},
			{"issue", "start", "LIT-1"},
			{"issue", "comment", "LIT-1", "--body", "Looks good"},
			{"issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"},
			{"issue", "comments", "LIT-1"},
			{"issue", "close", "LIT-1"},
			{"project", "list"},
			{"project", "get", "project-id"},
			{"project", "members", "project-id"},
			{"project", "updates", "project-id"},
			{"project-milestone", "list", "project-id"},
			{"project-milestone", "get", "project-milestone-id"},
			{"project-milestone", "create", "project-id", "--name", "Created milestone"},
			{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"},
			{"project", "create", "--name", "Created project"},
			{"project", "update", "project-id", "--name", "Updated project"},
			{"project", "archive", "project-id"},
			{"document", "list"},
			{"document", "get", "document-id"},
			{"label", "list"},
			{"label", "get", "label-id"},
			{"team", "list"},
			{"team", "get", "team-id"},
			{"team", "members", "team-id"},
			{"user", "list"},
			{"user", "get", "user-id"},
			{"user", "me"},
			{"user", "assigned-issues", "user-id"},
			{"user", "created-issues", "user-id"},
			{"user", "delegated-issues", "user-id"},
			{"user", "team-memberships", "user-id"},
			{"user", "teams", "user-id"},
			{"user", "my-assigned-issues"},
			{"user", "my-created-issues"},
			{"user", "my-delegated-issues"},
			{"user", "my-team-memberships"},
			{"user", "my-teams"},
			{"custom-view", "subscribers", "custom-view-id"},
			{"custom-view", "initiatives", "custom-view-id"},
			{"custom-view", "issues", "custom-view-id"},
			{"custom-view", "organization-preferences", "custom-view-id"},
			{"custom-view", "organization-preferences", "values", "custom-view-id"},
			{"custom-view", "projects", "custom-view-id"},
			{"custom-view", "user-preferences", "custom-view-id"},
			{"custom-view", "user-preferences", "values", "custom-view-id"},
			{"custom-view", "preference-values", "custom-view-id"},
			{"sla-configuration", "list", "team-id"},
		}
		for _, args := range commands {
			t.Run(strings.Join(args, " "), func(t *testing.T) {
				original := buildCommandRuntime
				buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, errors.New("runtime failed")
				}
				defer func() {
					buildCommandRuntime = original
				}()
				command := NewRootCommand(context.Background(), BuildInfo{})
				command.SetArgs(args)

				err := command.ExecuteContext(context.Background())

				require.Error(t, err)
				require.Contains(t, err.Error(), "runtime failed")
			})
		}
	})

	t.Run("writeIssues returns writer errors", func(t *testing.T) {
		command := &cobra.Command{}
		command.SetOut(commandFailingWriter{})

		err := writeIssues(command, &rootOptions{}, []client.IssueSummary{{Identifier: "LIT-1", Title: "Broken", State: "Todo"}})

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("doctor returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"doctor"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "list"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project members returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "members", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project updates returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project", "updates", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("project milestone list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"project-milestone", "list", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("release search returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"release", "search", "mobile"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("release search reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "release", "search", "mobile"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("SLA configuration list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"sla-configuration", "list", "team-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("SLA configuration list reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "sla-configuration", "list", "team-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("semantic search returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"semantic-search", "agent search"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("semantic search reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "semantic-search", "agent search"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("typed search returns writer errors", func(t *testing.T) {
		tests := [][]string{
			{"search", "documents", "agent search"},
			{"search", "issues", "agent search"},
			{"search", "projects", "agent search"},
		}
		for _, args := range tests {
			t.Run(strings.Join(args[:2], " "), func(t *testing.T) {
				restore := useCommandRuntime(t, commandFlowFakeClient{})
				defer restore()
				command := NewRootCommand(context.Background(), BuildInfo{})
				command.SetOut(commandFailingWriter{})
				command.SetArgs(args)

				err := command.ExecuteContext(context.Background())

				require.Error(t, err)
				require.Contains(t, err.Error(), "write line")
			})
		}
	})

	t.Run("typed search reports sort errors", func(t *testing.T) {
		tests := [][]string{
			{"--sort", "missing", "search", "documents", "agent search"},
			{"--sort", "missing", "search", "issues", "agent search"},
			{"--sort", "missing", "search", "projects", "agent search"},
		}
		for _, args := range tests {
			t.Run(strings.Join(args[2:4], " "), func(t *testing.T) {
				restore := useCommandRuntime(t, commandFlowFakeClient{})
				defer restore()
				command := NewRootCommand(context.Background(), BuildInfo{})
				command.SetArgs(args)

				err := command.ExecuteContext(context.Background())

				require.Error(t, err)
				require.Contains(t, err.Error(), `sort field "missing" is not present`)
			})
		}
	})

	t.Run("issue child list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"issue", "history", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("issue child list reports sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "issue", "children", "LIT-1"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("document list returns writer errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"document", "list"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})

	t.Run("usage returns writer errors", func(t *testing.T) {
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(commandFailingWriter{})
		command.SetArgs([]string{"usage"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "write line")
	})
}

func Test_CommandFlows_print_json_for_read_and_comment_commands(t *testing.T) {
	tests := [][]string{
		{"--json", "target"},
		{"--json", "doctor"},
		{"--json", "whoami"},
		{"--json", "organization", "exists", "kyanite"},
		{"--json", "--fields", "id,name,type,team_key", "organization", "templates", "--limit", "1"},
		{"--json", "rate-limit", "status"},
		{"--json", "notification", "list", "--limit", "1"},
		{"--json", "notification", "get", "notification-id"},
		{"--json", "notification", "subscription", "list", "--limit", "1"},
		{"--json", "notification", "subscription", "get", "notification-subscription-id"},
		{"--json", "--fields", "id,team_key,action", "triage-responsibility", "list", "--limit", "1"},
		{"--json", "triage-responsibility", "get", "triage-responsibility-id"},
		{"--json", "triage-responsibility", "manual-selection", "triage-responsibility-id"},
		{"--json", "release-pipeline", "list", "--limit", "1"},
		{"--json", "release-pipeline", "get", "release-pipeline-id"},
		{"--json", "--fields", "id,pipeline_id,stage_id", "release-pipeline", "releases", "release-pipeline-id", "--limit", "1"},
		{"--json", "--fields", "id,pipeline_id,type", "release-pipeline", "stages", "release-pipeline-id", "--limit", "1"},
		{"--json", "release-stage", "list", "--limit", "1"},
		{"--json", "release-stage", "get", "release-stage-id"},
		{"--json", "--fields", "id,pipeline_id,stage_id", "release-stage", "releases", "release-stage-id", "--limit", "1"},
		{"--json", "release", "list", "--limit", "1"},
		{"--json", "release", "search", "mobile", "--limit", "1"},
		{"--json", "release", "get", "release-id"},
		{"--json", "--fields", "id,release_id,entry_count", "release", "history", "release-id", "--limit", "1"},
		{"--json", "--fields", "id,label,url", "release", "links", "release-id", "--limit", "1"},
		{"--json", "--fields", "id,label,url", "external-link", "get", "release-link-id"},
		{"--json", "release-note", "list", "--limit", "1"},
		{"--json", "release-note", "get", "release-note-id"},
		{"--json", "--fields", "id,client_id,name", "application", "info", "app-client-id"},
		{"--json", "--fields", "id,content_type,agent_session_id", "agent-activity", "list", "--limit", "1"},
		{"--json", "agent-activity", "get", "agent-activity-id"},
		{"--json", "--fields", "id,title,shared", "agent-skill", "list", "--limit", "1"},
		{"--json", "agent-skill", "get", "agent-skill-id"},
		{"--json", "--fields", "id,display_name,last_seen", "external-user", "list", "--limit", "1"},
		{"--json", "external-user", "get", "external-user-id"},
		{"--json", "--fields", "type,description", "audit-entry", "types"},
		{"--json", "--fields", "type,id,key,title", "semantic-search", "agent search", "--limit", "2"},
		{"--json", "--fields", "id,title,parent_type", "search", "documents", "agent search", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "search", "issues", "agent search", "--limit", "1"},
		{"--json", "--fields", "id,name,slug_id", "search", "projects", "agent search", "--limit", "1"},
		{"--json", "next", "--dry-run"},
		{"--json", "issue", "list", "--limit", "1"},
		{"--json", "issue", "search", "needle", "--limit", "1"},
		{"--json", "issue", "deps", "LIT-1", "--limit", "2"},
		{"--json", "issue", "attachments", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "issue", "children", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,title,parent_type", "issue", "documents", "LIT-1", "--limit", "1"},
		{"--json", "issue", "former-attachments", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,issue_id,updated_description", "issue", "history", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,type,issue_identifier,related_issue_identifier", "issue", "inverse-relations", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,name,color", "issue", "labels", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,type,issue_identifier,related_issue_identifier", "issue", "relations", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,name,version", "issue", "releases", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,type,issue_identifier,related_issue_identifier", "issue-relation", "list", "--limit", "1"},
		{"--json", "issue-relation", "get", "issue-relation-id"},
		{"--json", "--fields", "id,issue_id,release_id", "issue-to-release", "list", "--limit", "1"},
		{"--json", "issue-to-release", "get", "issue-to-release-id"},
		{"--json", "issue", "pr", "LIT-1"},
		{"--json", "issue", "start", "LIT-1"},
		{"--json", "issue", "comment", "LIT-1", "--body", "Looks good"},
		{"--json", "issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"},
		{"--json", "--fields", "id,display_name", "issue", "comments", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "id,display_name", "comment", "list", "--limit", "1"},
		{"--json", "comment", "get", "comment-id"},
		{"--json", "project", "list", "--limit", "1"},
		{"--json", "project", "members", "project-id", "--limit", "1"},
		{"--json", "--fields", "id,health,display_name", "project", "updates", "project-id", "--limit", "1"},
		{"--json", "--fields", "id,health,project_id", "project-update", "list", "--limit", "1"},
		{"--json", "project-update", "get", "project-update-id"},
		{"--json", "--fields", "id,name,status", "project-milestone", "list", "project-id", "--limit", "1"},
		{"--json", "project-status", "list", "--limit", "1"},
		{"--json", "project-status", "get", "project-status-id"},
		{"--json", "project-label", "list", "--limit", "1"},
		{"--json", "project-label", "children", "project-label-id", "--limit", "1"},
		{"--json", "project-label", "projects", "project-label-id", "--limit", "1"},
		{"--json", "--fields", "id,type,project_id,related_project_id", "project-relation", "list", "--limit", "1"},
		{"--json", "project-relation", "get", "project-relation-id"},
		{"--json", "--fields", "id,title,parent_type", "document", "list", "--limit", "1"},
		{"--json", "--fields", "id,name,color", "label", "list", "--limit", "1"},
		{"--json", "--fields", "id,key,name", "team", "list", "--limit", "1"},
		{"--json", "--fields", "id,name,status", "team", "cycles", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "team", "issues", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,name,color", "team", "labels", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,display_name,email", "team", "members", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,team_key,user_id,owner", "team", "memberships", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,name,status", "team", "projects", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,name,type", "team", "release-pipelines", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,name,type", "team", "states", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,name,type,team_key", "team", "templates", "team-id", "--limit", "1"},
		{"--json", "--fields", "id,team_key,user_id,owner", "team-membership", "list", "--limit", "1"},
		{"--json", "team-membership", "get", "team-membership-id"},
		{"--json", "--fields", "id,display_name,email", "user", "list", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "user", "assigned-issues", "user-id", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "user", "created-issues", "user-id", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "user", "delegated-issues", "user-id", "--limit", "1"},
		{"--json", "--fields", "id,team_key,user_id,owner", "user", "team-memberships", "user-id", "--limit", "1"},
		{"--json", "--fields", "id,key,name", "user", "teams", "user-id", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "user", "my-assigned-issues", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "user", "my-created-issues", "--limit", "1"},
		{"--json", "--fields", "id,identifier,title", "user", "my-delegated-issues", "--limit", "1"},
		{"--json", "--fields", "id,team_key,user_id,owner", "user", "my-team-memberships", "--limit", "1"},
		{"--json", "--fields", "id,key,name", "user", "my-teams", "--limit", "1"},
		{"--json", "time-schedule", "list", "--limit", "1"},
		{"--json", "time-schedule", "get", "time-schedule-id"},
		{"--json", "--fields", "id,name,sla_type", "sla-configuration", "list", "team-id"},
		{"--json", "--fields", "id,name,type,team_key", "template", "list", "--limit", "1"},
		{"--json", "template", "get", "template-id"},
		{"--json", "initiative", "list", "--limit", "1"},
		{"--json", "initiative", "get", "initiative-id"},
		{"--json", "--fields", "id,initiative_id,entry_count", "initiative", "history", "initiative-id", "--limit", "1"},
		{"--json", "--fields", "id,label,url", "initiative", "links", "initiative-id", "--limit", "1"},
		{"--json", "--fields", "id,name,status", "initiative", "sub-initiatives", "initiative-id", "--limit", "1"},
		{"--json", "--fields", "id,health,initiative_id", "initiative", "updates", "initiative-id", "--limit", "1"},
		{"--json", "--fields", "id,parent_initiative_id,related_initiative_id", "initiative-relation", "list", "--limit", "1"},
		{"--json", "initiative-relation", "get", "initiative-relation-id"},
		{"--json", "--fields", "id,initiative_id,project_id", "initiative-to-project", "list", "--limit", "1"},
		{"--json", "initiative-to-project", "get", "initiative-to-project-id"},
		{"--json", "--fields", "id,health,initiative_id", "initiative-update", "list", "--limit", "1"},
		{"--json", "initiative-update", "get", "initiative-update-id"},
		{"--json", "roadmap", "list", "--limit", "1"},
		{"--json", "roadmap", "get", "roadmap-id"},
		{"--json", "--fields", "id,roadmap_id,project_id", "roadmap-to-project", "list", "--limit", "1"},
		{"--json", "roadmap-to-project", "get", "roadmap-to-project-id"},
		{"--json", "custom-view", "list", "--limit", "1"},
		{"--json", "custom-view", "subscribers", "custom-view-id"},
		{"--json", "custom-view", "get", "custom-view-id"},
		{"--json", "--fields", "id,name,status", "custom-view", "initiatives", "custom-view-id", "--limit", "1"},
		{"--json", "--fields", "identifier,title,state", "custom-view", "issues", "custom-view-id", "--limit", "1"},
		{"--json", "custom-view", "organization-preferences", "custom-view-id"},
		{"--json", "custom-view", "organization-preferences", "values", "custom-view-id"},
		{"--json", "--fields", "id,name,status", "custom-view", "projects", "custom-view-id", "--limit", "1"},
		{"--json", "custom-view", "user-preferences", "custom-view-id"},
		{"--json", "custom-view", "user-preferences", "values", "custom-view-id"},
		{"--json", "custom-view", "preference-values", "custom-view-id"},
		{"--json", "customer", "list", "--limit", "1"},
		{"--json", "customer", "get", "customer-id"},
		{"--json", "customer-need", "list", "--limit", "1"},
		{"--json", "customer-need", "get", "customer-need-id"},
		{"--json", "customer-status", "list", "--limit", "1"},
		{"--json", "customer-status", "get", "customer-status-id"},
		{"--json", "customer-tier", "list", "--limit", "1"},
		{"--json", "customer-tier", "get", "customer-tier-id"},
		{"--json", "favorite", "list", "--limit", "1"},
		{"--json", "favorite", "children", "favorite-folder-id", "--limit", "1"},
		{"--json", "favorite", "get", "favorite-id"},
		{"--json", "emoji", "list", "--limit", "1"},
		{"--json", "emoji", "get", "emoji-id"},
		{"--json", "attachment", "list", "--limit", "1"},
		{"--json", "attachment", "url", "https://github.com/kyanite/linctl/pull/1", "--limit", "1"},
		{"--json", "attachment", "get", "attachment-id"},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), "{")
		})
	}
}

func Test_CommandFlows_print_compact_json_when_compact_flag_is_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "--compact", "issue", "get", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `{"id":"issue-id"`)
	require.NotContains(t, output.String(), "\n  ")
}

func Test_CommandFlows_project_json_fields_when_fields_flag_is_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "--fields", "identifier,title,state", "issue", "get", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"identifier": "LIT-1"`)
	require.Contains(t, output.String(), `"title": "Detail issue"`)
	require.Contains(t, output.String(), `"state": "Todo"`)
	require.NotContains(t, output.String(), `"url"`)
	require.NotContains(t, output.String(), `"project_id"`)
}

func Test_CommandFlows_print_only_id_when_id_only_flag_is_set(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		output string
	}{
		{name: "issue get", args: []string{"--id-only", "issue", "get", "LIT-1"}, output: "issue-id\n"},
		{name: "issue history", args: []string{"--id-only", "issue", "history", "LIT-1"}, output: "issue-history-id\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, test.output, output.String())
		})
	}
}

func Test_CommandFlows_suppress_success_output_when_quiet_flag_is_set(t *testing.T) {
	tests := [][]string{
		{"--quiet", "doctor"},
		{"--quiet", "issue", "get", "LIT-1"},
		{"--quiet", "issue", "history", "LIT-1"},
	}

	for _, args := range tests {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Empty(t, output.String())
		})
	}
}

func Test_CommandFlows_fail_on_empty_list_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyIssueList: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "issue", "list"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_fail_on_empty_issue_child_list_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyIssueChildren: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "issue", "children", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_issue_child_list_reports_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"issue", "children", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime failed")
}

func Test_CommandFlows_fail_on_empty_project_updates_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyProjectUpdates: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "project", "updates", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_allow_empty_project_updates_without_fail_on_empty(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyProjectUpdates: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"project", "updates", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Empty(t, output.String())
}

func Test_CommandFlows_report_project_updates_sort_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--sort", "missing", "project", "updates", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), `sort field "missing" is not present`)
}

func Test_CommandFlows_project_comment_children_omit_body_from_json(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "project comments", args: []string{"project", "comments", "project-id", "--json"}},
		{name: "project update comments", args: []string{"project-update", "comments", "project-update-id", "--json"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), `"comments"`)
			require.NotContains(t, output.String(), `"body"`)
		})
	}
}

func Test_CommandFlows_project_child_reads_cover_json_and_sort_branches(t *testing.T) {
	t.Run("project milestone issues json", func(t *testing.T) {
		output := bytes.Buffer{}
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(&output)
		command.SetArgs([]string{"project-milestone", "issues", "project-milestone-id", "--json"})

		err := command.ExecuteContext(context.Background())

		require.NoError(t, err)
		require.Contains(t, output.String(), `"project_milestone_id"`)
		require.Contains(t, output.String(), `"issues"`)
	})

	t.Run("project comments sort errors", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--sort", "missing", "project", "comments", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), `sort field "missing" is not present`)
	})

	t.Run("project comments text output", func(t *testing.T) {
		output := bytes.Buffer{}
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(&output)
		command.SetArgs([]string{"project", "comments", "project-id"})

		err := command.ExecuteContext(context.Background())

		require.NoError(t, err)
		require.Contains(t, output.String(), "comment-id Omer 2026-06-19T12:00:00Z")
	})

	t.Run("release search json", func(t *testing.T) {
		output := bytes.Buffer{}
		restore := useCommandRuntime(t, commandFlowFakeClient{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetOut(&output)
		command.SetArgs([]string{"release", "search", "mobile", "--json"})

		err := command.ExecuteContext(context.Background())

		require.NoError(t, err)
		require.Contains(t, output.String(), `"releases"`)
	})

	t.Run("release search fail on empty", func(t *testing.T) {
		restore := useCommandRuntime(t, commandFlowFakeClient{emptyReleaseSearch: true})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"--fail-on-empty", "release", "search", "mobile"})

		err := command.ExecuteContext(context.Background())

		require.Error(t, err)
		require.Contains(t, err.Error(), "empty result")
	})
}

func Test_CommandFlows_label_child_reads_cover_json_pages(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name:     "label children json",
			args:     []string{"label", "children", "label-id", "--json"},
			contains: `"label_name": "Bug"`,
		},
		{
			name:     "label issues json",
			args:     []string{"label", "issues", "label-id", "--json"},
			contains: `"identifier": "LIT-1"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), test.contains)
		})
	}
}

func Test_CommandFlows_fail_on_empty_project_milestones_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyProjectMilestones: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "project-milestone", "list", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_allow_empty_project_milestones_without_fail_on_empty(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{emptyProjectMilestones: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"project-milestone", "list", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Empty(t, output.String())
}

func Test_CommandFlows_report_project_milestone_sort_errors(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--sort", "missing", "project-milestone", "list", "project-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), `sort field "missing" is not present`)
}

func Test_CommandFlows_fail_on_empty_sla_configurations_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptySLAConfigurations: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "sla-configuration", "list", "team-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_fail_on_empty_semantic_search_when_fail_on_empty_flag_is_set(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{emptySemanticSearch: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"--fail-on-empty", "semantic-search", "agent search"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "empty result")
}

func Test_CommandFlows_fail_on_empty_typed_search_when_fail_on_empty_flag_is_set(t *testing.T) {
	tests := []struct {
		name string
		args []string
		fake commandFlowFakeClient
	}{
		{
			name: "documents",
			args: []string{"--fail-on-empty", "search", "documents", "agent search"},
			fake: commandFlowFakeClient{emptySearchDocuments: true},
		},
		{
			name: "issues",
			args: []string{"--fail-on-empty", "search", "issues", "agent search"},
			fake: commandFlowFakeClient{emptySearchIssues: true},
		},
		{
			name: "projects",
			args: []string{"--fail-on-empty", "search", "projects", "agent search"},
			fake: commandFlowFakeClient{emptySearchProjects: true},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "empty result")
		})
	}
}

func Test_CommandFlows_semantic_search_honors_id_only_and_quiet(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		output string
	}{
		{name: "id only", args: []string{"--id-only", "semantic-search", "agent search"}, output: "issue-id\n"},
		{name: "quiet", args: []string{"--quiet", "semantic-search", "agent search"}, output: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, test.output, output.String())
		})
	}
}

func Test_CommandFlows_typed_search_writers_emit_json(t *testing.T) {
	tests := []struct {
		name  string
		write func(*cobra.Command, *rootOptions) error
		want  string
	}{
		{
			name: "document",
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSearchDocument(command, options, client.SearchDocumentSummary{
					ID:    "search-document-id",
					Title: "Search spec",
				})
			},
			want: `{"id":"search-document-id","title":"Search spec","slug_id":"","url":""}`,
		},
		{
			name: "issue",
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSearchIssue(command, options, client.SearchIssueSummary{
					ID:         "search-issue-id",
					Identifier: "LIT-30",
					Title:      "Search issue",
				})
			},
			want: `{"id":"search-issue-id","identifier":"LIT-30","title":"Search issue","url":"","team_id":"","team_key":"","team_name":"","state_id":"","state_name":"","state_type":""}`,
		},
		{
			name: "project",
			write: func(command *cobra.Command, options *rootOptions) error {
				return writeSearchProject(command, options, client.SearchProjectSummary{
					ID:   "search-project-id",
					Name: "Search project",
				})
			},
			want: `{"id":"search-project-id","name":"Search project","slug_id":"","url":"","status":{"id":"","name":"","type":""},"teams":null}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			command := &cobra.Command{}
			command.SetOut(&output)

			err := test.write(command, &rootOptions{json: true})

			require.NoError(t, err)
			require.JSONEq(t, test.want, output.String())
		})
	}
}

func Test_CommandFlows_typed_search_honors_id_only_and_quiet(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		output string
	}{
		{name: "documents id only", args: []string{"--id-only", "search", "documents", "agent search"}, output: "search-document-id\n"},
		{name: "documents quiet", args: []string{"--quiet", "search", "documents", "agent search"}, output: ""},
		{name: "issues id only", args: []string{"--id-only", "search", "issues", "agent search"}, output: "search-issue-id\n"},
		{name: "issues quiet", args: []string{"--quiet", "search", "issues", "agent search"}, output: ""},
		{name: "projects id only", args: []string{"--id-only", "search", "projects", "agent search"}, output: "search-project-id\n"},
		{name: "projects quiet", args: []string{"--quiet", "search", "projects", "agent search"}, output: ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, test.output, output.String())
		})
	}
}

func Test_CommandFlows_user_drafts_honor_list_controls(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		fake   commandFlowFakeClient
		output string
	}{
		{name: "id only", args: []string{"--id-only", "user", "drafts"}, output: "draft-id\n"},
		{name: "quiet", args: []string{"--quiet", "user", "drafts"}, output: ""},
		{
			name:   "sort",
			args:   []string{"--sort", "parent_key", "--order", "desc", "user", "drafts"},
			output: "draft-id issue LIT-3 Draft issue\n",
		},
		{
			name: "empty",
			args: []string{"--fail-on-empty", "user", "drafts"},
			fake: commandFlowFakeClient{emptyViewerDrafts: true},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			if test.name == "empty" {
				require.Error(t, err)
				require.Contains(t, err.Error(), "empty result")
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.output, output.String())
		})
	}
}

func Test_CommandFlows_user_drafts_json_uses_projected_page(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "--sort", "parent_key", "user", "drafts"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"drafts"`)
	require.Contains(t, output.String(), `"parent_key": "LIT-3"`)
}

func Test_CommandFlows_get_project_milestone(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"project-milestone", "get", "project-milestone-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "project-milestone-id Launch milestone [next]")
}

func Test_CommandFlows_get_project_milestone_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "project-milestone", "get", "project-milestone-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"id": "project-milestone-id"`)
	require.Contains(t, output.String(), `"status": "next"`)
}

func Test_CommandFlows_report_project_milestone_get_runtime_error(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"project-milestone", "get", "project-milestone-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime failed")
}

func Test_CommandFlows_report_project_milestone_write_runtime_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "create", args: []string{"project-milestone", "create", "project-id", "--name", "Created milestone"}},
		{name: "update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := buildCommandRuntime
			buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
				return commandRuntime{}, errors.New("runtime failed")
			}
			defer func() {
				buildCommandRuntime = original
			}()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "runtime failed")
		})
	}
}

func Test_CommandFlows_report_project_milestone_get_writer_error(t *testing.T) {
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"project-milestone", "get", "project-milestone-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func Test_CommandFlows_report_project_milestone_write_writer_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "create", args: []string{"project-milestone", "create", "project-id", "--name", "Created milestone"}},
		{name: "update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(commandFailingWriter{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "write failed")
		})
	}
}

func Test_CommandFlows_sort_issue_list_when_sort_flags_are_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{multiIssueList: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--sort", "title", "--order", "desc", "issue", "list"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Less(t, strings.Index(output.String(), "Zebra issue"), strings.Index(output.String(), "Alpha issue"))
}

func Test_CommandFlows_print_minimal_human_output_when_format_flag_is_set(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--format", "minimal", "issue", "get", "LIT-1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Equal(t, "LIT-1\n", output.String())
}

func Test_CommandFlows_print_workflow_state_list_as_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, commandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "workflow-state", "list", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"workflow_states": [`)
	require.Contains(t, output.String(), `"team_key": "LIT"`)
}

func Test_CommandFlows_report_operation_errors(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		operation string
		contains  string
	}{
		{name: "target resolve", args: []string{"target"}, operation: "Teams", contains: "resolve teams"},
		{name: "doctor target resolve", args: []string{"doctor"}, operation: "Teams", contains: "resolve teams"},
		{name: "whoami resolve", args: []string{"whoami"}, operation: "Viewer", contains: "resolve viewer"},
		{name: "application info", args: []string{"application", "info", "app-client-id"}, operation: "applicationInfo", contains: "get application info app-client-id"},
		{name: "agent activity list", args: []string{"agent-activity", "list"}, operation: "agentActivities", contains: "list agent activities"},
		{name: "agent activity get", args: []string{"agent-activity", "get", "agent-activity-id"}, operation: "agentActivity", contains: "get agent activity agent-activity-id"},
		{name: "agent skill list", args: []string{"agent-skill", "list"}, operation: "agentSkills", contains: "list agent skills"},
		{name: "agent skill get", args: []string{"agent-skill", "get", "agent-skill-id"}, operation: "agentSkill", contains: "get agent skill agent-skill-id"},
		{name: "external user list", args: []string{"external-user", "list"}, operation: "externalUsers", contains: "list external users"},
		{name: "external user get", args: []string{"external-user", "get", "external-user-id"}, operation: "externalUser", contains: "get external user external-user-id"},
		{name: "audit entry types", args: []string{"audit-entry", "types"}, operation: "auditEntryTypes", contains: "list audit entry types"},
		{name: "organization exists", args: []string{"organization", "exists", "kyanite"}, operation: "organizationExists", contains: "operation failed"},
		{name: "organization labels", args: []string{"organization", "labels"}, operation: "organization_labels", contains: "list organization labels"},
		{name: "organization project labels", args: []string{"organization", "project-labels"}, operation: "organization_projectLabels", contains: "list organization project labels"},
		{name: "organization teams", args: []string{"organization", "teams"}, operation: "organization_teams", contains: "list organization teams"},
		{name: "organization templates", args: []string{"organization", "templates"}, operation: "organization_templates", contains: "list organization templates"},
		{name: "organization users", args: []string{"organization", "users"}, operation: "organization_users", contains: "list organization users"},
		{name: "rate limit status", args: []string{"rate-limit", "status"}, operation: "rateLimitStatus", contains: "operation failed"},
		{name: "notification list", args: []string{"notification", "list"}, operation: "notifications", contains: "list notifications"},
		{name: "notification get", args: []string{"notification", "get", "notification-id"}, operation: "notification", contains: "get notification notification-id"},
		{name: "notification subscription list", args: []string{"notification", "subscription", "list"}, operation: "notificationSubscriptions", contains: "list notification subscriptions"},
		{name: "notification subscription get", args: []string{"notification", "subscription", "get", "notification-subscription-id"}, operation: "notificationSubscription", contains: "get notification subscription notification-subscription-id"},
		{name: "triage responsibility list", args: []string{"triage-responsibility", "list"}, operation: "triageResponsibilities", contains: "list triage responsibilities"},
		{name: "triage responsibility get", args: []string{"triage-responsibility", "get", "triage-responsibility-id"}, operation: "triageResponsibility", contains: "get triage responsibility triage-responsibility-id"},
		{name: "triage responsibility manual selection", args: []string{"triage-responsibility", "manual-selection", "triage-responsibility-id"}, operation: "triageResponsibility_manualSelection", contains: "get triage responsibility manual selection triage-responsibility-id"},
		{name: "SLA configuration list", args: []string{"sla-configuration", "list", "team-id"}, operation: "slaConfigurations", contains: "list SLA configurations team-id"},
		{name: "semantic search", args: []string{"semantic-search", "agent search"}, operation: "semanticSearch", contains: "semantic search"},
		{name: "search documents", args: []string{"search", "documents", "agent search"}, operation: "searchDocuments", contains: "search documents"},
		{name: "search issues", args: []string{"search", "issues", "agent search"}, operation: "searchIssues", contains: "search issues"},
		{name: "search projects", args: []string{"search", "projects", "agent search"}, operation: "searchProjects", contains: "search projects"},
		{name: "release pipeline list", args: []string{"release-pipeline", "list"}, operation: "releasePipelines", contains: "list release pipelines"},
		{name: "release pipeline get", args: []string{"release-pipeline", "get", "release-pipeline-id"}, operation: "releasePipeline", contains: "get release pipeline release-pipeline-id"},
		{name: "release pipeline releases", args: []string{"release-pipeline", "releases", "release-pipeline-id"}, operation: "releasePipeline_releases", contains: "list release pipeline releases release-pipeline-id"},
		{name: "release pipeline stages", args: []string{"release-pipeline", "stages", "release-pipeline-id"}, operation: "releasePipeline_stages", contains: "list release pipeline stages release-pipeline-id"},
		{name: "release pipeline teams", args: []string{"release-pipeline", "teams", "release-pipeline-id"}, operation: "releasePipeline_teams", contains: "list release pipeline teams release-pipeline-id"},
		{name: "release stage list", args: []string{"release-stage", "list"}, operation: "releaseStages", contains: "list release stages"},
		{name: "release stage get", args: []string{"release-stage", "get", "release-stage-id"}, operation: "releaseStage", contains: "get release stage release-stage-id"},
		{name: "release stage releases", args: []string{"release-stage", "releases", "release-stage-id"}, operation: "releaseStage_releases", contains: "list release stage releases release-stage-id"},
		{name: "release list", args: []string{"release", "list"}, operation: "releases", contains: "list releases"},
		{name: "release search", args: []string{"release", "search", "mobile"}, operation: "releaseSearch", contains: "search releases"},
		{name: "release get", args: []string{"release", "get", "release-id"}, operation: "release", contains: "get release release-id"},
		{name: "release history", args: []string{"release", "history", "release-id"}, operation: "release_history", contains: "list release history release-id"},
		{name: "release documents", args: []string{"release", "documents", "release-id"}, operation: "release_documents", contains: "list release documents release-id"},
		{name: "release issues", args: []string{"release", "issues", "release-id"}, operation: "release_issues", contains: "list release issues release-id"},
		{name: "release links", args: []string{"release", "links", "release-id"}, operation: "release_links", contains: "list release links release-id"},
		{name: "external link get", args: []string{"external-link", "get", "release-link-id"}, operation: "entityExternalLink", contains: "get external link release-link-id"},
		{name: "release note list", args: []string{"release-note", "list"}, operation: "releaseNotes", contains: "list release notes"},
		{name: "release note get", args: []string{"release-note", "get", "release-note-id"}, operation: "releaseNote", contains: "get release note release-note-id"},
		{name: "issue to release list", args: []string{"issue-to-release", "list"}, operation: "issueToReleases", contains: "list issue to releases"},
		{name: "issue to release get", args: []string{"issue-to-release", "get", "issue-to-release-id"}, operation: "issueToRelease", contains: "get issue to release issue-to-release-id"},
		{name: "next target resolve", args: []string{"next", "--dry-run"}, operation: "Teams", contains: "resolve teams"},
		{name: "next issues", args: []string{"next", "--dry-run"}, operation: "NextIssuesByTeam", contains: "list next issues"},
		{name: "issue list target resolve", args: []string{"issue", "list"}, operation: "Teams", contains: "resolve teams"},
		{name: "issue list", args: []string{"issue", "list"}, operation: "IssuesByTeam", contains: "list issues"},
		{name: "issue list project filter", args: []string{"issue", "list", "--project", "project-id"}, operation: "IssuesByTeamProject", contains: "list issues"},
		{name: "issue list mine filter", args: []string{"issue", "list", "--mine"}, operation: "IssuesByTeamAssignee", contains: "list issues"},
		{name: "issue list assignee filter", args: []string{"issue", "list", "--assignee", "assignee-id"}, operation: "IssuesByTeamAssignee", contains: "list issues"},
		{name: "issue list label filter", args: []string{"issue", "list", "--label", "label-id"}, operation: "IssuesByTeamLabel", contains: "list issues"},
		{name: "issue list cycle filter", args: []string{"issue", "list", "--cycle", "cycle-id"}, operation: "IssuesByTeamCycle", contains: "list issues"},
		{name: "issue list created-after filter", args: []string{"issue", "list", "--created-after", "2026-06-01"}, operation: "IssuesByTeamCreatedAfter", contains: "list issues"},
		{name: "issue list created-since filter", args: []string{"issue", "list", "--created-since", "2026-06-01"}, operation: "IssuesByTeamCreatedAfter", contains: "list issues"},
		{name: "issue list created-before filter", args: []string{"issue", "list", "--created-before", "2026-06-30"}, operation: "IssuesByTeamCreatedBefore", contains: "list issues"},
		{name: "issue list has blockers filter", args: []string{"issue", "list", "--has-blockers"}, operation: "IssuesByTeamHasBlockers", contains: "list issues"},
		{name: "issue list blocks filter", args: []string{"issue", "list", "--blocks"}, operation: "IssuesByTeamBlocks", contains: "list issues"},
		{name: "issue list blocked by filter", args: []string{"issue", "list", "--blocked-by", "LIT-1"}, operation: "IssueBlockedIssues", contains: "list issues"},
		{name: "issue list all teams", args: []string{"issue", "list", "--all-teams"}, operation: "issues", contains: "list issues"},
		{name: "issue search target resolve", args: []string{"issue", "search", "needle"}, operation: "Teams", contains: "resolve teams"},
		{name: "issue search", args: []string{"issue", "search", "needle"}, operation: "issueSearch", contains: "search issues"},
		{name: "issue get", args: []string{"issue", "get", "LIT-1"}, operation: "issue", contains: "get issue LIT-1"},
		{name: "issue deps", args: []string{"issue", "deps", "LIT-1"}, operation: "IssueDependencies", contains: "get issue dependencies LIT-1"},
		{name: "issue attachments", args: []string{"issue", "attachments", "LIT-1"}, operation: "issue_attachments", contains: "list issue attachments LIT-1"},
		{name: "issue children", args: []string{"issue", "children", "LIT-1"}, operation: "issue_children", contains: "list issue children LIT-1"},
		{name: "issue documents", args: []string{"issue", "documents", "LIT-1"}, operation: "issue_documents", contains: "list issue documents LIT-1"},
		{name: "issue former attachments", args: []string{"issue", "former-attachments", "LIT-1"}, operation: "issue_formerAttachments", contains: "list issue former attachments LIT-1"},
		{name: "issue history", args: []string{"issue", "history", "LIT-1"}, operation: "issue_history", contains: "list issue history LIT-1"},
		{name: "issue inverse relations", args: []string{"issue", "inverse-relations", "LIT-1"}, operation: "issue_inverseRelations", contains: "list issue inverse relations LIT-1"},
		{name: "issue labels", args: []string{"issue", "labels", "LIT-1"}, operation: "issue_labels", contains: "list issue labels LIT-1"},
		{name: "issue relations", args: []string{"issue", "relations", "LIT-1"}, operation: "issue_relations", contains: "list issue relations LIT-1"},
		{name: "issue releases", args: []string{"issue", "releases", "LIT-1"}, operation: "issue_releases", contains: "list issue releases LIT-1"},
		{name: "issue relation list", args: []string{"issue-relation", "list"}, operation: "issueRelations", contains: "list issue relations"},
		{name: "issue relation get", args: []string{"issue-relation", "get", "issue-relation-id"}, operation: "issueRelation", contains: "get issue relation issue-relation-id"},
		{name: "issue pr", args: []string{"issue", "pr", "LIT-1"}, operation: "issue", contains: "get issue LIT-1"},
		{name: "issue create", args: []string{"issue", "create", "--title", "Created issue"}, operation: "IssueCreate", contains: "create issue"},
		{name: "issue update", args: []string{"issue", "update", "LIT-1", "--title", "Updated issue"}, operation: "IssueUpdate", contains: "update issue LIT-1"},
		{name: "issue start state", args: []string{"issue", "start", "LIT-1"}, operation: "StartedWorkflowStates", contains: "list started workflow states"},
		{name: "issue start update", args: []string{"issue", "start", "LIT-1"}, operation: "IssueUpdate", contains: "start issue LIT-1"},
		{name: "issue comment", args: []string{"issue", "comment", "LIT-1", "--body", "Looks good"}, operation: "IssueCommentCreate", contains: "comment on issue LIT-1"},
		{name: "issue reply", args: []string{"issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"}, operation: "IssueCommentCreate", contains: "comment on issue LIT-1"},
		{name: "comment list", args: []string{"comment", "list"}, operation: "comments", contains: "list comments"},
		{name: "comment get", args: []string{"comment", "get", "comment-id"}, operation: "comment", contains: "get comment comment-id"},
		{name: "issue close", args: []string{"issue", "close", "LIT-1"}, operation: "IssueClose", contains: "close issue LIT-1"},
		{name: "project list target resolve", args: []string{"project", "list"}, operation: "Teams", contains: "resolve teams"},
		{name: "project list", args: []string{"project", "list"}, operation: "Projects", contains: "list projects"},
		{name: "project get", args: []string{"project", "get", "project-id"}, operation: "project", contains: "get project project-id"},
		{name: "project attachments", args: []string{"project", "attachments", "project-id"}, operation: "project_attachments", contains: "list project attachments project-id"},
		{name: "project documents", args: []string{"project", "documents", "project-id"}, operation: "project_documents", contains: "list project documents project-id"},
		{name: "project external links", args: []string{"project", "external-links", "project-id"}, operation: "project_externalLinks", contains: "list project external links project-id"},
		{name: "project history", args: []string{"project", "history", "project-id"}, operation: "project_history", contains: "list project history project-id"},
		{name: "project initiative links", args: []string{"project", "initiative-links", "project-id"}, operation: "project_initiativeToProjects", contains: "list project initiative associations project-id"},
		{name: "project initiatives", args: []string{"project", "initiatives", "project-id"}, operation: "project_initiatives", contains: "list project initiatives project-id"},
		{name: "project inverse relations", args: []string{"project", "inverse-relations", "project-id"}, operation: "project_inverseRelations", contains: "list project inverse relations project-id"},
		{name: "project issues", args: []string{"project", "issues", "project-id"}, operation: "project_issues", contains: "list project issues project-id"},
		{name: "project comments", args: []string{"project", "comments", "project-id"}, operation: "project_comments", contains: "list project comments project-id"},
		{name: "project labels", args: []string{"project", "labels", "project-id"}, operation: "project_labels", contains: "list project labels project-id"},
		{name: "project members", args: []string{"project", "members", "project-id"}, operation: "project_members", contains: "list project members project-id"},
		{name: "project needs", args: []string{"project", "needs", "project-id"}, operation: "project_needs", contains: "list project customer needs project-id"},
		{name: "project relations", args: []string{"project", "relations", "project-id"}, operation: "project_relations", contains: "list project relations project-id"},
		{name: "project teams", args: []string{"project", "teams", "project-id"}, operation: "project_teams", contains: "list project teams project-id"},
		{name: "project updates", args: []string{"project", "updates", "project-id"}, operation: "ProjectUpdates", contains: "list project updates project-id"},
		{name: "project update list", args: []string{"project-update", "list"}, operation: "projectUpdates", contains: "list project updates"},
		{name: "project update get", args: []string{"project-update", "get", "project-update-id"}, operation: "projectUpdate", contains: "get project update project-update-id"},
		{name: "project update comments", args: []string{"project-update", "comments", "project-update-id"}, operation: "projectUpdate_comments", contains: "list project update comments project-update-id"},
		{name: "project milestone list", args: []string{"project-milestone", "list", "project-id"}, operation: "project_projectMilestones", contains: "list project milestones project-id"},
		{name: "project milestone get", args: []string{"project-milestone", "get", "project-milestone-id"}, operation: "projectMilestone", contains: "get project milestone project-milestone-id"},
		{name: "project milestone issues", args: []string{"project-milestone", "issues", "project-milestone-id"}, operation: "projectMilestone_issues", contains: "list project milestone issues project-milestone-id"},
		{name: "project milestone create", args: []string{"project-milestone", "create", "project-id", "--name", "Created milestone"}, operation: "ProjectMilestoneCreate", contains: "create project milestone"},
		{name: "project milestone update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated milestone"}, operation: "ProjectMilestoneUpdate", contains: "update project milestone project-milestone-id"},
		{name: "project create", args: []string{"project", "create", "--name", "Created project"}, operation: "ProjectCreate", contains: "create project"},
		{name: "project update", args: []string{"project", "update", "project-id", "--name", "Updated project"}, operation: "ProjectUpdate", contains: "update project project-id"},
		{name: "project archive", args: []string{"project", "archive", "project-id"}, operation: "ProjectArchive", contains: "archive project project-id"},
		{name: "document list", args: []string{"document", "list"}, operation: "Documents", contains: "list documents"},
		{name: "document get", args: []string{"document", "get", "document-id"}, operation: "document", contains: "get document document-id"},
		{name: "label list", args: []string{"label", "list"}, operation: "IssueLabels", contains: "list labels"},
		{name: "label get", args: []string{"label", "get", "label-id"}, operation: "issueLabel", contains: "get label label-id"},
		{name: "label children", args: []string{"label", "children", "label-id"}, operation: "issueLabel_children", contains: "list label children label-id"},
		{name: "label issues", args: []string{"label", "issues", "label-id"}, operation: "issueLabel_issues", contains: "list label issues label-id"},
		{name: "team list", args: []string{"team", "list"}, operation: "Teams", contains: "list teams"},
		{name: "team get", args: []string{"team", "get", "team-id"}, operation: "team", contains: "get team team-id"},
		{name: "team cycles", args: []string{"team", "cycles", "team-id"}, operation: "team_cycles", contains: "list team cycles team-id"},
		{name: "team issues", args: []string{"team", "issues", "team-id"}, operation: "team_issues", contains: "list team issues team-id"},
		{name: "team labels", args: []string{"team", "labels", "team-id"}, operation: "team_labels", contains: "list team labels team-id"},
		{name: "team members", args: []string{"team", "members", "team-id"}, operation: "team_members", contains: "list team members team-id"},
		{name: "team memberships", args: []string{"team", "memberships", "team-id"}, operation: "team_memberships", contains: "list team memberships team-id"},
		{name: "team projects", args: []string{"team", "projects", "team-id"}, operation: "team_projects", contains: "list team projects team-id"},
		{name: "team release pipelines", args: []string{"team", "release-pipelines", "team-id"}, operation: "team_releasePipelines", contains: "list team release pipelines team-id"},
		{name: "team states", args: []string{"team", "states", "team-id"}, operation: "team_states", contains: "list team states team-id"},
		{name: "team templates", args: []string{"team", "templates", "team-id"}, operation: "team_templates", contains: "list team templates team-id"},
		{name: "user list", args: []string{"user", "list"}, operation: "users", contains: "list users"},
		{name: "user get", args: []string{"user", "get", "user-id"}, operation: "user", contains: "get user user-id"},
		{name: "user me", args: []string{"user", "me"}, operation: "viewer", contains: "get viewer user"},
		{name: "user drafts", args: []string{"user", "drafts"}, operation: "viewer_drafts", contains: "list viewer drafts"},
		{name: "user assigned issues", args: []string{"user", "assigned-issues", "user-id"}, operation: "user_assignedIssues", contains: "list user assigned issues user-id"},
		{name: "user created issues", args: []string{"user", "created-issues", "user-id"}, operation: "user_createdIssues", contains: "list user created issues user-id"},
		{name: "user delegated issues", args: []string{"user", "delegated-issues", "user-id"}, operation: "user_delegatedIssues", contains: "list user delegated issues user-id"},
		{name: "user team memberships", args: []string{"user", "team-memberships", "user-id"}, operation: "user_teamMemberships", contains: "list user team memberships user-id"},
		{name: "user teams", args: []string{"user", "teams", "user-id"}, operation: "user_teams", contains: "list user teams user-id"},
		{name: "user my assigned issues", args: []string{"user", "my-assigned-issues"}, operation: "viewer_assignedIssues", contains: "list viewer assigned issues"},
		{name: "user my created issues", args: []string{"user", "my-created-issues"}, operation: "viewer_createdIssues", contains: "list viewer created issues"},
		{name: "user my delegated issues", args: []string{"user", "my-delegated-issues"}, operation: "viewer_delegatedIssues", contains: "list viewer delegated issues"},
		{name: "user my team memberships", args: []string{"user", "my-team-memberships"}, operation: "viewer_teamMemberships", contains: "list viewer team memberships"},
		{name: "user my teams", args: []string{"user", "my-teams"}, operation: "viewer_teams", contains: "list viewer teams"},
		{name: "workflow state list", args: []string{"workflow-state", "list"}, operation: "workflowStates", contains: "list workflow states"},
		{name: "workflow state get", args: []string{"workflow-state", "get", "workflow-state-id"}, operation: "workflowState", contains: "get workflow state workflow-state-id"},
		{name: "time schedule list", args: []string{"time-schedule", "list"}, operation: "timeSchedules", contains: "list time schedules"},
		{name: "time schedule get", args: []string{"time-schedule", "get", "time-schedule-id"}, operation: "timeSchedule", contains: "get time schedule time-schedule-id"},
		{name: "template list", args: []string{"template", "list"}, operation: "templates", contains: "list templates"},
		{name: "template get", args: []string{"template", "get", "template-id"}, operation: "template", contains: "get template template-id"},
		{name: "initiative list", args: []string{"initiative", "list"}, operation: "initiatives", contains: "list initiatives"},
		{name: "initiative get", args: []string{"initiative", "get", "initiative-id"}, operation: "initiative", contains: "get initiative initiative-id"},
		{name: "initiative history", args: []string{"initiative", "history", "initiative-id"}, operation: "initiative_history", contains: "list initiative history initiative-id"},
		{name: "initiative links", args: []string{"initiative", "links", "initiative-id"}, operation: "initiative_links", contains: "list initiative links initiative-id"},
		{name: "initiative sub-initiatives", args: []string{"initiative", "sub-initiatives", "initiative-id"}, operation: "initiative_subInitiatives", contains: "list initiative sub-initiatives initiative-id"},
		{name: "initiative updates", args: []string{"initiative", "updates", "initiative-id"}, operation: "initiative_initiativeUpdates", contains: "list initiative updates initiative-id"},
		{name: "initiative documents", args: []string{"initiative", "documents", "initiative-id"}, operation: "initiative_documents", contains: "list initiative documents initiative-id"},
		{name: "initiative projects", args: []string{"initiative", "projects", "initiative-id"}, operation: "initiative_projects", contains: "list initiative projects initiative-id"},
		{name: "initiative relation list", args: []string{"initiative-relation", "list"}, operation: "initiativeRelations", contains: "list initiative relations"},
		{name: "initiative relation get", args: []string{"initiative-relation", "get", "initiative-relation-id"}, operation: "initiativeRelation", contains: "get initiative relation initiative-relation-id"},
		{name: "initiative to project list", args: []string{"initiative-to-project", "list"}, operation: "initiativeToProjects", contains: "list initiative to projects"},
		{name: "initiative to project get", args: []string{"initiative-to-project", "get", "initiative-to-project-id"}, operation: "initiativeToProject", contains: "get initiative to project initiative-to-project-id"},
		{name: "initiative update list", args: []string{"initiative-update", "list"}, operation: "initiativeUpdates", contains: "list initiative updates"},
		{name: "initiative update get", args: []string{"initiative-update", "get", "initiative-update-id"}, operation: "initiativeUpdate", contains: "get initiative update initiative-update-id"},
		{name: "roadmap list", args: []string{"roadmap", "list"}, operation: "roadmaps", contains: "list roadmaps"},
		{name: "roadmap get", args: []string{"roadmap", "get", "roadmap-id"}, operation: "roadmap", contains: "get roadmap roadmap-id"},
		{name: "roadmap to project list", args: []string{"roadmap-to-project", "list"}, operation: "roadmapToProjects", contains: "list roadmap to projects"},
		{name: "roadmap to project get", args: []string{"roadmap-to-project", "get", "roadmap-to-project-id"}, operation: "roadmapToProject", contains: "get roadmap to project roadmap-to-project-id"},
		{name: "custom view list", args: []string{"custom-view", "list"}, operation: "customViews", contains: "list custom views"},
		{name: "custom view subscribers", args: []string{"custom-view", "subscribers", "custom-view-id"}, operation: "customViewHasSubscribers", contains: "get custom view subscribers custom-view-id"},
		{name: "custom view get", args: []string{"custom-view", "get", "custom-view-id"}, operation: "customView", contains: "get custom view custom-view-id"},
		{name: "custom view initiatives", args: []string{"custom-view", "initiatives", "custom-view-id"}, operation: "customView_initiatives", contains: "list custom view initiatives custom-view-id"},
		{name: "custom view issues", args: []string{"custom-view", "issues", "custom-view-id"}, operation: "customView_issues", contains: "list custom view issues custom-view-id"},
		{name: "custom view organization preferences", args: []string{"custom-view", "organization-preferences", "custom-view-id"}, operation: "customView_organizationViewPreferences", contains: "get custom view organization preferences custom-view-id"},
		{name: "custom view organization preference values", args: []string{"custom-view", "organization-preferences", "values", "custom-view-id"}, operation: "customView_organizationViewPreferences_preferences", contains: "get custom view organization preference values custom-view-id"},
		{name: "custom view projects", args: []string{"custom-view", "projects", "custom-view-id"}, operation: "customView_projects", contains: "list custom view projects custom-view-id"},
		{name: "custom view user preferences", args: []string{"custom-view", "user-preferences", "custom-view-id"}, operation: "customView_userViewPreferences", contains: "get custom view user preferences custom-view-id"},
		{name: "custom view user preference values", args: []string{"custom-view", "user-preferences", "values", "custom-view-id"}, operation: "customView_userViewPreferences_preferences", contains: "get custom view user preference values custom-view-id"},
		{name: "custom view preference values", args: []string{"custom-view", "preference-values", "custom-view-id"}, operation: "customView_viewPreferencesValues", contains: "get custom view preference values custom-view-id"},
		{name: "customer list", args: []string{"customer", "list"}, operation: "customers", contains: "list customers"},
		{name: "customer get", args: []string{"customer", "get", "customer-id"}, operation: "customer", contains: "get customer customer-id"},
		{name: "customer need list", args: []string{"customer-need", "list"}, operation: "customerNeeds", contains: "list customer needs"},
		{name: "customer need get", args: []string{"customer-need", "get", "customer-need-id"}, operation: "customerNeed", contains: "get customer need customer-need-id"},
		{name: "customer status list", args: []string{"customer-status", "list"}, operation: "customerStatuses", contains: "list customer statuses"},
		{name: "customer status get", args: []string{"customer-status", "get", "customer-status-id"}, operation: "customerStatus", contains: "get customer status customer-status-id"},
		{name: "customer tier list", args: []string{"customer-tier", "list"}, operation: "customerTiers", contains: "list customer tiers"},
		{name: "customer tier get", args: []string{"customer-tier", "get", "customer-tier-id"}, operation: "customerTier", contains: "get customer tier customer-tier-id"},
		{name: "favorite list", args: []string{"favorite", "list"}, operation: "favorites", contains: "list favorites"},
		{name: "favorite children", args: []string{"favorite", "children", "favorite-folder-id"}, operation: "favorite_children", contains: "list favorite children favorite-folder-id"},
		{name: "favorite get", args: []string{"favorite", "get", "favorite-id"}, operation: "favorite", contains: "get favorite favorite-id"},
		{name: "emoji list", args: []string{"emoji", "list"}, operation: "emojis", contains: "list emojis"},
		{name: "emoji get", args: []string{"emoji", "get", "emoji-id"}, operation: "emoji", contains: "get emoji emoji-id"},
		{name: "attachment list", args: []string{"attachment", "list"}, operation: "attachments", contains: "list attachments"},
		{name: "attachment url", args: []string{"attachment", "url", "https://github.com/kyanite/linctl/pull/1"}, operation: "attachmentsForURL", contains: "list attachments for url https://github.com/kyanite/linctl/pull/1"},
		{name: "attachment get", args: []string{"attachment", "get", "attachment-id"}, operation: "attachment", contains: "get attachment attachment-id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: test.operation})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.contains)
		})
	}
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()

	command := exec.Command("git", args...)
	command.Dir = dir
	output, err := command.CombinedOutput()
	require.NoError(t, err, string(output))
}

func writeTempTextFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "body.md")
	err := os.WriteFile(path, []byte(content), 0o600)
	require.NoError(t, err)

	return path
}

type commandFailingWriter struct{}

func (writer commandFailingWriter) Write(_ []byte) (int, error) {
	return 0, errors.New("write failed")
}

func useCommandRuntime(t *testing.T, graphqlClient graphql.Client) func() {
	t.Helper()

	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return testCommandRuntime(graphqlClient), nil
	}

	return func() {
		buildCommandRuntime = original
	}
}

func testCommandRuntime(graphqlClient graphql.Client) commandRuntime {
	return commandRuntime{
		config: config.Resolved{
			Token: "test-token",
			Target: config.Target{
				OrgID:     "org-id",
				TeamKey:   "LIT",
				TeamID:    "team-id",
				ProjectID: "project-id",
			},
		},
		graphqlClient: graphqlClient,
	}
}

type commandFlowFakeClient struct {
	emptyIssueList              bool
	emptyIssueChildren          bool
	emptyIssueComments          bool
	emptyIssueProject           bool
	emptyIssueMine              bool
	emptyIssueLabel             bool
	emptyIssueCycle             bool
	emptyIssueCreatedAfter      bool
	emptyIssueCreatedBefore     bool
	emptyIssueHasBlockers       bool
	emptyIssueBlocks            bool
	emptyIssueBlockedBy         bool
	emptyIssueAllTeams          bool
	emptyIssueSearch            bool
	emptyNextIssues             bool
	rankedNextIssues            bool
	expectedStateType           string
	expectedProjectID           string
	expectedAssigneeID          string
	expectedLabelID             string
	expectedCycleID             string
	expectedCreatedAfter        string
	expectedCreatedBefore       string
	expectedBlockedBy           string
	expectedIssueDeps           string
	expectedSearchQuery         string
	expectedReleaseSearchTerm   string
	expectedSemanticSearchQuery string
	expectedTypedSearchTerm     string
	emptyReleaseSearch          bool
	emptyProjectList            bool
	emptyProjectMembers         bool
	emptyProjectUpdates         bool
	emptyProjectMilestones      bool
	emptySLAConfigurations      bool
	emptySemanticSearch         bool
	emptySearchDocuments        bool
	emptySearchIssues           bool
	emptySearchProjects         bool
	emptyViewerDrafts           bool
	expectedCommentBody         string
	expectedCommentParentID     string
	expectedCreateDescription   string
	expectedUpdateDescription   string
	expectedStartAssigneeID     string
	expectedStartStateID        string
	expectedOrganizationURLKey  string
	expectedApplicationClientID string
	failOperation               string
	multiIssueList              bool
}

func (client commandFlowFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if request.OpName == client.failOperation {
		return errors.New("operation failed")
	}
	if err := client.requireExpectedVariables(request); err != nil {
		return err
	}

	payload, err := commandFlowPayload(request.OpName, client)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

func (client commandFlowFakeClient) requireExpectedVariables(request *graphql.Request) error {
	if client.expectedCreateDescription != "" && request.OpName == "IssueCreate" {
		return requireRequestVariable(
			request,
			[]string{"input", "description"},
			client.expectedCreateDescription,
			"create description",
		)
	}
	if client.expectedCommentBody != "" && request.OpName == "IssueCommentCreate" {
		return requireRequestVariable(request, []string{"input", "body"}, client.expectedCommentBody, "comment body")
	}
	if client.expectedCommentParentID != "" && request.OpName == "IssueCommentCreate" {
		return requireRequestVariable(
			request,
			[]string{"input", "parentId"},
			client.expectedCommentParentID,
			"comment parent id",
		)
	}
	if client.expectedUpdateDescription != "" && request.OpName == "IssueUpdate" {
		return requireRequestVariable(
			request,
			[]string{"input", "description"},
			client.expectedUpdateDescription,
			"update description",
		)
	}
	if err := client.requireExpectedIssueListVariables(request); err != nil {
		return err
	}
	if err := client.requireExpectedSearchVariables(request); err != nil {
		return err
	}
	if err := client.requireExpectedOrganizationVariables(request); err != nil {
		return err
	}
	if client.expectedApplicationClientID != "" && request.OpName == "applicationInfo" {
		return requireRequestVariable(request, []string{"clientId"}, client.expectedApplicationClientID, "application client id")
	}
	return client.requireExpectedIssueStartVariables(request)
}

func (client commandFlowFakeClient) requireExpectedSearchVariables(request *graphql.Request) error {
	if client.expectedSearchQuery != "" && request.OpName == "issueSearch" {
		return requireRequestVariable(request, []string{"query"}, client.expectedSearchQuery, "search query")
	}
	if client.expectedReleaseSearchTerm != "" && request.OpName == "releaseSearch" {
		return requireRequestVariable(request, []string{"term"}, client.expectedReleaseSearchTerm, "release search term")
	}
	if client.expectedSemanticSearchQuery != "" && request.OpName == "semanticSearch" {
		return requireRequestVariable(request, []string{"query"}, client.expectedSemanticSearchQuery, "semantic search query")
	}
	if client.expectedTypedSearchTerm != "" &&
		(request.OpName == "searchDocuments" ||
			request.OpName == "searchIssues" ||
			request.OpName == "searchProjects") {
		return requireRequestVariable(request, []string{"term"}, client.expectedTypedSearchTerm, "typed search term")
	}
	if client.expectedIssueDeps != "" && request.OpName == "IssueDependencies" {
		return requireRequestVariable(request, []string{"id"}, client.expectedIssueDeps, "issue deps id")
	}

	return nil
}

func (client commandFlowFakeClient) requireExpectedOrganizationVariables(request *graphql.Request) error {
	if client.expectedOrganizationURLKey != "" && request.OpName == "organizationExists" {
		return requireRequestVariable(request, []string{"urlKey"}, client.expectedOrganizationURLKey, "organization url key")
	}

	return nil
}

func (client commandFlowFakeClient) requireExpectedIssueStartVariables(request *graphql.Request) error {
	if client.expectedStartAssigneeID != "" && request.OpName == "IssueUpdate" {
		if err := requireRequestVariable(
			request,
			[]string{"input", "assigneeId"},
			client.expectedStartAssigneeID,
			"start assignee id",
		); err != nil {
			return err
		}
	}
	if client.expectedStartStateID != "" && request.OpName == "IssueUpdate" {
		return requireRequestVariable(request, []string{"input", "stateId"}, client.expectedStartStateID, "start state id")
	}

	return nil
}

func (client commandFlowFakeClient) requireExpectedIssueListVariables(request *graphql.Request) error {
	if client.expectedStateType != "" && request.OpName == "IssuesByTeamState" {
		return requireRequestVariable(request, []string{"stateType"}, client.expectedStateType, "state type")
	}
	if client.expectedProjectID != "" && request.OpName == "IssuesByTeamProject" {
		return requireRequestVariable(request, []string{"projectID"}, client.expectedProjectID, "project id")
	}
	if client.expectedAssigneeID != "" && request.OpName == "IssuesByTeamAssignee" {
		return requireRequestVariable(request, []string{"assigneeID"}, client.expectedAssigneeID, "assignee id")
	}
	if client.expectedLabelID != "" && request.OpName == "IssuesByTeamLabel" {
		return requireRequestVariable(request, []string{"labelID"}, client.expectedLabelID, "label id")
	}
	if client.expectedCycleID != "" && request.OpName == "IssuesByTeamCycle" {
		return requireRequestVariable(request, []string{"cycleID"}, client.expectedCycleID, "cycle id")
	}

	return client.requireExpectedDependencyIssueListVariables(request)
}

func (client commandFlowFakeClient) requireExpectedDependencyIssueListVariables(request *graphql.Request) error {
	if client.expectedCreatedAfter != "" && request.OpName == "IssuesByTeamCreatedAfter" {
		return requireRequestVariable(request, []string{"createdAfter"}, client.expectedCreatedAfter, "created after")
	}
	if client.expectedCreatedBefore != "" && request.OpName == "IssuesByTeamCreatedBefore" {
		return requireRequestVariable(request, []string{"createdBefore"}, client.expectedCreatedBefore, "created before")
	}
	if client.expectedBlockedBy != "" && request.OpName == "IssueBlockedIssues" {
		return requireRequestVariable(request, []string{"id"}, client.expectedBlockedBy, "blocked by issue")
	}

	return nil
}

func requireRequestVariable(request *graphql.Request, keys []string, expected string, label string) error {
	actual, err := requestVariable[string](request, keys...)
	if err != nil {
		return err
	}
	if actual != expected {
		return errors.New(label + " = " + actual)
	}

	return nil
}

func commandFlowPayload(operation string, fake commandFlowFakeClient) (string, error) {
	if payload, ok := commandFlowBasePayload(operation); ok {
		return payload, nil
	}

	if payload, ok := commandFlowTeamMembershipPayload(operation); ok {
		return payload, nil
	}
	if payload, ok := commandFlowAttachmentIssuePayload(operation); ok {
		return payload, nil
	}
	if payload, ok := commandFlowIssueVCSBranchPayload(operation); ok {
		return payload, nil
	}
	if payload, ok := commandFlowIssuePayload(operation, fake); ok {
		return payload, nil
	}
	if payload, ok := commandFlowProjectPayload(operation, fake); ok {
		return payload, nil
	}
	if payload, ok := commandFlowPeopleAndReferencePayload(operation, fake); ok {
		return payload, nil
	}
	if payload, ok := commandFlowOrganizationPayload(operation); ok {
		return payload, nil
	}

	return "", errors.New("missing fake response for " + operation)
}

func commandFlowBasePayload(operation string) (string, bool) {
	switch operation {
	case "Viewer":
		return `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`, true
	case "Teams":
		return `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "TargetProject":
		return `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`, true
	case "applicationInfo":
		return commandApplicationInfoPayload(), true
	case "agentActivities":
		return `{"agentActivities":{"nodes":[` + commandAgentActivityJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "agentActivity":
		return `{"agentActivity":` + commandAgentActivityJSON() + `}`, true
	case "agentSkills":
		return `{"agentSkills":{"nodes":[` + commandAgentSkillJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "agentSkill":
		return `{"agentSkill":` + commandAgentSkillJSON() + `}`, true
	case "externalUsers":
		return `{"externalUsers":{"nodes":[` + commandExternalUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "externalUser":
		return `{"externalUser":` + commandExternalUserJSON() + `}`, true
	case "rateLimitStatus":
		return commandRateLimitStatusPayload(), true
	default:
		return "", false
	}
}

func commandFlowOrganizationPayload(operation string) (string, bool) {
	switch operation {
	case "organizationExists":
		return `{"organizationExists":{"success":true,"exists":true}}`, true
	case "organization_labels":
		return `{"organization":{"labels":{"nodes":[` + commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "organization_projectLabels":
		return `{"organization":{"projectLabels":{"nodes":[` + commandProjectLabelJSON("project-label-id", "Roadmap", "#f2c94c") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "organization_teams":
		return `{"organization":{"teams":{"nodes":[` + commandTeamJSON(false) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "organization_templates":
		return `{"organization":{"templates":{"nodes":[` + commandTemplateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "organization_users":
		return `{"organization":{"users":{"nodes":[` + commandUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowTeamMembershipPayload(operation string) (string, bool) {
	switch operation {
	case "teamMemberships":
		return `{"teamMemberships":{"nodes":[` + commandTeamMembershipJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "teamMembership":
		return `{"teamMembership":` + commandTeamMembershipJSON() + `}`, true
	default:
		return "", false
	}
}

func commandRateLimitStatusPayload() string {
	return `{"rateLimitStatus":{"identifier":"api-key","kind":"api","limits":[{"type":"complexity","requestedAmount":1,"allowedAmount":1000,"period":60000,"remainingAmount":900,"reset":1720000000000}]}}`
}

func commandApplicationInfoPayload() string {
	return `{"applicationInfo":{"id":"app-id","clientId":"app-client-id","name":"Demo App","description":"Demo authorization app","developer":"Kyanite","developerUrl":"https://example.com","imageUrl":"https://example.com/app.png"}}`
}

func commandTeamMembershipJSON() string {
	return `{
		"id":"team-membership-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"owner":true,
		"sortOrder":1.5,
		"user":{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com","active":true,"guest":false,"admin":false},
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandAgentActivityJSON() string {
	return `{
		"id":"agent-activity-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"signal":"continue",
		"ephemeral":false,
		"agentSession":{"id":"agent-session-id"},
		"sourceComment":{"id":"comment-id"},
		"user":{"id":"user-id"},
		"content":{
			"__typename":"AgentActivityActionContent",
			"type":"action",
			"action":"read_file",
			"parameter":"README.md",
			"result":"Read file"
		}
	}`
}

func commandAgentSkillJSON() string {
	return `{
		"id":"agent-skill-id",
		"title":"Triage Helper",
		"body":"Use this skill for triage.",
		"description":"Helps triage issues",
		"slugId":"triage-helper",
		"teamId":"team-id",
		"shared":true,
		"icon":"sparkles",
		"color":"#5e6ad2",
		"recentUsageCount":3,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"lastUsedAt":"2026-06-20T12:00:00Z",
		"owner":{"id":"owner-id"},
		"creator":{"id":"creator-id"},
		"lastUpdatedBy":{"id":"updater-id"}
	}`
}

func commandExternalUserJSON() string {
	return `{
		"id":"external-user-id",
		"name":"External User",
		"displayName":"@external",
		"avatarUrl":"https://example.com/avatar.png",
		"lastSeen":"2026-06-19T12:00:00Z",
		"createdAt":"2026-06-18T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null
	}`
}

func commandFlowTeamChildPayload(operation string) (string, bool) {
	switch operation {
	case "team_cycles":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","cycles":{"nodes":[` +
			commandCycleJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_issues":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "state-id", "Todo", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_labels":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","labels":{"nodes":[` +
			commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_memberships":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","memberships":{"nodes":[` +
			commandTeamMembershipJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_projects":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_releasePipelines":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","releasePipelines":{"nodes":[` +
			commandReleasePipelineJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_states":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","states":{"nodes":[` +
			commandWorkflowStateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_templates":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","templates":{"nodes":[` +
			commandTemplateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowPeopleAndReferencePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowTeamChildPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowUserChildPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowLabelChildPayload(operation); ok {
		return payload, true
	}

	switch operation {
	case "Documents":
		return `{"documents":{"nodes":[` + commandDocumentJSON(
			"Spec",
			`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
		) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "document":
		return `{"document":` + commandDocumentJSON(
			"Team note",
			`"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null`,
		) + `}`, true
	case "IssueLabels":
		return `{"issueLabels":{"nodes":[` + commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issueLabel":
		return `{"issueLabel":` + commandLabelJSON("") + `}`, true
	case "team":
		return `{"team":` + commandTeamJSON(true) + `}`, true
	case "team_members":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","members":{"nodes":[` + commandUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "users":
		return `{"users":{"nodes":[` + commandUserJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "user":
		return `{"user":` + commandUserJSON() + `}`, true
	case "viewer":
		return `{"viewer":` + commandUserJSON() + `}`, true
	case "viewer_drafts":
		if fake.emptyViewerDrafts {
			return `{"viewer":{"drafts":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"viewer":{"drafts":{"nodes":[` + commandDraftJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	}

	return commandFlowStateAndCommentPayload(operation, fake)
}

func commandFlowUserChildPayload(operation string) (string, bool) {
	switch operation {
	case "user_assignedIssues":
		return commandFlowUserIssueListPayload("user", "assignedIssues"), true
	case "user_createdIssues":
		return commandFlowUserIssueListPayload("user", "createdIssues"), true
	case "user_delegatedIssues":
		return commandFlowUserIssueListPayload("user", "delegatedIssues"), true
	case "user_teamMemberships":
		return `{"user":{"teamMemberships":{"nodes":[` + commandTeamMembershipJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "user_teams":
		return `{"user":{"teams":{"nodes":[` + commandTeamJSON(false) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "viewer_assignedIssues":
		return commandFlowUserIssueListPayload("viewer", "assignedIssues"), true
	case "viewer_createdIssues":
		return commandFlowUserIssueListPayload("viewer", "createdIssues"), true
	case "viewer_delegatedIssues":
		return commandFlowUserIssueListPayload("viewer", "delegatedIssues"), true
	case "viewer_teamMemberships":
		return `{"viewer":{"teamMemberships":{"nodes":[` + commandTeamMembershipJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "viewer_teams":
		return `{"viewer":{"teams":{"nodes":[` + commandTeamJSON(false) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowStateAndCommentPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "workflowStates":
		return `{"workflowStates":{"nodes":[` + commandWorkflowStateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "workflowState":
		return `{"workflowState":` + commandWorkflowStateJSON() + `}`, true
	}

	return commandFlowInitiativePayload(operation, fake)
}

//nolint:gocyclo // The table-driven command-flow fake is intentionally centralized by operation name.
func commandFlowInitiativePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "initiatives":
		return `{"initiatives":{"nodes":[` + commandInitiativeJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "initiative":
		return `{"initiative":` + commandInitiativeJSON() + `}`, true
	case "initiative_history":
		return `{"initiative":{"history":{"nodes":[` + commandInitiativeHistoryJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_links":
		return `{"initiative":{"links":{"nodes":[` + commandEntityExternalLinkJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_subInitiatives":
		return `{"initiative":{"subInitiatives":{"nodes":[` + commandSubInitiativeJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_initiativeUpdates":
		return `{"initiative":{"initiativeUpdates":{"nodes":[` + commandInitiativeUpdateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_documents":
		return `{"initiative":{"documents":{"nodes":[` + commandDocumentJSON(
			"Spec",
			`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
		) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiative_projects":
		return `{"initiative":{"projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "initiativeRelations":
		return `{"initiativeRelations":{"nodes":[` + commandInitiativeRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "initiativeRelation":
		return `{"initiativeRelation":` + commandInitiativeRelationJSON() + `}`, true
	case "initiativeToProjects":
		return `{"initiativeToProjects":{"nodes":[` + commandInitiativeToProjectJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "initiativeToProject":
		return `{"initiativeToProject":` + commandInitiativeToProjectJSON() + `}`, true
	case "roadmapToProjects":
		return `{"roadmapToProjects":{"nodes":[` + commandRoadmapToProjectJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "roadmapToProject":
		return `{"roadmapToProject":` + commandRoadmapToProjectJSON() + `}`, true
	case "initiativeUpdates":
		return `{"initiativeUpdates":{"nodes":[` + commandInitiativeUpdateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "initiativeUpdate":
		return `{"initiativeUpdate":` + commandInitiativeUpdateJSON() + `}`, true
	case "comments":
		return `{"comments":{"nodes":[` + commandTopLevelCommentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "comment":
		return `{"comment":` + commandTopLevelCommentJSON() + `}`, true
	case "comment_botActor":
		return `{"comment":{"id":"comment-id","botActor":` + commandActorBotJSON() + `}}`, true
	case "comment_children":
		return `{"comment":{"id":"comment-id","children":{"nodes":[` +
			commandCommentMetadataJSONWithID("child-comment-id", "comment-id", "", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "comment_createdIssues":
		return `{"comment":{"id":"comment-id","createdIssues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	}

	return commandFlowExtraReadPayload(operation, fake)
}

//nolint:gocyclo // The table-driven command-flow fake is intentionally centralized by operation name.
func commandFlowExtraReadPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "auditEntryTypes":
		return `{"auditEntryTypes":[{"type":"user_login","description":"User logged in"}]}`, true
	case "notifications":
		return `{"notifications":{"nodes":[` + commandNotificationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "notification":
		return `{"notification":` + commandNotificationJSON() + `}`, true
	case "notificationSubscriptions":
		return `{"notificationSubscriptions":{"nodes":[` + commandNotificationSubscriptionJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "notificationSubscription":
		return `{"notificationSubscription":` + commandNotificationSubscriptionJSON() + `}`, true
	case "triageResponsibilities":
		return `{"triageResponsibilities":{"nodes":[` + commandTriageResponsibilityJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "triageResponsibility":
		return `{"triageResponsibility":` + commandTriageResponsibilityJSON() + `}`, true
	case "triageResponsibility_manualSelection":
		return `{"triageResponsibility":{"id":"triage-responsibility-id","manualSelection":{"userIds":["user-id","other-user-id"]}}}`, true
	case "slaConfigurations":
		if fake.emptySLAConfigurations {
			return `{"slaConfigurations":[]}`, true
		}
		return `{"slaConfigurations":[` + commandSLAConfigurationJSON() + `]}`, true
	case "semanticSearch":
		if fake.emptySemanticSearch {
			return `{"semanticSearch":{"results":[]}}`, true
		}
		return `{"semanticSearch":{"results":[` + commandSemanticSearchResultJSON() + `]}}`, true
	case "searchDocuments":
		if fake.emptySearchDocuments {
			return `{"searchDocuments":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":0}}`, true
		}
		return `{"searchDocuments":{"nodes":[` + commandSearchDocumentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":1}}`, true
	case "searchIssues":
		if fake.emptySearchIssues {
			return `{"searchIssues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":0}}`, true
		}
		return `{"searchIssues":{"nodes":[` + commandSearchIssueJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":1}}`, true
	case "searchProjects":
		if fake.emptySearchProjects {
			return `{"searchProjects":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":0}}`, true
		}
		return `{"searchProjects":{"nodes":[` + commandSearchProjectJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null},"totalCount":1}}`, true
	case "releasePipelines":
		return `{"releasePipelines":{"nodes":[` + commandReleasePipelineJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "releasePipeline":
		return `{"releasePipeline":` + commandReleasePipelineJSON() + `}`, true
	case "releasePipeline_releases":
		return `{"releasePipeline":{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "releasePipeline_stages":
		return `{"releasePipeline":{"stages":{"nodes":[` + commandReleaseStageJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "releasePipeline_teams":
		return `{"releasePipeline":{"teams":{"nodes":[` + commandTeamJSON(true) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "releaseStages":
		return `{"releaseStages":{"nodes":[` + commandReleaseStageJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "releaseStage":
		return `{"releaseStage":` + commandReleaseStageJSON() + `}`, true
	case "releaseStage_releases":
		return `{"releaseStage":{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "releases":
		return `{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "releaseSearch":
		if fake.emptyReleaseSearch {
			return `{"releaseSearch":[]}`, true
		}
		return `{"releaseSearch":[` + commandReleaseJSON() + `]}`, true
	case "release":
		return `{"release":` + commandReleaseJSON() + `}`, true
	case "release_history":
		return `{"release":{"history":{"nodes":[` + commandReleaseHistoryJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "release_documents":
		return `{"release":{"documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "release_issues":
		return `{"release":{"issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "release_links":
		return `{"release":{"links":{"nodes":[` + commandEntityExternalLinkJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "entityExternalLink":
		return `{"entityExternalLink":` + commandEntityExternalLinkJSON() + `}`, true
	case "releaseNotes":
		return `{"releaseNotes":{"nodes":[` + commandReleaseNoteJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "releaseNote":
		return `{"releaseNote":` + commandReleaseNoteJSON() + `}`, true
	case "issueToReleases":
		return `{"issueToReleases":{"nodes":[` + commandIssueToReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issueToRelease":
		return `{"issueToRelease":` + commandIssueToReleaseJSON() + `}`, true
	case "timeSchedules":
		return `{"timeSchedules":{"nodes":[` + commandTimeScheduleJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "timeSchedule":
		return `{"timeSchedule":` + commandTimeScheduleJSON() + `}`, true
	case "templates":
		return `{"templates":[` + commandTemplateJSON() + `]}`, true
	case "template":
		return `{"template":` + commandTemplateJSON() + `}`, true
	case "roadmaps":
		return `{"roadmaps":{"nodes":[` + commandRoadmapJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "roadmap":
		return `{"roadmap":` + commandRoadmapJSON() + `}`, true
	case "customViews":
		return `{"customViews":{"nodes":[` + commandCustomViewJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customViewHasSubscribers":
		return `{"customViewHasSubscribers":{"hasSubscribers":true}}`, true
	case "customView":
		return `{"customView":` + commandCustomViewJSON() + `}`, true
	case "customView_initiatives":
		return `{"customView":{"initiatives":{"nodes":[` + commandInitiativeJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "customView_issues":
		return `{"customView":{"issues":{"nodes":[` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "customView_organizationViewPreferences":
		return `{"customView":{"organizationViewPreferences":` + commandCustomViewPreferencesJSON("priority", "list") + `}}`, true
	case "customView_organizationViewPreferences_preferences":
		return `{"customView":{"organizationViewPreferences":{"preferences":` + commandCustomViewPreferenceValuesJSON("priority", "list") + `}}}`, true
	case "customView_projects":
		return `{"customView":{"projects":{"nodes":[` + commandProjectJSON("Pinned project", "Backlog", "backlog") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "customView_userViewPreferences":
		return `{"customView":{"userViewPreferences":` + commandCustomViewScopedPreferencesJSON("user", "updatedAt", "board") + `}}`, true
	case "customView_userViewPreferences_preferences":
		return `{"customView":{"userViewPreferences":{"preferences":` + commandCustomViewPreferenceValuesJSON("updatedAt", "board") + `}}}`, true
	case "customView_viewPreferencesValues":
		return `{"customView":{"viewPreferencesValues":` + commandCustomViewPreferenceValuesJSON("updatedAt", "board") + `}}`, true
	case "customers":
		return `{"customers":{"nodes":[` + commandCustomerJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customer":
		return `{"customer":` + commandCustomerJSON() + `}`, true
	case "customerNeeds":
		return `{"customerNeeds":{"nodes":[` + commandCustomerNeedJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customerNeed":
		return `{"customerNeed":` + commandCustomerNeedJSON() + `}`, true
	case "customerStatuses":
		return `{"customerStatuses":{"nodes":[` + commandCustomerStatusJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customerStatus":
		return `{"customerStatus":` + commandCustomerStatusJSON() + `}`, true
	case "customerTiers":
		return `{"customerTiers":{"nodes":[` + commandCustomerTierJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "customerTier":
		return `{"customerTier":` + commandCustomerTierJSON() + `}`, true
	case "favorites":
		return `{"favorites":{"nodes":[` + commandFavoriteJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "favorite_children":
		return `{"favorite":{"children":{"nodes":[` + commandFavoriteChildJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "favorite":
		return `{"favorite":` + commandFavoriteJSON() + `}`, true
	case "emojis":
		return `{"emojis":{"nodes":[` + commandEmojiJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "emoji":
		return `{"emoji":` + commandEmojiJSON() + `}`, true
	case "attachments":
		return `{"attachments":{"nodes":[` + commandAttachmentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "attachmentsForURL":
		return `{"attachmentsForURL":{"nodes":[` + commandAttachmentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "attachment":
		return `{"attachment":` + commandAttachmentJSON() + `}`, true
	default:
		return "", false
	}
}

func commandFlowLabelChildPayload(operation string) (string, bool) {
	switch operation {
	case "issueLabel_children":
		return `{"issueLabel":{"id":"label-id","name":"Bug","children":{"nodes":[` +
			commandNamedLabelJSON("child-label-id", "Mobile", "#56ccf2", "child label") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueLabel_issues":
		return `{"issueLabel":{"id":"label-id","name":"Bug","issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssuePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowIssueReadPayload(operation, fake); ok {
		return payload, true
	}

	return commandFlowIssueWritePayload(operation, fake)
}

func commandFlowAttachmentIssuePayload(operation string) (string, bool) {
	if !strings.HasPrefix(operation, "attachmentIssue") {
		return "", false
	}

	switch operation {
	case "attachmentIssue":
		return `{"attachmentIssue":` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`}`, true
	case "attachmentIssue_attachments":
		return `{"attachmentIssue":{"attachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_botActor":
		return `{"attachmentIssue":{"id":"issue-id","botActor":` + commandActorBotJSON() + `}}`, true
	case "attachmentIssue_children":
		return `{"attachmentIssue":{"children":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_documents":
		return `{"attachmentIssue":{"documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_formerAttachments":
		return `{"attachmentIssue":{"formerAttachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_history":
		return `{"attachmentIssue":{"history":{"nodes":[` +
			commandIssueHistoryJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_inverseRelations":
		return `{"attachmentIssue":{"inverseRelations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_labels":
		return `{"attachmentIssue":{"labels":{"nodes":[` +
			commandLabelJSON("label body") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_relations":
		return `{"attachmentIssue":{"relations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_releases":
		return `{"attachmentIssue":{"releases":{"nodes":[` +
			commandReleaseJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_stateHistory":
		return `{"attachmentIssue":{"id":"issue-id","stateHistory":{"nodes":[` +
			commandIssueStateSpanJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_subscribers":
		return `{"attachmentIssue":{"id":"issue-id","subscribers":{"nodes":[` +
			commandUserJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueVCSBranchPayload(operation string) (string, bool) {
	if !strings.HasPrefix(operation, "issueVcsBranchSearch") {
		return "", false
	}

	switch operation {
	case "issueVcsBranchSearch":
		return `{"issueVcsBranchSearch":` +
			commandIssueJSON("LIT-40", "Branch issue", "todo-state", "Todo", "unstarted") +
			`}`, true
	case "issueVcsBranchSearch_attachments":
		return `{"issueVcsBranchSearch":{"attachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_botActor":
		return `{"issueVcsBranchSearch":{"id":"issue-id","botActor":` + commandActorBotJSON() + `}}`, true
	case "issueVcsBranchSearch_children":
		return `{"issueVcsBranchSearch":{"children":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_documents":
		return `{"issueVcsBranchSearch":{"documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_formerAttachments":
		return `{"issueVcsBranchSearch":{"formerAttachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_history":
		return `{"issueVcsBranchSearch":{"history":{"nodes":[` +
			commandIssueHistoryJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_inverseRelations":
		return `{"issueVcsBranchSearch":{"inverseRelations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_labels":
		return `{"issueVcsBranchSearch":{"labels":{"nodes":[` +
			commandLabelJSON("label body") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_relations":
		return `{"issueVcsBranchSearch":{"relations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_releases":
		return `{"issueVcsBranchSearch":{"releases":{"nodes":[` +
			commandReleaseJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_stateHistory":
		return `{"issueVcsBranchSearch":{"id":"issue-id","stateHistory":{"nodes":[` +
			commandIssueStateSpanJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_subscribers":
		return `{"issueVcsBranchSearch":{"id":"issue-id","subscribers":{"nodes":[` +
			commandUserJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueReadPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowIssueListPayload(operation, fake); ok {
		return payload, true
	}
	if payload, ok := commandFlowIssueChildPayload(operation, fake); ok {
		return payload, true
	}

	switch operation {
	case "issueSearch":
		if fake.emptyIssueSearch {
			return `{"issueSearch":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"issueSearch":{"nodes":[` + commandIssueJSON("LIT-3", "Search result", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "NextIssuesByTeam":
		if fake.emptyNextIssues {
			return emptyCommandIssuesPayload(), true
		}
		if fake.rankedNextIssues {
			return `{"issues":{"nodes":[` +
				commandIssueWithNextRankJSON("LIT-28", "Low priority standalone", 4, "Low", "2026-05-01T12:00:00Z", 0) + `,` +
				commandIssueWithNextRankJSON("LIT-29", "Urgent standalone", 1, "Urgent", "2026-06-01T12:00:00Z", 0) + `,` +
				commandIssueWithNextRankJSON("LIT-30", "Unblocks checkout", 2, "High", "2026-06-10T12:00:00Z", 2) +
				`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"issues":{"nodes":[` + commandIssueWithNextRankJSON("LIT-27", "Next issue", 0, "No priority", "2026-06-01T12:00:00Z", 0) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issue":
		return `{"issue":` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `}`, true
	case "IssueDependencies":
		return commandFlowIssueDependenciesPayload(), true
	case "issueRelations":
		return `{"issueRelations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issueRelation":
		return `{"issueRelation":` + commandIssueRelationJSON() + `}`, true
	case "issue_comments":
		if fake.emptyIssueComments {
			return `{"issue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[{"id":"comment-id","body":"First comment","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","parentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueChildPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "issue_attachments":
		return `{"issue":{"attachments":{"nodes":[` + commandAttachmentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_botActor":
		return `{"issue":{"id":"issue-id","botActor":` + commandActorBotJSON() + `}}`, true
	case "issue_children":
		if fake.emptyIssueChildren {
			return `{"issue":{"children":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"issue":{"children":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_documents":
		return `{"issue":{"documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":null,"team":null,"issue":{"id":"issue-id","identifier":"LIT-1","title":"Detail issue"},"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_formerAttachments":
		return `{"issue":{"formerAttachments":{"nodes":[` + commandAttachmentJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_history":
		return `{"issue":{"history":{"nodes":[` + commandIssueHistoryJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_inverseRelations":
		return `{"issue":{"inverseRelations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_labels":
		return `{"issue":{"labels":{"nodes":[` + commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_relations":
		return `{"issue":{"relations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_releases":
		return `{"issue":{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_stateHistory":
		return `{"issue":{"id":"issue-id","stateHistory":{"nodes":[` +
			commandIssueStateSpanJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_subscribers":
		return `{"issue":{"id":"issue-id","subscribers":{"nodes":[` +
			commandUserJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueDependenciesPayload() string {
	return `{"issue":{
		"id":"issue-id",
		"identifier":"LIT-1",
		"parent":` + commandIssueJSON("LIT-25", "Parent issue", "todo-state", "Todo", "unstarted") + `,
		"children":{
			"nodes":[` + commandIssueJSON("LIT-26", "Child issue", "todo-state", "Todo", "unstarted") + `],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		},
		"relations":{
			"nodes":[{"id":"blocks-relation","type":"blocks","relatedIssue":` + commandIssueJSON("LIT-23", "Blocked issue", "todo-state", "Todo", "unstarted") + `}],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		},
		"inverseRelations":{
			"nodes":[{"id":"blocked-by-relation","type":"blocks","issue":` + commandIssueJSON("LIT-24", "Blocker issue", "todo-state", "Todo", "unstarted") + `}],
			"pageInfo":{"hasNextPage":false,"endCursor":null}
		}
	}}`
}

func commandFlowIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowBroadIssueListPayload(operation, fake); ok {
		return payload, true
	}
	if payload, ok := commandFlowDependencyIssueListPayload(operation, fake); ok {
		return payload, true
	}

	switch operation {
	case "IssuesByTeamState":
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-2", "Started issue", "started-state", "Started", fake.expectedStateType) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamProject":
		if fake.emptyIssueProject {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-4", "Project issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamAssignee":
		return commandFlowAssigneeIssueListPayload(fake), true
	case "IssuesByTeamLabel":
		if fake.emptyIssueLabel {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-7", "Labeled issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCycle":
		if fake.emptyIssueCycle {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-8", "Cycle issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCreatedAfter":
		if fake.emptyIssueCreatedAfter {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-9", "Recent issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamCreatedBefore":
		if fake.emptyIssueCreatedBefore {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-19", "Older issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	default:
		return "", false
	}
}

func commandFlowDependencyIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "IssuesByTeamHasBlockers":
		if fake.emptyIssueHasBlockers {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-21", "Blocked issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeamBlocks":
		if fake.emptyIssueBlocks {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-22", "Blocking issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssueBlockedIssues":
		if fake.emptyIssueBlockedBy {
			return `{"issue":{"id":"issue-id","identifier":"LIT-1","relations":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","relations":{"nodes":[{"id":"relation-id","type":"blocks","relatedIssue":` + commandIssueJSON("LIT-23", "Blocked by issue", "todo-state", "Todo", "unstarted") + `}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowBroadIssueListPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "issues":
		if fake.emptyIssueAllTeams {
			return emptyCommandIssuesPayload(), true
		}
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-20", "All-team issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "IssuesByTeam":
		return commandFlowUnfilteredIssueListPayload(fake), true
	default:
		return "", false
	}
}

func commandFlowUnfilteredIssueListPayload(fake commandFlowFakeClient) string {
	if fake.emptyIssueList {
		return emptyCommandIssuesPayload()
	}
	if fake.multiIssueList {
		return `{"issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Alpha issue", "todo-state", "Todo", "unstarted") + `,` +
			commandIssueJSON("LIT-2", "Zebra issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
	}

	return `{"issues":{"nodes":[` + commandIssueJSON("LIT-1", "Listed issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func commandFlowAssigneeIssueListPayload(fake commandFlowFakeClient) string {
	if fake.emptyIssueMine {
		return emptyCommandIssuesPayload()
	}
	if fake.expectedAssigneeID == "assignee-id" {
		return `{"issues":{"nodes":[` + commandIssueJSON("LIT-6", "Assigned issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
	}

	return `{"issues":{"nodes":[` + commandIssueJSON("LIT-5", "Mine issue", "todo-state", "Todo", "unstarted") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func emptyCommandIssuesPayload() string {
	return `{"issues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
}

func commandFlowIssueWritePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "IssueCreate":
		return `{"issueCreate":{"success":true,"issue":` + commandIssueJSON("LIT-2", "Created issue", "todo-state", "Todo", "unstarted") + `}}`, true
	case "IssueUpdate":
		if fake.expectedStartStateID != "" {
			return `{"issueUpdate":{"success":true,"issue":` +
				commandIssueJSON("LIT-1", "Started issue", "started-state", "Started", "started") + `}}`, true
		}
		return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSON("LIT-1", "Updated issue", "todo-state", "Todo", "unstarted") + `}}`, true
	case "IssueCommentCreate":
		return `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"Looks good","url":"https://linear.app/comment/comment-id","issue":` + commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") + `}}}`, true
	case "CompletedWorkflowStates":
		return `{"workflowStates":{"nodes":[{"id":"done-state","name":"Done","type":"completed","position":1}]}}`, true
	case "StartedWorkflowStates":
		return `{"workflowStates":{"nodes":[{"id":"started-state","name":"Started","type":"started","position":1}]}}`, true
	case "IssueClose":
		return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSON("LIT-1", "Closed issue", "done-state", "Done", "completed") + `}}`, true
	default:
		return "", false
	}
}

func commandFlowProjectPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	if payload, ok := commandFlowProjectStatusPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowProjectLabelPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowProjectRelationPayload(operation); ok {
		return payload, true
	}
	if payload, ok := commandFlowProjectReadPayload(operation, fake); ok {
		return payload, true
	}

	return commandFlowProjectWritePayload(operation)
}

//nolint:gocyclo // The fake payload switch mirrors the project command operation surface.
func commandFlowProjectReadPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "Projects":
		if fake.emptyProjectList {
			return `{"team":{"projects":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"team":{"projects":{"nodes":[` + commandProjectJSON("Listed project", "Backlog", "backlog") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project":
		return `{"project":` + commandProjectJSON("Detail project", "Backlog", "backlog") + `}`, true
	case "project_attachments":
		return `{"project":{"id":"project-id","name":"Detail project","attachments":{"nodes":[` +
			commandAttachmentJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_documents":
		return `{"project":{"id":"project-id","name":"Detail project","documents":{"nodes":[` +
			commandDocumentJSON(
				"Spec",
				`"project":{"id":"project-id","name":"Pinned project"},"team":null,"issue":null,"cycle":null`,
			) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_externalLinks":
		return `{"project":{"id":"project-id","name":"Detail project","externalLinks":{"nodes":[` +
			commandEntityExternalLinkJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_history":
		return `{"project":{"id":"project-id","name":"Detail project","history":{"nodes":[` +
			commandProjectHistoryJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_initiativeToProjects":
		return `{"project":{"id":"project-id","name":"Detail project","initiativeToProjects":{"nodes":[` +
			commandInitiativeToProjectJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_initiatives":
		return `{"project":{"id":"project-id","name":"Detail project","initiatives":{"nodes":[` +
			commandInitiativeJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_inverseRelations":
		return `{"project":{"id":"project-id","name":"Detail project","inverseRelations":{"nodes":[` +
			commandProjectRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_issues":
		return `{"project":{"id":"project-id","name":"Detail project","issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_comments":
		return `{"project":{"id":"project-id","name":"Detail project","comments":{"nodes":[` +
			commandCommentMetadataJSON("project-id", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_labels":
		return `{"project":{"id":"project-id","name":"Detail project","labels":{"nodes":[` +
			commandProjectLabelJSON("project-label-id", "Roadmap", "#f2c94c") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_members":
		if fake.emptyProjectMembers {
			return `{"project":{"id":"project-id","name":"Detail project","members":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com"}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_needs":
		return `{"project":{"id":"project-id","name":"Detail project","needs":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_relations":
		return `{"project":{"id":"project-id","name":"Detail project","relations":{"nodes":[` +
			commandProjectRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_teams":
		return `{"project":{"id":"project-id","name":"Detail project","teams":{"nodes":[` +
			commandTeamJSON(true) +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "ProjectUpdates":
		if fake.emptyProjectUpdates {
			return `{"project":{"id":"project-id","name":"Detail project","projectUpdates":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","projectUpdates":{"nodes":[{"id":"project-update-id","body":"First update","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projectUpdates":
		if fake.emptyProjectUpdates {
			return `{"projectUpdates":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"projectUpdates":{"nodes":[` + commandProjectUpdateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectUpdate":
		return `{"projectUpdate":` + commandProjectUpdateJSON() + `}`, true
	case "projectUpdate_comments":
		return `{"projectUpdate":{"id":"project-update-id","comments":{"nodes":[` +
			commandCommentMetadataJSON("", "project-update-id") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "project_projectMilestones":
		if fake.emptyProjectMilestones {
			return `{"project":{"id":"project-id","name":"Detail project","projectMilestones":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projectMilestone":
		return `{"projectMilestone":` + commandProjectMilestoneJSON("Launch milestone", "next") + `}`, true
	case "projectMilestone_issues":
		return `{"projectMilestone":{"id":"project-milestone-id","name":"Launch milestone","issues":{"nodes":[` +
			commandIssueJSON("LIT-2", "Milestone issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowProjectStatusPayload(operation string) (string, bool) {
	switch operation {
	case "projectStatuses":
		return `{"projectStatuses":{"nodes":[` +
			commandProjectStatusJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectStatus":
		return `{"projectStatus":` + commandProjectStatusJSON() + `}`, true
	default:
		return "", false
	}
}

func commandFlowProjectLabelPayload(operation string) (string, bool) {
	switch operation {
	case "projectLabels":
		return `{"projectLabels":{"nodes":[` +
			commandProjectLabelJSON("project-label-id", "Roadmap", "#f2c94c") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectLabel":
		return `{"projectLabel":` + commandProjectLabelJSON("project-label-id", "Roadmap", "#f2c94c") + `}`, true
	case "projectLabel_children":
		return `{"projectLabel":{"id":"project-label-id","name":"Roadmap","children":{"nodes":[` +
			commandProjectLabelJSON("child-project-label-id", "Mobile", "#56ccf2") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projectLabel_projects":
		return `{"projectLabel":{"id":"project-label-id","name":"Roadmap","projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowProjectRelationPayload(operation string) (string, bool) {
	switch operation {
	case "projectRelations":
		return `{"projectRelations":{"nodes":[` +
			commandProjectRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "projectRelation":
		return `{"projectRelation":` + commandProjectRelationJSON() + `}`, true
	default:
		return "", false
	}
}

func commandFlowProjectWritePayload(operation string) (string, bool) {
	switch operation {
	case "ProjectMilestoneCreate":
		return `{"projectMilestoneCreate":{"success":true,"projectMilestone":` + commandProjectMilestoneJSON("Created milestone", "next") + `}}`, true
	case "ProjectMilestoneUpdate":
		return `{"projectMilestoneUpdate":{"success":true,"projectMilestone":` + commandProjectMilestoneJSON("Updated milestone", "done") + `}}`, true
	case "ProjectCreate":
		return `{"projectCreate":{"success":true,"project":` + commandProjectJSON("Created project", "Backlog", "backlog") + `}}`, true
	case "ProjectUpdate":
		return `{"projectUpdate":{"success":true,"project":` + commandProjectJSON("Updated project", "Started", "started") + `}}`, true
	case "ProjectArchive":
		return `{"projectArchive":{"success":true,"entity":` + commandProjectJSON("Archived project", "Canceled", "canceled") + `}}`, true
	default:
		return "", false
	}
}

func commandIssueJSON(identifier string, title string, stateID string, state string, stateType string) string {
	return `{
		"id":"issue-id",
		"description":"Existing description",
		"identifier":"` + identifier + `",
		"title":"` + title + `",
		"branchName":"` + strings.ToLower(identifier) + `-` + strings.ToLower(strings.ReplaceAll(title, " ", "-")) + `",
		"url":"https://linear.app/kyanite/issue/` + identifier + `",
		"priority":0,
		"priorityLabel":"No priority",
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"state":{"id":"` + stateID + `","name":"` + state + `","type":"` + stateType + `"},
		"assignee":null,
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

func commandIssueRelationJSON() string {
	return `{
		"id":"issue-relation-id",
		"type":"blocks",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},
		"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}
	}`
}

func commandIssueToReleaseJSON() string {
	return `{
		"id":"issue-to-release-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"issue":{"id":"issue-id"},
		"release":{"id":"release-id"}
	}`
}

func commandIssueHistoryJSON() string {
	return `{
		"id":"issue-history-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"actorId":"user-id",
		"updatedDescription":true,
		"issue":{"id":"issue-id"}
	}`
}

func commandIssueStateSpanJSON() string {
	return `{
		"id":"issue-state-span-id",
		"stateId":"started-state",
		"startedAt":"2026-06-19T12:00:00Z",
		"endedAt":null,
		"state":{"id":"started-state","name":"Started","type":"started"}
	}`
}

func commandCycleJSON() string {
	return `{
		"id":"cycle-id",
		"number":7,
		"name":"Planning cycle",
		"description":"Cycle body",
		"startsAt":"2026-06-01T00:00:00Z",
		"endsAt":"2026-07-15T00:00:00Z",
		"completedAt":null,
		"progress":0.5,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandIssueWithNextRankJSON(
	identifier string,
	title string,
	priority int,
	priorityLabel string,
	createdAt string,
	unblocksCount int,
) string {
	return strings.TrimSuffix(commandIssueJSON(identifier, title, "todo-state", "Todo", "unstarted"), "\n\t}") +
		`,
		"priority":` + strconv.Itoa(priority) + `,
		"priorityLabel":"` + priorityLabel + `",
		"createdAt":"` + createdAt + `",
		"relations":{"nodes":[` + commandBlockingRelationsJSON(unblocksCount) + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}
	}`
}

func commandBlockingRelationsJSON(count int) string {
	relations := make([]string, 0, count)
	for i := range count {
		relations = append(relations, fmt.Sprintf(`{"type":"blocks","relatedIssue":{"id":"blocked-%d","state":{"type":"unstarted"}}}`, i))
	}

	return strings.Join(relations, ",")
}

func commandProjectJSON(name string, status string, statusType string) string {
	return `{
		"id":"project-id",
		"name":"` + name + `",
		"description":"description",
		"slugId":"` + name + `",
		"url":"https://linear.app/kyanite/project/project-id",
		"priority":0,
		"status":{"id":"status-id","name":"` + status + `","type":"` + statusType + `"},
		"lead":null,
		"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl"}]}
	}`
}

func commandProjectUpdateJSON() string {
	return `{
		"id":"project-update-id",
		"body":"First update",
		"health":"onTrack",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"url":"https://linear.app/project-update/project-update-id",
		"project":{"id":"project-id","name":"Pinned project"},
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandProjectStatusJSON() string {
	return `{
		"id":"project-status-id",
		"name":"Backlog",
		"description":"Ready for planning",
		"type":"backlog",
		"color":"#bec2c8",
		"position":1,
		"archivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z"
	}`
}

func commandProjectLabelJSON(id string, name string, color string) string {
	return `{
		"id":"` + id + `",
		"name":"` + name + `",
		"description":"Project label",
		"color":"` + color + `",
		"isGroup":false,
		"lastAppliedAt":"2026-06-19T12:00:00Z",
		"retiredAt":null,
		"archivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"parent":null
	}`
}

func commandProjectRelationJSON() string {
	return `{
		"id":"project-relation-id",
		"type":"blocks",
		"anchorType":"project",
		"relatedAnchorType":"project",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"project":{"id":"project-id","name":"Pinned project"},
		"projectMilestone":null,
		"relatedProject":{"id":"related-project-id","name":"Related project"},
		"relatedProjectMilestone":null,
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandProjectHistoryJSON() string {
	return `{
		"id":"project-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"status","from":"Backlog","to":"Started"}],
		"project":{"id":"project-id"}
	}`
}

func commandInitiativeUpdateJSON() string {
	return `{
		"id":"initiative-update-id",
		"body":"First initiative update",
		"health":"onTrack",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"url":"https://linear.app/initiative-update/initiative-update-id",
		"slugId":"initiative-update-slug",
		"commentCount":1,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandInitiativeRelationJSON() string {
	return `{
		"id":"initiative-relation-id",
		"sortOrder":1.5,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"relatedInitiative":{"id":"child-initiative-id","name":"Child initiative"},
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandInitiativeToProjectJSON() string {
	return `{
		"id":"initiative-to-project-id",
		"sortOrder":"1",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"initiative":{"id":"initiative-id","name":"Platform"},
		"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}
	}`
}

func commandRoadmapToProjectJSON() string {
	return `{
		"id":"roadmap-to-project-id",
		"sortOrder":"1",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"roadmap":{"id":"roadmap-id","name":"Platform roadmap"},
		"project":{"id":"project-id","name":"Pinned project","slugId":"pinned-project","url":"https://linear.app/project/project-id"}
	}`
}

func commandProjectMilestoneJSON(name string, status string) string {
	return `{
		"id":"project-milestone-id",
		"name":"` + name + `",
		"description":"milestone body",
		"targetDate":"2026-06-30",
		"status":"` + status + `",
		"progress":0.5,
		"sortOrder":1,
		"project":` + commandProjectJSON("Pinned project", "Backlog", "backlog") + `
	}`
}

func commandDocumentJSON(title string, parents string) string {
	return `{
		"id":"document-id",
		"title":"` + title + `",
		"slugId":"document-slug",
		"archivedAt":null,
		` + parents + `
	}`
}

func commandLabelJSON(description string) string {
	return commandNamedLabelJSON("label-id", "Bug", "#ff0000", description)
}

func commandNamedLabelJSON(id string, name string, color string, description string) string {
	descriptionPayload := "null"
	if description != "" {
		descriptionPayload = `"` + description + `"`
	}

	return `{
		"id":"` + id + `",
		"name":"` + name + `",
		"description":` + descriptionPayload + `,
		"color":"` + color + `",
		"isGroup":false,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandTeamJSON(includeDescription bool) string {
	descriptionPayload := "null"
	if includeDescription {
		descriptionPayload = `"team body"`
	}

	return `{
		"id":"team-id",
		"key":"LIT",
		"name":"linctl",
		"description":` + descriptionPayload + `,
		"archivedAt":null,
		"organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}
	}`
}

func commandUserJSON() string {
	return `{
		"id":"user-id",
		"name":"omer",
		"displayName":"Omer",
		"email":"omer@example.com",
		"active":true,
		"guest":false,
		"admin":true
	}`
}

func commandDraftJSON() string {
	return `{
		"id":"draft-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"issue":{"id":"issue-id","identifier":"LIT-3","title":"Draft issue"},
		"project":null,
		"projectUpdate":null,
		"initiative":null,
		"initiativeUpdate":null,
		"parentComment":null,
		"customerNeed":null,
		"team":null
	}`
}

func commandFlowUserIssueListPayload(parent string, field string) string {
	return `{"` + parent + `":{"` + field + `":{"nodes":[` +
		commandIssueJSON("LIT-1", "Detail issue", "todo-state", "Todo", "unstarted") +
		`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`
}

func commandWorkflowStateJSON() string {
	return `{
		"id":"workflow-state-id",
		"name":"Started",
		"type":"started",
		"color":"#f2c94c",
		"position":2,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}

func commandInitiativeJSON() string {
	return `{
		"id":"initiative-id",
		"name":"Platform",
		"description":"Platform initiative",
		"status":"Active",
		"priority":2,
		"targetDate":"2026-12-31",
		"slugId":"platform-init",
		"url":"https://linear.app/kyanite/initiative/platform-init"
	}`
}

func commandSubInitiativeJSON() string {
	return `{
		"id":"child-initiative-id",
		"name":"Child platform",
		"description":"Child initiative",
		"status":"Planned",
		"priority":1,
		"targetDate":"2026-11-30",
		"slugId":"child-platform",
		"url":"https://linear.app/kyanite/initiative/child-platform"
	}`
}

func commandInitiativeHistoryJSON() string {
	return `{
		"id":"initiative-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"status","from":"Planned","to":"Active"}],
		"initiative":{"id":"initiative-id"}
	}`
}

func commandCustomViewJSON() string {
	return `{
		"id":"custom-view-id",
		"name":"My issues",
		"description":"Saved issue view",
		"modelName":"Issue",
		"shared":true,
		"color":"#5e6ad2",
		"slugId":"my-issues"
	}`
}

func commandCustomViewPreferencesJSON(ordering string, layout string) string {
	return commandCustomViewScopedPreferencesJSON("organization", ordering, layout)
}

func commandCustomViewScopedPreferencesJSON(scope string, ordering string, layout string) string {
	return `{
		"id":"view-preferences-id",
		"createdAt":"2026-06-01T12:00:00Z",
		"updatedAt":"2026-06-01T12:01:00Z",
		"archivedAt":null,
		"type":"` + scope + `",
		"viewType":"customView",
		"preferences":` + commandCustomViewPreferenceValuesJSON(ordering, layout) + `
	}`
}

func commandCustomViewPreferenceValuesJSON(ordering string, layout string) string {
	return `{
		"layout":"` + layout + `",
		"viewOrdering":"` + ordering + `",
		"viewOrderingDirection":"Descending",
		"issueGrouping":"status",
		"issueSubGrouping":"priority",
		"showCompletedIssues":"all",
		"showArchivedItems":true,
		"showEmptyGroups":true,
		"hiddenColumns":["column-id"],
		"hiddenRows":["row-id"],
		"hiddenGroupsList":["group-id"],
		"columnOrderBoard":["board-column-id"],
		"columnOrderList":["list-column-id"],
		"projectLayout":"timeline",
		"projectViewOrdering":"priority",
		"projectGrouping":"status",
		"projectSubGrouping":"lead",
		"projectShowEmptyGroups":"all",
		"projectShowEmptySubGroups":"all"
	}`
}

func commandCustomerJSON() string {
	return `{
		"id":"customer-id",
		"name":"Acme",
		"domains":["acme.example"],
		"externalIds":["crm-acme"],
		"slackChannelId":"slack-channel-id",
		"status":{"id":"status-id","name":"Active"},
		"tier":{"id":"tier-id","name":"Enterprise"},
		"owner":{"id":"user-id","displayName":"Omer"},
		"revenue":120000,
		"size":42,
		"approximateNeedCount":3,
		"slugId":"acme",
		"url":"https://linear.app/kyanite/customer/acme"
	}`
}

func commandRoadmapJSON() string {
	return `{
		"id":"roadmap-id",
		"name":"Platform roadmap",
		"description":"Roadmap body",
		"color":"#5e6ad2",
		"slugId":"platform-roadmap",
		"sortOrder":1,
		"archivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"url":"https://linear.app/kyanite/roadmap/platform-roadmap",
		"creator":{"id":"user-id","displayName":"Omer"},
		"owner":{"id":"owner-id","displayName":"Owner"}
	}`
}

func commandTimeScheduleJSON() string {
	return `{
		"id":"time-schedule-id",
		"name":"Primary on-call",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"externalId":"pd-primary",
		"externalUrl":"https://example.com/schedule",
		"integration":{"id":"integration-id"},
		"entries":[{"startsAt":"2026-06-20T00:00:00Z","endsAt":"2026-06-21T00:00:00Z","userId":"user-id","userEmail":"omer@example.com"}]
	}`
}

func commandTemplateJSON() string {
	return `{
		"id":"template-id",
		"name":"Bug report",
		"type":"issue",
		"description":"Bug report template",
		"icon":"bug",
		"color":"#ff0000",
		"sortOrder":1,
		"lastAppliedAt":"2026-06-19T12:00:00Z",
		"createdAt":"2026-06-18T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"pipeline":{"id":"pipeline-id"},
		"creator":{"id":"creator-id"},
		"lastUpdatedBy":{"id":"updated-by-id"},
		"inheritedFrom":{"id":"parent-template-id"}
	}`
}

func commandCustomerNeedJSON() string {
	return `{
		"id":"customer-need-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"priority":1,
		"body":"Need body",
		"content":"Need content",
		"url":"https://example.com/need",
		"customer":{"id":"customer-id","name":"Acme"},
		"issue":{"id":"issue-id","identifier":"LIT-1","title":"Need issue"},
		"project":{"id":"project-id","name":"Customer project"}
	}`
}

func commandCustomerStatusJSON() string {
	return `{
		"id":"customer-status-id",
		"name":"active",
		"displayName":"Active",
		"color":"#00ff00",
		"description":"Active customers",
		"position":1,
		"archivedAt":null
	}`
}

func commandCustomerTierJSON() string {
	return `{
		"id":"customer-tier-id",
		"name":"enterprise",
		"displayName":"Enterprise",
		"color":"#0000ff",
		"description":"Enterprise customers",
		"position":2,
		"archivedAt":null
	}`
}

func commandFavoriteJSON() string {
	return `{
		"id":"favorite-id",
		"type":"issue",
		"folderName":null,
		"url":"https://linear.app/kyanite/issue/LIT-1"
	}`
}

func commandFavoriteChildJSON() string {
	return `{
		"id":"favorite-child-id",
		"type":"project",
		"folderName":null,
		"url":"https://linear.app/kyanite/project/project-id"
	}`
}

func commandEmojiJSON() string {
	return `{
		"id":"emoji-id",
		"name":"party",
		"url":"https://linear.app/kyanite/emoji/party.png",
		"source":"custom"
	}`
}

func commandAttachmentJSON() string {
	return `{
		"id":"attachment-id",
		"title":"Linked PR",
		"subtitle":"feat: add thing",
		"url":"https://github.com/kyanite/linctl/pull/1",
		"sourceType":"github"
	}`
}

func commandNotificationJSON() string {
	return `{
		"__typename":"IssueNotification",
		"id":"notification-id",
		"type":"issueMention",
		"category":"mentions",
		"title":"Mentioned you",
		"subtitle":"LIT-1",
		"url":"https://linear.app/kyanite/issue/LIT-1",
		"inboxUrl":"https://linear.app/kyanite/inbox/notification-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"readAt":null,
		"emailedAt":null,
		"snoozedUntilAt":null,
		"unsnoozedAt":null,
		"user":{"id":"user-id","displayName":"Omer"},
		"actor":{"id":"actor-id","displayName":"Ada"},
		"externalUserActor":null
	}`
}

func commandNotificationSubscriptionJSON() string {
	return `{
		"__typename":"ProjectNotificationSubscription",
		"id":"notification-subscription-id",
		"active":true,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"contextViewType":null,
		"userContextViewType":null,
		"subscriber":{"id":"user-id","displayName":"Omer"},
		"customer":null,
		"customView":null,
		"cycle":null,
		"initiative":null,
		"label":null,
		"project":{"id":"project-id","name":"Roadmap"},
		"team":null,
		"user":null
	}`
}

func commandTriageResponsibilityJSON() string {
	return `{
		"id":"triage-responsibility-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"action":"notify",
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"timeSchedule":{"id":"time-schedule-id","name":"Primary rotation"},
		"currentUser":{"id":"user-id","displayName":"Omer"},
		"manualSelection":{"userIds":["user-id","other-user-id"]}
	}`
}

func commandSLAConfigurationJSON() string {
	return `{
		"id":"sla-configuration-id",
		"name":"First response",
		"conditions":{"priority":{"eq":1}},
		"sla":3600000,
		"slaType":"all",
		"removesSla":false
	}`
}

func commandSemanticSearchResultJSON() string {
	return `{
		"id":"issue-id",
		"type":"issue",
		"issue":{"id":"issue-id","identifier":"LIT-3","title":"Search result","url":"https://linear.app/kyanite/issue/LIT-3"},
		"project":null,
		"initiative":null,
		"document":null
	}`
}

func commandSearchDocumentJSON() string {
	return `{
		"id":"search-document-id",
		"title":"Search spec",
		"slugId":"search-spec",
		"url":"https://linear.app/kyanite/document/search-spec",
		"project":null,
		"initiative":null,
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"issue":null,
		"release":null,
		"cycle":null
	}`
}

func commandSearchIssueJSON() string {
	return `{
		"id":"search-issue-id",
		"identifier":"LIT-30",
		"title":"Search issue",
		"url":"https://linear.app/kyanite/issue/LIT-30",
		"team":{"id":"team-id","key":"LIT","name":"linctl"},
		"state":{"id":"state-id","name":"Todo","type":"unstarted"},
		"project":{"id":"project-id","name":"Pinned project"}
	}`
}

func commandSearchProjectJSON() string {
	return `{
		"id":"search-project-id",
		"name":"Search project",
		"slugId":"search-project",
		"url":"https://linear.app/kyanite/project/search-project",
		"status":{"id":"status-id","name":"Backlog","type":"backlog"},
		"lead":{"id":"user-id","name":"omer","displayName":"Omer"},
		"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl"}]}
	}`
}

func commandReleasePipelineJSON() string {
	return `{
		"id":"release-pipeline-id",
		"name":"Production",
		"slugId":"production",
		"type":"scheduled",
		"isProduction":true,
		"autoGenerateReleaseNotesOnCompletion":true,
		"approximateReleaseCount":4,
		"url":"https://linear.app/kyanite/releases/production",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"trashed":null,
		"includePathPatterns":["services/api/**"],
		"releaseNoteTemplate":{"id":"template-id"},
		"latestReleaseNote":{"id":"release-note-id"}
	}`
}

func commandReleaseStageJSON() string {
	return `{
		"id":"release-stage-id",
		"name":"Started",
		"color":"#00ff00",
		"type":"started",
		"position":2,
		"frozen":false,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"pipeline":{"id":"release-pipeline-id","name":"Production","slugId":"production"}
	}`
}

func commandReleaseJSON() string {
	return `{
		"id":"release-id",
		"name":"Mobile 1.2.3",
		"slugId":"mobile-1-2-3",
		"version":"v1.2.3",
		"description":"Release body",
		"commitSha":"abc123",
		"issueCount":3,
		"trashed":null,
		"url":"https://linear.app/kyanite/release/mobile-1-2-3",
		"startDate":"2026-06-20",
		"targetDate":"2026-06-30",
		"startedAt":"2026-06-20T12:00:00Z",
		"completedAt":null,
		"canceledAt":null,
		"autoArchivedAt":null,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-20T12:00:00Z",
		"archivedAt":null,
		"pipeline":{"id":"release-pipeline-id","name":"Production","slugId":"production"},
		"stage":{"id":"release-stage-id","name":"Started","type":"started"},
		"releaseNotes":[{"id":"release-note-id","title":"Launch notes","slugId":"launch-notes"}],
		"creator":{"id":"user-id","displayName":"Omer"}
	}`
}

func commandReleaseHistoryJSON() string {
	return `{
		"id":"release-history-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"entries":[{"type":"stage","from":"planned","to":"started"}],
		"release":{"id":"release-id"}
	}`
}

func commandEntityExternalLinkJSON() string {
	return `{
		"id":"release-link-id",
		"createdAt":"2026-06-03T12:00:00Z",
		"updatedAt":"2026-06-03T12:01:00Z",
		"archivedAt":null,
		"url":"https://example.com/runbook",
		"label":"Runbook",
		"sortOrder":1.5,
		"creator":{"id":"user-id","displayName":"Omer"},
		"initiative":null,
		"project":null
	}`
}

func commandReleaseNoteJSON() string {
	return `{
		"id":"release-note-id",
		"title":"Launch notes",
		"slugId":"launch-notes",
		"generationStatus":"completed",
		"releaseCount":2,
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-20T12:00:00Z",
		"archivedAt":null,
		"pipeline":{"id":"release-pipeline-id","name":"Production","slugId":"production"},
		"firstRelease":{"id":"release-id","name":"Mobile 1.2.2","version":"v1.2.2"},
		"lastRelease":{"id":"release-id","name":"Mobile 1.2.3","version":"v1.2.3"}
	}`
}

func commandTopLevelCommentJSON() string {
	return `{
		"id":"comment-id",
		"body":"First comment",
		"url":"https://linear.app/comment/comment-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"editedAt":null,
		"resolvedAt":null,
		"parentId":null,
		"issueId":"issue-id",
		"projectId":null,
		"projectUpdateId":null,
		"initiativeId":null,
		"initiativeUpdateId":null,
		"documentContentId":null,
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandCommentMetadataJSON(projectID string, projectUpdateID string) string {
	return commandCommentMetadataJSONWithID("comment-id", "", projectID, projectUpdateID)
}

func commandCommentMetadataJSONWithID(
	id string,
	parentID string,
	projectID string,
	projectUpdateID string,
) string {
	return `{
		"id":"` + id + `",
		"url":"https://linear.app/comment/` + id + `",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"editedAt":null,
		"resolvedAt":null,
		"parentId":` + commandNullableStringJSON(parentID) + `,
		"issueId":null,
		"projectId":` + commandNullableStringJSON(projectID) + `,
		"projectUpdateId":` + commandNullableStringJSON(projectUpdateID) + `,
		"initiativeId":null,
		"initiativeUpdateId":null,
		"documentContentId":null,
		"user":{"id":"user-id","name":"omer","displayName":"Omer"}
	}`
}

func commandActorBotJSON() string {
	return `{
		"id":"bot-actor-id",
		"type":"github",
		"subType":"issue",
		"name":"GitHub",
		"userDisplayName":"octocat",
		"avatarUrl":"https://example.com/github.png"
	}`
}

func commandNullableStringJSON(value string) string {
	if value == "" {
		return `null`
	}

	return `"` + value + `"`
}

var _ graphql.Client = commandFlowFakeClient{}

func requestVariable[T comparable](request *graphql.Request, keys ...string) (T, error) {
	var zero T
	payload, err := json.Marshal(request.Variables)
	if err != nil {
		return zero, err
	}
	var variables map[string]any
	if err := json.Unmarshal(payload, &variables); err != nil {
		return zero, err
	}
	current := any(variables)
	for _, key := range keys {
		object, ok := current.(map[string]any)
		if !ok {
			return zero, errors.New("request variable is not an object")
		}
		value, ok := object[key]
		if !ok {
			return zero, errors.New("request variable missing " + key)
		}
		current = value
	}
	value, ok := current.(T)
	if !ok {
		return zero, errors.New("request variable has unexpected type")
	}

	return value, nil
}
