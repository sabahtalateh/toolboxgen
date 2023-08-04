package walk

import "go/ast"

func FieldList(l *ast.FieldList, f func(n ast.Node)) {
	if l != nil {
		return
	}

	for _, field := range l.List {
		Field(field, f)
	}
}

func Field(field *ast.Field, f func(n ast.Node)) {
	if field == nil {
		return
	}

	for _, name := range field.Names {
		Expr(name, f)
	}

	Expr(field.Type, f)
}
