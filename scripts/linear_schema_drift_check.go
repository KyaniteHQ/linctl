//go:build ignore

// Package main reports semantic drift between the vendored Linear GraphQL
// schema and the upstream Linear SDK schema: types, fields, and enum values
// added or removed, and field type signature changes on shared types. It does
// not compare descriptions or declaration order, since the two files are
// different renderings of the same schema (live introspection vs. the SDK's
// generated file) and a textual diff would be noise.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator"
)

func main() {
	upstreamDir := flag.String("upstream", "", "required upstream Linear SDK checkout")
	vendoredPath := flag.String("vendored", "internal/client/schema.graphql", "vendored Linear schema path")
	maxReport := flag.Int("max-report", 50, "maximum named items listed per drift category")
	flag.Parse()

	if *upstreamDir == "" {
		exitError(fmt.Errorf("missing -upstream: set it to a Linear SDK checkout path"))
	}

	upstreamPath := filepath.Join(*upstreamDir, "packages/sdk/src/schema.graphql")
	upstreamSchema := loadSchema(upstreamPath)
	vendoredSchema := loadSchema(*vendoredPath)

	report := diffSchemas(vendoredSchema, upstreamSchema)
	if report.empty() {
		fmt.Println("schema drift: none")
		return
	}

	report.print(*maxReport)
	os.Exit(1)
}

func loadSchema(path string) *ast.Schema {
	source, err := os.ReadFile(path)
	if err != nil {
		exitError(fmt.Errorf("read %s: %w", path, err))
	}
	schema, err := validator.LoadSchema(validator.Prelude, &ast.Source{Name: path, Input: string(source)})
	if err != nil {
		exitError(fmt.Errorf("load %s: %w", path, err))
	}
	return schema
}

// drift holds the categorized comparison result. Named items are stored
// unsorted; print sorts them for deterministic output.
type drift struct {
	typesAdded        []string
	typesRemoved      []string
	fieldsAdded       []string
	fieldsRemoved     []string
	fieldTypesChanged []string
	enumValuesAdded   []string
	enumValuesRemoved []string
}

func (d *drift) empty() bool {
	return len(d.typesAdded) == 0 && len(d.typesRemoved) == 0 &&
		len(d.fieldsAdded) == 0 && len(d.fieldsRemoved) == 0 &&
		len(d.fieldTypesChanged) == 0 &&
		len(d.enumValuesAdded) == 0 && len(d.enumValuesRemoved) == 0
}

func (d *drift) print(maxReport int) {
	fmt.Printf(
		"schema drift: %d types added upstream, %d removed, %d fields added, %d fields removed, %d field types changed, %d enum values added, %d enum values removed\n",
		len(d.typesAdded), len(d.typesRemoved),
		len(d.fieldsAdded), len(d.fieldsRemoved), len(d.fieldTypesChanged),
		len(d.enumValuesAdded), len(d.enumValuesRemoved),
	)
	printCategory("types added upstream", d.typesAdded, maxReport)
	printCategory("types removed upstream", d.typesRemoved, maxReport)
	printCategory("fields added", d.fieldsAdded, maxReport)
	printCategory("fields removed", d.fieldsRemoved, maxReport)
	printCategory("field types changed", d.fieldTypesChanged, maxReport)
	printCategory("enum values added", d.enumValuesAdded, maxReport)
	printCategory("enum values removed", d.enumValuesRemoved, maxReport)
}

func printCategory(label string, items []string, maxReport int) {
	if len(items) == 0 {
		return
	}
	sorted := append([]string(nil), items...)
	sort.Strings(sorted)
	fmt.Printf("\n%s (%d):\n", label, len(sorted))
	shown := sorted
	if len(shown) > maxReport {
		shown = shown[:maxReport]
	}
	for _, item := range shown {
		fmt.Printf("  %s\n", item)
	}
	if len(sorted) > maxReport {
		fmt.Printf("  ... %d more\n", len(sorted)-maxReport)
	}
}

// diffSchemas compares vendored against upstream and reports what upstream
// has that vendored doesn't (added), what vendored has that upstream doesn't
// (removed), and field/enum-value differences on types present in both.
func diffSchemas(vendored, upstream *ast.Schema) *drift {
	report := &drift{}

	for name, upstreamType := range upstream.Types {
		if skipType(upstreamType) {
			continue
		}
		vendoredType, ok := vendored.Types[name]
		if !ok {
			report.typesAdded = append(report.typesAdded, name)
			continue
		}
		diffType(report, name, vendoredType, upstreamType)
	}

	for name, vendoredType := range vendored.Types {
		if skipType(vendoredType) {
			continue
		}
		if upstreamType, ok := upstream.Types[name]; !ok || skipType(upstreamType) {
			report.typesRemoved = append(report.typesRemoved, name)
		}
	}

	return report
}

func skipType(def *ast.Definition) bool {
	return def.BuiltIn || strings.HasPrefix(def.Name, "__")
}

func diffType(report *drift, typeName string, vendored, upstream *ast.Definition) {
	upstreamFields := fieldsByName(upstream.Fields)
	vendoredFields := fieldsByName(vendored.Fields)

	for fieldName, upstreamField := range upstreamFields {
		vendoredField, ok := vendoredFields[fieldName]
		if !ok {
			report.fieldsAdded = append(report.fieldsAdded, typeName+"."+fieldName)
			continue
		}
		if upstreamField.Type.String() != vendoredField.Type.String() {
			report.fieldTypesChanged = append(report.fieldTypesChanged, fmt.Sprintf(
				"%s.%s: %s -> %s", typeName, fieldName, vendoredField.Type.String(), upstreamField.Type.String(),
			))
		}
	}
	for fieldName := range vendoredFields {
		if _, ok := upstreamFields[fieldName]; !ok {
			report.fieldsRemoved = append(report.fieldsRemoved, typeName+"."+fieldName)
		}
	}

	upstreamValues := enumValuesByName(upstream.EnumValues)
	vendoredValues := enumValuesByName(vendored.EnumValues)
	for value := range upstreamValues {
		if _, ok := vendoredValues[value]; !ok {
			report.enumValuesAdded = append(report.enumValuesAdded, typeName+"."+value)
		}
	}
	for value := range vendoredValues {
		if _, ok := upstreamValues[value]; !ok {
			report.enumValuesRemoved = append(report.enumValuesRemoved, typeName+"."+value)
		}
	}
}

func fieldsByName(fields ast.FieldList) map[string]*ast.FieldDefinition {
	byName := make(map[string]*ast.FieldDefinition, len(fields))
	for _, field := range fields {
		byName[field.Name] = field
	}
	return byName
}

func enumValuesByName(values ast.EnumValueList) map[string]*ast.EnumValueDefinition {
	byName := make(map[string]*ast.EnumValueDefinition, len(values))
	for _, value := range values {
		byName[value.Name] = value
	}
	return byName
}

func exitError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
