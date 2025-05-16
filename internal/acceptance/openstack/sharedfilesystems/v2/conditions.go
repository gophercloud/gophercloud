package v2

import (
	"os"
	"testing"
)

// RequireManilaReplicas will restrict a test to only be run with enabled
// manila replicas.
func RequireManilaReplicas(t *testing.T) {
	if os.Getenv("OS_MANILA_REPLICAS") != "true" {
		t.Skip("manila replicas must be enabled to run this test")
	}
}
