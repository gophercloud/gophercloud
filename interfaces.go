package gophercloud

// AccessProvider instances encapsulate a Keystone authentication interface.
type AccessProvider interface {
	// FirstEndpointUrlByCriteria searches through the service catalog for the first
	// matching entry endpoint fulfilling the provided criteria.  If nothing found,
	// return "".  Otherwise, return either the public or internal URL for the
	// endpoint, depending on both its existence and the setting of the ApiCriteria.UrlChoice
	// field.
	FirstEndpointUrlByCriteria(ApiCriteria) string
}

// ComputeProvider instances encapsulate a Cloud Servers API, should one exist in the service catalog
// for your provider.
type ComputeProvider interface {
	ListServers() ([]Server, error)
}
