package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/syntax"
	"go/ast"
)

func (c *Converter) Call(ctx Context, call *ast.CallExpr) (string, error) {
	cc := syntax.ParseCall(ctx.files, call)
	println(cc)

	println(123)
	return "", nil
}
