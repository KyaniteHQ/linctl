//go:build ignore

// Package main checks that every Go struct genqlient binds to a Linear GraphQL
// input type sends only fields the vendored schema defines for that input.
// genqlient and the operation check never read bound structs, so a schema
// refresh that removes an input field passes both and fails only at Linear.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	gqlast "github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator"
	"gopkg.in/yaml.v2"
)

type genqlientConfig struct {
	Schema   string `yaml:"schema"`
	Bindings map[string]struct {
		Type string `yaml:"type"`
	} `yaml:"bindings"`
}

func main() {
	configPath := flag.String("config", "internal/client/genqlient.yaml", "genqlient config with the bindings")
	modelsPath := flag.String("models", "internal/client/internal/gqlmodel/inputs.go", "Go file with the bound structs")
	flag.Parse()

	config := genqlientConfig{}
	if err := yaml.Unmarshal(mustReadFile(*configPath), &config); err != nil {
		exitError(fmt.Errorf("parse %s: %w", *configPath, err))
	}
	schemaPath := filepath.Join(filepath.Dir(*configPath), config.Schema)
	schema, err := validator.LoadSchema(
		validator.Prelude,
		&gqlast.Source{Name: schemaPath, Input: string(mustReadFile(schemaPath))},
	)
	if err != nil {
		exitError(fmt.Errorf("load schema %s: %w", schemaPath, err))
	}
	structFields, err := readStructFields(*modelsPath)
	if err != nil {
		exitError(err)
	}

	var problems []string
	for _, inputName := range sortedKeys(config.Bindings) {
		definition := schema.Types[inputName]
		if definition == nil {
			problems = append(problems, fmt.Sprintf("binding %s names no type in %s", inputName, schemaPath))
			continue
		}
		if definition.Kind != gqlast.InputObject {
			continue
		}
		goType := config.Bindings[inputName].Type
		structName := goType[strings.LastIndex(goType, ".")+1:]
		fields, ok := structFields[structName]
		if !ok {
			problems = append(problems, fmt.Sprintf("binding %s: struct %s is not in %s", inputName, structName, *modelsPath))
			continue
		}
		for _, field := range fields {
			if definition.Fields.ForName(field) == nil {
				problems = append(problems, fmt.Sprintf(
					"%s sends %q, which input %s in %s does not define", structName, field, inputName, schemaPath,
				))
			}
		}
	}
	if len(problems) > 0 {
		exitError(fmt.Errorf("bound input structs disagree with the schema:\n  %s", strings.Join(problems, "\n  ")))
	}
}

func readStructFields(path string) (map[string][]string, error) {
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	structs := map[string][]string{}
	ast.Inspect(file, func(node ast.Node) bool {
		spec, ok := node.(*ast.TypeSpec)
		if !ok {
			return true
		}
		structType, ok := spec.Type.(*ast.StructType)
		if !ok {
			return false
		}
		var fields []string
		for _, field := range structType.Fields.List {
			if field.Tag == nil {
				continue
			}
			tag, err := strconv.Unquote(field.Tag.Value)
			if err != nil {
				continue
			}
			name, _, _ := strings.Cut(reflect.StructTag(tag).Get("json"), ",")
			if name != "" && name != "-" {
				fields = append(fields, name)
			}
		}
		structs[spec.Name.Name] = fields
		return false
	})
	return structs, nil
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
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
