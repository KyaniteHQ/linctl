package client

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ClientReadScenarios_return_user_settings(t *testing.T) {
	graphqlClient := fakeGraphQLClient{
		"userSettings": `{"userSettings":` + userSettingsJSON() + `}`,
		"userSettings_notificationCategoryPreferences": `{"userSettings":{"notificationCategoryPreferences":` +
			notificationCategoriesJSON() + `}}`,
		"userSettings_notificationChannelPreferences": `{"userSettings":{"notificationChannelPreferences":` +
			notificationChannelJSON() + `}}`,
		"userSettings_notificationDeliveryPreferences": `{"userSettings":{"notificationDeliveryPreferences":` +
			notificationDeliveryPreferencesJSON() + `}}`,
		"userSettings_notificationDeliveryPreferences_mobile": `{"userSettings":{"notificationDeliveryPreferences":{"mobile":` +
			notificationDeliveryChannelJSON() + `}}}`,
		"userSettings_notificationDeliveryPreferences_mobile_schedule": `{"userSettings":{"notificationDeliveryPreferences":{"mobile":{"schedule":` +
			notificationDeliveryScheduleJSON() + `}}}}`,
		"userSettings_theme":                `{"userSettings":{"theme":` + userSettingsThemeJSON(true) + `}}`,
		"userSettings_theme_custom":         `{"userSettings":{"theme":{"custom":` + userSettingsCustomThemeJSON(true) + `}}}`,
		"userSettings_theme_custom_sidebar": `{"userSettings":{"theme":{"custom":{"sidebar":` + userSettingsCustomSidebarThemeJSON() + `}}}}`,
	}

	settings, err := GetUserSettings(context.Background(), graphqlClient)
	require.NoError(t, err)
	require.Equal(t, "settings-id", settings.ID)
	require.Equal(t, "daily", settings.FeedSummarySchedule)
	require.True(t, settings.NotificationChannelPreferences.Desktop)

	categories, err := GetUserSettingsNotificationCategoryPreferences(context.Background(), graphqlClient)
	require.NoError(t, err)
	require.True(t, categories.Assignments.Mobile)

	for _, category := range []string{
		"apps-and-integrations",
		"assignments",
		"billing",
		"comments-and-replies",
		"customers",
		"document-changes",
		"feed",
		"mentions",
		"posts-and-updates",
		"reactions",
		"reminders",
		"reviews",
		"status-changes",
		"subscriptions",
		"system",
		"triage",
	} {
		preference, err := GetUserSettingsNotificationCategoryPreference(context.Background(), graphqlClient, category)
		require.NoError(t, err)
		require.True(t, preference.Desktop)
	}

	channel, err := GetUserSettingsNotificationChannelPreferences(context.Background(), graphqlClient)
	require.NoError(t, err)
	require.True(t, channel.Slack)

	delivery, err := GetUserSettingsNotificationDeliveryPreferences(context.Background(), graphqlClient)
	require.NoError(t, err)
	require.NotNil(t, delivery.Mobile)

	mobile, err := GetUserSettingsMobileDeliveryPreferences(context.Background(), graphqlClient)
	require.NoError(t, err)
	require.NotNil(t, mobile.Schedule)

	schedule, err := GetUserSettingsMobileSchedule(context.Background(), graphqlClient)
	require.NoError(t, err)
	require.Equal(t, "09:00", *schedule.Monday.Start)

	for _, day := range []string{"friday", "monday", "saturday", "sunday", "thursday", "tuesday", "wednesday"} {
		window, err := GetUserSettingsMobileScheduleDay(context.Background(), graphqlClient, day)
		require.NoError(t, err)
		require.Equal(t, "18:00", *window.End)
	}

	theme, err := GetUserSettingsTheme(context.Background(), graphqlClient, "desktop", "dark")
	require.NoError(t, err)
	require.Equal(t, "custom", theme.Preset)

	customTheme, err := GetUserSettingsCustomTheme(context.Background(), graphqlClient, "mobile-web", "light")
	require.NoError(t, err)
	require.Equal(t, 50, customTheme.Contrast)

	sidebar, err := GetUserSettingsCustomSidebarTheme(context.Background(), graphqlClient, "mobileWeb", "light")
	require.NoError(t, err)
	require.Equal(t, 70, sidebar.Contrast)

	_, err = GetUserSettingsNotificationCategoryPreference(context.Background(), graphqlClient, "unknown")
	require.Error(t, err)
	_, err = GetUserSettingsMobileScheduleDay(context.Background(), graphqlClient, "funday")
	require.Error(t, err)
	_, err = GetUserSettingsTheme(context.Background(), graphqlClient, "watch", "light")
	require.Error(t, err)
	_, err = GetUserSettingsTheme(context.Background(), graphqlClient, "desktop", "sepia")
	require.Error(t, err)
	_, err = GetUserSettingsCustomTheme(context.Background(), graphqlClient, "watch", "light")
	require.Error(t, err)
	_, err = GetUserSettingsCustomTheme(context.Background(), graphqlClient, "desktop", "sepia")
	require.Error(t, err)
	_, err = GetUserSettingsCustomSidebarTheme(context.Background(), graphqlClient, "watch", "light")
	require.Error(t, err)
	_, err = GetUserSettingsCustomSidebarTheme(context.Background(), graphqlClient, "desktop", "sepia")
	require.Error(t, err)
}

