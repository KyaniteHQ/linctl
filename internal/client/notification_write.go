package client

import (
	"context"
	"fmt"
	"time"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// MarkNotificationRead marks one notification as read after Viewer-Scoped
// recipient comparison against the authenticated actor.
func MarkNotificationRead(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	notificationID string,
) (NotificationSummary, error) {
	if notificationID == "" {
		return NotificationSummary{}, requiredFieldError("notification id")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return NotificationSummary{}, err
	}

	return guard.markNotificationRead(ctx, notificationID)
}

// ArchiveNotification archives one notification after Viewer-Scoped recipient
// comparison against the authenticated actor.
func ArchiveNotification(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	notificationID string,
) (NotificationSummary, error) {
	if notificationID == "" {
		return NotificationSummary{}, requiredFieldError("notification id")
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return NotificationSummary{}, err
	}

	return guard.archiveNotification(ctx, notificationID)
}

func (guard *guardedClient) markNotificationRead(
	ctx context.Context,
	notificationID string,
) (NotificationSummary, error) {
	if err := guard.requireNotification(ctx, notificationID); err != nil {
		return NotificationSummary{}, err
	}

	readAt := time.Now().UTC().Format(time.RFC3339)
	updated, err := gql.NotificationUpdate(ctx, guard.graphqlClient, notificationID, LinearNotificationUpdateInput{
		ReadAt: &readAt,
	})
	if err != nil {
		return NotificationSummary{}, fmt.Errorf("mark notification read %s: %w", notificationID, err)
	}
	if !updated.NotificationUpdate.Success {
		return NotificationSummary{}, fmt.Errorf(
			"%w: notificationUpdate returned no notification",
			ErrMutationFailed,
		)
	}

	return notificationSummary(updated.NotificationUpdate.Notification), nil
}

func (guard *guardedClient) archiveNotification(
	ctx context.Context,
	notificationID string,
) (NotificationSummary, error) {
	if err := guard.requireNotification(ctx, notificationID); err != nil {
		return NotificationSummary{}, err
	}

	archived, err := gql.NotificationArchive(ctx, guard.graphqlClient, notificationID)
	if err != nil {
		return NotificationSummary{}, fmt.Errorf("archive notification %s: %w", notificationID, err)
	}
	if !archived.NotificationArchive.Success || archived.NotificationArchive.Entity == nil {
		return NotificationSummary{}, fmt.Errorf(
			"%w: notificationArchive returned no notification",
			ErrMutationFailed,
		)
	}

	return notificationSummary(*archived.NotificationArchive.Entity), nil
}
