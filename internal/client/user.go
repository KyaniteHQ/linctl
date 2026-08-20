package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
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
	Users []UserSummary `json:"users"`
	Page
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
	Drafts []DraftSummary `json:"drafts"`
	Page
}

//nolint:lll
type usersNode = gql.XUsersUsersUserConnectionNodesUser

//nolint:lll
type viewerDraftsNode = gql.XViewer_draftsViewerUserDraftsDraftConnectionNodesDraft

//nolint:lll
type userAssignedIssuesNode = gql.XUser_assignedIssuesUserAssignedIssuesIssueConnectionNodesIssue

//nolint:lll
type userCreatedIssuesNode = gql.XUser_createdIssuesUserCreatedIssuesIssueConnectionNodesIssue

//nolint:lll
type userDelegatedIssuesNode = gql.XUser_delegatedIssuesUserDelegatedIssuesIssueConnectionNodesIssue

//nolint:lll
type userTeamMembershipsNode = gql.XUser_teamMembershipsUserTeamMembershipsTeamMembershipConnectionNodesTeamMembership

//nolint:lll
type userTeamsNode = gql.XUser_teamsUserTeamsTeamConnectionNodesTeam

//nolint:lll
type viewerAssignedIssuesNode = gql.XViewer_assignedIssuesViewerUserAssignedIssuesIssueConnectionNodesIssue

//nolint:lll
type viewerCreatedIssuesNode = gql.XViewer_createdIssuesViewerUserCreatedIssuesIssueConnectionNodesIssue

//nolint:lll
type viewerDelegatedIssuesNode = gql.XViewer_delegatedIssuesViewerUserDelegatedIssuesIssueConnectionNodesIssue

//nolint:lll
type viewerTeamMembershipsNode = gql.XViewer_teamMembershipsViewerUserTeamMembershipsTeamMembershipConnectionNodesTeamMembership

//nolint:lll
type viewerTeamsNode = gql.XViewer_teamsViewerUserTeamsTeamConnectionNodesTeam

type usersQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type viewerDraftsQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

type userScopedQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

type viewerQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListUsers returns visible users.
func ListUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (UserList, error) {
	query := usersQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list users", limit, defaultListPageSize,
		query.page,
		usersNodeSummary,
	)
	if err != nil {
		return UserList{}, err
	}

	return UserList{Users: page.Items, Page: page.Page}, nil
}

// GetUserByID returns one User by id.
func GetUserByID(ctx context.Context, graphqlClient graphql.Client, id string) (UserSummary, error) {
	userResult, err := gql.XUser(ctx, graphqlClient, id)
	if err != nil {
		return UserSummary{}, fmt.Errorf("get user %s: %w", id, err)
	}

	return userSummary(userResult.User.UserSummaryFields), nil
}

// GetViewerUser returns the authenticated User.
func GetViewerUser(ctx context.Context, graphqlClient graphql.Client) (UserSummary, error) {
	userResult, err := gql.XViewer(ctx, graphqlClient)
	if err != nil {
		return UserSummary{}, fmt.Errorf("get viewer user: %w", err)
	}

	return userSummary(userResult.Viewer.UserSummaryFields), nil
}

// ListViewerDrafts returns the authenticated user's saved draft metadata.
func ListViewerDrafts(ctx context.Context, graphqlClient graphql.Client, limit int) (DraftList, error) {
	query := viewerDraftsQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list viewer drafts", limit, defaultListPageSize,
		query.page,
		viewerDraftsNodeSummary,
	)
	if err != nil {
		return DraftList{}, err
	}

	return DraftList{Drafts: page.Items, Page: page.Page}, nil
}

// ListUserAssignedIssues returns issues assigned to one User.
func ListUserAssignedIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	query := userScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list user assigned issues "+id, limit, defaultListPageSize,
		query.assignedIssues,
		userAssignedIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, Page: page.Page}, nil
}

// ListUserCreatedIssues returns issues created by one User.
func ListUserCreatedIssues(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (IssueList, error) {
	query := userScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list user created issues "+id, limit, defaultListPageSize,
		query.createdIssues,
		userCreatedIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, Page: page.Page}, nil
}

// ListUserDelegatedIssues returns issues delegated to one User.
func ListUserDelegatedIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	query := userScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list user delegated issues "+id, limit, defaultListPageSize,
		query.delegatedIssues,
		userDelegatedIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, Page: page.Page}, nil
}

