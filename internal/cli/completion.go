package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// completionListLimit caps how many entities a dynamic completion fetches.
const completionListLimit = 100

// completionLoader produces completion candidates from a built runtime.
type completionLoader func(context.Context, commandRuntime) ([]string, error)

// completionValues builds a runtime and runs load, degrading to no candidates on
// any error so shell completion never surfaces an error or a stack trace.
func completionValues(
	ctx context.Context,
	options *rootOptions,
	load completionLoader,
) ([]string, cobra.ShellCompDirective) {
	runtime, err := buildCommandRuntime(ctx, options)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	values, err := load(ctx, runtime)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return values, cobra.ShellCompDirectiveNoFileComp
}

func teamKeyCandidates(ctx context.Context, runtime commandRuntime) ([]string, error) {
	teams, err := client.ListTeams(ctx, runtime.graphqlClient, completionListLimit)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(teams.Teams))
	for _, team := range teams.Teams {
		keys = append(keys, team.Key)
	}

	return keys, nil
}

func projectIDCandidates(ctx context.Context, runtime commandRuntime) ([]string, error) {
	projects, err := client.ListProjects(ctx, runtime.graphqlClient, completionListLimit)
	if err != nil {
		return nil, err
	}
	candidates := make([]string, 0, len(projects.Projects))
	for _, project := range projects.Projects {
		// "id\tdescription" — cobra renders the text after the tab as the hint.
		candidates = append(candidates, project.ID+"\t"+project.Name)
	}

	return candidates, nil
}

func workflowStateTypeCandidates(ctx context.Context, runtime commandRuntime) ([]string, error) {
	states, err := client.ListWorkflowStates(ctx, runtime.graphqlClient, completionListLimit)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	types := make([]string, 0, len(states.WorkflowStates))
	for _, state := range states.WorkflowStates {
		if state.Type == "" || seen[state.Type] {
			continue
		}
		seen[state.Type] = true
		types = append(types, state.Type)
	}

	return types, nil
}

// flagCompletion wraps a loader as a cobra flag completion function.
func flagCompletion(
	ctx context.Context,
	options *rootOptions,
	load completionLoader,
) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return completionValues(ctx, options, load)
	}
}

// firstArgCompletion wraps a loader as a cobra ValidArgsFunction that completes
// only the first positional argument.
func firstArgCompletion(
	ctx context.Context,
	options *rootOptions,
	load completionLoader,
) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		return completionValues(ctx, options, load)
	}
}

// registerFlagCompletion attaches a dynamic completion to a flag. The only error
// RegisterFlagCompletionFunc returns is for an undefined flag, which is a
// programming error covered by tests, so it is ignored here.
func registerFlagCompletion(
	command *cobra.Command,
	flag string,
	completion func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective),
) {
	_ = command.RegisterFlagCompletionFunc(flag, completion) //nolint:errcheck // flag is defined before registration
}

// registerGlobalCompletions wires dynamic completion for the persistent
// --team and --project flags.
func registerGlobalCompletions(ctx context.Context, command *cobra.Command, options *rootOptions) {
	registerFlagCompletion(command, "team", flagCompletion(ctx, options, teamKeyCandidates))
	registerFlagCompletion(command, "project", flagCompletion(ctx, options, projectIDCandidates))
}

// registerStateCompletion wires dynamic completion for a command's --state flag.
func registerStateCompletion(ctx context.Context, command *cobra.Command, options *rootOptions) {
	registerFlagCompletion(command, "state", flagCompletion(ctx, options, workflowStateTypeCandidates))
}
