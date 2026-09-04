package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthOptionsFromCloudV3PasswordInferred(t *testing.T) {
	cloud := clouds.Cloud{
		AuthInfo: &clouds.AuthInfo{
			AuthURL:  "http://example.com:5000",
			Username: "testuser",
			Password: "testpass",
		},
	}

	opts, err := auth.AuthOptionsFromCloud(cloud)
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromCloudV2ViaIdentityAPIVersion(t *testing.T) {
	cloud := clouds.Cloud{
		IdentityAPIVersion: "2.0",
		AuthInfo: &clouds.AuthInfo{
			AuthURL:  "http://example.com:5000",
			Username: "testuser",
			Password: "testpass",
		},
	}

	opts, err := auth.AuthOptionsFromCloud(cloud)
	th.AssertNoErr(t, err)

	v2Opts, ok := opts.(*auth.AuthOptionsV2)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v2Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass"}, v2Opts.Auth)
}

func TestAuthOptionsFromCloudMissingAuthURL(t *testing.T) {
	cloud := clouds.Cloud{
		AuthInfo: &clouds.AuthInfo{
			Username: "testuser",
			Password: "testpass",
		},
	}

	_, err := auth.AuthOptionsFromCloud(cloud)
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromCloudNilAuthInfo(t *testing.T) {
	cloud := clouds.Cloud{}

	_, err := auth.AuthOptionsFromCloud(cloud)
	th.AssertErr(t, err)

	_, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
}

func TestAuthOptionsFromCloudNoUsableCredentials(t *testing.T) {
	cloud := clouds.Cloud{
		AuthInfo: &clouds.AuthInfo{
			AuthURL: "http://example.com:5000",
		},
	}

	_, err := auth.AuthOptionsFromCloud(cloud)
	th.AssertErr(t, err)

	missingInput, ok := err.(gophercloud.ErrMissingInput)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "Auth", missingInput.Argument)
}
