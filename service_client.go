package gophercloud

import (
	"io"
	"net/http"
	"strings"

	"github.com/Masterminds/semver"
)

// ServiceClient stores details required to interact with a specific service API implemented by a provider.
// Generally, you'll acquire these by calling the appropriate `New` method on a ProviderClient.
type ServiceClient struct {
	// ProviderClient is a reference to the provider that implements this service.
	*ProviderClient

	// Endpoint is the base URL of the service's API, acquired from a service catalog.
	// It MUST end with a /.
	Endpoint string

	// Endpoint data except URL, i.e. service type, service name, microversion, etc.
	EndpointExtraInfo EndpointExtraInfo

	// ResourceBase is the base URL shared by the resources within a service's API. It should include
	// the API version and, like Endpoint, MUST end with a / if set. If not set, the Endpoint is used
	// as-is, instead.
	ResourceBase string

	Microversion string
}

// OpenStackVersion is the OpenStack version
type openStackVersion string

const (
	// OpenStackV20 is the OpenStack v2.0 version
	openStackV20 openStackVersion = "v2.0"
	// OpenStackV30 is the OpenStack v3.0 version
	openStackV30 openStackVersion = "v3.0"
)

type openStackV2Endpoint struct {
	tenantID    string
	publicURL   string
	internalURL string
	adminURL    string
	region      string
	versionID   string
	versionInfo string
	versionList string
}

type openStackV3Endpoint struct {
	id                string
	region            string
	endpointInterface string
	unnormalisedURL   string
}

type openStackV3EndpointExtraInfo struct {
	microversion *semver.Version
}

type endpointExtraInfoStatus int

const (
	uninitialised endpointExtraInfoStatus = iota
	serviceTypeSet
	// microversionSet status is only for OpenStackV30 and for OpenStack service type that support microversions
	microversionSet
)

// EndpointExtraInfo stores endpoint data except URL, i.e. service type, service name, microversion, etc.
type EndpointExtraInfo struct {
	status              endpointExtraInfoStatus
	serviceType         string
	serviceName         string
	version             openStackVersion
	v2Endpoint          openStackV2Endpoint
	v3Endpoint          openStackV3Endpoint
	v3EndpointExtraInfo openStackV3EndpointExtraInfo
}

// ResourceBaseURL returns the base URL of any resources used by this service. It MUST end with a /.
func (client *ServiceClient) ResourceBaseURL() string {
	if client.ResourceBase != "" {
		return client.ResourceBase
	}
	return client.Endpoint
}

// ServiceURL constructs a URL for a resource belonging to this provider.
func (client *ServiceClient) ServiceURL(parts ...string) string {
	return client.ResourceBaseURL() + strings.Join(parts, "/")
}

// Get calls `Request` with the "GET" HTTP verb.
func (client *ServiceClient) Get(url string, JSONResponse interface{}, opts *RequestOpts) (*http.Response, error) {
	if opts == nil {
		opts = &RequestOpts{}
	}
	if JSONResponse != nil {
		opts.JSONResponse = JSONResponse
	}

	if opts.MoreHeaders == nil {
		opts.MoreHeaders = make(map[string]string)
	}
	opts.MoreHeaders["X-OpenStack-Nova-API-Version"] = client.Microversion

	return client.Request("GET", url, opts)
}

// Post calls `Request` with the "POST" HTTP verb.
func (client *ServiceClient) Post(url string, JSONBody interface{}, JSONResponse interface{}, opts *RequestOpts) (*http.Response, error) {
	if opts == nil {
		opts = &RequestOpts{}
	}

	if v, ok := (JSONBody).(io.Reader); ok {
		opts.RawBody = v
	} else if JSONBody != nil {
		opts.JSONBody = JSONBody
	}

	if JSONResponse != nil {
		opts.JSONResponse = JSONResponse
	}

	if opts.MoreHeaders == nil {
		opts.MoreHeaders = make(map[string]string)
	}
	opts.MoreHeaders["X-OpenStack-Nova-API-Version"] = client.Microversion

	return client.Request("POST", url, opts)
}

// Put calls `Request` with the "PUT" HTTP verb.
func (client *ServiceClient) Put(url string, JSONBody interface{}, JSONResponse interface{}, opts *RequestOpts) (*http.Response, error) {
	if opts == nil {
		opts = &RequestOpts{}
	}

	if v, ok := (JSONBody).(io.Reader); ok {
		opts.RawBody = v
	} else if JSONBody != nil {
		opts.JSONBody = JSONBody
	}

	if JSONResponse != nil {
		opts.JSONResponse = JSONResponse
	}

	if opts.MoreHeaders == nil {
		opts.MoreHeaders = make(map[string]string)
	}
	opts.MoreHeaders["X-OpenStack-Nova-API-Version"] = client.Microversion

	return client.Request("PUT", url, opts)
}

// Patch calls `Request` with the "PATCH" HTTP verb.
func (client *ServiceClient) Patch(url string, JSONBody interface{}, JSONResponse interface{}, opts *RequestOpts) (*http.Response, error) {
	if opts == nil {
		opts = &RequestOpts{}
	}

	if v, ok := (JSONBody).(io.Reader); ok {
		opts.RawBody = v
	} else if JSONBody != nil {
		opts.JSONBody = JSONBody
	}

	if JSONResponse != nil {
		opts.JSONResponse = JSONResponse
	}

	if opts.MoreHeaders == nil {
		opts.MoreHeaders = make(map[string]string)
	}
	opts.MoreHeaders["X-OpenStack-Nova-API-Version"] = client.Microversion

	return client.Request("PATCH", url, opts)
}

// Delete calls `Request` with the "DELETE" HTTP verb.
func (client *ServiceClient) Delete(url string, opts *RequestOpts) (*http.Response, error) {
	if opts == nil {
		opts = &RequestOpts{}
	}

	if opts.MoreHeaders == nil {
		opts.MoreHeaders = make(map[string]string)
	}
	opts.MoreHeaders["X-OpenStack-Nova-API-Version"] = client.Microversion

	return client.Request("DELETE", url, opts)
}
