package cli

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

func Test_CommandFlows_initiative_update_create(t *testing.T) {
	restore := useCommandRuntime(t, initiativeUpdateCreateFake{})
	defer restore()

	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{
		"initiative-update", "create", "initiative-id",
		"--body", "Status note",
		"--health", "on-track",
	})
	require.NoError(t, command.ExecuteContext(context.Background()))
}

func Test_CommandFlows_initiative_update_create_validation_errors(t *testing.T) {
	t.Run("bad health", func(t *testing.T) {
		restore := useCommandRuntime(t, initiativeUpdateCreateFake{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"initiative-update", "create", "initiative-id", "--health", "nope"})
		require.Error(t, command.ExecuteContext(context.Background()))
	})
	t.Run("missing body file", func(t *testing.T) {
		restore := useCommandRuntime(t, initiativeUpdateCreateFake{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{
			"initiative-update", "create", "initiative-id",
			"--body-file", "/no/such/file.md",
		})
		require.Error(t, command.ExecuteContext(context.Background()))
	})
	t.Run("body and body-file exclusive", func(t *testing.T) {
		restore := useCommandRuntime(t, initiativeUpdateCreateFake{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{
			"initiative-update", "create", "initiative-id",
			"--body", "x",
			"--body-file", "also.md",
		})
		require.Error(t, command.ExecuteContext(context.Background()))
	})
	t.Run("body stdin read failure", func(t *testing.T) {
		restore := useCommandRuntime(t, initiativeUpdateCreateFake{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetIn(initiativeBodyFailingReader{})
		command.SetArgs([]string{
			"initiative-update", "create", "initiative-id",
			"--body", "-",
			"--health", "on-track",
		})
		require.Error(t, command.ExecuteContext(context.Background()))
	})
}

type initiativeBodyFailingReader struct{}

func (initiativeBodyFailingReader) Read([]byte) (int, error) {
	return 0, context.DeadlineExceeded
}

func Test_CommandFlows_notification_ack_writes(t *testing.T) {
	t.Run("mark-read", func(t *testing.T) {
		restore := useCommandRuntime(t, notificationAckFake{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"notification", "mark-read", "notification-id"})
		require.NoError(t, command.ExecuteContext(context.Background()))
	})
	t.Run("archive", func(t *testing.T) {
		restore := useCommandRuntime(t, notificationAckFake{})
		defer restore()
		command := NewRootCommand(context.Background(), BuildInfo{})
		command.SetArgs([]string{"notification", "archive", "notification-id"})
		require.NoError(t, command.ExecuteContext(context.Background()))
	})
}

type initiativeUpdateCreateFake struct {
	commandFlowFakeClient
}

func (client initiativeUpdateCreateFake) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	switch request.OpName {
	case "initiative":
		payload := `{"initiative":` + commandInitiativeJSON() + `}`

		return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
	case "InitiativeUpdateCreate":
		payload := `{"initiativeUpdateCreate":{"success":true,"initiativeUpdate":` +
			commandInitiativeUpdateJSON() + `}}`

		return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
	default:
		return client.commandFlowFakeClient.MakeRequest(ctx, request, response)
	}
}

type notificationAckFake struct {
	commandFlowFakeClient
}

func (client notificationAckFake) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	payload := `{
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
		"readAt":"2026-06-20T00:00:00Z",
		"emailedAt":null,
		"snoozedUntilAt":null,
		"unsnoozedAt":null,
		"user":{"id":"user-id","displayName":"Omer"},
		"actor":null,
		"externalUserActor":null
	}`
	switch request.OpName {
	case "notification":
		return json.Unmarshal([]byte(`{"data":{"notification":`+payload+`}}`), response)
	case "NotificationUpdate":
		return json.Unmarshal(
			[]byte(`{"data":{"notificationUpdate":{"success":true,"notification":`+payload+`}}}`),
			response,
		)
	case "NotificationArchive":
		return json.Unmarshal(
			[]byte(`{"data":{"notificationArchive":{"success":true,"entity":`+payload+`}}}`),
			response,
		)
	default:
		return client.commandFlowFakeClient.MakeRequest(ctx, request, response)
	}
}
