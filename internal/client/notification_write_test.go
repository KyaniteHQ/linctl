package client

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/config"
)

func notificationWriteJSON(userID string, readAt string) string {
	return `{
		"__typename":"IssueNotification",
		"id":"notification-id",
		"type":"issueAssignedToYou",
		"category":"assignments",
		"title":"Assigned",
		"subtitle":"LIT-1",
		"url":"https://linear.app/n/notification-id",
		"inboxUrl":"https://linear.app/inbox",
		"createdAt":"2026-06-19T12:00:00Z",
		"updatedAt":"2026-06-19T12:00:00Z",
		"archivedAt":null,
		"readAt":` + nullOrString(readAt) + `,
		"emailedAt":null,
		"snoozedUntilAt":null,
		"unsnoozedAt":null,
		"user":{"id":"` + userID + `","displayName":"Omer"},
		"actor":null,
		"externalUserActor":null
	}`
}

func nullOrString(value string) string {
	if value == "" {
		return "null"
	}

	return `"` + value + `"`
}

func Test_MarkNotificationRead_sets_readAt_when_viewer_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"notification": `{"notification":` + notificationWriteJSON("user-id", "") + `}`,
		"NotificationUpdate": `{"notificationUpdate":{"success":true,"notification":` +
			notificationWriteJSON("user-id", "2026-06-20T00:00:00Z") + `}}`,
	})}

	notification, err := MarkNotificationRead(
		context.Background(), recorder, matchingTarget(), "notification-id",
	)

	require.NoError(t, err)
	require.Equal(t, "notification-id", notification.ID)
	require.Equal(t, "user-id", notification.UserID)
	require.NotEmpty(t, notification.ReadAt)
	vars := map[string]any{}
	require.NoError(t, json.Unmarshal(recorder.variablesFor(t, "NotificationUpdate"), &vars))
	require.Equal(t, "notification-id", vars["id"])
	input, ok := vars["input"].(map[string]any)
	require.True(t, ok)
	require.NotEmpty(t, input["readAt"])
	require.Len(t, input, 1)
}

func Test_MarkNotificationRead_refuses_foreign_recipient(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"notification": `{"notification":` + notificationWriteJSON("other-user", "") + `}`,
	})}

	_, err := MarkNotificationRead(context.Background(), recorder, matchingTarget(), "notification-id")
	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("NotificationUpdate"))
}

func Test_ArchiveNotification_returns_summary_when_viewer_matches(t *testing.T) {
	recorder := &recordingGraphQLClient{inner: projectWriteFakeClient(map[string]string{
		"notification": `{"notification":` + notificationWriteJSON("user-id", "") + `}`,
		"NotificationArchive": `{"notificationArchive":{"success":true,"entity":` +
			notificationWriteJSON("user-id", "2026-06-20T00:00:00Z") + `}}`,
	})}

	notification, err := ArchiveNotification(
		context.Background(), recorder, matchingTarget(), "notification-id",
	)

	require.NoError(t, err)
	require.Equal(t, "notification-id", notification.ID)
	require.JSONEq(t, `{"id":"notification-id"}`, string(recorder.variablesFor(t, "NotificationArchive")))
}

func Test_ArchiveNotification_refuses_foreign_recipient(t *testing.T) {
	recorder := &mutationRecordingClient{inner: projectWriteFakeClient(map[string]string{
		"notification": `{"notification":` + notificationWriteJSON("other-user", "") + `}`,
	})}

	_, err := ArchiveNotification(context.Background(), recorder, matchingTarget(), "notification-id")
	require.ErrorIs(t, err, ErrTargetMismatch)
	require.False(t, recorder.sentOperation("NotificationArchive"))
}

func Test_notification_ack_requires_id(t *testing.T) {
	_, err := MarkNotificationRead(context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), "")
	require.ErrorIs(t, err, ErrWriteInvalid)
	_, err = ArchiveNotification(context.Background(), projectWriteFakeClient(map[string]string{}), matchingTarget(), "")
	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_notification_ack_refuses_unresolved_target(t *testing.T) {
	bad := config.Target{OrgID: "org-id", TeamKey: "WRONG", TeamID: "wrong"}
	_, err := MarkNotificationRead(context.Background(), projectWriteFakeClient(map[string]string{}), bad, "notification-id")
	require.ErrorIs(t, err, ErrTargetMismatch)
	_, err = ArchiveNotification(context.Background(), projectWriteFakeClient(map[string]string{}), bad, "notification-id")
	require.ErrorIs(t, err, ErrTargetMismatch)
}

func Test_MarkNotificationRead_fails_when_mutation_reports_no_success(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"notification":       `{"notification":` + notificationWriteJSON("user-id", "") + `}`,
		"NotificationUpdate": `{"notificationUpdate":{"success":false,"notification":` + notificationWriteJSON("user-id", "") + `}}`,
	})
	_, err := MarkNotificationRead(context.Background(), graphqlClient, matchingTarget(), "notification-id")
	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_ArchiveNotification_fails_when_entity_missing(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"notification":        `{"notification":` + notificationWriteJSON("user-id", "") + `}`,
		"NotificationArchive": `{"notificationArchive":{"success":true,"entity":null}}`,
	})
	_, err := ArchiveNotification(context.Background(), graphqlClient, matchingTarget(), "notification-id")
	require.ErrorIs(t, err, ErrMutationFailed)
}

func Test_MarkNotificationRead_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"notification": `{"notification":` + notificationWriteJSON("user-id", "") + `}`,
	})
	_, err := MarkNotificationRead(context.Background(), graphqlClient, matchingTarget(), "notification-id")
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_ArchiveNotification_wraps_mutation_error(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{
		"notification": `{"notification":` + notificationWriteJSON("user-id", "") + `}`,
	})
	_, err := ArchiveNotification(context.Background(), graphqlClient, matchingTarget(), "notification-id")
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}

func Test_notification_ack_wraps_missing_notification_read(t *testing.T) {
	graphqlClient := projectWriteFakeClient(map[string]string{})
	_, err := MarkNotificationRead(context.Background(), graphqlClient, matchingTarget(), "notification-id")
	require.Error(t, err)
	require.NotErrorIs(t, err, ErrTargetMismatch)
}
