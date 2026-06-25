package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addNotificationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	parentCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.NotificationList, client.NotificationSummary]{
			Use:           "notification",
			Short:         "Read Linear notifications",
			ListShort:     "List authenticated user notifications",
			LimitHelp:     "maximum notifications to return",
			GetUse:        "get NOTIFICATION_ID",
			GetShort:      "Get one notification by id",
			LoadList:      loadNotificationList,
			PageWithItems: notificationPageWithItems,
			LoadGet:       loadNotification,
			WriteItem:     writeNotification,
		},
	)

	addReadListGetCommand(
		ctx,
		parentCommand,
		options,
		readListGetSpec[client.NotificationSubscriptionList, client.NotificationSubscriptionSummary]{
			Use:           "subscription",
			Short:         "Read Linear notification subscriptions",
			ListShort:     "List authenticated user notification subscriptions",
			LimitHelp:     "maximum notification subscriptions to return",
			GetUse:        "get NOTIFICATION_SUBSCRIPTION_ID",
			GetShort:      "Get one notification subscription by id",
			LoadList:      loadNotificationSubscriptionList,
			PageWithItems: notificationSubscriptionPageWithItems,
			LoadGet:       loadNotificationSubscription,
			WriteItem:     writeNotificationSubscription,
		},
	)
}

func writeNotification(command *cobra.Command, options *rootOptions, notification client.NotificationSummary) error {
	return writeItem(command, options, notification, notification.ID,
		func(command *cobra.Command, _ *rootOptions, notification client.NotificationSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s [%s] %s",
				notification.ID,
				notification.Type,
				notification.Category,
				notification.Title,
			)
		})
}

func writeNotificationSubscription(
	command *cobra.Command,
	options *rootOptions,
	subscription client.NotificationSubscriptionSummary,
) error {
	return writeItem(command, options, subscription, subscription.ID,
		func(command *cobra.Command, _ *rootOptions, subscription client.NotificationSubscriptionSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s active %t",
				subscription.ID,
				emptyDash(subscription.TargetType),
				emptyDash(subscription.TargetName),
				subscription.Active,
			)
		})
}

func loadNotificationList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.NotificationList, []client.NotificationSummary, error) {
	notifications, err := client.ListNotifications(ctx, runtime.graphqlClient, limit)
	return notifications, notifications.Notifications, err
}

func loadNotification(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.NotificationSummary, error) {
	return client.GetNotificationByID(ctx, runtime.graphqlClient, id)
}

func notificationPageWithItems(
	page client.NotificationList,
	notifications []client.NotificationSummary,
) client.NotificationList {
	page.Notifications = notifications
	return page
}

func loadNotificationSubscriptionList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.NotificationSubscriptionList, []client.NotificationSubscriptionSummary, error) {
	subscriptions, err := client.ListNotificationSubscriptions(ctx, runtime.graphqlClient, limit)
	return subscriptions, subscriptions.Subscriptions, err
}

func loadNotificationSubscription(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.NotificationSubscriptionSummary, error) {
	return client.GetNotificationSubscriptionByID(ctx, runtime.graphqlClient, id)
}

func notificationSubscriptionPageWithItems(
	page client.NotificationSubscriptionList,
	subscriptions []client.NotificationSubscriptionSummary,
) client.NotificationSubscriptionList {
	page.Subscriptions = subscriptions
	return page
}
