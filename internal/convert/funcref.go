package convert

// func (c *Converter) convFuncRef(ctx context.Context, ref parse.FuncCallRef) (*tool.FuncRef, *errors.PositionedErr) {
// 	// if ref.PkgAlias == "" && builtin.Func(ref.FuncName) {
// 	// 	return nil, errors.BuiltinFunctionErr(ref.Position(), ref.FuncName)
// 	// }
// 	//
// 	// pkg, pathErr := c.packagePath(ctx, ref.PkgAlias)
// 	// if pathErr != nil {
// 	// 	return nil, errors.ParseError(ref.Position(), pathErr)
// 	// }
// 	//
// 	// fDef, err := c.findFuncDef(pkg, ref.FuncName, ref.Position())
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	//
// 	// if len(ref.TypeParams) != len(fDef.TypeParams) {
// 	// 	return nil, errors.InconsistentTypeParamsErr(ref.Position())
// 	// }
// 	//
// 	// fRef := tool.FuncRefFromDef(ref.Position(), fDef)
// 	// fRef.Code = ref.Code()
// 	// for i, tParam := range ref.TypeParams {
// 	// 	param, pErr := fRef.NthTypeParam(i)
// 	// 	if err != nil {
// 	// 		return nil, errors.ParseError(param.Position, pErr)
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
