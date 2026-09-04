package openstack

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
)

// NewClient prepares an unauthenticated ProviderClient instance.
// Most users will probably prefer using the AuthenticatedClient function
// instead.
//
// This is useful if you wish to explicitly control the version of the identity
// service that's used for authentication explicitly, for example.
//
// A basic example of using this would be:
//
//	ao, err := openstack.AuthOptionsFromEnv()
//	provider, err := openstack.NewClient(ao.IdentityEndpoint)
//	client, err := openstack.NewIdentityV3(ctx, provider, gophercloud.EndpointOpts{})
func NewClient(endpoint string) (*gophercloud.ProviderClient, error) {
	base, err := utils.BaseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	endpoint = gophercloud.NormalizeURL(endpoint)
	base = gophercloud.NormalizeURL(base)

	p := new(gophercloud.ProviderClient)
	p.IdentityBase = base
	p.IdentityEndpoint = endpoint
	p.UseTokenLock()

	return p, nil
}

// AuthenticatedClient logs in to an OpenStack cloud found at the identity endpoint
// specified by the options, acquires a token, and returns a Provider Client
// instance that's ready to operate.
//
// If the full path to a versioned identity endpoint was specified  (example:
// http://example.com:5000/v3), that path will be used as the endpoint to query.
//
// If a versionless endpoint was specified (example: http://example.com:5000/),
// the endpoint will be queried to determine which versions of the identity service
// are available, then chooses the most recent or most supported version.
//
// Example:
//
//	ao, err := auth.AuthOptionsFromEnv()
//	provider, err := openstack.AuthenticatedClient(ctx, ao)
//	client, err := openstack.NewNetworkV2(ctx, provider, gophercloud.EndpointOpts{
//		Region: os.Getenv("OS_REGION_NAME"),
//	})
func AuthenticatedClient(ctx context.Context, options auth.Authenticator) (*gophercloud.ProviderClient, error) {
	client, err := NewClient(options.GetAuthURL())
	if err != nil {
		return nil, err
	}

	err = Authenticate(ctx, client, options)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Authenticate authenticates or re-authenticates against the most
// recent identity service supported at the provided endpoint.
func Authenticate(ctx context.Context, client *gophercloud.ProviderClient, options auth.Authenticator) error {
	result, err := options.Authenticate(ctx, &client.HTTPClient)
	if err != nil {
		return err
	}
	if err := client.SetTokenAndAuthResult(result); err != nil {
		return err
	}
	client.EndpointLocator = result.EndpointLocator()

	if result.CanReauth {
		client.ReauthFunc = func(ctx context.Context) error {
			result, err := options.Authenticate(ctx, &client.HTTPClient)
			if err != nil {
				return err
			}
			if err := client.SetTokenAndAuthResult(result); err != nil {
				return err
			}
			client.EndpointLocator = result.EndpointLocator()
			return nil
		}
	}
	return nil
}

// NewIdentityV2 creates a ServiceClient that may be used to interact with the
// v2 identity service.
func NewIdentityV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	endpoint := client.IdentityBase + "v2.0/"
	clientType := "identity"
	var err error
	if !reflect.DeepEqual(eo, gophercloud.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(ctx, eo)
		if err != nil {
			return nil, err
		}
	}

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

// NewIdentityV3 creates a ServiceClient that may be used to access the v3
// identity service.
func NewIdentityV3(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	endpoint := client.IdentityBase + "v3/"
	clientType := "identity"
	var err error
	if !reflect.DeepEqual(eo, gophercloud.EndpointOpts{}) {
		eo.ApplyDefaults(clientType)
		endpoint, err = client.EndpointLocator(ctx, eo)
		if err != nil {
			return nil, err
		}
	}

	// Ensure endpoint still has a suffix of v3.
	// This is because EndpointLocator might have found a versionless
	// endpoint or the published endpoint is still /v2.0. In both
	// cases, we need to fix the endpoint to point to /v3.
	base, err := utils.BaseEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	base = gophercloud.NormalizeURL(base)

	endpoint = base + "v3/"

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       endpoint,
		Type:           clientType,
	}, nil
}

// TODO(stephenfin): Allow passing aliases to all New${SERVICE}V${VERSION} methods in v3
func initClientOpts(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts, clientType string, version int) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)

	eo.ApplyDefaults(clientType)
	if eo.Version != 0 && eo.Version != version {
		return sc, errors.New("conflict between requested service major version and manually set version")
	}
	eo.Version = version

	url, err := client.EndpointLocator(ctx, eo)
	if err != nil {
		return sc, err
	}

	sc.ProviderClient = client
	sc.Endpoint = url
	sc.Type = clientType
	return sc, nil
}

