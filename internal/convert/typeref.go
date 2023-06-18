package convert

import (
	"go/ast"
	"go/token"

	"github.com/life4/genesis/slices"

	"github.com/sabahtalateh/toolboxgen/internal/builtin"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

func (c *Converter) TypeRef(
	astType syntax.TypeRef,
	imports []*ast.ImportSpec,
	currPkg string,
	files *token.FileSet,
	rename []tool.TypeParam,
) (tool.TypeRef, *errors.PositionedErr) {
	pkg := ""
	// calc package path for non-builtin type
	if !(astType.PkgAlias == "" && builtin.Type(astType.TypeName)) {
		var pathErr error
		pkg, pathErr = c.packagePath(astType.PkgAlias, currPkg, imports)
		if pathErr != nil {
			return nil, errors.Error(astType.Position, pathErr)
		}
	}

	foundType, err := c.types.Find(pkg, astType.TypeName, astType.Position, files)
	if err != nil {
		return nil, err
	}

	switch x := foundType.(type) {
	case *tool.BuiltinRef:
		x.Pointer = astType.Star
	case *tool.StructRef:
		x.Pointer = astType.Star
	case *tool.InterfaceRef:
		x.Pointer = astType.Star
	default:
		panic("not implemented")
	}

	if parametrized, ok := foundType.(tool.Parametrized); ok {
		if parametrized.NumberOfParams() != len(astType.TypeParams) {
			return nil, errors.InconsistentTypeParamsErr(astType.Position)
		}

		renameParams := map[string]struct{}{}
		slices.Each(rename, func(el tool.TypeParam) { renameParams[el.Name] = struct{}{} })

		for i := 0; i < len(astType.TypeParams); i++ {
			actualTypeParam := astType.TypeParams[i]
			param, err := parametrized.NthTypeParam(i)
			if err != nil {
				return nil, errors.Error(actualTypeParam.Position, err)
			}

			renameParam := false
			if actualTypeParam.PkgAlias == "" {
				_, renameParam = renameParams[actualTypeParam.TypeName]
			}

			if renameParam {
				oldParam, err := parametrized.NthTypeParam(i)
				if err != nil {
					return nil, errors.InconsistentTypeParamsErr(astType.Position)
				}
				parametrized.RenameTypeParam(oldParam.Name, actualTypeParam.TypeName)
			} else {
				effective, pErr := c.TypeRef(astType.TypeParams[i], imports, currPkg, files, rename)
				if err != nil {
					return nil, pErr
				}
				parametrized.SetEffectiveParam(param.Name, effective)
			}
		}
	}

	return foundType, nil
}
