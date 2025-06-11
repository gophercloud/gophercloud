package utils

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

// Version is a supported API version, corresponding to a vN package within the appropriate service.
type Version struct {
	ID       string
	Suffix   string
	Priority int
}

var goodStatus = map[string]bool{
	"current":   true,
	"supported": true,
	"stable":    true,
}

// ChooseVersion queries the base endpoint of an API to choose the identity service version.
// It will pick a version among the recognized, taking into account the priority and avoiding
// experimental alternatives from the published versions. However, if the client specifies a full
// endpoint that is among the recognized versions, it will be used regardless of priority.
// It returns the highest-Priority Version, OR exact match with client endpoint,
// among the alternatives that are provided, as well as its corresponding endpoint.
func ChooseVersion(ctx context.Context, client *gophercloud.ProviderClient, recognized []*Version) (*Version, string, error) {
	type linkResp struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	}

	type valueResp struct {
		ID     string     `json:"id"`
		Status string     `json:"status"`
		Links  []linkResp `json:"links"`
	}

	type versionsResp struct {
		Values []valueResp `json:"values"`
	}

	type response struct {
		Versions versionsResp `json:"versions"`
	}

	normalize := func(endpoint string) string {
		if !strings.HasSuffix(endpoint, "/") {
			return endpoint + "/"
		}
		return endpoint
	}
	identityEndpoint := normalize(client.IdentityEndpoint)

	// If a full endpoint is specified, check version suffixes for a match first.
	for _, v := range recognized {
		if strings.HasSuffix(identityEndpoint, v.Suffix) {
			return v, identityEndpoint, nil
		}
	}

	var resp response
	_, err := client.Request(ctx, "GET", client.IdentityBase, &gophercloud.RequestOpts{
		JSONResponse: &resp,
		OkCodes:      []int{200, 300},
	})

	if err != nil {
		return nil, "", err
	}

	var highest *Version
	var endpoint string

	for _, value := range resp.Versions.Values {
		href := ""
		for _, link := range value.Links {
			if link.Rel == "self" {
				href = normalize(link.Href)
			}
		}

		for _, version := range recognized {
			if strings.Contains(value.ID, version.ID) {
				// Prefer a version that exactly matches the provided endpoint.
				if href == identityEndpoint {
					if href == "" {
						return nil, "", fmt.Errorf("Endpoint missing in version %s response from %s", value.ID, client.IdentityBase)
					}
					return version, href, nil
				}

				// Otherwise, find the highest-priority version with a whitelisted status.
				if goodStatus[strings.ToLower(value.Status)] {
					if highest == nil || version.Priority > highest.Priority {
						highest = version
						endpoint = href
					}
				}
			}
		}
	}

	if highest == nil {
		return nil, "", fmt.Errorf("No supported version available from endpoint %s", client.IdentityBase)
	}
	if endpoint == "" {
		return nil, "", fmt.Errorf("Endpoint missing in version %s response from %s", highest.ID, client.IdentityBase)
	}

	return highest, endpoint, nil
}

type Status string

const (
	StatusCurrent      Status = "CURRENT"
	StatusSupported    Status = "SUPPORTED"
	StatusDeprecated   Status = "DEPRECATED"
	StatusExperimental Status = "EXPERIMENTAL"
	StatusUnknown      Status = "UNKNOWN"
)

// SupportedVersion stores a normalized form of the API version data. It handles APIs that
// support microversions as well as those that do not.
type SupportedVersion struct {
	// Major is the major version number of the API
	Major int
	// Minor is the minor version number of the API
	Minor int
	// Status is the status of the API
	Status Status
	SupportedMicroversions
}

// SupportedMicroversions stores a normalized form of the maximum and minimum API microversions
// supported by a given service.
type SupportedMicroversions struct {
	// MaxMajor is the major version number of the maximum supported API microversion
	MaxMajor int
	// MaxMinor is the minor version number of the maximum supported API microversion
	MaxMinor int
	// MinMajor is the major version number of the minimum supported API microversion
	MinMajor int
	// MinMinor is the minor version number of the minimum supported API microversion
	MinMinor int
}

