package clients

import (
	"os"
	"testing"
)

// RequireAdmin will restrict a test to only be run by admin users.
func RequireAdmin(t *testing.T) {
	if os.Getenv("OS_USERNAME") != "admin" {
		t.Skip("must be admin to run this test")
	}
}

// RequireNovaNetwork will restrict a test to only be run in
// environments that support nova-network.
func RequireNovaNetwork(t *testing.T) {
	if os.Getenv("OS_NOVANET") == "" {
		t.Skip("this test requires nova-network and to set OS_NOVANET to 1")
	}
}
