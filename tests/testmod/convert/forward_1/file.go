package forward_1

type ABC[T any] struct {
}

type XYZ[YYY any] struct {
	Y ABC[YYY]
}

type A[P1, P2 any] interface {
	Func2(XYZ[XYZ[XYZ[P2]]]) P2
}

type AA[X any] A[string, X]

type BBB[B100, D100, Z999 any] AA[Z999]
