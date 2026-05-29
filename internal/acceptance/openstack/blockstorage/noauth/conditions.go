package noauth

import (
	"os"
	"testing"
)

// RequireCinderNoAuth will restrict a test to be only run in environments that
// have Cinder using noauth.
func RequireCinderNoAuth(t *testing.T) {
	if os.Getenv("CINDER_ENDPOINT") == "" || os.Getenv("OS_USERNAME") == "" {
		t.Skip("this test requires Cinder using noauth, set OS_USERNAME and CINDER_ENDPOINT")
	}
}