func Test_ClientReadScenarios_cover_user_settings_error_and_null_branches(t *testing.T) {
	errorClient := errorGraphQLClient{err: errors.New("graphql failed")}

	_, err := GetUserSettings(context.Background(), errorClient)
	require.Error(t, err)
	_, err = GetUserSettingsNotificationCategoryPreferences(context.Background(), errorClient)
	require.Error(t, err)
	_, err = GetUserSettingsNotificationChannelPreferences(context.Background(), errorClient)
	require.Error(t, err)
	_, err = GetUserSettingsNotificationDeliveryPreferences(context.Background(), errorClient)
	require.Error(t, err)
	_, err = GetUserSettingsMobileDeliveryPreferences(context.Background(), errorClient)
	require.Error(t, err)
	_, err = GetUserSettingsMobileSchedule(context.Background(), errorClient)
	require.Error(t, err)
	_, err = GetUserSettingsTheme(context.Background(), errorClient, "desktop", "light")
	require.Error(t, err)
	_, err = GetUserSettingsCustomTheme(context.Background(), errorClient, "desktop", "light")
	require.Error(t, err)
	_, err = GetUserSettingsCustomSidebarTheme(context.Background(), errorClient, "desktop", "light")
	require.Error(t, err)

	for _, category := range []string{
		"apps-and-integrations",
		"assignments",
		"billing",
		"comments-and-replies",
		"customers",
		"document-changes",
		"feed",
		"mentions",
		"posts-and-updates",
		"reactions",
		"reminders",
		"reviews",
		"status-changes",
		"subscriptions",
		"system",
		"triage",
	} {
		_, err = GetUserSettingsNotificationCategoryPreference(context.Background(), errorClient, category)
		require.Error(t, err)
	}

	for _, day := range []string{"friday", "monday", "saturday", "sunday", "thursday", "tuesday", "wednesday"} {
		_, err = GetUserSettingsMobileScheduleDay(context.Background(), errorClient, day)
		require.Error(t, err)
	}

	nullClient := fakeGraphQLClient{
		"userSettings_notificationDeliveryPreferences":                 `{"userSettings":{"notificationDeliveryPreferences":{"mobile":null}}}`,
		"userSettings_notificationDeliveryPreferences_mobile":          `{"userSettings":{"notificationDeliveryPreferences":{"mobile":null}}}`,
		"userSettings_notificationDeliveryPreferences_mobile_schedule": `{"userSettings":{"notificationDeliveryPreferences":{"mobile":null}}}`,
		"userSettings_theme":                `{"userSettings":{"theme":null}}`,
		"userSettings_theme_custom":         `{"userSettings":{"theme":{"custom":null}}}`,
		"userSettings_theme_custom_sidebar": `{"userSettings":{"theme":{"custom":null}}}`,
	}
	delivery, err := GetUserSettingsNotificationDeliveryPreferences(context.Background(), nullClient)
	require.NoError(t, err)
	require.Nil(t, delivery.Mobile)
	mobile, err := GetUserSettingsMobileDeliveryPreferences(context.Background(), nullClient)
	require.NoError(t, err)
	require.Nil(t, mobile)
	schedule, err := GetUserSettingsMobileSchedule(context.Background(), nullClient)
	require.NoError(t, err)
	require.Nil(t, schedule)
	theme, err := GetUserSettingsTheme(context.Background(), nullClient, "desktop", "light")
	require.NoError(t, err)
	require.Nil(t, theme)
	customTheme, err := GetUserSettingsCustomTheme(context.Background(), nullClient, "desktop", "light")
	require.NoError(t, err)
	require.Nil(t, customTheme)
	sidebar, err := GetUserSettingsCustomSidebarTheme(context.Background(), nullClient, "desktop", "light")
	require.NoError(t, err)
	require.Nil(t, sidebar)

	scheduleNullClient := fakeGraphQLClient{
		"userSettings_notificationDeliveryPreferences_mobile_schedule": `{"userSettings":{"notificationDeliveryPreferences":{"mobile":{"schedule":null}}}}`,
	}
	schedule, err = GetUserSettingsMobileSchedule(context.Background(), scheduleNullClient)
	require.NoError(t, err)
	require.Nil(t, schedule)

	scheduleDay, err := GetUserSettingsMobileScheduleDay(context.Background(), nullClient, "monday")
	require.NoError(t, err)
	require.Empty(t, scheduleDay)
	scheduleDay, err = GetUserSettingsMobileScheduleDay(context.Background(), scheduleNullClient, "monday")
	require.NoError(t, err)
	require.Empty(t, scheduleDay)
}
