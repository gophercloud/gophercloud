package clouds

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud/v2/auth"
)

// Clouds represents a collection of Cloud entries in a clouds.yaml file.
// The format of clouds.yaml is documented at
// https://docs.openstack.org/os-client-config/latest/user/configuration.html.
type Clouds struct {
	Clouds map[string]Cloud `yaml:"clouds" json:"clouds"`
}

type PublicClouds struct {
	Clouds map[string]Cloud `yaml:"public-clouds" json:"public-clouds"`
}

// Cloud represents an entry in a clouds.yaml/public-clouds.yaml/secure.yaml file.
type Cloud struct {
	Cloud   string `yaml:"cloud,omitempty" json:"cloud,omitempty"`
	Profile string `yaml:"profile,omitempty" json:"profile,omitempty"`

	// Auth holds the auth: section of a clouds.yaml entry verbatim. Its
	// shape depends on AuthType; see Cloud.AuthOptions in auth.go for how
	// it's turned into a concrete auth mechanism.
	Auth map[string]any `yaml:"auth,omitempty" json:"auth,omitempty"`

	AuthType   auth.AuthType `yaml:"auth_type,omitempty" json:"auth_type,omitempty"`
	RegionName string        `yaml:"region_name,omitempty" json:"region_name,omitempty"`
	Regions    []Region      `yaml:"regions,omitempty" json:"regions,omitempty"`

	// EndpointType and Interface both specify whether to use the public, internal,
	// or admin interface of a service. They should be considered synonymous, but
	// EndpointType will take precedence when both are specified.
	EndpointType string `yaml:"endpoint_type,omitempty" json:"endpoint_type,omitempty"`
	Interface    string `yaml:"interface,omitempty" json:"interface,omitempty"`

	// API Version overrides.
	IdentityAPIVersion string `yaml:"identity_api_version,omitempty" json:"identity_api_version,omitempty"`
	VolumeAPIVersion   string `yaml:"volume_api_version,omitempty" json:"volume_api_version,omitempty"`

	// Verify whether or not SSL API requests should be verified.
	Verify *bool `yaml:"verify,omitempty" json:"verify,omitempty"`

	// CACertFile a path to a CA Cert bundle that can be used as part of
	// verifying SSL API requests.
	CACertFile string `yaml:"cacert,omitempty" json:"cacert,omitempty"`

	// ClientCertFile a path to a client certificate to use as part of the SSL
	// transaction.
	ClientCertFile string `yaml:"cert,omitempty" json:"cert,omitempty"`

	// ClientKeyFile a path to a client key to use as part of the SSL
	// transaction.
	ClientKeyFile string `yaml:"key,omitempty" json:"key,omitempty"`
}

// Region represents a region included as part of cloud in clouds.yaml
// According to Python-based openstacksdk, this can be either a struct (as defined)
// or a plain string. Custom unmarshallers handle both cases.
type Region struct {
	Name   string `yaml:"name,omitempty" json:"name,omitempty"`
	Values Cloud  `yaml:"values,omitempty" json:"values,omitempty"`
}

// UnmarshalJSON handles either a plain string acting as the Name property or
// a struct, mimicking the Python-based openstacksdk.
func (r *Region) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err == nil {
		r.Name = name
		return nil
	}

	type region Region
	var tmp region
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	r.Name = tmp.Name
	r.Values = tmp.Values

	return nil
}

// UnmarshalYAML handles either a plain string acting as the Name property or
// a struct, mimicking the Python-based openstacksdk.
func (r *Region) UnmarshalYAML(unmarshal func(any) error) error {
	var name string
	if err := unmarshal(&name); err == nil {
		r.Name = name
		return nil
	}

	type region Region
	var tmp region
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	r.Name = tmp.Name
	r.Values = tmp.Values

	return nil
}
