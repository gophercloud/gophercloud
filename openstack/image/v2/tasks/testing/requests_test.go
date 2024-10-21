package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/tasks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fakeclient "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, TasksListResult)
	})

	count := 0

	err := tasks.List(fakeclient.ServiceClient(), tasks.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := tasks.ExtractTasks(page)
		if err != nil {
			t.Errorf("Failed to extract tasks: %v", err)
			return false, nil
		}

		expected := []tasks.Task{
			Task1,
			Task2,
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/tasks/1252f636-1246-4319-bfba-c47cde0efbe0", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, TasksGetResult)
	})

	s, err := tasks.Get(context.TODO(), fakeclient.ServiceClient(), "1252f636-1246-4319-bfba-c47cde0efbe0").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Status, string(tasks.TaskStatusPending))
	th.AssertEquals(t, s.CreatedAt, time.Date(2018, 7, 25, 8, 59, 13, 0, time.UTC))
	th.AssertEquals(t, s.UpdatedAt, time.Date(2018, 7, 25, 8, 59, 14, 0, time.UTC))
	th.AssertEquals(t, s.Self, "/v2/tasks/1252f636-1246-4319-bfba-c47cde0efbe0")
	th.AssertEquals(t, s.Owner, "424e7cf0243c468ca61732ba45973b3e")
	th.AssertEquals(t, s.Message, "")
	th.AssertEquals(t, s.Type, "import")
	th.AssertEquals(t, s.ID, "1252f636-1246-4319-bfba-c47cde0efbe0")
	th.AssertEquals(t, s.Schema, "/v2/schemas/task")
	th.AssertDeepEquals(t, s.Result, map[string]any(nil))
	th.AssertDeepEquals(t, s.Input, map[string]any{
		"import_from":        "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img",
		"import_from_format": "raw",
		"image_properties": map[string]any{
			"container_format": "bare",
			"disk_format":      "raw",
		},
	})
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)
		th.TestJSONRequest(t, r, TaskCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, TaskCreateResult)
	})

	opts := tasks.CreateOpts{
		Type: "import",
		Input: map[string]any{
			"image_properties": map[string]any{
				"container_format": "bare",
				"disk_format":      "raw",
			},
			"import_from_format": "raw",
			"import_from":        "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img",
		},
	}
	s, err := tasks.Create(context.TODO(), fakeclient.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Status, string(tasks.TaskStatusPending))
	th.AssertEquals(t, s.CreatedAt, time.Date(2018, 7, 25, 11, 7, 54, 0, time.UTC))
	th.AssertEquals(t, s.UpdatedAt, time.Date(2018, 7, 25, 11, 7, 54, 0, time.UTC))
	th.AssertEquals(t, s.Self, "/v2/tasks/d550c87d-86ed-430a-9895-c7a1f5ce87e9")
	th.AssertEquals(t, s.Owner, "fb57277ef2f84a0e85b9018ec2dedbf7")
	th.AssertEquals(t, s.Message, "")
	th.AssertEquals(t, s.Type, "import")
	th.AssertEquals(t, s.ID, "d550c87d-86ed-430a-9895-c7a1f5ce87e9")
	th.AssertEquals(t, s.Schema, "/v2/schemas/task")
	th.AssertDeepEquals(t, s.Result, map[string]any(nil))
	th.AssertDeepEquals(t, s.Input, map[string]any{
		"import_from":        "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img",
		"import_from_format": "raw",
		"image_properties": map[string]any{
			"container_format": "bare",
			"disk_format":      "raw",
		},
	})
}
