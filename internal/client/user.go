package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// UserSummary is the compact User model used by user commands.
type UserSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Active      bool   `json:"active"`
	Guest       bool   `json:"guest"`
	Admin       bool   `json:"admin"`
}

// UserList is a page of users.
type UserList struct {
	Users       []UserSummary `json:"users"`
	HasNextPage bool          `json:"has_next_page"`
	EndCursor   *string       `json:"end_cursor,omitempty"`
}

// DraftSummary is the compact saved draft model used by viewer-scoped draft reads.
type DraftSummary struct {
	ID          string `json:"id"`
	ParentType  string `json:"parent_type"`
	ParentID    string `json:"parent_id"`
	ParentKey   string `json:"parent_key,omitempty"`
	ParentTitle string `json:"parent_title,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ArchivedAt  string `json:"archived_at,omitempty"`
}

// DraftList is a page of the authenticated user's saved drafts.
type DraftList struct {
	Drafts      []DraftSummary `json:"drafts"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// ListUsers returns visible users.
func ListUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (UserList, error) {
	userPage, err := users(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true), boolPtr(true))
	if err != nil {
		return UserList{}, fmt.Errorf("list users: %w", err)
	}

	summaries := make([]UserSummary, 0, len(userPage.Users.Nodes))
	for _, user := range userPage.Users.Nodes {
		summaries = append(summaries, userSummary(user.UserSummaryFields))
	}

	return UserList{
		Users:       summaries,
		HasNextPage: userPage.Users.PageInfo.HasNextPage,
		EndCursor:   userPage.Users.PageInfo.EndCursor,
	}, nil
}

// GetUserByID returns one User by id.
func GetUserByID(ctx context.Context, graphqlClient graphql.Client, id string) (UserSummary, error) {
	userResult, err := user(ctx, graphqlClient, id)
	if err != nil {
		return UserSummary{}, fmt.Errorf("get user %s: %w", id, err)
	}

	return userSummary(userResult.User.UserSummaryFields), nil
}

// GetViewerUser returns the authenticated User.
func GetViewerUser(ctx context.Context, graphqlClient graphql.Client) (UserSummary, error) {
	userResult, err := viewer(ctx, graphqlClient)
	if err != nil {
		return UserSummary{}, fmt.Errorf("get viewer user: %w", err)
	}

	return userSummary(userResult.Viewer.UserSummaryFields), nil
}

// ListViewerDrafts returns the authenticated user's saved draft metadata.
func ListViewerDrafts(ctx context.Context, graphqlClient graphql.Client, limit int) (DraftList, error) {
	draftPage, err := viewer_drafts(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return DraftList{}, fmt.Errorf("list viewer drafts: %w", err)
	}

	summaries := make([]DraftSummary, 0, len(draftPage.Viewer.Drafts.Nodes))
	for _, draft := range draftPage.Viewer.Drafts.Nodes {
		summaries = append(summaries, draftSummary(draft.DraftSummaryFields))
	}

	return DraftList{
		Drafts:      summaries,
		HasNextPage: draftPage.Viewer.Drafts.PageInfo.HasNextPage,
		EndCursor:   draftPage.Viewer.Drafts.PageInfo.EndCursor,
	}, nil
}

// ListUserAssignedIssues returns issues assigned to one User.
func ListUserAssignedIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	result, err := user_assignedIssues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list user assigned issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.User.AssignedIssues.Nodes))
	for _, issue := range result.User.AssignedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.User.AssignedIssues.PageInfo.HasNextPage,
		EndCursor:   result.User.AssignedIssues.PageInfo.EndCursor,
	}, nil
}

// ListUserCreatedIssues returns issues created by one User.
func ListUserCreatedIssues(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (IssueList, error) {
	result, err := user_createdIssues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list user created issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.User.CreatedIssues.Nodes))
	for _, issue := range result.User.CreatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.User.CreatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.User.CreatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListUserDelegatedIssues returns issues delegated to one User.
func ListUserDelegatedIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	result, err := user_delegatedIssues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list user delegated issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.User.DelegatedIssues.Nodes))
	for _, issue := range result.User.DelegatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.User.DelegatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.User.DelegatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListUserTeamMemberships returns TeamMemberships associated with one User.
