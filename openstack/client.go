package openstack

import (
	"fmt"
	"net/url"

	"github.com/gophercloud/gophercloud"
	tokens2 "github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/utils"
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

	endpoint = gophercloud.NormalizeURL(endpoint)
	base = gophercloud.NormalizeURL(base)

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
		{ID: v20, Priority: 20, Suffix: "/v2.0/"},
		{ID: v30, Priority: 30, Suffix: "/v3/"},
	}

	chosen, endpoint, err := utils.ChooseVersion(client, versions)
	if err != nil {
		return err
	}

	switch chosen.ID {
	case v20:
		return v2auth(client, endpoint, options, gophercloud.EndpointOpts{})
	case v30:
		return v3auth(client, endpoint, options, gophercloud.EndpointOpts{})
	default:
		// The switch statement must be out of date from the versions list.
		return fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
	}
}

// AuthenticateV2 explicitly authenticates against the identity v2 endpoint.
func AuthenticateV2(client *gophercloud.ProviderClient, options gophercloud.AuthOptions, eo gophercloud.EndpointOpts) error {
	return v2auth(client, "", options, eo)
}

func v2auth(client *gophercloud.ProviderClient, endpoint string, options gophercloud.AuthOptions, eo gophercloud.EndpointOpts) error {
	v2Client, err := NewIdentityV2(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v2Client.Endpoint = endpoint
	}

	v2Opts := tokens2.AuthOptions{
		IdentityEndpoint: options.IdentityEndpoint,
		Username:         options.Username,
		Password:         options.Password,
		TenantID:         options.TenantID,
		TenantName:       options.TenantName,
		AllowReauth:      options.AllowReauth,
		TokenID:          options.TokenID,
	}

	result := tokens2.Create(v2Client, v2Opts)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	if options.AllowReauth {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v2auth(client, endpoint, options, eo)
		}
	}
	client.TokenID = token.ID
	client.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
		return V2EndpointURL(catalog, opts)
	}

	return nil
}

// AuthenticateV3 explicitly authenticates against the identity v3 service.
func AuthenticateV3(client *gophercloud.ProviderClient, options gophercloud.AuthOptions, eo gophercloud.EndpointOpts) error {
	return v3auth(client, "", options, eo)
}

func v3auth(client *gophercloud.ProviderClient, endpoint string, options gophercloud.AuthOptions, eo gophercloud.EndpointOpts) error {
	// Override the generated service endpoint with the one returned by the version endpoint.
	v3Client, err := NewIdentityV3(client, eo)
	if err != nil {
		return err
	}

	if endpoint != "" {
		v3Client.Endpoint = endpoint
	}

	// copy the auth options to a local variable that we can change. `options`
	// needs to stay as-is for reauth purposes
	v3Options := options

	var scope *tokens3.Scope
	if options.TenantID != "" {
		scope = &tokens3.Scope{
			ProjectID: options.TenantID,
		}
		v3Options.TenantID = ""
		v3Options.TenantName = ""
	} else {
		if options.TenantName != "" {
			scope = &tokens3.Scope{
				ProjectName: options.TenantName,
				DomainID:    options.DomainID,
				DomainName:  options.DomainName,
			}
			v3Options.TenantName = ""
		}
	}

	v3Opts := tokens3.AuthOptions{
		IdentityEndpoint: options.IdentityEndpoint,
		Username:         options.Username,
		UserID:           options.UserID,
		Password:         options.Password,
		DomainID:         options.DomainID,
		DomainName:       options.DomainName,
		TenantID:         options.TenantID,
		TenantName:       options.TenantName,
		AllowReauth:      options.AllowReauth,
		TokenID:          options.TokenID,
	}

	result := tokens3.Create(v3Client, v3Opts, scope)

	token, err := result.ExtractToken()
	if err != nil {
		return err
	}

	catalog, err := result.ExtractServiceCatalog()
	if err != nil {
		return err
	}

	client.TokenID = token.ID

	if options.AllowReauth {
		client.ReauthFunc = func() error {
			client.TokenID = ""
			return v3auth(client, endpoint, options, eo)
		}
	}
	client.EndpointLocator = func(opts gophercloud.EndpointOpts) (string, error) {
		return V3EndpointURL(catalog, opts)
	}

	return nil
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the v2 identity service.
func NewIdentityV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	v2Endpoint := client.IdentityBase + "v2.0/"
	/*
		eo.ApplyDefaults("identity")
		url, err := client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	*/

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       v2Endpoint,
		//Endpoint: url,
	}, nil
}

// NewIdentityV3 creates a ServiceClient that may be used to access the v3 identity service.
func NewIdentityV3(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	v3Endpoint := client.IdentityBase + "v3/"
	/*
		eo.ApplyDefaults("identity")
		url, err := client.EndpointLocator(eo)
		if err != nil {
			return nil, err
		}
	*/

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       v3Endpoint,
		//Endpoint: url,
	}, nil
}

// NewObjectStorageV1 creates a ServiceClient that may be used with the v1 object storage package.
func NewObjectStorageV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("object-store")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}

// NewComputeV2 creates a ServiceClient that may be used with the v2 compute package.
func NewComputeV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("compute")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}

// NewNetworkV2 creates a ServiceClient that may be used with the v2 network package.
func NewNetworkV2(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("network")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       url,
		ResourceBase:   url + "v2.0/",
	}, nil
}

// NewBlockStorageV1 creates a ServiceClient that may be used to access the v1 block storage service.
func NewBlockStorageV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("volume")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}

// NewCDNV1 creates a ServiceClient that may be used to access the OpenStack v1
// CDN service.
func NewCDNV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("cdn")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}

// NewOrchestrationV1 creates a ServiceClient that may be used to access the v1 orchestration service.
func NewOrchestrationV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("orchestration")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}

// NewDBV1 creates a ServiceClient that may be used to access the v1 DB service.
func NewDBV1(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	eo.ApplyDefaults("database")
	url, err := client.EndpointLocator(eo)
	if err != nil {
		return nil, err
	}
	return &gophercloud.ServiceClient{ProviderClient: client, Endpoint: url}, nil
}
