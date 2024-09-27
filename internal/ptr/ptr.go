package ptr

func To[T any](v T) *T {
	return &v
}

func From[T any](p *T) T {
	if p == nil {
		var v T
		return v
	}
	return *p
}
