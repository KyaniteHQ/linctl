// Package main generates the Linear API coverage ledger.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

type rootField struct {
	Name        string
	Kind        string
	ReturnType  string
	Description string
}

type sdkMethod struct {
	Name string
	Kind string
}

type localOperation struct {
	Name       string
	Kind       string
	RootFields []string
	Path       string
}

type domainCommand struct {
	Section string
	Command string
	Backing string
	Scope   string
}

func main() {
	upstreamDir := flag.String("upstream", "/tmp/linctl-upstream-linear", "upstream linear repo checkout")
	outputPath := flag.String("output", "docs/linear-api-coverage.md", "coverage ledger path")
	operationAuditPath := flag.String("operation-audit", "", "optional SDK operation audit output path")
	flag.Parse()

	upstreamSchemaPath := filepath.Join(*upstreamDir, "packages/sdk/src/schema.graphql")
	upstreamSDKPath := filepath.Join(*upstreamDir, "packages/sdk/src/_generated_sdk.ts")
	upstreamDocumentsPath := filepath.Join(*upstreamDir, "packages/sdk/src/_generated_documents.graphql")
	localOperationsPattern := "internal/client/operations/*.graphql"
	localGeneratedPath := "internal/client/generated.go"
	domainMapPath := "docs/domain-map.md"

	upstreamQueries := mustRootFields(upstreamSchemaPath, "Query")
	upstreamMutations := mustRootFields(upstreamSchemaPath, "Mutation")
	sdkMethods := mustSDKMethods(upstreamSDKPath)
	localOperations := mustLocalOperations(localOperationsPattern)
	localGenerated := mustGeneratedOperations(localGeneratedPath)
	domainCommands := mustDomainCommands(domainMapPath)
	sdkOperations := mustSDKOperations(upstreamDocumentsPath)

	implementedRoots := implementedRootSet(localOperations)
	localOperationNames := mapSet(localGenerated)
	commandNames := domainCommandSet(domainCommands)

	var output bytes.Buffer
	writeHeader(&output, *upstreamDir, upstreamSchemaPath, upstreamSDKPath)
	writeSummary(
		&output,
		sdkMethods,
		upstreamQueries,
		upstreamMutations,
		localGenerated,
		domainCommands,
		implementedRoots,
	)
	writeSDKTable(&output, sdkMethods, commandNames, implementedRoots)
	writeRootTable(&output, "Upstream Query Root Fields", upstreamQueries, implementedRoots)
	writeRootTable(&output, "Upstream Mutation Root Fields", upstreamMutations, implementedRoots)
	writeLocalOperationsTable(&output, localOperations, localOperationNames)
	writeDomainCommandTable(&output, domainCommands, commandNames)

	// #nosec G306 -- this generated markdown ledger is intended to be world-readable repo documentation.
	if err := os.WriteFile(*outputPath, output.Bytes(), 0o644); err != nil {
		fail(err)
	}
	if *operationAuditPath != "" {
		var audit bytes.Buffer
		writeSDKOperationAudit(&audit, upstreamDocumentsPath, sdkOperations, localOperations)
		// #nosec G306 -- this generated markdown audit is intended to be world-readable repo documentation.
		if err := os.WriteFile(*operationAuditPath, audit.Bytes(), 0o644); err != nil {
			fail(err)
		}
	}
}

func writeHeader(output *bytes.Buffer, upstreamDir string, upstreamSchemaPath string, upstreamSDKPath string) {
	commit := strings.TrimSpace(runGit(upstreamDir, "rev-parse", "--short", "HEAD"))
	fmt.Fprintf(output, "# Linear API coverage ledger\n\n")
	fmt.Fprintf(output, "Generated from current local sources and upstream Linear SDK commit `%s`.\n\n", commit)
	fmt.Fprintf(output, "Sources:\n\n")
	fmt.Fprintf(output, "- Upstream SDK methods: `%s`\n", upstreamSDKPath)
	fmt.Fprintf(output, "- Upstream schema roots: `%s`\n", upstreamSchemaPath)
	fmt.Fprintf(output, "- Local generated operations: `internal/client/generated.go`\n")
	fmt.Fprintf(output, "- Local GraphQL operations: `internal/client/operations/*.graphql`\n")
	fmt.Fprintf(output, "- Repo domain map: `docs/domain-map.md`\n\n")
	fmt.Fprintf(
		output,
		"Statuses: `implemented`, `accepted_gap`, `safe_candidate`, "+
			"`blocked_needs_design`, `intentionally_excluded`.\n\n",
	)
}

