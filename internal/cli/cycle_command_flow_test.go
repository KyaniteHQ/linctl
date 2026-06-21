package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

func Test_CycleCommandFlows_list_cycles(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"cycle", "list", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "cycle-id Cycle 12 [active]")
}

func Test_CycleCommandFlows_list_cycles_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "cycle", "list", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"cycles": [`)
	require.Contains(t, output.String(), `"id": "cycle-id"`)
}

func Test_CycleCommandFlows_get_cycle(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"cycle", "get", "cycle-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "cycle-id Named cycle [active]")
}

func Test_CycleCommandFlows_get_cycle_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "cycle", "get", "cycle-id"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"id": "cycle-id"`)
	require.Contains(t, output.String(), `"name": "Named cycle"`)
}

func Test_CycleCommandFlows_write_cycles(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{
			name: "create",
			args: []string{
				"cycle",
				"create",
				"--starts-at",
				"2026-07-01T00:00:00Z",
				"--ends-at",
				"2026-07-15T00:00:00Z",
				"--name",
				"Created cycle",
			},
			contains: "cycle-id Created cycle [active]",
		},
		{
			name:     "update",
			args:     []string{"cycle", "update", "cycle-id", "--name", "Updated cycle"},
			contains: "cycle-id Updated cycle [active]",
		},
		{
			name:     "archive",
			args:     []string{"cycle", "archive", "cycle-id"},
			contains: "cycle-id Archived cycle [active]",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := bytes.Buffer{}
			restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&output)
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Contains(t, output.String(), test.contains)
		})
	}
}

func Test_CycleCommandFlows_get_current_sprint(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"sprint", "current"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "cycle-id Cycle 12 [active]")
}

func Test_CycleCommandFlows_get_current_sprint_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "sprint", "current"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"id": "cycle-id"`)
	require.Contains(t, output.String(), `"status": "active"`)
}

func Test_CycleCommandFlows_report_current_sprint_runtime_error(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"sprint", "current"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime failed")
}

func Test_CycleCommandFlows_report_current_sprint_edges(t *testing.T) {
	tests := []struct {
		name        string
		fake        cycleCommandFlowFakeClient
		wantMessage string
	}{
		{
			name:        "resolve target",
			fake:        cycleCommandFlowFakeClient{failOperation: "Viewer"},
			wantMessage: "viewer failed",
		},
		{
			name:        "no active cycle",
			fake:        cycleCommandFlowFakeClient{emptyCycles: true},
			wantMessage: "current sprint: no active Cycle",
		},
		{
			name:        "list operation",
			fake:        cycleCommandFlowFakeClient{failOperation: "cycles"},
			wantMessage: "current sprint: list cycles",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs([]string{"sprint", "current"})

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.wantMessage)
		})
	}
}

func Test_CycleCommandFlows_report_current_sprint_writer_error(t *testing.T) {
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"sprint", "current"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func Test_CycleCommandFlows_report_sprint(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"sprint", "report", "cycle-id", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "cycle-id Current sprint [active]")
	require.Contains(t, output.String(), "LIT-1 Ship report [Started]")
}

func Test_CycleCommandFlows_report_sprint_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "sprint", "report", "cycle-id", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"cycle": {`)
	require.Contains(t, output.String(), `"identifier": "LIT-1"`)
}

func Test_CycleCommandFlows_report_sprint_edges(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		fake        cycleCommandFlowFakeClient
		wantMessage string
	}{
		{
			name:        "runtime",
			args:        []string{"sprint", "report", "cycle-id"},
			wantMessage: "runtime failed",
		},
		{
			name:        "operation",
			args:        []string{"sprint", "report", "cycle-id"},
			fake:        cycleCommandFlowFakeClient{failOperation: "CycleReport"},
			wantMessage: "cyclereport failed",
		},
		{
			name:        "empty",
			args:        []string{"--fail-on-empty", "sprint", "report", "cycle-id"},
			fake:        cycleCommandFlowFakeClient{emptyReport: true},
			wantMessage: "empty result",
		},
		{
			name:        "sort",
			args:        []string{"--sort", "missing", "sprint", "report", "cycle-id"},
			wantMessage: `sort field "missing" is not present`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var restore func()
			if test.name == "runtime" {
				original := buildCommandRuntime
				buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, errors.New("runtime failed")
				}
				restore = func() {
					buildCommandRuntime = original
				}
			} else {
				restore = useCommandRuntime(t, test.fake)
			}
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.wantMessage)
		})
	}
}

