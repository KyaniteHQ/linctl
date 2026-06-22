package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func templateContentResponseJSON(templateData string) string {
	return `{"template":{"id":"template-id","name":"Bug report","templateData":` + templateData + `}}`
}

func Test_GetIssueTemplateContent_extracts_title_and_description_from_object(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"templateContent": templateContentResponseJSON(`{"title":"Bug report","description":"## Steps\n\n1. "}`),
	})

	content, err := GetIssueTemplateContent(context.Background(), graphqlClient, "template-id")

	require.NoError(t, err)
	require.Equal(t, "Bug report", content.Title)
	require.Contains(t, content.Description, "Steps")
}

func Test_GetIssueTemplateContent_decodes_json_encoded_string_template_data(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"templateContent": templateContentResponseJSON(`"{\"title\":\"Encoded\",\"description\":\"body\"}"`),
	})

	content, err := GetIssueTemplateContent(context.Background(), graphqlClient, "template-id")

	require.NoError(t, err)
	require.Equal(t, "Encoded", content.Title)
	require.Equal(t, "body", content.Description)
}

func Test_GetIssueTemplateContent_returns_empty_when_template_data_absent(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"templateContent": `{"template":{"id":"template-id","name":"Bug report"}}`,
	})

	content, err := GetIssueTemplateContent(context.Background(), graphqlClient, "template-id")

	require.NoError(t, err)
	require.Empty(t, content.Title)
	require.Empty(t, content.Description)
}

func Test_GetIssueTemplateContent_returns_empty_for_empty_encoded_string(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"templateContent": templateContentResponseJSON(`""`),
	})

	content, err := GetIssueTemplateContent(context.Background(), graphqlClient, "template-id")

	require.NoError(t, err)
	require.Empty(t, content.Title)
}

func Test_GetIssueTemplateContent_ignores_non_string_fields(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"templateContent": templateContentResponseJSON(`{"title":123,"description":"ok"}`),
	})

	content, err := GetIssueTemplateContent(context.Background(), graphqlClient, "template-id")

	require.NoError(t, err)
	require.Empty(t, content.Title)
	require.Equal(t, "ok", content.Description)
}

func Test_GetIssueTemplateContent_requires_template_id(t *testing.T) {
	_, err := GetIssueTemplateContent(context.Background(), fakeGraphQLClient(map[string]string{}), "")

	require.ErrorIs(t, err, ErrWriteInvalid)
}

func Test_GetIssueTemplateContent_wraps_read_error(t *testing.T) {
	_, err := GetIssueTemplateContent(context.Background(), fakeGraphQLClient(map[string]string{}), "template-id")

	require.Error(t, err)
	require.Contains(t, err.Error(), "get template content")
}

func Test_GetIssueTemplateContent_errors_on_non_object_template_data(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"templateContent": templateContentResponseJSON(`123`),
	})

	_, err := GetIssueTemplateContent(context.Background(), graphqlClient, "template-id")

	require.Error(t, err)
	require.Contains(t, err.Error(), "decode template")
}

func Test_GetIssueTemplateContent_errors_on_unparsable_encoded_string(t *testing.T) {
	graphqlClient := fakeGraphQLClient(map[string]string{
		"templateContent": templateContentResponseJSON(`"not json"`),
	})

	_, err := GetIssueTemplateContent(context.Background(), graphqlClient, "template-id")

	require.Error(t, err)
	require.Contains(t, err.Error(), "decode template")
}
