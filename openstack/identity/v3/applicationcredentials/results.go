package applicationcredentials

import (
	"encoding/json"
	"time"

	"github.com/bizflycloud/gophercloud"
	"github.com/bizflycloud/gophercloud/pagination"
)

type Role struct {
	// DomainID is the domain ID the role belongs to.
	DomainID string `json:"domain_id,omitempty"`
	// ID is the unique ID of the role.
	ID string `json:"id,omitempty"`
	// Name is the role name
	Name string `json:"name,omitempty"`
}

// ApplicationCredential represents the access rule object
type AccessRule struct {
	// The ID of the access rule
	ID string `json:"id,omitempty"`
	// The API path that the application credential is permitted to access
	Path string `json:"path,omitempty"`
	// The request method that the application credential is permitted to use for a
	// given API endpoint
	Method string `json:"method,omitempty"`
	// The service type identifier for the service that the application credential
	// is permitted to access
	Service string `json:"service,omitempty"`
}

// ApplicationCredential represents the application credential object
type ApplicationCredential struct {
	// The ID of the application credential.
	ID string `json:"id"`
	// The name of the application credential.
	Name string `json:"name"`
	// A description of the application credential’s purpose.
	Description string `json:"description"`
	// A flag indicating whether the application credential may be used for creation or destruction of other application credentials or trusts.
	// Defaults to false
	Unrestricted bool `json:"unrestricted"`
	// The secret for the application credential, either generated by the server or provided by the user.
	// This is only ever shown once in the response to a create request. It is not stored nor ever shown again.
	// If the secret is lost, a new application credential must be created.
	Secret string `json:"secret"`
	// The ID of the project the application credential was created for and that authentication requests using this application credential will be scoped to.
	ProjectID string `json:"project_id"`
	// A list of one or more roles that this application credential has associated with its project.
	// A token using this application credential will have these same roles.
	Roles []Role `json:"roles"`
	// The expiration time of the application credential, if one was specified.
	ExpiresAt time.Time `json:"-"`
	// A list of access rules objects.
	AccessRules []AccessRule `json:"access_rules,omitempty"`
	// Links contains referencing links to the application credential.
	Links map[string]interface{} `json:"links"`
}

func (r *ApplicationCredential) UnmarshalJSON(b []byte) error {
	type tmp ApplicationCredential
	var s struct {
		tmp
		ExpiresAt gophercloud.JSONRFC3339MilliNoZ `json:"expires_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = ApplicationCredential(s.tmp)

	r.ExpiresAt = time.Time(s.ExpiresAt)

	return nil
}

type applicationCredentialResult struct {
	gophercloud.Result
}

// GetResult is the response from a Get operation. Call its Extract method
// to interpret it as an ApplicationCredential.
type GetResult struct {
	applicationCredentialResult
}

// CreateResult is the response from a Create operation. Call its Extract method
// to interpret it as an ApplicationCredential.
type CreateResult struct {
	applicationCredentialResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr to
// determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// an ApplicationCredentialPage is a single page of an ApplicationCredential results.
type ApplicationCredentialPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a an ApplicationCredentialPage contains any results.
func (r ApplicationCredentialPage) IsEmpty() (bool, error) {
	applicationCredentials, err := ExtractApplicationCredentials(r)
	return len(applicationCredentials) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r ApplicationCredentialPage) NextPageURL() (string, error) {
	var s struct {
		Links struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Links.Next, err
}

// Extractan ApplicationCredentials returns a slice of ApplicationCredentials contained in a single page of results.
func ExtractApplicationCredentials(r pagination.Page) ([]ApplicationCredential, error) {
	var s struct {
		ApplicationCredentials []ApplicationCredential `json:"application_credentials"`
	}
	err := (r.(ApplicationCredentialPage)).ExtractInto(&s)
	return s.ApplicationCredentials, err
}

// Extract interprets any application_credential results as an ApplicationCredential.
func (r applicationCredentialResult) Extract() (*ApplicationCredential, error) {
	var s struct {
		ApplicationCredential *ApplicationCredential `json:"application_credential"`
	}
	err := r.ExtractInto(&s)
	return s.ApplicationCredential, err
}

// GetAccessRuleResult is the response from a Get operation. Call its Extract method
// to interpret it as an AccessRule.
type GetAccessRuleResult struct {
	gophercloud.Result
}

// an AccessRulePage is a single page of an AccessRule results.
type AccessRulePage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a an AccessRulePage contains any results.
func (r AccessRulePage) IsEmpty() (bool, error) {
	accessRules, err := ExtractAccessRules(r)
	return len(accessRules) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r AccessRulePage) NextPageURL() (string, error) {
	var s struct {
		Links struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Links.Next, err
}

// ExtractAccessRules returns a slice of AccessRules contained in a single page of results.
func ExtractAccessRules(r pagination.Page) ([]AccessRule, error) {
	var s struct {
		AccessRules []AccessRule `json:"access_rules"`
	}
	err := (r.(AccessRulePage)).ExtractInto(&s)
	return s.AccessRules, err
}

// Extract interprets any access_rule results as an AccessRule.
func (r GetAccessRuleResult) Extract() (*AccessRule, error) {
	var s struct {
		AccessRule *AccessRule `json:"access_rule"`
	}
	err := r.ExtractInto(&s)
	return s.AccessRule, err
}
