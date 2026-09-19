package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// AttachmentLinkRequest describes a guarded URL attachment on an issue.
type AttachmentLinkRequest struct {
	IssueID  string
	URL      string
	Title    string
	Subtitle string
}

// LinkIssueAttachment attaches a URL to an issue after resolving and comparing
// the pinned write target. The issue must belong to the resolved team.
func LinkIssueAttachment(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request AttachmentLinkRequest,
) (AttachmentSummary, error) {
	if request.IssueID == "" {
		return AttachmentSummary{}, requiredFieldError("issue id")
	}
	if request.URL == "" {
		return AttachmentSummary{}, requiredFieldError("url")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return AttachmentSummary{}, err
	}

	return guard.linkIssueAttachment(ctx, request)
}

func (guard *guardedClient) linkIssueAttachment(
	ctx context.Context,
	request AttachmentLinkRequest,
) (AttachmentSummary, error) {
	issue, err := guard.requireIssue(ctx, request.IssueID)
	if err != nil {
		return AttachmentSummary{}, err
	}

	created, err := gql.AttachmentLinkURL(ctx, guard.graphqlClient, LinearAttachmentCreateInput{
		Title:    optionalString(request.Title),
		Subtitle: optionalString(request.Subtitle),
		URL:      request.URL,
		IssueID:  issue.ID,
	})
	if err != nil {
		return AttachmentSummary{}, fmt.Errorf("link attachment to issue %s: %w", request.IssueID, err)
	}
	if err := mutationSuccess(created.AttachmentCreate.Success, "attachmentCreate"); err != nil {
		return AttachmentSummary{}, err
	}

	return attachmentSummary(created.AttachmentCreate.Attachment.AttachmentSummaryFields), nil
}

// DeleteAttachment removes an issue attachment after resolving the attachment's
// issue and comparing the pinned target. Attachment delete is irreversible
// through linctl.
func DeleteAttachment(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	attachmentID string,
) (string, error) {
	if attachmentID == "" {
		return "", requiredFieldError("attachment id")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return "", err
	}

	return guard.deleteAttachment(ctx, attachmentID)
}

func (guard *guardedClient) deleteAttachment(ctx context.Context, attachmentID string) (string, error) {
	if err := guard.requireAttachmentTarget(ctx, attachmentID); err != nil {
		return "", err
	}

	deleted, err := gql.AttachmentDelete(ctx, guard.graphqlClient, attachmentID)
	if err != nil {
		return "", fmt.Errorf("delete attachment %s: %w", attachmentID, err)
	}
	if err := mutationSuccess(deleted.AttachmentDelete.Success, "attachmentDelete"); err != nil {
		return "", err
	}

	return deleted.AttachmentDelete.EntityId, nil
}

func (guard *guardedClient) requireAttachmentTarget(ctx context.Context, attachmentID string) error {
	issue, err := GetAttachmentIssue(ctx, guard.graphqlClient, attachmentID)
	if err != nil {
		return err
	}
	if issue.ID == "" {
		return fmt.Errorf(
			"%w: attachment %s is not attached to an issue; only issue attachments are guarded",
			ErrWriteInvalid,
			attachmentID,
		)
	}
	_, err = guard.requireIssue(ctx, issue.ID)

	return err
}
