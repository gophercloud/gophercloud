package clients

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

// RequiredSystemScope will restrict a test to only be run by system scope.
func RequiredSystemScope(t *testing.T) {
	if os.Getenv("OS_SYSTEM_SCOPE") != "all" {
		t.Skip("must use system scope to run this test")
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

// RequireLong will ensure long-running tests can run.
func RequireLong(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
}

func getReleaseFromEnv(t *testing.T) string {
	current := strings.TrimPrefix(os.Getenv("OS_BRANCH"), "stable/")
	if current == "" {
		t.Fatal("this test requires OS_BRANCH to be set but it wasn't")
	}
	return current
}

// SkipRelease will have the test be skipped on a certain release.
// Releases are named such as 'stable/dalmatian', master, etc.
func SkipRelease(t *testing.T, release string) {
	current := getReleaseFromEnv(t)
	if current == strings.TrimPrefix(release, "stable/") {
		t.Skipf("this is not supported in %s", release)
	}
}

// SkipReleasesBelow will have the test be skipped on releases below a certain one.
// Releases are named such as 'stable/dalmatian', master, etc.
func SkipReleasesBelow(t *testing.T, release string) {
	current := getReleaseFromEnv(t)

	if IsCurrentBelow(t, release) {
		t.Skipf("this is not supported below %s, testing in %s", release, current)
	}
}

// SkipReleasesAbove will have the test be skipped on releases above a certain one.
// The test is always skipped on master release.
// Releases are named such as 'stable/dalmatian', master, etc.
func SkipReleasesAbove(t *testing.T, release string) {
	current := getReleaseFromEnv(t)

	if IsCurrentAbove(t, release) {
		t.Skipf("this is not supported above %s, testing in %s", release, current)
	}
}

func isReleaseNumeral(release string) bool {
	_, err := strconv.Atoi(release[0:1])
	return err == nil
}

// IsCurrentAbove will return true on releases above a certain one.
// The result is always true on master release.
// Releases are named such as 'stable/dalmatian', master, etc.
func IsCurrentAbove(t *testing.T, release string) bool {
	current := getReleaseFromEnv(t)
	release = strings.TrimPrefix(release, "stable/")

	if release != "master" {
		// Assume master is always too new
		if current == "master" {
			return true
		}
		// Numeral releases are always newer than non-numeral ones
		if isReleaseNumeral(current) && !isReleaseNumeral(release) {
			return true
		}
		if current > release && !(!isReleaseNumeral(current) && isReleaseNumeral(release)) {
			return true
		}
	}
	t.Logf("Target release %s is below the current branch %s", release, current)
	return false
}

// IsCurrentBelow will return true on releases below a certain one.
// The result is always false on master release.
// Releases are named such as 'stable/dalmatian', master, etc.
func IsCurrentBelow(t *testing.T, release string) bool {
	current := getReleaseFromEnv(t)
	release = strings.TrimPrefix(release, "stable/")

	if current != "master" {
		// Assume master is always too new
		if release == "master" {
			return true
		}
		// Numeral releases are always newer than non-numeral ones
		if isReleaseNumeral(release) && !isReleaseNumeral(current) {
			return true
		}
		if release > current && !(!isReleaseNumeral(release) && isReleaseNumeral(current)) {
			return true
		}
	}
	t.Logf("Target release %s is above the current branch %s", release, current)
	return false
}
