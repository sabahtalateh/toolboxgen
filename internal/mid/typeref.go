// package mid is a short-word for intermediate
// includes functions to parse ast into more convenient types

package mid

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

type (
	TypeRef interface {
		typeRef()
		Error() error
	}

	Type struct {
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams TypeRefs
		Position   token.Position
		Declared   string
		error      error
	}

	Map struct {
		Modifiers Modifiers
		Key       TypeRef
		Value     TypeRef
		Position  token.Position
		Declared  string
		error     error
	}

	Chan struct {
		Modifiers Modifiers
		Value     TypeRef
		Position  token.Position
		Declared  string
		error     error
	}

	FuncType struct {
		Modifiers Modifiers
		Params    Fields
		Results   Fields
		Position  token.Position
		Declared  string
		error     error
	}

	StructType struct {
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
		Declared  string
		error     error
	}

	InterfaceType struct {
		Modifiers Modifiers
		Fields    Fields
		Position  token.Position
		Declared  string
		error     error
	}

	TypeRefs []TypeRef
)

func (x *Type) typeRef()          {}
func (x *Map) typeRef()           {}
func (x *Chan) typeRef()          {}
func (x *FuncType) typeRef()      {}
func (x *StructType) typeRef()    {}
func (x *InterfaceType) typeRef() {}

type (
	Modifier interface {
		modifier()
	}

	Pointer struct {
		Position token.Position
	}

	Array struct {
		Sized    bool
		Position token.Position
	}

	Ellipsis struct {
		Position token.Position
	}

	Modifiers []Modifier
)

func (p *Pointer) modifier()  {}
func (a *Array) modifier()    {}
func (a *Ellipsis) modifier() {}

type (
	Field struct {
		Name     string
		Type     TypeRef
		Position token.Position
		Declared string
		error    error
	}

	Fields []*Field
)

type kind int

const (
	kUnknown kind = iota + 1
	kType
	kMap
	kChan
	kFuncType
	kStructType
	kInterfaceType
)

func ParseTypeRef(files *token.FileSet, e ast.Expr) TypeRef {
	v := newRefVisitor(files)
	v.visitExpr(e)

	switch v.kind {
	case kType:
		var (
			Package  string
			typeName string
		)

		selLen := len(v.selector)
		switch selLen {
		case 1:
			typeName = v.selector[0]
		case 2:
			Package = v.selector[0]
			typeName = v.selector[1]
		default:
			v.errorf(e.Pos(), "malformed type selector: %s", strings.Join(v.selector, "."))
		}

		return &Type{
			Modifiers:  v.modifiers,
			Package:    Package,
			TypeName:   typeName,
			TypeParams: v.typeParams,
			Position:   files.Position(e.Pos()),
			Declared:   code.OfNode(e),
			error:      v.err,
		}
	case kMap:
		return &Map{
			Modifiers: v.modifiers,
			Key:       ParseTypeRef(files, v.key),
			Value:     ParseTypeRef(files, v.value),
			Position:  files.Position(e.Pos()),
			Declared:  code.OfNode(e),
			error:     v.err,
		}
	case kChan:
		return &Chan{
			Modifiers: v.modifiers,
			Value:     ParseTypeRef(files, v.value),
			Position:  files.Position(e.Pos()),
			Declared:  code.OfNode(e),
			error:     v.err,
		}
	case kFuncType:
		return &FuncType{
			Modifiers: v.modifiers,
			Params:    v.params,
			Results:   v.results,
			Position:  files.Position(e.Pos()),
			Declared:  code.OfNode(e),
			error:     v.err,
		}
	case kStructType:
		return &StructType{
			Modifiers: v.modifiers,
			Fields:    v.fields,
			Position:  files.Position(e.Pos()),
			Declared:  code.OfNode(e),
			error:     v.err,
		}
	case kInterfaceType:
		return &InterfaceType{
			Modifiers: v.modifiers,
			Fields:    v.fields,
			Position:  files.Position(e.Pos()),
			Declared:  code.OfNode(e),
			error:     v.err,
		}
	default:
		panic("unsupported ref type")
	}
}

type refVisitor struct {
	kind kind

	// common
	modifiers Modifiers
	ellipsis  bool

	// kType
	selector   []string
	typeParams []TypeRef

	// kFuncType only
	params  Fields
	results Fields

	// kMap only
	key ast.Expr
	// kMap + kChan
	value ast.Expr

	// kStructType + kInterfaceType
	fields Fields

	files *token.FileSet
	err   error
}

func newRefVisitor(files *token.FileSet) *refVisitor {
	return &refVisitor{kind: kUnknown, files: files}
}

