//go:build ignore

package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// statusOrder lists accounting statuses from most to least settled for stable output.
var statusOrder = []string{
	"implemented",
	"intentionally_excluded",
	"blocked_needs_design",
	"accepted_gap",
	"safe_candidate",
}

type accountedOperation struct {
	Name      string
	Kind      string
	Status    string
	Rationale string
}

func writeSDKOperationAudit(
	output *bytes.Buffer,
	upstreamDocumentsPath string,
	sdkOperations []sdkMethod,
	localOperations []localOperation,
) {
	localOperationNames := operationNameSet(localOperations)
	accounted := accountOperations(sdkOperations, localOperationNames)
	statusCounts := statusCountSet(accounted)
	byStatus := operationsByStatus(accounted)

	fmt.Fprintf(output, "# linctl SDK operation coverage audit\n\n")
	fmt.Fprintf(output, "Generated from `%s`.\n\n", upstreamDocumentsPath)

	fmt.Fprintf(output, "## Accounting summary\n\n")
	fmt.Fprintf(output, "- Official SDK operation total: %d\n", len(sdkOperations))
	fmt.Fprintf(output, "- Current linctl operation total: %d\n", len(localOperations))
	fmt.Fprintf(
		output,
		"- Accounted (every operation carries a documented status and rationale): %d (100%%)\n",
		len(accounted),
	)
	for _, status := range statusOrder {
		fmt.Fprintf(output, "  - %s: %d\n", status, statusCounts[status])
	}
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, "## What \"accounted\" means\n\n")
	fmt.Fprintf(
		output,
		"Every official SDK operation holds exactly one accounting status, so the SDK surface is "+
			"fully accounted even though it is deliberately not fully implemented.\n\n"+
			"- `implemented`: backed by a local GraphQL operation in `internal/client/operations`.\n"+
			"- `intentionally_excluded`: admin, auth, integration, and internal surfaces that sit "+
			"outside an agent-safe control surface and are not planned.\n"+
			"- `blocked_needs_design`: writes and state changes that stay closed until an explicit "+
			"target-pinned guard and a mismatch test exist.\n"+
			"- `accepted_gap` and `safe_candidate`: reads that may join a future slice but are "+
			"deferred under the control-surface safety model.\n\n",
	)

	fmt.Fprintf(output, "## Implementation order\n\n")
	fmt.Fprintf(
		output,
		"1. Read-only operations in existing CLI domains: issue, comment, project, "+
			"ProjectMilestone, Cycle, document, label, team, user, viewer, organization.\n",
	)
	fmt.Fprintf(
		output,
		"2. Read-only operations in adjacent execution domains: attachment, initiative, "+
			"roadmap, release, customer, custom view, notification.\n",
	)
	fmt.Fprintf(
		output,
		"3. Safe create/update operations only after the resource has an explicit "+
			"team-scoped or resource-scoped guard and a target-mismatch test.\n",
	)
	fmt.Fprintf(
		output,
		"4. Destructive, admin, integration, auth, security, and organization-wide "+
			"operations stay blocked until their guard model is documented.\n\n",
	)

	fmt.Fprintf(output, "## Current linctl operations\n\n")
	for _, operation := range localOperations {
		fmt.Fprintf(output, "- `%s %s` - `%s`\n", operation.Kind, operation.Name, operation.Path)
	}
	fmt.Fprintf(output, "\n")

	fmt.Fprintf(output, "## Operations by status\n\n")
	for _, status := range statusOrder {
		operations := byStatus[status]
		if len(operations) == 0 {
			continue
		}
		queryCount := accountedKindCount(operations, "query")
		mutationCount := accountedKindCount(operations, "mutation")
		fmt.Fprintf(
			output,
			"### %s (%d: %d queries, %d mutations)\n\n",
			status,
			len(operations),
			queryCount,
			mutationCount,
		)
		for _, operation := range operations {
			fmt.Fprintf(output, "- `%s %s` - %s\n", operation.Kind, operation.Name, operation.Rationale)
		}
		fmt.Fprintf(output, "\n")
	}
}

func accountOperations(sdkOperations []sdkMethod, localOperationNames map[string]bool) []accountedOperation {
	accounted := make([]accountedOperation, 0, len(sdkOperations))
	for _, operation := range sdkOperations {
		status, rationale := classifyOperation(operation, localOperationNames)
		accounted = append(accounted, accountedOperation{
			Name:      operation.Name,
			Kind:      operation.Kind,
			Status:    status,
			Rationale: rationale,
		})
	}
	return accounted
}

func classifyOperation(operation sdkMethod, localOperationNames map[string]bool) (string, string) {
	if localOperationImplemented(operation.Name, localOperationNames) {
		return "implemented", "backed by a local GraphQL operation in internal/client/operations"
	}
	return classifyLoose(operation.Name, operation.Kind)
}

func localOperationImplemented(name string, localOperationNames map[string]bool) bool {
	if localOperationNames[name] || localOperationNames[strings.ToLower(name)] {
		return true
	}
	if name == "user_drafts" && localOperationNames["viewer_drafts"] {
		return true
	}
	return false
}

func statusCountSet(operations []accountedOperation) map[string]int {
	counts := map[string]int{}
	for _, operation := range operations {
		counts[operation.Status]++
	}
	return counts
}

func operationsByStatus(operations []accountedOperation) map[string][]accountedOperation {
	byStatus := map[string][]accountedOperation{}
	for _, operation := range operations {
		byStatus[operation.Status] = append(byStatus[operation.Status], operation)
	}
	for _, group := range byStatus {
		sort.Slice(group, func(left int, right int) bool {
			if group[left].Name == group[right].Name {
				return group[left].Kind < group[right].Kind
			}
			return group[left].Name < group[right].Name
		})
	}
	return byStatus
}

func countWhere[T any](items []T, predicate func(T) bool) int {
	count := 0
	for _, item := range items {
		if predicate(item) {
			count++
		}
	}
	return count
}

func accountedKindCount(operations []accountedOperation, kind string) int {
	return countWhere(operations, func(operation accountedOperation) bool {
		return operation.Kind == kind
	})
}
