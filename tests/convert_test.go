package tests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kr/pretty"

	"github.com/sabahtalateh/toolboxgen/internal/convert"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"github.com/sabahtalateh/toolboxgen/internal/types"
)

func convertFile(file string) convertOut {
	dir, err := os.Getwd()
	check(err)

	m, err := mod.LookupDir(filepath.Dir(filepath.Join(dir, file)), true)
	check(err)

	conv, err := convert.New(m)
	check(err)

	files := token.NewFileSet()
	pkgs, err := parser.ParseDir(files, filepath.Join(dir, file), nil, parser.ParseComments)
	check(err)

	res := convertOut{types: map[string]types.Type{}, funcs: map[string]*types.Function{}}
	for _, pkg := range pkgs {
		for filePath, f := range pkg.Files {
			Package := filepath.Dir(filePath)
			if strings.HasPrefix(Package, m.Dir) {
				Package = strings.Replace(Package, m.Dir, m.Path, 1)
			}

			ast.Inspect(f, func(node ast.Node) bool {
				ctx := convert.NewContext(Package, f.Imports, files, token.Position{}, nil)
				switch n := node.(type) {
				case *ast.TypeSpec:
					t, err := conv.Type(ctx.WithPos(n.Pos()), n)
					check(err)
					res.types[typeID(t)] = t
					return false
					// case *ast.FuncDecl:
					// 	f, err := conv.Function(ctx.WithPos(n.Pos()), n)
					// 	check(err)
					// 	res.funcs[fmt.Sprintf("%s.%s", f.Package, f.FuncName)] = f
					// 	println(f)
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
		{
			name: "struct-1",
			dir:  "testmod/struct_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/struct_1.A": &types.Struct{
						Package:  "testmod/struct_1",
						TypeName: "A",
					},
				},
			},
		},
		{
			name: "interface-1",
			dir:  "testmod/interface_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/struct_1.A": &types.Interface{
						Package:  "testmod/interface_1",
						TypeName: "A",
						Methods: []*types.Field{
							{
								Name: "Method",
								Type: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "xyz",
			dir:  "testmod/xyz",
			want: convertOut{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertFile(tt.dir)
			if !got.equal(tt.want) {
				t.Errorf("convertFile() = %s\n\n\nwant = %s", pretty.Sprint(got), pretty.Sprint(tt.want))
			}
		})
	}
}
