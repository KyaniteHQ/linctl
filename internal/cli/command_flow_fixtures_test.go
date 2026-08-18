package cli

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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
			Target: config.Target{
				OrgID:     "org-id",
				TeamKey:   "LIT",
				TeamID:    "team-id",
				ProjectID: "project-id",
			},
		},
		repoConfig:    repoConfigSelection{Path: "/repo/.linctl.toml", Status: repoConfigLoaded},
		graphqlClient: graphqlClient,
	}
}

type commandFlowFakeClient struct {
	emptyIssueList                bool
	emptyIssueChildren            bool
	emptyIssueComments            bool
	truncatedExport               bool
	invalidExportLeaf             bool
	emptyIssueProject             bool
	emptyIssueMine                bool
	emptyIssueLabel               bool
	orgWideLabel                  bool
	otherOrgProjectLabel          bool
	otherOrgInitiativeLabel       bool
	emptyIssueCycle               bool
	emptyIssueCreatedAfter        bool
	emptyIssueCreatedBefore       bool
	emptyIssueUpdatedAfter        bool
	emptyIssueUpdatedBefore       bool
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
	expectedUpdatedAfter          string
	expectedUpdatedBefore         string
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
	// "team" is the direct-lookup fast path for the same team resolution that
	// "Teams" performs by scanning, so a test forcing "Teams" to fail must also
	// fail the fast path, or resolution would silently succeed through it.
	if client.failOperation == "Teams" && request.OpName == "team" {
		return errors.New("operation failed")
	}
	if err := client.requireExpectedVariables(request); err != nil {
		return err
	}

	operation := request.OpName
	if operation == "IssuesByTeamFiltered" {
		operation = filteredIssueListPayloadKey(request)
	}
	payload, err := commandFlowPayload(operation, client)
	if err != nil {
		return err
	}
	if crossTeamPayload, ok := commandFlowCrossTeamPayload(request); ok {
		payload = crossTeamPayload
	}

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

// commandFlowCrossTeamPayload returns destination-aware fixtures for the
// project add-team, issue move-team and issue move-project happy paths. The
// default team/issue/project fixtures stay on the pinned LIT team and project,
// so a move that lands elsewhere needs a response that names that destination
// or the post-write check fails closed.
func commandFlowCrossTeamPayload(request *graphql.Request) (string, bool) {
	switch request.OpName {
	case "team":
		id, err := requestVariable[string](request, "id")
		if err != nil || id != "ops-team-id" {
			return "", false
		}

		return `{"team":` + commandDestinationTeamJSON("ops-team-id", "OPS") + `}`, true
	case "project":
		id, err := requestVariable[string](request, "id")
		if err != nil || id != commandDestinationProjectID {
			return "", false
		}

		return `{"project":` + commandDestinationProjectJSON() + `}`, true
	case "IssueUpdate":
		if projectID, err := requestVariable[string](request, "input", "projectId"); err == nil {
			if projectID != commandDestinationProjectID {
				return "", false
			}

			return `{"issueUpdate":{"success":true,"issue":` + commandIssueJSONWithProject(
				"LIT-1", "Moved issue", commandDestinationProjectID, "EOIR Case Scraper",
			) + `}}`, true
		}
		teamID, err := requestVariable[string](request, "input", "teamId")
		if err != nil || teamID != "ops-team-id" {
			return "", false
		}

		return `{"issueUpdate":{"success":true,"issue":` +
			commandIssueJSONWithTeam("LIT-1", "Moved issue", "todo-state", "Todo", "unstarted", "ops-team-id", "OPS") +
			`}}`, true
	case "ProjectUpdate":
		if !requestHasTeamID(request, "ops-team-id") {
			return "", false
		}

		return `{"projectUpdate":{"success":true,"project":` + commandProjectJSONWithTeams(
			"Multi-team project", "Started", "started",
			[]commandProjectTeam{
				{ID: "team-id", Key: "LIT", Name: "linctl"},
				{ID: "ops-team-id", Key: "OPS", Name: "OPS"},
			},
		) + `}}`, true
	default:
		return "", false
	}
}

func requestHasTeamID(request *graphql.Request, teamID string) bool {
	payload, err := json.Marshal(request.Variables)
	if err != nil {
		return false
	}
	var variables struct {
		Input struct {
			TeamIDs []string `json:"teamIds"`
		} `json:"input"`
	}
	if err := json.Unmarshal(payload, &variables); err != nil {
		return false
	}
	for _, id := range variables.Input.TeamIDs {
		if id == teamID {
			return true
		}
	}

	return false
}

// filteredIssueListPayloadKey maps a composed IssuesByTeamFiltered request to a
// per-clause payload key, so the canned fixtures stay distinct per filter now
// that one operation serves every team-scoped issue listing. The key is built
// from the clauses actually present in the request, so combined or unmapped
// clauses surface as a missing-fixture failure instead of the wrong payload.
func filteredIssueListPayloadKey(request *graphql.Request) string {
	payload, err := json.Marshal(request.Variables)
	if err != nil {
		return "IssuesByTeamFiltered:unreadable-variables"
	}
	var variables struct {
		Filter map[string]json.RawMessage `json:"filter"`
	}
	if err := json.Unmarshal(payload, &variables); err != nil {
		return "IssuesByTeamFiltered:unreadable-variables"
	}

	clauses := make([]string, 0, len(variables.Filter))
	for clause, value := range variables.Filter {
		if clause == "team" {
			continue
		}
		if clause == "createdAt" {
			clause = dateComparatorClauseKey(value, "createdAfter", "createdBefore")
		}
		if clause == "updatedAt" {
			clause = dateComparatorClauseKey(value, "updatedAfter", "updatedBefore")
		}
		clauses = append(clauses, clause)
	}
	if len(clauses) == 0 {
		return "IssuesByTeamFiltered:unfiltered"
	}
	sort.Strings(clauses)

	return "IssuesByTeamFiltered:" + strings.Join(clauses, "+")
}

func dateComparatorClauseKey(window json.RawMessage, afterKey string, beforeKey string) string {
	var comparator struct {
		Gte *string `json:"gte"`
	}
	if json.Unmarshal(window, &comparator) == nil && comparator.Gte != nil {
		return afterKey
	}

	return beforeKey
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
	if request.OpName != "IssuesByTeamFiltered" {
		if client.expectedBlockedBy != "" && request.OpName == "IssueBlockedIssues" {
			return requireRequestVariable(request, []string{"id"}, client.expectedBlockedBy, "blocked by issue")
		}

		return nil
	}
	checks := []struct {
		expected string
		keys     []string
		label    string
	}{
		{client.expectedStateType, []string{"filter", "state", "type", "eq"}, "state type"},
		{client.expectedProjectID, []string{"filter", "project", "id", "eq"}, "project id"},
		{client.expectedAssigneeID, []string{"filter", "assignee", "id", "eq"}, "assignee id"},
		{client.expectedLabelID, []string{"filter", "labels", "some", "id", "eq"}, "label id"},
		{client.expectedCycleID, []string{"filter", "cycle", "id", "eq"}, "cycle id"},
		{client.expectedCreatedAfter, []string{"filter", "createdAt", "gte"}, "created after"},
		{client.expectedCreatedBefore, []string{"filter", "createdAt", "lte"}, "created before"},
		{client.expectedUpdatedAfter, []string{"filter", "updatedAt", "gte"}, "updated after"},
		{client.expectedUpdatedBefore, []string{"filter", "updatedAt", "lte"}, "updated before"},
	}
	for _, check := range checks {
		if check.expected == "" {
			continue
		}
		if err := requireRequestVariable(request, check.keys, check.expected, check.label); err != nil {
			return err
		}
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
