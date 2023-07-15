package set

type Set[T comparable] map[T]struct{}

func Make[T comparable]() Set[T] {
	return make(map[T]struct{})
}

func (s Set[T]) Put(x T) bool {
	_, ok := s[x]
	s[x] = struct{}{}
	return ok
}

func (s Set[T]) Has(x T) bool {
	_, ok := s[x]
	return ok
}

func FromSlice[T comparable](s []T) Set[T] {
	res := Make[T]()

	for _, v := range s {
		res.Put(v)
	}

	return res
}

func FromMapKeys[K comparable, V any](m map[K]V) Set[K] {
	res := Make[K]()

	for k := range m {
		res.Put(k)
	}

	return res
}

func FromMapValues[K, V comparable](m map[K]V) Set[V] {
	res := Make[V]()

	for _, v := range m {
		res.Put(v)
	}

	return res
}
