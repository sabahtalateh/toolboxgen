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
	dir := tests.Unwrap(os.Getwd())
	m := tests.Unwrap(mod.LookupDir(filepath.Dir(filepath.Join(dir, file)), true))
	conv := tests.Unwrap(convert.New(m))

	files := token.NewFileSet()
	pkgs := tests.Unwrap(parser.ParseDir(files, filepath.Join(dir, file), nil, parser.ParseComments))

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
					t := tests.Unwrap(conv.Type(ctx.WithPos(n.Pos()), n))
					i := inspect.New(inspect.Config{TrimPackage: trim}).Type(t)
					for k, v := range i {
						res[k] = v
					}
				case *ast.FuncDecl:
					f := tests.Unwrap(conv.Function(ctx.WithPos(n.Pos()), n))
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
	type testCase struct {
		name string
		dir  string
	}
	ttt := []testCase{
		{name: "struct", dir: "mod/convert/struct"},
		{name: "interface", dir: "mod/convert/interface"},
		{name: "interface-2", dir: "mod/convert/interface2"},
		{name: "typedef", dir: "mod/convert/typedef"},
		{name: "typealias", dir: "mod/convert/typealias"},
		{name: "builtin", dir: "mod/convert/builtin"},
		{name: "map", dir: "mod/convert/map"},
		{name: "chan", dir: "mod/convert/chan"},
		{name: "functype", dir: "mod/convert/functype"},
		{name: "structtype", dir: "mod/convert/structtype"},
		{name: "interfacetype", dir: "mod/convert/interfacetype"},
		{name: "structref", dir: "mod/convert/structref"},
		{name: "interfaceref", dir: "mod/convert/interfaceref"},
		{name: "typedefref", dir: "mod/convert/typedefref"},
		{name: "typealiasref", dir: "mod/convert/typealiasref"},
		{name: "typeparamref", dir: "mod/convert/typeparamref"},
		{name: "complex", dir: "mod/convert/complex"},
		{name: "typeparams", dir: "mod/convert/typeparams"},
		{name: "func", dir: "mod/convert/func"},
	}
	for _, tt := range ttt {
		t.Run(tt.name, func(t *testing.T) {
			dir := tests.Unwrap(os.Getwd())
			bb := tests.Unwrap(os.ReadFile(filepath.Join(dir, tt.dir, "want.yaml")))

			var want map[string]any
			tests.Check(yaml3.Unmarshal(bb, &want))

			got := convertFile(tt.dir, tt.dir)

			if !reflect.DeepEqual(got, want) {
				g := tests.Unwrap(yaml2.Marshal(ordered(got)))
				w := tests.Unwrap(yaml2.Marshal(ordered(want)))

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
