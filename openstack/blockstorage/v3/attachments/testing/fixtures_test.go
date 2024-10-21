package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/attachments"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

var (
	attachedAt, _      = time.Parse(gophercloud.RFC3339MilliNoZ, "2015-09-16T09:28:52.000000")
	detachedAt, _      = time.Parse(gophercloud.RFC3339MilliNoZ, "2015-09-16T09:28:52.000000")
	expectedAttachment = &attachments.Attachment{
		ID:             "05551600-a936-4d4a-ba42-79a037c1-c91a",
		VolumeID:       "289da7f8-6440-407c-9fb4-7db01ec49164",
		Instance:       "83ec2e3b-4321-422b-8706-a84185f52a0a",
		AttachMode:     "rw",
		Status:         "attaching",
		AttachedAt:     attachedAt,
		DetachedAt:     detachedAt,
		ConnectionInfo: map[string]any{},
	}
)

func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc("/attachments/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

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
    "attachments": [
        {
            "status": "attaching",
            "detached_at": "2015-09-16T09:28:52.000000",
            "connection_info": {},
            "attached_at": "2015-09-16T09:28:52.000000",
            "attach_mode": "rw",
            "instance": "83ec2e3b-4321-422b-8706-a84185f52a0a",
            "volume_id": "289da7f8-6440-407c-9fb4-7db01ec49164",
            "id": "05551600-a936-4d4a-ba42-79a037c1-c91a"
        }
    ],
	"attachments_links": [
	{
		"href": "%s/attachments/detail?marker=1",
		"rel": "next"
	}
    ]
}
  `, th.Server.URL)
		case "1":
			fmt.Fprintf(w, `{"volumes": []}`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc("/attachments/05551600-a936-4d4a-ba42-79a037c1-c91a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
{
    "attachment": {
        "status": "attaching",
        "detached_at": "2015-09-16T09:28:52.000000",
        "connection_info": {},
        "attached_at": "2015-09-16T09:28:52.000000",
        "attach_mode": "rw",
        "instance": "83ec2e3b-4321-422b-8706-a84185f52a0a",
        "volume_id": "289da7f8-6440-407c-9fb4-7db01ec49164",
        "id": "05551600-a936-4d4a-ba42-79a037c1-c91a"
    }
}
      `)
	})
}

func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc("/attachments", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "attachment": {
        "instance_uuid": "83ec2e3b-4321-422b-8706-a84185f52a0a",
        "connector": {
            "initiator": "iqn.1993-08.org.debian: 01: cad181614cec",
            "ip": "192.168.1.20",
            "platform": "x86_64",
            "host": "tempest-1",
            "os_type": "linux2",
            "multipath": false,
            "mountpoint": "/dev/vdb",
            "mode": "rw"
        },
        "volume_uuid": "289da7f8-6440-407c-9fb4-7db01ec49164"
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "attachment": {
        "status": "attaching",
        "detached_at": "2015-09-16T09:28:52.000000",
        "connection_info": {},
        "attached_at": "2015-09-16T09:28:52.000000",
        "attach_mode": "rw",
        "instance": "83ec2e3b-4321-422b-8706-a84185f52a0a",
        "volume_id": "289da7f8-6440-407c-9fb4-7db01ec49164",
        "id": "05551600-a936-4d4a-ba42-79a037c1-c91a"
    }
}
    `)
	})
}

func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/attachments/05551600-a936-4d4a-ba42-79a037c1-c91a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusOK)
	})
}

func MockUpdateResponse(t *testing.T) {
	th.Mux.HandleFunc("/attachments/05551600-a936-4d4a-ba42-79a037c1-c91a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
    "attachment": {
        "connector": {
            "initiator": "iqn.1993-08.org.debian: 01: cad181614cec",
            "ip": "192.168.1.20",
            "platform": "x86_64",
            "host": "tempest-1",
            "os_type": "linux2",
            "multipath": false,
            "mountpoint": "/dev/vdb",
            "mode": "rw"
        }
    }
}
      `)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
{
    "attachment": {
        "status": "attaching",
        "detached_at": "2015-09-16T09:28:52.000000",
        "connection_info": {},
        "attached_at": "2015-09-16T09:28:52.000000",
        "attach_mode": "rw",
        "instance": "83ec2e3b-4321-422b-8706-a84185f52a0a",
        "volume_id": "289da7f8-6440-407c-9fb4-7db01ec49164",
        "id": "05551600-a936-4d4a-ba42-79a037c1-c91a"
    }
}
        `)
	})
}

func MockUpdateEmptyResponse(t *testing.T) {
	th.Mux.HandleFunc("/attachments/05551600-a936-4d4a-ba42-79a037c1-c91a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
    "attachment": {
        "connector": null
    }
}
      `)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
{
    "attachment": {
        "status": "attaching",
        "detached_at": "2015-09-16T09:28:52.000000",
        "connection_info": {},
        "attached_at": "2015-09-16T09:28:52.000000",
        "attach_mode": "rw",
        "instance": "83ec2e3b-4321-422b-8706-a84185f52a0a",
        "volume_id": "289da7f8-6440-407c-9fb4-7db01ec49164",
        "id": "05551600-a936-4d4a-ba42-79a037c1-c91a"
    }
}
        `)
	})
}

var completeRequest = `
{
    "os-complete": null
}
`

func MockCompleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/attachments/05551600-a936-4d4a-ba42-79a037c1-c91a/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, completeRequest)
		w.WriteHeader(http.StatusNoContent)
	})
}
