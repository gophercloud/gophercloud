package hosts_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/instanceha/v1/hosts"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const fixture = `{"id":7,"uuid":"host","name":"compute-1","type":"COMPUTE_HOST","control_attributes":"SSH","reserved":true,"on_maintenance":false,"failover_segment_id":"segment"}`

func check(t *testing.T, actual, expected string) {
	t.Helper()
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestListAndGet(t *testing.T) {
	calls := 0
	server := th.SetupHTTP()
	defer server.Teardown()
	server.Mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		check(t, req.Method, "GET")
		check(t, req.Header.Get("OpenStack-API-Version"), "instance-ha 1.3")
		w.Header().Set("Content-Type", "application/json")
		if req.URL.Path == "/segments/segment/hosts/resource" {
			fmt.Fprint(w, `{"host":`+fixture+`}`)
			return
		}
		check(t, req.URL.Path, "/segments/segment/hosts")
		calls++
		q := req.URL.Query()
		check(t, q.Get("limit"), "1")
		check(t, q.Get("sort_key"), "id")
		check(t, q.Get("sort_dir"), "asc")
		check(t, q.Get("name"), "compute-1")
		check(t, q.Get("type"), "COMPUTE_HOST")
		check(t, q.Get("control_attributes"), "SSH")
		check(t, q.Get("reserved"), "false")
		check(t, q.Get("on_maintenance"), "false")
		switch q.Get("marker") {
		case "":
			fmt.Fprint(w, `{"hosts":[`+fixture+`]}`)
		case "7":
			fmt.Fprint(w, `{"hosts":[]}`)
		default:
			t.Errorf("unexpected marker %q", q.Get("marker"))
		}
	})
	client := &gophercloud.ServiceClient{ProviderClient: &gophercloud.ProviderClient{}, Endpoint: server.Endpoint(), Type: "instance-ha", Microversion: "1.3"}
	disabled := false
	opts := hosts.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc", Name: "compute-1", Type: "COMPUTE_HOST", ControlAttributes: "SSH", Reserved: &disabled, OnMaintenance: &disabled}
	pages, err := hosts.List(client, "segment", opts).AllPages(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	items, err := hosts.ExtractHosts(pages)
	if err != nil {
		t.Fatal(err)
	}
	var expected hosts.Host
	if err := json.Unmarshal([]byte(fixture), &expected); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || !reflect.DeepEqual(items[0], expected) {
		t.Fatalf("unexpected resources: %#v", items)
	}
	if calls != 2 {
		t.Fatalf("expected two pages, got %d", calls)
	}
	item, err := hosts.Get(context.Background(), client, "segment", "resource").Extract()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*item, expected) {
		t.Fatalf("unexpected resource: %#v", item)
	}
}

func TestListErrors(t *testing.T) {
	for _, tc := range []struct {
		name, body string
		status     int
		wantError  bool
	}{
		{"forbidden", "{}", 403, true},
		{"invalid JSON", "{", 200, true},
		{"invalid resource", `{"hosts":[{"id":"invalid"}]}`, 200, true},
		{"missing marker", `{"hosts":[{}]}`, 200, true},
		{"repeated marker", `{"hosts":[{"id":7}]}`, 200, true},
		{"empty", `{"hosts":[]}`, 200, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := th.SetupHTTP()
			defer server.Teardown()
			server.Mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.body)
			})
			client := &gophercloud.ServiceClient{ProviderClient: &gophercloud.ProviderClient{}, Endpoint: server.Endpoint()}
			_, err := hosts.List(client, "segment", nil).AllPages(context.Background())
			if (err != nil) != tc.wantError {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
