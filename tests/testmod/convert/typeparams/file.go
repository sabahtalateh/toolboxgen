package typeparams

type C[A, B any] *interface {
	F(a A) (B, error)
}

type A[X, Y, Z any] struct {
	f  *X
	f2 []C[Y, uint]
	f3 *Z
}

// type B[T any] A[[]string, []T, []struct{ a A[**T, struct{}, **T] }]

// type D C[string, struct{}]
