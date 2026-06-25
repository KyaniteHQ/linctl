package cli

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// projectMilestoneCreator is the Command Port the ProjectMilestone create command depends on.
type projectMilestoneCreator interface {
	CreateProjectMilestone(
		ctx context.Context,
		request client.ProjectMilestoneCreateRequest,
	) (client.ProjectMilestoneSummary, error)
}

// projectMilestoneUpdater is the Command Port the ProjectMilestone update command depends on.
type projectMilestoneUpdater interface {
	UpdateProjectMilestone(
		ctx context.Context,
		request client.ProjectMilestoneUpdateRequest,
	) (client.ProjectMilestoneSummary, error)
}
