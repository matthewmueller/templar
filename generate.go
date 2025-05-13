package templar

import (
	"bytes"
	"fmt"
	"io"

	"github.com/a-h/templ/parser/v2"
)

func Literals(w io.Writer, template *parser.TemplateFile, visitors ...parser.Visitor) error {
	return nil
}

func Generate(w io.Writer, template *parser.TemplateFile, visitors ...parser.Visitor) error {
	generator := &Generator{
		Filename: template.Filepath,
		Runtime: &Import{
			Name: "templruntime",
			Path: "github.com/a-h/templ/runtime",
		},
	}
	return generator.Generate(w, template, visitors...)
}

type Generator struct {
	Filename string
	Runtime  *Import
	Live     bool
}

type Import struct {
	Name string
	Path string
}

func (g *Generator) Generate(w io.Writer, template *parser.TemplateFile, visitors ...parser.Visitor) error {
	if g.Runtime == nil {
		return fmt.Errorf("runtime is not set")
	}
	generator := &generator{w, g.Filename, g.Runtime, map[string]string{}, 0, 0, 0}
	if err := generator.VisitTemplateFile(template); err != nil {
		return fmt.Errorf("error visiting template file: %w", err)
	}
	return nil
}

type generator struct {
	w         io.Writer
	filename  string
	runtime   *Import
	variables map[string]string
	count     int
	rows      int
	indents   int
}

func (g *generator) setVar(name string) string {
	if _, ok := g.variables[name]; !ok {
		g.variables[name] = "_" + name
		return g.variables[name]
	}
	g.count++
	g.variables[name] = fmt.Sprintf("_%s%d", name, g.count)
	return g.variables[name]
}

func (g *generator) getVar(name string) (string, bool) {
	v, ok := g.variables[name]
	return v, ok
}

func (g *generator) Writef(s string, args ...any) {
	g.w.Write([]byte(fmt.Sprintf(s, args...)))
}

func (g *generator) indent(delta int) {
	g.indents += delta
}

func (g *generator) nextRow() int {
	g.rows++
	return g.rows
}

func (g *generator) tabs() string {
	return string(bytes.Repeat([]byte("\t"), g.indents))
}

var _ parser.Visitor = (*generator)(nil)

func (g *generator) VisitTemplateFile(n *parser.TemplateFile) error {
	for _, header := range n.Header {
		if !header.BeforePackage {
			continue
		}
		if err := g.VisitTemplateFileGoExpression(header); err != nil {
			return fmt.Errorf("error visiting template file go expression: %w", err)
		}
	}

	if err := g.VisitPackage(&n.Package); err != nil {
		return fmt.Errorf("error visiting package: %w", err)
	}

	for _, header := range n.Header {
		if header.BeforePackage {
			continue
		}
		if err := g.VisitTemplateFileGoExpression(header); err != nil {
			return fmt.Errorf("error visiting template file go expression: %w", err)
		}
	}

	g.Writef("import %q\n", "github.com/a-h/templ")
	g.Writef("import %s %q\n", g.runtime.Name, g.runtime.Path)

	for _, node := range n.Nodes {
		if err := node.Visit(g); err != nil {
			return fmt.Errorf("error visiting node: %w", err)
		}
	}

	return nil
}

func (g *generator) VisitTemplateFileGoExpression(n *parser.TemplateFileGoExpression) error {
	g.w.Write([]byte(n.Expression.Value))
	return nil
}

func (g *generator) VisitPackage(n *parser.Package) error {
	g.w.Write([]byte(n.Expression.Value))
	g.w.Write([]byte{'\n'})
	return nil
}

func (g *generator) VisitWhitespace(n *parser.Whitespace) error {
	g.w.Write([]byte(n.Value))
	return nil
}

func (g *generator) VisitCSSTemplate(n *parser.CSSTemplate) error {
	return fmt.Errorf("VisitCSSTemplate not implemented")
}

func (g *generator) VisitConstantCSSProperty(n *parser.ConstantCSSProperty) error {
	return fmt.Errorf("VisitConstantCSSProperty not implemented")
}

func (g *generator) VisitExpressionCSSProperty(n *parser.ExpressionCSSProperty) error {
	return fmt.Errorf("VisitExpressionCSSProperty not implemented")
}

func (g *generator) VisitDocType(n *parser.DocType) error {
	return fmt.Errorf("VisitDocType not implemented")
}

