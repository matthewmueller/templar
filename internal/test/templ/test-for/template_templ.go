package testfor

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func render(items []string) templ.Component {
	return templruntime.GeneratedTemplate(func(_in templruntime.GeneratedComponentInput) (_err error) {
		_w, _ctx := _in.Writer, _in.Context
		if _ctx.Err() != nil {
			return _ctx.Err()
		}
		_buf, _is_buf := templruntime.GetBuffer(_w)
		if !_is_buf {
			defer func() {
				_err2 := templruntime.ReleaseBuffer(_buf)
				if _err == nil {
					_err = _err2
				}
			}()
		}
		_ctx = templ.InitializeContext(_ctx)
		_children := templ.GetChildren(_ctx)
		if _children == nil {
			_children = templ.NopComponent
		}
		_ctx = templ.ClearChildren(_ctx)
		for _, item := range items {
			_err = templruntime.WriteString(_buf, 1, "<div>")
			if _err != nil {
				return _err
			}
			_var, _err := templ.JoinStringErrs(item)
			if _err != nil {
				return templ.Error{Err: _err, FileName: "internal/test/templ/test-for/template.templ", Line: 5, Col: 13}
			}
			_, _err = _buf.WriteString(templ.EscapeString(_var))
			if _err != nil {
				return _err
			}
			_err = templruntime.WriteString(_buf, 2, "</div>")
			if _err != nil {
				return _err
			}

		}

		return nil
	})
}
