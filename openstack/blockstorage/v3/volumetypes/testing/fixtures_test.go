package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func MockListResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `
{
    "volume_types": [
        {
            "name": "SSD",
            "qos_specs_id": null,
            "os-volume-type-access:is_public": true,
            "extra_specs": {
              "volume_backend_name": "lvmdriver-1"
            },
            "is_public": true,
            "id": "6685584b-1eac-4da6-b5c3-555430cf68ff",
            "description": null
        },
        {
            "name": "SATA",
            "qos_specs_id": null,
            "os-volume-type-access:is_public": true,
            "extra_specs": {
                "volume_backend_name": "lvmdriver-1"
            },
            "is_public": true,
            "id": "8eb69a46-df97-4e41-9586-9a40a7533803",
            "description": null
        }
    ],
	"volume_type_links": [
        {
            "href": "%s/types?marker=1",
            "rel": "next"
        }
    ]
}
  `, fakeServer.Server.URL)
		case "1":
			fmt.Fprint(w, `{"volume_types": []}`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

func MockGetResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
{
    "volume_type": {
        "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
        "name": "vol-type-001",
        "os-volume-type-access:is_public": true,
        "qos_specs_id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
        "description": "volume type 001",
        "is_public": true,
        "extra_specs": {
            "capabilities": "gpu"
        }
    }
}
`)
	})
}

func MockCreateResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "volume_type": {
        "name": "test_type",
        "os-volume-type-access:is_public": true,
        "description": "test_type_desc",
        "extra_specs": {
            "capabilities": "gpu"
        }
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "volume_type": {
        "name": "test_type",
        "extra_specs": {},
        "is_public": true,
        "os-volume-type-access:is_public": true,
        "id": "6d0ff92a-0007-4780-9ece-acfe5876966a",
        "description": "test_type_desc",
        "extra_specs": {
            "capabilities": "gpu"
        }
    }
}
    `)
	})
}

func MockDeleteResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

func MockUpdateResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
{
    "volume_type": {
        "name": "vol-type-002",
        "description": "volume type 0001",
        "is_public": true,
        "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22"
    }
}`)
	})
}

// ExtraSpecsGetBody provides a GET result of the extra_specs for a volume type
const ExtraSpecsGetBody = `
{
    "extra_specs" : {
        "capabilities": "gpu",
        "volume_backend_name": "ssd"
    }
}
`

// GetExtraSpecBody provides a GET result of a particular extra_spec for a volume type
const GetExtraSpecBody = `
{
    "capabilities": "gpu"
}
`

// UpdatedExtraSpecBody provides an PUT result of a particular updated extra_spec for a volume type
const UpdatedExtraSpecBody = `
{
    "capabilities": "gpu-2"
}
`

// ExtraSpecs is the expected extra_specs returned from GET on a volume type's extra_specs
var ExtraSpecs = map[string]string{
	"capabilities":        "gpu",
	"volume_backend_name": "ssd",
}

// ExtraSpec is the expected extra_spec returned from GET on a volume type's extra_specs
var ExtraSpec = map[string]string{
	"capabilities": "gpu",
}

// UpdatedExtraSpec is the expected extra_spec returned from PUT on a volume type's extra_specs
var UpdatedExtraSpec = map[string]string{
	"capabilities": "gpu-2",
}

func HandleListIsPublicParam(t *testing.T, fakeServer th.FakeServer, values map[string]string) {
	fakeServer.Mux.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestFormValues(t, r, values)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"volume_types": []}`)
	})
}

