package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// SinglePage sample and test cases.

type SinglePageResult struct {
	pagination.SinglePageBase
}

func (r SinglePageResult) IsEmpty() (bool, error) {
	is, err := ExtractSingleInts(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractSingleInts(r pagination.Page) ([]int, error) {
	var s struct {
		Ints []int `json:"ints"`
	}
	err := (r.(SinglePageResult)).ExtractInto(&s)
	return s.Ints, err
}

func setupSinglePaged(fakeServer th.FakeServer) pagination.Pager {
	client := client.ServiceClient(fakeServer)

	fakeServer.Mux.HandleFunc("/only", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{ "ints": [1, 2, 3] }`)
	})

	createPage := func(r pagination.PageResult) pagination.Page {
		return SinglePageResult{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, fakeServer.Server.URL+"/only", createPage)
}

func TestEnumerateSinglePaged(t *testing.T) {
	callCount := 0
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	pager := setupSinglePaged(fakeServer)

	err := pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		callCount++

		expected := []int{1, 2, 3}
		actual, err := ExtractSingleInts(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, expected, actual)
		return true, nil
	})
	th.CheckNoErr(t, err)
	th.CheckEquals(t, 1, callCount)
}

func TestAllPagesSingle(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	pager := setupSinglePaged(fakeServer)

	page, err := pager.AllPages(context.TODO())
	th.AssertNoErr(t, err)

	expected := []int{1, 2, 3}
	actual, err := ExtractSingleInts(page)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}
