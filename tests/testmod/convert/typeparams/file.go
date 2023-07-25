package typeparams

// type A4[U, V any] A3[[]A2[A1[*U], V]]
//

type A1[T any] struct {
	t []T
	w *string
}

type Alias = A3[string]

type A2[A, B any] interface {
	Int2
	M1(a *A1[*A], b int) (A1[A1[[]B]], error)
}

type A4[U, V any] A3[[]A2[A1[*U], V]]

type A3[X any] A1[*X]

type Int2 interface {
	F2()
}

// Actual A3
//
// type A3[X any] -> A2[*X, []struct{}]
// ->
// A2[*X, []struct{}] interface {
//	M1(*A1[**X]) A1[A1[[][]struct{}]]
// }

// Actual A4
//
// A4[U, V any] -> A3[[]A2[A1[*U], V]]
// ->
// A3[[]A2[A1[*U], V]] -> A2[*[]A2[A1[*U], V], []struct{}]
// ->
// type A2[*[]A2[A1[*U], V], []struct{}] interface {
//	M1(*A1[**[]A2[A1[*U], V]]) A1[A1[[][]struct{}]]
// }