// ListUserTeamMemberships returns TeamMemberships associated with one User.
func ListUserTeamMemberships(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (TeamMembershipList, error) {
	query := userScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list user team memberships "+id, limit, defaultListPageSize,
		query.teamMemberships,
		userTeamMembershipsNodeSummary,
	)
	if err != nil {
		return TeamMembershipList{}, err
	}

	return TeamMembershipList{Memberships: page.Items, Page: page.Page}, nil
}

// ListUserTeams returns Teams associated with one User.
func ListUserTeams(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (TeamList, error) {
	query := userScopedQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, err := listConnection(
		"list user teams "+id, limit, defaultListPageSize,
		query.teams,
		userTeamsNodeSummary,
	)
	if err != nil {
		return TeamList{}, err
	}

	return TeamList{Teams: page.Items, Page: page.Page}, nil
}

// ListViewerAssignedIssues returns issues assigned to the authenticated User.
func ListViewerAssignedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	query := viewerQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list viewer assigned issues", limit, defaultListPageSize,
		query.assignedIssues,
		viewerAssignedIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, Page: page.Page}, nil
}

// ListViewerCreatedIssues returns issues created by the authenticated User.
func ListViewerCreatedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	query := viewerQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list viewer created issues", limit, defaultListPageSize,
		query.createdIssues,
		viewerCreatedIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, Page: page.Page}, nil
}

// ListViewerDelegatedIssues returns issues delegated to the authenticated User.
func ListViewerDelegatedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	query := viewerQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list viewer delegated issues", limit, defaultListPageSize,
		query.delegatedIssues,
		viewerDelegatedIssuesNodeSummary,
	)
	if err != nil {
		return IssueList{}, err
	}

	return IssueList{Issues: page.Items, Page: page.Page}, nil
}

// ListViewerTeamMemberships returns TeamMemberships associated with the authenticated User.
func ListViewerTeamMemberships(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (TeamMembershipList, error) {
	query := viewerQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list viewer team memberships", limit, defaultListPageSize,
		query.teamMemberships,
		viewerTeamMembershipsNodeSummary,
	)
	if err != nil {
		return TeamMembershipList{}, err
	}

	return TeamMembershipList{Memberships: page.Items, Page: page.Page}, nil
}

// ListViewerTeams returns Teams associated with the authenticated User.
func ListViewerTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	query := viewerQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list viewer teams", limit, defaultListPageSize,
		query.teams,
		viewerTeamsNodeSummary,
	)
	if err != nil {
		return TeamList{}, err
	}

	return TeamList{Teams: page.Items, Page: page.Page}, nil
}

func (query usersQuery) page(pageSize int, after *string) ([]usersNode, bool, *string, error) {
	result, err := gql.XUsers(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(true), boolPtr(true))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Users.Nodes, result.Users.PageInfo.HasNextPage, result.Users.PageInfo.EndCursor, nil
}

