package classes

import (
	"go/ast"
	goparser "go/parser"
	"go/token"

	"github.com/a-h/templ/parser/v2"
	"github.com/a-h/templ/parser/v2/visitor"
)

var headElements = map[string]bool{
	"head":     true,
	"link":     true,
	"meta":     true,
	"style":    true,
	"title":    true,
	"base":     true,
	"script":   true,
	"noscript": true,
}

// Prepend adds a class to elements in the script
func Prepend(className string, templ *parser.HTMLTemplate) error {
	v := visitor.New()

	// Visit the element and its children
	visitElement := v.Element
	v.Element = func(el *parser.Element) error {
		// ignore anything inside <head>
		if el.Name == "head" {
			return visitElement(el)
		}
		if len(el.Name) == 0 || !isLower(el.Name[0]) || headElements[el.Name] {
			return visitElement(el)
		}
		attrs, err := updateAttributes(className, el.Attributes)
		if err != nil {
			return err
		}
		el.Attributes = attrs
		return visitElement(el)
	}

	// Start visiting
	return templ.Visit(v)
}

func updateAttributes(className string, attrs []parser.Attribute) ([]parser.Attribute, error) {
	hasClass := false
	for i, attr := range attrs {
		ok := false
		switch a := attr.(type) {
		case *parser.BoolConstantAttribute:
			attrs[i], ok = updateBoolConstantAttribute(a, className)
		case *parser.ConstantAttribute:
			attrs[i], ok = updateConstantAttribute(a, className)
		case *parser.ExpressionAttribute:
			attrs[i], ok = updateExpressionAttribute(a, className)
		case *parser.ConditionalAttribute:
			attrs[i], ok = updateConditionalAttribute(a, className)
		}
		hasClass = hasClass || ok
	}
	if !hasClass {
		attrs = append(attrs, &parser.ConstantAttribute{
			Key:   attributeKey("class"),
			Value: className,
		})
	}
	return attrs, nil
}

func updateBoolConstantAttribute(attr *parser.BoolConstantAttribute, className string) (parser.Attribute, bool) {
	if !isClass(attr.Key) {
		return attr, false
	}
	return &parser.ConstantAttribute{
		Key:   attributeKey("class"),
		Value: className,
	}, true
}

func updateConstantAttribute(attr *parser.ConstantAttribute, className string) (parser.Attribute, bool) {
	if !isClass(attr.Key) {
		return attr, false
	}
	attr.Value = className + " " + attr.Value
	return attr, true
}

func updateExpressionAttribute(attr *parser.ExpressionAttribute, className string) (parser.Attribute, bool) {
	if !isClass(attr.Key) {
		return attr, false
	}
	expr, err := goparser.ParseExpr(attr.Expression.Value)
	if err != nil {
		return attr, false
	}
	switch x := expr.(type) {
	case *ast.BasicLit:
		if x.Kind == token.STRING {
			attr.Expression.Value = prefixExpr(className, attr.Expression.Value)
			return attr, true
		}
	case *ast.Ident:
		if !isGoKeyword(x) {
			attr.Expression.Value = prefixExpr(className, attr.Expression.Value)
			return attr, true
		}
	case *ast.BinaryExpr:
		attr.Expression.Value = prefixExpr(className, "("+attr.Expression.Value+")")
		return attr, true
	}
	// fmt.Println(expr)
	return attr, false
}

// Always add class to both then and else to make sure the class is always present
// TODO: we could be smarter here, where we only add the class if the class
// is not present outside of the conditional attribute.
// Note: The browser chooses the first class it finds and ignores the rest.
func updateConditionalAttribute(attr *parser.ConditionalAttribute, className string) (parser.Attribute, bool) {
	thens, err := updateAttributes(className, attr.Then)
	if err != nil {
		return nil, false
	}
	elses, err := updateAttributes(className, attr.Else)
	if err != nil {
		return nil, false
	}
	attr.Then = thens
	attr.Else = elses
	return attr, true
}

func isClass(name parser.AttributeKey) bool {
	return name.String() == "class" || name.String() == "className"
}

func isLower(ch byte) bool {
	return ch >= 'a' && ch <= 'z'
}

func attributeKey(name string) parser.AttributeKey {
	return parser.ConstantAttributeKey{Name: name}
}

func isGoKeyword(ident *ast.Ident) bool {
	if ident == nil {
		return false
	}
	return token.Lookup(ident.Name).IsKeyword()
}

func prefixExpr(prefix string, expr string) string {
	return "\"" + prefix + " \" + " + expr
}
