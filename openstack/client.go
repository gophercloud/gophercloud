package openstack

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/rackspace/gophercloud"
	identity2 "github.com/rackspace/gophercloud/openstack/identity/v2"
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

	return &gophercloud.ProviderClient{
		IdentityEndpoint: normalized,
		Reauthenticate: func() error {
			return errors.New("Unable to reauthenticate before authenticating the first time.")
		},
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
		fmt.Printf("Endpoint is: %s\n", endpoint)

		result, err := identity2.Authenticate(v2Client, options)
		if err != nil {
			return err
		}

		token, err := identity2.GetToken(result)
		if err != nil {
			return err
		}
		client.TokenID = token.ID

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

		return nil
	default:
		// The switch statement must be out of date from the versions list.
		return fmt.Errorf("Unrecognized identity version: %s", chosen.ID)
	}
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