func preHtmlTemplate(fn, runtime, childrenVar string) []byte {
	return []byte(fmt.Sprintf(`func %[1]s templ.Component {
	return %[2]s.GeneratedTemplate(func(_in %[2]s.GeneratedComponentInput) (_err error) {
		_w, _ctx := _in.Writer, _in.Context
		if _ctx.Err() != nil {
			return _ctx.Err()
		}
		_buf, _is_buf := %[2]s.GetBuffer(_w)
		if !_is_buf {
			defer func() {
				_err2 := %[2]s.ReleaseBuffer(_buf)
				if _err == nil {
					_err = _err2
				}
			}()
		}
		_ctx = templ.InitializeContext(_ctx)
		%[3]s := templ.GetChildren(_ctx)
		if %[3]s == nil {
			%[3]s = templ.NopComponent
		}
		_ctx = templ.ClearChildren(_ctx)
`, fn, runtime, childrenVar))
}

func postHtmlTemplate() []byte {
	return []byte(`		return nil
	})
}`)
}

func (g *generator) VisitHTMLTemplate(n *parser.HTMLTemplate) error {
	g.w.Write(preHtmlTemplate(n.Expression.Value, g.runtime.Name, g.setVar("children")))
	g.indent(2)
	for _, child := range n.Children {
		if err := child.Visit(g); err != nil {
			return fmt.Errorf("error visiting child: %w", err)
		}
	}
	g.indent(-2)
	g.w.Write(postHtmlTemplate())
	return nil
}

func (g *generator) VisitText(n *parser.Text) error {
	fmt.Println("VisitText")
	return nil
}

func (g *generator) VisitElement(n *parser.Element) error {
	g.Writef(`%s_err = %s.WriteString(_buf, %d, "<%s`, g.tabs(), g.runtime.Name, g.nextRow(), n.Name)
	// g.w.Write([]byte(fmt.Sprintf(`%s.Write([]byte(w, "`, g.runtime.Name)))
	// g.w.Write([]byte("<"))
	// g.w.Write([]byte(n.Name))
	// g.w.Write([]byte(`")`))
	// g.w.Write([]byte("\n"))
	for _, attr := range n.Attributes {
		g.Writef(" ")
		if err := attr.Visit(g); err != nil {
			return fmt.Errorf("error visiting attribute: %w", err)
		}
	}
	g.Writef(`>`)
	g.Writef(`")` + "\n")
	g.Writef("%sif _err != nil {\n", g.tabs())
	g.indent(1)
	g.Writef("%sreturn _err\n", g.tabs())
	g.indent(-1)
	g.Writef("%s}\n", g.tabs())

	for _, child := range n.Children {
		if err := child.Visit(g); err != nil {
			return fmt.Errorf("error visiting child: %w", err)
		}
	}

	g.Writef(`%s_err = %s.WriteString(_buf, %d, "</%s>")`+"\n", g.tabs(), g.runtime.Name, g.nextRow(), n.Name)
	g.Writef("%sif _err != nil {\n", g.tabs())
	g.indent(1)
	g.Writef("%sreturn _err\n", g.tabs())
	g.indent(-1)
	g.Writef("%s}\n", g.tabs())

	// g.w.Write([]byte("\n"))
	return nil
}

func (g *generator) VisitScriptElement(n *parser.ScriptElement) error {
	fmt.Println("VisitScriptElement")
	return fmt.Errorf("VisitScriptElement not implemented")
}

func (g *generator) VisitRawElement(n *parser.RawElement) error {
	return fmt.Errorf("VisitRawElement not implemented")
}

func (g *generator) VisitBoolConstantAttribute(n *parser.BoolConstantAttribute) error {
	fmt.Println("VisitBoolConstantAttribute")
	return fmt.Errorf("VisitBoolConstantAttribute not implemented")
}

func (g *generator) VisitConstantAttribute(n *parser.ConstantAttribute) error {
	fmt.Println("VisitConstantAttribute")
	return fmt.Errorf("VisitConstantAttribute not implemented")
}

func (g *generator) VisitBoolExpressionAttribute(n *parser.BoolExpressionAttribute) error {
	fmt.Println("VisitBoolExpressionAttribute")
	return fmt.Errorf("VisitBoolExpressionAttribute not implemented")
}

