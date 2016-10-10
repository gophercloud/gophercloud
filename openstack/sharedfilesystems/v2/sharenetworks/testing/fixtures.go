package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func createReq(name, description, network, subnetwork string) string {
	return fmt.Sprintf(`{
		"share_network": {
			"name": "%s",
			"description": "%s",
			"neutron_net_id": "%s",
			"neutron_subnet_id": "%s"
		}
	}`, name, description, network, subnetwork)
}

func createResp(name, description, network, subnetwork string) string {
	return fmt.Sprintf(`
	{
		"share_network": {
			"name": "%s",
			"description": "%s",
			"segmentation_id": null,
			"created_at": "2015-09-07T14:37:00.583656",
			"updated_at": null,
			"id": "77eb3421-4549-4789-ac39-0d5185d68c29",
			"neutron_net_id": "%s",
			"neutron_subnet_id": "%s",
			"ip_version": null,
			"nova_net_id": null,
			"cidr": null,
			"project_id": "e10a683c20da41248cfd5e1ab3d88c62",
			"network_type": null
		}
	}`, name, description, network, subnetwork)
}

func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc("/share-networks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, createReq("my_network",
			"This is my share network",
			"998b42ee-2cee-4d36-8b95-67b5ca1f2109",
			"53482b62-2c84-4a53-b6ab-30d9d9800d06"))

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, createResp("my_network",
			"This is my share network",
			"998b42ee-2cee-4d36-8b95-67b5ca1f2109",
			"53482b62-2c84-4a53-b6ab-30d9d9800d06"))
	})
}

func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/share-networks/fa158a3d-6d9f-4187-9ca5-abbb82646eb2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}
