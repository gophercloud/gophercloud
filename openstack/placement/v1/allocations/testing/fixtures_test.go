package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocations"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const ConsumerUUID = "ba8f2c8e-0bf7-4a32-aacf-7c11f7f9a321"
const EmptyConsumerUUID = "00000000-0000-0000-0000-000000000000"
const ConflictConsumerUUID = "ffffffff-ffff-ffff-ffff-ffffffffffff"
const NotFoundConsumerUUID = "11111111-1111-1111-1111-111111111111"
const ManageConsumerUUID1 = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
const ManageConsumerUUID2 = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
const ProviderUUID1 = "7d4f1abe-2f91-4f7a-8872-b70d9fb5c3dd"
const ProviderUUID2 = "f3f97e00-13e1-4c88-a7cd-db3bb4f99357"
const ProjectID = "42a2b0fa980d4f7f873e8f0d8b4e1b0e"
const UserID = "b29d880b2c114d5d9a55748b26b2e41e"

var GetAllocationsBody = fmt.Sprintf(`
{
    "allocations": {
        "%s": {
            "generation": 3,
            "resources": {
                "VCPU": 2,
                "MEMORY_MB": 2048
            }
        },
        "%s": {
            "generation": 7,
            "resources": {
                "DISK_GB": 50
            }
        }
    },
    "project_id": "%s",
    "user_id": "%s",
    "consumer_generation": 2,
    "consumer_type": "INSTANCE"
}
`, ProviderUUID1, ProviderUUID2, ProjectID, UserID)

// GetEmptyAllocationsBody is the response for a consumer UUID that has no
// allocations. Per the Placement API, only the "allocations" key is present;
// project_id, user_id, consumer_generation, and consumer_type are absent.
const GetEmptyAllocationsBody = `
{
    "allocations": {}
}
`

var consumerGeneration = 2
var consumerType = "INSTANCE"
var projectID = ProjectID
var userID = UserID

var ExpectedAllocations = allocations.Allocations{
	Allocations: map[string]allocations.ProviderAllocations{
		ProviderUUID1: {
			Generation: 3,
			Resources: map[string]int{
				"VCPU":      2,
				"MEMORY_MB": 2048,
			},
		},
		ProviderUUID2: {
			Generation: 7,
			Resources: map[string]int{
				"DISK_GB": 50,
			},
		},
	},
	ProjectID:          &projectID,
	UserID:             &userID,
	ConsumerGeneration: &consumerGeneration,
	ConsumerType:       &consumerType,
}

var ExpectedEmptyAllocations = allocations.Allocations{
	Allocations: map[string]allocations.ProviderAllocations{},
}

var UpdateAllocationsRequest = fmt.Sprintf(`
{
    "allocations": {
        "%s": {
            "resources": {
                "VCPU": 2,
                "MEMORY_MB": 2048
            }
        }
    },
    "project_id": "%s",
    "user_id": "%s",
    "consumer_generation": null
}
`, ProviderUUID1, ProjectID, UserID)

var UpdateAllocationsExistingRequest = fmt.Sprintf(`
{
    "allocations": {
        "%s": {
            "resources": {
                "VCPU": 2,
                "MEMORY_MB": 2048
            }
        }
    },
    "project_id": "%s",
    "user_id": "%s",
    "consumer_generation": 1
}
`, ProviderUUID1, ProjectID, UserID)

var GetAllocationsAfterUpdateBody = fmt.Sprintf(`
{
    "allocations": {
        "%s": {
            "generation": 4,
            "resources": {
                "VCPU": 2,
                "MEMORY_MB": 2048
            }
        }
    },
    "project_id": "%s",
    "user_id": "%s",
    "consumer_generation": 1
}
`, ProviderUUID1, ProjectID, UserID)

var consumerGenerationAfterUpdate = 1

var ExpectedAllocationsAfterUpdate = allocations.Allocations{
	Allocations: map[string]allocations.ProviderAllocations{
		ProviderUUID1: {
			Generation: 4,
			Resources: map[string]int{
				"VCPU":      2,
				"MEMORY_MB": 2048,
			},
		},
	},
	ProjectID:          &projectID,
	UserID:             &userID,
	ConsumerGeneration: &consumerGenerationAfterUpdate,
}

var ManageAllocationsRequest = fmt.Sprintf(`
{
    "%s": {
        "allocations": {
            "%s": {
                "resources": {
                    "VCPU": 1
                }
            }
        },
        "project_id": "%s",
        "user_id": "%s",
        "consumer_generation": null
    },
    "%s": {
        "allocations": {
            "%s": {
                "resources": {
                    "VCPU": 1
                }
            }
        },
        "project_id": "%s",
        "user_id": "%s",
        "consumer_generation": null
    }
}
`, ManageConsumerUUID1, ProviderUUID1, ProjectID, UserID, ManageConsumerUUID2, ProviderUUID1, ProjectID, UserID)

