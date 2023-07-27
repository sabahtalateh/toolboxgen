package tests

import (
	"go/ast"
	"go/parser"
	"go/token"
	yaml2 "gopkg.in/yaml.v2"
	yaml3 "gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sabahtalateh/toolboxgen/internal/convert"
	"github.com/sabahtalateh/toolboxgen/internal/inspect"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"golang.org/x/exp/maps"
)

func convertFile(file, trim string) map[string]any {
	dir, err := os.Getwd()
	check(err)

	m, err := mod.LookupDir(filepath.Dir(filepath.Join(dir, file)), true)
	check(err)

	conv, err := convert.New(m)
	check(err)

	files := token.NewFileSet()
	pkgs, err := parser.ParseDir(files, filepath.Join(dir, file), nil, parser.ParseComments)
	check(err)

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
					t, err := conv.Type(ctx.WithPos(n.Pos()), n)
					check(err)

					i := inspect.New(inspect.Config{TrimPackage: trim}).Type(t)
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
	tests := []testCase{
		{name: "struct", dir: "testmod/convert/struct"},
		{name: "interface", dir: "testmod/convert/interface"},
		{name: "interface-2", dir: "testmod/convert/interface2"},
		{name: "typedef", dir: "testmod/convert/typedef"},
		{name: "typealias", dir: "testmod/convert/typealias"},
		{name: "builtin", dir: "testmod/convert/builtin"},
		{name: "map", dir: "testmod/convert/map"},
		{name: "chan", dir: "testmod/convert/chan"},
		{name: "functype", dir: "testmod/convert/functype"},
		{name: "structtype", dir: "testmod/convert/structtype"},
		{name: "interfacetype", dir: "testmod/convert/interfacetype"},
		{name: "structref", dir: "testmod/convert/structref"},
		{name: "interfaceref", dir: "testmod/convert/interfaceref"},
		{name: "typedefref", dir: "testmod/convert/typedefref"},
		{name: "typealiasref", dir: "testmod/convert/typealiasref"},
		{name: "typeparamref", dir: "testmod/convert/typeparamref"},
		{name: "complex", dir: "testmod/convert/complex"},
		{name: "typeparams", dir: "testmod/convert/typeparams"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := check2(os.Getwd())
			bb := check2(os.ReadFile(filepath.Join(dir, tt.dir, "want.yaml")))

			var want map[string]any
			check(yaml3.Unmarshal(bb, &want))

			got := convertFile(tt.dir, tt.dir)

			if !reflect.DeepEqual(got, want) {
				g := check2(yaml2.Marshal(ordered(got)))
				w := check2(yaml2.Marshal(ordered(want)))

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
	case "struct", "interface", "typedef", "typealias":
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
