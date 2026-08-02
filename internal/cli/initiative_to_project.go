package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func addInitiativeToProjectCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	addReadListGetCommand(
		ctx,
		root,
		options,
		readListGetSpec[client.InitiativeToProjectList, client.InitiativeToProjectSummary]{
			Use:       "initiative-to-project",
			Short:     "Read Linear Initiative-to-Project associations",
			ListShort: "List visible Initiative-to-Project associations",
			LimitHelp: "maximum Initiative-to-Project associations to return",
			GetUse:    "get INITIATIVE_TO_PROJECT_ID",
			GetShort:  "Get one Initiative-to-Project association by id",
			LoadList:  clientList(client.ListInitiativeToProjects),
			LoadGet:   clientGet(client.GetInitiativeToProjectByID),
			WriteItem: writeInitiativeToProject,
		},
	)
}

func writeInitiativeToProject(
	command *cobra.Command,
	options *rootOptions,
	association client.InitiativeToProjectSummary,
) error {
	return writeItemLine(
		command, options, association, association.ID,
		"%s %s -> %s order %s",
		association.ID,
		association.InitiativeName,
		association.ProjectName,
		association.SortOrder,
	)
}
