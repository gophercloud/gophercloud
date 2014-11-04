package vips

// VIP represents a Virtual IP API resource.
type VIP struct {
	Address string `json:"address,omitempty"`
	ID      int    `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Version string `json:"ipVersion,omitempty" mapstructure:"ipVersion"`
}
