// This package contains utilities which eases working with Gophercloud's OpenStack APIs.
package utils

import (
	"fmt"
	"github.com/rackspace/gophercloud/openstack/identity"
	"os"
)

var nilOptions = identity.AuthOptions{}

// ErrNoAuthUrl, ErrNoUsername, and ErrNoPassword errors indicate of the required OS_AUTH_URL, OS_USERNAME, or OS_PASSWORD
// environment variables, respectively, remain undefined.  See the AuthOptions() function for more details.
var (
	ErrNoAuthUrl  = fmt.Errorf("Environment variable OS_AUTH_URL needs to be set.")
	ErrNoUsername = fmt.Errorf("Environment variable OS_USERNAME needs to be set.")
	ErrNoPassword = fmt.Errorf("Environment variable OS_PASSWORD needs to be set.")
)

// AuthOptions fills out an identity.AuthOptions structure with the settings found on the various OpenStack
// OS_* environment variables.  The following variables provide sources of truth: OS_AUTH_URL, OS_USERNAME,
// OS_PASSWORD, OS_TENANT_ID, and OS_TENANT_NAME.  Of these, OS_USERNAME, OS_PASSWORD, and OS_AUTH_URL must
// have settings, or an error will result.  OS_TENANT_ID and OS_TENANT_NAME are optional.
func AuthOptions() (identity.AuthOptions, error) {
	authUrl := os.Getenv("OS_AUTH_URL")
	username := os.Getenv("OS_USERNAME")
	password := os.Getenv("OS_PASSWORD")
	tenantId := os.Getenv("OS_TENANT_ID")
	tenantName := os.Getenv("OS_TENANT_NAME")

	if authUrl == "" {
		return nilOptions, ErrNoAuthUrl
	}

	if username == "" {
		return nilOptions, ErrNoUsername
	}

	if password == "" {
		return nilOptions, ErrNoPassword
	}

	ao := identity.AuthOptions{
		Endpoint:   authUrl,
		Username:   username,
		Password:   password,
		TenantId:   tenantId,
		TenantName: tenantName,
	}

	return ao, nil
}

func BuildQuery(params map[string]string) string {
	query := "?"
	for k, v := range params {
		query += k + "=" + v + "&"
	}
	query = query[:len(query)-1]
	return query
}
