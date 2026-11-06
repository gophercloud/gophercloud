package utils

import (
	"context"
	"regexp"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
)

var versionedServiceTypeAliasRegexp = regexp.MustCompile(`^.*v(\d)$`)

func extractServiceTypeVersion(serviceType string) int {
	matches := versionedServiceTypeAliasRegexp.FindAllStringSubmatch(serviceType, 1)
	if matches == nil {
		return 0
	}

	version, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return 0
	}
	return version
}

// EndpointSupportsVersion reports whether an endpoint supports the requested
// major API version. A nil client limits the check to versions implied by
// legacy service type aliases such as volumev2 and volumev3.
func EndpointSupportsVersion(ctx context.Context, client *gophercloud.ProviderClient, serviceType, endpointURL string, expectedVersion int) (bool, error) {
	// Swift doesn't support version discovery :(
	if expectedVersion == 0 || serviceType == "object-store" {
		return true, nil
	}

	// Repeating verbatim from keystoneauth1 [1]:
	//
	// > The sins of our fathers become the blood on our hands.
	// > If a user requests an old-style service type such as volumev2, then they
	// > are inherently requesting the major API version 2. It's not a good
	// > interface, but it's the one that was imposed on the world years ago
	// > because the client libraries all hid the version discovery document.
	// > In order to be able to ensure that a user who requests volumev2 does not
	// > get a block-storage endpoint that only provides v3 of the block-storage
	// > service, we need to pull the version out of the service_type. The
	// > service-types-authority will prevent the growth of new monstrosities such
	// > as this, but in order to move forward without breaking people, we have
	// > to just cry in the corner while striking ourselves with thorned branches.
	// > That said, for sure only do this hack for officially known service_types.
	//
	// So yeah, what mordred said.
	//
	// https://github.com/openstack/keystoneauth/blob/5.10.0/keystoneauth1/discover.py#L270-L290

	// Legacy service types such as volumev2 imply a specific major version.
	impliedVersion := extractServiceTypeVersion(serviceType)
	if impliedVersion != 0 && impliedVersion != expectedVersion {
		return false, nil
	}
	if client == nil {
		return true, nil
	}

	endpointURL, err := BaseVersionedEndpoint(endpointURL)
	if err != nil {
		return false, err
	}

	supportedVersions, err := GetServiceVersions(ctx, client, endpointURL, false)
	if err != nil {
		return false, err
	}

	for _, supportedVersion := range supportedVersions {
		if supportedVersion.Major == expectedVersion {
			return true, nil
		}
	}

	return false, nil
}
