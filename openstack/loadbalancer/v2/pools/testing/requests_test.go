package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListPools(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePoolListSuccessfully(t, fakeServer)

	pages := 0
	err := pools.List(fake.ServiceClient(fakeServer), pools.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := pools.ExtractPools(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 pools, got %d", len(actual))
		}
		th.CheckDeepEquals(t, PoolWeb, actual[0])
		th.CheckDeepEquals(t, PoolDb, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllPools(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePoolListSuccessfully(t, fakeServer)

	allPages, err := pools.List(fake.ServiceClient(fakeServer), pools.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := pools.ExtractPools(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, PoolWeb, actual[0])
	th.CheckDeepEquals(t, PoolDb, actual[1])
}

func TestCreatePool(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePoolCreationSuccessfully(t, fakeServer, SinglePoolBody)

	actual, err := pools.Create(context.TODO(), fake.ServiceClient(fakeServer), pools.CreateOpts{
		LBMethod:       pools.LBMethodRoundRobin,
		Protocol:       "HTTP",
		Name:           "Example pool",
		ProjectID:      "2ffc6e22aae24e4795f87155d24c896f",
		LoadbalancerID: "79e05663-7f03-45d2-a092-8b94062f22ab",
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, PoolDb, *actual)
}

func TestGetPool(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePoolGetSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := pools.Get(context.TODO(), client, "c3741b06-df4d-4715-b142-276b6bce75ab").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, PoolDb, *actual)
}

func TestDeletePool(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePoolDeletionSuccessfully(t, fakeServer)

	res := pools.Delete(context.TODO(), fake.ServiceClient(fakeServer), "c3741b06-df4d-4715-b142-276b6bce75ab")
	th.AssertNoErr(t, res.Err)
}

func TestUpdatePool(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePoolUpdateSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	name := "NewPoolName"
	actual, err := pools.Update(context.TODO(), client, "c3741b06-df4d-4715-b142-276b6bce75ab", pools.UpdateOpts{
		Name:     &name,
		LBMethod: pools.LBMethodLeastConnections,
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, PoolUpdated, *actual)
}

func TestRequiredPoolCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := pools.Create(context.TODO(), fake.ServiceClient(fakeServer), pools.CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = pools.Create(context.TODO(), fake.ServiceClient(fakeServer), pools.CreateOpts{
		LBMethod:       pools.LBMethod("invalid"),
		Protocol:       pools.ProtocolHTTPS,
		LoadbalancerID: "69055154-f603-4a28-8951-7cc2d9e54a9a",
	})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}

	res = pools.Create(context.TODO(), fake.ServiceClient(fakeServer), pools.CreateOpts{
		LBMethod:       pools.LBMethodRoundRobin,
		Protocol:       pools.Protocol("invalid"),
		LoadbalancerID: "69055154-f603-4a28-8951-7cc2d9e54a9a",
	})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}

	res = pools.Create(context.TODO(), fake.ServiceClient(fakeServer), pools.CreateOpts{
		LBMethod: pools.LBMethodRoundRobin,
		Protocol: pools.ProtocolHTTPS,
	})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}
}

func TestListMembers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleMemberListSuccessfully(t, fakeServer)

	pages := 0
	err := pools.ListMembers(fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", pools.ListMembersOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := pools.ExtractMembers(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 members, got %d", len(actual))
		}
		th.CheckDeepEquals(t, MemberWeb, actual[0])
		th.CheckDeepEquals(t, MemberDb, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllMembers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleMemberListSuccessfully(t, fakeServer)

	allPages, err := pools.ListMembers(fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", pools.ListMembersOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := pools.ExtractMembers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, MemberWeb, actual[0])
	th.CheckDeepEquals(t, MemberDb, actual[1])
}

func TestCreateMember(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleMemberCreationSuccessfully(t, fakeServer, SingleMemberBody)

	weight := 10
	actual, err := pools.CreateMember(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", pools.CreateMemberOpts{
		Name:         "db",
		SubnetID:     "1981f108-3c48-48d2-b908-30f7d28532c9",
		ProjectID:    "2ffc6e22aae24e4795f87155d24c896f",
		Address:      "10.0.2.11",
		ProtocolPort: 80,
		Weight:       &weight,
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, MemberDb, *actual)
}

func TestRequiredMemberCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := pools.CreateMember(context.TODO(), fake.ServiceClient(fakeServer), "", pools.CreateMemberOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = pools.CreateMember(context.TODO(), fake.ServiceClient(fakeServer), "", pools.CreateMemberOpts{Address: "1.2.3.4", ProtocolPort: 80})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}
	res = pools.CreateMember(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", pools.CreateMemberOpts{ProtocolPort: 80})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}
	res = pools.CreateMember(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", pools.CreateMemberOpts{Address: "1.2.3.4"})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}
}

func TestGetMember(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleMemberGetSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := pools.GetMember(context.TODO(), client, "332abe93-f488-41ba-870b-2ac66be7f853", "2a280670-c202-4b0b-a562-34077415aabf").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, MemberDb, *actual)
}

func TestDeleteMember(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleMemberDeletionSuccessfully(t, fakeServer)

	res := pools.DeleteMember(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", "2a280670-c202-4b0b-a562-34077415aabf")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateMember(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleMemberUpdateSuccessfully(t, fakeServer)

	weight := 4
	client := fake.ServiceClient(fakeServer)
	name := "newMemberName"
	actual, err := pools.UpdateMember(context.TODO(), client, "332abe93-f488-41ba-870b-2ac66be7f853", "2a280670-c202-4b0b-a562-34077415aabf", pools.UpdateMemberOpts{
		Name:   &name,
		Weight: &weight,
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, MemberUpdated, *actual)
}

func TestBatchUpdateMembers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleMembersUpdateSuccessfully(t, fakeServer)

	name_1 := "web-server-1"
	weight_1 := 20
	subnetID := "bbb35f84-35cc-4b2f-84c2-a6a29bba68aa"
	member1 := pools.BatchUpdateMemberOpts{
		Address:      "192.0.2.16",
		ProtocolPort: 80,
		Name:         &name_1,
		SubnetID:     &subnetID,
		Weight:       &weight_1,
	}

	name_2 := "web-server-2"
	weight_2 := 10
	member2 := pools.BatchUpdateMemberOpts{
		Address:      "192.0.2.17",
		ProtocolPort: 80,
		Name:         &name_2,
		Weight:       &weight_2,
		SubnetID:     &subnetID,
	}
	members := []pools.BatchUpdateMemberOpts{member1, member2}

	res := pools.BatchUpdateMembers(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", members)
	th.AssertNoErr(t, res.Err)
}

func TestEmptyBatchUpdateMembers(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleEmptyMembersUpdateSuccessfully(t, fakeServer)

	res := pools.BatchUpdateMembers(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", []pools.BatchUpdateMemberOpts{})
	th.AssertNoErr(t, res.Err)
}

func TestRequiredBatchUpdateMemberOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	name := "web-server-1"
	res := pools.BatchUpdateMembers(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", []pools.BatchUpdateMemberOpts{
		{
			Name: &name,
		},
	})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}

	res = pools.BatchUpdateMembers(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", []pools.BatchUpdateMemberOpts{
		{
			Address: "192.0.2.17",
			Name:    &name,
		},
	})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}

	res = pools.BatchUpdateMembers(context.TODO(), fake.ServiceClient(fakeServer), "332abe93-f488-41ba-870b-2ac66be7f853", []pools.BatchUpdateMemberOpts{
		{
			ProtocolPort: 80,
			Name:         &name,
		},
	})
	if res.Err == nil {
		t.Fatalf("Expected error, but got none")
	}
}
