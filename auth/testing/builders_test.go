package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestV2PasswordOptsToAuthBody(t *testing.T) {
	opts := auth.V2PasswordOpts{
		Username:   "testuser",
		Password:   "testpass",
		TenantID:   "tenant-id",
		TenantName: "tenant-name",
	}

	body, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]map[string]any{
		"password": {
			"tenantId":   "tenant-id",
			"tenantName": "tenant-name",
			"passwordCredentials": map[string]any{
				"username": "testuser",
				"password": "testpass",
			},
		},
	}, body)
}

func TestV2PasswordOptsRequiresCredentials(t *testing.T) {
	tests := []struct {
		name string
		opts auth.V2PasswordOpts
	}{
		{name: "username", opts: auth.V2PasswordOpts{Password: "testpass"}},
		{name: "password", opts: auth.V2PasswordOpts{Username: "testuser"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.opts.ToAuthBody()
			th.AssertErr(t, err)
			_, ok := err.(gophercloud.ErrMissingInput)
			th.AssertEquals(t, true, ok)
		})
	}
}

func TestV2TokenOptsToAuthBody(t *testing.T) {
	opts := auth.V2TokenOpts{Token: "testtoken", TenantName: "tenant-name"}

	body, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]map[string]any{
		"token": {
			"tenantName":       "tenant-name",
			"tokenCredentials": map[string]any{"id": "testtoken"},
		},
	}, body)
}

func TestV2TokenOptsRequiresToken(t *testing.T) {
	_, err := (auth.V2TokenOpts{}).ToAuthBody()
	th.AssertErr(t, err)

	missing, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "Token", missing.Argument)
}

func TestV3PasswordOptsToAuthBody(t *testing.T) {
	tests := []struct {
		name string
		opts auth.V3PasswordOpts
		want map[string]map[string]any
	}{
		{
			name: "username and domain ID",
			opts: auth.V3PasswordOpts{Username: "testuser", Password: "testpass", UserDomainID: "domain-id"},
			want: map[string]map[string]any{"password": {"user": map[string]any{
				"name": "testuser", "password": "testpass", "domain": map[string]any{"id": "domain-id"},
			}}},
		},
		{
			name: "username and domain name",
			opts: auth.V3PasswordOpts{Username: "testuser", Password: "testpass", UserDomainName: "domain-name"},
			want: map[string]map[string]any{"password": {"user": map[string]any{
				"name": "testuser", "password": "testpass", "domain": map[string]any{"name": "domain-name"},
			}}},
		},
		{
			name: "user ID",
			opts: auth.V3PasswordOpts{UserID: "user-id", Password: "testpass"},
			want: map[string]map[string]any{"password": {"user": map[string]any{
				"id": "user-id", "password": "testpass",
			}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := test.opts.ToAuthBody()
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, test.want, body)
		})
	}
}

func TestV3PasswordOptsValidation(t *testing.T) {
	tests := []struct {
		name string
		opts auth.V3PasswordOpts
		kind any
	}{
		{name: "missing password", opts: auth.V3PasswordOpts{UserID: "user-id"}, kind: gophercloud.ErrMissingPassword{}},
		{name: "missing user", opts: auth.V3PasswordOpts{Password: "testpass"}, kind: gophercloud.ErrUsernameOrUserID{}},
		{name: "both username and user ID", opts: auth.V3PasswordOpts{Username: "testuser", UserID: "user-id", Password: "testpass"}, kind: gophercloud.ErrUsernameOrUserID{}},
		{name: "username without domain", opts: auth.V3PasswordOpts{Username: "testuser", Password: "testpass"}, kind: gophercloud.ErrDomainIDOrDomainName{}},
		{name: "username with both domains", opts: auth.V3PasswordOpts{Username: "testuser", Password: "testpass", UserDomainID: "domain-id", UserDomainName: "domain-name"}, kind: gophercloud.ErrDomainIDOrDomainName{}},
		{name: "user ID with domain ID", opts: auth.V3PasswordOpts{UserID: "user-id", Password: "testpass", UserDomainID: "domain-id"}, kind: gophercloud.ErrDomainIDWithUserID{}},
		{name: "user ID with domain name", opts: auth.V3PasswordOpts{UserID: "user-id", Password: "testpass", UserDomainName: "domain-name"}, kind: gophercloud.ErrDomainNameWithUserID{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.opts.ToAuthBody()
			th.AssertErr(t, err)
			th.AssertEquals(t, test.kind, err)
		})
	}
}

