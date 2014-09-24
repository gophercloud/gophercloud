package gophercloud

// MaybeString takes a string that might be a zero-value, and either returns a
// pointer to its address or a nil value (i.e. empty pointer). This is useful
// for converting zero values in options structs when the end-user hasn't
// defined values. Those zero values need to be nil in order for the JSON
// serialization to ignore them.
func MaybeString(original string) *string {
	if original != "" {
		return &original
	}
	return nil
}
