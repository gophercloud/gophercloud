package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/racker/perigee"
)

// Version is a supported API version, corresponding to a vN package within the appropriate service.
type Version struct {
	ID       string
	Priority int
}

var goodStatus = map[string]bool{
	"current":   true,
	"supported": true,
	"stable":    true,
}

// ChooseVersion queries the base endpoint of a API to choose the most recent non-experimental alternative from a service's
// published versions.
// It returns the highest-Priority Version among the alternatives that are provided, as well as its corresponding endpoint.
func ChooseVersion(identityEndpoint string, recognized []*Version) (*Version, string, error) {
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

	// Normalize the identity endpoint that's provided by trimming any path, query or fragment from the URL.
	u, err := url.Parse(identityEndpoint)
	if err != nil {
		return nil, "", err
	}
	u.Path, u.RawQuery, u.Fragment = "", "", ""
	normalized := u.String()

	var resp response
	_, err = perigee.Request("GET", normalized, perigee.Options{
		Results: &resp,
		OkCodes: []int{200, 300},
	})

	if err != nil {
		return nil, "", err
	}

	byID := make(map[string]*Version)
	for _, version := range recognized {
		byID[version.ID] = version
	}

	var highest *Version
	var endpoint string

	normalize := func(endpoint string) string {
		if !strings.HasSuffix(endpoint, "/") {
			return endpoint + "/"
		}
		return endpoint
	}
	normalizedGiven := normalize(identityEndpoint)

	for _, value := range resp.Versions.Values {
		href := ""
		for _, link := range value.Links {
			if link.Rel == "self" {
				href = normalize(link.Href)
			}
		}

		if matching, ok := byID[value.ID]; ok {
			// Prefer a version that exactly matches the provided endpoint.
			if href == normalizedGiven {
				if href == "" {
					return nil, "", fmt.Errorf("Endpoint missing in version %s response from %s", value.ID, normalized)
				}
				return matching, href, nil
			}

			// Otherwise, find the highest-priority version with a whitelisted status.
			if goodStatus[strings.ToLower(value.Status)] {
				if highest == nil || matching.Priority > highest.Priority {
					highest = matching
					endpoint = href
				}
			}
		}
	}

	if highest == nil {
		return nil, "", fmt.Errorf("No supported version available from endpoint %s", normalized)
	}
	if endpoint == "" {
		return nil, "", fmt.Errorf("Endpoint missing in version %s response from %s", highest.ID, normalized)
	}

	return highest, endpoint, nil
}
