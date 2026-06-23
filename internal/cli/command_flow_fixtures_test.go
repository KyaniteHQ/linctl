package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

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

	return useCommandRuntimeWithFiles(t, graphqlClient, http.DefaultClient)
}

func useCommandRuntimeWithFiles(t *testing.T, graphqlClient graphql.Client, fileClient httpDoer) func() {
	t.Helper()

	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		runtime := testCommandRuntime(graphqlClient)
		runtime.fileClient = fileClient
		return runtime, nil
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
	emptyIssueList                bool
	emptyIssueChildren            bool
	emptyIssueComments            bool
	truncatedExport               bool
	emptyIssueProject             bool
	emptyIssueMine                bool
	emptyIssueLabel               bool
	emptyIssueCycle               bool
	emptyIssueCreatedAfter        bool
	emptyIssueCreatedBefore       bool
	emptyIssueHasBlockers         bool
	emptyIssueBlocks              bool
	emptyIssueBlockedBy           bool
	emptyIssueAllTeams            bool
	emptyIssueSearch              bool
	emptyIssueFigmaSearch         bool
	emptyNextIssues               bool
	rankedNextIssues              bool
	expectedStateType             string
	expectedProjectID             string
	expectedAssigneeID            string
	expectedLabelID               string
	expectedCycleID               string
	expectedCreatedAfter          string
	expectedCreatedBefore         string
	expectedBlockedBy             string
	expectedIssueDeps             string
	expectedSearchQuery           string
	expectedIssueFigmaFileKey     string
	expectedIssueFilterPrompt     string
	expectedIssueFilterTeamID     string
	expectedIssueTitleRequest     string
	expectedReleaseSearchTerm     string
	expectedSemanticSearchQuery   string
	expectedTypedSearchTerm       string
	emptyReleaseSearch            bool
	emptyProjectList              bool
	emptyProjectMembers           bool
	emptyProjectUpdates           bool
	emptyProjectMilestones        bool
	emptySLAConfigurations        bool
	emptySemanticSearch           bool
	emptySearchDocuments          bool
	emptySearchIssues             bool
	emptySearchProjects           bool
	emptyViewerDrafts             bool
	expectedCommentBody           string
	expectedCommentParentID       string
	expectedCreateDescription     string
	expectedUpdateDescription     string
	expectedCreateTitle           string
	expectedUpdateTitle           string
	expectedProjectCreateName     string
	expectedProjectUpdateName     string
	expectedMilestoneCreateName   string
	expectedMilestoneUpdateName   string
	expectedStartAssigneeID       string
	expectedStartStateID          string
	expectedOrganizationURLKey    string
	expectedApplicationClientID   string
	missingCustomerNeedAttachment bool
	failOperation                 string
	multiIssueList                bool
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
	if err := client.requireExpectedWriteVariables(request); err != nil {
		return err
	}
	return client.requireExpectedIssueStartVariables(request)
}

// requireExpectedWriteVariables asserts that guarded-write commands forward the
// user-supplied title/name flag into the GraphQL input, so a silently dropped
// flag value fails the test instead of passing on the output substring alone.
func (client commandFlowFakeClient) requireExpectedWriteVariables(request *graphql.Request) error {
	if client.expectedCreateTitle != "" && request.OpName == "IssueCreate" {
		return requireRequestVariable(request, []string{"input", "title"}, client.expectedCreateTitle, "create title")
	}
	if client.expectedUpdateTitle != "" && request.OpName == "IssueUpdate" {
		return requireRequestVariable(request, []string{"input", "title"}, client.expectedUpdateTitle, "update title")
	}
	if client.expectedProjectCreateName != "" && request.OpName == "ProjectCreate" {
		return requireRequestVariable(request, []string{"input", "name"}, client.expectedProjectCreateName, "project create name")
	}
	if client.expectedProjectUpdateName != "" && request.OpName == "ProjectUpdate" {
		return requireRequestVariable(request, []string{"input", "name"}, client.expectedProjectUpdateName, "project update name")
	}
	if client.expectedMilestoneCreateName != "" && request.OpName == "ProjectMilestoneCreate" {
		return requireRequestVariable(request, []string{"input", "name"}, client.expectedMilestoneCreateName, "milestone create name")
	}
	if client.expectedMilestoneUpdateName != "" && request.OpName == "ProjectMilestoneUpdate" {
		return requireRequestVariable(request, []string{"input", "name"}, client.expectedMilestoneUpdateName, "milestone update name")
	}

	return nil
}