func TestV3TOTPOptsToAuthBody(t *testing.T) {
	tests := []struct {
		name string
		opts auth.V3TOTPOpts
		want map[string]map[string]any
	}{
		{
			name: "username and domain ID",
			opts: auth.V3TOTPOpts{Username: "testuser", Passcode: "123456", UserDomainID: "domain-id"},
			want: map[string]map[string]any{"totp": {"user": map[string]any{
				"name": "testuser", "passcode": "123456", "domain": map[string]any{"id": "domain-id"},
			}}},
		},
		{
			name: "username and domain name",
			opts: auth.V3TOTPOpts{Username: "testuser", Passcode: "123456", UserDomainName: "domain-name"},
			want: map[string]map[string]any{"totp": {"user": map[string]any{
				"name": "testuser", "passcode": "123456", "domain": map[string]any{"name": "domain-name"},
			}}},
		},
		{
			name: "user ID",
			opts: auth.V3TOTPOpts{UserID: "user-id", Passcode: "123456"},
			want: map[string]map[string]any{"totp": {"user": map[string]any{
				"id": "user-id", "passcode": "123456",
			}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := test.opts.ToAuthBody()
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, test.want, body)
		})
	}
}

func TestV3TOTPOptsValidation(t *testing.T) {
	tests := []struct {
		name string
		opts auth.V3TOTPOpts
		kind any
	}{
		{name: "missing passcode", opts: auth.V3TOTPOpts{UserID: "user-id"}, kind: gophercloud.ErrMissingPasscode{}},
		{name: "missing user", opts: auth.V3TOTPOpts{Passcode: "123456"}, kind: gophercloud.ErrUsernameOrUserID{}},
		{name: "both username and user ID", opts: auth.V3TOTPOpts{Username: "testuser", UserID: "user-id", Passcode: "123456"}, kind: gophercloud.ErrUsernameOrUserID{}},
		{name: "username without domain", opts: auth.V3TOTPOpts{Username: "testuser", Passcode: "123456"}, kind: gophercloud.ErrDomainIDOrDomainName{}},
		{name: "username with both domains", opts: auth.V3TOTPOpts{Username: "testuser", Passcode: "123456", UserDomainID: "domain-id", UserDomainName: "domain-name"}, kind: gophercloud.ErrDomainIDOrDomainName{}},
		{name: "user ID with domain ID", opts: auth.V3TOTPOpts{UserID: "user-id", Passcode: "123456", UserDomainID: "domain-id"}, kind: gophercloud.ErrDomainIDWithUserID{}},
		{name: "user ID with domain name", opts: auth.V3TOTPOpts{UserID: "user-id", Passcode: "123456", UserDomainName: "domain-name"}, kind: gophercloud.ErrDomainNameWithUserID{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.opts.ToAuthBody()
			th.AssertErr(t, err)
			th.AssertEquals(t, test.kind, err)
		})
	}
}

func TestV3ApplicationCredentialOptsToAuthBody(t *testing.T) {
	tests := []struct {
		name string
		opts auth.V3ApplicationCredentialOpts
		want map[string]map[string]any
	}{
		{
			name: "credential ID",
			opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "credential-id", ApplicationCredentialSecret: "secret"},
			want: map[string]map[string]any{"application_credential": {
				"id": "credential-id", "secret": "secret", "user": map[string]any{},
			}},
		},
		{
			name: "credential name and user ID",
			opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret", UserID: "user-id"},
			want: map[string]map[string]any{"application_credential": {
				"name": "credential-name", "secret": "secret", "user": map[string]any{"id": "user-id"},
			}},
		},
		{
			name: "credential name, username, and domain name",
			opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret", Username: "testuser", UserDomainName: "domain-name"},
			want: map[string]map[string]any{"application_credential": {
				"name": "credential-name", "secret": "secret", "user": map[string]any{
					"name": "testuser", "domain": map[string]any{"name": "domain-name"},
				},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := test.opts.ToAuthBody()
			th.AssertNoErr(t, err)
			th.AssertDeepEquals(t, test.want, body)
		})
	}
}

