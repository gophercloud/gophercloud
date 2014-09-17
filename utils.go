package utils

func MaybeString(original string) *string {
	if original != "" {
		return &original
	}
	return nil
}

func MaybeInt(original int) {
	if original != 0 {
		return &original
	}
	return nil
}
