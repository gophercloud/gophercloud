package openstack

import (
	"context"
	"slices"

	"github.com/gophercloud/gophercloud/v2"
	tokens2 "github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tokens"
	tokens3 "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/utils"
)

func endpointSupportsVersion(ctx context.Context, client *gophercloud.ProviderClient, serviceType, endpointURL string, expectedVersion int) (bool, error) {
	// Swift doesn't support version discovery :(
	if expectedVersion == 0 || serviceType == "object-store" {
		return true, nil
	}

	endpointURL, err := utils.BaseVersionedEndpoint(endpointURL)
	if err != nil {
		return false, err
	}

	supportedMicroversions, err := utils.GetServiceVersions(ctx, client, endpointURL)
	if err != nil {
		return false, err
	}

	return supportedMicroversions.MinMajor == 0 || supportedMicroversions.MinMajor == expectedVersion, nil
}

/*
V2Endpoint discovers the endpoint URL for a specific service from a
ServiceCatalog acquired during the v2 identity service.

The specified EndpointOpts are used to identify a unique, unambiguous endpoint
to return. It's an error both when multiple endpoints match the provided
criteria and when none do. The minimum that can be specified is a Type, but you
will also often need to specify a Name and/or a Region depending on what's
available on your OpenStack deployment.
*/
func V2Endpoint(ctx context.Context, client *gophercloud.ProviderClient, catalog *tokens2.ServiceCatalog, opts gophercloud.EndpointOpts) (string, error) {
	// Extract Endpoints from the catalog entries that match the requested Type, Name if provided, and Region if provided.
	//
	// If multiple endpoints are found, we return the first result and disregard the rest.
	// This behavior matches the Python library. See GH-1764.
	for _, entry := range catalog.Entries {
		if (slices.Contains(opts.Types(), entry.Type)) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Region != "" && endpoint.Region != opts.Region {
					continue
				}

				var endpointURL string
				switch opts.Availability {
				case gophercloud.AvailabilityPublic:
					endpointURL = gophercloud.NormalizeURL(endpoint.PublicURL)
				case gophercloud.AvailabilityInternal:
					endpointURL = gophercloud.NormalizeURL(endpoint.InternalURL)
				case gophercloud.AvailabilityAdmin:
					endpointURL = gophercloud.NormalizeURL(endpoint.AdminURL)
				default:
					err := &ErrInvalidAvailabilityProvided{}
					err.Argument = "Availability"
					err.Value = opts.Availability
					return "", err
				}

				endpointSupportsVersion, err := endpointSupportsVersion(ctx, client, entry.Type, endpointURL, opts.Version)
				if err != nil {
					return "", err
				}
				if !endpointSupportsVersion {
					continue
				}

				return endpointURL, nil
			}
		}
	}

	// Report an error if there were no matching endpoints.
	err := &gophercloud.ErrEndpointNotFound{}
	return "", err
}

/*
V3Endpoint discovers the endpoint URL for a specific service from a Catalog
acquired during the v3 identity service.

The specified EndpointOpts are used to identify a unique, unambiguous endpoint
to return. It's an error both when multiple endpoints match the provided
criteria and when none do. The minimum that can be specified is a Type, but you
will also often need to specify a Name and/or a Region depending on what's
available on your OpenStack deployment.
*/
func V3Endpoint(ctx context.Context, client *gophercloud.ProviderClient, catalog *tokens3.ServiceCatalog, opts gophercloud.EndpointOpts) (string, error) {
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
	//
	// If multiple endpoints are found, we return the first result and disregard the rest.
	// This behavior matches the Python library. See GH-1764.
	for _, entry := range catalog.Entries {
		if (slices.Contains(opts.Types(), entry.Type)) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Availability != gophercloud.Availability(endpoint.Interface) {
					continue
				}
				if opts.Region != "" && endpoint.Region != opts.Region && endpoint.RegionID != opts.Region {
					continue
				}

				endpointURL := gophercloud.NormalizeURL(endpoint.URL)

				endpointSupportsVersion, err := endpointSupportsVersion(ctx, client, entry.Type, endpointURL, opts.Version)
				if err != nil {
					return "", err
				}
				if !endpointSupportsVersion {
					continue
				}

				return endpointURL, nil
			}
		}
	}

	// Report an error if there were no matching endpoints.
	err := &gophercloud.ErrEndpointNotFound{}
	return "", err
}
