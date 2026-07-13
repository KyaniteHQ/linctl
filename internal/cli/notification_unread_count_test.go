package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

type notificationUnreadCountFakeClient struct {
	fail bool
}

func (client notificationUnreadCountFakeClient) MakeRequest(
	_ context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if client.fail {
		return errors.New("operation failed")
	}
	if request.OpName != "notificationsUnreadCount" {
		return errors.New("unexpected operation")
	}

	return json.Unmarshal([]byte(`{"data":{"notificationsUnreadCount":7}}`), response)
}

func Test_NotificationUnreadCount_command_outputs_supported_modes(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "human", args: []string{"notification", "unread-count"}, want: "7\n"},
		{name: "json", args: []string{"--json", "notification", "unread-count"}, want: "{\n  \"unread_count\": 7\n}\n"},
		{name: "quiet", args: []string{"--quiet", "notification", "unread-count"}, want: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, notificationUnreadCountFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, test.want, output.String())
		})
	}
}

func Test_NotificationUnreadCount_command_returns_operation_errors(t *testing.T) {
	restore := useCommandRuntime(t, notificationUnreadCountFakeClient{fail: true})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"notification", "unread-count"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "get notifications unread count")
}

func Test_NotificationUnreadCount_command_returns_runtime_errors(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	t.Cleanup(func() { buildCommandRuntime = original })
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"notification", "unread-count"})

	err := command.ExecuteContext(context.Background())

	require.EqualError(t, err, "runtime failed")
}
