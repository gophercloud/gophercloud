package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocationcandidates"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAllocationCandidatesSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocationCandidatesSuccess(t, fakeServer)

	page, err := allocationcandidates.List(client.ServiceClient(fakeServer), allocationcandidates.ListOpts{
		Resources: "VCPU:1,MEMORY_MB:1024",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationCandidates, *actual)
}

func TestListAllocationCandidatesPre134Success(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocationCandidatesPre134Success(t, fakeServer)

	page, err := allocationcandidates.List(client.ServiceClient(fakeServer), allocationcandidates.ListOpts{
		Resources: "VCPU:1,MEMORY_MB:1024",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationCandidatesPre134, *actual)
}

func TestListAllocationCandidatesPre129Success(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocationCandidatesPre129Success(t, fakeServer)

	page, err := allocationcandidates.List(client.ServiceClient(fakeServer), allocationcandidates.ListOpts{
		Resources: "VCPU:1,MEMORY_MB:1024",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationCandidatesPre129, *actual)
}

func TestListAllocationCandidatesPre117Success(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocationCandidatesPre117Success(t, fakeServer)

	page, err := allocationcandidates.List(client.ServiceClient(fakeServer), allocationcandidates.ListOpts{
		Resources: "VCPU:1,MEMORY_MB:1024",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationCandidatesPre117, *actual)
}

func TestListAllocationCandidates110Success(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocationCandidates110Success(t, fakeServer)

	page, err := allocationcandidates.List(client.ServiceClient(fakeServer), allocationcandidates.ListOpts{
		Resources: "VCPU:1,MEMORY_MB:1024",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := allocationcandidates.ExtractAllocationCandidates110(page)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationCandidates110, *actual)
}

func TestListAllocationCandidatesEmptySuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocationCandidatesEmptySuccess(t, fakeServer)

	page, err := allocationcandidates.List(client.ServiceClient(fakeServer), allocationcandidates.ListOpts{
		Resources: "VCPU:1,MEMORY_MB:1024",
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(actual.AllocationRequests))
	th.AssertEquals(t, 0, len(actual.ProviderSummaries))
}

func TestListAllocationCandidatesBadRequest(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocationCandidatesBadRequest(t, fakeServer)

	_, err := allocationcandidates.List(client.ServiceClient(fakeServer), allocationcandidates.ListOpts{
		Resources: "INVALID",
	}).AllPages(context.TODO())
	th.AssertErr(t, err)
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusBadRequest))
}

func TestListAllocationCandidatesWithFullQuerySuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListAllocationCandidatesWithFullQuerySuccess(t, fakeServer)

	page, err := allocationcandidates.List(client.ServiceClient(fakeServer), allocationcandidates.ListOpts{
		Resources:   "VCPU:1,MEMORY_MB:1024",
		Required:    []string{"HW_CPU_X86_SSE", "!HW_CPU_X86_AVX2"},
		Limit:       5,
		GroupPolicy: "isolate",
		ResourceGroups: map[string]allocationcandidates.ResourceGroup{
			"1": {
				Resources: "SRIOV_NET_VF:1",
				Required:  []string{"CUSTOM_PHYSNET1"},
			},
		},
	}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := allocationcandidates.ExtractAllocationCandidates(page)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationCandidates, *actual)
}
