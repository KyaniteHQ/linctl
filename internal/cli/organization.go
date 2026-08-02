package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addOrganizationCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	command := newGroupCommand("organization", "Read Linear organization metadata")
	addOrganizationLabelsCommand(ctx, command, options)
	addOrganizationProjectLabelsCommand(ctx, command, options)
	addOrganizationTeamsCommand(ctx, command, options)
	addOrganizationUsersCommand(ctx, command, options)
	addReadGetCommand(ctx, command, options, readGetSpec[client.OrganizationExistsStatus]{
		Use:   "exists URL_KEY",
		Short: "Check whether a Linear organization URL key exists already",
		Load: func(
			ctx context.Context, runtime commandRuntime, id string,
		) (client.OrganizationExistsStatus, error) {
			return client.CheckOrganizationExists(ctx, runtime.graphqlClient, id)
		},
		Write: writeOrganizationExists,
	})
	addListCommand(ctx, command, options, listCommandSpec[client.TemplateList, client.TemplateSummary]{
		Use:       "templates",
		Short:     "List organization-level Linear templates",
		LimitHelp: "organization templates",
		Args:      cobra.NoArgs,
		Load:      loadOrganizationTemplateList,
		WriteItem: writeTemplate,
	})
	root.AddCommand(command)
}

func addOrganizationLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.LabelList, client.LabelSummary]{
		Use:       "labels",
		Short:     "List organization-level issue labels",
		LimitHelp: "labels",
		Args:      cobra.NoArgs,
		Load:      loadOrganizationLabels,
		WriteItem: writeLabel,
	})
}

func addOrganizationProjectLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectLabelList, client.ProjectLabelSummary]{
		Use:       "project-labels",
		Short:     "List organization-level project labels",
		LimitHelp: "project labels",
		Args:      cobra.NoArgs,
		Load:      loadOrganizationProjectLabels,
		WriteItem: writeProjectLabel,
	})
}

func addOrganizationTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.TeamList, client.TeamSummary]{
		Use:       "teams",
		Short:     "List teams visible to the authenticated user",
		LimitHelp: "teams",
		Args:      cobra.NoArgs,
		Load:      loadOrganizationTeams,
		WriteItem: writeTeam,
	})
}

func addOrganizationUsersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.UserList, client.UserSummary]{
		Use:       "users",
		Short:     "List active users visible to the authenticated user",
		LimitHelp: "users",
		Args:      cobra.NoArgs,
		Load:      loadOrganizationUsers,
		WriteItem: writeUser,
	})
}

func writeOrganizationExists(
	command *cobra.Command,
	options *rootOptions,
	status client.OrganizationExistsStatus,
) error {
	return writeItemLine(
		command, options, status, status.URLKey,
		"%s exists %t success %t", status.URLKey, status.Exists, status.Success,
	)
}

func loadOrganizationLabels(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.LabelList, error) {
	labels, err := client.ListOrganizationLabels(ctx, runtime.graphqlClient, limit)
	return labels, err
}

func loadOrganizationProjectLabels(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectLabelList, error) {
	labels, err := client.ListOrganizationProjectLabels(ctx, runtime.graphqlClient, limit)
	return labels, err
}

func loadOrganizationTeams(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TeamList, error) {
	teams, err := client.ListOrganizationTeams(ctx, runtime.graphqlClient, limit)
	return teams, err
}

func loadOrganizationUsers(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.UserList, error) {
	users, err := client.ListOrganizationUsers(ctx, runtime.graphqlClient, limit)
	return users, err
}

func loadOrganizationTemplateList(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.TemplateList, error) {
	templates, err := client.ListOrganizationTemplates(ctx, runtime.graphqlClient, limit)
	return templates, err
}
