package client

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
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
	result, err := auditEntryTypes(ctx, graphqlClient)
	if err != nil {
		return AuditEntryTypeList{}, fmt.Errorf("list audit entry types: %w", err)
	}

	types := make([]AuditEntryTypeSummary, 0, len(result.AuditEntryTypes))
	for _, entryType := range result.AuditEntryTypes {
		types = append(types, AuditEntryTypeSummary(entryType))
	}

	return AuditEntryTypeList{AuditEntryTypes: types}, nil
}
