package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const clusterUUID = "746e779a-751a-456b-a3e9-c883d734946f"
const clusterUUID2 = "846e779a-751a-456b-a3e9-c883d734946f"
const requestUUID = "req-781e9bdc-4163-46eb-91c9-786c53188bbb"

var ClusterCreateResponse = fmt.Sprintf(`
										{
											"uuid":"%s"
										}`, clusterUUID)

var ExpectedCluster = clusters.Cluster{
	APIAddress:        "https://172.24.4.6:6443",
	COEVersion:        "v1.2.0",
	ClusterTemplateID: "0562d357-8641-4759-8fed-8173f02c9633",
	CreateTimeout:     60,
	CreatedAt:         time.Date(2016, 8, 29, 6, 51, 31, 0, time.UTC),
	DiscoveryURL:      "https://discovery.etcd.io/cbeb580da58915809d59ee69348a84f3",
	Links: []gophercloud.Link{
		{
			Href: "http://10.164.180.104:9511/v1/clusters/746e779a-751a-456b-a3e9-c883d734946f",
			Rel:  "self",
		},
		{
			Href: "http://10.164.180.104:9511/clusters/746e779a-751a-456b-a3e9-c883d734946f",
			Rel:  "bookmark",
		},
	},
	KeyPair:            "my-keypair",
	MasterAddresses:    []string{"172.24.4.6"},
	MasterCount:        1,
	Name:               "k8s",
	NodeAddresses:      []string{"172.24.4.13"},
	NodeCount:          1,
	StackID:            "9c6f1169-7300-4d08-a444-d2be38758719",
	Status:             "CREATE_COMPLETE",
	StatusReason:       "Stack CREATE completed successfully",
	UpdatedAt:          time.Date(2016, 8, 29, 6, 53, 24, 0, time.UTC),
	UUID:               clusterUUID,
	FloatingIPEnabled:  true,
	MasterLBEnabled:    true,
	FixedNetwork:       "private_network",
	FixedSubnet:        "private_subnet",
	HealthStatus:       "HEALTHY",
	HealthStatusReason: map[string]any{"api": "ok"},
}

var ExpectedCluster2 = clusters.Cluster{
	APIAddress:        "https://172.24.4.6:6443",
	COEVersion:        "v1.2.0",
	ClusterTemplateID: "0562d357-8641-4759-8fed-8173f02c9633",
	CreateTimeout:     60,
	CreatedAt:         time.Time{},
	DiscoveryURL:      "https://discovery.etcd.io/cbeb580da58915809d59ee69348a84f3",
	Links: []gophercloud.Link{
		{
			Href: "http://10.164.180.104:9511/v1/clusters/746e779a-751a-456b-a3e9-c883d734946f",
			Rel:  "self",
		},
		{
			Href: "http://10.164.180.104:9511/clusters/746e779a-751a-456b-a3e9-c883d734946f",
			Rel:  "bookmark",
		},
	},
	KeyPair:            "my-keypair",
	MasterAddresses:    []string{"172.24.4.6"},
	MasterCount:        1,
	Name:               "k8s",
	NodeAddresses:      []string{"172.24.4.13"},
	NodeCount:          1,
	StackID:            "9c6f1169-7300-4d08-a444-d2be38758719",
	Status:             "CREATE_COMPLETE",
	StatusReason:       "Stack CREATE completed successfully",
	UpdatedAt:          time.Date(2016, 8, 29, 6, 53, 24, 0, time.UTC),
	UUID:               clusterUUID2,
	FloatingIPEnabled:  true,
	MasterLBEnabled:    true,
	FixedNetwork:       "private_network",
	FixedSubnet:        "private_subnet",
	HealthStatus:       "HEALTHY",
	HealthStatusReason: map[string]any{"api": "ok"},
}

var ExpectedClusterUUID = clusterUUID

func HandleCreateClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-Id", requestUUID)
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, ClusterCreateResponse)
	})
}

func HandleGetClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ClusterGetResponse)
	})
}

var ClusterGetResponse = fmt.Sprintf(`
{
		"status":"CREATE_COMPLETE",
		"uuid":"%s",
		"links":[
		  {
			 "href":"http://10.164.180.104:9511/v1/clusters/746e779a-751a-456b-a3e9-c883d734946f",
			 "rel":"self"
		  },
		  {
			 "href":"http://10.164.180.104:9511/clusters/746e779a-751a-456b-a3e9-c883d734946f",
			 "rel":"bookmark"
		  }
		],
		"stack_id":"9c6f1169-7300-4d08-a444-d2be38758719",
		"created_at":"2016-08-29T06:51:31+00:00",
		"api_address":"https://172.24.4.6:6443",
		"discovery_url":"https://discovery.etcd.io/cbeb580da58915809d59ee69348a84f3",
		"updated_at":"2016-08-29T06:53:24+00:00",
		"master_count":1,
		"coe_version": "v1.2.0",
		"keypair":"my-keypair",
		"cluster_template_id":"0562d357-8641-4759-8fed-8173f02c9633",
		"master_addresses":[
		  "172.24.4.6"
		],
		"node_count":1,
		"node_addresses":[
		  "172.24.4.13"
		],
		"status_reason":"Stack CREATE completed successfully",
		"create_timeout":60,
		"name":"k8s",
		"floating_ip_enabled": true,
		"master_lb_enabled": true,
		"fixed_network": "private_network",
		"fixed_subnet": "private_subnet",
		"health_status": "HEALTHY",
		"health_status_reason": {"api": "ok"}
}`, clusterUUID)

