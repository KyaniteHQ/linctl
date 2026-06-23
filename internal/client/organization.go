package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ListOrganizationLabels returns organization-wide issue labels.
func ListOrganizationLabels(ctx context.Context, graphqlClient graphql.Client, limit int) (LabelList, error) {
	result, err := organization_labels(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return LabelList{}, fmt.Errorf("list organization labels: %w", err)
	}

	labels := make([]LabelSummary, 0, len(result.Organization.Labels.Nodes))
	for _, label := range result.Organization.Labels.Nodes {
		labels = append(labels, labelSummary(label.IssueLabelSummaryFields))
	}

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
	result, err := organization_projectLabels(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return ProjectLabelList{}, fmt.Errorf("list organization project labels: %w", err)
	}

	labels := make([]ProjectLabelSummary, 0, len(result.Organization.ProjectLabels.Nodes))
	for _, label := range result.Organization.ProjectLabels.Nodes {
		labels = append(labels, projectLabelSummary(label.ProjectLabelSummaryFields))
	}

	return ProjectLabelList{
		ProjectLabels: labels,
		HasNextPage:   result.Organization.ProjectLabels.PageInfo.HasNextPage,
		EndCursor:     result.Organization.ProjectLabels.PageInfo.EndCursor,
	}, nil
}

// ListOrganizationTeams returns teams visible to the authenticated user.
func ListOrganizationTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	result, err := organization_teams(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamList{}, fmt.Errorf("list organization teams: %w", err)
	}

	teams := make([]TeamSummary, 0, len(result.Organization.Teams.Nodes))
	for _, team := range result.Organization.Teams.Nodes {
		teams = append(teams, teamSummary(team.TeamSummaryFields))
	}

	return TeamList{
		Teams:       teams,
		HasNextPage: result.Organization.Teams.PageInfo.HasNextPage,
		EndCursor:   result.Organization.Teams.PageInfo.EndCursor,
	}, nil
}

// ListOrganizationUsers returns active users visible to the authenticated user.
func ListOrganizationUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (UserList, error) {
	result, err := organization_users(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return UserList{}, fmt.Errorf("list organization users: %w", err)
	}

	users := make([]UserSummary, 0, len(result.Organization.Users.Nodes))
	for _, user := range result.Organization.Users.Nodes {
		users = append(users, userSummary(user.UserSummaryFields))
	}

	return UserList{
		Users:       users,
		HasNextPage: result.Organization.Users.PageInfo.HasNextPage,
		EndCursor:   result.Organization.Users.PageInfo.EndCursor,
	}, nil
}
