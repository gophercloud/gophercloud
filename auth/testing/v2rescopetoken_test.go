package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestV2RescopeTokenOptsToAuthBodyTenantID(t *testing.T) {
	opts := auth.V2RescopeTokenOpts{Token: "the-token", TenantID: "tenant-id"}

	authData, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)

	expected := map[string]map[string]any{
		"": {
			"tenantId":         "tenant-id",
			"tokenCredentials": map[string]any{"id": "the-token"},
		},
	}
	th.AssertDeepEquals(t, expected, authData)
}

func TestV2RescopeTokenOptsToAuthBodyTenantName(t *testing.T) {
	opts := auth.V2RescopeTokenOpts{Token: "the-token", TenantName: "tenant-name"}

	authData, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)

	expected := map[string]map[string]any{
		"": {
			"tenantName":       "tenant-name",
			"tokenCredentials": map[string]any{"id": "the-token"},
		},
	}
	th.AssertDeepEquals(t, expected, authData)
}

func TestV2RescopeTokenOptsMissingTenant(t *testing.T) {
	opts := auth.V2RescopeTokenOpts{Token: "the-token"}

	_, err := opts.ToAuthBody()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrScopeEmpty)
	th.AssertEquals(t, true, ok)
}

func TestV2RescopeTokenOptsMissingToken(t *testing.T) {
	opts := auth.V2RescopeTokenOpts{TenantID: "tenant-id"}

	_, err := opts.ToAuthBody()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestV2RescopeTokenOptsCanReauth(t *testing.T) {
	opts := auth.V2RescopeTokenOpts{Token: "t", TenantID: "tenant-id"}
	th.AssertEquals(t, false, opts.CanReauth())
}
