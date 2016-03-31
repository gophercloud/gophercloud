package gophercloud

type PasswordCredentialsV2 struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type TokenCredentialsV2 struct {
	ID string `json:"id,omitempty" required:"true"`
}

// AuthOptionsV2 wraps a gophercloud AuthOptions in order to adhere to the AuthOptionsBuilder
// interface.
type AuthOptionsV2 struct {
	PasswordCredentials *PasswordCredentialsV2 `json:"passwordCredentials,omitempty" xor:"TokenCredentials"`

	// The TenantID and TenantName fields are optional for the Identity V2 API.
	// Some providers allow you to specify a TenantName instead of the TenantId.
	// Some require both. Your provider's authentication policies will determine
	// how these fields influence authentication.
	TenantID   string `json:"tenantId,omitempty"`
	TenantName string `json:"tenantName,omitempty"`

	// TokenCredentials allows users to authenticate (possibly as another user) with an
	// authentication token ID.
	TokenCredentials *TokenCredentialsV2 `json:"token,omitempty" xor:"PasswordCredentials"`
}
