package forward_1

type A[T any] map[string]A[T]

// type CDE[T any] *interface {
// 	BB() [][]T
// }
//
// type ABC[T comparable] map[T]CDE[**T]
//
// type XYZ[YYY any] struct {
// 	Y ABC[uintptr]
// }
//
// type A[P1, P2 any] interface {
// 	Func2(XYZ[XYZ[XYZ[P2]]]) P2
// }
//
// type AA[X any] A[string, X]
//
// type BBB[B100, D100, Z999 any] AA[Z999]
