package complex

type A[A, B, C, D any] struct {
	d []D
	s *string
}

// type B[AA, BB, CC comparable] = []A[*interface{ Func(b *BB) string }, CC, []float32, **BB]
// ->
// []A[*interface{ Func(b *BB) string }, CC, []float32, **BB] = []struct {
//	d []**BB
//	s *string
// }

type B[AA, BB, CC comparable] []A[*interface{ Func(b *BB) string }, CC, []float32, **BB]
