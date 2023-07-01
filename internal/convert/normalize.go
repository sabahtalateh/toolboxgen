package convert

import (
	"fmt"

	"github.com/sabahtalateh/toolboxgen/internal/tool"
)

func makeNormalizeTypeParametersList(fd *tool.FuncDef) map[string]string {
	var origTypeParams []tool.TypeParam
	if len(fd.TypeParams) != 0 {
		origTypeParams = fd.TypeParams
	} else {
		origTypeParams = fd.Receiver.TypeParams
	}

	res := map[string]string{}
	for i, origTypeParam := range origTypeParams {
		newName := fmt.Sprintf("T%d", i+1)
		res[origTypeParam.Name] = newName
	}

	return res
}

func normalizeDefinedTypeParams(fd *tool.FuncDef) {
	var origTypeParams []tool.TypeParam
	if len(fd.TypeParams) != 0 {
		origTypeParams = fd.TypeParams
	} else {
		origTypeParams = fd.Receiver.TypeParams
	}

	for i := range origTypeParams {
		origTypeParams[i].Name = fmt.Sprintf("T%d", i+1)
	}
}
