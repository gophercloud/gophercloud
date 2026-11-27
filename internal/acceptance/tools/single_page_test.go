package tools

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/secgroups"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/volumeattach"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/securityservices"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCheckSinglePage(t *testing.T) {
	for _, tc := range []struct {
		name     string
		path     string
		body     string
		expected []string
		list     func(*gophercloud.ServiceClient) pagination.Pager
		extract  func(pagination.Page) ([]string, error)
	}{
		{
			name:     "security group associations",
			path:     "servers/server/os-security-groups",
			body:     `{"security_groups":[{"id":"first"},{"id":"second"}]}`,
			expected: []string{"second", "first"},
			list: func(client *gophercloud.ServiceClient) pagination.Pager {
				return secgroups.ListByServer(client, "server")
			},
			extract: func(page pagination.Page) ([]string, error) {
				items, err := secgroups.ExtractSecurityGroups(page)
				ids := make([]string, len(items))
				for i, item := range items {
					ids[i] = item.ID
				}
				return ids, err
			},
		},
		{
			name:     "volume attachments",
			path:     "servers/server/os-volume_attachments",
			body:     `{"volumeAttachments":[{"volumeId":"first"},{"volumeId":"second"}]}`,
			expected: []string{"second", "first"},
			list:     func(client *gophercloud.ServiceClient) pagination.Pager { return volumeattach.List(client, "server") },
			extract: func(page pagination.Page) ([]string, error) {
				items, err := volumeattach.ExtractVolumeAttachments(page)
				ids := make([]string, len(items))
				for i, item := range items {
					ids[i] = item.VolumeID
				}
				return ids, err
			},
		},
		{
			name:     "Manila security services",
			path:     "security-services/detail",
			body:     `{"security_services":[{"id":"first"},{"id":"second"}]}`,
			expected: []string{"second", "first"},
			list:     func(client *gophercloud.ServiceClient) pagination.Pager { return securityservices.List(client, nil) },
			extract: func(page pagination.Page) ([]string, error) {
				items, err := securityservices.ExtractSecurityServices(page)
				ids := make([]string, len(items))
				for i, item := range items {
					ids[i] = item.ID
				}
				return ids, err
			},
		},
		{
			name: "empty collection",
			path: "servers/server/os-volume_attachments",
			body: `{"volumeAttachments":[]}`,
			list: func(client *gophercloud.ServiceClient) pagination.Pager { return volumeattach.List(client, "server") },
			extract: func(page pagination.Page) ([]string, error) {
				items, err := volumeattach.ExtractVolumeAttachments(page)
				th.CheckEquals(t, 0, len(items))
				return nil, err
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			fakeServer.Mux.HandleFunc("/"+tc.path, func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, "GET")
				th.CheckEquals(t, "", r.URL.RawQuery)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, tc.body)
			})
			client := &gophercloud.ServiceClient{ProviderClient: &gophercloud.ProviderClient{}, Endpoint: fakeServer.Endpoint()}
			CheckSinglePage(t, context.Background(), client, client.ServiceURL(tc.path), func() pagination.Pager {
				return tc.list(client)
			}, tc.extract, tc.expected)
			th.AssertTrue(t, client.HTTPClient.Transport == nil)
		})
	}
}

func TestSinglePageTransportRejectsSDKContinuation(t *testing.T) {
	for _, allPages := range []bool{false, true} {
		name := "EachPage"
		if allPages {
			name = "AllPages"
		}
		t.Run(name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()
			var forwarded atomic.Int64
			fakeServer.Mux.HandleFunc("/os-security-groups", func(w http.ResponseWriter, r *http.Request) {
				forwarded.Add(1)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"security_groups":[{"id":"only"}]}`)
			})
			client := &gophercloud.ServiceClient{ProviderClient: &gophercloud.ProviderClient{}, Endpoint: fakeServer.Endpoint()}
			collection, err := url.Parse(client.ServiceURL("os-security-groups"))
			th.AssertNoErr(t, err)
			rt := &singlePageTransport{collection: *collection, next: http.DefaultTransport}
			client.HTTPClient.Transport = rt
			// The top-level list deliberately supplies a real offset pager. The
			// guard must reject its continuation even if the response were empty.
			pager := secgroups.List(client)
			if allPages {
				_, err = pager.AllPages(context.Background())
			} else {
				err = pager.EachPage(context.Background(), func(context.Context, pagination.Page) (bool, error) {
					return true, nil
				})
			}
			th.AssertErr(t, err)
			th.AssertTrue(t, strings.Contains(err.Error(), "unexpected continuation request"))
			th.AssertEquals(t, int64(1), forwarded.Load())
			th.AssertEquals(t, int64(2), rt.requests.Load())
		})
	}
}

func TestSinglePageTransportRejectsContinuation(t *testing.T) {
	collection, err := url.Parse("https://example.com/items")
	th.AssertNoErr(t, err)
	forwarded := 0
	rt := &singlePageTransport{collection: *collection, next: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		forwarded++
		return &http.Response{StatusCode: http.StatusOK, Request: req}, nil
	})}
	for _, rawURL := range []string{collection.String(), collection.String() + "?offset=2"} {
		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		th.AssertNoErr(t, err)
		_, err = rt.RoundTrip(req)
		if rawURL == collection.String() {
			th.AssertNoErr(t, err)
		} else {
			th.AssertErr(t, err)
		}
	}
	th.AssertEquals(t, 1, forwarded)
	th.AssertEquals(t, int64(2), rt.requests.Load())
	for _, rawURL := range []string{"https://other.example.com/items", "https://example.com/items/child"} {
		req, err := http.NewRequest(http.MethodGet, rawURL, nil)
		th.AssertNoErr(t, err)
		_, err = rt.RoundTrip(req)
		th.AssertNoErr(t, err)
	}
	req, err := http.NewRequest(http.MethodPost, collection.String(), nil)
	th.AssertNoErr(t, err)
	_, err = rt.RoundTrip(req)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 4, forwarded)
	th.AssertEquals(t, int64(2), rt.requests.Load())
}