var ClusterListResponse = fmt.Sprintf(`
{
	"clusters": [
		{
			"api_address":"https://172.24.4.6:6443",
			"cluster_template_id":"0562d357-8641-4759-8fed-8173f02c9633",
			"coe_version": "v1.2.0",
			"create_timeout":60,
			"created_at":"2016-08-29T06:51:31+00:00",
			"discovery_url":"https://discovery.etcd.io/cbeb580da58915809d59ee69348a84f3",
			"keypair":"my-keypair",
			"links":[
			  {
				 "href":"http://10.164.180.104:9511/v1/clusters/746e779a-751a-456b-a3e9-c883d734946f",
				 "rel":"self"
			  },
			  {
				 "href":"http://10.164.180.104:9511/clusters/746e779a-751a-456b-a3e9-c883d734946f",
				 "rel":"bookmark"
			  }
			],
			"master_addresses":[
			  "172.24.4.6"
			],
			"master_count":1,
			"name":"k8s",
			"node_addresses":[
			  "172.24.4.13"
			],
			"node_count":1,
			"stack_id":"9c6f1169-7300-4d08-a444-d2be38758719",
			"status":"CREATE_COMPLETE",
			"status_reason":"Stack CREATE completed successfully",
			"updated_at":"2016-08-29T06:53:24+00:00",
			"uuid":"%s",
			"floating_ip_enabled": true,
			"master_lb_enabled": true,
			"fixed_network": "private_network",
			"fixed_subnet": "private_subnet",
			"health_status": "HEALTHY",
			"health_status_reason": {"api": "ok"}
		},
		{
			"api_address":"https://172.24.4.6:6443",
			"cluster_template_id":"0562d357-8641-4759-8fed-8173f02c9633",
			"coe_version": "v1.2.0",
			"create_timeout":60,
			"created_at":null,
			"discovery_url":"https://discovery.etcd.io/cbeb580da58915809d59ee69348a84f3",
			"keypair":"my-keypair",
			"links":[
			  {
				 "href":"http://10.164.180.104:9511/v1/clusters/746e779a-751a-456b-a3e9-c883d734946f",
				 "rel":"self"
			  },
			  {
				 "href":"http://10.164.180.104:9511/clusters/746e779a-751a-456b-a3e9-c883d734946f",
				 "rel":"bookmark"
			  }
			],
			"master_addresses":[
			  "172.24.4.6"
			],
			"master_count":1,
			"name":"k8s",
			"node_addresses":[
			  "172.24.4.13"
			],
			"node_count":1,
			"stack_id":"9c6f1169-7300-4d08-a444-d2be38758719",
			"status":"CREATE_COMPLETE",
			"status_reason":"Stack CREATE completed successfully",
			"updated_at":null,
			"uuid":"%s",
			"floating_ip_enabled": true,
			"master_lb_enabled": true,
			"fixed_network": "private_network",
			"fixed_subnet": "private_subnet",
			"health_status": "HEALTHY",
			"health_status_reason": {"api": "ok"}
		}
	]
}`, clusterUUID, clusterUUID2)

var ExpectedClusters = []clusters.Cluster{ExpectedCluster}

func HandleListClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-Id", requestUUID)
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ClusterListResponse)
	})
}

func HandleListDetailClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-Id", requestUUID)
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ClusterListResponse)
	})
}

var UpdateResponse = fmt.Sprintf(`
{
	"uuid":"%s"
}`, clusterUUID)

func HandleUpdateClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-Id", requestUUID)
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})
}

var UpgradeResponse = fmt.Sprintf(`
{
	"uuid":"%s"
}`, clusterUUID)

func HandleUpgradeClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID+"/actions/upgrade", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-Id", requestUUID)
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, UpgradeResponse)
	})
}

func HandleDeleteClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-Id", requestUUID)
		w.WriteHeader(http.StatusNoContent)
	})
}

var ResizeResponse = fmt.Sprintf(`
{
	"uuid": "%s"
}`, clusterUUID)

func HandleResizeClusterSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID+"/actions/resize", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-OpenStack-Request-Id", requestUUID)
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprint(w, ResizeResponse)
	})
}
