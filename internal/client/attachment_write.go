package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/config"
)

// AttachmentLinkRequest describes a guarded URL attachment on an issue.
type AttachmentLinkRequest struct {
	IssueID  string
	URL      string
	Title    string
	Subtitle string
}

// LinearAttachmentCreateInput is the sparse Linear attachmentCreate payload linctl supports.
type LinearAttachmentCreateInput struct {
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
	URL      string  `json:"url"`
	IssueID  string  `json:"issueId"`
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
		return AttachmentSummary{}, fmt.Errorf("%w: issue id is required", ErrWriteInvalid)
	}
	if request.URL == "" {
		return AttachmentSummary{}, fmt.Errorf("%w: url is required", ErrWriteInvalid)
	}

	return guardedMutation(ctx, graphqlClient, expected, func(guard writeGuard) (AttachmentSummary, error) {
		issue, err := guard.requireIssue(ctx, graphqlClient, request.IssueID)
		if err != nil {
			return AttachmentSummary{}, err
		}

		created, err := AttachmentLinkURL(ctx, graphqlClient, LinearAttachmentCreateInput{
			Title:    optionalString(request.Title),
			Subtitle: optionalString(request.Subtitle),
			URL:      request.URL,
			IssueID:  issue.ID,
		})
		if err != nil {
			return AttachmentSummary{}, fmt.Errorf("link attachment to issue %s: %w", request.IssueID, err)
		}
		if !created.AttachmentCreate.Success {
			return AttachmentSummary{}, fmt.Errorf("%w: attachmentCreate reported no success", ErrMutationFailed)
		}

		return attachmentSummary(created.AttachmentCreate.Attachment.AttachmentSummaryFields), nil
	})
}
