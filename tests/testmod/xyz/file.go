package xyz

type XX[T any] struct{}

type String2[A, B, C, D any] XX[C]
type String3[A2, B2, C2 any] String2[C2, string, B2, B2]

// func A[TT, YY any](x *String3[*[]YY, TT, int8]) {
// }

type SS[AAAA, BBBB any] struct {
}

func (s *SS[TT, YY]) A(x *String3[*[]YY, TT, int8]) {
}
