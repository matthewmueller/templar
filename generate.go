package templar

import (
	"fmt"
	"html"
	"io"
	"strconv"
	"strings"

	"github.com/a-h/templ/generator"
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
	visitor := &visitor{generator.NewRangeWriter(w), g.Filename, g.Runtime, map[string]string{}, 0, 0, 0, 0, ""}
	if err := visitor.VisitTemplateFile(template); err != nil {
		return fmt.Errorf("error visiting template file: %w", err)
	}
	return nil
}

type visitor struct {
	w           *generator.RangeWriter
	filename    string
	runtime     *Import
	variables   map[string]string
	count       int
	rows        int
	indents     int
	variableID  int
	childrenVar string
}

func (g *visitor) writeTemplBuffer(indentLevel int) (err error) {
	// templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
	if _, err = g.w.WriteIndent(indentLevel, "templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)\n"); err != nil {
		return err
	}
	// if !templ_7745c5c3_IsBuffer {
	//	defer func() {
	//		templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
	//		if templ_7745c5c3_Err == nil {
	//			templ_7745c5c3_Err = templ_7745c5c3_BufErr
	//		}
	//	}()
	// }
	if _, err = g.w.WriteIndent(indentLevel, "if !templ_7745c5c3_IsBuffer {\n"); err != nil {
		return err
	}
	{
		indentLevel++
		if _, err = g.w.WriteIndent(indentLevel, "defer func() {\n"); err != nil {
			return err
		}
		{
			indentLevel++
			if _, err = g.w.WriteIndent(indentLevel, "templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)\n"); err != nil {
				return err
			}
			if _, err = g.w.WriteIndent(indentLevel, "if templ_7745c5c3_Err == nil {\n"); err != nil {
				return err
			}
			{
				indentLevel++
				if _, err = g.w.WriteIndent(indentLevel, "templ_7745c5c3_Err = templ_7745c5c3_BufErr\n"); err != nil {
					return err
				}
				indentLevel--
			}
			if _, err = g.w.WriteIndent(indentLevel, "}\n"); err != nil {
				return err
			}
			indentLevel--
		}
		if _, err = g.w.WriteIndent(indentLevel, "}()\n"); err != nil {
			return err
		}
		indentLevel--
	}
	if _, err = g.w.WriteIndent(indentLevel, "}\n"); err != nil {
		return err
	}
	return
}

func (g *visitor) createVariableName() string {
	g.variableID++
	return "templ_7745c5c3_Var" + strconv.Itoa(g.variableID)
}

func (g *visitor) writeExpressionErrorHandler(indentLevel int, expression parser.Expression) (err error) {
	_, err = g.w.WriteIndent(indentLevel, "if templ_7745c5c3_Err != nil {\n")
	if err != nil {
		return err
	}
	indentLevel++
	line := int(expression.Range.To.Line + 1)
	col := int(expression.Range.To.Col)
	_, err = g.w.WriteIndent(indentLevel, "return	templ.Error{Err: templ_7745c5c3_Err, FileName: "+createGoString(g.filename)+", Line: "+strconv.Itoa(line)+", Col: "+strconv.Itoa(col)+"}\n")
	if err != nil {
		return err
	}
	indentLevel--
	_, err = g.w.WriteIndent(indentLevel, "}\n")
	if err != nil {
		return err
	}
	return err
}

func (g *visitor) writeErrorHandler(indentLevel int) (err error) {
	_, err = g.w.WriteIndent(indentLevel, "if templ_7745c5c3_Err != nil {\n")
	if err != nil {
		return err
	}
	indentLevel++
	_, err = g.w.WriteIndent(indentLevel, "return templ_7745c5c3_Err\n")
	if err != nil {
		return err
	}
	indentLevel--
	_, err = g.w.WriteIndent(indentLevel, "}\n")
	if err != nil {
		return err
	}
	return err
}

// func (g *visitor) setVar(name string) string {
// 	if _, ok := g.variables[name]; !ok {
// 		g.variables[name] = "_" + name
// 		return g.variables[name]
// 	}
// 	g.count++
// 	g.variables[name] = fmt.Sprintf("_%s%d", name, g.count)
// 	return g.variables[name]
// }

// func (g *visitor) getVar(name string) (string, bool) {
// 	v, ok := g.variables[name]
// 	return v, ok
// }

// func (g *visitor) Writef(s string, args ...any) {
// 	g.w.Write([]byte(fmt.Sprintf(s, args...)))
// }

// func (g *visitor) indent(delta int) {
// 	g.indents += delta
// }