func (client commandFlowFakeClient) requireExpectedSearchVariables(request *graphql.Request) error {
	if client.expectedSearchQuery != "" && request.OpName == "issueSearch" {
		return requireRequestVariable(request, []string{"query"}, client.expectedSearchQuery, "search query")
	}
	if err := client.requireExpectedIssueUtilityVariables(request); err != nil {
		return err
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

func (client commandFlowFakeClient) requireExpectedIssueUtilityVariables(request *graphql.Request) error {
	if client.expectedIssueFigmaFileKey != "" && request.OpName == "issueFigmaFileKeySearch" {
		return requireRequestVariable(request, []string{"fileKey"}, client.expectedIssueFigmaFileKey, "figma file key")
	}
	if client.expectedIssueFilterPrompt != "" && request.OpName == "issueFilterSuggestion" {
		if err := requireRequestVariable(
			request,
			[]string{"prompt"},
			client.expectedIssueFilterPrompt,
			"issue filter prompt",
		); err != nil {
			return err
		}
	}
	if client.expectedIssueFilterTeamID != "" && request.OpName == "issueFilterSuggestion" {
		return requireRequestVariable(request, []string{"teamId"}, client.expectedIssueFilterTeamID, "issue filter team id")
	}
	if client.expectedIssueTitleRequest != "" && request.OpName == "issueTitleSuggestionFromCustomerRequest" {
		return requireRequestVariable(
			request,
			[]string{"request"},
			client.expectedIssueTitleRequest,
			"issue title request",
		)
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
	if payload, ok := commandFlowCommentPayload(operation); ok {
		return payload, nil
	}
	if payload, ok := commandFlowFilePayload(operation); ok {
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
	case "team_gitAutomationStates":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","gitAutomationStates":{"nodes":[` +
			commandGitAutomationStateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "team_templates":
		return `{"team":{"id":"team-id","key":"LIT","name":"linctl","templates":{"nodes":[` +
			commandTemplateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

//nolint:gocyclo // The command-flow fake is intentionally centralized by operation name.
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
			`"project":{"id":"project-id","name":"Pinned project"},`+
				`"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null`,
		) + `}`, true
	case "DocumentCreate":
		return `{"documentCreate":{"success":true,"document":` + commandDocumentJSON(
			"Created doc",
			`"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null`,
		) + `}}`, true
	case "DocumentUpdate":
		return `{"documentUpdate":{"success":true,"document":` + commandDocumentJSON(
			"Updated doc",
			`"project":null,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issue":null,"cycle":null`,
		) + `}}`, true
	case "ProjectUpdateCreate":
		return `{"projectUpdateCreate":{"success":true,"projectUpdate":` + commandProjectUpdateJSON() + `}}`, true
	case "document_comments":
		return `{"document":{"id":"document-id","comments":{"nodes":[` +
			commandCommentMetadataJSON("", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
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
	if payload, ok := commandFlowUserSettingsPayload(operation); ok {
		return payload, true
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

//nolint:gocyclo // Each branch mirrors a distinct official UserSettings operation.
func commandFlowUserSettingsPayload(operation string) (string, bool) {
	switch operation {
	case "userSettings":
		return `{"userSettings":` + commandUserSettingsJSON() + `}`, true
	case "userSettings_notificationCategoryPreferences":
		return `{"userSettings":{"notificationCategoryPreferences":` + commandNotificationCategoriesJSON() + `}}`, true
	case "userSettings_notificationCategoryPreferences_appsAndIntegrations":
		return commandUserSettingsCategoryPayload("appsAndIntegrations"), true
	case "userSettings_notificationCategoryPreferences_assignments":
		return commandUserSettingsCategoryPayload("assignments"), true
	case "userSettings_notificationCategoryPreferences_billing":
		return commandUserSettingsCategoryPayload("billing"), true
	case "userSettings_notificationCategoryPreferences_commentsAndReplies":
		return commandUserSettingsCategoryPayload("commentsAndReplies"), true
	case "userSettings_notificationCategoryPreferences_customers":
		return commandUserSettingsCategoryPayload("customers"), true
	case "userSettings_notificationCategoryPreferences_documentChanges":
		return commandUserSettingsCategoryPayload("documentChanges"), true
	case "userSettings_notificationCategoryPreferences_feed":
		return commandUserSettingsCategoryPayload("feed"), true
	case "userSettings_notificationCategoryPreferences_mentions":
		return commandUserSettingsCategoryPayload("mentions"), true
	case "userSettings_notificationCategoryPreferences_postsAndUpdates":
		return commandUserSettingsCategoryPayload("postsAndUpdates"), true
	case "userSettings_notificationCategoryPreferences_reactions":
		return commandUserSettingsCategoryPayload("reactions"), true
	case "userSettings_notificationCategoryPreferences_reminders":
		return commandUserSettingsCategoryPayload("reminders"), true
	case "userSettings_notificationCategoryPreferences_reviews":
		return commandUserSettingsCategoryPayload("reviews"), true
	case "userSettings_notificationCategoryPreferences_statusChanges":
		return commandUserSettingsCategoryPayload("statusChanges"), true
	case "userSettings_notificationCategoryPreferences_subscriptions":
		return commandUserSettingsCategoryPayload("subscriptions"), true
	case "userSettings_notificationCategoryPreferences_system":
		return commandUserSettingsCategoryPayload("system"), true
	case "userSettings_notificationCategoryPreferences_triage":
		return commandUserSettingsCategoryPayload("triage"), true
	case "userSettings_notificationChannelPreferences":
		return `{"userSettings":{"notificationChannelPreferences":` + commandNotificationChannelJSON() + `}}`, true
	case "userSettings_notificationDeliveryPreferences":
		return `{"userSettings":{"notificationDeliveryPreferences":` + commandNotificationDeliveryPreferencesJSON() + `}}`, true
	case "userSettings_notificationDeliveryPreferences_mobile":
		return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":` + commandNotificationDeliveryChannelJSON() + `}}}`, true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule":
		return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":{"schedule":` + commandNotificationDeliveryScheduleJSON() + `}}}}`, true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_friday":
		return commandUserSettingsScheduleDayPayload("friday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_monday":
		return commandUserSettingsScheduleDayPayload("monday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_saturday":
		return commandUserSettingsScheduleDayPayload("saturday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_sunday":
		return commandUserSettingsScheduleDayPayload("sunday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_thursday":
		return commandUserSettingsScheduleDayPayload("thursday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_tuesday":
		return commandUserSettingsScheduleDayPayload("tuesday"), true
	case "userSettings_notificationDeliveryPreferences_mobile_schedule_wednesday":
		return commandUserSettingsScheduleDayPayload("wednesday"), true
	case "userSettings_theme":
		return `{"userSettings":{"theme":` + commandUserSettingsThemeJSON(true) + `}}`, true
	case "userSettings_theme_custom":
		return `{"userSettings":{"theme":{"custom":` + commandUserSettingsCustomThemeJSON(true) + `}}}`, true
	case "userSettings_theme_custom_sidebar":
		return `{"userSettings":{"theme":{"custom":{"sidebar":` + commandUserSettingsCustomSidebarThemeJSON() + `}}}}`, true
	default:
		return "", false
	}
}

func commandUserSettingsCategoryPayload(category string) string {
	return `{"userSettings":{"notificationCategoryPreferences":{"` + category + `":` + commandNotificationChannelJSON() + `}}}`
}

func commandUserSettingsScheduleDayPayload(day string) string {
	return `{"userSettings":{"notificationDeliveryPreferences":{"mobile":{"schedule":{"` + day + `":` +
		commandNotificationDeliveryDayJSON() + `}}}}}`
}

func commandFlowStateAndCommentPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "workflowStates":
		return `{"workflowStates":{"nodes":[` + commandWorkflowStateJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "workflowState":
		return `{"workflowState":` + commandWorkflowStateJSON() + `}`, true
	case "workflowState_issues":
		return `{"workflowState":{"id":"workflow-state-id","name":"Started","issues":{"nodes":[` +
			commandIssueJSON("LIT-1", "Detail issue", "state-id", "Todo", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	}

	return commandFlowInitiativePayload(operation, fake)
}

func commandFlowFilePayload(operation string) (string, bool) {
	if operation != "fileUpload" {
		return "", false
	}

	return `{"fileUpload":{"success":true,"uploadFile":{` +
		`"filename":"upload.txt","contentType":"text/plain","size":11,` +
		`"uploadUrl":"https://uploads.example/put","assetUrl":"https://assets.example/file.txt",` +
		`"headers":[{"key":"x-test","value":"1"}]}}}`, true
}

func commandFlowCommentPayload(operation string) (string, bool) {
	switch operation {
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
	case "CommentUpdate":
		return `{"commentUpdate":{"success":true,"comment":` + commandTopLevelCommentJSON() + `}}`, true
	case "CommentDelete":
		return `{"commentDelete":{"success":true,"entityId":"comment-id"}}`, true
	default:
		return "", false
	}
}

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
	}

	return commandFlowInitiativeUpdatePayload(operation, fake)
}

func commandFlowInitiativeUpdatePayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
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
	case "initiativeUpdate_comments":
		return `{"initiativeUpdate":{"id":"initiative-update-id","comments":{"nodes":[` +
			commandCommentMetadataJSON("", "") +
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
	case "templateContent":
		return `{"template":{"id":"template-id","name":"Bug report","templateData":` +
			`{"title":"Template title","description":"## Steps\n\nReproduce here"}}}`, true
	case "roadmaps":
		return `{"roadmaps":{"nodes":[` + commandRoadmapJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "roadmap":
		return `{"roadmap":` + commandRoadmapJSON() + `}`, true
	case "roadmap_projects":
		return `{"roadmap":{"id":"roadmap-id","name":"Platform roadmap","projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
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
	case "customerNeed_projectAttachment":
		if fake.missingCustomerNeedAttachment {
			return `{"customerNeed":{"id":"customer-need-id","projectAttachment":null}}`, true
		}
		return `{"customerNeed":{"id":"customer-need-id","projectAttachment":` + commandAttachmentJSON() + `}}`, true
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

//nolint:gocyclo // The command-flow fake is intentionally centralized by operation name.
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
	case "attachmentIssue_comments":
		return `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[` +
			commandCommentMetadataJSON("issue-id", "") +
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
	case "attachmentIssue_formerNeeds":
		return `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","formerNeeds":{"nodes":[` +
			commandCustomerNeedJSON() +
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
	case "attachmentIssue_needs":
		return `{"attachmentIssue":{"id":"issue-id","identifier":"LIT-1","needs":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_relations":
		return `{"attachmentIssue":{"relations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_releases":
		return `{"attachmentIssue":{"releases":{"nodes":[` +
			commandReleaseJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "attachmentIssue_sharedAccess":
		return commandIssueSharedAccessPayload("attachmentIssue"), true
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

//nolint:gocyclo // The command-flow fake is intentionally centralized by operation name.
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
	case "issueVcsBranchSearch_comments":
		return `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[` +
			commandCommentMetadataJSON("issue-id", "") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_formerNeeds":
		return `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","formerNeeds":{"nodes":[` +
			commandCustomerNeedJSON() +
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
	case "issueVcsBranchSearch_needs":
		return `{"issueVcsBranchSearch":{"id":"issue-id","identifier":"LIT-1","needs":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_relations":
		return `{"issueVcsBranchSearch":{"relations":{"nodes":[` +
			commandIssueRelationJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_releases":
		return `{"issueVcsBranchSearch":{"releases":{"nodes":[` +
			commandReleaseJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issueVcsBranchSearch_sharedAccess":
		return commandIssueSharedAccessPayload("issueVcsBranchSearch"), true
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

func commandFlowIssueRelationPayload(operation string) (string, bool) {
	switch operation {
	case "issueRelations":
		return `{"issueRelations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issueRelation":
		return `{"issueRelation":` + commandIssueRelationJSON() + `}`, true
	case "IssueRelationCreate":
		return `{"issueRelationCreate":{"success":true,"issueRelation":{` +
			`"id":"issue-relation-id","type":"related",` +
			`"createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","archivedAt":null,` +
			`"issue":{"id":"issue-id","identifier":"LIT-1","title":"Source issue"},` +
			`"relatedIssue":{"id":"related-issue-id","identifier":"LIT-2","title":"Related issue"}}}}`, true
	case "IssueRelationDelete":
		return `{"issueRelationDelete":{"success":true,"entityId":"issue-relation-id"}}`, true
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
	if payload, ok := commandFlowIssueUtilityPayload(operation, fake); ok {
		return payload, true
	}
	if payload, ok := commandFlowIssueRelationPayload(operation); ok {
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
	case "issue_comments":
		if fake.emptyIssueComments {
			return `{"issue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		hasNextPage := "false"
		if fake.truncatedExport {
			hasNextPage = "true"
		}
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","comments":{"nodes":[{"id":"comment-id","body":"First comment","url":"https://linear.app/comment/comment-id","createdAt":"2026-06-19T12:00:00Z","parentId":null,"user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":` + hasNextPage + `,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func commandFlowIssueUtilityPayload(operation string, fake commandFlowFakeClient) (string, bool) {
	switch operation {
	case "issueFigmaFileKeySearch":
		if fake.emptyIssueFigmaSearch {
			return `{"issueFigmaFileKeySearch":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"issueFigmaFileKeySearch":{"nodes":[` +
			commandIssueJSON("LIT-41", "Figma issue", "todo-state", "Todo", "unstarted") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "issuePriorityValues":
		return `{"issuePriorityValues":[{"priority":1,"label":"Urgent"},{"priority":0,"label":"No priority"}]}`, true
	case "issueFilterSuggestion":
		return `{"issueFilterSuggestion":{"filter":{"state":{"type":{"eq":"started"}}},"logId":"issue-filter-log-id"}}`, true
	case "issueTitleSuggestionFromCustomerRequest":
		return `{"issueTitleSuggestionFromCustomerRequest":{"title":"Improve exports","logId":"title-log-id"}}`, true
	default:
		return "", false
	}
}

//nolint:gocyclo // The command-flow fake is intentionally centralized by operation name.
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
	case "issue_formerNeeds":
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","formerNeeds":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_history":
		return `{"issue":{"history":{"nodes":[` + commandIssueHistoryJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_inverseRelations":
		return `{"issue":{"inverseRelations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_labels":
		return `{"issue":{"labels":{"nodes":[` + commandLabelJSON("label body") + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_needs":
		return `{"issue":{"id":"issue-id","identifier":"LIT-1","needs":{"nodes":[` +
			commandCustomerNeedJSON() +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_relations":
		return `{"issue":{"relations":{"nodes":[` + commandIssueRelationJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_releases":
		return `{"issue":{"releases":{"nodes":[` + commandReleaseJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "issue_sharedAccess":
		return commandIssueSharedAccessPayload("issue"), true
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

func commandIssueSharedAccessPayload(root string) string {
	return `{"` + root + `":{"id":"issue-id","identifier":"LIT-1","sharedAccess":` +
		commandIssueSharedAccessJSON() + `}}`
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
	case "WorkflowStatesByType":
		return `{"workflowStates":{"nodes":[{"id":"type-state-id","name":"TypeState","type":"unstarted","position":1}]}}`, true
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
	case "projects":
		if fake.emptyProjectList {
			return `{"projects":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"projects":{"nodes":[` +
			commandProjectJSON("Listed project", "Backlog", "backlog") +
			`],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
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
	case "project_projectUpdates":
		if fake.emptyProjectUpdates {
			return `{"project":{"id":"project-id","name":"Detail project","projectUpdates":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"project":{"id":"project-id","name":"Detail project","projectUpdates":{"nodes":[{"id":"project-update-id","health":"onTrack","createdAt":"2026-06-19T12:00:00Z","updatedAt":"2026-06-19T12:00:00Z","url":"https://linear.app/project-update/project-update-id","user":{"id":"user-id","name":"omer","displayName":"Omer"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "projectFilterSuggestion":
		return `{"projectFilterSuggestion":{"filter":{"status":{"type":{"eq":"started"}}},"logId":"filter-log-id"}}`, true
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
	case "projectMilestones":
		if fake.emptyProjectMilestones {
			return `{"projectMilestones":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"projectMilestones":{"nodes":[{"id":"project-milestone-id","name":"Launch milestone","description":"milestone body","targetDate":"2026-06-30","status":"next","progress":0.5,"sortOrder":1}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
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
	case "projectStatusProjectCount":
		return `{"projectStatusProjectCount":{"count":12,"privateCount":2,"archivedTeamCount":1}}`, true
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

func commandIssueSharedAccessJSON() string {
	return `{
		"isShared":true,
		"viewerHasOnlySharedAccess":false,
		"sharedWithCount":2,
		"disallowedIssueFields":["description","priority"]
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

func commandGitAutomationStateJSON() string {
	return `{
		"id":"git-automation-state-id",
		"event":"review",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"state":{"id":"workflow-state-id","name":"Started","type":"started"},
		"targetBranch":{"id":"target-branch-id","branchPattern":"main","isRegex":false}
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

func commandUserSettingsJSON() string {
	return `{
		"id":"settings-id",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:01:00Z",
		"archivedAt":null,
		"autoAssignToSelf":true,
		"feedLastSeenTime":null,
		"feedSummarySchedule":"daily",
		"showFullUserNames":false,
		"subscribedToChangelog":true,
		"subscribedToDPA":false,
		"subscribedToInviteAccepted":true,
		"subscribedToPrivacyLegalUpdates":true,
		"user":{"id":"user-id"},
		"notificationCategoryPreferences":` + commandNotificationCategoriesJSON() + `,
		"notificationChannelPreferences":` + commandNotificationChannelJSON() + `,
		"notificationDeliveryPreferences":` + commandNotificationDeliveryPreferencesJSON() + `
	}`
}

func commandNotificationCategoriesJSON() string {
	channel := commandNotificationChannelJSON()
	return `{
		"appsAndIntegrations":` + channel + `,
		"assignments":` + channel + `,
		"billing":` + channel + `,
		"commentsAndReplies":` + channel + `,
		"customers":` + channel + `,
		"documentChanges":` + channel + `,
		"feed":` + channel + `,
		"mentions":` + channel + `,
		"postsAndUpdates":` + channel + `,
		"reactions":` + channel + `,
		"reminders":` + channel + `,
		"reviews":` + channel + `,
		"statusChanges":` + channel + `,
		"subscriptions":` + channel + `,
		"system":` + channel + `,
		"triage":` + channel + `
	}`
}

func commandNotificationChannelJSON() string {
	return `{"desktop":true,"email":false,"mobile":true,"slack":false}`
}

func commandNotificationDeliveryPreferencesJSON() string {
	return `{"mobile":` + commandNotificationDeliveryChannelJSON() + `}`
}

func commandNotificationDeliveryChannelJSON() string {
	return `{"notificationsDisabled":false,"schedule":` + commandNotificationDeliveryScheduleJSON() + `}`
}

func commandNotificationDeliveryScheduleJSON() string {
	day := commandNotificationDeliveryDayJSON()
	return `{
		"disabled":false,
		"friday":` + day + `,
		"monday":` + day + `,
		"saturday":` + day + `,
		"sunday":` + day + `,
		"thursday":` + day + `,
		"tuesday":` + day + `,
		"wednesday":` + day + `
	}`
}

func commandNotificationDeliveryDayJSON() string {
	return `{"start":"09:00","end":"18:00"}`
}

func commandUserSettingsThemeJSON(includeCustom bool) string {
	custom := "null"
	if includeCustom {
		custom = commandUserSettingsCustomThemeJSON(true)
	}

	return `{"preset":"custom","custom":` + custom + `}`
}

func commandUserSettingsCustomThemeJSON(includeSidebar bool) string {
	sidebar := "null"
	if includeSidebar {
		sidebar = commandUserSettingsCustomSidebarThemeJSON()
	}

	return `{"accent":[50.5,20.5,10.5],"base":[90.5,0,0],"contrast":50,"sidebar":` + sidebar + `}`
}

func commandUserSettingsCustomSidebarThemeJSON() string {
	return `{"accent":[60.5,20.5,10.5],"base":[20.5,0,0],"contrast":70}`
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
