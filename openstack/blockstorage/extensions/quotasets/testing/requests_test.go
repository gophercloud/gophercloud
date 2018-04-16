package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

var emptyQuotaSet = quotasets.QuotaSet{}
var emptyQuotaUsageSet = quotasets.QuotaUsageSet{}

func testSuccessTestCase(t *testing.T, httpMethod, uriPath, jsonBody string, uriQueryParams map[string]string, expectedQuotaSet quotasets.QuotaSet, expectedQuotaUsageSet quotasets.QuotaUsageSet) error {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSuccessfulRequest(t, httpMethod, uriPath, jsonBody, uriQueryParams)

	if expectedQuotaSet != emptyQuotaSet {
		actual, err := quotasets.Get(client.ServiceClient(), FirstTenantID).Extract()
		if err != nil {
			return err
		}
		th.CheckDeepEquals(t, &expectedQuotaSet, actual)
	} else if expectedQuotaUsageSet != emptyQuotaUsageSet {
		actual, err := quotasets.GetUsage(client.ServiceClient(), FirstTenantID).Extract()
		if err != nil {
			return err
		}
		th.CheckDeepEquals(t, expectedQuotaUsageSet, actual)
	}
	return nil
}

func TestSuccessTestCases(t *testing.T) {
	for _, tt := range successTestCases {
		err := testSuccessTestCase(t, tt.httpMethod, tt.uriPath, tt.jsonBody, tt.uriQueryParams, tt.expectedQuotaSet, tt.expectedQuotaUsageSet)
		if err != nil {
			t.Fatalf("Test case '%s' failed with error:\n%s", tt.name, err)
		}
	}
}
