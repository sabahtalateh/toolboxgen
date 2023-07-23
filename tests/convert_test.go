package tests

import (
	"github.com/sabahtalateh/toolboxgen/internal/inspect"
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
					ins := inspect.Type(
						inspect.EmptyContext().WithTrimPackage("testmod/convert/typeparams"),
						t,
					)
					println(ins)
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
					// type A[A, B, C, D any] struct {
					//	d []D
					//	s *string
					// }
					// type A[A, B, C, D any]
					"testmod/convert/complex.A": &types.Struct{
						Package:    "testmod/convert/complex",
						TypeName:   "A",
						TypeParams: []*types.TypeParam{{Order: 0}, {Order: 1}, {Order: 2}, {Order: 3}},
						// d []D
						// s *string
						Fields: []*types.Field{
							// d []D
							{
								Name: "d",
								Type: &types.TypeParamRef{
									Order:     3,
									Modifiers: []types.Modifier{&types.Array{}},
								},
							},
							// s *string
							{
								Name: "s",
								Type: &types.BuiltinRef{
									Modifiers: []types.Modifier{&types.Pointer{}},
									TypeName:  "string",
								},
							},
						},
					},
					// type B[AA, BB, CC comparable] []A[*interface{ Func(b *BB) string }, CC, []float32, **BB]
					"testmod/convert/complex.B": &types.TypeDef{
						Package:    "testmod/convert/complex",
						TypeName:   "B",
						TypeParams: []*types.TypeParam{{Order: 0}, {Order: 1}, {Order: 2}},
						// []A[*interface{ Func(b *BB) string }, CC, []float32, **BB] = []struct {
						//	d []**BB
						//	s *string
						// }
						Type: &types.StructRef{
							Modifiers: []types.Modifier{&types.Array{}},
							Package:   "testmod/convert/complex",
							TypeName:  "A",
							TypeParams: []types.TypeRef{
								// *interface{ Func(b *BB) string }
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
								// CC
								&types.TypeParamRef{Order: 2},
								// []float32
								&types.BuiltinRef{
									Modifiers: []types.Modifier{&types.Array{}},
									TypeName:  "float32",
								},
								// **BB
								&types.TypeParamRef{
									Modifiers: []types.Modifier{&types.Pointer{}, &types.Pointer{}},
									Order:     1,
								},
							},
							// d []**BB
							// s *string
							Fields: []*types.Field{
								{
									Name: "d",
									Type: &types.TypeParamRef{
										Modifiers: []types.Modifier{&types.Array{}, &types.Pointer{}, &types.Pointer{}},
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
		{
			name: "typeparams",
			dir:  "testmod/convert/typeparams",
			want: convertOut{
				types: map[string]types.Type{
					// type A1[T any] struct {
					//	t []T
					// }
					"testmod/convert/typeparams.A1": &types.Struct{
						TypeParams: []*types.TypeParam{{Order: 0}},
						Package:    "testmod/convert/typeparams",
						TypeName:   "A1",
						// t []T
						Fields: []*types.Field{{
							Name: "t",
							Type: &types.TypeParamRef{Order: 0, Modifiers: []types.Modifier{&types.Array{}}},
						}},
					},
					// type A2[A, B any] interface {
					// 	M1(*A1[*A]) A1[A1[[]B]]
					// }
					"testmod/convert/typeparams.A2": &types.Interface{
						Package:    "testmod/convert/typeparams",
						TypeName:   "A2",
						TypeParams: []*types.TypeParam{{Order: 0}, {Order: 1}},
						Methods: []*types.Field{
							// M1(*A1[*A]) A1[A1[[]B]]
							{
								Name: "M1",
								Type: &types.FuncTypeRef{
									Params: []*types.Field{
										// *A1[*A]
										{
											Type: &types.StructRef{
												Package:   "testmod/convert/typeparams",
												TypeName:  "A1",
												Modifiers: []types.Modifier{&types.Pointer{}},
												TypeParams: []types.TypeRef{
													// *A
													&types.TypeParamRef{
														Modifiers: []types.Modifier{&types.Pointer{}},
														Order:     0,
													},
												},
											},
										},
									},
									Results: []*types.Field{
										// A1[A1[[]B]]
										{
											Type: &types.StructRef{
												Package:  "testmod/convert/typeparams",
												TypeName: "A1",
												TypeParams: []types.TypeRef{
													// A1[[]B]
													&types.StructRef{
														Package:  "testmod/convert/typeparams",
														TypeName: "A1",
														TypeParams: []types.TypeRef{
															// []B
															&types.TypeParamRef{
																Order:     1,
																Modifiers: []types.Modifier{&types.Array{}},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					// type A3[X any] -> A2[*X, []struct{}]
					// ->
					// A2[*X, []struct{}] interface {
					//	M1(*A1[**X]) A1[A1[[][]struct{}]]
					// }
					"testmod/convert/typeparams.A3":
					// type A3[X any]
					&types.TypeDef{
						Package:    "testmod/convert/typeparams",
						TypeName:   "A3",
						TypeParams: []*types.TypeParam{{Order: 0}},
						// A2[*X, []struct{}] interface {
						//	M1(*A1[**X]) A1[A1[[][]struct{}]]
						// }
						Type: &types.InterfaceRef{
							Package:  "testmod/convert/typeparams",
							TypeName: "A2",
							// [*X, []struct{}]
							TypeParams: []types.TypeRef{
								// *X
								&types.TypeParamRef{Modifiers: []types.Modifier{&types.Pointer{}}, Order: 0},
								// []struct{}
								&types.StructTypeRef{Modifiers: []types.Modifier{&types.Array{}}},
							},
							// M1(*A1[**X]) A1[A1[[][]struct{}]]
							Methods: []*types.Field{{
								Name: "M1",
								Type: &types.FuncTypeRef{
									Params: []*types.Field{{
										// *A1[**X]
										// ->
										// A1[**X] struct {
										// 	t []**X
										// }
										Type: &types.StructRef{
											Package:   "testmod/convert/typeparams",
											TypeName:  "A1",
											Modifiers: []types.Modifier{&types.Pointer{}},
											// [**X]
											TypeParams: []types.TypeRef{
												// **X
												&types.TypeParamRef{
													Modifiers: []types.Modifier{&types.Pointer{}, &types.Pointer{}},
													Order:     0,
												},
											},
											Fields: []*types.Field{
												// t []**X
												{
													Name: "t",
													Type: &types.TypeParamRef{
														Order: 0,
														Modifiers: []types.Modifier{
															&types.Array{}, &types.Pointer{}, &types.Pointer{},
														},
													},
												},
											},
										},
									}},
									Results: []*types.Field{{
										// A1[A1[[][]struct{}]] -> type A1[A1[[][]struct{}]] struct {
										//	t []A1[[][]struct{}]
										// }
										Type: &types.StructRef{
											Package:  "testmod/convert/typeparams",
											TypeName: "A1",
											TypeParams: []types.TypeRef{
												// A1[[][]struct{}] -> type A1[[][]struct{}] struct {
												// 	t [][][]struct{}
												// }
												&types.StructRef{
													Package:  "testmod/convert/typeparams",
													TypeName: "A1",
													TypeParams: []types.TypeRef{
														// [][]struct{}
														&types.StructTypeRef{
															Modifiers: []types.Modifier{&types.Array{}, &types.Array{}},
														},
													},
													// t [][][]struct{}
													Fields: []*types.Field{{
														Name: "t",
														Type: &types.StructTypeRef{
															Modifiers: []types.Modifier{&types.Array{}, &types.Array{}},
														},
													}},
												},
											},
											// t []A1[[][]struct{}]
											Fields: []*types.Field{
												// t []A1[[][]struct{}]
												{
													Name: "t",
													Type: &types.StructRef{
														Package:   "testmod/convert/typeparams",
														TypeName:  "A1",
														Modifiers: []types.Modifier{&types.Array{}},
														TypeParams: []types.TypeRef{
															&types.StructTypeRef{
																Modifiers: []types.Modifier{
																	&types.Array{}, &types.Array{},
																},
															},
														},
													},
												},
											},
										},
									}},
								},
							}},
						},
					},
					"testmod/convert/typeparams.A4":
					// type A4[U, V any] A3[[]A2[A1[*U], V]]
					&types.TypeDef{
						Package:    "testmod/convert/typeparams",
						TypeName:   "A4",
						TypeParams: []*types.TypeParam{{Order: 0}, {Order: 1}},
						// A3[[]A2[A1[*U], V]] -> A2[*[]A2[A1[*U], V], []struct{}]
						Type: &types.TypeDefRef{
							Package:  "testmod/convert/typeparams",
							TypeName: "A3",
							// A2[*[]A2[A1[*U], V], []struct{}] -> type A2[*[]A2[A1[*U], []struct{}] interface {
							//	M1(*A1[**[]A2[A1[*U]]) A1[A1[[][]struct{}]]
							// }
							Type: &types.InterfaceRef{
								Package:  "testmod/convert/typeparams",
								TypeName: "A2",
								Methods: []*types.Field{
									// M1(*A1[**[]A2[A1[*U]]) A1[A1[[][]struct{}]]
									{
										Name: "M1",
										Type: &types.FuncTypeRef{
											// *A1[**[]A2[A1[*U]]
											Params: []*types.Field{
												{
													// *A1[**[]A2[A1[*U]] -> type A1[**[]A2[A1[*U]] struct {
													//	t []**[]A2[A1[*U]
													// }
													Type: &types.StructRef{
														Package:  "testmod/convert/typeparams",
														TypeName: "A1",
														Fields: []*types.Field{
															// t []**[]A2[A1[*U]]
															{
																Name: "t",
																Type: &types.InterfaceRef{
																	Package:  "testmod/convert/typeparams",
																	TypeName: "A2",
																	Methods:  nil, // <--
																	Modifiers: []types.Modifier{
																		&types.Array{},
																		&types.Pointer{},
																		&types.Pointer{},
																		&types.Array{},
																	},
																	// [A1[*U]]
																	TypeParams: []types.TypeRef{
																		// A1[*U] -> type A1[*U] struct {
																		//	t []*U
																		// }
																		&types.StructRef{
																			Package:  "testmod/convert/typeparams",
																			TypeName: "A1",
																			// t []*U
																			Fields: []*types.Field{
																				// t []*U
																				{
																					Name: "t",
																					Type: &types.TypeParamRef{
																						Order: 0,
																						Modifiers: []types.Modifier{
																							&types.Array{},
																							&types.Pointer{},
																						},
																					},
																				},
																			},
																			TypeParams: []types.TypeRef{
																				&types.TypeParamRef{
																					Order: 0,
																					Modifiers: []types.Modifier{
																						&types.Pointer{},
																					},
																				},
																			},
																		},
																		&types.TypeParamRef{Order: 1},
																	},
																},
															},
														},
														Modifiers: []types.Modifier{&types.Pointer{}},
														// **[]A2[A1[*U]
														TypeParams: nil, // <--
													},
													Position: token.Position{},
													Declared: "",
												},
											}, // <--
											// A1[A1[[][]struct{}]]
											Results: nil, // <--
										},
									},
								},
								// *[]A2[A1[*U], V], []struct{}
								TypeParams: []types.TypeRef{
									// *[]A2[A1[*U], V]
									&types.InterfaceRef{
										Package:    "testmod/convert/typeparams",
										TypeName:   "A2",
										Modifiers:  []types.Modifier{&types.Pointer{}, &types.Array{}},
										TypeParams: []types.TypeRef{}, // <--
										Methods:    []*types.Field{},  // <--
									},
									// []struct{}
									&types.StructTypeRef{Modifiers: []types.Modifier{&types.Array{}}},
								},
							},
							TypeParams: []types.TypeRef{
								// []A2[A1[*U], V]
								&types.InterfaceRef{
									Package:   "testmod/convert/typeparams",
									TypeName:  "A2",
									Modifiers: []types.Modifier{&types.Array{}},
									// [A1[*U], V]
									TypeParams: []types.TypeRef{}, // <--
									Methods:    []*types.Field{},  // <--
								},
							},
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
