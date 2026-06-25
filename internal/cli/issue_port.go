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
