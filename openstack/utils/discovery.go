package utils

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

type SupportedMicroversions struct {
	MaxMajor int
	MaxMinor int
	MinMajor int
	MinMinor int
}

// GetServiceVersions returns the minimum and maximum microversion that is supported by the ServiceClient Endpoint.
func GetServiceVersions(ctx context.Context, client *gophercloud.ProviderClient, endpointURL string) (SupportedMicroversions, error) {
	type valueResp struct {
		ID         string `json:"id"`
		Status     string `json:"status"`
		Version    string `json:"version"`
		MinVersion string `json:"min_version"`
	}

	type response struct {
		Version  valueResp   `json:"version"`
		Versions []valueResp `json:"versions"`
	}
	var minVersion, maxVersion string
	var supportedMicroversions SupportedMicroversions
	var resp response
	_, err := client.Request(ctx, "GET", endpointURL, &gophercloud.RequestOpts{
		JSONResponse: &resp,
		OkCodes:      []int{200, 300},
	})

	if err != nil {
		return supportedMicroversions, err
	}

	if len(resp.Versions) > 0 {
		// We are dealing with an unversioned endpoint
		// We only handle the case when there is exactly one, and assume it is the correct one
		if len(resp.Versions) > 1 {
			return supportedMicroversions, fmt.Errorf("unversioned endpoint with multiple alternatives not supported")
		}
		minVersion = resp.Versions[0].MinVersion
		maxVersion = resp.Versions[0].Version
	} else {
		minVersion = resp.Version.MinVersion
		maxVersion = resp.Version.Version
	}

	// Return early if the endpoint does not support microversions
	if minVersion == "" && maxVersion == "" {
		return supportedMicroversions, fmt.Errorf("microversions not supported by endpoint")
	}

	supportedMicroversions.MinMajor, supportedMicroversions.MinMinor, err = ParseMicroversion(minVersion)
	if err != nil {
		return supportedMicroversions, err
	}

	supportedMicroversions.MaxMajor, supportedMicroversions.MaxMinor, err = ParseMicroversion(maxVersion)
	if err != nil {
		return supportedMicroversions, err
	}

	return supportedMicroversions, nil
}

// GetSupportedMicroversions returns the minimum and maximum microversion that is supported by the ServiceClient Endpoint.
func GetSupportedMicroversions(ctx context.Context, client *gophercloud.ServiceClient) (SupportedMicroversions, error) {
	return GetServiceVersions(ctx, client.ProviderClient, client.Endpoint)
}

// RequireMicroversion checks that the required microversion is supported and
// returns a ServiceClient with the microversion set.
func RequireMicroversion(ctx context.Context, client gophercloud.ServiceClient, required string) (gophercloud.ServiceClient, error) {
	supportedMicroversions, err := GetSupportedMicroversions(ctx, &client)
	if err != nil {
		return client, fmt.Errorf("unable to determine supported microversions: %w", err)
	}
	supported, err := supportedMicroversions.IsSupported(required)
	if err != nil {
		return client, err
	}
	if !supported {
		return client, fmt.Errorf("microversion %s not supported. Supported versions: %v", required, supportedMicroversions)
	}
	client.Microversion = required
	return client, nil
}

// IsSupported checks if a microversion falls in the supported interval.
// It returns true if the version is within the interval and false otherwise.
func (supported SupportedMicroversions) IsSupported(version string) (bool, error) {
	// Parse the version X.Y into X and Y integers that are easier to compare.
	vMajor, vMinor, err := ParseMicroversion(version)
	if err != nil {
		return false, err
	}

	// Check that the major version number is supported.
	if (vMajor < supported.MinMajor) || (vMajor > supported.MaxMajor) {
		return false, nil
	}

	// Check that the minor version number is supported
	if (vMinor <= supported.MaxMinor) && (vMinor >= supported.MinMinor) {
		return true, nil
	}

	return false, nil
}

// ParseMicroversion parses the version major.minor into separate integers major and minor.
// For example, "2.53" becomes 2 and 53.
func ParseMicroversion(version string) (major int, minor int, err error) {
	parts := strings.Split(version, ".")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid microversion format: %q", version)
	}
	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	minor, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return major, minor, nil
}
