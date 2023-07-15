package convert

import (
	"github.com/sabahtalateh/toolboxgen/internal/maps"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/sabahtalateh/toolboxgen/internal/code"
	"github.com/sabahtalateh/toolboxgen/internal/errors"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func (c *Converter) Type(ctx Context, t *ast.TypeSpec) (types.Type, error) {
	switch t.Type.(type) {
	case *ast.StructType:
		return structFromSpec(ctx, t), nil
	case *ast.InterfaceType:
		return interfaceFromSpec(ctx, t), nil
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
		return nil, errors.Error(ctx.CurrentPosition(), err)
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

	if spec == nil && Package == ctx.Package {
		if t, ok := c.builtin.Types[Type]; ok {
			return &types.Builtin{Declared: t.Declared, TypeName: Type}, nil
		}
	}

	if spec == nil {
		return nil, errors.Errorf(ctx.CurrentPosition(), "type %s not found at %s", Type, Package)
	}

	ctx = ctx.WithPackage(Package).WithImports(imports).WithFiles(files).WithPos(spec.Pos())
	return c.Type(ctx, spec)
}

func structFromSpec(ctx Context, t *ast.TypeSpec) *types.Struct {
	typ := &types.Struct{
		Declared:   code.OfNode(t),
		Package:    ctx.Package,
		TypeName:   t.Name.Name,
		TypeParams: TypeParams(ctx, t.TypeParams),
		Position:   ctx.Position(t.Pos()),
	}
	typ.TypeParams = TypeParams(ctx, t.TypeParams)

	return typ
}

func interfaceFromSpec(ctx Context, t *ast.TypeSpec) *types.Interface {
	typ := &types.Interface{
		Declared:   code.OfNode(t),
		Package:    ctx.Package,
		TypeName:   t.Name.Name,
		TypeParams: TypeParams(ctx, t.TypeParams),
		Position:   ctx.Position(t.Pos()),
	}

	return typ
}

func (c *Converter) typeDefFromSpec(ctx Context, t *ast.TypeSpec) (*types.TypeDef, error) {
	var (
		typ *types.TypeDef
		err error
	)

	typ = &types.TypeDef{
		Declared:   code.OfNode(t),
		Package:    ctx.Package,
		TypeName:   t.Name.Name,
		TypeParams: TypeParams(ctx, t.TypeParams),
		Position:   ctx.Position(t.Pos()),
	}

	defined := maps.FromSlice(typ.TypeParams, func(v *types.TypeParam) string { return v.Original })
	typ.Type, err = c.TypeRef(ctx.WithDefined(defined), t.Type)
	if err != nil {
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
		Declared: code.OfNode(t),
		Package:  ctx.Package,
		TypeName: t.Name.Name,
		Position: ctx.Position(t.Pos()),
	}

	typ.Type, err = c.TypeRef(ctx, t.Type)
	if err != nil {
		return nil, err
	}

	return typ, nil
}
