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
