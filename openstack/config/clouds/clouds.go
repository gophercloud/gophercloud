// package clouds provides a parser for OpenStack credentials stored in a clouds.yaml file.
//
// Example use:
//
//	ctx := context.Background()
//	ao, eo, tlsConfig, err := clouds.Parse()
//	if err != nil {
//		panic(err)
//	}
//
//	providerClient, err := config.NewProviderClient(ctx, ao, config.WithTLSConfig(tlsConfig))
//	if err != nil {
//		panic(err)
//	}
//
//	networkClient, err := openstack.NewNetworkV2(ctx, providerClient, eo)
//	if err != nil {
//		panic(err)
//	}
package clouds

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"

	"github.com/gophercloud/gophercloud/v2"
	"go.yaml.in/yaml/v3"
)

// Parse fetches a clouds.yaml file from disk and returns the parsed
// credentials.
//
// By default this function mimics the behaviour of python-openstackclient, which is:
//
//   - if the environment variable `OS_CLIENT_CONFIG_FILE` is set and points to a
//     clouds.yaml, use that location as the only search location for `clouds.yaml` and `secure.yaml`;
//   - otherwise, the search locations for `clouds.yaml` and `secure.yaml` are:
//     1. the current working directory (on Linux: `./`)
//     2. the directory `openstack` under the standard user config location for
//     the operating system (on Linux: `${XDG_CONFIG_HOME:-$HOME/.config}/openstack/`)
//     3. on Linux, `/etc/openstack/`
//
// Once `clouds.yaml` is found in a search location, the same location is used to search for `secure.yaml`.
//
// Like in python-openstackclient, relative paths in the `clouds.yaml` section
// `cacert` are interpreted as relative the the current directory, and not to
// the `clouds.yaml` location.
//
// Search locations, as well as individual `clouds.yaml` properties, can be
// overwritten with functional options.
func Parse(opts ...ParseOption) (gophercloud.AuthOptions, gophercloud.EndpointOpts, *tls.Config, error) {
	options := cloudOpts{
		cloudName:    os.Getenv("OS_CLOUD"),
		region:       os.Getenv("OS_REGION_NAME"),
		endpointType: os.Getenv("OS_INTERFACE"),
		locations: func() map[CloudsType][]string {
			out := make(map[CloudsType][]string)
			if path := os.Getenv("OS_CLIENT_CONFIG_FILE"); path != "" {
				out[Default] = []string{path}
			}
			return out
		}(),
		readers: make(map[CloudsType]io.Reader),
	}

	for _, apply := range opts {
		apply(&options)
	}

	if options.cloudName == "" {
		return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("the empty string \"\" is not a valid cloud name")
	}

	clouds, err := readClouds(&options, Default)
	if err != nil {
		return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, err
	}

	cloud, ok := clouds[options.cloudName]
	if !ok {
		return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("cloud %q not found in clouds.yaml", options.cloudName)
	}

	secureClouds, err := readClouds(&options, Secure)
	if err != nil {
		// For secure.yaml we ignore if there is not secure.yaml file found
		if _, ok := err.(ErrFileNotFound); !ok {
			return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, err
		}

	}
	secure, ok := secureClouds[options.cloudName]
	if ok {
		if !reflect.DeepEqual(cloud, secure) {
			cloud, err = mergeClouds(secure, cloud)
			if err != nil {
				return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("unable to merge information from clouds.yaml and secure.yaml: %w", err)
			}
		}
	}

	var profile string
	if cloud.Profile != "" {
		profile = cloud.Profile
	} else if cloud.Cloud != "" {
		profile = cloud.Cloud
	}

	if profile != "" {
		publicsCloud, err := readClouds(&options, Public)
		if err != nil {
			return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, err
		}

		public, ok := publicsCloud[profile]
		if ok {
			cloud, err = mergeClouds(cloud, public)
			if err != nil {
				return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, err
			}
		}
	}

	tlsConfig, err := computeTLSConfig(cloud, options)
	if err != nil {
		return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("unable to compute TLS configuration: %w", err)
	}

	endpointType := coalesce(options.endpointType, cloud.EndpointType, cloud.Interface)

	var scope *gophercloud.AuthScope
	if trustID := cloud.AuthInfo.TrustID; trustID != "" {
		scope = &gophercloud.AuthScope{
			TrustID: trustID,
		}
	}

	return gophercloud.AuthOptions{
			IdentityEndpoint:            coalesce(options.authURL, cloud.AuthInfo.AuthURL),
			Username:                    coalesce(options.username, cloud.AuthInfo.Username),
			UserID:                      coalesce(options.userID, cloud.AuthInfo.UserID),
			Password:                    coalesce(options.password, cloud.AuthInfo.Password),
			DomainID:                    coalesce(options.domainID, cloud.AuthInfo.UserDomainID, cloud.AuthInfo.ProjectDomainID, cloud.AuthInfo.DomainID),
			DomainName:                  coalesce(options.domainName, cloud.AuthInfo.UserDomainName, cloud.AuthInfo.ProjectDomainName, cloud.AuthInfo.DomainName),
			TenantID:                    coalesce(options.projectID, cloud.AuthInfo.ProjectID),
			TenantName:                  coalesce(options.projectName, cloud.AuthInfo.ProjectName),
			TokenID:                     coalesce(options.token, cloud.AuthInfo.Token),
			Scope:                       coalesce(options.scope, scope),
			ApplicationCredentialID:     coalesce(options.applicationCredentialID, cloud.AuthInfo.ApplicationCredentialID),
			ApplicationCredentialName:   coalesce(options.applicationCredentialName, cloud.AuthInfo.ApplicationCredentialName),
			ApplicationCredentialSecret: coalesce(options.applicationCredentialSecret, cloud.AuthInfo.ApplicationCredentialSecret),
		}, gophercloud.EndpointOpts{
			Region:       coalesce(options.region, cloud.RegionName),
			Availability: computeAvailability(endpointType),
		},
		tlsConfig,
		nil
}

