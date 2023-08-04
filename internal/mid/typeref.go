// package mid is a short-word for intermediate
// includes functions to parse ast into more convenient types

package mid

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	poserrors "github.com/sabahtalateh/toolboxgen/internal/errors"
)

var (
	ErrUnknownExpr = errors.New("unknown expression")
)

func ParseTypeRef(files *token.FileSet, e ast.Expr) TypeRef {
	v := &typeRefVisitor{kind: kUnknown, files: files}
	v.visitExpr(e)

	switch v.kind {
	case kType:
		t := &Type{
			Modifiers:  v.modifiers,
			TypeParams: v.typeParams,
			Position:   files.Position(e.Pos()),
			Code:       code.OfNode(e),
			error:      v.err,
		}

		switch len(v.selector) {
		case 1:
			t.TypeName = v.selector[0]
		case 2:
			t.Package = v.selector[0]
			t.TypeName = v.selector[1]
		default:
			v.err = poserrors.Errorf(v.files.Position(e.Pos()), "malformed selector: %s", strings.Join(v.selector, "."))
		}

		return t
	case kMap:
		return &Map{
			Modifiers: v.modifiers,
			Key:       ParseTypeRef(files, v.key),
			Value:     ParseTypeRef(files, v.value),
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
			error:     v.err,
		}
	case kChan:
		return &Chan{
			Modifiers: v.modifiers,
			Value:     ParseTypeRef(files, v.value),
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
			error:     v.err,
		}
	case kFuncType:
		return &FuncType{
			Modifiers: v.modifiers,
			Params:    v.params,
			Results:   v.results,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
			error:     v.err,
		}
	case kStructType:
		return &StructType{
			Modifiers: v.modifiers,
			Fields:    v.fields,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
			error:     v.err,
		}
	case kInterfaceType:
		return &InterfaceType{
			Modifiers: v.modifiers,
			Fields:    v.fields,
			Position:  files.Position(e.Pos()),
			Code:      code.OfNode(e),
			error:     v.err,
		}
	default:
		return &UnknownExpr{Expr: v.value, Position: files.Position(e.Pos())}
	}
}

type (
	TypeRef interface {
		typeRef()
		Error() error

		Get() GetFromTypeRef
	}

	Type struct {
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams TypeRefs
		Position   token.Position
		Code       string
		error      error
	}

	Map struct {
		Modifiers Modifiers
		Key       TypeRef
		Value     TypeRef
		Position  token.Position
		Code      string
		error     error
	}

	Chan struct {
		Modifiers Modifiers
		Value     TypeRef
		Position  token.Position
		Code      string
		error     error
	}

	FuncType struct {
		Modifiers Modifiers
		Params    Fields
		Results   Fields
		Position  token.Position
		Code      string
		error     error
	}

	StructType struct {
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
		Code      string
		error     error
	}

	InterfaceType struct {
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
		Code      string
		error     error
	}

	UnknownExpr struct {
		Expr     ast.Expr
		Position token.Position
	}

	TypeRefs []TypeRef
)

func (t *Type) typeRef()          {}
func (t *Map) typeRef()           {}
func (t *Chan) typeRef()          {}
func (t *FuncType) typeRef()      {}
func (t *StructType) typeRef()    {}
func (t *InterfaceType) typeRef() {}
func (t *UnknownExpr) typeRef()   {}

func (t *Type) Error() error {
	if err := t.TypeParams.Error(); err != nil {
		return err
	}

	return t.error
}

func (t *Map) Error() error {
	if err := t.Key.Error(); err != nil {
		return err
	}

	if err := t.Value.Error(); err != nil {
		return err
	}

	return t.error
}

func (t *Chan) Error() error {
	if err := t.Value.Error(); err != nil {
		return err
	}

	return t.error
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

	return t.error
}

func (t *StructType) Error() error {
	if err := t.Fields.Error(); err != nil {
		return err
	}

	return t.error
}

func (t *InterfaceType) Error() error {
	if err := t.Fields.Error(); err != nil {
		return err
	}

	return t.error
}

func (t *UnknownExpr) Error() error {
	return errors.Join(
		ErrUnknownExpr,
		fmt.Errorf("unknown expression at\n\t%s", t.Position),
	)
}

func (x TypeRefs) Error() error {
	for _, xx := range x {
		if err := xx.Error(); err != nil {
			return err
		}
	}

	return nil
}

