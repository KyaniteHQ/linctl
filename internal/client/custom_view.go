package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// CustomViewSummary is the compact custom view model used by read-only commands.
type CustomViewSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ModelName   string `json:"model_name"`
	Shared      bool   `json:"shared"`
	Color       string `json:"color,omitempty"`
	SlugID      string `json:"slug_id"`
}

// CustomViewList is a page of custom views.
type CustomViewList struct {
	CustomViews []CustomViewSummary `json:"custom_views"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// CustomViewSubscriberStatus reports whether a custom view has active subscribers.
type CustomViewSubscriberStatus struct {
	ID             string `json:"id"`
	HasSubscribers bool   `json:"has_subscribers"`
}

// CustomViewPreferences is the compact organization preference model for one custom view.
type CustomViewPreferences struct {
	CustomViewID string                      `json:"custom_view_id"`
	ID           string                      `json:"id,omitempty"`
	Type         string                      `json:"type,omitempty"`
	ViewType     string                      `json:"view_type,omitempty"`
	CreatedAt    string                      `json:"created_at,omitempty"`
	UpdatedAt    string                      `json:"updated_at,omitempty"`
	ArchivedAt   string                      `json:"archived_at,omitempty"`
	Values       CustomViewPreferencesValues `json:"values"`
}

// CustomViewPreferencesValues is the compact display-settings model for one custom view.
type CustomViewPreferencesValues struct {
	CustomViewID                string   `json:"custom_view_id,omitempty"`
	Layout                      string   `json:"layout,omitempty"`
	ViewOrdering                string   `json:"view_ordering,omitempty"`
	ViewOrderingDirection       string   `json:"view_ordering_direction,omitempty"`
	IssueGrouping               string   `json:"issue_grouping,omitempty"`
	IssueSubGrouping            string   `json:"issue_sub_grouping,omitempty"`
	ShowCompletedIssues         string   `json:"show_completed_issues,omitempty"`
	ShowArchivedItems           bool     `json:"show_archived_items,omitempty"`
	ShowEmptyGroups             bool     `json:"show_empty_groups,omitempty"`
	HiddenColumns               []string `json:"hidden_columns,omitempty"`
	HiddenRows                  []string `json:"hidden_rows,omitempty"`
	HiddenGroupsList            []string `json:"hidden_groups_list,omitempty"`
	ColumnOrderBoard            []string `json:"column_order_board,omitempty"`
	ColumnOrderList             []string `json:"column_order_list,omitempty"`
	ProjectLayout               string   `json:"project_layout,omitempty"`
	ProjectViewOrdering         string   `json:"project_view_ordering,omitempty"`
	ProjectGrouping             string   `json:"project_grouping,omitempty"`
	ProjectSubGrouping          string   `json:"project_sub_grouping,omitempty"`
	ProjectShowEmptyGroups      string   `json:"project_show_empty_groups,omitempty"`
	ProjectShowEmptySubGroups   string   `json:"project_show_empty_sub_groups,omitempty"`
	HasOrganizationPreferences  bool     `json:"has_organization_preferences,omitempty"`
	HasUserPreferences          bool     `json:"has_user_preferences,omitempty"`
	HasEffectivePreferenceValue bool     `json:"has_effective_preference_value,omitempty"`
}

// ListCustomViews returns visible custom views.
func ListCustomViews(ctx context.Context, graphqlClient graphql.Client, limit int) (CustomViewList, error) {
	result, err := customViews(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return CustomViewList{}, fmt.Errorf("list custom views: %w", err)
	}

	summaries := make([]CustomViewSummary, 0, len(result.CustomViews.Nodes))
	for _, node := range result.CustomViews.Nodes {
		summaries = append(summaries, customViewSummary(node.CustomViewSummaryFields))
	}

	return CustomViewList{
		CustomViews: summaries,
		HasNextPage: result.CustomViews.PageInfo.HasNextPage,
		EndCursor:   result.CustomViews.PageInfo.EndCursor,
	}, nil
}

// GetCustomViewByID returns one custom view by Linear id or slug.
func GetCustomViewByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewSummary, error) {
	result, err := customView(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewSummary{}, fmt.Errorf("get custom view %s: %w", id, err)
	}

	return customViewSummary(result.CustomView.CustomViewSummaryFields), nil
}

// GetCustomViewSubscriberStatus returns whether a custom view has active subscribers.
func GetCustomViewSubscriberStatus(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewSubscriberStatus, error) {
	result, err := customViewHasSubscribers(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewSubscriberStatus{}, fmt.Errorf("get custom view subscribers %s: %w", id, err)
	}

	return CustomViewSubscriberStatus{
		ID:             id,
		HasSubscribers: result.CustomViewHasSubscribers.HasSubscribers,
	}, nil
}

// ListCustomViewInitiatives returns initiatives matching one custom view's initiative filter.
func ListCustomViewInitiatives(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (InitiativeList, error) {
	result, err := customView_initiatives(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return InitiativeList{}, fmt.Errorf("list custom view initiatives %s: %w", id, err)
	}

	initiatives := make([]InitiativeSummary, 0, len(result.CustomView.Initiatives.Nodes))
	for _, node := range result.CustomView.Initiatives.Nodes {
		initiatives = append(initiatives, initiativeSummary(node.InitiativeSummaryFields))
	}

	return InitiativeList{
		Initiatives: initiatives,
		HasNextPage: result.CustomView.Initiatives.PageInfo.HasNextPage,
		EndCursor:   result.CustomView.Initiatives.PageInfo.EndCursor,
	}, nil
}

// ListCustomViewIssues returns issues matching one custom view's issue filter.
func ListCustomViewIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	result, err := customView_issues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list custom view issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.CustomView.Issues.Nodes))
	for _, node := range result.CustomView.Issues.Nodes {
		issues = append(issues, issueSummaryFromFields(node.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.CustomView.Issues.PageInfo.HasNextPage,
		EndCursor:   result.CustomView.Issues.PageInfo.EndCursor,
	}, nil
}

// GetCustomViewOrganizationPreferences returns organization defaults for one custom view.
func GetCustomViewOrganizationPreferences(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewPreferences, error) {
	result, err := customView_organizationViewPreferences(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewPreferences{}, fmt.Errorf("get custom view organization preferences %s: %w", id, err)
	}
	if result.CustomView.OrganizationViewPreferences == nil {
		return CustomViewPreferences{CustomViewID: id}, nil
	}

	fields := result.CustomView.OrganizationViewPreferences.CustomViewPreferencesFields
	return buildCustomViewPreferences(id, fields), nil
}

func buildCustomViewPreferences(id string, fields CustomViewPreferencesFields) CustomViewPreferences {
	return CustomViewPreferences{
		CustomViewID: id,
		ID:           fields.Id,
		Type:         fields.Type,
		ViewType:     fields.ViewType,
		CreatedAt:    fields.CreatedAt,
		UpdatedAt:    fields.UpdatedAt,
		ArchivedAt:   stringValue(fields.ArchivedAt),
		Values:       customViewPreferencesValues(id, fields.Preferences.CustomViewPreferencesValueFields),
	}
}

// GetCustomViewOrganizationPreferenceValues returns organization default display settings for one custom view.
func GetCustomViewOrganizationPreferenceValues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewPreferencesValues, error) {
	result, err := customView_organizationViewPreferences_preferences(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewPreferencesValues{}, fmt.Errorf(
			"get custom view organization preference values %s: %w",
			id,
			err,
		)
	}
	if result.CustomView.OrganizationViewPreferences == nil {
		return CustomViewPreferencesValues{CustomViewID: id}, nil
	}

	values := customViewPreferencesValues(
		id,
		result.CustomView.OrganizationViewPreferences.Preferences.CustomViewPreferencesValueFields,
	)
	values.HasOrganizationPreferences = true
	return values, nil
}

// ListCustomViewProjects returns projects matching one custom view's project filter.
func ListCustomViewProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectList, error) {
	result, err := customView_projects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return ProjectList{}, fmt.Errorf("list custom view projects %s: %w", id, err)
	}

	projects := make([]ProjectSummary, 0, len(result.CustomView.Projects.Nodes))
	for _, node := range result.CustomView.Projects.Nodes {
		projects = append(projects, projectSummaryFromFields(node.ProjectSummaryFields))
	}

	return ProjectList{
		Projects:    projects,
		HasNextPage: result.CustomView.Projects.PageInfo.HasNextPage,
		EndCursor:   result.CustomView.Projects.PageInfo.EndCursor,
	}, nil
}

// GetCustomViewUserPreferences returns current-user preferences for one custom view.
func GetCustomViewUserPreferences(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewPreferences, error) {
	result, err := customView_userViewPreferences(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewPreferences{}, fmt.Errorf("get custom view user preferences %s: %w", id, err)
	}
	if result.CustomView.UserViewPreferences == nil {
		return CustomViewPreferences{CustomViewID: id}, nil
	}

	fields := result.CustomView.UserViewPreferences.CustomViewPreferencesFields
	return buildCustomViewPreferences(id, fields), nil
}

// GetCustomViewUserPreferenceValues returns current-user display settings for one custom view.
func GetCustomViewUserPreferenceValues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewPreferencesValues, error) {
	result, err := customView_userViewPreferences_preferences(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewPreferencesValues{}, fmt.Errorf("get custom view user preference values %s: %w", id, err)
	}
	if result.CustomView.UserViewPreferences == nil {
		return CustomViewPreferencesValues{CustomViewID: id}, nil
	}

	values := customViewPreferencesValues(
		id,
		result.CustomView.UserViewPreferences.Preferences.CustomViewPreferencesValueFields,
	)
	values.HasUserPreferences = true
	return values, nil
}

// GetCustomViewPreferenceValues returns effective display settings for one custom view.
func GetCustomViewPreferenceValues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (CustomViewPreferencesValues, error) {
	result, err := customView_viewPreferencesValues(ctx, graphqlClient, id)
	if err != nil {
		return CustomViewPreferencesValues{}, fmt.Errorf("get custom view preference values %s: %w", id, err)
	}

	values := customViewPreferencesValues(id, result.CustomView.ViewPreferencesValues.CustomViewPreferencesValueFields)
	values.HasEffectivePreferenceValue = true
	return values, nil
}

func customViewSummary(fields CustomViewSummaryFields) CustomViewSummary {
	return CustomViewSummary{
		ID:          fields.Id,
		Name:        fields.Name,
		Description: stringValue(fields.Description),
		ModelName:   fields.ModelName,
		Shared:      fields.Shared,
		Color:       stringValue(fields.Color),
		SlugID:      fields.SlugId,
	}
}

func customViewPreferencesValues(
	id string,
	fields CustomViewPreferencesValueFields,
) CustomViewPreferencesValues {
	return CustomViewPreferencesValues{
		CustomViewID:              id,
		Layout:                    stringValue(fields.Layout),
		ViewOrdering:              stringValue(fields.ViewOrdering),
		ViewOrderingDirection:     stringValue(fields.ViewOrderingDirection),
		IssueGrouping:             stringValue(fields.IssueGrouping),
		IssueSubGrouping:          stringValue(fields.IssueSubGrouping),
		ShowCompletedIssues:       stringValue(fields.ShowCompletedIssues),
		ShowArchivedItems:         boolValue(fields.ShowArchivedItems),
		ShowEmptyGroups:           boolValue(fields.ShowEmptyGroups),
		HiddenColumns:             fields.HiddenColumns,
		HiddenRows:                fields.HiddenRows,
		HiddenGroupsList:          fields.HiddenGroupsList,
		ColumnOrderBoard:          fields.ColumnOrderBoard,
		ColumnOrderList:           fields.ColumnOrderList,
		ProjectLayout:             stringValue(fields.ProjectLayout),
		ProjectViewOrdering:       stringValue(fields.ProjectViewOrdering),
		ProjectGrouping:           stringValue(fields.ProjectGrouping),
		ProjectSubGrouping:        stringValue(fields.ProjectSubGrouping),
		ProjectShowEmptyGroups:    stringValue(fields.ProjectShowEmptyGroups),
		ProjectShowEmptySubGroups: stringValue(fields.ProjectShowEmptySubGroups),
	}
}
