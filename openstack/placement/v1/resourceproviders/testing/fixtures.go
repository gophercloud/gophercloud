package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/placement/v1/resourceproviders"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

const ResourceProviderTestID = "99c09379-6e52-4ef8-9a95-b9ce6f68452e"

const ResourceProvidersBody = `
{
  "resource_providers": [
    {
      "generation": 1,
      "uuid": "99c09379-6e52-4ef8-9a95-b9ce6f68452e",
      "links": [
        {
          "href": "/resource_providers/99c09379-6e52-4ef8-9a95-b9ce6f68452e",
          "rel": "self"
        }
      ],
      "name": "vgr.localdomain",
      "parent_provider_uuid": "542df8ed-9be2-49b9-b4db-6d3183ff8ec8",
      "root_provider_uuid": "542df8ed-9be2-49b9-b4db-6d3183ff8ec8"
    },
    {
      "generation": 2,
      "uuid": "d0b381e9-8761-42de-8e6c-bba99a96d5f5",
      "links": [
        {
          "href": "/resource_providers/d0b381e9-8761-42de-8e6c-bba99a96d5f5",
          "rel": "self"
        }
      ],
      "name": "pony1",
      "parent_provider_uuid": null,
      "root_provider_uuid": "d0b381e9-8761-42de-8e6c-bba99a96d5f5"
    }
  ]
}
`

const ResourceProviderCreateBody = `
{
  "generation": 1,
  "uuid": "99c09379-6e52-4ef8-9a95-b9ce6f68452e",
  "links": [
	{
	  "href": "/resource_providers/99c09379-6e52-4ef8-9a95-b9ce6f68452e",
	  "rel": "self"
	}
  ],
  "name": "vgr.localdomain",
  "parent_provider_uuid": "542df8ed-9be2-49b9-b4db-6d3183ff8ec8",
  "root_provider_uuid": "542df8ed-9be2-49b9-b4db-6d3183ff8ec8"
}
`

const UsagesBody = `
{
    "resource_provider_generation": 1,
    "usages": {
        "DISK_GB": 1,
        "MEMORY_MB": 512,
        "VCPU": 1
    }
}
`

const InventoriesBody = `
{
    "inventories": {
        "DISK_GB": {
            "allocation_ratio": 1.0,
            "max_unit": 35,
            "min_unit": 1,
            "reserved": 0,
            "step_size": 1,
            "total": 35
        },
        "MEMORY_MB": {
            "allocation_ratio": 1.5,
            "max_unit": 5825,
            "min_unit": 1,
            "reserved": 512,
            "step_size": 1,
            "total": 5825
        },
        "VCPU": {
            "allocation_ratio": 16.0,
            "max_unit": 4,
            "min_unit": 1,
            "reserved": 0,
            "step_size": 1,
            "total": 4
        }
    },
    "resource_provider_generation": 7
}
`

const TraitsBody = `
{
    "resource_provider_generation": 1,
    "traits": [
        "CUSTOM_HW_FPGA_CLASS1",
        "CUSTOM_HW_FPGA_CLASS3"
    ]
}
`

var ExpectedResourceProvider1 = resourceproviders.ResourceProvider{
	Generation: 1,
	UUID:       "99c09379-6e52-4ef8-9a95-b9ce6f68452e",
	Links: []resourceproviders.ResourceProviderLinks{
		{
			Href: "/resource_providers/99c09379-6e52-4ef8-9a95-b9ce6f68452e",
			Rel:  "self",
		},
	},
	Name:               "vgr.localdomain",
	ParentProviderUUID: "542df8ed-9be2-49b9-b4db-6d3183ff8ec8",
	RootProviderUUID:   "542df8ed-9be2-49b9-b4db-6d3183ff8ec8",
}

var ExpectedResourceProvider2 = resourceproviders.ResourceProvider{
	Generation: 2,
	UUID:       "d0b381e9-8761-42de-8e6c-bba99a96d5f5",
	Links: []resourceproviders.ResourceProviderLinks{
		{
			Href: "/resource_providers/d0b381e9-8761-42de-8e6c-bba99a96d5f5",
			Rel:  "self",
		},
	},
	Name:               "pony1",
	ParentProviderUUID: "",
	RootProviderUUID:   "d0b381e9-8761-42de-8e6c-bba99a96d5f5",
}

var ExpectedResourceProviders = []resourceproviders.ResourceProvider{
	ExpectedResourceProvider1,
	ExpectedResourceProvider2,
}

var ExpectedUsages = resourceproviders.ResourceProviderUsage{
	ResourceProviderGeneration: 1,
	Usages: map[string]int{
		"DISK_GB":   1,
		"MEMORY_MB": 512,
		"VCPU":      1,
	},
}

var ExpectedInventories = resourceproviders.ResourceProviderInventories{
	ResourceProviderGeneration: 7,
	Inventories: map[string]resourceproviders.Inventory{
		"DISK_GB": resourceproviders.Inventory{
			AllocationRatio: 1.0,
			MaxUnit:         35,
			MinUnit:         1,
			Reserved:        0,
			StepSize:        1,
			Total:           35,
		},
		"MEMORY_MB": resourceproviders.Inventory{
			AllocationRatio: 1.5,
			MaxUnit:         5825,
			MinUnit:         1,
			Reserved:        512,
			StepSize:        1,
			Total:           5825,
		},
		"VCPU": resourceproviders.Inventory{
			AllocationRatio: 16.0,
			MaxUnit:         4,
			MinUnit:         1,
			Reserved:        0,
			StepSize:        1,
			Total:           4,
		},
	},
}

var ExpectedTraits = resourceproviders.ResourceProviderTraits{
	ResourceProviderGeneration: 1,
	Traits: []string{
		"CUSTOM_HW_FPGA_CLASS1",
		"CUSTOM_HW_FPGA_CLASS3",
	},
}

func HandleResourceProviderList(t *testing.T) {
	th.Mux.HandleFunc("/resource_providers",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, ResourceProvidersBody)
		})
}

func HandleResourceProviderCreate(t *testing.T) {
	th.Mux.HandleFunc("/resource_providers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, ResourceProviderCreateBody)
	})
}

func HandleResourceProviderGetUsages(t *testing.T) {
	usageTestUrl := fmt.Sprintf("/resource_providers/%s/usages", ResourceProviderTestID)

	th.Mux.HandleFunc(usageTestUrl,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, UsagesBody)
		})
}

func HandleResourceProviderGetInventories(t *testing.T) {
	inventoriesTestUrl := fmt.Sprintf("/resource_providers/%s/inventories", ResourceProviderTestID)

	th.Mux.HandleFunc(inventoriesTestUrl,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, InventoriesBody)
		})
}

func HandleResourceProviderGetTraits(t *testing.T) {
	traitsTestUrl := fmt.Sprintf("/resource_providers/%s/traits", ResourceProviderTestID)

	th.Mux.HandleFunc(traitsTestUrl,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, TraitsBody)
		})
}