// func (g *visitor) nextRow() int {
// 	g.rows++
// 	return g.rows
// }

// func (g *visitor) tabs() string {
// 	return string(bytes.Repeat([]byte("\t"), g.indents))
// }

var _ parser.Visitor = (*visitor)(nil)

func (g *visitor) VisitTemplateFile(n *parser.TemplateFile) error {
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

	if _, err := g.w.Write("import \"github.com/a-h/templ\"\n"); err != nil {
		return err
	}
	if _, err := g.w.Write(fmt.Sprintf("import %s %q\n", g.runtime.Name, g.runtime.Path)); err != nil {
		return err
	}

	for _, node := range n.Nodes {
		if err := node.Visit(g); err != nil {
			return fmt.Errorf("error visiting node: %w", err)
		}
	}

	return nil
}

func (g *visitor) VisitTemplateFileGoExpression(n *parser.TemplateFileGoExpression) error {
	if _, err := g.w.Write(n.Expression.Value); err != nil {
		return err
	}
	v := n.Expression.Value
	lineSlice := strings.Split(v, "\n")
	lastLine := lineSlice[len(lineSlice)-1]
	if strings.HasPrefix(lastLine, "//") {
		if _, err := g.w.WriteIndent(0, "\n"); err != nil {
			return err
		}
		return nil
	}
	if _, err := g.w.WriteIndent(0, "\n\n"); err != nil {
		return err
	}
	return nil
}

func (g *visitor) VisitPackage(n *parser.Package) error {
	// package ...
	if _, err := g.w.Write(n.Expression.Value + "\n\n"); err != nil {
		return err
	}
	if _, err := g.w.Write("//lint:file-ignore SA4006 This context is only used if a nested component is present.\n\n"); err != nil {
		return err
	}
	return nil
}

func (g *visitor) VisitWhitespace(n *parser.Whitespace) error {
	// g.w.Write([]byte(n.Value))
	return fmt.Errorf("VisitWhitespace not implemented")
}

func (g *visitor) VisitCSSTemplate(n *parser.CSSTemplate) error {
	return fmt.Errorf("VisitCSSTemplate not implemented")
}

func (g *visitor) VisitConstantCSSProperty(n *parser.ConstantCSSProperty) error {
	return fmt.Errorf("VisitConstantCSSProperty not implemented")
}

func (g *visitor) VisitExpressionCSSProperty(n *parser.ExpressionCSSProperty) error {
	return fmt.Errorf("VisitExpressionCSSProperty not implemented")
}

func (g *visitor) VisitDocType(n *parser.DocType) error {
	return fmt.Errorf("VisitDocType not implemented")
}

// func preHtmlTemplate(fn, runtime, childrenVar string) []byte {
// 	return []byte(fmt.Sprintf(`func %[1]s templ.Component {
// 	return %[2]s.GeneratedTemplate(func(_in %[2]s.GeneratedComponentInput) (_err error) {
// 		_w, _ctx := _in.Writer, _in.Context
// 		if _ctx.Err() != nil {
// 			return _ctx.Err()
// 		}
// 		_buf, _is_buf := %[2]s.GetBuffer(_w)
// 		if !_is_buf {
// 			defer func() {
// 				_err2 := %[2]s.ReleaseBuffer(_buf)
// 				if _err == nil {
// 					_err = _err2
// 				}
// 			}()
// 		}
// 		_ctx = templ.InitializeContext(_ctx)
// 		%[3]s := templ.GetChildren(_ctx)
// 		if %[3]s == nil {
// 			%[3]s = templ.NopComponent
// 		}
// 		_ctx = templ.ClearChildren(_ctx)
// `, fn, runtime, childrenVar))
// }

// func postHtmlTemplate() []byte {
// 	return []byte(`		return nil
// 	})
// }`)
// }

