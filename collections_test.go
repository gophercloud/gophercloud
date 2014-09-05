package gophercloud

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/rackspace/gophercloud/testhelper"
)

// SinglePage sample and test cases.

type SinglePageCollection struct {
	results []int
}

func (c SinglePageCollection) Pager() Pager {
	return SinglePager{}
}

func AsSingleInts(c Collection) []int {
	return c.(SinglePageCollection).results
}

var single = SinglePageCollection{
	results: []int{1, 2, 3},
}

func TestEnumerateSinglePaged(t *testing.T) {
	callCount := 0
	EachPage(single, func(page Collection) bool {
		callCount++

		expected := []int{1, 2, 3}
		actual := AsSingleInts(page)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %v, but was %v", expected, actual)
		}
		return true
	})

	if callCount != 1 {
		t.Errorf("Callback was invoked %d times", callCount)
	}
}

func TestAllSinglePaged(t *testing.T) {
	r, err := AllPages(single)
	if err != nil {
		t.Fatalf("Unexpected error when iterating pages: %v", err)
	}

	expected := []int{1, 2, 3}
	actual := AsSingleInts(r)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, but was %v", expected, actual)
	}
}

// LinkedPager sample and test cases.

type LinkedCollection struct {
	PaginationLinks

	service *ServiceClient
	results []int
}

func (c LinkedCollection) Links() PaginationLinks {
	return c.PaginationLinks
}

func (c LinkedCollection) Service() *ServiceClient {
	return c.service
}

func (c LinkedCollection) Interpret(response interface{}) (LinkCollection, error) {
	fmt.Printf("Interpreting result: %#v\n", response)
	casted, ok := response.([]interface{})
	if ok {
		asInts := make([]int, len(casted))
		for index, item := range casted {
			f := item.(float64)
			asInts[index] = int(f)
		}

		var nextURL *string
		switch asInts[0] {
		case 4:
			u := testhelper.Server.URL + "/foo?page=3&perPage=3"
			nextURL = &u
		case 7:
			// Leave nextURL as nil.
		default:
			return nil, fmt.Errorf("Unexpected resultset: %#v", asInts)
		}

		result := LinkedCollection{
			PaginationLinks: PaginationLinks{Next: nextURL},
			service:         c.service,
			results:         asInts,
		}
		if nextURL != nil {
			fmt.Printf("Returning result: %s\n", *nextURL)
		} else {
			fmt.Printf("No next link")
		}
		return result, nil
	}
	return nil, errors.New("Wat")
}

func (c LinkedCollection) Pager() Pager {
	return NewLinkPager(c)
}

func AsLinkedInts(results Collection) []int {
	return results.(LinkedCollection).results
}

func createLinked() LinkedCollection {
	nextURL := testhelper.Server.URL + "/foo?page=2&perPage=3"
	return LinkedCollection{
		PaginationLinks: PaginationLinks{Next: &nextURL},
		service: &ServiceClient{
			Provider: &ProviderClient{TokenID: "1234"},
			Endpoint: testhelper.Endpoint(),
		},
		results: []int{1, 2, 3},
	}
}

func setupLinkedResponses(t *testing.T) {
	testhelper.Mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", "1234")
		w.Header().Add("Content-Type", "application/json")

		r.ParseForm()

		pages := r.Form["page"]
		if len(pages) != 1 {
			t.Errorf("Endpoint called with unexpected page: %#v", r.Form)
		}

		switch pages[0] {
		case "2":
			fmt.Fprintf(w, `[4, 5, 6]`)
		case "3":
			fmt.Fprintf(w, `[7, 8, 9]`)
		default:
			t.Errorf("Endpoint called with unexpected page: %s", pages[0])
		}
	})
}

func TestEnumerateLinked(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	setupLinkedResponses(t)
	lc := createLinked()

	callCount := 0
	err := EachPage(lc, func(page Collection) bool {
		actual := AsLinkedInts(page)
		t.Logf("Handler invoked with %v", actual)

		var expected []int
		switch callCount {
		case 0:
			expected = []int{1, 2, 3}
		case 1:
			expected = []int{4, 5, 6}
		case 2:
			expected = []int{7, 8, 9}
		default:
			t.Fatalf("Unexpected call count: %d", callCount)
			return false
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Call %d: Expected %#v, but was %#v", callCount, expected, actual)
		}

		callCount++
		return true
	})
	if err != nil {
		t.Errorf("Unexpected error for page iteration: %v", err)
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, but was %d", callCount)
	}
}
