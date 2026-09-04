package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthOptionsFromEnvV2Password(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Password(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v2Opts, ok := opts.(*auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v2.0", v2Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass"}, v2Opts.Auth)
}

func TestAuthOptionsFromEnvV2Token(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Token(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v2Opts, ok := opts.(*auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v2.0", v2Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2TokenOpts{Token: "testtoken"}, v2Opts.Auth)
}

func TestAuthOptionsFromEnvV3Password(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Password(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3Token(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Token(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3TokenOpts{Token: "testtoken", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3TOTP(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3TOTP(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3TOTPOpts{Passcode: "123456", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3AppCred(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3AppCred(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "app-cred-id"}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3AppCredName(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3AppCredName(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "app-cred-name"}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3ExplicitAuthType(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3ExplicitAuthType(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvV3WithAuthMethods(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3WithAuthMethods(t)

	opts, err := auth.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000/v3", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromEnvMissingAuthURL(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvMissingAuthURL(t)

	_, err := auth.AuthOptionsFromEnv()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingEnvironmentVariable)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromEnvUnsupportedAuthType(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvUnsupportedAuthType(t)

	_, err := auth.AuthOptionsFromEnv()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrUnsupportedAuthType)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromEnvNoCredentials(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvNoCredentials(t)

	_, err := auth.AuthOptionsFromEnv()
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrUnsupportedAuthType)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromEnvV2Direct(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV2Password(t)

	opts, err := auth.AuthOptionsFromEnvV2()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "http://example.com:5000/v2.0", opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass"}, opts.Auth)
}

func TestAuthOptionsFromEnvV3Direct(t *testing.T) {
	defer CleanupEnv(t)
	SetupEnvV3Password(t)

	opts, err := auth.AuthOptionsFromEnvV3()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "http://example.com:5000/v3", opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, opts.Auth)
}
