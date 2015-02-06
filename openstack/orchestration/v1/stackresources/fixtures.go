package stackresources

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

var FindExpected = []Resource{
	Resource{
		Name: "hello_world",
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
				Rel:  "self",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
				Rel:  "stack",
			},
		},
		LogicalID:    "hello_world",
		StatusReason: "state changed",
		UpdatedTime:  time.Date(2015, 2, 5, 21, 33, 11, 0, time.UTC),
		RequiredBy:   []interface{}{},
		Status:       "CREATE_IN_PROGRESS",
		PhysicalID:   "49181cd6-169a-4130-9455-31185bbfc5bf",
		Type:         "OS::Nova::Server",
	},
}

const FindOutput = `
{
  "resources": [
  {
    "resource_name": "hello_world",
    "links": [
      {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
      "rel": "self"
      },
      {
        "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
        "rel": "stack"
      }
    ],
    "logical_resource_id": "hello_world",
    "resource_status_reason": "state changed",
    "updated_time": "2015-02-05T21:33:11Z",
    "required_by": [],
    "resource_status": "CREATE_IN_PROGRESS",
    "physical_resource_id": "49181cd6-169a-4130-9455-31185bbfc5bf",
    "resource_type": "OS::Nova::Server"
  }
  ]
}`

// HandleFindSuccessfully creates an HTTP handler at `/stacks/hello_world/resources`
// on the test handler mux that responds with a `Find` response.
func HandleFindSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks/hello_world/resources", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}

var ListExpected = []Resource{
	Resource{
		Name: "hello_world",
		Links: []gophercloud.Link{
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
				Rel:  "self",
			},
			gophercloud.Link{
				Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
				Rel:  "stack",
			},
		},
		LogicalID:    "hello_world",
		StatusReason: "state changed",
		UpdatedTime:  time.Date(2015, 2, 5, 21, 33, 11, 0, time.UTC),
		RequiredBy:   []interface{}{},
		Status:       "CREATE_IN_PROGRESS",
		PhysicalID:   "49181cd6-169a-4130-9455-31185bbfc5bf",
		Type:         "OS::Nova::Server",
	},
}

const ListOutput = `{
  "resources": [
  {
    "resource_name": "hello_world",
    "links": [
    {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b/resources/hello_world",
      "rel": "self"
    },
    {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/postman_stack/5f57cff9-93fc-424e-9f78-df0515e7f48b",
      "rel": "stack"
    }
    ],
    "logical_resource_id": "hello_world",
    "resource_status_reason": "state changed",
    "updated_time": "2015-02-05T21:33:11Z",
    "required_by": [],
    "resource_status": "CREATE_IN_PROGRESS",
    "physical_resource_id": "49181cd6-169a-4130-9455-31185bbfc5bf",
    "resource_type": "OS::Nova::Server"
  }
]
}`

// HandleListSuccessfully creates an HTTP handler at `/stacks/hello_world/49181cd6-169a-4130-9455-31185bbfc5bf/resources`
// on the test handler mux that responds with a `List` response.
func HandleListSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks/hello_world/49181cd6-169a-4130-9455-31185bbfc5bf/resources", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		r.ParseForm()
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, output)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

var GetExpected = &Resource{
	Name: "wordpress_instance",
	Links: []gophercloud.Link{
		gophercloud.Link{
			Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance",
			Rel:  "self",
		},
		gophercloud.Link{
			Href: "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e",
			Rel:  "stack",
		},
	},
	LogicalID:    "wordpress_instance",
	StatusReason: "state changed",
	UpdatedTime:  time.Date(2014, 12, 10, 18, 34, 35, 0, time.UTC),
	RequiredBy:   []interface{}{},
	Status:       "CREATE_COMPLETE",
	PhysicalID:   "00e3a2fe-c65d-403c-9483-4db9930dd194",
	Type:         "OS::Nova::Server",
}

const GetOutput = `
{
  "resource": {
    "resource_name": "wordpress_instance",
    "description": "",
    "links": [
    {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance",
      "rel": "self"
    },
    {
      "href": "http://166.78.160.107:8004/v1/98606384f58d4ad0b3db7d0d779549ac/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e",
      "rel": "stack"
    }
    ],
    "logical_resource_id": "wordpress_instance",
    "resource_status": "CREATE_COMPLETE",
    "updated_time": "2014-12-10T18:34:35Z",
    "required_by": [],
    "resource_status_reason": "state changed",
    "physical_resource_id": "00e3a2fe-c65d-403c-9483-4db9930dd194",
    "resource_type": "OS::Nova::Server"
  }
}`

// HandleGetSuccessfully creates an HTTP handler at `/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance`
// on the test handler mux that responds with a `Get` response.
func HandleGetSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}

var MetadataExpected = map[string]string{
	"number": "7",
	"animal": "auk",
}

const MetadataOutput = `
{
    "metadata": {
      "number": "7",
      "animal": "auk"
    }
}`

// HandleMetadataSuccessfully creates an HTTP handler at `/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance/metadata`
// on the test handler mux that responds with a `Metadata` response.
func HandleMetadataSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/stacks/teststack/0b1771bd-9336-4f2b-ae86-a80f971faf1e/resources/wordpress_instance/metadata", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}

var ListTypesExpected = []string{
	"OS::Nova::Server",
	"OS::Heat::RandomString",
	"OS::Swift::Container",
	"OS::Trove::Instance",
	"OS::Nova::FloatingIPAssociation",
	"OS::Cinder::VolumeAttachment",
	"OS::Nova::FloatingIP",
	"OS::Nova::KeyPair",
}

const ListTypesOutput = `
{
  "resource_types": [
    "OS::Nova::Server",
    "OS::Heat::RandomString",
    "OS::Swift::Container",
    "OS::Trove::Instance",
    "OS::Nova::FloatingIPAssociation",
    "OS::Cinder::VolumeAttachment",
    "OS::Nova::FloatingIP",
    "OS::Nova::KeyPair"
  ]
}`

// HandleListTypesSuccessfully creates an HTTP handler at `/resource_types`
// on the test handler mux that responds with a `ListTypes` response.
func HandleListTypesSuccessfully(t *testing.T, output string) {
	th.Mux.HandleFunc("/resource_types", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, output)
	})
}
