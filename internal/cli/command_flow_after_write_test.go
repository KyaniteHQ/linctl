package cli

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/Khan/genqlient/graphql"
)

type commandFlowAfterWriteClient struct {
	inner         graphql.Client
	afterIssue    string
	afterRelation string
	startStateID  string
}

func wrapCommandFlowAfterWrite(inner graphql.Client) graphql.Client {
	if _, ok := inner.(*commandFlowAfterWriteClient); ok {
		return inner
	}
	wrapper := &commandFlowAfterWriteClient{inner: inner}
	if fake, ok := inner.(commandFlowFakeClient); ok {
		wrapper.startStateID = fake.expectedStartStateID
	}

	return wrapper
}

func (client *commandFlowAfterWriteClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName == "issue" && client.afterIssue != "" {
		return json.Unmarshal([]byte(`{"data":`+client.afterIssue+`}`), response)
	}
	if request.OpName == "issueRelation" && client.afterRelation != "" {
		return json.Unmarshal([]byte(`{"data":`+client.afterRelation+`}`), response)
	}
	err := client.inner.MakeRequest(ctx, request, response)
	if err != nil {
		return err
	}
	client.record(request)

	return nil
}

func (client *commandFlowAfterWriteClient) record(request *graphql.Request) {
	switch request.OpName {
	case "IssueUpdate", "IssueCreate", "IssueClose":
		stateID, err := requestVariable[string](request, "input", "stateId")
		if err != nil || stateID == "" {
			return
		}
		name, stateType := commandFlowStateByID(stateID)
		title := "Updated issue"
		identifier := "LIT-1"
		if request.OpName == "IssueCreate" {
			title = "Created issue"
			identifier = "LIT-2"
		}
		if request.OpName == "IssueClose" {
			title = "Closed issue"
		}
		if request.OpName == "IssueUpdate" && client.startStateID != "" {
			title = "Started issue"
		}
		client.afterIssue = `{"issue":` + commandIssueJSON(identifier, title, stateID, name, stateType) + `}`
	case "IssueRelationCreate":
		client.afterRelation = `{"issueRelation":` + strings.Replace(
			commandIssueRelationJSON(),
			`"type":"blocks"`,
			`"type":"related"`,
			1,
		) + `}`
	}
}
