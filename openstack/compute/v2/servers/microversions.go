package servers

// Tags represents a servers tags.
type Tags []string

// ExtractTags will extract the tags of a server.
// This requires the client to be set to microversion 2.26 or later.
func (r serverResult) ExtractTags() (Tags, error) {
	var s struct {
		Tags Tags `json:"tags"`
	}
	err := r.ExtractInto(&s)
	return s.Tags, err
}
