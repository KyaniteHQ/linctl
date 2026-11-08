package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// ProjectInitiativeLinkRequest names the project and Initiative of one
// Initiative-to-Project association. SortOrder only applies to a new link.
type ProjectInitiativeLinkRequest struct {
	ProjectID    string
	InitiativeID string
	SortOrder    *float64
}

// LinkProjectInitiative adds a project to an Initiative after comparing the
// project against the pinned target and the Initiative against the pinned
// organization. A project already linked to the Initiative returns the
// existing link unchanged.
func LinkProjectInitiative(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectInitiativeLinkRequest,
) (InitiativeToProjectSummary, error) {
	guard, err := newProjectInitiativeGuard(ctx, graphqlClient, expected, request)
	if err != nil {
		return InitiativeToProjectSummary{}, err
	}

	return guard.linkProjectInitiative(ctx, request)
}

func (guard *guardedClient) linkProjectInitiative(
	ctx context.Context,
	request ProjectInitiativeLinkRequest,
) (InitiativeToProjectSummary, error) {
	existing, found, err := guard.projectInitiativeLink(ctx, request)
	if err != nil || found {
		return existing, err
	}

	created, err := gql.InitiativeToProjectCreate(ctx, guard.graphqlClient, gql.InitiativeToProjectCreateInput{
		InitiativeId: request.InitiativeID,
		ProjectId:    request.ProjectID,
		SortOrder:    request.SortOrder,
	})
	if err != nil {
		return InitiativeToProjectSummary{}, fmt.Errorf(
			"link project %s to initiative %s: %w", request.ProjectID, request.InitiativeID, err,
		)
	}
	payload := created.InitiativeToProjectCreate
	if err := mutationSuccess(payload.Success, "initiativeToProjectCreate"); err != nil {
		return InitiativeToProjectSummary{}, err
	}

	return initiativeToProjectSummary(payload.InitiativeToProject.InitiativeToProjectSummaryFields), nil
}

// UnlinkProjectInitiative removes a project from an Initiative after the same
// comparisons as LinkProjectInitiative. It deletes only the association and
// returns it; the project and the Initiative stay unchanged.
func UnlinkProjectInitiative(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectInitiativeLinkRequest,
) (InitiativeToProjectSummary, error) {
	guard, err := newProjectInitiativeGuard(ctx, graphqlClient, expected, request)
	if err != nil {
		return InitiativeToProjectSummary{}, err
	}

	return guard.unlinkProjectInitiative(ctx, request)
}

func (guard *guardedClient) unlinkProjectInitiative(
	ctx context.Context,
	request ProjectInitiativeLinkRequest,
) (InitiativeToProjectSummary, error) {
	link, found, err := guard.projectInitiativeLink(ctx, request)
	if err != nil {
		return InitiativeToProjectSummary{}, err
	}
	if !found {
		return InitiativeToProjectSummary{}, fmt.Errorf(
			"%w: project %s is not linked to initiative %s",
			ErrWriteInvalid, request.ProjectID, request.InitiativeID,
		)
	}

	deleted, err := gql.InitiativeToProjectDelete(ctx, guard.graphqlClient, link.ID)
	if err != nil {
		return InitiativeToProjectSummary{}, fmt.Errorf(
			"unlink project %s from initiative %s: %w", request.ProjectID, request.InitiativeID, err,
		)
	}
	if err := mutationSuccess(deleted.InitiativeToProjectDelete.Success, "initiativeToProjectDelete"); err != nil {
		return InitiativeToProjectSummary{}, err
	}

	return link, nil
}

func newProjectInitiativeGuard(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request ProjectInitiativeLinkRequest,
) (*guardedClient, error) {
	if request.ProjectID == "" || request.InitiativeID == "" {
		return nil, fmt.Errorf("%w: project id and initiative id are required", ErrWriteInvalid)
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return nil, err
	}
	if err := guard.requireProject(ctx, request.ProjectID); err != nil {
		return nil, err
	}
	if err := guard.requireInitiative(ctx, request.InitiativeID); err != nil {
		return nil, err
	}

	return guard, nil
}

// projectInitiativeLink finds the live association between the project and
// the Initiative. It refuses to guess when the project has more links than
// one page holds.
func (guard *guardedClient) projectInitiativeLink(
	ctx context.Context,
	request ProjectInitiativeLinkRequest,
) (InitiativeToProjectSummary, bool, error) {
	links, err := ListProjectInitiativeToProjects(ctx, guard.graphqlClient, request.ProjectID, defaultListPageSize)
	if err != nil {
		return InitiativeToProjectSummary{}, false, err
	}
	for _, link := range links.Associations {
		if link.InitiativeID == request.InitiativeID && link.ArchivedAt == "" {
			return link, true, nil
		}
	}
	if links.HasNextPage {
		return InitiativeToProjectSummary{}, false, fmt.Errorf(
			"%w: project %s has more than %d initiative links; cannot find the link to initiative %s",
			ErrWriteInvalid, request.ProjectID, defaultListPageSize, request.InitiativeID,
		)
	}

	return InitiativeToProjectSummary{}, false, nil
}
