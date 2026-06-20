package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
	"github.com/KyaniteHQ/linctl/internal/render"
)

type doctorReport struct {
	Config string              `json:"config"`
	Token  string              `json:"token"`
	Target doctorTargetReport  `json:"target"`
	Viewer client.TargetViewer `json:"viewer"`
}

type doctorTargetReport struct {
	Status    string                  `json:"status"`
	OrgID     string                  `json:"org_id"`
	TeamKey   string                  `json:"team_key"`
	TeamID    string                  `json:"team_id"`
	ProjectID string                  `json:"project_id,omitempty"`
	Expected  map[string]string       `json:"expected"`
	Resolved  map[string]string       `json:"resolved"`
	Project   *client.ResolvedProject `json:"project,omitempty"`
}

func addDoctorCommand(ctx context.Context, root *cobra.Command, options *rootOptions) {
	root.AddCommand(&cobra.Command{
		Use:   "doctor",
		Short: "Check linctl config, token, and target health",
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
			if options.quiet {
				return nil
			}

			report := newDoctorReport(target)
			if options.json {
				return writeJSONValue(command, options, report)
			}

			return render.WriteLine(
				command.OutOrStdout(),
				"config %s\n token %s\n target %s %s/%s project %s",
				report.Config,
				report.Token,
				report.Target.Status,
				report.Target.TeamKey,
				report.Target.TeamID,
				report.Target.ProjectID,
			)
		},
	})
}

func newDoctorReport(target client.ResolvedTarget) doctorReport {
	return doctorReport{
		Config: "ok",
		Token:  "set",
		Target: doctorTargetReport{
			Status:    "confirmed",
			OrgID:     target.Org.ID,
			TeamKey:   target.Team.Key,
			TeamID:    target.Team.ID,
			ProjectID: projectID(target.Project),
			Expected:  targetMap(target.Expected),
			Resolved:  targetMap(target.Resolved),
			Project:   target.Project,
		},
		Viewer: target.Viewer,
	}
}

func targetMap(target config.Target) map[string]string {
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
