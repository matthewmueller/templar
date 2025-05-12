package templar

import (
	"io"

	"github.com/a-h/templ/parser/v2"
)

func Generate(w io.Writer, template *parser.TemplateFile, visitors ...parser.Visitor) error {
	return nil
}

func Literals(w io.Writer, template *parser.TemplateFile, visitors ...parser.Visitor) error {
	return nil
}
