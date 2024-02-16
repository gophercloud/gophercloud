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
//	networkClient, err := openstack.NewNetworkV2(providerClient, eo)
//	if err != nil {
//		panic(err)
//	}
package clouds

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/gophercloud/gophercloud/v2"
	"gopkg.in/yaml.v2"
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
//     2. the directory `openstack` under the standatd user config location for
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
		locations: func() []string {
			if path := os.Getenv("OS_CLIENT_CONFIG_FILE"); path != "" {
				return []string{path}
			}
			return nil
		}(),
	}

	for _, apply := range opts {
		apply(&options)
	}

	if options.cloudName == "" {
		return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("the empty string \"\" is not a valid cloud name")
	}

	// Set the defaults and open the files for reading. This code only runs
	// if no override has been set, because it is fallible.
	if options.cloudsyamlReader == nil {
		if len(options.locations) < 1 {
			cwd, err := os.Getwd()
			if err != nil {
				return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("failed to get the current working directory: %w", err)
			}
			userConfig, err := os.UserConfigDir()
			if err != nil {
				return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("failed to get the user config directory: %w", err)
			}
			options.locations = []string{path.Join(cwd, "clouds.yaml"), path.Join(userConfig, "openstack", "clouds.yaml"), path.Join("/etc", "openstack", "clouds.yaml")}
		}

		for _, cloudsPath := range options.locations {
			var errNotFound *os.PathError
			f, err := os.Open(cloudsPath)
			if err != nil && !errors.As(err, &errNotFound) {
				return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("failed to open %q: %w", cloudsPath, err)
			}
			if err == nil {
				defer f.Close()
				options.cloudsyamlReader = f

				if options.secureyamlReader == nil {
					securePath := path.Join(path.Base(cloudsPath), "secure.yaml")
					secureF, err := os.Open(securePath)
					if err != nil && !errors.As(err, &errNotFound) {
						return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("failed to open %q: %w", securePath, err)
					}
					if err == nil {
						defer secureF.Close()
						options.secureyamlReader = secureF
					}
				}
			}
		}
		if options.cloudsyamlReader == nil {
			return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("clouds file not found. Search locations were: %v", options.locations)
		}
	}

	// Parse the YAML payloads.
	var clouds Clouds
	if err := yaml.NewDecoder(options.cloudsyamlReader).Decode(&clouds); err != nil {
		return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, err
	}

	cloud, ok := clouds.Clouds[options.cloudName]
	if !ok {
		return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("cloud %q not found in clouds.yaml", options.cloudName)
	}

	if options.secureyamlReader != nil {
		var secureClouds Clouds
		if err := yaml.NewDecoder(options.secureyamlReader).Decode(&secureClouds); err != nil {
			return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("failed to parse secure.yaml: %w", err)
		}

		if secureCloud, ok := secureClouds.Clouds[options.cloudName]; ok {
			// If secureCloud has content and it differs from the cloud entry,
			// merge the two together.
			if !reflect.DeepEqual((gophercloud.AuthOptions{}), secureClouds) && !reflect.DeepEqual(clouds, secureClouds) {
				var err error
				cloud, err = mergeClouds(secureCloud, cloud)
				if err != nil {
					return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("unable to merge information from clouds.yaml and secure.yaml")
				}
			}
		}
	}

	tlsConfig, err := computeTLSConfig(cloud, options)
	if err != nil {
		return gophercloud.AuthOptions{}, gophercloud.EndpointOpts{}, nil, fmt.Errorf("unable to compute TLS configuration: %w", err)
	}

	endpointType := coalesce(options.endpointType, cloud.EndpointType, cloud.Interface)

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
			Scope:                       options.scope,
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

// coalesce returns the first argument that is not the empty string, or the
// empty string.
func coalesce(items ...string) string {
	for _, item := range items {
		if item != "" {
			return item
		}
	}
	return ""
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
	var overrideInterface interface{}
	err = json.Unmarshal(overrideJson, &overrideInterface)
	if err != nil {
		return Cloud{}, err
	}
	var cloudInterface interface{}
	err = json.Unmarshal(cloudJson, &cloudInterface)
	if err != nil {
		return Cloud{}, err
	}
	var mergedCloud Cloud
	mergedInterface := mergeInterfaces(overrideInterface, cloudInterface)
	mergedJson, err := json.Marshal(mergedInterface)
	err = json.Unmarshal(mergedJson, &mergedCloud)
	if err != nil {
		return Cloud{}, err
	}
	return mergedCloud, nil
}

// merges two interfaces. In cases where a value is defined for both 'overridingInterface' and
// 'inferiorInterface' the value in 'overridingInterface' will take precedence.
func mergeInterfaces(overridingInterface, inferiorInterface interface{}) interface{} {
	switch overriding := overridingInterface.(type) {
	case map[string]interface{}:
		interfaceMap, ok := inferiorInterface.(map[string]interface{})
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
	case []interface{}:
		list, ok := inferiorInterface.([]interface{})
		if !ok {
			return overriding
		}

		return append(overriding, list...)
	case nil:
		// mergeClouds(nil, map[string]interface{...}) -> map[string]interface{...}
		v, ok := inferiorInterface.(map[string]interface{})
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
