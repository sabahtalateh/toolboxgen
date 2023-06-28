package syntax

// type FuncDef struct {
// 	files *token.FileSet
//
// 	FuncDecl     *ast.FuncDecl
// 	Code         string
// 	PkgAlias     string
// 	FuncName     string
// 	Presented bool
// 	Receiver     TypeRef
// 	TypeParams   []TypeParam
// 	Args         []TypeRef
// 	Results      []TypeRef
// 	Position     token.Position
// 	Err          *errors.PositionedErr
// }
//
// func ParseFuncDef(def *ast.FuncDecl, files *token.FileSet) FuncDef {
// 	// fd := newFuncDef(def, files)
// 	// fd.visitFuncDecl(def)
// 	// return *fd
// 	return
// }
//
// func newFuncDef(def *ast.FuncDecl, files *token.FileSet) *FuncDef {
// 	fd := &FuncDef{files: files}
// 	fd.Code, fd.Err = code(def, files.Position(def.Pos()))
//
// 	// remove function body. leave just function signature
// 	var body string
// 	body, fd.Err = code(def.Body, files.Position(def.Pos()))
// 	fd.Code = strings.TrimSpace(strings.TrimSuffix(fd.Code, body))
//
// 	return fd
// }
//
// func (f *FuncDef) visitFuncDecl(decl *ast.FuncDecl) {
// 	f.FuncName = decl.Name.Name
//
// 	if decl.Recv != nil {
// 		f.Presented = true
// 		f.Receiver = ParseTypeRef(decl.Recv.List[0].Type, f.files)
// 	}
//
// 	if decl.Type.TypeParams != nil {
// 		for _, tp := range decl.Type.TypeParams.List {
// 			f.TypeParams = append(f.TypeParams, ParseTypeParams(tp, f.files)...)
// 		}
// 	}
//
// 	if decl.Type.Params != nil {
// 		for _, tp := range decl.Type.Params.List {
// 			f.Args = append(f.Args, ParseTypeRef(tp.Type, f.files))
// 		}
// 	}
//
// 	if decl.Type.Results != nil {
// 		for _, r := range decl.Type.Results.List {
// 			f.Results = append(f.Results, ParseTypeRef(r.Type, f.files))
// 		}
// 	}
// }
