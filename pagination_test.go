package gophercloud

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/testhelper"
)

func createClient() *ServiceClient {
	return &ServiceClient{
		Provider: &ProviderClient{TokenID: "abc123"},
		Endpoint: testhelper.Endpoint(),
	}
}

// SinglePage sample and test cases.

type SinglePageResult struct {
	SinglePageBase
}

func (r SinglePageResult) IsEmpty() (bool, error) {
	is, err := ExtractSingleInts(r)
	if err != nil {
		return true, err
	}
	return len(is) == 0, nil
}

func ExtractSingleInts(page Page) ([]int, error) {
	var response struct {
		Ints []int `mapstructure:"ints"`
	}

	err := mapstructure.Decode(page.(SinglePageResult).Body, &response)
	if err != nil {
		return nil, err
	}

	return response.Ints, nil
}

func setupSinglePaged() Pager {
	testhelper.SetupHTTP()
	client := createClient()

	testhelper.Mux.HandleFunc("/only", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{ "ints": [1, 2, 3] }`)
	})

	createPage := func(r LastHTTPResponse) Page {
		return SinglePageResult{SinglePageBase(r)}
	}

	return NewSinglePager(client, testhelper.Server.URL+"/only", createPage)
}

func TestEnumerateSinglePaged(t *testing.T) {
	callCount := 0
	pager := setupSinglePaged()
	defer testhelper.TeardownHTTP()

	err := pager.EachPage(func(page Page) (bool, error) {
		callCount++

		expected := []int{1, 2, 3}
		actual, err := ExtractSingleInts(page)
		testhelper.AssertNoErr(t, err)
		testhelper.CheckDeepEquals(t, expected, actual)
		return true, nil
	})
	testhelper.CheckNoErr(t, err)
	testhelper.CheckEquals(t, 1, callCount)
}

// LinkedPager sample and test cases.

type LinkedPageResult struct {
	LinkedPageBase
}

func (r LinkedPageResult) IsEmpty() (bool, error) {
	is, err := ExtractLinkedInts(r)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

func ExtractLinkedInts(page Page) ([]int, error) {
	var response struct {
		Ints []int `mapstructure:"ints"`
	}

	err := mapstructure.Decode(page.(LinkedPageResult).Body, &response)
	if err != nil {
		return nil, err
	}

	return response.Ints, nil
}

func createLinked(t *testing.T) Pager {
	testhelper.SetupHTTP()

	testhelper.Mux.HandleFunc("/page1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{ "ints": [1, 2, 3], "links": { "next": "%s/page2" } }`, testhelper.Server.URL)
	})

	testhelper.Mux.HandleFunc("/page2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{ "ints": [4, 5, 6], "links": { "next": "%s/page3" } }`, testhelper.Server.URL)
	})

	testhelper.Mux.HandleFunc("/page3", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{ "ints": [7, 8, 9], "links": { "next": null } }`)
	})

	client := createClient()

	createPage := func(r LastHTTPResponse) Page {
		return LinkedPageResult{LinkedPageBase(r)}
	}

	return NewLinkedPager(client, testhelper.Server.URL+"/page1", createPage)
}

func TestEnumerateLinked(t *testing.T) {
	pager := createLinked(t)
	defer testhelper.TeardownHTTP()

	callCount := 0
	err := pager.EachPage(func(page Page) (bool, error) {
		actual, err := ExtractLinkedInts(page)
		if err != nil {
			return false, err
		}

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
			return false, nil
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Call %d: Expected %#v, but was %#v", callCount, expected, actual)
		}

		callCount++
		return true, nil
	})
	if err != nil {
		t.Errorf("Unexpected error for page iteration: %v", err)
	}

	if callCount != 3 {
		t.Errorf("Expected 3 calls, but was %d", callCount)
	}
}

// MarkerPager sample and test cases.

type MarkerPageResult struct {
	MarkerPageBase
}

func (r MarkerPageResult) IsEmpty() (bool, error) {
	results, err := ExtractMarkerStrings(r)
	if err != nil {
		return true, err
	}
	return len(results) == 0, err
}

func (r MarkerPageResult) LastMark() (string, error) {
	results, err := ExtractMarkerStrings(r)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", nil
	}
	return results[len(results)-1], nil
}

func createMarkerPaged(t *testing.T) Pager {
	testhelper.SetupHTTP()

	testhelper.Mux.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ms := r.Form["marker"]
		switch {
		case len(ms) == 0:
			fmt.Fprintf(w, "aaa\nbbb\nccc")
		case len(ms) == 1 && ms[0] == "ccc":
			fmt.Fprintf(w, "ddd\neee\nfff")
		case len(ms) == 1 && ms[0] == "fff":
			fmt.Fprintf(w, "ggg\nhhh\niii")
		case len(ms) == 1 && ms[0] == "iii":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("Request with unexpected marker: [%v]", ms)
		}
	})

	client := createClient()

	createPage := func(r LastHTTPResponse) MarkerPage {
		p := MarkerPageResult{MarkerPageBase{LastHTTPResponse: r}}
		p.MarkerPageBase.Self = p
		return p
	}

	return NewMarkerPager(client, testhelper.Server.URL+"/page", createPage)
}

func ExtractMarkerStrings(page Page) ([]string, error) {
	content := page.(MarkerPageResult).Body.([]uint8)
	parts := strings.Split(string(content), "\n")
	results := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) > 0 {
			results = append(results, part)
		}
	}
	return results, nil
}

func TestEnumerateMarker(t *testing.T) {
	pager := createMarkerPaged(t)
	defer testhelper.TeardownHTTP()

	callCount := 0
	err := pager.EachPage(func(page Page) (bool, error) {
		actual, err := ExtractMarkerStrings(page)
		if err != nil {
			return false, err
		}

		t.Logf("Handler invoked with %v", actual)

		var expected []string
		switch callCount {
		case 0:
			expected = []string{"aaa", "bbb", "ccc"}
		case 1:
			expected = []string{"ddd", "eee", "fff"}
		case 2:
			expected = []string{"ggg", "hhh", "iii"}
		default:
			t.Fatalf("Unexpected call count: %d", callCount)
			return false, nil
		}

		testhelper.CheckDeepEquals(t, expected, actual)

		callCount++
		return true, nil
	})
	testhelper.AssertNoErr(t, err)
	testhelper.AssertEquals(t, callCount, 3)
}
