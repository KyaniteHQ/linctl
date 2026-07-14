//go:build ignore

// exportgen exposes only the lower-case genqlient declarations referenced by
// the parent client package. The genqlient output itself remains untouched.
package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	goformat "go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	gofumpt "mvdan.cc/gofumpt/format"
)

type declaration struct {
	kind     token.Token
	function *ast.FuncDecl
	typeSpec *ast.TypeSpec
}

type generatedPackage struct {
	name         string
	declarations map[string]declaration
	imports      map[string]*ast.ImportSpec
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: exportgen GENERATED OUTPUT PARENT_DIR GQL_IMPORT_PATH")
		os.Exit(2)
	}

	if err := generateExports(args[0], args[1], args[2], args[3]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func generateExports(generatedPath string, outputPath string, parentDir string, gqlImportPath string) error {
	generated, err := parseGeneratedPackage(generatedPath)
	if err != nil {
		return err
	}
	required, err := referencedSelectors(parentDir, gqlImportPath)
	if err != nil {
		return err
	}

	declarations, imports, err := exportDeclarations(generated, required)
	if err != nil {
		return err
	}
	content, err := formatExports(generated.name, declarations, imports)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, content, 0o644); err != nil {
		return fmt.Errorf("write exports %s: %w", outputPath, err)
	}

	return nil
}

func parseGeneratedPackage(filePath string) (generatedPackage, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.SkipObjectResolution)
	if err != nil {
		return generatedPackage{}, fmt.Errorf("parse generated client %s: %w", filePath, err)
	}

	generated := generatedPackage{
		name:         file.Name.Name,
		declarations: make(map[string]declaration),
		imports:      make(map[string]*ast.ImportSpec),
	}
	for _, importSpec := range file.Imports {
		localName, err := importLocalName(importSpec)
		if err != nil {
			return generatedPackage{}, fmt.Errorf("generated client import: %w", err)
		}
		if localName == "_" || localName == "." {
			continue
		}
		generated.imports[localName] = importSpec
	}
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			if decl.Recv == nil {
				if err := addDeclaration(generated.declarations, decl.Name.Name, declaration{
					kind:     token.FUNC,
					function: decl,
				}); err != nil {
					return generatedPackage{}, err
				}
			}
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					if err := addDeclaration(generated.declarations, spec.Name.Name, declaration{
						kind:     token.TYPE,
						typeSpec: spec,
					}); err != nil {
						return generatedPackage{}, err
					}
				case *ast.ValueSpec:
					for _, name := range spec.Names {
						if err := addDeclaration(generated.declarations, name.Name, declaration{kind: decl.Tok}); err != nil {
							return generatedPackage{}, err
						}
					}
				}
			}
		}
	}

	return generated, nil
}

func addDeclaration(declarations map[string]declaration, name string, decl declaration) error {
	if _, exists := declarations[name]; exists {
		return fmt.Errorf("generated client declares %s more than once", name)
	}
	declarations[name] = decl
	return nil
}

func referencedSelectors(parentDir string, gqlImportPath string) ([]string, error) {
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, fmt.Errorf("read parent source directory %s: %w", parentDir, err)
	}

	required := make(map[string]struct{})
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		filePath := filepath.Join(parentDir, entry.Name())
		file, err := parser.ParseFile(token.NewFileSet(), filePath, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parse parent source %s: %w", filePath, err)
		}
		aliases, err := gqlImportAliases(file, gqlImportPath)
		if err != nil {
			return nil, fmt.Errorf("resolve parent import in %s: %w", filePath, err)
		}
		if len(aliases) == 0 {
			continue
		}
		ast.Inspect(file, func(node ast.Node) bool {
			selector, ok := node.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			qualifier, ok := selector.X.(*ast.Ident)
			if !ok || qualifier.Obj != nil {
				return true
			}
			if _, ok := aliases[qualifier.Name]; ok {
				required[selector.Sel.Name] = struct{}{}
			}
			return true
		})
	}

	names := make([]string, 0, len(required))
	for name := range required {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

func gqlImportAliases(file *ast.File, gqlImportPath string) (map[string]struct{}, error) {
	aliases := make(map[string]struct{})
	for _, importSpec := range file.Imports {
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			return nil, err
		}
		if importPath != gqlImportPath {
			continue
		}
		localName := path.Base(importPath)
		if importSpec.Name != nil {
			localName = importSpec.Name.Name
		}
		if localName == "." {
			return nil, errors.New("dot imports of the generated GraphQL package are unsupported")
		}
		if localName != "_" {
			aliases[localName] = struct{}{}
		}
	}
	return aliases, nil
}

