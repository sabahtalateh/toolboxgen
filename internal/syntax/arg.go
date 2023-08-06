package syntax

import (
	"go/ast"
	"go/token"
)

type (
	Arg interface {
		arg()
	}
)

func ParseArg(files *token.FileSet, e *ast.CallExpr) Arg {
	return nil
}
