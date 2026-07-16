package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// IssueLabelAssociationRequest describes a guarded label attach/detach on an issue.
type IssueLabelAssociationRequest struct {
	IssueID string
	LabelID string
}

// AddIssueLabel attaches an IssueLabel to an issue after resolving and
// comparing both the issue and the label against the pinned target. A
// team-scoped label must match the resolved team; an organization-wide label
// is always attachable within the resolved organization.
func AddIssueLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueLabelAssociationRequest,
) (IssueSummary, error) {
	if err := validateIssueLabelAssociationRequest(request); err != nil {
		return IssueSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}

	return guard.addIssueLabel(ctx, request)
}

func (guard *guardedClient) addIssueLabel(
	ctx context.Context,
	request IssueLabelAssociationRequest,
) (IssueSummary, error) {
	issue, err := guard.resolveIssueLabelAssociation(ctx, request)
	if err != nil {
		return IssueSummary{}, err
	}

	updated, err := gql.IssueAddLabel(ctx, guard.graphqlClient, issue.ID, request.LabelID)
	if err != nil {
		return IssueSummary{}, fmt.Errorf("add label to issue %s: %w", request.IssueID, err)
	}
	succeeded := updated.IssueAddLabel.Success && updated.IssueAddLabel.Issue != nil
	if err := mutationSuccess(succeeded, "issueAddLabel"); err != nil {
		return IssueSummary{}, err
	}

	return issueSummaryFromFields(updated.IssueAddLabel.Issue.IssueSummaryFields), nil
}

// RemoveIssueLabel detaches an IssueLabel from an issue after resolving and
// comparing both the issue and the label against the pinned target.
func RemoveIssueLabel(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request IssueLabelAssociationRequest,
) (IssueSummary, error) {
	if err := validateIssueLabelAssociationRequest(request); err != nil {
		return IssueSummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return IssueSummary{}, err
	}

	return guard.removeIssueLabel(ctx, request)
}

func (guard *guardedClient) removeIssueLabel(
	ctx context.Context,
	request IssueLabelAssociationRequest,
) (IssueSummary, error) {
	issue, err := guard.resolveIssueLabelAssociation(ctx, request)
	if err != nil {
		return IssueSummary{}, err
	}

	updated, err := gql.IssueRemoveLabel(ctx, guard.graphqlClient, issue.ID, request.LabelID)
	if err != nil {
		return IssueSummary{}, fmt.Errorf("remove label from issue %s: %w", request.IssueID, err)
	}
	succeeded := updated.IssueRemoveLabel.Success && updated.IssueRemoveLabel.Issue != nil
	if err := mutationSuccess(succeeded, "issueRemoveLabel"); err != nil {
		return IssueSummary{}, err
	}

	return issueSummaryFromFields(updated.IssueRemoveLabel.Issue.IssueSummaryFields), nil
}

// resolveIssueLabelAssociation resolves the issue and confirms the label is
// attachable within the resolved team, shared by AddIssueLabel and
// RemoveIssueLabel before either mutation is sent.
func (guard *guardedClient) resolveIssueLabelAssociation(
	ctx context.Context,
	request IssueLabelAssociationRequest,
) (IssueSummary, error) {
	issue, err := guard.requireIssue(ctx, request.IssueID)
	if err != nil {
		return IssueSummary{}, err
	}
	if err := guard.requireAttachableLabel(ctx, request.LabelID); err != nil {
		return IssueSummary{}, err
	}

	return issue, nil
}

func validateIssueLabelAssociationRequest(request IssueLabelAssociationRequest) error {
	if request.IssueID == "" || request.LabelID == "" {
		return fmt.Errorf("%w: issue id and label id are required", ErrWriteInvalid)
	}

	return nil
}