func (g *visitor) VisitHTMLTemplate(n *parser.HTMLTemplate) (err error) {
	indentLevel := 0

	// func
	if _, err := g.w.Write("func "); err != nil {
		return err
	}

	// (r *Receiver) Name(params []string)
	if _, err := g.w.Write(n.Expression.Value); err != nil {
		return err
	}

	// templ.Component {
	if _, err = g.w.Write(" templ.Component {\n"); err != nil {
		return err
	}
	indentLevel++
	// return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
	if _, err = g.w.WriteIndent(indentLevel, "return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {\n"); err != nil {
		return err
	}
	{
		indentLevel++
		if _, err = g.w.WriteIndent(indentLevel, "templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context\n"); err != nil {
			return err
		}
		if _, err = g.w.WriteIndent(indentLevel, "if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {\n"); err != nil {
			return err
		}
		{
			indentLevel++
			if _, err = g.w.WriteIndent(indentLevel, "return templ_7745c5c3_CtxErr"); err != nil {
				return err
			}
			indentLevel--
		}
		if _, err = g.w.WriteIndent(indentLevel, "}\n"); err != nil {
			return err
		}
		if err := g.writeTemplBuffer(indentLevel); err != nil {
			return err
		}
		// ctx = templ.InitializeContext(ctx)
		if _, err = g.w.WriteIndent(indentLevel, "ctx = templ.InitializeContext(ctx)\n"); err != nil {
			return err
		}
		g.childrenVar = g.createVariableName()
		// templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		// if templ_7745c5c3_Var1 == nil {
		//  	templ_7745c5c3_Var1 = templ.NopComponent
		// }
		if _, err = g.w.WriteIndent(indentLevel, fmt.Sprintf("%s := templ.GetChildren(ctx)\n", g.childrenVar)); err != nil {
			return err
		}
		if _, err = g.w.WriteIndent(indentLevel, fmt.Sprintf("if %s == nil {\n", g.childrenVar)); err != nil {
			return err
		}
		{
			indentLevel++
			if _, err = g.w.WriteIndent(indentLevel, fmt.Sprintf("%s = templ.NopComponent\n", g.childrenVar)); err != nil {
				return err
			}
			indentLevel--
		}
		if _, err = g.w.WriteIndent(indentLevel, "}\n"); err != nil {
			return err
		}
		// ctx = templ.ClearChildren(children)
		if _, err = g.w.WriteIndent(indentLevel, "ctx = templ.ClearChildren(ctx)\n"); err != nil {
			return err
		}
		// Nodes.
		for _, child := range stripWhitespace(n.Children) {
			if err := child.Visit(g); err != nil {
				return fmt.Errorf("error visiting child: %w", err)
			}
		}
		// if err = g.writeNodes(indentLevel, stripWhitespace(n.Children), nil); err != nil {
		// 	return err
		// }
		// return nil
		if _, err = g.w.WriteIndent(indentLevel, "return nil\n"); err != nil {
			return err
		}
		indentLevel--
	}
	// })
	if _, err = g.w.WriteIndent(indentLevel, "})\n"); err != nil {
		return err
	}
	indentLevel--
	// }

	// Note: gofmt wants to remove a single empty line at the end of a file
	// so we have to make sure we don't output one if this is the last node.
	closingBrace := "}\n\n"
	// if nodeIdx+1 >= len(g.tf.Nodes) {
	// 	closingBrace = "}\n"
	// }

	if _, err := g.w.WriteIndent(indentLevel, closingBrace); err != nil {
		return err
	}

	// Keep a track of symbol ranges for the LSP.
	// tgtSymbolRange.To = r.To
	// g.sourceMap.AddSymbolRange(t.Range, tgtSymbolRange)

	return nil
}

func (g *visitor) VisitText(n *parser.Text) error {
	fmt.Println("VisitText")
	return fmt.Errorf("VisitText not implemented")
}

func (g *visitor) VisitElement(n *parser.Element) error {
	if len(n.Attributes) == 0 {
		// <div>
		if _, err := g.w.WriteStringLiteral(g.indents, fmt.Sprintf(`<%s>`, html.EscapeString(n.Name))); err != nil {
			return err
		}
	} else {
		// attrs := copyAttributes(n.Attributes)
		for _, attr := range n.Attributes {
			if err := attr.Visit(g); err != nil {
				return fmt.Errorf("error visiting attribute: %w", err)
			}
		}
		// TODO: figure out how to add
		// // <style type="text/css"></style>
		// if err := g.writeElementCSS(g.indents, attrs); err != nil {
		// 	return err
		// }
		// // <script></script>
		// if err := g.writeElementScript(g.indents, attrs); err != nil {
		// 	return err
		// }
		// <div
		if _, err := g.w.WriteStringLiteral(g.indents, fmt.Sprintf(`<%s`, html.EscapeString(n.Name))); err != nil {
			return err
		}
		// Visit attributes.
		for _, attr := range n.Attributes {
			if err := attr.Visit(g); err != nil {
				return fmt.Errorf("error visiting attribute: %w", err)
			}
		}
		// >
		if _, err := g.w.WriteStringLiteral(g.indents, `>`); err != nil {
			return err
		}
	}
	// Skip children and close tag for void elements.
	if n.IsVoidElement() && len(n.Children) == 0 {
		return nil
	}
	// Children.
	for _, child := range stripWhitespace(n.Children) {
		if err := child.Visit(g); err != nil {
			return fmt.Errorf("error visiting child: %w", err)
		}
	}
	// </div>
	if _, err := g.w.WriteStringLiteral(g.indents, fmt.Sprintf(`</%s>`, html.EscapeString(n.Name))); err != nil {
		return err
	}
	return nil
}

