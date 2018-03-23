package testing

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/nodes"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreateNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCreateSuccessfully(t)

	createOpts := nodes.CreateOpts{
		ClusterID: "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
		Metadata: map[string]interface{}{
			"foo": "bar",
			"test": map[string]interface{}{
				"nil_interface": interface{}(nil),
				"float_value":   float64(123.3),
				"string_value":  "test_string",
				"bool_value":    false,
			},
		},
		Name:      "node-e395be1e-002",
		ProfileID: "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		Role:      "",
	}

	res := nodes.Create(fake.ServiceClient(), createOpts)
	th.AssertNoErr(t, res.Err)

	requestID := res.Header.Get("X-Openstack-Request-Id")
	th.AssertEquals(t, "req-3791a089-9d46-4671-a3f9-55e95e55d2b4", requestID)

	location := res.Header.Get("Location")
	th.AssertEquals(t, "http://senlin.cloud.blizzard.net:8778/v1/actions/ffd94dd8-6266-4887-9a8c-5b78b72136da", location)

	locationFields := strings.Split(location, "actions/")
	actionID := locationFields[1]
	th.AssertEquals(t, "ffd94dd8-6266-4887-9a8c-5b78b72136da", actionID)

	actual, err := res.Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedCreate, *actual)
}

func TestListNodes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListSuccessfully(t)

	count := 0
	err := nodes.List(fake.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		actual, err := nodes.ExtractNodes(page)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, ExpectedList, actual)
		count++
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, count, 1)
}

func TestDeleteNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleDeleteSuccessfully(t)

	deleteResult := nodes.Delete(fake.ServiceClient(), "6dc6d336e3fc4c0a951b5698cd1236ee")
	th.AssertNoErr(t, deleteResult.ExtractErr())
}

func TestGetNode(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes/573aa1ba-bf45-49fd-907d-6b5d6e6adfd3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": "2016-05-13T07:02:20Z",
				"data": {},
				"dependents": {},
				"domain": null,
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": "2016-05-13T07:02:04Z",
				"metadata": {"foo": "bar"},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": "2016-05-13T07:02:20Z",
				"user": "ab79b9647d074e46ac223a8fa297b846"
			}
		}`)
	})

	initAt, _ := time.Parse(time.RFC3339, "2016-05-13T07:02:04Z")
	createdAt, _ := time.Parse(time.RFC3339, "2016-05-13T07:02:20Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2016-05-13T07:02:20Z")
	expected := nodes.Node{
		ClusterID:    "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
		CreatedAt:    createdAt,
		Data:         map[string]interface{}{},
		Dependents:   map[string]interface{}{},
		Domain:       "",
		ID:           "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
		Index:        2,
		InitAt:       initAt,
		Metadata:     map[string]interface{}{"foo": "bar"},
		Name:         "node-e395be1e-002",
		PhysicalID:   "66a81d68-bf48-4af5-897b-a3bfef7279a8",
		ProfileID:    "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
		ProfileName:  "pcirros",
		ProjectID:    "eee0b7c083e84501bdd50fb269d2a10e",
		Role:         "",
		Status:       "ACTIVE",
		StatusReason: "Creation succeeded",
		UpdatedAt:    updatedAt,
		User:         "ab79b9647d074e46ac223a8fa297b846",
	}

	actual, err := nodes.Get(fake.ServiceClient(), "573aa1ba-bf45-49fd-907d-6b5d6e6adfd3").Extract()
	if err != nil {
		t.Errorf("Failed Get nodes. %v", err)
	} else {
		th.AssertDeepEquals(t, expected, *actual)
	}
}

func TestGetNodeInvalidTimeFloat(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes/573aa1ba-bf45-49fd-907d-6b5d6e6adfd3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": 123456789.0,
				"data": {},
				"dependents": {},
				"domain": null,
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": 123456789.0,
				"metadata": {"foo": "bar"},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": 123456789.0,
				"user": "ab79b9647d074e46ac223a8fa297b846"
			}
		}`)
	})

	_, err := nodes.Get(fake.ServiceClient(), "573aa1ba-bf45-49fd-907d-6b5d6e6adfd3").Extract()
	th.AssertEquals(t, false, err == nil)
}

func TestGetNodeInvalidTimeString(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/nodes/573aa1ba-bf45-49fd-907d-6b5d6e6adfd3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
		{
			"node": {
				"cluster_id": "e395be1e-8d8e-43bb-bd6c-943eccf76a6d",
				"created_at": "invalid",
				"data": {},
				"dependents": {},
				"domain": null,
				"id": "82fe28e0-9fcb-42ca-a2fa-6eb7dddd75a1",
				"index": 2,
				"init_at": "invalid",
				"metadata": {"foo": "bar"},
				"name": "node-e395be1e-002",
				"physical_id": "66a81d68-bf48-4af5-897b-a3bfef7279a8",
				"profile_id": "d8a48377-f6a3-4af4-bbbb-6e8bcaa0cbc0",
				"profile_name": "pcirros",
				"project_id": "eee0b7c083e84501bdd50fb269d2a10e",
				"role": "",
				"status": "ACTIVE",
				"status_reason": "Creation succeeded",
				"updated_at": "invalid",
				"user": "ab79b9647d074e46ac223a8fa297b846"
			}
		}`)
	})

	_, err := nodes.Get(fake.ServiceClient(), "573aa1ba-bf45-49fd-907d-6b5d6e6adfd3").Extract()
	th.AssertEquals(t, false, err == nil)
}
