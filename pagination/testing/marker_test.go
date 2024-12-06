package testing

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// MarkerPager sample and test cases.

type MarkerPageResult struct {
	pagination.MarkerPageBase
}

func (r MarkerPageResult) IsEmpty() (bool, error) {
	results, err := ExtractMarkerStrings(r)
	if err != nil {
		return true, err
	}
	return len(results) == 0, err
}

func (r MarkerPageResult) LastMarker() (string, error) {
	results, err := ExtractMarkerStrings(r)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", nil
	}
	return results[len(results)-1], nil
}

func createMarkerPaged(t *testing.T) pagination.Pager {
	th.SetupHTTP()

	th.Mux.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
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

	createPage := func(r pagination.PageResult) pagination.Page {
		p := MarkerPageResult{pagination.MarkerPageBase{PageResult: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

	return pagination.NewPager(client, th.Server.URL+"/page", createPage)
}

func ExtractMarkerStrings(page pagination.Page) ([]string, error) {
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
	defer th.TeardownHTTP()

	callCount := 0
	err := pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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

		th.CheckDeepEquals(t, expected, actual)

		callCount++
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, callCount, 3)
}

func TestAllPagesMarker(t *testing.T) {
	pager := createMarkerPaged(t)
	defer th.TeardownHTTP()

	page, err := pager.AllPages(context.TODO())
	th.AssertNoErr(t, err)

	expected := []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff", "ggg", "hhh", "iii"}
	actual, err := ExtractMarkerStrings(page)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}
