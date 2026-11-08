//go:build integration

package client

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
	"github.com/KyaniteHQ/linctl/internal/oauth"
)

type liveIntegrationConfig struct {
	OrgID     string `json:"org_id"`
	TeamKey   string `json:"team_key"`
	TeamID    string `json:"team_id"`
	ProjectID string `json:"project_id"`
}

const (
	liveWriteIntegrationEnv  = "LINCTL_TEST_ENABLE_WRITES"
	liveInitiativeFixtureEnv = "LINCTL_TEST_INITIATIVE_ID"
)

func Test_Integration_generatedViewerOrganizationTeams_whenTokenConfigured(t *testing.T) {
	// Given
	config := readLiveIntegrationConfig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	transport := newLiveIntegrationTransport(t, 5*time.Second)
	first := 50
	includeArchived := true

	// When
	viewer, viewerErr := gql.Viewer(ctx, transport)
	organization, organizationErr := gql.Organization(ctx, transport)
	teams, teamsErr := gql.Teams(ctx, transport, &first, nil, &includeArchived)

	// Then
	require.NoError(t, viewerErr)
	require.NoError(t, organizationErr)
	require.NoError(t, teamsErr)
	require.Equal(t, config.OrgID, viewer.Viewer.Organization.Id)
	require.Equal(t, config.OrgID, organization.Organization.Id)
	require.True(t, containsTeamID(teams.Teams.Nodes, config.TeamID))
}

func Test_Integration_issueWriteRoundTrip_whenTargetPinned(t *testing.T) {
	// Given
	requireLiveWriteIntegration(t)
	fixture := readLiveIntegrationConfig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	transport := newLiveIntegrationTransport(t, 10*time.Second)
	target := config.Target{
		OrgID:     fixture.OrgID,
		TeamKey:   fixture.TeamKey,
		TeamID:    fixture.TeamID,
		ProjectID: fixture.ProjectID,
	}
	title := "linctl-it-" + time.Now().UTC().Format("20060102T150405")

	// When
	created, createErr := CreateIssue(ctx, transport, target, IssueCreateRequest{
		Title:       title,
		Description: "created by linctl integration test",
	})
	require.NoError(t, createErr)
	defer func() {
		_, err := archiveIntegrationIssue(context.Background(), transport, created.ID)
		require.NoError(t, err)
	}()
	read, readErr := GetIssueByID(ctx, transport, created.Identifier)
	updated, updateErr := UpdateIssue(ctx, transport, target, IssueUpdateRequest{
		ID:    created.Identifier,
		Title: title + " updated",
	})
	comment, commentErr := CommentOnIssue(ctx, transport, target, IssueCommentRequest{
		ID:   created.Identifier,
		Body: "linctl integration comment",
	})
	closed, closeErr := CloseIssue(ctx, transport, target, created.Identifier)

	// Then
	require.NoError(t, readErr)
	require.NoError(t, updateErr)
	require.NoError(t, commentErr)
	require.NoError(t, closeErr)
	require.Equal(t, created.Identifier, read.Identifier)
	require.Equal(t, title+" updated", updated.Title)
	require.Equal(t, created.Identifier, comment.Issue.Identifier)
	require.Equal(t, "completed", closed.StateType)
}

func Test_Integration_issueClose_whenCloseStateConfigured(t *testing.T) {
	// Given
	requireLiveWriteIntegration(t)
	fixture := readLiveIntegrationConfig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	transport := newLiveIntegrationTransport(t, 10*time.Second)
	states, listErr := listTeamWorkflowStates(ctx, transport, fixture.TeamID)
	require.NoError(t, listErr)
	var want *workflowStateCandidate
	for index := range states {
		// The highest-position completed state is never the one a bare close picks
		// when the team has several, so the configured name has to win.
		if states[index].Type == "completed" && (want == nil || states[index].Position > want.Position) {
			want = &states[index]
		}
	}
	require.NotNil(t, want, "the team has no completed workflow state")
	target := config.Target{
		OrgID:     fixture.OrgID,
		TeamKey:   fixture.TeamKey,
		TeamID:    fixture.TeamID,
		ProjectID: fixture.ProjectID,
		States:    config.States{Close: want.Name},
	}

	// When
	created, createErr := CreateIssue(ctx, transport, target, IssueCreateRequest{
		Title:       "linctl-it-" + time.Now().UTC().Format("20060102T150405"),
		Description: "created by linctl integration test",
	})
	require.NoError(t, createErr)
	defer func() {
		_, err := archiveIntegrationIssue(context.Background(), transport, created.ID)
		require.NoError(t, err)
	}()
	closed, closeErr := CloseIssue(ctx, transport, target, created.Identifier)

	// Then
	require.NoError(t, closeErr)
	require.Equal(t, want.ID, closed.StateID)
	require.Equal(t, "completed", closed.StateType)
}

