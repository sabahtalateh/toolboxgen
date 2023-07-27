package tests

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func unwrap[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}
