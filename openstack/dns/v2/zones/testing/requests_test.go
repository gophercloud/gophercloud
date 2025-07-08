package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	count := 0
	err := zones.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := zones.ExtractZones(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedZonesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListAllPages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	allPages, err := zones.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	allZones, err := zones.ExtractZones(allPages)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 2, len(allZones))
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer)

	actual, err := zones.Get(context.TODO(), client.ServiceClient(fakeServer), "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstZone, actual)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer)

	createOpts := zones.CreateOpts{
		Name:        "example.org.",
		Email:       "joe@example.org",
		Type:        "PRIMARY",
		TTL:         7200,
		Description: "This is an example zone.",
	}

	actual, err := zones.Create(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedZone, actual)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSuccessfully(t, fakeServer)

	var description = "Updated Description"
	updateOpts := zones.UpdateOpts{
		TTL:         600,
		Description: &description,
	}

	UpdatedZone := CreatedZone
	UpdatedZone.Status = "PENDING"
	UpdatedZone.Action = "UPDATE"
	UpdatedZone.TTL = 600
	UpdatedZone.Description = "Updated Description"

	actual, err := zones.Update(context.TODO(), client.ServiceClient(fakeServer), UpdatedZone.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &UpdatedZone, actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)

	DeletedZone := CreatedZone
	DeletedZone.Status = "PENDING"
	DeletedZone.Action = "DELETE"
	DeletedZone.TTL = 600
	DeletedZone.Description = "Updated Description"

	actual, err := zones.Delete(context.TODO(), client.ServiceClient(fakeServer), DeletedZone.ID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &DeletedZone, actual)
}

func TestShare(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/zones/zone-id/shares", func(w http.ResponseWriter, r *http.Request) {
		th.AssertEquals(t, r.Method, "POST")

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		th.AssertNoErr(t, err)

		var reqBody map[string]string
		err = json.Unmarshal(body, &reqBody)
		th.AssertNoErr(t, err)
		expectedBody := map[string]string{"target_project_id": "project-id"}
		th.CheckDeepEquals(t, expectedBody, reqBody)

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, ShareZoneResponse)
	})

	opts := zones.ShareZoneOpts{TargetProjectID: "project-id"}
	zone, err := zones.Share(context.TODO(), client.ServiceClient(fakeServer), "zone-id", opts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ShareZone, *zone)
}

func TestUnshare(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/zones/zone-id/shares/share-id", func(w http.ResponseWriter, r *http.Request) {
		th.AssertEquals(t, r.Method, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	err := zones.Unshare(context.TODO(), client.ServiceClient(fakeServer), "zone-id", "share-id").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListShares(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/zones/zone-id/shares", func(w http.ResponseWriter, r *http.Request) {
		th.AssertEquals(t, r.Method, "GET")
		th.AssertEquals(t, "true", r.Header.Get("X-Auth-All-Projects"))

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, ListSharesResponse)
	})

	opts := zones.ListSharesOpts{
		AllProjects: true,
	}
	pages, err := zones.ListShares(client.ServiceClient(fakeServer), "zone-id", opts).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := zones.ExtractZoneShares(pages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ListZoneShares, actual)
}
