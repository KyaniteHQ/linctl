package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/render"
)

func addProjectCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	projectCommand := newGroupCommand("project", "Read and write Linear projects")
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
	addProjectAddTeamCommand(ctx, projectCommand, options)
	addProjectAddLabelCommand(ctx, projectCommand, options)
	addProjectRemoveLabelCommand(ctx, projectCommand, options)
	addProjectOpenCommand(ctx, projectCommand, options)
	addDomainUsageCommand(projectCommand, options, "project")
	root.AddCommand(projectCommand)
}

func addProjectListCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectList, client.ProjectSummary]{
		Use:       "list",
		Short:     "List projects for the resolved team",
		LimitHelp: "projects",
		Args:      cobra.NoArgs,
		Load:      loadProjectsByTeam,
		WriteItem: writeProject,
	})
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
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectList, client.ProjectSummary]{
		Use:       "all",
		Short:     "List visible Linear projects across the organization",
		LimitHelp: "projects",
		Args:      cobra.NoArgs,
		Load:      loadProjectsAll,
		WriteItem: writeProject,
	})
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
	addReadGetCommand(ctx, root, options, readGetSpec[client.ProjectSummary]{
		Use:   "get PROJECT_ID",
		Short: "Get one project by id or slug",
		Configure: func(command *cobra.Command) {
			command.ValidArgsFunction = firstArgCompletion(ctx, options, projectIDCandidates)
		},
		Load: func(ctx context.Context, runtime commandRuntime, id string) (client.ProjectSummary, error) {
			return client.GetProjectByID(ctx, runtime.graphqlClient, id)
		},
		Write: writeProject,
	})
}

func addProjectAttachmentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectAttachmentList, client.AttachmentSummary]{
		Use:       "attachments PROJECT_ID",
		Short:     "List project attachments",
		LimitHelp: "attachments",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectAttachmentList, []client.AttachmentSummary, error) {
			list, err := client.ListProjectAttachments(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Attachments, err
		},
		WriteItem: writeAttachment,
	})
}

func addProjectDocumentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectDocumentList, client.DocumentSummary]{
		Use:       "documents PROJECT_ID",
		Short:     "List project documents",
		LimitHelp: "documents",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectDocumentList, []client.DocumentSummary, error) {
			list, err := client.ListProjectDocuments(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Documents, err
		},
		WriteItem: writeDocument,
	})
}

func addProjectExternalLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	spec := listCommandSpec[client.ProjectExternalLinkList, client.EntityExternalLinkSummary]{
		Use:       "external-links PROJECT_ID",
		Short:     "List project external links",
		LimitHelp: "external links",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectExternalLinkList, []client.EntityExternalLinkSummary, error) {
			list, err := client.ListProjectExternalLinks(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Links, err
		},
		WriteItem: writeEntityExternalLink,
	}
	addListCommand(ctx, root, options, spec)
}

func addProjectHistoryCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectHistoryList, client.ProjectHistorySummary]{
		Use:       "history PROJECT_ID",
		Short:     "List project history",
		LimitHelp: "history entries",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectHistoryList, []client.ProjectHistorySummary, error) {
			list, err := client.ListProjectHistory(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.History, err
		},
		WriteItem: writeProjectHistory,
	})
}

func addProjectInitiativeLinksCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	spec := listCommandSpec[client.ProjectInitiativeToProjectList, client.InitiativeToProjectSummary]{
		Use:       "initiative-links PROJECT_ID",
		Short:     "List project initiative associations",
		LimitHelp: "initiative associations",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectInitiativeToProjectList, []client.InitiativeToProjectSummary, error) {
			list, err := client.ListProjectInitiativeToProjects(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Associations, err
		},
		WriteItem: writeInitiativeToProject,
	}
	addListCommand(ctx, root, options, spec)
}

func addProjectInitiativesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectInitiativeList, client.InitiativeSummary]{
		Use:       "initiatives PROJECT_ID",
		Short:     "List project initiatives",
		LimitHelp: "initiatives",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectInitiativeList, []client.InitiativeSummary, error) {
			list, err := client.ListProjectInitiatives(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Initiatives, err
		},
		WriteItem: writeInitiative,
	})
}

func addProjectInverseRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"inverse-relations PROJECT_ID",
		"List project inverse relations",
		"inverse relations",
		client.ListProjectInverseRelations,
		projectRelationListItems,
		writeProjectRelation,
	)
}

func addProjectIssuesCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectIssueList, client.IssueSummary]{
		Use:       "issues PROJECT_ID",
		Short:     "List project issues",
		LimitHelp: "issues",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectIssueList, []client.IssueSummary, error) {
			list, err := client.ListProjectIssues(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Issues, err
		},
		WriteItem: writeIssue,
	})
}

func addProjectCommentsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectCommentList, client.CommentMetadataSummary]{
		Use:       "comments PROJECT_ID",
		Short:     "List project comments without body content",
		LimitHelp: "comments",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectCommentList, []client.CommentMetadataSummary, error) {
			list, err := client.ListProjectComments(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Comments, err
		},
		WriteItem: writeCommentMetadata,
	})
}

func addProjectLabelsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectProjectLabelList, client.ProjectLabelSummary]{
		Use:       "labels PROJECT_ID",
		Short:     "List project labels",
		LimitHelp: "labels",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectProjectLabelList, []client.ProjectLabelSummary, error) {
			list, err := client.ListLabelsForProject(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.ProjectLabels, err
		},
		WriteItem: writeProjectLabel,
	})
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
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectMemberList, client.ProjectMember]{
		Use:       "members PROJECT_ID",
		Short:     "List project members",
		LimitHelp: "members",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectMemberList, []client.ProjectMember, error) {
			list, err := client.ListProjectMembers(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Members, err
		},
		WriteItem: writeProjectMember,
	})
}

func addProjectNeedsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectCustomerNeedList, client.CustomerNeedSummary]{
		Use:       "needs PROJECT_ID",
		Short:     "List project customer needs",
		LimitHelp: "customer needs",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectCustomerNeedList, []client.CustomerNeedSummary, error) {
			list, err := client.ListProjectNeeds(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Needs, err
		},
		WriteItem: writeCustomerNeed,
	})
}

func addProjectRelationsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addChildListCommand(
		ctx,
		root,
		options,
		"relations PROJECT_ID",
		"List project relations",
		"relations",
		client.ListProjectRelationsForProject,
		projectRelationListItems,
		writeProjectRelation,
	)
}

func addProjectTeamsCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectTeamList, client.TeamSummary]{
		Use:       "teams PROJECT_ID",
		Short:     "List project teams",
		LimitHelp: "teams",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectTeamList, []client.TeamSummary, error) {
			list, err := client.ListProjectTeams(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Teams, err
		},
		WriteItem: writeTeam,
	})
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
	addListCommand(ctx, root, options, listCommandSpec[client.ProjectUpdateList, client.ProjectUpdateSummary]{
		Use:       "updates PROJECT_ID",
		Short:     "List project status updates",
		LimitHelp: "project updates",
		Args:      cobra.ExactArgs(1),
		Load: func(
			ctx context.Context, runtime commandRuntime, args []string, limit int,
		) (client.ProjectUpdateList, []client.ProjectUpdateSummary, error) {
			list, err := client.ListProjectUpdates(ctx, runtime.graphqlClient, args[0], limit)
			return list, list.Updates, err
		},
		WriteItem: writeProjectChildUpdate,
	})
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

func projectRelationListItems(list client.ProjectProjectRelationList) []client.ProjectRelationSummary {
	return list.Relations
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
