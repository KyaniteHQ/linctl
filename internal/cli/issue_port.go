package cli

import (
	"context"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// The Command Port is the narrow, domain-typed interface an issue command depends
// on to reach Linear, decoupled from the GraphQL transport. Each interface is
// defined by its consumer and returns domain summaries rather than GraphQL
// responses, so command logic is tested through it without canned transport
// payloads. issueClientAdapter is the production adapter over the client package;
// tests supply an in-memory fake. The guarded-write comparison stays in
// internal/client — the adapter only forwards.

// issueTemplateReader reads a Linear template's issue defaults (a free read).
type issueTemplateReader interface {
	GetIssueTemplateContent(ctx context.Context, templateID string) (client.IssueTemplateContent, error)
}

// issueCreator is the port the issue create command depends on.
type issueCreator interface {
	issueTemplateReader
	CreateIssue(ctx context.Context, request client.IssueCreateRequest) (client.IssueSummary, error)
}

// issueCloser is the port the issue close command depends on.
type issueCloser interface {
	CloseIssue(ctx context.Context, issueID string) (client.IssueSummary, error)
}

// issueUpdater is the port the issue update command depends on.
type issueUpdater interface {
	UpdateIssue(ctx context.Context, request client.IssueUpdateRequest) (client.IssueSummary, error)
}

// issueCommenter is the port the issue comment and reply commands depend on.
type issueCommenter interface {
	CommentOnIssue(ctx context.Context, request client.IssueCommentRequest) (client.IssueCommentResult, error)
}

// issueAttachmentLinker is the port the issue link command depends on.
type issueAttachmentLinker interface {
	LinkIssueAttachment(ctx context.Context, request client.AttachmentLinkRequest) (client.AttachmentSummary, error)
}

// issueRelationCreator is the port the issue relate command depends on.
type issueRelationCreator interface {
	CreateIssueRelation(
		ctx context.Context,
		request client.IssueRelationCreateRequest,
	) (client.IssueRelationSummary, error)
}

// issueRelationDeleter is the port the issue unrelate command depends on.
type issueRelationDeleter interface {
	DeleteIssueRelation(ctx context.Context, relationID string) (string, error)
}

// issueReader is the port the issue list command depends on for its dispatch:
// it either lists across teams or resolves the pinned team and lists with the
// assembled filters.
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

// issueClientAdapter satisfies the issue command ports by forwarding to the
// client package's guarded free functions with the runtime's transport and
// pinned target. It is a pass-through adapter: large surface, trivial body.
type issueClientAdapter struct {
	graphqlClient graphql.Client
	target        config.Target
}

// issueAdapterFor builds the production issue port from a resolved runtime.
func issueAdapterFor(runtime commandRuntime) issueClientAdapter {
	return issueClientAdapter{graphqlClient: runtime.graphqlClient, target: runtime.config.Target}
}

func (adapter issueClientAdapter) CreateIssue(
	ctx context.Context,
	request client.IssueCreateRequest,
) (client.IssueSummary, error) {
	return client.CreateIssue(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter issueClientAdapter) CloseIssue(ctx context.Context, issueID string) (client.IssueSummary, error) {
	return client.CloseIssue(ctx, adapter.graphqlClient, adapter.target, issueID)
}

func (adapter issueClientAdapter) GetIssueTemplateContent(
	ctx context.Context,
	templateID string,
) (client.IssueTemplateContent, error) {
	return client.GetIssueTemplateContent(ctx, adapter.graphqlClient, templateID)
}

func (adapter issueClientAdapter) UpdateIssue(
	ctx context.Context,
	request client.IssueUpdateRequest,
) (client.IssueSummary, error) {
	return client.UpdateIssue(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter issueClientAdapter) CommentOnIssue(
	ctx context.Context,
	request client.IssueCommentRequest,
) (client.IssueCommentResult, error) {
	return client.CommentOnIssue(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter issueClientAdapter) StartIssue(ctx context.Context, issueID string) (client.IssueSummary, error) {
	return client.StartIssue(ctx, adapter.graphqlClient, adapter.target, issueID)
}

func (adapter issueClientAdapter) LinkIssueAttachment(
	ctx context.Context,
	request client.AttachmentLinkRequest,
) (client.AttachmentSummary, error) {
	return client.LinkIssueAttachment(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter issueClientAdapter) CreateIssueRelation(
	ctx context.Context,
	request client.IssueRelationCreateRequest,
) (client.IssueRelationSummary, error) {
	return client.CreateIssueRelation(ctx, adapter.graphqlClient, adapter.target, request)
}

func (adapter issueClientAdapter) DeleteIssueRelation(ctx context.Context, relationID string) (string, error) {
	return client.DeleteIssueRelation(ctx, adapter.graphqlClient, adapter.target, relationID)
}

func (adapter issueClientAdapter) ResolveTarget(ctx context.Context) (client.ResolvedTarget, error) {
	return client.ResolveTarget(ctx, adapter.graphqlClient, adapter.target)
}

func (adapter issueClientAdapter) ListIssues(ctx context.Context, limit int) (client.IssueList, error) {
	return client.ListIssues(ctx, adapter.graphqlClient, limit)
}

func (adapter issueClientAdapter) ListIssuesByTeam(
	ctx context.Context,
	teamID string,
	limit int,
	filters client.IssueListFilters,
) (client.IssueList, error) {
	return client.ListIssuesByTeam(ctx, adapter.graphqlClient, teamID, limit, filters)
}

func (adapter issueClientAdapter) GetIssueByID(ctx context.Context, issueID string) (client.IssueSummary, error) {
	return client.GetIssueByID(ctx, adapter.graphqlClient, issueID)
}

func (adapter issueClientAdapter) GetIssueDependencies(
	ctx context.Context,
	issueID string,
	limit int,
) (client.IssueDependencyGraph, error) {
	return client.GetIssueDependencies(ctx, adapter.graphqlClient, issueID, limit)
}
