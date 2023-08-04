// package mid is a short-word for intermediate
// includes functions to parse ast into more convenient types

package mid

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

type call struct {
	code     string
	selector []string
	// args     []E
	err error
}

type callVisitor struct {
	files    *token.FileSet
	stack    []*call
	complete []*call
}

func (v *callVisitor) push(expr *ast.CallExpr) {
	v.stack = append([]*call{{code: code.OfNode(expr)}}, v.stack...)
}

func (v *callVisitor) pop() {
	v.complete = append(v.complete, v.stack[0])
	v.stack = v.stack[1:]
}

func (v *callVisitor) head() *call {
	return v.stack[0]
}

func (v *callVisitor) errorf(pos token.Pos, format string, a ...any) {
	v.stack[0].err = errors.Errorf(v.files.Position(pos), format, a...)
}

func (v *callVisitor) visitExpr(e ast.Expr) {
	switch ee := e.(type) {
	case *ast.CallExpr:
		v.visitCallExpr(ee)
	case *ast.Ident:
		v.visitIdent(ee)
	case *ast.SelectorExpr:
		v.visitSelectorExpr(ee)
	case *ast.IndexExpr:
		v.visitIndexExpr(ee)
	case *ast.IndexListExpr:
		v.visitIndexListExpr(ee)
	default:
		v.errorf(e.Pos(), "unknown")
	}
}

func (v *callVisitor) visitCallExpr(e *ast.CallExpr) {
	v.push(e)
	v.visitExpr(e.Fun)

	// h := v.head()
	for _, arg := range e.Args {
		println(arg)
		// hh := ParseExpr(v.files, arg)
		// println(hh)
		// h.args = append(h.args, )
	}

	v.pop()
}

// visitIdent final step
func (v *callVisitor) visitIdent(id *ast.Ident) {
	h := v.head()
	h.selector = append(h.selector, id.Name)
}

func (v *callVisitor) visitIndexListExpr(e *ast.IndexListExpr) {
	// for _, index := range ee.Indices {
	// v.addTypeParam(ParseTypeRef(v.files, index))
	// }
	v.visitExpr(e.X)
}

func (v *callVisitor) visitIndexExpr(e *ast.IndexExpr) {
	// v.addTypeParam(ParseTypeRef(v.files, ee.Index))
	v.visitExpr(e.X)
}

func (v *callVisitor) visitSelectorExpr(e *ast.SelectorExpr) {
	v.visitExpr(e.X)
	v.visitExpr(e.Sel)
}

type (
	Call struct {
		Code     string
		PkgAlias string
		FuncName string
		// Args     []E
		// TypeParams []*ParseTypeRef
		Err      *errors.PositionedErr
		Position token.Position
	}

	Calls []*Call
)

func ParseFuncCalls(ce *ast.CallExpr, files *token.FileSet) Calls {
	cc := &callVisitor{files: files}
	cc.visitCallExpr(ce)
	c := cc.complete
	println(c)

	return nil
}

func (x *Call) Path() string {
	if x.PkgAlias == "" {
		return x.FuncName
	}
	return strings.Join([]string{x.PkgAlias, x.FuncName}, ".")
}
