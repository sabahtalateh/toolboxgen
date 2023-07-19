package convert

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (c *Converter) Type(ctx Context, t *ast.TypeSpec) (types.Type, error) {
	if builtin, ok := c.builtin.Types[t.Name.Name]; ok {
		return &types.Builtin{
			TypeName: builtin.TypeName,
			Declared: builtin.Declared,
		}, nil
	}

	switch typ := t.Type.(type) {
	case *ast.StructType:
		switch t.Assign {
		case token.NoPos:
			return c.structFromSpec(ctx, t, typ)
		default:
			return c.typeAliasFromSpec(ctx, t)
		}
	case *ast.InterfaceType:
		switch t.Assign {
		case token.NoPos:
			return c.interfaceFromSpec(ctx, t, typ)
		default:
			return c.typeAliasFromSpec(ctx, t)
		}
	default:
		switch t.Assign {
		case token.NoPos:
			return c.typeDefFromSpec(ctx, t)
		default:
			return c.typeAliasFromSpec(ctx, t)
		}
	}
}

func (c *Converter) findType(ctx Context, Package, Type string) (types.Type, error) {
	pkgDir, err := c.pkgDir.Dir(Package)
	if err != nil {
		return nil, errors.Error(ctx.Position(), err)
	}

	files := token.NewFileSet()
	pkgs, err := parser.ParseDir(files, pkgDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var (
		spec    *ast.TypeSpec
		imports []*ast.ImportSpec
	)

	for pkgName, pkg := range pkgs {
		if strings.HasSuffix(pkgName, "_test") {
			continue
		}
		for _, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				if err != nil || spec != nil {
					return false
				}

				switch n := node.(type) {
				case *ast.TypeSpec:
					if n.Name.Name == Type {
						spec = n
						imports = file.Imports
					}
				}
				return true
			})
		}
	}

	if err != nil {
		return nil, err
	}

	if spec == nil && Package == ctx.Package() {
		if t, ok := c.builtin.Types[Type]; ok {
			return &types.Builtin{TypeName: Type, Declared: t.Declared}, nil
		}
	}

	if spec == nil {
		return nil, errors.Errorf(ctx.Position(), "type %s not found at %s", Type, Package)
	}

	ctx = ctx.WithPackage(Package).WithImports(imports).WithFiles(files).WithPos(spec.Pos())
	return c.Type(ctx, spec)
}

func (c *Converter) structFromSpec(ctx Context, spec *ast.TypeSpec, typ *ast.StructType) (*types.Struct, error) {
	var (
		res *types.Struct
		err error
	)

	res = &types.Struct{
		Package:    ctx.Package(),
		TypeName:   spec.Name.Name,
		TypeParams: TypeParams(ctx, spec.TypeParams),
		Position:   ctx.NodePosition(spec),
		Declared:   code.OfNode(spec),
	}

	if res.Fields, err = c.Fields(ctx.WithDefined(res.TypeParams), typ.Fields); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) interfaceFromSpec(ctx Context, spec *ast.TypeSpec, typ *ast.InterfaceType) (*types.Interface, error) {
	var (
		res *types.Interface
		err error
	)

	res = &types.Interface{
		Package:    ctx.Package(),
		TypeName:   spec.Name.Name,
		TypeParams: TypeParams(ctx, spec.TypeParams),
		Position:   ctx.NodePosition(spec),
		Declared:   code.OfNode(spec),
	}

	if res.Methods, err = c.Fields(ctx.WithDefined(res.TypeParams), typ.Methods); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) typeDefFromSpec(ctx Context, t *ast.TypeSpec) (*types.TypeDef, error) {
	var (
		typ *types.TypeDef
		err error
	)

	typ = &types.TypeDef{
		Package:    ctx.Package(),
		TypeName:   t.Name.Name,
		TypeParams: TypeParams(ctx, t.TypeParams),
		Position:   ctx.NodePosition(t),
		Declared:   code.OfNode(t),
	}

	if typ.Type, err = c.TypeRef(ctx.WithDefined(typ.TypeParams), t.Type); err != nil {
		return nil, err
	}

	return typ, err
}

func (c *Converter) typeAliasFromSpec(ctx Context, t *ast.TypeSpec) (*types.TypeAlias, error) {
	var (
		typ *types.TypeAlias
		err error
	)

	typ = &types.TypeAlias{
		Package:  ctx.Package(),
		TypeName: t.Name.Name,
		Position: ctx.NodePosition(t),
		Declared: code.OfNode(t),
	}

	if typ.Type, err = c.TypeRef(ctx, t.Type); err != nil {
		return nil, err
	}

	return typ, nil
}
