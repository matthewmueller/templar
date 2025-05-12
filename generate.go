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

// func generate(w io.StringWriter, tf *parser.TemplateFile) error {
// 	runtime := &Import{
// 		Name: "tr",
// 		Path: "github.com/a-h/templ/runtime",
// 	}
// 	g := &generator{w, 0, runtime}
// 	if err := g.VisitTemplateFile(tf); err != nil {
// 		return fmt.Errorf("error visiting template file: %w", err)
// 	}
// 	return nil
// }

// type Import struct {
// 	Name string
// 	Path string
// }

// type generator struct {
// 	w       io.StringWriter
// 	indent  int
// 	runtime *Import
// }

// var _ parser.Visitor = (*generator)(nil)

// func (g *generator) VisitTemplateFile(n *parser.TemplateFile) error {
// 	// fmt.Println(valast.String(n))

// 	for _, header := range n.Header {
// 		if !header.BeforePackage {
// 			continue
// 		}
// 		if err := g.VisitTemplateFileGoExpression(header); err != nil {
// 			return fmt.Errorf("error visiting template file go expression: %w", err)
// 		}
// 	}

// 	if err := g.VisitPackage(&n.Package); err != nil {
// 		return fmt.Errorf("error visiting package: %w", err)
// 	}

// 	for _, header := range n.Header {
// 		if header.BeforePackage {
// 			continue
// 		}
// 		if err := g.VisitTemplateFileGoExpression(header); err != nil {
// 			return fmt.Errorf("error visiting template file go expression: %w", err)
// 		}
// 	}

// 	for _, node := range n.Nodes {
// 		if err := node.Visit(g); err != nil {
// 			return fmt.Errorf("error visiting node: %w", err)
// 		}
// 	}

// 	return nil
// }

// func (g *generator) VisitTemplateFileGoExpression(n *parser.TemplateFileGoExpression) error {
// 	g.w.WriteString(n.Expression.Value)
// 	return nil
// }

// func (g *generator) VisitPackage(n *parser.Package) error {
// 	g.w.WriteString(n.Expression.Value)
// 	g.w.WriteString("\n")
// 	return nil
// }

// func (g *generator) VisitWhitespace(n *parser.Whitespace) error {
// 	g.w.WriteString(n.Value)
// 	return nil
// }

// func (g *generator) VisitCSSTemplate(n *parser.CSSTemplate) error {
// 	fmt.Println("VisitCSSTemplate")
// 	return nil
// }

// func (g *generator) VisitConstantCSSProperty(n *parser.ConstantCSSProperty) error {
// 	fmt.Println("VisitConstantCSSProperty")
// 	return nil
// }

// func (g *generator) VisitExpressionCSSProperty(n *parser.ExpressionCSSProperty) error {
// 	fmt.Println("VisitExpressionCSSProperty")
// 	return nil
// }

// func (g *generator) VisitDocType(n *parser.DocType) error {
// 	fmt.Println("VisitDocType")
// 	return nil
// }

// func (g *generator) VisitHTMLTemplate(n *parser.HTMLTemplate) error {
// 	g.w.WriteString("\n\nfunc ")
// 	g.w.WriteString(n.Expression.Value)
// 	g.w.WriteString(" " + g.runtime.Name + ".Component {")
// 	g.w.WriteString("\n\t")
// 	g.w.WriteString(fmt.Sprintf(`return %[1]s.GeneratedTemplate(func(in %[1]s.GeneratedComponentInput) error {`, g.runtime.Name))
// 	g.w.WriteString("\n")
// 	for _, child := range n.Children {
// 		if err := child.Visit(g); err != nil {
// 			return fmt.Errorf("error visiting child: %w", err)
// 		}
// 	}
// 	g.w.WriteString("\treturn nil")
// 	g.w.WriteString("\n\t})")
// 	g.w.WriteString("\n}")
// 	return nil
// }

// func (g *generator) VisitText(n *parser.Text) error {
// 	fmt.Println("VisitText")
// 	return nil
// }

