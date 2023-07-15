package tool

// type BuiltinDef struct {
// 	TypeName  string
// 	Modifiers []parse.TypeRefModifier
// 	Position  token.Position
// }
//
// func (b BuiltinDef) typDef() {}
//
// type BuiltinRef struct {
// 	TypeName string
// 	Mods     []parse.TypeRefModifier
// 	Position token.Position
// }
//
// func BuiltinRefFromDef(d *BuiltinDef) *BuiltinRef {
// 	return &BuiltinRef{
// 		TypeName: d.TypeName,
// 		Mods:     d.Modifiers,
// 		Position: d.Position,
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
