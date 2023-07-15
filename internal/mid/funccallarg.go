// package mid is a short-word for intermediate
// includes functions to parse ast into more convenient types

package mid

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/utils"
)

type FunctionCallArg struct {
	Code     string
	Val      FunctionCallArgValue
	Position token.Position
	Err      *errors.PositionedErr

	files *token.FileSet
	expr  ast.Expr
}

func (a *FunctionCallArg) Error() *errors.PositionedErr {
	return a.Err
}

type FunctionCallArgValue interface {
	functionCallArgument()

	setPosition(p token.Position)
	setCode(code string)
	setErr(err *errors.PositionedErr)

	Code() string
	Position() token.Position
	Error() *errors.PositionedErr
}

type FuncCallRef struct {
	PkgAlias string
	FuncName string
	// TypeParams []*ParseTypeRef
	code     string
	position token.Position
	err      *errors.PositionedErr
}

func (c *FuncCallRef) functionCallArgument() {}

func (c *FuncCallRef) setPosition(p token.Position) {
	c.position = p
}

func (c *FuncCallRef) setCode(code string) {
	c.code = code
}

func (c *FuncCallRef) setErr(err *errors.PositionedErr) {
	c.err = err
}

func (c *FuncCallRef) Code() string {
	return c.code
}

func (c *FuncCallRef) Position() token.Position {
	return c.position
}

type String struct {
	Val      string
	code     string
	position token.Position
	err      *errors.PositionedErr
}

func (s *String) functionCallArgument() {}

func (s *String) setPosition(p token.Position) {
	s.position = p
}

func (s *String) setCode(code string) {
	s.code = code
}

func (s *String) setErr(err *errors.PositionedErr) {
	s.err = err
}

func (s *String) Code() string {
	return s.code
}

func (s *String) Position() token.Position {
	return s.position
}

type Int struct {
	Val      int
	code     string
	position token.Position
	err      *errors.PositionedErr
}

func (i *Int) functionCallArgument() {}

func (i *Int) setPosition(p token.Position) {
	i.position = p
}

func (i *Int) setCode(code string) {
	i.code = code
}

func (i *Int) setErr(err *errors.PositionedErr) {
	i.err = err
}

func (i *Int) Code() string {
	return i.code
}

func (i *Int) Position() token.Position {
	return i.position
}

func ParseCallArg(files *token.FileSet, expr ast.Expr) FunctionCallArg {
	a := newCallArg(files, expr)
	a.visitExpr(expr)

	switch v := a.Val.(type) {
	case *FuncCallRef:
		if v.FuncName == "" {
			v.FuncName = v.PkgAlias
			v.PkgAlias = ""
		}
	}

	if a.Val != nil {
		a.Val.setCode(a.Code)
		a.Val.setPosition(a.Position)
		a.Val.setErr(a.Err)
	}

	return *a
}

func (a *FunctionCallArg) addCallArgRefSelector(pos token.Pos, sel string) {
	switch ref := a.Val.(type) {
	case *FuncCallRef:
		if ref.PkgAlias == "" {
			ref.PkgAlias = sel
		} else if ref.FuncName == "" {
			ref.FuncName = sel
		} else {
			a.errorf(pos, "malformed selector")
		}

		if !a.Position.IsValid() {
			a.Position = a.files.Position(pos)
		}
	default:
		a.errorf(pos, "not supported")
	}
}

func newCallArg(files *token.FileSet, expr ast.Expr) *FunctionCallArg {
	c := &FunctionCallArg{files: files, expr: expr}
	c.Code = code.OfNode(expr)
	return c
}

func (a *FunctionCallArg) errorf(pos token.Pos, format string, x ...any) {
	a.Err = errors.Errorf(a.files.Position(pos), format, x...)
}

func (a *FunctionCallArg) visitExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		a.visitBasicLit(e)
	case *ast.Ident:
		a.visitIdent(e)
	case *ast.SelectorExpr:
		a.visitSelectorExpr(e)
	case *ast.IndexExpr:
		a.visitIndexExpr(e)
	case *ast.IndexListExpr:
		a.visitIndexListExpr(e)
	case *ast.FuncLit:
		a.errorf(e.Pos(), "function literal not supported")
	case *ast.CompositeLit:
		a.errorf(e.Pos(), "struct not supported")
	case *ast.CallExpr:
		a.errorf(e.Pos(), "function call not supported")
	default:
		a.errorf(e.Pos(), "not supported")
	}
}

func (a *FunctionCallArg) visitBasicLit(lit *ast.BasicLit) {
	switch lit.Kind {
	case token.INT:
		if a.Val == nil {
			a.Val = new(Int)
			a.Position = a.files.Position(lit.Pos())
		}
		switch intArg := a.Val.(type) {
		case *Int:
			var err error
			intArg.Val, err = strconv.Atoi(lit.Value)
			if err != nil {
				a.errorf(lit.Pos(), "error converting int")
			}
		default:
			a.errorf(lit.Pos(), "not supported")
		}
	case token.STRING:
		if a.Val == nil {
			a.Val = new(String)
			a.Position = a.files.Position(lit.Pos())
		}
		switch stringArg := a.Val.(type) {
		case *String:
			stringArg.Val = utils.Unquote(lit.Value)
		default:
			a.errorf(lit.Pos(), "not supported")
		}
	default:
		a.errorf(lit.Pos(), "unsupported literal type. only `int` and `string`")
	}
}

func (a *FunctionCallArg) visitIdent(id *ast.Ident) {
	if a.Val == nil {
		a.Val = new(FuncCallRef)
	}

	switch a.Val.(type) {
	case *FuncCallRef:
		a.addCallArgRefSelector(id.Pos(), id.Name)
	default:
		a.errorf(id.Pos(), "not supported")
	}
}

func (a *FunctionCallArg) visitSelectorExpr(sel *ast.SelectorExpr) {
	if a.Val == nil {
		a.Val = new(FuncCallRef)
	}

	switch x := sel.X.(type) {
	case *ast.Ident:
		a.visitIdent(x)
	case *ast.SelectorExpr:
		a.visitSelectorExpr(x)
	default:
		a.errorf(sel.Pos(), "not supported")
	}

	switch a.Val.(type) {
	case *FuncCallRef:
		a.addCallArgRefSelector(sel.Sel.Pos(), sel.Sel.Name)
	default:
		a.errorf(sel.Pos(), "not supported")
	}
}

func (a *FunctionCallArg) visitIndexExpr(ind *ast.IndexExpr) {
	switch x := ind.X.(type) {
	case *ast.Ident:
		a.visitIdent(x)
	case *ast.SelectorExpr:
		a.visitSelectorExpr(x)
	default:
		a.errorf(x.Pos(), "not supported")
	}

	// switch refArg := a.Val.(type) {
	// case *FuncCallRef:
	// refArg.TypeParams = append(refArg.TypeParams, ParseTypeRef(a.files, ind.Index))
	// default:
	// 	a.errorf(ind.Pos(), "not supported")
	// }
}

func (a *FunctionCallArg) visitIndexListExpr(ind *ast.IndexListExpr) {
	switch x := ind.X.(type) {
	case *ast.Ident:
		a.visitIdent(x)
	case *ast.SelectorExpr:
		a.visitSelectorExpr(x)
	default:
		a.errorf(x.Pos(), "not supported")
	}

	// switch refArg := a.Val.(type) {
	// case *FuncCallRef:
	// 	for _, index := range ind.Indices {
	// 		refArg.TypeParams = append(refArg.TypeParams, ParseTypeRef(a.files, index))
	// 	}
	// default:
	// 	a.errorf(ind.Pos(), "not supported")
	// }
}
