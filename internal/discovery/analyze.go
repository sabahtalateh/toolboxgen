package discovery

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

type callSequence struct {
	pakage string
	calls  []call
}

type call struct {
	pakage        string
	typed         bool
	typeParameter typeParameter
	funcName      string
	arguments     []argument
	error         error
	position      token.Position
}

type typeParameter struct {
	pakage   string
	typ      string
	error    error
	position token.Position
}

type argumentType string

const (
	stringType argumentType = "string"
	refType    argumentType = "ref"
)

type argument struct {
	typ       argumentType
	stringArg struct {
		val string
	}
	refArg struct {
		pakage string
		elem   string
	}
	error    error
	position token.Position
}

type parseError struct {
	err      error
	position token.Position
}

func (c *call) firstError() *parseError {
	if c.error != nil {
		return &parseError{err: c.error, position: c.position}
	}

	if c.typeParameter.error != nil {
		return &parseError{err: c.typeParameter.error, position: c.typeParameter.position}
	}

	for _, a := range c.arguments {
		if a.error != nil {
			return &parseError{err: a.error, position: a.position}
		}
	}

	return nil
}

func analyzeCallExpr(e *ast.CallExpr, fset *token.FileSet, imports []*ast.ImportSpec) (callSequence, *parseError) {
	var calls []call
	analyzeCallExprRec(&calls, e, fset, imports)
	var seq callSequence

	if len(calls) > 0 {
		seq.pakage = calls[0].pakage
		seq.calls = calls
	}

	var err *parseError
	for _, c := range calls {
		if err = c.firstError(); err != nil {
			break
		}
	}

	return seq, err
}

func analyzeCallExprRec(calls *[]call, e *ast.CallExpr, fset *token.FileSet, imports []*ast.ImportSpec) {
	var fun ast.Expr
	var pkg string
	var err error
	switch n := e.Fun.(type) {
	case *ast.SelectorExpr:
		switch x := n.X.(type) {
		case *ast.CallExpr:
			analyzeCallExprRec(calls, x, fset, imports)
			fun = n.Sel
		case *ast.Ident:
			pkg, err = resolvePackage(x.Name, imports)
			fun = n.Sel
		default:
			*calls = append(*calls, call{
				error:    errors.New("only function call supported"),
				position: fset.Position(x.Pos()),
			})
			return
		}
	default:
		fun = e.Fun
	}

	switch n := fun.(type) {
	case *ast.Ident:
		c := callFromIdent(pkg, n, e.Args, fset, imports)
		if c.error == nil {
			c.error = err
		}
		*calls = append(*calls, c)
	case *ast.IndexExpr:
		*calls = append(*calls, analyzeIndexExpr(n, e.Args, fset, imports)...)
	default:
		*calls = append(*calls, call{
			error:    errors.New("not supported"),
			position: fset.Position(n.Pos()),
		})
	}
}

func analyzeIndexExpr(i *ast.IndexExpr, args []ast.Expr, fset *token.FileSet, imports []*ast.ImportSpec) []call {
	var res []call
	c := call{
		typed:    true,
		position: fset.Position(i.Pos()),
	}
	switch x := i.X.(type) {
	case *ast.Ident:
		c = callFromIdent("", x, args, fset, imports)
		c.typed = true
		c.position = fset.Position(i.Pos())
	case *ast.SelectorExpr:
		switch xx := x.X.(type) {
		case *ast.CallExpr:
			analyzeCallExprRec(&res, xx, fset, imports)
		case *ast.Ident:
			var err error
			c.pakage, err = resolvePackage(xx.Name, imports)
			if err != nil {
				c.error = err
			}
		default:
			c.error = errors.New("not supported")
		}
		c.funcName = x.Sel.Name
		c.arguments = analyzeArguments(args, fset, imports)
		c.position = fset.Position(x.Sel.Pos())
	case *ast.CallExpr:
		c.error = errors.New("function call not supported")
	default:
		c.error = errors.New("not supported")
	}

	c.typeParameter.position = fset.Position(i.Index.Pos())
	switch index := i.Index.(type) {
	case *ast.Ident:
		c.typeParameter.typ = index.Name
	case *ast.SelectorExpr:
		c.typeParameter.typ = index.Sel.Name
		switch x := index.X.(type) {
		case *ast.Ident:
			c.typeParameter.pakage = x.Name
		default:
			c.typeParameter.error = errors.New("not supported")
		}
	case *ast.CallExpr:
		c.typeParameter.error = errors.New("function call not supported")
	default:
		c.typeParameter.error = errors.New("not supported")
	}

	return append(res, c)
}

