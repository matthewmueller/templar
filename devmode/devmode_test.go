package devmode_test

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a-h/templ/generator"
	"github.com/a-h/templ/parser/v2"
	"github.com/hexops/valast"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
	"github.com/matthewmueller/templar/devmode"
)

var update = flag.Bool("update", false, "update the expected output")

func testData(dir string) (dirs []string, err error) {
	des, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, d := range des {
		if !d.IsDir() || strings.HasPrefix(d.Name(), ".") || strings.HasPrefix(d.Name(), "_") {
			continue
		}
		dirs = append(dirs, filepath.Join(dir, d.Name()))
	}
	return dirs, nil
}

func genGoPath(templPath string) string {
	extless := strings.TrimSuffix(templPath, ".templ")
	return extless + "_templ.go"
}

func genTxtPath(templPath string) string {
	extless := strings.TrimSuffix(templPath, ".templ")
	return extless + "_templ.txt"
}

// Test compares the actual output with the expected output. Use -update to
// update the expected output.
func golden(path string, actual any) error {
	expect, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if *update {
				if err := updateFile(path, actual); err != nil {
					return fmt.Errorf("failed to update file %s: %v", path, err)
				}
				return nil
			} else {
				return diff.String(formatString(actual), "")
			}
		}
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}
	if *update {
		if err := updateFile(path, actual); err != nil {
			return fmt.Errorf("failed to update file %s: %v", path, err)
		}
		return nil
	}
	return diff.String(formatString(actual), string(expect))
}

func formatString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return valast.StringWithOptions(v, &valast.Options{
		Unqualify: true,
	})
}

func updateFile(path string, actual any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(formatString(actual)), 0644)
}

func generate(templPath, templCode string) ([]byte, string, error) {
	templAst, err := parser.ParseString(templCode)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse %s: %w", templPath, err)
	}
	generated := new(bytes.Buffer)
	out, err := generator.Generate(templAst, generated, generator.WithFileName(templPath))
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate %s: %w", templPath, err)
	}
	literals := strings.Join(out.Literals, "\n")
	formatted, err := format.Source(generated.Bytes())
	if err != nil {
		return nil, "", fmt.Errorf("failed to format %s: %w", templPath, err)
	}
	return formatted, literals, nil
}

func TestData(t *testing.T) {
	is := is.New(t)
	dirs, err := testData("testdata")
	is.NoErr(err)
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			is := is.New(t)
			templPaths, err := filepath.Glob(filepath.Join(dir, "*.templ"))
			is.NoErr(err)
			for _, templPath := range templPaths {
				templCode, err := os.ReadFile(templPath)
				is.NoErr(err)
				generated, literals, err := generate(templPath, string(templCode))
				is.NoErr(err)

				// Patch with devmode
				modified, err := devmode.Transform(templPath, generated)
				is.NoErr(err)

				is.NoErr(golden(genGoPath(templPath), string(modified)))
				is.NoErr(golden(genTxtPath(templPath), literals))
			}
		})
	}
}

func TestChanges(t *testing.T) {
	is := is.New(t)

	const one = `package testcall

templ showAll() {
	@a()
	@b(c("C"))
	@d()
	@showOne(e())
	@wrapChildren() {
		<div>Child content</div>
	}
}

templ a() {
	<div>A</div>
}

templ b(child templ.Component) {
	<div>B</div>
	@child
}

templ c(text string) {
	<div>{ text }</div>
}

templ d() {
	<div>Legacy call style</div>
}

templ e() {
	e
}

templ showOne(component templ.Component) {
	<div>
		@component
	</div>
}

templ wrapChildren() {
	<div id="wrapper">
		{ children... }
	</div>
}
`

	const two = `package testcall

templ showAll() {
	@a()
	@b(c("C"))
	@d()
	@showOne(e())
	@wrapChildren() {
		<div>Child content!</div>
	}
}

templ a() {
	<div>A</div>
}

templ b(child templ.Component) {
	<div>B</div>
	@child
}

templ c(text string) {
	<div>{ text }<a>hi</a></div>
}

templ d() {
	<div>Legacy call style!!!</div>
}

templ e() {
	e!!!
}

templ showOne(component templ.Component) {
	<div class="nice">
		@component
	</div>
}

templ wrapChildren() {
	<div id="wrapperz">
		{ children... }
		!
	</div>
}
`

	oneGen, oneLits, err := generate("testcall.templ", one)
	is.NoErr(err)

	twoGen, twoLits, err := generate("testcall.templ", two)
	is.NoErr(err)

	is.True(!bytes.Equal(oneGen, twoGen)) // generated code should be different
	is.True(oneLits != twoLits)           // literals should be different

	oneMod, err := devmode.Transform("testcall.templ", oneGen)
	is.NoErr(err)

	twoMod, err := devmode.Transform("testcall.templ", twoGen)
	is.NoErr(err)

	is.True(bytes.Equal(oneMod, twoMod)) // modified code should be the same
}