// func (g *generator) VisitElement(n *parser.Element) error {
// 	g.w.WriteString(fmt.Sprintf(`%s.WriteString(w, "`, g.runtime.Name))
// 	g.w.WriteString("<")
// 	g.w.WriteString(n.Name)
// 	g.w.WriteString(`")`)
// 	g.w.WriteString("\n")
// 	for _, attr := range n.Attributes {
// 		g.w.WriteString(` `)
// 		if err := attr.Visit(g); err != nil {
// 			return fmt.Errorf("error visiting attribute: %w", err)
// 		}
// 	}
// 	g.w.WriteString(`")`)
// 	g.w.WriteString("\n")
// 	return nil
// }

// func (g *generator) VisitScriptElement(n *parser.ScriptElement) error {
// 	fmt.Println("VisitScriptElement")
// 	return nil
// }

// func (g *generator) VisitRawElement(n *parser.RawElement) error {
// 	fmt.Println("VisitRawElement")
// 	return nil
// }

// func (g *generator) VisitBoolConstantAttribute(n *parser.BoolConstantAttribute) error {
// 	fmt.Println("VisitBoolConstantAttribute")
// 	return nil
// }

// func (g *generator) VisitConstantAttribute(n *parser.ConstantAttribute) error {
// 	fmt.Println("VisitConstantAttribute")
// 	return nil
// }

// func (g *generator) VisitBoolExpressionAttribute(n *parser.BoolExpressionAttribute) error {
// 	fmt.Println("VisitBoolExpressionAttribute")
// 	return nil
// }

// func (g *generator) VisitExpressionAttribute(n *parser.ExpressionAttribute) error {
// 	g.w.WriteString(fmt.Sprintf(`%s.WriteString(w, %s=%s)`, g.runtime.Name, n.Name, n.Expression.Value))
// 	fmt.Println("VisitExpressionAttribute")
// 	return nil
// }

// func (g *generator) VisitSpreadAttributes(n *parser.SpreadAttributes) error {
// 	fmt.Println("VisitSpreadAttributes")
// 	return nil
// }

// func (g *generator) VisitConditionalAttribute(n *parser.ConditionalAttribute) error {
// 	fmt.Println("VisitConditionalAttribute")
// 	return nil
// }

// func (g *generator) VisitGoComment(n *parser.GoComment) error {
// 	fmt.Println("VisitGoComment")
// 	return nil
// }

// func (g *generator) VisitHTMLComment(n *parser.HTMLComment) error {
// 	fmt.Println("VisitHTMLComment")
// 	return nil
// }

// func (g *generator) VisitCallTemplateExpression(n *parser.CallTemplateExpression) error {
// 	fmt.Println("VisitCallTemplateExpression")
// 	return nil
// }

// func (g *generator) VisitTemplElementExpression(n *parser.TemplElementExpression) error {
// 	fmt.Println("VisitTemplElementExpression")
// 	return nil
// }

// func (g *generator) VisitChildrenExpression(n *parser.ChildrenExpression) error {
// 	fmt.Println("VisitChildrenExpression")
// 	return nil
// }

// func (g *generator) VisitIfExpression(n *parser.IfExpression) error {
// 	fmt.Println("VisitIfExpression")
// 	return nil
// }

// func (g *generator) VisitSwitchExpression(n *parser.SwitchExpression) error {
// 	fmt.Println("VisitSwitchExpression")
// 	return nil
// }

// func (g *generator) VisitForExpression(n *parser.ForExpression) error {
// 	fmt.Println("VisitForExpression")
// 	return nil
// }

// func (g *generator) VisitGoCode(n *parser.GoCode) error {
// 	fmt.Println("VisitGoCode")
// 	return nil
// }

// func (g *generator) VisitStringExpression(n *parser.StringExpression) error {
// 	fmt.Println("VisitStringExpression")
// 	return nil
// }

// func (g *generator) VisitScriptTemplate(n *parser.ScriptTemplate) error {
// 	fmt.Println("VisitScriptTemplate")
// 	return nil
// }
