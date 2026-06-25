package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		{"--json", "--fields", "id,slug_id,status", "agent-session", "list", "--limit", "1"},
		{"--json", "agent-session", "get", "agent-session-id"},
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
		{"--json", "--fields", "id,display_name", "document", "comments", "document-id", "--limit", "1"},
		{"--json", "issue", "pr", "LIT-1"},
		{"--json", "issue", "start", "LIT-1"},
		{"--json", "issue", "comment", "LIT-1", "--body", "Looks good"},
		{"--json", "issue", "reply", "LIT-1", "comment-id", "--body", "Reply body"},
		{"--json", "issue", "relate", "LIT-1", "LIT-2", "--type", "related"},
		{"--json", "issue", "unrelate", "issue-relation-id"},
		{"--json", "issue", "link", "https://example.com/pr/1", "LIT-1"},
		{"--json", "--fields", "id,display_name", "issue", "comments", "LIT-1", "--limit", "1"},
		{"--json", "--fields", "identifier,title", "issue", "figma-file-key-search", "figma-key", "--limit", "1"},
		{"--json", "issue", "priority-values"},
		{"--json", "issue", "filter-suggestion", "started issues"},
		{"--json", "issue", "title-suggestion", "Customer asks for faster exports"},
		{"--json", "--fields", "id,display_name", "comment", "list", "--limit", "1"},
		{"--json", "comment", "get", "comment-id"},
		{"--json", "comment", "update", "comment-id", "--body", "New body"},
		{"--json", "comment", "delete", "comment-id"},
		{"--json", "--fields", "id,display_name", "initiative-update", "comments", "initiative-update-id", "--limit", "1"},
		{"--json", "project", "list", "--limit", "1"},
		{"--json", "project", "all", "--limit", "1"},
		{"--json", "project", "members", "project-id", "--limit", "1"},
		{"--json", "--fields", "id,health,display_name", "project", "updates", "project-id", "--limit", "1"},
		{"--json", "project", "filter-suggestion", "started projects"},
		{"--json", "--fields", "id,health,project_id", "project-update", "list", "--limit", "1"},
		{"--json", "project-update", "get", "project-update-id"},
		{"--json", "project-update", "create", "project-id", "--health", "on-track", "--body", "Posted update"},
		{"--json", "--fields", "id,name,status", "project-milestone", "all", "--limit", "1"},
		{"--json", "--fields", "id,name,status", "project-milestone", "list", "project-id", "--limit", "1"},
		{"--json", "project-status", "list", "--limit", "1"},
		{"--json", "project-status", "get", "project-status-id"},
		{"--json", "project-label", "list", "--limit", "1"},
		{"--json", "project-label", "children", "project-label-id", "--limit", "1"},
		{"--json", "project-label", "projects", "project-label-id", "--limit", "1"},
		{"--json", "--fields", "id,type,project_id,related_project_id", "project-relation", "list", "--limit", "1"},
		{"--json", "project-relation", "get", "project-relation-id"},
		{"--json", "--fields", "id,title,parent_type", "document", "list", "--limit", "1"},
		{"--json", "document", "create", "--title", "Created doc"},
		{"--json", "document", "update", "document-id", "--title", "Updated doc"},
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
		{"--json", "team", "git-automation-states", "team-id", "--limit", "1"},
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
		{"--json", "--fields", "id,name,status", "roadmap", "projects", "roadmap-id", "--limit", "1"},
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
		{"--json", "customer-need", "project-attachment", "customer-need-id"},
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