func Test_CycleCommandFlows_report_sprint_writer_error(t *testing.T) {
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"sprint", "report", "cycle-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func Test_CycleCommandFlows_report_sprint_issue_writer_error(t *testing.T) {
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&cycleCommandFailingSecondWriter{})
	command.SetArgs([]string{"sprint", "report", "cycle-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "second write failed")
}

func Test_CycleCommandFlows_list_cycle_issues(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"cycle", "issues", "cycle-id", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "LIT-1 Cycle issue [Started]")
}

func Test_CycleCommandFlows_list_cycle_issues_json(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"--json", "cycle", "issues", "cycle-id", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), `"cycle": {`)
	require.Contains(t, output.String(), `"identifier": "LIT-1"`)
}

func Test_CycleCommandFlows_list_cycle_uncompleted_issues(t *testing.T) {
	output := bytes.Buffer{}
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(&output)
	command.SetArgs([]string{"cycle", "uncompleted-issues", "cycle-id", "--limit", "1"})

	err := command.ExecuteContext(context.Background())

	require.NoError(t, err)
	require.Contains(t, output.String(), "LIT-2 Carry issue [Todo]")
}

func Test_CycleCommandFlows_list_cycle_issue_edges(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		fake        cycleCommandFlowFakeClient
		wantMessage string
	}{
		{
			name:        "runtime",
			args:        []string{"cycle", "issues", "cycle-id"},
			wantMessage: "runtime failed",
		},
		{
			name:        "operation",
			args:        []string{"cycle", "issues", "cycle-id"},
			fake:        cycleCommandFlowFakeClient{failOperation: "cycle_issues"},
			wantMessage: "cycle_issues failed",
		},
		{
			name:        "uncompleted operation",
			args:        []string{"cycle", "uncompleted-issues", "cycle-id"},
			fake:        cycleCommandFlowFakeClient{failOperation: "cycle_uncompletedIssuesUponClose"},
			wantMessage: "cycle_uncompletedissuesuponclose failed",
		},
		{
			name:        "empty",
			args:        []string{"--fail-on-empty", "cycle", "issues", "cycle-id"},
			fake:        cycleCommandFlowFakeClient{emptyReport: true},
			wantMessage: "empty result",
		},
		{
			name:        "sort",
			args:        []string{"--sort", "missing", "cycle", "issues", "cycle-id"},
			wantMessage: `sort field "missing" is not present`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var restore func()
			if test.name == "runtime" {
				original := buildCommandRuntime
				buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, errors.New("runtime failed")
				}
				restore = func() {
					buildCommandRuntime = original
				}
			} else {
				restore = useCommandRuntime(t, test.fake)
			}
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.wantMessage)
		})
	}
}

func Test_CycleCommandFlows_list_cycle_issue_writer_error(t *testing.T) {
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"cycle", "issues", "cycle-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func Test_CycleCommandFlows_report_cycle_get_runtime_error(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"cycle", "get", "cycle-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime failed")
}

func Test_CycleCommandFlows_report_cycle_get_operation_error(t *testing.T) {
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{failOperation: "cycle"})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"cycle", "get", "cycle-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "cycle failed")
}

func Test_CycleCommandFlows_report_cycle_get_writer_error(t *testing.T) {
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"cycle", "get", "cycle-id"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func Test_CycleCommandFlows_report_cycle_write_edges(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		fake        cycleCommandFlowFakeClient
		wantMessage string
	}{
		{
			name:        "create runtime",
			args:        []string{"cycle", "create", "--starts-at", "2026-07-01T00:00:00Z", "--ends-at", "2026-07-15T00:00:00Z"},
			wantMessage: "runtime failed",
		},
		{
			name:        "create operation",
			args:        []string{"cycle", "create", "--starts-at", "2026-07-01T00:00:00Z", "--ends-at", "2026-07-15T00:00:00Z"},
			fake:        cycleCommandFlowFakeClient{failOperation: "CycleCreate"},
			wantMessage: "cyclecreate failed",
		},
		{
			name:        "update operation",
			args:        []string{"cycle", "update", "cycle-id", "--name", "Updated cycle"},
			fake:        cycleCommandFlowFakeClient{failOperation: "CycleUpdate"},
			wantMessage: "cycleupdate failed",
		},
		{
			name:        "archive operation",
			args:        []string{"cycle", "archive", "cycle-id"},
			fake:        cycleCommandFlowFakeClient{failOperation: "CycleArchive"},
			wantMessage: "cyclearchive failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.wantMessage == "runtime failed" {
				original := buildCommandRuntime
				buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
					return commandRuntime{}, errors.New("runtime failed")
				}
				defer func() {
					buildCommandRuntime = original
				}()
			} else {
				restore := useCommandRuntime(t, test.fake)
				defer restore()
			}
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.wantMessage)
		})
	}
}