type (
	Field struct {
		Name     string
		Type     TypeRef
		Position token.Position
		Code     string
		error    error
	}

	Fields []*Field
)

func (x *Field) Error() error {
	if err := x.Type.Error(); err != nil {
		return err
	}

	return x.error
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

type typeRefVisitor struct {
	kind kind

	// common
	modifiers Modifiers

	// kType
	selector   []string
	typeParams []TypeRef

	// kFuncType only
	params  Fields
	results Fields

	// kMap only
	key ast.Expr
	// kMap + kChan + kUnknown
	value ast.Expr

	// kStructType + kInterfaceType
	fields Fields

	files *token.FileSet
	err   error
}

// visitExpr
//
// https://github.com/golang/go/blob/release-branch.go1.21/src/go/ast/ast.go#L548
// "+" - can be part of type ref
// "-" - can NOT be part of type ref
// - ast.BadExpr
// + ast.Ident				ex: *.TypeName
// + ast.Ellipsis			ex: ...TypeName
// - ast.BasicLit			ex: 1, "hello", true, ..
// - ast.FuncLit			ex: func(x string) {}
// - ast.CompositeLit 		ex: TypeName{a: 1}, []int{1,2}, map[string]int{"1": 1}. https://go.dev/ref/spec#Composite_literals
// - ast.ParenExpr			ex: (a + b)
// + ast.SelectorExpr		ex: *.*
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
func (v *typeRefVisitor) visitExpr(e ast.Expr) {
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
		v.visitUnknown(ee)
	}
}

func (v *typeRefVisitor) visitIdent(e *ast.Ident) {
	v.kind = kType
	v.selector = append(v.selector, e.Name)
}

func (v *typeRefVisitor) visitEllipsis(e *ast.Ellipsis) {
	v.modifiers = append(v.modifiers, &Ellipsis2{Position: v.files.Position(e.Pos())})
	v.visitExpr(e.Elt)
}

func (v *typeRefVisitor) visitSelectorExpr(e *ast.SelectorExpr) {
	v.visitExpr(e.X)
	v.visitExpr(e.Sel)
}

func (v *typeRefVisitor) visitIndexExpr(e *ast.IndexExpr) {
	v.typeParams = append(v.typeParams, ParseTypeRef(v.files, e.Index))
	v.visitExpr(e.X)
}

func (v *typeRefVisitor) visitIndexListExpr(e *ast.IndexListExpr) {
	for _, index := range e.Indices {
		v.typeParams = append(v.typeParams, ParseTypeRef(v.files, index))
	}
	v.visitExpr(e.X)
}

func (v *typeRefVisitor) visitStarExpr(e *ast.StarExpr) {
	v.modifiers = append(v.modifiers, &Pointer{Position: v.files.Position(e.Pos())})
	v.visitExpr(e.X)
}

func (v *typeRefVisitor) visitArrayType(e *ast.ArrayType) {
	mod := &Array{Position: v.files.Position(e.Pos())}
	if e.Len != nil {
		mod.Sized = true
	}

	v.modifiers = append(v.modifiers, mod)
	v.visitExpr(e.Elt)
}

func (v *typeRefVisitor) visitStructType(e *ast.StructType) {
	v.kind = kStructType
	v.fields = v.fieldList(e.Fields)
}

func (v *typeRefVisitor) visitFuncType(e *ast.FuncType) {
	v.kind = kFuncType
	v.params = v.fieldList(e.Params)
	v.results = v.fieldList(e.Results)
}

func (v *typeRefVisitor) visitInterfaceType(e *ast.InterfaceType) {
	v.kind = kInterfaceType
	v.fields = v.fieldList(e.Methods)
}

func (v *typeRefVisitor) visitMapType(e *ast.MapType) {
	v.kind = kMap
	v.key = e.Key
	v.value = e.Value
}

func (v *typeRefVisitor) visitChanType(e *ast.ChanType) {
	v.kind = kChan
	v.value = e.Value
}

func (v *typeRefVisitor) visitUnknown(e ast.Expr) {
	v.kind = kUnknown
	v.value = e
}

func (v *typeRefVisitor) fieldList(fields *ast.FieldList) Fields {
	if fields == nil {
		return nil
	}

	var res Fields

	for _, f := range fields.List {
		tr := ParseTypeRef(v.files, f.Type)
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
