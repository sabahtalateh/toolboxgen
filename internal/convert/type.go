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

func (c *Converter) Type(ctx Context, typ *ast.TypeSpec) (types.Type, error) {
	if b, ok := c.builtin.Types[typ.Name.Name]; ok {
		return &types.Builtin{
			TypeName: b.TypeName,
			Declared: b.Declared,
		}, nil
	}

	switch t := typ.Type.(type) {
	case *ast.StructType:
		switch typ.Assign {
		case token.NoPos:
			return c.structFromSpec(ctx, typ, t)
		default:
			return c.typeAliasFromSpec(ctx, typ)
		}
	case *ast.InterfaceType:
		switch typ.Assign {
		case token.NoPos:
			return c.interfaceFromSpec(ctx, typ, t)
		default:
			return c.typeAliasFromSpec(ctx, typ)
		}
	default:
		switch typ.Assign {
		case token.NoPos:
			return c.typeDefFromSpec(ctx, typ)
		default:
			return c.typeAliasFromSpec(ctx, typ)
		}
	}
}

func (c *Converter) findType(ctx Context, Package, Type string) (types.Type, error) {
	if typ, ok := c.types[typeKey(Package, Type)]; ok {
		return typ, nil
	}

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

	return c.Type(ctx.WithPackage(Package).WithImports(imports).WithFiles(files).WithPos(spec.Pos()), spec)
}

func (c *Converter) structFromSpec(ctx Context, spec *ast.TypeSpec, typ *ast.StructType) (*types.Struct, error) {
	var (
		res *types.Struct
		err error
	)

	res = &types.Struct{
		Package:      ctx.Package(),
		TypeName:     spec.Name.Name,
		TypeParams:   TypeParams(ctx, spec.TypeParams),
		Position:     ctx.NodePosition(spec),
		TypePosition: ctx.NodePosition(spec.Type),
		Declared:     code.OfNode(spec),
	}

	c.putType(res)

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
		Package:      ctx.Package(),
		TypeName:     spec.Name.Name,
		TypeParams:   TypeParams(ctx, spec.TypeParams),
		Position:     ctx.NodePosition(spec),
		TypePosition: ctx.NodePosition(spec.Type),
		Declared:     code.OfNode(spec),
	}

	c.putType(res)

	if res.Methods, err = c.Fields(ctx.WithDefined(res.TypeParams), typ.Methods); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) typeDefFromSpec(ctx Context, spec *ast.TypeSpec) (*types.TypeDef, error) {
	var (
		res *types.TypeDef
		err error
	)

	res = &types.TypeDef{
		Package:      ctx.Package(),
		TypeName:     spec.Name.Name,
		TypeParams:   TypeParams(ctx, spec.TypeParams),
		Position:     ctx.NodePosition(spec),
		TypePosition: ctx.NodePosition(spec.Type),
		Declared:     code.OfNode(spec),
	}

	c.putType(res)

	if res.Type, err = c.TypeRef(ctx.WithDefined(res.TypeParams), spec.Type); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Converter) typeAliasFromSpec(ctx Context, spec *ast.TypeSpec) (*types.TypeAlias, error) {
	var (
		res *types.TypeAlias
		err error
	)

	res = &types.TypeAlias{
		Package:      ctx.Package(),
		TypeName:     spec.Name.Name,
		Position:     ctx.NodePosition(spec),
		TypePosition: ctx.NodePosition(spec.Type),
		Declared:     code.OfNode(spec),
	}

	c.putType(res)

	if res.Type, err = c.TypeRef(ctx, spec.Type); err != nil {
		return nil, err
	}

	return res, nil
}
