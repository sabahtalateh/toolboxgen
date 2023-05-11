package discovery

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

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

func firstError(cc []call) *parseError {
	for _, c := range cc {
		err := c.firstError()
		if err != nil {
			return err
		}
	}

	return nil
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

func (d *discovery) discoverToolsInFile(fset *token.FileSet, file *ast.File, currPkg string) ([]Tool, error) {
	var err error
	var tools []Tool

	ast.Inspect(file, func(node ast.Node) bool {
		if err != nil {
			return false
		}
		switch n := node.(type) {
		case *ast.FuncDecl:
			if !isInit(n) {
				return true
			}
		LOOP:
			for _, stmt := range n.Body.List {
				switch exprStmt := stmt.(type) {
				case *ast.ExprStmt:
					switch expr := exprStmt.X.(type) {
					case *ast.CallExpr:
						var calls []call
						parseCalls(expr, fset, file.Imports, &calls, currPkg)
						var t Tool
						t, err = d.makeTool(calls)
						if err != nil {
							break LOOP
						}
						if t != nil {
							tools = append(tools, t)
						}
					}
				}
			}
		}
		return true
	})

	return tools, err
}

func parseCalls(e *ast.CallExpr, fset *token.FileSet, imports []*ast.ImportSpec, calls *[]call, currPkg string) {
	var fun ast.Expr
	var pkg string
	var err error
	switch n := e.Fun.(type) {
	case *ast.SelectorExpr:
		switch x := n.X.(type) {
		case *ast.CallExpr:
			parseCalls(x, fset, imports, calls, currPkg)
			fun = n.Sel
		case *ast.Ident:
			pkg, err = resolvePackage(x.Name, imports, currPkg)
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
		c := callFromIdent(pkg, n, e.Args, fset, imports, currPkg)
		if c.error == nil {
			c.error = err
		}
		*calls = append(*calls, c)
	case *ast.IndexExpr:
		*calls = append(*calls, parseIndexExpr(n, e.Args, fset, imports, currPkg)...)
	default:
		*calls = append(*calls, call{
			error:    errors.New("not supported"),
			position: fset.Position(n.Pos()),
		})
	}
}

func parseIndexExpr(i *ast.IndexExpr, args []ast.Expr, fset *token.FileSet, imports []*ast.ImportSpec, currPkg string) []call {
	var calls []call
	c := call{
		typed:    true,
		position: fset.Position(i.Pos()),
	}
	switch x := i.X.(type) {
	case *ast.Ident:
		c = callFromIdent("", x, args, fset, imports, currPkg)
		c.typed = true
		c.position = fset.Position(i.Pos())
	case *ast.SelectorExpr:
		switch xx := x.X.(type) {
		case *ast.CallExpr:
			parseCalls(xx, fset, imports, &calls, currPkg)
		case *ast.Ident:
			var err error
			c.pakage, err = resolvePackage(xx.Name, imports, currPkg)
			if err != nil {
				c.error = err
			}
		default:
			c.error = errors.New("not supported")
		}
		c.funcName = x.Sel.Name
		c.arguments = parseArguments(args, fset, imports, currPkg)
		c.position = fset.Position(x.Sel.Pos())
	case *ast.CallExpr:
		c.error = errors.New("function call not supported")
	default:
		c.error = errors.New("not supported")
	}

	c.typeParameter.position = fset.Position(i.Index.Pos())
	switch index := i.Index.(type) {
	case *ast.Ident:
		c.typeParameter.pakage = currPkg
		c.typeParameter.typ = index.Name
	case *ast.SelectorExpr:
		c.typeParameter.typ = index.Sel.Name
		switch x := index.X.(type) {
		case *ast.Ident:
			var err error
			c.typeParameter.pakage, err = resolvePackage(x.Name, imports, currPkg)
			if err != nil {
				c.error = err
			}
		default:
			c.typeParameter.error = errors.New("not supported")
		}
	case *ast.CallExpr:
		c.typeParameter.error = errors.New("function call as type parameter not supported")
	default:
		c.typeParameter.error = errors.New("not supported")
	}

	return append(calls, c)
}

func callFromIdent(pkg string, i *ast.Ident, args []ast.Expr, fset *token.FileSet, imports []*ast.ImportSpec, currPkg string) call {
	return call{
		pakage:    pkg,
		funcName:  i.Name,
		arguments: parseArguments(args, fset, imports, currPkg),
		position:  fset.Position(i.Pos()),
	}
}

func parseArguments(args []ast.Expr, fset *token.FileSet, imports []*ast.ImportSpec, currPkg string) []argument {
	var res []argument

	for _, arg := range args {
		res = append(res, parseArgument(arg, fset, imports, currPkg))
	}

	return res
}

func argumentFromError(error error, pos token.Position) argument {
	return argument{
		error:    error,
		position: pos,
	}
}

func parseArgument(arg ast.Expr, fset *token.FileSet, imports []*ast.ImportSpec, currPkg string) argument {
	pos := fset.Position(arg.Pos())
	switch a := arg.(type) {
	case *ast.BasicLit:
		if a.Kind != token.STRING {
			return argumentFromError(errors.New("only string literal supported"), pos)
		} else {
			return stringArgument(unquote(a.Value), pos)
		}
	case *ast.Ident:
		return refArgument(currPkg, a.Name, pos)
	case *ast.SelectorExpr:
		var pkg string
		var err error
		funcName := a.Sel.Name

		switch x := a.X.(type) {
		case *ast.Ident:
			pkg, err = resolvePackage(x.Name, imports, currPkg)
			if err != nil {
				return argumentFromError(fmt.Errorf("error analyzing argument: %w", err), pos)
			}
		default:
			return argumentFromError(errors.New("should be in form of `package.ConstructorFunctionName`"), pos)
		}

		return refArgument(pkg, funcName, pos)
	case *ast.CallExpr:
		return argumentFromError(errors.New("function call not supported"), pos)
	default:
		return argumentFromError(errors.New("not supported"), pos)
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

func refArgument(pkg, name string, pos token.Position) argument {
	return argument{
		typ: refType,
		refArg: struct {
			pakage string
			elem   string
		}{pakage: pkg, elem: name},
		position: pos,
	}
}

func resolvePackage(pkgAlias string, imports []*ast.ImportSpec, currPkg string) (string, error) {
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
