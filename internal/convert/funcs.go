package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/pkgdir"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/builtin"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

type funcs struct {
	mod       *mod.Module
	types     *types
	converter *Converter
	pkgDir    *pkgdir.PkgDir
}

func (f *funcs) Find(pakage string, funcName string, pos token.Position, ff *token.FileSet) (*tool.FunctionRef, *errors.PositionedErr) {
	if pakage == "" && builtin.Type(funcName) {
		return nil, errors.BuiltinFunctionErr(pos, funcName)
	}

	var (
		fun  *tool.FunctionRef
		pErr *errors.PositionedErr
	)

	pkgDir, err := f.pkgDir.Dir(pakage)
	if err != nil {
		return nil, errors.Error(pos, err)
	}

	files := token.NewFileSet()
	pkgs, err := parser.ParseDir(files, pkgDir, nil, parser.ParseComments)
	for pkgName, pkg := range pkgs {
		if strings.HasSuffix(pkgName, "_test") {
			continue
		}
		for _, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if err != nil || pErr != nil || fun != nil {
					return false
				}

				switch n := node.(type) {
				case *ast.FuncDecl:
					if n.Name.Name == funcName {
						fun = &tool.FunctionRef{
							Package:  pakage,
							FuncName: n.Name.Name,
							Position: pos,
						}

						if n.Type.TypeParams != nil {
							for _, tp := range n.Type.TypeParams.List {
								for _, name := range tp.Names {
									fun.TypeParams.Params = append(fun.TypeParams.Params, tool.TypeParam{
										Name:     name.Name,
										Position: ff.Position(tp.Pos()),
									})
									fun.TypeParams.Effective = append(fun.TypeParams.Effective, nil)
								}
							}
						}

						if n.Type.Params != nil {
							for _, funcParam := range n.Type.Params.List {
								var paramType tool.TypeRef
								paramType, pErr = f.converter.TypeRef(
									syntax.ParseTypeRef(funcParam.Type, files),
									file.Imports,
									pakage,
									files,
									fun.TypeParams.Params,
								)
								if pErr != nil {
									return false
								}
								fun.Parameters = append(fun.Parameters, paramType)
							}
						}

						if n.Type.Results != nil {
							for _, funcRes := range n.Type.Results.List {
								var resType tool.TypeRef
								resType, pErr = f.converter.TypeRef(
									syntax.ParseTypeRef(funcRes.Type, files),
									file.Imports,
									pakage,
									files,
									fun.TypeParams.Params,
								)
								if pErr != nil {
									return false
								}
								fun.Results = append(fun.Results, resType)
							}
						}
					}
				}

				return true
			})
		}
	}

	if pErr != nil {
		return nil, pErr
	}

	if fun == nil {
		return nil, errors.Errorf(pos, "function `%s` not found at `%s` package", funcName, pakage)
	}

	return fun, pErr
}
