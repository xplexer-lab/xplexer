package utils

func ResultOk[T any](result T) (T, error) {
	return result, nil
}
