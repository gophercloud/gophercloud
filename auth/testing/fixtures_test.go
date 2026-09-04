package testing

import (
	"os"
	"testing"
)

// SetupEnvV2Password sets up environment variables for v2 password auth
func SetupEnvV2Password(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_IDENTITY_API_VERSION", "2.0")
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v2.0")
	t.Setenv("OS_USERNAME", "testuser")
	t.Setenv("OS_PASSWORD", "testpass")
}

// SetupEnvV2Token sets up environment variables for v2 token auth
func SetupEnvV2Token(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_IDENTITY_API_VERSION", "2.0")
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v2.0")
	t.Setenv("OS_TOKEN", "testtoken")
}

// SetupEnvV3Password sets up environment variables for v3 password auth
func SetupEnvV3Password(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
	t.Setenv("OS_USERNAME", "testuser")
	t.Setenv("OS_PASSWORD", "testpass")
}

// SetupEnvV3Token sets up environment variables for v3 token auth
func SetupEnvV3Token(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
	t.Setenv("OS_TOKEN", "testtoken")
}

// SetupEnvV3TOTP sets up environment variables for v3 TOTP auth
func SetupEnvV3TOTP(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
	t.Setenv("OS_PASSCODE", "123456")
}

// SetupEnvV3AppCred sets up environment variables for v3 application credential auth
func SetupEnvV3AppCred(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
	t.Setenv("OS_APPLICATION_CREDENTIAL_ID", "app-cred-id")
}

// SetupEnvV3AppCredName sets up environment variables for v3 application credential auth with name
func SetupEnvV3AppCredName(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
	t.Setenv("OS_APPLICATION_CREDENTIAL_NAME", "app-cred-name")
}

// SetupEnvV3ExplicitAuthType sets up environment variables for v3 with explicit auth type
func SetupEnvV3ExplicitAuthType(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
	t.Setenv("OS_AUTH_TYPE", "v3password")
	t.Setenv("OS_USERNAME", "testuser")
	t.Setenv("OS_PASSWORD", "testpass")
}

// SetupEnvV3WithAuthMethods sets up environment variables for v3 with auth methods
func SetupEnvV3WithAuthMethods(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
	t.Setenv("OS_AUTH_METHODS", "password,totp")
	t.Setenv("OS_USERNAME", "testuser")
	t.Setenv("OS_PASSWORD", "testpass")
}

// SetupEnvMissingAuthURL sets up environment variables missing auth URL
func SetupEnvMissingAuthURL(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_USERNAME", "testuser")
	t.Setenv("OS_PASSWORD", "testpass")
}

// SetupEnvUnsupportedAuthType sets up environment variables with unsupported auth type
func SetupEnvUnsupportedAuthType(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
	t.Setenv("OS_AUTH_TYPE", "unsupported")
}

// SetupEnvNoCredentials sets up environment variables with no credentials
func SetupEnvNoCredentials(t *testing.T) {
	CleanupEnv(t)
	t.Setenv("OS_AUTH_URL", "http://example.com:5000/v3")
}

// CleanupEnv cleans up all auth-related environment variables
func CleanupEnv(t *testing.T) {
	envVars := []string{
		"OS_IDENTITY_API_VERSION",
		"OS_AUTH_URL",
		"OS_AUTH_TYPE",
		"OS_AUTH_METHODS",
		"OS_USERNAME",
		"OS_USERID",
		"OS_PASSWORD",
		"OS_TOKEN",
		"OS_PASSCODE",
		"OS_TENANT_ID",
		"OS_TENANT_NAME",
		"OS_USER_DOMAIN_ID",
		"OS_USER_DOMAIN_NAME",
		"OS_DOMAIN_ID",
		"OS_DOMAIN_NAME",
		"OS_PROJECT_ID",
		"OS_PROJECT_NAME",
		"OS_PROJECT_DOMAIN_ID",
		"OS_PROJECT_DOMAIN_NAME",
		"OS_APPLICATION_CREDENTIAL_ID",
		"OS_APPLICATION_CREDENTIAL_NAME",
		"OS_APPLICATION_CREDENTIAL_SECRET",
	}
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}
