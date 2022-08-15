package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const VolumeListBody = `
	{
  "connectors": [
    {
      "href": "http://127.0.0.1:6385/v1/volume/connectors",
      "rel": "self"
    },
    {
      "href": "http://127.0.0.1:6385/volume/connectors",
      "rel": "bookmark"
    }
  ],
  "links": [
    {
      "href": "http://127.0.0.1:6385/v1/volume/",
      "rel": "self"
    },
    {
      "href": "http://127.0.0.1:6385/volume/",
      "rel": "bookmark"
    }
  ],
  "targets": [
    {
      "href": "http://127.0.0.1:6385/v1/volume/targets",
      "rel": "self"
    },
    {
      "href": "http://127.0.0.1:6385/volume/targets",
      "rel": "bookmark"
    }
  ]
}
`

const ConnectorListBody = `{
  "connectors": [
    {
      "connector_id": "iqn.2017-07.org.openstack:01:d9a51732c3f",
      "links": [
        {
          "href": "http://127.0.0.1:6385/v1/volume/connectors/9bf93e01-d728-47a3-ad4b-5e66a835037c",
          "rel": "self"
        },
        {
          "href": "http://127.0.0.1:6385/volume/connectors/9bf93e01-d728-47a3-ad4b-5e66a835037c",
          "rel": "bookmark"
        }
      ],
      "node_uuid": "6d85703a-565d-469a-96ce-30b6de53079d",
      "type": "iqn",
      "uuid": "9bf93e01-d728-47a3-ad4b-5e66a835037c"
    }
  ]
}`

const TargetListBody = `
	{
  "targets": [
    {
      "boot_index": 0,
      "links": [
        {
          "href": "http://127.0.0.1:6385/v1/volume/targets/bd4d008c-7d31-463d-abf9-6c23d9d55f7f",
          "rel": "self"
        },
        {
          "href": "http://127.0.0.1:6385/volume/targets/bd4d008c-7d31-463d-abf9-6c23d9d55f7f",
          "rel": "bookmark"
        }
      ],
      "node_uuid": "6d85703a-565d-469a-96ce-30b6de53079d",
      "uuid": "bd4d008c-7d31-463d-abf9-6c23d9d55f7f",
      "volume_id": "04452bed-5367-4202-8bf5-de4335ac56d2",
      "volume_type": "iscsi"
    }
  ]
}
`

const SingleConnectorBody = `
	{
  "connector_id": "iqn.2017-07.org.openstack:01:d9a51732c3f",
  "created_at": "2016-08-18T22:28:48.643434+11:11",
  "extra": {},
  "links": [
    {
      "href": "http://127.0.0.1:6385/v1/volume/connectors/9bf93e01-d728-47a3-ad4b-5e66a835037c",
      "rel": "self"
    },
    {
      "href": "http://127.0.0.1:6385/volume/connectors/9bf93e01-d728-47a3-ad4b-5e66a835037c",
      "rel": "bookmark"
    }
  ],
  "node_uuid": "6d85703a-565d-469a-96ce-30b6de53079d",
  "type": "iqn",
  "updated_at": null,
  "uuid": "9bf93e01-d728-47a3-ad4b-5e66a835037c"
}
`

const SingleTargetBody = `
	{
  "boot_index": 0,
  "created_at": "2016-08-18T22:28:48.643434+11:11",
  "extra": {},
  "links": [
    {
      "href": "http://127.0.0.1:6385/v1/volume/targets/bd4d008c-7d31-463d-abf9-6c23d9d55f7f",
      "rel": "self"
    },
    {
      "href": "http://127.0.0.1:6385/volume/targets/bd4d008c-7d31-463d-abf9-6c23d9d55f7f",
      "rel": "bookmark"
    }
  ],
  "node_uuid": "6d85703a-565d-469a-96ce-30b6de53079d",
  "properties": {},
  "updated_at": null,
  "uuid": "bd4d008c-7d31-463d-abf9-6c23d9d55f7f",
  "volume_id": "04452bed-5367-4202-8bf5-de4335ac56d2",
  "volume_type": "iscsi"
}
`

func HandleVolumeListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/volume", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "appliction/json")
		fmt.Fprintf(w, VolumeListBody)
	})
}

func HandleConnectorListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/volume/connectors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "appliction/json")
		fmt.Fprintf(w, ConnectorListBody)
	})
}

func HandleConnectorCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/volume/connectors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
          "node_uuid": "6d85703a-565d-469a-96ce-30b6de53079d",
          "type": "iqn",
          "connector_id": "iqn.2017-07.org.openstack:01:d9a51732c3f"
        }`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

func HandleConnectorDeletionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/volume/connectors/9bf93e01-d728-47a3-ad4b-5e66a835037c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}
func HandleConnectorGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/volume/connectors/9bf93e01-d728-47a3-ad4b-5e66a835037c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		fmt.Fprintf(w, SingleConnectorBody)
	})
}

func HandleConnectorUpdateSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/volume/connectors/9bf93e01-d728-47a3-ad4b-5e66a835037c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `[{"op": "replace", "path": "/connector_id", "value": "iqn.2017-07.org.openstack:01:66666666666"}]`)
		fmt.Fprintf(w, response)
	})
}

func HandleTargetListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/volume/targets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "appliction/json")
		fmt.Fprintf(w, TargetListBody)
	})
}

func HandleTargetCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/volume/targets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
          "boot_index": "0",
          "node_uuid": "6d85703a-565d-469a-96ce-30b6de53079d",
          "volume_id": "04452bed-5367-4202-8bf5-de4335ac56d2",
          "volume_type":"iscsi"
        }`)
		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

func HandleTargetDeletionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/volume/targets/bd4d008c-7d31-463d-abf9-6c23d9d55f7f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}
func HandleTargetGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/volume/targets/bd4d008c-7d31-463d-abf9-6c23d9d55f7f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		fmt.Fprintf(w, SingleTargetBody)
	})
}

func HandleTargetUpdateSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/volume/targets/bd4d008c-7d31-463d-abf9-6c23d9d55f7f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `[{"op": "replace", "path": "/volume_id", "value": "06666bed-5367-4202-8bf5-de4335ac56d2"}]`)
		fmt.Fprintf(w, response)
	})
}
