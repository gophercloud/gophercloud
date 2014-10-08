package openstack

import (
	"strings"
	"testing"

	"github.com/rackspace/gophercloud"
	tokens2 "github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
	th "github.com/rackspace/gophercloud/testhelper"
)

// Service catalog fixtures take too much vertical space!
var catalog2 = tokens2.ServiceCatalog{
	Entries: []tokens2.CatalogEntry{
		tokens2.CatalogEntry{
			Type: "same",
			Name: "same",
			Endpoints: []tokens2.Endpoint{
				tokens2.Endpoint{
					Region:      "same",
					PublicURL:   "https://public.correct.com/",
					InternalURL: "https://internal.correct.com/",
					AdminURL:    "https://admin.correct.com/",
				},
				tokens2.Endpoint{
					Region:    "different",
					PublicURL: "https://badregion.com/",
				},
			},
		},
		tokens2.CatalogEntry{
			Type: "same",
			Name: "different",
			Endpoints: []tokens2.Endpoint{
				tokens2.Endpoint{
					Region:    "same",
					PublicURL: "https://badname.com/",
				},
				tokens2.Endpoint{
					Region:    "different",
					PublicURL: "https://badname.com/+badregion",
				},
			},
		},
		tokens2.CatalogEntry{
			Type: "different",
			Name: "different",
			Endpoints: []tokens2.Endpoint{
				tokens2.Endpoint{
					Region:    "same",
					PublicURL: "https://badtype.com/+badname",
				},
				tokens2.Endpoint{
					Region:    "different",
					PublicURL: "https://badtype.com/+badregion+badname",
				},
			},
		},
	},
}

func TestV2EndpointExact(t *testing.T) {
	expectedURLs := map[gophercloud.Availability]string{
		gophercloud.AvailabilityPublic:   "https://public.correct.com/",
		gophercloud.AvailabilityAdmin:    "https://admin.correct.com/",
		gophercloud.AvailabilityInternal: "https://internal.correct.com/",
	}

	for availability, expected := range expectedURLs {
		actual, err := V2EndpointURL(&catalog2, gophercloud.EndpointOpts{
			Type:         "same",
			Name:         "same",
			Region:       "same",
			Availability: availability,
		})
		th.AssertNoErr(t, err)
		th.CheckEquals(t, expected, actual)
	}
}

func TestV2EndpointNone(t *testing.T) {
	_, err := V2EndpointURL(&catalog2, gophercloud.EndpointOpts{
		Type:         "nope",
		Availability: gophercloud.AvailabilityPublic,
	})
	th.CheckEquals(t, gophercloud.ErrEndpointNotFound, err)
}

func TestV2EndpointMultiple(t *testing.T) {
	_, err := V2EndpointURL(&catalog2, gophercloud.EndpointOpts{
		Type:         "same",
		Region:       "same",
		Availability: gophercloud.AvailabilityPublic,
	})
	if !strings.HasPrefix(err.Error(), "Discovered 2 matching endpoints:") {
		t.Errorf("Received unexpected error: %v", err)
	}
}

func TestV2EndpointBadAvailability(t *testing.T) {
	_, err := V2EndpointURL(&catalog2, gophercloud.EndpointOpts{
		Type:         "same",
		Name:         "same",
		Region:       "same",
		Availability: "wat",
	})
	th.CheckEquals(t, err.Error(), "Unexpected availability in endpoint query: wat")
}
