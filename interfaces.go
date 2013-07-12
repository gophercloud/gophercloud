package gophercloud

// AccessProvider instances encapsulate a Keystone authentication interface.
type AccessProvider interface {
	// FirstEndpointUrlByCriteria searches through the service catalog for the first
	// matching entry endpoint fulfilling the provided criteria.  If nothing found,
	// return "".  Otherwise, return either the public or internal URL for the
	// endpoint, depending on both its existence and the setting of the ApiCriteria.UrlChoice
	// field.
	FirstEndpointUrlByCriteria(ApiCriteria) string

	// TODO(sfalvo): get Token() to automatically renew the authentication token if it's near expiry.

	// AuthToken provides a copy of the current authentication token for the user's credentials.
	AuthToken() string
}

// CloudServersProvider instances encapsulate a Cloud Servers API, should one exist in the service catalog
// for your provider.
type CloudServersProvider interface {
  // Servers

	ListServers() ([]Server, error)
	ServerById(id string) (*Server, error)
	CreateServer(ns NewServer) (*NewServer, error)

  // Images

  ListImages() ([]Image, error)

  // Flavors

  ListFlavors() ([]Flavor, error)
}

