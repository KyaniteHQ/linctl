package cli

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// projectLabelCreator is the Command Port the project-label create command depends on.
type projectLabelCreator interface {
	CreateProjectLabel(
		ctx context.Context,
		request client.ProjectLabelCreateRequest,
	) (client.ProjectLabelSummary, error)
}

// projectLabelUpdater is the Command Port the project-label update command depends on.
type projectLabelUpdater interface {
	UpdateProjectLabel(
		ctx context.Context,
		request client.ProjectLabelUpdateRequest,
	) (client.ProjectLabelSummary, error)
}

// projectLabelRetirer is the Command Port the project-label retire command depends on.
type projectLabelRetirer interface {
	RetireProjectLabel(ctx context.Context, id string, orgWide bool) (client.ProjectLabelSummary, error)
}

// projectLabelRestorer is the Command Port the project-label restore command depends on.
type projectLabelRestorer interface {
	RestoreProjectLabel(ctx context.Context, id string, orgWide bool) (client.ProjectLabelSummary, error)
}
