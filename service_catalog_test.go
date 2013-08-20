package gophercloud

import (
	"testing"
)

func TestFindFirstEndpointByCriteria(t *testing.T) {
	endpoint := FindFirstEndpointByCriteria([]CatalogEntry{}, ApiCriteria{Name: "test"})
	if endpoint.PublicURL != "" {
		t.Error("Not expecting to find anything in an empty service catalog.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		[]CatalogEntry{
			{Name: "test"},
		},
		ApiCriteria{Name: "test"},
	)
	if endpoint.PublicURL != "" {
		t.Error("Even though we have a matching entry, no endpoints exist")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "http://localhost", "", ""),
		ApiCriteria{Name: "test"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("Looking for an endpoint by name but without region or version ID should match first entry endpoint.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "http://localhost", "", ""),
		ApiCriteria{Name: "test", Region: "RGN"},
	)
	if endpoint.PublicURL != "" {
		t.Error("If provided, the Region qualifier must exclude endpoints with missing or mismatching regions.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "http://localhost", "rgn", ""),
		ApiCriteria{Name: "test", Region: "RGN"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("Regions are case insensitive.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "http://localhost", "rgn", ""),
		ApiCriteria{Name: "test", Region: "RGN", VersionId: "2"},
	)
	if endpoint.PublicURL != "" {
		t.Error("Missing version ID means no match.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "http://localhost", "rgn", "3"),
		ApiCriteria{Name: "test", Region: "RGN", VersionId: "2"},
	)
	if endpoint.PublicURL != "" {
		t.Error("Mismatched version ID means no match.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "http://localhost", "rgn", "2"),
		ApiCriteria{Name: "test", Region: "RGN", VersionId: "2"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("All search criteria met; endpoint expected.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "http://localhost", "ord", "2"),
		ApiCriteria{Name: "test", VersionId: "2"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("Sometimes, you might not care what region your stuff is in.")
		return
	}
}

func catalog(name, url, region, version string) []CatalogEntry {
	return []CatalogEntry{
		{
			Name: name,
			Endpoints: []EntryEndpoint{
				{
					PublicURL: url,
					Region:    region,
					VersionId: version,
				},
			},
		},
	}
}
