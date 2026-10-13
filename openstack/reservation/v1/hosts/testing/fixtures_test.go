package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/reservation/v1/hosts"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const HostsListResult = `
{
  "hosts": [
    {
      "id": "18",
      "vcpus": 20,
      "cpu_info": "{\"arch\": \"x86_64\", \"model\": \"Broadwell-IBRS\", \"vendor\": \"Intel\", \"topology\": {\"cells\": 2, \"sockets\": 1, \"cores\": 10, \"threads\": 1}, \"maxphysaddr\": {\"mode\": \"emulate\", \"bits\": 46}, \"features\": [\"lm\", \"pcid\", \"sse4.2\", \"arat\", \"fpu\", \"ssbd\", \"md-clear\", \"acpi\", \"sse2\", \"fsgsbase\", \"sse4.1\", \"pae\", \"smap\", \"invpcid\", \"pge\", \"movbe\", \"tsc\", \"mtrr\", \"f16c\", \"vme\", \"tsc_adjust\", \"msr\", \"monitor\", \"ds_cpl\", \"est\", \"syscall\", \"avx\", \"nx\", \"adx\", \"clflush\", \"pni\", \"lahf_lm\", \"hle\", \"pdcm\", \"pat\", \"fxsr\", \"bmi1\", \"de\", \"bmi2\", \"popcnt\", \"mmx\", \"tm\", \"smep\", \"pse\", \"pbe\", \"apic\", \"fma\", \"rdseed\", \"smx\", \"sep\", \"xsave\", \"ssse3\", \"cmov\", \"mce\", \"ds\", \"abm\", \"stibp\", \"sse\", \"invtsc\", \"spec-ctrl\", \"xtpr\", \"ss\", \"pclmuldq\", \"avx2\", \"xsaveopt\", \"erms\", \"3dnowprefetch\", \"intel-pt\", \"aes\", \"rdtscp\", \"cx8\", \"tsc-deadline\", \"rtm\", \"mca\", \"pse36\", \"ht\", \"rdrand\", \"vmx\", \"flush-l1d\", \"dtes64\", \"dca\", \"cx16\", \"tm2\", \"x2apic\", \"pdpe1gb\"]}",
      "hypervisor_type": "QEMU",
      "hypervisor_version": 9000000,
      "hypervisor_hostname": "compute-1.example.com",
      "service_name": "compute-1.example.com",
      "memory_mb": 128297,
      "local_gb": 793,
      "status": null,
      "availability_zone": "nova",
      "trust_id": "e07b8b0f5c1e4b4a9d3e2c1f0a9b8c7d",
      "reservable": true,
      "created_at": "2026-05-27 13:01:19",
      "updated_at": null
    }
  ]
}
`

const HostsListWithCapabilitiesResult = `
{
  "hosts": [
    {
      "id": "18",
      "vcpus": 20,
      "cpu_info": "{\"arch\": \"x86_64\", \"model\": \"Broadwell-IBRS\", \"vendor\": \"Intel\", \"topology\": {\"cells\": 2, \"sockets\": 1, \"cores\": 10, \"threads\": 1}, \"maxphysaddr\": {\"mode\": \"emulate\", \"bits\": 46}, \"features\": [\"lm\", \"pcid\", \"sse4.2\", \"arat\", \"fpu\", \"ssbd\", \"md-clear\", \"acpi\", \"sse2\", \"fsgsbase\", \"sse4.1\", \"pae\", \"smap\", \"invpcid\", \"pge\", \"movbe\", \"tsc\", \"mtrr\", \"f16c\", \"vme\", \"tsc_adjust\", \"msr\", \"monitor\", \"ds_cpl\", \"est\", \"syscall\", \"avx\", \"nx\", \"adx\", \"clflush\", \"pni\", \"lahf_lm\", \"hle\", \"pdcm\", \"pat\", \"fxsr\", \"bmi1\", \"de\", \"bmi2\", \"popcnt\", \"mmx\", \"tm\", \"smep\", \"pse\", \"pbe\", \"apic\", \"fma\", \"rdseed\", \"smx\", \"sep\", \"xsave\", \"ssse3\", \"cmov\", \"mce\", \"ds\", \"abm\", \"stibp\", \"sse\", \"invtsc\", \"spec-ctrl\", \"xtpr\", \"ss\", \"pclmuldq\", \"avx2\", \"xsaveopt\", \"erms\", \"3dnowprefetch\", \"intel-pt\", \"aes\", \"rdtscp\", \"cx8\", \"tsc-deadline\", \"rtm\", \"mca\", \"pse36\", \"ht\", \"rdrand\", \"vmx\", \"flush-l1d\", \"dtes64\", \"dca\", \"cx16\", \"tm2\", \"x2apic\", \"pdpe1gb\"]}",
      "hypervisor_type": "QEMU",
      "hypervisor_version": 9000000,
      "hypervisor_hostname": "compute-1.example.com",
      "service_name": "compute-1.example.com",
      "memory_mb": 128297,
      "local_gb": 793,
      "status": null,
      "availability_zone": "nova",
      "trust_id": "e07b8b0f5c1e4b4a9d3e2c1f0a9b8c7d",
      "reservable": true,
      "created_at": "2026-05-27 13:01:19",
      "updated_at": null,
      "gpu": "a100",
      "rack": "b12"
    }
  ]
}
`

const HostGetResult = `
{
  "host": {
    "id": "18",
    "vcpus": 20,
    "cpu_info": "{\"arch\": \"x86_64\", \"model\": \"Broadwell-IBRS\", \"vendor\": \"Intel\", \"topology\": {\"cells\": 2, \"sockets\": 1, \"cores\": 10, \"threads\": 1}, \"maxphysaddr\": {\"mode\": \"emulate\", \"bits\": 46}, \"features\": [\"lm\", \"pcid\", \"sse4.2\", \"arat\", \"fpu\", \"ssbd\", \"md-clear\", \"acpi\", \"sse2\", \"fsgsbase\", \"sse4.1\", \"pae\", \"smap\", \"invpcid\", \"pge\", \"movbe\", \"tsc\", \"mtrr\", \"f16c\", \"vme\", \"tsc_adjust\", \"msr\", \"monitor\", \"ds_cpl\", \"est\", \"syscall\", \"avx\", \"nx\", \"adx\", \"clflush\", \"pni\", \"lahf_lm\", \"hle\", \"pdcm\", \"pat\", \"fxsr\", \"bmi1\", \"de\", \"bmi2\", \"popcnt\", \"mmx\", \"tm\", \"smep\", \"pse\", \"pbe\", \"apic\", \"fma\", \"rdseed\", \"smx\", \"sep\", \"xsave\", \"ssse3\", \"cmov\", \"mce\", \"ds\", \"abm\", \"stibp\", \"sse\", \"invtsc\", \"spec-ctrl\", \"xtpr\", \"ss\", \"pclmuldq\", \"avx2\", \"xsaveopt\", \"erms\", \"3dnowprefetch\", \"intel-pt\", \"aes\", \"rdtscp\", \"cx8\", \"tsc-deadline\", \"rtm\", \"mca\", \"pse36\", \"ht\", \"rdrand\", \"vmx\", \"flush-l1d\", \"dtes64\", \"dca\", \"cx16\", \"tm2\", \"x2apic\", \"pdpe1gb\"]}",
    "hypervisor_type": "QEMU",
    "hypervisor_version": 9000000,
    "hypervisor_hostname": "compute-1.example.com",
    "service_name": "compute-1.example.com",
    "memory_mb": 128297,
    "local_gb": 793,
    "status": null,
    "availability_zone": "nova",
    "trust_id": "e07b8b0f5c1e4b4a9d3e2c1f0a9b8c7d",
    "reservable": true,
    "created_at": "2026-05-27 13:01:19",
    "updated_at": null,
    "gpu": "a100",
    "rack": "b12"
  }
}
`

const HostCreateRequest = `
{
  "name": "compute-1.example.com",
  "gpu": "a100",
  "rack": "b12"
}
`

const HostCreateResult = `
{
  "host": {
    "id": "18",
    "vcpus": 20,
    "cpu_info": "{\"arch\": \"x86_64\", \"model\": \"Broadwell-IBRS\", \"vendor\": \"Intel\", \"topology\": {\"cells\": 2, \"sockets\": 1, \"cores\": 10, \"threads\": 1}, \"maxphysaddr\": {\"mode\": \"emulate\", \"bits\": 46}, \"features\": [\"lm\", \"pcid\", \"sse4.2\", \"arat\", \"fpu\", \"ssbd\", \"md-clear\", \"acpi\", \"sse2\", \"fsgsbase\", \"sse4.1\", \"pae\", \"smap\", \"invpcid\", \"pge\", \"movbe\", \"tsc\", \"mtrr\", \"f16c\", \"vme\", \"tsc_adjust\", \"msr\", \"monitor\", \"ds_cpl\", \"est\", \"syscall\", \"avx\", \"nx\", \"adx\", \"clflush\", \"pni\", \"lahf_lm\", \"hle\", \"pdcm\", \"pat\", \"fxsr\", \"bmi1\", \"de\", \"bmi2\", \"popcnt\", \"mmx\", \"tm\", \"smep\", \"pse\", \"pbe\", \"apic\", \"fma\", \"rdseed\", \"smx\", \"sep\", \"xsave\", \"ssse3\", \"cmov\", \"mce\", \"ds\", \"abm\", \"stibp\", \"sse\", \"invtsc\", \"spec-ctrl\", \"xtpr\", \"ss\", \"pclmuldq\", \"avx2\", \"xsaveopt\", \"erms\", \"3dnowprefetch\", \"intel-pt\", \"aes\", \"rdtscp\", \"cx8\", \"tsc-deadline\", \"rtm\", \"mca\", \"pse36\", \"ht\", \"rdrand\", \"vmx\", \"flush-l1d\", \"dtes64\", \"dca\", \"cx16\", \"tm2\", \"x2apic\", \"pdpe1gb\"]}",
    "hypervisor_type": "QEMU",
    "hypervisor_version": 9000000,
    "hypervisor_hostname": "compute-1.example.com",
    "service_name": "compute-1.example.com",
    "memory_mb": 128297,
    "local_gb": 793,
    "status": null,
    "availability_zone": "nova",
    "trust_id": "e07b8b0f5c1e4b4a9d3e2c1f0a9b8c7d",
    "reservable": true,
    "created_at": "2026-05-27 13:01:19",
    "updated_at": null,
    "gpu": "a100",
    "rack": "b12"
  }
}
`

const HostUpdateRequest = `
{
  "gpu": "h100"
}
`

const HostUpdateResult = `
{
  "host": {
    "id": "18",
    "vcpus": 20,
    "cpu_info": "{\"arch\": \"x86_64\", \"model\": \"Broadwell-IBRS\", \"vendor\": \"Intel\", \"topology\": {\"cells\": 2, \"sockets\": 1, \"cores\": 10, \"threads\": 1}, \"maxphysaddr\": {\"mode\": \"emulate\", \"bits\": 46}, \"features\": [\"lm\", \"pcid\", \"sse4.2\", \"arat\", \"fpu\", \"ssbd\", \"md-clear\", \"acpi\", \"sse2\", \"fsgsbase\", \"sse4.1\", \"pae\", \"smap\", \"invpcid\", \"pge\", \"movbe\", \"tsc\", \"mtrr\", \"f16c\", \"vme\", \"tsc_adjust\", \"msr\", \"monitor\", \"ds_cpl\", \"est\", \"syscall\", \"avx\", \"nx\", \"adx\", \"clflush\", \"pni\", \"lahf_lm\", \"hle\", \"pdcm\", \"pat\", \"fxsr\", \"bmi1\", \"de\", \"bmi2\", \"popcnt\", \"mmx\", \"tm\", \"smep\", \"pse\", \"pbe\", \"apic\", \"fma\", \"rdseed\", \"smx\", \"sep\", \"xsave\", \"ssse3\", \"cmov\", \"mce\", \"ds\", \"abm\", \"stibp\", \"sse\", \"invtsc\", \"spec-ctrl\", \"xtpr\", \"ss\", \"pclmuldq\", \"avx2\", \"xsaveopt\", \"erms\", \"3dnowprefetch\", \"intel-pt\", \"aes\", \"rdtscp\", \"cx8\", \"tsc-deadline\", \"rtm\", \"mca\", \"pse36\", \"ht\", \"rdrand\", \"vmx\", \"flush-l1d\", \"dtes64\", \"dca\", \"cx16\", \"tm2\", \"x2apic\", \"pdpe1gb\"]}",
    "hypervisor_type": "QEMU",
    "hypervisor_version": 9000000,
    "hypervisor_hostname": "compute-1.example.com",
    "service_name": "compute-1.example.com",
    "memory_mb": 128297,
    "local_gb": 793,
    "status": null,
    "availability_zone": "nova",
    "trust_id": "e07b8b0f5c1e4b4a9d3e2c1f0a9b8c7d",
    "reservable": true,
    "created_at": "2026-05-27 13:01:19",
    "updated_at": "2026-05-28 09:12:44",
    "gpu": "h100",
    "rack": "b12"
  }
}
`

// Blazar removes an extra capability when its value is explicitly null.
const HostRemoveCapabilityRequest = `
{
  "gpu": null
}
`

const HostRemoveCapabilityResult = `
{
  "host": {
    "id": "18",
    "vcpus": 20,
    "cpu_info": "{\"arch\": \"x86_64\", \"model\": \"Broadwell-IBRS\", \"vendor\": \"Intel\", \"topology\": {\"cells\": 2, \"sockets\": 1, \"cores\": 10, \"threads\": 1}, \"maxphysaddr\": {\"mode\": \"emulate\", \"bits\": 46}, \"features\": [\"lm\", \"pcid\", \"sse4.2\", \"arat\", \"fpu\", \"ssbd\", \"md-clear\", \"acpi\", \"sse2\", \"fsgsbase\", \"sse4.1\", \"pae\", \"smap\", \"invpcid\", \"pge\", \"movbe\", \"tsc\", \"mtrr\", \"f16c\", \"vme\", \"tsc_adjust\", \"msr\", \"monitor\", \"ds_cpl\", \"est\", \"syscall\", \"avx\", \"nx\", \"adx\", \"clflush\", \"pni\", \"lahf_lm\", \"hle\", \"pdcm\", \"pat\", \"fxsr\", \"bmi1\", \"de\", \"bmi2\", \"popcnt\", \"mmx\", \"tm\", \"smep\", \"pse\", \"pbe\", \"apic\", \"fma\", \"rdseed\", \"smx\", \"sep\", \"xsave\", \"ssse3\", \"cmov\", \"mce\", \"ds\", \"abm\", \"stibp\", \"sse\", \"invtsc\", \"spec-ctrl\", \"xtpr\", \"ss\", \"pclmuldq\", \"avx2\", \"xsaveopt\", \"erms\", \"3dnowprefetch\", \"intel-pt\", \"aes\", \"rdtscp\", \"cx8\", \"tsc-deadline\", \"rtm\", \"mca\", \"pse36\", \"ht\", \"rdrand\", \"vmx\", \"flush-l1d\", \"dtes64\", \"dca\", \"cx16\", \"tm2\", \"x2apic\", \"pdpe1gb\"]}",
    "hypervisor_type": "QEMU",
    "hypervisor_version": 9000000,
    "hypervisor_hostname": "compute-1.example.com",
    "service_name": "compute-1.example.com",
    "memory_mb": 128297,
    "local_gb": 793,
    "status": null,
    "availability_zone": "nova",
    "trust_id": "e07b8b0f5c1e4b4a9d3e2c1f0a9b8c7d",
    "reservable": true,
    "created_at": "2026-05-27 13:01:19",
    "updated_at": "2026-05-28 09:12:44",
    "rack": "b12"
  }
}
`

// CPUInfo is a sample cpu_info string.
const CPUInfo = `{"arch": "x86_64", "model": "Broadwell-IBRS", "vendor": "Intel", "topology": {"cells": 2, "sockets": 1, "cores": 10, "threads": 1}, "maxphysaddr": {"mode": "emulate", "bits": 46}, "features": ["lm", "pcid", "sse4.2", "arat", "fpu", "ssbd", "md-clear", "acpi", "sse2", "fsgsbase", "sse4.1", "pae", "smap", "invpcid", "pge", "movbe", "tsc", "mtrr", "f16c", "vme", "tsc_adjust", "msr", "monitor", "ds_cpl", "est", "syscall", "avx", "nx", "adx", "clflush", "pni", "lahf_lm", "hle", "pdcm", "pat", "fxsr", "bmi1", "de", "bmi2", "popcnt", "mmx", "tm", "smep", "pse", "pbe", "apic", "fma", "rdseed", "smx", "sep", "xsave", "ssse3", "cmov", "mce", "ds", "abm", "stibp", "sse", "invtsc", "spec-ctrl", "xtpr", "ss", "pclmuldq", "avx2", "xsaveopt", "erms", "3dnowprefetch", "intel-pt", "aes", "rdtscp", "cx8", "tsc-deadline", "rtm", "mca", "pse36", "ht", "rdrand", "vmx", "flush-l1d", "dtes64", "dca", "cx16", "tm2", "x2apic", "pdpe1gb"]}`

var ExpectedHost = hosts.Host{
	ID:                 "18",
	HypervisorHostname: "compute-1.example.com",
	HypervisorType:     "QEMU",
	HypervisorVersion:  9000000,
	ServiceName:        "compute-1.example.com",
	VCPUs:              20,
	CPUInfo:            CPUInfo,
	MemoryMB:           128297,
	LocalGB:            793,
	Status:             "",
	AvailabilityZone:   "nova",
	TrustID:            "e07b8b0f5c1e4b4a9d3e2c1f0a9b8c7d",
	Reservable:         true,
	CreatedAt:          time.Date(2026, 5, 27, 13, 1, 19, 0, time.UTC),
	UpdatedAt:          nil,
	ExtraCapabilities:  map[string]any{},
}

var ExpectedHostsList = []hosts.Host{ExpectedHost}

var ExpectedHostWithCapabilities = hosts.Host{
	ID:                 "18",
	HypervisorHostname: "compute-1.example.com",
	HypervisorType:     "QEMU",
	HypervisorVersion:  9000000,
	ServiceName:        "compute-1.example.com",
	VCPUs:              20,
	CPUInfo:            CPUInfo,
	MemoryMB:           128297,
	LocalGB:            793,
	Status:             "",
	AvailabilityZone:   "nova",
	TrustID:            "e07b8b0f5c1e4b4a9d3e2c1f0a9b8c7d",
	Reservable:         true,
	CreatedAt:          time.Date(2026, 5, 27, 13, 1, 19, 0, time.UTC),
	UpdatedAt:          nil,
	ExtraCapabilities: map[string]any{
		"gpu":  "a100",
		"rack": "b12",
	},
}

var ExpectedHostsListWithCapabilities = []hosts.Host{ExpectedHostWithCapabilities}

func HandleListHosts(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, HostsListResult)
		})
}

func HandleCreateHost(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, HostCreateRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			fmt.Fprint(w, HostCreateResult)
		})
}

func HandleUpdateHost(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts/18",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, HostUpdateRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, HostUpdateResult)
		})
}

func HandleRemoveHostCapability(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts/18",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, HostRemoveCapabilityRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, HostRemoveCapabilityResult)
		})
}

func HandleDeleteHost(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts/18",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}

func HandleGetHost(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts/18",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, HostGetResult)
		})
}

func HandleListHostsWithCapabilities(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, HostsListWithCapabilitiesResult)
		})
}

const AllocationsListResult = `
{
  "allocations": [
    {
      "reservations": [
        {
          "end_date": "2026-09-30T19:35:00.000000",
          "id": "c04d56e0-6b31-40b2-a57e-7283958100bd",
          "lease_id": "98c3544d-0afe-4251-8556-700c847a127f",
          "start_date": "2026-09-30T18:35:00.000000"
        }
      ],
      "resource_id": "18"
    }
  ]
}
`

const AllocationGetResult = `
{
  "allocation": {
    "reservations": [
      {
        "end_date": "2026-09-30T19:35:00.000000",
        "id": "c04d56e0-6b31-40b2-a57e-7283958100bd",
        "lease_id": "98c3544d-0afe-4251-8556-700c847a127f",
        "start_date": "2026-09-30T18:35:00.000000"
      }
    ],
    "resource_id": "18"
  }
}
`

const ResourcePropertiesListResult = `
{
  "resource_properties": [
    {
      "property": "gophercloud_test"
    }
  ]
}
`

const ResourcePropertiesListDetailResult = `
{
  "resource_properties": [
    {
      "private": false,
      "property": "gophercloud_test",
      "values": [
        "true"
      ]
    }
  ]
}
`

const ResourcePropertyUpdateRequest = `
{
  "private": true
}
`

const ResourcePropertyUpdateResult = `
{
  "resource_property": {
    "created_at": "2026-09-30T18:38:10.000000",
    "id": "cd3c026a-accc-4013-b119-34032ae8b23c",
    "private": true,
    "property_name": "gophercloud_test",
    "resource_type": "physical:host",
    "updated_at": "2026-09-30T18:38:21.000000"
  }
}
`

var ExpectedAllocation = hosts.Allocation{
	ResourceID: "18",
	Reservations: []hosts.AllocationReservation{
		{
			ID:        "c04d56e0-6b31-40b2-a57e-7283958100bd",
			LeaseID:   "98c3544d-0afe-4251-8556-700c847a127f",
			StartDate: time.Date(2026, 9, 30, 18, 35, 0, 0, time.UTC),
			EndDate:   time.Date(2026, 9, 30, 19, 35, 0, 0, time.UTC),
		},
	},
}

var ExpectedAllocationsList = []hosts.Allocation{ExpectedAllocation}

var propertyPublic = false

var ExpectedResourcePropertiesDetail = []hosts.ResourceProperty{
	{
		Property: "gophercloud_test",
		Private:  &propertyPublic,
		Values:   []string{"true"},
	},
}

var propertyUpdatedAt = time.Date(2026, 9, 30, 18, 38, 21, 0, time.UTC)

var ExpectedUpdatedResourceProperty = hosts.UpdatedResourceProperty{
	ID:           "cd3c026a-accc-4013-b119-34032ae8b23c",
	ResourceType: "physical:host",
	PropertyName: "gophercloud_test",
	Private:      true,
	CreatedAt:    time.Date(2026, 9, 30, 18, 38, 10, 0, time.UTC),
	UpdatedAt:    &propertyUpdatedAt,
}

func HandleListAllocations(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts/allocations",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestFormValues(t, r, map[string]string{
				"lease_id":       "98c3544d-0afe-4251-8556-700c847a127f",
				"reservation_id": "c04d56e0-6b31-40b2-a57e-7283958100bd",
			})

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationsListResult)
		})
}

func HandleGetAllocation(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts/18/allocation",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestFormValues(t, r, map[string]string{
				"lease_id":       "98c3544d-0afe-4251-8556-700c847a127f",
				"reservation_id": "c04d56e0-6b31-40b2-a57e-7283958100bd",
			})

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, AllocationGetResult)
		})
}

func HandleListResourceProperties(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts/properties",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			if r.URL.Query().Get("detail") == "true" {
				th.TestFormValues(t, r, map[string]string{"all": "true", "detail": "true"})
				fmt.Fprint(w, ResourcePropertiesListDetailResult)
				return
			}

			th.TestFormValues(t, r, map[string]string{})
			fmt.Fprint(w, ResourcePropertiesListResult)
		})
}

func HandleUpdateResourceProperty(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-hosts/properties/gophercloud_test",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, ResourcePropertyUpdateRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprint(w, ResourcePropertyUpdateResult)
		})
}
