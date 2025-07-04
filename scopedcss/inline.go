package scopedcss

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/a-h/templ/parser/v2"
	"github.com/a-h/templ/parser/v2/visitor"
	"github.com/matthewmueller/templar/internal/classes"
)

// Inline processes a template file looking for <style scoped> elements.
// When it finds them, it processes the styles and moves them to the top of the
// template file, only including them once.
func Inline(path string, tf *parser.TemplateFile, replace func(css string) string) (err error) {
	v := visitor.New()
	prefix := "css-"

	var onces []*parser.TemplateFileGoExpression

	v.HTMLTemplate = func(templ *parser.HTMLTemplate) error {
		class := ""
		styles := new(bytes.Buffer)

		// Look for <style scoped> in the children
		var children []parser.Node
		for _, child := range templ.Children {
			el, ok := child.(*parser.RawElement)
			if !ok || !isScopedStyle(el) {
				children = append(children, child)
				continue
			}
			css, cls, err := updateScopedStyle(path, prefix, el)
			if err != nil {
				return err
			}
			class = cls
			styles.WriteString(css)
			styles.WriteString("\n")
		}

		// If we didn't find a <style scoped> element, continue
		if class == "" {
			return nil
		}

		// Prepend the class to the class attribute
		if err := classes.Prepend(class, templ); err != nil {
			return err
		}

		funcName := getFuncName(templ.Expression)
		if funcName == "" {
			return fmt.Errorf("%s: expected %s to be a function", path, templ.Expression.Value)
		}

		// Create a new once expression for the scoped CSS
		onceVar := fmt.Sprintf("scopedcss_%s", funcName)
		onces = append(onces, toTemplateFileGoExpression(
			fmt.Sprintf("//Generated by scopedcss\nvar %s = templ.NewOnceHandle()", onceVar),
		))

		children = append([]parser.Node{inlineStyle(onceVar, replace(styles.String()))}, children...)

		// Update elements children to exclude the <style scoped> elements
		templ.Children = children
		return nil
	}

	// Visit the template file and collect all scoped CSS styles
	if err := tf.Visit(v); err != nil {
		return err
	}

	// Append the once expressions to the template file header
	for _, once := range onces {
		tf.Nodes = append(tf.Nodes, once)
	}

	return nil
}

func getFuncName(expr parser.Expression) string {
	index := strings.Index(expr.Value, "(")
	if index == -1 {
		return ""
	}
	return expr.Value[:index]
}

func toTemplateFileGoExpression(value string) *parser.TemplateFileGoExpression {
	return &parser.TemplateFileGoExpression{
		BeforePackage: false,
		Expression: parser.Expression{
			Value: value,
		},
	}
}

func inlineStyle(onceVar string, styles string) parser.Node {
	return &parser.TemplElementExpression{
		Expression: parser.Expression{
			Value: fmt.Sprintf("%s.Once()", onceVar),
		},
		Children: []parser.Node{
			&parser.RawElement{
				Name:     "style",
				Contents: styles,
			},
		},
	}
}
