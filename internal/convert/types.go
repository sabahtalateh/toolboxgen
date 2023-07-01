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
	renameTypeParams map[string]string,
) (tool.TypeRef, *errors.PositionedErr) {
	tr, err := c.typeRefRec(ctx, astTypeRef, definedTypeParams)
	if err != nil {
		return nil, err
	}

	switch t := tr.(type) {
	case tool.ParametrizedRef:
		if renameTypeParams != nil {
			for from, to := range renameTypeParams {
				t.RenameTypeParamRecursive(from, to)
			}
		}
	case *tool.TypeParamRef:
		if newName, ok := renameTypeParams[t.Name]; ok {
			t.Name = newName
		}
	}

	return tr, nil
}

func (c *Converter) typeRefRec(
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
				Mods:     astTypeRef.Modifiers,
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

	typeDef, err := c.findTypeDef(ctx.WithPackage(pkg), astTypeRef)
	if err != nil {
		return nil, err
	}

	typeRef := tool.TypeRefFromDef(typeDef)

	if parametrized, ok := typeRef.(tool.ParametrizedRef); ok {
		if parametrized.NumberOfTypeParams() != len(astTypeRef.TypeParams) {
			return nil, errors.InconsistentTypeParamsErr(astTypeRef.Position)
		}

		for i := 0; i < len(astTypeRef.TypeParams); i++ {
			actualTypeParam := astTypeRef.TypeParams[i]
			param, err := parametrized.NthTypeParam(i)
			if err != nil {
				return nil, errors.Error(actualTypeParam.Position, err)
			}

			definedTypeParam := false
			if actualTypeParam.PkgAlias == "" {
				_, definedTypeParam = definedTypeParamsNames[actualTypeParam.TypeName]
			}

			if definedTypeParam {
				oldParam, err := parametrized.NthTypeParam(i)
				if err != nil {
					return nil, errors.InconsistentTypeParamsErr(astTypeRef.Position)
				}
				parametrized.RenameTypeParam(oldParam.Name, actualTypeParam.TypeName)
			}

			effective, pErr := c.typeRefRec(ctx, actualTypeParam, definedTypeParams)
			if err != nil {
				return nil, pErr
			}
			parametrized.SetEffectiveParamRecursive(param.Name, effective)
		}
	}

	return typeRef, nil
}

func (c *Converter) findTypeDef(ctx context.Context, ref syntax.TypeRef) (tool.TypeDef, *errors.PositionedErr) {
	var (
		typ  tool.TypeDef
		pErr *errors.PositionedErr
	)

	if ctx.Package == "" && builtin.Type(ref.TypeName) {
		typ = &tool.BuiltinDef{TypeName: ref.TypeName, Position: ref.Position}
	} else {
		pkgDir, err := c.pkgDir.Dir(ctx.Package)
		if err != nil {
			return nil, errors.Error(ref.Position, err)
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
						if n.Name.Name == ref.TypeName {
							switch n.Type.(type) {
							case *ast.StructType:
								strukt := &tool.StructDef{
									Code:     ref.Code,
									Package:  ctx.Package,
									TypeName: ref.TypeName,
									Position: ref.Position,
								}
								if n.TypeParams != nil {
									for _, tp := range n.TypeParams.List {
										for _, name := range tp.Names {
											strukt.TypeParams = append(
												strukt.TypeParams,
												tool.TypeParam{
													Name:     name.Name,
													Position: ctx.Position(tp.Pos()),
												},
											)
										}
									}
								}
								typ = strukt
								return false
							case *ast.InterfaceType:
								iface := &tool.InterfaceDef{
									Code:     ref.Code,
									Package:  ctx.Package,
									TypeName: ref.TypeName,
									Position: ref.Position,
								}
								if n.TypeParams != nil {
									for _, tp := range n.TypeParams.List {
										for _, name := range tp.Names {
											iface.TypeParams = append(
												iface.TypeParams,
												tool.TypeParam{
													Name:     name.Name,
													Position: ctx.Position(tp.Pos()),
												},
											)
										}
									}
								}
								typ = iface
								return false
							default:
								pErr = errors.Errorf(ref.Position, "unsupported type")
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
		return nil, errors.Errorf(ref.Position, "type `%s` not found at `%s` package", ref.TypeName, ctx.Package)
	}

	switch x := typ.(type) {
	case *tool.BuiltinDef:
		x.Modifiers = ref.Modifiers
	case *tool.StructDef:
		x.Modifiers = ref.Modifiers
	case *tool.InterfaceDef:
		x.Modifiers = ref.Modifiers
	default:
		return nil, errors.Errorf(ref.Position, "not implemented")
	}

	return typ, nil
}
