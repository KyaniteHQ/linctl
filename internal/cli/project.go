//nolint:dupl // Project child read commands intentionally share the same list-command shape.
package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"
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
	addProjectAttachmentsCommand(ctx, projectCommand, options)
	addProjectDocumentsCommand(ctx, projectCommand, options)
	addProjectExternalLinksCommand(ctx, projectCommand, options)
	addProjectHistoryCommand(ctx, projectCommand, options)
	addProjectInitiativeLinksCommand(ctx, projectCommand, options)
	addProjectInitiativesCommand(ctx, projectCommand, options)
	addProjectInverseRelationsCommand(ctx, projectCommand, options)
	addProjectIssuesCommand(ctx, projectCommand, options)
	addProjectLabelsCommand(ctx, projectCommand, options)
	addProjectMembersCommand(ctx, projectCommand, options)
	addProjectNeedsCommand(ctx, projectCommand, options)
	addProjectRelationsCommand(ctx, projectCommand, options)
	addProjectTeamsCommand(ctx, projectCommand, options)
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

func addProjectAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"attachments PROJECT_ID",
		"List project attachments",
		"attachments",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectAttachmentList, error) {
			return client.ListProjectAttachments(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectAttachmentList) int {
			return len(list.Attachments)
		},
		func(list client.ProjectAttachmentList) (client.ProjectAttachmentList, error) {
			items, err := sortByJSONField(list.Attachments, options.sortField, options.sortOrder)
			list.Attachments = items
			return list, err
		},
		func(command *cobra.Command, item client.AttachmentSummary) error {
			return writeAttachment(command, options, item)
		},
		func(list client.ProjectAttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addProjectDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"documents PROJECT_ID",
		"List project documents",
		"documents",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectDocumentList, error) {
			return client.ListProjectDocuments(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectDocumentList) int {
			return len(list.Documents)
		},
		func(list client.ProjectDocumentList) (client.ProjectDocumentList, error) {
			items, err := sortByJSONField(list.Documents, options.sortField, options.sortOrder)
			list.Documents = items
			return list, err
		},
		func(command *cobra.Command, item client.DocumentSummary) error {
			return writeDocument(command, options, item)
		},
		func(list client.ProjectDocumentList) []client.DocumentSummary {
			return list.Documents
		},
	)
}

func addProjectExternalLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"external-links PROJECT_ID",
		"List project external links",
		"external links",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectExternalLinkList, error) {
			return client.ListProjectExternalLinks(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectExternalLinkList) int {
			return len(list.Links)
		},
		func(list client.ProjectExternalLinkList) (client.ProjectExternalLinkList, error) {
			items, err := sortByJSONField(list.Links, options.sortField, options.sortOrder)
			list.Links = items
			return list, err
		},
		func(command *cobra.Command, item client.EntityExternalLinkSummary) error {
			return writeEntityExternalLink(command, options, item)
		},
		func(list client.ProjectExternalLinkList) []client.EntityExternalLinkSummary {
			return list.Links
		},
	)
}

func addProjectHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"history PROJECT_ID",
		"List project history",
		"history entries",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectHistoryList, error) {
			return client.ListProjectHistory(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectHistoryList) int {
			return len(list.History)
		},
		func(list client.ProjectHistoryList) (client.ProjectHistoryList, error) {
			items, err := sortByJSONField(list.History, options.sortField, options.sortOrder)
			list.History = items
			return list, err
		},
		func(command *cobra.Command, item client.ProjectHistorySummary) error {
			return writeProjectHistory(command, options, item)
		},
		func(list client.ProjectHistoryList) []client.ProjectHistorySummary {
			return list.History
		},
	)
}

func addProjectInitiativeLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"initiative-links PROJECT_ID",
		"List project initiative associations",
		"initiative associations",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectInitiativeToProjectList, error) {
			return client.ListProjectInitiativeToProjects(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectInitiativeToProjectList) int {
			return len(list.Associations)
		},
		func(list client.ProjectInitiativeToProjectList) (client.ProjectInitiativeToProjectList, error) {
			items, err := sortByJSONField(list.Associations, options.sortField, options.sortOrder)
			list.Associations = items
			return list, err
		},
		func(command *cobra.Command, item client.InitiativeToProjectSummary) error {
			return writeInitiativeToProject(command, options, item)
		},
		func(list client.ProjectInitiativeToProjectList) []client.InitiativeToProjectSummary {
			return list.Associations
		},
	)
}

func addProjectInitiativesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"initiatives PROJECT_ID",
		"List project initiatives",
		"initiatives",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectInitiativeList, error) {
			return client.ListProjectInitiatives(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectInitiativeList) int {
			return len(list.Initiatives)
		},
		func(list client.ProjectInitiativeList) (client.ProjectInitiativeList, error) {
			items, err := sortByJSONField(list.Initiatives, options.sortField, options.sortOrder)
			list.Initiatives = items
			return list, err
		},
		func(command *cobra.Command, item client.InitiativeSummary) error {
			return writeInitiative(command, options, item)
		},
		func(list client.ProjectInitiativeList) []client.InitiativeSummary {
			return list.Initiatives
		},
	)
}

func addProjectInverseRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectRelationChildListCommand(
		ctx,
		root,
		options,
		"inverse-relations PROJECT_ID",
		"List project inverse relations",
		"inverse relations",
		client.ListProjectInverseRelations,
	)
}

func addProjectIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"issues PROJECT_ID",
		"List project issues",
		"issues",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectIssueList, error) {
			return client.ListProjectIssues(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectIssueList) int {
			return len(list.Issues)
		},
		func(list client.ProjectIssueList) (client.ProjectIssueList, error) {
			items, err := sortByJSONField(list.Issues, options.sortField, options.sortOrder)
			list.Issues = items
			return list, err
		},
		func(command *cobra.Command, item client.IssueSummary) error {
			return writeIssue(command, options, item)
		},
		func(list client.ProjectIssueList) []client.IssueSummary {
			return list.Issues
		},
	)
}

func addProjectLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"labels PROJECT_ID",
		"List project labels",
		"labels",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectProjectLabelList, error) {
			return client.ListLabelsForProject(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectProjectLabelList) int {
			return len(list.ProjectLabels)
		},
		func(list client.ProjectProjectLabelList) (client.ProjectProjectLabelList, error) {
			items, err := sortByJSONField(list.ProjectLabels, options.sortField, options.sortOrder)
			list.ProjectLabels = items
			return list, err
		},
		func(command *cobra.Command, item client.ProjectLabelSummary) error {
			return writeProjectLabel(command, options, item)
		},
		func(list client.ProjectProjectLabelList) []client.ProjectLabelSummary {
			return list.ProjectLabels
		},
	)
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

func addProjectNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"needs PROJECT_ID",
		"List project customer needs",
		"customer needs",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectCustomerNeedList, error) {
			return client.ListProjectNeeds(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectCustomerNeedList) int {
			return len(list.Needs)
		},
		func(list client.ProjectCustomerNeedList) (client.ProjectCustomerNeedList, error) {
			items, err := sortByJSONField(list.Needs, options.sortField, options.sortOrder)
			list.Needs = items
			return list, err
		},
		func(command *cobra.Command, item client.CustomerNeedSummary) error {
			return writeCustomerNeed(command, options, item)
		},
		func(list client.ProjectCustomerNeedList) []client.CustomerNeedSummary {
			return list.Needs
		},
	)
}

func addProjectRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectRelationChildListCommand(
		ctx,
		root,
		options,
		"relations PROJECT_ID",
		"List project relations",
		"relations",
		client.ListProjectRelationsForProject,
	)
}

func addProjectTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		"teams PROJECT_ID",
		"List project teams",
		"teams",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectTeamList, error) {
			return client.ListProjectTeams(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectTeamList) int {
			return len(list.Teams)
		},
		func(list client.ProjectTeamList) (client.ProjectTeamList, error) {
			items, err := sortByJSONField(list.Teams, options.sortField, options.sortOrder)
			list.Teams = items
			return list, err
		},
		func(command *cobra.Command, item client.TeamSummary) error {
			return writeTeam(command, options, item)
		},
		func(list client.ProjectTeamList) []client.TeamSummary {
			return list.Teams
		},
	)
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

func addProjectRelationChildListCommand(
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch func(context.Context, graphql.Client, string, int) (client.ProjectProjectRelationList, error),
) {
	addProjectChildListCommand(
		ctx,
		root,
		options,
		use,
		short,
		limitHelp,
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectProjectRelationList, error) {
			return fetch(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectProjectRelationList) int {
			return len(list.Relations)
		},
		func(list client.ProjectProjectRelationList) (client.ProjectProjectRelationList, error) {
			items, err := sortByJSONField(list.Relations, options.sortField, options.sortOrder)
			list.Relations = items
			return list, err
		},
		func(command *cobra.Command, item client.ProjectRelationSummary) error {
			return writeProjectRelation(command, options, item)
		},
		func(list client.ProjectProjectRelationList) []client.ProjectRelationSummary {
			return list.Relations
		},
	)
}

func addProjectChildListCommand[List any, Item any](
	ctx context.Context,
	root *cobra.Command,
	options *rootOptions,
	use string,
	short string,
	limitHelp string,
	fetch func(commandRuntime, string, int) (List, error),
	count func(List) int,
	sortList func(List) (List, error),
	writeItem func(*cobra.Command, Item) error,
	items func(List) []Item,
) {
	limit := 50
	command := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			list, err := fetch(runtime, args[0], limit)
			if err != nil {
				return err
			}
			if err := ensureNonEmpty(options, count(list)); err != nil {
				return err
			}
			list, err = sortList(list)
			if err != nil {
				return err
			}
			if options.json {
				return writeJSONValue(command, options, list)
			}
			for _, item := range items(list) {
				if err := writeItem(command, item); err != nil {
					return err
				}
			}

			return nil
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum "+limitHelp+" to return")
	root.AddCommand(command)
}

func writeProjectHistory(command *cobra.Command, options *rootOptions, history client.ProjectHistorySummary) error {
	if wrote, err := writeIDOnly(command, options, history.ID); wrote || err != nil {
		return err
	}
	if options.quiet {
		return nil
	}
	if options.json {
		return writeJSONValue(command, options, history)
	}

	return render.WriteLine(
		command.OutOrStdout(),
		"%s project %s entries %d",
		history.ID,
		history.ProjectID,
		history.EntryCount,
	)
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
