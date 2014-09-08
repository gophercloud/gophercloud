package openstack

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/rackspace/gophercloud"
	identity2 "github.com/rackspace/gophercloud/openstack/identity/v2"
	endpoints3 "github.com/rackspace/gophercloud/openstack/identity/v3/endpoints"
	services3 "github.com/rackspace/gophercloud/openstack/identity/v3/services"
	tokens3 "github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"github.com/rackspace/gophercloud/openstack/utils"
)

const (
	v20 = "v2.0"
	v30 = "v3.0"
)

// NewClient prepares an unauthenticated ProviderClient instance.
// Most users will probably prefer using the AuthenticatedClient function instead.
// This is useful if you wish to explicitly control the version of the identity service that's used for authentication explicitly,
// for example.
func NewClient(endpoint string) (*gophercloud.ProviderClient, error) {
	// Normalize the identity endpoint that's provided by trimming any path, query or fragment from the URL.
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	u.Path, u.RawQuery, u.Fragment = "", "", ""
	normalized := u.String()

	return &gophercloud.ProviderClient{IdentityEndpoint: normalized}, nil
}

// AuthenticatedClient logs in to an OpenStack cloud found at the identity endpoint specified by options, acquires a token, and
// returns a Client instance that's ready to operate.
// It first queries the root identity endpoint to determine which versions of the identity service are supported, then chooses
// the most recent identity service available to proceed.
func AuthenticatedClient(options gophercloud.AuthOptions) (*gophercloud.ProviderClient, error) {
	client, err := NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	err = Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authenticate or re-authenticate against the most recent identity service supported at the provided endpoint.
func Authenticate(client *gophercloud.ProviderClient, options gophercloud.AuthOptions) error {
	versions := []*utils.Version{
		&utils.Version{ID: v20, Priority: 20},
		&utils.Version{ID: v30, Priority: 30},
	}

	chosen, endpoint, err := utils.ChooseVersion(client.IdentityEndpoint, versions)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(endpoint, "/") {
		endpoint = endpoint + "/"
	}

	switch chosen.ID {
	case v20:
		v2Client := NewIdentityV2(client)
		v2Client.Endpoint = endpoint

		result, err := identity2.Authenticate(v2Client, options)
		if err != nil {
			return err
		}

		token, err := identity2.GetToken(result)
		if err != nil {
			return err
		}

		client.TokenID = token.ID
		client.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
			return v2endpointLocator(v2Client, opts)
		}

		return nil
	case v30:
		// Override the generated service endpoint with the one returned by the version endpoint.
		v3Client := NewIdentityV3(client)
		v3Client.Endpoint = endpoint

		result, err := tokens3.Create(v3Client, options, nil)
		if err != nil {
			return err
		}

		client.TokenID, err = result.TokenID()
		if err != nil {
			return err
		}
		client.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
			return v3endpointLocator(v3Client, opts)
		}

		return nil
	default:
		// The switch statement must be out of date from the versions list.
		return fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
	}
}

func v2endpointLocator(v2Client *gophercloud.ServiceClient, opts gophercloud.EndpointOpts) (string, error) {
	return "", gophercloud.ErrEndpointNotFound
}

func v3endpointLocator(v3Client *gophercloud.ServiceClient, opts gophercloud.EndpointOpts) (string, error) {
	// Transform URLType into an Interface.
	var endpointInterface = endpoints3.InterfacePublic
	switch opts.URLType {
	case "", "public":
		endpointInterface = endpoints3.InterfacePublic
	case "internal":
		endpointInterface = endpoints3.InterfaceInternal
	case "admin":
		endpointInterface = endpoints3.InterfaceAdmin
	default:
		return "", fmt.Errorf("Unrecognized URLType: %s", opts.URLType)
	}

	// Discover the service we're interested in.
	computeResults, err := services3.List(v3Client, services3.ListOpts{ServiceType: opts.Type})
	if err != nil {
		return "", err
	}

	serviceResults, err := gophercloud.AllPages(computeResults)
	if err != nil {
		return "", err
	}
	allServices := services3.AsServices(serviceResults)

	if opts.Name != "" {
		filtered := make([]services3.Service, 1)
		for _, service := range allServices {
			if service.Name == opts.Name {
				filtered = append(filtered, service)
			}
		}
		allServices = filtered
	}

	if len(allServices) == 0 {
		return "", gophercloud.ErrEndpointNotFound
	}
	if len(allServices) > 1 {
		return "", fmt.Errorf("Discovered %d matching services: %#v", len(allServices), allServices)
	}

	service := allServices[0]

	// Enumerate the endpoints available for this service.
	endpointResults, err := endpoints3.List(v3Client, endpoints3.ListOpts{
		Interface: endpointInterface,
		ServiceID: service.ID,
	})
	if err != nil {
		return "", err
	}

	allEndpoints, err := gophercloud.AllPages(endpointResults)
	if err != nil {
		return "", err
	}

	endpoints := endpoints3.AsEndpoints(allEndpoints)

	if opts.Name != "" {
		filtered := make([]endpoints3.Endpoint, 1)
		for _, endpoint := range endpoints {
			if endpoint.Region == opts.Region {
				filtered = append(filtered, endpoint)
			}
		}
		endpoints = filtered
	}

	if len(endpoints) == 0 {
		return "", gophercloud.ErrEndpointNotFound
	}
	if len(endpoints) > 1 {
		return "", fmt.Errorf("Discovered %d matching endpoints: %#v", len(endpoints), endpoints)
	}

	endpoint := endpoints[0]

	return endpoint.URL, nil
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the v2 identity service.
func NewIdentityV2(client *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	v2Endpoint := client.IdentityEndpoint + "/v2.0/"

	return &gophercloud.ServiceClient{
		Provider: client,
		Endpoint: v2Endpoint,
	}
}

// NewIdentityV3 creates a ServiceClient that may be used to access the v3 identity service.
func NewIdentityV3(client *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	v3Endpoint := client.IdentityEndpoint + "/v3/"

	return &gophercloud.ServiceClient{
		Provider: client,
		Endpoint: v3Endpoint,
	}
}
