package utils

import (
	"fmt"

	"github.com/racker/perigee"
)

// Version is a supported API version, corresponding to a vN package within the appropriate service.
type Version struct {
	ID       string
	Priority int
}

// ChooseVersion queries the base endpoint of a API to choose the most recent non-experimental alternative from a service's
// published versions.
// It returns the highest-Priority Version among the alternatives that are provided, as well as its corresponding endpoint.
func ChooseVersion(baseEndpoint string, recognized []*Version) (*Version, string, error) {
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

	var resp response
	_, err := perigee.Request("GET", baseEndpoint, perigee.Options{
		Results: &resp,
		OkCodes: []int{200},
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

	for _, value := range resp.Versions.Values {
		if matching, ok := byID[value.ID]; ok {
			if highest == nil || matching.Priority > highest.Priority {
				highest = matching

				found := false
				for _, link := range value.Links {
					if link.Rel == "self" {
						found = true
						endpoint = link.Href
					}
				}

				if !found {
					return nil, "", fmt.Errorf("Endpoint missing in version %s response from %s", value.ID, baseEndpoint)
				}
			}
		}
	}

	if highest == nil || endpoint == "" {
		return nil, "", fmt.Errorf("No supported version available from endpoint %s", baseEndpoint)
	}

	return highest, endpoint, nil
}
