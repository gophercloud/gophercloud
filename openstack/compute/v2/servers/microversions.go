package servers

// ExtractDescription will extract the description of a server.
// This requires the client to be set to microversion 2.19 or later.
func (r serverResult) ExtractDescription() (string, error) {
	var s struct {
		Description string `json:"description"`
	}
	err := r.ExtractInto(&s)
	return s.Description, err
}