func (g *generator) VisitExpressionAttribute(n *parser.ExpressionAttribute) error {
	g.w.Write([]byte(fmt.Sprintf(`%s.Write([]byte(w, %s=%s)`, g.runtime.Name, n.Name, n.Expression.Value)))
	fmt.Println("VisitExpressionAttribute")
	return fmt.Errorf("VisitExpressionAttribute not implemented")
}

func (g *generator) VisitSpreadAttributes(n *parser.SpreadAttributes) error {
	fmt.Println("VisitSpreadAttributes")
	return fmt.Errorf("VisitSpreadAttributes not implemented")
}

func (g *generator) VisitConditionalAttribute(n *parser.ConditionalAttribute) error {
	fmt.Println("VisitConditionalAttribute")
	return fmt.Errorf("VisitConditionalAttribute not implemented")
}

func (g *generator) VisitGoComment(n *parser.GoComment) error {
	fmt.Println("VisitGoComment")
	return fmt.Errorf("VisitGoComment not implemented")
}

func (g *generator) VisitHTMLComment(n *parser.HTMLComment) error {
	fmt.Println("VisitHTMLComment")
	return fmt.Errorf("VisitHTMLComment not implemented")
}

func (g *generator) VisitCallTemplateExpression(n *parser.CallTemplateExpression) error {
	fmt.Println("VisitCallTemplateExpression")
	return fmt.Errorf("VisitCallTemplateExpression not implemented")
}

func (g *generator) VisitTemplElementExpression(n *parser.TemplElementExpression) error {
	fmt.Println("VisitTemplElementExpression")
	return fmt.Errorf("VisitTemplElementExpression not implemented")
}

func (g *generator) VisitChildrenExpression(n *parser.ChildrenExpression) error {
	fmt.Println("VisitChildrenExpression")
	return fmt.Errorf("VisitChildrenExpression not implemented")
}

func (g *generator) VisitIfExpression(n *parser.IfExpression) error {
	fmt.Println("VisitIfExpression")
	return fmt.Errorf("VisitIfExpression not implemented")
}

func (g *generator) VisitSwitchExpression(n *parser.SwitchExpression) error {
	fmt.Println("VisitSwitchExpression")
	return fmt.Errorf("VisitSwitchExpression not implemented")
}

func (g *generator) VisitForExpression(n *parser.ForExpression) error {
	g.Writef("%sfor %s {\n", g.tabs(), n.Expression.Value)
	g.indent(1)
	for _, child := range n.Children {
		if err := child.Visit(g); err != nil {
			return fmt.Errorf("error visiting child: %w", err)
		}
	}
	g.indent(-1)
	g.Writef("\n%s}\n", g.tabs())
	return nil
}

func (g *generator) VisitGoCode(n *parser.GoCode) error {
	fmt.Println("VisitGoCode")
	return nil
}

func (g *generator) VisitStringExpression(n *parser.StringExpression) error {
	// var templ_7745c5c3_Var2 string
	// templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(item)
	// if templ_7745c5c3_Err != nil {
	// 	return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/test/templ/test-for/template.templ`, Line: 5, Col: 13}
	// }
	// _, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
	// if templ_7745c5c3_Err != nil {
	// 	return templ_7745c5c3_Err
	// }

	out := g.setVar("var")
	g.Writef("%s%s, _err := templ.JoinStringErrs(%s)\n", g.tabs(), out, n.Expression.Value)
	g.Writef("%sif _err != nil {\n", g.tabs())
	g.indent(1)
	g.Writef("%sreturn templ.Error{Err: _err, FileName: %q, Line: %d, Col: %d}\n", g.tabs(), g.filename, n.Expression.Range.From.Line+1, n.Expression.Range.To.Col)
	g.indent(-1)
	g.Writef("%s}\n", g.tabs())
	g.Writef("%s_, _err = _buf.WriteString(templ.EscapeString(%s))\n", g.tabs(), out)
	g.Writef("%sif _err != nil {\n", g.tabs())
	g.indent(1)
	g.Writef("%sreturn _err\n", g.tabs())
	g.indent(-1)
	g.Writef("%s}\n", g.tabs())

	return nil
}

func (g *generator) VisitScriptTemplate(n *parser.ScriptTemplate) error {
	fmt.Println("VisitScriptTemplate")
	return nil
}
