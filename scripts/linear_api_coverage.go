//go:build ignore

// Package main generates the Linear API coverage ledger. It is a standalone
// maintenance tool, excluded from the module build/vet/lint/test graph via the
// ignore build tag (and from scripts/go-packages.sh).
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

const defaultUpstreamDir = "/tmp/linctl-upstream-linear"

func main() {
	upstreamDir := flag.String("upstream", defaultUpstreamDir, "upstream linear repo checkout")
	outputPath := flag.String("output", "docs/linear-api-coverage.md", "coverage ledger path")
	operationAuditPath := flag.String("operation-audit", "", "optional SDK operation audit output path")
	flag.Parse()

	upstreamSchemaPath := filepath.Join(*upstreamDir, "packages/sdk/src/schema.graphql")
	upstreamSDKPath := filepath.Join(*upstreamDir, "packages/sdk/src/_generated_sdk.ts")
	upstreamDocumentsPath := filepath.Join(*upstreamDir, "packages/sdk/src/_generated_documents.graphql")
	localOperationsPattern := "internal/client/operations/*.graphql"
	localGeneratedPath := "internal/client/generated.go"
	domainMapPath := "docs/domain-map.md"

	mustValidateUpstreamCheckout(
		*upstreamDir,
		upstreamSchemaPath,
		upstreamSDKPath,
		upstreamDocumentsPath,
	)
	upstreamQueries, upstreamMutations := mustRootFieldSets(upstreamSchemaPath)
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

	ledger := append(bytes.TrimRight(output.Bytes(), "\n"), '\n')
	// #nosec G306 -- this generated markdown ledger is intended to be world-readable repo documentation.
	if err := os.WriteFile(*outputPath, ledger, 0o644); err != nil {
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
	sdkRel, err := filepath.Rel(upstreamDir, upstreamSDKPath)
	if err != nil {
		sdkRel = upstreamSDKPath
	}
	schemaRel, err := filepath.Rel(upstreamDir, upstreamSchemaPath)
	if err != nil {
		schemaRel = upstreamSchemaPath
	}
	fmt.Fprintf(output, "Sources (paths relative to the upstream Linear SDK checkout):\n\n")
	fmt.Fprintf(output, "- Upstream SDK methods: `%s`\n", sdkRel)
	fmt.Fprintf(output, "- Upstream schema roots: `%s`\n", schemaRel)
	fmt.Fprintf(output, "- Local generated operations: `internal/client/generated.go`\n")
	fmt.Fprintf(output, "- Local GraphQL operations: `internal/client/operations/*.graphql`\n")
	fmt.Fprintf(output, "- Public CLI commands: Cobra command inventory enriched by `docs/domain-map.md`\n\n")
	fmt.Fprintf(
		output,
		"Status vocabulary is surface-specific: upstream SDK/root tables use `generated_operation` "+
			"for local GraphQL operation coverage, local operation rows use `generated`, and "+
			"public CLI rows use `public_command` or `guarded_write_command` only when a "+
			"registered command exposes the operation. Generated operations alone are not "+
			"counted as public CLI coverage. Planning statuses remain `accepted_gap`, "+
			"`safe_candidate`, `blocked_needs_design`, and `intentionally_excluded`.\n\n",
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
	fmt.Fprintf(output, "| Surface | Total | Covered/exposed | Classified |\n")
	fmt.Fprintf(output, "| --- | ---: | ---: | ---: |\n")
	fmt.Fprintf(
		output,
		"| Upstream SDK root methods with generated local operations | %d | %d | %d |\n",
		len(sdkMethods),
		countImplementedSDK(sdkMethods, implementedRoots),
		len(sdkMethods),
	)
	fmt.Fprintf(
		output,
		"| Upstream Query root fields used by generated local operations | %d | %d | %d |\n",
		len(queries),
		implementedQueryCount,
		len(queries),
	)
	fmt.Fprintf(
		output,
		"| Upstream Mutation root fields used by generated local operations | %d | %d | %d |\n",
		len(mutations),
		implementedMutationCount,
		len(mutations),
	)
	fmt.Fprintf(
		output,
		"| Local generated Go operations declared in GraphQL files | %d | %d | %d |\n",
		len(localGenerated),
		len(localGenerated),
		len(localGenerated),
	)
	fmt.Fprintf(
		output,
		"| Public CLI commands from command inventory | %d | %d | %d |\n\n",
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
			status = "generated"
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
		status, evidence := classifyDomainLedgerCommand(
			command,
			commandInventory,
			implementedRoots,
			operationRoots,
		)
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

func classifyDomainLedgerCommand(
	command domainCommand,
	commandInventory map[string]cli.CommandInfo,
	implementedRoots map[string]bool,
	operationRoots map[string][]string,
) (string, string) {
	status := "accepted_gap"
	evidence := "planned in `docs/domain-map.md`"
	if commandInfo, ok := commandInventory[command.Command]; ok {
		status, evidence = classifyDomainCommand(commandInfo, implementedRoots, operationRoots)
	}
	if status == "public_command" || status == "guarded_write_command" {
		return status, evidence
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

	return status, evidence
}

func classifyDomainCommand(
	command cli.CommandInfo,
	implementedRoots map[string]bool,
	operationRoots map[string][]string,
) (string, string) {
	references := commandGraphQLReferences(command)
	if len(references) == 0 {
		return "public_command", "`linctl --help` / public CLI tests; no direct GraphQL root in backing"
	}
	for _, reference := range references {
		if implementedRoots[reference.RootKey] {
			return commandCoverageStatus(command), "`linctl --help`, `docs/domain-map.md`, and local GraphQL root"
		}
		for _, key := range operationRoots[reference.OperationName] {
			if implementedRoots[key] {
				return commandCoverageStatus(command), "`linctl --help`, `docs/domain-map.md`, and local GraphQL operation/root"
			}
		}
	}

	return "accepted_gap", "public command exists, but domain-map backing is not matched to local GraphQL roots"
}

func commandCoverageStatus(command cli.CommandInfo) string {
	if command.Safety == cli.CommandSafetyWrite {
		return "guarded_write_command"
	}
	return "public_command"
}
