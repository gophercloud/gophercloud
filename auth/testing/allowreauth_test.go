package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCanReauthRespectsAllowReauth(t *testing.T) {
	th.AssertEquals(t, false, auth.V2PasswordOpts{}.CanReauth())
	th.AssertEquals(t, true, auth.V2PasswordOpts{AllowReauth: true}.CanReauth())

	th.AssertEquals(t, false, auth.V2TokenOpts{}.CanReauth())
	th.AssertEquals(t, true, auth.V2TokenOpts{AllowReauth: true}.CanReauth())

	th.AssertEquals(t, false, auth.V3TokenOpts{}.CanReauth())
	th.AssertEquals(t, true, auth.V3TokenOpts{AllowReauth: true}.CanReauth())

	th.AssertEquals(t, false, auth.V3PasswordOpts{}.CanReauth())
	th.AssertEquals(t, true, auth.V3PasswordOpts{AllowReauth: true}.CanReauth())

	th.AssertEquals(t, false, auth.V3ApplicationCredentialOpts{}.CanReauth())
	th.AssertEquals(t, true, auth.V3ApplicationCredentialOpts{AllowReauth: true}.CanReauth())
}

// TestCanReauthHardcodedFalseMechanismsIgnoreAllowReauth guards against a
// future edit accidentally enabling automatic reauthentication for these
// mechanisms.
func TestCanReauthHardcodedFalseMechanismsIgnoreAllowReauth(t *testing.T) {
	th.AssertEquals(t, false, auth.V3TOTPOpts{}.CanReauth())
	th.AssertEquals(t, false, auth.V3MultifactorOpts{}.CanReauth())
}