func (g *visitor) VisitScriptElement(n *parser.ScriptElement) error {
	fmt.Println("VisitScriptElement")
	return fmt.Errorf("VisitScriptElement not implemented")
}

func (g *visitor) VisitRawElement(n *parser.RawElement) error {
	return fmt.Errorf("VisitRawElement not implemented")
}

func (g *visitor) VisitBoolConstantAttribute(n *parser.BoolConstantAttribute) error {
	fmt.Println("VisitBoolConstantAttribute")
	return fmt.Errorf("VisitBoolConstantAttribute not implemented")
}

func (g *visitor) VisitConstantAttribute(n *parser.ConstantAttribute) error {
	fmt.Println("VisitConstantAttribute")
	return fmt.Errorf("VisitConstantAttribute not implemented")
}

func (g *visitor) VisitBoolExpressionAttribute(n *parser.BoolExpressionAttribute) error {
	fmt.Println("VisitBoolExpressionAttribute")
	return fmt.Errorf("VisitBoolExpressionAttribute not implemented")
}

func (g *visitor) VisitExpressionAttribute(n *parser.ExpressionAttribute) error {
	// g.w.Write([]byte(fmt.Sprintf(`%s.Write([]byte(w, %s=%s)`, g.runtime.Name, n.Name, n.Expression.Value)))
	// fmt.Println("VisitExpressionAttribute")
	return fmt.Errorf("VisitExpressionAttribute not implemented")
}

func (g *visitor) VisitSpreadAttributes(n *parser.SpreadAttributes) error {
	fmt.Println("VisitSpreadAttributes")
	return fmt.Errorf("VisitSpreadAttributes not implemented")
}

func (g *visitor) VisitConditionalAttribute(n *parser.ConditionalAttribute) error {
	fmt.Println("VisitConditionalAttribute")
	return fmt.Errorf("VisitConditionalAttribute not implemented")
}

func (g *visitor) VisitGoComment(n *parser.GoComment) error {
	fmt.Println("VisitGoComment")
	return fmt.Errorf("VisitGoComment not implemented")
}

func (g *visitor) VisitHTMLComment(n *parser.HTMLComment) error {
	fmt.Println("VisitHTMLComment")
	return fmt.Errorf("VisitHTMLComment not implemented")
}

func (g *visitor) VisitCallTemplateExpression(n *parser.CallTemplateExpression) error {
	fmt.Println("VisitCallTemplateExpression")
	return fmt.Errorf("VisitCallTemplateExpression not implemented")
}

func (g *visitor) VisitTemplElementExpression(n *parser.TemplElementExpression) error {
	fmt.Println("VisitTemplElementExpression")
	return fmt.Errorf("VisitTemplElementExpression not implemented")
}

func (g *visitor) VisitChildrenExpression(n *parser.ChildrenExpression) error {
	fmt.Println("VisitChildrenExpression")
	return fmt.Errorf("VisitChildrenExpression not implemented")
}

func (g *visitor) VisitIfExpression(n *parser.IfExpression) error {
	fmt.Println("VisitIfExpression")
	return fmt.Errorf("VisitIfExpression not implemented")
}

func (g *visitor) VisitSwitchExpression(n *parser.SwitchExpression) error {
	fmt.Println("VisitSwitchExpression")
	return fmt.Errorf("VisitSwitchExpression not implemented")
}

func (g *visitor) VisitForExpression(n *parser.ForExpression) error {
	// for
	if _, err := g.w.WriteIndent(g.indents, `for `); err != nil {
		return err
	}
	// i, v := range p.Stuff
	if _, err := g.w.Write(n.Expression.Value); err != nil {
		return err
	}

	// {
	if _, err := g.w.Write(` {` + "\n"); err != nil {
		return err
	}
	// Children.
	g.indents++
	for _, child := range stripLeadingAndTrailingWhitespace(n.Children) {
		if err := child.Visit(g); err != nil {
			return fmt.Errorf("error visiting child: %w", err)
		}
	}
	g.indents--
	// }
	if _, err := g.w.WriteIndent(g.indents, `}`+"\n"); err != nil {
		return err
	}

	return nil
}

