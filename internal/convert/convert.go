package convert

import (
	"fmt"
	"github.com/sabahtalateh/toolboxgen/internal/types"
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/life4/genesis/slices"

	"github.com/sabahtalateh/toolboxgen/internal/builtin"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"github.com/sabahtalateh/toolboxgen/internal/pkgdir"
	"github.com/sabahtalateh/toolboxgen/internal/utils"
)

type Converter struct {
	pkgDir  *pkgdir.PkgDir
	builtin *builtin.Builtin
	types   map[string]types.Type
}

func New(mod *mod.Module) (*Converter, error) {
	var (
		c   = &Converter{types: map[string]types.Type{}}
		err error
	)

	if c.pkgDir, err = pkgdir.Init(mod); err != nil {
		return nil, err
	}

	c.builtin = builtin.Init()
	return c, nil
}

func (c *Converter) putType(typ types.Type) types.Type {
	c.types[typeKeyFromType(typ)] = typ
	return typ
}

func typeKeyFromType(typ types.Type) string {
	var Package, typeName string
	switch t := typ.(type) {
	case *types.Builtin:
		typeName = t.TypeName
	case *types.Struct:
		Package, typeName = t.Package, t.TypeName
	case *types.Interface:
		Package, typeName = t.Package, t.TypeName
	case *types.TypeDef:
		Package, typeName = t.Package, t.TypeName
	case *types.TypeAlias:
		Package, typeName = t.Package, t.TypeName
	default:
		panic("unknown type")
	}

	return fmt.Sprintf("%s.%s", Package, typeName)
}

func typeKey(Package string, typeName string) string {
	key := typeName
	if Package != "" {
		key = fmt.Sprintf("%s.%s", Package, key)
	}
	return key
}

func packagePath(ctx Context, alias string) (Package string, err error) {
	if alias == "" {
		Package = ctx.Package()
	} else {
		var imp *ast.ImportSpec
		for _, spec := range ctx.Imports() {
			if spec.Name != nil {
				if spec.Name.Name == alias {
					imp = spec
					break
				}
			} else {
				parts := strings.Split(utils.Unquote(spec.Path.Value), "/")
				if len(parts) == 0 {
					err = fmt.Errorf("empty import")
					break
				}
				lastPart := parts[len(parts)-1]
				if alias == lastPart {
					imp = spec
					break
				}
			}
		}

		if imp == nil {
			err = fmt.Errorf("package `%s` not resolved", alias)
		} else {
			Package = imp.Path.Value
		}
	}

	return utils.Unquote(Package), err
}

func (c *Converter) dir(mod *mod.Module, pakage string) (string, error) {
	if strings.HasPrefix(pakage, mod.Path) {
		suffix := strings.TrimPrefix(pakage, mod.Path)
		suffParts := strings.Split(suffix, "/")
		suffParts = slices.Filter(suffParts, func(el string) bool {
			return el != ""
		})
		d := filepath.Join(append([]string{mod.Dir}, suffParts...)...)
		return d, nil
	}
	return "", fmt.Errorf("failed to determine directory for `%s` package", pakage)
}