func ListUserTeamMemberships(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (TeamMembershipList, error) {
	result, err := user_teamMemberships(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamMembershipList{}, fmt.Errorf("list user team memberships %s: %w", id, err)
	}

	memberships := make([]TeamMembershipSummary, 0, len(result.User.TeamMemberships.Nodes))
	for _, membership := range result.User.TeamMemberships.Nodes {
		memberships = append(memberships, teamMembershipSummary(membership.TeamMembershipSummaryFields))
	}

	return TeamMembershipList{
		Memberships: memberships,
		HasNextPage: result.User.TeamMemberships.PageInfo.HasNextPage,
		EndCursor:   result.User.TeamMemberships.PageInfo.EndCursor,
	}, nil
}

// ListUserTeams returns Teams associated with one User.
func ListUserTeams(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (TeamList, error) {
	result, err := user_teams(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamList{}, fmt.Errorf("list user teams %s: %w", id, err)
	}

	teams := make([]TeamSummary, 0, len(result.User.Teams.Nodes))
	for _, team := range result.User.Teams.Nodes {
		teams = append(teams, teamSummary(team.TeamSummaryFields))
	}

	return TeamList{
		Teams:       teams,
		HasNextPage: result.User.Teams.PageInfo.HasNextPage,
		EndCursor:   result.User.Teams.PageInfo.EndCursor,
	}, nil
}

// ListViewerAssignedIssues returns issues assigned to the authenticated User.
func ListViewerAssignedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	result, err := viewer_assignedIssues(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list viewer assigned issues: %w", err)
	}

	issues := make([]IssueSummary, 0, len(result.Viewer.AssignedIssues.Nodes))
	for _, issue := range result.Viewer.AssignedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Viewer.AssignedIssues.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.AssignedIssues.PageInfo.EndCursor,
	}, nil
}

// ListViewerCreatedIssues returns issues created by the authenticated User.
func ListViewerCreatedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	result, err := viewer_createdIssues(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list viewer created issues: %w", err)
	}

	issues := make([]IssueSummary, 0, len(result.Viewer.CreatedIssues.Nodes))
	for _, issue := range result.Viewer.CreatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Viewer.CreatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.CreatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListViewerDelegatedIssues returns issues delegated to the authenticated User.
func ListViewerDelegatedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	result, err := viewer_delegatedIssues(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list viewer delegated issues: %w", err)
	}

	issues := make([]IssueSummary, 0, len(result.Viewer.DelegatedIssues.Nodes))
	for _, issue := range result.Viewer.DelegatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Viewer.DelegatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.DelegatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListViewerTeamMemberships returns TeamMemberships associated with the authenticated User.
func ListViewerTeamMemberships(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (TeamMembershipList, error) {
	result, err := viewer_teamMemberships(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamMembershipList{}, fmt.Errorf("list viewer team memberships: %w", err)
	}

	memberships := make([]TeamMembershipSummary, 0, len(result.Viewer.TeamMemberships.Nodes))
	for _, membership := range result.Viewer.TeamMemberships.Nodes {
		memberships = append(memberships, teamMembershipSummary(membership.TeamMembershipSummaryFields))
	}

	return TeamMembershipList{
		Memberships: memberships,
		HasNextPage: result.Viewer.TeamMemberships.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.TeamMemberships.PageInfo.EndCursor,
	}, nil
}

// ListViewerTeams returns Teams associated with the authenticated User.
func ListViewerTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	result, err := viewer_teams(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamList{}, fmt.Errorf("list viewer teams: %w", err)
	}

	teams := make([]TeamSummary, 0, len(result.Viewer.Teams.Nodes))
	for _, team := range result.Viewer.Teams.Nodes {
		teams = append(teams, teamSummary(team.TeamSummaryFields))
	}

	return TeamList{
		Teams:       teams,
		HasNextPage: result.Viewer.Teams.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.Teams.PageInfo.EndCursor,
	}, nil
}

func userSummary(user UserSummaryFields) UserSummary {
	return UserSummary{
		ID:          user.Id,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Active:      user.Active,
		Guest:       user.Guest,
		Admin:       user.Admin,
	}
}

func draftSummary(draft DraftSummaryFields) DraftSummary {
	summary := DraftSummary{
		ID:         draft.Id,
		CreatedAt:  draft.CreatedAt,
		UpdatedAt:  draft.UpdatedAt,
		ArchivedAt: stringValue(draft.ArchivedAt),
	}
	switch {
	case draft.Issue != nil:
		summary.ParentType = "issue"
		summary.ParentID = draft.Issue.Id
		summary.ParentKey = draft.Issue.Identifier
		summary.ParentTitle = draft.Issue.Title
	case draft.Project != nil:
		summary.ParentType = "project"
		summary.ParentID = draft.Project.Id
		summary.ParentTitle = draft.Project.Name
	case draft.ProjectUpdate != nil:
		summary.ParentType = "project_update"
		summary.ParentID = draft.ProjectUpdate.Id
	case draft.Initiative != nil:
		summary.ParentType = "initiative"
		summary.ParentID = draft.Initiative.Id
		summary.ParentTitle = draft.Initiative.Name
	case draft.InitiativeUpdate != nil:
		summary.ParentType = "initiative_update"
		summary.ParentID = draft.InitiativeUpdate.Id
	case draft.ParentComment != nil:
		summary.ParentType = "comment"
		summary.ParentID = draft.ParentComment.Id
	case draft.CustomerNeed != nil:
		summary.ParentType = "customer_need"
		summary.ParentID = draft.CustomerNeed.Id
	case draft.Team != nil:
		summary.ParentType = "team"
		summary.ParentID = draft.Team.Id
		summary.ParentKey = draft.Team.Key
		summary.ParentTitle = draft.Team.Name
	default:
		summary.ParentType = "unknown"
	}

	return summary
}
