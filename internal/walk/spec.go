package walk

import "go/ast"

func Spec(s ast.Spec, f func(n ast.Node)) {
	f(s)

	if s == nil {
		return
	}

	switch ss := s.(type) {
	case *ast.ImportSpec:
		ImportSpec(ss, f)
	case *ast.ValueSpec:
		ValueSpec(ss, f)
	case *ast.TypeSpec:
		TypeSpec(ss, f)
	}
}

func ImportSpec(s *ast.ImportSpec, f func(n ast.Node)) {
	Expr(s.Name, f)
	Expr(s.Path, f)
}

func ValueSpec(s *ast.ValueSpec, f func(n ast.Node)) {
	for _, name := range s.Names {
		Expr(name, f)
	}

	Expr(s.Type, f)

	for _, value := range s.Values {
		Expr(value, f)
	}
}

func TypeSpec(s *ast.TypeSpec, f func(n ast.Node)) {
	Expr(s.Name, f)
	FieldList(s.TypeParams, f)
	Expr(s.Type, f)
}