func (g *visitor) VisitGoCode(n *parser.GoCode) error {
	fmt.Println("VisitGoCode")
	return fmt.Errorf("VisitGoCode not implemented")
}

func (g *visitor) VisitStringExpression(n *parser.StringExpression) error {
	if strings.TrimSpace(n.Expression.Value) == "" {
		return nil
	}
	vn := g.createVariableName()
	// var vn string
	if _, err := g.w.WriteIndent(g.indents, "var "+vn+" string\n"); err != nil {
		return err
	}
	// vn, templ_7745c5c3_Err = templ.JoinStringErrs(
	if _, err := g.w.WriteIndent(g.indents, vn+", templ_7745c5c3_Err = templ.JoinStringErrs("); err != nil {
		return err
	}
	// p.Name()
	if _, err := g.w.Write(n.Expression.Value); err != nil {
		return err
	}
	// g.sourceMap.Add(e, r)
	// )
	if _, err := g.w.Write(")\n"); err != nil {
		return err
	}

	// String expression error handler.
	err := g.writeExpressionErrorHandler(g.indents, n.Expression)
	if err != nil {
		return err
	}

	// _, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(vn)
	if _, err = g.w.WriteIndent(g.indents, "_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString("+vn+"))\n"); err != nil {
		return err
	}
	if err := g.writeErrorHandler(g.indents); err != nil {
		return err
	}
	return nil
	// var templ_7745c5c3_Var2 string
	// templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(item)
	// if templ_7745c5c3_Err != nil {
	// 	return templ.Error{Err: templ_7745c5c3_Err, FileName: `internal/test/templ/test-for/template.templ`, Line: 5, Col: 13}
	// }
	// _, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
	// if templ_7745c5c3_Err != nil {
	// 	return templ_7745c5c3_Err
	// }

	// out := g.setVar("var")
	// g.Writef("%s%s, _err := templ.JoinStringErrs(%s)\n", g.tabs(), out, n.Expression.Value)
	// g.Writef("%sif _err != nil {\n", g.tabs())
	// g.indent(1)
	// g.Writef("%sreturn templ.Error{Err: _err, FileName: %q, Line: %d, Col: %d}\n", g.tabs(), g.filename, n.Expression.Range.From.Line+1, n.Expression.Range.To.Col)
	// g.indent(-1)
	// g.Writef("%s}\n", g.tabs())
	// g.Writef("%s_, _err = _buf.WriteString(templ.EscapeString(%s))\n", g.tabs(), out)
	// g.Writef("%sif _err != nil {\n", g.tabs())
	// g.indent(1)
	// g.Writef("%sreturn _err\n", g.tabs())
	// g.indent(-1)
	// g.Writef("%s}\n", g.tabs())

	return fmt.Errorf("VisitStringExpression not implemented")
}

func (g *visitor) VisitScriptTemplate(n *parser.ScriptTemplate) error {
	fmt.Println("VisitScriptTemplate")
	return fmt.Errorf("VisitScriptTemplate not implemented")
}

func stripWhitespace(input []parser.Node) (output []parser.Node) {
	for i, n := range input {
		if _, isWhiteSpace := n.(*parser.Whitespace); !isWhiteSpace {
			output = append(output, input[i])
		}
	}
	return output
}

func stripLeadingWhitespace(nodes []parser.Node) []parser.Node {
	for i, n := range nodes {
		if _, isWhiteSpace := n.(*parser.Whitespace); !isWhiteSpace {
			return nodes[i:]
		}
	}
	return []parser.Node{}
}

func stripTrailingWhitespace(nodes []parser.Node) []parser.Node {
	for i := len(nodes) - 1; i >= 0; i-- {
		n := nodes[i]
		if _, isWhiteSpace := n.(*parser.Whitespace); !isWhiteSpace {
			return nodes[0 : i+1]
		}
	}
	return []parser.Node{}
}

func stripLeadingAndTrailingWhitespace(nodes []parser.Node) []parser.Node {
	return stripTrailingWhitespace(stripLeadingWhitespace(nodes))
}

func createGoString(s string) string {
	var sb strings.Builder
	sb.WriteRune('`')
	sects := strings.Split(s, "`")
	for i, sect := range sects {
		sb.WriteString(sect)
		if len(sects) > i+1 {
			sb.WriteString("` + \"`\" + `")
		}
	}
	sb.WriteRune('`')
	return sb.String()
}
