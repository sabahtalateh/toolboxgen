package convert

import (
	"go/ast"
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/builtin"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

func (c *Converter) FuncRef(
	ref syntax.Ref,
	currPkg string,
	imports []*ast.ImportSpec,
	files *token.FileSet,
) (*tool.FunctionRef, *errors.PositionedErr) {
	if ref.PkgAlias == "" && builtin.Func(ref.FuncName) {
		return nil, errors.BuiltinFunctionErr(ref.Position(), ref.FuncName)
	}

	pkg, pathErr := c.packagePath(ref.PkgAlias, currPkg, imports)
	if pathErr != nil {
		return nil, errors.Error(ref.Position(), pathErr)
	}

	fun, err := c.funcs.Find(pkg, ref.FuncName, ref.Position(), files)
	if err != nil {
		return nil, err
	}

	if len(fun.TypeParams.Params) != len(ref.TypeParams) {
		return nil, errors.InconsistentTypeParamsErr(ref.Position())
	}

	for i, tParam := range ref.TypeParams {
		param, pErr := fun.NthTypeParam(i)
		if err != nil {
			return nil, errors.Error(param.Position, pErr)
		}

		effective, err := c.TypeRef(tParam, imports, currPkg, files, nil)
		if err != nil {
			return nil, err
		}
		fun.SetEffectiveParam(param.Name, effective)
	}

	return fun, nil
}
