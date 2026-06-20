package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addOrganizationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "organization",
		Short: "Read Linear organization metadata",
	}
	addOrganizationLabelsCommand(ctx, command, options)
	addOrganizationProjectLabelsCommand(ctx, command, options)
	addOrganizationTeamsCommand(ctx, command, options)
	addOrganizationUsersCommand(ctx, command, options)
	command.AddCommand(&cobra.Command{
		Use:   "exists URL_KEY",
		Short: "Check whether a Linear organization URL key exists",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			status, err := client.CheckOrganizationExists(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeOrganizationExists(command, options, status)
		},
	})
	templatesCommand := &cobra.Command{
		Use:   "templates",
		Short: "List workspace-level Linear templates",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadOrganizationTemplateList,
				templatePageWithItems,
				writeTemplate,
			)
		},
	}
	templatesCommand.Flags().IntVar(&limit, "limit", limit, "maximum organization templates to return")
	command.AddCommand(templatesCommand)
	root.AddCommand(command)
}

func addOrganizationLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "labels",
		Short: "List workspace-level issue labels",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadOrganizationLabels,
				labelPageWithItems,
				writeLabel,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum labels to return")
	root.AddCommand(command)
}

func addOrganizationProjectLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "project-labels",
		Short: "List workspace-level project labels",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadOrganizationProjectLabels,
				projectLabelPageWithItems,
				writeProjectLabel,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum project labels to return")
	root.AddCommand(command)
}

func addOrganizationTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "teams",
		Short: "List teams visible in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadOrganizationTeams,
				teamPageWithItems,
				writeTeam,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum teams to return")
	root.AddCommand(command)
}

func addOrganizationUsersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "users",
		Short: "List active users visible in the workspace",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadOrganizationUsers,
				userPageWithItems,
				writeUser,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum users to return")
	root.AddCommand(command)
}

func writeOrganizationExists(
	command *cobra.Command,
	options *rootOptions,
	status client.OrganizationExistsStatus,
) error {
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, status)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s exists %t success %t",
		status.URLKey,
		status.Exists,
		status.Success,
	)
}

func loadOrganizationLabels(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.LabelList, []client.LabelSummary, error) {
	labels, err := client.ListOrganizationLabels(ctx, runtime.graphqlClient, limit)
	return labels, labels.Labels, err
}

func loadOrganizationProjectLabels(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectLabelList, []client.ProjectLabelSummary, error) {
	labels, err := client.ListOrganizationProjectLabels(ctx, runtime.graphqlClient, limit)
	return labels, labels.ProjectLabels, err
}

func loadOrganizationTeams(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamList, []client.TeamSummary, error) {
	teams, err := client.ListOrganizationTeams(ctx, runtime.graphqlClient, limit)
	return teams, teams.Teams, err
}

func loadOrganizationUsers(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.UserList, []client.UserSummary, error) {
	users, err := client.ListOrganizationUsers(ctx, runtime.graphqlClient, limit)
	return users, users.Users, err
}

func loadOrganizationTemplateList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TemplateList, []client.TemplateSummary, error) {
	templates, err := client.ListOrganizationTemplates(ctx, runtime.graphqlClient, limit)
	return templates, templates.Templates, err
}
