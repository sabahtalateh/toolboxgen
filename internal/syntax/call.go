package syntax

import (
	"go/ast"
	"go/token"
)

type (
	Call interface {
		call()
		Error() error
	}

	BaseCall []*C

	CallExpr struct {
		Expr     *ast.CallExpr
		Position token.Position
	}

	C struct {
		Position token.Position
		err      error
	}
)

func (c BaseCall) call()  {}
func (c *CallExpr) call() {}

func (c BaseCall) Error() error {
	for _, cc := range c {
		if cc.err != nil {
			return cc.err
		}
	}

	return nil
}

func (c *CallExpr) Error() error {
	return nil
}

func ParseCall(files *token.FileSet, e *ast.CallExpr) Call {
	if !isBaseCall(e) {
		return &CallExpr{Expr: e, Position: files.Position(e.Pos())}
	}

	cc := &callVisitor{files: files}
	cc.visitCallExpr(e)

	return nil
}

type (
	call struct {
		selector []string
		typeArgs TypeExprs
		args     Args
	}

	callVisitor struct {
		calls []*call
		files *token.FileSet
	}
)

func (v *callVisitor) visitExpr(e ast.Expr, c *call) {
	switch ee := e.(type) {
	case *ast.CallExpr:
		v.visitCallExpr(ee)
	case *ast.Ident:
		v.visitIdent(ee, c)
	case *ast.SelectorExpr:
		v.visitSelectorExpr(ee, c)
	case *ast.IndexExpr:
		v.visitIndexExpr(ee, c)
	case *ast.IndexListExpr:
		v.visitIndexListExpr(ee, c)
	default:
		// can not get here because if `isBaseCall` function in `ParseCall`
	}
}

func (v *callVisitor) visitCallExpr(e *ast.CallExpr) {
	c := new(call)
	for _, arg := range e.Args {
		c.args = append(c.args, ParseArg(v.files, arg))
	}
	v.visitExpr(e.Fun, c)
	v.calls = append(v.calls, c)
}

func (v *callVisitor) visitIdent(id *ast.Ident, c *call) {
	c.selector = append(c.selector, id.Name)
}

func (v *callVisitor) visitIndexExpr(e *ast.IndexExpr, c *call) {
	c.typeArgs = append(c.typeArgs, ParseTypeExpr(v.files, e.Index))
	v.visitExpr(e.X, c)
}

func (v *callVisitor) visitIndexListExpr(e *ast.IndexListExpr, c *call) {
	for _, index := range e.Indices {
		c.typeArgs = append(c.typeArgs, ParseTypeExpr(v.files, index))
	}
	v.visitExpr(e.X, c)
}

func (v *callVisitor) visitSelectorExpr(e *ast.SelectorExpr, c *call) {
	v.visitExpr(e.X, c)
	v.visitExpr(e.Sel, c)
}

func isBaseCall(e *ast.CallExpr) bool {
	fun := e.Fun
	switch f := fun.(type) {
	case *ast.IndexExpr:
		fun = f.X
	case *ast.IndexListExpr:
		fun = f.X
	}

	switch ee := fun.(type) {
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return isBaseSelector(ee)
	}

	return false
}

func isBaseSelector(e *ast.SelectorExpr) bool {
	switch x := e.X.(type) {
	case *ast.Ident:
		return true
	case *ast.SelectorExpr:
		return isBaseSelector(x)
	case *ast.CallExpr:
		return isBaseCall(x)
	}

	return false
}
