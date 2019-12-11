package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/nodegroups"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

const (
	clusterUUID    = "bda75056-3a57-4ada-b943-658ac27beea0"
	badClusterUUID = "252e2f37-d83e-4848-be39-eed1b41211ac"

	nodeGroup1UUID   = "b2e581be-2eec-45b8-921a-c85fbc23aaa3"
	nodeGroup2UUID   = "2457febf-520f-4be3-abb9-96b892d7b5a0"
	badNodeGroupUUID = "4973f3aa-40a2-4857-bf9e-c15faffb08c8"
)

var (
	nodeGroup1Created, _ = time.Parse(time.RFC3339, "2019-10-18T14:03:37+00:00")
	nodeGroup1Updated, _ = time.Parse(time.RFC3339, "2019-10-18T14:18:35+00:00")
)

var expectedNodeGroup1 = nodegroups.NodeGroup{
	ID:               9,
	UUID:             nodeGroup1UUID,
	Name:             "default-master",
	ClusterID:        clusterUUID,
	ProjectID:        "e91d02d561374de6b49960a27b3f08d0",
	DockerVolumeSize: nil,
	Labels: map[string]string{
		"kube_tag": "v1.14.7",
	},
	Links: []gophercloud.Link{
		{
			Href: "http://123.456.789.0:9511/v1/clusters/bda75056-3a57-4ada-b943-658ac27beea0/nodegroups/b2e581be-2eec-45b8-921a-c85fbc23aaa3",
			Rel:  "self",
		},
		{
			Href: "http://123.456.789.0:9511/clusters/bda75056-3a57-4ada-b943-658ac27beea0/nodegroups/b2e581be-2eec-45b8-921a-c85fbc23aaa3",
			Rel:  "bookmark",
		},
	},
	FlavorID:      "",
	ImageID:       "Fedora-AtomicHost-29-20190820.0.x86_64",
	NodeAddresses: []string{"172.24.4.19"},
	NodeCount:     1,
	Role:          "master",
	MinNodeCount:  1,
	MaxNodeCount:  nil,
	IsDefault:     true,
	StackID:       "3cd55bb0-1115-4838-8eca-cefc13f7a21b",
	Status:        "UPDATE_COMPLETE",
	StatusReason:  "Stack UPDATE completed successfully",
	Version:       "",
	CreatedAt:     nodeGroup1Created,
	UpdatedAt:     nodeGroup1Updated,
}

func handleGetNodeGroupSuccess(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID+"/nodegroups/"+nodeGroup1UUID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, nodeGroupGetResponse)
	})
}

func handleGetNodeGroupNotFound(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID+"/nodegroups/"+badNodeGroupUUID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprint(w, nodeGroupGetNotFoundResponse)
	})
}

func handleGetNodeGroupClusterNotFound(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+badClusterUUID+"/nodegroups/"+badNodeGroupUUID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprint(w, nodeGroupGetClusterNotFoundResponse)
	})
}

func handleListNodeGroupsSuccess(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID+"/nodegroups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, nodeGroupListResponse)
	})
}

func handleListNodeGroupsLimitSuccess(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+clusterUUID+"/nodegroups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		r.ParseForm()
		if marker, ok := r.Form["marker"]; !ok {
			// No marker, this is the first request.
			th.TestFormValues(t, r, map[string]string{"limit": "1"})
			fmt.Fprintf(w, nodeGroupListLimitResponse1, th.Endpoint())
		} else {
			switch marker[0] {
			case nodeGroup1UUID:
				// Marker is the UUID of the first node group, return the second.
				fmt.Fprintf(w, nodeGroupListLimitResponse2, th.Endpoint())
			case nodeGroup2UUID:
				// Marker is the UUID of the second node group, there are no more to return.
				fmt.Fprint(w, nodeGroupListLimitResponse3)
			}
		}
	})
}

func handleListNodeGroupsClusterNotFound(t *testing.T) {
	th.Mux.HandleFunc("/v1/clusters/"+badClusterUUID+"/nodegroups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprint(w, nodeGroupListClusterNotFoundResponse)
	})
}

var nodeGroupGetResponse = fmt.Sprintf(`
{
  "links":[
    {
      "href":"http://123.456.789.0:9511/v1/clusters/bda75056-3a57-4ada-b943-658ac27beea0/nodegroups/b2e581be-2eec-45b8-921a-c85fbc23aaa3",
      "rel":"self"
    },
    {
      "href":"http://123.456.789.0:9511/clusters/bda75056-3a57-4ada-b943-658ac27beea0/nodegroups/b2e581be-2eec-45b8-921a-c85fbc23aaa3",
      "rel":"bookmark"
    }
  ],
  "labels":{
    "kube_tag":"v1.14.7"
  },
  "updated_at":"2019-10-18T14:18:35+00:00",
  "cluster_id":"%s",
  "min_node_count":1,
  "id":9,
  "uuid":"%s",
  "version":null,
  "role":"master",
  "node_count":1,
  "project_id":"e91d02d561374de6b49960a27b3f08d0",
  "status":"UPDATE_COMPLETE",
  "docker_volume_size":null,
  "max_node_count":null,
  "is_default":true,
  "image_id":"Fedora-AtomicHost-29-20190820.0.x86_64",
  "node_addresses":[
    "172.24.4.19"
  ],
  "status_reason":"Stack UPDATE completed successfully",
  "name":"default-master",
  "stack_id":"3cd55bb0-1115-4838-8eca-cefc13f7a21b",
  "created_at":"2019-10-18T14:03:37+00:00",
  "flavor_id":null
}`, clusterUUID, nodeGroup1UUID)

