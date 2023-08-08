package syntax

import (
	"errors"
	"fmt"
	"github.com/sabahtalateh/toolboxgen/internal/code"
	"go/ast"
	"go/token"
)

var (
	ErrUnexpectedExpr = errors.New("unexpected expression")
)

type (
	TypeExpr interface {
		typeExpr()

		Error() error
		Get() GetFromTypeExpr
	}

	Type struct {
		Modifiers Modifiers
		Package   string
		TypeName  string
		TypeArgs  TypeExprs
		Position  token.Position
		Code      string
	}

	Map struct {
		Modifiers Modifiers
		Key       TypeExpr
		Value     TypeExpr
		Position  token.Position
		Code      string
	}

	Chan struct {
		Modifiers Modifiers
		Value     TypeExpr
		Position  token.Position
		Code      string
	}

	FuncType struct {
		Modifiers Modifiers
		Params    Fields
		Results   Fields
		Position  token.Position
		Code      string
	}

	StructType struct {
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
		Code      string
	}

	InterfaceType struct {
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
		Code      string
	}

	TypeExprs []TypeExpr
)

func (t *Type) typeExpr()          {}
func (t *Map) typeExpr()           {}
func (t *Chan) typeExpr()          {}
func (t *FuncType) typeExpr()      {}
func (t *StructType) typeExpr()    {}
func (t *InterfaceType) typeExpr() {}

func (t *Type) Error() error {
	if err := t.TypeArgs.Error(); err != nil {
		return err
	}

	return nil
}

func (t *Map) Error() error {
	if err := t.Key.Error(); err != nil {
		return err
	}

	if err := t.Value.Error(); err != nil {
		return err
	}

	return nil
}

func (t *Chan) Error() error {
	if err := t.Value.Error(); err != nil {
		return err
	}

	return nil
}

func (t *FuncType) Error() error {
	for _, p := range t.Params {
		if err := p.Error(); err != nil {
			return err
		}
	}

	for _, r := range t.Results {
		if err := r.Error(); err != nil {
			return err
		}
	}

	return nil
}

func (t *StructType) Error() error {
	if err := t.Fields.Error(); err != nil {
		return err
	}

	return nil
}

func (t *InterfaceType) Error() error {
	if err := t.Fields.Error(); err != nil {
		return err
	}

	return nil
}

func (t TypeExprs) Error() error {
	for _, e := range t {
		if err := e.Error(); err != nil {
			return err
		}
	}

	return nil
}

type (
	Field struct {
		Name     string
		Type     TypeExpr
		Position token.Position
		Code     string
		err      error
	}

	Fields []*Field
)

func (x *Field) Error() error {
	if err := x.Type.Error(); err != nil {
		return err
	}

	return x.err
}

func (x Fields) Error() error {
	for _, field := range x {
		if err := field.Type.Error(); err != nil {
			return err
		}
		if err := field.Error(); err != nil {
			return err
		}
	}

	return nil
}

func ParseTypeExpr(files *token.FileSet, e ast.Expr) TypeExpr {
	v := &typeExprVisitor{kind: kUnexpected, files: files}
	v.visitExpr(e)

	switch v.kind {
	case kType:
		t := &Type{
			Modifiers: v.modifiers,
			TypeArgs:  v.typeArgs,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
		}

		switch len(v.selector) {
		case 1:
			t.TypeName = v.selector[0]
		case 2:
			t.Package = v.selector[0]
			t.TypeName = v.selector[1]
		}

		return t
	case kMap:
		return &Map{
			Modifiers: v.modifiers,
			Key:       v.key,
			Value:     v.value,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
		}
	case kChan:
		return &Chan{
			Modifiers: v.modifiers,
			Value:     v.value,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
		}
	case kFuncType:
		return &FuncType{
			Modifiers: v.modifiers,
			Params:    v.params,
			Results:   v.results,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
		}
	case kStructType:
		return &StructType{
			Modifiers: v.modifiers,
			Fields:    v.fields,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
		}
	case kInterfaceType:
		return &InterfaceType{
			Modifiers: v.modifiers,
			Fields:    v.fields,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
		}
	default:
		return &UnexpectedExpr{Expr: v.unexpected, Position: files.Position(e.Pos())}
	}
}

type typeExprVisitor struct {
	kind kind

	// common
	modifiers Modifiers

	// kType
	selector []string
	typeArgs TypeExprs

	// kFuncType only
	params  Fields
	results Fields

	// kMap only
	key TypeExpr
	// kMap + kChan
	value TypeExpr

	// kStructType + kInterfaceType
	fields Fields

	// kUnexpected
	unexpected ast.Expr

	files *token.FileSet
}

