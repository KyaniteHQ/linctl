package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CommandFlows_print_json_for_read_and_comment_commands(t *testing.T) {
	tests := []struct {
		args        []string
		keys        []string
		wantStderr  string
		nonObjectOK bool
	}{
		{args: []string{"--json", "target"}, keys: []string{"viewer", "org", "team", "expected", "resolved", "confirmed"}},
		{args: []string{"--json", "doctor"}, keys: []string{"config", "token", "target", "viewer"}},
		{args: []string{"--json", "whoami"}, keys: []string{"id", "display_name", "email"}},
		{args: []string{"--json", "organization", "exists", "kyanite"}, keys: []string{"exists"}},
		{args: []string{"--json", "--fields", "id,name,type,team_key", "organization", "templates", "--limit", "1"}, keys: []string{"templates"}},
		{args: []string{"--json", "rate-limit", "status"}, keys: []string{"identifier", "kind", "limits"}},
		{args: []string{"--json", "notification", "list", "--limit", "1"}, keys: []string{"notifications", "has_next_page"}},
		{args: []string{"--json", "notification", "get", "notification-id"}, keys: []string{"id"}},
		{args: []string{"--json", "notification", "subscription", "list", "--limit", "1"}, keys: []string{"notification_subscriptions", "has_next_page"}},
		{args: []string{"--json", "notification", "subscription", "get", "notification-subscription-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,team_key,action", "triage-responsibility", "list", "--limit", "1"}, keys: []string{"triage_responsibilities"}},
		{args: []string{"--json", "triage-responsibility", "get", "triage-responsibility-id"}, keys: []string{"id"}},
		{args: []string{"--json", "triage-responsibility", "manual-selection", "triage-responsibility-id"}, keys: []string{"id"}},
		{args: []string{"--json", "release-pipeline", "list", "--limit", "1"}, keys: []string{"release_pipelines", "has_next_page"}},
		{args: []string{"--json", "release-pipeline", "get", "release-pipeline-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,pipeline_id,stage_id", "release-pipeline", "releases", "release-pipeline-id", "--limit", "1"}, keys: []string{"releases"}},
		{args: []string{"--json", "--fields", "id,pipeline_id,type", "release-pipeline", "stages", "release-pipeline-id", "--limit", "1"}, keys: []string{"release_stages"}},
		{args: []string{"--json", "release-stage", "list", "--limit", "1"}, keys: []string{"release_stages", "has_next_page"}},
		{args: []string{"--json", "release-stage", "get", "release-stage-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,pipeline_id,stage_id", "release-stage", "releases", "release-stage-id", "--limit", "1"}, keys: []string{"releases"}},
		{args: []string{"--json", "release", "list", "--limit", "1"}, keys: []string{"releases", "has_next_page"}},
		{args: []string{"--json", "release", "search", "mobile", "--limit", "1"}, keys: []string{"releases", "has_next_page"}},
		{args: []string{"--json", "release", "get", "release-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,release_id,entry_count", "release", "history", "release-id", "--limit", "1"}, keys: []string{"history"}},
		{args: []string{"--json", "--fields", "id,label,url", "release", "links", "release-id", "--limit", "1"}, keys: []string{"links"}},
		{args: []string{"--json", "--fields", "id,label,url", "external-link", "get", "release-link-id"}, keys: []string{"id", "label", "url"}},
		{args: []string{"--json", "release-note", "list", "--limit", "1"}, keys: []string{"release_notes", "has_next_page"}},
		{args: []string{"--json", "release-note", "get", "release-note-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,client_id,name", "application", "info", "app-client-id"}, keys: []string{"id", "client_id", "name"}},
		{args: []string{"--json", "--fields", "id,content_type,agent_session_id", "agent-activity", "list", "--limit", "1"}, keys: []string{"agent_activities"}},
		{args: []string{"--json", "agent-activity", "get", "agent-activity-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,title,shared", "agent-skill", "list", "--limit", "1"}, keys: []string{"agent_skills"}},
		{args: []string{"--json", "agent-skill", "get", "agent-skill-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,slug_id,status", "agent-session", "list", "--limit", "1"}, keys: []string{"agent_sessions"}},
		{args: []string{"--json", "agent-session", "get", "agent-session-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,display_name,last_seen", "external-user", "list", "--limit", "1"}, keys: []string{"external_users"}},
		{args: []string{"--json", "external-user", "get", "external-user-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "type,description", "audit-entry", "types"}, keys: []string{"audit_entry_types"}},
		{args: []string{"--json", "--fields", "type,id,key,title", "semantic-search", "agent search", "--limit", "2"}, keys: []string{"results"}},
		{args: []string{"--json", "--fields", "id,title,parent_type", "search", "documents", "agent search", "--limit", "1"}, keys: []string{"documents"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "search", "issues", "agent search", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,name,slug_id", "search", "projects", "agent search", "--limit", "1"}, keys: []string{"projects"}},
		{args: []string{"--json", "next", "--dry-run"}, keys: []string{"id", "identifier", "title"}},
		{args: []string{"--json", "issue", "list", "--limit", "1"}, keys: []string{"issues", "has_next_page"}},
		{args: []string{"--json", "issue", "search", "needle", "--limit", "1"}, keys: []string{"issues", "has_next_page"}},
		{args: []string{"--json", "issue", "deps", "LIT-1", "--limit", "2"}, keys: []string{"identifier", "children", "blocks", "blocked_by", "has_next_page"}},
		{args: []string{"--json", "issue", "attachments", "LIT-1", "--limit", "1"}, keys: []string{"attachments", "has_next_page"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "issue", "children", "LIT-1", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,title,parent_type", "issue", "documents", "LIT-1", "--limit", "1"}, keys: []string{"documents"}},
		{args: []string{"--json", "issue", "former-attachments", "LIT-1", "--limit", "1"}, keys: []string{"attachments", "has_next_page"}},
		{args: []string{"--json", "--fields", "id,issue_id,updated_description", "issue", "history", "LIT-1", "--limit", "1"}, keys: []string{"history"}},
		{args: []string{"--json", "--fields", "id,type,issue_identifier,related_issue_identifier", "issue", "inverse-relations", "LIT-1", "--limit", "1"}, keys: []string{"relations"}},
		{args: []string{"--json", "--fields", "id,name,color", "issue", "labels", "LIT-1", "--limit", "1"}, keys: []string{"labels"}},
		{args: []string{"--json", "--fields", "id,type,issue_identifier,related_issue_identifier", "issue", "relations", "LIT-1", "--limit", "1"}, keys: []string{"relations"}},
		{args: []string{"--json", "--fields", "id,name,version", "issue", "releases", "LIT-1", "--limit", "1"}, keys: []string{"releases"}},
		{args: []string{"--json", "--fields", "id,type,issue_identifier,related_issue_identifier", "issue-relation", "list", "--limit", "1"}, keys: []string{"relations"}},
		{args: []string{"--json", "issue-relation", "get", "issue-relation-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,issue_id,release_id", "issue-to-release", "list", "--limit", "1"}, keys: []string{"associations"}},
		{args: []string{"--json", "issue-to-release", "get", "issue-to-release-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,display_name", "document", "comments", "document-id", "--limit", "1"}, keys: []string{"comments"}},
		{args: []string{"--json", "issue", "pr", "LIT-1"}, keys: []string{"title", "body", "command"}},
		{args: []string{"--json", "issue", "start", "LIT-1"}, keys: []string{"id", "identifier", "title"}},
		{args: []string{"--json", "issue", "comment", "LIT-1", "--body", "Looks good"}, keys: []string{"id", "body", "issue"}},
		{args: []string{"--json", "issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"}, keys: []string{"id", "body", "issue"}},
		{args: []string{"--json", "issue", "relate", "LIT-1", "LIT-2", "--type", "related"}, keys: []string{"id", "type"}},
		{args: []string{"--json", "issue", "unrelate", "issue-relation-id"}, keys: []string{"id", "status"}},
		{args: []string{"--json", "issue", "link", "https://example.com/pr/1", "LIT-1"}, keys: []string{"id", "title", "url"}},
		{args: []string{"--json", "--fields", "id,display_name", "issue", "comments", "LIT-1", "--limit", "1"}, keys: []string{"comments"}},
		{args: []string{"--json", "--fields", "identifier,title", "issue", "figma-file-key-search", "figma-key", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "issue", "priority-values"}, nonObjectOK: true},
		{args: []string{"--json", "issue", "filter-suggestion", "started issues"}, keys: []string{"filter", "log_id"}},
		{args: []string{"--json", "issue", "title-suggestion", "Customer asks for faster exports"}, keys: []string{"title"}},
		{args: []string{"--json", "--fields", "id,display_name", "comment", "list", "--limit", "1"}, keys: []string{"comments"}},
		{args: []string{"--json", "comment", "get", "comment-id"}, keys: []string{"id"}},
		{args: []string{"--json", "comment", "update", "comment-id", "--body", "New body"}, keys: []string{"id", "body"}},
		{args: []string{"--json", "comment", "delete", "comment-id"}, keys: []string{"id", "status"}},
		{args: []string{"--json", "--fields", "id,display_name", "initiative-update", "comments", "initiative-update-id", "--limit", "1"}, keys: []string{"comments"}},
		{args: []string{"--json", "project", "list", "--limit", "1"}, keys: []string{"projects", "has_next_page"}},
		{args: []string{"--json", "project", "all", "--limit", "1"}, keys: []string{"projects", "has_next_page"}},
		{args: []string{"--json", "project", "members", "project-id", "--limit", "1"}, keys: []string{"project_id", "project_name", "members", "has_next_page"}},
		{args: []string{"--json", "--fields", "id,health,display_name", "project", "updates", "project-id", "--limit", "1"}, keys: []string{"updates"}},
		{args: []string{"--json", "project", "filter-suggestion", "started projects"}, keys: []string{"filter", "log_id"}},
		{args: []string{"--json", "--fields", "id,health,project_id", "project-update", "list", "--limit", "1"}, keys: []string{"updates"}},
		{args: []string{"--json", "project-update", "get", "project-update-id"}, keys: []string{"id"}},
		{
			args:       []string{"--json", "project-update", "create", "project-id", "--health", "on-track", "--body", "Posted update"},
			keys:       []string{"id", "health", "body"},
			wantStderr: `note: health "on-track" normalized to "onTrack"` + "\n",
		},
		{args: []string{"--json", "--fields", "id,name,status", "project-milestone", "all", "--limit", "1"}, keys: []string{"milestones"}},
		{args: []string{"--json", "--fields", "id,name,status", "project-milestone", "list", "project-id", "--limit", "1"}, keys: []string{"milestones"}},
		{args: []string{"--json", "project-status", "list", "--limit", "1"}, keys: []string{"project_statuses", "has_next_page"}},
		{args: []string{"--json", "project-status", "get", "project-status-id"}, keys: []string{"id"}},
		{args: []string{"--json", "project-label", "list", "--limit", "1"}, keys: []string{"project_labels", "has_next_page"}},
		{args: []string{"--json", "project-label", "children", "project-label-id", "--limit", "1"}, keys: []string{"project_labels", "has_next_page"}},
		{args: []string{"--json", "project-label", "projects", "project-label-id", "--limit", "1"}, keys: []string{"projects", "has_next_page"}},
		{args: []string{"--json", "--fields", "id,type,project_id,related_project_id", "project-relation", "list", "--limit", "1"}, keys: []string{"relations"}},
		{args: []string{"--json", "project-relation", "get", "project-relation-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,title,parent_type", "document", "list", "--limit", "1"}, keys: []string{"documents"}},
		{args: []string{"--json", "document", "create", "--title", "Created doc"}, keys: []string{"id", "title"}},
		{args: []string{"--json", "document", "update", "document-id", "--title", "Updated doc"}, keys: []string{"id", "title"}},
		{args: []string{"--json", "--fields", "id,name,color", "label", "list", "--limit", "1"}, keys: []string{"labels"}},
		{args: []string{"--json", "--fields", "id,key,name", "team", "list", "--limit", "1"}, keys: []string{"teams"}},
		{args: []string{"--json", "--fields", "id,name,status", "team", "cycles", "team-id", "--limit", "1"}, keys: []string{"cycles"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "team", "issues", "team-id", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,name,color", "team", "labels", "team-id", "--limit", "1"}, keys: []string{"labels"}},
		{args: []string{"--json", "--fields", "id,display_name,email", "team", "members", "team-id", "--limit", "1"}, keys: []string{"members"}},
		{args: []string{"--json", "--fields", "id,team_key,user_id,owner", "team", "memberships", "team-id", "--limit", "1"}, keys: []string{"memberships"}},
		{args: []string{"--json", "--fields", "id,name,status", "team", "projects", "team-id", "--limit", "1"}, keys: []string{"projects"}},
		{args: []string{"--json", "--fields", "id,name,type", "team", "release-pipelines", "team-id", "--limit", "1"}, keys: []string{"release_pipelines"}},
		{args: []string{"--json", "--fields", "id,name,type", "team", "states", "team-id", "--limit", "1"}, keys: []string{"workflow_states"}},
		{args: []string{"--json", "team", "git-automation-states", "team-id", "--limit", "1"}, keys: []string{"git_automation_states", "has_next_page"}},
		{args: []string{"--json", "--fields", "id,name,type,team_key", "team", "templates", "team-id", "--limit", "1"}, keys: []string{"templates"}},
		{args: []string{"--json", "--fields", "id,team_key,user_id,owner", "team-membership", "list", "--limit", "1"}, keys: []string{"memberships"}},
		{args: []string{"--json", "team-membership", "get", "team-membership-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,display_name,email", "user", "list", "--limit", "1"}, keys: []string{"users"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "user", "assigned-issues", "user-id", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "user", "created-issues", "user-id", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "user", "delegated-issues", "user-id", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,team_key,user_id,owner", "user", "team-memberships", "user-id", "--limit", "1"}, keys: []string{"memberships"}},
		{args: []string{"--json", "--fields", "id,key,name", "user", "teams", "user-id", "--limit", "1"}, keys: []string{"teams"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "user", "my-assigned-issues", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "user", "my-created-issues", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,identifier,title", "user", "my-delegated-issues", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "--fields", "id,team_key,user_id,owner", "user", "my-team-memberships", "--limit", "1"}, keys: []string{"memberships"}},
		{args: []string{"--json", "--fields", "id,key,name", "user", "my-teams", "--limit", "1"}, keys: []string{"teams"}},
		{args: []string{"--json", "time-schedule", "list", "--limit", "1"}, keys: []string{"time_schedules", "has_next_page"}},
		{args: []string{"--json", "time-schedule", "get", "time-schedule-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,name,sla_type", "sla-configuration", "list", "team-id"}, keys: []string{"sla_configurations"}},
		{args: []string{"--json", "--fields", "id,name,type,team_key", "template", "list", "--limit", "1"}, keys: []string{"templates"}},
		{args: []string{"--json", "template", "get", "template-id"}, keys: []string{"id"}},
		{args: []string{"--json", "initiative", "list", "--limit", "1"}, keys: []string{"initiatives", "has_next_page"}},
		{args: []string{"--json", "initiative", "get", "initiative-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,initiative_id,entry_count", "initiative", "history", "initiative-id", "--limit", "1"}, keys: []string{"history"}},
		{args: []string{"--json", "--fields", "id,label,url", "initiative", "links", "initiative-id", "--limit", "1"}, keys: []string{"links"}},
		{args: []string{"--json", "--fields", "id,name,status", "initiative", "sub-initiatives", "initiative-id", "--limit", "1"}, keys: []string{"initiatives"}},
		{args: []string{"--json", "--fields", "id,health,initiative_id", "initiative", "updates", "initiative-id", "--limit", "1"}, keys: []string{"updates"}},
		{args: []string{"--json", "--fields", "id,parent_initiative_id,related_initiative_id", "initiative-relation", "list", "--limit", "1"}, keys: []string{"relations"}},
		{args: []string{"--json", "initiative-relation", "get", "initiative-relation-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,initiative_id,project_id", "initiative-to-project", "list", "--limit", "1"}, keys: []string{"associations"}},
		{args: []string{"--json", "initiative-to-project", "get", "initiative-to-project-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,health,initiative_id", "initiative-update", "list", "--limit", "1"}, keys: []string{"updates"}},
		{args: []string{"--json", "initiative-update", "get", "initiative-update-id"}, keys: []string{"id"}},
		{args: []string{"--json", "roadmap", "list", "--limit", "1"}, keys: []string{"roadmaps", "has_next_page"}},
		{args: []string{"--json", "roadmap", "get", "roadmap-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,name,status", "roadmap", "projects", "roadmap-id", "--limit", "1"}, keys: []string{"projects"}},
		{args: []string{"--json", "--fields", "id,roadmap_id,project_id", "roadmap-to-project", "list", "--limit", "1"}, keys: []string{"associations"}},
		{args: []string{"--json", "roadmap-to-project", "get", "roadmap-to-project-id"}, keys: []string{"id"}},
		{args: []string{"--json", "custom-view", "list", "--limit", "1"}, keys: []string{"custom_views", "has_next_page"}},
		{args: []string{"--json", "custom-view", "subscribers", "custom-view-id"}, keys: []string{"id", "has_subscribers"}},
		{args: []string{"--json", "custom-view", "get", "custom-view-id"}, keys: []string{"id"}},
		{args: []string{"--json", "--fields", "id,name,status", "custom-view", "initiatives", "custom-view-id", "--limit", "1"}, keys: []string{"initiatives"}},
		{args: []string{"--json", "--fields", "identifier,title,state", "custom-view", "issues", "custom-view-id", "--limit", "1"}, keys: []string{"issues"}},
		{args: []string{"--json", "custom-view", "organization-preferences", "custom-view-id"}, keys: []string{"id", "type", "values"}},
		{args: []string{"--json", "custom-view", "organization-preferences", "values", "custom-view-id"}, keys: []string{"layout", "view_ordering"}},
		{args: []string{"--json", "--fields", "id,name,status", "custom-view", "projects", "custom-view-id", "--limit", "1"}, keys: []string{"projects"}},
		{args: []string{"--json", "custom-view", "user-preferences", "custom-view-id"}, keys: []string{"id", "type", "values"}},
		{args: []string{"--json", "custom-view", "user-preferences", "values", "custom-view-id"}, keys: []string{"layout", "view_ordering"}},
		{args: []string{"--json", "custom-view", "preference-values", "custom-view-id"}, keys: []string{"layout", "view_ordering"}},
		{args: []string{"--json", "customer", "list", "--limit", "1"}, keys: []string{"customers", "has_next_page"}},
		{args: []string{"--json", "customer", "get", "customer-id"}, keys: []string{"id"}},
		{args: []string{"--json", "customer-need", "list", "--limit", "1"}, keys: []string{"customer_needs", "has_next_page"}},
		{args: []string{"--json", "customer-need", "get", "customer-need-id"}, keys: []string{"id"}},
		{args: []string{"--json", "customer-need", "project-attachment", "customer-need-id"}, keys: []string{"customer_need_id", "attachment"}},
		{args: []string{"--json", "customer-status", "list", "--limit", "1"}, keys: []string{"customer_statuses", "has_next_page"}},
		{args: []string{"--json", "customer-status", "get", "customer-status-id"}, keys: []string{"id"}},
		{args: []string{"--json", "customer-tier", "list", "--limit", "1"}, keys: []string{"customer_tiers", "has_next_page"}},
		{args: []string{"--json", "customer-tier", "get", "customer-tier-id"}, keys: []string{"id"}},
		{args: []string{"--json", "favorite", "list", "--limit", "1"}, keys: []string{"favorites", "has_next_page"}},
		{args: []string{"--json", "favorite", "children", "favorite-folder-id", "--limit", "1"}, keys: []string{"favorites", "has_next_page"}},
		{args: []string{"--json", "favorite", "get", "favorite-id"}, keys: []string{"id"}},
		{args: []string{"--json", "emoji", "list", "--limit", "1"}, keys: []string{"emojis", "has_next_page"}},
		{args: []string{"--json", "emoji", "get", "emoji-id"}, keys: []string{"id"}},
		{args: []string{"--json", "attachment", "list", "--limit", "1"}, keys: []string{"attachments", "has_next_page"}},
		{args: []string{"--json", "attachment", "url", "https://github.com/kyanite/linctl/pull/1", "--limit", "1"}, keys: []string{"attachments", "has_next_page"}},
		{args: []string{"--json", "attachment", "get", "attachment-id"}, keys: []string{"id"}},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.args, " "), func(t *testing.T) {
			output := bytes.Buffer{}
			stderr := bytes.Buffer{}
			restore := useCommandRuntime(t, commandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetErr(&stderr)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, test.wantStderr, stderr.String())
			if test.nonObjectOK {
				requireJSONOutputValue(t, output.String())
				return
			}
			requireJSONOutputObject(t, output.String(), test.keys...)
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

func requireJSONOutputObject(t *testing.T, output string, requiredKeys ...string) map[string]any {
	t.Helper()
	value := requireJSONOutputValue(t, output)
	envelope, ok := value.(map[string]any)
	require.Truef(t, ok, "expected JSON object, got %T", value)
	require.NotEmpty(t, envelope)
	for _, key := range requiredKeys {
		require.Contains(t, envelope, key)
	}

	return envelope
}

func requireJSONOutputValue(t *testing.T, output string) any {
	t.Helper()
	trimmed := strings.TrimSpace(output)
	require.NotEmpty(t, trimmed)
	require.True(t, strings.HasSuffix(output, "\n"), "expected JSON output to be newline terminated")

	var value any
	require.NoError(t, json.Unmarshal([]byte(trimmed), &value))

	return value
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
		{name: "issue unrelate", args: []string{"--id-only", "issue", "unrelate", "issue-relation-id"}, output: "issue-relation-id\n"},
		{name: "issue link", args: []string{"--id-only", "issue", "link", "https://example.com/pr/1", "LIT-1"}, output: "attachment-id\n"},
		{name: "agent session get", args: []string{"--id-only", "agent-session", "get", "agent-session-id"}, output: "agent-session-id\n"},
		{name: "team git automation states", args: []string{"--id-only", "team", "git-automation-states", "team-id"}, output: "git-automation-state-id\n"},
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
		{"--quiet", "issue", "unrelate", "issue-relation-id"},
		{"--quiet", "issue", "link", "https://example.com/pr/1", "LIT-1"},
		{"--quiet", "agent-session", "get", "agent-session-id"},
		{"--quiet", "team", "git-automation-states", "team-id"},
		{"--quiet", "customer-need", "project-attachment", "customer-need-id"},
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
