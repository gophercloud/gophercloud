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
