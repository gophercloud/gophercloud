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
		catalog("test", "compute", "http://localhost", "", ""),
		ApiCriteria{Name: "test"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("Looking for an endpoint by name but without region or version ID should match first entry endpoint.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "compute", "http://localhost", "", ""),
		ApiCriteria{Type: "compute"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("Looking for an endpoint by type but without region or version ID should match first entry endpoint.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "compute", "http://localhost", "", ""),
		ApiCriteria{Type: "identity"},
	)
	if endpoint.PublicURL != "" {
		t.Error("Returned mismatched type.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "compute", "http://localhost", "", ""),
		ApiCriteria{Name: "test", Region: "RGN"},
	)
	if endpoint.PublicURL != "" {
		t.Error("If provided, the Region qualifier must exclude endpoints with missing or mismatching regions.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "compute", "http://localhost", "rgn", ""),
		ApiCriteria{Name: "test", Region: "RGN"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("Regions are case insensitive.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "compute", "http://localhost", "rgn", ""),
		ApiCriteria{Name: "test", Region: "RGN", VersionId: "2"},
	)
	if endpoint.PublicURL != "" {
		t.Error("Missing version ID means no match.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "compute", "http://localhost", "rgn", "3"),
		ApiCriteria{Name: "test", Region: "RGN", VersionId: "2"},
	)
	if endpoint.PublicURL != "" {
		t.Error("Mismatched version ID means no match.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "compute", "http://localhost", "rgn", "2"),
		ApiCriteria{Name: "test", Region: "RGN", VersionId: "2"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("All search criteria met; endpoint expected.")
		return
	}

	endpoint = FindFirstEndpointByCriteria(
		catalog("test", "compute", "http://localhost", "ord", "2"),
		ApiCriteria{Name: "test", VersionId: "2"},
	)
	if endpoint.PublicURL != "http://localhost" {
		t.Error("Sometimes, you might not care what region your stuff is in.")
		return
	}
}

func catalog(name, entry_type, url, region, version string) []CatalogEntry {
	return []CatalogEntry{
		{
			Name: name,
			Type: entry_type,
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
