package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// IssueTemplateContent is the issue title and description pre-filled by a Linear template.
type IssueTemplateContent struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// GetIssueTemplateContent reads a Linear template and extracts the issue title and
// description from its templateData. It is a free read: the template is fetched
// without a pinned target so a draft can be assembled before any guarded write.
func GetIssueTemplateContent(
	ctx context.Context,
	graphqlClient graphql.Client,
	templateID string,
) (IssueTemplateContent, error) {
	if templateID == "" {
		return IssueTemplateContent{}, requiredFieldError("template id")
	}
	result, err := gql.XTemplateContent(ctx, graphqlClient, templateID)
	if err != nil {
		return IssueTemplateContent{}, fmt.Errorf("get template content %s: %w", templateID, err)
	}
	fields, err := decodeTemplateData(result.Template.TemplateData)
	if err != nil {
		return IssueTemplateContent{}, fmt.Errorf("decode template %s data: %w", templateID, err)
	}

	return IssueTemplateContent{
		Title:       templateStringField(fields, "title"),
		Description: templateStringField(fields, "description"),
	}, nil
}

// decodeTemplateData parses templateData, which Linear returns either as a JSON
// object or as a JSON-encoded string wrapping that object. Absent data and an
// encoded empty string are empty field maps.
func decodeTemplateData(raw json.RawMessage) (map[string]json.RawMessage, error) {
	if len(raw) == 0 {
		return map[string]json.RawMessage{}, nil
	}
	value, err := decodeTemplateEnvelope(raw)
	if err != nil {
		return nil, err
	}
	if encoded, ok := value.(string); ok && encoded == "" {
		return map[string]json.RawMessage{}, nil
	}

	return templateDataFieldMap(value)
}

func templateDataFieldMap(value any) (map[string]json.RawMessage, error) {
	fields := map[string]json.RawMessage{}
	if err := json.Unmarshal(mustJSON(value), &fields); err != nil {
		return nil, err
	}

	return fields, nil
}

// templateStringField returns the named field as a string, or "" when it is
// absent or not a JSON string.
func templateStringField(fields map[string]json.RawMessage, key string) string {
	raw, ok := fields[key]
	if !ok {
		return ""
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return ""
	}

	return value
}
