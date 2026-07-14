package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/KyaniteHQ/linctl/internal/client/internal/gql"
)

// AuditEntryTypeSummary is the compact audit entry type catalog model.
type AuditEntryTypeSummary struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// AuditEntryTypeList is the audit entry type catalog.
type AuditEntryTypeList struct {
	AuditEntryTypes []AuditEntryTypeSummary `json:"audit_entry_types"`
}

// ListAuditEntryTypes returns the audit entry type catalog.
func ListAuditEntryTypes(ctx context.Context, graphqlClient graphql.Client) (AuditEntryTypeList, error) {
	result, err := gql.XAuditEntryTypes(ctx, graphqlClient)
	if err != nil {
		return AuditEntryTypeList{}, fmt.Errorf("list audit entry types: %w", err)
	}

	types := mapNodes(result.AuditEntryTypes, func(
		entryType gql.XAuditEntryTypesAuditEntryTypesAuditEntryType,
	) AuditEntryTypeSummary {
		return AuditEntryTypeSummary(entryType)
	})

	return AuditEntryTypeList{AuditEntryTypes: types}, nil
}
