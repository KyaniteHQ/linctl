package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// issueTemplateReader isolates the optional template lookup performed before
// issue creation. The create mutation itself crosses the client boundary
// directly after the command has assembled the final request.
type issueTemplateReader interface {
	GetIssueTemplateContent(ctx context.Context, templateID string) (client.IssueTemplateContent, error)
}

// issueReader isolates list dispatch between all visible teams and the pinned
// team selected after target resolution.
type issueReader interface {
	ResolveTarget(ctx context.Context) (client.ResolvedTarget, error)
	ListIssues(ctx context.Context, limit int) (client.IssueList, error)
	ListIssuesByTeam(
		ctx context.Context,
		teamID string,
		limit int,
		filters client.IssueListFilters,
	) (client.IssueList, error)
}

// issueSearcher isolates target resolution followed by team-scoped search.
type issueSearcher interface {
	ResolveTarget(ctx context.Context) (client.ResolvedTarget, error)
	SearchIssuesByTeam(ctx context.Context, teamID string, query string, limit int) (client.IssueList, error)
}

// nextIssuePicker isolates target resolution, candidate listing, optional
// checkout, and the guarded start of the selected issue.
type nextIssuePicker interface {
	ResolveTarget(ctx context.Context) (client.ResolvedTarget, error)
	ListNextIssuesByTeam(ctx context.Context, teamID string, limit int) (client.IssueList, error)
	StartIssue(ctx context.Context, issueID string) (client.IssueSummary, error)
}

type issueTemplateClient struct {
	graphqlClient graphql.Client
}

func (reader issueTemplateClient) GetIssueTemplateContent(
	ctx context.Context,
	templateID string,
) (client.IssueTemplateContent, error) {
	return client.GetIssueTemplateContent(ctx, reader.graphqlClient, templateID)
}

type issueListClient struct {
	graphqlClient graphql.Client
	target        config.Target
}

func issueListClientFor(runtime commandRuntime) issueListClient {
	return issueListClient{graphqlClient: runtime.graphqlClient, target: runtime.config.Target}
}

func (reader issueListClient) ResolveTarget(ctx context.Context) (client.ResolvedTarget, error) {
	return client.ResolveTarget(ctx, reader.graphqlClient, reader.target)
}

func (reader issueListClient) ListIssues(ctx context.Context, limit int) (client.IssueList, error) {
	return client.ListIssues(ctx, reader.graphqlClient, limit)
}

func (reader issueListClient) ListIssuesByTeam(
	ctx context.Context,
	teamID string,
	limit int,
	filters client.IssueListFilters,
) (client.IssueList, error) {
	return client.ListIssuesByTeam(ctx, reader.graphqlClient, teamID, limit, filters)
}

type issueSearchClient struct {
	graphqlClient graphql.Client
	target        config.Target
}

func issueSearchClientFor(runtime commandRuntime) issueSearchClient {
	return issueSearchClient{graphqlClient: runtime.graphqlClient, target: runtime.config.Target}
}

func (searcher issueSearchClient) ResolveTarget(ctx context.Context) (client.ResolvedTarget, error) {
	return client.ResolveTarget(ctx, searcher.graphqlClient, searcher.target)
}

func (searcher issueSearchClient) SearchIssuesByTeam(
	ctx context.Context,
	teamID string,
	query string,
	limit int,
) (client.IssueList, error) {
	return client.SearchIssuesByTeam(ctx, searcher.graphqlClient, teamID, query, limit)
}
