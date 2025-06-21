package testutil

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ/generator"
	"github.com/a-h/templ/parser/v2"
	"github.com/hexops/valast"
	"github.com/matthewmueller/diff"
)

var update = flag.Bool("update", false, "update the expected output")

func TestData(dir string) (dirs []string, err error) {
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

// Golden compares the actual output with the expected output. Use -update to
// update the expected output.
func Golden(path string, actual any) error {
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

func GoPath(templPath string) string {
	extless := strings.TrimSuffix(templPath, ".templ")
	return extless + "_templ.go"
}

func TxtPath(templPath string) string {
	extless := strings.TrimSuffix(templPath, ".templ")
	return extless + "_templ.txt"
}

func Generate(templPath, templCode string, visitors ...parser.Visitor) ([]byte, string, error) {
	// Parse the template code
	templAst, err := parser.ParseString(templCode)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse %s: %w", templPath, err)
	}

	// Run visitors on the AST
	for _, visitor := range visitors {
		if err := templAst.Visit(visitor); err != nil {
			return nil, "", fmt.Errorf("failed to visit %s: %w", templPath, err)
		}
	}

	// Generate the Go code from the template AST
	generated := new(bytes.Buffer)
	out, err := generator.Generate(templAst, generated, generator.WithFileName(templPath))
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate %s: %w", templPath, err)
	}
	literals := strings.Join(out.Literals, "\n")

	// Format the generated Go code
	formatted, err := format.Source(generated.Bytes())
	if err != nil {
		return nil, "", fmt.Errorf("failed to format %s: %w", templPath, err)
	}
	return formatted, literals, nil
}