func callFromIdent(pkg string, i *ast.Ident, args []ast.Expr, fset *token.FileSet, imports []*ast.ImportSpec) call {
	return call{
		pakage:    pkg,
		funcName:  i.Name,
		arguments: analyzeArguments(args, fset, imports),
		position:  fset.Position(i.Pos()),
	}
}

func stringArgument(val string, pos token.Position) argument {
	return argument{
		typ: stringType,
		stringArg: struct {
			val string
		}{val: val},
		position: pos,
	}
}

func funcArgument(pkg, name string, pos token.Position) argument {
	return argument{
		typ: refType,
		refArg: struct {
			pakage string
			elem   string
		}{pakage: pkg, elem: name},
		position: pos,
	}
}

func analyzeArguments(args []ast.Expr, fset *token.FileSet, imports []*ast.ImportSpec) []argument {
	var res []argument

	for _, arg := range args {
		res = append(res, analyzeArgument(arg, fset, imports))
	}

	return res
}

func argumentFromError(error error, pos token.Position) argument {
	return argument{
		error:    error,
		position: pos,
	}
}

func analyzeArgument(arg ast.Expr, fset *token.FileSet, imports []*ast.ImportSpec) argument {
	pos := fset.Position(arg.Pos())
	switch a := arg.(type) {
	case *ast.BasicLit:
		if a.Kind != token.STRING {
			return argumentFromError(errors.New("only string literal supported"), pos)
		} else {
			return stringArgument(unquote(a.Value), pos)
		}
	case *ast.Ident:
		return funcArgument("", a.Name, pos)
	case *ast.SelectorExpr:
		var pkg string
		var err error
		funcName := a.Sel.Name

		switch x := a.X.(type) {
		case *ast.Ident:
			pkg, err = resolvePackage(x.Name, imports)
			if err != nil {
				return argumentFromError(fmt.Errorf("error analyzing argument: %w", err), pos)
			}
		default:
			return argumentFromError(errors.New("should be in form of `package.ConstructorFunctionName`"), pos)
		}

		return funcArgument(pkg, funcName, pos)
	case *ast.CallExpr:
		return argumentFromError(errors.New("function call not supported"), pos)
	default:
		return argumentFromError(errors.New("not supported"), pos)
	}
}

func resolvePackage(pkgAlias string, imports []*ast.ImportSpec) (string, error) {
	if pkgAlias == "" {
		return "", nil
	}

	for _, imp := range imports {
		path := unquote(imp.Path.Value)
		var impAlias string
		if imp.Name != nil {
			impAlias = imp.Name.Name
		} else {
			parts := strings.Split(path, "/")
			impAlias = parts[len(parts)-1]
		}

		if impAlias == pkgAlias {
			return path, nil
		}
	}

	return "", fmt.Errorf("unresolved import alias: %s", pkgAlias)
}

func isInit(f *ast.FuncDecl) bool {
	return f.Name.Name == "init" &&
		(f.Type.Params == nil || len(f.Type.Params.List) == 0) &&
		(f.Type.Results == nil || len(f.Type.Results.List) == 0)
}

func unquote(x string) string {
	x = strings.Trim(x, "\"")
	x = strings.Trim(x, "`")
	x = strings.Trim(x, "'")

	return x
}
