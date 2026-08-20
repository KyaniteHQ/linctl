package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
	"github.com/KyaniteHQ/linctl/internal/render"
)

const (
	initPinFileName   = ".linctl.toml"
	initTeamListLimit = 250
)

type initRequest struct {
	TeamKey   string
	TeamID    string
	ProjectID string
	Path      string
}

type initResult struct {
	Path   string            `json:"path"`
	Target map[string]string `json:"target"`
}

func initTargetJSON(target config.Target) map[string]string {
	values := map[string]string{
		"org_id":   target.OrgID,
		"team_key": target.TeamKey,
		"team_id":  target.TeamID,
	}
	if target.ProjectID != "" {
		values["project_id"] = target.ProjectID
	}

	return values
}

func addInitCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	request := initRequest{}
	command := &cobra.Command{
		Use:   "init",
		Short: "Create a target-only .linctl.toml pin from the active credential",
		Long: "Write .linctl.toml in the current working directory from the active " +
			"OAuth credential. The file contains only [target] org_id, team_key, " +
			"team_id, and optional project_id. Never writes auth material. Refuses " +
			"to overwrite an existing file. When multiple teams are visible, pass " +
			"--team KEY or --team-id ID (discover with linctl team list).",
		Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			runtime, err := buildCommandRuntime(ctx, options)
			if err != nil {
				return err
			}
			if request.Path == "" {
				request.Path = initPinFileName
			}
			result, err := runInit(ctx, runtime.graphqlClient, request)
			if err != nil {
				return err
			}
			if options.quiet {
				return nil
			}
			if options.json {
				return writeJSONValue(command, options, result)
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"wrote %s org_id=%s team_key=%s team_id=%s project_id=%s",
				result.Path,
				result.Target["org_id"],
				result.Target["team_key"],
				result.Target["team_id"],
				result.Target["project_id"],
			)
		},
	}
	command.Flags().StringVar(&request.TeamKey, "team", "", "team key when multiple teams are visible")
	command.Flags().StringVar(&request.TeamID, "team-id", "", "team id when multiple teams are visible")
	command.Flags().StringVar(&request.ProjectID, "project", "", "optional project id to pin after team selection")
	registerFlagCompletion(command, "team", flagCompletion(ctx, options, teamKeyCandidates))
	registerFlagCompletion(command, "project", flagCompletion(ctx, options, projectIDCandidates))
	addCommandWithSafety(root, CommandSafetyLocal, command)
}

func runInit(
	ctx context.Context,
	graphqlClient graphql.Client,
	request initRequest,
) (initResult, error) {
	if request.TeamKey != "" && request.TeamID != "" {
		return initResult{}, fmt.Errorf(
			"%w: use only one of --team or --team-id",
			client.ErrWriteInvalid,
		)
	}

	teams, err := client.ListTeams(ctx, graphqlClient, initTeamListLimit)
	if err != nil {
		return initResult{}, err
	}
	selected, err := selectInitTeam(teams, request.TeamKey, request.TeamID)
	if err != nil {
		return initResult{}, err
	}

	target := config.Target{
		OrgID:   selected.OrgID,
		TeamKey: selected.Key,
		TeamID:  selected.ID,
	}
	if request.ProjectID != "" {
		if err := verifyInitProject(ctx, graphqlClient, selected, request.ProjectID); err != nil {
			return initResult{}, err
		}
		target.ProjectID = request.ProjectID
	}

	path := request.Path
	if !filepath.IsAbs(path) {
		// Keep relative paths as written so the human message matches cwd discovery.
		path = filepath.Clean(path)
	}
	if err := config.WritePin(path, target); err != nil {
		return initResult{}, err
	}

	return initResult{Path: path, Target: initTargetJSON(target)}, nil
}

func selectInitTeam(
	teams client.TeamList,
	teamKey string,
	teamID string,
) (client.TeamSummary, error) {
	if len(teams.Teams) == 0 {
		return client.TeamSummary{}, fmt.Errorf(
			"%w: no teams visible to the active credential; check auth and Linear access",
			client.ErrWriteInvalid,
		)
	}

	if teamID != "" {
		for _, team := range teams.Teams {
			if team.ID == teamID {
				return team, nil
			}
		}

		return client.TeamSummary{}, fmt.Errorf(
			"%w: team id %s is not visible; run linctl team list",
			client.ErrWriteInvalid,
			teamID,
		)
	}
	if teamKey != "" {
		for _, team := range teams.Teams {
			if team.Key == teamKey {
				return team, nil
			}
		}

		return client.TeamSummary{}, fmt.Errorf(
			"%w: team key %s is not visible; run linctl team list",
			client.ErrWriteInvalid,
			teamKey,
		)
	}

	if len(teams.Teams) == 1 && !teams.HasNextPage {
		return teams.Teams[0], nil
	}

	return client.TeamSummary{}, fmt.Errorf(
		"%w: multiple teams are visible; pass --team KEY or --team-id ID (discover with linctl team list)",
		client.ErrWriteInvalid,
	)
}

func verifyInitProject(
	ctx context.Context,
	graphqlClient graphql.Client,
	team client.TeamSummary,
	projectID string,
) error {
	project, err := client.GetProjectByID(ctx, graphqlClient, projectID)
	if err != nil {
		return err
	}
	for _, projectTeam := range project.Teams {
		if projectTeam.ID == team.ID || projectTeam.Key == team.Key {
			return nil
		}
	}
	if project.TeamsTruncated {
		return fmt.Errorf(
			"%w: project %s team membership could not be fully verified against team %s; refuse to pin",
			client.ErrWriteInvalid,
			projectID,
			team.Key,
		)
	}

	return fmt.Errorf(
		"%w: project %s is not attached to team %s/%s",
		client.ErrWriteInvalid,
		projectID,
		team.Key,
		team.ID,
	)
}
