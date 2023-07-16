package tool

// type TypeParamRef struct {
// 	Name     string
// 	// Mods     []parse.TypeRefModifier
// 	NodePosition token.NodePosition
// }
//
// func (t *TypeParamRef) Params(name string) {
// 	t.Name = name
// }
//
// func (t *TypeParamRef) typRef() {
// }
//
// func (t *TypeParamRef) Equals(tp TypeRef) bool {
// 	switch t2 := tp.(type) {
// 	case *TypeParamRef:
// 		if !parse.ModifiersEquals(t.Modifiers(), t2.Modifiers()) {
// 			return false
// 		}
//
// 		if t.Name != t2.Name {
// 			return false
// 		}
//
// 		return true
// 	default:
// 		return false
// 	}
// }
//
// func (t *TypeParamRef) Modifiers() []parse.TypeRefModifier {
// 	return t.Mods
// }
