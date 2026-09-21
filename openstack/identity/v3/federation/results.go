package federation

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type UserType string

const (
	// UserTypeEphemeral represents a federated user that does not exist in the
	// Identity service.
	UserTypeEphemeral UserType = "ephemeral"

	// UserTypeLocal represents an existing user in the Identity service.
	UserTypeLocal UserType = "local"
)

// IdentityProvider represents an identity provider trusted by the Identity
// service.
type IdentityProvider struct {
	// AuthorizationTTL is the length of validity, in minutes, for group
	// memberships carried over through mapping and persisted in the database.
	// It is nil when the Identity service configured default is used.
	AuthorizationTTL *int `json:"authorization_ttl"`

	// DomainID is the ID of the domain associated with the identity provider.
	DomainID string `json:"domain_id"`

	// Description describes the identity provider. It is nil when no
	// description is set.
	Description *string `json:"description"`

	// Enabled indicates whether the identity provider accepts federated
	// authentication requests.
	Enabled bool `json:"enabled"`

	// ID is the unique ID of the identity provider.
	ID string `json:"id"`

	// Links contains links related to the identity provider.
	Links map[string]any `json:"links"`

	// RemoteIDs is the list of unique IDs used by the remote identity provider.
	RemoteIDs []string `json:"remote_ids"`
}

// Mapping a set of rules to map federation protocol attributes to
// Identity API objects.
type Mapping struct {
	// The Federation Mapping unique ID
	ID string `json:"id"`

	// Links contains links related to the mapping.
	Links map[string]any `json:"links"`

	// The list of rules used to map remote users into local users
	Rules []MappingRule `json:"rules"`
}

type MappingRule struct {
	// References a local Identity API resource, such as a group or user to which the remote attributes will be mapped.
	Local []RuleLocal `json:"local"`

	// Each object contains a rule for mapping remote attributes to Identity API concepts.
	Remote []RuleRemote `json:"remote"`
}

type RuleRemote struct {
	// Type represents an assertion type keyword.
	Type string `json:"type"`

	// If true, then each string will be evaluated as a regular expression search against the remote attribute type.
	Regex *bool `json:"regex,omitempty"`

	// The rule is matched only if any of the specified strings appear in the remote attribute type.
	// This is mutually exclusive with NotAnyOf.
	AnyOneOf []string `json:"any_one_of,omitempty"`

	// The rule is not matched if any of the specified strings appear in the remote attribute type.
	// This is mutually exclusive with AnyOneOf.
	NotAnyOf []string `json:"not_any_of,omitempty"`

	// The rule works as a filter, removing any specified strings that are listed there from the remote attribute type.
	// This is mutually exclusive with Whitelist.
	Blacklist []string `json:"blacklist,omitempty"`

	// The rule works as a filter, allowing only the specified strings in the remote attribute type to be passed ahead.
	// This is mutually exclusive with Blacklist.
	Whitelist []string `json:"whitelist,omitempty"`
}

type RuleLocal struct {
	// Domain to which the remote attributes will be matched.
	Domain *Domain `json:"domain,omitempty"`

	// Group to which the remote attributes will be matched.
	Group *Group `json:"group,omitempty"`

	// Group IDs to which the remote attributes will be matched.
	GroupIDs string `json:"group_ids,omitempty"`

	// Groups to which the remote attributes will be matched.
	Groups string `json:"groups,omitempty"`

	// Projects to which the remote attributes will be matched.
	Projects []RuleProject `json:"projects,omitempty"`

	// User to which the remote attributes will be matched.
	User *RuleUser `json:"user,omitempty"`
}

type Domain struct {
	// Domain ID
	// This is mutually exclusive with Name.
	ID string `json:"id,omitempty"`

	// Domain Name
	// This is mutually exclusive with ID.
	Name string `json:"name,omitempty"`
}

type Group struct {
	// Group ID to which the rule should match.
	// This is mutually exclusive with Name and Domain.
	ID string `json:"id,omitempty"`

	// Group Name to which the rule should match.
	// This is mutually exclusive with ID.
	Name string `json:"name,omitempty"`

	// Group Domain to which the rule should match.
	// This is mutually exclusive with ID.
	Domain *Domain `json:"domain,omitempty"`
}

type RuleProject struct {
	// Project name
	Name string `json:"name,omitempty"`

	// Project roles
	Roles []RuleProjectRole `json:"roles,omitempty"`
}

type RuleProjectRole struct {
	// Role name
	Name string `json:"name,omitempty"`
}