func HandleListWithNameFilter(t *testing.T, fakeServer th.FakeServer, values map[string]string) {
	fakeServer.Mux.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestFormValues(t, r, values)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
{
    "volume_types": [
        {
            "name": "test-type",
            "qos_specs_id": null,
            "os-volume-type-access:is_public": true,
            "extra_specs": {
                "storage_protocol": "nfs"
            },
            "is_public": true,
            "id": "996af3df-92fd-4814-a0ee-ba5f899aa1ec",
            "description": "test"
        }
    ]
}
`)
	})
}

func HandleListWithDescriptionFilter(t *testing.T, fakeServer th.FakeServer, values map[string]string) {
	fakeServer.Mux.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestFormValues(t, r, values)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `
{
    "volume_types": [
        {
            "name": "test-type",
            "qos_specs_id": null,
            "os-volume-type-access:is_public": true,
            "extra_specs": {
                "multiattach": "<is> True"
            },
            "is_public": true,
            "id": "ab948f0a-13ed-47c8-b9be-cade0beb0706",
            "description": "test"
        }
    ]
}
`)
	})
}

func HandleListWithExtraSpecsFilter(t *testing.T, fakeServer th.FakeServer, values map[string]string) {
	fakeServer.Mux.HandleFunc("/types", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestFormValues(t, r, values)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		extraSpecsFilter := r.Form.Get("extra_specs")

		// Determine which volume type to return based on extra_specs filter
		switch extraSpecsFilter {
		case "{'storage_protocol':'nfs'}":
			// Return only nfs-type
			fmt.Fprint(w, `
{
    "volume_types": [
        {
            "name": "nfs-type",
            "qos_specs_id": null,
            "os-volume-type-access:is_public": true,
            "extra_specs": {
                "storage_protocol": "nfs"
            },
            "is_public": true,
            "id": "6b0cfee7-48b6-41b7-9d68-0d74cbdc08de",
            "description": "NFS storage type"
        }
    ]
}
`)
		case "{'multiattach':'<is> True', 'replication_enabled':'<is> True', 'RESKEY:availability_zones':'zone'}":
			// Return only multiattach-type
			fmt.Fprint(w, `
{
    "volume_types": [
        {
            "name": "multiattach-type",
            "qos_specs_id": null,
            "os-volume-type-access:is_public": true,
            "extra_specs": {
                "multiattach": "<is> True",
                "replication_enabled": "<is> True",
                "RESKEY:availability_zones": "zone"
            },
            "is_public": true,
            "id": "e1fc0553-0357-4206-af30-23137ee5f743",
            "description": "Multiattach volume type"
        }
    ]
}
`)
		default:
			// Default: return empty list
			fmt.Fprint(w, `{"volume_types": []}`)
		}
	})
}

func HandleExtraSpecsListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/1/extra_specs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, ExtraSpecsGetBody)
	})
}

func HandleExtraSpecGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/1/extra_specs/capabilities", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, GetExtraSpecBody)
	})
}

func HandleExtraSpecsCreateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/1/extra_specs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `{
				"extra_specs": {
					"capabilities":        "gpu",
                    "volume_backend_name": "ssd"
				}
			}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, ExtraSpecsGetBody)
	})
}

func HandleExtraSpecUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/1/extra_specs/capabilities", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `{
				"capabilities":        "gpu-2"
			}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, UpdatedExtraSpecBody)
	})
}

func HandleExtraSpecDeleteSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/1/extra_specs/capabilities", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})
}

func MockEncryptionCreateResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/a5082c24-2a27-43a4-b48e-fcec1240e36b/encryption", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "encryption": {
        "key_size": 256,
        "provider": "luks",
        "control_location": "front-end",
        "cipher": "aes-xts-plain64"
   }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "encryption": {
        "volume_type_id": "a5082c24-2a27-43a4-b48e-fcec1240e36b",
        "control_location": "front-end",
        "encryption_id": "81e069c6-7394-4856-8df7-3b237ca61f74",
        "key_size": 256,
        "provider": "luks",
        "cipher": "aes-xts-plain64"
    }
}
    `)
	})
}

func MockDeleteEncryptionResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/a5082c24-2a27-43a4-b48e-fcec1240e36b/encryption/81e069c6-7394-4856-8df7-3b237ca61f74", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

func MockEncryptionUpdateResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/a5082c24-2a27-43a4-b48e-fcec1240e36b/encryption/81e069c6-7394-4856-8df7-3b237ca61f74", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "encryption": {
        "key_size": 256,
        "provider": "luks",
        "control_location": "front-end",
        "cipher": "aes-xts-plain64"
   }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "encryption": {
        "control_location": "front-end",
        "key_size": 256,
        "provider": "luks",
        "cipher": "aes-xts-plain64"
    }
}
    `)
	})
}

func MockEncryptionGetResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/a5082c24-2a27-43a4-b48e-fcec1240e36b/encryption", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "volume_type_id": "a5082c24-2a27-43a4-b48e-fcec1240e36b",
    "control_location": "front-end",
    "deleted": false,
    "created_at": "2016-12-28T02:32:25.000000",
    "updated_at": null,
    "encryption_id": "81e069c6-7394-4856-8df7-3b237ca61f74",
    "key_size": 256,
    "provider": "luks",
    "deleted_at": null,
    "cipher": "aes-xts-plain64"
}
    `)
	})
}

func MockEncryptionGetSpecResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/types/a5082c24-2a27-43a4-b48e-fcec1240e36b/encryption/cipher", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `
{
    "cipher": "aes-xts-plain64"
}
    `)
	})
}
