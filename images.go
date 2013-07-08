package gophercloud

// ImageLink provides a reference to a image by either ID or by direct URL.
// Some services use just the ID, others use just the URL.
// This structure provides a common means of expressing both in a single field.
type ImageLink struct {
	Id    string `json:"id"`
	Links []Link `json:"links"`
}
