package tool

// func NodePosition(t TypeRef) token.NodePosition {
// 	switch x := t.(type) {
// 	case *BuiltinRef:
// 		return x.NodePosition
// 	case *StructRef:
// 		return x.NodePosition
// 	case *InterfaceRef:
// 		return x.NodePosition
// 	case *TypeParamRef:
// 		return x.NodePosition
// 	default:
// 		panic("not implemented")
// 	}
// }
//
// func ModifiersToString(mm []parse.TypeRefModifier) string {
// 	strs := slices.Map(mm, func(el parse.TypeRefModifier) string { return el.String() })
// 	return slices.Reduce(strs, "", func(el string, acc string) string { return acc + el })
// }
//
// func prependModifiers(mm []parse.TypeRefModifier, t TypeRef) TypeRef {
// 	switch tt := t.(type) {
// 	case *StructRef:
// 		tt.Mods = append(mm, tt.Mods...)
// 	case *InterfaceRef:
// 		tt.Mods = append(mm, tt.Mods...)
// 	case *BuiltinRef:
// 		tt.Mods = append(mm, tt.Mods...)
// 	case *TypeParamRef:
// 		tt.Mods = append(mm, tt.Mods...)
// 	}
//
// 	return t
// }
