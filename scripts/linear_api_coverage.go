//go:build ignore

// Package main generates the Linear API coverage ledger. It is a standalone
// maintenance tool, excluded from the module build/vet/lint/test graph via the
// ignore build tag (and from scripts/go-packages.sh).
package main

import (
	"bufio"
	"bytes"
	"context"
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

	"github.com/KyaniteHQ/linctl/internal/cli"
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

type domainReference struct {
	Kind          string
	Field         string
	RootKey       string
	OperationName string
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
	commandInventory := commandInventoryByName(domainCommands)

	implementedRoots := implementedRootSet(localOperations)
	operationRoots := operationRootSet(localOperations)
	rootKinds := rootKindSet(upstreamQueries, upstreamMutations)
	localOperationNames := mapSet(localGenerated)

	var output bytes.Buffer
	writeHeader(&output, *upstreamDir, upstreamSchemaPath, upstreamSDKPath)
	writeSummary(
		&output,
		sdkMethods,
		upstreamQueries,
		upstreamMutations,
		localGenerated,
		domainCommands,
		commandInventory,
		implementedRoots,
		operationRoots,
	)
	writeSDKTable(&output, sdkMethods, commandInventory, implementedRoots, rootKinds)
	writeRootTable(&output, "Upstream Query Root Fields", upstreamQueries, implementedRoots)
	writeRootTable(&output, "Upstream Mutation Root Fields", upstreamMutations, implementedRoots)
	writeLocalOperationsTable(&output, localOperations, localOperationNames)
	writeDomainCommandTable(&output, domainCommands, commandInventory, implementedRoots, operationRoots)

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
	commandInventory map[string]cli.CommandInfo,
	implementedRoots map[string]bool,
	operationRoots map[string][]string,
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
		countImplementedDomain(domainCommands, commandInventory, implementedRoots, operationRoots),
		len(domainCommands),
	)
}