func Test_Integration_projectWriteRoundTrip_whenTargetPinned(t *testing.T) {
	// Given
	requireLiveWriteIntegration(t)
	fixture := readLiveIntegrationConfig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	transport := newLiveIntegrationTransport(t, 10*time.Second)
	target := config.Target{
		OrgID:   fixture.OrgID,
		TeamKey: fixture.TeamKey,
		TeamID:  fixture.TeamID,
	}
	name := "linctl-it-" + time.Now().UTC().Format("20060102T150405")

	// When
	listed, listErr := ListProjectsByTeam(ctx, transport, fixture.TeamID, 100)
	require.NoError(t, listErr)
	require.NotEmpty(t, listed.Projects)
	wrongProjectID := listed.Projects[0].ID
	created, createErr := CreateProject(ctx, transport, target, ProjectCreateRequest{
		Name:        name,
		Description: "created by linctl integration test",
	})
	require.NoError(t, createErr)
	defer func() {
		_, err := ArchiveProject(context.Background(), transport, config.Target{
			OrgID:     fixture.OrgID,
			TeamKey:   fixture.TeamKey,
			TeamID:    fixture.TeamID,
			ProjectID: created.ID,
		}, created.ID)
		require.NoError(t, err)
	}()
	read, readErr := GetProjectByID(ctx, transport, created.ID)
	updated, updateErr := UpdateProject(ctx, transport, config.Target{
		OrgID:     fixture.OrgID,
		TeamKey:   fixture.TeamKey,
		TeamID:    fixture.TeamID,
		ProjectID: created.ID,
	}, ProjectUpdateRequest{
		ID:   created.ID,
		Name: name + " updated",
	})
	members, membersErr := ListProjectMembers(ctx, transport, created.ID, 10)
	_, wrongProjectErr := ArchiveProject(ctx, transport, config.Target{
		OrgID:     fixture.OrgID,
		TeamKey:   fixture.TeamKey,
		TeamID:    fixture.TeamID,
		ProjectID: wrongProjectID,
	}, created.ID)

	// Then
	require.NoError(t, readErr)
	require.NoError(t, updateErr)
	require.NoError(t, membersErr)
	require.ErrorIs(t, wrongProjectErr, ErrTargetMismatch)
	require.Equal(t, created.ID, read.ID)
	require.Equal(t, name+" updated", updated.Name)
	require.Equal(t, created.ID, members.ProjectID)
}

