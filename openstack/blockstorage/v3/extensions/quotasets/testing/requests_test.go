package testing

import (
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

var emptyQuotaSet = quotasets.QuotaSet{}
var emptyQuotaDetailSet = quotasets.QuotaDetailSet{}
var emptyUpdateOpts = quotasets.UpdateOpts{}

func testSuccessTestCase(t *testing.T, jsonBody,
	uriPath, httpMethod string,
	updateOpts quotasets.UpdateOpts,
	expectedQuotaSet quotasets.QuotaSet,
	expectedQuotaDetailSet quotasets.QuotaDetailSet) error {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSuccessfulRequest(t, httpMethod, uriPath, jsonBody)

	if updateOpts != emptyUpdateOpts {
		actual, err := quotasets.Update(client.ServiceClient(),
			FirstTenantID, updateOpts).Extract()
		if err != nil {
			return err
		}
		th.CheckDeepEquals(t, &expectedQuotaSet, actual)
	} else if expectedQuotaSet != emptyQuotaSet {
		actual, err := quotasets.Get(client.ServiceClient(),
			FirstTenantID).Extract()
		if err != nil {
			return err
		}
		th.CheckDeepEquals(t, &expectedQuotaSet, actual)
	} else if expectedQuotaDetailSet != emptyQuotaDetailSet {
		actual, err := quotasets.GetDetail(client.ServiceClient(),
			FirstTenantID).Extract()
		if err != nil {
			return err
		}
		th.CheckDeepEquals(t, expectedQuotaDetailSet, actual)
	}
	return nil
}

func TestSuccessTestCases(t *testing.T) {
	for _, tt := range successTestCases {
		err := testSuccessTestCase(t, tt.jsonBody, tt.uriPath, tt.httpMethod,
			tt.updateOpts, tt.expectedQuotaSet, tt.expectedQuotaDetailSet)
		if err != nil {
			t.Fatalf("Test case '%s' failed with error:\n%s", tt.name, err)
		}
	}
}

type ErrorUpdateOpts quotasets.UpdateOpts

func (opts ErrorUpdateOpts) ToBlockStorageQuotaUpdateMap() (map[string]interface{}, error) {
	return nil, errors.New("This is an error")
}

func TestErrorInToBlockStorageQuotaUpdateMap(t *testing.T) {
	opts := &ErrorUpdateOpts{}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSuccessfulRequest(t, "PUT", "/os-quota-sets/"+FirstTenantID, "")
	_, err := quotasets.Update(client.ServiceClient(), FirstTenantID, opts).Extract()
	if err == nil {
		t.Fatal("Error handling failed")
	}
}
