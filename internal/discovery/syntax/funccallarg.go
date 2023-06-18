package syntax

import (
	"github.com/sabahtalateh/toolboxgen/internal/utils/strings"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

type callArg struct {
	files    *token.FileSet
	expr     ast.Expr
	code     string
	val      FunctionCallArgument
	position token.Position
	err      *errors.PositionedErr
}

type FunctionCallArgument interface {
	functionCallArgument()

	setPosition(p token.Position)
	setCode(code string)
	setErr(err *errors.PositionedErr)

	Code() string
	Position() token.Position
	FirstError() *errors.PositionedErr
}

type Ref struct {
	PkgAlias   string
	FuncName   string
	TypeParams []TypeRef
	code       string
	position   token.Position
	err        *errors.PositionedErr
}

func (c *Ref) functionCallArgument() {}

func (c *Ref) setPosition(p token.Position) {
	c.position = p
}

func (c *Ref) setCode(code string) {
	c.code = code
}

func (c *Ref) setErr(err *errors.PositionedErr) {
	c.err = err
}

func (c *Ref) Code() string {
	return c.code
}

func (c *Ref) Position() token.Position {
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

func ParseCallArg(files *token.FileSet, expr ast.Expr) FunctionCallArgument {
	a := newCallArg(files, expr)
	a.visitExpr(expr)

	res := a.val

	switch v := res.(type) {
	case *Ref:
		if v.FuncName == "" {
			v.FuncName = v.PkgAlias
			v.PkgAlias = ""
		}
	}

	res.setCode(a.code)
	res.setPosition(a.position)
	res.setErr(a.err)

	return a.val
}

func (a *callArg) addCallArgRefSelector(pos token.Pos, sel string) {
	switch ref := a.val.(type) {
	case *Ref:
		if ref.PkgAlias == "" {
			ref.PkgAlias = sel
		} else if ref.FuncName == "" {
			ref.FuncName = sel
		} else {
			a.errorf(pos, "malformed selector")
		}

		if !a.position.IsValid() {
			a.position = a.files.Position(pos)
		}
	default:
		a.errorf(pos, "not supported")
	}
}

func newCallArg(files *token.FileSet, expr ast.Expr) *callArg {
	c := &callArg{files: files, expr: expr}
	c.code, c.err = nodeToString(expr, files.Position(expr.Pos()))
	return c
}

func (a *callArg) errorf(pos token.Pos, format string, x ...any) {
	a.err = errors.Errorf(a.files.Position(pos), format, x...)
}

func (a *callArg) visitExpr(expr ast.Expr) {
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

func (a *callArg) visitBasicLit(lit *ast.BasicLit) {
	switch lit.Kind {
	case token.INT:
		if a.val == nil {
			a.val = new(Int)
			a.position = a.files.Position(lit.Pos())
		}
		switch intArg := a.val.(type) {
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
		if a.val == nil {
			a.val = new(String)
			a.position = a.files.Position(lit.Pos())
		}
		switch stringArg := a.val.(type) {
		case *String:
			stringArg.Val = strings.Unquote(lit.Value)
		default:
			a.errorf(lit.Pos(), "not supported")
		}
	default:
		a.errorf(lit.Pos(), "unsupported literal type. only `int` and `string`")
	}
}

func (a *callArg) visitIdent(id *ast.Ident) {
	if a.val == nil {
		a.val = new(Ref)
	}

	switch a.val.(type) {
	case *Ref:
		a.addCallArgRefSelector(id.Pos(), id.Name)
	default:
		a.errorf(id.Pos(), "not supported")
	}
}

func (a *callArg) visitSelectorExpr(sel *ast.SelectorExpr) {
	if a.val == nil {
		a.val = new(Ref)
	}

	switch x := sel.X.(type) {
	case *ast.Ident:
		a.visitIdent(x)
	case *ast.SelectorExpr:
		a.visitSelectorExpr(x)
	default:
		a.errorf(sel.Pos(), "not supported")
	}

	switch a.val.(type) {
	case *Ref:
		a.addCallArgRefSelector(sel.Sel.Pos(), sel.Sel.Name)
	default:
		a.errorf(sel.Pos(), "not supported")
	}
}

func (a *callArg) visitIndexExpr(ind *ast.IndexExpr) {
	switch x := ind.X.(type) {
	case *ast.Ident:
		a.visitIdent(x)
	case *ast.SelectorExpr:
		a.visitSelectorExpr(x)
	default:
		a.errorf(x.Pos(), "not supported")
	}

	switch refArg := a.val.(type) {
	case *Ref:
		refArg.TypeParams = append(refArg.TypeParams, ParseTypeRef(ind.Index, a.files))
	default:
		a.errorf(ind.Pos(), "not supported")
	}
}

func (a *callArg) visitIndexListExpr(ind *ast.IndexListExpr) {
	switch x := ind.X.(type) {
	case *ast.Ident:
		a.visitIdent(x)
	case *ast.SelectorExpr:
		a.visitSelectorExpr(x)
	default:
		a.errorf(x.Pos(), "not supported")
	}

	switch refArg := a.val.(type) {
	case *Ref:
		for _, index := range ind.Indices {
			refArg.TypeParams = append(refArg.TypeParams, ParseTypeRef(index, a.files))
		}
	default:
		a.errorf(ind.Pos(), "not supported")
	}
}
