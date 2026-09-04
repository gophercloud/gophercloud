package clouds

import (
	"io"
)

type cloudOpts struct {
	cloudName              string
	locations              []string
	publicLocations        []string
	cloudsyamlReader       io.Reader
	secureyamlReader       io.Reader
	cloudsPublicyamlReader io.Reader

	authURL      string
	endpointType string
	region       string

	caCertPath     string
	clientCertPath string
	clientKeyPath  string
	insecure       *bool
}

// ParseOption one of parse configuration returned by With* modifier
type ParseOption = func(*cloudOpts)

// WithCloudName allows to override the environment variable `OS_CLOUD`.
func WithCloudName(osCloud string) ParseOption {
	return func(co *cloudOpts) {
		co.cloudName = osCloud
	}
}

// WithLocations is a functional option that sets the search locations for the
// clouds.yaml file (and its optional companion secure.yaml). Each location is
// a file path pointing to a possible `clouds.yaml`.
func WithLocations(locations ...string) ParseOption {
	return func(co *cloudOpts) {
		co.locations = locations
	}
}

// WithPublicLocations is a functional option that sets the search location for
// the clouds-public.yaml. Each location is a file path pointing to a possible
// `clouds-public.yaml`
func WithPublicLocations(locations ...string) ParseOption {
	return func(co *cloudOpts) {
		co.publicLocations = locations
	}
}

// WithCloudsYAML is a functional option that lets you pass a clouds.yaml file
// as an io.Reader interface. When this option is passed, FromCloudsYaml will
// not attempt to fetch any file from the file system. To add a secure.yaml,
// use in conjunction with WithSecureYAML.
func WithCloudsYAML(clouds io.Reader) ParseOption {
	return func(co *cloudOpts) {
		co.cloudsyamlReader = clouds
	}
}

// WithSecureYAML is a functional option that lets you pass a secure.yaml file
// as an io.Reader interface, to complement the clouds.yaml that is either
// fetched from the filesystem, or passed with WithCloudsYAML.
func WithSecureYAML(secure io.Reader) ParseOption {
	return func(co *cloudOpts) {
		co.secureyamlReader = secure
	}
}

// WithCloudsPublicYAML is a functional option that lets you pass
// clouds-public.yaml file as an io.Reader interface
func WithCloudsPublicYAML(public io.Reader) ParseOption {
	return func(co *cloudOpts) {
		co.cloudsPublicyamlReader = public
	}
}

func WithIdentityEndpoint(authURL string) ParseOption {
	return func(co *cloudOpts) {
		co.authURL = authURL
	}
}

// WithRegion allows to override the endpoint type set in clouds.yaml or in the
// environment variable `OS_INTERFACE`.
func WithEndpointType(endpointType string) ParseOption {
	return func(co *cloudOpts) {
		co.endpointType = endpointType
	}
}

// WithRegion allows to override the region set in clouds.yaml or in the
// environment variable `OS_REGION_NAME`
func WithRegion(region string) ParseOption {
	return func(co *cloudOpts) {
		co.region = region
	}
}

func WithCACertPath(caCertPath string) ParseOption {
	return func(co *cloudOpts) {
		co.caCertPath = caCertPath
	}
}

func WithClientCertPath(clientCertPath string) ParseOption {
	return func(co *cloudOpts) {
		co.clientCertPath = clientCertPath
	}
}

func WithClientKeyPath(clientKeyPath string) ParseOption {
	return func(co *cloudOpts) {
		co.clientKeyPath = clientKeyPath
	}
}

func WithInsecure(insecure bool) ParseOption {
	return func(co *cloudOpts) {
		co.insecure = &insecure
	}
}
