package notifications_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/instanceha/v1/notifications"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const fixture = `{"id":7,"notification_uuid":"notification","source_host_uuid":"host","failover_segment_uuid":null,"type":"COMPUTE_HOST","status":"running","generated_time":"2026-10-25T12:00:00.123456","created_at":"2026-10-25T12:00:01","updated_at":null,"payload":{"event":"STOPPED"},"message":null}`

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
		if req.URL.Path == "/notifications/resource" {
			fmt.Fprint(w, `{"notification":`+fixture+`}`)
			return
		}
		check(t, req.URL.Path, "/notifications")
		calls++
		q := req.URL.Query()
		check(t, q.Get("limit"), "1")
		check(t, q.Get("sort_key"), "id")
		check(t, q.Get("sort_dir"), "asc")
		check(t, q.Get("type"), "COMPUTE_HOST")
		check(t, q.Get("status"), "running")
		check(t, q.Get("source_host_uuid"), "host")
		check(t, q.Get("generated-since"), "2026-10-25T12:00:00Z")
		switch q.Get("marker") {
		case "":
			fmt.Fprint(w, `{"notifications":[`+fixture+`]}`)
		case "7":
			fmt.Fprint(w, `{"notifications":[]}`)
		default:
			t.Errorf("unexpected marker %q", q.Get("marker"))
		}
	})
	client := &gophercloud.ServiceClient{ProviderClient: &gophercloud.ProviderClient{}, Endpoint: server.Endpoint(), Type: "instance-ha", Microversion: "1.3"}

	opts := notifications.ListOpts{Limit: 1, SortKey: "id", SortDir: "asc", Type: "COMPUTE_HOST", Status: "running", SourceHostUUID: "host", GeneratedSince: "2026-10-25T12:00:00Z"}
	pages, err := notifications.List(client, opts).AllPages(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	items, err := notifications.ExtractNotifications(pages)
	if err != nil {
		t.Fatal(err)
	}
	var expected notifications.Notification
	if err := json.Unmarshal([]byte(fixture), &expected); err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || !reflect.DeepEqual(items[0], expected) {
		t.Fatalf("unexpected resources: %#v", items)
	}
	if calls != 2 {
		t.Fatalf("expected two pages, got %d", calls)
	}
	item, err := notifications.Get(context.Background(), client, "resource").Extract()
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
		{"invalid resource", `{"notifications":[{"id":"invalid"}]}`, 200, true},
		{"missing marker", `{"notifications":[{}]}`, 200, true},
		{"repeated marker", `{"notifications":[{"id":7}]}`, 200, true},
		{"empty", `{"notifications":[]}`, 200, false},
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
			_, err := notifications.List(client, nil).AllPages(context.Background())
			if (err != nil) != tc.wantError {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
