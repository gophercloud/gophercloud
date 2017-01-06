package securityservices

import "github.com/gophercloud/gophercloud"

// SecurityService contains all the information associated with an OpenStack
// SecurityService.
type SecurityService struct {
	// The security service ID
	ID string `json:"id"`
	// The UUID of the project where the security service was created
	ProjectID string `json:"project_id"`
	// The security service domain
	Domain string `json:"domain"`
	// The security service status
	Status string `json:"status"`
	// The security service type. A valid value is ldap, kerberos, or active_directory
	Type string `json:"type"`
	// The security service name
	Name string `json:"name"`
	// The security service description
	Description string `json:"description"`
	// The DNS IP address that is used inside the tenant network
	DNSIP string `json:"dns_ip"`
	// The security service user or group name that is used by the tenant
	User string `json:"user"`
	// The user password, if you specify a user
	Password string `json:"password"`
	// The security service host name or IP address
	Server string `json:"server"`
	// The date and time stamp when the security service was created
	CreatedAt gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
	// The date and time stamp when the security service was updated
	UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the SecurityService object out of the commonResult object.
func (r commonResult) Extract() (*SecurityService, error) {
	var s struct {
		SecurityService *SecurityService `json:"security_service"`
	}
	err := r.ExtractInto(&s)
	return s.SecurityService, err
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}
