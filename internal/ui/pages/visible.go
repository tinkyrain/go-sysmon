package pages

func visible[T any](elements []T, count int) []T {
	if count < 1 {
		return nil
	}

	if len(elements) > count {
		return elements[:count]
	}

	return elements
}
