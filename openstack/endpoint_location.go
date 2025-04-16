package openstack

import (
	"slices"

	"github.com/gophercloud/gophercloud/v2"
	tokens2 "github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	tokens3 "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

/*
V2EndpointURL discovers the endpoint URL for a specific service from a
ServiceCatalog acquired during the v2 identity service.

The specified EndpointOpts are used to identify a unique, unambiguous endpoint
to return. It's an error both when multiple endpoints match the provided
criteria and when none do. The minimum that can be specified is a Type, but you
will also often need to specify a Name and/or a Region depending on what's
available on your OpenStack deployment.
*/
func V2EndpointURL(catalog *tokens2.ServiceCatalog, opts gophercloud.EndpointOpts) (string, error) {
	// Extract Endpoints from the catalog entries that match the requested Type, Name if provided, and Region if provided.
	var endpoints = make([]tokens2.Endpoint, 0, 1)
	for _, entry := range catalog.Entries {
		if (slices.Contains(opts.Types(), entry.Type)) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Region == "" || endpoint.Region == opts.Region {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if there were no matching endpoints.
	if len(endpoints) == 0 {
		err := &gophercloud.ErrEndpointNotFound{}
		return "", err
	}

	// Extract the appropriate URL from the matching Endpoint.
	//
	// If multiple endpoints were found, use the first result and disregard the other endpoints.
	// This behavior matches the Python library. See GH-1764.
	switch opts.Availability {
	case gophercloud.AvailabilityPublic:
		return gophercloud.NormalizeURL(endpoints[0].PublicURL), nil
	case gophercloud.AvailabilityInternal:
		return gophercloud.NormalizeURL(endpoints[0].InternalURL), nil
	case gophercloud.AvailabilityAdmin:
		return gophercloud.NormalizeURL(endpoints[0].AdminURL), nil
	default:
		err := &ErrInvalidAvailabilityProvided{}
		err.Argument = "Availability"
		err.Value = opts.Availability
		return "", err
	}
}

/*
V3EndpointURL discovers the endpoint URL for a specific service from a Catalog
acquired during the v3 identity service.

The specified EndpointOpts are used to identify a unique, unambiguous endpoint
to return. It's an error both when multiple endpoints match the provided
criteria and when none do. The minimum that can be specified is a Type, but you
will also often need to specify a Name and/or a Region depending on what's
available on your OpenStack deployment.
*/
func V3EndpointURL(catalog *tokens3.ServiceCatalog, opts gophercloud.EndpointOpts) (string, error) {
	if opts.Availability != gophercloud.AvailabilityAdmin &&
		opts.Availability != gophercloud.AvailabilityPublic &&
		opts.Availability != gophercloud.AvailabilityInternal {
		err := &ErrInvalidAvailabilityProvided{}
		err.Argument = "Availability"
		err.Value = opts.Availability
		return "", err
	}

	// Extract Endpoints from the catalog entries that match the requested Type, Interface,
	// Name if provided, and Region if provided.
	var endpoints = make([]tokens3.Endpoint, 0, 1)
	for _, entry := range catalog.Entries {
		if (slices.Contains(opts.Types(), entry.Type)) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if (opts.Availability == gophercloud.Availability(endpoint.Interface)) &&
					(opts.Region == "" || endpoint.Region == opts.Region || endpoint.RegionID == opts.Region) {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if there were no matching endpoints.
	if len(endpoints) == 0 {
		err := &gophercloud.ErrEndpointNotFound{}
		return "", err
	}

	// Extract the URL from the matching Endpoint.
	//
	// If multiple endpoints were found, use the first result and disregard the other endpoints.
	// This behavior matches the Python library. See GH-1764.
	return gophercloud.NormalizeURL(endpoints[0].URL), nil
}
