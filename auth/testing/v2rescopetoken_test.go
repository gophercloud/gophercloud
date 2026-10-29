package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestV2TokenOptsRescopeTenantID(t *testing.T) {
	opts := auth.V2TokenOpts{Token: "the-token", TenantID: "tenant-id"}
	body, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]map[string]any{
		"token": {"tenantId": "tenant-id", "token": map[string]any{"id": "the-token"}},
	}, body)
}

func TestV2TokenOptsRescopeTenantName(t *testing.T) {
	opts := auth.V2TokenOpts{Token: "the-token", TenantName: "tenant-name"}
	body, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]map[string]any{
		"token": {"tenantName": "tenant-name", "token": map[string]any{"id": "the-token"}},
	}, body)
}

func TestV2TokenOptsUnscoped(t *testing.T) {
	body, err := (auth.V2TokenOpts{Token: "the-token"}).ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]map[string]any{
		"token": {"token": map[string]any{"id": "the-token"}},
	}, body)
}

func TestV2TokenOptsRescopeMissingToken(t *testing.T) {
	_, err := (auth.V2TokenOpts{TenantID: "tenant-id"}).ToAuthBody()
	th.AssertErr(t, err)
	missing, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "Token", missing.Argument)
}

func TestV2TokenOptsRescopeCannotReauthByDefault(t *testing.T) {
	opts := auth.V2TokenOpts{Token: "t", TenantID: "tenant-id"}
	th.AssertEquals(t, false, opts.CanReauth())
}
