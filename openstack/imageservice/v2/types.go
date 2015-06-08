package v2

type ImageStatus string
const (
	ImageStatusActive ImageStatus = "active"
	// TODO
)

type ImageVisibility string
const (
	ImageVisibilityPublic ImageVisibility = "public"
	ImageVisibilityPrivate ImageVisibility = "private"
)
