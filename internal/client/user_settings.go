package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

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
	GetFriday() gql.NotificationDeliveryPreferencesScheduleFieldsFridayNotificationDeliveryPreferencesDay
	GetMonday() gql.NotificationDeliveryPreferencesScheduleFieldsMondayNotificationDeliveryPreferencesDay
	GetSaturday() gql.NotificationDeliveryPreferencesScheduleFieldsSaturdayNotificationDeliveryPreferencesDay
	GetSunday() gql.NotificationDeliveryPreferencesScheduleFieldsSundayNotificationDeliveryPreferencesDay
	GetThursday() gql.NotificationDeliveryPreferencesScheduleFieldsThursdayNotificationDeliveryPreferencesDay
	GetTuesday() gql.NotificationDeliveryPreferencesScheduleFieldsTuesdayNotificationDeliveryPreferencesDay
	GetWednesday() gql.NotificationDeliveryPreferencesScheduleFieldsWednesdayNotificationDeliveryPreferencesDay
}

type notificationDeliveryChannelSource interface {
	GetNotificationsDisabled() *bool
	GetSchedule() *gql.NotificationDeliveryPreferencesChannelFieldsScheduleNotificationDeliveryPreferencesSchedule
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
	GetSidebar() *gql.UserSettingsCustomThemeFieldsSidebarUserSettingsCustomSidebarTheme
}

type userSettingsThemeSource interface {
	GetPreset() gql.UserSettingsThemePreset
	GetCustom() *gql.UserSettingsThemeFieldsCustomUserSettingsCustomTheme
}

// GetUserSettings returns the authenticated user's compact settings.
func GetUserSettings(ctx context.Context, graphqlClient graphql.Client) (UserSettingsSummary, error) {
	result, err := gql.XUserSettings(ctx, graphqlClient)
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
	result, err := gql.XUserSettings_notificationCategoryPreferences(ctx, graphqlClient)
	if err != nil {
		return NotificationCategoryPreferences{}, fmt.Errorf("get user settings notification categories: %w", err)
	}

	return notificationCategoryPreferences(
		result.UserSettings.NotificationCategoryPreferences.NotificationCategoryPreferencesFields,
	), nil
}

// GetUserSettingsNotificationCategoryPreference returns one notification category preference.
func GetUserSettingsNotificationCategoryPreference(
	ctx context.Context,
	graphqlClient graphql.Client,
	category string,
) (NotificationChannelPreference, error) {
	selectPreference, ok := notificationCategoryPreferenceSelectors[normalizedUserSettingsKey(category)]
	if !ok {
		return NotificationChannelPreference{}, fmt.Errorf("unknown user settings notification category %q", category)
	}

	result, err := gql.XUserSettings_notificationCategoryPreferences(ctx, graphqlClient)
	if err != nil {
		return NotificationChannelPreference{}, fmt.Errorf("get user settings category %s: %w", category, err)
	}
	preferences := notificationCategoryPreferences(
		result.UserSettings.NotificationCategoryPreferences.NotificationCategoryPreferencesFields,
	)

	return selectPreference(preferences), nil
}

// notificationCategoryPreferenceSelectors picks one category out of the full
// preference matrix; the per-category read fetches the plural query once.
var notificationCategoryPreferenceSelectors = map[string]func(
	NotificationCategoryPreferences,
) NotificationChannelPreference{
	"apps-and-integrations": func(p NotificationCategoryPreferences) NotificationChannelPreference {
		return p.AppsAndIntegrations
	},
	"assignments": func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Assignments },
	"billing":     func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Billing },
	"comments-and-replies": func(p NotificationCategoryPreferences) NotificationChannelPreference {
		return p.CommentsAndReplies
	},
	"customers": func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Customers },
	"document-changes": func(p NotificationCategoryPreferences) NotificationChannelPreference {
		return p.DocumentChanges
	},
	"feed":     func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Feed },
	"mentions": func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Mentions },
	"posts-and-updates": func(p NotificationCategoryPreferences) NotificationChannelPreference {
		return p.PostsAndUpdates
	},
	"reactions": func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Reactions },
	"reminders": func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Reminders },
	"reviews":   func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Reviews },
	"status-changes": func(p NotificationCategoryPreferences) NotificationChannelPreference {
		return p.StatusChanges
	},
	"subscriptions": func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Subscriptions },
	"system":        func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.System },
	"triage":        func(p NotificationCategoryPreferences) NotificationChannelPreference { return p.Triage },
}

