package tool

type (
	Tool interface {
		Tool()
	}

	TypeDef interface {
		typDef()
	}

	TypeRef interface {
		typRef()

		// String() string
		Equals(t TypeRef) bool
		// Modifiers() []parse.TypeRefModifier
	}
)

// func TypeRefFromDef(d TypeDef) TypeRef {
// 	switch dd := d.(type) {
// 	case *BuiltinDef:
// 		return BuiltinRefFromDef(dd)
// 	case *StructDef:
// 		return StructRefFromDef(dd)
// 	case *InterfaceDef:
// 		return InterfaceRefFromDef(dd)
// 	default:
// 		panic("not supported")
// 	}
// }
