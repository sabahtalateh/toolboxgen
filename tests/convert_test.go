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
				ctx := convert.NewContext().WithPackage(Package).WithImports(f.Imports).WithFiles(files)
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
			name: "struct",
			dir:  "testmod/convert/struct",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/struct.A": &types.Struct{
						Package:  "testmod/convert/struct",
						TypeName: "A",
					},
				},
			},
		},
		{
			name: "interface",
			dir:  "testmod/convert/interface",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/interface.A": &types.Interface{
						Package:  "testmod/convert/interface",
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
			dir:  "testmod/convert/interface2",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/interface2.A": &types.Interface{
						Package:  "testmod/convert/interface2",
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
			name: "typedef",
			dir:  "testmod/convert/typedef",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typedef.A": &types.TypeDef{
						Package:  "testmod/convert/typedef",
						TypeName: "A",
						Type:     &types.BuiltinRef{TypeName: "int"},
					},
				},
			},
		},
		{
			name: "typealias",
			dir:  "testmod/convert/typealias",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typealias.A": &types.TypeAlias{
						Package:  "testmod/convert/typealias",
						TypeName: "A",
						Type:     &types.BuiltinRef{TypeName: "int"},
					},
				},
			},
		},
		{
			name: "builtin",
			dir:  "testmod/convert/builtin",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/builtin.Bool": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Bool",
						Type:     &types.BuiltinRef{TypeName: "bool"},
					},
					"testmod/convert/builtin.Uint8": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Uint8",
						Type:     &types.BuiltinRef{TypeName: "uint8"},
					},
					"testmod/convert/builtin.Uint16": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Uint16",
						Type:     &types.BuiltinRef{TypeName: "uint16"},
					},
					"testmod/convert/builtin.Uint32": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Uint32",
						Type:     &types.BuiltinRef{TypeName: "uint32"},
					},
					"testmod/convert/builtin.Uint64": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Uint64",
						Type:     &types.BuiltinRef{TypeName: "uint64"},
					},
					"testmod/convert/builtin.Int8": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Int8",
						Type:     &types.BuiltinRef{TypeName: "int8"},
					},
					"testmod/convert/builtin.Int16": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Int16",
						Type:     &types.BuiltinRef{TypeName: "int16"},
					},
					"testmod/convert/builtin.Int32": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Int32",
						Type:     &types.BuiltinRef{TypeName: "int32"},
					},
					"testmod/convert/builtin.Int64": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Int64",
						Type:     &types.BuiltinRef{TypeName: "int64"},
					},
					"testmod/convert/builtin.Float32": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Float32",
						Type:     &types.BuiltinRef{TypeName: "float32"},
					},
					"testmod/convert/builtin.Float64": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Float64",
						Type:     &types.BuiltinRef{TypeName: "float64"},
					},
					"testmod/convert/builtin.Complex64": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Complex64",
						Type:     &types.BuiltinRef{TypeName: "complex64"},
					},
					"testmod/convert/builtin.Complex128": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Complex128",
						Type:     &types.BuiltinRef{TypeName: "complex128"},
					},
					"testmod/convert/builtin.String": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "String",
						Type:     &types.BuiltinRef{TypeName: "string"},
					},
					"testmod/convert/builtin.Int": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Int",
						Type:     &types.BuiltinRef{TypeName: "int"},
					},
					"testmod/convert/builtin.Uint": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Uint",
						Type:     &types.BuiltinRef{TypeName: "uint"},
					},
					"testmod/convert/builtin.Uintptr": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Uintptr",
						Type:     &types.BuiltinRef{TypeName: "uintptr"},
					},
					"testmod/convert/builtin.Byte": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Byte",
						Type:     &types.BuiltinRef{TypeName: "byte"},
					},
					"testmod/convert/builtin.Rune": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Rune",
						Type:     &types.BuiltinRef{TypeName: "rune"},
					},
					"testmod/convert/builtin.Any": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Any",
						Type:     &types.BuiltinRef{TypeName: "any"},
					},
					"testmod/convert/builtin.Comparable": &types.TypeDef{
						Package:  "testmod/convert/builtin",
						TypeName: "Comparable",
						Type:     &types.BuiltinRef{TypeName: "comparable"},
					},
				},
			},
		},
		{
			name: "map",
			dir:  "testmod/convert/map",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/map.A": &types.TypeDef{
						Package:  "testmod/convert/map",
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
			name: "chan",
			dir:  "testmod/convert/chan",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/chan.A": &types.TypeDef{
						Package:  "testmod/convert/chan",
						TypeName: "A",
						Type:     &types.ChanRef{Value: &types.BuiltinRef{TypeName: "string"}},
					},
				},
			},
		},
		{
			name: "functype",
			dir:  "testmod/convert/functype",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/functype.A": &types.TypeDef{
						Package:  "testmod/convert/functype",
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
			name: "structtype",
			dir:  "testmod/convert/structtype",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/structtype.A": &types.TypeAlias{
						Package:  "testmod/convert/structtype",
						TypeName: "A",
						Type: &types.StructTypeRef{
							Fields: []*types.Field{{Name: "a", Type: &types.BuiltinRef{TypeName: "string"}}},
						},
					},
				},
			},
		},
		{
			name: "interfacetype",
			dir:  "testmod/convert/interfacetype",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/interfacetype.A": &types.TypeAlias{
						Package:  "testmod/convert/interfacetype",
						TypeName: "A",
						Type: &types.InterfaceTypeRef{
							Fields: []*types.Field{{
								Name: "Method",
								Type: &types.FuncTypeRef{
									Params:  []*types.Field{{Name: "b", Type: &types.BuiltinRef{TypeName: "bool"}}},
									Results: []*types.Field{{Type: &types.BuiltinRef{TypeName: "error"}}},
								},
							}},
						},
					},
				},
			},
		},
		{
			name: "structref",
			dir:  "testmod/convert/structref",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/structref.B": &types.TypeDef{
						Package:  "testmod/convert/structref",
						TypeName: "B",
						Type:     &types.StructRef{Package: "testmod/convert/structref", TypeName: "A"},
					},
				},
			},
		},
		{
			name: "interfaceref",
			dir:  "testmod/convert/interfaceref",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/interfaceref.B": &types.TypeDef{
						Package:  "testmod/convert/interfaceref",
						TypeName: "B",
						Type:     &types.InterfaceRef{Package: "testmod/convert/interfaceref", TypeName: "A"},
					},
				},
			},
		},
		{
			name: "typedefref",
			dir:  "testmod/convert/typedefref",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typedefref.B": &types.TypeDef{
						Package:  "testmod/convert/typedefref",
						TypeName: "B",
						Type:     &types.TypeDefRef{Package: "testmod/convert/typedefref", TypeName: "A"},
					},
				},
			},
		},
		{
			name: "typealiasref",
			dir:  "testmod/convert/typealiasref",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typealiasref.B": &types.TypeDef{
						Package:  "testmod/convert/typealiasref",
						TypeName: "B",
						Type:     &types.TypeAliasRef{Package: "testmod/convert/typealiasref", TypeName: "A"},
					},
				},
			},
		},
		{
			name: "typeparamref",
			dir:  "testmod/convert/typeparamref",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/typeparamref.B": &types.TypeDef{
						Package:    "testmod/convert/typeparamref",
						TypeName:   "B",
						TypeParams: []*types.TypeParam{{Order: 0}, {Order: 1}},
						Type: &types.StructRef{
							Package:    "testmod/convert/typeparamref",
							TypeName:   "A",
							TypeParams: []types.TypeRef{&types.TypeParamRef{Order: 1}},
						},
					},
				},
			},
		},
		{
			name: "complex",
			dir:  "testmod/convert/complex",
			want: convertOut{
				types: map[string]types.Type{
					"testmod/convert/complex.A": &types.Struct{
						Package:    "testmod/convert/complex",
						TypeName:   "A",
						TypeParams: []*types.TypeParam{{Order: 0}, {Order: 1}, {Order: 2}, {Order: 3}},
						Fields: []*types.Field{
							{
								Name: "d",
								Type: &types.TypeParamRef{
									Order:     3,
									Modifiers: []types.Modifier{&types.Array{}},
								},
							},
							{
								Name: "s",
								Type: &types.BuiltinRef{
									Modifiers: []types.Modifier{&types.Pointer{}},
									TypeName:  "string",
								},
							},
						},
					},
					"testmod/convert/complex.B": &types.TypeDef{
						Package:    "testmod/convert/complex",
						TypeName:   "B",
						TypeParams: []*types.TypeParam{{Order: 0}, {Order: 1}, {Order: 2}},
						Type: &types.StructRef{
							Modifiers: []types.Modifier{&types.Array{}},
							Package:   "testmod/convert/complex",
							TypeName:  "A",
							TypeParams: []types.TypeRef{
								&types.InterfaceTypeRef{
									Modifiers: []types.Modifier{&types.Pointer{}},
									Fields: []*types.Field{{
										Name: "Func",
										Type: &types.FuncTypeRef{
											Params: []*types.Field{
												{
													Name: "b",
													Type: &types.TypeParamRef{
														Modifiers: []types.Modifier{&types.Pointer{}},
														Order:     1,
													},
												},
											},
											Results: []*types.Field{{Type: &types.BuiltinRef{TypeName: "string"}}},
										},
									}},
								},
								&types.TypeParamRef{Order: 2},
								&types.BuiltinRef{
									Modifiers: []types.Modifier{&types.Array{}},
									TypeName:  "float32",
								},
								&types.TypeParamRef{
									Modifiers: []types.Modifier{&types.Array{}, &types.Pointer{}, &types.Pointer{}},
									Order:     1,
								},
							},
							Fields: []*types.Field{
								{
									Name: "d",
									Type: &types.TypeParamRef{
										Modifiers: []types.Modifier{&types.Array{}, &types.Pointer{}, &types.Pointer{}},
										Name:      "T2",
										Order:     1,
									},
								},
								{
									Name: "s",
									Type: &types.BuiltinRef{
										Modifiers: []types.Modifier{&types.Pointer{}},
										TypeName:  "string",
									},
								},
							},
						},
					},
				},
			},
		},
		// {
		// 	name: "recursive_type",
		// 	dir:  "testmod/convert/forward",
		// 	want: convertOut{},
		// },
		// {
		// 	name: "recursive_type",
		// 	dir:  "testmod/convert/forward",
		// 	want: convertOut{},
		// },
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
