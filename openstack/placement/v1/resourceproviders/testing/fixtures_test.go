package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
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

const ResourceProviderUpdateResponse = `
{
  "generation": 1,
  "uuid": "4e8e5957-649f-477b-9e5b-f1f75b21c03c",
  "links": [
	{
	  "href": "/resource_providers/4e8e5957-649f-477b-9e5b-f1f75b21c03c",
	  "rel": "self"
	}
  ],
  "name": "new_name",
  "parent_provider_uuid": "b99b3ab4-3aa6-4fba-b827-69b88b9c544a",
  "root_provider_uuid": "542df8ed-9be2-49b9-b4db-6d3183ff8ec8"
}
`

const ResourceProviderUpdateRequest = `
{
  "name": "new_name",
  "parent_provider_uuid": "b99b3ab4-3aa6-4fba-b827-69b88b9c544a"
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

const AllocationsBody = `
{
    "allocations": {
        "56785a3f-6f1c-4fec-af0b-0faf075b1fcb": {
            "resources": {
                "MEMORY_MB": 256,
                "VCPU": 1
            }
        },
        "9afd5aeb-d6b9-4dea-a588-1e6327a91834": {
            "resources": {
                "MEMORY_MB": 512,
                "VCPU": 2
            }
        },
        "9d16a611-e7f9-4ef3-be26-c61ed01ecefb": {
            "resources": {
                "MEMORY_MB": 1024,
                "VCPU": 1
            }
        }
    },
    "resource_provider_generation": 12
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
		"DISK_GB": {
			AllocationRatio: 1.0,
			MaxUnit:         35,
			MinUnit:         1,
			Reserved:        0,
			StepSize:        1,
			Total:           35,
		},
		"MEMORY_MB": {
			AllocationRatio: 1.5,
			MaxUnit:         5825,
			MinUnit:         1,
			Reserved:        512,
			StepSize:        1,
			Total:           5825,
		},
		"VCPU": {
			AllocationRatio: 16.0,
			MaxUnit:         4,
			MinUnit:         1,
			Reserved:        0,
			StepSize:        1,
			Total:           4,
		},
	},
}

var ExpectedAllocations = resourceproviders.ResourceProviderAllocations{
	ResourceProviderGeneration: 12,
	Allocations: map[string]resourceproviders.Allocation{
		"56785a3f-6f1c-4fec-af0b-0faf075b1fcb": {
			Resources: map[string]int{
				"MEMORY_MB": 256,
				"VCPU":      1,
			},
		},
		"9afd5aeb-d6b9-4dea-a588-1e6327a91834": {
			Resources: map[string]int{
				"MEMORY_MB": 512,
				"VCPU":      2,
			},
		},
		"9d16a611-e7f9-4ef3-be26-c61ed01ecefb": {
			Resources: map[string]int{
				"MEMORY_MB": 1024,
				"VCPU":      1,
			},
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

func HandleResourceProviderGet(t *testing.T) {
	th.Mux.HandleFunc("/resource_providers/99c09379-6e52-4ef8-9a95-b9ce6f68452e", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, ResourceProviderCreateBody)
	})
}

func HandleResourceProviderDelete(t *testing.T) {
	th.Mux.HandleFunc("/resource_providers/b99b3ab4-3aa6-4fba-b827-69b88b9c544a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleResourceProviderUpdate(t *testing.T) {
	th.Mux.HandleFunc("/resource_providers/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, ResourceProviderUpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, ResourceProviderUpdateResponse)
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

func HandleResourceProviderGetAllocations(t *testing.T) {
	allocationsTestUrl := fmt.Sprintf("/resource_providers/%s/allocations", ResourceProviderTestID)

	th.Mux.HandleFunc(allocationsTestUrl,
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, AllocationsBody)
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
