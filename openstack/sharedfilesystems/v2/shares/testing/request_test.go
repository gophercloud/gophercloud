package testing

import (
	"testing"
	"time"

	"encoding/json"
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/gophercloud/gophercloud/pagination"
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
		ProjectID:          "16e1ab15c35a457e9c2b2aa189f544e1",
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

func TestListAllShort(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	pages, err := shares.List(client.ServiceClient(), &shares.ListOpts{}, false).AllPages()
	th.AssertNoErr(t, err)
	act, err := shares.ExtractShares(pages)
	th.AssertNoErr(t, err)
	shortList := []shares.Share{
		{
			ID:   "d94a8548-2079-4be0-b21c-0a887acd31ca",
			Name: "My_share",
			Links: []map[string]string{
				{
					"href": "http://172.18.198.54:8786/v1/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
					"rel":  "self",
				},
				{
					"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
					"rel":  "bookmark",
				},
			},
		},
		{
			ID:   "406ea93b-32e9-4907-a117-148b3945749f",
			Name: "Share1",
			Links: []map[string]string{
				{
					"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/d94a8548-2079-4be0-b21c-0a887acd31ca",
					"rel":  "bookmark",
				},
			},
		},
	}
	th.CheckDeepEquals(t, shortList, act)
}

func TestListPaginateDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListDetailResponse(t)

	opts := &shares.ListOpts{
		Offset: 0,
		Limit:  1,
	}
	count := 0
	err := shares.List(client.ServiceClient(), opts, true).EachPage(func(page pagination.Page) (bool, error) {
		s, err := shares.ExtractShares(page)
		if err != nil {
			t.Errorf("Unable to extract shares: %v", err)
			return false, err
		}
		for i, share := range s {
			asJSON, _ := json.MarshalIndent(share, "", " ")
			fmt.Printf("share %d: \n%s\n", i, asJSON)
		}
		count++
		return true, nil
	})
	th.AssertEquals(t, 3, count)
	th.AssertNoErr(t, err)
}
