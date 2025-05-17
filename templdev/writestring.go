// templdev is a runtime package intended to be used with the devmode package to
// allow for live reloading of templates. It differs from the original templ
// runtime package, where it doesn't rely on environment variables. It also
// doesn't change the Go file if it's just a text change. This makes it better
// suited to integrate into build systems and third-party file watchers.
package templdev

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cache = sync.Map{}

type entry struct {
	ModTime time.Time
	Strings []string
}

func openLiterals(txtPath string) ([]string, error) {
	val, ok := cache.Load(txtPath)
	if !ok {
		return cacheLiterals(txtPath)
	}

	entry, ok := val.(entry)
	if !ok {
		return cacheLiterals(txtPath)
	}

	// This is a lot of stat-ing for large templates
	// TODO: create a gen-id to only do it once per generate
	stat, err := os.Stat(txtPath)
	if err != nil {
		return nil, fmt.Errorf("templdev: failed to stat %s: %w", txtPath, err)
	}

	if stat.ModTime().After(entry.ModTime) {
		return cacheLiterals(txtPath)
	}

	return entry.Strings, nil
}

func cacheLiterals(txtPath string) ([]string, error) {
	txtFile, err := os.Open(txtPath)
	if err != nil {
		return nil, fmt.Errorf("templdev: failed to open %s: %w", txtPath, err)
	}
	defer txtFile.Close()

	stat, err := txtFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("templdev: failed to stat %s: %w", txtPath, err)
	}

	all, err := io.ReadAll(txtFile)
	if err != nil {
		return nil, fmt.Errorf("templdev: failed to read %s: %w", txtPath, err)
	}

	literals := strings.Split(string(all), "\n")
	cache.Store(txtPath, entry{
		ModTime: stat.ModTime(),
		Strings: literals,
	})

	return literals, nil
}

func txtFilePath(goPath string) string {
	extless := strings.TrimSuffix(goPath, ".go")
	return extless + ".txt"
}

// WriteString writes the string at the index in the _templ.txt file. This is
// intended to be used with the devmode package to allow for live reloading of
// templates. It is not intended to be used in production code.
func WriteString(w io.Writer, lineNum int, _ string) error {
	_, path, _, _ := runtime.Caller(1)
	if !strings.HasSuffix(path, "_templ.go") {
		return errors.New("templdev: attempt to use WriteString from a non templ file")
	}
	txtPath := txtFilePath(path)
	literals, err := openLiterals(txtPath)
	if err != nil {
		return fmt.Errorf("templdev: failed to get watched strings for %q: %w", path, err)
	}
	if lineNum > len(literals) {
		return fmt.Errorf("templ: failed to find line %d in %s", lineNum, txtPath)
	}
	s, err := strconv.Unquote(`"` + literals[lineNum-1] + `"`)
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, s)
	return err
}