func (query viewerDraftsQuery) page(pageSize int, after *string) ([]viewerDraftsNode, bool, *string, error) {
	result, err := gql.XViewer_drafts(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Viewer.Drafts.Nodes,
		result.Viewer.Drafts.PageInfo.HasNextPage,
		result.Viewer.Drafts.PageInfo.EndCursor,
		nil
}

func (query userScopedQuery) assignedIssues(
	pageSize int,
	after *string,
) ([]userAssignedIssuesNode, bool, *string, error) {
	result, err := gql.XUser_assignedIssues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.User.AssignedIssues.Nodes,
		result.User.AssignedIssues.PageInfo.HasNextPage,
		result.User.AssignedIssues.PageInfo.EndCursor,
		nil
}

func (query userScopedQuery) createdIssues(
	pageSize int,
	after *string,
) ([]userCreatedIssuesNode, bool, *string, error) {
	result, err := gql.XUser_createdIssues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.User.CreatedIssues.Nodes,
		result.User.CreatedIssues.PageInfo.HasNextPage,
		result.User.CreatedIssues.PageInfo.EndCursor,
		nil
}

func (query userScopedQuery) delegatedIssues(
	pageSize int,
	after *string,
) ([]userDelegatedIssuesNode, bool, *string, error) {
	result, err := gql.XUser_delegatedIssues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.User.DelegatedIssues.Nodes,
		result.User.DelegatedIssues.PageInfo.HasNextPage,
		result.User.DelegatedIssues.PageInfo.EndCursor,
		nil
}

func (query userScopedQuery) teamMemberships(
	pageSize int,
	after *string,
) ([]userTeamMembershipsNode, bool, *string, error) {
	result, err := gql.XUser_teamMemberships(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.User.TeamMemberships.Nodes,
		result.User.TeamMemberships.PageInfo.HasNextPage,
		result.User.TeamMemberships.PageInfo.EndCursor,
		nil
}

func (query userScopedQuery) teams(pageSize int, after *string) ([]userTeamsNode, bool, *string, error) {
	result, err := gql.XUser_teams(query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(false))
	if err != nil {
		return nil, false, nil, err
	}

	return result.User.Teams.Nodes,
		result.User.Teams.PageInfo.HasNextPage,
		result.User.Teams.PageInfo.EndCursor,
		nil
}

func (query viewerQuery) assignedIssues(
	pageSize int,
	after *string,
) ([]viewerAssignedIssuesNode, bool, *string, error) {
	result, err := gql.XViewer_assignedIssues(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Viewer.AssignedIssues.Nodes,
		result.Viewer.AssignedIssues.PageInfo.HasNextPage,
		result.Viewer.AssignedIssues.PageInfo.EndCursor,
		nil
}

func (query viewerQuery) createdIssues(
	pageSize int,
	after *string,
) ([]viewerCreatedIssuesNode, bool, *string, error) {
	result, err := gql.XViewer_createdIssues(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Viewer.CreatedIssues.Nodes,
		result.Viewer.CreatedIssues.PageInfo.HasNextPage,
		result.Viewer.CreatedIssues.PageInfo.EndCursor,
		nil
}

func (query viewerQuery) delegatedIssues(
	pageSize int,
	after *string,
) ([]viewerDelegatedIssuesNode, bool, *string, error) {
	result, err := gql.XViewer_delegatedIssues(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Viewer.DelegatedIssues.Nodes,
		result.Viewer.DelegatedIssues.PageInfo.HasNextPage,
		result.Viewer.DelegatedIssues.PageInfo.EndCursor,
		nil
}

func (query viewerQuery) teamMemberships(
	pageSize int,
	after *string,
) ([]viewerTeamMembershipsNode, bool, *string, error) {
	result, err := gql.XViewer_teamMemberships(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Viewer.TeamMemberships.Nodes,
		result.Viewer.TeamMemberships.PageInfo.HasNextPage,
		result.Viewer.TeamMemberships.PageInfo.EndCursor,
		nil
}

func (query viewerQuery) teams(pageSize int, after *string) ([]viewerTeamsNode, bool, *string, error) {
	result, err := gql.XViewer_teams(query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false))
	if err != nil {
		return nil, false, nil, err
	}

	return result.Viewer.Teams.Nodes,
		result.Viewer.Teams.PageInfo.HasNextPage,
		result.Viewer.Teams.PageInfo.EndCursor,
		nil
}

func usersNodeSummary(user usersNode) UserSummary {
	return userSummary(user.UserSummaryFields)
}

func viewerDraftsNodeSummary(draft viewerDraftsNode) DraftSummary {
	return draftSummary(draft.DraftSummaryFields)
}

func userAssignedIssuesNodeSummary(issue userAssignedIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func userCreatedIssuesNodeSummary(issue userCreatedIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func userDelegatedIssuesNodeSummary(issue userDelegatedIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func userTeamMembershipsNodeSummary(membership userTeamMembershipsNode) TeamMembershipSummary {
	return teamMembershipSummary(membership.TeamMembershipSummaryFields)
}

func userTeamsNodeSummary(team userTeamsNode) TeamSummary {
	return teamSummary(team.TeamSummaryFields)
}

func viewerAssignedIssuesNodeSummary(issue viewerAssignedIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func viewerCreatedIssuesNodeSummary(issue viewerCreatedIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func viewerDelegatedIssuesNodeSummary(issue viewerDelegatedIssuesNode) IssueSummary {
	return issueSummaryFromFields(issue.IssueSummaryFields)
}

func viewerTeamMembershipsNodeSummary(membership viewerTeamMembershipsNode) TeamMembershipSummary {
	return teamMembershipSummary(membership.TeamMembershipSummaryFields)
}

func viewerTeamsNodeSummary(team viewerTeamsNode) TeamSummary {
	return teamSummary(team.TeamSummaryFields)
}

func userSummary(user gql.UserSummaryFields) UserSummary {
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

func draftSummary(draft gql.DraftSummaryFields) DraftSummary {
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
