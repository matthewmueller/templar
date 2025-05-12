package cli

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ/generator"
	"github.com/a-h/templ/parser/v2"
	"github.com/livebud/cli"
	"golang.org/x/sync/errgroup"
)

func Run() int {
	ctx := context.Background()
	args := os.Args[1:]
	if err := Parse(ctx, args...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func Parse(ctx context.Context, args ...string) error {
	return New().Parse(ctx, args...)
}

func New() *CLI {
	return &CLI{
		Dir:    ".",
		Env:    os.Environ(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
}

type CLI struct {
	Dir    string
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func (c *CLI) Parse(ctx context.Context, args ...string) error {
	cli := cli.New("templar", "templ extensions")

	{ // templar generate [paths...]
		in := &Generate{}
		cmd := cli.Command("generate", "generate code from templates")
		cmd.Flag("out", "output directory").Short('o').String(&in.OutDir).Default(".")
		cmd.Args("paths", "paths to templates").Strings(&in.Paths)
		cmd.Run(func(ctx context.Context) error { return c.Generate(ctx, in) })
	}

	return cli.Parse(ctx, args...)
}

// Walk the directory tree and find all .templ files
func (c *CLI) resolveTree(dir, path string) (resolved []string, sourceDir string, err error) {
	sourceDir = filepath.Join(dir, path)
	err = filepath.WalkDir(sourceDir, func(p string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		base := filepath.Base(p)
		if base == "." {
			return nil
		} else if de.IsDir() {
			if p == base {
				return nil
			}
			// Skip directories that start with _ or .
			if base[0] == '_' || base[0] == '.' {
				return filepath.SkipDir
			}
			return nil
		}
		// Skip files that start with _ or .
		if base[0] == '_' || base[0] == '.' {
			return nil
		}
		if filepath.Ext(p) == ".templ" {
			resolved = append(resolved, p)
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	return resolved, sourceDir, nil
}

func (c *CLI) resolveDir(dir, path string) (resolved []string, sourceDir string, err error) {
	sourceDir = filepath.Join(dir, path)
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, "", err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		base := file.Name()
		// Skip files that start with _ or .
		if base[0] == '_' || base[0] == '.' {
			continue
		}
		if filepath.Ext(base) == ".templ" {
			resolved = append(resolved, filepath.Join(dir, path, base))
		}
	}
	return resolved, sourceDir, nil
}

const ellipses = string(filepath.Separator) + "..."

func (c *CLI) resolvePath(dir string, path string) ([]string, string, error) {
	sourceDir := filepath.Join(dir, path)
	if strings.HasSuffix(path, ellipses) {
		return c.resolveTree(dir, strings.TrimSuffix(path, ellipses))
	}
	base := filepath.Base(path)
	if base[0] == '_' || base[0] == '.' {
		return nil, "", nil
	}
	stat, err := os.Stat(sourceDir)
	if err != nil {
		return nil, "", err
	}
	if stat.IsDir() {
		return c.resolveDir(dir, path)
	}
	if filepath.Ext(path) == ".templ" {
		return []string{path}, filepath.Dir(path), nil
	}
	return nil, sourceDir, nil
}

func (c *CLI) resolve(dir string, path string) (matches []string, sourceDir string, err error) {
	resolved, sourceDir, err := c.resolvePath(dir, path)
	if err != nil {
		return nil, "", err
	}
	matches = append(matches, resolved...)
	return matches, sourceDir, nil
}

type Generate struct {
	Paths  []string
	OutDir string
}

func resolveDir(dir string, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(dir, path)
}

func (c *CLI) Generate(ctx context.Context, in *Generate) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, path := range in.Paths {
		path := filepath.Clean(path)
		eg.Go(func() error {
			return c.generate(ctx, in, path)
		})
	}
	return eg.Wait()
}

func (c *CLI) generate(ctx context.Context, in *Generate, path string) error {
	matches, sourceDir, err := c.resolve(c.Dir, path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	outDir := resolveDir(c.Dir, in.OutDir)
	eg, ctx := errgroup.WithContext(ctx)
	for _, path := range matches {
		path := path
		eg.Go(func() error {
			return c.generateFile(ctx, sourceDir, outDir, path)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func extless(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}

func (c *CLI) generateFile(_ context.Context, sourceDir, targetDir, path string) error {
	code, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	tf, err := parser.ParseString(string(code))
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	generated := new(bytes.Buffer)
	if _, err := generator.Generate(tf, generated, generator.WithFileName(path)); err != nil {
		return fmt.Errorf("error generating code: %w", err)
	}

	formatted, err := format.Source(generated.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting code: %w", err)
	}

	targetPath := path
	if filepath.Clean(targetDir) != "." {
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("error getting relative path: %w", err)
		}
		targetPath = filepath.Join(targetDir, relPath)
	}
	outDir, outPath := filepath.Split(targetPath)
	outBase := extless(outPath) + "_templ.go"
	outPath = filepath.Join(outDir, outBase)

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", outDir, err)
	}
	if err := os.WriteFile(outPath, formatted, 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", outPath, err)
	}

	return nil
}
