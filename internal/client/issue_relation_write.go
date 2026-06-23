package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

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

// LinearIssueRelationCreateInput is the sparse Linear issueRelationCreate payload linctl supports.
type LinearIssueRelationCreateInput struct {
	Type           string `json:"type"`
	IssueID        string `json:"issueId"`
	RelatedIssueID string `json:"relatedIssueId"`
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

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (IssueRelationSummary, error) {
		resolved, err := requireIssuePair(ctx, graphqlClient, guard, request.IssueID, request.RelatedIssueID)
		if err != nil {
			return IssueRelationSummary{}, err
		}
		issue, related := resolved[0], resolved[1]
		if err := guardBlockingCycle(ctx, graphqlClient, request.Type, issue, related); err != nil {
			return IssueRelationSummary{}, err
		}

		created, err := IssueRelationCreate(ctx, graphqlClient, LinearIssueRelationCreateInput{
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
	})
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

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (string, error) {
		relation, err := GetIssueRelationByID(ctx, graphqlClient, relationID)
		if err != nil {
			return "", err
		}
		if _, err := requireIssuePair(
			ctx, graphqlClient, guard, relation.IssueIdentifier, relation.RelatedIssueIdentifier,
		); err != nil {
			return "", err
		}

		deleted, err := IssueRelationDelete(ctx, graphqlClient, relationID)
		if err != nil {
			return "", fmt.Errorf("delete issue relation %s: %w", relationID, err)
		}
		if !deleted.IssueRelationDelete.Success {
			return "", fmt.Errorf("%w: issueRelationDelete reported no success", ErrMutationFailed)
		}

		return relation.ID, nil
	})
}

// requireIssuePair resolves both endpoints of a relation through the guard,
// confirming each issue belongs to the resolved team before any mutation.
func requireIssuePair(
	ctx context.Context,
	graphqlClient graphql.Client,
	guard writeGuard,
	firstID string,
	secondID string,
) ([2]IssueSummary, error) {
	var resolved [2]IssueSummary
	for index, id := range [2]string{firstID, secondID} {
		summary, err := guard.requireIssue(ctx, graphqlClient, id)
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
func guardBlockingCycle(
	ctx context.Context,
	graphqlClient graphql.Client,
	relationType string,
	issue IssueSummary,
	related IssueSummary,
) error {
	if relationType != "blocks" {
		return nil
	}
	dependencies, err := GetIssueDependencies(ctx, graphqlClient, issue.ID, dependencyCheckLimit)
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

	return nil
}
