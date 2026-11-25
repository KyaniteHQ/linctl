package cli

import (
	"cmp"
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addNotificationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	parentCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.NotificationList, client.NotificationSummary]{
			Use:       "notification",
			Short:     "Read Linear notifications",
			ListShort: "List the notifications of the authenticated user",
			LimitHelp: "maximum notifications to return",
			GetUse:    "get NOTIFICATION_ID",
			GetShort:  "Get one notification by id",
			LoadList:  clientList(client.ListNotifications),
			LoadGet:   clientGet(client.GetNotificationByID),
			WriteItem: writeNotification,
		},
	)
	addNotificationMarkReadCommand(ctx, parentCommand, options)
	addNotificationArchiveCommand(ctx, parentCommand, options)

	addReadListGetCommand(
		ctx,
		parentCommand,
		options,
		readListGetSpec[client.NotificationSubscriptionList, client.NotificationSubscriptionSummary]{
			Use:       "subscription",
			Short:     "Read Linear notification subscriptions",
			ListShort: "List the notification subscriptions of the authenticated user",
			LimitHelp: "maximum notification subscriptions to return",
			GetUse:    "get NOTIFICATION_SUBSCRIPTION_ID",
			GetShort:  "Get one notification subscription by id",
			LoadList:  clientList(client.ListNotificationSubscriptions),
			LoadGet:   clientGet(client.GetNotificationSubscriptionByID),
			WriteItem: writeNotificationSubscription,
		},
	)
}

func addNotificationMarkReadCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.NotificationSummary]{
		Use:   "mark-read NOTIFICATION_ID",
		Short: "Mark one notification read for the authenticated actor's inbox",
		Args:  cobra.ExactArgs(1),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.NotificationSummary, error) {
			return client.MarkNotificationRead(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeNotification,
	})
}

func addNotificationArchiveCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addGuardedWriteCommand(ctx, root, options, guardedWriteSpec[client.NotificationSummary]{
		Use:   "archive NOTIFICATION_ID",
		Short: "Archive one notification for the authenticated actor's inbox",
		Args:  cobra.ExactArgs(1),
		Run: func(
			ctx context.Context, _ *cobra.Command, runtime commandRuntime, args []string,
		) (client.NotificationSummary, error) {
			return client.ArchiveNotification(ctx, runtime.graphqlClient, runtime.config.Target, args[0])
		},
		Write: writeNotification,
	})
}

func writeNotification(command *cobra.Command, options *rootOptions, notification client.NotificationSummary) error {
	return writeItemLine(
		command, options, notification, notification.ID,
		"%s %s [%s] %s %s", notification.ID, notification.Type, notification.Category,
		emptyDash(notification.TargetType), emptyDash(cmp.Or(notification.TargetName, notification.TargetID)),
	)
}

func writeNotificationSubscription(
	command *cobra.Command,
	options *rootOptions,
	subscription client.NotificationSubscriptionSummary,
) error {
	return writeItemLine(
		command, options, subscription, subscription.ID,
		"%s %s %s active %t",
		subscription.ID,
		emptyDash(subscription.TargetType),
		emptyDash(subscription.TargetName),
		subscription.Active,
	)
}
