package walk

import "go/ast"

func Decl(d ast.Decl, f func(n ast.Node)) {
	f(d)

	if d == nil {
		return
	}

	switch dd := d.(type) {
	case *ast.BadDecl:
		BadDecl(dd)
	case *ast.GenDecl:
		GenDecl(dd, f)
	case *ast.FuncDecl:
		FuncDecl(dd, f)
	}
}

func BadDecl(_ *ast.BadDecl) {}

func GenDecl(d *ast.GenDecl, f func(n ast.Node)) {
	for _, spec := range d.Specs {
		Spec(spec, f)
	}
}

func FuncDecl(d *ast.FuncDecl, f func(n ast.Node)) {
	FieldList(d.Recv, f)
	Expr(d.Name, f)
	Expr(d.Type, f)
}
