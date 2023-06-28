package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/context"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/builtin"
	"github.com/sabahtalateh/toolboxgen/internal/discovery/syntax"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

func (c *Converter) ConvertFuncDef(ctx context.Context, f *ast.FuncDecl) (*tool.FuncDef, *errors.PositionedErr) {
	funcDef := &tool.FuncDef{
		Package:  ctx.Package,
		FuncName: f.Name.Name,
		Position: ctx.Position(f.Pos()),
	}

	var err *errors.PositionedErr
	funcDef.Receiver, err = c.convFuncDefRecv(ctx, f)
	if err != nil {
		return nil, err
	}

	funcDef.TypeParams = c.convFuncDefTypeParams(ctx, f)

	funcDef.Parameters, err = c.convFuncDefParams(ctx, f, append(funcDef.TypeParams, funcDef.Receiver.TypeParams...))
	if err != nil {
		return nil, err
	}

	funcDef.Results, err = c.convFuncDefResults(ctx, f, append(funcDef.TypeParams, funcDef.Receiver.TypeParams...))
	if err != nil {
		return nil, err
	}

	return funcDef, nil
}

func (c *Converter) findFuncDef(
	Package string,
	funcName string,
	pos token.Position,
) (*tool.FuncDef, *errors.PositionedErr) {
	if Package == "" && builtin.Type(funcName) {
		return nil, errors.BuiltinFunctionErr(pos, funcName)
	}

	pkgDir, err := c.pkgDir.Dir(Package)
	if err != nil {
		return nil, errors.Error(pos, err)
	}

	files := token.NewFileSet()
	pkgs, err := parser.ParseDir(files, pkgDir, nil, parser.ParseComments)
	if err != nil {
		return nil, errors.Error(pos, err)
	}

	var (
		funcDef *tool.FuncDef
		pErr    *errors.PositionedErr
	)
	for pkgName, pkg := range pkgs {
		if strings.HasSuffix(pkgName, "_test") {
			continue
		}
		for _, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if pErr != nil || funcDef != nil {
					return false
				}

				switch n := node.(type) {
				case *ast.FuncDecl:
					if n.Name.Name == funcName {
						funcDef, pErr = c.ConvertFuncDef(context.New(Package, file.Imports, files), n)
					}
					return false
				}

				return true
			})
		}
	}

	if pErr != nil {
		return nil, pErr
	}

	if funcDef == nil {
		return nil, errors.Errorf(pos, "function `%s` not found at `%s` package", funcName, Package)
	}

	return funcDef, nil
}

func (c *Converter) convFuncDefRecv(ctx context.Context, f *ast.FuncDecl) (tool.Receiver, *errors.PositionedErr) {
	var (
		recv tool.Receiver
		err  *errors.PositionedErr
	)

	if f.Recv != nil {
		recv.Presented = true
		recv.Position = ctx.Position(f.Recv.List[0].Pos())
		astTypeRef := syntax.ParseTypeRef(f.Recv.List[0].Type, ctx.Files)
		if err = astTypeRef.Error(); err != nil {
			return recv, err
		}

		for _, tp := range astTypeRef.TypeParams {
			recv.TypeParams = append(recv.TypeParams, tool.TypeParam{Name: tp.TypeName, Position: tp.Position})
		}
		recv.Type, err = c.TypeRef(ctx, astTypeRef, recv.TypeParams)
		if err != nil {
			return recv, err
		}
	}

	return recv, err
}

func (c *Converter) convFuncDefTypeParams(
	ctx context.Context,
	f *ast.FuncDecl,
) (res []tool.TypeParam) {
	if f.Type.TypeParams != nil {
		for _, tp := range f.Type.TypeParams.List {
			for _, name := range tp.Names {
				res = append(res, tool.TypeParam{
					Name:     name.Name,
					Position: ctx.Position(tp.Pos()),
				})
			}
		}
	}

	return res
}

func (c *Converter) convFuncDefParams(
	ctx context.Context,
	f *ast.FuncDecl,
	defTypeParams []tool.TypeParam,
) (res []tool.FuncParam, err *errors.PositionedErr) {
	if f.Type.Params != nil {
		for _, funcParam := range f.Type.Params.List {
			for _, name := range funcParam.Names {
				paramType, err := c.TypeRef(ctx, syntax.ParseTypeRef(funcParam.Type, ctx.Files), defTypeParams)
				if err != nil {
					return nil, err
				}
				res = append(res, tool.FuncParam{
					Name:     name.Name,
					Type:     paramType,
					Position: ctx.Position(funcParam.Pos())})
			}
		}
	}

	return res, err
}
func (c *Converter) convFuncDefResults(
	ctx context.Context,
	f *ast.FuncDecl,
	defTypeParams []tool.TypeParam,
) (res []tool.TypeRef, err *errors.PositionedErr) {
	if f.Type.Results != nil {
		for _, funcRes := range f.Type.Results.List {
			resType, err := c.TypeRef(ctx, syntax.ParseTypeRef(funcRes.Type, ctx.Files), defTypeParams)
			if err != nil {
				return nil, err
			}
			res = append(res, resType)
		}
	}

	return res, err
}

func (c *Converter) convFuncRef(ctx context.Context, ref syntax.Ref) (*tool.FuncRef, *errors.PositionedErr) {
	if ref.PkgAlias == "" && builtin.Func(ref.FuncName) {
		return nil, errors.BuiltinFunctionErr(ref.Position(), ref.FuncName)
	}

	pkg, pathErr := c.packagePath(ctx, ref.PkgAlias)
	if pathErr != nil {
		return nil, errors.Error(ref.Position(), pathErr)
	}

	fDef, err := c.findFuncDef(pkg, ref.FuncName, ref.Position())
	if err != nil {
		return nil, err
	}

	if len(ref.TypeParams) != len(fDef.TypeParams) {
		return nil, errors.InconsistentTypeParamsErr(ref.Position())
	}

	fRef := tool.FuncRefFromDef(ref.Position(), fDef)
	for i, tParam := range ref.TypeParams {
		param, pErr := fRef.NthTypeParam(i)
		if err != nil {
			return nil, errors.Error(param.Position, pErr)
		}

		effective, err := c.TypeRef(ctx, tParam, nil)
		if err != nil {
			return nil, err
		}
		fRef.SetEffectiveParam(param.Name, effective)
	}

	return fRef, nil
}
