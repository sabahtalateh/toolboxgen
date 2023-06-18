package syntax

import (
	"go/ast"
	"go/token"

	"github.com/sabahtalateh/toolboxgen/internal/errors"
)

// TODO посмотреть что можно перевести на TypeParam
type TypeParam struct {
	Code     string
	Name     string
	Position token.Position
	Error    *errors.PositionedErr
}

type typeParams struct {
	files  *token.FileSet
	result []TypeParam
}

func ParseTypeParams(field *ast.Field, files *token.FileSet) []TypeParam {
	t := newTypeParam(files)
	t.visitField(field)
	return t.result
}

func newTypeParam(files *token.FileSet) *typeParams {
	return &typeParams{files: files}
}

func (t *typeParams) visitField(f *ast.Field) {
	for _, name := range f.Names {
		t.visitIdent(name)
	}
}

func (t *typeParams) visitIdent(id *ast.Ident) {
	param := TypeParam{
		Name:     id.Name,
		Position: t.files.Position(id.Pos()),
	}

	param.Code, param.Error = nodeToString(id, t.files.Position(id.Pos()))
	t.result = append(t.result, param)
}
