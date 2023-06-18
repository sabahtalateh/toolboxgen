package syntax

import (
	"go/ast"
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

type TypeRef struct {
	files *token.FileSet

	Expr       ast.Expr
	Code       string
	PkgAlias   string
	TypeName   string
	Star       bool
	TypeParams []TypeRef
	Position   token.Position
	Err        *errors.PositionedErr
}

func ParseTypeRef(expr ast.Expr, files *token.FileSet) TypeRef {
	t := newTypeRef(files, expr)
	t.visitExpr(expr)
	if t.TypeName == "" {
		t.TypeName = t.PkgAlias
		t.PkgAlias = ""
	}
	return *t
}

func newTypeRef(files *token.FileSet, expr ast.Expr) *TypeRef {
	t := &TypeRef{files: files, Expr: expr}
	t.Code, t.Err = code(expr, files.Position(expr.Pos()))
	return t
}

func (t *TypeRef) addSelector(pos token.Pos, sel string) {
	if t.PkgAlias == "" {
		t.PkgAlias = sel
	} else if t.TypeName == "" {
		t.TypeName = sel
	} else {
		t.errorf(pos, "malformed selector")
	}

	if !t.Position.IsValid() {
		t.Position = t.files.Position(pos)
	}
}

func (t *TypeRef) errorf(pos token.Pos, format string, a ...any) {
	t.Err = errors.Errorf(t.files.Position(pos), format, a...)
}

func (t *TypeRef) visitExpr(expr ast.Expr) {
	switch e := expr.(type) {
	case *ast.Ident:
		t.visitIdent(e)
	case *ast.StarExpr:
		t.visitStarExpr(e)
	case *ast.SelectorExpr:
		t.visitSelectorExpr(e)
	case *ast.IndexExpr:
		t.visitIndexExpr(e)
	case *ast.IndexListExpr:
		t.visitIndexListExpr(e)
	default:
		t.errorf(e.Pos(), "not supported")
	}
}

func (t *TypeRef) visitIdent(id *ast.Ident) {
	t.addSelector(id.Pos(), id.Name)
}

func (t *TypeRef) visitStarExpr(se *ast.StarExpr) {
	t.Star = true
	t.visitExpr(se.X)
}

func (t *TypeRef) visitSelectorExpr(sel *ast.SelectorExpr) {
	switch e := sel.X.(type) {
	case *ast.Ident:
		t.visitIdent(e)
	case *ast.SelectorExpr:
		t.visitSelectorExpr(e)
	default:
		t.errorf(e.Pos(), "not supported")
	}

	t.addSelector(sel.Sel.Pos(), sel.Sel.Name)
}

func (t *TypeRef) visitIndexExpr(ind *ast.IndexExpr) {
	switch x := ind.X.(type) {
	case *ast.Ident:
		t.visitIdent(x)
	case *ast.SelectorExpr:
		t.visitSelectorExpr(x)
	default:
		t.errorf(x.Pos(), "not supported")
	}
	t.TypeParams = append(t.TypeParams, ParseTypeRef(ind.Index, t.files))
}

func (t *TypeRef) visitIndexListExpr(ii *ast.IndexListExpr) {
	switch x := ii.X.(type) {
	case *ast.Ident:
		t.visitIdent(x)
	case *ast.SelectorExpr:
		t.visitSelectorExpr(x)
	default:
		t.errorf(x.Pos(), "not supported")
	}

	for _, exp := range ii.Indices {
		t.TypeParams = append(t.TypeParams, ParseTypeRef(exp, t.files))
	}
}