func writeSummary(
	output *bytes.Buffer,
	sdkMethods []sdkMethod,
	queries []rootField,
	mutations []rootField,
	localGenerated []string,
	domainCommands []domainCommand,
	implementedRoots map[string]bool,
) {
	implementedQueryCount := countImplemented(queries, implementedRoots)
	implementedMutationCount := countImplemented(mutations, implementedRoots)

	fmt.Fprintf(output, "## Summary\n\n")
	fmt.Fprintf(output, "| Surface | Total | Implemented/root-backed | Classified |\n")
	fmt.Fprintf(output, "| --- | ---: | ---: | ---: |\n")
	fmt.Fprintf(
		output,
		"| Upstream SDK root methods | %d | %d | %d |\n",
		len(sdkMethods),
		countImplementedSDK(sdkMethods, implementedRoots),
		len(sdkMethods),
	)
	fmt.Fprintf(
		output,
		"| Upstream Query root fields | %d | %d | %d |\n",
		len(queries),
		implementedQueryCount,
		len(queries),
	)
	fmt.Fprintf(
		output,
		"| Upstream Mutation root fields | %d | %d | %d |\n",
		len(mutations),
		implementedMutationCount,
		len(mutations),
	)
	fmt.Fprintf(
		output,
		"| Local generated Go operations | %d | %d | %d |\n",
		len(localGenerated),
		len(localGenerated),
		len(localGenerated),
	)
	fmt.Fprintf(
		output,
		"| Domain-map commands | %d | %d | %d |\n\n",
		len(domainCommands),
		countImplementedDomain(domainCommands),
		len(domainCommands),
	)
}

func writeSDKTable(
	output *bytes.Buffer,
	methods []sdkMethod,
	commandNames map[string]bool,
	implementedRoots map[string]bool,
) {
	fmt.Fprintf(output, "## Upstream SDK Root Methods\n\n")
	fmt.Fprintf(output, "| Method | Kind | Status | Evidence |\n")
	fmt.Fprintf(output, "| --- | --- | --- | --- |\n")
	for _, method := range methods {
		status, evidence := classifyName(method.Name, method.Kind, commandNames, implementedRoots)
		fmt.Fprintf(output, "| `%s` | %s | %s | %s |\n", method.Name, method.Kind, status, evidence)
	}
	fmt.Fprintf(output, "\n")
}

func writeRootTable(output *bytes.Buffer, title string, fields []rootField, implementedRoots map[string]bool) {
	fmt.Fprintf(output, "## %s\n\n", title)
	fmt.Fprintf(output, "| Field | Return type | Status | Evidence |\n")
	fmt.Fprintf(output, "| --- | --- | --- | --- |\n")
	for _, field := range fields {
		status, evidence := classifyRoot(field, implementedRoots)
		fmt.Fprintf(output, "| `%s` | `%s` | %s | %s |\n", field.Name, field.ReturnType, status, evidence)
	}
	fmt.Fprintf(output, "\n")
}

func writeLocalOperationsTable(output *bytes.Buffer, operations []localOperation, localOperationNames map[string]bool) {
	fmt.Fprintf(output, "## Local Generated Go Operations\n\n")
	fmt.Fprintf(output, "| Operation | Kind | Root fields | Status | Evidence |\n")
	fmt.Fprintf(output, "| --- | --- | --- | --- | --- |\n")
	for _, operation := range operations {
		status := "accepted_gap"
		evidence := "operation is declared but generated function not found"
		if localOperationNames[operation.Name] {
			status = "implemented"
			evidence = "`internal/client/generated.go`"
		}
		fmt.Fprintf(
			output,
			"| `%s` | %s | `%s` | %s | %s |\n",
			operation.Name,
			operation.Kind,
			strings.Join(operation.RootFields, "`, `"),
			status,
			evidence,
		)
	}
	fmt.Fprintf(output, "\n")
}

