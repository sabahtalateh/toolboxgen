package syntax

import (
	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/code"
	"go/ast"
	"go/token"
)

type (
	Arg interface {
		arg()
	}

	BasicLit struct {
		Lit      *ast.BasicLit
		Position token.Position
		Code     string
	}

	Selector struct {
		Path     []string
		TypeArgs TypeExprs
		Position token.Position
		Code     string
	}

	ArgExpr struct {
		Expr     ast.Expr
		Position token.Position
		Code     string
	}

	Args []Arg
)

func (a *BasicLit) arg() {}
func (a *Selector) arg() {}
func (a *ArgExpr) arg()  {}

func ParseArg(files *token.FileSet, e ast.Expr) Arg {
	var (
		orig     = e
		typeArgs TypeExprs
	)
	e, typeArgs = extractTypeArgs(files, e)

	switch ee := e.(type) {
	case *ast.BasicLit:
		return &BasicLit{Lit: ee, Position: files.Position(orig.Pos()), Code: code.OfNode(orig)}
	case *ast.Ident:
		return &Selector{
			Path:     []string{ee.Name},
			TypeArgs: typeArgs,
			Position: files.Position(orig.Pos()),
			Code:     code.OfNode(orig),
		}
	case *ast.SelectorExpr:
		return argSelectorExpr(files, ee, typeArgs)
	default:
		return &ArgExpr{Expr: orig, Position: files.Position(orig.Pos()), Code: code.OfNode(orig)}
	}
}

func extractTypeArgs(files *token.FileSet, e ast.Expr) (ast.Expr, TypeExprs) {
	var typeArgs TypeExprs
	switch ee := e.(type) {
	case *ast.IndexExpr:
		typeArgs = append(typeArgs, ParseTypeExpr(files, ee.Index))
		e = ee.X
	case *ast.IndexListExpr:
		for _, index := range ee.Indices {
			typeArgs = append(typeArgs, ParseTypeExpr(files, index))
		}
		e = ee.X
	}

	return e, typeArgs
}

func argSelectorExpr(files *token.FileSet, e *ast.SelectorExpr, typeArgs TypeExprs) Arg {
	var (
		orig      = e
		identOnly = true
		path      []string
	)
LOOP:
	for {
		path = append(path, e.Sel.Name)
		switch x := e.X.(type) {
		case *ast.Ident:
			path = append(path, x.Name)
			break LOOP
		case *ast.SelectorExpr:
			e = x
		default:
			identOnly = false
		}
	}

	if !identOnly {
		return &ArgExpr{Expr: orig, Position: files.Position(orig.Pos()), Code: code.OfNode(orig)}
	}

	return &Selector{
		Path:     slices.Reverse(path),
		TypeArgs: typeArgs,
		Position: files.Position(orig.Pos()),
		Code:     code.OfNode(orig),
	}
}
