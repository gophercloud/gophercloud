package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/ptr"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocationcandidates"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const AllocationCandidatesBody = `
{
  "allocation_requests": [
    {
      "allocations": {
        "rp-uuid-1": {
          "resources": {
            "VCPU": 1,
            "MEMORY_MB": 1024
          }
        },
        "rp-uuid-2": {
          "resources": {
            "DISK_GB": 100
          }
        }
      },
      "mappings": {
        "": ["rp-uuid-1"],
        "_DISK": ["rp-uuid-2"]
      }
    }
  ],
  "provider_summaries": {
    "rp-uuid-1": {
      "resources": {
        "VCPU": {
          "capacity": 16,
          "used": 0
        },
        "MEMORY_MB": {
          "capacity": 32768,
          "used": 0
        }
      },
      "traits": ["HW_CPU_X86_AVX2", "HW_CPU_X86_SSE"],
      "parent_provider_uuid": null,
      "root_provider_uuid": "rp-uuid-1"
    },
    "rp-uuid-2": {
      "resources": {
        "DISK_GB": {
          "capacity": 1900,
          "used": 0
        }
      },
      "traits": ["MISC_SHARES_VIA_AGGREGATE"],
      "parent_provider_uuid": null,
      "root_provider_uuid": "rp-uuid-2"
    }
  }
}
`

const AllocationCandidatesBodyPre134 = `
{
  "allocation_requests": [
    {
      "allocations": {
        "rp-uuid-1": {
          "resources": {
            "VCPU": 1,
            "MEMORY_MB": 1024
          }
        },
        "rp-uuid-2": {
          "resources": {
            "DISK_GB": 100
          }
        }
      }
    }
  ],
  "provider_summaries": {
    "rp-uuid-1": {
      "resources": {
        "VCPU": {
          "capacity": 16,
          "used": 0
        },
        "MEMORY_MB": {
          "capacity": 32768,
          "used": 0
        }
      },
      "traits": ["HW_CPU_X86_AVX2", "HW_CPU_X86_SSE"],
      "parent_provider_uuid": null,
      "root_provider_uuid": "rp-uuid-1"
    },
    "rp-uuid-2": {
      "resources": {
        "DISK_GB": {
          "capacity": 1900,
          "used": 0
        }
      },
      "traits": ["MISC_SHARES_VIA_AGGREGATE"],
      "parent_provider_uuid": null,
      "root_provider_uuid": "rp-uuid-2"
    }
  }
}
`

const AllocationCandidatesBodyPre129 = `
{
  "allocation_requests": [
    {
      "allocations": {
        "rp-uuid-1": {
          "resources": {
            "VCPU": 1,
            "MEMORY_MB": 1024
          }
        },
        "rp-uuid-2": {
          "resources": {
            "DISK_GB": 100
          }
        }
      }
    }
  ],
  "provider_summaries": {
    "rp-uuid-1": {
      "resources": {
        "VCPU": {
          "capacity": 16,
          "used": 0
        },
        "MEMORY_MB": {
          "capacity": 32768,
          "used": 0
        }
      },
      "traits": ["HW_CPU_X86_AVX2", "HW_CPU_X86_SSE"]
    },
    "rp-uuid-2": {
      "resources": {
        "DISK_GB": {
          "capacity": 1900,
          "used": 0
        }
      },
      "traits": ["MISC_SHARES_VIA_AGGREGATE"]
    }
  }
}
`

const AllocationCandidatesBodyPre117 = `
{
  "allocation_requests": [
    {
      "allocations": {
        "rp-uuid-1": {
          "resources": {
            "VCPU": 1,
            "MEMORY_MB": 1024
          }
        },
        "rp-uuid-2": {
          "resources": {
            "DISK_GB": 100
          }
        }
      }
    }
  ],
  "provider_summaries": {
    "rp-uuid-1": {
      "resources": {
        "VCPU": {
          "capacity": 16,
          "used": 0
        },
        "MEMORY_MB": {
          "capacity": 32768,
          "used": 0
        }
      }
    },
    "rp-uuid-2": {
      "resources": {
        "DISK_GB": {
          "capacity": 1900,
          "used": 0
        }
      }
    }
  }
}
`

