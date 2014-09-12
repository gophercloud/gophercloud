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
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	hadPath := u.Path != ""
	u.Path, u.RawQuery, u.Fragment = "", "", ""
	base := u.String()

	endpoint = normalizeURL(endpoint)
	base = normalizeURL(base)

	if hadPath {
		return &gophercloud.ProviderClient{
			IdentityBase:     base,
			IdentityEndpoint: endpoint,
		}, nil
	}

	return &gophercloud.ProviderClient{
		IdentityBase:     base,
		IdentityEndpoint: "",
	}, nil
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
		&utils.Version{ID: v20, Priority: 20, Suffix: "/v2.0/"},
		&utils.Version{ID: v30, Priority: 30, Suffix: "/v3/"},
	}

	chosen, endpoint, err := utils.ChooseVersion(client.IdentityBase, client.IdentityEndpoint, versions)
	if err != nil {
		return err
	}

	switch chosen.ID {
	case v20:
		return v2auth(client, endpoint, options)
	case v30:
		return v3auth(client, endpoint, options)
	default:
		// The switch statement must be out of date from the versions list.
		return fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
	}
}

// AuthenticateV2 explicitly authenticates against the identity v2 endpoint.
func AuthenticateV2(client *gophercloud.ProviderClient, options gophercloud.AuthOptions) error {
	return v2auth(client, "", options)
}

func v2auth(client *gophercloud.ProviderClient, endpoint string, options gophercloud.AuthOptions) error {
	v2Client := NewIdentityV2(client)
	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

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
		return v2endpointLocator(result, opts)
	}

	return nil
}

func v2endpointLocator(authResults identity2.AuthResults, opts gophercloud.EndpointOpts) (string, error) {
	catalog, err := identity2.GetServiceCatalog(authResults)
	if err != nil {
		return "", err
	}

	entries, err := catalog.CatalogEntries()
	if err != nil {
		return "", err
	}

	// Extract Endpoints from the catalog entries that match the requested Type, Name if provided, and Region if provided.
	var endpoints = make([]identity2.Endpoint, 0, 1)
	for _, entry := range entries {
		if (entry.Type == opts.Type) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Region == "" || endpoint.Region == opts.Region {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if the options were ambiguous.
	if len(endpoints) == 0 {
		return "", gophercloud.ErrEndpointNotFound
	}
	if len(endpoints) > 1 {
		return "", fmt.Errorf("Discovered %d matching endpoints: %#v", len(endpoints), endpoints)
	}

	// Extract the appropriate URL from the matching Endpoint.
	for _, endpoint := range endpoints {
		switch opts.Availability {
		case gophercloud.AvailabilityPublic:
			return normalizeURL(endpoint.PublicURL), nil
		case gophercloud.AvailabilityInternal:
			return normalizeURL(endpoint.InternalURL), nil
		default:
			return "", fmt.Errorf("Unexpected availability in endpoint query: %s", opts.Availability)
		}
	}

	return "", gophercloud.ErrEndpointNotFound
}

// AuthenticateV3 explicitly authenticates against the identity v3 service.
func AuthenticateV3(client *gophercloud.ProviderClient, options gophercloud.AuthOptions) error {
	return v3auth(client, "", options)
}

func v3auth(client *gophercloud.ProviderClient, endpoint string, options gophercloud.AuthOptions) error {
	// Override the generated service endpoint with the one returned by the version endpoint.
	v3Client := NewIdentityV3(client)
	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

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
}

func v3endpointLocator(v3Client *gophercloud.ServiceClient, opts gophercloud.EndpointOpts) (string, error) {
	// Discover the service we're interested in.
	var services = make([]services3.Service, 0, 1)
	var err error
	services3.List(v3Client, services3.ListOpts{ServiceType: opts.Type}).EachPage(func(page gophercloud.Page) bool {
		part, err := services3.ExtractServices(page)
		if err != nil {
			return false
		}

		for _, service := range part {
			if service.Name == opts.Name {
				services = append(services, service)
			}
		}

		return true
	})
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", gophercloud.ErrServiceNotFound
	}
	if len(services) > 1 {
		return "", fmt.Errorf("Discovered %d matching services: %#v", len(services), services)
	}
	service := services[0]

	// Enumerate the endpoints available for this service.
	var endpoints []endpoints3.Endpoint
	endpoints3.List(v3Client, endpoints3.ListOpts{
		Availability: opts.Availability,
		ServiceID:    service.ID,
	}).EachPage(func(page gophercloud.Page) bool {
		part, err := endpoints3.ExtractEndpoints(page)
		if err != nil {
			return false
		}

		for _, endpoint := range part {
			if opts.Region == "" || endpoint.Region == opts.Region {
				endpoints = append(endpoints, endpoint)
			}
		}

		return true
	})
	if err != nil {
		return "", err
	}

	if len(endpoints) == 0 {
		return "", gophercloud.ErrEndpointNotFound
	}
	if len(endpoints) > 1 {
		return "", fmt.Errorf("Discovered %d matching endpoints: %#v", len(endpoints), endpoints)
	}
	endpoint := endpoints[0]

	return normalizeURL(endpoint.URL), nil
}

// normalizeURL ensures that each endpoint URL has a closing `/`, as expected by ServiceClient.
func normalizeURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	}
	return url
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the v2 identity service.
func NewIdentityV2(client *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	v2Endpoint := client.IdentityBase + "v2.0/"

	return &gophercloud.ServiceClient{
		Provider: client,
		Endpoint: v2Endpoint,
	}
}

// NewIdentityV3 creates a ServiceClient that may be used to access the v3 identity service.
func NewIdentityV3(client *gophercloud.ProviderClient) *gophercloud.ServiceClient {
	v3Endpoint := client.IdentityBase + "v3/"

	return &gophercloud.ServiceClient{
		Provider: client,
		Endpoint: v3Endpoint,
	}
}

// NewStorageV1 creates a ServiceClient that may be used with the v1 object storage package.
func NewStorageV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("object-store")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{Provider: client, Endpoint: url}, nil
}
