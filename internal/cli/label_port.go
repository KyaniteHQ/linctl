package cli

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// labelCreator is the Command Port the label create command depends on.
type labelCreator interface {
	CreateLabel(ctx context.Context, request client.LabelCreateRequest) (client.LabelSummary, error)
}

// labelUpdater is the Command Port the label update command depends on.
type labelUpdater interface {
	UpdateLabel(ctx context.Context, request client.LabelUpdateRequest) (client.LabelSummary, error)
}

// labelRetirer is the Command Port the label retire command depends on.
type labelRetirer interface {
	RetireLabel(ctx context.Context, id string, orgWide bool) (client.LabelSummary, error)
}

// labelRestorer is the Command Port the label restore command depends on.
type labelRestorer interface {
	RestoreLabel(ctx context.Context, id string, orgWide bool) (client.LabelSummary, error)
}
