package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func Test_CommandFlows_print_workflow_state_pages_as_json(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
		also     string
	}{
		{
			name:     "list",
			args:     []string{"--json", "workflow-state", "list", "--limit", "1"},
			contains: `"workflow_states": [`,
			also:     `"team_key": "LIT"`,
		},
		{
			name:     "issues",
			args:     []string{"--json", "workflow-state", "issues", "workflow-state-id", "--limit", "1"},
			contains: `"issues": [`,
			also:     `"identifier": "LIT-1"`,
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
			require.Contains(t, output.String(), test.also)
		})
	}
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
		{name: "issue figma file key search", args: []string{"issue", "figma-file-key-search", "figma-key"}, operation: "issueFigmaFileKeySearch", contains: "search issues by Figma file key"},
		{name: "issue priority values", args: []string{"issue", "priority-values"}, operation: "issuePriorityValues", contains: "list issue priority values"},
		{name: "issue filter suggestion", args: []string{"issue", "filter-suggestion", "started issues"}, operation: "issueFilterSuggestion", contains: "get issue filter suggestion"},
		{name: "issue title suggestion", args: []string{"issue", "title-suggestion", "Customer asks for faster exports"}, operation: "issueTitleSuggestionFromCustomerRequest", contains: "get issue title suggestion"},
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
		{name: "issue create from template", args: []string{"issue", "create", "--template", "template-id"}, operation: "templateContent", contains: "get template content"},
		{name: "issue update", args: []string{"issue", "update", "LIT-1", "--title", "Updated issue"}, operation: "IssueUpdate", contains: "update issue LIT-1"},
		{name: "issue start state", args: []string{"issue", "start", "LIT-1"}, operation: "StartedWorkflowStates", contains: "list started workflow states"},
		{name: "issue start update", args: []string{"issue", "start", "LIT-1"}, operation: "IssueUpdate", contains: "start issue LIT-1"},
		{name: "issue comment", args: []string{"issue", "comment", "LIT-1", "--body", "Looks good"}, operation: "IssueCommentCreate", contains: "comment on issue LIT-1"},
		{name: "issue reply", args: []string{"issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"}, operation: "IssueCommentCreate", contains: "comment on issue LIT-1"},
		{name: "comment list", args: []string{"comment", "list"}, operation: "comments", contains: "list comments"},
		{name: "comment get", args: []string{"comment", "get", "comment-id"}, operation: "comment", contains: "get comment comment-id"},
		{name: "comment update", args: []string{"comment", "update", "comment-id", "--body", "New body"}, operation: "CommentUpdate", contains: "update comment comment-id"},
		{name: "comment delete", args: []string{"comment", "delete", "comment-id"}, operation: "CommentDelete", contains: "delete comment comment-id"},
		{name: "issue close", args: []string{"issue", "close", "LIT-1"}, operation: "IssueClose", contains: "close issue LIT-1"},
		{name: "issue relate", args: []string{"issue", "relate", "LIT-1", "LIT-2", "--type", "related"}, operation: "IssueRelationCreate", contains: "create issue relation"},
		{name: "issue unrelate", args: []string{"issue", "unrelate", "issue-relation-id"}, operation: "IssueRelationDelete", contains: "delete issue relation issue-relation-id"},
		{name: "project list target resolve", args: []string{"project", "list"}, operation: "Teams", contains: "resolve teams"},
		{name: "project list", args: []string{"project", "list"}, operation: "Projects", contains: "list projects"},
		{name: "project all", args: []string{"project", "all"}, operation: "projects", contains: "list projects"},
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
		{name: "project updates", args: []string{"project", "updates", "project-id"}, operation: "project_projectUpdates", contains: "list project updates project-id"},
		{name: "project filter suggestion", args: []string{"project", "filter-suggestion", "started projects"}, operation: "projectFilterSuggestion", contains: "get project filter suggestion"},
		{name: "project update list", args: []string{"project-update", "list"}, operation: "projectUpdates", contains: "list project updates"},
		{name: "project update get", args: []string{"project-update", "get", "project-update-id"}, operation: "projectUpdate", contains: "get project update project-update-id"},
		{name: "project update comments", args: []string{"project-update", "comments", "project-update-id"}, operation: "projectUpdate_comments", contains: "list project update comments project-update-id"},
		{name: "project update create", args: []string{"project-update", "create", "project-id", "--health", "on-track", "--body", "Posted update"}, operation: "ProjectUpdateCreate", contains: "create project update"},
		{name: "project milestone all", args: []string{"project-milestone", "all"}, operation: "projectMilestones", contains: "list project milestones"},
		{name: "project status project count", args: []string{"project-status", "project-count", "project-status-id"}, operation: "projectStatusProjectCount", contains: "get project status project count project-status-id"},
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
		{name: "document comments", args: []string{"document", "comments", "document-id"}, operation: "document_comments", contains: "list document comments document-id"},
		{name: "document create", args: []string{"document", "create", "--title", "Created doc"}, operation: "DocumentCreate", contains: "create document"},
		{name: "document update", args: []string{"document", "update", "document-id", "--title", "Updated doc"}, operation: "DocumentUpdate", contains: "update document document-id"},
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
		{name: "team git automation states", args: []string{"team", "git-automation-states", "team-id"}, operation: "team_gitAutomationStates", contains: "list team git automation states team-id"},
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
		{name: "workflow state issues", args: []string{"workflow-state", "issues", "workflow-state-id"}, operation: "workflowState_issues", contains: "list workflow state issues workflow-state-id"},
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
		{name: "initiative update comments", args: []string{"initiative-update", "comments", "initiative-update-id"}, operation: "initiativeUpdate_comments", contains: "list initiative update comments initiative-update-id"},
		{name: "roadmap list", args: []string{"roadmap", "list"}, operation: "roadmaps", contains: "list roadmaps"},
		{name: "roadmap get", args: []string{"roadmap", "get", "roadmap-id"}, operation: "roadmap", contains: "get roadmap roadmap-id"},
		{name: "roadmap projects", args: []string{"roadmap", "projects", "roadmap-id"}, operation: "roadmap_projects", contains: "list roadmap projects roadmap-id"},
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
		{name: "customer need project attachment", args: []string{"customer-need", "project-attachment", "customer-need-id"}, operation: "customerNeed_projectAttachment", contains: "get customer need project attachment customer-need-id"},
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
