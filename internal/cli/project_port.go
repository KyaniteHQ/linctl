package cli

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// projectCreator is the Command Port the project create command depends on.
type projectCreator interface {
	CreateProject(ctx context.Context, request client.ProjectCreateRequest) (client.ProjectSummary, error)
}

// projectUpdater is the Command Port the project update command depends on.
type projectUpdater interface {
	UpdateProject(ctx context.Context, request client.ProjectUpdateRequest) (client.ProjectSummary, error)
}

// projectArchiver is the Command Port the project archive command depends on.
type projectArchiver interface {
	ArchiveProject(ctx context.Context, projectID string) (client.ProjectSummary, error)
}
