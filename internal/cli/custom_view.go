//nolint:dupl // Minimal read-command glue is intentionally uniform across domains via addReadListGetCommand.
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
			Use:           "custom-view",
			Short:         "Read Linear custom views",
			ListShort:     "List visible custom views",
			LimitHelp:     "maximum custom views to return",
			GetUse:        "get CUSTOM_VIEW_ID",
			GetShort:      "Get one custom view by id or slug",
			LoadList:      loadCustomViewList,
			PageWithItems: customViewPageWithItems,
			LoadGet:       loadCustomView,
			WriteItem:     writeCustomView,
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
	command := &cobra.Command{
		Use:   "subscribers CUSTOM_VIEW_ID",
		Short: "Report whether a custom view has subscribers",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			status, err := client.GetCustomViewSubscriberStatus(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCustomViewSubscriberStatus(command, options, status)
		},
	}
	root.AddCommand(command)
}

func addCustomViewInitiativesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "initiatives CUSTOM_VIEW_ID",
		Short: "List initiatives matching a custom view",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadCustomViewInitiatives,
				initiativePageWithItems,
				writeInitiative,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum initiatives to return")
	root.AddCommand(command)
}

func addCustomViewIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "issues CUSTOM_VIEW_ID",
		Short: "List issues matching a custom view",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadCustomViewIssues,
				issuePageWithItems,
				writeIssue,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum issues to return")
	root.AddCommand(command)
}

func addCustomViewOrganizationPreferencesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "organization-preferences CUSTOM_VIEW_ID",
		Short: "Read organization default view preferences for a custom view",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			preferences, err := client.GetCustomViewOrganizationPreferences(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCustomViewPreferences(command, options, preferences)
		},
	}

	valuesCommand := &cobra.Command{
		Use:   "values CUSTOM_VIEW_ID",
		Short: "Read organization default view preference values for a custom view",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			values, err := client.GetCustomViewOrganizationPreferenceValues(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCustomViewPreferenceValues(command, options, values)
		},
	}

	command.AddCommand(valuesCommand)
	root.AddCommand(command)
}

func addCustomViewProjectsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "projects CUSTOM_VIEW_ID",
		Short: "List projects matching a custom view",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return runReadListCommand(
				ctx,
				command,
				args,
				options,
				limit,
				loadCustomViewProjects,
				projectPageWithItems,
				writeProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum projects to return")
	root.AddCommand(command)
}

func addCustomViewUserPreferencesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "user-preferences CUSTOM_VIEW_ID",
		Short: "Read current-user view preferences for a custom view",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			preferences, err := client.GetCustomViewUserPreferences(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCustomViewScopedPreferences(command, options, "user", preferences)
		},
	}

	valuesCommand := &cobra.Command{
		Use:   "values CUSTOM_VIEW_ID",
		Short: "Read current-user view preference values for a custom view",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			values, err := client.GetCustomViewUserPreferenceValues(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCustomViewPreferenceValues(command, options, values)
		},
	}

	command.AddCommand(valuesCommand)
	root.AddCommand(command)
}

func addCustomViewPreferenceValuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := &cobra.Command{
		Use:   "preference-values CUSTOM_VIEW_ID",
		Short: "Read effective view preference values for a custom view",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			values, err := client.GetCustomViewPreferenceValues(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeCustomViewPreferenceValues(command, options, values)
		},
	}
	root.AddCommand(command)
}

func writeCustomView(command *cobra.Command, options *rootOptions, view client.CustomViewSummary) error {
	return writeItem(command, options, view, view.ID,
		func(command *cobra.Command, _ *rootOptions, view client.CustomViewSummary) error {
			return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", view.ID, view.Name, view.ModelName)
		})
}

func writeCustomViewSubscriberStatus(
	command *cobra.Command,
	options *rootOptions,
	status client.CustomViewSubscriberStatus,
) error {
	return writeItem(command, options, status, status.ID,
		func(command *cobra.Command, _ *rootOptions, status client.CustomViewSubscriberStatus) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s has_subscribers %t",
				status.ID,
				status.HasSubscribers,
			)
		})
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

func loadCustomViewList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.CustomViewList, []client.CustomViewSummary, error) {
	views, err := client.ListCustomViews(ctx, runtime.graphqlClient, limit)
	return views, views.CustomViews, err
}

func loadCustomView(
	ctx context.Context,
	runtime commandRuntime,
	id string,
) (client.CustomViewSummary, error) {
	return client.GetCustomViewByID(ctx, runtime.graphqlClient, id)
}

func loadCustomViewInitiatives(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.InitiativeList, []client.InitiativeSummary, error) {
	initiatives, err := client.ListCustomViewInitiatives(ctx, runtime.graphqlClient, args[0], limit)
	return initiatives, initiatives.Initiatives, err
}

func loadCustomViewIssues(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.IssueList, []client.IssueSummary, error) {
	issues, err := client.ListCustomViewIssues(ctx, runtime.graphqlClient, args[0], limit)
	return issues, issues.Issues, err
}

func loadCustomViewProjects(
	ctx context.Context,
	runtime commandRuntime,
	args []string,
	limit int,
) (client.ProjectList, []client.ProjectSummary, error) {
	projects, err := client.ListCustomViewProjects(ctx, runtime.graphqlClient, args[0], limit)
	return projects, projects.Projects, err
}

func customViewPageWithItems(
	page client.CustomViewList,
	views []client.CustomViewSummary,
) client.CustomViewList {
	page.CustomViews = views
	return page
}
