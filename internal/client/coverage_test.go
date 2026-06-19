package client

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func Test_ClientReadScenarios_return_compact_lists_details_and_members(t *testing.T) {
	// Given
	endCursor := "cursor-1"
	graphqlClient := fakeGraphQLClient{
		"IssuesByTeam": `{"issues":{"nodes":[` + issueJSON(issueFixture{
			Identifier: "LIT-10",
			Title:      "listed issue",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}`,
		"IssueByID": `{"issue":` + issueJSON(issueFixture{
			Identifier: "LIT-11",
			Title:      "detail issue",
			StateID:    "done",
			State:      "Done",
			StateType:  "completed",
		}) + `}`,
		"Projects": `{"team":{"projects":{"nodes":[` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "listed",
			Status: "Backlog",
		}) + `],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
		"ProjectByID": `{"project":` + projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "detail",
			Status: "Started",
		}) + `}`,
		"ProjectMembers": `{"project":{"id":"project-id","name":"detail","members":{"nodes":[{"id":"user-id","name":"omer","displayName":"Omer","email":"omer@example.com"}],"pageInfo":{"hasNextPage":true,"endCursor":"` + endCursor + `"}}}}`,
	}

	// When
	issues, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	issue, err := GetIssueByID(context.Background(), graphqlClient, "LIT-11")
	require.NoError(t, err)
	projects, err := ListProjectsByTeam(context.Background(), graphqlClient, "team-id", 2)
	require.NoError(t, err)
	project, err := GetProjectByID(context.Background(), graphqlClient, "project-id")
	require.NoError(t, err)
	members, err := ListProjectMembers(context.Background(), graphqlClient, "project-id", 2)
	require.NoError(t, err)

	// Then
	require.True(t, issues.HasNextPage)
	require.Equal(t, "LIT-10", issues.Issues[0].Identifier)
	require.Equal(t, "LIT-11", issue.Identifier)
	require.True(t, projects.HasNextPage)
	require.Equal(t, "listed", projects.Projects[0].Name)
	require.Equal(t, "detail", project.Name)
	require.Equal(t, "Omer", members.Members[0].DisplayName)
	require.Equal(t, &endCursor, members.EndCursor)
}

func Test_ClientWriteScenarios_guard_writes_and_report_results(t *testing.T) {
	// Given
	t.Run("invalid requests fail before network", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{})

		_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{ID: "LIT-1"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{Title: "missing id"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{ID: "LIT-1"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{Body: "body"})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = CloseIssue(context.Background(), graphqlClient, matchingTarget(), "")
		require.Error(t, err)

		_, err = CreateProject(context.Background(), graphqlClient, matchingTarget(), ProjectCreateRequest{})
		require.ErrorIs(t, err, ErrWriteInvalid)

		_, err = UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{ID: "project-id"})
		require.ErrorIs(t, err, ErrWriteInvalid)
	})

	t.Run("issue comment succeeds", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-12",
				Title:      "comment target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueCommentCreate": `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"hello","url":"https://linear.app/comment/comment-id","issue":` + issueJSON(issueFixture{
				Identifier: "LIT-12",
				Title:      "comment target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}}}`,
		})

		comment, err := CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{
			ID:   "LIT-12",
			Body: "hello",
		})

		require.NoError(t, err)
		require.Equal(t, "comment-id", comment.ID)
		require.Equal(t, "LIT-12", comment.Issue.Identifier)
	})

	t.Run("issue update succeeds", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-21",
				Title:      "update target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueUpdate": `{"issueUpdate":{"success":true,"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-21",
				Title:      "updated",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}}`,
		})

		issue, err := UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{
			ID:    "LIT-21",
			Title: "updated",
		})

		require.NoError(t, err)
		require.Equal(t, "updated", issue.Title)
	})

	t.Run("project update and archive succeed", func(t *testing.T) {
		graphqlClient := projectWriteFakeClient(map[string]string{
			"ProjectByID": `{"project":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "fixture",
				Status: "Backlog",
			}) + `}`,
			"ProjectUpdate": `{"projectUpdate":{"success":true,"project":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "updated",
				Status: "Started",
			}) + `}}`,
			"ProjectArchive": `{"projectArchive":{"success":true,"entity":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "updated",
				Status: "Canceled",
			}) + `}}`,
		})

		project, err := UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{
			ID:   "project-id",
			Name: "updated",
		})
		require.NoError(t, err)
		require.Equal(t, "updated", project.Name)

		project, err = ArchiveProject(context.Background(), graphqlClient, matchingTarget(), "project-id")
		require.NoError(t, err)
		require.Equal(t, "Canceled", project.Status.Name)
	})
}

func Test_SummaryMappingScenarios_preserve_optional_people(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"IssueByID": `{"issue":` + issueJSONWithAssignee(issueFixture{
			Identifier: "LIT-30",
			Title:      "assigned",
			ProjectID:  "project-id",
			Project:    "fixture",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}, "Omer") + `}`,
		"ProjectByID": `{"project":` + projectJSONWithLead(projectFixture{
			ID:     "project-id",
			Name:   "led",
			Status: "Backlog",
		}, "Omer") + `}`,
	}

	issue, err := GetIssueByID(context.Background(), graphqlClient, "LIT-30")
	require.NoError(t, err)
	require.Equal(t, "Omer", issue.Assignee)

	project, err := GetProjectByID(context.Background(), graphqlClient, "project-id")
	require.NoError(t, err)
	require.Equal(t, "Omer", project.Lead)
}

func Test_ClientFailureScenarios_wrap_read_and_mutation_errors(t *testing.T) {
	t.Run("read operations wrap graphql errors", func(t *testing.T) {
		graphqlClient := errorGraphQLClient{err: errors.New("network down")}

		_, err := ListIssuesByTeam(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list issues")

		_, err = GetIssueByID(context.Background(), graphqlClient, "LIT-1")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get issue LIT-1")

		_, err = ListProjectsByTeam(context.Background(), graphqlClient, "team-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list projects")

		_, err = GetProjectByID(context.Background(), graphqlClient, "project-id")
		require.Error(t, err)
		require.Contains(t, err.Error(), "get project project-id")

		_, err = ListProjectMembers(context.Background(), graphqlClient, "project-id", 1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "list project members project-id")
	})

	t.Run("issue mutations fail when payload omits entity", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{
			"IssueCreate": `{"issueCreate":{"success":false,"issue":null}}`,
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-20",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueUpdate":        `{"issueUpdate":{"success":false,"issue":null}}`,
			"IssueCommentCreate": `{"commentCreate":{"success":true,"comment":{"id":"comment-id","body":"body","url":"url","issue":null}}}`,
			"CompletedWorkflowStates": `{"workflowStates":{"nodes":[
				{"id":"done-state","name":"Done","type":"completed","position":1}
			]}}`,
			"IssueClose": `{"issueUpdate":{"success":false,"issue":null}}`,
		})

		_, err := CreateIssue(context.Background(), graphqlClient, matchingTarget(), IssueCreateRequest{Title: "title"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = UpdateIssue(context.Background(), graphqlClient, matchingTarget(), IssueUpdateRequest{ID: "LIT-20", Title: "title"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = CommentOnIssue(context.Background(), graphqlClient, matchingTarget(), IssueCommentRequest{ID: "LIT-20", Body: "body"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = CloseIssue(context.Background(), graphqlClient, matchingTarget(), "LIT-20")
		require.ErrorIs(t, err, ErrMutationFailed)
	})

	t.Run("project mutations fail when payload omits entity", func(t *testing.T) {
		graphqlClient := projectWriteFakeClient(map[string]string{
			"ProjectCreate": `{"projectCreate":{"success":false,"project":null}}`,
			"ProjectByID": `{"project":` + projectJSON(projectFixture{
				ID:     "project-id",
				Name:   "fixture",
				Status: "Backlog",
			}) + `}`,
			"ProjectUpdate":  `{"projectUpdate":{"success":false,"project":null}}`,
			"ProjectArchive": `{"projectArchive":{"success":false,"entity":null}}`,
		})

		_, err := CreateProject(context.Background(), graphqlClient, matchingTarget(), ProjectCreateRequest{Name: "name"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = UpdateProject(context.Background(), graphqlClient, matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, ErrMutationFailed)

		_, err = ArchiveProject(context.Background(), graphqlClient, matchingTarget(), "project-id")
		require.ErrorIs(t, err, ErrMutationFailed)
	})

	t.Run("write operations wrap graphql operation errors", func(t *testing.T) {
		operationErr := errors.New("linear unavailable")

		_, err := CreateIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueCreate": "",
		}).withError(operationErr), matchingTarget(), IssueCreateRequest{Title: "title"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "create issue")

		_, err = UpdateIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-40",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueUpdate": "",
		}).withError(operationErr), matchingTarget(), IssueUpdateRequest{ID: "LIT-40", Title: "title"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "update issue LIT-40")

		_, err = CommentOnIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-40",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"IssueCommentCreate": "",
		}).withError(operationErr), matchingTarget(), IssueCommentRequest{ID: "LIT-40", Body: "body"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "comment on issue LIT-40")

		_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-40",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"CompletedWorkflowStates": `{"workflowStates":{"nodes":[{"id":"done-state","name":"Done","type":"completed","position":1}]}}`,
			"IssueClose":              "",
		}).withError(operationErr), matchingTarget(), "LIT-40")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "close issue LIT-40")

		_, err = CreateProject(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectCreate": "",
		}).withError(operationErr), matchingTarget(), ProjectCreateRequest{Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "create project")

		_, err = UpdateProject(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectByID":   `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectUpdate": "",
		}).withError(operationErr), matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "update project project-id")

		_, err = ArchiveProject(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectByID":    `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			"ProjectArchive": "",
		}).withError(operationErr), matchingTarget(), "project-id")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "archive project project-id")
	})

	t.Run("write operations return guard read errors", func(t *testing.T) {
		operationErr := errors.New("guard read failed")

		_, err := UpdateIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": "",
		}).withError(operationErr), matchingTarget(), IssueUpdateRequest{ID: "LIT-50", Title: "title"})
		require.ErrorIs(t, err, operationErr)

		_, err = CommentOnIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": "",
		}).withError(operationErr), matchingTarget(), IssueCommentRequest{ID: "LIT-50", Body: "body"})
		require.ErrorIs(t, err, operationErr)

		_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": "",
		}).withError(operationErr), matchingTarget(), "LIT-50")
		require.ErrorIs(t, err, operationErr)

		_, err = CloseIssue(context.Background(), issueWriteFakeClient(map[string]string{
			"IssueByID": `{"issue":` + issueJSON(issueFixture{
				Identifier: "LIT-51",
				Title:      "target",
				ProjectID:  "project-id",
				Project:    "fixture",
				StateID:    "todo",
				State:      "Todo",
				StateType:  "unstarted",
			}) + `}`,
			"CompletedWorkflowStates": "",
		}).withError(operationErr), matchingTarget(), "LIT-51")
		require.ErrorIs(t, err, operationErr)
		require.Contains(t, err.Error(), "list completed workflow states")

		_, err = UpdateProject(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectByID": "",
		}).withError(operationErr), matchingTarget(), ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, operationErr)

		_, err = ArchiveProject(context.Background(), projectWriteFakeClient(map[string]string{
			"ProjectByID": "",
		}).withError(operationErr), matchingTarget(), "project-id")
		require.ErrorIs(t, err, operationErr)
	})

	t.Run("write operations refuse unpinned targets", func(t *testing.T) {
		graphqlClient := issueWriteFakeClient(map[string]string{})
		emptyTarget := config.Target{}

		_, err := CreateIssue(context.Background(), graphqlClient, emptyTarget, IssueCreateRequest{Title: "title"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = UpdateIssue(context.Background(), graphqlClient, emptyTarget, IssueUpdateRequest{ID: "LIT-1", Title: "title"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = CommentOnIssue(context.Background(), graphqlClient, emptyTarget, IssueCommentRequest{ID: "LIT-1", Body: "body"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = CloseIssue(context.Background(), graphqlClient, emptyTarget, "LIT-1")
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = CreateProject(context.Background(), graphqlClient, emptyTarget, ProjectCreateRequest{Name: "name"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = UpdateProject(context.Background(), graphqlClient, emptyTarget, ProjectUpdateRequest{ID: "project-id", Name: "name"})
		require.ErrorIs(t, err, ErrTargetMismatch)

		_, err = ArchiveProject(context.Background(), graphqlClient, emptyTarget, "project-id")
		require.ErrorIs(t, err, ErrTargetMismatch)
	})
}

func Test_TargetFailureScenarios_refuse_unpinned_or_mismatched_targets(t *testing.T) {
	_, err := ResolveTarget(context.Background(), fakeGraphQLClient{}, config.Target{})
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer": `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":  "",
	}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve teams")

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"other-org","name":"Other","urlKey":"other"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`,
	}, matchingTarget())
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = ResolveTarget(context.Background(), fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": "",
	}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve project")

	graphqlClient := fakeGraphQLClient{
		"Viewer":        `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`,
		"Teams":         `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`,
		"TargetProject": `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"other-team","key":"ABC","name":"other","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`,
	}

	_, err = ResolveTarget(context.Background(), graphqlClient, matchingTarget())
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = firstCompletedStateID(context.Background(), fakeGraphQLClient{
		"CompletedWorkflowStates": `{"workflowStates":{"nodes":[]}}`,
	}, "team-id")
	require.ErrorIs(t, err, ErrWriteInvalid)

	err = requireTargetMatch(config.Target{OrgID: "org-id", TeamID: "team-id", TeamKey: "LIT"}, config.Target{
		OrgID:   "other-org",
		TeamID:  "team-id",
		TeamKey: "LIT",
	})
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_TransportScenarios_return_actionable_errors(t *testing.T) {
	require.Equal(t, "fallback", firstNonEmpty("", "fallback"))
	require.Equal(t, "primary", firstNonEmpty("primary", "fallback"))
	require.Equal(t, 3*time.Second, defaultDuration(3*time.Second, time.Second))
	require.Equal(t, time.Second, defaultDuration(0, time.Second))
	require.Equal(t, 200*time.Millisecond, retryDelay("", 1))
	require.Equal(t, 100*time.Millisecond, retryDelay("not-a-number", 0))
	require.Equal(t, 2*time.Second, retryDelay("2", 0))

	response := graphql.Response{}
	err := decodeGraphQLResponse([]byte("not json"), http.StatusOK, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "decode graphql response")

	err = decodeGraphQLResponse([]byte("server down"), http.StatusBadGateway, &response)
	require.Error(t, err)
	require.Contains(t, err.Error(), "graphql http status 502")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = waitForRateLimitRetry(ctx, http.StatusTooManyRequests, http.Header{}, 0, 1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "wait for retry")
}

type errorGraphQLClient struct {
	err error
}

func (client errorGraphQLClient) MakeRequest(
	_ context.Context,
	_ *graphql.Request,
	_ *graphql.Response,
) error {
	return client.err
}

type operationErrorFakeClient struct {
	responses map[string]string
	err       error
}

func (client operationErrorFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if client.responses[request.OpName] == "" {
		return client.err
	}

	return fakeGraphQLClient(client.responses).MakeRequest(ctx, request, response)
}

func (client issueWriteFakeClient) withError(err error) operationErrorFakeClient {
	return operationErrorFakeClient{
		responses: client.withTargetResponses(),
		err:       err,
	}
}

func (client projectWriteFakeClient) withError(err error) operationErrorFakeClient {
	return operationErrorFakeClient{
		responses: client.withTargetResponses(),
		err:       err,
	}
}

func Test_WriteGuardScenarios_refuse_mismatched_resources(t *testing.T) {
	guard := writeGuard{
		target: ResolvedTarget{
			Team: TargetTeam{ID: "team-id", Key: "LIT"},
		},
	}
	graphqlClient := fakeGraphQLClient{
		"IssueByID": `{"issue":` + strings.ReplaceAll(issueJSON(issueFixture{
			Identifier: "ABC-1",
			Title:      "wrong team",
			StateID:    "todo",
			State:      "Todo",
			StateType:  "unstarted",
		}), `"key":"LIT"`, `"key":"ABC"`) + `}`,
		"ProjectByID": `{"project":` + strings.ReplaceAll(projectJSON(projectFixture{
			ID:     "project-id",
			Name:   "wrong-team",
			Status: "Backlog",
		}), `"key":"LIT"`, `"key":"ABC"`) + `}`,
	}

	_, err := guard.requireIssue(context.Background(), graphqlClient, "ABC-1")
	require.ErrorIs(t, err, ErrTargetMismatch)

	err = guard.requireProject(context.Background(), graphqlClient, "project-id")
	require.ErrorIs(t, err, ErrTargetMismatch)

	_, err = newWriteGuard(context.Background(), errorGraphQLClient{err: errors.New("resolve failed")}, matchingTarget())
	require.Error(t, err)
	require.Contains(t, err.Error(), "resolve failed")

	_, err = guard.requireIssue(context.Background(), errorGraphQLClient{err: errors.New("read issue failed")}, "LIT-1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "read issue failed")

	err = guard.requireProject(context.Background(), errorGraphQLClient{err: errors.New("read project failed")}, "project-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "read project failed")
}

func Test_FakeGraphQLClient_respects_context_and_missing_operations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := fakeGraphQLClient{}.MakeRequest(ctx, &graphql.Request{OpName: "Viewer"}, &graphql.Response{})
	require.Error(t, err)

	err = fakeGraphQLClient{}.MakeRequest(context.Background(), &graphql.Request{OpName: "Viewer"}, &graphql.Response{})
	require.Error(t, err)
	require.True(t, errors.Is(err, errors.New("missing fake response for Viewer")) || strings.Contains(err.Error(), "missing fake response"))
}

func Test_TargetScenarios_allow_unpinned_project_and_matching_team(t *testing.T) {
	require.Nil(t, optionalString(""))
	require.Equal(t, "value", *optionalString("value"))
	require.Equal(t, "value", *stringPtr("value"))
	require.Equal(t, 7, *intPtr(7))
	require.True(t, *boolPtr(true))
	require.True(t, projectHasTeam(ProjectSummary{Teams: []ProjectTeam{{ID: "team-id", Key: "LIT"}}}, "team-id", "LIT"))
	require.False(t, projectHasTeam(ProjectSummary{Teams: []ProjectTeam{{ID: "team-id", Key: "ABC"}}}, "team-id", "LIT"))

	guard, err := newWriteGuard(context.Background(), issueWriteFakeClient(map[string]string{}), config.Target{
		OrgID:   "org-id",
		TeamKey: "LIT",
		TeamID:  "team-id",
	})

	require.NoError(t, err)
	require.Nil(t, guard.target.Project)

	err = validateProjectUpdateRequest(ProjectUpdateRequest{Name: "missing id"})
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func issueJSONWithAssignee(issue issueFixture, assignee string) string {
	return strings.ReplaceAll(issueJSON(issue), `"assignee":null`, `"assignee":{"id":"user-id","name":"omer","displayName":"`+assignee+`"}`)
}

func projectJSONWithLead(project projectFixture, lead string) string {
	return strings.Replace(projectJSON(project), `"lead":null`, `"lead":{"id":"user-id","name":"omer","displayName":"`+lead+`"}`, 1)
}