type RuleUser struct {
	// User domain
	Domain *Domain `json:"domain,omitempty"`

	// User email
	Email string `json:"email,omitempty"`

	// User ID
	ID string `json:"id,omitempty"`

	// User name
	Name string `json:"name,omitempty"`

	// User type
	Type *UserType `json:"type,omitempty"`
}

type mappingResult struct {
	gophercloud.Result
}

// Extract interprets any mappingResult as a Mapping.
func (c mappingResult) Extract() (*Mapping, error) {
	var s struct {
		Mapping *Mapping `json:"mapping"`
	}
	err := c.ExtractInto(&s)
	return s.Mapping, err
}

// CreateMappingResult is the response from a CreateMapping operation.
// Call its Extract method to interpret it as a Mapping.
type CreateMappingResult struct {
	mappingResult
}

// GetMappingResult is the response from a GetMapping operation.
// Call its Extract method to interpret it as a Mapping.
type GetMappingResult struct {
	mappingResult
}

// UpdateMappingResult is the response from a UpdateMapping operation.
// Call its Extract method to interpret it as a Mapping.
type UpdateMappingResult struct {
	mappingResult
}

// DeleteMappingResult is the response from a DeleteMapping operation.
// Call its ExtractErr to determine if the request succeeded or failed.
type DeleteMappingResult struct {
	gophercloud.ErrResult
}

type identityProviderResult struct {
	gophercloud.Result
}

// Extract interprets an identity provider result as an IdentityProvider.
func (r identityProviderResult) Extract() (*IdentityProvider, error) {
	var s struct {
		IdentityProvider *IdentityProvider `json:"identity_provider"`
	}
	err := r.ExtractInto(&s)
	return s.IdentityProvider, err
}

// CreateIdentityProviderResult is the response from a CreateIdentityProvider
// operation.
type CreateIdentityProviderResult struct {
	identityProviderResult
}

// GetIdentityProviderResult is the response from a GetIdentityProvider
// operation.
type GetIdentityProviderResult struct {
	identityProviderResult
}

// UpdateIdentityProviderResult is the response from an UpdateIdentityProvider
// operation.
type UpdateIdentityProviderResult struct {
	identityProviderResult
}

// DeleteIdentityProviderResult is the response from a DeleteIdentityProvider
// operation.
type DeleteIdentityProviderResult struct {
	gophercloud.ErrResult
}

// IdentityProvidersPage is a single page of IdentityProvider results.
type IdentityProvidersPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether an IdentityProvidersPage contains any results.
func (r IdentityProvidersPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	identityProviders, err := ExtractIdentityProviders(r)
	return len(identityProviders) == 0, err
}