func Test_Integration_projectInitiativeLinkRoundTrip_whenTargetPinned(t *testing.T) {
	requireLiveWriteIntegration(t)
	fixture := readLiveIntegrationConfig(t)
	initiativeID := os.Getenv(liveInitiativeFixtureEnv)
	if initiativeID == "" {
		t.Skip(liveInitiativeFixtureEnv + " must name the test organization's fixture Initiative")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	transport := newLiveIntegrationTransport(t, 10*time.Second)
	teamTarget := config.Target{OrgID: fixture.OrgID, TeamKey: fixture.TeamKey, TeamID: fixture.TeamID}
	listed, listErr := ListProjectsByTeam(ctx, transport, fixture.TeamID, 100)
	require.NoError(t, listErr)
	require.NotEmpty(t, listed.Projects)
	wrongProjectID := listed.Projects[0].ID
	created, createErr := CreateProject(ctx, transport, teamTarget, ProjectCreateRequest{
		Name:        "linctl-it-" + time.Now().UTC().Format("20060102T150405"),
		Description: "created by linctl integration test",
	})
	require.NoError(t, createErr)
	projectTarget := teamTarget
	projectTarget.ProjectID = created.ID
	request := ProjectInitiativeLinkRequest{ProjectID: created.ID, InitiativeID: initiativeID}
	t.Cleanup(func() {
		_, err := ArchiveProject(context.Background(), transport, projectTarget, created.ID)
		require.NoError(t, err)
	})
	t.Cleanup(func() {
		_, err := UnlinkProjectInitiative(context.Background(), transport, projectTarget, request)
		if !errors.Is(err, ErrWriteInvalid) {
			require.NoError(t, err)
		}
	})
	sortOrder := 7.0
	sorted := request
	sorted.SortOrder = &sortOrder

	linked, linkErr := LinkProjectInitiative(ctx, transport, projectTarget, sorted)
	relinked, relinkErr := LinkProjectInitiative(ctx, transport, projectTarget, request)
	afterLink, afterLinkErr := ListProjectInitiatives(ctx, transport, created.ID, 10)
	wrongTarget := teamTarget
	wrongTarget.ProjectID = wrongProjectID
	_, wrongProjectErr := UnlinkProjectInitiative(ctx, transport, wrongTarget, request)
	unlinked, unlinkErr := UnlinkProjectInitiative(ctx, transport, projectTarget, request)
	afterUnlink, afterUnlinkErr := ListProjectInitiatives(ctx, transport, created.ID, 10)
	_, secondUnlinkErr := UnlinkProjectInitiative(ctx, transport, projectTarget, request)
	linkedAgain, linkAgainErr := LinkProjectInitiative(ctx, transport, projectTarget, request)
	unlinkedAgain, unlinkAgainErr := UnlinkProjectInitiative(ctx, transport, projectTarget, request)

	require.NoError(t, linkErr)
	require.NoError(t, relinkErr)
	require.NoError(t, afterLinkErr)
	require.NoError(t, unlinkErr)
	require.NoError(t, afterUnlinkErr)
	require.NoError(t, linkAgainErr)
	require.NoError(t, unlinkAgainErr)
	require.ErrorIs(t, wrongProjectErr, ErrTargetMismatch)
	require.ErrorIs(t, secondUnlinkErr, ErrWriteInvalid)
	require.Equal(t, initiativeID, linked.InitiativeID)
	require.Equal(t, created.ID, linked.ProjectID)
	linkedOrder, parseErr := strconv.ParseFloat(linked.SortOrder, 64)
	require.NoError(t, parseErr)
	require.InDelta(t, sortOrder, linkedOrder, 0)
	require.Equal(t, linked.ID, relinked.ID)
	require.Equal(t, linked.ID, unlinked.ID)
	require.NotEqual(t, linked.ID, linkedAgain.ID)
	require.Equal(t, linkedAgain.ID, unlinkedAgain.ID)
	require.Contains(t, initiativeIDs(afterLink.Initiatives), initiativeID)
	require.NotContains(t, initiativeIDs(afterUnlink.Initiatives), initiativeID)
}

func initiativeIDs(initiatives []InitiativeSummary) []string {
	ids := make([]string, 0, len(initiatives))
	for _, initiative := range initiatives {
		ids = append(ids, initiative.ID)
	}

	return ids
}

func Test_Integration_issueCoordinationWriteRoundTrip_whenTargetPinned(t *testing.T) {
	// Given
	requireLiveWriteIntegration(t)
	fixture := readLiveIntegrationConfig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	transport := newLiveIntegrationTransport(t, 10*time.Second)
	target := config.Target{
		OrgID:     fixture.OrgID,
		TeamKey:   fixture.TeamKey,
		TeamID:    fixture.TeamID,
		ProjectID: fixture.ProjectID,
	}
	runID := "linctl-it-" + time.Now().UTC().Format("20060102T150405")

	source, sourceErr := CreateIssue(ctx, transport, target, IssueCreateRequest{
		Title:       runID + " coordination source",
		Description: "created by linctl coordination integration test",
	})
	require.NoError(t, sourceErr)
	registerLiveIssueCleanup(t, transport, source.ID)
	related, relatedErr := CreateIssue(ctx, transport, target, IssueCreateRequest{
		Title:       runID + " coordination related",
		Description: "created by linctl coordination integration test",
	})
	require.NoError(t, relatedErr)
	registerLiveIssueCleanup(t, transport, related.ID)

	// When
	started, startErr := StartIssue(ctx, transport, target, source.Identifier)
	if errors.Is(startErr, ErrViewerNotAssignable) {
		require.Equal(t, "app", os.Getenv("LINCTL_OAUTH_EXPECTED_ACTOR"))
		started, startErr = UpdateIssue(ctx, transport, target, IssueUpdateRequest{
			ID:        source.Identifier,
			StateType: "started",
		})
	}
	require.NoError(t, startErr)
	comment, commentErr := CommentOnIssue(ctx, transport, target, IssueCommentRequest{
		ID:   source.Identifier,
		Body: "linctl integration coordination comment",
	})
	require.NoError(t, commentErr)
	updatedComment, updateCommentErr := UpdateComment(ctx, transport, target, CommentUpdateRequest{
		ID:   comment.ID,
		Body: "linctl integration coordination comment updated",
	})
	require.NoError(t, updateCommentErr)
	relation, relateErr := CreateIssueRelation(ctx, transport, target, IssueRelationCreateRequest{
		IssueID:        source.Identifier,
		RelatedIssueID: related.Identifier,
		Type:           "related",
	})
	require.NoError(t, relateErr)
	deletedRelationID, unrelateErr := DeleteIssueRelation(
		ctx, transport, target, IssueRelationDeleteRequest{RelationID: relation.ID},
	)
	require.NoError(t, unrelateErr)
	deletedCommentID, deleteCommentErr := DeleteComment(ctx, transport, target, comment.ID)
	require.NoError(t, deleteCommentErr)

	// Then
	require.Equal(t, source.Identifier, started.Identifier)
	require.Equal(t, "started", started.StateType)
	require.Equal(t, comment.ID, updatedComment.ID)
	require.Equal(t, "linctl integration coordination comment updated", updatedComment.Body)
	require.Equal(t, relation.ID, deletedRelationID)
	require.Equal(t, comment.ID, deletedCommentID)
}

func readLiveIntegrationConfig(t *testing.T) liveIntegrationConfig {
	t.Helper()

	// Prefer env vars (the CI path, fed by repo secrets) so the suite is not
	// tied to a committed workspace file. Fall back to the local untracked
	// config for developer runs.
	if config, ok := liveIntegrationConfigFromEnv(); ok {
		return config
	}

	data, err := os.ReadFile("../../test/integration-config.json")
	require.NoError(t, err)

	var config liveIntegrationConfig
	require.NoError(t, json.Unmarshal(data, &config))

	return config
}

func liveIntegrationConfigFromEnv() (liveIntegrationConfig, bool) {
	config := liveIntegrationConfig{
		OrgID:     os.Getenv("LINCTL_TEST_ORG_ID"),
		TeamKey:   os.Getenv("LINCTL_TEST_TEAM_KEY"),
		TeamID:    os.Getenv("LINCTL_TEST_TEAM_ID"),
		ProjectID: os.Getenv("LINCTL_TEST_PROJECT_ID"),
	}
	if config.OrgID == "" || config.TeamKey == "" || config.TeamID == "" {
		return liveIntegrationConfig{}, false
	}

	return config, true
}

func newLiveIntegrationTransport(t *testing.T, timeout time.Duration) *Transport {
	t.Helper()

	clientID := os.Getenv("LINCTL_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("LINCTL_OAUTH_CLIENT_SECRET")
	scopes := splitLiveOAuthScopes(os.Getenv("LINCTL_OAUTH_SCOPES"))
	if clientID == "" || clientSecret == "" || len(scopes) == 0 {
		t.Skip("OAuth fixture env is required for integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	grant, err := oauth.NewClient(oauth.ClientConfig{}).ClientCredentials(ctx, oauth.ClientCredentialsRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
	})
	require.NoError(t, err)

	return NewTransport(TransportConfig{
		Token:      OAuthAccessToken(grant.AccessToken),
		Timeout:    timeout,
		MaxRetries: 1,
	})
}

func splitLiveOAuthScopes(value string) []string {
	return strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n'
	})
}

func registerLiveIssueCleanup(t *testing.T, transport *Transport, issueID string) {
	t.Helper()

	t.Cleanup(func() {
		_, err := archiveIntegrationIssue(context.Background(), transport, issueID)
		require.NoError(t, err)
	})
}

func requireLiveWriteIntegration(t *testing.T) {
	t.Helper()

	if os.Getenv(liveWriteIntegrationEnv) == "1" {
		return
	}

	t.Skip(liveWriteIntegrationEnv + "=1 is required for live write integration tests")
}

func archiveIntegrationIssue(ctx context.Context, transport *Transport, issueID string) (IssueSummary, error) {
	archived, err := gql.IssueArchive(ctx, transport, issueID, boolPtr(false))
	if err != nil {
		return IssueSummary{}, err
	}
	if !archived.IssueArchive.Success || archived.IssueArchive.Entity == nil {
		return IssueSummary{}, ErrMutationFailed
	}

	return issueSummaryFromFields(archived.IssueArchive.Entity.IssueSummaryFields), nil
}

func containsTeamID(teams []gql.TeamsTeamsTeamConnectionNodesTeam, teamID string) bool {
	for _, team := range teams {
		if team.Id == teamID {
			return true
		}
	}

	return false
}
