package convert

type component struct {
	converter *Converter
}

// func (c *component) convert(
// 	ctx context.Context,
// 	calls []parse.FunctionCall,
// ) (comp *di.Component, err *errors.PositionedErr) {
// 	firstCall := calls[0]
// 	if err = errors.Check(
// 		func() *errors.PositionedErr { return validation.FunctionCallHasExactlyNArguments(firstCall, 1) },
// 	); err != nil {
// 		return nil, err
// 	}
//
// 	comp = &di.Component{With: map[string][]*di.With{}}
// 	arg := firstCall.Args[0]
// 	var funcRef *tool.FuncRef
// 	switch val := arg.Val.(type) {
// 	case *parse.FuncCallRef:
// 		funcRef, err = c.converter.convFuncRef(ctx, *val)
// 		if err != nil {
// 			return nil, err
// 		}
// 	default:
// 		return nil, errors.FunctionRefExpectedErr(arg.Val.NodePosition())
// 	}
//
// 	if err = comp.EnrichWithFunction(funcRef); err != nil {
// 		return nil, err
// 	}
//
// 	for _, call := range calls[1:] {
// 		switch call.FuncName {
// 		case "name":
// 			if err = c.setName(comp, call); err != nil {
// 				return nil, err
// 			}
// 		case "With":
// 			key, ww, err := c.convertWith(ctx, call)
// 			if err != nil {
// 				return nil, err
// 			}
// 			comp.With[key] = ww
// 		default:
// 			return nil, errors.Errorf(call.NodePosition, "unsupported function")
// 		}
// 	}
//
// 	return comp, nil
// }
//
// func (c *component) setName(comp *di.Component, nameCall parse.FunctionCall) *errors.PositionedErr {
// 	if err := errors.Check(
// 		func() *errors.PositionedErr { return validation.FunctionCallHasExactlyNArguments(nameCall, 1) },
// 	); err != nil {
// 		return err
// 	}
//
// 	nameArg := nameCall.Args[0]
// 	switch name := nameArg.Val.(type) {
// 	case *parse.String:
// 		if name.Val == "" {
// 			return errors.Errorf(nameArg.Val.NodePosition(), "component name should not be empty")
// 		}
// 		if comp.name != "" {
// 			return errors.Errorf(nameArg.Val.NodePosition(), "component name already set")
// 		}
//
// 		comp.name = name.Val
// 		return nil
// 	default:
// 		return errors.StringLitExpectedErr(nameArg.Val.NodePosition())
// 	}
// }
//
// func (c *component) convertWith(
// 	ctx context.Context,
// 	call parse.FunctionCall,
// ) (string, []*di.With, *errors.PositionedErr) {
// 	if err := errors.Check(
// 		func() *errors.PositionedErr { return validation.FunctionCallHasAtMinimumNArguments(call, 2) },
// 	); err != nil {
// 		return "", nil, err
// 	}
//
// 	key, err := convertWithKey(call.Args[0].Val)
// 	if err != nil {
// 		return "", nil, err
// 	}
//
// 	var res []*di.With
//
// 	for _, arg := range call.Args[1:] {
// 		switch ref := arg.Val.(type) {
// 		case *parse.FuncCallRef:
// 			fRef, err := c.converter.convFuncRef(ctx, *ref)
// 			if err != nil {
// 				return "", nil, err
// 			}
//
// 			w, err := convertWithValue(fRef)
// 			w.Key = key
// 			if err != nil {
// 				return "", nil, err
// 			}
//
// 			res = append(res, w)
// 		default:
// 			return "", nil, errors.FunctionRefExpectedErr(arg.Val.NodePosition())
// 		}
// 	}
//
// 	return key, res, nil
// }
//
// func convertWithKey(arg parse.FunctionCallArgValue) (string, *errors.PositionedErr) {
// 	switch ka := arg.(type) {
// 	case *parse.String:
// 		if ka.Val == "" {
// 			return "", errors.Errorf(ka.NodePosition(), "`With` key can not be empty")
// 		}
//
// 		return ka.Val, nil
// 	default:
// 		return "", errors.Errorf(ka.NodePosition(), "`With` key should be string literal")
// 	}
// }
//
// func convertWithValue(f *tool.FuncRef) (*di.With, *errors.PositionedErr) {
// 	if err := errors.Check(
// 		func() *errors.PositionedErr { return validation.FuncDefHasAtMinimumNResults(f.Def, 1) },
// 		func() *errors.PositionedErr { return validation.FuncDefHasAtMaximumNResults(f.Def, 2) },
// 	); err != nil {
// 		return nil, err
// 	}
//
// 	w := &di.With{Function: f, Type: f.Results[0]}
// 	if len(f.Results) == 2 {
// 		switch res := f.Results[1].(type) {
// 		case *tool.BuiltinRef:
// 			if !res.IsError() {
// 				return nil, errors.ErrorExpectedErr(res.NodePosition)
// 			}
// 			w.WithError = true
// 		default:
// 			return nil, errors.ErrorExpectedErr(tool.NodePosition(res))
// 		}
// 	}
//
// 	return w, nil
// }
