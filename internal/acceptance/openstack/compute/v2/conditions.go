package v2

import (
	"os"
	"testing"
)

// RequireGuestAgent will restrict a test to only be run in
// environments that support the QEMU guest agent.
func RequireGuestAgent(t *testing.T) {
	if os.Getenv("OS_GUEST_AGENT") == "" {
		t.Skip("this test requires support for qemu guest agent and to set OS_GUEST_AGENT to 1")
	}
}

// RequireLiveMigration will restrict a test to only be run in
// environments that support live migration.
func RequireLiveMigration(t *testing.T) {
	if os.Getenv("OS_LIVE_MIGRATE") == "" {
		t.Skip("this test requires support for live migration and to set OS_LIVE_MIGRATE to 1")
	}
}
