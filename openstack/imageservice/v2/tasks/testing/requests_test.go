package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/tasks"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fakeclient "github.com/gophercloud/gophercloud/testhelper/client"
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

	tasks.List(fakeclient.ServiceClient(), tasks.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
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

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
