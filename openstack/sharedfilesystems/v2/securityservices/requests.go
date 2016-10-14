package securityservices

import "github.com/gophercloud/gophercloud"

type SecurityServiceType string

// Valid security service types
const (
	LDAP            SecurityServiceType = "ldap"
	Kerberos        SecurityServiceType = "kerberos"
	ActiveDirectory SecurityServiceType = "active_directory"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToSecurityServiceCreateMap() (map[string]interface{}, error)
}

// CreateOpts contains options for creating a SecurityService. This object is
// passed to the securityservices.Create function. For more information about
// these parameters, see the SecurityService object.
type CreateOpts struct {
	// The security service type. A valid value is ldap, kerberos, or active_directory
	Type SecurityServiceType `json:"type" required:"true"`
	// The security service name
	Name string `json:"name,omitempty"`
	// The security service description
	Description string `json:"description,omitempty"`
	// The DNS IP address that is used inside the tenant network
	DNSIP string `json:"dns_ip,omitempty"`
	// The security service user or group name that is used by the tenant
	User string `json:"user,omitempty"`
	// The user password, if you specify a user
	Password string `json:"password,omitempty"`
	// The security service domain
	Domain string `json:"domain,omitempty"`
	// The security service host name or IP address
	Server string `json:"server,omitempty"`
}

// ToSecurityServicesCreateMap assembles a request body based on the contents of a
// CreateOpts.
func (opts CreateOpts) ToSecurityServiceCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "security_service")
}

// Create will create a new SecurityService based on the values in CreateOpts. To
// extract the SecurityService object from the response, call the Extract method
// on the CreateResult.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSecurityServiceCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete will delete the existing SecurityService with the provided ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}
