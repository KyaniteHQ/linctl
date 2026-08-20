package client

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

//nolint:lll
type projectAttachmentsNode = gql.XProject_attachmentsProjectAttachmentsProjectAttachmentConnectionNodesProjectAttachment

//nolint:lll
type projectDocumentsNode = gql.XProject_documentsProjectDocumentsDocumentConnectionNodesDocument

//nolint:lll
type projectExternalLinksNode = gql.XProject_externalLinksProjectExternalLinksEntityExternalLinkConnectionNodesEntityExternalLink

//nolint:lll
type projectHistoryNode = gql.XProject_historyProjectHistoryProjectHistoryConnectionNodesProjectHistory

//nolint:lll
type initiativeToProjectNode = gql.XProject_initiativeToProjectsProjectInitiativeToProjectsInitiativeToProjectConnectionNodesInitiativeToProject

//nolint:lll
type projectInitiativesNode = gql.XProject_initiativesProjectInitiativesInitiativeConnectionNodesInitiative

//nolint:lll
type projectInverseRelationsNode = gql.XProject_inverseRelationsProjectInverseRelationsProjectRelationConnectionNodesProjectRelation

//nolint:lll
type projectIssuesNode = gql.XProject_issuesProjectIssuesIssueConnectionNodesIssue

//nolint:lll
type projectCommentsNode = gql.XProject_commentsProjectCommentsCommentConnectionNodesComment

//nolint:lll
type projectLabelsForProjectNode = gql.XProject_labelsProjectLabelsProjectLabelConnectionNodesProjectLabel

//nolint:lll
type projectNeedsNode = gql.XProject_needsProjectNeedsCustomerNeedConnectionNodesCustomerNeed

//nolint:lll
type projectRelationsForProjectNode = gql.XProject_relationsProjectRelationsProjectRelationConnectionNodesProjectRelation

//nolint:lll
type projectTeamsNode = gql.XProject_teamsProjectTeamsTeamConnectionNodesTeam

//nolint:lll
type projectMembersNode = gql.XProject_membersProjectMembersUserConnectionNodesUser

type projectChildQuery struct {
	ctx           context.Context
	graphqlClient graphql.Client
	id            string
}

// ListProjectAttachments returns Attachments associated with one Project.
func ListProjectAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectAttachmentList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project attachments "+id, limit, defaultListPageSize,
		query.attachments,
		projectAttachmentNodeSummary,
	)
	if err != nil {
		return ProjectAttachmentList{}, err
	}

	return ProjectAttachmentList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Attachments: page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectDocuments returns Documents associated with one Project.
func ListProjectDocuments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectDocumentList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project documents "+id, limit, defaultListPageSize,
		query.documents,
		projectDocumentNodeSummary,
	)
	if err != nil {
		return ProjectDocumentList{}, err
	}

	return ProjectDocumentList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Documents:   page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectExternalLinks returns external links associated with one Project.
func ListProjectExternalLinks(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectExternalLinkList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project external links "+id, limit, defaultListPageSize,
		query.externalLinks,
		projectExternalLinkNodeSummary,
	)
	if err != nil {
		return ProjectExternalLinkList{}, err
	}

	return ProjectExternalLinkList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Links:       page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectHistory returns history records associated with one Project.
func ListProjectHistory(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectHistoryList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project history "+id, limit, defaultListPageSize,
		query.history,
		projectHistoryNodeSummary,
	)
	if err != nil {
		return ProjectHistoryList{}, err
	}

	return ProjectHistoryList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		History:     page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectInitiativeToProjects returns Initiative-to-Project associations for one Project.
func ListProjectInitiativeToProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectInitiativeToProjectList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project initiative associations "+id, limit, defaultListPageSize,
		query.initiativeToProjects,
		projectInitiativeToProjectNodeSummary,
	)
	if err != nil {
		return ProjectInitiativeToProjectList{}, err
	}

	return ProjectInitiativeToProjectList{
		ProjectID:    parent.projectID,
		ProjectName:  parent.projectName,
		Associations: page.Items,
		Page:         page.Page,
	}, nil
}

// ListProjectInitiatives returns Initiatives associated with one Project.
func ListProjectInitiatives(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectInitiativeList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project initiatives "+id, limit, defaultListPageSize,
		query.initiatives,
		projectInitiativeNodeSummary,
	)
	if err != nil {
		return ProjectInitiativeList{}, err
	}

	return ProjectInitiativeList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Initiatives: page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectInverseRelations returns inverse project relations associated with one Project.
func ListProjectInverseRelations(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectProjectRelationList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project inverse relations "+id, limit, defaultListPageSize,
		query.inverseRelations,
		projectInverseRelationNodeSummary,
	)
	if err != nil {
		return ProjectProjectRelationList{}, err
	}

	return ProjectProjectRelationList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Relations:   page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectIssues returns Issues associated with one Project.
func ListProjectIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectIssueList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project issues "+id, limit, defaultListPageSize,
		query.issues,
		projectIssueNodeSummary,
	)
	if err != nil {
		return ProjectIssueList{}, err
	}

	return ProjectIssueList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Issues:      page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectComments returns body-free comments associated with one Project.
func ListProjectComments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectCommentList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project comments "+id, limit, defaultListPageSize,
		query.comments,
		projectCommentNodeSummary,
	)
	if err != nil {
		return ProjectCommentList{}, err
	}

	return ProjectCommentList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Comments:    page.Items,
		Page:        page.Page,
	}, nil
}

// ListLabelsForProject returns ProjectLabels associated with one Project.
func ListLabelsForProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectProjectLabelList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project labels "+id, limit, defaultListPageSize,
		query.labels,
		projectLabelForProjectNodeSummary,
	)
	if err != nil {
		return ProjectProjectLabelList{}, err
	}

	return ProjectProjectLabelList{
		ProjectID:     parent.projectID,
		ProjectName:   parent.projectName,
		ProjectLabels: page.Items,
		Page:          page.Page,
	}, nil
}

// ListProjectNeeds returns customer needs associated with one Project.
func ListProjectNeeds(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectCustomerNeedList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project customer needs "+id, limit, defaultListPageSize,
		query.needs,
		projectNeedNodeSummary,
	)
	if err != nil {
		return ProjectCustomerNeedList{}, err
	}

	return ProjectCustomerNeedList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Needs:       page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectRelationsForProject returns project relations associated with one Project.
func ListProjectRelationsForProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectProjectRelationList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project relations "+id, limit, defaultListPageSize,
		query.relations,
		projectRelationForProjectNodeSummary,
	)
	if err != nil {
		return ProjectProjectRelationList{}, err
	}

	return ProjectProjectRelationList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Relations:   page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectTeams returns Teams associated with one Project.
func ListProjectTeams(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectTeamList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project teams "+id, limit, defaultListPageSize,
		query.teams,
		projectTeamNodeSummary,
	)
	if err != nil {
		return ProjectTeamList{}, err
	}

	return ProjectTeamList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Teams:       page.Items,
		Page:        page.Page,
	}, nil
}

// ListProjectMembers returns members for one project.
func ListProjectMembers(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectMemberList, error) {
	query := &projectChildQuery{ctx: ctx, graphqlClient: graphqlClient, id: id}
	page, parent, err := listConnectionWithParent(
		"list project members "+id, limit, defaultListPageSize,
		query.members,
		projectMemberNodeSummary,
	)
	if err != nil {
		return ProjectMemberList{}, err
	}

	return ProjectMemberList{
		ProjectID:   parent.projectID,
		ProjectName: parent.projectName,
		Members:     page.Items,
		Page:        page.Page,
	}, nil
}

