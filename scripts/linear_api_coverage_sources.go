//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func mustRootFieldSets(path string) ([]rootField, []rootField) {
	source := mustRead(path)
	document, err := parser.ParseSchema(&ast.Source{Name: path, Input: string(source)})
	if err != nil {
		fail(err)
	}
	return rootFields(document, path, "Query"), rootFields(document, path, "Mutation")
}

func rootFields(document *ast.SchemaDocument, path string, typeName string) []rootField {
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
