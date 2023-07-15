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

type (
	TypeRef interface {
		ref()
		ParseError() error
	}

	Type struct {
		Declared   string
		Modifiers  Modifiers
		Package    string
		TypeName   string
		TypeParams []TypeRef
		Position   token.Position
		Error      error
	}

	Map struct {
		Declared  string
		Modifiers Modifiers
		Key       TypeRef
		Value     TypeRef
		Position  token.Position
		Error     error
	}

	Chan struct {
		Declared  string
		Modifiers Modifiers
		Value     TypeRef
		Position  token.Position
		Error     error
	}

	FuncType struct {
		Declared  string
		Modifiers Modifiers
		Params    []TypeRef
		Results   []TypeRef
		Position  token.Position
		Error     error
	}

	StructType struct {
		Declared  string
		Modifiers Modifiers
		Fields    map[string]TypeRef
		Position  token.Position
		Error     error
	}
)

func (x *Type) ref()       {}
func (x *Map) ref()        {}
func (x *Chan) ref()       {}
func (x *FuncType) ref()   {}
func (x *StructType) ref() {}

type (
	Modifier interface {
		modifier()
	}

	Modifiers []Modifier

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
)

func (p *Pointer) modifier()  {}
func (a *Array) modifier()    {}
func (a *Ellipsis) modifier() {}

type kind int

const (
	kUnknown kind = iota + 1
	kType
	kMap
	kChan
	kFuncType
	kStructType
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
			Declared:   code.OfNode(e),
			Modifiers:  v.modifiers,
			Package:    Package,
			TypeName:   typeName,
			TypeParams: v.typeParams,
			Position:   files.Position(e.Pos()),
			Error:      v.err,
		}
	case kMap:
		return &Map{
			Declared:  code.OfNode(e),
			Modifiers: v.modifiers,
			Key:       ParseTypeRef(files, v.key),
			Value:     ParseTypeRef(files, v.value),
			Position:  files.Position(e.Pos()),
			Error:     v.err,
		}
	case kChan:
		return &Chan{
			Declared:  code.OfNode(e),
			Modifiers: v.modifiers,
			Value:     ParseTypeRef(files, v.value),
			Position:  files.Position(e.Pos()),
			Error:     v.err,
		}
	case kFuncType:
		return &FuncType{
			Declared:  code.OfNode(e),
			Modifiers: v.modifiers,
			Params:    v.params,
			Results:   v.results,
			Position:  files.Position(e.Pos()),
			Error:     v.err,
		}
	case kStructType:
		structType := &StructType{
			Declared:  code.OfNode(e),
			Modifiers: v.modifiers,
			Fields:    map[string]TypeRef{},
			Position:  files.Position(e.Pos()),
			Error:     v.err,
		}
		for fName, expr := range v.fields {
			structType.Fields[fName] = ParseTypeRef(files, expr)
		}

		return structType
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
	params  []TypeRef
	results []TypeRef

	// kMap only
	key ast.Expr
	// kMap + kChan
	value ast.Expr

	// kStructType
	fields map[string]ast.Expr

	files *token.FileSet
	err   error
}

func newRefVisitor(files *token.FileSet) *refVisitor {
	return &refVisitor{kind: kUnknown, files: files, fields: map[string]ast.Expr{}}
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
// - InterfaceType		ex: interface {Method()}
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
	if e.Params != nil && e.Params.List != nil {
		for _, par := range e.Params.List {
			v.params = append(v.params, ParseTypeRef(v.files, par.Type))
		}
	}

	if e.Results != nil && e.Results.List != nil {
		for _, res := range e.Results.List {
			v.results = append(v.results, ParseTypeRef(v.files, res.Type))
		}
	}

	v.kind = kFuncType
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
	if ex.Fields != nil {
		for _, field := range ex.Fields.List {
			for _, name := range field.Names {
				v.fields[name.Name] = field.Type
			}
		}
	}
	v.kind = kStructType
}

func (v *refVisitor) errorf(pos token.Pos, format string, a ...any) {
	v.err = errors.Errorf(v.files.Position(pos), format, a...)
}

func (x *Type) ParseError() error {
	for _, param := range x.TypeParams {
		if err := param.ParseError(); err != nil {
			return err
		}
	}

	return x.Error
}

func (x *Map) ParseError() error {
	if err := x.Key.ParseError(); err != nil {
		return err
	}

	if err := x.Value.ParseError(); err != nil {
		return err
	}

	return x.Error
}

func (x *Chan) ParseError() error {
	if err := x.Value.ParseError(); err != nil {
		return err
	}

	return x.Error
}

func (x *FuncType) ParseError() error {
	for _, p := range x.Params {
		if err := p.ParseError(); err != nil {
			return err
		}
	}

	for _, r := range x.Results {
		if err := r.ParseError(); err != nil {
			return err
		}
	}

	return x.Error
}

func (x *StructType) ParseError() error {
	return x.Error
}

func Position(tr TypeRef) token.Position {
	switch x := tr.(type) {
	case *Type:
		return x.Position
	case *Map:
		return x.Position
	case *Chan:
		return x.Position
	case *FuncType:
		return x.Position
	case *StructType:
		return x.Position
	default:
		panic("unknown type reference")
	}
}