func Test_CycleCommandFlows_report_cycle_write_writer_errors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "create",
			args: []string{"cycle", "create", "--starts-at", "2026-07-01T00:00:00Z", "--ends-at", "2026-07-15T00:00:00Z"},
		},
		{name: "update", args: []string{"cycle", "update", "cycle-id", "--name", "Updated cycle"}},
		{name: "archive", args: []string{"cycle", "archive", "cycle-id"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(commandFailingWriter{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), "write failed")
		})
	}
}

func Test_CycleCommandFlows_report_cycle_list_edges(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		fake        cycleCommandFlowFakeClient
		wantMessage string
	}{
		{
			name:        "empty",
			args:        []string{"--fail-on-empty", "cycle", "list"},
			fake:        cycleCommandFlowFakeClient{emptyCycles: true},
			wantMessage: "empty result",
		},
		{
			name:        "sort",
			args:        []string{"--sort", "missing", "cycle", "list"},
			wantMessage: `sort field "missing" is not present`,
		},
		{
			name:        "resolve target",
			args:        []string{"cycle", "list"},
			fake:        cycleCommandFlowFakeClient{failOperation: "Viewer"},
			wantMessage: "viewer failed",
		},
		{
			name:        "list operation",
			args:        []string{"cycle", "list"},
			fake:        cycleCommandFlowFakeClient{failOperation: "cycles"},
			wantMessage: "cycles failed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, test.fake)
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.Error(t, err)
			require.Contains(t, err.Error(), test.wantMessage)
		})
	}
}

