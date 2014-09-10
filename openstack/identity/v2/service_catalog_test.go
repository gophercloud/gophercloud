package v2

import (
	"encoding/json"
	"testing"
)

func TestServiceCatalog(t *testing.T) {
	authResults := make(map[string]interface{})
	err := json.Unmarshal([]byte(authResultsOK), &authResults)
	if err != nil {
		t.Error(err)
		return
	}

	sc, err := GetServiceCatalog(authResults)
	if err != nil {
		panic(err)
	}

	if sc.NumberOfServices() != 3 {
		t.Errorf("Expected 3 services; got %d", sc.NumberOfServices())
	}

	ces, err := sc.CatalogEntries()
	if err != nil {
		t.Error(err)
		return
	}
	for _, ce := range ces {
		if strNotInStrList(ce.Name, "Cloud Servers", "Cloud Files", "DNS-as-a-Service") {
			t.Errorf("Expected \"%s\" to be one of Cloud Servers, Cloud Files, or DNS-as-a-Service", ce.Name)
			return
		}

		if strNotInStrList(ce.Type, "dnsextension:dns", "object-store", "compute") {
			t.Errorf("Expected \"%s\" to be one of dnsextension:dns, object-store, or compute")
			return
		}
	}

	eps := endpointsFor(ces, "compute")
	if len(eps) != 2 {
		t.Errorf("Expected 2 endpoints for compute service")
		return
	}
	for _, ep := range eps {
		if strNotInStrList(ep.VersionID, "1", "1.1", "1.1") {
			t.Errorf("Expected versionID field of compute resource to be one of 1 or 1.1")
			return
		}
	}

	eps = endpointsFor(ces, "object-store")
	if len(eps) != 2 {
		t.Errorf("Expected 2 endpoints for object-store service")
		return
	}
	for _, ep := range eps {
		if ep.VersionID != "1" {
			t.Errorf("Expected only version 1 object store API version")
			return
		}
	}

	eps = endpointsFor(ces, "dnsextension:dns")
	if len(eps) != 1 {
		t.Errorf("Expected 1 endpoint for DNS-as-a-Service service")
		return
	}
	if eps[0].VersionID != "2.0" {
		t.Errorf("Expected version 2.0 of DNS-as-a-Service service")
		return
	}
}

func endpointsFor(ces []CatalogEntry, t string) []Endpoint {
	for _, ce := range ces {
		if ce.Type == t {
			return ce.Endpoints
		}
	}
	panic("Precondition violated")
}

func strNotInStrList(needle, haystack1, haystack2, haystack3 string) bool {
	if (needle != haystack1) && (needle != haystack2) && (needle != haystack3) {
		return true
	}
	return false
}
