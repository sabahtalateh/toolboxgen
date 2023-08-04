package walk

import "go/ast"

func Stmt(s ast.Stmt, f func(n ast.Node)) {
	f(s)

	if s == nil {
		return
	}

	switch ss := s.(type) {
	case *ast.BadStmt:
		BadStmt(ss, f)
	case *ast.DeclStmt:
		DeclStmt(ss, f)
	case *ast.EmptyStmt:
		EmptyStmt(ss, f)
	case *ast.LabeledStmt:
		LabeledStmt(ss, f)
	case *ast.ExprStmt:
		ExprStmt(ss, f)
	case *ast.SendStmt:
		SendStmt(ss, f)
	case *ast.IncDecStmt:
		IncDecStmt(ss, f)
	case *ast.AssignStmt:
		AssignStmt(ss, f)
	case *ast.GoStmt:
		GoStmt(ss, f)
	case *ast.DeferStmt:
		DeferStmt(ss, f)
	case *ast.ReturnStmt:
		ReturnStmt(ss, f)
	case *ast.BranchStmt:
		BranchStmt(ss, f)
	case *ast.BlockStmt:
		BlockStmt(ss, f)
	case *ast.IfStmt:
		IfStmt(ss, f)
	case *ast.CaseClause:
		CaseClause(ss, f)
	case *ast.SwitchStmt:
		SwitchStmt(ss, f)
	case *ast.TypeSwitchStmt:
		TypeSwitchStmt(ss, f)
	case *ast.CommClause:
		CommClause(ss, f)
	case *ast.SelectStmt:
		SelectStmt(ss, f)
	case *ast.ForStmt:
		ForStmt(ss, f)
	case *ast.RangeStmt:
		RangeStmt(ss, f)
	}
}

func BadStmt(_ *ast.BadStmt, _ func(n ast.Node)) {}

func DeclStmt(s *ast.DeclStmt, f func(n ast.Node)) {
	Decl(s.Decl, f)
}

func EmptyStmt(_ *ast.EmptyStmt, _ func(n ast.Node)) {}

func LabeledStmt(s *ast.LabeledStmt, f func(n ast.Node)) {
	Expr(s.Label, f)
	Stmt(s.Stmt, f)
}

func ExprStmt(s *ast.ExprStmt, f func(n ast.Node)) {
	Expr(s.X, f)
}

func SendStmt(s *ast.SendStmt, f func(n ast.Node)) {
	Expr(s.Chan, f)
	Expr(s.Value, f)
}

func IncDecStmt(s *ast.IncDecStmt, f func(n ast.Node)) {
	Expr(s.X, f)
}

func AssignStmt(s *ast.AssignStmt, f func(n ast.Node)) {
	for _, lh := range s.Lhs {
		Expr(lh, f)
	}
	for _, rh := range s.Rhs {
		Expr(rh, f)
	}
}

func GoStmt(s *ast.GoStmt, f func(n ast.Node)) {
	Expr(s.Call, f)
}

func DeferStmt(s *ast.DeferStmt, f func(n ast.Node)) {
	Expr(s.Call, f)
}

func ReturnStmt(s *ast.ReturnStmt, f func(n ast.Node)) {
	for _, result := range s.Results {
		Expr(result, f)
	}
}

func BranchStmt(s *ast.BranchStmt, f func(n ast.Node)) {
	Expr(s.Label, f)
}

func BlockStmt(s *ast.BlockStmt, f func(n ast.Node)) {
	for _, stmt := range s.List {
		Stmt(stmt, f)
	}
}

func IfStmt(s *ast.IfStmt, f func(n ast.Node)) {
	Stmt(s.Init, f)
	Expr(s.Cond, f)
	Stmt(s.Body, f)
	Stmt(s.Else, f)
}

func CaseClause(s *ast.CaseClause, f func(n ast.Node)) {
	for _, expr := range s.List {
		Expr(expr, f)
	}
	for _, stmt := range s.Body {
		Stmt(stmt, f)
	}
}

func SwitchStmt(s *ast.SwitchStmt, f func(n ast.Node)) {
	Stmt(s.Init, f)
	Expr(s.Tag, f)
	Stmt(s.Body, f)
}

func TypeSwitchStmt(s *ast.TypeSwitchStmt, f func(n ast.Node)) {
	Stmt(s.Init, f)
	Stmt(s.Assign, f)
	Stmt(s.Body, f)
}

func CommClause(s *ast.CommClause, f func(n ast.Node)) {
	Stmt(s.Comm, f)
	for _, stmt := range s.Body {
		Stmt(stmt, f)
	}
}

func SelectStmt(s *ast.SelectStmt, f func(n ast.Node)) {
	Stmt(s.Body, f)
}

func ForStmt(s *ast.ForStmt, f func(n ast.Node)) {
	Stmt(s.Init, f)
	Expr(s.Cond, f)
	Stmt(s.Post, f)
	Stmt(s.Body, f)
}

func RangeStmt(s *ast.RangeStmt, f func(n ast.Node)) {
	Expr(s.Key, f)
	Expr(s.Value, f)
	Expr(s.X, f)
	Stmt(s.Body, f)
}
