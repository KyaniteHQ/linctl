package gql

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const fixtureImportPath = "example.com/generated/gql"

func TestExportgenSelectsReferencedWrappersTypesAndValues(t *testing.T) {
	generated := `package fixture

import (
	"context"
	graph "example.com/graphql"
)

type response struct{ Value string }
type unusedResponse struct{}
const operationName = "lowerCaseOperation"
var cachedValue = "cached"

func lowerCaseOperation(ctx context.Context, client graph.Client, ids ...string) (*response, error) {
	return nil, nil
}

func unusedOperation() *unusedResponse { return nil }
`
	parent := `package parent

import gql "example.com/generated/gql"

var _ gql.XResponse
var _ = gql.XOperationName
var _ = gql.XCachedValue

func use() { _, _ = gql.XLowerCaseOperation(nil, nil, "id") }
`
	paths := writeExportgenWorkspace(t, generated, parent)
	original := readFile(t, paths.generated)
	runExportgen(t, paths, true)

	got := readFile(t, paths.output)
	for _, want := range []string{
		`graph "example.com/graphql"`,
		`"context"`,
		"type XResponse = response",
		`const XOperationName = operationName`,
		`var XCachedValue = cachedValue`,
		"func XLowerCaseOperation(ctx context.Context, client graph.Client, ids ...string) (*response, error)",
		"return lowerCaseOperation(ctx, client, ids...)",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("generated exports missing %q:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"XUnusedResponse", "XUnusedOperation"} {
		if strings.Contains(got, unwanted) {
			t.Fatalf("generated exports unexpectedly contain %q:\n%s", unwanted, got)
		}
	}
	if after := readFile(t, paths.generated); after != original {
		t.Fatalf("export generation changed genqlient output:\nbefore:\n%s\nafter:\n%s", original, after)
	}
}

func TestExportgenRejectsUnresolvedSelectorWithoutWriting(t *testing.T) {
	paths := writeExportgenWorkspace(t, "package fixture\n", `package parent
import gql "example.com/generated/gql"
var _ gql.XMissing
`)
	requireWriteFile(t, paths.output, "sentinel")

	output := runExportgen(t, paths, false)
	if !strings.Contains(output, "referenced generated symbol XMissing cannot be resolved") {
		t.Fatalf("unexpected unresolved-symbol output: %s", output)
	}
	if got := readFile(t, paths.output); got != "sentinel" {
		t.Fatalf("failed generation changed output: %q", got)
	}
}

func TestExportgenRejectsCapitalizationCollisionWithoutWriting(t *testing.T) {
	paths := writeExportgenWorkspace(t, `package fixture
type response struct{}
type XResponse struct{}
`, `package parent
import gql "example.com/generated/gql"
var _ gql.XResponse
`)
	requireWriteFile(t, paths.output, "sentinel")

	output := runExportgen(t, paths, false)
	if !strings.Contains(output, "generated declarations XResponse and response collide") {
		t.Fatalf("unexpected collision output: %s", output)
	}
	if got := readFile(t, paths.output); got != "sentinel" {
		t.Fatalf("collision changed output: %q", got)
	}
}

func TestExportgenResolvesDefaultImportAndIgnoresShadowedQualifier(t *testing.T) {
	paths := writeExportgenWorkspace(t, `package fixture
type response struct{}
`, `package parent
import "example.com/generated/gql"
var _ gql.XResponse
func shadow() {
	gql := struct{ XMissing string }{}
	_ = gql.XMissing
}
`)

	runExportgen(t, paths, true)
	got := readFile(t, paths.output)
	if !strings.Contains(got, "type XResponse = response") || strings.Contains(got, "XMissing") {
		t.Fatalf("unexpected selector resolution:\n%s", got)
	}
}

func TestExportgenIsIdempotentAndReproducible(t *testing.T) {
	paths := writeExportgenWorkspace(t, `package fixture
func operation(value string) string { return value }
`, `package parent
import gql "example.com/generated/gql"
var _ = gql.XOperation("value")
`)

	runExportgen(t, paths, true)
	first := readFile(t, paths.output)
	runExportgen(t, paths, true)
	if second := readFile(t, paths.output); second != first {
		t.Fatalf("second generation changed output:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

type exportgenPaths struct {
	generated string
	output    string
	parent    string
}

func writeExportgenWorkspace(t *testing.T, generated string, parent string) exportgenPaths {
	t.Helper()
	root := t.TempDir()
	parentDir := filepath.Join(root, "parent")
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	paths := exportgenPaths{
		generated: filepath.Join(root, "generated.go"),
		output:    filepath.Join(root, "exports.go"),
		parent:    parentDir,
	}
	requireWriteFile(t, paths.generated, generated)
	requireWriteFile(t, filepath.Join(parentDir, "use.go"), parent)
	return paths
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func requireWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runExportgen(t *testing.T, paths exportgenPaths, wantSuccess bool) string {
	t.Helper()
	command := exec.Command(
		"go", "run", "./exportgen.go", "--",
		paths.generated,
		paths.output,
		paths.parent,
		fixtureImportPath,
	)
	output, err := command.CombinedOutput()
	if wantSuccess && err != nil {
		t.Fatalf("exportgen failed: %v\n%s", err, output)
	}
	if !wantSuccess && err == nil {
		t.Fatalf("exportgen unexpectedly succeeded")
	}
	return string(output)
}
