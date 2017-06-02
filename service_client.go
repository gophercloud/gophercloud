package gophercloud

import (
	"fmt"
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
}

// OpenStackVersion is the OpenStack version
type openStackVersion string

const (
	// OpenStackV20 is the OpenStack v2.0 version
	openStackV20 openStackVersion = "v2.0"
	// OpenStackV30 is the OpenStack v3.0 version
	openStackV30 openStackVersion = "v3.0"
)

const (
	// OpenStackNovaV2ServiceType is the OpenStack NovaV2 Type value
	OpenStackNovaV2ServiceType = "compute"
	// OpenStackManilaV2ServiceType is the OpenStack ManilaV2 Type valuse
	OpenStackManilaV2ServiceType = "sharev2"
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

// InitializeV2 initializes an EndpointExtraInfo instance with OpenStack v2.0 endpoint data
func (e *EndpointExtraInfo) InitializeV2(serviceType, serviceName, adminURL, internalURL, publicURL, region, tenantID, versionID, versionInfo, versionList string) {
	e.version = openStackV20
	e.serviceType = serviceType
	e.serviceName = serviceName
	e.v2Endpoint.adminURL = adminURL
	e.v2Endpoint.internalURL = internalURL
	e.v2Endpoint.publicURL = publicURL
	e.v2Endpoint.region = region
	e.v2Endpoint.tenantID = tenantID
	e.v2Endpoint.versionID = versionID
	e.v2Endpoint.versionInfo = versionInfo
	e.v2Endpoint.versionList = versionList
	e.status = serviceTypeSet
	return
}

// InitializeV3 initializes an EndpointExtraInfo instance with OpenStack v3.0 endpoint data
func (e *EndpointExtraInfo) InitializeV3(serviceType, serviceName, id, endpointInterface, region, unnormalisedURL string) {
	e.version = openStackV30
	e.serviceType = serviceType
	e.serviceName = serviceName
	e.v3Endpoint.id = id
	e.v3Endpoint.endpointInterface = endpointInterface
	e.v3Endpoint.region = region
	e.v3Endpoint.unnormalisedURL = unnormalisedURL
	e.status = serviceTypeSet
	return
}

// IsMicroversionImplemented returns true in case microversion is implemented for the e.serviceType
func (e *EndpointExtraInfo) IsMicroversionImplemented() bool {
	if e.serviceType == OpenStackManilaV2ServiceType {
		return true
	}
	return false
}

func (e *EndpointExtraInfo) validateForMicroversion() error {
	if e.status == uninitialised {
		err := BaseError{}
		err.Info = fmt.Sprintf("Internal error: cannot set microversion, EndpointExtraInfo struct is not initialised")
		return err
	}
	if e.version != openStackV30 {
		err := BaseError{}
		err.Info = fmt.Sprintf("Internal error: cannot set microversion, microversion support is not implemented for (%v)", e.version)
		return err
	}
	if !e.IsMicroversionImplemented() {
		err := BaseError{}
		err.Info = fmt.Sprintf("Internal error: cannot set microversion, service type %q does not have microversions implemented", e.serviceType)
		return err
	}
	return nil
}

// SetMicroversion sets value of microversion for service types that have microversion implemented
func (e *EndpointExtraInfo) SetMicroversion(microversion string) error {
	if err := e.validateForMicroversion(); err != nil {
		return err
	}
	var err error
	e.v3EndpointExtraInfo.microversion, err = semver.NewVersion(microversion)
	if err != nil {
		retErr := BaseError{}
		retErr.Info = fmt.Sprintf("The microversion %q has invalid format: %q", microversion, err.Error())
		return retErr
	}
	e.status = microversionSet
	return nil
}

// GetMicroversion returns value of microversion for service types that have microversion implemented
func (e *EndpointExtraInfo) GetMicroversion(microversion string) (*semver.Version, error) {
	if err := e.validateForMicroversion(); err != nil {
		return nil, err
	}
	return e.v3EndpointExtraInfo.microversion, nil
}

func (e *EndpointExtraInfo) getHTTPAPIVersionHeader() (string, error) {
	switch e.serviceType {
	case OpenStackManilaV2ServiceType:
		return "X-Openstack-Manila-Api-Version", nil
	default:
		err := BaseError{}
		err.Info = fmt.Sprintf("Internal error: the (%v) service type doesn't have microversions implemented", e.serviceType)
		return "", err
	}
}

// GetAPIVersionForHTTPHeader returns:
// - in case the EndpointExtraInfo contains valid data: (HTTPMicroversionHeader, Value (for the HTTPMicroversionHeader), nil)
// - in case the EndpointExtraInfo contains invalid data or an error occured: (N/A, N/A, error)
func (e *EndpointExtraInfo) GetAPIVersionForHTTPHeader() (string, string, error) {
	if err := e.validateForMicroversion(); err != nil {
		return "", "", err
	}
	header, err := e.getHTTPAPIVersionHeader()
	if err != nil {
		return "", "", err
	}
	value := e.v3EndpointExtraInfo.microversion.String()
	return header, value, nil
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
	if header, value, err := client.EndpointExtraInfo.GetAPIVersionForHTTPHeader(); err == nil {
		opts.MoreHeaders[header] = value
	}

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
	if header, value, err := client.EndpointExtraInfo.GetAPIVersionForHTTPHeader(); err == nil {
		opts.MoreHeaders[header] = value
	}

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
	if header, value, err := client.EndpointExtraInfo.GetAPIVersionForHTTPHeader(); err == nil {
		opts.MoreHeaders[header] = value
	}

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
	if header, value, err := client.EndpointExtraInfo.GetAPIVersionForHTTPHeader(); err == nil {
		opts.MoreHeaders[header] = value
	}

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
	if header, value, err := client.EndpointExtraInfo.GetAPIVersionForHTTPHeader(); err == nil {
		opts.MoreHeaders[header] = value
	}

	return client.Request("DELETE", url, opts)
}
