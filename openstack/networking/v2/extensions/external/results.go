package external

// NetworkExt represents a decorated form of a Network with based on the
// "external-net" extension.
type NetworkExt struct {
	// Specifies whether the network is an external network or not.
	External bool `json:"router:external"`
}
