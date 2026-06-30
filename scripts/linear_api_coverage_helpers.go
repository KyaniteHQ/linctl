//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/KyaniteHQ/linctl/internal/cli"
)

func countImplementedDomain(
	commands []domainCommand,
	commandInventory map[string]cli.CommandInfo,
	implementedRoots map[string]bool,
	operationRoots map[string][]string,
) int {
	return countWhere(commands, func(command domainCommand) bool {
		status, _ := classifyDomainLedgerCommand(
			command,
			commandInventory,
			implementedRoots,
			operationRoots,
		)
		return status == "public_command" || status == "guarded_write_command"
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

func mustValidateUpstreamCheckout(upstreamDir string, requiredPaths ...string) {
	if _, err := os.Stat(filepath.Join(upstreamDir, ".git")); err != nil {
		fail(fmt.Errorf("upstream Linear SDK checkout not found at %s: %w", upstreamDir, err))
	}
	for _, path := range requiredPaths {
		if _, err := os.Stat(path); err != nil {
			fail(fmt.Errorf("upstream Linear SDK checkout at %s is missing %s: %w", upstreamDir, path, err))
		}
	}
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
