package gophercloud

// Pagination stores information that's necessary to enumerate through pages of results.
type Pagination struct {

	// Next is the full URL to the next page of results, or nil if this is the last page.
	Next *string `json:"next,omitempty"`

	// Previous is the full URL to the previous page of results, or nil if this is the first page.s
	Previous *string `json:"previous,omitempty"`
}
