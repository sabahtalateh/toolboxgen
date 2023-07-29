package interfaceref

type A interface {
	Method(b bool) error
}

type B A
