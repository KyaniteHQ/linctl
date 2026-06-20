package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// NotificationSummary is the compact notification model used by read-only commands.
type NotificationSummary struct {
	ID                  string `json:"id"`
	Type                string `json:"type"`
	Category            string `json:"category"`
	Title               string `json:"title"`
	Subtitle            string `json:"subtitle"`
	URL                 string `json:"url"`
	InboxURL            string `json:"inbox_url"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	ArchivedAt          string `json:"archived_at,omitempty"`
	ReadAt              string `json:"read_at,omitempty"`
	EmailedAt           string `json:"emailed_at,omitempty"`
	SnoozedUntilAt      string `json:"snoozed_until_at,omitempty"`
	UnsnoozedAt         string `json:"unsnoozed_at,omitempty"`
	UserID              string `json:"user_id"`
	UserDisplayName     string `json:"user_display_name"`
	ActorID             string `json:"actor_id,omitempty"`
	ActorDisplayName    string `json:"actor_display_name,omitempty"`
	ExternalUserActorID string `json:"external_user_actor_id,omitempty"`
}

// NotificationList is a page of Linear notifications.
type NotificationList struct {
	Notifications []NotificationSummary `json:"notifications"`
	HasNextPage   bool                  `json:"has_next_page"`
	EndCursor     *string               `json:"end_cursor,omitempty"`
}

// NotificationSubscriptionSummary is the compact subscription model used by read-only commands.
type NotificationSubscriptionSummary struct {
	ID                  string `json:"id"`
	Active              bool   `json:"active"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	ArchivedAt          string `json:"archived_at,omitempty"`
	ContextViewType     string `json:"context_view_type,omitempty"`
	UserContextViewType string `json:"user_context_view_type,omitempty"`
	SubscriberID        string `json:"subscriber_id"`
	SubscriberName      string `json:"subscriber_name"`
	TargetType          string `json:"target_type,omitempty"`
	TargetID            string `json:"target_id,omitempty"`
	TargetName          string `json:"target_name,omitempty"`
}

// NotificationSubscriptionList is a page of Linear notification subscriptions.
type NotificationSubscriptionList struct {
	Subscriptions []NotificationSubscriptionSummary `json:"notification_subscriptions"`
	HasNextPage   bool                              `json:"has_next_page"`
	EndCursor     *string                           `json:"end_cursor,omitempty"`
}

// ListNotifications returns the authenticated user's notifications.
func ListNotifications(ctx context.Context, graphqlClient graphql.Client, limit int) (NotificationList, error) {
	result, err := notifications(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return NotificationList{}, fmt.Errorf("list notifications: %w", err)
	}

	summaries := make([]NotificationSummary, 0, len(result.Notifications.Nodes))
	for _, node := range result.Notifications.Nodes {
		summaries = append(summaries, notificationSummary(node))
	}

	return NotificationList{
		Notifications: summaries,
		HasNextPage:   result.Notifications.PageInfo.HasNextPage,
		EndCursor:     result.Notifications.PageInfo.EndCursor,
	}, nil
}

// GetNotificationByID returns one notification by id.
func GetNotificationByID(ctx context.Context, graphqlClient graphql.Client, id string) (NotificationSummary, error) {
	result, err := notification(ctx, graphqlClient, id)
	if err != nil {
		return NotificationSummary{}, fmt.Errorf("get notification %s: %w", id, err)
	}

	return notificationSummary(result.Notification), nil
}

// ListNotificationSubscriptions returns the authenticated user's notification subscriptions.
func ListNotificationSubscriptions(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (NotificationSubscriptionList, error) {
	result, err := notificationSubscriptions(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true))
	if err != nil {
		return NotificationSubscriptionList{}, fmt.Errorf("list notification subscriptions: %w", err)
	}

	summaries := make([]NotificationSubscriptionSummary, 0, len(result.NotificationSubscriptions.Nodes))
	for _, node := range result.NotificationSubscriptions.Nodes {
		summaries = append(summaries, notificationSubscriptionSummary(node))
	}

	return NotificationSubscriptionList{
		Subscriptions: summaries,
		HasNextPage:   result.NotificationSubscriptions.PageInfo.HasNextPage,
		EndCursor:     result.NotificationSubscriptions.PageInfo.EndCursor,
	}, nil
}