// NextPageURL extracts the next link from an IdentityProvidersPage.
func (r IdentityProvidersPage) NextPageURL(endpointURL string) (string, error) {
	var s struct {
		Links struct {
			Next string `json:"next"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	return s.Links.Next, err
}

// ExtractIdentityProviders extracts a slice of IdentityProvider values from a
// page returned by ListIdentityProviders.
func ExtractIdentityProviders(r pagination.Page) ([]IdentityProvider, error) {
	var s struct {
		IdentityProviders []IdentityProvider `json:"identity_providers"`
	}
	err := (r.(IdentityProvidersPage)).ExtractInto(&s)
	return s.IdentityProviders, err
}

// MappingsPage is a single page of Mapping results.
type MappingsPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of Mappings contains any results.
func (c MappingsPage) IsEmpty() (bool, error) {
	if c.StatusCode == 204 {
		return true, nil
	}

	mappings, err := ExtractMappings(c)
	return len(mappings) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (c MappingsPage) NextPageURL(endpointURL string) (string, error) {
	var s struct {
		Links struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"links"`
	}
	err := c.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Links.Next, err
}

// ExtractMappings returns a slice of Mappings contained in a single page of
// results.
func ExtractMappings(r pagination.Page) ([]Mapping, error) {
	var s struct {
		Mappings []Mapping `json:"mappings"`
	}
	err := (r.(MappingsPage)).ExtractInto(&s)
	return s.Mappings, err
}

// Protocol represents a federation protocol associated with an identity
// provider.
type Protocol struct {
	// ID is the unique ID of the protocol.
	ID string `json:"id"`

	// Links contains links related to the protocol.
	Links map[string]any `json:"links"`

	// MappingID is the ID of the mapping used by the protocol.
	MappingID string `json:"mapping_id"`

	// RemoteIDAttribute is the name of the attribute used to obtain the remote
	// identity provider ID. It is nil when the Identity service omits the
	// attribute from the response.
	RemoteIDAttribute *string `json:"remote_id_attribute"`
}

type protocolResult struct {
	gophercloud.Result
}

// Extract interprets a protocol result as a Protocol.
func (r protocolResult) Extract() (*Protocol, error) {
	var s struct {
		Protocol *Protocol `json:"protocol"`
	}
	err := r.ExtractInto(&s)
	return s.Protocol, err
}

// CreateProtocolResult is the response from a CreateProtocol operation.
type CreateProtocolResult struct {
	protocolResult
}

// GetProtocolResult is the response from a GetProtocol operation.
type GetProtocolResult struct {
	protocolResult
}

// UpdateProtocolResult is the response from an UpdateProtocol operation.
type UpdateProtocolResult struct {
	protocolResult
}

// DeleteProtocolResult is the response from a DeleteProtocol operation.
type DeleteProtocolResult struct {
	gophercloud.ErrResult
}

// ProtocolsPage is a single page of Protocol results.
type ProtocolsPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether a ProtocolsPage contains any results.
func (r ProtocolsPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	protocols, err := ExtractProtocols(r)
	return len(protocols) == 0, err
}

// NextPageURL extracts the next link from a ProtocolsPage.
func (r ProtocolsPage) NextPageURL(endpointURL string) (string, error) {
	var s struct {
		Links struct {
			Next string `json:"next"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	return s.Links.Next, err
}

// ExtractProtocols extracts a slice of Protocol values from a page returned by
// ListProtocols.
func ExtractProtocols(r pagination.Page) ([]Protocol, error) {
	var s struct {
		Protocols []Protocol `json:"protocols"`
	}
	err := (r.(ProtocolsPage)).ExtractInto(&s)
	return s.Protocols, err
}

// ServiceProvider represents a service provider trusted by the Identity
// service.
type ServiceProvider struct {
	// AuthURL is the protected URL where tokens can be retrieved after the user
	// is authenticated.
	AuthURL string `json:"auth_url"`

	// Description describes the service provider. It is nil when no
	// description is set.
	Description *string `json:"description"`

	// Enabled indicates whether bursting into the service provider is enabled.
	Enabled bool `json:"enabled"`

	// ID is the unique ID of the service provider.
	ID string `json:"id"`

	// Links contains links related to the service provider.
	Links map[string]any `json:"links"`

	// RelayStatePrefix is the relay state prefix used in ECP-wrapped SAML
	// messages. It is nil when no relay state prefix is set.
	RelayStatePrefix *string `json:"relay_state_prefix"`

	// SPURL is the URL at the remote peer where the assertion is sent.
	SPURL string `json:"sp_url"`
}

type serviceProviderResult struct {
	gophercloud.Result
}

// Extract interprets a service provider result as a ServiceProvider.
func (r serviceProviderResult) Extract() (*ServiceProvider, error) {
	var s struct {
		ServiceProvider *ServiceProvider `json:"service_provider"`
	}
	err := r.ExtractInto(&s)
	return s.ServiceProvider, err
}

// CreateServiceProviderResult is the response from a CreateServiceProvider
// operation.
type CreateServiceProviderResult struct {
	serviceProviderResult
}

// GetServiceProviderResult is the response from a GetServiceProvider
// operation.
type GetServiceProviderResult struct {
	serviceProviderResult
}

// UpdateServiceProviderResult is the response from an UpdateServiceProvider
// operation.
type UpdateServiceProviderResult struct {
	serviceProviderResult
}

// DeleteServiceProviderResult is the response from a DeleteServiceProvider
// operation.
type DeleteServiceProviderResult struct {
	gophercloud.ErrResult
}

// ServiceProvidersPage is a single page of ServiceProvider results.
type ServiceProvidersPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether a ServiceProvidersPage contains any results.
func (r ServiceProvidersPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	serviceProviders, err := ExtractServiceProviders(r)
	return len(serviceProviders) == 0, err
}

// NextPageURL extracts the next link from a ServiceProvidersPage.
func (r ServiceProvidersPage) NextPageURL(endpointURL string) (string, error) {
	var s struct {
		Links struct {
			Next string `json:"next"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	return s.Links.Next, err
}

// ExtractServiceProviders extracts a slice of ServiceProvider values from a
// page returned by ListServiceProviders.
func ExtractServiceProviders(r pagination.Page) ([]ServiceProvider, error) {
	var s struct {
		ServiceProviders []ServiceProvider `json:"service_providers"`
	}
	err := (r.(ServiceProvidersPage)).ExtractInto(&s)
	return s.ServiceProviders, err
}
