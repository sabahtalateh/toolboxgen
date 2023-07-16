package convert

// func (c *Converter) convFuncRef(ctx context.Context, ref parse.FuncCallRef) (*tool.FuncRef, *errors.PositionedErr) {
// 	// if ref.PkgAlias == "" && builtin.Func(ref.FuncName) {
// 	// 	return nil, errors.BuiltinFunctionErr(ref.NodePosition(), ref.FuncName)
// 	// }
// 	//
// 	// pkg, pathErr := c.packagePath(ctx, ref.PkgAlias)
// 	// if pathErr != nil {
// 	// 	return nil, errors.ParseError(ref.NodePosition(), pathErr)
// 	// }
// 	//
// 	// fDef, err := c.findFuncDef(pkg, ref.FuncName, ref.NodePosition())
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	//
// 	// if len(ref.TypeParams) != len(fDef.TypeParams) {
// 	// 	return nil, errors.InconsistentTypeParamsErr(ref.NodePosition())
// 	// }
// 	//
// 	// fRef := tool.FuncRefFromDef(ref.NodePosition(), fDef)
// 	// fRef.Code = ref.Code()
// 	// for i, tParam := range ref.TypeParams {
// 	// 	param, pErr := fRef.NthTypeParam(i)
// 	// 	if err != nil {
// 	// 		return nil, errors.ParseError(param.NodePosition, pErr)
// 	// 	}
// 	//
// 	// 	effective, err := c.TypeRef(ctx, tParam, nil)
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	fRef.SetEffectiveParamRecursive(param.Name, effective)
// 	// }
//
// 	// return fRef, nil
//
// 	return nil, nil
// }
