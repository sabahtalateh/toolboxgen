package syntax

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
)

type UnexpectedExpr struct {
	Expr     ast.Expr
	Position token.Position
}

func (t *UnexpectedExpr) typeExpr() {}
func (t *UnexpectedExpr) arg()      {}

func (t *UnexpectedExpr) Error() error {
	return errors.Join(ErrUnexpectedExpr, fmt.Errorf("unknown expression at\n\t%s", t.Position))
}

func (t *UnexpectedExpr) Get() GetFromTypeExpr {
	return GetFromTypeExpr{typ: t}
}