// GetNotificationSubscriptionByID returns one notification subscription by id.
func GetNotificationSubscriptionByID(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
) (NotificationSubscriptionSummary, error) {
	result, err := notificationSubscription(ctx, graphqlClient, id)
	if err != nil {
		return NotificationSubscriptionSummary{}, fmt.Errorf("get notification subscription %s: %w", id, err)
	}

	return notificationSubscriptionSummary(result.NotificationSubscription), nil
}

func notificationSummary(fields NotificationSummaryFields) NotificationSummary {
	user := fields.GetUser()
	summary := NotificationSummary{
		ID:              fields.GetId(),
		Type:            fields.GetType(),
		Category:        string(fields.GetCategory()),
		Title:           fields.GetTitle(),
		Subtitle:        fields.GetSubtitle(),
		URL:             fields.GetUrl(),
		InboxURL:        fields.GetInboxUrl(),
		CreatedAt:       fields.GetCreatedAt(),
		UpdatedAt:       fields.GetUpdatedAt(),
		ArchivedAt:      stringValue(fields.GetArchivedAt()),
		ReadAt:          stringValue(fields.GetReadAt()),
		EmailedAt:       stringValue(fields.GetEmailedAt()),
		SnoozedUntilAt:  stringValue(fields.GetSnoozedUntilAt()),
		UnsnoozedAt:     stringValue(fields.GetUnsnoozedAt()),
		UserID:          user.Id,
		UserDisplayName: user.DisplayName,
	}
	if actor := fields.GetActor(); actor != nil {
		summary.ActorID = actor.Id
		summary.ActorDisplayName = actor.DisplayName
	}
	if actor := fields.GetExternalUserActor(); actor != nil {
		summary.ExternalUserActorID = actor.Id
	}

	return summary
}

func notificationSubscriptionSummary(fields NotificationSubscriptionSummaryFields) NotificationSubscriptionSummary {
	subscriber := fields.GetSubscriber()
	summary := NotificationSubscriptionSummary{
		ID:                  fields.GetId(),
		Active:              fields.GetActive(),
		CreatedAt:           fields.GetCreatedAt(),
		UpdatedAt:           fields.GetUpdatedAt(),
		ArchivedAt:          stringValue(fields.GetArchivedAt()),
		ContextViewType:     contextViewTypeValue(fields.GetContextViewType()),
		UserContextViewType: userContextViewTypeValue(fields.GetUserContextViewType()),
		SubscriberID:        subscriber.Id,
		SubscriberName:      subscriber.DisplayName,
	}
	setNotificationSubscriptionTarget(&summary, fields)

	return summary
}

func contextViewTypeValue(value *ContextViewType) string {
	if value == nil {
		return ""
	}
	return string(*value)
}

func userContextViewTypeValue(value *UserContextViewType) string {
	if value == nil {
		return ""
	}
	return string(*value)
}

func setNotificationSubscriptionTarget(
	summary *NotificationSubscriptionSummary,
	fields NotificationSubscriptionSummaryFields,
) {
	if target := fields.GetCustomer(); target != nil {
		setNotificationSubscriptionTargetValues(summary, "customer", target.Id, target.Name)
	}
	if target := fields.GetCustomView(); target != nil {
		setNotificationSubscriptionTargetValues(summary, "custom_view", target.Id, target.Name)
	}
	if target := fields.GetCycle(); target != nil {
		setNotificationSubscriptionTargetValues(summary, "cycle", target.Id, stringValue(target.Name))
	}
	if target := fields.GetInitiative(); target != nil {
		setNotificationSubscriptionTargetValues(summary, "initiative", target.Id, target.Name)
	}
	if target := fields.GetLabel(); target != nil {
		setNotificationSubscriptionTargetValues(summary, "label", target.Id, target.Name)
	}
	if target := fields.GetProject(); target != nil {
		setNotificationSubscriptionTargetValues(summary, "project", target.Id, target.Name)
	}
	if target := fields.GetTeam(); target != nil {
		setNotificationSubscriptionTargetValues(summary, "team", target.Id, target.Key)
	}
	if target := fields.GetUser(); target != nil {
		setNotificationSubscriptionTargetValues(summary, "user", target.Id, target.DisplayName)
	}
}

func setNotificationSubscriptionTargetValues(
	summary *NotificationSubscriptionSummary,
	targetType string,
	targetID string,
	targetName string,
) {
	summary.TargetType = targetType
	summary.TargetID = targetID
	summary.TargetName = targetName
}
