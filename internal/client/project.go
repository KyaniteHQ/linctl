package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// ProjectSummary is the compact project model used by project commands.
type ProjectSummary struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	SlugID      string        `json:"slug_id"`
	URL         string        `json:"url"`
	Priority    int           `json:"priority"`
	Status      ProjectStatus `json:"status"`
	Lead        string        `json:"lead,omitempty"`
	Teams       []ProjectTeam `json:"teams"`
}

// ProjectStatus is the compact project lifecycle status.
type ProjectStatus struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ProjectTeam is a project-associated team.
type ProjectTeam struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// ProjectList is a page of projects scoped to a team.
type ProjectList struct {
	Projects    []ProjectSummary `json:"projects"`
	HasNextPage bool             `json:"has_next_page"`
	EndCursor   *string          `json:"end_cursor,omitempty"`
}

// ProjectMember is a project member.
type ProjectMember struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

// ProjectMemberList is a page of project members.
type ProjectMemberList struct {
	ProjectID   string          `json:"project_id"`
	ProjectName string          `json:"project_name"`
	Members     []ProjectMember `json:"members"`
	HasNextPage bool            `json:"has_next_page"`
	EndCursor   *string         `json:"end_cursor,omitempty"`
}

// ProjectUpdateSummary is one project status update.
type ProjectUpdateSummary struct {
	ID          string `json:"id"`
	Body        string `json:"body,omitempty"`
	Health      string `json:"health"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	URL         string `json:"url"`
	ProjectID   string `json:"project_id,omitempty"`
	ProjectName string `json:"project_name,omitempty"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// ProjectUpdateList is a page of project status updates.
type ProjectUpdateList struct {
	ProjectID   string                 `json:"project_id"`
	ProjectName string                 `json:"project_name"`
	Updates     []ProjectUpdateSummary `json:"updates"`
	HasNextPage bool                   `json:"has_next_page"`
	EndCursor   *string                `json:"end_cursor,omitempty"`
}

// ProjectFilterSuggestion is an AI-generated project filter suggestion.
type ProjectFilterSuggestion struct {
	Filter json.RawMessage `json:"filter,omitempty"`
	LogID  string          `json:"log_id,omitempty"`
}

// ProjectUpdateCommentList is a page of body-free Comments associated with one ProjectUpdate.
type ProjectUpdateCommentList struct {
	ProjectUpdateID string                   `json:"project_update_id"`
	Comments        []CommentMetadataSummary `json:"comments"`
	HasNextPage     bool                     `json:"has_next_page"`
	EndCursor       *string                  `json:"end_cursor,omitempty"`
}

// ProjectAttachmentList is a page of Attachments associated with one Project.
type ProjectAttachmentList struct {
	ProjectID   string              `json:"project_id"`
	ProjectName string              `json:"project_name"`
	Attachments []AttachmentSummary `json:"attachments"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// ProjectDocumentList is a page of Documents associated with one Project.
type ProjectDocumentList struct {
	ProjectID   string            `json:"project_id"`
	ProjectName string            `json:"project_name"`
	Documents   []DocumentSummary `json:"documents"`
	HasNextPage bool              `json:"has_next_page"`
	EndCursor   *string           `json:"end_cursor,omitempty"`
}

// ProjectExternalLinkList is a page of external links associated with one Project.
type ProjectExternalLinkList struct {
	ProjectID   string                      `json:"project_id"`
	ProjectName string                      `json:"project_name"`
	Links       []EntityExternalLinkSummary `json:"links"`
	HasNextPage bool                        `json:"has_next_page"`
	EndCursor   *string                     `json:"end_cursor,omitempty"`
}

// ProjectHistorySummary is the compact project history model used by read-only commands.
type ProjectHistorySummary struct {
	ID         string          `json:"id"`
	ProjectID  string          `json:"project_id"`
	EntryCount int             `json:"entry_count"`
	Entries    json.RawMessage `json:"entries"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
	ArchivedAt string          `json:"archived_at,omitempty"`
}

// ProjectHistoryList is a page of Linear project history records.
type ProjectHistoryList struct {
	ProjectID   string                  `json:"project_id"`
	ProjectName string                  `json:"project_name"`
	History     []ProjectHistorySummary `json:"history"`
	HasNextPage bool                    `json:"has_next_page"`
	EndCursor   *string                 `json:"end_cursor,omitempty"`
}

// ProjectInitiativeToProjectList is a page of initiative associations for one Project.
type ProjectInitiativeToProjectList struct {
	ProjectID    string                       `json:"project_id"`
	ProjectName  string                       `json:"project_name"`
	Associations []InitiativeToProjectSummary `json:"associations"`
	HasNextPage  bool                         `json:"has_next_page"`
	EndCursor    *string                      `json:"end_cursor,omitempty"`
}

// ProjectInitiativeList is a page of Initiatives associated with one Project.
type ProjectInitiativeList struct {
	ProjectID   string              `json:"project_id"`
	ProjectName string              `json:"project_name"`
	Initiatives []InitiativeSummary `json:"initiatives"`
	HasNextPage bool                `json:"has_next_page"`
	EndCursor   *string             `json:"end_cursor,omitempty"`
}

// ProjectIssueList is a page of Issues associated with one Project.
type ProjectIssueList struct {
	ProjectID   string         `json:"project_id"`
	ProjectName string         `json:"project_name"`
	Issues      []IssueSummary `json:"issues"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// ProjectCommentList is a page of body-free Comments associated with one Project.
type ProjectCommentList struct {
	ProjectID   string                   `json:"project_id"`
	ProjectName string                   `json:"project_name"`
	Comments    []CommentMetadataSummary `json:"comments"`
	HasNextPage bool                     `json:"has_next_page"`
	EndCursor   *string                  `json:"end_cursor,omitempty"`
}

// ProjectProjectLabelList is a page of ProjectLabels associated with one Project.
type ProjectProjectLabelList struct {
	ProjectID     string                `json:"project_id"`
	ProjectName   string                `json:"project_name"`
	ProjectLabels []ProjectLabelSummary `json:"project_labels"`
	HasNextPage   bool                  `json:"has_next_page"`
	EndCursor     *string               `json:"end_cursor,omitempty"`
}

// ProjectCustomerNeedList is a page of customer needs associated with one Project.
type ProjectCustomerNeedList struct {
	ProjectID   string                `json:"project_id"`
	ProjectName string                `json:"project_name"`
	Needs       []CustomerNeedSummary `json:"customer_needs"`
	HasNextPage bool                  `json:"has_next_page"`
	EndCursor   *string               `json:"end_cursor,omitempty"`
}

// ProjectProjectRelationList is a page of project relations associated with one Project.
type ProjectProjectRelationList struct {
	ProjectID   string                   `json:"project_id"`
	ProjectName string                   `json:"project_name"`
	Relations   []ProjectRelationSummary `json:"relations"`
	HasNextPage bool                     `json:"has_next_page"`
	EndCursor   *string                  `json:"end_cursor,omitempty"`
}

// ProjectTeamList is a page of Teams associated with one Project.
type ProjectTeamList struct {
	ProjectID   string        `json:"project_id"`
	ProjectName string        `json:"project_name"`
	Teams       []TeamSummary `json:"teams"`
	HasNextPage bool          `json:"has_next_page"`
	EndCursor   *string       `json:"end_cursor,omitempty"`
}

// ListProjectsByTeam returns projects scoped to a resolved team.
func ListProjectsByTeam(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamID string,
	limit int,
) (ProjectList, error) {
	projects, err := Projects(ctx, graphqlClient, teamID, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectList{}, fmt.Errorf("list projects: %w", err)
	}

	summaries := make([]ProjectSummary, 0, len(projects.Team.Projects.Nodes))
	for _, project := range projects.Team.Projects.Nodes {
		summaries = append(summaries, projectSummaryFromFields(project.ProjectSummaryFields))
	}

	return ProjectList{
		Projects:    summaries,
		HasNextPage: projects.Team.Projects.PageInfo.HasNextPage,
		EndCursor:   projects.Team.Projects.PageInfo.EndCursor,
	}, nil
}

// GetProjectByID returns a project by Linear id or slug.
func GetProjectByID(ctx context.Context, graphqlClient graphql.Client, id string) (ProjectSummary, error) {
	projectResult, err := project(ctx, graphqlClient, id)
	if err != nil {
		return ProjectSummary{}, fmt.Errorf("get project %s: %w", id, err)
	}

	return projectSummaryFromFields(projectResult.Project.ProjectSummaryFields), nil
}

// ListProjectAttachments returns Attachments associated with one Project.
func ListProjectAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectAttachmentList, error) {
	result, err := project_attachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectAttachmentList{}, fmt.Errorf("list project attachments %s: %w", id, err)
	}

	attachments := make([]AttachmentSummary, 0, len(result.Project.Attachments.Nodes))
	for _, node := range result.Project.Attachments.Nodes {
		attachments = append(attachments, projectAttachmentSummary(node.ProjectAttachmentSummaryFields))
	}

	return ProjectAttachmentList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Attachments: attachments,
		HasNextPage: result.Project.Attachments.PageInfo.HasNextPage,
		EndCursor:   result.Project.Attachments.PageInfo.EndCursor,
	}, nil
}

// ListProjectDocuments returns Documents associated with one Project.
func ListProjectDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectDocumentList, error) {
	result, err := project_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectDocumentList{}, fmt.Errorf("list project documents %s: %w", id, err)
	}

	documents := make([]DocumentSummary, 0, len(result.Project.Documents.Nodes))
	for _, node := range result.Project.Documents.Nodes {
		documents = append(documents, documentSummary(node.DocumentSummaryFields))
	}

	return ProjectDocumentList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Documents:   documents,
		HasNextPage: result.Project.Documents.PageInfo.HasNextPage,
		EndCursor:   result.Project.Documents.PageInfo.EndCursor,
	}, nil
}

// ListProjectExternalLinks returns external links associated with one Project.
func ListProjectExternalLinks(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectExternalLinkList, error) {
	result, err := project_externalLinks(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectExternalLinkList{}, fmt.Errorf("list project external links %s: %w", id, err)
	}

	links := make([]EntityExternalLinkSummary, 0, len(result.Project.ExternalLinks.Nodes))
	for _, node := range result.Project.ExternalLinks.Nodes {
		links = append(links, entityExternalLinkSummary(node.EntityExternalLinkSummaryFields))
	}

	return ProjectExternalLinkList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Links:       links,
		HasNextPage: result.Project.ExternalLinks.PageInfo.HasNextPage,
		EndCursor:   result.Project.ExternalLinks.PageInfo.EndCursor,
	}, nil
}

// ListProjectHistory returns history records associated with one Project.
func ListProjectHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectHistoryList, error) {
	result, err := project_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectHistoryList{}, fmt.Errorf("list project history %s: %w", id, err)
	}

	history := make([]ProjectHistorySummary, 0, len(result.Project.History.Nodes))
	for _, node := range result.Project.History.Nodes {
		history = append(history, projectHistorySummary(node.ProjectHistorySummaryFields))
	}

	return ProjectHistoryList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		History:     history,
		HasNextPage: result.Project.History.PageInfo.HasNextPage,
		EndCursor:   result.Project.History.PageInfo.EndCursor,
	}, nil
}

// ListProjectInitiativeToProjects returns Initiative-to-Project associations for one Project.
func ListProjectInitiativeToProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectInitiativeToProjectList, error) {
	result, err := project_initiativeToProjects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectInitiativeToProjectList{}, fmt.Errorf("list project initiative associations %s: %w", id, err)
	}

	associations := make([]InitiativeToProjectSummary, 0, len(result.Project.InitiativeToProjects.Nodes))
	for _, node := range result.Project.InitiativeToProjects.Nodes {
		associations = append(associations, initiativeToProjectSummary(node.InitiativeToProjectSummaryFields))
	}

	return ProjectInitiativeToProjectList{
		ProjectID:    result.Project.Id,
		ProjectName:  result.Project.Name,
		Associations: associations,
		HasNextPage:  result.Project.InitiativeToProjects.PageInfo.HasNextPage,
		EndCursor:    result.Project.InitiativeToProjects.PageInfo.EndCursor,
	}, nil
}

// ListProjectInitiatives returns Initiatives associated with one Project.
func ListProjectInitiatives(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectInitiativeList, error) {
	result, err := project_initiatives(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectInitiativeList{}, fmt.Errorf("list project initiatives %s: %w", id, err)
	}

	initiatives := make([]InitiativeSummary, 0, len(result.Project.Initiatives.Nodes))
	for _, node := range result.Project.Initiatives.Nodes {
		initiatives = append(initiatives, initiativeSummary(node.InitiativeSummaryFields))
	}

	return ProjectInitiativeList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Initiatives: initiatives,
		HasNextPage: result.Project.Initiatives.PageInfo.HasNextPage,
		EndCursor:   result.Project.Initiatives.PageInfo.EndCursor,
	}, nil
}

// ListProjectInverseRelations returns inverse project relations associated with one Project.
func ListProjectInverseRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectProjectRelationList, error) {
	result, err := project_inverseRelations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectProjectRelationList{}, fmt.Errorf("list project inverse relations %s: %w", id, err)
	}

	relations := make([]ProjectRelationSummary, 0, len(result.Project.InverseRelations.Nodes))
	for _, node := range result.Project.InverseRelations.Nodes {
		relations = append(relations, projectRelationSummary(node.ProjectRelationSummaryFields))
	}

	return ProjectProjectRelationList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Relations:   relations,
		HasNextPage: result.Project.InverseRelations.PageInfo.HasNextPage,
		EndCursor:   result.Project.InverseRelations.PageInfo.EndCursor,
	}, nil
}

// ListProjectIssues returns Issues associated with one Project.
func ListProjectIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectIssueList, error) {
	result, err := project_issues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectIssueList{}, fmt.Errorf("list project issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.Project.Issues.Nodes))
	for _, node := range result.Project.Issues.Nodes {
		issues = append(issues, issueSummaryFromFields(node.IssueSummaryFields))
	}

	return ProjectIssueList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Issues:      issues,
		HasNextPage: result.Project.Issues.PageInfo.HasNextPage,
		EndCursor:   result.Project.Issues.PageInfo.EndCursor,
	}, nil
}

// ListProjectComments returns body-free comments associated with one Project.
func ListProjectComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectCommentList, error) {
	result, err := project_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectCommentList{}, fmt.Errorf("list project comments %s: %w", id, err)
	}

	comments := make([]CommentMetadataSummary, 0, len(result.Project.Comments.Nodes))
	for _, node := range result.Project.Comments.Nodes {
		comments = append(comments, commentMetadataSummary(node.CommentMetadataFields))
	}

	return ProjectCommentList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Comments:    comments,
		HasNextPage: result.Project.Comments.PageInfo.HasNextPage,
		EndCursor:   result.Project.Comments.PageInfo.EndCursor,
	}, nil
}

// ListLabelsForProject returns ProjectLabels associated with one Project.
func ListLabelsForProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectProjectLabelList, error) {
	result, err := project_labels(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectProjectLabelList{}, fmt.Errorf("list project labels %s: %w", id, err)
	}

	labels := make([]ProjectLabelSummary, 0, len(result.Project.Labels.Nodes))
	for _, node := range result.Project.Labels.Nodes {
		labels = append(labels, projectLabelSummary(node.ProjectLabelSummaryFields))
	}

	return ProjectProjectLabelList{
		ProjectID:     result.Project.Id,
		ProjectName:   result.Project.Name,
		ProjectLabels: labels,
		HasNextPage:   result.Project.Labels.PageInfo.HasNextPage,
		EndCursor:     result.Project.Labels.PageInfo.EndCursor,
	}, nil
}

// ListProjectNeeds returns customer needs associated with one Project.
func ListProjectNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectCustomerNeedList, error) {
	result, err := project_needs(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectCustomerNeedList{}, fmt.Errorf("list project customer needs %s: %w", id, err)
	}

	needs := make([]CustomerNeedSummary, 0, len(result.Project.Needs.Nodes))
	for _, node := range result.Project.Needs.Nodes {
		needs = append(needs, customerNeedSummary(node.CustomerNeedSummaryFields))
	}

	return ProjectCustomerNeedList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Needs:       needs,
		HasNextPage: result.Project.Needs.PageInfo.HasNextPage,
		EndCursor:   result.Project.Needs.PageInfo.EndCursor,
	}, nil
}

// ListProjectRelationsForProject returns project relations associated with one Project.
func ListProjectRelationsForProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectProjectRelationList, error) {
	result, err := project_relations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectProjectRelationList{}, fmt.Errorf("list project relations %s: %w", id, err)
	}

	relations := make([]ProjectRelationSummary, 0, len(result.Project.Relations.Nodes))
	for _, node := range result.Project.Relations.Nodes {
		relations = append(relations, projectRelationSummary(node.ProjectRelationSummaryFields))
	}

	return ProjectProjectRelationList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Relations:   relations,
		HasNextPage: result.Project.Relations.PageInfo.HasNextPage,
		EndCursor:   result.Project.Relations.PageInfo.EndCursor,
	}, nil
}

// ListProjectTeams returns Teams associated with one Project.
func ListProjectTeams(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectTeamList, error) {
	result, err := project_teams(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectTeamList{}, fmt.Errorf("list project teams %s: %w", id, err)
	}

	teams := make([]TeamSummary, 0, len(result.Project.Teams.Nodes))
	for _, node := range result.Project.Teams.Nodes {
		teams = append(teams, teamSummary(node.TeamSummaryFields))
	}

	return ProjectTeamList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		Teams:       teams,
		HasNextPage: result.Project.Teams.PageInfo.HasNextPage,
		EndCursor:   result.Project.Teams.PageInfo.EndCursor,
	}, nil
}

// ListProjectMembers returns members for one project.
func ListProjectMembers(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectMemberList, error) {
	project, err := project_members(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectMemberList{}, fmt.Errorf("list project members %s: %w", id, err)
	}

	members := make([]ProjectMember, 0, len(project.Project.Members.Nodes))
	for _, member := range project.Project.Members.Nodes {
		members = append(members, ProjectMember{
			ID:          member.Id,
			Name:        member.Name,
			DisplayName: member.DisplayName,
			Email:       member.Email,
		})
	}

	return ProjectMemberList{
		ProjectID:   project.Project.Id,
		ProjectName: project.Project.Name,
		Members:     members,
		HasNextPage: project.Project.Members.PageInfo.HasNextPage,
		EndCursor:   project.Project.Members.PageInfo.EndCursor,
	}, nil
}

// ListProjectUpdates returns status updates for one project.
func ListProjectUpdates(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectUpdateList, error) {
	project, err := project_projectUpdates(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectUpdateList{}, fmt.Errorf("list project updates %s: %w", id, err)
	}

	updates := make([]ProjectUpdateSummary, 0, len(project.Project.ProjectUpdates.Nodes))
	for _, update := range project.Project.ProjectUpdates.Nodes {
		updates = append(updates, projectScopedProjectUpdateSummary(update))
	}

	return ProjectUpdateList{
		ProjectID:   project.Project.Id,
		ProjectName: project.Project.Name,
		Updates:     updates,
		HasNextPage: project.Project.ProjectUpdates.PageInfo.HasNextPage,
		EndCursor:   project.Project.ProjectUpdates.PageInfo.EndCursor,
	}, nil
}

// GetProjectFilterSuggestion returns a JSON project filter suggestion for a prompt.
func GetProjectFilterSuggestion(
	ctx context.Context,
	graphqlClient graphql.Client,
	prompt string,
	teamID string,
) (ProjectFilterSuggestion, error) {
	suggestion, err := projectFilterSuggestion(ctx, graphqlClient, prompt, optionalString(teamID))
	if err != nil {
		return ProjectFilterSuggestion{}, fmt.Errorf("get project filter suggestion: %w", err)
	}

	filter := json.RawMessage(nil)
	if suggestion.ProjectFilterSuggestion.Filter != nil {
		filter = *suggestion.ProjectFilterSuggestion.Filter
	}

	return ProjectFilterSuggestion{
		Filter: filter,
		LogID:  stringValue(suggestion.ProjectFilterSuggestion.LogId),
	}, nil
}

// ListAllProjectUpdates returns visible project status updates across projects.
func ListAllProjectUpdates(ctx context.Context, graphqlClient graphql.Client, limit int) (ProjectUpdateList, error) {
	updatesResponse, err := projectUpdates(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectUpdateList{}, fmt.Errorf("list project updates: %w", err)
	}

	updates := make([]ProjectUpdateSummary, 0, len(updatesResponse.ProjectUpdates.Nodes))
	for _, update := range updatesResponse.ProjectUpdates.Nodes {
		updates = append(updates, projectUpdateSummary(update.TopLevelProjectUpdateSummaryFields))
	}

	return ProjectUpdateList{
		Updates:     updates,
		HasNextPage: updatesResponse.ProjectUpdates.PageInfo.HasNextPage,
		EndCursor:   updatesResponse.ProjectUpdates.PageInfo.EndCursor,
	}, nil
}

// GetProjectUpdateByID returns one project update by Linear id.
func GetProjectUpdateByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (ProjectUpdateSummary, error) {
	update, err := projectUpdate(ctx, graphqlClient, id)
	if err != nil {
		return ProjectUpdateSummary{}, fmt.Errorf("get project update %s: %w", id, err)
	}

	return projectUpdateSummary(update.ProjectUpdate.TopLevelProjectUpdateSummaryFields), nil
}

// ListProjectUpdateComments returns body-free comments associated with one ProjectUpdate.
func ListProjectUpdateComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectUpdateCommentList, error) {
	result, err := projectUpdate_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectUpdateCommentList{}, fmt.Errorf("list project update comments %s: %w", id, err)
	}

	comments := make([]CommentMetadataSummary, 0, len(result.ProjectUpdate.Comments.Nodes))
	for _, node := range result.ProjectUpdate.Comments.Nodes {
		comments = append(comments, commentMetadataSummary(node.CommentMetadataFields))
	}

	return ProjectUpdateCommentList{
		ProjectUpdateID: result.ProjectUpdate.Id,
		Comments:        comments,
		HasNextPage:     result.ProjectUpdate.Comments.PageInfo.HasNextPage,
		EndCursor:       result.ProjectUpdate.Comments.PageInfo.EndCursor,
	}, nil
}

func projectScopedProjectUpdateSummary(
	update project_projectUpdatesProjectProjectUpdatesProjectUpdateConnectionNodesProjectUpdate,
) ProjectUpdateSummary {
	return ProjectUpdateSummary{
		ID:          update.Id,
		Health:      string(update.Health),
		CreatedAt:   update.CreatedAt,
		UpdatedAt:   update.UpdatedAt,
		URL:         update.Url,
		UserID:      update.User.Id,
		Name:        update.User.Name,
		DisplayName: update.User.DisplayName,
	}
}

func projectUpdateSummary(update TopLevelProjectUpdateSummaryFields) ProjectUpdateSummary {
	return ProjectUpdateSummary{
		ID:          update.Id,
		Body:        update.Body,
		Health:      string(update.Health),
		CreatedAt:   update.CreatedAt,
		UpdatedAt:   update.UpdatedAt,
		URL:         update.Url,
		ProjectID:   update.Project.Id,
		ProjectName: update.Project.Name,
		UserID:      update.User.Id,
		Name:        update.User.Name,
		DisplayName: update.User.DisplayName,
	}
}

func projectHistorySummary(fields ProjectHistorySummaryFields) ProjectHistorySummary {
	entryCount := 0
	var entries []json.RawMessage
	if err := json.Unmarshal(fields.Entries, &entries); err == nil {
		entryCount = len(entries)
	}

	return ProjectHistorySummary{
		ID:         fields.Id,
		ProjectID:  fields.Project.Id,
		EntryCount: entryCount,
		Entries:    fields.Entries,
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
	}
}

func projectAttachmentSummary(fields ProjectAttachmentSummaryFields) AttachmentSummary {
	return AttachmentSummary{
		ID:         fields.Id,
		Title:      fields.Title,
		Subtitle:   stringValue(fields.Subtitle),
		URL:        fields.Url,
		SourceType: stringValue(fields.SourceType),
	}
}

func projectSummaryFromFields(project ProjectSummaryFields) ProjectSummary {
	lead := ""
	if project.Lead != nil {
		lead = project.Lead.DisplayName
	}

	teams := make([]ProjectTeam, 0, len(project.Teams.Nodes))
	for _, team := range project.Teams.Nodes {
		teams = append(teams, ProjectTeam{
			ID:   team.Id,
			Key:  team.Key,
			Name: team.Name,
		})
	}

	return ProjectSummary{
		ID:          project.Id,
		Name:        project.Name,
		Description: project.Description,
		SlugID:      project.SlugId,
		URL:         project.Url,
		Priority:    project.Priority,
		Status: ProjectStatus{
			ID:   project.Status.Id,
			Name: project.Status.Name,
			Type: string(project.Status.Type),
		},
		Lead:  lead,
		Teams: teams,
	}
}
