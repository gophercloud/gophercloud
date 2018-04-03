package testing

import (
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

var emptyQuotaSet = quotasets.QuotaSet{}

func testSuccessTestCase(t *testing.T, jsonBody,
	uriPath, httpMethod string,
	expectedQuotaSet quotasets.QuotaSet,
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleSuccessfulRequest(t, httpMethod, uriPath, jsonBody)

	actual, err := quotasets.Get(client.ServiceClient(),
		FirstTenantID).Extract()
	if err != nil {
		return err
	}
	th.CheckDeepEquals(t, &expectedQuotaSet, actual)
}

func TestSuccessTestCases(t *testing.T) {
	for _, tt := range successTestCases {
		err := testSuccessTestCase(t, tt.jsonBody, tt.uriPath, tt.httpMethod,
			tt.expectedQuotaSet)
		if err != nil {
			t.Fatalf("Test case '%s' failed with error:\n%s", tt.name, err)
		}
	}
}
