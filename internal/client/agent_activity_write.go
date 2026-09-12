package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
	"github.com/KyaniteHQ/linctl/internal/config"
)

// AgentActivityCreateRequest describes a guarded AgentActivity create. The
// activity is emitted into an AgentSession by the agent that owns it, and its
// content follows Linear's typed activity payloads: thought, elicitation,
// response, and error carry a Body; action carries Action, Parameter, and an
// optional Result.
type AgentActivityCreateRequest struct {
	AgentSessionID string
	Type           string
	Body           string
	Action         string
	Parameter      string
	Result         string
	Signal         string
	Ephemeral      bool
}

// AgentActivityBodyTypes lists the activity content types that carry a body.
var AgentActivityBodyTypes = []string{"thought", "elicitation", "response", "error"}

// AgentActivityTypes lists every activity content type linctl can emit.
var AgentActivityTypes = append(append([]string{}, AgentActivityBodyTypes...), "action")

// AgentActivitySignals lists the AgentActivitySignal values Linear accepts.
var AgentActivitySignals = []string{"auth", "continue", "select", "stop"}

// CreateAgentActivity emits an AgentActivity into an AgentSession after resolving
// the session, requiring it to be attached to an issue, and comparing that issue
// against the pinned target. An AgentSession has no team of its own, so its
// issue is the whole scope: a session on an issue outside the pinned team or
// pinned project is a hard stop.
func CreateAgentActivity(
	ctx context.Context,
	graphqlClient graphql.Client,
	expected config.Target,
	request AgentActivityCreateRequest,
) (AgentActivitySummary, error) {
	content, err := agentActivityContent(request)
	if err != nil {
		return AgentActivitySummary{}, err
	}

	guard, err := newGuardedClient(ctx, graphqlClient, expected)
	if err != nil {
		return AgentActivitySummary{}, err
	}

	return guard.createAgentActivity(ctx, request, content)
}

func agentActivityContent(request AgentActivityCreateRequest) (json.RawMessage, error) {
	if request.AgentSessionID == "" {
		return nil, requiredFieldError("agent session id")
	}
	if request.Type == "" {
		return nil, requiredFieldError("type")
	}
	if !containsString(AgentActivityTypes, request.Type) {
		return nil, fmt.Errorf(
			"%w: type must be one of %s", ErrWriteInvalid, strings.Join(AgentActivityTypes, ", "),
		)
	}
	if request.Signal != "" && !containsString(AgentActivitySignals, request.Signal) {
		return nil, fmt.Errorf(
			"%w: signal must be one of %s", ErrWriteInvalid, strings.Join(AgentActivitySignals, ", "),
		)
	}

	if request.Type == "action" {
		return agentActionContent(request)
	}
	if request.Body == "" {
		return nil, requiredFieldError("body")
	}

	return mustJSON(map[string]string{"type": request.Type, "body": request.Body}), nil
}

func agentActionContent(request AgentActivityCreateRequest) (json.RawMessage, error) {
	if request.Action == "" {
		return nil, requiredFieldError("action")
	}
	if request.Parameter == "" {
		return nil, requiredFieldError("parameter")
	}
	content := map[string]string{"type": "action", "action": request.Action, "parameter": request.Parameter}
	if request.Result != "" {
		content["result"] = request.Result
	}

	return mustJSON(content), nil
}

func containsString(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}

	return false
}

func (guard *guardedClient) createAgentActivity(
	ctx context.Context,
	request AgentActivityCreateRequest,
	content json.RawMessage,
) (AgentActivitySummary, error) {
	session, err := GetAgentSessionByID(ctx, guard.graphqlClient, request.AgentSessionID)
	if err != nil {
		return AgentActivitySummary{}, err
	}
	if session.IssueID == "" {
		return AgentActivitySummary{}, fmt.Errorf(
			"%w: agent session %s is not attached to an issue; only issue sessions are guarded",
			ErrWriteInvalid, request.AgentSessionID,
		)
	}
	if _, err := guard.requireIssue(ctx, session.IssueID); err != nil {
		return AgentActivitySummary{}, err
	}

	input := LinearAgentActivityCreateInput{
		AgentSessionID: request.AgentSessionID,
		Content:        content,
		Signal:         optionalString(request.Signal),
	}
	if request.Ephemeral {
		input.Ephemeral = &request.Ephemeral
	}

	created, err := gql.AgentActivityCreate(ctx, guard.graphqlClient, input)
	if err != nil {
		return AgentActivitySummary{}, fmt.Errorf("create agent activity: %w", err)
	}
	if err := mutationSuccess(created.AgentActivityCreate.Success, "agentActivityCreate"); err != nil {
		return AgentActivitySummary{}, err
	}

	return agentActivitySummary(created.AgentActivityCreate.AgentActivity.AgentActivitySummaryFields), nil
}
