package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// SLAConfigurationSummary is the compact SLA rule model used by read-only commands.
type SLAConfigurationSummary struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Conditions json.RawMessage `json:"conditions"`
	SLA        float64         `json:"sla,omitempty"`
	SLAType    string          `json:"sla_type,omitempty"`
	RemovesSLA bool            `json:"removes_sla"`
}

// SLAConfigurationList is the active SLA configuration set for one team.
type SLAConfigurationList struct {
	TeamIDOrKey       string                    `json:"team_id_or_key"`
	SLAConfigurations []SLAConfigurationSummary `json:"sla_configurations"`
}

// ListSLAConfigurations returns active SLA rules that can apply to one team.
func ListSLAConfigurations(
	ctx context.Context,
	graphqlClient graphql.Client,
	teamIDOrKey string,
) (SLAConfigurationList, error) {
	result, err := slaConfigurations(ctx, graphqlClient, teamIDOrKey)
	if err != nil {
		return SLAConfigurationList{}, fmt.Errorf("list SLA configurations %s: %w", teamIDOrKey, err)
	}

	configurations := make([]SLAConfigurationSummary, 0, len(result.SlaConfigurations))
	for _, configuration := range result.SlaConfigurations {
		configurations = append(configurations, slaConfigurationSummary(configuration.SlaConfigurationSummaryFields))
	}

	return SLAConfigurationList{
		TeamIDOrKey:       teamIDOrKey,
		SLAConfigurations: configurations,
	}, nil
}

func slaConfigurationSummary(fields SlaConfigurationSummaryFields) SLAConfigurationSummary {
	summary := SLAConfigurationSummary{
		ID:         fields.Id,
		Name:       fields.Name,
		Conditions: fields.Conditions,
		RemovesSLA: fields.RemovesSla,
	}
	if fields.Sla != nil {
		summary.SLA = *fields.Sla
	}
	if fields.SlaType != nil {
		summary.SLAType = string(*fields.SlaType)
	}

	return summary
}
