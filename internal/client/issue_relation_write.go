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
	IssueID           string
	RelatedIssueID    string
	Type              string
	AllowedProjectIDs []string
}

// IssueRelationDeleteRequest describes a guarded issue-relation delete.
type IssueRelationDeleteRequest struct {
	RelationID        string
	AllowedProjectIDs []string
}

// IssueRelationWriteResult is the relation plus both endpoints after readback.
type IssueRelationWriteResult struct {
	IssueRelationSummary
	Issue        IssueSummary `json:"issue"`
	RelatedIssue IssueSummary `json:"related_issue"`
}

// CreateIssueRelation links two issues after resolving and comparing the pinned
// target for both endpoints. Each issue must belong to the resolved team. For
// blocks relations it refuses to close a direct cycle. AllowedProjectIDs names
// every project the caller permits; it does not widen the pin implicitly.
func CreateIssueRelation(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueRelationCreateRequest,
) (IssueRelationWriteResult, error) {
	if err := validateIssueRelationCreateRequest(request); err != nil {
		return IssueRelationWriteResult{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueRelationWriteResult{}, err
	}

	return guard.createIssueRelation(ctx, request)
}

func (guard *guardedClient) createIssueRelation(
	ctx context.Context,
	request IssueRelationCreateRequest,
) (IssueRelationWriteResult, error) {
	issue, related, err := guard.requireRelationIssues(
		ctx, request.IssueID, request.RelatedIssueID, request.AllowedProjectIDs,
	)
	if err != nil {
		return IssueRelationWriteResult{}, err
	}
	if issue.Summary.ID == related.Summary.ID {
		return IssueRelationWriteResult{}, fmt.Errorf("%w: an issue cannot relate to itself", ErrWriteInvalid)
	}
	if err := guard.guardBlockingCycle(ctx, request.Type, issue.Summary, related.Summary); err != nil {
		return IssueRelationWriteResult{}, err
	}
	if existing, found, existingErr := guard.existingIssueRelation(
		ctx, issue.Summary, related.Summary, request.Type,
	); existingErr != nil {
		return IssueRelationWriteResult{}, existingErr
	} else if found {
		return guard.readIssueRelationResult(ctx, existing, issue.Summary, related.Summary)
	}

	return guard.writeIssueRelation(ctx, request.Type, issue, related)
}

// DeleteIssueRelation removes an existing relation after resolving the relation
// and comparing the same organization, team, and allowed-project boundary as
// CreateIssueRelation for both linked issues.
func DeleteIssueRelation(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueRelationDeleteRequest,
) (string, error) {
	if request.RelationID == "" {
		return "", requiredFieldError("relation id")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return "", err
	}

	return guard.deleteIssueRelation(ctx, request)
}

func (guard *guardedClient) deleteIssueRelation(
	ctx context.Context,
	request IssueRelationDeleteRequest,
) (string, error) {
	relation, err := GetIssueRelationByID(ctx, guard.graphqlClient, request.RelationID)
	if err != nil {
		return "", err
	}
	if _, _, err := guard.requireRelationIssues(
		ctx, relation.IssueID, relation.RelatedIssueID, request.AllowedProjectIDs,
	); err != nil {
		return "", err
	}

	deleted, err := gql.IssueRelationDelete(ctx, guard.graphqlClient, request.RelationID)
	if err != nil {
		return "", fmt.Errorf("delete issue relation %s: %w", request.RelationID, err)
	}
	if err := mutationSuccess(deleted.IssueRelationDelete.Success, "issueRelationDelete"); err != nil {
		return "", err
	}

	return relation.ID, nil
}

func validateIssueRelationCreateRequest(request IssueRelationCreateRequest) error {
	if request.IssueID == "" || request.RelatedIssueID == "" {
		return fmt.Errorf("%w: issue id and related issue id are required", ErrWriteInvalid)
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
