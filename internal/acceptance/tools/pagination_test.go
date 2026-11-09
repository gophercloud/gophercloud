package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func keypairIDs(page pagination.Page) ([]string, error) {
	items, err := keypairs.ExtractKeyPairs(page)
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.Name
	}
	return ids, err
}

func paginationFixture(t *testing.T, pages [][]string, cycle bool) (*gophercloud.ServiceClient, *atomic.Int64) {
	t.Helper()
	fakeServer := th.SetupHTTP()
	t.Cleanup(fakeServer.Teardown)
	requests := new(atomic.Int64)
	fakeServer.Mux.HandleFunc("/os-keypairs", func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		index := 0
		if marker := r.URL.Query().Get("marker"); marker != "" {
			var err error
			index, err = strconv.Atoi(marker)
			th.AssertNoErr(t, err)
		}
		var items []map[string]any
		if index < len(pages) {
			for _, id := range pages[index] {
				items = append(items, map[string]any{"keypair": map[string]string{"name": id}})
			}
		}
		links := []gophercloud.Link{}
		if cycle || index+1 < len(pages) {
			next := index + 1
			if cycle {
				next = 0
			}
			links = append(links, gophercloud.Link{Rel: "next", Href: fmt.Sprintf("%sos-keypairs?marker=%d", fakeServer.Endpoint(), next)})
		}
		w.Header().Set("Content-Type", "application/json")
		th.CheckNoErr(t, json.NewEncoder(w).Encode(map[string]any{"keypairs": items, "keypairs_links": links}))
	})
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       fakeServer.Endpoint(),
	}, requests
}

func TestCheckPagination(t *testing.T) {
	client, requests := paginationFixture(t, [][]string{{"other"}, {"first"}, {"second"}}, false)
	CheckPagination(t, context.Background(), func() pagination.Pager {
		return keypairs.List(client, nil)
	}, keypairIDs, []string{"first", "second"})
	th.AssertEquals(t, int64(6), requests.Load())
}

func TestPaginationContract(t *testing.T) {
	for _, tc := range []struct {
		name     string
		pages    [][]string
		expected []string
		err      string
	}{
		{name: "complete", pages: [][]string{{"first"}, {"second"}}, expected: []string{"first", "second"}},
		{name: "missing second page", pages: [][]string{{"first"}}, expected: []string{"first", "second"}, err: "omitted expected ID"},
		{name: "repeated expected ID", pages: [][]string{{"first"}, {"first"}}, expected: []string{"first", "second"}, err: "repeated ID"},
		{name: "repeated unrelated ID", pages: [][]string{{"other"}, {"other"}, {"first"}, {"second"}}, expected: []string{"first", "second"}, err: "repeated ID"},
		{name: "empty ID", pages: [][]string{{""}, {"second"}}, expected: []string{"first", "second"}, err: "empty ID"},
		{name: "empty collection", expected: []string{"first", "second"}, err: "omitted expected ID"},
		{name: "one expected ID", expected: []string{"first"}, err: "at least two"},
		{name: "duplicate expectations", expected: []string{"first", "first"}, err: "nonempty and distinct"},
		{name: "empty expectation", expected: []string{"first", ""}, err: "nonempty and distinct"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client, _ := paginationFixture(t, tc.pages, false)
			for _, allPages := range []bool{false, true} {
				t.Run(fmt.Sprintf("AllPages=%t", allPages), func(t *testing.T) {
					err := checkPagination(context.Background(), keypairs.List(client, nil), keypairIDs, tc.expected, allPages)
					if tc.err == "" {
						th.AssertNoErr(t, err)
					} else {
						th.AssertErr(t, err)
						th.AssertTrue(t, strings.Contains(err.Error(), tc.err))
					}
				})
			}
		})
	}
}

func TestPaginationRejectsOversizedPage(t *testing.T) {
	client, _ := paginationFixture(t, [][]string{{"first", "second"}}, false)
	err := checkPagination(context.Background(), keypairs.List(client, nil), keypairIDs, []string{"first", "second"}, false)
	th.AssertErr(t, err)
	th.AssertTrue(t, strings.Contains(err.Error(), "expected page size 1, got 2"))
}

func TestPaginationPropagatesErrors(t *testing.T) {
	client, _ := paginationFixture(t, [][]string{{"first"}, {"second"}}, false)
	failure := errors.New("extraction failed")
	for _, allPages := range []bool{false, true} {
		err := checkPagination(context.Background(), keypairs.List(client, nil), func(pagination.Page) ([]string, error) {
			return nil, failure
		}, []string{"first", "second"}, allPages)
		th.AssertErrIs(t, err, failure)

		err = checkPagination(context.Background(), pagination.Pager{Err: failure}, keypairIDs, []string{"first", "second"}, allPages)
		th.AssertErrIs(t, err, failure)
	}
}

func TestPaginationCycleIsBounded(t *testing.T) {
	client, _ := paginationFixture(t, [][]string{{"first"}}, true)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := checkPagination(ctx, keypairs.List(client, nil), keypairIDs, []string{"first", "second"}, false)
	th.AssertErr(t, err)
	th.AssertTrue(t, strings.Contains(err.Error(), "repeated ID"))

	ctx, cancel = context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()
	err = checkPagination(ctx, keypairs.List(client, nil), keypairIDs, []string{"first", "second"}, true)
	th.AssertErrIs(t, err, context.DeadlineExceeded)
}
