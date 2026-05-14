package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGetQuotaSet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetQuotaSetSuccessfully(t)

	actual, err := quotasets.Get(client.ServiceClient(), tenantID).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ExpectedInitialQuotaSet, actual)
}

func TestUpdateQuotaSet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateQuotaSetSuccessfully(t)

	actual, err := quotasets.Update(client.ServiceClient(), tenantID, quotasets.UpdateOpts{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}).Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedUpdatedQuotaSet, actual)
}

func TestGetByShareType(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetByShareTypeSuccessfully(t)

	actual, err := quotasets.GetByShareType(client.ServiceClient(), tenantID, ShareType).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ExpectedInitialQuotaSet, actual)
}

func TestUpdateByShareType(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateByShareTypeSuccessfully(t)

	actual, err := quotasets.UpdateByShareType(client.ServiceClient(), tenantID, ShareType, quotasets.UpdateOpts{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ExpectedUpdatedQuotaSet, actual)
}

func TestGetByUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetByUserSuccessfully(t)

	actual, err := quotasets.GetByUser(client.ServiceClient(), tenantID, userID).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ExpectedInitialQuotaSet, actual)
}

func TestUpdateByUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateByUserSuccessfully(t)

	actual, err := quotasets.UpdateByUser(client.ServiceClient(), tenantID, userID, quotasets.UpdateOpts{
		Gigabytes:             gophercloud.IntToPointer(100),
		Snapshots:             gophercloud.IntToPointer(100),
		Shares:                gophercloud.IntToPointer(100),
		SnapshotGigabytes:     gophercloud.IntToPointer(100),
		ShareNetworks:         gophercloud.IntToPointer(100),
		ShareGroups:           gophercloud.IntToPointer(100),
		ShareGroupSnapshots:   gophercloud.IntToPointer(100),
		ShareReplicas:         gophercloud.IntToPointer(100),
		ShareReplicaGigabytes: gophercloud.IntToPointer(100),
		PerShareGigabytes:     gophercloud.IntToPointer(100),
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ExpectedUpdatedQuotaSet, actual)
}
