package client

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const projectInitiativeLinkJSON = `{
	"id":"link-id",
	"sortOrder":"2.5",
	"createdAt":"2026-11-08T00:00:00Z",
	"updatedAt":"2026-11-08T00:00:00Z",
	"archivedAt":null,
	"initiative":{"id":"created-initiative-id","name":"Platform"},
	"project":{"id":"project-id","name":"fixture","slugId":"fixture","url":"https://linear.app/project/fixture"}
}`

func projectInitiativeLinksPayload(hasNextPage bool, links ...string) string {
	nodes := ""
	for index, link := range links {
		if index > 0 {
			nodes += ","
		}
		nodes += link
	}
	next := "false"
	if hasNextPage {
		next = "true"
	}

	return `{"project":{"id":"project-id","name":"fixture","initiativeToProjects":{"nodes":[` + nodes +
		`],"pageInfo":{"hasNextPage":` + next + `,"endCursor":"cursor-1"}}}}`
}

func otherInitiativeLinks(count int) []string {
	links := make([]string, count)
	for index := range links {
		links[index] = strings.Replace(projectInitiativeLinkJSON, `"id":"created-initiative-id"`,
			`"id":"other-initiative-`+strconv.Itoa(index)+`"`, 1)
	}

	return links
}

func projectInitiativeResponses(initiativeOrgID string, links string) map[string]string {
	return map[string]string{
		"project":                      `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
		"initiative":                   `{"initiative":` + initiativeJSON(initiativeOrgID) + `}`,
		"project_initiativeToProjects": links,
	}
}

func linkRequest() ProjectInitiativeLinkRequest {
	return ProjectInitiativeLinkRequest{ProjectID: "project-id", InitiativeID: "created-initiative-id"}
}

func Test_LinkProjectInitiative_creates_link_with_sort_order(t *testing.T) {
	responses := projectInitiativeResponses("org-id", projectInitiativeLinksPayload(false))
	responses["InitiativeToProjectCreate"] = `{"initiativeToProjectCreate":{"success":true,"initiativeToProject":` +
		projectInitiativeLinkJSON + `}}`
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(responses)}
	sortOrder := 2.5
	request := linkRequest()
	request.SortOrder = &sortOrder

	link, err := LinkProjectInitiative(context.Background(), recorder, matchingTarget(), request)

	require.NoError(t, err)
	require.Equal(t, "link-id", link.ID)
	require.Equal(t, "created-initiative-id", link.InitiativeID)
	require.JSONEq(
		t,
		`{"input":{"id":null,"initiativeId":"created-initiative-id","projectId":"project-id","sortOrder":2.5}}`,
		string(recorder.variablesFor(t, "InitiativeToProjectCreate")),
	)
}

func Test_LinkProjectInitiative_returns_existing_link_without_creating(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(
		projectInitiativeResponses("org-id", projectInitiativeLinksPayload(false, projectInitiativeLinkJSON)),
	)}

	link, err := LinkProjectInitiative(context.Background(), recorder, matchingTarget(), linkRequest())

	require.NoError(t, err)
	require.Equal(t, "link-id", link.ID)
	require.Zero(t, recorder.countOf("InitiativeToProjectCreate"))
}

func Test_UnlinkProjectInitiative_deletes_the_live_link(t *testing.T) {
	archivedLink := `{"id":"archived-link-id","sortOrder":"1","createdAt":"2026-01-01T00:00:00Z",` +
		`"updatedAt":"2026-01-01T00:00:00Z","archivedAt":"2026-02-01T00:00:00Z",` +
		`"initiative":{"id":"created-initiative-id","name":"Platform"},` +
		`"project":{"id":"project-id","name":"fixture","slugId":"fixture","url":"https://linear.app/project/fixture"}}`
	responses := projectInitiativeResponses(
		"org-id", projectInitiativeLinksPayload(false, archivedLink, projectInitiativeLinkJSON),
	)
	responses["InitiativeToProjectDelete"] = `{"initiativeToProjectDelete":{"success":true}}`
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(responses)}

	link, err := UnlinkProjectInitiative(context.Background(), recorder, matchingTarget(), linkRequest())

	require.NoError(t, err)
	require.Equal(t, "link-id", link.ID)
	require.JSONEq(t, `{"id":"link-id"}`, string(recorder.variablesFor(t, "InitiativeToProjectDelete")))
}

func Test_ProjectInitiativeWrites_refuse_before_any_mutation(t *testing.T) {
	otherProject := map[string]string{
		"project": `{"project":` + projectJSON(projectFixture{ID: "other-project", Name: "other", Status: "Backlog"}) + `}`,
	}
	cases := []struct {
		name      string
		responses map[string]string
		request   ProjectInitiativeLinkRequest
		want      error
	}{
		{
			name:      "missing initiative id",
			responses: map[string]string{},
			request:   ProjectInitiativeLinkRequest{ProjectID: "project-id"},
			want:      ErrWriteInvalid,
		},
		{
			name:      "project outside the pinned project",
			responses: otherProject,
			request:   ProjectInitiativeLinkRequest{ProjectID: "other-project", InitiativeID: "created-initiative-id"},
			want:      ErrTargetMismatch,
		},
		{
			name:      "initiative in another organization",
			responses: projectInitiativeResponses("other-org", projectInitiativeLinksPayload(false)),
			request:   linkRequest(),
			want:      ErrTargetMismatch,
		},
		{
			name:      "more links than one page without a match",
			responses: projectInitiativeResponses("org-id", projectInitiativeLinksPayload(true, otherInitiativeLinks(defaultListPageSize)...)),
			request:   linkRequest(),
			want:      ErrWriteInvalid,
		},
	}
	writes := map[string]func(ProjectInitiativeLinkRequest, *recordingGraphQLClient) error{
		"link": func(request ProjectInitiativeLinkRequest, recorder *recordingGraphQLClient) error {
			_, err := LinkProjectInitiative(context.Background(), recorder, matchingTarget(), request)
			return err
		},
		"unlink": func(request ProjectInitiativeLinkRequest, recorder *recordingGraphQLClient) error {
			_, err := UnlinkProjectInitiative(context.Background(), recorder, matchingTarget(), request)
			return err
		},
	}
	for _, testCase := range cases {
		for writeName, write := range writes {
			t.Run(writeName+" "+testCase.name, func(t *testing.T) {
				recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(testCase.responses)}

				err := write(testCase.request, recorder)

				require.ErrorIs(t, err, testCase.want)
				require.Zero(t, recorder.countOf("InitiativeToProjectCreate"))
				require.Zero(t, recorder.countOf("InitiativeToProjectDelete"))
			})
		}
	}
}

func Test_UnlinkProjectInitiative_refuses_when_project_is_not_linked(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(
		projectInitiativeResponses("org-id", projectInitiativeLinksPayload(false)),
	)}

	_, err := UnlinkProjectInitiative(context.Background(), recorder, matchingTarget(), linkRequest())

	require.ErrorIs(t, err, ErrWriteInvalid)
	require.ErrorContains(t, err, "is not linked")
	require.Zero(t, recorder.countOf("InitiativeToProjectDelete"))
}

func Test_ProjectInitiativeWrites_report_mutation_failures(t *testing.T) {
	cases := []struct {
		name     string
		mutation map[string]string
		write    func(*recordingGraphQLClient) error
		want     error
	}{
		{
			name:     "link reports no success",
			mutation: map[string]string{"InitiativeToProjectCreate": `{"initiativeToProjectCreate":{"success":false,"initiativeToProject":null}}`},
			write:    linkWith,
			want:     ErrMutationFailed,
		},
		{
			name:     "unlink reports no success",
			mutation: map[string]string{"InitiativeToProjectDelete": `{"initiativeToProjectDelete":{"success":false}}`},
			write:    unlinkWith,
			want:     ErrMutationFailed,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			links := projectInitiativeLinksPayload(false)
			if testCase.name == "unlink reports no success" {
				links = projectInitiativeLinksPayload(false, projectInitiativeLinkJSON)
			}
			responses := projectInitiativeResponses("org-id", links)
			for operation, payload := range testCase.mutation {
				responses[operation] = payload
			}

			err := testCase.write(&recordingGraphQLClient{inner: projectWriteFakeClient(responses)})

			require.ErrorIs(t, err, testCase.want)
		})
	}
}

func Test_ProjectInitiativeWrites_wrap_transport_errors(t *testing.T) {
	cases := map[string]struct {
		responses map[string]string
		write     func(*recordingGraphQLClient) error
		contains  string
	}{
		"link create fails": {
			responses: projectInitiativeResponses("org-id", projectInitiativeLinksPayload(false)),
			write:     linkWith,
			contains:  "link project project-id to initiative created-initiative-id",
		},
		"unlink delete fails": {
			responses: projectInitiativeResponses("org-id", projectInitiativeLinksPayload(false, projectInitiativeLinkJSON)),
			write:     unlinkWith,
			contains:  "unlink project project-id from initiative created-initiative-id",
		},
		"link lookup fails": {
			responses: map[string]string{
				"project":    `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
				"initiative": `{"initiative":` + initiativeJSON("org-id") + `}`,
			},
			write:    unlinkWith,
			contains: "list project initiative associations project-id",
		},
		"initiative lookup fails": {
			responses: map[string]string{
				"project": `{"project":` + projectJSON(projectFixture{ID: "project-id", Name: "fixture", Status: "Backlog"}) + `}`,
			},
			write:    linkWith,
			contains: "get initiative created-initiative-id",
		},
	}
	for name, testCase := range cases {
		t.Run(name, func(t *testing.T) {
			err := testCase.write(&recordingGraphQLClient{inner: projectWriteFakeClient(testCase.responses)})

			require.ErrorContains(t, err, testCase.contains)
			require.NotErrorIs(t, err, ErrTargetMismatch)
		})
	}
}

func linkWith(recorder *recordingGraphQLClient) error {
	_, err := LinkProjectInitiative(context.Background(), recorder, matchingTarget(), linkRequest())
	return err
}

func unlinkWith(recorder *recordingGraphQLClient) error {
	_, err := UnlinkProjectInitiative(context.Background(), recorder, matchingTarget(), linkRequest())
	return err
}

func Test_LinkProjectInitiative_refuses_when_pinned_team_does_not_resolve(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(
		projectInitiativeResponses("org-id", projectInitiativeLinksPayload(false)),
	)}
	target := matchingTarget()
	target.TeamKey, target.TeamID = "WRONG", "wrong-id"

	_, err := LinkProjectInitiative(context.Background(), recorder, target, linkRequest())

	require.ErrorIs(t, err, ErrTargetMismatch)
	require.Zero(t, recorder.countOf("project"))
}
