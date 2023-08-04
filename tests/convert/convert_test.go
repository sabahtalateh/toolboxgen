package convert

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"golang.org/x/exp/maps"
	yaml2 "gopkg.in/yaml.v2"
	yaml3 "gopkg.in/yaml.v3"

	"github.com/sabahtalateh/toolboxgen/internal/convert"
	"github.com/sabahtalateh/toolboxgen/internal/inspect"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"github.com/sabahtalateh/toolboxgen/tests"
)

func convertFile(file, trim string) map[string]any {
	dir := tutils.Unwrap(os.Getwd())
	m := tutils.Unwrap(mod.LookupDir(filepath.Join(dir, "mod"), true))
	conv := tutils.Unwrap(convert.New(m))

	files := token.NewFileSet()
	pkgs := tutils.Unwrap(parser.ParseDir(files, filepath.Join(dir, "mod", file), nil, parser.ParseComments))

	res := map[string]any{}
	for _, pkg := range pkgs {
		for filePath, f := range pkg.Files {
			Package := filepath.Dir(filePath)
			if strings.HasPrefix(Package, m.Dir) {
				Package = strings.Replace(Package, m.Dir, m.Path, 1)
			}

			ast.Inspect(f, func(node ast.Node) bool {
				ctx := convert.NewContext().WithPackage(Package).WithImports(f.Imports).WithFiles(files)
				switch n := node.(type) {
				case *ast.TypeSpec:
					t := tutils.Unwrap(conv.Type(ctx.WithPos(n.Pos()), n))
					i := inspect.New(inspect.Config{TrimPackage: trim}).Type(t)
					for k, v := range i {
						res[k] = v
					}
				case *ast.FuncDecl:
					f := tutils.Unwrap(conv.Function(ctx.WithPos(n.Pos()), n))
					i := inspect.New(inspect.Config{TrimPackage: trim}).Function(f)
					for k, v := range i {
						res[k] = v
					}
				}
				return true
			})
		}
	}

	return res
}

func TestConvert(t *testing.T) {
	tests := []string{
		"struct",
		"interface",
		"interface_2",
		"typedef",
		"typealias",
		"builtin",
		"map",
		"chan",
		"functype",
		"structtype",
		"interfacetype",
		"structref",
		"interfaceref",
		"typedefref",
		"typealiasref",
		"typeparamref",
		"complex",
		"typeparams",
		"func",
	}
	for _, tName := range tests {
		t.Run(tName, func(t *testing.T) {
			dir := tutils.Unwrap(os.Getwd())
			bb := tutils.Unwrap(os.ReadFile(filepath.Join(dir, "mod", tName, "want.yaml")))

			var want map[string]any
			tutils.Check(yaml3.Unmarshal(bb, &want))

			got := convertFile(tName, filepath.Join("mod", tName))

			if !reflect.DeepEqual(got, want) {
				g := tutils.Unwrap(yaml2.Marshal(ordered(got)))
				w := tutils.Unwrap(yaml2.Marshal(ordered(want)))

				t.Errorf("\ngot:\n\n%s\nwant\n\n%s", g, w)
			}
		})
	}
}

func ordered(m map[string]any) yaml2.MapSlice {
	keys := maps.Keys(m)
	sort.SliceStable(keys, func(i, j int) bool {
		o1 := ord(keys[i])
		o2 := ord(keys[j])
		if o1 == o2 {
			return keys[i] < keys[j]
		}

		return o1 < o2
	})

	y := yaml2.MapSlice{}
	for _, k := range keys {
		v := m[k]
		switch vv := v.(type) {
		case map[string]any:
			y = append(y, yaml2.MapItem{Key: k, Value: ordered(vv)})
		case []any:
			var res []any
			for _, v := range vv {
				switch vv := v.(type) {
				case map[string]any:
					uu := ordered(vv)
					res = append(res, uu)
				default:
					res = append(res, vv)
				}
			}
			y = append(y, yaml2.MapItem{Key: k, Value: res})
		default:
			y = append(y, yaml2.MapItem{Key: k, Value: vv})
		}
	}

	return y
}

func ord(x string) int {
	switch x {
	case "struct", "interface", "typedef", "typealias", "receiver":
		return 3
	case "name":
		return 0
	case "modifiers":
		return 1
	case "actual":
		return 4
	case "intro":
		return 99
	case "map":
		return 100
	case "key":
		return 101
	case "value":
		return 102
	default:
		return 999
	}
}
