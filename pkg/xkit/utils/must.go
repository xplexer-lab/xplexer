package utils

func Must[T any](r T, err error) T {
	if err != nil {
		panic(err)
	}

	return r
}

func NewMust1[R, T1 any](fn func(T1) (R, error)) func(T1) R {
	return func(t1 T1) R {
		res, err := fn(t1)

		if err != nil {
			panic(err)
		}

		return res
	}
}

func NewMust2[R, T1, T2 any](fn func(T1, T2) (R, error)) func(T1, T2) R {
	return func(t1 T1, t2 T2) R {
		res, err := fn(t1, t2)

		if err != nil {
			panic(err)
		}

		return res
	}
}
