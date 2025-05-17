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
	visitor := &visitor{generator.NewRangeWriter(w), g.Filename, g.Runtime, map[string]string{}, parser.NewSourceMap(), "", 0, 0, 0, 0, ""}
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
	sourceMap   *parser.SourceMap
	elementName string
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
	if len(n.Value) == 0 {
		return nil
	}
	// _, err = templ_7745c5c3_Buffer.WriteString(` `)
	if _, err := g.w.WriteStringLiteral(g.indents, " "); err != nil {
		return err
	}
	return nil
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
	_, err := g.w.WriteStringLiteral(g.indents, escapeQuotes(n.Value))
	return err
}

func (g *visitor) writeAttributeCSS(indentLevel int, attr *parser.ExpressionAttribute) (result *parser.ExpressionAttribute, ok bool, err error) {
	var r parser.Range
	name := html.EscapeString(attr.Name)
	if name != "class" {
		ok = false
		return
	}
	// Create a class name for the style.
	// The expression can either be expecting a templ.Classes call, or an expression that returns
	// var templ_7745c5c3_CSSClasses = []any{
	classesName := g.createVariableName()
	if _, err = g.w.WriteIndent(indentLevel, "var "+classesName+" = []any{"); err != nil {
		return
	}
	// p.Name()
	if r, err = g.w.Write(attr.Expression.Value); err != nil {
		return
	}
	g.sourceMap.Add(attr.Expression, r)
	// }\n
	if _, err = g.w.Write("}\n"); err != nil {
		return
	}
	// Render the CSS before the element if required.
	// templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_CSSClasses...)
	if _, err = g.w.WriteIndent(indentLevel, "templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, "+classesName+"...)\n"); err != nil {
		return
	}
	if err = g.writeErrorHandler(indentLevel); err != nil {
		return
	}
	// Rewrite the ExpressionAttribute to point at the new variable.
	newAttr := &parser.ExpressionAttribute{
		Name:      attr.Name,
		NameRange: attr.NameRange,
		Expression: parser.Expression{
			Value: "templ.CSSClasses(" + classesName + ").String()",
		},
	}
	return newAttr, true, nil
}

func (g *visitor) writeAttributesCSS(indentLevel int, attrs []parser.Attribute) (err error) {
	for i, attr := range attrs {
		if attr, ok := attr.(*parser.ExpressionAttribute); ok {
			attr, ok, err = g.writeAttributeCSS(indentLevel, attr)
			if err != nil {
				return err
			}
			if ok {
				attrs[i] = attr
			}
		}
		if cattr, ok := attr.(*parser.ConditionalAttribute); ok {
			err = g.writeAttributesCSS(indentLevel, cattr.Then)
			if err != nil {
				return err
			}
			err = g.writeAttributesCSS(indentLevel, cattr.Else)
			if err != nil {
				return err
			}
			attrs[i] = cattr
		}
	}
	return nil
}

func (g *visitor) writeElementCSS(indentLevel int, attrs []parser.Attribute) (err error) {
	return g.writeAttributesCSS(indentLevel, attrs)
}

func (g *visitor) writeElementScript(indentLevel int, attrs []parser.Attribute) (err error) {
	var scriptExpressions []string
	for _, attr := range attrs {
		scriptExpressions = append(scriptExpressions, getAttributeScripts(attr)...)
	}
	if len(scriptExpressions) == 0 {
		return
	}
	// Render the scripts before the element if required.
	// templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, a, b, c)
	if _, err = g.w.WriteIndent(indentLevel, "templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, "+strings.Join(scriptExpressions, ", ")+")\n"); err != nil {
		return err
	}
	if err = g.writeErrorHandler(indentLevel); err != nil {
		return err
	}
	return err
}