func writeDomainCommandTable(output *bytes.Buffer, commands []domainCommand, commandNames map[string]bool) {
	fmt.Fprintf(output, "## Repo Domain-Map Commands\n\n")
	fmt.Fprintf(output, "| Domain | Command | Backing | Scope | Status | Evidence |\n")
	fmt.Fprintf(output, "| --- | --- | --- | --- | --- | --- |\n")
	for _, command := range commands {
		status := "accepted_gap"
		evidence := "planned in `docs/domain-map.md`"
		if commandNames[command.Command] {
			status = "implemented"
			evidence = "`linctl --help` / public CLI tests"
		}
		if domainCommandBlocked(command.Command) {
			status = "blocked_needs_design"
			evidence = "write command needs explicit target and safety semantics"
		}
		if strings.Contains(command.Command, "delete") {
			status = "blocked_needs_design"
			evidence = "destructive command needs explicit safety semantics"
		}
		isSprintNonReport := strings.Contains(command.Command, "sprint ") &&
			!strings.Contains(command.Command, "current") &&
			!strings.Contains(command.Command, "report")
		if isSprintNonReport {
			status = "intentionally_excluded"
			evidence = "Sprint is a read-only alias over Cycle"
		}
		fmt.Fprintf(
			output,
			"| %s | `%s` | %s | %s | %s | %s |\n",
			command.Section,
			command.Command,
			escapePipes(command.Backing),
			escapePipes(command.Scope),
			status,
			evidence,
		)
	}
	fmt.Fprintf(output, "\n")
}

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
	if localOperationNames[operation.Name] || localOperationNames[strings.ToLower(operation.Name)] {
		return "implemented", "backed by a local GraphQL operation in internal/client/operations"
	}
	return classifyLoose(operation.Name, operation.Kind)
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

func accountedKindCount(operations []accountedOperation, kind string) int {
	count := 0
	for _, operation := range operations {
		if operation.Kind == kind {
			count++
		}
	}
	return count
}

func mustRootFields(path string, typeName string) []rootField {
	source := mustRead(path)
	document, err := parser.ParseSchema(&ast.Source{Name: path, Input: string(source)})
	if err != nil {
		fail(err)
	}
	definition := document.Definitions.ForName(typeName)
	if definition == nil {
		fail(fmt.Errorf("%s not found in %s", typeName, path))
	}

	fields := make([]rootField, 0, len(definition.Fields))
	for _, field := range definition.Fields {
		fields = append(fields, rootField{
			Name:        field.Name,
			Kind:        strings.ToLower(typeName),
			ReturnType:  field.Type.String(),
			Description: field.Description,
		})
	}
	sort.Slice(fields, func(left int, right int) bool {
		return fields[left].Name < fields[right].Name
	})
	return fields
}

func mustSDKMethods(path string) []sdkMethod {
	input := string(mustRead(path))
	start := strings.Index(input, "export class LinearSdk extends Request")
	if start < 0 {
		fail(fmt.Errorf("LinearSdk class not found in %s", path))
	}
	input = input[start:]
	end := strings.Index(input, "\n}\nexport {")
	if end < 0 {
		fail(fmt.Errorf("LinearSdk class end not found in %s", path))
	}
	input = input[:end]

	methodPattern := regexp.MustCompile(`(?m)^\s*public\s+(get\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*[\(:]`)
	matches := methodPattern.FindAllStringSubmatch(input, -1)
	methods := make([]sdkMethod, 0, len(matches))
	for _, match := range matches {
		kind := "method"
		if strings.TrimSpace(match[1]) == "get" {
			kind = "getter"
		}
		methods = append(methods, sdkMethod{Name: match[2], Kind: kind})
	}
	sort.Slice(methods, func(left int, right int) bool {
		return methods[left].Name < methods[right].Name
	})
	return methods
}

func mustSDKOperations(path string) []sdkMethod {
	source := mustRead(path)
	document, err := parser.ParseQuery(&ast.Source{Name: path, Input: string(source)})
	if err != nil {
		fail(err)
	}
	operations := make([]sdkMethod, 0, len(document.Operations))
	for _, operation := range document.Operations {
		operations = append(operations, sdkMethod{
			Name: operation.Name,
			Kind: string(operation.Operation),
		})
	}
	sort.Slice(operations, func(left int, right int) bool {
		if operations[left].Name == operations[right].Name {
			return operations[left].Kind < operations[right].Kind
		}
		return operations[left].Name < operations[right].Name
	})
	return operations
}

func mustLocalOperations(pattern string) []localOperation {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		fail(err)
	}
	if len(paths) == 0 {
		fail(fmt.Errorf("no operation files match %s", pattern))
	}

	operations := []localOperation{}
	for _, path := range paths {
		source := mustRead(path)
		document, err := parser.ParseQuery(&ast.Source{Name: path, Input: string(source)})
		if err != nil {
			fail(err)
		}
		for _, operation := range document.Operations {
			fields := make([]string, 0, len(operation.SelectionSet))
			for _, selection := range operation.SelectionSet {
				if field, ok := selection.(*ast.Field); ok {
					fields = append(fields, field.Name)
				}
			}
			operations = append(operations, localOperation{
				Name:       operation.Name,
				Kind:       string(operation.Operation),
				RootFields: fields,
				Path:       path,
			})
		}
	}
	sort.Slice(operations, func(left int, right int) bool {
		return operations[left].Name < operations[right].Name
	})
	return operations
}

