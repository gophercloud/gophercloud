package v2

import (
	"os"
	"testing"
)

// RequireIdentityV2 will restrict a test to only be run in
// environments that support the Identity V2 API.
func RequireIdentityV2(t *testing.T) {
	if os.Getenv("OS_IDENTITY_API_VERSION") != "2.0" {
		t.Skip("this test requires support for the identity v2 API")
	}
}
