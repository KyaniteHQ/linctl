package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addUserSettingsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	settingsCommand := &cobra.Command{
		Use:   "settings",
		Short: "Read authenticated User settings",
	}
	addUserSettingsGetCommand(ctx, settingsCommand, options)
	addUserSettingsNotificationCategoriesCommand(ctx, settingsCommand, options)
	addUserSettingsNotificationCategoryCommand(ctx, settingsCommand, options)
	addUserSettingsNotificationChannelsCommand(ctx, settingsCommand, options)
	addUserSettingsNotificationDeliveryCommand(ctx, settingsCommand, options)
	addUserSettingsMobileDeliveryCommand(ctx, settingsCommand, options)
	addUserSettingsMobileScheduleCommand(ctx, settingsCommand, options)
	addUserSettingsMobileScheduleDayCommand(ctx, settingsCommand, options)
	addUserSettingsThemeCommand(ctx, settingsCommand, options, "theme")
	addUserSettingsThemeCommand(ctx, settingsCommand, options, "custom-theme")
	addUserSettingsThemeCommand(ctx, settingsCommand, options, "custom-sidebar-theme")
	root.AddCommand(settingsCommand)
}

func addUserSettingsGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get authenticated User settings",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			settings, err := client.GetUserSettings(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeUserSettings(command, options, settings)
		},
	})
}

func addUserSettingsNotificationCategoriesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "notification-categories",
		Short: "Get User notification category preferences",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			preferences, err := client.GetUserSettingsNotificationCategoryPreferences(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeUserSettingsValue(command, options, preferences, "notification categories")
		},
	})
}

func addUserSettingsNotificationCategoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "notification-category CATEGORY",
		Short: "Get one User notification category preference",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			preference, err := client.GetUserSettingsNotificationCategoryPreference(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeUserSettingsValue(command, options, preference, notificationChannelsText(args[0], preference))
		},
	})
}

func addUserSettingsNotificationChannelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "notification-channels",
		Short: "Get User notification channel preferences",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			preference, err := client.GetUserSettingsNotificationChannelPreferences(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeUserSettingsValue(
				command,
				options,
				preference,
				notificationChannelsText("channels", preference),
			)
		},
	})
}

func addUserSettingsNotificationDeliveryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "notification-delivery",
		Short: "Get User notification delivery preferences",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			preferences, err := client.GetUserSettingsNotificationDeliveryPreferences(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeUserSettingsValue(command, options, preferences, "notification delivery")
		},
	})
}

func addUserSettingsMobileDeliveryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "mobile-delivery",
		Short: "Get User mobile notification delivery preferences",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			channel, err := client.GetUserSettingsMobileDeliveryPreferences(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeUserSettingsNullableValue(command, options, channel, "mobile delivery")
		},
	})
}

func addUserSettingsMobileScheduleCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "mobile-schedule",
		Short: "Get User mobile notification schedule",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			schedule, err := client.GetUserSettingsMobileSchedule(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeUserSettingsNullableValue(command, options, schedule, "mobile schedule")
		},
	})
}

func addUserSettingsMobileScheduleDayCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "mobile-schedule-day DAY",
		Short: "Get one User mobile notification schedule day",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			day, err := client.GetUserSettingsMobileScheduleDay(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeUserSettingsValue(command, options, day, notificationDayText(args[0], day))
		},
	})
}

func addUserSettingsThemeCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	name string,
) {
	deviceType := "desktop"
	mode := "light"
	command := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("Get User %s settings", name),
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			return runUserSettingsThemeCommand(ctx, command, options, runtime, name, deviceType, mode)
		},
	}
	command.Flags().StringVar(&deviceType, "device-type", deviceType, "theme device type: desktop or mobile-web")
	command.Flags().StringVar(&mode, "mode", mode, "theme mode: light or dark")
	root.AddCommand(command)
}

func runUserSettingsThemeCommand(
	ctx context.Context,
	command *cobra.Command,
	options *rootOptions,
	runtime commandRuntime,
	name string,
	deviceType string,
	mode string,
) error {
	switch name {
	case "theme":
		theme, err := client.GetUserSettingsTheme(ctx, runtime.graphqlClient, deviceType, mode)
		if err != nil {
			return err
		}
		return writeUserSettingsNullableValue(command, options, theme, userSettingsThemeText(name, deviceType, mode))
	case "custom-theme":
		theme, err := client.GetUserSettingsCustomTheme(ctx, runtime.graphqlClient, deviceType, mode)
		if err != nil {
			return err
		}
		return writeUserSettingsNullableValue(command, options, theme, userSettingsThemeText(name, deviceType, mode))
	case "custom-sidebar-theme":
		theme, err := client.GetUserSettingsCustomSidebarTheme(ctx, runtime.graphqlClient, deviceType, mode)
		if err != nil {
			return err
		}
		return writeUserSettingsNullableValue(command, options, theme, userSettingsThemeText(name, deviceType, mode))
	default:
		return fmt.Errorf("unknown user settings theme command %q", name)
	}
}

func writeUserSettings(command *cobra.Command, options *rootOptions, settings client.UserSettingsSummary) error {
	if wrote, err := writeIDOnly(command, options, settings.ID); wrote || err != nil {
		return err
	}

	return writeUserSettingsValue(
		command,
		options,
		settings,
		fmt.Sprintf(
			"%s user=%s auto_assign=%t full_names=%t",
			settings.ID,
			settings.UserID,
			settings.AutoAssignToSelf,
			settings.ShowFullUserNames,
		),
	)
}

func writeUserSettingsValue(command *cobra.Command, options *rootOptions, value any, human string) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, value)
	}

	return render.WriteLine(command.OutOrStdout(), "%s", human)
}

func writeUserSettingsNullableValue(command *cobra.Command, options *rootOptions, value any, human string) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return render.WriteJSON(command.OutOrStdout(), value, options.compact)
	}
	if value == nil {
		return render.WriteLine(command.OutOrStdout(), "%s none", human)
	}

	return render.WriteLine(command.OutOrStdout(), "%s", human)
}

func notificationChannelsText(category string, preference client.NotificationChannelPreference) string {
	return fmt.Sprintf(
		"%s desktop=%t email=%t mobile=%t slack=%t",
		category,
		preference.Desktop,
		preference.Email,
		preference.Mobile,
		preference.Slack,
	)
}

func notificationDayText(day string, preference client.NotificationDeliveryDay) string {
	return fmt.Sprintf(
		"%s start=%s end=%s",
		day,
		defaultString(pointerString(preference.Start), "-"),
		defaultString(pointerString(preference.End), "-"),
	)
}

func userSettingsThemeText(name string, deviceType string, mode string) string {
	return fmt.Sprintf("%s device_type=%s mode=%s", name, deviceType, mode)
}

func pointerString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
