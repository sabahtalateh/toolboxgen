package syntax

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

type FunctionCall struct {
	CallExpr   *ast.CallExpr
	Code       string
	PkgAlias   string
	FuncName   string
	TypeParams []TypeRef
	Args       []FunctionCallArg
	Err        *errors.PositionedErr
	Position   token.Position
}

func ParseFuncCalls(ce *ast.CallExpr, files *token.FileSet) []FunctionCall {
	state := newCalls(files)
	state.visitCallExpr(ce)
	cc := state.complete
	return cc
}

func (c FunctionCall) Path() string {
	if c.PkgAlias == "" {
		return c.FuncName
	}
	return strings.Join([]string{c.PkgAlias, c.FuncName}, ".")
}

type calls struct {
	files    *token.FileSet
	stack    []FunctionCall
	complete []FunctionCall
}

func newCalls(files *token.FileSet) *calls {
	return &calls{files: files}
}

func (s *calls) push(expr *ast.CallExpr) {
	c := FunctionCall{CallExpr: expr}
	c.Code, c.Err = code(expr, s.files.Position(expr.Pos()))
	s.stack = append([]FunctionCall{c}, s.stack...)
}

func (s *calls) pop() {
	complete := s.stack[0]
	if complete.FuncName == "" {
		complete.FuncName = complete.PkgAlias
		complete.PkgAlias = ""
	}

	s.complete = append(s.complete, complete)
	s.stack = s.stack[1:]
}

func (s *calls) errorf(pos token.Pos, format string, a ...any) {
	s.stack[0].Err = errors.Errorf(s.files.Position(pos), format, a...)
}

func (s *calls) addSelector(pos token.Pos, sel string) {
	if s.stack[0].PkgAlias == "" {
		s.stack[0].PkgAlias = sel
	} else if s.stack[0].FuncName == "" {
		s.stack[0].FuncName = sel
	} else {
		s.errorf(pos, "malformed selector")
	}

	if !s.stack[0].Position.IsValid() {
		s.stack[0].Position = s.files.Position(pos)
	}
}

func (s *calls) addArg(a FunctionCallArg) {
	s.stack[0].Args = append(s.stack[0].Args, a)
}

func (s *calls) addTypeParam(t TypeRef) {
	s.stack[0].TypeParams = append(s.stack[0].TypeParams, t)
}

func (s *calls) visitCallExpr(callExpr *ast.CallExpr) {
	s.push(callExpr)
	defer s.pop()

	switch fun := callExpr.Fun.(type) {
	case *ast.Ident:
		s.visitIdent(fun)
	case *ast.SelectorExpr:
		s.visitSelectorExpr(fun)
	case *ast.IndexExpr:
		s.visitIndexExpr(fun)
	case *ast.IndexListExpr:
		s.visitIndexListExpr(fun)
	default:
		s.errorf(fun.Pos(), "not supported")
	}

	for _, arg := range callExpr.Args {
		s.addArg(ParseCallArg(s.files, arg))
	}
}

func (s *calls) visitIdent(id *ast.Ident) {
	s.addSelector(id.Pos(), id.Name)
}

func (s *calls) visitIndexListExpr(ind *ast.IndexListExpr) {
	switch x := ind.X.(type) {
	case *ast.Ident:
		s.visitIdent(x)
	case *ast.SelectorExpr:
		s.visitSelectorExpr(x)
	default:
		s.errorf(x.Pos(), "not supported")
	}

	for _, index := range ind.Indices {
		s.addTypeParam(ParseTypeRef(index, s.files))
	}
}

func (s *calls) visitIndexExpr(ind *ast.IndexExpr) {
	switch x := ind.X.(type) {
	case *ast.Ident:
		s.visitIdent(x)
	case *ast.SelectorExpr:
		s.visitSelectorExpr(x)
	default:
		s.errorf(x.Pos(), "not supported")
	}
	s.addTypeParam(ParseTypeRef(ind.Index, s.files))
}

func (s *calls) visitSelectorExpr(sel *ast.SelectorExpr) {
	switch x := sel.X.(type) {
	case *ast.Ident:
		s.visitIdent(x)
	case *ast.SelectorExpr:
		s.visitSelectorExpr(x)
	case *ast.CallExpr:
		s.visitCallExpr(x)
	default:
		s.errorf(x.Pos(), "not supported")
	}

	s.addSelector(sel.Sel.Pos(), sel.Sel.Name)
}