func (query *projectChildQuery) attachments(
	pageSize int,
	after *string,
) ([]projectAttachmentsNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_attachments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Attachments.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Attachments.PageInfo.HasNextPage,
		result.Project.Attachments.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) documents(
	pageSize int,
	after *string,
) ([]projectDocumentsNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_documents(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Documents.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Documents.PageInfo.HasNextPage,
		result.Project.Documents.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) externalLinks(
	pageSize int,
	after *string,
) ([]projectExternalLinksNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_externalLinks(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.ExternalLinks.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.ExternalLinks.PageInfo.HasNextPage,
		result.Project.ExternalLinks.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) history(
	pageSize int,
	after *string,
) ([]projectHistoryNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_history(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.History.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.History.PageInfo.HasNextPage,
		result.Project.History.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) initiativeToProjects(
	pageSize int,
	after *string,
) ([]initiativeToProjectNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_initiativeToProjects(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.InitiativeToProjects.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.InitiativeToProjects.PageInfo.HasNextPage,
		result.Project.InitiativeToProjects.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) initiatives(
	pageSize int,
	after *string,
) ([]projectInitiativesNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_initiatives(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Initiatives.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Initiatives.PageInfo.HasNextPage,
		result.Project.Initiatives.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) inverseRelations(
	pageSize int,
	after *string,
) ([]projectInverseRelationsNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_inverseRelations(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.InverseRelations.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.InverseRelations.PageInfo.HasNextPage,
		result.Project.InverseRelations.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) issues(
	pageSize int,
	after *string,
) ([]projectIssuesNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_issues(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Issues.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Issues.PageInfo.HasNextPage,
		result.Project.Issues.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) comments(
	pageSize int,
	after *string,
) ([]projectCommentsNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_comments(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Comments.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Comments.PageInfo.HasNextPage,
		result.Project.Comments.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) labels(
	pageSize int,
	after *string,
) ([]projectLabelsForProjectNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_labels(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Labels.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Labels.PageInfo.HasNextPage,
		result.Project.Labels.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) needs(
	pageSize int,
	after *string,
) ([]projectNeedsNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_needs(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Needs.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Needs.PageInfo.HasNextPage,
		result.Project.Needs.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) relations(
	pageSize int,
	after *string,
) ([]projectRelationsForProjectNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_relations(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Relations.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Relations.PageInfo.HasNextPage,
		result.Project.Relations.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) teams(
	pageSize int,
	after *string,
) ([]projectTeamsNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_teams(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Teams.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Teams.PageInfo.HasNextPage,
		result.Project.Teams.PageInfo.EndCursor,
		nil
}

func (query *projectChildQuery) members(
	pageSize int,
	after *string,
) ([]projectMembersNode, projectParent, bool, *string, error) {
	result, err := gql.XProject_members(
		query.ctx, query.graphqlClient, query.id, intPtr(pageSize), after, boolPtr(true),
	)
	if err != nil {
		return nil, projectParent{}, false, nil, err
	}

	return result.Project.Members.Nodes,
		projectParent{projectID: result.Project.Id, projectName: result.Project.Name},
		result.Project.Members.PageInfo.HasNextPage,
		result.Project.Members.PageInfo.EndCursor,
		nil
}

func projectAttachmentNodeSummary(node projectAttachmentsNode) AttachmentSummary {
	return projectAttachmentSummary(node.ProjectAttachmentSummaryFields)
}

func projectDocumentNodeSummary(node projectDocumentsNode) DocumentSummary {
	return documentSummary(node.DocumentSummaryFields)
}

func projectExternalLinkNodeSummary(node projectExternalLinksNode) EntityExternalLinkSummary {
	return entityExternalLinkSummary(node.EntityExternalLinkSummaryFields)
}

func projectHistoryNodeSummary(node projectHistoryNode) ProjectHistorySummary {
	return projectHistorySummary(node.ProjectHistorySummaryFields)
}

func projectInitiativeToProjectNodeSummary(node initiativeToProjectNode) InitiativeToProjectSummary {
	return initiativeToProjectSummary(node.InitiativeToProjectSummaryFields)
}

func projectInitiativeNodeSummary(node projectInitiativesNode) InitiativeSummary {
	return initiativeSummary(node.InitiativeSummaryFields)
}

func projectInverseRelationNodeSummary(node projectInverseRelationsNode) ProjectRelationSummary {
	return projectRelationSummary(node.ProjectRelationSummaryFields)
}

func projectIssueNodeSummary(node projectIssuesNode) IssueSummary {
	return issueSummaryFromFields(node.IssueSummaryFields)
}

func projectCommentNodeSummary(node projectCommentsNode) CommentMetadataSummary {
	return commentMetadataSummary(node.CommentMetadataFields)
}

func projectLabelForProjectNodeSummary(node projectLabelsForProjectNode) ProjectLabelSummary {
	return projectLabelSummary(node.ProjectLabelSummaryFields)
}

func projectNeedNodeSummary(node projectNeedsNode) CustomerNeedSummary {
	return customerNeedSummary(node.CustomerNeedSummaryFields)
}

func projectRelationForProjectNodeSummary(node projectRelationsForProjectNode) ProjectRelationSummary {
	return projectRelationSummary(node.ProjectRelationSummaryFields)
}

func projectTeamNodeSummary(node projectTeamsNode) TeamSummary {
	return teamSummary(node.TeamSummaryFields)
}

func projectMemberNodeSummary(member projectMembersNode) ProjectMember {
	return ProjectMember{
		ID:          member.Id,
		Name:        member.Name,
		DisplayName: member.DisplayName,
		Email:       member.Email,
	}
}

func projectHistorySummary(fields gql.ProjectHistorySummaryFields) ProjectHistorySummary {
	return ProjectHistorySummary{
		ID:         fields.Id,
		ProjectID:  fields.Project.Id,
		EntryCount: countJSONArrayEntries(fields.Entries),
		Entries:    fields.Entries,
		CreatedAt:  fields.CreatedAt,
		UpdatedAt:  fields.UpdatedAt,
		ArchivedAt: stringValue(fields.ArchivedAt),
	}
}

func projectAttachmentSummary(fields gql.ProjectAttachmentSummaryFields) AttachmentSummary {
	return AttachmentSummary{
		ID:         fields.Id,
		Title:      fields.Title,
		Subtitle:   stringValue(fields.Subtitle),
		URL:        fields.Url,
		SourceType: stringValue(fields.SourceType),
	}
}
