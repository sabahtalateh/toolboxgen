package convert

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/life4/genesis/slices"
	"github.com/sabahtalateh/toolboxgen/internal/builtin"
	"github.com/sabahtalateh/toolboxgen/internal/context"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

func (c *Converter) TypeRef(
	ctx context.Context,
	astTypeRef syntax.TypeRef,
	definedTypeParams []tool.TypeParam,
) (tool.TypeRef, *errors.PositionedErr) {
	if err := astTypeRef.Error(); err != nil {
		return nil, err
	}

	definedTypeParamsNames := map[string]struct{}{}
	slices.Each(definedTypeParams, func(el tool.TypeParam) { definedTypeParamsNames[el.Name] = struct{}{} })

	if astTypeRef.PkgAlias == "" {
		if _, ok := definedTypeParamsNames[astTypeRef.TypeName]; ok {
			return &tool.TypeParamRef{
				Name:     astTypeRef.TypeName,
				Pointer:  astTypeRef.Star,
				Position: astTypeRef.Position,
			}, nil
		}
	}

	pkg := ""
	// calc package path for non-builtin type
	if !(astTypeRef.PkgAlias == "" && builtin.Type(astTypeRef.TypeName)) {
		var pathErr error
		pkg, pathErr = c.packagePath(ctx, astTypeRef.PkgAlias)
		if pathErr != nil {
			return nil, errors.Error(astTypeRef.Position, pathErr)
		}
	}

	foundType, err := c.findTypeRef(ctx.WithPackage(pkg), astTypeRef.TypeName, astTypeRef.Position)
	if err != nil {
		return nil, err
	}

	switch x := foundType.(type) {
	case *tool.BuiltinRef:
		x.Pointer = astTypeRef.Star
	case *tool.StructRef:
		x.Pointer = astTypeRef.Star
	case *tool.InterfaceRef:
		return nil, errors.Errorf(astTypeRef.Position, "interface reference not allowed")
	default:
		panic("not implemented")
	}

	if parametrized, ok := foundType.(tool.ParametrizedRef); ok {
		if parametrized.NumberOfTypeParams() != len(astTypeRef.TypeParams) {
			return nil, errors.InconsistentTypeParamsErr(astTypeRef.Position)
		}

		for i := 0; i < len(astTypeRef.TypeParams); i++ {
			actualTypeParam := astTypeRef.TypeParams[i]
			param, err := parametrized.NthTypeParam(i)
			if err != nil {
				return nil, errors.Error(actualTypeParam.Position, err)
			}

			definedParam := false
			if actualTypeParam.PkgAlias == "" {
				_, definedParam = definedTypeParamsNames[actualTypeParam.TypeName]
			}

			if definedParam {
				oldParam, err := parametrized.NthTypeParam(i)
				if err != nil {
					return nil, errors.InconsistentTypeParamsErr(astTypeRef.Position)
				}
				parametrized.RenameTypeParam(oldParam.Name, actualTypeParam.TypeName)
			} else {
				effective, pErr := c.TypeRef(ctx, astTypeRef.TypeParams[i], definedTypeParams)
				if err != nil {
					return nil, pErr
				}
				parametrized.SetEffectiveParam(param.Name, effective)
			}
		}
	}

	return foundType, nil
}

func (c *Converter) findTypeRef(
	ctx context.Context,
	typeName string,
	pos token.Position,
) (tool.TypeRef, *errors.PositionedErr) {
	var (
		typ  tool.TypeRef
		pErr *errors.PositionedErr
	)

	if ctx.Package == "" && builtin.Type(typeName) {
		typ = &tool.BuiltinRef{TypeName: typeName, Position: pos}
	} else {
		pkgDir, err := c.pkgDir.Dir(ctx.Package)
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
									Package:  ctx.Package,
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
													Position: ctx.Position(tp.Pos()),
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
									Package:  ctx.Package,
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
													Position: ctx.Position(tp.Pos()),
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
		return nil, errors.Errorf(pos, "type `%s` not found at `%s` package", typeName, ctx.Package)
	}

	return typ, nil
}
