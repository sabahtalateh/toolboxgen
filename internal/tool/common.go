package tool

// func Position(t TypeRef) token.Position {
// 	switch x := t.(type) {
// 	case *BuiltinRef:
// 		return x.Position
// 	case *StructRef:
// 		return x.Position
// 	case *InterfaceRef:
// 		return x.Position
// 	case *TypeParamRef:
// 		return x.Position
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
