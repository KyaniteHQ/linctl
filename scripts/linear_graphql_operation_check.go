//go:build ignore

// Package main validates local Linear GraphQL operations against the upstream
// Linear SDK schema without mutating generated code or vendored schema files.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
)

func main() {
	upstreamDir := flag.String("upstream", "", "required upstream Linear repo checkout")
	operationsPattern := flag.String("operations", "internal/client/operations/*.graphql", "local GraphQL operations glob")
	flag.Parse()

	if *upstreamDir == "" {
		exitError(fmt.Errorf("missing -upstream: set it to a Linear SDK checkout path"))
	}

	schemaPath := filepath.Join(*upstreamDir, "packages/sdk/src/schema.graphql")
	schema, err := validator.LoadSchema(
		validator.Prelude,
		&ast.Source{
			Name:  schemaPath,
			Input: string(mustReadFile(schemaPath)),
		},
	)
	if err != nil {
		exitError(fmt.Errorf("load upstream schema: %w", err))
	}

	operations, err := loadOperations(*operationsPattern)
	if err != nil {
		exitError(err)
	}
	if gqlErrors := validator.Validate(schema, operations); len(gqlErrors) > 0 {
		exitError(fmt.Errorf("validate operations against %s: %v", schemaPath, gqlErrors))
	}
}

func loadOperations(pattern string) (*ast.QueryDocument, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob operations %s: %w", pattern, err)
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no operation files match %s", pattern)
	}

	merged := &ast.QueryDocument{}
	for _, path := range paths {
		source := mustReadFile(path)
		document, err := parser.ParseQuery(&ast.Source{Name: path, Input: string(source)})
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		merged.Operations = append(merged.Operations, document.Operations...)
		merged.Fragments = append(merged.Fragments, document.Fragments...)
	}

	return merged, nil
}

func mustReadFile(path string) []byte {
	source, err := os.ReadFile(path)
	if err != nil {
		exitError(fmt.Errorf("read %s: %w", path, err))
	}
	return source
}

func exitError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