func exportDeclarations(
	generated generatedPackage,
	required []string,
) ([]ast.Decl, map[string]*ast.ImportSpec, error) {
	declarations := make([]ast.Decl, 0, len(required))
	imports := make(map[string]*ast.ImportSpec)
	for _, exportedName := range required {
		direct, directExists := generated.declarations[exportedName]
		originalName, canReverse := unexportedName(exportedName)
		original, originalExists := generated.declarations[originalName]
		if directExists && originalExists && canReverse {
			return nil, nil, fmt.Errorf(
				"cannot expose %s: generated declarations %s and %s collide",
				exportedName,
				exportedName,
				originalName,
			)
		}
		if directExists {
			if !ast.IsExported(exportedName) {
				return nil, nil, fmt.Errorf("referenced generated symbol %s is not exported", exportedName)
			}
			_ = direct
			continue
		}
		if !canReverse || !originalExists || ast.IsExported(originalName) {
			return nil, nil, fmt.Errorf("referenced generated symbol %s cannot be resolved", exportedName)
		}

		var decl ast.Decl
		switch original.kind {
		case token.FUNC:
			wrapper, err := functionWrapper(exportedName, originalName, original.function)
			if err != nil {
				return nil, nil, err
			}
			decl = wrapper
			if err := collectSignatureImports(wrapper.Type, generated.imports, imports); err != nil {
				return nil, nil, err
			}
		case token.TYPE:
			decl = &ast.GenDecl{Tok: token.TYPE, Specs: []ast.Spec{&ast.TypeSpec{
				Name:   ast.NewIdent(exportedName),
				Assign: token.Pos(1),
				Type:   ast.NewIdent(originalName),
			}}}
		case token.CONST, token.VAR:
			decl = &ast.GenDecl{Tok: original.kind, Specs: []ast.Spec{&ast.ValueSpec{
				Names:  []*ast.Ident{ast.NewIdent(exportedName)},
				Values: []ast.Expr{ast.NewIdent(originalName)},
			}}}
		default:
			return nil, nil, fmt.Errorf("unsupported declaration kind for %s", exportedName)
		}
		declarations = append(declarations, decl)
	}

	return declarations, imports, nil
}

func unexportedName(exported string) (string, bool) {
	if !strings.HasPrefix(exported, "X") || len(exported) == 1 {
		return "", false
	}
	r, size := utf8.DecodeRuneInString(exported[1:])
	return string(unicode.ToLower(r)) + exported[1+size:], true
}

func functionWrapper(exportedName string, originalName string, original *ast.FuncDecl) (*ast.FuncDecl, error) {
	if original == nil || original.Type == nil {
		return nil, fmt.Errorf("generated function %s has no signature", originalName)
	}
	arguments := make([]ast.Expr, 0)
	variadic := false
	if original.Type.Params != nil {
		for fieldIndex, field := range original.Type.Params.List {
			if len(field.Names) == 0 {
				return nil, fmt.Errorf("generated function %s has an unnamed parameter", originalName)
			}
			for _, name := range field.Names {
				arguments = append(arguments, ast.NewIdent(name.Name))
			}
			if _, ok := field.Type.(*ast.Ellipsis); ok {
				if fieldIndex != len(original.Type.Params.List)-1 || len(field.Names) != 1 {
					return nil, fmt.Errorf("generated function %s has an invalid variadic parameter", originalName)
				}
				variadic = true
			}
		}
	}
	call := &ast.CallExpr{Fun: ast.NewIdent(originalName), Args: arguments}
	if variadic {
		call.Ellipsis = token.Pos(1)
	}
	var statement ast.Stmt = &ast.ExprStmt{X: call}
	if original.Type.Results != nil && len(original.Type.Results.List) > 0 {
		statement = &ast.ReturnStmt{Results: []ast.Expr{call}}
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent(exportedName),
		Type: original.Type,
		Body: &ast.BlockStmt{List: []ast.Stmt{statement}},
	}, nil
}

func collectSignatureImports(
	signature *ast.FuncType,
	available map[string]*ast.ImportSpec,
	selected map[string]*ast.ImportSpec,
) error {
	var collectErr error
	ast.Inspect(signature, func(node ast.Node) bool {
		if collectErr != nil {
			return false
		}
		selector, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		qualifier, ok := selector.X.(*ast.Ident)
		if !ok {
			return true
		}
		importSpec, ok := available[qualifier.Name]
		if !ok {
			collectErr = fmt.Errorf("generated signature references unresolved import %s", qualifier.Name)
			return false
		}
		selected[qualifier.Name] = importSpec
		return true
	})
	return collectErr
}

func formatExports(packageName string, declarations []ast.Decl, imports map[string]*ast.ImportSpec) ([]byte, error) {
	file := &ast.File{Name: ast.NewIdent(packageName)}
	if len(imports) > 0 {
		aliases := make([]string, 0, len(imports))
		for alias := range imports {
			aliases = append(aliases, alias)
		}
		sort.Strings(aliases)
		specs := make([]ast.Spec, 0, len(aliases))
		for _, alias := range aliases {
			original := imports[alias]
			clone := &ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: original.Path.Value}}
			defaultName, err := importLocalName(clone)
			if err != nil {
				return nil, err
			}
			if alias != defaultName {
				clone.Name = ast.NewIdent(alias)
			}
			specs = append(specs, clone)
		}
		file.Decls = append(file.Decls, &ast.GenDecl{Tok: token.IMPORT, Specs: specs})
	}
	file.Decls = append(file.Decls, declarations...)

	var output bytes.Buffer
	output.WriteString("// Code generated by exportgen; DO NOT EDIT.\n\n")
	if err := goformat.Node(&output, token.NewFileSet(), file); err != nil {
		return nil, fmt.Errorf("format generated exports: %w", err)
	}
	formatted, err := gofumpt.Source(output.Bytes(), gofumpt.Options{})
	if err != nil {
		return nil, fmt.Errorf("gofumpt generated exports: %w", err)
	}
	return formatted, nil
}

func importLocalName(importSpec *ast.ImportSpec) (string, error) {
	if importSpec.Name != nil {
		return importSpec.Name.Name, nil
	}
	importPath, err := strconv.Unquote(importSpec.Path.Value)
	if err != nil {
		return "", err
	}
	return path.Base(importPath), nil
}
