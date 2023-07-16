package tool

// type BuiltinDef struct {
// 	TypeName  string
// 	Modifiers []parse.TypeRefModifier
// 	NodePosition  token.NodePosition
// }
//
// func (b BuiltinDef) typDef() {}
//
// type BuiltinRef struct {
// 	TypeName string
// 	Mods     []parse.TypeRefModifier
// 	NodePosition token.NodePosition
// }
//
// func BuiltinRefFromDef(d *BuiltinDef) *BuiltinRef {
// 	return &BuiltinRef{
// 		TypeName: d.TypeName,
// 		Mods:     d.Modifiers,
// 		NodePosition: d.NodePosition,
// 	}
// }
//
// func (b *BuiltinRef) typRef() {}
//
// func (b *BuiltinRef) String() string {
// 	return fmt.Sprintf("%s%s", ModifiersToString(b.Modifiers()), b.TypeName)
// }
//
// func (b *BuiltinRef) Equals(t TypeRef) bool {
// 	switch t2 := t.(type) {
// 	case *BuiltinRef:
// 		if b.TypeName != t2.TypeName {
// 			return false
// 		}
//
// 		return true
// 	default:
// 		return false
// 	}
// }
//
// func (b *BuiltinRef) Modifiers() []parse.TypeRefModifier {
// 	return b.Mods
// }
//
// func (b *BuiltinRef) IsError() bool {
// 	return b.TypeName == "error"
// }
