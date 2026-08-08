package federation

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListIdentityProviders enumerates the identity providers.
func ListIdentityProviders(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, identityProvidersRootURL(client), func(r pagination.PageResult) pagination.Page {
		return IdentityProvidersPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateIdentityProviderOptsBuilder allows extensions to add additional
// parameters to the CreateIdentityProvider request.
type CreateIdentityProviderOptsBuilder interface {
	ToIdentityProviderCreateMap() (map[string]any, error)
}

// CreateIdentityProviderOpts provides options for creating an identity
// provider.
type CreateIdentityProviderOpts struct {
	// AuthorizationTTL is the length of validity, in minutes, for group
	// memberships carried over through mapping and persisted in the database.
	AuthorizationTTL *int `json:"authorization_ttl,omitempty"`

	// DomainID is the ID of the domain associated with the identity provider.
	// If omitted, the Identity service creates a domain automatically.
	DomainID string `json:"domain_id,omitempty"`

	// Description describes the identity provider.
	Description string `json:"description,omitempty"`

	// Enabled indicates whether the identity provider accepts federated
	// authentication requests.
	Enabled *bool `json:"enabled,omitempty"`

	// RemoteIDs is the list of unique IDs used by the remote identity provider.
	RemoteIDs []string `json:"remote_ids,omitempty"`
}

// ToIdentityProviderCreateMap formats CreateIdentityProviderOpts into a create
// request.
func (opts CreateIdentityProviderOpts) ToIdentityProviderCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "identity_provider")
}

