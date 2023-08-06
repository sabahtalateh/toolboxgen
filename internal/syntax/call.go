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

	CompositeCall struct {
		Expr     *ast.CallExpr
		Position token.Position
	}

	C struct {
		Position token.Position
		err      error
	}
)

func (c BaseCall) call()       {}
func (c *CompositeCall) call() {}

func (c BaseCall) Error() error {
	for _, cc := range c {
		if cc.err != nil {
			return cc.err
		}
	}

	return nil
}

func (c *CompositeCall) Error() error {
	return nil
}

func ParseCall(files *token.FileSet, e *ast.CallExpr) Call {
	if !isBase(e) {
		return &CompositeCall{Expr: e, Position: files.Position(e.Pos())}
	}

	cc := &callVisitor{files: files}
	cc.visitCallExpr(e)
	println(123)

	return nil
}

type (
	call struct {
		selector   []string
		typeParams TypeExprs
		err        error
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
	}
}

func (v *callVisitor) visitCallExpr(e *ast.CallExpr) {
	c := new(call)
	v.visitExpr(e.Fun, c)
	v.calls = append(v.calls, c)
}

func (v *callVisitor) visitIdent(id *ast.Ident, c *call) {
	c.selector = append(c.selector, id.Name)
}

func (v *callVisitor) visitIndexExpr(e *ast.IndexExpr, c *call) {
	c.typeParams = append(c.typeParams, ParseTypeExpr(v.files, e.Index))
	v.visitExpr(e.X, c)
}

func (v *callVisitor) visitIndexListExpr(e *ast.IndexListExpr, c *call) {
	for _, index := range e.Indices {
		c.typeParams = append(c.typeParams, ParseTypeExpr(v.files, index))
	}
	v.visitExpr(e.X, c)
}

func (v *callVisitor) visitSelectorExpr(e *ast.SelectorExpr, c *call) {
	v.visitExpr(e.X, c)
	v.visitExpr(e.Sel, c)
}

func isBase(e *ast.CallExpr) bool {
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
		switch ee.X.(type) {
		case *ast.Ident:
			return true
		}
	}

	return false
}
