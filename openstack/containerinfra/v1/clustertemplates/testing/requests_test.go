package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"
	fake "github.com/gophercloud/gophercloud/openstack/containerinfra/v1/common"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clustertemplates", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "clustertemplates": [
    {
      "insecure_registry": null,
      "links": [
        {
          "href": "http://65.61.151.130:9511/v1/clustertemplates/5b793604-fc76-4886-a834-ed522812cdcb",
          "rel": "self"
        },
        {
          "href": "http://65.61.151.130:9511/clustertemplates/5b793604-fc76-4886-a834-ed522812cdcb",
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
      "name": "k8sclustertemplate",
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

	results := clustertemplates.List(client, clustertemplates.ListOpts{})

	err := results.EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clustertemplates.ExtractClusterTemplates(page)
		if err != nil {
			t.Errorf("Failed to extract clustertemplates: %v", err)
			return false, err
		}

		expected := []clustertemplates.ClusterTemplate{
			{
				Name:        "k8sclustertemplate",
				ID:          "5b793604-fc76-4886-a834-ed522812cdcb",
				COE:         "kubernetes",
				ServerType:  "vm",
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

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clustertemplates/kubernetes-dev", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "insecure_registry": null,
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/clustertemplates/472807c2-f175-4946-9765-149701a5aba7",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/clustertemplates/472807c2-f175-4946-9765-149701a5aba7",
      "rel": "bookmark"
    }
  ],
  "http_proxy": null,
  "updated_at": null,
  "fixed_subnet": null,
  "master_flavor_id": null,
  "uuid": "472807c2-f175-4946-9765-149701a5aba7",
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
  "name": "kubernetes-dev",
  "created_at": "2016-08-10T13:47:01+00:00",
  "network_driver": "flannel",
  "fixed_network": null,
  "coe": "kubernetes",
  "flavor_id": "m1.small",
  "master_lb_enabled": false,
  "dns_nameserver": "8.8.8.8"
}
			`)
	})

	m, err := clustertemplates.Get(fake.ServiceClient(), "kubernetes-dev").Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, "kubernetes-dev", m.Name)
	th.AssertEquals(t, "kubernetes", m.COE)
	th.AssertEquals(t, "vm", m.ServerType)
	th.AssertEquals(t, "m1.small", m.FlavorID)
	th.AssertEquals(t, "fedora-atomic-latest", m.ImageID)
	th.AssertEquals(t, "testkey", m.KeyPairID)
	th.AssertEquals(t, false, m.TLSDisabled)
}

func TestGetFailed(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clustertemplates/duplicatename", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, `
{
  "errors": [
    {
      "status": 409,
      "code": "client",
      "links": [],
      "title": "Multiple clustertemplates exist with same name",
      "detail": "Multiple clustertemplates exist with same name. Please use the clustertemplate uuid instead.",
      "request_id": ""
    }
  ]
}
		`)
	})

	res := clustertemplates.Get(fake.ServiceClient(), "duplicatename")

	th.AssertEquals(t, "Multiple clustertemplates exist with same name. Please use the clustertemplate uuid instead.", res.Err.Error())

	er, ok := res.Err.(*fake.ErrorResponse)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, http.StatusConflict, er.Actual)
}
