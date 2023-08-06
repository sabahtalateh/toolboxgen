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

func (c *Converter) Type(ctx Context, s *ast.TypeSpec) (types.Type, error) {
	if b, ok := c.builtin.Types[s.Name.Name]; ok {
		return &types.Builtin{
			TypeName: b.TypeName,
			Code:     b.Code,
		}, nil
	}

	switch t := s.Type.(type) {
	case *ast.StructType:
		switch s.Assign {
		case token.NoPos:
			return c.convertStruct(ctx, s, t)
		default:
			return c.convertTypeAlias(ctx, s)
		}
	case *ast.InterfaceType:
		switch s.Assign {
		case token.NoPos:
			return c.convertInterface(ctx, s, t)
		default:
			return c.convertTypeAlias(ctx, s)
		}
	default:
		switch s.Assign {
		case token.NoPos:
			return c.convertTypeDef(ctx, s)
		default:
			return c.convertTypeAlias(ctx, s)
		}
	}
}

func (c *Converter) convertStruct(ctx Context, s *ast.TypeSpec, t *ast.StructType) (*types.Struct, error) {
	var (
		typ *types.Struct
		err error
	)

	typ = &types.Struct{
		Package:      ctx.Package(),
		TypeName:     s.Name.Name,
		TypeParams:   TypeParams(ctx, s.TypeParams),
		Position:     ctx.NodePosition(s),
		TypePosition: ctx.NodePosition(s.Type),
		Code:         code.OfNode(s),
	}

	c.putType(typ)

	if typ.Fields, err = c.Fields(ctx.WithDefined(typ.TypeParams), t.Fields); err != nil {
		return nil, err
	}

	return typ, nil
}

func (c *Converter) convertInterface(ctx Context, s *ast.TypeSpec, t *ast.InterfaceType) (*types.Interface, error) {
	var (
		typ *types.Interface
		err error
	)

	typ = &types.Interface{
		Package:      ctx.Package(),
		TypeName:     s.Name.Name,
		TypeParams:   TypeParams(ctx, s.TypeParams),
		Position:     ctx.NodePosition(s),
		TypePosition: ctx.NodePosition(s.Type),
		Code:         code.OfNode(s),
	}

	c.putType(typ)

	if typ.Fields, err = c.Fields(ctx.WithDefined(typ.TypeParams), t.Methods); err != nil {
		return nil, err
	}

	return typ, nil
}

func (c *Converter) convertTypeDef(ctx Context, s *ast.TypeSpec) (*types.TypeDef, error) {
	var (
		typ *types.TypeDef
		err error
	)

	typ = &types.TypeDef{
		Package:      ctx.Package(),
		TypeName:     s.Name.Name,
		TypeParams:   TypeParams(ctx, s.TypeParams),
		Position:     ctx.NodePosition(s),
		TypePosition: ctx.NodePosition(s.Type),
		Code:         code.OfNode(s),
	}

	c.putType(typ)

	if typ.Type, err = c.TypeExpr(ctx.WithDefined(typ.TypeParams), s.Type); err != nil {
		return nil, err
	}

	return typ, nil
}

func (c *Converter) convertTypeAlias(ctx Context, s *ast.TypeSpec) (*types.TypeAlias, error) {
	var (
		typ *types.TypeAlias
		err error
	)

	typ = &types.TypeAlias{
		Package:      ctx.Package(),
		TypeName:     s.Name.Name,
		Position:     ctx.NodePosition(s),
		TypePosition: ctx.NodePosition(s.Type),
		Code:         code.OfNode(s),
	}

	c.putType(typ)

	if typ.Type, err = c.TypeExpr(ctx, s.Type); err != nil {
		return nil, err
	}

	return typ, nil
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
