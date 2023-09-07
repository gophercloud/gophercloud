package clients

import (
	"os"
	"testing"
)

// RequiredSystemScope will restrict a test to only be run by system scope.
func RequiredSystemScope(t *testing.T) {
	if os.Getenv("OS_SYSTEM_SCOPE") != "all" {
		t.Skip("must use system scope to run this test")
	}
}

// RequireManilaReplicas will restrict a test to only be run with enabled
// manila replicas.
func RequireManilaReplicas(t *testing.T) {
	if os.Getenv("OS_MANILA_REPLICAS") != "true" {
		t.Skip("manila replicas must be enabled to run this test")
	}
}

// RequireAdmin will restrict a test to only be run by admin users.
func RequireAdmin(t *testing.T) {
	if os.Getenv("OS_USERNAME") != "admin" {
		t.Skip("must be admin to run this test")
	}
}

// RequireNonAdmin will restrict a test to only be run by non-admin users.
func RequireNonAdmin(t *testing.T) {
	if os.Getenv("OS_USERNAME") == "admin" {
		t.Skip("must be a non-admin to run this test")
	}
}

// RequirePortForwarding will restrict a test to only be run in environments
// that support port forwarding
func RequirePortForwarding(t *testing.T) {
	if os.Getenv("OS_PORTFORWARDING_ENVIRONMENT") == "" {
		t.Skip("this test requires support for port forwarding")
	}
}

// RequireGuestAgent will restrict a test to only be run in
// environments that support the QEMU guest agent.
func RequireGuestAgent(t *testing.T) {
	if os.Getenv("OS_GUEST_AGENT") == "" {
		t.Skip("this test requires support for qemu guest agent and to set OS_GUEST_AGENT to 1")
	}
}

// RequireIdentityV2 will restrict a test to only be run in
// environments that support the Identity V2 API.
func RequireIdentityV2(t *testing.T) {
	if os.Getenv("OS_IDENTITY_API_VERSION") != "2.0" {
		t.Skip("this test requires support for the identity v2 API")
	}
}

// RequireLiveMigration will restrict a test to only be run in
// environments that support live migration.
func RequireLiveMigration(t *testing.T) {
	if os.Getenv("OS_LIVE_MIGRATE") == "" {
		t.Skip("this test requires support for live migration and to set OS_LIVE_MIGRATE to 1")
	}
}

// RequireLong will ensure long-running tests can run.
func RequireLong(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
}

// RequireNovaNetwork will restrict a test to only be run in
// environments that support nova-network.
func RequireNovaNetwork(t *testing.T) {
	if os.Getenv("OS_NOVANET") == "" {
		t.Skip("this test requires nova-network and to set OS_NOVANET to 1")
	}
}

// RequireIronicHTTPBasic will restric a test to be only run in
// environments that have Ironic using http_basic.
func RequireIronicHTTPBasic(t *testing.T) {
	if os.Getenv("IRONIC_ENDPOINT") == "" || os.Getenv("OS_USERNAME") == "" || os.Getenv("OS_PASSWORD") == "" {
		t.Skip("this test requires Ironic using http_basic, set OS_USERNAME, OS_PASSWORD and IRONIC_ENDPOINT")
	}
}

func getReleaseFromEnv(t *testing.T) string {
	current_branch := os.Getenv("OS_BRANCH")
	if current_branch == "" {
		t.Fatal("this test requires OS_BRANCH to be set but it wasn't")
	}
	return current_branch
}

// SkipRelease will have the test be skipped on a certain
// release. Releases are named such as 'stable/mitaka', master, etc.
func SkipRelease(t *testing.T, release string) {
	current_branch := getReleaseFromEnv(t)
	if current_branch == release {
		t.Skipf("this is not supported in %s", release)
	}
}

// SkipReleasesBelow will have the test be skipped on releases below a certain
// one. Releases are named such as 'stable/mitaka', master, etc.
func SkipReleasesBelow(t *testing.T, release string) {
	current_branch := getReleaseFromEnv(t)

	if IsCurrentBelow(t, release) {
		t.Skipf("this is not supported below %s, testing in %s", release, current_branch)
	}
}

// SkipReleasesAbove will have the test be skipped on releases above a certain
// one. The test is always skipped on master release. Releases are named such
// as 'stable/mitaka', master, etc.
func SkipReleasesAbove(t *testing.T, release string) {
	current_branch := getReleaseFromEnv(t)

	// Assume master is always too new
	if IsCurrentAbove(t, release) {
		t.Skipf("this is not supported above %s, testing in %s", release, current_branch)
	}
}

// IsCurrentAbove will return true on releases above a certain
// one. The result is always true on master release. Releases are named such
// as 'stable/mitaka', master, etc.
func IsCurrentAbove(t *testing.T, release string) bool {
	current_branch := getReleaseFromEnv(t)

	// Assume master is always too new
	if current_branch == "master" || current_branch > release {
		return true
	}
	t.Logf("Target release %s is below the current branch %s", release, current_branch)
	return false
}

// IsCurrentBelow will return true on releases below a certain
// one. Releases are named such as 'stable/mitaka', master, etc.
func IsCurrentBelow(t *testing.T, release string) bool {
	current_branch := getReleaseFromEnv(t)

	if current_branch != "master" || current_branch < release {
		return true
	}
	t.Logf("Target release %s is above the current branch %s", release, current_branch)
	return false
}
