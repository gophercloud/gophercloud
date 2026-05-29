package noauth

import (
	"os"
	"testing"
)

// RequireIronicNoAuth will restrict a test to be only run in environments that

// have Ironic using noauth.
func RequireIronicNoAuth(t *testing.T) {
	if os.Getenv("IRONIC_ENDPOINT") == "" || os.Getenv("OS_USERNAME") == "" {
		t.Skip("this test requires IRONIC using noauth, set OS_USERNAME and IRONIC_ENDPOINT")
	}
}
