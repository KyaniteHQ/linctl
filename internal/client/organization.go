package client

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

//nolint:lll
type organizationLabelsNode = gql.XOrganization_labelsOrganizationLabelsIssueLabelConnectionNodesIssueLabel

//nolint:lll
type organizationProjectLabelsNode = gql.XOrganization_projectLabelsOrganizationProjectLabelsProjectLabelConnectionNodesProjectLabel

//nolint:lll
type organizationTeamsNode = gql.XOrganization_teamsOrganizationTeamsTeamConnectionNodesTeam

//nolint:lll
type organizationUsersNode = gql.XOrganization_usersOrganizationUsersUserConnectionNodesUser

type organizationQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
}

// ListOrganizationLabels returns organization-wide issue labels.
func ListOrganizationLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (LabelList, error) {
	query := organizationQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list organization labels", limit, defaultListPageSize,
		query.labels,
		organizationLabelsNodeSummary,
	)
	if err != nil {
		return LabelList{}, err
	}

	return LabelList{Labels: page.Items, Page: page.Page}, nil
}

// ListOrganizationProjectLabels returns organization-wide project labels.
func ListOrganizationProjectLabels(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (ProjectLabelList, error) {
	query := organizationQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list organization project labels", limit, defaultListPageSize,
		query.projectLabels,
		organizationProjectLabelsNodeSummary,
	)
	if err != nil {
		return ProjectLabelList{}, err
	}

	return ProjectLabelList{
		ProjectLabels: page.Items,
		Page:          page.Page,
	}, nil
}

// ListOrganizationTeams returns teams visible to the authenticated user.
func ListOrganizationTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	query := organizationQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list organization teams", limit, defaultListPageSize,
		query.teams,
		organizationTeamsNodeSummary,
	)
	if err != nil {
		return TeamList{}, err
	}

	return TeamList{Teams: page.Items, Page: page.Page}, nil
}

// ListOrganizationUsers returns active users visible to the authenticated user.
func ListOrganizationUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (UserList, error) {
	query := organizationQuery{ctx: ctx, graphqlClient: graphqlClient}
	page, err := listConnection(
		"list organization users", limit, defaultListPageSize,
		query.users,
		organizationUsersNodeSummary,
	)
	if err != nil {
		return UserList{}, err
	}

	return UserList{Users: page.Items, Page: page.Page}, nil
}

func (query organizationQuery) labels(
	pageSize int,
	after *string,
) ([]organizationLabelsNode, bool, *string, error) {
	result, err := gql.XOrganization_labels(
		query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Organization.Labels.Nodes,
		result.Organization.Labels.PageInfo.HasNextPage,
		result.Organization.Labels.PageInfo.EndCursor,
		nil
}

func (query organizationQuery) projectLabels(
	pageSize int,
	after *string,
) ([]organizationProjectLabelsNode, bool, *string, error) {
	result, err := gql.XOrganization_projectLabels(
		query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Organization.ProjectLabels.Nodes,
		result.Organization.ProjectLabels.PageInfo.HasNextPage,
		result.Organization.ProjectLabels.PageInfo.EndCursor,
		nil
}

func (query organizationQuery) teams(
	pageSize int,
	after *string,
) ([]organizationTeamsNode, bool, *string, error) {
	result, err := gql.XOrganization_teams(
		query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Organization.Teams.Nodes,
		result.Organization.Teams.PageInfo.HasNextPage,
		result.Organization.Teams.PageInfo.EndCursor,
		nil
}

func (query organizationQuery) users(
	pageSize int,
	after *string,
) ([]organizationUsersNode, bool, *string, error) {
	result, err := gql.XOrganization_users(
		query.ctx, query.graphqlClient, intPtr(pageSize), after, boolPtr(false),
	)
	if err != nil {
		return nil, false, nil, err
	}

	return result.Organization.Users.Nodes,
		result.Organization.Users.PageInfo.HasNextPage,
		result.Organization.Users.PageInfo.EndCursor,
		nil
}

func organizationLabelsNodeSummary(label organizationLabelsNode) LabelSummary {
	return labelSummary(label.IssueLabelSummaryFields)
}

func organizationProjectLabelsNodeSummary(label organizationProjectLabelsNode) ProjectLabelSummary {
	return projectLabelSummary(label.ProjectLabelSummaryFields)
}

func organizationTeamsNodeSummary(team organizationTeamsNode) TeamSummary {
	return teamSummary(team.TeamSummaryFields)
}

func organizationUsersNodeSummary(user organizationUsersNode) UserSummary {
	return userSummary(user.UserSummaryFields)
}
