package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// issueRelationTypes is the schema-aligned set of IssueRelationType values.
var issueRelationTypes = map[string]bool{
	"blocks":    true,
	"duplicate": true,
	"related":   true,
	"similar":   true,
}

// dependencyCheckLimit bounds the blocked-by scan used for the cycle pre-check.
const dependencyCheckLimit = 50

// IssueRelationCreateRequest describes a guarded issue-relation create.
type IssueRelationCreateRequest struct {
	IssueID        string
	RelatedIssueID string
	Type           string
}

// CreateIssueRelation links two issues after resolving and comparing the pinned
// target for both endpoints. Each issue must belong to the resolved team. For
// blocks relations it refuses to close a direct cycle.
func CreateIssueRelation(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueRelationCreateRequest,
) (IssueRelationSummary, error) {
	if err := validateIssueRelationCreateRequest(request); err != nil {
		return IssueRelationSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueRelationSummary{}, err
	}

	return guard.createIssueRelation(ctx, request)
}

func (guard *guardedClient) createIssueRelation(
	ctx context.Context,
	request IssueRelationCreateRequest,
) (IssueRelationSummary, error) {
	resolved, err := guard.requireIssuePair(ctx, request.IssueID, request.RelatedIssueID)
	if err != nil {
		return IssueRelationSummary{}, err
	}
	issue, related := resolved[0], resolved[1]
	if err := guard.guardBlockingCycle(ctx, request.Type, issue, related); err != nil {
		return IssueRelationSummary{}, err
	}

	created, err := gql.IssueRelationCreate(ctx, guard.graphqlClient, LinearIssueRelationCreateInput{
		Type:           request.Type,
		IssueID:        issue.ID,
		RelatedIssueID: related.ID,
	})
	if err != nil {
		return IssueRelationSummary{}, fmt.Errorf("create issue relation: %w", err)
	}
	if !created.IssueRelationCreate.Success {
		return IssueRelationSummary{}, fmt.Errorf("%w: issueRelationCreate returned no relation", ErrMutationFailed)
	}

	return issueRelationSummary(created.IssueRelationCreate.IssueRelation.IssueRelationSummaryFields), nil
}

// DeleteIssueRelation removes an existing relation after resolving the relation
// and comparing the pinned target for both linked issues.
func DeleteIssueRelation(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	relationID string,
) (string, error) {
	if relationID == "" {
		return "", fmt.Errorf("%w: relation id is required", ErrWriteInvalid)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return "", err
	}

	return guard.deleteIssueRelation(ctx, relationID)
}

func (guard *guardedClient) deleteIssueRelation(ctx context.Context, relationID string) (string, error) {
	relation, err := GetIssueRelationByID(ctx, guard.graphqlClient, relationID)
	if err != nil {
		return "", err
	}
	if _, err := guard.requireIssuePair(
		ctx, relation.IssueIdentifier, relation.RelatedIssueIdentifier,
	); err != nil {
		return "", err
	}

	deleted, err := gql.IssueRelationDelete(ctx, guard.graphqlClient, relationID)
	if err != nil {
		return "", fmt.Errorf("delete issue relation %s: %w", relationID, err)
	}
	if !deleted.IssueRelationDelete.Success {
		return "", fmt.Errorf("%w: issueRelationDelete reported no success", ErrMutationFailed)
	}

	return relation.ID, nil
}

// requireIssuePair resolves both endpoints of a relation through the guard,
// confirming each issue belongs to the resolved team before any mutation.
func (guard *guardedClient) requireIssuePair(
	ctx context.Context,
	firstID string,
	secondID string,
) ([2]IssueSummary, error) {
	var resolved [2]IssueSummary
	for index, id := range [2]string{firstID, secondID} {
		summary, err := guard.requireIssue(ctx, id)
		if err != nil {
			return resolved, err
		}
		resolved[index] = summary
	}

	return resolved, nil
}

func validateIssueRelationCreateRequest(request IssueRelationCreateRequest) error {
	if request.IssueID == "" || request.RelatedIssueID == "" {
		return fmt.Errorf("%w: issue id and related issue id are required", ErrWriteInvalid)
	}
	if request.IssueID == request.RelatedIssueID {
		return fmt.Errorf("%w: an issue cannot relate to itself", ErrWriteInvalid)
	}
	if !issueRelationTypes[request.Type] {
		return fmt.Errorf(
			"%w: unknown relation type %q: use blocks/duplicate/related/similar",
			ErrWriteInvalid,
			request.Type,
		)
	}

	return nil
}

// guardBlockingCycle refuses a blocks relation that would close a direct cycle:
// when the related issue already blocks the issue, adding issue->blocks->related
// makes them block each other. Non-blocks relation types are always allowed.
func (guard *guardedClient) guardBlockingCycle(
	ctx context.Context,
	relationType string,
	issue IssueSummary,
	related IssueSummary,
) error {
	if relationType != "blocks" {
		return nil
	}
	dependencies, err := GetIssueDependencies(ctx, guard.graphqlClient, issue.ID, dependencyCheckLimit)
	if err != nil {
		return err
	}
	for _, blocker := range dependencies.BlockedBy {
		if blocker.ID == related.ID {
			return fmt.Errorf(
				"%w: %s already blocks %s; the reverse relation would create a cycle",
				ErrWriteInvalid,
				related.Identifier,
				issue.Identifier,
			)
		}
	}
	if dependencies.HasNextPage {
		return fmt.Errorf(
			"%w: issue %s has more than %d relations; cannot verify the blocks relation would not close a cycle",
			ErrWriteInvalid,
			issue.Identifier,
			dependencyCheckLimit,
		)
	}

	return nil
}
