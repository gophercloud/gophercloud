package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockCreateResponse(t, fakeServer)

	options := &shares.CreateOpts{
		Size:       1,
		Name:       "my_test_share",
		ShareProto: "NFS",
		SchedulerHints: &shares.SchedulerHints{
			SameHost:      "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
			DifferentHost: "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
		},
	}
	n, err := shares.Create(context.TODO(), client.ServiceClient(fakeServer), options).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Name, "my_test_share")
	th.AssertEquals(t, n.Size, 1)
	th.AssertEquals(t, n.ShareProto, "NFS")
	th.AssertEquals(t, n.Metadata["__affinity_same_host"], "e268f4aa-d571-43dd-9ab3-f49ad06ffaef")
	th.AssertEquals(t, n.Metadata["__affinity_different_host"], "e268f4aa-d571-43dd-9ab3-f49ad06ffaef")
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockUpdateResponse(t, fakeServer)

	name := "my_new_test_share"
	description := ""
	iFalse := false
	options := &shares.UpdateOpts{
		DisplayName:        &name,
		DisplayDescription: &description,
		IsPublic:           &iFalse,
	}
	n, err := shares.Update(context.TODO(), client.ServiceClient(fakeServer), shareID, options).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, n.Name, "my_new_test_share")
	th.AssertEquals(t, n.Description, "")
	th.AssertEquals(t, n.IsPublic, false)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockDeleteResponse(t, fakeServer)

	result := shares.Delete(context.TODO(), client.ServiceClient(fakeServer), shareID)
	th.AssertNoErr(t, result.Err)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetResponse(t, fakeServer)

	s, err := shares.Get(context.TODO(), client.ServiceClient(fakeServer), shareID).Extract()
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
		Status:                         "available",
		Description:                    "My custom share London",
		Host:                           "manila2@generic1#GENERIC1",
		HasReplicas:                    false,
		ReplicationType:                "",
		TaskState:                      "",
		SnapshotSupport:                true,
		CreateShareFromSnapshotSupport: true,
		Name:                           "my_test_share",
		CreatedAt:                      time.Date(2015, time.September, 18, 10, 25, 24, 0, time.UTC),
		ShareProto:                     "NFS",
		VolumeType:                     "default",
		SourceCgsnapshotMemberID:       "",
		IsPublic:                       true,
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

func TestListDetail(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockListDetailResponse(t, fakeServer)

	allPages, err := shares.ListDetail(client.ServiceClient(fakeServer), &shares.ListOpts{}).AllPages(context.TODO())

	th.AssertNoErr(t, err)

	actual, err := shares.ExtractShares(allPages)
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, actual, []shares.Share{
		{
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
			Status:                         "available",
			Description:                    "My custom share London",
			Host:                           "manila2@generic1#GENERIC1",
			HasReplicas:                    false,
			ReplicationType:                "",
			TaskState:                      "",
			SnapshotSupport:                true,
			CreateShareFromSnapshotSupport: true,
			Name:                           "my_test_share",
			CreatedAt:                      time.Date(2015, time.September, 18, 10, 25, 24, 0, time.UTC),
			ShareProto:                     "NFS",
			VolumeType:                     "default",
			SourceCgsnapshotMemberID:       "",
			IsPublic:                       true,
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
		},
	})
}

func TestListExportLocationsSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockListExportLocationsResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for List Export Locations is 2.9
	c.Microversion = "2.9"

	s, err := shares.ListExportLocations(context.TODO(), c, shareID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s, []shares.ExportLocation{
		{
			Path:            "127.0.0.1:/var/lib/manila/mnt/share-9a922036-ad26-4d27-b955-7a1e285fa74d",
			ShareInstanceID: "011d21e2-fbc3-4e4a-9993-9ea223f73264",
			IsAdminOnly:     false,
			ID:              "80ed63fc-83bc-4afc-b881-da4a345ac83d",
			Preferred:       false,
		},
	})
}

func TestGetExportLocationSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetExportLocationResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for Get Export Location is 2.9
	c.Microversion = "2.9"

	s, err := shares.GetExportLocation(context.TODO(), c, shareID, "80ed63fc-83bc-4afc-b881-da4a345ac83d").Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s, &shares.ExportLocation{
		Path:            "127.0.0.1:/var/lib/manila/mnt/share-9a922036-ad26-4d27-b955-7a1e285fa74d",
		ShareInstanceID: "011d21e2-fbc3-4e4a-9993-9ea223f73264",
		IsAdminOnly:     false,
		ID:              "80ed63fc-83bc-4afc-b881-da4a345ac83d",
		Preferred:       false,
	})
}

func TestGrantAcessSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGrantAccessResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for Grant Access is 2.7
	c.Microversion = "2.7"

	var grantAccessReq shares.GrantAccessOpts
	grantAccessReq.AccessType = "ip"
	grantAccessReq.AccessTo = "0.0.0.0/0"
	grantAccessReq.AccessLevel = "rw"

	s, err := shares.GrantAccess(context.TODO(), c, shareID, grantAccessReq).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s, &shares.AccessRight{
		ShareID:     "011d21e2-fbc3-4e4a-9993-9ea223f73264",
		AccessType:  "ip",
		AccessTo:    "0.0.0.0/0",
		AccessKey:   "",
		AccessLevel: "rw",
		State:       "new",
		ID:          "a2f226a5-cee8-430b-8a03-78a59bd84ee8",
	})
}

func TestRevokeAccessSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockRevokeAccessResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for Revoke Access is 2.7
	c.Microversion = "2.7"

	options := &shares.RevokeAccessOpts{AccessID: "a2f226a5-cee8-430b-8a03-78a59bd84ee8"}

	err := shares.RevokeAccess(context.TODO(), c, shareID, options).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListAccessRightsSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockListAccessRightsResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for Grant Access is 2.7
	c.Microversion = "2.7"

	s, err := shares.ListAccessRights(context.TODO(), c, shareID).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s, []shares.AccessRight{
		{
			ShareID:     "011d21e2-fbc3-4e4a-9993-9ea223f73264",
			AccessType:  "ip",
			AccessTo:    "0.0.0.0/0",
			AccessKey:   "",
			AccessLevel: "rw",
			State:       "new",
			ID:          "a2f226a5-cee8-430b-8a03-78a59bd84ee8",
		},
	})
}

func TestExtendSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockExtendResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for Grant Access is 2.7
	c.Microversion = "2.7"

	err := shares.Extend(context.TODO(), c, shareID, &shares.ExtendOpts{NewSize: 2}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestShrinkSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockShrinkResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for Grant Access is 2.7
	c.Microversion = "2.7"

	err := shares.Shrink(context.TODO(), c, shareID, &shares.ShrinkOpts{NewSize: 1}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetMetadataSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetMetadataResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)

	actual, err := shares.GetMetadata(context.TODO(), c, shareID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]string{"foo": "bar"}, actual)
}

func TestGetMetadatumSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetMetadatumResponse(t, fakeServer, "foo")

	c := client.ServiceClient(fakeServer)

	actual, err := shares.GetMetadatum(context.TODO(), c, shareID, "foo").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]string{"foo": "bar"}, actual)
}

func TestSetMetadataSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockSetMetadataResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)

	actual, err := shares.SetMetadata(context.TODO(), c, shareID, &shares.SetMetadataOpts{Metadata: map[string]string{"foo": "bar"}}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]string{"foo": "bar"}, actual)
}

func TestUpdateMetadataSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockUpdateMetadataResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)

	actual, err := shares.UpdateMetadata(context.TODO(), c, shareID, &shares.UpdateMetadataOpts{Metadata: map[string]string{"foo": "bar"}}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]string{"foo": "bar"}, actual)
}

func TestUnsetMetadataSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockDeleteMetadatumResponse(t, fakeServer, "foo")

	c := client.ServiceClient(fakeServer)

	err := shares.DeleteMetadatum(context.TODO(), c, shareID, "foo").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestRevertSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockRevertResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for Revert is 2.27
	c.Microversion = "2.27"

	err := shares.Revert(context.TODO(), c, shareID, &shares.RevertOpts{SnapshotID: "ddeac769-9742-497f-b985-5bcfa94a3fd6"}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestResetStatusSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockResetStatusResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for ResetStatus is 2.7
	c.Microversion = "2.7"

	err := shares.ResetStatus(context.TODO(), c, shareID, &shares.ResetStatusOpts{Status: "error"}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestForceDeleteSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockForceDeleteResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for ForceDelete is 2.7
	c.Microversion = "2.7"

	err := shares.ForceDelete(context.TODO(), c, shareID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUnmanageSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockUnmanageResponse(t, fakeServer)

	c := client.ServiceClient(fakeServer)
	// Client c must have Microversion set; minimum supported microversion for Unmanage is 2.7
	c.Microversion = "2.7"

	err := shares.Unmanage(context.TODO(), c, shareID).ExtractErr()
	th.AssertNoErr(t, err)
}
