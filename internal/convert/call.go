package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/syntax"
	"go/ast"
)

func (c *Converter) Call(ctx Context, call *ast.CallExpr) {
	cc := syntax.ParseFuncCalls(call, ctx.files)

	println(cc)
}
