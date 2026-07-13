package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

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

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (IssueSummary, error) {
		issue, err := resolveIssueLabelAssociation(ctx, graphqlClient, guard, request)
		if err != nil {
			return IssueSummary{}, err
		}

		updated, err := IssueAddLabel(ctx, graphqlClient, issue.ID, request.LabelID)
		if err != nil {
			return IssueSummary{}, fmt.Errorf("add label to issue %s: %w", request.IssueID, err)
		}
		if !updated.IssueAddLabel.Success || updated.IssueAddLabel.Issue == nil {
			return IssueSummary{}, fmt.Errorf("%w: issueAddLabel reported no success", ErrMutationFailed)
		}

		return issueSummaryFromFields(updated.IssueAddLabel.Issue.IssueSummaryFields), nil
	})
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

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (IssueSummary, error) {
		issue, err := resolveIssueLabelAssociation(ctx, graphqlClient, guard, request)
		if err != nil {
			return IssueSummary{}, err
		}

		updated, err := IssueRemoveLabel(ctx, graphqlClient, issue.ID, request.LabelID)
		if err != nil {
			return IssueSummary{}, fmt.Errorf("remove label from issue %s: %w", request.IssueID, err)
		}
		if !updated.IssueRemoveLabel.Success || updated.IssueRemoveLabel.Issue == nil {
			return IssueSummary{}, fmt.Errorf("%w: issueRemoveLabel reported no success", ErrMutationFailed)
		}

		return issueSummaryFromFields(updated.IssueRemoveLabel.Issue.IssueSummaryFields), nil
	})
}

// resolveIssueLabelAssociation resolves the issue and confirms the label is
// attachable within the resolved team, shared by AddIssueLabel and
// RemoveIssueLabel before either mutation is sent.
func resolveIssueLabelAssociation(
	ctx context.Context,
	graphqlClient graphql.Client,
	guard writeGuard,
	request IssueLabelAssociationRequest,
) (IssueSummary, error) {
	issue, err := guard.requireIssue(ctx, graphqlClient, request.IssueID)
	if err != nil {
		return IssueSummary{}, err
	}
	if err := guard.requireAttachableLabel(ctx, graphqlClient, request.LabelID); err != nil {
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
