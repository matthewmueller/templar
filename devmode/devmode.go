package devmode

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

// Transform turns the generated _templ.go file into a file that works better
// with a file watcher.
func Transform(path string, src []byte) ([]byte, error) {
	transformer := &Transformer{
		Runtime: "github.com/matthewmueller/templar/templdev",
	}
	return transformer.Transform(path, src)
}

// Transformer for devmode
type Transformer struct {
	// Runtime is the import path for the templdev package.
	Runtime string
}

// Transform turns the generated _templ.go file into a file that works better
// with a file watcher.
func (t *Transformer) Transform(path string, src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("devmode: failed to parse file: %w", err)
	}

	modified := false

	// Replace templruntime.WriteString with templdev.WriteString and replace the
	// literal string with an empty string so that during generation, the Go file
	// is modified only when it's a non text change.
	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) != 3 {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "WriteString" {
			return true
		}

		pkg, ok := sel.X.(*ast.Ident)
		if !ok || pkg.Name != "templruntime" {
			return true
		}

		// Replace with templdev
		pkg.Name = "templdev"
		call.Args[2] = &ast.BasicLit{
			Kind:  token.STRING,
			Value: `""`,
		}
		modified = true
		return true
	})

	if !modified {
		return src, nil
	}

	// Add templdev import if missing
	if !hasImport(file, t.Runtime) {
		astutil.AddImport(fset, file, t.Runtime)
	}

	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, file)
	if err != nil {
		return nil, fmt.Errorf("devmode: failed to print file %q: %w", path, err)
	}

	return buf.Bytes(), nil
}

func hasImport(f *ast.File, path string) bool {
	for _, imp := range f.Imports {
		if strings.Trim(imp.Path.Value, `"`) == path {
			return true
		}
	}
	return false
}
