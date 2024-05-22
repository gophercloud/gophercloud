package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/schedulerstats"
	"github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const PoolsListBody = `
{
    "pools": [
        {
            "name": "opencloud@alpha#ALPHA_pool",
            "host": "opencloud",
            "backend": "alpha",
            "pool": "ALPHA_pool"
        },
        {
            "name": "opencloud@beta#BETA_pool",
            "host": "opencloud",
            "backend": "beta",
            "pool": "BETA_pool"
        },
        {
            "name": "opencloud@gamma#GAMMA_pool",
            "host": "opencloud",
            "backend": "gamma",
            "pool": "GAMMA_pool"
        },
        {
            "name": "opencloud@delta#DELTA_pool",
            "host": "opencloud",
            "backend": "delta",
            "pool": "DELTA_pool"
        }
    ]
}
`

const PoolsListBodyDetail = `
{
    "pools": [
        {
            "name": "opencloud@alpha#ALPHA_pool",
            "host": "opencloud",
            "backend": "alpha",
            "pool": "ALPHA_pool",
            "capabilities": {
                "pool_name": "ALPHA_pool",
                "total_capacity_gb": 1230.0,
                "free_capacity_gb": 1210.0,
                "reserved_percentage": 0,
                "share_backend_name": "ALPHA",
                "storage_protocol": "NFS_CIFS",
                "vendor_name": "Open Source",
                "driver_version": "1.0",
                "timestamp": "2019-05-07T00:28:02.935569",
                "driver_handles_share_servers": true,
                "snapshot_support": true,
                "create_share_from_snapshot_support": true,
                "revert_to_snapshot_support": true,
                "mount_snapshot_support": true,
                "dedupe": false,
                "compression": false,
                "replication_type": null,
                "replication_domain": null,
                "sg_consistent_snapshot_support": "pool",
                "ipv4_support": true,
                "ipv6_support": false
            }
        },
        {
            "name": "opencloud@beta#BETA_pool",
            "host": "opencloud",
            "backend": "beta",
            "pool": "BETA_pool",
            "capabilities": {
                "pool_name": "BETA_pool",
                "total_capacity_gb": 1230.0,
                "free_capacity_gb": 1210.0,
                "reserved_percentage": 0,
                "share_backend_name": "BETA",
                "storage_protocol": "NFS_CIFS",
                "vendor_name": "Open Source",
                "driver_version": "1.0",
                "timestamp": "2019-05-07T00:28:02.817309",
                "driver_handles_share_servers": true,
                "snapshot_support": true,
                "create_share_from_snapshot_support": true,
                "revert_to_snapshot_support": true,
                "mount_snapshot_support": true,
                "dedupe": false,
                "compression": false,
                "replication_type": null,
                "replication_domain": null,
                "sg_consistent_snapshot_support": "pool",
                "ipv4_support": true,
                "ipv6_support": false
            }
        },
        {
            "name": "opencloud@gamma#GAMMA_pool",
            "host": "opencloud",
            "backend": "gamma",
            "pool": "GAMMA_pool",
            "capabilities": {
                "pool_name": "GAMMA_pool",
                "total_capacity_gb": 1230.0,
                "free_capacity_gb": 1210.0,
                "reserved_percentage": 0,
                "replication_type": "readable",
                "share_backend_name": "GAMMA",
                "storage_protocol": "NFS_CIFS",
                "vendor_name": "Open Source",
                "driver_version": "1.0",
                "timestamp": "2019-05-07T00:28:02.899888",
                "driver_handles_share_servers": false,
                "snapshot_support": true,
                "create_share_from_snapshot_support": true,
                "revert_to_snapshot_support": true,
                "mount_snapshot_support": true,
                "dedupe": false,
                "compression": false,
                "sg_consistent_snapshot_support": "pool",
                "ipv4_support": true,
                "ipv6_support": false
            }
        },
        {
            "name": "opencloud@delta#DELTA_pool",
            "host": "opencloud",
            "backend": "delta",
            "pool": "DELTA_pool",
            "capabilities": {
                "pool_name": "DELTA_pool",
                "total_capacity_gb": 1230.0,
                "free_capacity_gb": 1210.0,
                "reserved_percentage": 0,
                "replication_type": "readable",
                "share_backend_name": "DELTA",
                "storage_protocol": "NFS_CIFS",
                "vendor_name": "Open Source",
                "driver_version": "1.0",
                "timestamp": "2019-05-07T00:28:02.963660",
                "driver_handles_share_servers": false,
                "snapshot_support": true,
                "create_share_from_snapshot_support": true,
                "revert_to_snapshot_support": true,
                "mount_snapshot_support": true,
                "dedupe": false,
                "compression": false,
                "sg_consistent_snapshot_support": "pool",
                "ipv4_support": true,
                "ipv6_support": false
            }
        }
    ]
}
`

var (
	PoolFake1 = schedulerstats.Pool{
		Name:    "opencloud@alpha#ALPHA_pool",
		Host:    "opencloud",
		Backend: "alpha",
		Pool:    "ALPHA_pool",
	}

	PoolFake2 = schedulerstats.Pool{
		Name:    "opencloud@beta#BETA_pool",
		Host:    "opencloud",
		Backend: "beta",
		Pool:    "BETA_pool",
	}

	PoolFake3 = schedulerstats.Pool{
		Name:    "opencloud@gamma#GAMMA_pool",
		Host:    "opencloud",
		Backend: "gamma",
		Pool:    "GAMMA_pool",
	}

	PoolFake4 = schedulerstats.Pool{
		Name:    "opencloud@delta#DELTA_pool",
		Host:    "opencloud",
		Backend: "delta",
		Pool:    "DELTA_pool",
	}

	PoolDetailFake1 = schedulerstats.Pool{
		Name:    "opencloud@alpha#ALPHA_pool",
		Host:    "opencloud",
		Backend: "alpha",
		Pool:    "ALPHA_pool",
		Capabilities: schedulerstats.Capabilities{
			DriverVersion:             "1.0",
			FreeCapacityGB:            1210,
			StorageProtocol:           "NFS_CIFS",
			TotalCapacityGB:           1230,
			VendorName:                "Open Source",
			ShareBackendName:          "ALPHA",
			Timestamp:                 "2019-05-07T00:28:02.935569",
			DriverHandlesShareServers: true,
			SnapshotSupport:           true,
		},
	}

	PoolDetailFake2 = schedulerstats.Pool{
		Name:    "opencloud@beta#BETA_pool",
		Host:    "opencloud",
		Backend: "beta",
		Pool:    "BETA_pool",
		Capabilities: schedulerstats.Capabilities{
			DriverVersion:             "1.0",
			FreeCapacityGB:            1210,
			StorageProtocol:           "NFS_CIFS",
			TotalCapacityGB:           1230,
			VendorName:                "Open Source",
			ShareBackendName:          "BETA",
			Timestamp:                 "2019-05-07T00:28:02.817309",
			DriverHandlesShareServers: true,
			SnapshotSupport:           true,
		},
	}

	PoolDetailFake3 = schedulerstats.Pool{
		Name:    "opencloud@gamma#GAMMA_pool",
		Host:    "opencloud",
		Backend: "gamma",
		Pool:    "GAMMA_pool",
		Capabilities: schedulerstats.Capabilities{
			DriverVersion:             "1.0",
			FreeCapacityGB:            1210,
			StorageProtocol:           "NFS_CIFS",
			TotalCapacityGB:           1230,
			VendorName:                "Open Source",
			ShareBackendName:          "GAMMA",
			Timestamp:                 "2019-05-07T00:28:02.899888",
			DriverHandlesShareServers: false,
			SnapshotSupport:           true,
		},
	}

	PoolDetailFake4 = schedulerstats.Pool{
		Name:    "opencloud@delta#DELTA_pool",
		Host:    "opencloud",
		Backend: "delta",
		Pool:    "DELTA_pool",
		Capabilities: schedulerstats.Capabilities{
			DriverVersion:             "1.0",
			FreeCapacityGB:            1210,
			StorageProtocol:           "NFS_CIFS",
			TotalCapacityGB:           1230,
			VendorName:                "Open Source",
			ShareBackendName:          "DELTA",
			Timestamp:                 "2019-05-07T00:28:02.963660",
			DriverHandlesShareServers: false,
			SnapshotSupport:           true,
		},
	}
)

func HandlePoolsListSuccessfully(t *testing.T) {
	testhelper.Mux.HandleFunc("/scheduler-stats/pools", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		fmt.Fprintf(w, PoolsListBody)

	})
	testhelper.Mux.HandleFunc("/scheduler-stats/pools/detail", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		fmt.Fprintf(w, PoolsListBodyDetail)
	})
}
