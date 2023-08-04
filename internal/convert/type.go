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
			Code:     b.Code,
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
			return &types.Builtin{TypeName: Type, Code: t.Code}, nil
		}
	}

	if spec == nil {
		return nil, errors.Errorf(ctx.Position(), "type %s not found at %s", Type, Package)
	}

	return c.Type(ctx.WithPackage(Package).WithImports(imports).WithFiles(files).WithPos(spec.Pos()), spec)
}

func (c *Converter) structFromSpec(ctx Context, spec *ast.TypeSpec, ast *ast.StructType) (*types.Struct, error) {
	var (
		typ *types.Struct
		err error
	)

	typ = &types.Struct{
		Package:      ctx.Package(),
		TypeName:     spec.Name.Name,
		TypeParams:   TypeParams(ctx, spec.TypeParams),
		Position:     ctx.NodePosition(spec),
		TypePosition: ctx.NodePosition(spec.Type),
		Code:         code.OfNode(spec),
	}

	c.putType(typ)

	if typ.Fields, err = c.Fields(ctx.WithDefined(typ.TypeParams), ast.Fields); err != nil {
		return nil, err
	}

	return typ, nil
}

func (c *Converter) interfaceFromSpec(ctx Context, spec *ast.TypeSpec, ast *ast.InterfaceType) (*types.Interface, error) {
	var (
		typ *types.Interface
		err error
	)

	typ = &types.Interface{
		Package:      ctx.Package(),
		TypeName:     spec.Name.Name,
		TypeParams:   TypeParams(ctx, spec.TypeParams),
		Position:     ctx.NodePosition(spec),
		TypePosition: ctx.NodePosition(spec.Type),
		Code:         code.OfNode(spec),
	}

	c.putType(typ)

	if typ.Fields, err = c.Fields(ctx.WithDefined(typ.TypeParams), ast.Methods); err != nil {
		return nil, err
	}

	return typ, nil
}

func (c *Converter) typeDefFromSpec(ctx Context, spec *ast.TypeSpec) (*types.TypeDef, error) {
	var (
		typ *types.TypeDef
		err error
	)

	typ = &types.TypeDef{
		Package:      ctx.Package(),
		TypeName:     spec.Name.Name,
		TypeParams:   TypeParams(ctx, spec.TypeParams),
		Position:     ctx.NodePosition(spec),
		TypePosition: ctx.NodePosition(spec.Type),
		Code:         code.OfNode(spec),
	}

	c.putType(typ)

	if typ.Type, err = c.TypeRef(ctx.WithDefined(typ.TypeParams), spec.Type); err != nil {
		return nil, err
	}

	return typ, nil
}

func (c *Converter) typeAliasFromSpec(ctx Context, spec *ast.TypeSpec) (*types.TypeAlias, error) {
	var (
		typ *types.TypeAlias
		err error
	)

	typ = &types.TypeAlias{
		Package:      ctx.Package(),
		TypeName:     spec.Name.Name,
		Position:     ctx.NodePosition(spec),
		TypePosition: ctx.NodePosition(spec.Type),
		Code:         code.OfNode(spec),
	}

	c.putType(typ)

	if typ.Type, err = c.TypeRef(ctx, spec.Type); err != nil {
		return nil, err
	}

	return typ, nil
}
