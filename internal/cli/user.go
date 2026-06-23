package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addUserCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	userCommand := &cobra.Command{
		Use:   "user",
		Short: "Read Linear users",
	}
	addUserListCommand(ctx, userCommand, options)
	addUserGetCommand(ctx, userCommand, options)
	addUserMeCommand(ctx, userCommand, options)
	addUserDraftsCommand(ctx, userCommand, options)
	addUserSettingsCommand(ctx, userCommand, options)
	addUserAssignedIssuesCommand(ctx, userCommand, options)
	addUserCreatedIssuesCommand(ctx, userCommand, options)
	addUserDelegatedIssuesCommand(ctx, userCommand, options)
	addUserTeamMembershipsCommand(ctx, userCommand, options)
	addUserTeamsCommand(ctx, userCommand, options)
	addViewerAssignedIssuesCommand(ctx, userCommand, options)
	addViewerCreatedIssuesCommand(ctx, userCommand, options)
	addViewerDelegatedIssuesCommand(ctx, userCommand, options)
	addViewerTeamMembershipsCommand(ctx, userCommand, options)
	addViewerTeamsCommand(ctx, userCommand, options)
	root.AddCommand(userCommand)
}

func addUserListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "list",
		Short: "List visible users",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(ctx, command, nil, options, limit, loadUserList, userPageWithItems, writeUser)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum users to return")
	root.AddCommand(command)
}

func addUserGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get USER_ID",
		Short: "Get one User by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			user, err := client.GetUserByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeUser(command, options, user)
		},
	})
}

func addUserMeCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "me",
		Short: "Get the authenticated User",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			user, err := client.GetViewerUser(ctx, runtime.graphqlClient)
			if err != nil {
				return err
			}

			return writeUser(command, options, user)
		},
	})
}

func addUserDraftsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "drafts",
		Short: "List the authenticated user's saved draft metadata",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadViewerDraftList,
				draftPageWithItems,
				writeDraft,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum drafts to return")
	root.AddCommand(command)
}

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

func addUserAssignedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "assigned-issues USER_ID",
		Short: "List issues assigned to a User",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadUserAssignedIssues, issuePageWithItems, writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addUserCreatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "created-issues USER_ID",
		Short: "List issues created by a User",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadUserCreatedIssues, issuePageWithItems, writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addUserDelegatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "delegated-issues USER_ID",
		Short: "List issues delegated to a User",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadUserDelegatedIssues, issuePageWithItems, writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addUserTeamMembershipsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "team-memberships USER_ID",
		Short: "List a User's team memberships",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx, command, args, options, limit,
				loadUserTeamMemberships, teamMembershipPageWithItems, writeTeamMembership,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum memberships to return")
	root.AddCommand(command)
}

func addUserTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "teams USER_ID",
		Short: "List Teams for a User",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(ctx, command, args, options, limit, loadUserTeams, teamPageWithItems, writeTeam)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum teams to return")
	root.AddCommand(command)
}

func addViewerAssignedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "my-assigned-issues",
		Short: "List issues assigned to the authenticated User",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx, command, nil, options, limit,
				loadViewerAssignedIssues, issuePageWithItems, writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addViewerCreatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "my-created-issues",
		Short: "List issues created by the authenticated User",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx, command, nil, options, limit,
				loadViewerCreatedIssues, issuePageWithItems, writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addViewerDelegatedIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "my-delegated-issues",
		Short: "List issues delegated to the authenticated User",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx, command, nil, options, limit,
				loadViewerDelegatedIssues, issuePageWithItems, writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addViewerTeamMembershipsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "my-team-memberships",
		Short: "List team memberships for the authenticated User",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx, command, nil, options, limit,
				loadViewerTeamMemberships, teamMembershipPageWithItems, writeTeamMembership,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum memberships to return")
	root.AddCommand(command)
}

func addViewerTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "my-teams",
		Short: "List Teams for the authenticated User",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(ctx, command, nil, options, limit, loadViewerTeams, teamPageWithItems, writeTeam)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum teams to return")
	root.AddCommand(command)
}

func writeUser(command *cobra.Command, options *rootOptions, user client.UserSummary) error {
	if wrote, err := writeIDOnly(command, options, user.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, user)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s <%s>", user.ID, user.DisplayName, user.Email)
}

func writeDraft(command *cobra.Command, options *rootOptions, draft client.DraftSummary) error {
	if wrote, err := writeIDOnly(command, options, draft.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, draft)
	}

	parentKey := defaultString(draft.ParentKey, "-")
	parentTitle := defaultString(draft.ParentTitle, "-")
	return render.WriteLine(
		command.OutOrStdout(),
		"%s %s %s %s",
		draft.ID,
		draft.ParentType,
		parentKey,
		parentTitle,
	)
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

func loadUserList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.UserList, []client.UserSummary, error) {
	users, err := client.ListUsers(ctx, runtime.graphqlClient, limit)
	return users, users.Users, err
}

func loadViewerDraftList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.DraftList, []client.DraftSummary, error) {
	drafts, err := client.ListViewerDrafts(ctx, runtime.graphqlClient, limit)
	return drafts, drafts.Drafts, err
}

func userPageWithItems(page client.UserList, users []client.UserSummary) client.UserList {
	page.Users = users
	return page
}

func draftPageWithItems(page client.DraftList, drafts []client.DraftSummary) client.DraftList {
	page.Drafts = drafts
	return page
}

func loadUserAssignedIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListUserAssignedIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, issues.Issues, err
}

func loadUserCreatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListUserCreatedIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, issues.Issues, err
}

func loadUserDelegatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListUserDelegatedIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, issues.Issues, err
}

func loadUserTeamMemberships(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamMembershipList, []client.TeamMembershipSummary, error) {
	memberships, err := client.ListUserTeamMemberships(ctx, runtime.graphqlClient, args[0], limit)
	return memberships, memberships.Memberships, err
}

func loadUserTeams(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.TeamList, []client.TeamSummary, error) {
	teams, err := client.ListUserTeams(ctx, runtime.graphqlClient, args[0], limit)
	return teams, teams.Teams, err
}

func loadViewerAssignedIssues(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListViewerAssignedIssues(ctx, runtime.graphqlClient, limit)
	return issues, issues.Issues, err
}

func loadViewerCreatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListViewerCreatedIssues(ctx, runtime.graphqlClient, limit)
	return issues, issues.Issues, err
}

func loadViewerDelegatedIssues(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListViewerDelegatedIssues(ctx, runtime.graphqlClient, limit)
	return issues, issues.Issues, err
}

func loadViewerTeamMemberships(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamMembershipList, []client.TeamMembershipSummary, error) {
	memberships, err := client.ListViewerTeamMemberships(ctx, runtime.graphqlClient, limit)
	return memberships, memberships.Memberships, err
}

func loadViewerTeams(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamList, []client.TeamSummary, error) {
	teams, err := client.ListViewerTeams(ctx, runtime.graphqlClient, limit)
	return teams, teams.Teams, err
}
