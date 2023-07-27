package tests

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func check2[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}
