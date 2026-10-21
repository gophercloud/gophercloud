package noauth

import (
	"os"
	"testing"
)

// RequireCinderNoAuth will restrict a test to be only run in environments that
// have Cinder using noauth.
func RequireCinderNoAuth(t *testing.T) {
	if os.Getenv("CINDER_ENDPOINT") == "" || os.Getenv("OS_USER_ID") == "" || os.Getenv("OS_PROJECT_ID") == "" {
		t.Skip("this test requires Cinder v3 using noauth; set CINDER_ENDPOINT, OS_USER_ID, and OS_PROJECT_ID")
	}
}
