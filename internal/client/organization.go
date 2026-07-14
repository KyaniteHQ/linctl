package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// ListOrganizationLabels returns organization-wide issue labels.
func ListOrganizationLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (LabelList, error) {
	result, err := gql.XOrganization_labels(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return LabelList{}, fmt.Errorf("list organization labels: %w", err)
	}

	labels := mapNodes(result.Organization.Labels.Nodes, func(
		label gql.XOrganization_labelsOrganizationLabelsIssueLabelConnectionNodesIssueLabel,
	) LabelSummary {
		return labelSummary(label.IssueLabelSummaryFields)
	})

	return LabelList{
		Labels:      labels,
		HasNextPage: result.Organization.Labels.PageInfo.HasNextPage,
		EndCursor:   result.Organization.Labels.PageInfo.EndCursor,
	}, nil
}

// ListOrganizationProjectLabels returns organization-wide project labels.
func ListOrganizationProjectLabels(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (ProjectLabelList, error) {
	result, err := gql.XOrganization_projectLabels(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return ProjectLabelList{}, fmt.Errorf("list organization project labels: %w", err)
	}

	labels := mapNodes(result.Organization.ProjectLabels.Nodes, func(
		label gql.XOrganization_projectLabelsOrganizationProjectLabelsProjectLabelConnectionNodesProjectLabel,
	) ProjectLabelSummary {
		return projectLabelSummary(label.ProjectLabelSummaryFields)
	})

	return ProjectLabelList{
		ProjectLabels: labels,
		HasNextPage:   result.Organization.ProjectLabels.PageInfo.HasNextPage,
		EndCursor:     result.Organization.ProjectLabels.PageInfo.EndCursor,
	}, nil
}

// ListOrganizationTeams returns teams visible to the authenticated user.
func ListOrganizationTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	result, err := gql.XOrganization_teams(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamList{}, fmt.Errorf("list organization teams: %w", err)
	}

	teams := mapNodes(result.Organization.Teams.Nodes, func(
		team gql.XOrganization_teamsOrganizationTeamsTeamConnectionNodesTeam,
	) TeamSummary {
		return teamSummary(team.TeamSummaryFields)
	})

	return TeamList{
		Teams:       teams,
		HasNextPage: result.Organization.Teams.PageInfo.HasNextPage,
		EndCursor:   result.Organization.Teams.PageInfo.EndCursor,
	}, nil
}

// ListOrganizationUsers returns active users visible to the authenticated user.
func ListOrganizationUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (UserList, error) {
	result, err := gql.XOrganization_users(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return UserList{}, fmt.Errorf("list organization users: %w", err)
	}

	users := mapNodes(result.Organization.Users.Nodes, func(
		user gql.XOrganization_usersOrganizationUsersUserConnectionNodesUser,
	) UserSummary {
		return userSummary(user.UserSummaryFields)
	})

	return UserList{
		Users:       users,
		HasNextPage: result.Organization.Users.PageInfo.HasNextPage,
		EndCursor:   result.Organization.Users.PageInfo.EndCursor,
	}, nil
}
