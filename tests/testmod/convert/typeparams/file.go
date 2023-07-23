package typeparams

type A2[A, B any] interface {
	Int2
	M1(a *A1[*A], b int) (A1[A1[[]B]], error)
}

type Int2 interface {
	F2()
}

type Alias = A3[string]

type A1[T any] struct {
	d interface {
		F(string2 string)
		F2()
	}
	m  map[string][]T
	c  chan []string
	Al *Alias
	t  []T
	w  *string
}

type A4[U, V any] A3[[]A2[A1[*U], V]]

type A3[X any] A2[*X, []struct {
	a string
	b string
}]

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