func (g *visitor) VisitElement(n *parser.Element) error {
	// Set the element name for children to use
	// TODO: this might need to be a stack if we have nested elements.
	g.elementName = n.Name
	defer func() {
		g.elementName = ""
	}()

	if len(n.Attributes) == 0 {
		// <div>
		if _, err := g.w.WriteStringLiteral(g.indents, fmt.Sprintf(`<%s>`, html.EscapeString(n.Name))); err != nil {
			return err
		}
	} else {
		attrs := copyAttributes(n.Attributes)
		// <style type="text/css"></style>
		if err := g.writeElementCSS(g.indents, attrs); err != nil {
			return err
		}
		// <script></script>
		if err := g.writeElementScript(g.indents, attrs); err != nil {
			return err
		}
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
	for _, child := range n.Children {
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

func (g *visitor) VisitBoolConstantAttribute(attr *parser.BoolConstantAttribute) (err error) {
	name := html.EscapeString(attr.Name)
	if _, err = g.w.WriteStringLiteral(g.indents, fmt.Sprintf(` %s`, name)); err != nil {
		return err
	}
	return nil
}

func (g *visitor) VisitConstantAttribute(n *parser.ConstantAttribute) error {
	name := html.EscapeString(n.Name)
	value := html.EscapeString(n.Value)
	value = escapeQuotes(value)
	if _, err := g.w.WriteStringLiteral(g.indents, fmt.Sprintf(` %s=\"%s\"`, name, value)); err != nil {
		return err
	}
	return nil
}

func (g *visitor) VisitBoolExpressionAttribute(n *parser.BoolExpressionAttribute) (err error) {
	name := html.EscapeString(n.Name)
	// if
	if _, err = g.w.WriteIndent(g.indents, `if `); err != nil {
		return err
	}
	// x == y
	var r parser.Range
	if r, err = g.w.Write(n.Expression.Value); err != nil {
		return err
	}
	g.sourceMap.Add(n.Expression, r)
	// {
	if _, err = g.w.Write(` {` + "\n"); err != nil {
		return err
	}
	{
		g.indents++
		if _, err = g.w.WriteStringLiteral(g.indents, fmt.Sprintf(` %s`, name)); err != nil {
			return err
		}
		g.indents--
	}
	// }
	if _, err = g.w.WriteIndent(g.indents, `}`+"\n"); err != nil {
		return err
	}
	return nil
}

func (g *visitor) writeExpressionAttributeValueURL(indentLevel int, attr *parser.ExpressionAttribute) (err error) {
	vn := g.createVariableName()
	// var vn templ.SafeURL =
	if _, err = g.w.WriteIndent(indentLevel, "var "+vn+" templ.SafeURL = "); err != nil {
		return err
	}
	// p.Name()
	// var r parser.Range
	if _, err := g.w.Write(attr.Expression.Value); err != nil {
		return err
	}
	// g.sourceMap.Add(attr.Expression, r)
	if _, err = g.w.Write("\n"); err != nil {
		return err
	}
	if _, err = g.w.WriteIndent(indentLevel, "_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string("+vn+")))\n"); err != nil {
		return err
	}
	return g.writeErrorHandler(indentLevel)
}

func (g *visitor) writeExpressionAttributeValueScript(indentLevel int, attr *parser.ExpressionAttribute) (err error) {
	// It's a JavaScript handler, and requires special handling, because we expect a JavaScript expression.
	vn := g.createVariableName()
	// var vn templ.ComponentScript =
	if _, err = g.w.WriteIndent(indentLevel, "var "+vn+" templ.ComponentScript = "); err != nil {
		return err
	}
	// p.Name()
	var r parser.Range
	if r, err = g.w.Write(attr.Expression.Value); err != nil {
		return err
	}
	g.sourceMap.Add(attr.Expression, r)
	if _, err = g.w.Write("\n"); err != nil {
		return err
	}
	if _, err = g.w.WriteIndent(indentLevel, "_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("+vn+".Call)\n"); err != nil {
		return err
	}
	return g.writeErrorHandler(indentLevel)
}

func (g *visitor) writeExpressionAttributeValueDefault(indentLevel int, attr *parser.ExpressionAttribute) (err error) {
	var r parser.Range
	vn := g.createVariableName()
	// var vn string
	if _, err = g.w.WriteIndent(indentLevel, "var "+vn+" string\n"); err != nil {
		return err
	}
	// vn, templ_7745c5c3_Err = templ.JoinStringErrs(
	if _, err = g.w.WriteIndent(indentLevel, vn+", templ_7745c5c3_Err = templ.JoinStringErrs("); err != nil {
		return err
	}
	// p.Name()
	if r, err = g.w.Write(attr.Expression.Value); err != nil {
		return err
	}
	g.sourceMap.Add(attr.Expression, r)
	// )
	if _, err = g.w.Write(")\n"); err != nil {
		return err
	}
	// Attribute expression error handler.
	err = g.writeExpressionErrorHandler(indentLevel, attr.Expression)
	if err != nil {
		return err
	}

	// _, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(vn)
	if _, err = g.w.WriteIndent(indentLevel, "_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString("+vn+"))\n"); err != nil {
		return err
	}
	return g.writeErrorHandler(indentLevel)
}

func (g *visitor) writeExpressionAttributeValueStyle(indentLevel int, attr *parser.ExpressionAttribute) (err error) {
	var r parser.Range
	vn := g.createVariableName()
	// var vn string
	if _, err = g.w.WriteIndent(indentLevel, "var "+vn+" string\n"); err != nil {
		return err
	}
	// vn, templ_7745c5c3_Err = templruntime.SanitizeStyleAttributeValues(
	if _, err = g.w.WriteIndent(indentLevel, vn+", templ_7745c5c3_Err = templruntime.SanitizeStyleAttributeValues("); err != nil {
		return err
	}
	// value
	if r, err = g.w.Write(attr.Expression.Value); err != nil {
		return err
	}
	g.sourceMap.Add(attr.Expression, r)
	// )
	if _, err = g.w.Write(")\n"); err != nil {
		return err
	}
	// Attribute expression error handler.
	err = g.writeExpressionErrorHandler(indentLevel, attr.Expression)
	if err != nil {
		return err
	}

	// _, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(vn))
	if _, err = g.w.WriteIndent(indentLevel, "_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString("+vn+"))\n"); err != nil {
		return err
	}
	return g.writeErrorHandler(indentLevel)
}

func (g *visitor) VisitExpressionAttribute(n *parser.ExpressionAttribute) error {
	attrName := html.EscapeString(n.Name)
	// Name
	if _, err := g.w.WriteStringLiteral(g.indents, fmt.Sprintf(` %s=`, attrName)); err != nil {
		return err
	}
	// Open quote.
	if _, err := g.w.WriteStringLiteral(g.indents, `\"`); err != nil {
		return err
	}
	// Value.
	if (g.elementName == "a" && n.Name == "href") || (g.elementName == "form" && n.Name == "action") {
		if err := g.writeExpressionAttributeValueURL(g.indents, n); err != nil {
			return err
		}
	} else if isScriptAttribute(n.Name) {
		if err := g.writeExpressionAttributeValueScript(g.indents, n); err != nil {
			return err
		}
	} else if n.Name == "style" {
		if err := g.writeExpressionAttributeValueStyle(g.indents, n); err != nil {
			return err
		}
	} else {
		if err := g.writeExpressionAttributeValueDefault(g.indents, n); err != nil {
			return err
		}
	}
	// Close quote.
	if _, err := g.w.WriteStringLiteral(g.indents, `\"`); err != nil {
		return err
	}
	return nil
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
	// var r parser.Range
	// if
	if _, err := g.w.WriteIndent(g.indents, `if `); err != nil {
		return err
	}
	// x == y {
	if _, err := g.w.Write(n.Expression.Value); err != nil {
		return err
	}
	// g.sourceMap.Add(n.Expression, r)
	// {
	if _, err := g.w.Write(` {` + "\n"); err != nil {
		return err
	}
	{
		g.indents++
		for _, child := range stripLeadingAndTrailingWhitespace(n.Then) {
			if err := child.Visit(g); err != nil {
				return fmt.Errorf("error visiting thens: %w", err)
			}
		}
		g.indents--
	}
	for _, elseIf := range n.ElseIfs {
		// } else if {
		if _, err := g.w.WriteIndent(g.indents, `} else if `); err != nil {
			return err
		}
		// x == y {
		if _, err := g.w.Write(elseIf.Expression.Value); err != nil {
			return err
		}
		// g.sourceMap.Add(elseIf.Expression, r)
		// {
		if _, err := g.w.Write(` {` + "\n"); err != nil {
			return err
		}
		{
			g.indents++
			for _, child := range stripLeadingAndTrailingWhitespace(elseIf.Then) {
				if err := child.Visit(g); err != nil {
					return fmt.Errorf("error visiting else if: %w", err)
				}
			}
			g.indents--
		}
	}
	if len(n.Else) > 0 {
		// } else {
		if _, err := g.w.WriteIndent(g.indents, `} else {`+"\n"); err != nil {
			return err
		}
		{
			g.indents++
			for _, child := range stripLeadingAndTrailingWhitespace(n.Else) {
				if err := child.Visit(g); err != nil {
					return fmt.Errorf("error visiting else: %w", err)
				}
			}
			g.indents--
		}
	}
	// }
	if _, err := g.w.WriteIndent(g.indents, `}`+"\n"); err != nil {
		return err
	}
	return nil
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

func escapeQuotes(s string) string {
	quoted := strconv.Quote(s)
	return quoted[1 : len(quoted)-1]
}

func isScriptAttribute(name string) bool {
	for _, prefix := range []string{"on", "hx-on:"} {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func copyAttributes(attr []parser.Attribute) []parser.Attribute {
	o := make([]parser.Attribute, len(attr))
	for i, a := range attr {
		if c, ok := a.(*parser.ConditionalAttribute); ok {
			c.Then = copyAttributes(c.Then)
			c.Else = copyAttributes(c.Else)
			o[i] = c
			continue
		}
		o[i] = a
	}
	return o
}

func getAttributeScripts(attr parser.Attribute) (scripts []string) {
	if attr, ok := attr.(*parser.ConditionalAttribute); ok {
		for _, attr := range attr.Then {
			scripts = append(scripts, getAttributeScripts(attr)...)
		}
		for _, attr := range attr.Else {
			scripts = append(scripts, getAttributeScripts(attr)...)
		}
	}
	if attr, ok := attr.(*parser.ExpressionAttribute); ok {
		name := html.EscapeString(attr.Name)
		if isScriptAttribute(name) {
			scripts = append(scripts, attr.Expression.Value)
		}
	}
	return scripts
}
