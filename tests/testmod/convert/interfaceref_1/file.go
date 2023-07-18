package interfaceref_1

type A interface {
	Method(b bool) error
}

type B A
