package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/bays"
	fake "github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/common"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/bays", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "bays": [
    {
      "status": "CREATE_COMPLETE",
      "uuid": "a56a6cd8-0779-461b-b1eb-26cec904284a",
      "links": [
        {
          "href": "http://65.61.151.130:9511/v1/bays/a56a6cd8-0779-461b-b1eb-26cec904284a",
          "rel": "self"
        },
        {
          "href": "http://65.61.151.130:9511/bays/a56a6cd8-0779-461b-b1eb-26cec904284a",
          "rel": "bookmark"
        }
      ],
      "stack_id": "f8ef771f-1ffa-4ad5-99b8-651bf7669f80",
      "master_count": 1,
      "baymodel_id": "5b793604-fc76-4886-a834-ed522812cdcb",
      "node_count": 1,
      "bay_create_timeout": 0,
      "name": "k8sbay"
    }
  ]
}
			`)
	})

	client := fake.ServiceClient()
	count := 0

	results := bays.List(client, bays.ListOpts{})

	err := results.EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := bays.ExtractBays(page)
		if err != nil {
			t.Errorf("Failed to extract bays: %v", err)
			return false, err
		}

		expected := []bays.Bay{
			{
				Status:     "CREATE_COMPLETE",
				Name:       "k8sbay",
				ID:         "a56a6cd8-0779-461b-b1eb-26cec904284a",
				Masters:    1,
				Nodes:      1,
				BayModelID: "5b793604-fc76-4886-a834-ed522812cdcb",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/bays/a56a6cd8-0779-461b-b1eb-26cec904284a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "status": "CREATE_COMPLETE",
  "uuid": "a56a6cd8-0779-461b-b1eb-26cec904284a",
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/bays/a56a6cd8-0779-461b-b1eb-26cec904284a",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/bays/a56a6cd8-0779-461b-b1eb-26cec904284a",
      "rel": "bookmark"
    }
  ],
  "stack_id": "f8ef771f-1ffa-4ad5-99b8-651bf7669f80",
  "created_at": "2016-07-14T23:58:50+00:00",
  "api_address": "https://172.29.248.18:6443",
  "discovery_url": "https://discovery.etcd.io/ac7f669ebe467d061c59bfe5b6a5f6fe",
  "updated_at": "2016-07-15T00:02:53+00:00",
  "master_count": 1,
  "baymodel_id": "5b793604-fc76-4886-a834-ed522812cdcb",
  "master_addresses": [
    "172.29.248.18"
  ],
  "node_count": 1,
  "node_addresses": [
    "172.29.248.19"
  ],
  "status_reason": "Stack CREATE completed successfully",
  "bay_create_timeout": 0,
  "name": "k8sbay"
}
			`)
	})

	b, err := bays.Get(fake.ServiceClient(), "a56a6cd8-0779-461b-b1eb-26cec904284a").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, b.Status, "CREATE_COMPLETE")
	th.AssertEquals(t, b.Name, "k8sbay")
	th.AssertEquals(t, b.BayModelID, "5b793604-fc76-4886-a834-ed522812cdcb")
	th.AssertEquals(t, b.Nodes, 1)
	th.AssertEquals(t, b.ID, "a56a6cd8-0779-461b-b1eb-26cec904284a")
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/bays", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
  "node_count": 2,
  "baymodel_id": "5b793604-fc76-4886-a834-ed522812cdcb",
  "name": "mycluster"
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
  "status": "CREATE_IN_PROGRESS",
  "uuid": "39109e8a-516e-41a4-8b1d-22e9a56e4aa2",
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/bays/39109e8a-516e-41a4-8b1d-22e9a56e4aa2",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/bays/39109e8a-516e-41a4-8b1d-22e9a56e4aa2",
      "rel": "bookmark"
    }
  ],
  "stack_id": "f27f1581-fb2e-4033-af93-d2cf19bb8462",
  "created_at": "2016-08-08T16:45:18+00:00",
  "api_address": null,
  "discovery_url": "https://discovery.etcd.io/32ccf12b42e75b6822ac18c2c0391e5f",
  "updated_at": null,
  "master_count": 1,
  "baymodel_id": "5b793604-fc76-4886-a834-ed522812cdcb",
  "master_addresses": null,
  "node_count": 2,
  "node_addresses": null,
  "status_reason": null,
  "bay_create_timeout": 60,
  "name": "mycluster"
}
		`)
	})

	options := bays.CreateOpts{Name: "mycluster", Nodes: 2, BayModelID: "5b793604-fc76-4886-a834-ed522812cdcb"}
	b, err := bays.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, b.Status, "CREATE_IN_PROGRESS")
	th.AssertEquals(t, b.Name, "mycluster")
	th.AssertEquals(t, b.BayModelID, "5b793604-fc76-4886-a834-ed522812cdcb")
	th.AssertEquals(t, b.ID, "39109e8a-516e-41a4-8b1d-22e9a56e4aa2")
	th.AssertEquals(t, b.Masters, 1)
	th.AssertEquals(t, b.Nodes, 2)
}