// nodeGroupGetNotFoundResponse is the returned error when there is a cluster with the requested ID but it does not have the requested node group.
var nodeGroupGetNotFoundResponse = fmt.Sprintf(`
{
  "errors":[
    {
      "status":404,
      "code":"client",
      "links":[

      ],
      "title":"Nodegroup %s could not be found",
      "request_id":""
    }
  ]
}`, badNodeGroupUUID)

// nodeGroupGetClusterNotFoundResponse is the returned error when there is no cluster with the requested ID.
var nodeGroupGetClusterNotFoundResponse = fmt.Sprintf(`
{
  "errors":[
    {
      "status":404,
      "code":"client",
      "links":[

      ],
      "title":"Cluster %s could not be found",
      "request_id":""
    }
  ]
}`, badClusterUUID)

var nodeGroupListResponse = fmt.Sprintf(`
{
  "nodegroups":[
    {
      "status":"UPDATE_COMPLETE",
      "is_default":true,
      "uuid":"%s",
      "max_node_count":null,
      "stack_id":"3cd55bb0-1115-4838-8eca-cefc13f7a21b",
      "min_node_count":1,
      "image_id":"Fedora-AtomicHost-29-20190820.0.x86_64",
      "role":"master",
      "flavor_id":null,
      "node_count":1,
      "name":"default-master"
    },
    {
      "status":"UPDATE_COMPLETE",
      "is_default":true,
      "uuid":"%s",
      "max_node_count":null,
      "stack_id":"3cd55bb0-1115-4838-8eca-cefc13f7a21b",
      "min_node_count":1,
      "image_id":"Fedora-AtomicHost-29-20190820.0.x86_64",
      "role":"worker",
      "flavor_id":"m1.small",
      "node_count":1,
      "name":"default-worker"
    }
  ]
}`, nodeGroup1UUID, nodeGroup2UUID)

// nodeGroupListLimitResponse1 is the first response when requesting the list of node groups with a limit of 1.
// It returns the URL for the next page with the marker of the first node group.
var nodeGroupListLimitResponse1 = fmt.Sprintf(`
{
  "nodegroups":[
    {
      "status":"UPDATE_COMPLETE",
      "is_default":true,
      "name":"default-master",
      "max_node_count":null,
      "stack_id":"3cd55bb0-1115-4838-8eca-cefc13f7a21b",
      "min_node_count":1,
      "image_id":"Fedora-AtomicHost-29-20190820.0.x86_64",
      "cluster_id":"bda75056-3a57-4ada-b943-658ac27beea0",
      "flavor_id":null,
      "role":"master",
      "node_count":1,
      "uuid":"%s"
    }
  ],
  "next":"%%s/v1/clusters/bda75056-3a57-4ada-b943-658ac27beea0/nodegroups?sort_key=id&sort_dir=asc&limit=1&marker=%s"
}`, nodeGroup1UUID, nodeGroup1UUID)

// nodeGroupListLimitResponse2 is returned when making a request to the URL given by "next" in the first response.
var nodeGroupListLimitResponse2 = fmt.Sprintf(`
{
  "nodegroups":[
    {
      "status":"UPDATE_COMPLETE",
      "is_default":true,
      "name":"default-worker",
      "max_node_count":null,
      "stack_id":"3cd55bb0-1115-4838-8eca-cefc13f7a21b",
      "min_node_count":1,
      "image_id":"Fedora-AtomicHost-29-20190820.0.x86_64",
      "cluster_id":"bda75056-3a57-4ada-b943-658ac27beea0",
      "flavor_id":"m1.small",
      "role":"worker",
      "node_count":1,
      "uuid":"%s"
    }
  ],
  "next":"%%s/v1/clusters/bda75056-3a57-4ada-b943-658ac27beea0/nodegroups?sort_key=id&sort_dir=asc&limit=1&marker=%s"
}`, nodeGroup2UUID, nodeGroup2UUID)

// nodeGroupListLimitResponse3 is the response when listing node groups using a marker and all node groups have already been returned.
var nodeGroupListLimitResponse3 = `{"nodegroups": []}`

// nodeGroupListClusterNotFoundResponse is the error returned when the list operation fails because there is no cluster with the requested ID.
var nodeGroupListClusterNotFoundResponse = fmt.Sprintf(`
{
  "errors":[
    {
      "status":404,
      "code":"client",
      "links":[

      ],
      "title":"Cluster %s could not be found",
      "request_id":""
    }
  ]
}`, badClusterUUID)
