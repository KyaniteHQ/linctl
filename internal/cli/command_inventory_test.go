package cli

import (
	"context"
	"slices"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

func Test_CommandInventory_exposes_stable_public_command_surface(t *testing.T) {
	root := NewRootCommand(context.Background(), BuildInfo{})

	inventory := CommandInventory(root)
	require.NotEmpty(t, inventory)

	paths := make([]string, 0, len(inventory))
	pathsByName := map[string]bool{}
	aliasesByName := map[string]bool{}
	commandsByPath := map[string]CommandInfo{}
	for _, command := range inventory {
		require.NotEmpty(t, command.Path)
		require.NotContains(t, command.Path, "linctl ")
		require.NotEmpty(t, command.UseLine)
		require.NotEmpty(t, command.Short)
		require.NotEmpty(t, command.Entity)
		require.NotEmpty(t, command.DocCategory)
		require.NotEmpty(t, command.Safety)
		paths = append(paths, command.Path)
		pathsByName[command.Path] = true
		commandsByPath[command.Path] = command
		for _, alias := range command.Aliases {
			aliasesByName[alias] = true
		}
	}

	require.True(t, slices.IsSorted(paths), "command inventory must remain stable")
	require.True(t, pathsByName["issue list"])
	require.True(t, pathsByName["organization teams"])
	require.False(t, pathsByName["completion"])
	require.False(t, pathsByName["help"])
	require.True(t, aliasesByName["issue get ISSUE_ID"])

	attachmentList := commandsByPath["attachment list"]
	require.Equal(t, CommandSafetyRead, attachmentList.Safety)
	require.Equal(t, "attachment", attachmentList.Entity)
	require.Equal(t, "attachments", attachmentList.CollectionKey)

	documentCreate := commandsByPath["document create"]
	require.Equal(t, CommandSafetyWrite, documentCreate.Safety)
	require.Equal(t, "document", documentCreate.Entity)
}

func Test_EnrichCommandInventory_adds_graphql_backing_without_mutating_source(t *testing.T) {
	commands := []CommandInfo{
		{
			Path:   "issue list",
			Safety: CommandSafetyRead,
		},
		{
			Path:    "issue get",
			Aliases: []string{"issue get ISSUE_ID"},
			Safety:  CommandSafetyRead,
		},
	}
	graphqlRoots := []CommandGraphQLRoot{
		{Kind: "query", Field: "issues", Operation: "issues"},
	}
	backingByPath := map[string]CommandBacking{
		"issue list": {
			OperationBacking: "Query.issues",
			TargetScope:      "Read-only",
			GraphQLRoots:     graphqlRoots,
		},
		"issue get ISSUE_ID": {
			OperationBacking: "Query.issue",
			TargetScope:      "Read-only",
			GraphQLRoots: []CommandGraphQLRoot{
				{Kind: "query", Field: "issue", Operation: "issue"},
			},
		},
	}

	enriched := EnrichCommandInventory(commands, backingByPath)
	graphqlRoots[0].Field = "mutated"

	require.Empty(t, commands[0].OperationBacking)
	require.Equal(t, "Query.issues", enriched[0].OperationBacking)
	require.Equal(t, "Read-only", enriched[0].TargetScope)
	require.Equal(t, []CommandGraphQLRoot{
		{Kind: "query", Field: "issues", Operation: "issues"},
	}, enriched[0].GraphQLRoots)
	require.Equal(t, "Query.issue", enriched[1].OperationBacking)
}

func Test_collectionKeyForPage_uses_typed_list_envelope(t *testing.T) {
	require.Equal(t, "attachments", collectionKeyForPage[client.AttachmentList]())
	require.Equal(t, "users", collectionKeyForPage[client.UserList]())
}

func Test_commandAnnotationHelpers_ignore_empty_values(t *testing.T) {
	command := &cobra.Command{Use: "list"}

	annotateCommand(command, commandCollectionKeyAnnotation, "")
	require.Nil(t, command.Annotations)
	require.Empty(t, commandCollectionKey(command))

	annotateReadCollectionCommand(command, "issues")
	require.Equal(t, "issues", commandCollectionKey(command))
	require.Equal(t, CommandSafetyRead, commandSafety(command))

	annotateCommand(command, commandCollectionKeyAnnotation, "projects")
	require.Equal(t, "projects", commandCollectionKey(command))
}

func Test_commandSafety_classifies_annotations_paths_and_descriptions(t *testing.T) {
	for _, safety := range []CommandSafety{
		CommandSafetyRead,
		CommandSafetyWrite,
		CommandSafetyLocal,
		CommandSafetyUnknown,
	} {
		command := &cobra.Command{
			Use:         "annotated",
			Annotations: map[string]string{commandSafetyAnnotation: string(safety)},
		}
		require.Equal(t, safety, commandSafety(command))
	}

	require.Equal(t, CommandSafetyWrite, commandSafety(commandWithPath("document", "create", "Create document")))
	require.Equal(t, CommandSafetyRead, commandSafety(commandWithPath("issue", "get", "Get issue")))
	require.Equal(t, CommandSafetyLocal, commandSafety(commandWithPath("completion", "bash", "Generate completion")))
	require.Equal(t, CommandSafetyUnknown, commandSafety(commandWithPath("issue", "sync", "Synchronize issue")))
	require.Empty(t, commandEntity(&cobra.Command{}))
	require.Empty(t, commandDocCategory(&cobra.Command{}))
}

func Test_collectionKeyForPage_rejects_ambiguous_or_non_collection_shapes(t *testing.T) {
	type pointerPage struct {
		Items []string `json:"items,omitempty"`
	}
	type noJSONTagPage struct {
		Items []string
	}
	type ignoredJSONTagPage struct {
		Items []string `json:"-"`
	}
	type multipleSlicesPage struct {
		Items  []string `json:"items"`
		Labels []string `json:"labels"`
	}

	require.Equal(t, "items", collectionKeyForPage[*pointerPage]())
	require.Empty(t, collectionKeyForPage[string]())
	require.Empty(t, collectionKeyForPage[noJSONTagPage]())
	require.Empty(t, collectionKeyForPage[ignoredJSONTagPage]())
	require.Empty(t, collectionKeyForPage[multipleSlicesPage]())
}

func commandWithPath(parentUse string, childUse string, childShort string) *cobra.Command {
	parent := &cobra.Command{Use: "linctl"}
	group := &cobra.Command{Use: parentUse}
	child := &cobra.Command{Use: childUse, Short: childShort}

	parent.AddCommand(group)
	group.AddCommand(child)

	return child
}
