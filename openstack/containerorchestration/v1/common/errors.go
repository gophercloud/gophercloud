package common

// ErrorResponse represents failed response
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error represents a Magnum API error message
type Error struct {
	Status    int    `json:"status"`
	Code      string `json:"code"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	RequestID string `json:"request_id,omitempty"`
}

func (r ErrorResponse) Error() string {
	var msg string
	for _, e := range r.Errors {
		msg += e.Detail + "\n"
	}

	return msg
}
