package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	projectCommand := &cobra.Command{
		Use:   "project",
		Short: "Read and write Linear projects",
	}
	addProjectListCommand(ctx, projectCommand, options)
	addProjectGetCommand(ctx, projectCommand, options)
	addProjectMembersCommand(ctx, projectCommand, options)
	addProjectUpdatesCommand(ctx, projectCommand, options)
	addProjectCreateCommand(ctx, projectCommand, options)
	addProjectUpdateCommand(ctx, projectCommand, options)
	addProjectArchiveCommand(ctx, projectCommand, options)
	addDomainUsageCommand(projectCommand, options, "project")
	root.AddCommand(projectCommand)
}

func addProjectListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "list",
		Short: "List projects for the resolved team",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			target, err := runtime.resolveTarget(ctx)
			if err != nil {
				return err
			}
			projects, err := client.ListProjectsByTeam(ctx, runtime.graphqlClient, target.Team.ID, limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(projects.Projects)); err != nil {
				return err
			}
			projects.Projects, err = sortByJSONField(projects.Projects, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, projects)
			}
			for _, project := range projects.Projects {
				if err := writeProject(command, options, project); err != nil {
					return err
				}
			}

			return nil
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum projects to return")
	root.AddCommand(command)
}

func addProjectGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "get PROJECT_ID",
		Short: "Get one project by id or slug",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			project, err := client.GetProjectByID(ctx, runtime.graphqlClient, args[0])
			if err != nil {
				return err
			}

			return writeProject(command, options, project)
		},
	})
}

func addProjectMembersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "members PROJECT_ID",
		Short: "List project members",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			members, err := client.ListProjectMembers(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(members.Members)); err != nil {
				return err
			}
			members.Members, err = sortByJSONField(members.Members, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, members)
			}
			for _, member := range members.Members {
				if err := render.WriteLine(command.OutOrStdout(), "%s %s", member.ID, member.DisplayName); err != nil {
					return err
				}
			}

			return nil
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum members to return")
	root.AddCommand(command)
}

func addProjectUpdatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "updates PROJECT_ID",
		Short: "List project status updates",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			updates, err := client.ListProjectUpdates(ctx, runtime.graphqlClient, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, len(updates.Updates)); err != nil {
				return err
			}
			updates.Updates, err = sortByJSONField(updates.Updates, options.sortField, options.sortOrder)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, updates)
			}
			for _, update := range updates.Updates {
				if err := render.WriteLine(
					command.OutOrStdout(),
					"%s %s %s %s",
					update.ID,
					update.Health,
					update.DisplayName,
					update.Body,
				); err != nil {
					return err
				}
			}

			return nil
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum project updates to return")
	root.AddCommand(command)
}

func writeProject(command *cobra.Command, options *rootOptions, project client.ProjectSummary) error {
	if wrote, err := writeIDOnly(command, options, project.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, project)
	}

	format, err := normalizedHumanFormat(options)
	if err != nil {
		return err
	}
	if format == "minimal" {
		return render.WriteLine(command.OutOrStdout(), "%s", project.ID)
	}
	if format == "full" {
		return render.WriteLine(
			command.OutOrStdout(),
			"%s %s [%s] url=%s",
			project.ID,
			project.Name,
			project.Status.Name,
			project.URL,
		)
	}

	return render.WriteLine(command.OutOrStdout(), "%s %s [%s]", project.ID, project.Name, project.Status.Name)
}
