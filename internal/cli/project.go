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
	addProjectAllCommand(ctx, projectCommand, options)
	addProjectGetCommand(ctx, projectCommand, options)
	addProjectAttachmentsCommand(ctx, projectCommand, options)
	addProjectDocumentsCommand(ctx, projectCommand, options)
	addProjectExternalLinksCommand(ctx, projectCommand, options)
	addProjectHistoryCommand(ctx, projectCommand, options)
	addProjectInitiativeLinksCommand(ctx, projectCommand, options)
	addProjectInitiativesCommand(ctx, projectCommand, options)
	addProjectInverseRelationsCommand(ctx, projectCommand, options)
	addProjectIssuesCommand(ctx, projectCommand, options)
	addProjectCommentsCommand(ctx, projectCommand, options)
	addProjectLabelsCommand(ctx, projectCommand, options)
	addProjectMembersCommand(ctx, projectCommand, options)
	addProjectNeedsCommand(ctx, projectCommand, options)
	addProjectRelationsCommand(ctx, projectCommand, options)
	addProjectTeamsCommand(ctx, projectCommand, options)
	addProjectUpdatesCommand(ctx, projectCommand, options)
	addProjectFilterSuggestionCommand(ctx, projectCommand, options)
	addProjectCreateCommand(ctx, projectCommand, options)
	addProjectUpdateCommand(ctx, projectCommand, options)
	addProjectArchiveCommand(ctx, projectCommand, options)
	addProjectOpenCommand(ctx, projectCommand, options)
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
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadProjectsByTeam,
				projectPageWithItems,
				writeProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum projects to return")
	root.AddCommand(command)
}

func loadProjectsByTeam(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectList, []client.ProjectSummary, error) {
	target, err := runtime.resolveTarget(ctx)
	if err != nil {
		return client.ProjectList{}, nil, err
	}
	projects, err := client.ListProjectsByTeam(ctx, runtime.graphqlClient, target.Team.ID, limit)

	return projects, projects.Projects, err
}

func addProjectAllCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	limit := 50
	command := &cobra.Command{
		Use:   "all",
		Short: "List visible Linear projects across the organization",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return runReadListCommand(
				ctx,
				command,
				nil,
				options,
				limit,
				loadProjectsAll,
				projectPageWithItems,
				writeProject,
			)
		},
	}
	command.Flags().IntVar(&limit, "limit", limit, "maximum projects to return")
	root.AddCommand(command)
}

func loadProjectsAll(
	ctx context.Context,
	runtime commandRuntime,
	_ []string,
	limit int,
) (client.ProjectList, []client.ProjectSummary, error) {
	projects, err := client.ListProjects(ctx, runtime.graphqlClient, limit)

	return projects, projects.Projects, err
}

func addProjectGetCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:               "get PROJECT_ID",
		Short:             "Get one project by id or slug",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: firstArgCompletion(ctx, options, projectIDCandidates),
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
	addChildListCommand(
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
		writeAttachment,
		func(list client.ProjectAttachmentList) []client.AttachmentSummary {
			return list.Attachments
		},
	)
}

func addProjectDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
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
		writeDocument,
		func(list client.ProjectDocumentList) []client.DocumentSummary {
			return list.Documents
		},
	)
}

func addProjectExternalLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
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
		writeEntityExternalLink,
		func(list client.ProjectExternalLinkList) []client.EntityExternalLinkSummary {
			return list.Links
		},
	)
}

func addProjectHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
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
		writeProjectHistory,
		func(list client.ProjectHistoryList) []client.ProjectHistorySummary {
			return list.History
		},
	)
}

func addProjectInitiativeLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
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
		writeInitiativeToProject,
		func(list client.ProjectInitiativeToProjectList) []client.InitiativeToProjectSummary {
			return list.Associations
		},
	)
}

func addProjectInitiativesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
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
		writeInitiative,
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
	addChildListCommand(
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
		writeIssue,
		func(list client.ProjectIssueList) []client.IssueSummary {
			return list.Issues
		},
	)
}

func addProjectCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"comments PROJECT_ID",
		"List project comments without body content",
		"comments",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectCommentList, error) {
			return client.ListProjectComments(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectCommentList) int {
			return len(list.Comments)
		},
		func(list client.ProjectCommentList) (client.ProjectCommentList, error) {
			items, err := sortByJSONField(list.Comments, options.sortField, options.sortOrder)
			list.Comments = items
			return list, err
		},
		writeCommentMetadata,
		func(list client.ProjectCommentList) []client.CommentMetadataSummary {
			return list.Comments
		},
	)
}

func addProjectLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
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
		writeProjectLabel,
		func(list client.ProjectProjectLabelList) []client.ProjectLabelSummary {
			return list.ProjectLabels
		},
	)
}

