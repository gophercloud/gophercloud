package testing

import (
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
)

// Test for ErrEnvironmentVarsExpectedEqual to be returned where there
// is ambiguity between OS_USER_DOMAIN_X and OS_USER_PROJECT_X. Pairs
// are expected to be equal. OS_DOMAIN_X is added for backwards compatibility
func TestEnvVarsExpectedEqual(t *testing.T) {

	cases := []map[string]string{
		{
			"OS_TENANT_ID":  "tenant-id",
			"OS_PROJECT_ID": "project-id",
		},
		{
			"OS_TENANT_NAME":  "tenant-name",
			"OS_PROJECT_NAME": "project-name",
		},
		{
			"OS_USER_DOMAIN_ID":    "user-domain-id",
			"OS_PROJECT_DOMAIN_ID": "project-domain-id",
		},
		{
			"OS_USER_DOMAIN_ID": "user-domain-id",
			"OS_DOMAIN_ID":      "domain-id",
		},
		{
			"OS_DOMAIN_ID":         "domain-id",
			"OS_PROJECT_DOMAIN_ID": "project-domain-id",
		},
		{
			"OS_USER_DOMAIN_ID":    "user-domain-id",
			"OS_PROJECT_DOMAIN_ID": "domain-id",
			"OS_DOMAIN_ID":         "domain-id",
		},
		{
			"OS_USER_DOMAIN_ID":    "domain-id",
			"OS_PROJECT_DOMAIN_ID": "project-domain-id",
			"OS_DOMAIN_ID":         "domain-id",
		},
		{
			"OS_USER_DOMAIN_ID":    "user-domain-id",
			"OS_PROJECT_DOMAIN_ID": "project-domain-id",
			"OS_DOMAIN_ID":         "domain-id",
		},
		{
			"OS_USER_DOMAIN_NAME":    "user-domain",
			"OS_PROJECT_DOMAIN_NAME": "project-domain",
		},
		{
			"OS_USER_DOMAIN_NAME": "user-domain",
			"OS_DOMAIN_NAME":      "domain",
		},
		{
			"OS_DOMAIN_NAME":         "domain",
			"OS_PROJECT_DOMAIN_NAME": "project-domain",
		},
		{
			"OS_USER_DOMAIN_NAME":    "user-domain",
			"OS_PROJECT_DOMAIN_NAME": "project-domain",
			"OS_DOMAIN_NAME":         "domain",
		},
		{
			"OS_USER_DOMAIN_NAME":    "domain",
			"OS_PROJECT_DOMAIN_NAME": "project-domain",
			"OS_DOMAIN_NAME":         "domain",
		},
		{
			"OS_USER_DOMAIN_NAME":    "user-domain",
			"OS_PROJECT_DOMAIN_NAME": "domain",
			"OS_DOMAIN_NAME":         "domain",
		},
	}

	for _, c := range cases {

		for k, v := range c {
			os.Setenv(k, v)
			os.Setenv(k, v)
		}
		_, err := openstack.AuthOptionsFromEnv()

		// Didn't get at Error
		if err == nil {
			t.Errorf("Expected ErrEnvironmentVarsExpectedEqual, got %v", err)
		}

		// Didn't get the correct error
		if _, ok := err.(gophercloud.ErrEnvironmentVarsExpectedEqual); !ok {
			t.Errorf("Expected ErrEnvironmentVarsExpectedEqual, got %v", err)
		}

		// Clear ENV Vars for next iteration
		for k := range c {
			os.Unsetenv(k)
			os.Unsetenv(k)
		}
	}
}