func writeSDKTable(
	output *bytes.Buffer,
	methods []sdkMethod,
	commandInventory map[string]cli.CommandInfo,
	implementedRoots map[string]bool,
	rootKinds map[string]string,
) {
	fmt.Fprintf(output, "## Upstream SDK Root Methods\n\n")
	fmt.Fprintf(output, "| Method | Kind | Status | Evidence |\n")
	fmt.Fprintf(output, "| --- | --- | --- | --- |\n")
	for _, method := range methods {
		status, evidence := classifySDKMethod(method.Name, commandInventory, implementedRoots, rootKinds)
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

func writeDomainCommandTable(
	output *bytes.Buffer,
	commands []domainCommand,
	commandInventory map[string]cli.CommandInfo,
	implementedRoots map[string]bool,
	operationRoots map[string][]string,
) {
	fmt.Fprintf(output, "## Repo Domain-Map Commands\n\n")
	fmt.Fprintf(output, "| Domain | Command | Backing | Scope | Status | Evidence |\n")
	fmt.Fprintf(output, "| --- | --- | --- | --- | --- | --- |\n")
	for _, command := range commands {
		status := "accepted_gap"
		evidence := "planned in `docs/domain-map.md`"
		if commandInfo, ok := commandInventory[command.Command]; ok {
			status, evidence = classifyDomainCommand(commandInfo, implementedRoots, operationRoots)
		}
		if domainCommandBlocked(command.Command) {
			status = "blocked_needs_design"
			evidence = "write command needs explicit target and safety semantics"
		}
		if strings.HasPrefix(command.Scope, "Blocked:") {
			status = "blocked_needs_design"
			evidence = "blocked in `docs/domain-map.md` pending explicit safety semantics"
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

func classifyDomainCommand(
	command cli.CommandInfo,
	implementedRoots map[string]bool,
	operationRoots map[string][]string,
) (string, string) {
	references := commandGraphQLReferences(command)
	if len(references) == 0 {
		return "implemented", "`linctl --help` / public CLI tests; no direct GraphQL root in backing"
	}
	for _, reference := range references {
		if implementedRoots[reference.RootKey] {
			return "implemented", "`linctl --help`, `docs/domain-map.md`, and local GraphQL root"
		}
		for _, key := range operationRoots[reference.OperationName] {
			if implementedRoots[key] {
				return "implemented", "`linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root"
			}
		}
	}

	return "accepted_gap", "public command exists, but domain-map backing is not matched to local GraphQL roots"
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
			implemented[rootKey(operation.Kind, field)] = true
		}
	}
	return implemented
}

func operationRootSet(operations []localOperation) map[string][]string {
	roots := map[string][]string{}
	for _, operation := range operations {
		keys := make([]string, 0, len(operation.RootFields))
		for _, field := range operation.RootFields {
			keys = append(keys, rootKey(operation.Kind, field))
		}
		sort.Strings(keys)
		roots[operation.Name] = keys
		roots[strings.ToLower(operation.Name)] = keys
	}
	return roots
}

func rootKindSet(queries []rootField, mutations []rootField) map[string]string {
	roots := map[string]string{}
	for _, field := range queries {
		roots[strings.ToLower(field.Name)] = field.Kind
	}
	for _, field := range mutations {
		roots[strings.ToLower(field.Name)] = field.Kind
	}

	return roots
}

func mapSet(values []string) map[string]bool {
	set := map[string]bool{}
	for _, value := range values {
		set[value] = true
	}
	return set
}

func commandInventoryByName(domainCommands []domainCommand) map[string]cli.CommandInfo {
	root := cli.NewRootCommand(context.Background(), cli.BuildInfo{})
	inventory := cli.EnrichCommandInventory(
		cli.CommandInventory(root),
		domainCommandBacking(domainCommands),
	)
	commands := map[string]cli.CommandInfo{}
	for _, command := range inventory {
		commands[command.Path] = command
		for _, alias := range command.Aliases {
			commands[alias] = command
		}
	}
	if command, ok := commands["next"]; ok {
		command.Safety = cli.CommandSafetyRead
		command.TargetScope = "`--dry-run` read-only"
		commands["next --dry-run"] = command
	}

	return commands
}

func domainCommandBacking(commands []domainCommand) map[string]cli.CommandBacking {
	backingByPath := map[string]cli.CommandBacking{}
	for _, command := range commands {
		roots := make([]cli.CommandGraphQLRoot, 0, len(domainCommandReferences(command)))
		for _, reference := range domainCommandReferences(command) {
			roots = append(roots, cli.CommandGraphQLRoot{
				Kind:      reference.Kind,
				Field:     reference.Field,
				Operation: reference.OperationName,
			})
		}
		backingByPath[command.Command] = cli.CommandBacking{
			OperationBacking: command.Backing,
			TargetScope:      command.Scope,
			GraphQLRoots:     roots,
		}
	}
	return backingByPath
}

var blockedDomainCommands = map[string]bool{
	"document create":                                   true,
	"document update":                                   true,
	"comment resolve":                                   true,
	"comment unresolve":                                 true,
	"issue-relation create":                             true,
	"issue-relation update":                             true,
	"issue-relation delete":                             true,
	"project-update create":                             true,
	"project-update update":                             true,
	"project-update archive":                            true,
	"project-status create":                             true,
	"project-status update":                             true,
	"project-status archive":                            true,
	"project-status unarchive":                          true,
	"project-label create":                              true,
	"project-label update":                              true,
	"project-label delete":                              true,
	"project-label retire":                              true,
	"project-label restore":                             true,
	"project-relation create":                           true,
	"project-relation update":                           true,
	"project-relation delete":                           true,
	"label create":                                      true,
	"label update":                                      true,
	"team create":                                       true,
	"team update":                                       true,
	"team delete":                                       true,
	"team-membership create":                            true,
	"team-membership update":                            true,
	"team-membership delete":                            true,
	"workflow-state create":                             true,
	"workflow-state update":                             true,
	"workflow-state archive":                            true,
	"time-schedule create":                              true,
	"time-schedule update":                              true,
	"time-schedule delete":                              true,
	"time-schedule upsert-external":                     true,
	"template create":                                   true,
	"template update":                                   true,
	"template delete":                                   true,
	"initiative-relation create":                        true,
	"initiative-relation update":                        true,
	"initiative-relation delete":                        true,
	"initiative-to-project create":                      true,
	"initiative-to-project update":                      true,
	"initiative-to-project delete":                      true,
	"roadmap-to-project create":                         true,
	"roadmap-to-project update":                         true,
	"roadmap-to-project delete":                         true,
	"initiative-update create":                          true,
	"initiative-update update":                          true,
	"initiative-update archive":                         true,
	"initiative-update unarchive":                       true,
	"initiative create":                                 true,
	"initiative update":                                 true,
	"initiative archive":                                true,
	"roadmap create":                                    true,
	"roadmap update":                                    true,
	"roadmap archive":                                   true,
	"roadmap delete":                                    true,
	"custom-view create":                                true,
	"custom-view update":                                true,
	"customer create":                                   true,
	"customer update":                                   true,
	"customer archive":                                  true,
	"customer-need create":                              true,
	"customer-need update":                              true,
	"customer-need archive":                             true,
	"customer-need delete":                              true,
	"customer-status create":                            true,
	"customer-status update":                            true,
	"customer-status delete":                            true,
	"customer-tier create":                              true,
	"customer-tier update":                              true,
	"customer-tier delete":                              true,
	"favorite create":                                   true,
	"favorite update":                                   true,
	"emoji create":                                      true,
	"attachment create":                                 true,
	"attachment update":                                 true,
	"notification archive":                              true,
	"notification archive all":                          true,
	"notification update":                               true,
	"notification mark read all":                        true,
	"notification mark unread all":                      true,
	"notification snooze all":                           true,
	"notification unsnooze all":                         true,
	"notification category channel subscription update": true,
	"notification subscription create":                  true,
	"notification subscription update":                  true,
	"notification subscription delete":                  true,
	"release-pipeline create":                           true,
	"release-pipeline update":                           true,
	"release-pipeline archive":                          true,
	"release-pipeline unarchive":                        true,
	"release-pipeline delete":                           true,
	"release-stage create":                              true,
	"release-stage update":                              true,
	"release-stage archive":                             true,
	"release-stage unarchive":                           true,
	"release create":                                    true,
	"release update":                                    true,
	"release archive":                                   true,
	"release unarchive":                                 true,
	"release delete":                                    true,
	"release complete":                                  true,
	"release sync":                                      true,
	"release-note create":                               true,
	"release-note update":                               true,
	"release-note archive":                              true,
	"release-note delete":                               true,
	"issue-to-release create":                           true,
	"issue-to-release update":                           true,
	"issue-to-release delete":                           true,
}

func domainCommandBlocked(command string) bool {
	return blockedDomainCommands[command]
}

func classifySDKMethod(
	name string,
	commandInventory map[string]cli.CommandInfo,
	implementedRoots map[string]bool,
	rootKinds map[string]string,
) (string, string) {
	if sdkImplemented(name, implementedRoots) {
		return "implemented", "local operation or command exists"
	}
	if command, ok := commandInventory[commandLookupName(name)]; ok && command.OperationBacking != "" {
		return "implemented", "local operation or command exists"
	}
	if kind, ok := sdkRootKind(name, rootKinds); ok {
		return classifyLoose(name, kind)
	}
	if status, rationale, ok := explicitRiskClassification(strings.ToLower(name)); ok {
		return status, rationale
	}

	return "blocked_needs_design", "SDK method is not matched to a GraphQL root field; explicit classification required"
}

func classifyRoot(field rootField, implementedRoots map[string]bool) (string, string) {
	if implementedRoots[rootKey(field.Kind, field.Name)] {
		return "implemented", "root field used by local GraphQL operation"
	}
	return classifyLoose(field.Name, field.Kind)
}

func classifyLoose(name string, kind string) (string, string) {
	lower := strings.ToLower(name)
	if status, rationale, ok := explicitRiskClassification(lower); ok {
		return status, rationale
	}
	switch {
	case strings.Contains(lower, "latestreleasebyaccesskey"),
		strings.Contains(lower, "releasepipelinebyaccesskey"):
		return "intentionally_excluded", accessKeyReleaseRationale()
	case strings.Contains(lower, "documentcontent"),
		strings.Contains(lower, "archivepayload"),
		strings.Contains(lower, "externalthread"):
		return "blocked_needs_design", contentPayloadReadRationale()
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

type riskClassification struct {
	status    string
	rationale string
}

var explicitRiskClassifications = map[string]riskClassification{
	"auditentries": {
		status: "blocked_needs_design",
		rationale: "audit logs can expose actor, IP, country, and request metadata; " +
			"needs explicit admin/security output model",
	},
	"emailintakeaddress": {
		status:    "intentionally_excluded",
		rationale: "email intake administration sits outside the ordinary agent CLI read surface",
	},
	"emailintakeaddress_sesdomainidentity": {
		status:    "intentionally_excluded",
		rationale: "email domain identity administration sits outside the ordinary agent CLI read surface",
	},
	"attachmentlinkgithubissue": {
		status: "blocked_needs_design",
		rationale: "attachment-to-GitHub linking mutates third-party integration state; " +
			"needs explicit integration guard semantics",
	},
	"attachmentlinkjiraissue": {
		status: "blocked_needs_design",
		rationale: "attachment-to-Jira linking mutates third-party integration state; " +
			"needs explicit integration guard semantics",
	},
	"availableusers": {
		status: "intentionally_excluded",
		rationale: "available-user picker enumeration is a specialized product resolver; " +
			"`user list` is the supported user read surface",
	},
	"cycleshiftall": {
		status: "blocked_needs_design",
		rationale: "bulk Cycle date shifting is a state-changing organization operation that " +
			"needs target-pinned guard semantics",
	},
	"cyclestartupcomingcycletoday": {
		status: "blocked_needs_design",
		rationale: "starting an upcoming Cycle changes team planning state and needs " +
			"target-pinned guard semantics",
	},
	"issueaddlabel": {
		status:    "blocked_needs_design",
		rationale: "issue label mutation needs issue target pinning and target-mismatch tests",
	},
	"issueexternalsyncdisable": {
		status: "blocked_needs_design",
		rationale: "issue external-sync disable changes integration state and needs explicit " +
			"integration guard semantics",
	},
	"issueimportcheckcsv": {
		status: "blocked_needs_design",
		rationale: "CSV import validation can expose imported row payloads and needs an " +
			"explicit redaction/output model",
	},
	"issueimportchecksync": {
		status: "blocked_needs_design",
		rationale: "sync import validation can expose external tracker payloads and needs an " +
			"explicit redaction/output model",
	},
	"issueimportcreateasana": {
		status: "blocked_needs_design",
		rationale: "Asana issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreatecsvjira": {
		status: "blocked_needs_design",
		rationale: "CSV/Jira issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreateclubhouse": {
		status: "blocked_needs_design",
		rationale: "Clubhouse issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreategithub": {
		status: "blocked_needs_design",
		rationale: "GitHub issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreatejira": {
		status: "blocked_needs_design",
		rationale: "Jira issue import creation starts external import workflow state and " +
			"needs explicit integration guard semantics",
	},
	"issueimportjqlcheck": {
		status: "blocked_needs_design",
		rationale: "JQL import validation can expose external tracker payloads and needs an " +
			"explicit redaction/output model",
	},
	"issueimportprocess": {
		status: "blocked_needs_design",
		rationale: "issue import processing advances external import workflow state and needs " +
			"explicit integration guard semantics",
	},
	"issuelabelrestore": {
		status:    "blocked_needs_design",
		rationale: "issue label lifecycle restore needs explicit organization/admin safety semantics",
	},
	"issuelabelretire": {
		status:    "blocked_needs_design",
		rationale: "issue label lifecycle retire needs explicit organization/admin safety semantics",
	},
	"issuereminder": {
		status: "blocked_needs_design",
		rationale: "issue reminder mutation changes notification state and needs target-pinned " +
			"guard semantics",
	},
	"issuerepositorysuggestions": {
		status: "intentionally_excluded",
		rationale: "repository suggestion reads expose VCS integration metadata outside the " +
			"default Linear work CLI surface",
	},
	"issuedescriptionupdatefromfront": {
		status: "blocked_needs_design",
		rationale: "Front-origin description updates mutate issue content through integration state; " +
			"needs explicit integration guard semantics",
	},
	"issueimportcreatelinearv2": {
		status: "blocked_needs_design",
		rationale: "Linear v2 issue import creation starts import workflow state and needs explicit " +
			"import guard semantics",
	},
	"issueshare": {
		status:    "blocked_needs_design",
		rationale: "issue sharing changes access state and needs target-pinned guard semantics",
	},
	"issuesubscribe": {
		status: "blocked_needs_design",
		rationale: "issue subscription changes notification state and needs target-pinned " +
			"guard semantics",
	},
	"issueunshare": {
		status:    "blocked_needs_design",
		rationale: "issue unsharing changes access state and needs target-pinned guard semantics",
	},
	"issueunsubscribe": {
		status: "blocked_needs_design",
		rationale: "issue unsubscribe changes notification state and needs target-pinned " +
			"guard semantics",
	},
	"latestreleasebyaccesskey": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"latestreleasebyaccesskey_history": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"latestreleasebyaccesskey_links": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"initiativelabeladd": {
		status:    "blocked_needs_design",
		rationale: "initiative label mutation needs initiative target pinning and target-mismatch tests",
	},
	"initiativeaddlabel": {
		status:    "blocked_needs_design",
		rationale: "initiative label mutation needs initiative target pinning and target-mismatch tests",
	},
	"microsoftteamschannels": {
		status: "intentionally_excluded",
		rationale: "Microsoft Teams channel enumeration exposes chat integration metadata outside the " +
			"default Linear work CLI surface",
	},
	"organizationinvite": {
		status:    "intentionally_excluded",
		rationale: organizationInviteRationale(),
	},
	"organizationinvites": {
		status:    "intentionally_excluded",
		rationale: organizationInviteRationale(),
	},
	"organization_subscription": {
		status:    "intentionally_excluded",
		rationale: "organization subscription and billing state is outside the ordinary agent CLI surface",
	},
	"pushsubscriptiontest": {
		status: "intentionally_excluded",
		rationale: "push subscription diagnostics are notification-device integration plumbing " +
			"outside the CLI surface",
	},
	"projectlabelrestore": {
		status:    "blocked_needs_design",
		rationale: "project label lifecycle restore needs explicit organization/admin safety semantics",
	},
	"projectlabelretire": {
		status:    "blocked_needs_design",
		rationale: "project label lifecycle retire needs explicit organization/admin safety semantics",
	},
	"projectaddlabel": {
		status:    "blocked_needs_design",
		rationale: "project label mutation needs project target pinning and target-mismatch tests",
	},
	"projectexternalsyncdisable": {
		status: "blocked_needs_design",
		rationale: "project external-sync disable changes integration state and needs explicit " +
			"integration guard semantics",
	},
	"projectcreateslackchannel": {
		status: "blocked_needs_design",
		rationale: "project Slack channel creation mutates chat integration state and needs explicit " +
			"integration guard semantics",
	},
	"projectreassignstatus": {
		status: "blocked_needs_design",
		rationale: "project status reassignment mutates project workflow state and needs target-pinned " +
			"guard semantics",
	},
	"recentreleasesbyaccesskey": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"releasepipelinebyaccesskey": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"releasepipelinebyaccesskey_releases": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"releasepipelinebyaccesskey_stages": {
		status:    "intentionally_excluded",
		rationale: accessKeyReleaseRationale(),
	},
	"ssourlfromemail": {
		status:    "intentionally_excluded",
		rationale: "SSO discovery from email belongs to auth flow tooling, not the Linear work CLI",
	},
	"userchangerole": {
		status:    "intentionally_excluded",
		rationale: "user role changes are organization administration outside the ordinary agent CLI surface",
	},
	"userdiscordconnect": {
		status:    "intentionally_excluded",
		rationale: "Discord account connection belongs to user auth/integration setup, not work CLI reads",
	},
	"userexternaluserdisconnect": {
		status: "intentionally_excluded",
		rationale: "external-user disconnection is identity integration administration outside the " +
			"ordinary agent CLI surface",
	},
	"usersettingsflagsreset": {
		status: "intentionally_excluded",
		rationale: "user settings flag reset is internal preference administration outside the " +
			"ordinary agent CLI surface",
	},
	"userunlinkfromidentityprovider": {
		status:    "intentionally_excluded",
		rationale: "identity-provider unlinking is auth administration outside the ordinary agent CLI surface",
	},
	"verifygithubenterpriseserverinstallation": {
		status: "intentionally_excluded",
		rationale: "GitHub Enterprise installation verification is integration administration " +
			"outside the CLI surface",
	},
}

func explicitRiskClassification(lowerName string) (string, string, bool) {
	classification, ok := explicitRiskClassifications[lowerName]
	return classification.status, classification.rationale, ok
}

func accessKeyReleaseRationale() string {
	return "access-key release reads are unauthenticated sharing surfaces " +
		"outside the token-scoped agent CLI"
}

func contentPayloadReadRationale() string {
	return "content, thread, and archive payload reads can expose body/blob data; " +
		"needs explicit opt-in projection before CLI exposure"
}

func organizationInviteRationale() string {
	return "organization invite reads can expose invitee and admin metadata " +
		"outside an agent-safe CLI surface"
}

func countImplemented(fields []rootField, implementedRoots map[string]bool) int {
	return countWhere(fields, func(field rootField) bool {
		return implementedRoots[rootKey(field.Kind, field.Name)]
	})
}

func countImplementedSDK(methods []sdkMethod, implementedRoots map[string]bool) int {
	return countWhere(methods, func(method sdkMethod) bool {
		return sdkImplemented(method.Name, implementedRoots)
	})
}

func sdkImplemented(name string, implementedRoots map[string]bool) bool {
	if implementedRoots[rootKey("query", name)] || implementedRoots[rootKey("mutation", name)] {
		return true
	}
	for _, candidate := range sdkMutationRootCandidates(name) {
		if implementedRoots[rootKey("mutation", candidate)] {
			return true
		}
	}
	return false
}

func sdkRootKind(name string, rootKinds map[string]string) (string, bool) {
	if kind, ok := rootKinds[strings.ToLower(name)]; ok {
		return kind, true
	}
	for _, candidate := range sdkMutationRootCandidates(name) {
		if kind, ok := rootKinds[strings.ToLower(candidate)]; ok {
			return kind, true
		}
	}

	return "", false
}

func rootKey(kind string, name string) string {
	return strings.ToLower(kind) + ":" + name
}

func sdkMutationRootCandidates(name string) []string {
	candidates := []string{}
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
		if strings.HasPrefix(lowerName, prefix) || strings.HasSuffix(lowerName, prefix) {
			return true
		}
	}
	return false
}

func countImplementedDomain(
	commands []domainCommand,
	commandInventory map[string]cli.CommandInfo,
	implementedRoots map[string]bool,
	operationRoots map[string][]string,
) int {
	return countWhere(commands, func(command domainCommand) bool {
		commandInfo, ok := commandInventory[command.Command]
		if !ok {
			return false
		}
		status, _ := classifyDomainCommand(commandInfo, implementedRoots, operationRoots)
		return status == "implemented"
	})
}

var (
	kebabCasePattern      = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	commandLookupReplacer = strings.NewReplacer("-", " ", "_", " ")
	domainRootPattern     = regexp.MustCompile(`\b(Query|Mutation)\.([A-Za-z_][A-Za-z0-9_]*)`)
)

func domainCommandReferences(command domainCommand) []domainReference {
	matches := domainRootPattern.FindAllStringSubmatch(command.Backing, -1)
	references := make([]domainReference, 0, len(matches))
	seen := map[string]bool{}
	for _, match := range matches {
		kind := strings.ToLower(match[1])
		name := match[2]
		key := rootKey(kind, name)
		if seen[key] {
			continue
		}
		seen[key] = true
		references = append(references, domainReference{
			Kind:          kind,
			Field:         name,
			RootKey:       key,
			OperationName: name,
		})
	}
	sort.Slice(references, func(left int, right int) bool {
		return references[left].RootKey < references[right].RootKey
	})
	return references
}

func commandGraphQLReferences(command cli.CommandInfo) []domainReference {
	references := make([]domainReference, 0, len(command.GraphQLRoots))
	seen := map[string]bool{}
	for _, root := range command.GraphQLRoots {
		key := rootKey(root.Kind, root.Field)
		if seen[key] {
			continue
		}
		seen[key] = true
		references = append(references, domainReference{
			Kind:          root.Kind,
			Field:         root.Field,
			RootKey:       key,
			OperationName: root.Operation,
		})
	}
	sort.Slice(references, func(left int, right int) bool {
		return references[left].RootKey < references[right].RootKey
	})
	return references
}

func kebabCase(value string) string {
	return strings.ToLower(kebabCasePattern.ReplaceAllString(value, `${1}-${2}`))
}

func commandLookupName(value string) string {
	return commandLookupReplacer.Replace(kebabCase(value))
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
	// #nosec G204 -- arguments are fixed git metadata commands assembled by this generator.
	output, err := exec.Command("git", command...).Output()
	if err != nil {
		return "unknown"
	}
	return string(output)
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
