package devmode_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/templar/devmode"
	"github.com/matthewmueller/templar/internal/testutil"
)

func TestData(t *testing.T) {
	is := is.New(t)
	dirs, err := testutil.TestData("testdata")
	is.NoErr(err)
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			is := is.New(t)
			templPaths, err := filepath.Glob(filepath.Join(dir, "*.templ"))
			is.NoErr(err)
			for _, templPath := range templPaths {
				templCode, err := os.ReadFile(templPath)
				is.NoErr(err)
				generated, literals, err := testutil.Generate(templPath, string(templCode))
				is.NoErr(err)

				// Patch with devmode
				modified, err := devmode.Transform(templPath, generated)
				is.NoErr(err)

				is.NoErr(testutil.Golden(testutil.GoPath(templPath), string(modified)))
				is.NoErr(testutil.Golden(testutil.TxtPath(templPath), literals))
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

	oneGen, oneLits, err := testutil.Generate("testcall.templ", one)
	is.NoErr(err)

	twoGen, twoLits, err := testutil.Generate("testcall.templ", two)
	is.NoErr(err)

	is.True(!bytes.Equal(oneGen, twoGen)) // generated code should be different
	is.True(oneLits != twoLits)           // literals should be different

	oneMod, err := devmode.Transform("testcall.templ", oneGen)
	is.NoErr(err)

	twoMod, err := devmode.Transform("testcall.templ", twoGen)
	is.NoErr(err)

	is.True(bytes.Equal(oneMod, twoMod)) // modified code should be the same
}