const AllocationCandidatesBody110 = `
{
  "allocation_requests": [
    {
      "allocations": [
        {
          "resource_provider": {
            "uuid": "rp-uuid-1"
          },
          "resources": {
            "VCPU": 1,
            "MEMORY_MB": 1024
          }
        }
      ]
    }
  ],
  "provider_summaries": {
    "rp-uuid-1": {
      "resources": {
        "VCPU": {
          "capacity": 16,
          "used": 0
        },
        "MEMORY_MB": {
          "capacity": 32768,
          "used": 0
        }
      }
    }
  }
}
`

const AllocationCandidatesEmptyBody = `
{
  "allocation_requests": [],
  "provider_summaries": {}
}
`

var ExpectedAllocationCandidates = allocationcandidates.AllocationCandidates{
	AllocationRequests: []allocationcandidates.AllocationRequest{
		{
			Allocations: map[string]allocationcandidates.AllocationRequestResource{
				"rp-uuid-1": {
					Resources: map[string]int{
						"VCPU":      1,
						"MEMORY_MB": 1024,
					},
				},
				"rp-uuid-2": {
					Resources: map[string]int{
						"DISK_GB": 100,
					},
				},
			},
			Mappings: &map[string][]string{
				"":      {"rp-uuid-1"},
				"_DISK": {"rp-uuid-2"},
			},
		},
	},
	ProviderSummaries: map[string]allocationcandidates.ProviderSummary{
		"rp-uuid-1": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"VCPU": {
					Capacity: 16,
					Used:     0,
				},
				"MEMORY_MB": {
					Capacity: 32768,
					Used:     0,
				},
			},
			Traits:             ptr.To([]string{"HW_CPU_X86_AVX2", "HW_CPU_X86_SSE"}),
			ParentProviderUUID: nil,
			RootProviderUUID:   ptr.To("rp-uuid-1"),
		},
		"rp-uuid-2": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"DISK_GB": {
					Capacity: 1900,
					Used:     0,
				},
			},
			Traits:             ptr.To([]string{"MISC_SHARES_VIA_AGGREGATE"}),
			ParentProviderUUID: nil,
			RootProviderUUID:   ptr.To("rp-uuid-2"),
		},
	},
}

var ExpectedAllocationCandidatesPre134 = allocationcandidates.AllocationCandidates{
	AllocationRequests: []allocationcandidates.AllocationRequest{
		{
			Allocations: map[string]allocationcandidates.AllocationRequestResource{
				"rp-uuid-1": {
					Resources: map[string]int{
						"VCPU":      1,
						"MEMORY_MB": 1024,
					},
				},
				"rp-uuid-2": {
					Resources: map[string]int{
						"DISK_GB": 100,
					},
				},
			},
		},
	},
	ProviderSummaries: map[string]allocationcandidates.ProviderSummary{
		"rp-uuid-1": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"VCPU": {
					Capacity: 16,
					Used:     0,
				},
				"MEMORY_MB": {
					Capacity: 32768,
					Used:     0,
				},
			},
			Traits:             ptr.To([]string{"HW_CPU_X86_AVX2", "HW_CPU_X86_SSE"}),
			ParentProviderUUID: nil,
			RootProviderUUID:   ptr.To("rp-uuid-1"),
		},
		"rp-uuid-2": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"DISK_GB": {
					Capacity: 1900,
					Used:     0,
				},
			},
			Traits:             ptr.To([]string{"MISC_SHARES_VIA_AGGREGATE"}),
			ParentProviderUUID: nil,
			RootProviderUUID:   ptr.To("rp-uuid-2"),
		},
	},
}

