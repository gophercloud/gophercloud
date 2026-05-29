package httpbasic

import (
	"os"
	"testing"
)

// RequireIronicHTTPBasic will restrict a test to be only run in environments
// that have Ironic using http_basic.
func RequireIronicHTTPBasic(t *testing.T) {
	if os.Getenv("IRONIC_ENDPOINT") == "" || os.Getenv("OS_USERNAME") == "" || os.Getenv("OS_PASSWORD") == "" {
		t.Skip("this test requires Ironic using http_basic, set OS_USERNAME, OS_PASSWORD and IRONIC_ENDPOINT")
	}
}
