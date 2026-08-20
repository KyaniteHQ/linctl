package client

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/Khan/genqlient/graphql"
)

func workflowStatesByTeamJSON(nodes string) string {
	return `{"workflowStates":{"nodes":[` + nodes + `],"pageInfo":{"hasNextPage":false}}}`
}

func multipleStartedStatesJSON() string {
	return workflowStatesByTeamJSON(`
		{"id":"in-review-state","name":"In Review","type":"started","position":2},
		{"id":"started-state","name":"Started","type":"started","position":1},
		{"id":"in-progress-state","name":"In Progress","type":"started","position":0}
	`)
}

type issueAfterWriteFake struct {
	inner graphql.Client
	after string
	mu    sync.Mutex
	wrote bool
}

func withIssueAfterWrite(inner graphql.Client, after issueFixture) *issueAfterWriteFake {
	return withIssueAfterWriteJSON(inner, issueJSON(after))
}

func withIssueAfterWriteJSON(inner graphql.Client, issueBody string) *issueAfterWriteFake {
	return &issueAfterWriteFake{inner: inner, after: `{"issue":` + issueBody + `}`}
}

func (fake *issueAfterWriteFake) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	fake.mu.Lock()
	wrote := fake.wrote
	fake.mu.Unlock()
	if request.OpName == "issue" && wrote && fake.after != "" {
		return json.Unmarshal([]byte(`{"data":`+fake.after+`}`), response)
	}
	err := fake.inner.MakeRequest(ctx, request, response)
	switch request.OpName {
	case "IssueUpdate", "IssueCreate", "IssueClose", "IssueRelationCreate":
		fake.mu.Lock()
		fake.wrote = true
		fake.mu.Unlock()
	}

	return err
}
