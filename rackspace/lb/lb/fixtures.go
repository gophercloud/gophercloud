package lb

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

func mockListLBResponse(t *testing.T) {
	th.Mux.HandleFunc("/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "loadBalancers": [
    {
      "name": "lb-site1",
      "id": 71,
      "protocol": "HTTP",
      "port": 80,
      "algorithm": "RANDOM",
      "status": "ACTIVE",
      "nodeCount": 3,
      "virtualIps": [
        {
          "id": 403,
          "address": "206.55.130.1",
          "type": "PUBLIC",
          "ipVersion": "IPV4"
        }
      ],
      "created": {
        "time": "2010-11-30T03:23:42Z"
      },
      "updated": {
        "time": "2010-11-30T03:23:44Z"
      }
    }
  ]
}
  `)
	})
}

func mockCreateLBResponse(t *testing.T) {
	th.Mux.HandleFunc("/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		th.TestJSONRequest(t, r, `
{
  "loadBalancer": {
    "name": "a-new-loadbalancer",
    "port": 80,
    "protocol": "HTTP",
    "virtualIps": [
      {
        "id": 2341
      },
      {
        "id": 900001
      }
    ],
    "nodes": [
      {
        "address": "10.1.1.1",
        "port": 80,
        "condition": "ENABLED"
      }
    ]
  }
}
		`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "loadBalancer": {
    "name": "a-new-loadbalancer",
    "id": 144,
    "protocol": "HTTP",
    "halfClosed": false,
    "port": 83,
    "algorithm": "RANDOM",
    "status": "BUILD",
    "timeout": 30,
    "cluster": {
      "name": "ztm-n01.staging1.lbaas.rackspace.net"
    },
    "nodes": [
      {
        "address": "10.1.1.1",
        "id": 653,
        "port": 80,
        "status": "ONLINE",
        "condition": "ENABLED",
        "weight": 1
      }
    ],
    "virtualIps": [
      {
        "address": "206.10.10.210",
        "id": 39,
        "type": "PUBLIC",
        "ipVersion": "IPV4"
      },
      {
        "address": "2001:4801:79f1:0002:711b:be4c:0000:0021",
        "id": 900001,
        "type": "PUBLIC",
        "ipVersion": "IPV6"
      }
    ],
    "created": {
      "time": "2011-04-13T14:18:07Z"
    },
    "updated": {
      "time": "2011-04-13T14:18:07Z"
    },
    "connectionLogging": {
      "enabled": false
    }
  }
}
	`)
	})
}

func mockDeleteLBResponse(t *testing.T) {
	th.Mux.HandleFunc("/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}
