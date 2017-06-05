package openstack

import (
	"github.com/gophercloud/gophercloud"
	tokens2 "github.com/gophercloud/gophercloud/openstack/identity/v2/tokens"
	tokens3 "github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

// V2EndpointURL discovers the endpoint URL for a specific service from a ServiceCatalog acquired
// during the v2 identity service. The specified EndpointOpts are used to identify a unique,
// unambiguous endpoint to return. It's an error both when multiple endpoints match the provided
// criteria and when none do. The minimum that can be specified is a Type, but you will also often
// need to specify a Name and/or a Region depending on what's available on your OpenStack
// deployment.
func V2EndpointURL(catalog *tokens2.ServiceCatalog, opts gophercloud.EndpointOpts) (string, gophercloud.EndpointExtraInfo, error) {
	// Extract Endpoints from the catalog entries that match the requested Type, Name if provided, and Region if provided.
	var endpoints = make([]tokens2.Endpoint, 0, 1)
	var endpointExtraInfo gophercloud.EndpointExtraInfo
	for _, entry := range catalog.Entries {
		if (entry.Type == opts.Type) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Region == "" || endpoint.Region == opts.Region {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if the options were ambiguous.
	if len(endpoints) > 1 {
		err := &ErrMultipleMatchingEndpointsV2{}
		err.Endpoints = endpoints
		return "", gophercloud.EndpointExtraInfo{}, err
	}

	// Extract the appropriate URL from the matching Endpoint.
	for _, endpoint := range endpoints {
		switch opts.Availability {
		case gophercloud.AvailabilityPublic:
			endpointExtraInfo.InitializeV2(opts.Type, opts.Name, endpoint.AdminURL, endpoint.InternalURL, endpoint.PublicURL, endpoint.Region, endpoint.TenantID, endpoint.VersionID, endpoint.VersionInfo, endpoint.VersionList)
			return gophercloud.NormalizeURL(endpoint.PublicURL), endpointExtraInfo, nil
		case gophercloud.AvailabilityInternal:
			endpointExtraInfo.InitializeV2(opts.Type, opts.Name, endpoint.AdminURL, endpoint.InternalURL, endpoint.PublicURL, endpoint.Region, endpoint.TenantID, endpoint.VersionID, endpoint.VersionInfo, endpoint.VersionList)
			return gophercloud.NormalizeURL(endpoint.InternalURL), endpointExtraInfo, nil
		case gophercloud.AvailabilityAdmin:
			endpointExtraInfo.InitializeV2(opts.Type, opts.Name, endpoint.AdminURL, endpoint.InternalURL, endpoint.PublicURL, endpoint.Region, endpoint.TenantID, endpoint.VersionID, endpoint.VersionInfo, endpoint.VersionList)
			return gophercloud.NormalizeURL(endpoint.AdminURL), endpointExtraInfo, nil
		default:
			err := &ErrInvalidAvailabilityProvided{}
			err.Argument = "Availability"
			err.Value = opts.Availability
			return "", gophercloud.EndpointExtraInfo{}, err
		}
	}

	// Report an error if there were no matching endpoints.
	err := &gophercloud.ErrEndpointNotFound{}
	return "", gophercloud.EndpointExtraInfo{}, err
}

// V3EndpointURL discovers the endpoint URL for a specific service from a Catalog acquired
// during the v3 identity service. The specified EndpointOpts are used to identify a unique,
// unambiguous endpoint to return. It's an error both when multiple endpoints match the provided
// criteria and when none do. The minimum that can be specified is a Type, but you will also often
// need to specify a Name and/or a Region depending on what's available on your OpenStack
// deployment.
func V3EndpointURL(catalog *tokens3.ServiceCatalog, opts gophercloud.EndpointOpts) (string, gophercloud.EndpointExtraInfo, error) {
	// Extract Endpoints from the catalog entries that match the requested Type, Interface,
	// Name if provided, and Region if provided.
	var endpoints = make([]tokens3.Endpoint, 0, 1)
	var endpointExtraInfo gophercloud.EndpointExtraInfo
	for _, entry := range catalog.Entries {
		if (entry.Type == opts.Type) && (opts.Name == "" || entry.Name == opts.Name) {
			for _, endpoint := range entry.Endpoints {
				if opts.Availability != gophercloud.AvailabilityAdmin &&
					opts.Availability != gophercloud.AvailabilityPublic &&
					opts.Availability != gophercloud.AvailabilityInternal {
					err := &ErrInvalidAvailabilityProvided{}
					err.Argument = "Availability"
					err.Value = opts.Availability
					return "", gophercloud.EndpointExtraInfo{}, err
				}
				if (opts.Availability == gophercloud.Availability(endpoint.Interface)) &&
					(opts.Region == "" || endpoint.Region == opts.Region) {
					endpoints = append(endpoints, endpoint)
				}
			}
		}
	}

	// Report an error if the options were ambiguous.
	if len(endpoints) > 1 {
		return "", gophercloud.EndpointExtraInfo{}, ErrMultipleMatchingEndpointsV3{Endpoints: endpoints}
	}

	// Extract the URL from the matching Endpoint.
	for _, endpoint := range endpoints {
		url := gophercloud.NormalizeURL(endpoint.URL)
		endpointExtraInfo.InitializeV3(opts.Type, opts.Name, endpoint.ID, endpoint.Interface, endpoint.Region, endpoint.URL)
		return url, gophercloud.EndpointExtraInfo{}, nil
	}

	// Report an error if there were no matching endpoints.
	err := &gophercloud.ErrEndpointNotFound{}
	return "", gophercloud.EndpointExtraInfo{}, err
}
