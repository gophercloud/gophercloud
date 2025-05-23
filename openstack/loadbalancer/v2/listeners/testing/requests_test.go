package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListListeners(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListenerListSuccessfully(t, fakeServer)

	pages := 0
	err := listeners.List(fake.ServiceClient(fakeServer), listeners.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := listeners.ExtractListeners(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 listeners, got %d", len(actual))
		}
		th.CheckDeepEquals(t, ListenerWeb, actual[0])
		th.CheckDeepEquals(t, ListenerDb, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllListeners(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListenerListSuccessfully(t, fakeServer)

	allPages, err := listeners.List(fake.ServiceClient(fakeServer), listeners.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := listeners.ExtractListeners(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ListenerWeb, actual[0])
	th.CheckDeepEquals(t, ListenerDb, actual[1])
}

func TestCreateListener(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListenerCreationSuccessfully(t, fakeServer, SingleListenerBody)

	actual, err := listeners.Create(context.TODO(), fake.ServiceClient(fakeServer), listeners.CreateOpts{
		Protocol:               "TCP",
		Name:                   "db",
		LoadbalancerID:         "79e05663-7f03-45d2-a092-8b94062f22ab",
		AdminStateUp:           gophercloud.Enabled,
		DefaultTlsContainerRef: "2c433435-20de-4411-84ae-9cc8917def76",
		DefaultPoolID:          "41efe233-7591-43c5-9cf7-923964759f9e",
		ProtocolPort:           3306,
		InsertHeaders:          map[string]string{"X-Forwarded-For": "true"},
		AllowedCIDRs:           []string{"192.0.2.0/24", "198.51.100.0/24"},
		TLSVersions:            []listeners.TLSVersion{"TLSv1.2"},
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, ListenerDb, *actual)
}

func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	res := listeners.Create(context.TODO(), fake.ServiceClient(fakeServer), listeners.CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = listeners.Create(context.TODO(), fake.ServiceClient(fakeServer), listeners.CreateOpts{Name: "foo"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = listeners.Create(context.TODO(), fake.ServiceClient(fakeServer), listeners.CreateOpts{Name: "foo", ProjectID: "bar"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = listeners.Create(context.TODO(), fake.ServiceClient(fakeServer), listeners.CreateOpts{Name: "foo", ProjectID: "bar", Protocol: "bar"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = listeners.Create(context.TODO(), fake.ServiceClient(fakeServer), listeners.CreateOpts{Name: "foo", ProjectID: "bar", Protocol: "bar", ProtocolPort: 80})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGetListener(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListenerGetSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := listeners.Get(context.TODO(), client, "4ec89087-d057-4e2c-911f-60a3b47ee304").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, ListenerDb, *actual)
}

func TestDeleteListener(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListenerDeletionSuccessfully(t, fakeServer)

	res := listeners.Delete(context.TODO(), fake.ServiceClient(fakeServer), "4ec89087-d057-4e2c-911f-60a3b47ee304")
	th.AssertNoErr(t, res.Err)
}

func TestUpdateListener(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListenerUpdateSuccessfully(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	i1001 := 1001
	i181000 := 181000
	name := "NewListenerName"
	defaultPoolID := ""
	insertHeaders := map[string]string{
		"X-Forwarded-For":  "true",
		"X-Forwarded-Port": "false",
	}
	tlsVersions := []listeners.TLSVersion{"TLSv1.2", "TLSv1.3"}
	actual, err := listeners.Update(context.TODO(), client, "4ec89087-d057-4e2c-911f-60a3b47ee304", listeners.UpdateOpts{
		Name:                 &name,
		ConnLimit:            &i1001,
		DefaultPoolID:        &defaultPoolID,
		TimeoutMemberData:    &i181000,
		TimeoutClientData:    &i181000,
		TimeoutMemberConnect: &i181000,
		TimeoutTCPInspect:    &i181000,
		InsertHeaders:        &insertHeaders,
		TLSVersions:          &tlsVersions,
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, ListenerUpdated, *actual)
}

func TestGetListenerStatsTree(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListenerGetStatsTree(t, fakeServer)

	client := fake.ServiceClient(fakeServer)
	actual, err := listeners.GetStats(context.TODO(), client, "4ec89087-d057-4e2c-911f-60a3b47ee304").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, ListenerStatsTree, *actual)
}
