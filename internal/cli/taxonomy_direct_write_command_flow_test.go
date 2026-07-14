package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/Khan/genqlient/graphql"
	"github.com/stretchr/testify/require"
)

type expectedTaxonomyWriteVariable struct {
	path   []string
	value  any
	absent bool
}

type taxonomyWriteCaptureClient struct {
	operation string
	variables []expectedTaxonomyWriteVariable
	delegate  graphql.Client
	calls     int
}

func (client *taxonomyWriteCaptureClient) MakeRequest(
	ctx context.Context,
	request *graphql.Request,
	response *graphql.Response,
) error {
	if request.OpName != client.operation {
		return client.delegate.MakeRequest(ctx, request, response)
	}

	client.calls++
	for _, expected := range client.variables {
		if err := validateTaxonomyWriteVariable(request, expected); err != nil {
			return err
		}
	}

	return client.delegate.MakeRequest(ctx, request, response)
}

func validateTaxonomyWriteVariable(
	request *graphql.Request,
	expected expectedTaxonomyWriteVariable,
) error {
	actual, found, err := taxonomyWriteVariable(request, expected.path...)
	path := strings.Join(expected.path, ".")
	switch {
	case err != nil:
		return err
	case expected.absent && found:
		return fmt.Errorf("%s is present, want omitted", path)
	case expected.absent:
		return nil
	case !found:
		return errors.New("request variable missing " + path)
	case !reflect.DeepEqual(expected.value, actual):
		return fmt.Errorf("%s = %#v, want %#v", path, actual, expected.value)
	default:
		return nil
	}
}

func taxonomyWriteVariable(request *graphql.Request, path ...string) (any, bool, error) {
	payload, err := json.Marshal(request.Variables)
	if err != nil {
		return nil, false, err
	}
	var variables map[string]any
	if err := json.Unmarshal(payload, &variables); err != nil {
		return nil, false, err
	}

	current := any(variables)
	for _, key := range path {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false, errors.New("request variable is not an object")
		}
		value, ok := object[key]
		if !ok {
			return nil, false, nil
		}
		current = value
	}

	return current, true, nil
}