func TestV3ApplicationCredentialOptsValidation(t *testing.T) {
	tests := []struct {
		name string
		opts auth.V3ApplicationCredentialOpts
		kind any
	}{
		{name: "missing secret", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "credential-id"}, kind: gophercloud.ErrAppCredMissingSecret{}},
		{name: "missing credential identifier", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialSecret: "secret"}, kind: gophercloud.ErrAppCredNameOrAppCredID{}},
		{name: "both credential identifiers", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "credential-id", ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret"}, kind: gophercloud.ErrAppCredNameOrAppCredID{}},
		{name: "name without user", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret"}, kind: gophercloud.ErrUsernameOrUserID{}},
		{name: "both username and user ID", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret", Username: "testuser", UserID: "user-id"}, kind: gophercloud.ErrUsernameOrUserID{}},
		{name: "username without domain", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret", Username: "testuser"}, kind: gophercloud.ErrDomainIDOrDomainName{}},
		{name: "username with both domains", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret", Username: "testuser", UserDomainID: "domain-id", UserDomainName: "domain-name"}, kind: gophercloud.ErrDomainIDOrDomainName{}},
		{name: "user ID with domain ID", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret", UserID: "user-id", UserDomainID: "domain-id"}, kind: gophercloud.ErrDomainIDWithUserID{}},
		{name: "user ID with domain name", opts: auth.V3ApplicationCredentialOpts{ApplicationCredentialName: "credential-name", ApplicationCredentialSecret: "secret", UserID: "user-id", UserDomainName: "domain-name"}, kind: gophercloud.ErrDomainNameWithUserID{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.opts.ToAuthBody()
			th.AssertErr(t, err)
			th.AssertEquals(t, test.kind, err)
		})
	}
}

func TestV3TokenOptsToAuthBody(t *testing.T) {
	body, err := (auth.V3TokenOpts{Token: "testtoken"}).ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]map[string]any{"token": {"id": "testtoken"}}, body)
}

func TestV3TokenOptsRequiresToken(t *testing.T) {
	_, err := (auth.V3TokenOpts{}).ToAuthBody()
	th.AssertErr(t, err)

	missing, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "Token", missing.Argument)
}

func TestV3AuthTypes(t *testing.T) {
	tests := []struct {
		name string
		opts auth.AuthOptionsBuilderV3
		want auth.AuthType
	}{
		{name: "password", opts: auth.V3PasswordOpts{}, want: auth.AuthV3Password},
		{name: "TOTP", opts: auth.V3TOTPOpts{}, want: auth.AuthV3Totp},
		{name: "token", opts: auth.V3TokenOpts{}, want: auth.AuthV3Token},
		{name: "application credential", opts: auth.V3ApplicationCredentialOpts{}, want: auth.AuthV3ApplicationCredential},
		{name: "rescope token", opts: auth.V3RescopeTokenOpts{}, want: auth.AuthV3Token},
		{name: "multifactor", opts: auth.V3MultifactorOpts{}, want: auth.AuthV3MultiFactor},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			th.AssertEquals(t, test.want, test.opts.ToAuthType())
		})
	}
}

func TestV3MultifactorOptsRequiresAuthMethods(t *testing.T) {
	_, err := (auth.V3MultifactorOpts{}).ToAuthBody()
	th.AssertErr(t, err)

	missing, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "AuthMethods", missing.Argument)
}

func TestV3MultifactorOptsToAuthBody(t *testing.T) {
	opts := auth.V3MultifactorOpts{AuthMethods: []auth.AuthOptionsBuilderV3{
		auth.V3PasswordOpts{UserID: "user-id", Password: "testpass"},
		auth.V3TOTPOpts{UserID: "user-id", Passcode: "123456"},
		auth.V3TokenOpts{Token: "testtoken"},
		auth.V3ApplicationCredentialOpts{ApplicationCredentialID: "credential-id", ApplicationCredentialSecret: "secret"},
	}}

	body, err := opts.ToAuthBody()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, map[string]map[string]any{
		"password":               {"user": map[string]any{"id": "user-id", "password": "testpass"}},
		"totp":                   {"user": map[string]any{"id": "user-id", "passcode": "123456"}},
		"token":                  {"id": "testtoken"},
		"application_credential": {"id": "credential-id", "secret": "secret", "user": map[string]any{}},
	}, body)
}

func TestV3MultifactorOptsPropagatesMethodError(t *testing.T) {
	opts := auth.V3MultifactorOpts{AuthMethods: []auth.AuthOptionsBuilderV3{
		auth.V3PasswordOpts{UserID: "user-id"},
	}}

	_, err := opts.ToAuthBody()
	th.AssertErr(t, err)
	_, ok := err.(gophercloud.ErrMissingPassword)
	th.AssertEquals(t, true, ok)
}

func TestV3MultifactorOptsRejectsUnsupportedMethod(t *testing.T) {
	opts := auth.V3MultifactorOpts{AuthMethods: []auth.AuthOptionsBuilderV3{
		auth.V3RescopeTokenOpts{Token: "testtoken", Scope: &auth.Scope{ProjectID: "project-id"}},
	}}

	_, err := opts.ToAuthBody()
	th.AssertErr(t, err)
	_, ok := err.(gophercloud.ErrUnsupportedAuthType)
	th.AssertEquals(t, true, ok)
}
