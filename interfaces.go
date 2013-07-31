package gophercloud

// AccessProvider instances encapsulate a Keystone authentication interface.
type AccessProvider interface {
	// FirstEndpointUrlByCriteria searches through the service catalog for the first
	// matching entry endpoint fulfilling the provided criteria.  If nothing found,
	// return "".  Otherwise, return either the public or internal URL for the
	// endpoint, depending on both its existence and the setting of the ApiCriteria.UrlChoice
	// field.
	FirstEndpointUrlByCriteria(ApiCriteria) string

	// AuthToken provides a copy of the current authentication token for the user's credentials.
	// Note that AuthToken() will not automatically refresh an expired token.
	AuthToken() string

	// Revoke allows you to terminate any program's access to the OpenStack API by token ID.
	Revoke(string) error

	// Reauthenticate attempts to acquire a new authentication token, if the feature is enabled by
	// AuthOptions.AllowReauth.
	Reauthenticate() error
}

// CloudServersProvider instances encapsulate a Cloud Servers API, should one exist in the service catalog
// for your provider.
type CloudServersProvider interface {
  // Servers

	ListServers() ([]Server, error)
	ListServersLinksOnly() ([]Server, error)
	ServerById(id string) (*Server, error)
	CreateServer(ns NewServer) (*NewServer, error)
	DeleteServerById(id string) error
	SetAdminPassword(id string, pw string) error
	ResizeServer(id, newName, newFlavor, newDiskConfig string) error
	RevertResize(id string) error
	ConfirmResize(id string) error
	RebootServer(id string, hard bool) error
	RescueServer(id string) (string, error)
	UnrescueServer(id string) error
	UpdateServer(id string, newValues NewServerSettings) (*Server, error)

  // Images

  ListImages() ([]Image, error)

  // Flavors

  ListFlavors() ([]Flavor, error)
}

