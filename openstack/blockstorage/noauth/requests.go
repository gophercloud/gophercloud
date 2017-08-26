package noauth

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

// EndpointOpts specifies "noauth" Endpoints - e.g. CinderEndpoint
type EndpointOpts struct {
	// CinderEndpoint [required] is currently only used w/ "noauth" Cinder.
	// A cinder endpoint w/ "auth_strategy=noauth" is necessary - e.g. http://cinder:8776/v2
	CinderEndpoint string
}

/*
NewClient prepares an unauthenticated ProviderClient instance.
Most users will probably prefer using the UnAuthenticatedClient function
instead.
*/
func NewClient(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	client := &gophercloud.ProviderClient{
		IdentityBase: options.TenantName,
		TokenID:      fmt.Sprintf("%s:%s", options.Username, options.TenantName),
	}

	return client, nil
}

// UnAuthenticatedClient allows for standalone "noauth" Cinder usage
func UnAuthenticatedClient(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	if len(options.Username) == 0 {
		options.Username = "admin"
	}
	if len(options.TenantName) == 0 {
		options.TenantName = "admin"
	}
	client, err := NewClient(options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func initClientOpts(client *gophercloud.ProviderClient, eo EndpointOpts, clientType string) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	if len(eo.CinderEndpoint) > 0 {
		sc.Endpoint = getURL(client, eo.CinderEndpoint)
	} else {
		return nil, fmt.Errorf("Pass proper EndPointOpt")
	}
	sc.ProviderClient = client
	sc.Type = clientType
	return sc, nil
}

// NewBlockStorageV2 creates a ServiceClient that may be used to access a v2 "noauth"
// block storage service.
func NewBlockStorageV2(client *gophercloud.ProviderClient, eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(client, eo, "volumev2")
}