// NewBareMetalV1 creates a ServiceClient that may be used with the v1
// bare metal package.
func NewBareMetalV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "baremetal", 1)
	if err != nil {
		return sc, err
	}
	if !strings.HasSuffix(strings.TrimSuffix(sc.Endpoint, "/"), "v1") {
		sc.ResourceBase = sc.Endpoint + "v1/"
	}
	return sc, nil
}

// NewBareMetalIntrospectionV1 creates a ServiceClient that may be used with the v1
// bare metal introspection package.
func NewBareMetalIntrospectionV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "baremetal-introspection", 1)
}

// NewObjectStorageV1 creates a ServiceClient that may be used with the v1
// object storage package.
func NewObjectStorageV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "object-store", 1)
}

// NewComputeV2 creates a ServiceClient that may be used with the v2 compute
// package.
func NewComputeV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "compute", 2)
}

// NewNetworkV2 creates a ServiceClient that may be used with the v2 network
// package.
func NewNetworkV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "network", 2)
	if err != nil {
		return sc, err
	}
	sc.ResourceBase = sc.Endpoint + "v2.0/"
	return sc, nil
}

// TODO(stephenfin): Remove this in v3. We no longer support the V1 Block Storage service.
// NewBlockStorageV1 creates a ServiceClient that may be used to access the v1
// block storage service.
func NewBlockStorageV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "volume", 1)
}

// NewBlockStorageV2 creates a ServiceClient that may be used to access the v2
// block storage service.
func NewBlockStorageV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "block-storage", 2)
}

// NewBlockStorageV3 creates a ServiceClient that may be used to access the v3 block storage service.
func NewBlockStorageV3(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "block-storage", 3)
}

// NewSharedFileSystemV2 creates a ServiceClient that may be used to access the v2 shared file system service.
func NewSharedFileSystemV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "shared-file-system", 2)
}

// NewOrchestrationV1 creates a ServiceClient that may be used to access the v1
// orchestration service.
func NewOrchestrationV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "orchestration", 1)
}

// NewDBV1 creates a ServiceClient that may be used to access the v1 DB service.
func NewDBV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "database", 1)
}

// NewDNSV2 creates a ServiceClient that may be used to access the v2 DNS
// service.
func NewDNSV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "dns", 2)
	if err != nil {
		return sc, err
	}
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, nil
}

// NewImageV2 creates a ServiceClient that may be used to access the v2 image
// service.
func NewImageV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "image", 2)
	if err != nil {
		return sc, err
	}
	sc.ResourceBase = sc.Endpoint + "v2/"
	return sc, nil
}

// NewLoadBalancerV2 creates a ServiceClient that may be used to access the v2
// load balancer service.
func NewLoadBalancerV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "load-balancer", 2)
	if err != nil {
		return sc, err
	}

	// Fixes edge case having an OpenStack lb endpoint with trailing version number.
	endpoint := strings.ReplaceAll(sc.Endpoint, "v2.0/", "")

	sc.ResourceBase = endpoint + "v2.0/"
	return sc, nil
}

// NewMetricV1 creates a ServiceClient that may be used with the v1 metric-storage
// service (Aetos).
func NewMetricV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "metric-storage", 1)
	if err != nil {
		return sc, err
	}
	sc.ResourceBase = sc.Endpoint + "api/v1/"
	return sc, nil
}

// NewMessagingV2 creates a ServiceClient that may be used with the v2 messaging
// service.
func NewMessagingV2(ctx context.Context, client *gophercloud.ProviderClient, clientID string, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "message", 2)
	if err != nil {
		return sc, err
	}
	sc.MoreHeaders = map[string]string{"Client-ID": clientID}
	return sc, nil
}

// NewContainerV1 creates a ServiceClient that may be used with v1 container package
func NewContainerV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "application-container", 1)
}

// NewKeyManagerV1 creates a ServiceClient that may be used with the v1 key
// manager service.
func NewKeyManagerV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "key-manager", 1)
	if err != nil {
		return sc, err
	}
	sc.ResourceBase = sc.Endpoint + "v1/"
	return sc, nil
}

// NewContainerInfraV1 creates a ServiceClient that may be used with the v1 container infra management
// package.
func NewContainerInfraV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "container-infrastructure-management", 1)
}

// NewWorkflowV2 creates a ServiceClient that may be used with the v2 workflow management package.
func NewWorkflowV2(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "workflow", 2)
}

// NewPlacementV1 creates a ServiceClient that may be used with the placement package.
func NewPlacementV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(ctx, client, eo, "placement", 1)
}

// NewReservationV1 creates a ServiceClient that may be used with the reservation package.
func NewReservationV1(ctx context.Context, client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(ctx, client, eo, "reservation", 1)
	if err != nil {
		return sc, err
	}
	if !strings.HasSuffix(strings.TrimSuffix(sc.Endpoint, "/"), "v1") {
		sc.ResourceBase = sc.Endpoint + "v1/"
	}
	return sc, nil
}
