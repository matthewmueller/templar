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
				templAst, err := parser.Parse(templPath)
				is.NoErr(err)
				generated := new(bytes.Buffer)
				out, err := generator.Generate(templAst, generated, generator.WithFileName(templPath))
				is.NoErr(err)
				literals := strings.Join(out.Literals, "\n")

				formatted, err := format.Source(generated.Bytes())
				is.NoErr(err)

				// Patch with devmode
				// modified, err := devmode.Transform(templPath, formatted)
				// is.NoErr(err)

				is.NoErr(golden(genGoPath(templPath), string(formatted)))
				is.NoErr(golden(genTxtPath(templPath), literals))
			}
		})
	}
}
