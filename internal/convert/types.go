package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/builtin"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"github.com/sabahtalateh/toolboxgen/internal/pkgdir"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type types struct {
	mod    *mod.Module
	pkgDir *pkgdir.PkgDir
}

func (t *types) Find(pakage string, typeName string, pos token.Position, ff *token.FileSet) (tool.TypeRef, *errors.PositionedErr) {
	var (
		typ  tool.TypeRef
		pErr *errors.PositionedErr
	)

	if pakage == "" && builtin.Type(typeName) {
		typ = &tool.BuiltinRef{TypeName: typeName, Position: pos}
	} else {
		pkgDir, err := t.pkgDir.Dir(pakage)
		if err != nil {
			return nil, errors.Error(pos, err)
		}

		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, pkgDir, nil, parser.ParseComments)
		for pkgName, pkg := range pkgs {
			if strings.HasSuffix(pkgName, "_test") {
				continue
			}
			for _, file := range pkg.Files {
				ast.Inspect(file, func(node ast.Node) bool {
					if err != nil || pErr != nil || typ != nil {
						return false
					}

					switch n := node.(type) {
					case *ast.TypeSpec:
						if n.Name.Name == typeName {
							switch n.Type.(type) {
							case *ast.StructType:
								strukt := &tool.StructRef{
									Package:  pakage,
									TypeName: typeName,
									Position: pos,
								}
								if n.TypeParams != nil {
									for _, tp := range n.TypeParams.List {
										for _, name := range tp.Names {
											strukt.TypeParams.Params = append(
												strukt.TypeParams.Params,
												tool.TypeParam{
													Name:     name.Name,
													Position: ff.Position(tp.Pos()),
												},
											)
											strukt.TypeParams.Effective = append(strukt.TypeParams.Effective, nil)
										}
									}
								}
								typ = strukt
								return false
							case *ast.InterfaceType:
								iface := &tool.InterfaceRef{
									Package:  pakage,
									TypeName: typeName,
									Position: pos,
								}
								if n.TypeParams != nil {
									for _, tp := range n.TypeParams.List {
										for _, name := range tp.Names {
											iface.TypeParams.Params = append(
												iface.TypeParams.Params,
												tool.TypeParam{
													Name:     name.Name,
													Position: ff.Position(tp.Pos()),
												},
											)
											iface.TypeParams.Effective = append(iface.TypeParams.Effective, nil)
										}
									}
								}
								typ = iface
								return false
							default:
								pErr = errors.Errorf(pos, "unsupported type")
								return false
							}
						}
					}
					return true
				})
			}
		}
	}

	if pErr != nil {
		return nil, pErr
	}

	if typ == nil {
		return nil, errors.Errorf(pos, "type `%s` not found at `%s` package", typeName, pakage)
	}

	return typ, nil
}
