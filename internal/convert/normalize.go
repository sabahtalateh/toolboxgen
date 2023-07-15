package convert

// func makeNormalizeTypeParametersList(fd *tool.FuncDef) map[string]string {
// 	var origTypeParams []tool.TypeParam
// 	if len(fd.TypeParams) != 0 {
// 		origTypeParams = fd.TypeParams
// 	} else {
// 		// origTypeParams = fd.Receiver.TypeParams
// 	}
//
// 	res := map[string]string{}
// 	for i, origTypeParam := range origTypeParams {
// 		newName := fmt.Sprintf("T%d", i+1)
// 		res[origTypeParam.Name] = newName
// 	}
//
// 	return res
// }
//
// func normalizeDefinedTypeParams(fd *tool.FuncDef) {
// 	renameMap := makeRenameMap(fd)
//
// 	// if _, ok := fd.Receiver.Type.(tool.ParametrizedRef); ok {
// 	// 	fd.Receiver.Type.(tool.ParametrizedRef).NormalizeTypeParams(renameMap)
// 	// }
//
// 	// for i, tp := range fd.Receiver.TypeParams {
// 	// 	if newName, ok := renameMap[tp.Name]; ok {
// 	// 		fd.Receiver.TypeParams[i].Params(newName)
// 	// 	}
// 	// }
//
// 	for i, tp := range fd.TypeParams {
// 		if newName, ok := renameMap[tp.Name]; ok {
// 			fd.TypeParams[i].Params(newName)
// 		}
// 	}
//
// 	for i := range fd.Parameters {
// 		switch p := fd.Parameters[i].Type.(type) {
// 		case *tool.TypeParamRef:
// 			p.Params(fmt.Sprintf("T%d", i+1))
// 		case tool.ParametrizedRef:
// 			p.NormalizeTypeParams(renameMap)
// 		}
// 	}
//
// 	for i := range fd.Results {
// 		switch p := fd.Results[i].(type) {
// 		case *tool.TypeParamRef:
// 			p.Params(fmt.Sprintf("T%d", i+1))
// 		case tool.ParametrizedRef:
// 			p.NormalizeTypeParams(renameMap)
// 		}
// 	}
//
// 	println(fd)
// }
//
// func makeRenameMap(fd *tool.FuncDef) map[string]string {
// 	renameMap := make(map[string]string)
// 	var typeParams []tool.TypeParam
// 	if fd.Receiver.Presented {
// 		// typeParams = fd.Receiver.TypeParams
// 	} else {
// 		typeParams = fd.TypeParams
// 	}
// 	for i, p := range typeParams {
// 		renameMap[p.Name] = fmt.Sprintf("T%d", i+1)
// 	}
//
// 	return renameMap
// }
