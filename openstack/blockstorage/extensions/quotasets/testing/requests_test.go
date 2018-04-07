package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

var emptyQuotaSet = quotasets.QuotaSet{}

func testSuccessTestCase(t *testing.T, httpMethod, uriPath, jsonBody string, expectedQuotaSet quotasets.QuotaSet) error {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSuccessfulRequest(t, httpMethod, uriPath, jsonBody)

	actual, err := quotasets.Get(client.ServiceClient(), FirstTenantID).Extract()
	if err != nil {
		return err
	}
	th.CheckDeepEquals(t, &expectedQuotaSet, actual)
	return nil
}

func TestSuccessTestCases(t *testing.T) {
	for _, tt := range successTestCases {
		err := testSuccessTestCase(t, tt.httpMethod, tt.uriPath, tt.jsonBody, tt.expectedQuotaSet)
		if err != nil {
			t.Fatalf("Test case '%s' failed with error:\n%s", tt.name, err)
		}
	}
}
