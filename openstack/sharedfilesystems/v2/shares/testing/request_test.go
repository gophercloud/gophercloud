package testing

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &shares.CreateOpts{Size: 1, Name: "my_test_share", ShareProto: "NFS"}
	n, err := shares.Create(client.ServiceClient(), options).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Name, "my_test_share")
	th.AssertEquals(t, n.Size, 1)
	th.AssertEquals(t, n.ShareProto, "NFS")
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	result := shares.Delete(client.ServiceClient(), shareID)
	th.AssertNoErr(t, result.Err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	s, err := shares.Get(client.ServiceClient(), shareID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s, &shares.Share{
		AvailabilityZone:   "nova",
		ShareNetworkID:     "713df749-aac0-4a54-af52-10f6c991e80c",
		ShareServerID:      "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
		SnapshotID:         "",
		ID:                 shareID,
		Size:               1,
		ShareType:          "25747776-08e5-494f-ab40-a64b9d20d8f7",
		ShareTypeName:      "default",
		ConsistencyGroupID: "9397c191-8427-4661-a2e8-b23820dc01d4",
		ProjectID:          projectID,
		Metadata: map[string]string{
			"project": "my_app",
			"aim":     "doc",
		},
		Status:                   "available",
		Description:              "My custom share London",
		Host:                     "manila2@generic1#GENERIC1",
		HasReplicas:              false,
		ReplicationType:          "",
		TaskState:                "",
		SnapshotSupport:          true,
		Name:                     "my_test_share",
		CreatedAt:                time.Date(2015, time.September, 18, 10, 25, 24, 0, time.UTC),
		ShareProto:               "NFS",
		VolumeType:               "default",
		SourceCgsnapshotMemberID: "",
		IsPublic:                 true,
		Links: []map[string]string{
			{
				"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
				"rel":  "self",
			},
			{
				"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
				"rel":  "bookmark",
			},
		},
	})
}

func TestGetMicroversion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetMicroversionResponse(t)

	c := client.ServiceClient()
	// Adding a pattern
	c.Endpoint += "v2/"
	// Adding projectID that is removed by getMicroversionsURL
	c.Endpoint += projectID

	s, err := shares.GetMicroversion(c).ExtractMicroversion()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s, &[]shares.Version{
		{
			Status:  "CURRENT",
			Updated: time.Date(2015, time.August, 27, 11, 33, 21, 0, time.UTC),
			Links: []shares.Link{
				{
					HRef: "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
					Rel:  "self",
				},
				{
					HRef: "http://docs.openstack.org/",
					Type: "text/html",
					Rel:  "describedby",
				},
			},
			MinVersion: "2.0",
			Version:    "2.15",
			MediaTypes: []shares.MediaType{
				{
					Base: "application/json",
					Type: "application/vnd.openstack.share+json;version=1",
				},
			},
			ID: "v2.0",
		},
	})
}
