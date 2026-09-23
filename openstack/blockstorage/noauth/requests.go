package noauth

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
)

// EndpointOpts specifies a "noauth" Cinder Endpoint.
type EndpointOpts struct {
	// CinderEndpoint [required] is currently only used with "noauth" Cinder.
	// A cinder endpoint with auth_strategy=noauth is necessary, for example:
	// http://example.com:8776/v3.
	CinderEndpoint string
}

// NewClient prepares an unauthenticated ProviderClient instance.
func NewClient(options auth.NoAuthOpts) (*gophercloud.ProviderClient, error) {
	if options == nil {
		return nil, gophercloud.ErrMissingInput{Argument: "options"}
	}
	token, err := options.ToNoAuthToken()
	if err != nil {
		return nil, err
	}

	return &gophercloud.ProviderClient{TokenID: token}, nil
}

func initClientOpts(client *gophercloud.ProviderClient, eo EndpointOpts, clientType string) (*gophercloud.ServiceClient, error) {
	if eo.CinderEndpoint == "" {
		return nil, fmt.Errorf("CinderEndpoint is required")
	}

	userID, projectID, found := strings.Cut(client.TokenID, ":")
	if !found || userID == "" || projectID == "" || strings.Contains(projectID, ":") {
		return nil, fmt.Errorf("malformed noauth token")
	}

	endpoint := gophercloud.NormalizeURL(gophercloud.NormalizeURL(eo.CinderEndpoint) + projectID)
	return &gophercloud.ServiceClient{
		Endpoint:       endpoint,
		ProviderClient: client,
		Type:           clientType,
	}, nil
}

// NewBlockStorageNoAuthV2 creates a ServiceClient that may be used to access "noauth" v2 block storage service.
func NewBlockStorageNoAuthV2(client *gophercloud.ProviderClient, eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "block-storage")
}

// NewBlockStorageNoAuthV3 creates a ServiceClient that may be used to access "noauth" v3 block storage service.
func NewBlockStorageNoAuthV3(client *gophercloud.ProviderClient, eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "block-storage")
}
