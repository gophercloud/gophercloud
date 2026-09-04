package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// These tests call auth.AuthOptionsFromCloud directly, using a real
// clouds.Cloud as the auth.CloudSource, to prove the restored top-level
// function works standalone and not merely through Cloud.AuthOptions's
// wrapper (covered exhaustively in openstack/config/clouds/auth_test.go).

func TestAuthOptionsFromCloudDirect(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"username": "testuser",
			"password": "testpass",
		},
	}

	opts, err := auth.AuthOptionsFromCloud(cloud)
	th.AssertNoErr(t, err)

	v3Opts, ok := opts.(*auth.AuthOptionsV3)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, "http://example.com:5000", v3Opts.AuthURL)
	th.AssertDeepEquals(t, auth.V3PasswordOpts{Username: "testuser", Password: "testpass", Scope: &auth.Scope{}}, v3Opts.Auth)
}

func TestAuthOptionsFromCloudV2Direct(t *testing.T) {
	cloud := clouds.Cloud{
		AuthType: auth.AuthV2Password,
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"username": "testuser",
			"password": "testpass",
		},
	}

	v2Opts, err := auth.AuthOptionsFromCloudV2(cloud)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, auth.V2PasswordOpts{Username: "testuser", Password: "testpass"}, v2Opts.Auth)
}

func TestAuthOptionsFromCloudV3Direct(t *testing.T) {
	cloud := clouds.Cloud{
		Auth: map[string]any{
			"auth_url": "http://example.com:5000",
			"password": "secret",
		},
	}

	v3Opts, err := auth.AuthOptionsFromCloudV3(cloud, auth.WithUsername("Kris"))
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "Kris", v3Opts.Auth.(auth.V3PasswordOpts).Username)
}