func Test_TaxonomyDirectWriteCommandFlows_forward_exact_mutation_variables(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		operation string
		variables []expectedTaxonomyWriteVariable
		fake      commandFlowFakeClient
	}{
		{
			name: "organization-wide label create",
			args: []string{
				"label", "create", "--name", "Platform", "--color", "#123456",
				"--description", "Shared label", "--parent", "parent-label-id", "--org-wide",
			},
			operation: "IssueLabelCreate",
			variables: []expectedTaxonomyWriteVariable{
				{path: []string{"replaceTeamLabels"}, value: false},
				{path: []string{"input", "name"}, value: "Platform"},
				{path: []string{"input", "color"}, value: "#123456"},
				{path: []string{"input", "description"}, value: "Shared label"},
				{path: []string{"input", "parentId"}, value: "parent-label-id"},
				{path: []string{"input", "teamId"}, absent: true},
			},
			fake: commandFlowFakeClient{orgWideLabel: true},
		},
		{
			name: "label update",
			args: []string{
				"label", "update", "label-id", "--name", "Customer", "--color", "#654321",
				"--description", "Customer work",
			},
			operation: "IssueLabelUpdate",
			variables: []expectedTaxonomyWriteVariable{
				{path: []string{"id"}, value: "label-id"},
				{path: []string{"replaceTeamLabels"}, value: false},
				{path: []string{"input", "name"}, value: "Customer"},
				{path: []string{"input", "color"}, value: "#654321"},
				{path: []string{"input", "description"}, value: "Customer work"},
			},
		},
		{
			name:      "label retire",
			args:      []string{"label", "retire", "label-id"},
			operation: "IssueLabelRetire",
			variables: []expectedTaxonomyWriteVariable{{path: []string{"id"}, value: "label-id"}},
		},
		{
			name:      "label restore",
			args:      []string{"label", "restore", "label-id"},
			operation: "IssueLabelRestore",
			variables: []expectedTaxonomyWriteVariable{{path: []string{"id"}, value: "label-id"}},
		},
		{
			name: "project label create",
			args: []string{
				"project-label", "create", "--name", "Roadmap", "--color", "#abcdef",
				"--description", "Roadmap work", "--org-wide",
			},
			operation: "ProjectLabelCreate",
			variables: []expectedTaxonomyWriteVariable{
				{path: []string{"input", "name"}, value: "Roadmap"},
				{path: []string{"input", "color"}, value: "#abcdef"},
				{path: []string{"input", "description"}, value: "Roadmap work"},
			},
		},
		{
			name: "project label update",
			args: []string{
				"project-label", "update", "project-label-id", "--name", "Launch", "--color", "#fedcba",
				"--description", "Launch work", "--org-wide",
			},
			operation: "ProjectLabelUpdate",
			variables: []expectedTaxonomyWriteVariable{
				{path: []string{"id"}, value: "project-label-id"},
				{path: []string{"input", "name"}, value: "Launch"},
				{path: []string{"input", "color"}, value: "#fedcba"},
				{path: []string{"input", "description"}, value: "Launch work"},
			},
		},
		{
			name:      "project label retire",
			args:      []string{"project-label", "retire", "project-label-id", "--org-wide"},
			operation: "ProjectLabelRetire",
			variables: []expectedTaxonomyWriteVariable{{path: []string{"id"}, value: "project-label-id"}},
		},
		{
			name:      "project label restore",
			args:      []string{"project-label", "restore", "project-label-id", "--org-wide"},
			operation: "ProjectLabelRestore",
			variables: []expectedTaxonomyWriteVariable{{path: []string{"id"}, value: "project-label-id"}},
		},
		{
			name: "project milestone create",
			args: []string{
				"project-milestone", "create", "project-id", "--name", "Public beta",
				"--description", "Beta launch", "--target-date", "2026-08-01",
			},
			operation: "ProjectMilestoneCreate",
			variables: []expectedTaxonomyWriteVariable{
				{path: []string{"input", "projectId"}, value: "project-id"},
				{path: []string{"input", "name"}, value: "Public beta"},
				{path: []string{"input", "description"}, value: "Beta launch"},
				{path: []string{"input", "targetDate"}, value: "2026-08-01"},
			},
		},
		{
			name: "project milestone update",
			args: []string{
				"project-milestone", "update", "project-milestone-id", "--name", "General availability",
				"--description", "GA launch", "--target-date", "2026-08-15",
			},
			operation: "ProjectMilestoneUpdate",
			variables: []expectedTaxonomyWriteVariable{
				{path: []string{"id"}, value: "project-milestone-id"},
				{path: []string{"input", "name"}, value: "General availability"},
				{path: []string{"input", "description"}, value: "GA launch"},
				{path: []string{"input", "targetDate"}, value: "2026-08-15"},
			},
		},
		{
			name:      "project milestone delete",
			args:      []string{"project-milestone", "delete", "project-milestone-id"},
			operation: "ProjectMilestoneDelete",
			variables: []expectedTaxonomyWriteVariable{{path: []string{"id"}, value: "project-milestone-id"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fake := &taxonomyWriteCaptureClient{
				operation: test.operation,
				variables: test.variables,
				delegate:  test.fake,
			}
			restore := useCommandRuntime(t, fake)
			defer restore()

			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetOut(&bytes.Buffer{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.NoError(t, err)
			require.Equal(t, 1, fake.calls)
		})
	}
}

func Test_TaxonomyDirectWriteCommandFlows_propagate_mutation_errors(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		operation string
	}{
		{name: "label create", args: []string{"label", "create", "--name", "Created"}, operation: "IssueLabelCreate"},
		{name: "label update", args: []string{"label", "update", "label-id", "--name", "Updated"}, operation: "IssueLabelUpdate"},
		{name: "label retire", args: []string{"label", "retire", "label-id"}, operation: "IssueLabelRetire"},
		{name: "label restore", args: []string{"label", "restore", "label-id"}, operation: "IssueLabelRestore"},
		{name: "project label create", args: []string{"project-label", "create", "--name", "Created", "--org-wide"}, operation: "ProjectLabelCreate"},
		{name: "project label update", args: []string{"project-label", "update", "project-label-id", "--name", "Updated", "--org-wide"}, operation: "ProjectLabelUpdate"},
		{name: "project label retire", args: []string{"project-label", "retire", "project-label-id", "--org-wide"}, operation: "ProjectLabelRetire"},
		{name: "project label restore", args: []string{"project-label", "restore", "project-label-id", "--org-wide"}, operation: "ProjectLabelRestore"},
		{name: "project milestone create", args: []string{"project-milestone", "create", "project-id", "--name", "Created"}, operation: "ProjectMilestoneCreate"},
		{name: "project milestone update", args: []string{"project-milestone", "update", "project-milestone-id", "--name", "Updated"}, operation: "ProjectMilestoneUpdate"},
		{name: "project milestone delete", args: []string{"project-milestone", "delete", "project-milestone-id"}, operation: "ProjectMilestoneDelete"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			restore := useCommandRuntime(t, commandFlowFakeClient{failOperation: test.operation})
			defer restore()
			command := NewRootCommand(context.Background(), BuildInfo{})
			command.SetArgs(test.args)

			err := command.ExecuteContext(context.Background())

			require.ErrorContains(t, err, "operation failed")
		})
	}
}
