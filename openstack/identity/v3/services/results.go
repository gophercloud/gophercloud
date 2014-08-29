package services

// ServiceResult is the result of a list or information query.
type ServiceResult struct {
	Description *string `json:"description,omitempty"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
}
