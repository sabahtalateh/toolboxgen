package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/mid"
	"go/ast"
)

func (c *Converter) Call(ctx Context, call *ast.CallExpr) {
	cc := mid.ParseFuncCalls(call, ctx.files)

	println(cc)
}
