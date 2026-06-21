package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"
)

// UserSummary is the compact User model used by user commands.
type UserSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Active      bool   `json:"active"`
	Guest       bool   `json:"guest"`
	Admin       bool   `json:"admin"`
}

// UserList is a page of users.
type UserList struct {
	Users       []UserSummary `json:"users"`
	HasNextPage bool          `json:"has_next_page"`
	EndCursor   *string       `json:"end_cursor,omitempty"`
}

// DraftSummary is the compact saved draft model used by viewer-scoped draft reads.
type DraftSummary struct {
	ID          string `json:"id"`
	ParentType  string `json:"parent_type"`
	ParentID    string `json:"parent_id"`
	ParentKey   string `json:"parent_key,omitempty"`
	ParentTitle string `json:"parent_title,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ArchivedAt  string `json:"archived_at,omitempty"`
}

// DraftList is a page of the authenticated user's saved drafts.
type DraftList struct {
	Drafts      []DraftSummary `json:"drafts"`
	HasNextPage bool           `json:"has_next_page"`
	EndCursor   *string        `json:"end_cursor,omitempty"`
}

// NotificationChannelPreference is a compact channel preference set.
type NotificationChannelPreference struct {
	Desktop bool `json:"desktop"`
	Email   bool `json:"email"`
	Mobile  bool `json:"mobile"`
	Slack   bool `json:"slack"`
}

// NotificationCategoryPreferences is a compact notification preference matrix.
type NotificationCategoryPreferences struct {
	AppsAndIntegrations NotificationChannelPreference `json:"apps_and_integrations"`
	Assignments         NotificationChannelPreference `json:"assignments"`
	Billing             NotificationChannelPreference `json:"billing"`
	CommentsAndReplies  NotificationChannelPreference `json:"comments_and_replies"`
	Customers           NotificationChannelPreference `json:"customers"`
	DocumentChanges     NotificationChannelPreference `json:"document_changes"`
	Feed                NotificationChannelPreference `json:"feed"`
	Mentions            NotificationChannelPreference `json:"mentions"`
	PostsAndUpdates     NotificationChannelPreference `json:"posts_and_updates"`
	Reactions           NotificationChannelPreference `json:"reactions"`
	Reminders           NotificationChannelPreference `json:"reminders"`
	Reviews             NotificationChannelPreference `json:"reviews"`
	StatusChanges       NotificationChannelPreference `json:"status_changes"`
	Subscriptions       NotificationChannelPreference `json:"subscriptions"`
	System              NotificationChannelPreference `json:"system"`
	Triage              NotificationChannelPreference `json:"triage"`
}

// NotificationDeliveryDay is one mobile notification delivery window.
type NotificationDeliveryDay struct {
	Start *string `json:"start,omitempty"`
	End   *string `json:"end,omitempty"`
}

// NotificationDeliverySchedule is the compact weekly notification schedule.
type NotificationDeliverySchedule struct {
	Disabled  *bool                   `json:"disabled,omitempty"`
	Friday    NotificationDeliveryDay `json:"friday"`
	Monday    NotificationDeliveryDay `json:"monday"`
	Saturday  NotificationDeliveryDay `json:"saturday"`
	Sunday    NotificationDeliveryDay `json:"sunday"`
	Thursday  NotificationDeliveryDay `json:"thursday"`
	Tuesday   NotificationDeliveryDay `json:"tuesday"`
	Wednesday NotificationDeliveryDay `json:"wednesday"`
}

// NotificationDeliveryChannel is a compact notification delivery channel.
type NotificationDeliveryChannel struct {
	NotificationsDisabled *bool                         `json:"notifications_disabled,omitempty"`
	Schedule              *NotificationDeliverySchedule `json:"schedule,omitempty"`
}

// NotificationDeliveryPreferences is the compact notification delivery preference set.
type NotificationDeliveryPreferences struct {
	Mobile *NotificationDeliveryChannel `json:"mobile,omitempty"`
}

// UserSettingsSummary is the compact viewer-scoped settings model.
type UserSettingsSummary struct {
	ID                              string                          `json:"id"`
	UserID                          string                          `json:"user_id"`
	CreatedAt                       string                          `json:"created_at"`
	UpdatedAt                       string                          `json:"updated_at"`
	ArchivedAt                      *string                         `json:"archived_at,omitempty"`
	AutoAssignToSelf                bool                            `json:"auto_assign_to_self"`
	FeedLastSeenTime                *string                         `json:"feed_last_seen_time,omitempty"`
	FeedSummarySchedule             string                          `json:"feed_summary_schedule,omitempty"`
	ShowFullUserNames               bool                            `json:"show_full_user_names"`
	SubscribedToChangelog           bool                            `json:"subscribed_to_changelog"`
	SubscribedToDPA                 bool                            `json:"subscribed_to_dpa"`
	SubscribedToInviteAccepted      bool                            `json:"subscribed_to_invite_accepted"`
	SubscribedToPrivacyLegalUpdates bool                            `json:"subscribed_to_privacy_legal_updates"`
	NotificationCategoryPreferences NotificationCategoryPreferences `json:"notification_category_preferences"`
	NotificationChannelPreferences  NotificationChannelPreference   `json:"notification_channel_preferences"`
	NotificationDeliveryPreferences NotificationDeliveryPreferences `json:"notification_delivery_preferences"`
}

// UserSettingsCustomSidebarTheme is a compact custom sidebar theme.
type UserSettingsCustomSidebarTheme struct {
	Accent   []float64 `json:"accent"`
	Base     []float64 `json:"base"`
	Contrast int       `json:"contrast"`
}

// UserSettingsCustomTheme is a compact custom theme.
type UserSettingsCustomTheme struct {
	Accent   []float64                       `json:"accent"`
	Base     []float64                       `json:"base"`
	Contrast int                             `json:"contrast"`
	Sidebar  *UserSettingsCustomSidebarTheme `json:"sidebar,omitempty"`
}

// UserSettingsThemeSummary is a compact resolved theme.
type UserSettingsThemeSummary struct {
	Preset string                   `json:"preset"`
	Custom *UserSettingsCustomTheme `json:"custom,omitempty"`
}

type notificationChannelPreferenceSource interface {
	GetDesktop() bool
	GetEmail() bool
	GetMobile() bool
	GetSlack() bool
}

type notificationDeliveryDaySource interface {
	GetStart() *string
	GetEnd() *string
}

type notificationDeliveryScheduleSource interface {
	GetDisabled() *bool
	GetFriday() NotificationDeliveryPreferencesScheduleFieldsFridayNotificationDeliveryPreferencesDay
	GetMonday() NotificationDeliveryPreferencesScheduleFieldsMondayNotificationDeliveryPreferencesDay
	GetSaturday() NotificationDeliveryPreferencesScheduleFieldsSaturdayNotificationDeliveryPreferencesDay
	GetSunday() NotificationDeliveryPreferencesScheduleFieldsSundayNotificationDeliveryPreferencesDay
	GetThursday() NotificationDeliveryPreferencesScheduleFieldsThursdayNotificationDeliveryPreferencesDay
	GetTuesday() NotificationDeliveryPreferencesScheduleFieldsTuesdayNotificationDeliveryPreferencesDay
	GetWednesday() NotificationDeliveryPreferencesScheduleFieldsWednesdayNotificationDeliveryPreferencesDay
}

type notificationDeliveryChannelSource interface {
	GetNotificationsDisabled() *bool
	GetSchedule() *NotificationDeliveryPreferencesChannelFieldsScheduleNotificationDeliveryPreferencesSchedule
}

type userSettingsCustomSidebarThemeSource interface {
	GetAccent() []float64
	GetBase() []float64
	GetContrast() int
}

type userSettingsCustomThemeSource interface {
	GetAccent() []float64
	GetBase() []float64
	GetContrast() int
	GetSidebar() *UserSettingsCustomThemeFieldsSidebarUserSettingsCustomSidebarTheme
}

type userSettingsThemeSource interface {
	GetPreset() UserSettingsThemePreset
	GetCustom() *UserSettingsThemeFieldsCustomUserSettingsCustomTheme
}

//nolint:lll
type userSettingsFridayMobile = userSettings_notificationDeliveryPreferences_mobile_schedule_fridayUserSettingsNotificationDeliveryPreferencesMobileNotificationDeliveryPreferencesChannel

//nolint:lll
type userSettingsMondayMobile = userSettings_notificationDeliveryPreferences_mobile_schedule_mondayUserSettingsNotificationDeliveryPreferencesMobileNotificationDeliveryPreferencesChannel

//nolint:lll
type userSettingsSaturdayMobile = userSettings_notificationDeliveryPreferences_mobile_schedule_saturdayUserSettingsNotificationDeliveryPreferencesMobileNotificationDeliveryPreferencesChannel

//nolint:lll
type userSettingsSundayMobile = userSettings_notificationDeliveryPreferences_mobile_schedule_sundayUserSettingsNotificationDeliveryPreferencesMobileNotificationDeliveryPreferencesChannel

//nolint:lll
type userSettingsThursdayMobile = userSettings_notificationDeliveryPreferences_mobile_schedule_thursdayUserSettingsNotificationDeliveryPreferencesMobileNotificationDeliveryPreferencesChannel

//nolint:lll
type userSettingsTuesdayMobile = userSettings_notificationDeliveryPreferences_mobile_schedule_tuesdayUserSettingsNotificationDeliveryPreferencesMobileNotificationDeliveryPreferencesChannel

//nolint:lll
type userSettingsWednesdayMobile = userSettings_notificationDeliveryPreferences_mobile_schedule_wednesdayUserSettingsNotificationDeliveryPreferencesMobileNotificationDeliveryPreferencesChannel

// ListUsers returns visible users.
func ListUsers(ctx context.Context, graphqlClient graphql.Client, limit int) (UserList, error) {
	userPage, err := users(ctx, graphqlClient, intPtr(limit), nil, boolPtr(true), boolPtr(true))
	if err != nil {
		return UserList{}, fmt.Errorf("list users: %w", err)
	}

	summaries := make([]UserSummary, 0, len(userPage.Users.Nodes))
	for _, user := range userPage.Users.Nodes {
		summaries = append(summaries, userSummary(user.UserSummaryFields))
	}

	return UserList{
		Users:       summaries,
		HasNextPage: userPage.Users.PageInfo.HasNextPage,
		EndCursor:   userPage.Users.PageInfo.EndCursor,
	}, nil
}

// GetUserSettings returns the authenticated user's compact settings.
func GetUserSettings(ctx context.Context, graphqlClient graphql.Client) (UserSettingsSummary, error) {
	result, err := userSettings(ctx, graphqlClient)
	if err != nil {
		return UserSettingsSummary{}, fmt.Errorf("get user settings: %w", err)
	}

	return userSettingsSummary(result.UserSettings.UserSettingsSummaryFields), nil
}

// GetUserSettingsNotificationCategoryPreferences returns all notification category preferences.
func GetUserSettingsNotificationCategoryPreferences(
	ctx context.Context,
	graphqlClient graphql.Client,
) (NotificationCategoryPreferences, error) {
	result, err := userSettings_notificationCategoryPreferences(ctx, graphqlClient)
	if err != nil {
		return NotificationCategoryPreferences{}, fmt.Errorf("get user settings notification categories: %w", err)
	}

	return notificationCategoryPreferences(
		result.UserSettings.NotificationCategoryPreferences.NotificationCategoryPreferencesFields,
	), nil
}

// GetUserSettingsNotificationCategoryPreference returns one notification category preference.
//
//nolint:funlen,gocognit,gocyclo // Each case calls a distinct official SDK operation.
func GetUserSettingsNotificationCategoryPreference(
	ctx context.Context,
	graphqlClient graphql.Client,
	category string,
) (NotificationChannelPreference, error) {
	switch normalizedUserSettingsKey(category) {
	case "apps-and-integrations":
		result, err := userSettings_notificationCategoryPreferences_appsAndIntegrations(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(
			&result.UserSettings.NotificationCategoryPreferences.AppsAndIntegrations,
		), nil
	case "assignments":
		result, err := userSettings_notificationCategoryPreferences_assignments(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Assignments), nil
	case "billing":
		result, err := userSettings_notificationCategoryPreferences_billing(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Billing), nil
	case "comments-and-replies":
		result, err := userSettings_notificationCategoryPreferences_commentsAndReplies(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(
			&result.UserSettings.NotificationCategoryPreferences.CommentsAndReplies,
		), nil
	case "customers":
		result, err := userSettings_notificationCategoryPreferences_customers(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Customers), nil
	case "document-changes":
		result, err := userSettings_notificationCategoryPreferences_documentChanges(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.DocumentChanges), nil
	case "feed":
		result, err := userSettings_notificationCategoryPreferences_feed(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Feed), nil
	case "mentions":
		result, err := userSettings_notificationCategoryPreferences_mentions(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Mentions), nil
	case "posts-and-updates":
		result, err := userSettings_notificationCategoryPreferences_postsAndUpdates(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.PostsAndUpdates), nil
	case "reactions":
		result, err := userSettings_notificationCategoryPreferences_reactions(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Reactions), nil
	case "reminders":
		result, err := userSettings_notificationCategoryPreferences_reminders(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Reminders), nil
	case "reviews":
		result, err := userSettings_notificationCategoryPreferences_reviews(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Reviews), nil
	case "status-changes":
		result, err := userSettings_notificationCategoryPreferences_statusChanges(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.StatusChanges), nil
	case "subscriptions":
		result, err := userSettings_notificationCategoryPreferences_subscriptions(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Subscriptions), nil
	case "system":
		result, err := userSettings_notificationCategoryPreferences_system(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.System), nil
	case "triage":
		result, err := userSettings_notificationCategoryPreferences_triage(ctx, graphqlClient)
		if err != nil {
			return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
		}
		return notificationChannelPreference(&result.UserSettings.NotificationCategoryPreferences.Triage), nil
	default:
		return NotificationChannelPreference{}, fmt.Errorf("unknown user settings notification category %q", category)
	}
}

// GetUserSettingsNotificationChannelPreferences returns the top-level notification channel preferences.
func GetUserSettingsNotificationChannelPreferences(
	ctx context.Context,
	graphqlClient graphql.Client,
) (NotificationChannelPreference, error) {
	result, err := userSettings_notificationChannelPreferences(ctx, graphqlClient)
	if err != nil {
		return NotificationChannelPreference{}, fmt.Errorf("get user settings notification channels: %w", err)
	}

	return notificationChannelPreference(&result.UserSettings.NotificationChannelPreferences), nil
}

// GetUserSettingsNotificationDeliveryPreferences returns notification delivery preferences.
func GetUserSettingsNotificationDeliveryPreferences(
	ctx context.Context,
	graphqlClient graphql.Client,
) (NotificationDeliveryPreferences, error) {
	result, err := userSettings_notificationDeliveryPreferences(ctx, graphqlClient)
	if err != nil {
		return NotificationDeliveryPreferences{}, fmt.Errorf("get user settings notification delivery: %w", err)
	}

	return notificationDeliveryPreferences(
		result.UserSettings.NotificationDeliveryPreferences.NotificationDeliveryPreferencesFields,
	), nil
}

// GetUserSettingsMobileDeliveryPreferences returns mobile notification delivery preferences.
//
//nolint:nilnil // A nil channel is a valid nullable GraphQL result for this read.
func GetUserSettingsMobileDeliveryPreferences(
	ctx context.Context,
	graphqlClient graphql.Client,
) (*NotificationDeliveryChannel, error) {
	result, err := userSettings_notificationDeliveryPreferences_mobile(ctx, graphqlClient)
	if err != nil {
		return nil, fmt.Errorf("get user settings mobile delivery: %w", err)
	}
	if result.UserSettings.NotificationDeliveryPreferences.Mobile == nil {
		return nil, nil
	}

	return notificationDeliveryChannel(result.UserSettings.NotificationDeliveryPreferences.Mobile), nil
}

// GetUserSettingsMobileSchedule returns the mobile notification delivery schedule.
//
//nolint:nilnil // A nil schedule is a valid nullable GraphQL result for this read.
func GetUserSettingsMobileSchedule(
	ctx context.Context,
	graphqlClient graphql.Client,
) (*NotificationDeliverySchedule, error) {
	result, err := userSettings_notificationDeliveryPreferences_mobile_schedule(ctx, graphqlClient)
	if err != nil {
		return nil, fmt.Errorf("get user settings mobile schedule: %w", err)
	}
	mobile := result.UserSettings.NotificationDeliveryPreferences.Mobile
	if mobile == nil || mobile.Schedule == nil {
		return nil, nil
	}

	return notificationDeliverySchedule(mobile.Schedule), nil
}

// GetUserSettingsMobileScheduleDay returns one mobile notification delivery schedule day.
func GetUserSettingsMobileScheduleDay(
	ctx context.Context,
	graphqlClient graphql.Client,
	day string,
) (NotificationDeliveryDay, error) {
	switch normalizedUserSettingsKey(day) {
	case "friday":
		result, err := userSettings_notificationDeliveryPreferences_mobile_schedule_friday(ctx, graphqlClient)
		if err != nil {
			return NotificationDeliveryDay{}, fmt.Errorf("get user settings mobile schedule %s: %w", day, err)
		}
		return notificationDeliveryDayFromMobileFriday(result.UserSettings.NotificationDeliveryPreferences.Mobile), nil
	case "monday":
		result, err := userSettings_notificationDeliveryPreferences_mobile_schedule_monday(ctx, graphqlClient)
		if err != nil {
			return NotificationDeliveryDay{}, fmt.Errorf("get user settings mobile schedule %s: %w", day, err)
		}
		return notificationDeliveryDayFromMobileMonday(result.UserSettings.NotificationDeliveryPreferences.Mobile), nil
	case "saturday":
		result, err := userSettings_notificationDeliveryPreferences_mobile_schedule_saturday(ctx, graphqlClient)
		if err != nil {
			return NotificationDeliveryDay{}, fmt.Errorf("get user settings mobile schedule %s: %w", day, err)
		}
		return notificationDeliveryDayFromMobileSaturday(
			result.UserSettings.NotificationDeliveryPreferences.Mobile,
		), nil
	case "sunday":
		result, err := userSettings_notificationDeliveryPreferences_mobile_schedule_sunday(ctx, graphqlClient)
		if err != nil {
			return NotificationDeliveryDay{}, fmt.Errorf("get user settings mobile schedule %s: %w", day, err)
		}
		return notificationDeliveryDayFromMobileSunday(result.UserSettings.NotificationDeliveryPreferences.Mobile), nil
	case "thursday":
		result, err := userSettings_notificationDeliveryPreferences_mobile_schedule_thursday(ctx, graphqlClient)
		if err != nil {
			return NotificationDeliveryDay{}, fmt.Errorf("get user settings mobile schedule %s: %w", day, err)
		}
		return notificationDeliveryDayFromMobileThursday(
			result.UserSettings.NotificationDeliveryPreferences.Mobile,
		), nil
	case "tuesday":
		result, err := userSettings_notificationDeliveryPreferences_mobile_schedule_tuesday(ctx, graphqlClient)
		if err != nil {
			return NotificationDeliveryDay{}, fmt.Errorf("get user settings mobile schedule %s: %w", day, err)
		}
		return notificationDeliveryDayFromMobileTuesday(result.UserSettings.NotificationDeliveryPreferences.Mobile), nil
	case "wednesday":
		result, err := userSettings_notificationDeliveryPreferences_mobile_schedule_wednesday(ctx, graphqlClient)
		if err != nil {
			return NotificationDeliveryDay{}, fmt.Errorf("get user settings mobile schedule %s: %w", day, err)
		}
		return notificationDeliveryDayFromMobileWednesday(
			result.UserSettings.NotificationDeliveryPreferences.Mobile,
		), nil
	default:
		return NotificationDeliveryDay{}, fmt.Errorf("unknown user settings mobile schedule day %q", day)
	}
}

// GetUserSettingsTheme returns the user's theme for one device and mode.
//
//nolint:nilnil // A nil theme is a valid nullable GraphQL result for this read.
func GetUserSettingsTheme(
	ctx context.Context,
	graphqlClient graphql.Client,
	deviceType string,
	mode string,
) (*UserSettingsThemeSummary, error) {
	deviceTypeValue, modeValue, err := userSettingsThemeArgs(deviceType, mode)
	if err != nil {
		return nil, err
	}
	result, err := userSettings_theme(ctx, graphqlClient, deviceTypeValue, modeValue)
	if err != nil {
		return nil, fmt.Errorf("get user settings theme: %w", err)
	}
	if result.UserSettings.Theme == nil {
		return nil, nil
	}

	return userSettingsTheme(result.UserSettings.Theme), nil
}

// GetUserSettingsCustomTheme returns the user's custom theme for one device and mode.
//
//nolint:nilnil // A nil custom theme is a valid nullable GraphQL result for this read.
func GetUserSettingsCustomTheme(
	ctx context.Context,
	graphqlClient graphql.Client,
	deviceType string,
	mode string,
) (*UserSettingsCustomTheme, error) {
	deviceTypeValue, modeValue, err := userSettingsThemeArgs(deviceType, mode)
	if err != nil {
		return nil, err
	}
	result, err := userSettings_theme_custom(ctx, graphqlClient, deviceTypeValue, modeValue)
	if err != nil {
		return nil, fmt.Errorf("get user settings custom theme: %w", err)
	}
	if result.UserSettings.Theme == nil || result.UserSettings.Theme.Custom == nil {
		return nil, nil
	}

	return userSettingsCustomTheme(result.UserSettings.Theme.Custom), nil
}

// GetUserSettingsCustomSidebarTheme returns the user's custom sidebar theme for one device and mode.
//
//nolint:nilnil // A nil custom sidebar theme is a valid nullable GraphQL result for this read.
func GetUserSettingsCustomSidebarTheme(
	ctx context.Context,
	graphqlClient graphql.Client,
	deviceType string,
	mode string,
) (*UserSettingsCustomSidebarTheme, error) {
	deviceTypeValue, modeValue, err := userSettingsThemeArgs(deviceType, mode)
	if err != nil {
		return nil, err
	}
	result, err := userSettings_theme_custom_sidebar(ctx, graphqlClient, deviceTypeValue, modeValue)
	if err != nil {
		return nil, fmt.Errorf("get user settings custom sidebar theme: %w", err)
	}
	if result.UserSettings.Theme == nil ||
		result.UserSettings.Theme.Custom == nil ||
		result.UserSettings.Theme.Custom.Sidebar == nil {
		return nil, nil
	}

	return userSettingsCustomSidebarTheme(result.UserSettings.Theme.Custom.Sidebar), nil
}

// GetUserByID returns one User by id.
func GetUserByID(ctx context.Context, graphqlClient graphql.Client, id string) (UserSummary, error) {
	userResult, err := user(ctx, graphqlClient, id)
	if err != nil {
		return UserSummary{}, fmt.Errorf("get user %s: %w", id, err)
	}

	return userSummary(userResult.User.UserSummaryFields), nil
}

// GetViewerUser returns the authenticated User.
func GetViewerUser(ctx context.Context, graphqlClient graphql.Client) (UserSummary, error) {
	userResult, err := viewer(ctx, graphqlClient)
	if err != nil {
		return UserSummary{}, fmt.Errorf("get viewer user: %w", err)
	}

	return userSummary(userResult.Viewer.UserSummaryFields), nil
}

// ListViewerDrafts returns the authenticated user's saved draft metadata.
func ListViewerDrafts(ctx context.Context, graphqlClient graphql.Client, limit int) (DraftList, error) {
	draftPage, err := viewer_drafts(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return DraftList{}, fmt.Errorf("list viewer drafts: %w", err)
	}

	summaries := make([]DraftSummary, 0, len(draftPage.Viewer.Drafts.Nodes))
	for _, draft := range draftPage.Viewer.Drafts.Nodes {
		summaries = append(summaries, draftSummary(draft.DraftSummaryFields))
	}

	return DraftList{
		Drafts:      summaries,
		HasNextPage: draftPage.Viewer.Drafts.PageInfo.HasNextPage,
		EndCursor:   draftPage.Viewer.Drafts.PageInfo.EndCursor,
	}, nil
}

// ListUserAssignedIssues returns issues assigned to one User.
func ListUserAssignedIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	result, err := user_assignedIssues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list user assigned issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.User.AssignedIssues.Nodes))
	for _, issue := range result.User.AssignedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.User.AssignedIssues.PageInfo.HasNextPage,
		EndCursor:   result.User.AssignedIssues.PageInfo.EndCursor,
	}, nil
}

// ListUserCreatedIssues returns issues created by one User.
func ListUserCreatedIssues(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (IssueList, error) {
	result, err := user_createdIssues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list user created issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.User.CreatedIssues.Nodes))
	for _, issue := range result.User.CreatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.User.CreatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.User.CreatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListUserDelegatedIssues returns issues delegated to one User.
func ListUserDelegatedIssues(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (IssueList, error) {
	result, err := user_delegatedIssues(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list user delegated issues %s: %w", id, err)
	}

	issues := make([]IssueSummary, 0, len(result.User.DelegatedIssues.Nodes))
	for _, issue := range result.User.DelegatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.User.DelegatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.User.DelegatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListUserTeamMemberships returns TeamMemberships associated with one User.
func ListUserTeamMemberships(
	ctx context.Context,
	graphqlClient graphql.Client,
	id string,
	limit int,
) (TeamMembershipList, error) {
	result, err := user_teamMemberships(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamMembershipList{}, fmt.Errorf("list user team memberships %s: %w", id, err)
	}

	memberships := make([]TeamMembershipSummary, 0, len(result.User.TeamMemberships.Nodes))
	for _, membership := range result.User.TeamMemberships.Nodes {
		memberships = append(memberships, teamMembershipSummary(membership.TeamMembershipSummaryFields))
	}

	return TeamMembershipList{
		Memberships: memberships,
		HasNextPage: result.User.TeamMemberships.PageInfo.HasNextPage,
		EndCursor:   result.User.TeamMemberships.PageInfo.EndCursor,
	}, nil
}

// ListUserTeams returns Teams associated with one User.
func ListUserTeams(ctx context.Context, graphqlClient graphql.Client, id string, limit int) (TeamList, error) {
	result, err := user_teams(ctx, graphqlClient, id, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamList{}, fmt.Errorf("list user teams %s: %w", id, err)
	}

	teams := make([]TeamSummary, 0, len(result.User.Teams.Nodes))
	for _, team := range result.User.Teams.Nodes {
		teams = append(teams, teamSummary(team.TeamSummaryFields))
	}

	return TeamList{
		Teams:       teams,
		HasNextPage: result.User.Teams.PageInfo.HasNextPage,
		EndCursor:   result.User.Teams.PageInfo.EndCursor,
	}, nil
}

// ListViewerAssignedIssues returns issues assigned to the authenticated User.
func ListViewerAssignedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	result, err := viewer_assignedIssues(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list viewer assigned issues: %w", err)
	}

	issues := make([]IssueSummary, 0, len(result.Viewer.AssignedIssues.Nodes))
	for _, issue := range result.Viewer.AssignedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Viewer.AssignedIssues.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.AssignedIssues.PageInfo.EndCursor,
	}, nil
}

// ListViewerCreatedIssues returns issues created by the authenticated User.
func ListViewerCreatedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	result, err := viewer_createdIssues(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list viewer created issues: %w", err)
	}

	issues := make([]IssueSummary, 0, len(result.Viewer.CreatedIssues.Nodes))
	for _, issue := range result.Viewer.CreatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Viewer.CreatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.CreatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListViewerDelegatedIssues returns issues delegated to the authenticated User.
func ListViewerDelegatedIssues(ctx context.Context, graphqlClient graphql.Client, limit int) (IssueList, error) {
	result, err := viewer_delegatedIssues(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return IssueList{}, fmt.Errorf("list viewer delegated issues: %w", err)
	}

	issues := make([]IssueSummary, 0, len(result.Viewer.DelegatedIssues.Nodes))
	for _, issue := range result.Viewer.DelegatedIssues.Nodes {
		issues = append(issues, issueSummaryFromFields(issue.IssueSummaryFields))
	}

	return IssueList{
		Issues:      issues,
		HasNextPage: result.Viewer.DelegatedIssues.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.DelegatedIssues.PageInfo.EndCursor,
	}, nil
}

// ListViewerTeamMemberships returns TeamMemberships associated with the authenticated User.
func ListViewerTeamMemberships(
	ctx context.Context,
	graphqlClient graphql.Client,
	limit int,
) (TeamMembershipList, error) {
	result, err := viewer_teamMemberships(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamMembershipList{}, fmt.Errorf("list viewer team memberships: %w", err)
	}

	memberships := make([]TeamMembershipSummary, 0, len(result.Viewer.TeamMemberships.Nodes))
	for _, membership := range result.Viewer.TeamMemberships.Nodes {
		memberships = append(memberships, teamMembershipSummary(membership.TeamMembershipSummaryFields))
	}

	return TeamMembershipList{
		Memberships: memberships,
		HasNextPage: result.Viewer.TeamMemberships.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.TeamMemberships.PageInfo.EndCursor,
	}, nil
}

// ListViewerTeams returns Teams associated with the authenticated User.
func ListViewerTeams(ctx context.Context, graphqlClient graphql.Client, limit int) (TeamList, error) {
	result, err := viewer_teams(ctx, graphqlClient, intPtr(limit), nil, boolPtr(false))
	if err != nil {
		return TeamList{}, fmt.Errorf("list viewer teams: %w", err)
	}

	teams := make([]TeamSummary, 0, len(result.Viewer.Teams.Nodes))
	for _, team := range result.Viewer.Teams.Nodes {
		teams = append(teams, teamSummary(team.TeamSummaryFields))
	}

	return TeamList{
		Teams:       teams,
		HasNextPage: result.Viewer.Teams.PageInfo.HasNextPage,
		EndCursor:   result.Viewer.Teams.PageInfo.EndCursor,
	}, nil
}

func userSummary(user UserSummaryFields) UserSummary {
	return UserSummary{
		ID:          user.Id,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Active:      user.Active,
		Guest:       user.Guest,
		Admin:       user.Admin,
	}
}

func draftSummary(draft DraftSummaryFields) DraftSummary {
	summary := DraftSummary{
		ID:         draft.Id,
		CreatedAt:  draft.CreatedAt,
		UpdatedAt:  draft.UpdatedAt,
		ArchivedAt: stringValue(draft.ArchivedAt),
	}
	switch {
	case draft.Issue != nil:
		summary.ParentType = "issue"
		summary.ParentID = draft.Issue.Id
		summary.ParentKey = draft.Issue.Identifier
		summary.ParentTitle = draft.Issue.Title
	case draft.Project != nil:
		summary.ParentType = "project"
		summary.ParentID = draft.Project.Id
		summary.ParentTitle = draft.Project.Name
	case draft.ProjectUpdate != nil:
		summary.ParentType = "project_update"
		summary.ParentID = draft.ProjectUpdate.Id
	case draft.Initiative != nil:
		summary.ParentType = "initiative"
		summary.ParentID = draft.Initiative.Id
		summary.ParentTitle = draft.Initiative.Name
	case draft.InitiativeUpdate != nil:
		summary.ParentType = "initiative_update"
		summary.ParentID = draft.InitiativeUpdate.Id
	case draft.ParentComment != nil:
		summary.ParentType = "comment"
		summary.ParentID = draft.ParentComment.Id
	case draft.CustomerNeed != nil:
		summary.ParentType = "customer_need"
		summary.ParentID = draft.CustomerNeed.Id
	case draft.Team != nil:
		summary.ParentType = "team"
		summary.ParentID = draft.Team.Id
		summary.ParentKey = draft.Team.Key
		summary.ParentTitle = draft.Team.Name
	default:
		summary.ParentType = "unknown"
	}

	return summary
}

func userSettingsSummary(settings UserSettingsSummaryFields) UserSettingsSummary {
	feedSummarySchedule := ""
	if settings.FeedSummarySchedule != nil {
		feedSummarySchedule = string(*settings.FeedSummarySchedule)
	}

	return UserSettingsSummary{
		ID:                              settings.Id,
		UserID:                          settings.User.Id,
		CreatedAt:                       settings.CreatedAt,
		UpdatedAt:                       settings.UpdatedAt,
		ArchivedAt:                      settings.ArchivedAt,
		AutoAssignToSelf:                settings.AutoAssignToSelf,
		FeedLastSeenTime:                settings.FeedLastSeenTime,
		FeedSummarySchedule:             feedSummarySchedule,
		ShowFullUserNames:               settings.ShowFullUserNames,
		SubscribedToChangelog:           settings.SubscribedToChangelog,
		SubscribedToDPA:                 settings.SubscribedToDPA,
		SubscribedToInviteAccepted:      settings.SubscribedToInviteAccepted,
		SubscribedToPrivacyLegalUpdates: settings.SubscribedToPrivacyLegalUpdates,
		NotificationCategoryPreferences: notificationCategoryPreferences(
			settings.NotificationCategoryPreferences.NotificationCategoryPreferencesFields,
		),
		NotificationChannelPreferences: notificationChannelPreference(&settings.NotificationChannelPreferences),
		NotificationDeliveryPreferences: notificationDeliveryPreferences(
			settings.NotificationDeliveryPreferences.NotificationDeliveryPreferencesFields,
		),
	}
}

func notificationCategoryPreferences(
	preferences NotificationCategoryPreferencesFields,
) NotificationCategoryPreferences {
	return NotificationCategoryPreferences{
		AppsAndIntegrations: notificationChannelPreference(&preferences.AppsAndIntegrations),
		Assignments:         notificationChannelPreference(&preferences.Assignments),
		Billing:             notificationChannelPreference(&preferences.Billing),
		CommentsAndReplies:  notificationChannelPreference(&preferences.CommentsAndReplies),
		Customers:           notificationChannelPreference(&preferences.Customers),
		DocumentChanges:     notificationChannelPreference(&preferences.DocumentChanges),
		Feed:                notificationChannelPreference(&preferences.Feed),
		Mentions:            notificationChannelPreference(&preferences.Mentions),
		PostsAndUpdates:     notificationChannelPreference(&preferences.PostsAndUpdates),
		Reactions:           notificationChannelPreference(&preferences.Reactions),
		Reminders:           notificationChannelPreference(&preferences.Reminders),
		Reviews:             notificationChannelPreference(&preferences.Reviews),
		StatusChanges:       notificationChannelPreference(&preferences.StatusChanges),
		Subscriptions:       notificationChannelPreference(&preferences.Subscriptions),
		System:              notificationChannelPreference(&preferences.System),
		Triage:              notificationChannelPreference(&preferences.Triage),
	}
}

func notificationChannelPreference(source notificationChannelPreferenceSource) NotificationChannelPreference {
	return NotificationChannelPreference{
		Desktop: source.GetDesktop(),
		Email:   source.GetEmail(),
		Mobile:  source.GetMobile(),
		Slack:   source.GetSlack(),
	}
}

func notificationDeliveryPreferences(
	preferences NotificationDeliveryPreferencesFields,
) NotificationDeliveryPreferences {
	if preferences.Mobile == nil {
		return NotificationDeliveryPreferences{}
	}

	return NotificationDeliveryPreferences{Mobile: notificationDeliveryChannel(preferences.Mobile)}
}

func notificationDeliveryChannel(source notificationDeliveryChannelSource) *NotificationDeliveryChannel {
	channel := NotificationDeliveryChannel{
		NotificationsDisabled: source.GetNotificationsDisabled(),
	}
	if source.GetSchedule() != nil {
		channel.Schedule = notificationDeliverySchedule(source.GetSchedule())
	}

	return &channel
}

func notificationDeliverySchedule(source notificationDeliveryScheduleSource) *NotificationDeliverySchedule {
	friday := source.GetFriday()
	monday := source.GetMonday()
	saturday := source.GetSaturday()
	sunday := source.GetSunday()
	thursday := source.GetThursday()
	tuesday := source.GetTuesday()
	wednesday := source.GetWednesday()

	return &NotificationDeliverySchedule{
		Disabled:  source.GetDisabled(),
		Friday:    notificationDeliveryDay(&friday),
		Monday:    notificationDeliveryDay(&monday),
		Saturday:  notificationDeliveryDay(&saturday),
		Sunday:    notificationDeliveryDay(&sunday),
		Thursday:  notificationDeliveryDay(&thursday),
		Tuesday:   notificationDeliveryDay(&tuesday),
		Wednesday: notificationDeliveryDay(&wednesday),
	}
}

func notificationDeliveryDay(source notificationDeliveryDaySource) NotificationDeliveryDay {
	return NotificationDeliveryDay{
		Start: source.GetStart(),
		End:   source.GetEnd(),
	}
}

func notificationDeliveryDayFromMobileFriday(mobile *userSettingsFridayMobile) NotificationDeliveryDay {
	if mobile == nil || mobile.Schedule == nil {
		return NotificationDeliveryDay{}
	}

	return notificationDeliveryDay(&mobile.Schedule.Friday)
}

func notificationDeliveryDayFromMobileMonday(mobile *userSettingsMondayMobile) NotificationDeliveryDay {
	if mobile == nil || mobile.Schedule == nil {
		return NotificationDeliveryDay{}
	}

	return notificationDeliveryDay(&mobile.Schedule.Monday)
}

func notificationDeliveryDayFromMobileSaturday(mobile *userSettingsSaturdayMobile) NotificationDeliveryDay {
	if mobile == nil || mobile.Schedule == nil {
		return NotificationDeliveryDay{}
	}

	return notificationDeliveryDay(&mobile.Schedule.Saturday)
}

func notificationDeliveryDayFromMobileSunday(mobile *userSettingsSundayMobile) NotificationDeliveryDay {
	if mobile == nil || mobile.Schedule == nil {
		return NotificationDeliveryDay{}
	}

	return notificationDeliveryDay(&mobile.Schedule.Sunday)
}

func notificationDeliveryDayFromMobileThursday(mobile *userSettingsThursdayMobile) NotificationDeliveryDay {
	if mobile == nil || mobile.Schedule == nil {
		return NotificationDeliveryDay{}
	}

	return notificationDeliveryDay(&mobile.Schedule.Thursday)
}

func notificationDeliveryDayFromMobileTuesday(mobile *userSettingsTuesdayMobile) NotificationDeliveryDay {
	if mobile == nil || mobile.Schedule == nil {
		return NotificationDeliveryDay{}
	}

	return notificationDeliveryDay(&mobile.Schedule.Tuesday)
}

func notificationDeliveryDayFromMobileWednesday(mobile *userSettingsWednesdayMobile) NotificationDeliveryDay {
	if mobile == nil || mobile.Schedule == nil {
		return NotificationDeliveryDay{}
	}

	return notificationDeliveryDay(&mobile.Schedule.Wednesday)
}

func userSettingsTheme(source userSettingsThemeSource) *UserSettingsThemeSummary {
	theme := UserSettingsThemeSummary{Preset: string(source.GetPreset())}
	if source.GetCustom() != nil {
		theme.Custom = userSettingsCustomTheme(source.GetCustom())
	}

	return &theme
}

func userSettingsCustomTheme(source userSettingsCustomThemeSource) *UserSettingsCustomTheme {
	theme := UserSettingsCustomTheme{
		Accent:   source.GetAccent(),
		Base:     source.GetBase(),
		Contrast: source.GetContrast(),
	}
	if source.GetSidebar() != nil {
		theme.Sidebar = userSettingsCustomSidebarTheme(source.GetSidebar())
	}

	return &theme
}

func userSettingsCustomSidebarTheme(
	source userSettingsCustomSidebarThemeSource,
) *UserSettingsCustomSidebarTheme {
	return &UserSettingsCustomSidebarTheme{
		Accent:   source.GetAccent(),
		Base:     source.GetBase(),
		Contrast: source.GetContrast(),
	}
}

func userSettingsThemeArgs(
	deviceType string,
	mode string,
) (*UserSettingsThemeDeviceType, *UserSettingsThemeMode, error) {
	deviceTypeValue, err := parseUserSettingsThemeDeviceType(deviceType)
	if err != nil {
		return nil, nil, err
	}
	modeValue, err := parseUserSettingsThemeMode(mode)
	if err != nil {
		return nil, nil, err
	}

	return deviceTypeValue, modeValue, nil
}

func parseUserSettingsThemeDeviceType(value string) (*UserSettingsThemeDeviceType, error) {
	switch normalizedUserSettingsKey(value) {
	case "", "desktop":
		deviceType := UserSettingsThemeDeviceTypeDesktop
		return &deviceType, nil
	case "mobile-web", "mobileweb":
		deviceType := UserSettingsThemeDeviceTypeMobileweb
		return &deviceType, nil
	default:
		return nil, fmt.Errorf("invalid theme device type %q: use desktop or mobile-web", value)
	}
}

func parseUserSettingsThemeMode(value string) (*UserSettingsThemeMode, error) {
	switch normalizedUserSettingsKey(value) {
	case "", "light":
		mode := UserSettingsThemeModeLight
		return &mode, nil
	case "dark":
		mode := UserSettingsThemeModeDark
		return &mode, nil
	default:
		return nil, fmt.Errorf("invalid theme mode %q: use light or dark", value)
	}
}

func normalizedUserSettingsKey(value string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(value)), "_", "-")
}
