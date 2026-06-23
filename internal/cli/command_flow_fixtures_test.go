package cli

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
