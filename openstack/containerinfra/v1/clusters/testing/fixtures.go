package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

const clusterUUID = "746e779a-751a-456b-a3e9-c883d734946f"
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
	KeyPair:         "my-keypair",
	MasterAddresses: []string{"172.24.4.6"},
	MasterCount:     1,
	Name:            "k8s",
	NodeAddresses:   []string{"172.24.4.13"},
	NodeCount:       1,
	StackID:         "9c6f1169-7300-4d08-a444-d2be38758719",
	Status:          "CREATE_COMPLETE",
	StatusReason:    "Stack CREATE completed successfully",
	UpdatedAt:       time.Date(2016, 8, 29, 6, 53, 24, 0, time.UTC),
	UUID:            clusterUUID,
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
