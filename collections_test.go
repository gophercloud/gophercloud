package gophercloud

import (
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

func (c SinglePageCollection) NextPageURL() string {
	panic("NextPageURL should never be called on a single-paged collection.")
}

func (c SinglePageCollection) Concat(other Collection) Collection {
	panic("Concat should never be called on a single-paged collection.")
}

func ExtractSingleInts(c Collection) []int {
	return c.(SinglePageCollection).results
}

func setupSinglePaged() Pager {
	return NewSinglePager(func() (Collection, error) {
		return SinglePageCollection{results: []int{1, 2, 3}}, nil
	})
}

func TestEnumerateSinglePaged(t *testing.T) {
	callCount := 0
	err := setupSinglePaged().EachPage(func(page Collection) bool {
		callCount++

		expected := []int{1, 2, 3}
		actual := AsSingleInts(page)
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %v, but was %v", expected, actual)
		}
		return true
	})
	if err != nil {
		t.Fatalf("Unexpected error calling EachPage: %v", err)
	}

	if callCount != 1 {
		t.Errorf("Callback was invoked %d times", callCount)
	}
}

func TestAllSinglePaged(t *testing.T) {
	r, err := setupSinglePaged().AllPages()
	if err != nil {
		t.Fatalf("Unexpected error when iterating pages: %v", err)
	}

	expected := []int{1, 2, 3}
	actual := ExtractSingleInts(r)
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, but was %v", expected, actual)
	}
}

// LinkedPager sample and test cases.

type LinkedCollection struct {
	PaginationLinks

	results []int
}

func (page LinkedCollection) NextPageURL() string {
	n := page.PaginationLinks.Next
	if n == nil {
		return ""
	}
	return *n
}

func (page LinkedCollection) Concat(other Collection) Collection {
	return LinkedCollection{
		service: page.service,
		results: append(c.results, AsLinkedInts(other)...),
	}
}

func AsLinkedInts(results Collection) []int {
	return results.(LinkedCollection).results
}

func createLinked() Pager {
	nextURL := testhelper.Server.URL + "/foo?page=2&perPage=3"
	return CreatePager(func(url) Collection {
		LinkedCollection{
			PaginationLinks: PaginationLinks{Next: &nextURL},
			results:         []int{1, 2, 3},
		}
	})
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

func TestAllLinked(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	setupLinkedResponses(t)
	lc := createLinked()

	all, err := AllPages(lc)
	if err != nil {
		t.Fatalf("Unexpected error collection all linked pages: %v", err)
	}

	actual := AsLinkedInts(all)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, but was %v", expected, actual)
	}

	original := []int{1, 2, 3}
	if !reflect.DeepEqual(AsLinkedInts(lc), original) {
		t.Errorf("AllPages modified the original page, and now it contains: %v", lc)
	}
}
