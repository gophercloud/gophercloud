package testing

import (
	"context"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/segments"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetSegment(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/segments/"+SegmentID1, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(createResponse))
		th.AssertNoErr(t, err)
	})

	res, err := segments.Get(context.TODO(), fake.ServiceClient(fakeServer), SegmentID1).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Segment1, *res)
}

func TestListSegments(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/segments", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(SegmentsListBody))
		th.AssertNoErr(t, err)
	})

	count := 0
	pager := segments.List(fake.ServiceClient(fakeServer), nil)
	err := pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := segments.ExtractSegments(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, []segments.Segment{Segment1, Segment2}, actual)
		return true, nil
	})
	th.AssertNoErr(t, err)
	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreateSegment(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/segments", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, createRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(createResponse))
		th.AssertNoErr(t, err)
	})

	opts := segments.CreateOpts{
		NetworkID:       Segment1.NetworkID,
		NetworkType:     "flat",
		PhysicalNetwork: "public",
		Name:            "seg1",
		Description:     "desc",
	}
	actual, err := segments.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, Segment1, *actual)
}

func TestUpdateSegment(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/segments/"+SegmentID1, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestJSONRequest(t, r, updateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(updateResponse))
		th.AssertNoErr(t, err)
	})

	newName := "new-name"
	newDesc := "new-desc"
	opts := segments.UpdateOpts{
		Name:        &newName,
		Description: &newDesc,
	}
	actual, err := segments.Update(context.TODO(), fake.ServiceClient(fakeServer), SegmentID1, opts).Extract()
	th.AssertNoErr(t, err)

	expected := Segment1
	expected.Name = newName
	expected.Description = newDesc

	th.CheckDeepEquals(t, expected, *actual)
}

func TestDeleteSegment(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/segments/"+SegmentID1, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	err := segments.Delete(context.TODO(), fake.ServiceClient(fakeServer), SegmentID1).ExtractErr()
	th.AssertNoErr(t, err)
}