// visitExpr
// Expr implementations support list. - not supported. + supported
// - BadExpr
// + Ident				ex: *.TypeName
// + Ellipsis			ex: ...SomeType
// - BasicLit			ex: 1, "hello", true, ..
// - FuncLit			ex: _ = func() {}
// - CompositeLit 		ex: {a := 1}
// - ParenExpr			ex: (a + b)
// + SelectorExpr		ex: package.*
// + IndexExpr			ex: [A any]
// + IndexListExpr		ex: [A any, B any]
// - SliceExpr			ex: a[1:n]
// - TypeAssertExpr		ex: X.(type)
// - CallExpr			ex: a()
// + StarExpr			ex: *int
// - UnaryExpr			ex: !a
// - BinaryExpr			ex: a + b
// - KeyValueExpr		ex: one: "1"
// + ArrayType			ex: []int
// + StructType			ex: struct {a: string}
// + FuncType			ex: func(x string)
// + InterfaceType		ex: interface {Method()}
// + MapType			ex: map[..]..
// + ChanType			ex: chan ..
func (v *refVisitor) visitExpr(e ast.Expr) {
	switch ex := e.(type) {
	case *ast.Ellipsis:
		v.visitEllipsis(ex)
	case *ast.Ident:
		v.visitIdent(ex)
	case *ast.StarExpr:
		v.visitStarExpr(ex)
	case *ast.ArrayType:
		v.visitArrayType(ex)
	case *ast.SelectorExpr:
		v.visitSelectorExpr(ex)
	case *ast.IndexExpr:
		v.visitIndexExpr(ex)
	case *ast.IndexListExpr:
		v.visitIndexListExpr(ex)
	case *ast.FuncType:
		v.visitFuncType(ex)
	case *ast.MapType:
		v.visitMapType(ex)
	case *ast.ChanType:
		v.visitChanType(ex)
	case *ast.StructType:
		v.visitStructType(ex)
	case *ast.InterfaceType:
		v.visitInterfaceType(ex)
	default:
		v.errorf(e.Pos(), "type ref can not be %s", code.OfNode(e))
	}
}

func (v *refVisitor) visitEllipsis(e *ast.Ellipsis) {
	v.modifiers = append(v.modifiers, &Ellipsis{Position: v.files.Position(e.Pos())})
	v.visitExpr(e.Elt)
}

// visitIdent final step
func (v *refVisitor) visitIdent(e *ast.Ident) {
	v.selector = append(v.selector, e.Name)
	v.kind = kType
}

func (v *refVisitor) visitStarExpr(e *ast.StarExpr) {
	v.modifiers = append(v.modifiers, &Pointer{Position: v.files.Position(e.Pos())})
	v.visitExpr(e.X)
}

func (v *refVisitor) visitArrayType(e *ast.ArrayType) {
	mod := &Array{Position: v.files.Position(e.Pos())}
	if e.Len != nil {
		mod.Sized = true
	}

	v.modifiers = append(v.modifiers, mod)

	v.visitExpr(e.Elt)
}

func (v *refVisitor) visitSelectorExpr(e *ast.SelectorExpr) {
	v.visitExpr(e.X)
	v.visitExpr(e.Sel)
}

func (v *refVisitor) visitIndexExpr(e *ast.IndexExpr) {
	v.typeParams = append(v.typeParams, ParseTypeRef(v.files, e.Index))
	v.visitExpr(e.X)
}

func (v *refVisitor) visitIndexListExpr(e *ast.IndexListExpr) {
	for _, index := range e.Indices {
		v.typeParams = append(v.typeParams, ParseTypeRef(v.files, index))
	}
	v.visitExpr(e.X)
}

// visitFuncType final step
func (v *refVisitor) visitFuncType(e *ast.FuncType) {
	v.params = v.fieldList(e.Params)
	v.results = v.fieldList(e.Results)

	v.kind = kFuncType
}

func (v *refVisitor) fieldList(fields *ast.FieldList) Fields {
	if fields == nil {
		return nil
	}

	var res Fields

	for _, f := range fields.List {
		tr := ParseTypeRef(v.files, f.Type)
		position := v.files.Position(f.Pos())

		if len(f.Names) == 0 {
			res = append(res, &Field{Name: "", Type: tr, Position: position, Declared: code.OfNode(f)})
			continue
		}

		for _, name := range f.Names {
			declared := fmt.Sprintf("%s %s", name.Name, code.OfNode(f.Type))
			res = append(res, &Field{Name: name.Name, Type: tr, Position: position, Declared: declared})
		}
	}

	return res
}

// visitMapType final step
func (v *refVisitor) visitMapType(ex *ast.MapType) {
	v.key = ex.Key
	v.value = ex.Value

	v.kind = kMap
}

// visitChanType final step
func (v *refVisitor) visitChanType(ex *ast.ChanType) {
	v.value = ex.Value
	v.kind = kChan
}

// visitStructType final step
func (v *refVisitor) visitStructType(ex *ast.StructType) {
	v.fields = v.fieldList(ex.Fields)
	v.kind = kStructType
}

// visitInterfaceType final step
func (v *refVisitor) visitInterfaceType(ex *ast.InterfaceType) {
	v.fields = v.fieldList(ex.Methods)
	v.kind = kInterfaceType
}

func (v *refVisitor) errorf(pos token.Pos, format string, a ...any) {
	v.err = errors.Errorf(v.files.Position(pos), format, a...)
}

func (x *Type) Error() error {
	if err := x.TypeParams.Error(); err != nil {
		return err
	}

	return x.error
}

func (x *Map) Error() error {
	if err := x.Key.Error(); err != nil {
		return err
	}

	if err := x.Value.Error(); err != nil {
		return err
	}

	return x.error
}

func (x *Chan) Error() error {
	if err := x.Value.Error(); err != nil {
		return err
	}

	return x.error
}

func (x *FuncType) Error() error {
	for _, p := range x.Params {
		if err := p.Error(); err != nil {
			return err
		}
	}

	for _, r := range x.Results {
		if err := r.Error(); err != nil {
			return err
		}
	}

	return x.error
}

func (x *StructType) Error() error {
	if err := x.Fields.Error(); err != nil {
		return err
	}

	return x.error
}

func (x *InterfaceType) Error() error {
	if err := x.Fields.Error(); err != nil {
		return err
	}

	return x.error
}

func (x *Field) Error() error {
	if err := x.Type.Error(); err != nil {
		return err
	}

	return x.error
}

func (x TypeRefs) Error() error {
	for _, xx := range x {
		if err := xx.Error(); err != nil {
			return err
		}
	}

	return nil
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
