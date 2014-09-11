package networks

// User-defined options sent to the API when creating or updating a network.
type NetworkOpts struct {
	// The administrative state of the network, which is up (true) or down (false).
	AdminStateUp bool `json:"admin_state_up"`
	// The network name (optional)
	Name string `json:"name"`
	// Indicates whether this network is shared across all tenants. By default,
	// only administrative users can change this value.
	Shared bool `json:"shared"`
	// Admin-only. The UUID of the tenant that will own the network. This tenant
	// can be different from the tenant that makes the create network request.
	// However, only administrative users can specify a tenant ID other than their
	// own. You cannot change this value through authorization policies.
	TenantID string `json:"tenant_id"`
}
