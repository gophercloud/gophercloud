package clouds

import (
	"io"

	"github.com/gophercloud/gophercloud/v2"
)

type cloudOpts struct {
	cloudName        string
	locations        []string
	cloudsyamlReader io.Reader
	secureyamlReader io.Reader

	applicationCredentialID     string
	applicationCredentialName   string
	applicationCredentialSecret string
	authURL                     string
	domainID                    string
	domainName                  string
	endpointType                string
	password                    string
	projectID                   string
	projectName                 string
	region                      string
	scope                       *gophercloud.AuthScope
	token                       string
	userID                      string
	username                    string

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

func WithApplicationCredentialID(applicationCredentialID string) ParseOption {
	return func(co *cloudOpts) {
		co.applicationCredentialID = applicationCredentialID
	}
}

func WithApplicationCredentialName(applicationCredentialName string) ParseOption {
	return func(co *cloudOpts) {
		co.applicationCredentialName = applicationCredentialName
	}
}

func WithApplicationCredentialSecret(applicationCredentialSecret string) ParseOption {
	return func(co *cloudOpts) {
		co.applicationCredentialSecret = applicationCredentialSecret
	}
}

func WithIdentityEndpoint(authURL string) ParseOption {
	return func(co *cloudOpts) {
		co.authURL = authURL
	}
}

func WithDomainID(domainID string) ParseOption {
	return func(co *cloudOpts) {
		co.domainID = domainID
	}
}

func WithDomainName(domainName string) ParseOption {
	return func(co *cloudOpts) {
		co.domainName = domainName
	}
}

// WithRegion allows to override the endpoint type set in clouds.yaml or in the
// environment variable `OS_INTERFACE`.
func WithEndpointType(endpointType string) ParseOption {
	return func(co *cloudOpts) {
		co.endpointType = endpointType
	}
}

func WithPassword(password string) ParseOption {
	return func(co *cloudOpts) {
		co.password = password
	}
}

func WithProjectID(projectID string) ParseOption {
	return func(co *cloudOpts) {
		co.projectID = projectID
	}
}

func WithProjectName(projectName string) ParseOption {
	return func(co *cloudOpts) {
		co.projectName = projectName
	}
}

// WithRegion allows to override the region set in clouds.yaml or in the
// environment variable `OS_REGION_NAME`
func WithRegion(region string) ParseOption {
	return func(co *cloudOpts) {
		co.region = region
	}
}

func WithScope(scope *gophercloud.AuthScope) ParseOption {
	return func(co *cloudOpts) {
		co.scope = scope
	}
}

func WithToken(token string) ParseOption {
	return func(co *cloudOpts) {
		co.token = token
	}
}

func WithUserID(userID string) ParseOption {
	return func(co *cloudOpts) {
		co.userID = userID
	}
}

func WithUsername(username string) ParseOption {
	return func(co *cloudOpts) {
		co.username = username
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
