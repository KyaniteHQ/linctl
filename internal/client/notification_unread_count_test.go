package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetNotificationsUnreadCount_returns_fake_response(t *testing.T) {
	count, err := GetNotificationsUnreadCount(context.Background(), fakeGraphQLClient{
		"notificationsUnreadCount": `{"notificationsUnreadCount":7}`,
	})

	require.NoError(t, err)
	require.Equal(t, 7, count)
}

func Test_GetNotificationsUnreadCount_wraps_operation_errors(t *testing.T) {
	_, err := GetNotificationsUnreadCount(context.Background(), fakeGraphQLClient{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "get notifications unread count")
}
