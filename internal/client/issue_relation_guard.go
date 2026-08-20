package client

import (
	"context"
	"fmt"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

func (guard *guardedClient) requireRelationIssues(
	ctx context.Context,
	request IssueRelationCreateRequest,
) (IssueDetail, IssueDetail, error) {
	if len(request.AllowedProjectIDs) == 0 {
		return guard.requirePinnedRelationIssues(ctx, request.IssueID, request.RelatedIssueID)
	}
	allowed, err := guard.allowedProjectSet(ctx, request.AllowedProjectIDs)
	if err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}
	issue, err := guard.requireAllowedRelationIssue(ctx, request.IssueID, allowed)
	if err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}
	related, err := guard.requireAllowedRelationIssue(ctx, request.RelatedIssueID, allowed)
	if err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}

	return issue, related, nil
}

func (guard *guardedClient) requirePinnedRelationIssues(
	ctx context.Context,
	firstID string,
	secondID string,
) (IssueDetail, IssueDetail, error) {
	issue, err := guard.requirePinnedRelationIssue(ctx, firstID)
	if err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}
	related, err := guard.requirePinnedRelationIssue(ctx, secondID)
	if err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}

	return issue, related, nil
}

func (guard *guardedClient) requirePinnedRelationIssue(
	ctx context.Context,
	issueID string,
) (IssueDetail, error) {
	issue, err := guard.requireIssueOnTeam(ctx, issueID)
	if err != nil {
		return IssueDetail{}, err
	}
	if guard.target.Project != nil && issue.Summary.ProjectID != guard.target.Project.ID {
		return IssueDetail{}, guard.projectMismatchError("issue project_id", issue.Summary.ProjectID)
	}

	return issue, nil
}

func (guard *guardedClient) requireAllowedRelationIssue(
	ctx context.Context,
	issueID string,
	allowed map[string]struct{},
) (IssueDetail, error) {
	issue, err := guard.requireIssueOnTeam(ctx, issueID)
	if err != nil {
		return IssueDetail{}, err
	}
	if _, ok := allowed[issue.Summary.ProjectID]; !ok {
		return IssueDetail{}, fmt.Errorf(
			"%w: issue %s project_id=%s is outside the allowed projects",
			ErrTargetMismatch,
			issue.Summary.Identifier,
			issue.Summary.ProjectID,
		)
	}

	return issue, nil
}

func (guard *guardedClient) allowedProjectSet(
	ctx context.Context,
	projectIDs []string,
) (map[string]struct{}, error) {
	allowed := make(map[string]struct{}, len(projectIDs))
	for _, projectID := range projectIDs {
		if projectID == "" {
			return nil, fmt.Errorf("%w: allowed project id is required", ErrWriteInvalid)
		}
		if _, ok := allowed[projectID]; ok {
			continue
		}
		project, err := GetProjectByID(ctx, guard.graphqlClient, projectID)
		if err != nil {
			return nil, err
		}
		if err := guard.requireProjectTeam(project); err != nil {
			return nil, err
		}
		allowed[projectID] = struct{}{}
	}

	return allowed, nil
}

func (guard *guardedClient) existingIssueRelation(
	ctx context.Context,
	issue IssueSummary,
	related IssueSummary,
	relationType string,
) (IssueRelationSummary, bool, error) {
	result, err := gql.XIssue_relations(
		ctx, guard.graphqlClient, issue.ID, intPtr(workflowStatePageSize), nil, boolPtr(true),
	)
	if err != nil {
		return IssueRelationSummary{}, false, err
	}
	for _, node := range result.Issue.Relations.Nodes {
		relation := issueRelationSummary(node.IssueRelationSummaryFields)
		if relation.Type == relationType &&
			relation.IssueID == issue.ID &&
			relation.RelatedIssueID == related.ID {
			return relation, true, nil
		}
	}
	if result.Issue.Relations.PageInfo.HasNextPage {
		return IssueRelationSummary{}, false, fmt.Errorf(
			"%w: issue %s has more than %d relations; cannot reconcile before create",
			ErrWriteInvalid,
			issue.Identifier,
			workflowStatePageSize,
		)
	}

	return IssueRelationSummary{}, false, nil
}

func (guard *guardedClient) writeIssueRelation(
	ctx context.Context,
	relationType string,
	issue IssueDetail,
	related IssueDetail,
) (IssueRelationWriteResult, error) {
	created, err := gql.IssueRelationCreate(ctx, guard.graphqlClient, LinearIssueRelationCreateInput{
		Type:           relationType,
		IssueID:        issue.Summary.ID,
		RelatedIssueID: related.Summary.ID,
	})
	if err != nil {
		return guard.reconcileIssueRelation(
			ctx, issue.Summary, related.Summary, relationType,
			fmt.Errorf("create issue relation: %w", err),
		)
	}
	if !created.IssueRelationCreate.Success {
		return guard.reconcileIssueRelation(
			ctx, issue.Summary, related.Summary, relationType,
			fmt.Errorf("%w: issueRelationCreate returned no relation", ErrMutationFailed),
		)
	}
	summary := issueRelationSummary(created.IssueRelationCreate.IssueRelation.IssueRelationSummaryFields)

	return guard.readIssueRelationResult(ctx, summary, issue.Summary, related.Summary)
}

func (guard *guardedClient) reconcileIssueRelation(
	ctx context.Context,
	issue IssueSummary,
	related IssueSummary,
	relationType string,
	writeErr error,
) (IssueRelationWriteResult, error) {
	existing, found, err := guard.existingIssueRelation(ctx, issue, related, relationType)
	if err != nil || !found {
		return IssueRelationWriteResult{}, writeErr
	}
	result, readErr := guard.readIssueRelationResult(ctx, existing, issue, related)
	if readErr != nil {
		return IssueRelationWriteResult{}, readErr
	}

	return result, nil
}

func (guard *guardedClient) readIssueRelationResult(
	ctx context.Context,
	relation IssueRelationSummary,
	beforeIssue IssueSummary,
	beforeRelated IssueSummary,
) (IssueRelationWriteResult, error) {
	issue, err := GetIssueDetail(ctx, guard.graphqlClient, beforeIssue.ID)
	if err != nil {
		return IssueRelationWriteResult{}, err
	}
	related, err := GetIssueDetail(ctx, guard.graphqlClient, beforeRelated.ID)
	if err != nil {
		return IssueRelationWriteResult{}, err
	}
	readRelation, err := GetIssueRelationByID(ctx, guard.graphqlClient, relation.ID)
	if err != nil {
		return IssueRelationWriteResult{}, err
	}
	result := IssueRelationWriteResult{
		IssueRelationSummary: readRelation,
		Issue:                issue.Summary,
		RelatedIssue:         related.Summary,
	}
	if err := refuseMovedRelationProjects(result, beforeIssue.ProjectID, beforeRelated.ProjectID); err != nil {
		return IssueRelationWriteResult{}, err
	}

	return result, nil
}

func refuseMovedRelationProjects(
	result IssueRelationWriteResult,
	sourceProject string,
	relatedProject string,
) error {
	if result.Issue.ProjectID != sourceProject || result.RelatedIssue.ProjectID != relatedProject {
		return fmt.Errorf(
			"%w: relation changed an issue project; expected %s and %s resolved %s and %s",
			ErrWriteInvalid,
			sourceProject,
			relatedProject,
			result.Issue.ProjectID,
			result.RelatedIssue.ProjectID,
		)
	}

	return nil
}
