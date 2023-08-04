package walk

import "go/ast"

func Expr(e ast.Expr, f func(n ast.Node)) {
	f(e)

	if e == nil {
		return
	}

	switch ee := e.(type) {
	case *ast.BadExpr:
		BadExpr(ee, f)
	case *ast.Ident:
		Ident(ee, f)
	case *ast.Ellipsis:
		Ellipsis(ee, f)
	case *ast.BasicLit:
		BasicLit(ee, f)
	case *ast.FuncLit:
		FuncLit(ee, f)
	case *ast.CompositeLit:
		CompositeLit(ee, f)
	case *ast.ParenExpr:
		ParenExpr(ee, f)
	case *ast.SelectorExpr:
		SelectorExpr(ee, f)
	case *ast.IndexExpr:
		IndexExpr(ee, f)
	case *ast.IndexListExpr:
		IndexListExpr(ee, f)
	case *ast.SliceExpr:
		SliceExpr(ee, f)
	case *ast.TypeAssertExpr:
		TypeAssertExpr(ee, f)
	case *ast.CallExpr:
		CallExpr(ee, f)
	case *ast.StarExpr:
		StarExpr(ee, f)
	case *ast.UnaryExpr:
		UnaryExpr(ee, f)
	case *ast.BinaryExpr:
		BinaryExpr(ee, f)
	case *ast.KeyValueExpr:
		KeyValueExpr(ee, f)
	case *ast.ArrayType:
		ArrayType(ee, f)
	case *ast.StructType:
		StructType(ee, f)
	case *ast.FuncType:
		FuncType(ee, f)
	case *ast.InterfaceType:
		InterfaceType(ee, f)
	case *ast.MapType:
		MapType(ee, f)
	case *ast.ChanType:
		ChanType(ee, f)
	}
}

func BadExpr(_ *ast.BadExpr, _ func(n ast.Node)) {}

func Ident(_ *ast.Ident, _ func(n ast.Node)) {}

func Ellipsis(e *ast.Ellipsis, f func(n ast.Node)) {
	Expr(e.Elt, f)
}

func BasicLit(_ *ast.BasicLit, _ func(n ast.Node)) {}

func FuncLit(e *ast.FuncLit, f func(n ast.Node)) {
	Expr(e.Type, f)
	Stmt(e.Body, f)
}

func CompositeLit(e *ast.CompositeLit, f func(n ast.Node)) {
	Expr(e.Type, f)
	for _, elt := range e.Elts {
		Expr(elt, f)
	}
}

func ParenExpr(e *ast.ParenExpr, f func(n ast.Node)) {
	Expr(e.X, f)
}

func SelectorExpr(e *ast.SelectorExpr, f func(n ast.Node)) {
	Expr(e.X, f)
	Expr(e.Sel, f)
}

func IndexExpr(e *ast.IndexExpr, f func(n ast.Node)) {
	Expr(e.X, f)
	Expr(e.Index, f)
}

func IndexListExpr(e *ast.IndexListExpr, f func(n ast.Node)) {
	Expr(e.X, f)
	for _, index := range e.Indices {
		Expr(index, f)
	}
}

func SliceExpr(e *ast.SliceExpr, f func(n ast.Node)) {
	Expr(e.X, f)
	Expr(e.Low, f)
	Expr(e.High, f)
	Expr(e.Max, f)
}

func TypeAssertExpr(e *ast.TypeAssertExpr, f func(n ast.Node)) {
	Expr(e.X, f)
	Expr(e.Type, f)
}

func CallExpr(e *ast.CallExpr, f func(n ast.Node)) {
	Expr(e.Fun, f)
	for _, arg := range e.Args {
		Expr(arg, f)
	}
}

func StarExpr(e *ast.StarExpr, f func(n ast.Node)) {
	Expr(e.X, f)
}

func UnaryExpr(e *ast.UnaryExpr, f func(n ast.Node)) {
	Expr(e.X, f)
}

func BinaryExpr(e *ast.BinaryExpr, f func(n ast.Node)) {
	Expr(e.X, f)
	Expr(e.Y, f)
}

func KeyValueExpr(e *ast.KeyValueExpr, f func(n ast.Node)) {
	Expr(e.Key, f)
	Expr(e.Value, f)
}

func ArrayType(e *ast.ArrayType, f func(n ast.Node)) {
	Expr(e.Len, f)
	Expr(e.Elt, f)
}

func StructType(e *ast.StructType, f func(n ast.Node)) {
	FieldList(e.Fields, f)
}

func FuncType(e *ast.FuncType, f func(n ast.Node)) {
	FieldList(e.TypeParams, f)
	FieldList(e.Params, f)
	FieldList(e.Results, f)
}

func InterfaceType(e *ast.InterfaceType, f func(n ast.Node)) {
	FieldList(e.Methods, f)
}

func MapType(e *ast.MapType, f func(n ast.Node)) {
	Expr(e.Key, f)
	Expr(e.Value, f)
}

func ChanType(e *ast.ChanType, f func(n ast.Node)) {
	Expr(e.Value, f)
}