var GetAllocationsAfterManageBody = fmt.Sprintf(`
{
    "allocations": {
        "%s": {
            "generation": 1,
            "resources": {
                "VCPU": 1
            }
        }
    },
    "project_id": "%s",
    "user_id": "%s",
    "consumer_generation": 1
}
`, ProviderUUID1, ProjectID, UserID)

var ExpectedAllocationsAfterManage = allocations.Allocations{
	Allocations: map[string]allocations.ProviderAllocations{
		ProviderUUID1: {
			Generation: 1,
			Resources: map[string]int{
				"VCPU": 1,
			},
		},
	},
	ProjectID:          &projectID,
	UserID:             &userID,
	ConsumerGeneration: &consumerGenerationAfterUpdate,
}

func HandleGetAllocationsSuccess(t *testing.T, fakeServer th.FakeServer) {
	url := fmt.Sprintf("/allocations/%s", ConsumerUUID)

	fakeServer.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetAllocationsBody)
	})
}

func HandleGetEmptyAllocationsSuccess(t *testing.T, fakeServer th.FakeServer) {
	url := fmt.Sprintf("/allocations/%s", EmptyConsumerUUID)

	fakeServer.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetEmptyAllocationsBody)
	})
}

// HandleUpdateAndGetAllocationsSuccess handles PUT followed by GET on the same URL,
// branching by HTTP method.
func HandleUpdateAndGetAllocationsSuccess(t *testing.T, fakeServer th.FakeServer) {
	url := fmt.Sprintf("/allocations/%s", ConsumerUUID)

	fakeServer.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		switch r.Method {
		case http.MethodPut:
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, UpdateAllocationsRequest)
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, GetAllocationsAfterUpdateBody)
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	})
}

// HandleUpdateAllocationsNewConsumerSuccess handles PUT with nil consumer_generation
// (new consumer case), verifying that null is sent in the request body.
func HandleUpdateAllocationsNewConsumerSuccess(t *testing.T, fakeServer th.FakeServer) {
	url := fmt.Sprintf("/allocations/%s", EmptyConsumerUUID)

	fakeServer.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		switch r.Method {
		case http.MethodPut:
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, UpdateAllocationsRequest)
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, GetAllocationsAfterUpdateBody)
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	})
}

// HandleUpdateAllocationsConflict simulates a 409 when the consumer_generation
// in the request does not match the server's current value.
func HandleUpdateAllocationsConflict(t *testing.T, fakeServer th.FakeServer) {
	url := fmt.Sprintf("/allocations/%s", ConflictConsumerUUID)

	fakeServer.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{"errors":[{"status":409,"title":"Conflict","code":"placement.concurrent_update"}]}`)
	})
}

// HandleDeleteAndGetAllocationsSuccess handles DELETE followed by GET on the same URL.
// DELETE returns 204; the subsequent GET returns an empty allocations body.
func HandleDeleteAndGetAllocationsSuccess(t *testing.T, fakeServer th.FakeServer) {
	url := fmt.Sprintf("/allocations/%s", ConsumerUUID)

	fakeServer.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		switch r.Method {
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, GetEmptyAllocationsBody)
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	})
}

// HandleDeleteAllocationsNotFound simulates a 404 when the consumer does not exist.
func HandleDeleteAllocationsNotFound(t *testing.T, fakeServer th.FakeServer) {
	url := fmt.Sprintf("/allocations/%s", NotFoundConsumerUUID)

	fakeServer.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"errors":[{"status":404,"title":"Not Found","code":"placement.undefined_code"}]}`)
	})
}

// HandleManageAllocationsSuccess handles POST to /allocations, verifying the
// request body, then handles GET requests for each managed consumer.
func HandleManageAllocationsSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ManageAllocationsRequest)

		w.WriteHeader(http.StatusNoContent)
	})

	for _, uuid := range []string{ManageConsumerUUID1, ManageConsumerUUID2} {
		uuid := uuid
		url := fmt.Sprintf("/allocations/%s", uuid)
		fakeServer.Mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, GetAllocationsAfterManageBody)
		})
	}
}

// HandleManageAllocationsConflict simulates a 409 when a consumer_generation
// in the batch does not match the server's current value.
func HandleManageAllocationsConflict(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, `{"errors":[{"status":409,"title":"Conflict","code":"placement.concurrent_update"}]}`)
	})
}
