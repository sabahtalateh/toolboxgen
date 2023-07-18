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

	res := convertOut{types: map[string]types.Type{}}
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
			dir:  "testmod/convert/struct_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/struct_1.A": &types.Struct{
						Package:  "testmod/convert/struct_1",
						TypeName: "A",
					},
				},
			},
		},
		{
			name: "interface-1",
			dir:  "testmod/convert/interface_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/interface_1.A": &types.Interface{
						Package:  "testmod/convert/interface_1",
						TypeName: "A",
						Methods: []*types.Field{{
							Name: "Method",
							Type: &types.FuncTypeRef{
								Params:  []*types.Field{{Type: &types.BuiltinRef{TypeName: "string"}}},
								Results: []*types.Field{{Type: &types.BuiltinRef{TypeName: "string"}}},
							}},
						},
					},
				},
			},
		},
		{
			name: "interface-2",
			dir:  "testmod/convert/interface_2",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/interface_2.A": &types.Interface{
						Package:  "testmod/convert/interface_2",
						TypeName: "A",
						Methods: []*types.Field{
							{
								Name: "Method1",
								Type: &types.FuncTypeRef{
									Params: []*types.Field{{Name: "a", Type: &types.BuiltinRef{TypeName: "string"}}},
									Results: []*types.Field{
										{Name: "b", Type: &types.BuiltinRef{TypeName: "string"}},
										{Name: "c", Type: &types.BuiltinRef{TypeName: "error"}},
									},
								},
							},
							{
								Name: "Method2",
								Type: &types.FuncTypeRef{
									Params: []*types.Field{{Name: "d", Type: &types.BuiltinRef{TypeName: "string"}}},
									Results: []*types.Field{
										{Name: "e", Type: &types.BuiltinRef{TypeName: "string"}},
										{Name: "f", Type: &types.BuiltinRef{TypeName: "error"}},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "typedef-1",
			dir:  "testmod/convert/typedef_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typedef_1.A": &types.TypeDef{
						Package:  "testmod/convert/typedef_1",
						TypeName: "A",
						Type:     &types.BuiltinRef{TypeName: "int"},
					},
				},
			},
		},
		{
			name: "typealias-1",
			dir:  "testmod/convert/typealias_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typealias_1.A": &types.TypeAlias{
						Package:  "testmod/convert/typealias_1",
						TypeName: "A",
						Type:     &types.BuiltinRef{TypeName: "int"},
					},
				},
			},
		},
		{
			name: "builtin-1",
			dir:  "testmod/convert/builtin_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/builtin_1.Bool": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Bool",
						Type:     &types.BuiltinRef{TypeName: "bool"},
					},
					"testmod/convert/builtin_1.Uint8": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Uint8",
						Type:     &types.BuiltinRef{TypeName: "uint8"},
					},
					"testmod/convert/builtin_1.Uint16": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Uint16",
						Type:     &types.BuiltinRef{TypeName: "uint16"},
					},
					"testmod/convert/builtin_1.Uint32": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Uint32",
						Type:     &types.BuiltinRef{TypeName: "uint32"},
					},
					"testmod/convert/builtin_1.Uint64": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Uint64",
						Type:     &types.BuiltinRef{TypeName: "uint64"},
					},
					"testmod/convert/builtin_1.Int8": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Int8",
						Type:     &types.BuiltinRef{TypeName: "int8"},
					},
					"testmod/convert/builtin_1.Int16": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Int16",
						Type:     &types.BuiltinRef{TypeName: "int16"},
					},
					"testmod/convert/builtin_1.Int32": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Int32",
						Type:     &types.BuiltinRef{TypeName: "int32"},
					},
					"testmod/convert/builtin_1.Int64": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Int64",
						Type:     &types.BuiltinRef{TypeName: "int64"},
					},
					"testmod/convert/builtin_1.Float32": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Float32",
						Type:     &types.BuiltinRef{TypeName: "float32"},
					},
					"testmod/convert/builtin_1.Float64": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Float64",
						Type:     &types.BuiltinRef{TypeName: "float64"},
					},
					"testmod/convert/builtin_1.Complex64": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Complex64",
						Type:     &types.BuiltinRef{TypeName: "complex64"},
					},
					"testmod/convert/builtin_1.Complex128": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Complex128",
						Type:     &types.BuiltinRef{TypeName: "complex128"},
					},
					"testmod/convert/builtin_1.String": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "String",
						Type:     &types.BuiltinRef{TypeName: "string"},
					},
					"testmod/convert/builtin_1.Int": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Int",
						Type:     &types.BuiltinRef{TypeName: "int"},
					},
					"testmod/convert/builtin_1.Uint": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Uint",
						Type:     &types.BuiltinRef{TypeName: "uint"},
					},
					"testmod/convert/builtin_1.Uintptr": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Uintptr",
						Type:     &types.BuiltinRef{TypeName: "uintptr"},
					},
					"testmod/convert/builtin_1.Byte": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Byte",
						Type:     &types.BuiltinRef{TypeName: "byte"},
					},
					"testmod/convert/builtin_1.Rune": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Rune",
						Type:     &types.BuiltinRef{TypeName: "rune"},
					},
					"testmod/convert/builtin_1.Any": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Any",
						Type:     &types.BuiltinRef{TypeName: "any"},
					},
					"testmod/convert/builtin_1.Comparable": &types.TypeDef{
						Package:  "testmod/convert/builtin_1",
						TypeName: "Comparable",
						Type:     &types.BuiltinRef{TypeName: "comparable"},
					},
				},
			},
		},
		{
			name: "map-1",
			dir:  "testmod/convert/map_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/map_1.A": &types.TypeDef{
						Package:  "testmod/convert/map_1",
						TypeName: "A",
						Type: &types.MapRef{
							Key:   &types.BuiltinRef{TypeName: "string"},
							Value: &types.BuiltinRef{TypeName: "string"},
						},
					},
				},
			},
		},
		{
			name: "chan-1",
			dir:  "testmod/convert/chan_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/chan_1.A": &types.TypeDef{
						Package:  "testmod/convert/chan_1",
						TypeName: "A",
						Type:     &types.ChanRef{Value: &types.BuiltinRef{TypeName: "string"}},
					},
				},
			},
		},
		{
			name: "functype-1",
			dir:  "testmod/convert/functype_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/functype_1.A": &types.TypeDef{
						Package:  "testmod/convert/functype_1",
						TypeName: "A",
						Type: &types.FuncTypeRef{
							Params:  []*types.Field{{Name: "a", Type: &types.BuiltinRef{TypeName: "int8"}}},
							Results: []*types.Field{{Type: &types.BuiltinRef{TypeName: "int16"}}},
						},
					},
				},
			},
		},
		{
			name: "structtype-1",
			dir:  "testmod/convert/structtype_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/structtype_1.A": &types.TypeAlias{
						Package:  "testmod/convert/structtype_1",
						TypeName: "A",
						Type: &types.StructTypeRef{
							Fields: []*types.Field{{Name: "a", Type: &types.BuiltinRef{TypeName: "string"}}},
						},
					},
				},
			},
		},
		{
			name: "structref-1",
			dir:  "testmod/convert/structref_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/structref_1.B": &types.TypeDef{
						Package:  "testmod/convert/structref_1",
						TypeName: "B",
						Type:     &types.StructRef{Package: "testmod/convert/structref_1", TypeName: "A"},
					},
				},
			},
		},
		{
			name: "interfaceref-1",
			dir:  "testmod/convert/interfaceref_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/interfaceref_1.B": &types.TypeDef{
						Package:  "testmod/convert/interfaceref_1",
						TypeName: "B",
						Type:     &types.InterfaceRef{Package: "testmod/convert/interfaceref_1", TypeName: "A"},
					},
				},
			},
		},
		{
			name: "typedefref-1",
			dir:  "testmod/convert/typedefref_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typedefref_1.B": &types.TypeDef{
						Package:  "testmod/convert/typedefref_1",
						TypeName: "B",
						Type:     &types.TypeDefRef{Package: "testmod/convert/typedefref_1", TypeName: "A"},
					},
				},
			},
		},
		{
			name: "typealiasref-1",
			dir:  "testmod/convert/typealiasref_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typealiasref_1.B": &types.TypeDef{
						Package:  "testmod/convert/typealiasref_1",
						TypeName: "B",
						Type:     &types.TypeAliasRef{Package: "testmod/convert/typealiasref_1", TypeName: "A"},
					},
				},
			},
		},
		{
			name: "typeparamref-1",
			dir:  "testmod/convert/typeparamref_1",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typeparamref_1.B": &types.TypeDef{
						Package:  "testmod/convert/typeparamref_1",
						TypeName: "B",
						TypeParams: []*types.TypeParam{
							{Name: "T1", Order: 0},
							{Name: "T2", Order: 1},
						},
						Type: &types.StructRef{
							Package:    "testmod/convert/typeparamref_1",
							TypeName:   "A",
							TypeParams: []types.TypeRef{&types.TypeParamRef{Name: "T2"}},
						},
					},
				},
			},
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