// GetServiceVersions returns the versions supported by the ServiceClient Endpoint.
// If the endpoint resolves to an unversioned discovery API, this should return one or more supported versions.
// If the endpoint resolves to a versioned discovery API, this should return exactly one supported version.
func GetServiceVersions(ctx context.Context, client *gophercloud.ProviderClient, endpointURL string) ([]SupportedVersion, error) {
	type valueResp struct {
		ID         string `json:"id"`
		Status     string `json:"status"`
		Version    string `json:"version,omitempty"`
		MaxVersion string `json:"max_version,omitempty"`
		MinVersion string `json:"min_version"`
	}
	type response struct {
		Version  valueResp   `json:"version"`
		Versions []valueResp `json:"versions"`
	}

	var supportedVersions []SupportedVersion

	var resp response
	_, err := client.Request(ctx, "GET", endpointURL, &gophercloud.RequestOpts{
		JSONResponse: &resp,
		OkCodes:      []int{200, 300},
	})

	if err != nil {
		return supportedVersions, err
	}

	var versions []valueResp
	if len(resp.Versions) > 0 {
		versions = resp.Versions
	} else {
		versions = append(versions, resp.Version)
	}

	for _, version := range versions {
		majorVersion, minorVersion, err := ParseVersion(version.ID)
		if err != nil {
			return supportedVersions, err
		}

		status, err := ParseStatus(version.Status)
		if err != nil {
			return supportedVersions, err
		}

		supportedVersion := SupportedVersion{
			Major:  majorVersion,
			Minor:  minorVersion,
			Status: status,
		}

		// Only normalize the microversions if there are microversions to normalize
		if (version.Version != "" || version.MaxVersion != "") && version.MinVersion != "" {
			supportedVersion.MinMajor, supportedVersion.MinMinor, err = ParseMicroversion(version.MinVersion)
			if err != nil {
				return supportedVersions, err
			}

			maxVersion := version.Version
			if maxVersion == "" {
				maxVersion = version.MaxVersion
			}
			supportedVersion.MaxMajor, supportedVersion.MaxMinor, err = ParseMicroversion(maxVersion)
			if err != nil {
				return supportedVersions, err
			}
		}

		supportedVersions = append(supportedVersions, supportedVersion)
	}

	sort.Slice(supportedVersions, func(i, j int) bool {
		return supportedVersions[i].Major < supportedVersions[j].Major || (supportedVersions[i].Major == supportedVersions[j].Major &&
			supportedVersions[i].Minor < supportedVersions[j].Minor)
	})

	return supportedVersions, nil
}

// GetSupportedMicroversions returns the minimum and maximum microversion that is supported by the ServiceClient Endpoint.
func GetSupportedMicroversions(ctx context.Context, client *gophercloud.ServiceClient) (SupportedMicroversions, error) {
	var supportedMicroversions SupportedMicroversions

	supportedVersions, err := GetServiceVersions(ctx, client.ProviderClient, client.Endpoint)
	if err != nil {
		return supportedMicroversions, err
	}

	// If there are multiple versions then we were handed an unversioned endpoint. These don't
	// provide microversion information, so we need to fail. Likewise, if there are no versions
	// then something has gone wrong and we also need to fail.
	if len(supportedVersions) > 1 {
		return supportedMicroversions, fmt.Errorf("unversioned endpoint with multiple alternatives not supported")
	} else if len(supportedVersions) == 0 {
		return supportedMicroversions, fmt.Errorf("microversions not supported by endpoint")
	}

	supportedMicroversions = supportedVersions[0].SupportedMicroversions

	if supportedMicroversions.MaxMajor == 0 &&
		supportedMicroversions.MaxMinor == 0 &&
		supportedMicroversions.MinMajor == 0 &&
		supportedMicroversions.MinMinor == 0 {
		return supportedMicroversions, fmt.Errorf("microversions not supported by endpoint")
	}

	return supportedMicroversions, err
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

// ParseVersion parsed the version strings v{MAJOR} and v{MAJOR}.{MINOR} into separate integers
// major and minor.
// For example, "v2.1" becomes 2 and 1, "v3" becomes 3 and 0, and "1" becomes 1 and 0.
func ParseVersion(version string) (major, minor int, err error) {
	version = strings.TrimLeft(version, "v")

	parts := strings.Split(version, ".")
	if len(parts) == 1 {
		parts = append(parts, "0")
	} else if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format: %q", version)
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

func ParseStatus(status string) (Status, error) {
	switch strings.ToUpper(status) {
	case "CURRENT", "STABLE": // keystone uses STABLE instead of CURRENT
		return StatusCurrent, nil
	case "SUPPORTED":
		return StatusSupported, nil
	case "DEPRECATED":
		return StatusDeprecated, nil
	default:
		return StatusUnknown, fmt.Errorf("invalid status: %q", status)
	}
}
