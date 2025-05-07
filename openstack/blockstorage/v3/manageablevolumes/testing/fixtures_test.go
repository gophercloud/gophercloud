package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func MockManageExistingResponse(t *testing.T) {
	th.Mux.HandleFunc("/manageable_volumes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "volume": {
        "host": "host@lvm#LVM",
        "ref": {
            "source-name": "volume-73796b96-169f-4675-a5bc-73fc0f8f9a17"
        },
        "name": "New Volume",
        "availability_zone": "nova",
        "description": "Volume imported from existingLV",
        "volume_type": "lvm",
        "bootable": true,
        "metadata": {
            "key1": "value1",
            "key2": "value2"
        }
    }
}
		`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, `
{
    "volume": {
        "id": "23cf872b-c781-4cd4-847d-5f2ec8cbd91c",
        "status": "creating",
        "size": 0,
        "availability_zone": "nova",
		"created_at": "2025-03-20T11:58:05.000000",
        "updated_at": "2025-03-20T11:58:05.000000",
        "name": "New Volume",
        "description": "Volume imported from existingLV",
        "volume_type": "lvm",
        "snapshot_id": null,
        "source_volid": null,
        "metadata": {
            "key1": "value1",
            "key2": "value2"
        },
        "links": [
            {
                "href": "http://10.0.2.15:8776/v3/87c8522052ca4eed98bc672b4c1a3ddb/volumes/23cf872b-c781-4cd4-847d-5f2ec8cbd91c",
                "rel": "self"
            },
            {
                "href": "http://10.0.2.15:8776/87c8522052ca4eed98bc672b4c1a3ddb/volumes/23cf872b-c781-4cd4-847d-5f2ec8cbd91c",
                "rel": "bookmark"
            }
        ],
        "user_id": "eae1472b5fc5496998a3d06550929e7e",
        "bootable": "true",
        "encrypted": false,
		"replication_status": null,
		"consistencygroup_id": null,
		"multiattach": false,
        "attachments": [],
        "created_at": "2014-07-18T00:12:54.000000",
		"migration_status": null,
		"group_id": null,
		"provider_id": null,
		"shared_targets": true,
		"service_uuid": null,
		"cluster_name": null,
		"volume_type_id": "a218796e-605b-4b6f-9dfc-8be95a0d7d03",
		"consumes_quota": true,
		"os-vol-mig-status-attr:migstat": null,
		"os-vol-mig-status-attr:name_id": null,
        "os-vol-tenant-attr:tenant_id": "87c8522052ca4eed98bc672b4c1a3ddb",
		"os-vol-host-attr:host": "host@lvm#LVM"
	}
}
		`)
	})
}
