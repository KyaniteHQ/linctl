package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"

	"github.com/KyaniteHQ/linctl/internal/client"
)

type initiativeLabelCommandFakeClient struct{}

func (initiativeLabelCommandFakeClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	var payload string
	switch request.OpName {
	case "initiativeLabels":
		payload = `{"initiativeLabels":{"nodes":[` + initiativeLabelCommandJSON() + `],"pageInfo":{"hasNextPage":false,"endCursor":null}}}`
	case "initiativeLabel":
		payload = `{"initiativeLabel":` + initiativeLabelCommandJSON() + `}`
	default:
		return errors.New("missing fake response for " + request.OpName)
	}

	return json.Unmarshal([]byte(`{"data":`+payload+`}`), response)
}

func initiativeLabelCommandJSON() string {
	return commandInitiativeLabelJSON("Strategy")
}

func Test_InitiativeLabelCommandFlows_list_get_and_project_fields(t *testing.T) {
	restoreRuntime := useCommandRuntime(t, initiativeLabelCommandFakeClient{})
	defer restoreRuntime()

	tests := []struct {
		name     string
		args     []string
		contains string
	}{
		{name: "list", args: []string{"initiative-label", "list", "--limit", "1"}, contains: "initiative-label-id Strategy #5e6ad2"},
		{name: "get", args: []string{"initiative-label", "get", "initiative-label-id"}, contains: "initiative-label-id Strategy #5e6ad2"},
		{name: "fields", args: []string{"--json", "--fields", "id", "initiative-label", "list", "--limit", "1"}, contains: `"id": "initiative-label-id"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := NewRootCommand(context.Background(), BuildInfo{})
			output := bytes.Buffer{}
			command.SetOut(&output)
			command.SetErr(&output)
			command.SetArgs(test.args)

			require.NoError(t, command.Execute())
			require.Contains(t, output.String(), test.contains)
		})
	}
}

func Test_WriteInitiativeLabel_covers_human_and_machine_modes(t *testing.T) {
	label := client.InitiativeLabelSummary{
		ID:         "initiative-label-id",
		Name:       "Strategy",
		Color:      "#5e6ad2",
		ParentName: "Themes",
	}
	command := &cobra.Command{}
	output := bytes.Buffer{}
	command.SetOut(&output)

	require.NoError(t, writeInitiativeLabel(command, &rootOptions{format: "minimal"}, label))
	require.NoError(t, writeInitiativeLabel(command, &rootOptions{format: "full"}, label))
	require.NoError(t, writeInitiativeLabel(command, &rootOptions{idOnly: true}, label))
	require.NoError(t, writeInitiativeLabel(command, &rootOptions{quiet: true}, label))
	require.Error(t, writeInitiativeLabel(command, &rootOptions{format: "wide"}, label))

	require.Contains(t, output.String(), "initiative-label-id")
	require.Contains(t, output.String(), "parent=Themes")
}
