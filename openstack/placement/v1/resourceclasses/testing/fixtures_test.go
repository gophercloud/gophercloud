package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceclasses"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const PresentResourceClass = "CUSTOM_RESOURCE_CLASS"
const AbsentResourceClass = "NON_EXISTENT_RC"
const NewResourceClass = "CUSTOM_NEW_RC"

const ResourceClassGetResult = `
{
  "links": [
    {
      "href": "/placement/resource_classes/CUSTOM_RESOURCE_CLASS",
      "rel": "self"
    }
  ],
  "name": "CUSTOM_RESOURCE_CLASS"
}
`

const ResourceClassesListResult = `
{
  "resource_classes": [
    {
      "name": "VCPU",
      "links": [
        {
          "href": "/resource_classes/VCPU",
          "rel": "self"
        }
      ]
    },
    {
      "name": "CUSTOM_RESOURCE_CLASS",
      "links": [
        {
          "href": "/placement/resource_classes/CUSTOM_RESOURCE_CLASS",
          "rel": "self"
        }
      ]
    }
  ]
}
`

var ExpectedResourceClass = resourceclasses.ResourceClass{
	Name: "CUSTOM_RESOURCE_CLASS",
	Links: []resourceclasses.Link{
		{
			Href: "/placement/resource_classes/CUSTOM_RESOURCE_CLASS",
			Rel:  "self",
		},
	},
}

var ExpectedResourceClassesList = []resourceclasses.ResourceClass{
	{
		Name: "VCPU",
		Links: []resourceclasses.Link{
			{
				Href: "/resource_classes/VCPU",
				Rel:  "self",
			},
		},
	},
	{
		Name: "CUSTOM_RESOURCE_CLASS",
		Links: []resourceclasses.Link{
			{
				Href: "/placement/resource_classes/CUSTOM_RESOURCE_CLASS",
				Rel:  "self",
			},
		},
	},
}

func HandleListResourceClasses(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, ResourceClassesListResult)
		})
}

func HandleGetResourceClassSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes/"+PresentResourceClass,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, ResourceClassGetResult)
		})
}

func HandleGetResourceClassNotFound(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes/"+AbsentResourceClass,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNotFound)
		})
}

func HandleCreateResourceClassSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, `{"name": "CUSTOM_NEW_RC"}`)

			w.WriteHeader(http.StatusCreated)
		})
}

func HandleCreateResourceClassConflict(t *testing.T, fakeServer th.FakeServer) {
	// We simulate a conflict by trying to create a resource class that already exists.
	fakeServer.Mux.HandleFunc("/resource_classes",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusConflict)
		})
}

func HandleUpdateResourceClassSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes/"+NewResourceClass,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusCreated)
		})
}

func HandleUpdateResourceClassExists(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes/"+PresentResourceClass,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}

func HandleUpdateResourceClassNonCustom(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes/VCPU",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusBadRequest)
		})
}

func HandleDeleteResourceClassSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes/"+PresentResourceClass,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}

func HandleDeleteResourceClassNotFound(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes/"+AbsentResourceClass,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNotFound)
		})
}

func HandleDeleteResourceClassStandardClass(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/resource_classes/VCPU",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusBadRequest)
		})
}
