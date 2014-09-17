// Package utils contains utilities which eases working with Gophercloud's OpenStack APIs.
package utils

import (
	"fmt"
	"os"

	"github.com/rackspace/gophercloud"
)

var nilOptions = gophercloud.AuthOptions{}

// ErrNoAuthUrl, ErrNoUsername, and ErrNoPassword errors indicate of the required OS_AUTH_URL, OS_USERNAME, or OS_PASSWORD
// environment variables, respectively, remain undefined.  See the AuthOptions() function for more details.
var (
	ErrNoAuthURL  = fmt.Errorf("Environment variable OS_AUTH_URL needs to be set.")
	ErrNoUsername = fmt.Errorf("Environment variable OS_USERNAME needs to be set.")
	ErrNoPassword = fmt.Errorf("Environment variable OS_PASSWORD needs to be set.")
)

// AuthOptions fills out an identity.AuthOptions structure with the settings found on the various OpenStack
// OS_* environment variables.  The following variables provide sources of truth: OS_AUTH_URL, OS_USERNAME,
// OS_PASSWORD, OS_TENANT_ID, and OS_TENANT_NAME.  Of these, OS_USERNAME, OS_PASSWORD, and OS_AUTH_URL must
// have settings, or an error will result.  OS_TENANT_ID and OS_TENANT_NAME are optional.
func AuthOptions() (gophercloud.AuthOptions, error) {
	authURL := os.Getenv("OS_AUTH_URL")
	username := os.Getenv("OS_USERNAME")
	userID := os.Getenv("OS_USERID")
	password := os.Getenv("OS_PASSWORD")
	tenantID := os.Getenv("OS_TENANT_ID")
	tenantName := os.Getenv("OS_TENANT_NAME")
	domainID := os.Getenv("OS_DOMAIN_ID")
	domainName := os.Getenv("OS_DOMAIN_NAME")

	if authURL == "" {
		return nilOptions, ErrNoAuthURL
	}

	if username == "" && userID == "" {
		return nilOptions, ErrNoUsername
	}

	if password == "" {
		return nilOptions, ErrNoPassword
	}

	ao := gophercloud.AuthOptions{
		IdentityEndpoint: authURL,
		UserID:           userID,
		Username:         username,
		Password:         password,
		TenantID:         tenantID,
		TenantName:       tenantName,
		DomainID:         domainID,
		DomainName:       domainName,
	}

	return ao, nil
}

// BuildQuery constructs the query section of a URI from a map.
func BuildQuery(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	query := "?"
	for k, v := range params {
		query += k + "=" + v + "&"
	}
	query = query[:len(query)-1]
	return query
}

func MaybeString(original string) *string {
	if original != "" {
		return &original
	}
	return nil
}

func MaybeInt(original int) *int {
	if original != 0 {
		return &original
	}
	return nil
}
