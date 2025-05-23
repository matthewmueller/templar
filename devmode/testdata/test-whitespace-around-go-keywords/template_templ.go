// Code generated by templ - DO NOT EDIT.

package testwhitespacearoundgokeywords

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import (
	"github.com/a-h/templ"
	"github.com/matthewmueller/templar/templdev"
	templruntime "github.com/a-h/templ/runtime"
	"fmt"
)

func WhitespaceIsConsistentInIf(firstIf, secondIf bool) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 1, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if firstIf {
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 2, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else if secondIf {
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 3, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 4, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 5, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

const WhitespaceIsConsistentInTrueIfExpected = `<button>Start</button> <button>If</button> <button>End</button>`
const WhitespaceIsConsistentInTrueElseIfExpected = `<button>Start</button> <button>ElseIf</button> <button>End</button>`
const WhitespaceIsConsistentInTrueElseExpected = `<button>Start</button> <button>Else</button> <button>End</button>`

func WhitespaceIsConsistentInFalseIf() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 6, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if false {
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 7, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 8, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

const WhitespaceIsConsistentInFalseIfExpected = `<button>Start</button> <button>End</button>`

func WhitespaceIsConsistentInSwitch(i int) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 9, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		switch i {
		case 1:
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 10, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		default:
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 11, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 12, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

const WhitespaceIsConsistentInOneSwitchExpected = `<button>Start</button> <button>1</button> <button>End</button>`
const WhitespaceIsConsistentInDefaultSwitchExpected = `<button>Start</button> <button>default</button> <button>End</button>`

func WhitespaceIsConsistentInSwitchNoDefault() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 13, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		switch false {
		case true:
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 14, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 15, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

const WhitespaceIsConsistentInSwitchNoDefaultExpected = `<button>Start</button> <button>End</button>`

func WhitespaceIsConsistentInFor(i int) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 16, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for j := 0; j < i; j++ {
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 17, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var6 string
			templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprint(j))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `testdata/test-whitespace-around-go-keywords/template.templ`, Line: 59, Col: 25}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 18, "")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templdev.WriteString(templ_7745c5c3_Buffer, 19, "")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

const WhitespaceIsConsistentInForZeroExpected = `<button>Start</button> <button>End</button>`
const WhitespaceIsConsistentInForOneExpected = `<button>Start</button> <button>0</button> <button>End</button>`
const WhitespaceIsConsistentInForThreeExpected = `<button>Start</button> <button>0</button> <button>1</button> <button>2</button> <button>End</button>`

var _ = templruntime.GeneratedTemplate