func Test_CycleCommandFlows_report_cycle_list_runtime_error(t *testing.T) {
	original := buildCommandRuntime
	buildCommandRuntime = func(_ context.Context, _ *rootOptions) (commandRuntime, error) {
		return commandRuntime{}, errors.New("runtime failed")
	}
	defer func() {
		buildCommandRuntime = original
	}()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetArgs([]string{"cycle", "list"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "runtime failed")
}

func Test_CycleCommandFlows_report_cycle_list_writer_error(t *testing.T) {
	restore := useCommandRuntime(t, cycleCommandFlowFakeClient{})
	defer restore()
	command := NewRootCommand(context.Background(), BuildInfo{})
	command.SetOut(commandFailingWriter{})
	command.SetArgs([]string{"cycle", "list"})

	err := command.ExecuteContext(context.Background())

	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

type cycleCommandFlowFakeClient struct {
	emptyCycles   bool
	emptyReport   bool
	failOperation string
}

func (client cycleCommandFlowFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if client.failOperation == request.OpName {
		return errors.New(strings.ToLower(request.OpName) + " failed")
	}
	payload, ok := cycleCommandFlowPayload(request.OpName, client.emptyCycles, client.emptyReport)
	if !ok {
		return errors.New("missing fake response for " + request.OpName)
	}

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

var _ graphql.Client = cycleCommandFlowFakeClient{}

type cycleCommandFailingSecondWriter struct {
	writes int
}

func (writer *cycleCommandFailingSecondWriter) Write(payload []byte) (int, error) {
	writer.writes++
	if writer.writes == 1 {
		return len(payload), nil
	}
	return 0, errors.New("second write failed")
}

func cycleCommandFlowPayload(operation string, emptyCycles bool, emptyReport bool) (string, bool) {
	if payload, ok := cycleCommandFlowCyclePayload(operation, emptyCycles, emptyReport); ok {
		return payload, true
	}

	switch operation {
	case "Viewer":
		return `{"viewer":{"id":"user-id","name":"Omer","displayName":"Omer","email":"omer@example.com","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}}`, true
	case "Teams":
		return `{"teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "TargetProject":
		return `{"project":{"id":"project-id","name":"Pinned project","teams":{"nodes":[{"id":"team-id","key":"LIT","name":"linctl","organization":{"id":"org-id","name":"Kyanite","urlKey":"kyanite"}}]}}}`, true
	default:
		return "", false
	}
}

func cycleCommandFlowCyclePayload(operation string, emptyCycles bool, emptyReport bool) (string, bool) {
	switch operation {
	case "cycles":
		if emptyCycles {
			return `{"cycles":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
		}
		return `{"cycles":{"nodes":[{"id":"cycle-id","number":12,"name":null,"description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`, true
	case "cycle":
		return `{"cycle":{"id":"cycle-id","number":12,"name":"Named cycle","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"}}}`, true
	case "CycleCreate":
		return `{"cycleCreate":{"success":true,"cycle":` + cycleCommandCycleJSON("Created cycle") + `}}`, true
	case "CycleUpdate":
		return `{"cycleUpdate":{"success":true,"cycle":` + cycleCommandCycleJSON("Updated cycle") + `}}`, true
	case "CycleArchive":
		return `{"cycleArchive":{"success":true,"entity":` + cycleCommandCycleJSON("Archived cycle") + `}}`, true
	case "CycleReport":
		if emptyReport {
			return `{"cycle":{"id":"cycle-id","number":12,"name":"Current sprint","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"cycle":{"id":"cycle-id","number":12,"name":"Current sprint","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issues":{"nodes":[{"id":"issue-id","identifier":"LIT-1","title":"Ship report","branchName":"omer/ship-report","url":"https://linear.app/issue/LIT-1","priority":1,"priorityLabel":"Urgent","team":{"id":"team-id","key":"LIT","name":"linctl"},"state":{"id":"started","name":"Started","type":"started"},"assignee":{"id":"user-id","name":"omer","displayName":"Omer"},"project":{"id":"project-id","name":"Pinned project"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "cycle_issues":
		if emptyReport {
			return `{"cycle":{"id":"cycle-id","number":12,"name":"Current cycle","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issues":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"cycle":{"id":"cycle-id","number":12,"name":"Current cycle","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2099-01-01T00:00:00Z","completedAt":null,"progress":0.25,"team":{"id":"team-id","key":"LIT","name":"linctl"},"issues":{"nodes":[{"id":"issue-id","identifier":"LIT-1","title":"Cycle issue","branchName":"omer/cycle-issue","url":"https://linear.app/issue/LIT-1","priority":1,"priorityLabel":"Urgent","team":{"id":"team-id","key":"LIT","name":"linctl"},"state":{"id":"started","name":"Started","type":"started"},"assignee":{"id":"user-id","name":"omer","displayName":"Omer"},"project":{"id":"project-id","name":"Pinned project"}}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	case "cycle_uncompletedIssuesUponClose":
		if emptyReport {
			return `{"cycle":{"id":"cycle-id","number":12,"name":"Closed cycle","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2026-01-15T00:00:00Z","completedAt":"2026-01-15T00:00:00Z","progress":0.75,"team":{"id":"team-id","key":"LIT","name":"linctl"},"uncompletedIssuesUponClose":{"nodes":[],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
		}
		return `{"cycle":{"id":"cycle-id","number":12,"name":"Closed cycle","description":"cycle body","startsAt":"2026-01-01T00:00:00Z","endsAt":"2026-01-15T00:00:00Z","completedAt":"2026-01-15T00:00:00Z","progress":0.75,"team":{"id":"team-id","key":"LIT","name":"linctl"},"uncompletedIssuesUponClose":{"nodes":[{"id":"issue-id-2","identifier":"LIT-2","title":"Carry issue","branchName":"omer/carry-issue","url":"https://linear.app/issue/LIT-2","priority":2,"priorityLabel":"High","team":{"id":"team-id","key":"LIT","name":"linctl"},"state":{"id":"todo","name":"Todo","type":"unstarted"},"assignee":null,"project":null}],"pageInfo":{"hasNextPage":false,"endCursor":null}}}}`, true
	default:
		return "", false
	}
}

func cycleCommandCycleJSON(name string) string {
	return `{
		"id":"cycle-id",
		"number":12,
		"name":"` + name + `",
		"description":"cycle body",
		"startsAt":"2026-01-01T00:00:00Z",
		"endsAt":"2099-01-01T00:00:00Z",
		"completedAt":null,
		"progress":0.25,
		"team":{"id":"team-id","key":"LIT","name":"linctl"}
	}`
}