// visitExpr
//
// https://github.com/golang/go/blob/release-branch.go1.21/src/go/ast/ast.go#L548
// "+" - can be part of type expression
// "-" - can NOT be part of type expression
// - ast.BadExpr
// + ast.Ident				ex: *.TypeName
// + ast.Ellipsis			ex: ...TypeName
// - ast.BasicLit			ex: 1, "hello", true, ..
// - ast.FuncLit			ex: func(x string) {}
// - ast.CompositeLit 		ex: TypeName{a: 1}, []int{1,2}, map[string]int{"1": 1}. https://go.dev/ref/spec#Composite_literals
// - ast.ParenExpr			ex: (a + b)
// + ast.SelectorExpr		ex: *.abc
// + ast.IndexExpr			ex: ..[A any]
// + ast.IndexListExpr		ex: ..[A any, B any]
// - ast.SliceExpr			ex: a[1:n]
// - ast.TypeAssertExpr		ex: TypeName.(type)
// - ast.CallExpr			ex: a()
// + ast.StarExpr			ex: *int
// - ast.UnaryExpr			ex: !a
// - ast.BinaryExpr			ex: a + b
// - ast.KeyValueExpr		ex: one: "1"
// + ast.ArrayType			ex: []TypeName
// + ast.StructType			ex: struct {a: string}
// + ast.FuncType			ex: func(x string)
// + ast.InterfaceType		ex: interface {Method()}
// + ast.MapType			ex: map[..]..
// + ast.ChanType			ex: chan ..
func (v *typeExprVisitor) visitExpr(e ast.Expr) {
	switch ee := e.(type) {
	case *ast.Ident:
		v.visitIdent(ee)
	case *ast.Ellipsis:
		v.visitEllipsis(ee)
	case *ast.SelectorExpr:
		v.visitSelectorExpr(ee)
	case *ast.IndexExpr:
		v.visitIndexExpr(ee)
	case *ast.IndexListExpr:
		v.visitIndexListExpr(ee)
	case *ast.StarExpr:
		v.visitStarExpr(ee)
	case *ast.ArrayType:
		v.visitArrayType(ee)
	case *ast.StructType:
		v.visitStructType(ee)
	case *ast.FuncType:
		v.visitFuncType(ee)
	case *ast.InterfaceType:
		v.visitInterfaceType(ee)
	case *ast.MapType:
		v.visitMapType(ee)
	case *ast.ChanType:
		v.visitChanType(ee)
	default:
		v.visitUnexpected(ee)
	}
}

func (v *typeExprVisitor) visitIdent(e *ast.Ident) {
	v.kind = kType
	v.selector = append(v.selector, e.Name)
}

func (v *typeExprVisitor) visitEllipsis(e *ast.Ellipsis) {
	v.modifiers = append(v.modifiers, &Ellipsis{Position: v.files.Position(e.Pos())})
	v.visitExpr(e.Elt)
}

func (v *typeExprVisitor) visitSelectorExpr(e *ast.SelectorExpr) {
	v.visitExpr(e.X)
	v.visitExpr(e.Sel)
}

func (v *typeExprVisitor) visitIndexExpr(e *ast.IndexExpr) {
	v.typeArgs = append(v.typeArgs, ParseTypeExpr(v.files, e.Index))
	v.visitExpr(e.X)
}

func (v *typeExprVisitor) visitIndexListExpr(e *ast.IndexListExpr) {
	for _, index := range e.Indices {
		v.typeArgs = append(v.typeArgs, ParseTypeExpr(v.files, index))
	}
	v.visitExpr(e.X)
}

func (v *typeExprVisitor) visitStarExpr(e *ast.StarExpr) {
	v.modifiers = append(v.modifiers, &Pointer{Position: v.files.Position(e.Pos())})
	v.visitExpr(e.X)
}

func (v *typeExprVisitor) visitArrayType(e *ast.ArrayType) {
	mod := &Array{Position: v.files.Position(e.Pos())}
	if e.Len != nil {
		mod.Sized = true
	}

	v.modifiers = append(v.modifiers, mod)
	v.visitExpr(e.Elt)
}

func (v *typeExprVisitor) visitStructType(e *ast.StructType) {
	v.kind = kStructType
	v.fields = v.fieldList(e.Fields)
}

func (v *typeExprVisitor) visitFuncType(e *ast.FuncType) {
	v.kind = kFuncType
	v.params = v.fieldList(e.Params)
	v.results = v.fieldList(e.Results)
}

func (v *typeExprVisitor) visitInterfaceType(e *ast.InterfaceType) {
	v.kind = kInterfaceType
	v.fields = v.fieldList(e.Methods)
}

func (v *typeExprVisitor) visitMapType(e *ast.MapType) {
	v.kind = kMap
	v.key = ParseTypeExpr(v.files, e.Key)
	v.value = ParseTypeExpr(v.files, e.Value)
}

func (v *typeExprVisitor) visitChanType(e *ast.ChanType) {
	v.kind = kChan
	v.value = ParseTypeExpr(v.files, e.Value)
}

func (v *typeExprVisitor) visitUnexpected(e ast.Expr) {
	v.kind = kUnexpected
	v.unexpected = e
}

func (v *typeExprVisitor) fieldList(fields *ast.FieldList) Fields {
	if fields == nil {
		return nil
	}

	var res Fields

	for _, f := range fields.List {
		tr := ParseTypeExpr(v.files, f.Type)
		position := v.files.Position(f.Pos())

		if len(f.Names) == 0 {
			res = append(res, &Field{Name: "", Type: tr, Position: position, Code: code.OfNode(f)})
			continue
		}

		for _, name := range f.Names {
			c := fmt.Sprintf("%s %s", name.Name, code.OfNode(f.Type))
			res = append(res, &Field{Name: name.Name, Type: tr, Position: position, Code: c})
		}
	}

	return res
}