func mustGeneratedOperations(path string) []string {
	pattern := regexp.MustCompile(`^func ([A-Za-z_][A-Za-z0-9_]*)\(`)
	names := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(mustRead(path)))
	for scanner.Scan() {
		match := pattern.FindStringSubmatch(scanner.Text())
		if match != nil {
			names = append(names, match[1])
		}
	}
	if err := scanner.Err(); err != nil {
		fail(err)
	}
	sort.Strings(names)
	return names
}

func mustDomainCommands(path string) []domainCommand {
	commands := []domainCommand{}
	section := ""
	scanner := bufio.NewScanner(bytes.NewReader(mustRead(path)))
	pattern := regexp.MustCompile(`^\| ` + "`" + `([^` + "`" + `]+)` + "`" + ` \| ([^|]+) \| ([^|]+) \|$`)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## ") {
			section = strings.TrimSpace(strings.TrimPrefix(line, "## "))
			continue
		}
		match := pattern.FindStringSubmatch(line)
		if match == nil || match[1] == "Command" {
			continue
		}
		commands = append(commands, domainCommand{
			Section: section,
			Command: match[1],
			Backing: strings.TrimSpace(match[2]),
			Scope:   strings.TrimSpace(match[3]),
		})
	}
	if err := scanner.Err(); err != nil {
		fail(err)
	}
	return commands
}

func operationNameSet(operations []localOperation) map[string]bool {
	names := map[string]bool{}
	for _, operation := range operations {
		names[operation.Name] = true
		names[strings.ToLower(operation.Name)] = true
	}
	return names
}

func implementedRootSet(operations []localOperation) map[string]bool {
	implemented := map[string]bool{}
	for _, operation := range operations {
		for _, field := range operation.RootFields {
			implemented[field] = true
		}
	}
	return implemented
}

func mapSet(values []string) map[string]bool {
	set := map[string]bool{}
	for _, value := range values {
		set[value] = true
	}
	return set
}

func domainCommandSet(commands []domainCommand) map[string]bool {
	set := map[string]bool{}
	for _, command := range commands {
		if commandImplemented(command.Command) {
			set[command.Command] = true
		}
	}
	return set
}

func commandImplemented(command string) bool {
	implemented := map[string]bool{
		"whoami":                   true,
		"target":                   true,
		"issue list":               true,
		"issue search":             true,
		"issue get":                true,
		"issue deps":               true,
		"issue id":                 true,
		"issue title":              true,
		"issue url":                true,
		"issue branch":             true,
		"issue pr":                 true,
		"next --dry-run":           true,
		"done":                     true,
		"issue create":             true,
		"issue update":             true,
		"issue start":              true,
		"issue comment":            true,
		"issue reply":              true,
		"issue close":              true,
		"issue comments":           true,
		"comment list":             true,
		"comment get":              true,
		"cycle list":               true,
		"cycle get":                true,
		"cycle create":             true,
		"cycle update":             true,
		"cycle archive":            true,
		"sprint current":           true,
		"sprint report":            true,
		"project list":             true,
		"project get":              true,
		"project create":           true,
		"project update":           true,
		"project archive":          true,
		"project members":          true,
		"project updates":          true,
		"project-update list":      true,
		"project-update get":       true,
		"project-milestone list":   true,
		"project-milestone get":    true,
		"project-milestone create": true,
		"project-milestone update": true,
		"document list":            true,
		"document get":             true,
		"label list":               true,
		"label get":                true,
		"team list":                true,
		"team get":                 true,
		"team members":             true,
		"user list":                true,
		"user get":                 true,
		"user me":                  true,
		"workflow-state list":      true,
		"workflow-state get":       true,
		"initiative list":          true,
		"initiative get":           true,
		"custom-view list":         true,
		"custom-view get":          true,
		"favorite list":            true,
		"favorite get":             true,
		"emoji list":               true,
		"emoji get":                true,
	}
	return implemented[command]
}

func domainCommandBlocked(command string) bool {
	blocked := map[string]bool{
		"document create":        true,
		"document update":        true,
		"comment resolve":        true,
		"comment unresolve":      true,
		"project-update create":  true,
		"project-update update":  true,
		"project-update archive": true,
		"label create":           true,
		"label update":           true,
		"team create":            true,
		"team update":            true,
		"workflow-state create":  true,
		"workflow-state update":  true,
		"workflow-state archive": true,
		"initiative create":      true,
		"initiative update":      true,
		"initiative archive":     true,
		"custom-view create":     true,
		"custom-view update":     true,
		"favorite create":        true,
		"favorite update":        true,
		"emoji create":           true,
	}
	return blocked[command]
}

