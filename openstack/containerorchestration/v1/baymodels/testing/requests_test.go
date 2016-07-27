package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/baymodels"
	fake "github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/common"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/baymodels", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "baymodels": [
    {
      "insecure_registry": null,
      "links": [
        {
          "href": "http://65.61.151.130:9511/v1/baymodels/5b793604-fc76-4886-a834-ed522812cdcb",
          "rel": "self"
        },
        {
          "href": "http://65.61.151.130:9511/baymodels/5b793604-fc76-4886-a834-ed522812cdcb",
          "rel": "bookmark"
        }
      ],
      "http_proxy": null,
      "updated_at": null,
      "fixed_subnet": null,
      "master_flavor_id": null,
      "uuid": "5b793604-fc76-4886-a834-ed522812cdcb",
      "no_proxy": null,
      "https_proxy": null,
      "tls_disabled": false,
      "keypair_id": "testkey",
      "public": false,
      "labels": {},
      "docker_volume_size": 5,
      "server_type": "vm",
      "external_network_id": "public",
      "cluster_distro": "fedora-atomic",
      "image_id": "fedora-atomic-latest",
      "volume_driver": null,
      "registry_enabled": false,
      "docker_storage_driver": null,
      "apiserver_port": null,
      "name": "k8sbaymodel",
      "created_at": "2016-07-06T19:01:31+00:00",
      "network_driver": "flannel",
      "fixed_network": null,
      "coe": "kubernetes",
      "flavor_id": "m1.small",
      "master_lb_enabled": false,
      "dns_nameserver": "8.8.8.8"
    }
  ]
}`)
	})

	client := fake.ServiceClient()
	count := 0

	results := baymodels.List(client, baymodels.ListOpts{})

	err := results.EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := baymodels.ExtractBayModels(page)
		if err != nil {
			t.Errorf("Failed to extract baymodels: %v", err)
			return false, err
		}

		expected := []baymodels.BayModel{
			{
				Name:        "k8sbaymodel",
				ID:          "5b793604-fc76-4886-a834-ed522812cdcb",
				COE:         "kubernetes",
				FlavorID:    "m1.small",
				ImageID:     "fedora-atomic-latest",
				KeyPairID:   "testkey",
				TLSDisabled: false,
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