// CreateIdentityProvider registers a new identity provider.
func CreateIdentityProvider(ctx context.Context, client *gophercloud.ServiceClient, identityProviderID string, opts CreateIdentityProviderOptsBuilder) (r CreateIdentityProviderResult) {
	b, err := opts.ToIdentityProviderCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, identityProvidersResourceURL(client, identityProviderID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetIdentityProvider retrieves details about an identity provider.
func GetIdentityProvider(ctx context.Context, client *gophercloud.ServiceClient, identityProviderID string) (r GetIdentityProviderResult) {
	resp, err := client.Get(ctx, identityProvidersResourceURL(client, identityProviderID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateIdentityProviderOptsBuilder allows extensions to add additional
// parameters to the UpdateIdentityProvider request.
type UpdateIdentityProviderOptsBuilder interface {
	ToIdentityProviderUpdateMap() (map[string]any, error)
}

// UpdateIdentityProviderOpts provides options for updating an identity
// provider. The associated domain cannot be updated.
type UpdateIdentityProviderOpts struct {
	// AuthorizationTTL is the length of validity, in minutes, for group
	// memberships carried over through mapping and persisted in the database.
	// Set this to -1 to reset the value to the Identity service configured
	// default.
	AuthorizationTTL *int `json:"authorization_ttl,omitempty"`

	// Description describes the identity provider.
	// Set this to an empty string to remove the description.
	Description *string `json:"description,omitempty"`

	// Enabled indicates whether the identity provider accepts federated
	// authentication requests.
	Enabled *bool `json:"enabled,omitempty"`

	// RemoteIDs is the list of unique IDs used by the remote identity provider.
	// A pointer to an empty slice removes all remote IDs.
	RemoteIDs *[]string `json:"remote_ids,omitempty"`
}

// ToIdentityProviderUpdateMap formats UpdateIdentityProviderOpts into an update
// request.
func (opts UpdateIdentityProviderOpts) ToIdentityProviderUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "identity_provider")
	if err != nil {
		return nil, err
	}

	m := b["identity_provider"].(map[string]any)

	// Keystone distinguishes an omitted authorization_ttl from a null value.
	// The former leaves the current value unchanged, while the latter resets it
	// to the configured default. Zero is a valid value, so -1 is used as the
	// sentinel for null.
	if opts.AuthorizationTTL != nil && *opts.AuthorizationTTL == -1 {
		m["authorization_ttl"] = nil
	}

	// An empty string is used as the sentinel for a null description.
	if opts.Description != nil && *opts.Description == "" {
		m["description"] = nil
	}

	return b, nil
}

// UpdateIdentityProvider updates an existing identity provider.
func UpdateIdentityProvider(ctx context.Context, client *gophercloud.ServiceClient, identityProviderID string, opts UpdateIdentityProviderOptsBuilder) (r UpdateIdentityProviderResult) {
	b, err := opts.ToIdentityProviderUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, identityProvidersResourceURL(client, identityProviderID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteIdentityProvider deletes an identity provider.
func DeleteIdentityProvider(ctx context.Context, client *gophercloud.ServiceClient, identityProviderID string) (r DeleteIdentityProviderResult) {
	resp, err := client.Delete(ctx, identityProvidersResourceURL(client, identityProviderID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListMappings enumerates the mappings.
func ListMappings(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, mappingsRootURL(client), func(r pagination.PageResult) pagination.Page {
		return MappingsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateMappingOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateMappingOptsBuilder interface {
	ToMappingCreateMap() (map[string]any, error)
}

// CreateMappingOpts provides options for creating a mapping.
type CreateMappingOpts struct {
	// The list of rules used to map remote users into local users
	Rules []MappingRule `json:"rules"`
}

// ToMappingCreateMap formats a CreateMappingOpts into a create request.
func (opts CreateMappingOpts) ToMappingCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "mapping")
}

// CreateMapping creates a new Mapping.
func CreateMapping(ctx context.Context, client *gophercloud.ServiceClient, mappingID string, opts CreateMappingOptsBuilder) (r CreateMappingResult) {
	b, err := opts.ToMappingCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, mappingsResourceURL(client, mappingID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetMapping retrieves details on a single mapping, by ID.
func GetMapping(ctx context.Context, client *gophercloud.ServiceClient, mappingID string) (r GetMappingResult) {
	resp, err := client.Get(ctx, mappingsResourceURL(client, mappingID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateMappingOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateMappingOptsBuilder interface {
	ToMappingUpdateMap() (map[string]any, error)
}

// UpdateMappingOpts provides options for updating a mapping.
type UpdateMappingOpts struct {
	// The list of rules used to map remote users into local users
	Rules []MappingRule `json:"rules"`
}

// ToMappingUpdateMap formats an UpdateMappingOpts into an update request.
func (opts UpdateMappingOpts) ToMappingUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "mapping")
}

// UpdateMapping updates an existing mapping.
func UpdateMapping(ctx context.Context, client *gophercloud.ServiceClient, mappingID string, opts UpdateMappingOptsBuilder) (r UpdateMappingResult) {
	b, err := opts.ToMappingUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, mappingsResourceURL(client, mappingID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMapping deletes a mapping.
func DeleteMapping(ctx context.Context, client *gophercloud.ServiceClient, mappingID string) (r DeleteMappingResult) {
	resp, err := client.Delete(ctx, mappingsResourceURL(client, mappingID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListProtocols enumerates the federation protocols associated with an
// identity provider.
func ListProtocols(client *gophercloud.ServiceClient, identityProviderID string) pagination.Pager {
	return pagination.NewPager(client, protocolsRootURL(client, identityProviderID), func(r pagination.PageResult) pagination.Page {
		return ProtocolsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateProtocolOptsBuilder allows extensions to add additional parameters to
// the CreateProtocol request.
type CreateProtocolOptsBuilder interface {
	ToProtocolCreateMap() (map[string]any, error)
}

// CreateProtocolOpts provides options for adding a protocol to an identity
// provider.
type CreateProtocolOpts struct {
	// MappingID is the ID of the mapping used to process federated
	// authentication requests.
	MappingID string `json:"mapping_id" required:"true"`

	// RemoteIDAttribute is the name of the attribute used to obtain the remote
	// identity provider ID.
	RemoteIDAttribute string `json:"remote_id_attribute,omitempty"`
}

// ToProtocolCreateMap formats CreateProtocolOpts into a create request.
func (opts CreateProtocolOpts) ToProtocolCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "protocol")
}

// CreateProtocol adds a protocol to an identity provider.
func CreateProtocol(ctx context.Context, client *gophercloud.ServiceClient, identityProviderID, protocolID string, opts CreateProtocolOptsBuilder) (r CreateProtocolResult) {
	b, err := opts.ToProtocolCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, protocolsResourceURL(client, identityProviderID, protocolID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetProtocol retrieves details about a protocol associated with an identity
// provider.
func GetProtocol(ctx context.Context, client *gophercloud.ServiceClient, identityProviderID, protocolID string) (r GetProtocolResult) {
	resp, err := client.Get(ctx, protocolsResourceURL(client, identityProviderID, protocolID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateProtocolOptsBuilder allows extensions to add additional parameters to
// the UpdateProtocol request.
type UpdateProtocolOptsBuilder interface {
	ToProtocolUpdateMap() (map[string]any, error)
}

// UpdateProtocolOpts provides options for updating a protocol.
type UpdateProtocolOpts struct {
	// MappingID is the ID of the mapping used to process federated
	// authentication requests.
	MappingID string `json:"mapping_id,omitempty"`

	// RemoteIDAttribute is the name of the attribute used to obtain the remote
	// identity provider ID.
	RemoteIDAttribute *string `json:"remote_id_attribute,omitempty"`
}

// ToProtocolUpdateMap formats UpdateProtocolOpts into an update request.
func (opts UpdateProtocolOpts) ToProtocolUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "protocol")
}

// UpdateProtocol updates a protocol associated with an identity provider.
func UpdateProtocol(ctx context.Context, client *gophercloud.ServiceClient, identityProviderID, protocolID string, opts UpdateProtocolOptsBuilder) (r UpdateProtocolResult) {
	b, err := opts.ToProtocolUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, protocolsResourceURL(client, identityProviderID, protocolID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteProtocol deletes a protocol from an identity provider.
func DeleteProtocol(ctx context.Context, client *gophercloud.ServiceClient, identityProviderID, protocolID string) (r DeleteProtocolResult) {
	resp, err := client.Delete(ctx, protocolsResourceURL(client, identityProviderID, protocolID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListServiceProviders enumerates the service providers.
func ListServiceProviders(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, serviceProvidersRootURL(client), func(r pagination.PageResult) pagination.Page {
		return ServiceProvidersPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateServiceProviderOptsBuilder allows extensions to add additional
// parameters to the CreateServiceProvider request.
type CreateServiceProviderOptsBuilder interface {
	ToServiceProviderCreateMap() (map[string]any, error)
}

// CreateServiceProviderOpts provides options for creating a service provider.
type CreateServiceProviderOpts struct {
	// AuthURL is the protected URL where tokens can be retrieved after the user
	// is authenticated.
	AuthURL string `json:"auth_url" required:"true"`

	// Description describes the service provider.
	Description string `json:"description,omitempty"`

	// Enabled indicates whether bursting into the service provider is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// RelayStatePrefix is the relay state prefix used in ECP-wrapped SAML
	// messages.
	RelayStatePrefix string `json:"relay_state_prefix,omitempty"`

	// SPURL is the URL at the remote peer where the assertion is sent.
	SPURL string `json:"sp_url" required:"true"`
}

// ToServiceProviderCreateMap formats CreateServiceProviderOpts into a create
// request.
func (opts CreateServiceProviderOpts) ToServiceProviderCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "service_provider")
}

// CreateServiceProvider registers a new service provider.
func CreateServiceProvider(ctx context.Context, client *gophercloud.ServiceClient, serviceProviderID string, opts CreateServiceProviderOptsBuilder) (r CreateServiceProviderResult) {
	b, err := opts.ToServiceProviderCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, serviceProvidersResourceURL(client, serviceProviderID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetServiceProvider retrieves details about a service provider.
func GetServiceProvider(ctx context.Context, client *gophercloud.ServiceClient, serviceProviderID string) (r GetServiceProviderResult) {
	resp, err := client.Get(ctx, serviceProvidersResourceURL(client, serviceProviderID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateServiceProviderOptsBuilder allows extensions to add additional
// parameters to the UpdateServiceProvider request.
type UpdateServiceProviderOptsBuilder interface {
	ToServiceProviderUpdateMap() (map[string]any, error)
}

// UpdateServiceProviderOpts provides options for updating a service provider.
type UpdateServiceProviderOpts struct {
	// AuthURL is the protected URL where tokens can be retrieved after the user
	// is authenticated.
	AuthURL *string `json:"auth_url,omitempty"`

	// Description describes the service provider.
	// Set this to an empty string to remove the description.
	Description *string `json:"description,omitempty"`

	// Enabled indicates whether bursting into the service provider is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// RelayStatePrefix is the relay state prefix used in ECP-wrapped SAML
	// messages.
	// Set this to an empty string to reset the relay state prefix.
	RelayStatePrefix *string `json:"relay_state_prefix,omitempty"`

	// SPURL is the URL at the remote peer where the assertion is sent.
	SPURL *string `json:"sp_url,omitempty"`
}

// ToServiceProviderUpdateMap formats UpdateServiceProviderOpts into an update
// request.
func (opts UpdateServiceProviderOpts) ToServiceProviderUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "service_provider")
	if err != nil {
		return nil, err
	}

	m := b["service_provider"].(map[string]any)

	// Empty strings are used as sentinels for nullable fields.
	if opts.Description != nil && *opts.Description == "" {
		m["description"] = nil
	}
	if opts.RelayStatePrefix != nil && *opts.RelayStatePrefix == "" {
		m["relay_state_prefix"] = nil
	}

	return b, nil
}

// UpdateServiceProvider updates an existing service provider.
func UpdateServiceProvider(ctx context.Context, client *gophercloud.ServiceClient, serviceProviderID string, opts UpdateServiceProviderOptsBuilder) (r UpdateServiceProviderResult) {
	b, err := opts.ToServiceProviderUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, serviceProvidersResourceURL(client, serviceProviderID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteServiceProvider deletes a service provider.
func DeleteServiceProvider(ctx context.Context, client *gophercloud.ServiceClient, serviceProviderID string) (r DeleteServiceProviderResult) {
	resp, err := client.Delete(ctx, serviceProvidersResourceURL(client, serviceProviderID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