func writeCommentMetadata(command *cobra.Command, options *rootOptions, comment client.CommentMetadataSummary) error {
	return writeItem(command, options, comment, comment.ID,
		func(command *cobra.Command, _ *rootOptions, comment client.CommentMetadataSummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s %s %s",
				comment.ID,
				emptyDash(comment.DisplayName),
				comment.CreatedAt,
			)
		})
}

func writeProjectMember(command *cobra.Command, options *rootOptions, member client.ProjectMember) error {
	return writeItem(command, options, member, member.ID,
		func(command *cobra.Command, _ *rootOptions, member client.ProjectMember) error {
			return render.WriteLine(command.OutOrStdout(), "%s %s", member.ID, member.DisplayName)
		})
}

func addProjectMembersCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"members PROJECT_ID",
		"List project members",
		"members",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectMemberList, error) {
			return client.ListProjectMembers(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectMemberList) int {
			return len(list.Members)
		},
		func(list client.ProjectMemberList) (client.ProjectMemberList, error) {
			items, err := sortByJSONField(list.Members, options.sortField, options.sortOrder)
			list.Members = items
			return list, err
		},
		writeProjectMember,
		func(list client.ProjectMemberList) []client.ProjectMember {
			return list.Members
		},
	)
}

func addProjectNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
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
		writeCustomerNeed,
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
	addChildListCommand(
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
		writeTeam,
		func(list client.ProjectTeamList) []client.TeamSummary {
			return list.Teams
		},
	)
}

func addProjectFilterSuggestionCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	teamID := ""
	command := &cobra.Command{
		Use:   "filter-suggestion PROMPT",
		Short: "Suggest a project filter from a text prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			suggestion, err := client.GetProjectFilterSuggestion(ctx, runtime.graphqlClient, args[0], teamID)
			if err != nil {
				return err
			}

			return writeProjectFilterSuggestion(command, options, suggestion)
		},
	}
	command.Flags().StringVar(&teamID, "team-id", teamID, "optional team id for team-scoped project views")
	root.AddCommand(command)
}

func addProjectUpdatesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"updates PROJECT_ID",
		"List project status updates",
		"project updates",
		func(runtime commandRuntime, projectID string, limit int) (client.ProjectUpdateList, error) {
			return client.ListProjectUpdates(ctx, runtime.graphqlClient, projectID, limit)
		},
		func(list client.ProjectUpdateList) int {
			return len(list.Updates)
		},
		func(list client.ProjectUpdateList) (client.ProjectUpdateList, error) {
			items, err := sortByJSONField(list.Updates, options.sortField, options.sortOrder)
			list.Updates = items
			return list, err
		},
		writeProjectChildUpdate,
		func(list client.ProjectUpdateList) []client.ProjectUpdateSummary {
			return list.Updates
		},
	)
}

func writeProjectChildUpdate(command *cobra.Command, options *rootOptions, update client.ProjectUpdateSummary) error {
	return writeItem(command, options, update, update.ID,
		func(command *cobra.Command, _ *rootOptions, update client.ProjectUpdateSummary) error {
			return render.WriteLine(command.OutOrStdout(), "%s %s %s", update.ID, update.Health, update.DisplayName)
		})
}

func writeProjectFilterSuggestion(
	command *cobra.Command,
	options *rootOptions,
	suggestion client.ProjectFilterSuggestion,
) error {
	return writeItem(command, options, suggestion, suggestion.LogID,
		func(command *cobra.Command, _ *rootOptions, suggestion client.ProjectFilterSuggestion) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"log_id=%s filter=%s",
				emptyDash(suggestion.LogID),
				emptyDash(string(suggestion.Filter)),
			)
		})
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
	addChildListCommand(
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
		writeProjectRelation,
		func(list client.ProjectProjectRelationList) []client.ProjectRelationSummary {
			return list.Relations
		},
	)
}

func writeProjectHistory(command *cobra.Command, options *rootOptions, history client.ProjectHistorySummary) error {
	return writeItem(command, options, history, history.ID,
		func(command *cobra.Command, _ *rootOptions, history client.ProjectHistorySummary) error {
			return render.WriteLine(
				command.OutOrStdout(),
				"%s project %s entries %d",
				history.ID,
				history.ProjectID,
				history.EntryCount,
			)
		})
}

func writeProject(command *cobra.Command, options *rootOptions, project client.ProjectSummary) error {
	return writeItem(command, options, project, project.ID,
		func(command *cobra.Command, options *rootOptions, project client.ProjectSummary) error {
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
		})
}

func projectPageWithItems(page client.ProjectList, projects []client.ProjectSummary) client.ProjectList {
	page.Projects = projects
	return page
}
