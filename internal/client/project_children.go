package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// ListProjectAttachments returns Attachments associated with one Project.
func ListProjectAttachments(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectAttachmentList, error) {
	result, err := gql.XProject_attachments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectAttachmentList{}, fmt.Errorf("list project attachments %s: %w", id, err)
	}

	attachments := mapNodes(result.Project.Attachments.Nodes, func(
		node gql.XProject_attachmentsProjectAttachmentsProjectAttachmentConnectionNodesProjectAttachment,
	) AttachmentSummary {
		return projectAttachmentSummary(node.ProjectAttachmentSummaryFields)
	})

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
	result, err := gql.XProject_documents(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectDocumentList{}, fmt.Errorf("list project documents %s: %w", id, err)
	}

	documents := mapNodes(result.Project.Documents.Nodes, func(
		node gql.XProject_documentsProjectDocumentsDocumentConnectionNodesDocument,
	) DocumentSummary {
		return documentSummary(node.DocumentSummaryFields)
	})

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
	result, err := gql.XProject_externalLinks(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectExternalLinkList{}, fmt.Errorf("list project external links %s: %w", id, err)
	}

	links := mapNodes(result.Project.ExternalLinks.Nodes, func(
		node gql.XProject_externalLinksProjectExternalLinksEntityExternalLinkConnectionNodesEntityExternalLink,
	) EntityExternalLinkSummary {
		return entityExternalLinkSummary(node.EntityExternalLinkSummaryFields)
	})

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
	result, err := gql.XProject_history(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectHistoryList{}, fmt.Errorf("list project history %s: %w", id, err)
	}

	history := mapNodes(result.Project.History.Nodes, func(
		node gql.XProject_historyProjectHistoryProjectHistoryConnectionNodesProjectHistory,
	) ProjectHistorySummary {
		return projectHistorySummary(node.ProjectHistorySummaryFields)
	})

	return ProjectHistoryList{
		ProjectID:   result.Project.Id,
		ProjectName: result.Project.Name,
		History:     history,
		HasNextPage: result.Project.History.PageInfo.HasNextPage,
		EndCursor:   result.Project.History.PageInfo.EndCursor,
	}, nil
}

//nolint:lll
type initiativeToProjectNode = gql.XProject_initiativeToProjectsProjectInitiativeToProjectsInitiativeToProjectConnectionNodesInitiativeToProject

// ListProjectInitiativeToProjects returns Initiative-to-Project associations for one Project.
func ListProjectInitiativeToProjects(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (ProjectInitiativeToProjectList, error) {
	result, err := gql.XProject_initiativeToProjects(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectInitiativeToProjectList{}, fmt.Errorf("list project initiative associations %s: %w", id, err)
	}

	associations := mapNodes(result.Project.InitiativeToProjects.Nodes,
		func(node initiativeToProjectNode) InitiativeToProjectSummary {
			return initiativeToProjectSummary(node.InitiativeToProjectSummaryFields)
		})

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
	result, err := gql.XProject_initiatives(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectInitiativeList{}, fmt.Errorf("list project initiatives %s: %w", id, err)
	}

	initiatives := mapNodes(result.Project.Initiatives.Nodes, func(
		node gql.XProject_initiativesProjectInitiativesInitiativeConnectionNodesInitiative,
	) InitiativeSummary {
		return initiativeSummary(node.InitiativeSummaryFields)
	})

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
	result, err := gql.XProject_inverseRelations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectProjectRelationList{}, fmt.Errorf("list project inverse relations %s: %w", id, err)
	}

	relations := mapNodes(result.Project.InverseRelations.Nodes, func(
		node gql.XProject_inverseRelationsProjectInverseRelationsProjectRelationConnectionNodesProjectRelation,
	) ProjectRelationSummary {
		return projectRelationSummary(node.ProjectRelationSummaryFields)
	})

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
	result, err := gql.XProject_issues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectIssueList{}, fmt.Errorf("list project issues %s: %w", id, err)
	}

	issues := mapNodes(result.Project.Issues.Nodes, func(
		node gql.XProject_issuesProjectIssuesIssueConnectionNodesIssue,
	) IssueSummary {
		return issueSummaryFromFields(node.IssueSummaryFields)
	})

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
	result, err := gql.XProject_comments(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectCommentList{}, fmt.Errorf("list project comments %s: %w", id, err)
	}

	comments := mapNodes(result.Project.Comments.Nodes, func(
		node gql.XProject_commentsProjectCommentsCommentConnectionNodesComment,
	) CommentMetadataSummary {
		return commentMetadataSummary(node.CommentMetadataFields)
	})

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
	result, err := gql.XProject_labels(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectProjectLabelList{}, fmt.Errorf("list project labels %s: %w", id, err)
	}

	labels := mapNodes(result.Project.Labels.Nodes, func(
		node gql.XProject_labelsProjectLabelsProjectLabelConnectionNodesProjectLabel,
	) ProjectLabelSummary {
		return projectLabelSummary(node.ProjectLabelSummaryFields)
	})

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
	result, err := gql.XProject_needs(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectCustomerNeedList{}, fmt.Errorf("list project customer needs %s: %w", id, err)
	}

	needs := mapNodes(result.Project.Needs.Nodes, func(
		node gql.XProject_needsProjectNeedsCustomerNeedConnectionNodesCustomerNeed,
	) CustomerNeedSummary {
		return customerNeedSummary(node.CustomerNeedSummaryFields)
	})

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
	result, err := gql.XProject_relations(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectProjectRelationList{}, fmt.Errorf("list project relations %s: %w", id, err)
	}

	relations := mapNodes(result.Project.Relations.Nodes, func(
		node gql.XProject_relationsProjectRelationsProjectRelationConnectionNodesProjectRelation,
	) ProjectRelationSummary {
		return projectRelationSummary(node.ProjectRelationSummaryFields)
	})

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
	result, err := gql.XProject_teams(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectTeamList{}, fmt.Errorf("list project teams %s: %w", id, err)
	}

	teams := mapNodes(result.Project.Teams.Nodes, func(
		node gql.XProject_teamsProjectTeamsTeamConnectionNodesTeam,
	) TeamSummary {
		return teamSummary(node.TeamSummaryFields)
	})

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
	project, err := gql.XProject_members(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return ProjectMemberList{}, fmt.Errorf("list project members %s: %w", id, err)
	}

	members := mapNodes(project.Project.Members.Nodes, func(
		member gql.XProject_membersProjectMembersUserConnectionNodesUser,
	) ProjectMember {
		return ProjectMember{
			ID:          member.Id,
			Name:        member.Name,
			DisplayName: member.DisplayName,
			Email:       member.Email,
		}
	})

	return ProjectMemberList{
		ProjectID:   project.Project.Id,
		ProjectName: project.Project.Name,
		Members:     members,
		HasNextPage: project.Project.Members.PageInfo.HasNextPage,
		EndCursor:   project.Project.Members.PageInfo.EndCursor,
	}, nil
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
