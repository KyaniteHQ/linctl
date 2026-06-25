package cli

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// cycleCreator is the Command Port the Cycle create command depends on.
type cycleCreator interface {
	CreateCycle(ctx context.Context, request client.CycleCreateRequest) (client.CycleSummary, error)
}

// cycleUpdater is the Command Port the Cycle update command depends on.
type cycleUpdater interface {
	UpdateCycle(ctx context.Context, request client.CycleUpdateRequest) (client.CycleSummary, error)
}

// cycleArchiver is the Command Port the Cycle archive command depends on.
type cycleArchiver interface {
	ArchiveCycle(ctx context.Context, cycleID string) (client.CycleSummary, error)
}
