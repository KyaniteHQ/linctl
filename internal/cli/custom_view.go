package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addCustomViewCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	customViewCommand := addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.CustomViewList, client.CustomViewSummary]{
			Use:       "custom-view",
			Short:     "Read Linear custom views",
			ListShort: "List visible custom views",
			LimitHelp: "maximum custom views to return",
			GetUse:    "get CUSTOM_VIEW_ID",
			GetShort:  "Get one custom view by id or slug",
			LoadList:  clientList(client.ListCustomViews),
			LoadGet:   clientGet(client.GetCustomViewByID),
			WriteItem: writeCustomView,
		},
	)
	addCustomViewSubscribersCommand(ctx, customViewCommand, options)
	addCustomViewInitiativesCommand(ctx, customViewCommand, options)
	addCustomViewIssuesCommand(ctx, customViewCommand, options)
	addCustomViewOrganizationPreferencesCommand(ctx, customViewCommand, options)
	addCustomViewProjectsCommand(ctx, customViewCommand, options)
	addCustomViewUserPreferencesCommand(ctx, customViewCommand, options)
	addCustomViewPreferenceValuesCommand(ctx, customViewCommand, options)
}

func addCustomViewSubscribersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.CustomViewSubscriberStatus]{
		Use:   "subscribers CUSTOM_VIEW_ID",
		Short: "Show whether a custom view has subscribers",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.CustomViewSubscriberStatus, error) {
			return client.GetCustomViewSubscriberStatus(ctx, runtime.graphqlClient, id)
		},
		Write: writeCustomViewSubscriberStatus,
	})
}

func addCustomViewInitiativesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"initiatives CUSTOM_VIEW_ID",
		"List initiatives matching a custom view",
		"initiatives",
		client.ListCustomViewInitiatives,
		writeInitiative,
	)
}

func addCustomViewIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"issues CUSTOM_VIEW_ID",
		"List issues matching a custom view",
		"issues",
		client.ListCustomViewIssues,
		writeIssue,
	)
}

func addCustomViewOrganizationPreferencesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadGetCommand(ctx, root, options, readGetSpec[client.CustomViewPreferences]{
		Use:   "organization-preferences CUSTOM_VIEW_ID",
		Short: "Read organization default view preferences for a custom view",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.CustomViewPreferences, error) {
			return client.GetCustomViewOrganizationPreferences(ctx, runtime.graphqlClient, id)
		},
		Write: writeCustomViewPreferences,
	})

	addReadGetCommand(ctx, command, options, readGetSpec[client.CustomViewPreferencesValues]{
		Use:   "values CUSTOM_VIEW_ID",
		Short: "Read organization default view preference values for a custom view",
		Load: func(
			ctx context.Context, runtime commandRuntime, id string,
		) (client.CustomViewPreferencesValues, error) {
			return client.GetCustomViewOrganizationPreferenceValues(ctx, runtime.graphqlClient, id)
		},
		Write: writeCustomViewPreferenceValues,
	})
}

func addCustomViewProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"projects CUSTOM_VIEW_ID",
		"List projects matching a custom view",
		"projects",
		client.ListCustomViewProjects,
		writeProject,
	)
}

func addCustomViewUserPreferencesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := addReadGetCommand(ctx, root, options, readGetSpec[client.CustomViewPreferences]{
		Use:   "user-preferences CUSTOM_VIEW_ID",
		Short: "Read current-user view preferences for a custom view",
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.CustomViewPreferences, error) {
			return client.GetCustomViewUserPreferences(ctx, runtime.graphqlClient, id)
		},
		Write: func(command *cobra.Command, options *rootOptions, preferences client.CustomViewPreferences) error {
			return writeCustomViewScopedPreferences(command, options, "user", preferences)
		},
	})

	addReadGetCommand(ctx, command, options, readGetSpec[client.CustomViewPreferencesValues]{
		Use:   "values CUSTOM_VIEW_ID",
		Short: "Read current-user view preference values for a custom view",
		Load: func(
			ctx context.Context, runtime commandRuntime, id string,
		) (client.CustomViewPreferencesValues, error) {
			return client.GetCustomViewUserPreferenceValues(ctx, runtime.graphqlClient, id)
		},
		Write: writeCustomViewPreferenceValues,
	})
}

func addCustomViewPreferenceValuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadGetCommand(ctx, root, options, readGetSpec[client.CustomViewPreferencesValues]{
		Use:   "preference-values CUSTOM_VIEW_ID",
		Short: "Read effective view preference values for a custom view",
		Load: func(
			ctx context.Context, runtime commandRuntime, id string,
		) (client.CustomViewPreferencesValues, error) {
			return client.GetCustomViewPreferenceValues(ctx, runtime.graphqlClient, id)
		},
		Write: writeCustomViewPreferenceValues,
	})
}

func writeCustomView(command *cobra.Command, options *rootOptions, view client.CustomViewSummary) error {
	return writeItemLine(command, options, view, view.ID, "%s %s [%s]", view.ID, view.Name, view.ModelName)
}

func writeCustomViewSubscriberStatus(
	command *cobra.Command,
	options *rootOptions,
	status client.CustomViewSubscriberStatus,
) error {
	return writeItemLine(command, options, status, status.ID, "%s has_subscribers %t", status.ID, status.HasSubscribers)
}

func writeCustomViewPreferences(
	command *cobra.Command,
	options *rootOptions,
	preferences client.CustomViewPreferences,
) error {
	return writeCustomViewScopedPreferences(command, options, "organization", preferences)
}

func writeCustomViewScopedPreferences(
	command *cobra.Command,
	options *rootOptions,
	scope string,
	preferences client.CustomViewPreferences,
) error {
	return writeItem(command, options, preferences, preferences.ID,
		func(command *cobra.Command, _ *rootOptions, preferences client.CustomViewPreferences) error {
			if preferences.ID == "" {
				return render.WriteLine(command.OutOrStdout(), "%s %s preferences -", preferences.CustomViewID, scope)
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s preferences %s %s layout %s",
				preferences.CustomViewID,
				scope,
				preferences.Type,
				preferences.ViewType,
				emptyDash(preferences.Values.Layout),
			)
		})
}

func writeCustomViewPreferenceValues(
	command *cobra.Command,
	options *rootOptions,
	values client.CustomViewPreferencesValues,
) error {
	return writeItemNoID(command, options, values,
		func(command *cobra.Command, _ *rootOptions, values client.CustomViewPreferencesValues) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s preference values layout %s ordering %s",
				values.CustomViewID,
				emptyDash(values.Layout),
				emptyDash(values.ViewOrdering),
			)
		})
}