var ExpectedAllocationCandidatesPre129 = allocationcandidates.AllocationCandidates{
	AllocationRequests: []allocationcandidates.AllocationRequest{
		{
			Allocations: map[string]allocationcandidates.AllocationRequestResource{
				"rp-uuid-1": {
					Resources: map[string]int{
						"VCPU":      1,
						"MEMORY_MB": 1024,
					},
				},
				"rp-uuid-2": {
					Resources: map[string]int{
						"DISK_GB": 100,
					},
				},
			},
		},
	},
	ProviderSummaries: map[string]allocationcandidates.ProviderSummary{
		"rp-uuid-1": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"VCPU": {
					Capacity: 16,
					Used:     0,
				},
				"MEMORY_MB": {
					Capacity: 32768,
					Used:     0,
				},
			},
			Traits: ptr.To([]string{"HW_CPU_X86_AVX2", "HW_CPU_X86_SSE"}),
		},
		"rp-uuid-2": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"DISK_GB": {
					Capacity: 1900,
					Used:     0,
				},
			},
			Traits: ptr.To([]string{"MISC_SHARES_VIA_AGGREGATE"}),
		},
	},
}

var ExpectedAllocationCandidatesPre117 = allocationcandidates.AllocationCandidates{
	AllocationRequests: []allocationcandidates.AllocationRequest{
		{
			Allocations: map[string]allocationcandidates.AllocationRequestResource{
				"rp-uuid-1": {
					Resources: map[string]int{
						"VCPU":      1,
						"MEMORY_MB": 1024,
					},
				},
				"rp-uuid-2": {
					Resources: map[string]int{
						"DISK_GB": 100,
					},
				},
			},
		},
	},
	ProviderSummaries: map[string]allocationcandidates.ProviderSummary{
		"rp-uuid-1": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"VCPU": {
					Capacity: 16,
					Used:     0,
				},
				"MEMORY_MB": {
					Capacity: 32768,
					Used:     0,
				},
			},
		},
		"rp-uuid-2": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"DISK_GB": {
					Capacity: 1900,
					Used:     0,
				},
			},
		},
	},
}

var ExpectedAllocationCandidates110 = allocationcandidates.AllocationCandidates110{
	AllocationRequests: []allocationcandidates.AllocationRequest110{
		{
			Allocations: []allocationcandidates.AllocationRequest110Resource{
				{
					ResourceProvider: allocationcandidates.AllocationRequest110ResourceProvider{
						UUID: "rp-uuid-1",
					},
					Resources: map[string]int{
						"VCPU":      1,
						"MEMORY_MB": 1024,
					},
				},
			},
		},
	},
	ProviderSummaries: map[string]allocationcandidates.ProviderSummary110{
		"rp-uuid-1": {
			Resources: map[string]allocationcandidates.ProviderSummaryResource{
				"VCPU": {
					Capacity: 16,
					Used:     0,
				},
				"MEMORY_MB": {
					Capacity: 32768,
					Used:     0,
				},
			},
		},
	},
}

func HandleListAllocationCandidatesSuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocation_candidates",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationCandidatesBody)
		})
}

func HandleListAllocationCandidatesPre134Success(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocation_candidates",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationCandidatesBodyPre134)
		})
}

func HandleListAllocationCandidatesPre129Success(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocation_candidates",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationCandidatesBodyPre129)
		})
}

func HandleListAllocationCandidatesPre117Success(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocation_candidates",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationCandidatesBodyPre117)
		})
}

func HandleListAllocationCandidates110Success(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocation_candidates",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationCandidatesBody110)
		})
}

func HandleListAllocationCandidatesEmptySuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocation_candidates",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationCandidatesEmptyBody)
		})
}

func HandleListAllocationCandidatesBadRequest(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocation_candidates",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			fmt.Fprint(w, `{"errors": [{"status": 400, "detail": "Invalid resources parameter."}]}`)
		})
}

func HandleListAllocationCandidatesWithFullQuerySuccess(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/allocation_candidates",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			q := r.URL.Query()
			th.AssertEquals(t, "VCPU:1,MEMORY_MB:1024", q.Get("resources"))
			th.AssertEquals(t, "5", q.Get("limit"))
			th.AssertEquals(t, "isolate", q.Get("group_policy"))
			th.AssertEquals(t, "SRIOV_NET_VF:1", q.Get("resources1"))
			th.AssertEquals(t, "CUSTOM_PHYSNET1", q.Get("required1"))
			// required is repeated; order matches the ListOpts.Required slice order.
			th.AssertDeepEquals(t, []string{"HW_CPU_X86_SSE", "!HW_CPU_X86_AVX2"}, q["required"])

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationCandidatesBody)
		})
}
