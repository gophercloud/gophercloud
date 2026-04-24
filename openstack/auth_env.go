package openstack

import (
	"os"

	"github.com/gophercloud/gophercloud/v2"
)

var nilOptions = gophercloud.AuthOptions{}

/*
AuthOptionsFromEnv fills out an identity.AuthOptions structure with the
settings found on the various OpenStack OS_* environment variables.

The following variables provide sources of truth: OS_AUTH_URL, OS_USERNAME,
OS_PASSWORD and OS_PROJECT_ID.

Of these, OS_USERNAME, OS_PASSWORD, and OS_AUTH_URL must have settings,
or an error will result.  OS_PROJECT_ID, is optional.

OS_TENANT_ID and OS_TENANT_NAME are deprecated forms of OS_PROJECT_ID and
OS_PROJECT_NAME and the latter are expected against a v3 auth api.

If OS_PROJECT_ID and OS_PROJECT_NAME are set, they will still be referred
as "tenant" in Gophercloud.

If OS_PROJECT_NAME is set, it requires OS_PROJECT_DOMAIN_ID or OS_PROJECT_DOMAIN_NAME
to be set as well to handle projects not on the default domain. Previous versions of
gophercloud allowed OS_DOMAIN_ID and OS_DOMAIN_NAME to be used, so they are
still accepted for backwards compatibility.

Previous versions of gophercloud did not require OS_USER_DOMAIN_ID or OS_USER_DOMAIN_NAME
when using OS_USERNAME - instead it relied solely on OS_DOMAIN_NAME and OS_DOMAIN_ID. Now,
all the above are accepted for backwards compatibility.

To use this function, first set the OS_* environment variables (for example,
by sourcing an `openrc` file), then:

	opts, err := openstack.AuthOptionsFromEnv()
	provider, err := openstack.AuthenticatedClient(context.TODO(), opts)
*/
func AuthOptionsFromEnv() (gophercloud.AuthOptions, error) {
	authURL := os.Getenv("OS_AUTH_URL")
	username := os.Getenv("OS_USERNAME")
	userID := os.Getenv("OS_USERID")
	password := os.Getenv("OS_PASSWORD")
	passcode := os.Getenv("OS_PASSCODE")
	tenantID := os.Getenv("OS_TENANT_ID")
	tenantName := os.Getenv("OS_TENANT_NAME")
	domainID := os.Getenv("OS_DOMAIN_ID")
	domainName := os.Getenv("OS_DOMAIN_NAME")
	applicationCredentialID := os.Getenv("OS_APPLICATION_CREDENTIAL_ID")
	applicationCredentialName := os.Getenv("OS_APPLICATION_CREDENTIAL_NAME")
	applicationCredentialSecret := os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET")
	systemScope := os.Getenv("OS_SYSTEM_SCOPE")

	// Disambiguate between OS_TENANT_ID and OS_PROJECT_ID
	if v1, v2 := os.Getenv("OS_TENANT_ID"), os.Getenv("OS_PROJECT_ID"); v1 != v2 {
		// Allow setting OS_TENANT_ID and OS_PROJECT_ID if they are equal
		// Do nothing if one is set and the other isn't
		if v1 != "" && v2 != "" {
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_TENANT_ID", "OS_PROJECT_ID"},
			}
			return nilOptions, err
		}

	}

	// Disambiguate between OS_TENANT_NAME and OS_PROJECT_NAME
	if v1, v2 := os.Getenv("OS_TENANT_NAME"), os.Getenv("OS_PROJECT_NAME"); v1 != v2 {
		// Allow setting OS_TENANT_NAME and OS_PROJECT_NAME if they are equal
		// Do nothing if one is set and the other isn't
		if v1 != "" && v2 != "" {
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_TENANT_NAME", "OS_PROJECT_NAME"},
			}
			return nilOptions, err
		}

	}

	// If OS_PROJECT_ID is set, overwrite tenantID with the value.
	if v := os.Getenv("OS_PROJECT_ID"); v != "" {
		tenantID = v
	}

	// If OS_PROJECT_NAME is set, overwrite tenantName with the value.
	if v := os.Getenv("OS_PROJECT_NAME"); v != "" {
		tenantName = v
	}

	// Disambiguate between OS_USER_DOMAIN_ID/OS_PROJECT_DOMAIN_ID/OS_DOMAIN_ID
	// 		and OS_USER_DOMAIN_NAME/OS_PROJECT_DOMAIN_NAME/OS_DOMAIN_NAME
	// Allow setting OS_USER_DOMAIN_ID and OS_PROJECT_DOMAIN_ID if they are equal
	if v1, v2 := os.Getenv("OS_USER_DOMAIN_ID"), os.Getenv("OS_PROJECT_DOMAIN_ID"); v1 != v2 {
		// Do nothing if one is set and the other isn't
		if v1 != "" && v2 != "" {
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_USER_DOMAIN_ID", "OS_PROJECT_DOMAIN_ID"},
			}
			return nilOptions, err
		}
	}

	// Allow setting OS_USER_DOMAIN_ID and OS_DOMAIN_ID if they are equal
	if v1, v2 := os.Getenv("OS_USER_DOMAIN_ID"), os.Getenv("OS_DOMAIN_ID"); v1 != v2 {
		// Do nothing if one is set and the other isn't
		if v1 != "" && v2 != "" {
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_USER_DOMAIN_ID", "OS_DOMAIN_ID"},
			}
			return nilOptions, err
		}
	}

	// Allow setting OS_DOMAIN_ID and OS_PROJECT_DOMAIN_ID if they are equal
	if v1, v2 := os.Getenv("OS_DOMAIN_ID"), os.Getenv("OS_PROJECT_DOMAIN_ID"); v1 != v2 {
		// Do nothing if one is set and the other isn't
		if v1 != "" && v2 != "" {
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_DOMAIN_ID", "OS_PROJECT_DOMAIN_ID"},
			}
			return nilOptions, err
		}
	}

	// Allow setting OS_USER_DOMAIN_ID, OS_PROJECT_DOMAIN_ID, and OS_DOMAIN_ID if they are equal
	if v1, v2, v3 := os.Getenv("OS_USER_DOMAIN_ID"), os.Getenv("OS_PROJECT_DOMAIN_ID"), os.Getenv("OS_DOMAIN_ID"); v1 != v2 || v2 != v3 {
		// Do nothing if any one is not set
		if v1 != "" && v2 != "" && v3 != "" {
			// If you get here, it's because all three have been set and they aren't all equal
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_DOMAIN_ID", "OS_USER_DOMAIN_ID", "OS_PROJECT_DOMAIN_ID"},
			}
			return nilOptions, err
		}
	}

	// Allow setting OS_USER_DOMAIN_NAME and OS_PROJECT_DOMAIN_NAME if they are equal
	if v1, v2 := os.Getenv("OS_USER_DOMAIN_NAME"), os.Getenv("OS_PROJECT_DOMAIN_NAME"); v1 != v2 {
		// Do nothing if one is set and the other isn't
		if v1 != "" && v2 != "" {
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_USER_DOMAIN_NAME", "OS_PROJECT_DOMAIN_NAME"},
			}
			return nilOptions, err
		}
	}

	// Allow setting OS_USER_DOMAIN_NAME and OS_DOMAIN_NAME if they are equal
	if v1, v2 := os.Getenv("OS_USER_DOMAIN_NAME"), os.Getenv("OS_DOMAIN_NAME"); v1 != v2 {
		// Do nothing if one is set and the other isn't
		if v1 != "" && v2 != "" {
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_USER_DOMAIN_NAME", "OS_DOMAIN_NAME"},
			}
			return nilOptions, err
		}
	}

	// Allow setting OS_DOMAIN_NAME and OS_PROJECT_DOMAIN_NAME if they are equal
	if v1, v2 := os.Getenv("OS_DOMAIN_NAME"), os.Getenv("OS_PROJECT_DOMAIN_NAME"); v1 != v2 {
		// Do nothing if one is set and the other isn't
		if v1 != "" && v2 != "" {
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_DOMAIN_NAME", "OS_PROJECT_DOMAIN_NAME"},
			}
			return nilOptions, err
		}
	}

	// Allow setting OS_USER_DOMAIN_NAME, OS_PROJECT_DOMAIN_NAME, and OS_DOMAIN_NAME if they are equal
	if v1, v2, v3 := os.Getenv("OS_USER_DOMAIN_NAME"), os.Getenv("OS_PROJECT_DOMAIN_NAME"), os.Getenv("OS_DOMAIN_NAME"); v1 != v2 || v2 != v3 {

		// Do nothing if any one is not set
		if v1 != "" && v2 != "" && v3 != "" {
			// If you get here, it's because all three have been set and they aren't all equal
			err := gophercloud.ErrEnvironmentVarsExpectedEqual{
				EnvironmentVariables: []string{"OS_DOMAIN_NAME", "OS_USER_DOMAIN_NAME", "OS_PROJECT_DOMAIN_NAME"},
			}
			return nilOptions, err
		}
	}

	// Keep track of how domainID or domainName is overwritten
	domainSources := map[string]string{
		"domainID":   "OS_DOMAIN_ID",
		"domainName": "OS_DOMAIN_NAME",
	}

	userDomainID := os.Getenv("OS_USER_DOMAIN_ID")
	projectDomainID := os.Getenv("OS_PROJECT_DOMAIN_ID")

	// If OS_USER_DOMAIN_ID or OS_PROJECT_DOMAIN_ID is set, overwrite domainID with the value.
	if userDomainID != "" && domainID == "" {
		domainID = userDomainID
		domainSources["domainID"] = "OS_USER_DOMAIN_ID"
		// No conflicts, OS_PROJECT_DOMAIN_ID provided
	} else if projectDomainID != "" && domainID == "" {
		domainID = projectDomainID
		domainSources["domainID"] = "OS_PROJECT_DOMAIN_ID"
	}

	// Check for cross-type conflicts: OS_DOMAIN_NAME.
	// Mirrors downstream logic where ID and NAME cannot be used together
	if (userDomainID != "" || projectDomainID != "") && domainName != "" {
		var conflictingVar string
		if userDomainID != "" {
			conflictingVar = "OS_USER_DOMAIN_ID"
		} else {
			conflictingVar = "OS_PROJECT_DOMAIN_ID"
		}
		err := gophercloud.ErrAmbiguousEnvironmentVarsClash{
			EnvironmentVariables: []string{"OS_DOMAIN_NAME", conflictingVar},
		}
		return nilOptions, err
	}

	userDomainName := os.Getenv("OS_USER_DOMAIN_NAME")
	projectDomainName := os.Getenv("OS_PROJECT_DOMAIN_NAME")

	// If OS_USER_DOMAIN_NAME or OS_PROJECT_DOMAIN_NAME is set, overwrite domainName with the value.
	if userDomainName != "" && domainName == "" {
		domainName = userDomainName
		domainSources["domainName"] = "OS_USER_DOMAIN_NAME"
		// No conflicts, OS_PROJECT_DOMAIN_NAME provided
	} else if projectDomainName != "" && domainName == "" {
		domainName = projectDomainName
		domainSources["domainName"] = "OS_PROJECT_DOMAIN_NAME"

	}

	// Check for cross-type conflicts: OS_DOMAIN_ID
	// Mirrors downstream logic where ID and NAME cannot be used together
	if (userDomainName != "" || projectDomainName != "") && domainID != "" {
		var conflictingVar string
		if userDomainName != "" {
			conflictingVar = "OS_USER_DOMAIN_NAME"
		} else {
			conflictingVar = "OS_PROJECT_DOMAIN_NAME"
		}
		err := gophercloud.ErrAmbiguousEnvironmentVarsClash{
			EnvironmentVariables: []string{"OS_DOMAIN_ID", conflictingVar},
		}
		return nilOptions, err
	}

	// Flag case where both domainID and domainName are both set by either OS_USER_DOMAIN_X or OS_PROJECT_DOMAIN_X
	// Mirrors downstream logic where ID and NAME cannot be used together
	// Downstream error handling deals with the case where OS_DOMAIN_ID and OS_DOMAIN_NAME are both set
	if domainID != "" && domainName != "" &&
		domainSources["domainID"] != "OS_DOMAIN_ID" &&
		domainSources["domainName"] != "OS_DOMAIN_NAME" {
		// One of OS_USER_DOMAIN_ID/OS_PROJECT_DOMAIN_ID or
		// OS_USER_DOMAIN_NAME/OS_PROJECT_DOMAIN_NAME can be used
		err := gophercloud.ErrAmbiguousEnvironmentVarsClash{
			EnvironmentVariables: []string{domainSources["domainName"], domainSources["domainID"]},
		}
		return nilOptions, err
	}

	if authURL == "" {
		err := gophercloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "OS_AUTH_URL",
		}
		return nilOptions, err
	}

	if userID == "" && username == "" {
		// Empty username and userID could be ignored, when applicationCredentialID and applicationCredentialSecret are set
		if applicationCredentialID == "" && applicationCredentialSecret == "" {
			err := gophercloud.ErrMissingAnyoneOfEnvironmentVariables{
				EnvironmentVariables: []string{"OS_USERID", "OS_USERNAME"},
			}
			return nilOptions, err
		}
	}

	if password == "" && passcode == "" && applicationCredentialID == "" && applicationCredentialName == "" {
		err := gophercloud.ErrMissingEnvironmentVariable{
			// silently ignore TOTP passcode warning, since it is not a common auth method
			EnvironmentVariable: "OS_PASSWORD",
		}
		return nilOptions, err
	}

	if (applicationCredentialID != "" || applicationCredentialName != "") && applicationCredentialSecret == "" {
		err := gophercloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "OS_APPLICATION_CREDENTIAL_SECRET",
		}
		return nilOptions, err
	}

	if domainID == "" && domainName == "" && tenantID == "" && tenantName != "" {
		err := gophercloud.ErrMissingEnvironmentVariable{
			EnvironmentVariable: "OS_PROJECT_ID",
		}
		return nilOptions, err
	}

	if applicationCredentialID == "" && applicationCredentialName != "" && applicationCredentialSecret != "" {
		if userID == "" && username == "" {
			return nilOptions, gophercloud.ErrMissingAnyoneOfEnvironmentVariables{
				EnvironmentVariables: []string{"OS_USERID", "OS_USERNAME"},
			}
		}
		if username != "" && domainID == "" && domainName == "" {
			return nilOptions, gophercloud.ErrMissingAnyoneOfEnvironmentVariables{
				EnvironmentVariables: []string{"OS_DOMAIN_ID", "OS_DOMAIN_NAME"},
			}
		}
	}

	var scope *gophercloud.AuthScope
	if systemScope == "all" {
		scope = &gophercloud.AuthScope{
			System: true,
		}
	}

	ao := gophercloud.AuthOptions{
		IdentityEndpoint:            authURL,
		UserID:                      userID,
		Username:                    username,
		Password:                    password,
		Passcode:                    passcode,
		TenantID:                    tenantID,
		TenantName:                  tenantName,
		DomainID:                    domainID,
		DomainName:                  domainName,
		ApplicationCredentialID:     applicationCredentialID,
		ApplicationCredentialName:   applicationCredentialName,
		ApplicationCredentialSecret: applicationCredentialSecret,
		Scope:                       scope,
	}

	return ao, nil
}
