package clouds

import "github.com/gophercloud/gophercloud/v2/auth"

// GetAuthType, GetIdentityAPIVersion, and GetAuth make Cloud satisfy
// auth.CloudSource, so it can be passed directly to
// auth.AuthOptionsFromCloud/V2/V3.

// GetAuthType returns c.AuthType.
func (c Cloud) GetAuthType() auth.AuthType { return c.AuthType }

// GetIdentityAPIVersion returns c.IdentityAPIVersion.
func (c Cloud) GetIdentityAPIVersion() string { return c.IdentityAPIVersion }

// GetAuth returns c.Auth.
func (c Cloud) GetAuth() map[string]any { return c.Auth }

// AuthOptions builds an auth.Authenticator from a Cloud (as returned by
// Parse). See auth.AuthOptionsFromCloud for the mechanism-selection rules.
func (c Cloud) AuthOptions(opts ...auth.CloudOption) (auth.Authenticator, error) {
	return auth.AuthOptionsFromCloud(c, opts...)
}

// AuthOptionsV2 builds an *auth.AuthOptionsV2 from a Cloud, ignoring
// c.IdentityAPIVersion. See auth.AuthOptionsFromCloudV2.
func (c Cloud) AuthOptionsV2(opts ...auth.CloudOption) (*auth.AuthOptionsV2, error) {
	return auth.AuthOptionsFromCloudV2(c, opts...)
}

// AuthOptionsV3 builds an *auth.AuthOptionsV3 from a Cloud, ignoring
// c.IdentityAPIVersion. See auth.AuthOptionsFromCloudV3.
func (c Cloud) AuthOptionsV3(opts ...auth.CloudOption) (*auth.AuthOptionsV3, error) {
	return auth.AuthOptionsFromCloudV3(c, opts...)
}