func classifyName(
	name string,
	kind string,
	commandNames map[string]bool,
	implementedRoots map[string]bool,
) (string, string) {
	if sdkImplemented(name, implementedRoots) || commandNames[strings.ReplaceAll(kebabCase(name), "-", " ")] {
		return "implemented", "local operation or command exists"
	}
	return classifyLoose(name, kind)
}

func classifyRoot(field rootField, implementedRoots map[string]bool) (string, string) {
	if implementedRoots[field.Name] {
		return "implemented", "root field used by local GraphQL operation"
	}
	return classifyLoose(field.Name, field.Kind)
}

func classifyLoose(name string, kind string) (string, string) {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "delete"),
		strings.Contains(lower, "remove"),
		strings.Contains(lower, "revoke"),
		strings.Contains(lower, "suspend"):
		return "blocked_needs_design", "destructive or access-changing operation needs explicit safety model"
	case strings.Contains(lower, "admin"),
		strings.Contains(lower, "auth"),
		strings.Contains(lower, "oauth"),
		strings.Contains(lower, "session"),
		strings.Contains(lower, "webhook"),
		strings.Contains(lower, "integration"):
		return "intentionally_excluded", "admin/auth/internal integration surface outside ordinary agent CLI"
	case hasWritePrefix(lower):
		return "blocked_needs_design", "write operation needs guarded target semantics before exposure"
	case strings.Contains(lower, "resolve"):
		return "blocked_needs_design", "state-changing operation needs guarded target semantics before exposure"
	case strings.Contains(lower, "issue"),
		strings.Contains(lower, "project"),
		strings.Contains(lower, "cycle"),
		strings.Contains(lower, "document"),
		strings.Contains(lower, "label"),
		strings.Contains(lower, "team"),
		strings.Contains(lower, "user"),
		strings.Contains(lower, "comment"):
		return "accepted_gap", "repo-planned or likely useful CLI domain"
	default:
		if kind == "mutation" {
			return "blocked_needs_design", "mutation needs product and safety design"
		}
		return "safe_candidate", "read operation may fit future CLI coverage"
	}
}

func countImplemented(fields []rootField, implementedRoots map[string]bool) int {
	count := 0
	for _, field := range fields {
		if implementedRoots[field.Name] {
			count++
		}
	}
	return count
}

func countImplementedSDK(methods []sdkMethod, implementedRoots map[string]bool) int {
	count := 0
	for _, method := range methods {
		if sdkImplemented(method.Name, implementedRoots) {
			count++
		}
	}
	return count
}

func sdkImplemented(name string, implementedRoots map[string]bool) bool {
	for _, candidate := range sdkRootCandidates(name) {
		if implementedRoots[candidate] {
			return true
		}
	}
	return false
}

func sdkRootCandidates(name string) []string {
	candidates := []string{name}
	for _, prefix := range []string{"create", "update", "archive", "delete", "unarchive", "cancel"} {
		if strings.HasPrefix(name, prefix) && len(name) > len(prefix) {
			entity := lowerFirst(strings.TrimPrefix(name, prefix))
			candidates = append(candidates, entity+upperFirst(prefix))
		}
	}
	return candidates
}

func hasWritePrefix(lowerName string) bool {
	for _, prefix := range []string{
		"create",
		"update",
		"archive",
		"delete",
		"unarchive",
		"cancel",
		"mark",
		"move",
		"rotate",
	} {
		if strings.HasPrefix(lowerName, prefix) {
			return true
		}
	}
	return false
}

func countImplementedDomain(commands []domainCommand) int {
	count := 0
	for _, command := range commands {
		if commandImplemented(command.Command) {
			count++
		}
	}
	return count
}

func kebabCase(value string) string {
	pattern := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	return strings.ToLower(pattern.ReplaceAllString(value, `${1}-${2}`))
}

func lowerFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToLower(value[:1]) + value[1:]
}

func upperFirst(value string) string {
	if value == "" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func escapePipes(value string) string {
	return strings.ReplaceAll(value, "|", "\\|")
}

func mustRead(path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		fail(err)
	}
	return content
}

func runGit(dir string, args ...string) string {
	command := append([]string{"-C", dir}, args...)
	output, err := runGitCommand(command...)
	if err != nil {
		return "unknown"
	}
	return output
}

func runGitCommand(args ...string) (string, error) {
	// #nosec G204 -- arguments are fixed git metadata commands assembled by this generator.
	command := exec.Command("git", args...)
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
