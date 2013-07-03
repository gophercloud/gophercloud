package gophercloud

import "strings"

// globalContext is the, well, "global context."
// Most of this SDK is written in a manner to facilitate easier testing,
// which doesn't require all the configuration a real-world application would require.
// However, for real-world deployments, applications should be able to rely on a consistent configuration of providers, etc.
var globalContext *Context

// providers is the set of supported providers.
var providers = map[string]Provider{
	"rackspace-us": Provider{
		AuthEndpoint: "https://identity.api.rackspacecloud.com/v2.0/tokens",
	},
	"rackspace-uk": Provider{
		AuthEndpoint: "https://lon.identity.api.rackspacecloud.com/v2.0/tokens",
	},
}

// Initialize the global context to sane configuration.
// The Go runtime ensures this function is called before main(),
// thus guaranteeing proper configuration before your application ever runs.
func init() {
	globalContext = TestContext()
	for name, descriptor := range providers {
		globalContext.RegisterProvider(name, descriptor)
	}
}

// Authenticate() grants access to the OpenStack-compatible provider API.
//
// Providers are identified through a unique key string.
// Specifying an unsupported provider will result in an ErrProvider error.
//
// The supplied AuthOptions instance allows the client to specify only those credentials
// relevant for the authentication request.  At present, support exists for OpenStack
// Identity V2 API only; support for V3 will become available as soon as documentation for it
// becomes readily available.
//
// For Identity V2 API requirements, you must provide at least the Username and Password
// options.  The TenantId field is optional, and defaults to "".
func Authenticate(provider string, options AuthOptions) (*Access, error) {
	return globalContext.Authenticate(provider, options)
}

// ApiCriteria provides one or more criteria for the SDK to look for appropriate endpoints.
// Fields left unspecified or otherwise set to their zero-values are assumed to not be
// relevant, and do not participate in the endpoint search.
type ApiCriteria struct {
	// Name specifies the desired service catalog entry name.
	Name string

	// Region specifies the desired endpoint region.
	Region string

	// VersionId specifies the desired version of the endpoint.
	// Note that this field is matched exactly, and is (at present)
	// opaque to Gophercloud.  Thus, requesting a version 2
	// endpoint will _not_ match a version 3 endpoint.
	VersionId string

	// The UrlChoice field inidicates whether or not gophercloud
	// should use the public or internal endpoint URL if a
	// candidate endpoint is found.
	UrlChoice int
}

// The choices available for UrlChoice.  See the ApiCriteria structure for details.
const (
	PublicURL = iota
	InternalURL
)

// ComputeProvider instances encapsulate a Cloud Servers API, should one exist in the service catalog
// for your provider.
type ComputeProvider interface {
	ListServers() ([]Server, error)
}

// AccessProvider instances encapsulate a Keystone authentication interface.
type AccessProvider interface {
	// FirstEndpointUrlByCriteria searches through the service catalog for the first
	// matching entry endpoint fulfilling the provided criteria.  If nothing found,
	// return "".  Otherwise, return either the public or internal URL for the
	// endpoint, depending on both its existence and the setting of the ApiCriteria.UrlChoice
	// field.
	FirstEndpointUrlByCriteria(ApiCriteria) string
}

// genericCloudProvider structures provide the implementation for generic OpenStack-compatible
// ComputeProvider interfaces.
type genericCloudProvider struct {
	// endpoint refers to the provider's API endpoint base URL.  This will be used to construct
	// and issue queries.
	endpoint string
}

func ComputeApi(acc AccessProvider, criteria ApiCriteria) (ComputeProvider, error) {
	url := acc.FirstEndpointUrlByCriteria(criteria)
	if url == "" {
		return nil, ErrEndpoint
	}

	gcp := &genericCloudProvider{
		endpoint: url,
	}

	return gcp, nil
}

// See AccessProvider interface definition for details.
func (a *Access) FirstEndpointUrlByCriteria(ac ApiCriteria) string {
	ep := FindFirstEndpointByCriteria(a.ServiceCatalog, ac)
	urls := []string{ep.PublicURL, ep.InternalURL}
	return urls[ac.UrlChoice]
}

// Given a set of criteria to match on, locate the first candidate endpoint
// in the provided service catalog.
//
// If nothing found, the result will be a zero-valued EntryEndpoint (all URLs
// set to "").
func FindFirstEndpointByCriteria(entries []CatalogEntry, ac ApiCriteria) EntryEndpoint {
	rgn := strings.ToUpper(ac.Region)

	for _, entry := range entries {
		if (ac.Name != "") && (ac.Name != entry.Name) {
			continue
		}

		for _, endpoint := range entry.Endpoints {
			if (ac.Region != "") && (rgn != strings.ToUpper(endpoint.Region)) {
				continue
			}

			if (ac.VersionId != "") && (ac.VersionId != endpoint.VersionId) {
				continue
			}

			return endpoint
		}
	}
	return EntryEndpoint{}
}

// See the ComputeProvider interface for details.
func (gcp *genericCloudProvider) ListServers() ([]Server, error) {
	return nil, nil
}

// Server structures provide data about a server running in your provider's cloud.
type Server struct {
	Id string
}
