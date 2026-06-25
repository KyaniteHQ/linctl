package cli

import (
	"context"

	"github.com/KyaniteHQ/linctl/internal/client"
)

// projectUpdateCreator is the Command Port the project-update create command
// depends on to reach Linear, decoupled from the GraphQL transport. It returns a
// domain summary, so the command's body-resolution and health-normalization
// logic is tested through an in-memory fake rather than canned GraphQL JSON. The
// shared issueClientAdapter is the production adapter; the guarded-write
// comparison stays in internal/client and the adapter only forwards.
type projectUpdateCreator interface {
	CreateProjectUpdate(
		ctx context.Context,
		request client.ProjectUpdateCreateRequest,
	) (client.ProjectUpdateSummary, error)
}

// The single shared adapter satisfies the project-update port too, so no
// separate adapter struct is needed; the assertion fails the build if its
// forwarding (and thus the write guard) ever stops matching the port.
var _ projectUpdateCreator = issueClientAdapter{}

func (adapter issueClientAdapter) CreateProjectUpdate(
	ctx context.Context,
	request client.ProjectUpdateCreateRequest,
) (client.ProjectUpdateSummary, error) {
	return client.CreateProjectUpdate(ctx, adapter.graphqlClient, adapter.target, request)
}