// Test that setting OS_USER_DOMAIN_X with its OS_PROJECT_DOMAIN_X
// with the same value is ok and does not return an error. Includes OS_DOMAIN_X for backward
// compatibility
func TestUnambiguousEnvVars(t *testing.T) {

	cases := []map[string]string{
		{
			"OS_AUTH_URL":   "somedomain.com",
			"OS_USERNAME":   "someuser",
			"OS_PASSWORD":   "somepassword",
			"OS_TENANT_ID":  "project",
			"OS_PROJECT_ID": "project",
		},
		{
			"OS_AUTH_URL":          "somedomain.com",
			"OS_USERNAME":          "someuser",
			"OS_PASSWORD":          "somepassword",
			"OS_TENANT_NAME":       "project",
			"OS_PROJECT_NAME":      "project",
			"OS_PROJECT_DOMAIN_ID": "domain",
		},
		{
			"OS_AUTH_URL":          "somedomain.com",
			"OS_USERNAME":          "someuser",
			"OS_PASSWORD":          "somepassword",
			"OS_USER_DOMAIN_ID":    "domain",
			"OS_PROJECT_DOMAIN_ID": "domain",
		},
		{
			"OS_AUTH_URL":       "somedomain.com",
			"OS_USERNAME":       "someuser",
			"OS_PASSWORD":       "somepassword",
			"OS_USER_DOMAIN_ID": "domain",
			"OS_DOMAIN_ID":      "domain",
		},
		{
			"OS_AUTH_URL":          "somedomain.com",
			"OS_USERNAME":          "someuser",
			"OS_PASSWORD":          "somepassword",
			"OS_PROJECT_DOMAIN_ID": "domain",
			"OS_DOMAIN_ID":         "domain",
		},
		{
			"OS_AUTH_URL":          "somedomain.com",
			"OS_USERNAME":          "someuser",
			"OS_PASSWORD":          "somepassword",
			"OS_USER_DOMAIN_ID":    "domain",
			"OS_PROJECT_DOMAIN_ID": "domain",
			"OS_DOMAIN_ID":         "domain",
		},
		{
			"OS_AUTH_URL":            "somedomain.com",
			"OS_USERNAME":            "someuser",
			"OS_PASSWORD":            "somepassword",
			"OS_USER_DOMAIN_NAME":    "domain",
			"OS_PROJECT_DOMAIN_NAME": "domain",
		},
		{
			"OS_AUTH_URL":         "somedomain.com",
			"OS_USERNAME":         "someuser",
			"OS_PASSWORD":         "somepassword",
			"OS_USER_DOMAIN_NAME": "domain",
			"OS_DOMAIN_NAME":      "domain",
		},
		{
			"OS_AUTH_URL":            "somedomain.com",
			"OS_USERNAME":            "someuser",
			"OS_PASSWORD":            "somepassword",
			"OS_PROJECT_DOMAIN_NAME": "domain",
			"OS_DOMAIN_NAME":         "domain",
		},
		{
			"OS_AUTH_URL":            "somedomain.com",
			"OS_USERNAME":            "someuser",
			"OS_PASSWORD":            "somepassword",
			"OS_USER_DOMAIN_NAME":    "domain",
			"OS_PROJECT_DOMAIN_NAME": "domain",
			"OS_DOMAIN_NAME":         "domain",
		},
	}

	for _, c := range cases {

		for k, v := range c {
			os.Setenv(k, v)
			os.Setenv(k, v)
		}
		_, err := openstack.AuthOptionsFromEnv()

		// Got and error but shouldn't
		if err != nil {
			t.Errorf("Expected nil, got %v", err)
		}

		// Clear ENV Vars for next iteration
		for k := range c {
			os.Unsetenv(k)
			os.Unsetenv(k)
		}
	}
}

// Test various cases where there is ambiguity between passed environment
// varialbes and a ErrAmbiguousEnvironmentVarsClash is expected to be returned
func TestErrorAmbiguousEnvVarsClash(t *testing.T) {

	cases := []map[string]string{
		{
			"OS_USER_DOMAIN_ID":      "user-domain-id",
			"OS_PROJECT_DOMAIN_NAME": "project-domain",
		},
		{
			"OS_USER_DOMAIN_ID": "user-domain-id",
			"OS_DOMAIN_NAME":    "domain",
		},
		{
			"OS_USER_DOMAIN_ID":   "user-domain-id",
			"OS_USER_DOMAIN_NAME": "user-domain",
		},
		{
			"OS_PROJECT_DOMAIN_ID":   "project-domain-id",
			"OS_PROJECT_DOMAIN_NAME": "project-domain",
		},
		{
			"OS_PROJECT_DOMAIN_ID": "project-domain-id",
			"OS_DOMAIN_NAME":       "domain",
		},
		{
			"OS_PROJECT_DOMAIN_ID": "project-domain-id",
			"OS_USER_DOMAIN_NAME":  "user-domain",
		},
	}

	for _, c := range cases {

		for k, v := range c {
			os.Setenv(k, v)
			os.Setenv(k, v)
		}

		_, err := openstack.AuthOptionsFromEnv()

		// Didn't get at Error
		if err == nil {
			t.Errorf("Expected ErrAmbiguousEnvironmentVarsClash, got %v", err)
		}

		// Didn't get the correct error
		if _, ok := err.(gophercloud.ErrAmbiguousEnvironmentVarsClash); !ok {
			t.Errorf("Expected ErrAmbiguousEnvironmentVarsClash, got %v", err)
		}

		// Clear ENV Vars for next iteration
		for k := range c {
			os.Unsetenv(k)
			os.Unsetenv(k)
		}
	}
}
