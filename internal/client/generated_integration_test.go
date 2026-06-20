//go:build integration

package client

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

type liveIntegrationConfig struct {
	OrgID     string `json:"org_id"`
	TeamKey   string `json:"team_key"`
	TeamID    string `json:"team_id"`
	ProjectID string `json:"project_id"`
}

func Test_Integration_generatedViewerOrganizationTeams_whenTokenConfigured(t *testing.T) {
	// Given
	config := readLiveIntegrationConfig(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	transport := newLiveIntegrationTransport(t, 5*time.Second)
	first := 50
	includeArchived := true

	// When
	viewer, viewerErr := Viewer(ctx, transport)
	organization, organizationErr := Organization(ctx, transport)
	teams, teamsErr := Teams(ctx, transport, &first, nil, &includeArchived)

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

func Test_Integration_projectWriteRoundTrip_whenTargetPinned(t *testing.T) {
	// Given
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
		ProjectID: fixture.ProjectID,
	}, created.ID)

	// Then
	require.NoError(t, listErr)
	require.NoError(t, readErr)
	require.NoError(t, updateErr)
	require.NoError(t, membersErr)
	require.ErrorIs(t, wrongProjectErr, ErrTargetMismatch)
	require.NotEmpty(t, listed.Projects)
	require.Equal(t, created.ID, read.ID)
	require.Equal(t, name+" updated", updated.Name)
	require.Equal(t, created.ID, members.ProjectID)
}

func readLiveIntegrationConfig(t *testing.T) liveIntegrationConfig {
	t.Helper()

	data, err := os.ReadFile("../../test/integration-config.json")
	require.NoError(t, err)

	var config liveIntegrationConfig
	require.NoError(t, json.Unmarshal(data, &config))

	return config
}

func newLiveIntegrationTransport(t *testing.T, timeout time.Duration) *Transport {
	t.Helper()

	token := os.Getenv("LINCTL_TEST_TOKEN")
	if token == "" {
		t.Skip("LINCTL_TEST_TOKEN is required for integration tests")
	}

	return NewTransport(TransportConfig{
		Token:      PersonalAPIToken(token),
		Timeout:    timeout,
		MaxRetries: 1,
	})
}

func archiveIntegrationIssue(ctx context.Context, transport *Transport, issueID string) (IssueSummary, error) {
	archived, err := IssueArchive(ctx, transport, issueID, boolPtr(false))
	if err != nil {
		return IssueSummary{}, err
	}
	if !archived.IssueArchive.Success || archived.IssueArchive.Entity == nil {
		return IssueSummary{}, ErrMutationFailed
	}

	return issueSummaryFromFields(archived.IssueArchive.Entity.IssueSummaryFields), nil
}

func containsTeamID(teams []TeamsTeamsTeamConnectionNodesTeam, teamID string) bool {
	for _, team := range teams {
		if team.Id == teamID {
			return true
		}
	}

	return false
}
