package tests

import (
	"github.com/kr/pretty"
	"github.com/sabahtalateh/toolboxgen/internal/inspect"
	"go/ast"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sabahtalateh/toolboxgen/internal/convert"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
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

					i := inspect.New(inspect.Config{TrimPackage: trim, Introspective: true}).Type(t)
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
		want convertOut
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
			dir, err := os.Getwd()
			check(err)

			bb, err := os.ReadFile(filepath.Join(dir, tt.dir, "want.yaml"))
			check(err)

			var want map[string]any
			err = yaml.Unmarshal(bb, &want)
			check(err)

			got := convertFile(tt.dir, tt.dir)

			if !reflect.DeepEqual(want, got) {
				t.Errorf("convertFile() = %s\n\n\nwant = %s", pretty.Sprint(got), pretty.Sprint(want))
			}
		})
	}
}
