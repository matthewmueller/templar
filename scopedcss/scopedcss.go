package scopedcss

import (
	"bytes"

	"github.com/a-h/templ/parser/v2"
	"github.com/a-h/templ/parser/v2/visitor"
	"github.com/matthewmueller/css"
	"github.com/matthewmueller/css/scoper"
	"github.com/matthewmueller/templar/internal/classes"
	"github.com/matthewmueller/templar/internal/murmur"
)

func Rewrite(path string, tf *parser.TemplateFile) (css string, err error) {
	v := visitor.New()
	styles := new(bytes.Buffer)
	prefix := "jsx-"

	v.HTMLTemplate = func(templ *parser.HTMLTemplate) error {
		class := ""

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

		// Update elements children to exclude the <style scoped> elements
		templ.Children = children
		return nil
	}

	// Visit the template file and collect all scoped CSS styles
	if err := tf.Visit(v); err != nil {
		return "", err
	}

	return styles.String(), nil
}

func isScopedStyle(el *parser.RawElement) bool {
	if el.Name != "style" {
		return false
	}
	for _, attr := range el.Attributes {
		switch a := attr.(type) {
		case *parser.BoolConstantAttribute:
			if a.Key.String() != "scoped" {
				continue
			}
			return true
		case *parser.ConstantAttribute:
			if a.Key.String() != "scoped" {
				continue
			}
			return a.Value != ""
		}
	}
	return false
}

func updateScopedStyle(path, prefix string, style *parser.RawElement) (scopedCss, class string, err error) {
	styles := style.Contents
	if styles == "" {
		return "", "", nil
	}
	hash := murmur.Hash(styles)
	stylesheet, err := css.Parse(path, styles)
	if err != nil {
		return "", "", err
	}
	class = prefix + hash
	scoped, err := scoper.ScopeAST(path, "."+class, stylesheet)
	if err != nil {
		return "", "", err
	}
	return scoped.String(), class, nil
}
