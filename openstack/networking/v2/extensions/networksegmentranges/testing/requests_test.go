package testing

import (
	"context"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/networksegmentranges"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/network_segment_ranges", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`
{
    "network_segment_ranges": [
        {
            "id": "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
            "name": "range1",
            "default": false,
            "shared": false,
            "project_id": "3e2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
            "network_type": "vxlan",
            "physical_network": "",
            "minimum": 100,
            "maximum": 200,
            "used": {},
            "available": [100, 101, 102]
        }
    ]
}`))
		th.AssertNoErr(t, err)
	})

	count := 0
	err := networksegmentranges.List(fake.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := networksegmentranges.ExtractNetworkSegmentRanges(page)
		if err != nil {
			t.Errorf("Failed to extract network segment ranges: %v", err)
			return false, err
		}

		expected := []networksegmentranges.NetworkSegmentRange{
			{
				ID:              "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
				Name:            "range1",
				Default:         false,
				Shared:          false,
				ProjectID:       "3e2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
				NetworkType:     "vxlan",
				PhysicalNetwork: "",
				Minimum:         100,
				Maximum:         200,
				Used:            map[int]string{},
				Available:       []int{100, 101, 102},
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

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/network_segment_ranges", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "network_segment_range": {
        "name": "range1",
        "network_type": "vxlan",
        "minimum": 100,
        "maximum": 200
    }
}
		`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_, err := w.Write([]byte(`
{
    "network_segment_range": {
        "id": "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
        "name": "range1",
        "default": false,
        "shared": false,
        "project_id": "3e2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
        "network_type": "vxlan",
        "physical_network": "",
        "minimum": 100,
        "maximum": 200,
        "used": {},
        "available": []
    }
}
		`))
		th.AssertNoErr(t, err)
	})

	opts := networksegmentranges.CreateOpts{
		Name:        "range1",
		NetworkType: "vxlan",
		Minimum:     100,
		Maximum:     200,
	}
	n, err := networksegmentranges.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5", n.ID)
	th.AssertEquals(t, "range1", n.Name)
	th.AssertEquals(t, "vxlan", n.NetworkType)
	th.AssertEquals(t, 100, n.Minimum)
	th.AssertEquals(t, 200, n.Maximum)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/network_segment_ranges/59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`
{
    "network_segment_range": {
        "id": "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
        "name": "range1",
        "default": false,
        "shared": false,
        "project_id": "3e2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
        "network_type": "vxlan",
        "physical_network": "",
        "minimum": 100,
        "maximum": 200,
        "used": {},
        "available": [100, 101, 102]
    }
}
		`))
		th.AssertNoErr(t, err)
	})

	n, err := networksegmentranges.Get(context.TODO(), fake.ServiceClient(fakeServer), "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5", n.ID)
	th.AssertEquals(t, "range1", n.Name)
	th.AssertEquals(t, "vxlan", n.NetworkType)
	th.AssertEquals(t, 100, n.Minimum)
	th.AssertEquals(t, 200, n.Maximum)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/network_segment_ranges/59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "network_segment_range": {
        "name": "range1-updated",
        "minimum": 150,
        "maximum": 250
    }
}
		`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`
{
    "network_segment_range": {
        "id": "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
        "name": "range1-updated",
        "default": false,
        "shared": false,
        "project_id": "3e2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5",
        "network_type": "vxlan",
        "physical_network": "",
        "minimum": 150,
        "maximum": 250,
        "used": {},
        "available": []
    }
}
		`))
		th.AssertNoErr(t, err)
	})

	name := "range1-updated"
	minimum := 150
	maximum := 250
	opts := networksegmentranges.UpdateOpts{
		Name:    &name,
		Minimum: &minimum,
		Maximum: &maximum,
	}
	n, err := networksegmentranges.Update(context.TODO(), fake.ServiceClient(fakeServer), "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "range1-updated", n.Name)
	th.AssertEquals(t, 150, n.Minimum)
	th.AssertEquals(t, 250, n.Maximum)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/network_segment_ranges/59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := networksegmentranges.Delete(context.TODO(), fake.ServiceClient(fakeServer), "59b2f3a1-09aa-49e4-b9f7-8d7c81f7c7e5")
	th.AssertNoErr(t, res.Err)
}
