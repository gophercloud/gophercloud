package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

var emptyQuotaSet = quotasets.QuotaSet{}
var emptyQuotaDetailSet = quotasets.QuotaDetailSet{}

func testSuccessTestCase(t *testing.T,
	httpMethod, uriPath, jsonBody string,
	uriQueryParams map[string]string,
	expectedQuotaSet quotasets.QuotaSet,
	expectedQuotaDetailSet quotasets.QuotaDetailSet,
) error {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSuccessfulRequest(t, httpMethod, uriPath, jsonBody, uriQueryParams)

	if expectedQuotaSet != emptyQuotaSet {
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
		err := testSuccessTestCase(t,
			tt.httpMethod, tt.uriPath, tt.jsonBody, tt.uriQueryParams,
			tt.expectedQuotaSet,
			tt.expectedQuotaDetailSet,
		)
		if err != nil {
			t.Fatalf("Test case '%s' failed with error:\n%s", tt.name, err)
		}
	}
}
