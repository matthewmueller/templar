package scopedcss_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/templar/internal/testutil"
	"github.com/matthewmueller/templar/scopedcss"
)

func TestInline(t *testing.T) {
	is := is.New(t)
	dirs, err := testutil.TestData("testdata")
	is.NoErr(err)
	for _, dir := range dirs {
		t.Run(filepath.Base(dir), func(t *testing.T) {
			is := is.New(t)
			templPath := filepath.Join(dir, "input.templ")
			templCode, err := os.ReadFile(templPath)
			is.NoErr(err)

			templAst, err := testutil.Parse(templPath, string(templCode))
			is.NoErr(err)

			err = scopedcss.Inline(templPath, templAst)
			is.NoErr(err)

			actual, err := testutil.Format(templAst)
			is.NoErr(err)

			dir := filepath.Dir(templPath)
			is.NoErr(testutil.Golden(filepath.Join(dir, "inline.templ"), actual))
		})
	}
}
