package syntax

// import (
// 	"go/ast"
// 	"go/token"
//
// 	"github.com/life4/genesis/slices"
// )
//
// type (
// 	E interface {
// 		e()
// 	}
//
// 	BadExpr struct {
// 		Position token.Position
// 		Expr     *ast.BadExpr
// 	}
//
// 	Ident struct {
// 		Value    string
// 		Position token.Position
// 		Expr     *ast.Ident
// 	}
//
// 	Ellipsis struct {
// 		Elt      Expr
// 		Position token.Position
// 		Expr     *ast.Ellipsis
// 	}
//
// 	BasicLit struct {
// 		X        *ast.BasicLit
// 		Position token.Position
// 		Expr     *ast.BasicLit
// 	}
//
// 	FuncLit struct {
// 		Type     Expr
// 		Body     *ast.BlockStmt
// 		Position token.Position
// 		Expr     *ast.FuncLit
// 	}
//
// 	CompositeLit struct {
// 		Type     Expr
// 		Elts     []Expr
// 		Position token.Position
// 		Expr     *ast.CompositeLit
// 	}
//
// 	ParenExpr struct {
// 		X        Expr
// 		Position token.Position
// 		Expr     *ast.ParenExpr
// 	}
//
// 	IndexExpr struct {
// 		X        Expr
// 		Indices  []Expr
// 		Position token.Position
// 		Expr     ast.Expr
// 	}
//
// 	SliceExpr struct {
// 		X        Expr
// 		Low      Expr
// 		High     Expr
// 		Max      Expr
// 		Position token.Position
// 		Expr     *ast.SliceExpr
// 	}
//
// 	TypeAssertExpr struct {
// 		X        Expr
// 		Type     Expr
// 		Position token.Position
// 		Expr     *ast.TypeAssertExpr
// 	}
//
// 	CallExpr struct {
// 		Fun      Expr
// 		Args     []Expr
// 		Ellipsis bool
// 		Position token.Position
// 		Expr     *ast.CallExpr
// 	}
//
// 	StarExpr struct {
// 		X        Expr
// 		Position token.Position
// 		Expr     *ast.StarExpr
// 	}
//
// 	UnaryExpr struct {
// 		Op       token.Token
// 		X        Expr
// 		Position token.Position
// 		Expr     *ast.UnaryExpr
// 	}
//
// 	BinaryExpr struct {
// 		X        Expr
// 		Op       token.Token
// 		Y        Expr
// 		Position token.Position
// 		Expr     *ast.BinaryExpr
// 	}
//
// 	KeyValueExpr struct {
// 		Key      Expr
// 		Value    Expr
// 		Position token.Position
// 		Expr     *ast.KeyValueExpr
// 	}
//
// 	ArrayType struct {
// 		Len      Expr
// 		Elt      Expr
// 		Position token.Position
// 		Expr     *ast.ArrayType
// 	}
//
// 	StructType struct {
// 		Fields   []*Field
// 		Position token.Position
// 		Expr     *ast.StructType
// 	}
//
// 	FuncType struct {
// 		Params   []*Field
// 		Results  []*Field
// 		Position token.Position
// 		Expr     *ast.FuncType
// 	}
//
// 	InterfaceType struct {
// 		Methods  []*Field
// 		Position token.Position
// 		Expr     *ast.InterfaceType
// 	}
//
// 	MapType struct {
// 		Key      Expr
// 		Value    Expr
// 		Position token.Position
// 		Expr     *ast.MapType
// 	}
//
// 	ChanType struct {
// 		Dir      ast.ChanDir
// 		Value    Expr
// 		Position token.Position
// 		Expr     *ast.ChanType
// 	}
//
// 	Expr []E
//
// 	Field struct {
// 		Name     string
// 		Type     Expr
// 		Position token.Position
// 	}
// )
//
// func (e *BadExpr) e()        {}
// func (e *Ident) e()          {}
// func (e *Ellipsis) e()       {}
// func (e *BasicLit) e()       {}
// func (e *FuncLit) e()        {}
// func (e *CompositeLit) e()   {}
// func (e *ParenExpr) e()      {}
// func (e *IndexExpr) e()      {}
// func (e *SliceExpr) e()      {}
// func (e *TypeAssertExpr) e() {}
// func (e *CallExpr) e()       {}
// func (e *StarExpr) e()       {}
// func (e *UnaryExpr) e()      {}
// func (e *BinaryExpr) e()     {}
// func (e *KeyValueExpr) e()   {}
// func (e *ArrayType) e()      {}
// func (e *StructType) e()     {}
// func (e *FuncType) e()       {}
// func (e *InterfaceType) e()  {}
// func (e *MapType) e()        {}
// func (e *ChanType) e()       {}
// func (e *UnknownExpr) e()        {}
//
// func ParseExpr(files *token.FileSet, e ast.Expr) Expr {
// 	v := exprVisitor{files: files}
// 	v.visitExpr(e)
// 	return v.ee
// }
//
// func ParseExprs(files *token.FileSet, ee []ast.Expr) []Expr {
// 	return slices.Map(ee, func(e ast.Expr) Expr { return ParseExpr(files, e) })
// }
//
// type exprVisitor struct {
// 	ee    Expr
// 	files *token.FileSet
// 	err   error
// }
//
// func (v *exprVisitor) visitExpr(e ast.Expr) {
// 	if e == nil {
// 		v.ee = append(v.ee, nil)
// 		return
// 	}
//
// 	switch ee := e.(type) {
// 	case *ast.BadExpr:
// 		v.visitBadExpr(ee)
// 	case *ast.Ident:
// 		v.visitIdent(ee)
// 	case *ast.Ellipsis:
// 		v.visitEllipsis(ee)
// 	case *ast.BasicLit:
// 		v.visitBasicLit(ee)
// 	case *ast.FuncLit:
// 		v.visitFuncLit(ee)
// 	case *ast.CompositeLit:
// 		v.visitCompositeLit(ee)
// 	case *ast.ParenExpr:
// 		v.visitParenExpr(ee)
// 	case *ast.SelectorExpr:
// 		v.visitSelectorExpr(ee)
// 	case *ast.IndexExpr:
// 		v.visitIndexExpr(ee)
// 	case *ast.IndexListExpr:
// 		v.visitIndexListExpr(ee)
// 	case *ast.SliceExpr:
// 		v.visitSliceExpr(ee)
// 	case *ast.TypeAssertExpr:
// 		v.visitTypeAssertExpr(ee)
// 	case *ast.CallExpr:
// 		v.visitCallExpr(ee)
// 	case *ast.StarExpr:
// 		v.visitStarExpr(ee)
// 	case *ast.UnaryExpr:
// 		v.visitUnaryExpr(ee)
// 	case *ast.BinaryExpr:
// 		v.visitBinaryExpr(ee)
// 	case *ast.KeyValueExpr:
// 		v.visitKeyValueExpr(ee)
// 	case *ast.ArrayType:
// 		v.visitArrayType(ee)
// 	case *ast.StructType:
// 		v.visitStructType(ee)
// 	case *ast.FuncType:
// 		v.visitFuncType(ee)
// 	case *ast.InterfaceType:
// 		v.visitInterfaceType(ee)
// 	case *ast.MapType:
// 		v.visitMapType(ee)
// 	case *ast.ChanType:
// 		v.visitChanType(ee)
// 	default:
// 		v.visitUnknown(ee)
// 	}
// }
//
// func (v *exprVisitor) visitBadExpr(e *ast.BadExpr) {
// 	v.ee = append(v.ee, &BadExpr{Position: v.files.Position(e.Pos()), Expr: e})
// }
//
// func (v *exprVisitor) visitIdent(e *ast.Ident) {
// 	v.ee = append(v.ee, &Ident{Value: e.Name, Position: v.files.Position(e.Pos()), Expr: e})
// }
//
// func (v *exprVisitor) visitEllipsis(e *ast.Ellipsis) {
// 	v.ee = append(v.ee, &Ellipsis{
// 		// Elt:      ParseExpr(v.files, e.Elt),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
//
// 	v.visitExpr(e.Elt)
// }
//
// func (v *exprVisitor) visitBasicLit(e *ast.BasicLit) {
// 	v.ee = append(v.ee, &BasicLit{
// 		X:        e,
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitFuncLit(e *ast.FuncLit) {
// 	v.ee = append(v.ee, &FuncLit{
// 		Type:     ParseExpr(v.files, e.Type),
// 		Body:     e.Body,
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitCompositeLit(e *ast.CompositeLit) {
// 	v.ee = append(v.ee, &CompositeLit{
// 		Type:     ParseExpr(v.files, e.Type),
// 		Elts:     ParseExprs(v.files, e.Elts),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitParenExpr(e *ast.ParenExpr) {
// 	v.ee = append(v.ee, &ParenExpr{
// 		X:        ParseExpr(v.files, e.X),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitSelectorExpr(e *ast.SelectorExpr) {
// 	v.visitExpr(e.X)
// 	v.visitExpr(e.Sel)
// }
//
// func (v *exprVisitor) visitIndexExpr(e *ast.IndexExpr) {
// 	v.ee = append(v.ee, &IndexExpr{
// 		X:        ParseExpr(v.files, e.X),
// 		Indices:  []Expr{ParseExpr(v.files, e.Index)},
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitIndexListExpr(e *ast.IndexListExpr) {
// 	v.ee = append(v.ee, &IndexExpr{
// 		X:        ParseExpr(v.files, e.X),
// 		Indices:  ParseExprs(v.files, e.Indices),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitSliceExpr(e *ast.SliceExpr) {
// 	v.ee = append(v.ee, &SliceExpr{
// 		X:        ParseExpr(v.files, e.X),
// 		Low:      ParseExpr(v.files, e.Low),
// 		High:     ParseExpr(v.files, e.High),
// 		Max:      ParseExpr(v.files, e.Max),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitTypeAssertExpr(e *ast.TypeAssertExpr) {
// 	v.ee = append(v.ee, &TypeAssertExpr{
// 		X:        ParseExpr(v.files, e.X),
// 		Type:     ParseExpr(v.files, e.Type),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitCallExpr(e *ast.CallExpr) {
// 	v.visitExpr(e.Fun)
//
// 	v.ee = append(v.ee, &CallExpr{
// 		// Fun:      ParseExpr(v.files, e.Fun),
// 		Args:     ParseExprs(v.files, e.Args),
// 		Ellipsis: e.Ellipsis != token.NoPos,
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitStarExpr(e *ast.StarExpr) {
// 	v.ee = append(v.ee, &StarExpr{
// 		// X:        ParseExpr(v.files, e.X),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
//
// 	v.visitExpr(e.X)
// }
//
// func (v *exprVisitor) visitUnaryExpr(e *ast.UnaryExpr) {
// 	v.ee = append(v.ee, &UnaryExpr{
// 		Op:       e.Op,
// 		X:        ParseExpr(v.files, e.X),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitBinaryExpr(e *ast.BinaryExpr) {
// 	v.ee = append(v.ee, &BinaryExpr{
// 		X:        ParseExpr(v.files, e.X),
// 		Op:       e.Op,
// 		Y:        ParseExpr(v.files, e.Y),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitKeyValueExpr(e *ast.KeyValueExpr) {
// 	v.ee = append(v.ee, &KeyValueExpr{
// 		Key:      ParseExpr(v.files, e.Key),
// 		Value:    ParseExpr(v.files, e.Value),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitArrayType(e *ast.ArrayType) {
// 	v.ee = append(v.ee, &ArrayType{
// 		Len: ParseExpr(v.files, e.Len),
// 		// Elt:      ParseExpr(v.files, e.Elt),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
//
// 	v.visitExpr(e.Elt)
// }
//
// func (v *exprVisitor) visitStructType(e *ast.StructType) {
// 	v.ee = append(v.ee, &StructType{
// 		Fields:   v.fieldList(e.Fields),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitFuncType(e *ast.FuncType) {
// 	v.ee = append(v.ee, &FuncType{
// 		Params:   v.fieldList(e.Params),
// 		Results:  v.fieldList(e.Results),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitInterfaceType(e *ast.InterfaceType) {
// 	v.ee = append(v.ee, &InterfaceType{
// 		Methods:  v.fieldList(e.Methods),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitMapType(e *ast.MapType) {
// 	v.ee = append(v.ee, &MapType{
// 		Key:      ParseExpr(v.files, e.Key),
// 		Value:    ParseExpr(v.files, e.Value),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitChanType(e *ast.ChanType) {
// 	v.ee = append(v.ee, &ChanType{
// 		Dir:      e.Dir,
// 		Value:    ParseExpr(v.files, e.Value),
// 		Position: v.files.Position(e.Pos()),
// 		Expr:     e,
// 	})
// }
//
// func (v *exprVisitor) visitUnknown(e ast.Expr) {
// 	v.ee = append(v.ee, &UnknownExpr{Position: v.files.Position(e.Pos())})
// }
//
// func (v *exprVisitor) field(f *ast.Field) []*Field {
// 	if len(f.Names) == 0 {
// 		return []*Field{{Type: ParseExpr(v.files, f.Type), Position: v.files.Position(f.Pos())}}
// 	}
//
// 	var ff []*Field
// 	for _, name := range f.Names {
// 		ff = append(ff, &Field{Name: name.Name, Type: ParseExpr(v.files, f.Type), Position: v.files.Position(f.Pos())})
// 	}
// 	return ff
// }
//
// func (v *exprVisitor) fieldList(f *ast.FieldList) []*Field {
// 	if f == nil {
// 		return nil
// 	}
//
// 	var ff []*Field
// 	for _, field := range f.List {
// 		ff = append(ff, v.field(field)...)
// 	}
// 	return ff
// }
