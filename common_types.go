package gophercloud

// Link is used for JSON (un)marshalling.
// It provides RESTful links to a resource.
type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}