func readClouds(options *cloudOpts, ctype CloudsType) (map[string]Cloud, error) {
	// Set the defaults and open the files for reading. This code only runs
	// if no override has been set, because it is fallible.
	if options.readers[ctype] == nil {
		if len(options.locations[ctype]) < 1 {
			cwd, err := os.Getwd()
			if err != nil {
				return map[string]Cloud{}, fmt.Errorf("failed to get the current working directory: %w", err)
			}
			userConfig, err := getUserConfig()
			if err != nil {
				return map[string]Cloud{}, err
			}
			options.locations[ctype] = []string{path.Join(cwd, string(ctype)), path.Join(userConfig, "openstack", string(ctype)), path.Join("/etc", "openstack", string(ctype))}
		}

		for _, cloudsPath := range options.locations[ctype] {
			f, err := os.Open(cloudsPath)
			if err != nil {
				continue
			}
			defer f.Close()
			options.readers[ctype] = f
			break
		}
		if options.readers[ctype] == nil {
			return map[string]Cloud{}, ErrFileNotFound{
				file:            string(ctype),
				searchLocations: options.locations[ctype]}
		}
	}

	// Parse the YAML payloads.
	var cloudsMap map[string]Cloud
	if ctype == Public {
		var clouds PublicClouds
		if err := yaml.NewDecoder(options.readers[ctype]).Decode(&clouds); err != nil {
			return map[string]Cloud{}, err
		}
		cloudsMap = clouds.Clouds
	} else {
		var clouds Clouds
		if err := yaml.NewDecoder(options.readers[ctype]).Decode(&clouds); err != nil {
			return map[string]Cloud{}, err
		}
		cloudsMap = clouds.Clouds
	}

	return cloudsMap, nil
}

func getUserConfig() (string, error) {
	// Use XDG_CONFIG_HOME or fall back to ~/.config, matching the
	// OpenStack convention for clouds.yaml location on all platforms.
	userConfig := os.Getenv("XDG_CONFIG_HOME")
	if userConfig != "" {
		return userConfig, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get the user home directory: %w", err)
	}
	userConfig = path.Join(homeDir, ".config")
	return userConfig, nil
}

// computeAvailability is a helper method to determine the endpoint type
// requested by the user.
func computeAvailability(endpointType string) gophercloud.Availability {
	if endpointType == "internal" || endpointType == "internalURL" {
		return gophercloud.AvailabilityInternal
	}
	if endpointType == "admin" || endpointType == "adminURL" {
		return gophercloud.AvailabilityAdmin
	}
	return gophercloud.AvailabilityPublic
}

// coalesce returns the first argument that is not the zero value for its type,
// or the zero value for its type.
func coalesce[T comparable](items ...T) T {
	var t T
	for _, item := range items {
		if item != t {
			return item
		}
	}
	return t
}

// mergeClouds merges two Clouds recursively (the AuthInfo also gets merged).
// In case both Clouds define a value, the value in the 'override' cloud takes precedence
func mergeClouds(override, cloud Cloud) (Cloud, error) {
	overrideJson, err := json.Marshal(override)
	if err != nil {
		return Cloud{}, err
	}
	cloudJson, err := json.Marshal(cloud)
	if err != nil {
		return Cloud{}, err
	}
	var overrideInterface any
	err = json.Unmarshal(overrideJson, &overrideInterface)
	if err != nil {
		return Cloud{}, err
	}
	var cloudInterface any
	err = json.Unmarshal(cloudJson, &cloudInterface)
	if err != nil {
		return Cloud{}, err
	}
	var mergedCloud Cloud
	mergedInterface := mergeInterfaces(overrideInterface, cloudInterface)
	mergedJson, err := json.Marshal(mergedInterface)
	if err != nil {
		return Cloud{}, err
	}
	err = json.Unmarshal(mergedJson, &mergedCloud)
	if err != nil {
		return Cloud{}, err
	}
	return mergedCloud, nil
}

// merges two interfaces. In cases where a value is defined for both 'overridingInterface' and
// 'inferiorInterface' the value in 'overridingInterface' will take precedence.
func mergeInterfaces(overridingInterface, inferiorInterface any) any {
	switch overriding := overridingInterface.(type) {
	case map[string]any:
		interfaceMap, ok := inferiorInterface.(map[string]any)
		if !ok {
			return overriding
		}
		for k, v := range interfaceMap {
			if overridingValue, ok := overriding[k]; ok {
				overriding[k] = mergeInterfaces(overridingValue, v)
			} else {
				overriding[k] = v
			}
		}
	case []any:
		list, ok := inferiorInterface.([]any)
		if !ok {
			return overriding
		}

		return append(overriding, list...)
	case nil:
		// mergeClouds(nil, map[string]interface{...}) -> map[string]interface{...}
		v, ok := inferiorInterface.(map[string]any)
		if ok {
			return v
		}
	}
	// We don't want to override with empty values
	if reflect.DeepEqual(overridingInterface, nil) || reflect.DeepEqual(reflect.Zero(reflect.TypeOf(overridingInterface)).Interface(), overridingInterface) {
		return inferiorInterface
	} else {
		return overridingInterface
	}
}
