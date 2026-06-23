package client

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

type errorGraphQLClient struct {
	err error
}

func (client errorGraphQLClient) MakeRequest(
	_ context.Context,
	_ *graphql.Request,
	_ *graphql.Response,
) error {
	return client.err
}

type operationErrorFakeClient struct {
	responses map[string]string
	err       error
}

func (client operationErrorFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if client.responses[request.OpName] == "" {
		return client.err
	}

	return fakeGraphQLClient(client.responses).MakeRequest(ctx, request, response)
}

func (client issueWriteFakeClient) withError(err error) operationErrorFakeClient {
	return operationErrorFakeClient{
		responses: client.withTargetResponses(),
		err:       err,
	}
}

func (client projectWriteFakeClient) withError(err error) operationErrorFakeClient {
	return operationErrorFakeClient{
		responses: client.withTargetResponses(),
		err:       err,
	}
}
