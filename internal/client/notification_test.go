package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetNotificationByID_reports_the_public_target_of_each_notification_type(t *testing.T) {
	tests := []struct {
		name     string
		typename string
		target   string
		want     NotificationSummary
	}{
		{
			name:     "issue",
			typename: "IssueNotification",
			target:   `"issue":{"id":"issue-id","identifier":"LIT-1","title":"Mentioned you","url":"https://linear.app/kyanite/issue/LIT-1"}`,
			want: NotificationSummary{
				TargetType: "issue",
				TargetID:   "issue-id",
				TargetName: "LIT-1 Mentioned you",
				TargetURL:  "https://linear.app/kyanite/issue/LIT-1",
			},
		},
		{
			name:     "project",
			typename: "ProjectNotification",
			target:   `"project":{"id":"project-id","name":"Pinned project","url":"https://linear.app/kyanite/project/pinned"}`,
			want: NotificationSummary{
				TargetType: "project",
				TargetID:   "project-id",
				TargetName: "Pinned project",
				TargetURL:  "https://linear.app/kyanite/project/pinned",
			},
		},
		{
			name:     "initiative",
			typename: "InitiativeNotification",
			target:   `"initiative":{"id":"initiative-id","name":"Platform","url":"https://linear.app/kyanite/initiative/platform"}`,
			want: NotificationSummary{
				TargetType: "initiative",
				TargetID:   "initiative-id",
				TargetName: "Platform",
				TargetURL:  "https://linear.app/kyanite/initiative/platform",
			},
		},
		{
			name:     "initiative removed",
			typename: "InitiativeNotification",
			target:   `"initiative":null`,
		},
		{
			name:     "document",
			typename: "DocumentNotification",
			target:   `"documentId":"document-id"`,
			want:     NotificationSummary{TargetType: "document", TargetID: "document-id"},
		},
		{
			name:     "without a target",
			typename: "PullRequestNotification",
			target:   `"pullRequest":null`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := notificationJSONWithTarget(test.typename, test.target)
			graphqlClient := projectWriteFakeClient(map[string]string{"notification": `{"notification":` + payload + `}`})

			notification, err := GetNotificationByID(context.Background(), graphqlClient, "notification-id")

			require.NoError(t, err)
			require.Equal(t, "notification-id", notification.ID)
			require.Equal(t, test.want.TargetType, notification.TargetType)
			require.Equal(t, test.want.TargetID, notification.TargetID)
			require.Equal(t, test.want.TargetName, notification.TargetName)
			require.Equal(t, test.want.TargetURL, notification.TargetURL)
		})
	}
}
