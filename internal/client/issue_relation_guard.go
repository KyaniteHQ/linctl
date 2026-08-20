package client

import (
	"context"
	"fmt"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

const issueRelationReconcilePageSize = 50

func (guard *guardedClient) requireRelationIssues(
	ctx context.Context,
	request IssueRelationCreateRequest,
) (IssueDetail, IssueDetail, error) {
	allowed, err := guard.relationAllowedProjects(ctx, request.AllowedProjectIDs)
	if err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}
	issue, err := guard.requireIssueOnTeam(ctx, request.IssueID)
	if err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}
	related, err := guard.requireIssueOnTeam(ctx, request.RelatedIssueID)
	if err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}
	if issue.Summary.ProjectID != related.Summary.ProjectID && len(request.AllowedProjectIDs) == 0 {
		return IssueDetail{}, IssueDetail{}, fmt.Errorf(
			"%w: relating across projects needs --allowed-project for each project",
			ErrWriteInvalid,
		)
	}
	if err := guard.requireIssueInAllowedProjects(issue, allowed); err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}
	if err := guard.requireIssueInAllowedProjects(related, allowed); err != nil {
		return IssueDetail{}, IssueDetail{}, err
	}

	return issue, related, nil
}

func (guard *guardedClient) relationAllowedProjects(
	ctx context.Context,
	extraIDs []string,
) (map[string]struct{}, error) {
	allowed, err := guard.allowedProjectSet(ctx, extraIDs)
	if err != nil {
		return nil, err
	}
	if guard.target.Project != nil {
		allowed[guard.target.Project.ID] = struct{}{}
	}

	return allowed, nil
}

func (guard *guardedClient) requireIssueInAllowedProjects(
	issue IssueDetail,
	allowed map[string]struct{},
) error {
	if len(allowed) == 0 {
		return nil
	}
	if _, ok := allowed[issue.Summary.ProjectID]; !ok {
		return fmt.Errorf(
			"%w: issue %s project_id=%s is outside the allowed projects",
			ErrTargetMismatch,
			issue.Summary.Identifier,
			issue.Summary.ProjectID,
		)
	}

	return nil
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
		ctx, guard.graphqlClient, issue.ID, intPtr(issueRelationReconcilePageSize), nil, boolPtr(true),
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
			issueRelationReconcilePageSize,
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
	if err != nil {
		return IssueRelationWriteResult{}, err
	}
	if !found {
		return applyMutationRetryClass(
			IssueRelationCreateRetryClass(), IssueRelationWriteResult{}, false, writeErr,
		)
	}
	result, readErr := guard.readIssueRelationResult(ctx, existing, issue, related)
	if readErr != nil {
		return IssueRelationWriteResult{}, readErr
	}

	return applyMutationRetryClass(
		IssueRelationCreateRetryClass(), result, true, writeErr,
	)
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
