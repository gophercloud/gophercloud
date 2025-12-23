package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/availabilityzoneprofiles"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const AvailabilityZoneProfilesListBody = `
{
	"availability_zone_profiles": [
        {
            "id": "1d334061-d807-4997-8f34-9fe428ba37df",
            "name": "availability-zone-profile-first",
            "provider_name": "amphora",
            "availability_zone_data": "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}"
        },
		{
            "id": "56f45d00-86e4-4bea-8525-19e835776c4e",
            "name": "availability-zone-profile-second",
            "provider_name": "amphora",
            "availability_zone_data": "{\"compute_zone\": \"not_nova\", \"volume_zone\": \"not_nova\"}"
        }
    ]
}
`

const SingleAvailabilityZoneProfileBody = `
{
	"availability_zone_profile": {
		"id": "13be083b-f502-426e-8500-07600f98b91b",
		"name": "availability-zone-profile",
		"provider_name": "amphora",
		"availability_zone_data": "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}"
	}
}
`

const PostUpdateAvailabilityZoneFlavorBody = `
{
	"availability_zone_profile": {
		"id": "13be083b-f502-426e-8500-07600f98b91b",
		"name": "availability-zone-profile-updated",
		"provider_name": "amphora",
		"availability_zone_data": "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}"
	}
}
`

var (
	AvailabilityZoneProfileSingle = availabilityzoneprofiles.AvailabilityZoneProfile{
		ID:                   "1d334061-d807-4997-8f34-9fe428ba37df",
		Name:                 "availability-zone-profile-first",
		ProviderName:         "amphora",
		AvailabilityZoneData: "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}",
	}

	AvailabilityZoneProfileAct = availabilityzoneprofiles.AvailabilityZoneProfile{
		ID:                   "56f45d00-86e4-4bea-8525-19e835776c4e",
		Name:                 "availability-zone-profile-second",
		ProviderName:         "amphora",
		AvailabilityZoneData: "{\"compute_zone\": \"not_nova\", \"volume_zone\": \"not_nova\"}",
	}

	AvailabilityZoneProfileDb = availabilityzoneprofiles.AvailabilityZoneProfile{
		ID:                   "13be083b-f502-426e-8500-07600f98b91b",
		Name:                 "availability-zone-profile",
		ProviderName:         "amphora",
		AvailabilityZoneData: "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}",
	}

	AvailabilityZoneProfileUpdated = availabilityzoneprofiles.AvailabilityZoneProfile{
		ID:                   "13be083b-f502-426e-8500-07600f98b91b",
		Name:                 "availability-zone-profile-updated",
		ProviderName:         "amphora",
		AvailabilityZoneData: "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}",
	}
)

func HandleAvailabilityZoneProfileListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v2.0/lbaas/availabilityzoneprofiles", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprint(w, AvailabilityZoneProfilesListBody)
		case "56f45d00-86e4-4bea-8525-19e835776c4e":
			fmt.Fprint(w, `{ "availability_zone_profiles": [] }`)
		default:
			t.Fatalf("/v2.0/lbaas/availabilityzoneprofiles invoked with unexpected marker=[%s]", marker)
		}
	})
}

func HandleAvailabilityZoneProfileCreationSuccessfully(t *testing.T, fakeServer th.FakeServer, response string) {
	fakeServer.Mux.HandleFunc("/v2.0/lbaas/availabilityzoneprofiles", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
			"availability_zone_profile": {
				"name": "availability-zone-profile",
				"provider_name": "amphora",
				"availability_zone_data":  "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}"
			}
		}`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, response)
	})
}

func HandleAvailabilityZoneProfileGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v2.0/lbaas/availabilityzoneprofiles/13be083b-f502-426e-8500-07600f98b91b", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprint(w, SingleAvailabilityZoneProfileBody)
	})
}

func HandleAvailabilityZoneProfileDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v2.0/lbaas/availabilityzoneprofiles/13be083b-f502-426e-8500-07600f98b91b", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleAvailabilityZoneProfileUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v2.0/lbaas/availabilityzoneprofiles/dcd65be5-f117-4260-ab3d-b32cc5bd1272", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `{
			"availability_zone_profile": {
				"name": "availability-zone-profile-updated",
				"provider_name": "amphora",
				"availability_zone_data": "{\"compute_zone\": \"nova\", \"volume_zone\": \"nova\"}"
			}
		}`)

		fmt.Fprint(w, PostUpdateAvailabilityZoneFlavorBody)
	})
}
