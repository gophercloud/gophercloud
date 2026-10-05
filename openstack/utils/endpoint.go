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
	// Swift doesn't support version discovery.
	if expectedVersion == 0 || serviceType == "object-store" {
		return true, nil
	}

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
