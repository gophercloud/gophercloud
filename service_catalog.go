package gophercloud

import (
	"strings"
)

// ApiCriteria provides one or more criteria for the SDK to look for appropriate endpoints.
// Fields left unspecified or otherwise set to their zero-values are assumed to not be
// relevant, and do not participate in the endpoint search.
type ApiCriteria struct {
	// Name specifies the desired service catalog entry name.
	Name string

	// Type specifies the desired service catalog entry type.
	Type string

	// Region specifies the desired endpoint region.
	Region string

	// VersionId specifies the desired version of the endpoint.
	// Note that this field is matched exactly, and is (at present)
	// opaque to Gophercloud.  Thus, requesting a version 2
	// endpoint will _not_ match a version 3 endpoint.
	VersionId string

	// The UrlChoice field inidicates whether or not gophercloud
	// should use the public or internal endpoint URL if a
	// candidate endpoint is found.
	UrlChoice int
}

// The choices available for UrlChoice.  See the ApiCriteria structure for details.
const (
	PublicURL = iota
	InternalURL
)

// Given a set of criteria to match on, locate the first candidate endpoint
// in the provided service catalog.
//
// If nothing found, the result will be a zero-valued EntryEndpoint (all URLs
// set to "").
func FindFirstEndpointByCriteria(entries []CatalogEntry, ac ApiCriteria) EntryEndpoint {
	rgn := strings.ToUpper(ac.Region)

	for _, entry := range entries {
		if (ac.Name != "") && (ac.Name != entry.Name) {
			continue
		}

		if (ac.Type != "") && (ac.Type != entry.Type) {
			continue
		}

		for _, endpoint := range entry.Endpoints {
			if (ac.Region != "") && (rgn != strings.ToUpper(endpoint.Region)) {
				continue
			}

			if (ac.VersionId != "") && (ac.VersionId != endpoint.VersionId) {
				continue
			}

			return endpoint
		}
	}
	return EntryEndpoint{}
}
