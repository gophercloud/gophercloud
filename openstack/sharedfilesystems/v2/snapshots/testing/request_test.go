package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/snapshots"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &snapshots.CreateOpts{ShareID: shareID, Name: "test snapshot", Description: "test description"}
	n, err := snapshots.Create(context.TODO(), client.ServiceClient(), options).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Name, "test snapshot")
	th.AssertEquals(t, n.Description, "test description")
	th.AssertEquals(t, n.ShareProto, "NFS")
	th.AssertEquals(t, n.ShareSize, 1)
	th.AssertEquals(t, n.Size, 1)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateResponse(t)

	name := "my_new_test_snapshot"
	description := ""
	options := &snapshots.UpdateOpts{
		DisplayName:        &name,
		DisplayDescription: &description,
	}
	n, err := snapshots.Update(context.TODO(), client.ServiceClient(), snapshotID, options).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Name, "my_new_test_snapshot")
	th.AssertEquals(t, n.Description, "")
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	result := snapshots.Delete(context.TODO(), client.ServiceClient(), snapshotID)
	th.AssertNoErr(t, result.Err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	s, err := snapshots.Get(context.TODO(), client.ServiceClient(), snapshotID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s, &snapshots.Snapshot{
		ID:          snapshotID,
		Name:        "new_app_snapshot",
		Description: "",
		ShareID:     "19865c43-3b91-48c9-85a0-7ac4d6bb0efe",
		ShareProto:  "NFS",
		ShareSize:   1,
		Size:        1,
		Status:      "available",
		ProjectID:   "16e1ab15c35a457e9c2b2aa189f544e1",
		CreatedAt:   time.Date(2019, time.January, 06, 11, 11, 02, 0, time.UTC),
		Links: []map[string]string{
			{
				"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/bc082e99-3bdb-4400-b95e-b85c7a41622c",
				"rel":  "self",
			},
			{
				"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/bc082e99-3bdb-4400-b95e-b85c7a41622c",
				"rel":  "bookmark",
			},
		},
	})
}

func TestListDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListDetailResponse(t)

	allPages, err := snapshots.ListDetail(client.ServiceClient(), &snapshots.ListOpts{}).AllPages(context.TODO())

	th.AssertNoErr(t, err)

	actual, err := snapshots.ExtractSnapshots(allPages)
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, actual, []snapshots.Snapshot{
		{
			ID:          snapshotID,
			Name:        "new_app_snapshot",
			Description: "",
			ShareID:     "19865c43-3b91-48c9-85a0-7ac4d6bb0efe",
			ShareProto:  "NFS",
			ShareSize:   1,
			Size:        1,
			Status:      "available",
			ProjectID:   "16e1ab15c35a457e9c2b2aa189f544e1",
			CreatedAt:   time.Date(2019, time.January, 06, 11, 11, 02, 0, time.UTC),
			Links: []map[string]string{
				{
					"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/bc082e99-3bdb-4400-b95e-b85c7a41622c",
					"rel":  "self",
				},
				{
					"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/snapshots/bc082e99-3bdb-4400-b95e-b85c7a41622c",
					"rel":  "bookmark",
				},
			},
		},
	})
}

func TestResetStatusSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockResetStatusResponse(t)

	c := client.ServiceClient()

	err := snapshots.ResetStatus(context.TODO(), c, snapshotID, &snapshots.ResetStatusOpts{Status: "error"}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestForceDeleteSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockForceDeleteResponse(t)

	c := client.ServiceClient()

	err := snapshots.ForceDelete(context.TODO(), c, snapshotID).ExtractErr()
	th.AssertNoErr(t, err)
}