// GetUserSettingsNotificationChannelPreferences returns the top-level notification channel preferences.
func GetUserSettingsNotificationChannelPreferences(
	ctx context.Context,
	graphqlClient graphql.Client,
) (NotificationChannelPreference, error) {
	result, err := gql.XUserSettings_notificationChannelPreferences(ctx, graphqlClient)
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
	result, err := gql.XUserSettings_notificationDeliveryPreferences(ctx, graphqlClient)
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
	result, err := gql.XUserSettings_notificationDeliveryPreferences_mobile(ctx, graphqlClient)
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
	result, err := gql.XUserSettings_notificationDeliveryPreferences_mobile_schedule(ctx, graphqlClient)
	if err != nil {
		return nil, fmt.Errorf("get user settings mobile schedule: %w", err)
	}
	mobile := result.UserSettings.NotificationDeliveryPreferences.Mobile
	if mobile == nil || mobile.Schedule == nil {
		return nil, nil
	}

	return notificationDeliverySchedule(mobile.Schedule), nil
}

// mobileScheduleDayAccessors maps a normalized day key onto its field of the
// full-week mobile notification delivery schedule.
var mobileScheduleDayAccessors = map[string]func(NotificationDeliverySchedule) NotificationDeliveryDay{
	"friday":    func(schedule NotificationDeliverySchedule) NotificationDeliveryDay { return schedule.Friday },
	"monday":    func(schedule NotificationDeliverySchedule) NotificationDeliveryDay { return schedule.Monday },
	"saturday":  func(schedule NotificationDeliverySchedule) NotificationDeliveryDay { return schedule.Saturday },
	"sunday":    func(schedule NotificationDeliverySchedule) NotificationDeliveryDay { return schedule.Sunday },
	"thursday":  func(schedule NotificationDeliverySchedule) NotificationDeliveryDay { return schedule.Thursday },
	"tuesday":   func(schedule NotificationDeliverySchedule) NotificationDeliveryDay { return schedule.Tuesday },
	"wednesday": func(schedule NotificationDeliverySchedule) NotificationDeliveryDay { return schedule.Wednesday },
}

// GetUserSettingsMobileScheduleDay returns one mobile notification delivery schedule day.
func GetUserSettingsMobileScheduleDay(
	ctx context.Context,
	graphqlClient graphql.Client,
	day string,
) (NotificationDeliveryDay, error) {
	dayOfSchedule, ok := mobileScheduleDayAccessors[normalizedUserSettingsKey(day)]
	if !ok {
		return NotificationDeliveryDay{}, fmt.Errorf("unknown user settings mobile schedule day %q", day)
	}
	result, err := gql.XUserSettings_notificationDeliveryPreferences_mobile_schedule(ctx, graphqlClient)
	if err != nil {
		return NotificationDeliveryDay{}, fmt.Errorf("get user settings mobile schedule %s: %w", day, err)
	}
	mobile := result.UserSettings.NotificationDeliveryPreferences.Mobile
	if mobile == nil || mobile.Schedule == nil {
		return NotificationDeliveryDay{}, nil
	}

	return dayOfSchedule(*notificationDeliverySchedule(mobile.Schedule)), nil
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
	result, err := gql.XUserSettings_theme(ctx, graphqlClient, deviceTypeValue, modeValue)
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
	result, err := gql.XUserSettings_theme_custom(ctx, graphqlClient, deviceTypeValue, modeValue)
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
	result, err := gql.XUserSettings_theme_custom_sidebar(ctx, graphqlClient, deviceTypeValue, modeValue)
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

func userSettingsSummary(settings gql.UserSettingsSummaryFields) UserSettingsSummary {
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
	preferences gql.NotificationCategoryPreferencesFields,
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
	preferences gql.NotificationDeliveryPreferencesFields,
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
) (*gql.UserSettingsThemeDeviceType, *gql.UserSettingsThemeMode, error) {
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

func parseUserSettingsThemeDeviceType(value string) (*gql.UserSettingsThemeDeviceType, error) {
	switch normalizedUserSettingsKey(value) {
	case "", "desktop":
		deviceType := gql.UserSettingsThemeDeviceTypeDesktop
		return &deviceType, nil
	case "mobile-web", "mobileweb":
		deviceType := gql.UserSettingsThemeDeviceTypeMobileweb
		return &deviceType, nil
	default:
		return nil, fmt.Errorf("invalid theme device type %q: use desktop or mobile-web", value)
	}
}

func parseUserSettingsThemeMode(value string) (*gql.UserSettingsThemeMode, error) {
	switch normalizedUserSettingsKey(value) {
	case "", "light":
		mode := gql.UserSettingsThemeModeLight
		return &mode, nil
	case "dark":
		mode := gql.UserSettingsThemeModeDark
		return &mode, nil
	default:
		return nil, fmt.Errorf("invalid theme mode %q: use light or dark", value)
	}
}

func normalizedUserSettingsKey(value string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(value)), "_", "-")
}
