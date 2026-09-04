package auth

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
)

// AuthType respresents a valid method of authentication.
type AuthType string

const (
	// AuthPassword defines an unknown version of the password
	AuthPassword AuthType = "password"

	// AuthToken defined an unknown version of the token
	AuthToken AuthType = "token"

	// AuthV2Password defines version 2 of the password
	AuthV2Password AuthType = "v2password"

	// AuthV2Token defines version 2 of the token
	AuthV2Token AuthType = "v2token"

	// AuthV3Password defines version 3 of the password
	AuthV3Password AuthType = "v3password"

	// AuthV3Totp defines version 3 of the totp
	AuthV3Totp AuthType = "v3totp"

	// AuthV3ApplicationCredential defines version 3 of the application credential
	AuthV3ApplicationCredential AuthType = "v3applicationcredential"

	// AuthV3Token defines version 3 of the token
	AuthV3Token AuthType = "v3token"

	// AuthV3MultiFactor defines version 3 of the multifactor
	AuthV3MultiFactor AuthType = "v3multifactor"
)

// Helper that returns auth method used in request bodies
func (at AuthType) toAuthMethod() string {
	switch at {
	case AuthV3Password:
		return "password"
	case AuthV3Totp:
		return "totp"
	case AuthV3ApplicationCredential:
		return "application_credential"
	case AuthV3Token:
		return "token"
	case AuthV3MultiFactor:
		return ""
	default:
		return ""
	}
}

type AuthOptionsBuilder interface {
	ToAuthBody() (map[string]map[string]any, error)
	CanReauth() bool
}

type AuthOptionsBuilderV2 interface {
	AuthOptionsBuilder
}

type AuthOptionsBuilderV3 interface {
	AuthOptionsBuilder
	ToAuthHeaders() (map[string]any, error)
	ToAuthScope() (map[string]any, error)
}

type Authenticator interface {
	Authenticate(ctx context.Context, httpClient *http.Client) (*AuthResult, error)
	GetAuthURL() string
}

type AuthOptionsV2 struct {
	AuthURL string
	Auth    AuthOptionsBuilderV2
}

func (ao AuthOptionsV2) GetAuthURL() string {
	base := strings.TrimSuffix(ao.AuthURL, "/")
	return base + "/v2.0/"
}

func (ao AuthOptionsV2) Authenticate(ctx context.Context, httpClient *http.Client) (*AuthResult, error) {
	if ao.Auth == nil {
		return nil, gophercloud.ErrMissingInput{Argument: "Auth"}
	}

	authData, err := ao.Auth.ToAuthBody()
	if err != nil {
		return nil, err
	}

	var body map[string]any
	for _, v := range authData {
		body = v
	}
	if body == nil {
		return nil, gophercloud.ErrMissingInput{Argument: "Auth"}
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	client := &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{HTTPClient: *httpClient},
		Endpoint:       ao.GetAuthURL(),
	}

	var result gophercloud.Result
	resp, err := client.Post(ctx, client.ServiceURL("tokens"), map[string]any{"auth": body}, &result.Body, &gophercloud.RequestOpts{
		OkCodes:     []int{200, 203},
		OmitHeaders: []string{"X-Auth-Token"},
	})
	_, result.Header, result.Err = gophercloud.ParseResponse(resp, err)
	if result.Err != nil {
		return nil, result.Err
	}

	var respBody v2AccessBody
	if err := result.ExtractInto(&respBody); err != nil {
		return nil, err
	}

	return respBody.toAuthResult(ao.Auth.CanReauth()), nil
}

type AuthOptionsV3 struct {
	AuthURL string
	Auth    AuthOptionsBuilderV3
}

func (ao AuthOptionsV3) GetAuthURL() string {
	base := strings.TrimSuffix(ao.AuthURL, "/")
	return base + "/v3/"
}

func (ao AuthOptionsV3) Authenticate(ctx context.Context, httpClient *http.Client) (*AuthResult, error) {
	if ao.Auth == nil {
		return nil, gophercloud.ErrMissingInput{Argument: "Auth"}
	}

	authData, err := ao.Auth.ToAuthBody()
	if err != nil {
		return nil, err
	}

	// wire method names come from authData's keys, not the ToAuthBody
	// return value's internal semantic labels (e.g. "v3password")
	methods := slices.Collect(maps.Keys(authData))
	slices.Sort(methods) // map iteration order is non-deterministic; keep the request stable
	identity := map[string]any{"methods": methods}
	for k, v := range authData {
		identity[k] = v
	}

	body := map[string]any{"identity": identity}
	scopeMap, err := ao.Auth.ToAuthScope()
	if err != nil {
		return nil, err
	}
	if scopeMap != nil {
		body["scope"] = scopeMap
	}

	headers, err := ao.Auth.ToAuthHeaders()
	if err != nil {
		return nil, err
	}
	moreHeaders := make(map[string]string, len(headers))
	for k, v := range headers {
		moreHeaders[k] = fmt.Sprint(v)
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	client := &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{HTTPClient: *httpClient},
		Endpoint:       ao.GetAuthURL(),
	}

	var result gophercloud.Result
	resp, err := client.Post(ctx, client.ServiceURL("auth", "tokens"), map[string]any{"auth": body}, &result.Body, &gophercloud.RequestOpts{
		MoreHeaders: moreHeaders,
		OmitHeaders: []string{"X-Auth-Token"},
	})
	_, result.Header, result.Err = gophercloud.ParseResponse(resp, err)
	if result.Err != nil {
		return nil, result.Err
	}

	var respBody v3TokenBody
	if err := result.ExtractIntoStructPtr(&respBody, "token"); err != nil {
		return nil, err
	}

	return respBody.toAuthResult(result.Header.Get("X-Subject-Token"), ao.Auth.CanReauth()), nil
}

// AuthResult describes a completed Keystone authentication: the issued
// token plus everything needed to use it against a service, independent of
// how the token was acquired (v2 or v3, and by which mechanism).
type AuthResult struct {
	// TokenID is the ID of the issued token. Send this as the
	// "X-Auth-Token" header on requests to Keystone-middleware-protected
	// services.
	TokenID string

	// ExpiresAt is when the token stops being valid.
	ExpiresAt time.Time

	// IssuedAt is when the token was issued. Zero for v2.
	IssuedAt time.Time

	// Methods lists the auth methods the server recorded for this token
	// (e.g. "password", "totp"). Populated from the response, not the
	// request. Empty for v2.
	Methods []AuthType

	// AuditIDs is the chain of audit IDs used for tracing token
	// reissuance. Empty for v2.
	AuditIDs []string

	User                  User
	Project               *Project // nil if domain-scoped, system-scoped, or unscoped
	Domain                *Domain  // nil unless domain-scoped
	System                bool     // true if this token is system-scoped
	Roles                 []Role
	Trust                 *Trust                 // nil unless trust-scoped
	ApplicationCredential *ApplicationCredential // nil unless app-cred auth

	Catalog ServiceCatalog

	CanReauth bool
}

type User struct {
	ID, Name string
	Domain   *Domain // nil for v2
}

type Project struct {
	ID, Name string
	Domain   *Domain // nil for v2
}

type Domain struct {
	ID, Name string
}

type Role struct {
	ID, Name string // ID empty for v2
}

type Trust struct {
	ID            string
	Impersonation bool
	TrusteeUserID string
	TrustorUserID string
}

type ApplicationCredential struct {
	ID          string
	Name        string
	Restricted  bool
	AccessRules []AccessRule
}

type AccessRule struct {
	ID, Path, Method, Service string
}

type ServiceCatalog struct {
	Entries []CatalogEntry
}

type CatalogEntry struct {
	ID, Name, Type string
	Endpoints      []Endpoint
}

type Endpoint struct {
	ID, Region, RegionID, Interface, URL string
}

// Token returns the issued token ID.
func (r AuthResult) Token() string {
	return r.TokenID
}

// AuthenticatedHeaders returns the header(s) needed to authenticate a
// request against a Keystone-middleware-protected service using this
// token.
func (r AuthResult) AuthenticatedHeaders() map[string]string {
	return map[string]string{"X-Auth-Token": r.TokenID}
}

// Expired reports whether the token has already expired.
func (r AuthResult) Expired() bool {
	return r.ExpiresAt.Before(time.Now())
}

// WillExpireBy reports whether the token will expire within d.
func (r AuthResult) WillExpireBy(d time.Duration) bool {
	return r.ExpiresAt.Before(time.Now().Add(d))
}

// Endpoint discovers the endpoint URL for a specific service from this
// AuthResult's catalog. opts must contain enough information to
// unambiguously identify one, and only one, endpoint.
func (r AuthResult) Endpoint(opts gophercloud.EndpointOpts) (string, error) {
	availability := opts.Availability
	if availability == "" {
		availability = gophercloud.AvailabilityPublic
	}

	opts.ApplyDefaults(opts.Type)

	for _, entry := range r.Catalog.Entries {
		if !slices.Contains(opts.Types(), entry.Type) {
			continue
		}
		if opts.Name != "" && entry.Name != opts.Name {
			continue
		}
		for _, endpoint := range entry.Endpoints {
			if gophercloud.Availability(endpoint.Interface) != availability {
				continue
			}
			if opts.Region != "" && endpoint.Region != opts.Region && endpoint.RegionID != opts.Region {
				continue
			}
			return gophercloud.NormalizeURL(endpoint.URL), nil
		}
	}

	return "", &gophercloud.ErrEndpointNotFound{}
}

// EndpointLocator adapts Endpoint to the gophercloud.EndpointLocator
// function type, so it can be assigned directly to a
// gophercloud.ProviderClient's EndpointLocator field.
func (r AuthResult) EndpointLocator() gophercloud.EndpointLocator {
	return func(_ context.Context, opts gophercloud.EndpointOpts) (string, error) {
		return r.Endpoint(opts)
	}
}

// ExtractTokenID lets AuthResult satisfy gophercloud.AuthResult, the same
// way tokens2.CreateResult/tokens3.CreateResult do today.
func (r AuthResult) ExtractTokenID() (string, error) {
	return r.TokenID, nil
}

// v3TokenBody is the private wire shape of a Keystone v3 "token" response
// body, decoded via gophercloud.Result.ExtractIntoStructPtr(&v, "token")
// and mapped onto AuthResult by toAuthResult. It intentionally does not
// import openstack/identity/v3/tokens, per the auth/ layering rule.
type v3TokenBody struct {
	Methods               []string             `json:"methods"`
	ExpiresAt             time.Time            `json:"expires_at"`
	IssuedAt              time.Time            `json:"issued_at"`
	AuditIDs              []string             `json:"audit_ids"`
	User                  v3UserBody           `json:"user"`
	Project               *v3ProjectBody       `json:"project"`
	Domain                *v3DomainBody        `json:"domain"`
	System                map[string]any       `json:"system"`
	Roles                 []v3RoleBody         `json:"roles"`
	Trust                 *v3TrustBody         `json:"OS-TRUST:trust"`
	ApplicationCredential *v3AppCredBody       `json:"application_credential"`
	Catalog               []v3CatalogEntryBody `json:"catalog"`
}

type v3DomainBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type v3UserBody struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Domain *v3DomainBody `json:"domain"`
}

type v3ProjectBody struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Domain *v3DomainBody `json:"domain"`
}

type v3RoleBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type v3TrustUserBody struct {
	ID string `json:"id"`
}

type v3TrustBody struct {
	ID            string          `json:"id"`
	Impersonation bool            `json:"impersonation"`
	TrusteeUserID v3TrustUserBody `json:"trustee_user"`
	TrustorUserID v3TrustUserBody `json:"trustor_user"`
}

type v3AccessRuleBody struct {
	ID      string `json:"id"`
	Path    string `json:"path"`
	Method  string `json:"method"`
	Service string `json:"service"`
}

type v3AppCredBody struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Restricted  bool               `json:"restricted"`
	AccessRules []v3AccessRuleBody `json:"access_rules"`
}

type v3EndpointBody struct {
	ID        string `json:"id"`
	Region    string `json:"region"`
	RegionID  string `json:"region_id"`
	Interface string `json:"interface"`
	URL       string `json:"url"`
}

type v3CatalogEntryBody struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	Endpoints []v3EndpointBody `json:"endpoints"`
}

func v3DomainFromBody(d *v3DomainBody) *Domain {
	if d == nil {
		return nil
	}
	return &Domain{ID: d.ID, Name: d.Name}
}

func (b v3TokenBody) toAuthResult(tokenID string, canReauth bool) *AuthResult {
	result := AuthResult{
		TokenID:   tokenID,
		ExpiresAt: b.ExpiresAt,
		IssuedAt:  b.IssuedAt,
		AuditIDs:  b.AuditIDs,
		User: User{
			ID:     b.User.ID,
			Name:   b.User.Name,
			Domain: v3DomainFromBody(b.User.Domain),
		},
		Domain:    v3DomainFromBody(b.Domain),
		System:    b.System != nil,
		CanReauth: canReauth,
	}

	for _, m := range b.Methods {
		result.Methods = append(result.Methods, AuthType(m))
	}

	if b.Project != nil {
		result.Project = &Project{
			ID:     b.Project.ID,
			Name:   b.Project.Name,
			Domain: v3DomainFromBody(b.Project.Domain),
		}
	}

	for _, r := range b.Roles {
		result.Roles = append(result.Roles, Role{ID: r.ID, Name: r.Name})
	}

	if b.Trust != nil {
		result.Trust = &Trust{
			ID:            b.Trust.ID,
			Impersonation: b.Trust.Impersonation,
			TrusteeUserID: b.Trust.TrusteeUserID.ID,
			TrustorUserID: b.Trust.TrustorUserID.ID,
		}
	}

	if b.ApplicationCredential != nil {
		ac := &ApplicationCredential{
			ID:         b.ApplicationCredential.ID,
			Name:       b.ApplicationCredential.Name,
			Restricted: b.ApplicationCredential.Restricted,
		}
		for _, ar := range b.ApplicationCredential.AccessRules {
			ac.AccessRules = append(ac.AccessRules, AccessRule{
				ID: ar.ID, Path: ar.Path, Method: ar.Method, Service: ar.Service,
			})
		}
		result.ApplicationCredential = ac
	}

	for _, entry := range b.Catalog {
		ce := CatalogEntry{ID: entry.ID, Name: entry.Name, Type: entry.Type}
		for _, ep := range entry.Endpoints {
			ce.Endpoints = append(ce.Endpoints, Endpoint{
				ID: ep.ID, Region: ep.Region, RegionID: ep.RegionID,
				Interface: ep.Interface, URL: ep.URL,
			})
		}
		result.Catalog.Entries = append(result.Catalog.Entries, ce)
	}

	return &result
}

// v2AccessBody is the private wire shape of a Keystone v2 "access" response
// body. It intentionally does not import openstack/identity/v2/tokens, per
// the auth/ layering rule.
type v2AccessBody struct {
	Access struct {
		Token struct {
			ID      string                       `json:"id"`
			Expires gophercloud.JSONRFC3339Milli `json:"expires"`
			Tenant  *v2TenantBody                `json:"tenant"`
		} `json:"token"`
		User struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Roles []struct {
				Name string `json:"name"`
			} `json:"roles"`
		} `json:"user"`
		ServiceCatalog []v2CatalogEntryBody `json:"serviceCatalog"`
	} `json:"access"`
}

type v2TenantBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type v2EndpointBody struct {
	PublicURL   string `json:"publicURL"`
	InternalURL string `json:"internalURL"`
	AdminURL    string `json:"adminURL"`
	Region      string `json:"region"`
}

type v2CatalogEntryBody struct {
	Name      string           `json:"name"`
	Type      string           `json:"type"`
	Endpoints []v2EndpointBody `json:"endpoints"`
}

func (b v2AccessBody) toAuthResult(canReauth bool) *AuthResult {
	result := AuthResult{
		TokenID:   b.Access.Token.ID,
		ExpiresAt: time.Time(b.Access.Token.Expires),
		User: User{
			ID:   b.Access.User.ID,
			Name: b.Access.User.Name,
		},
		CanReauth: canReauth,
	}

	if b.Access.Token.Tenant != nil {
		result.Project = &Project{
			ID:   b.Access.Token.Tenant.ID,
			Name: b.Access.Token.Tenant.Name,
		}
	}

	for _, r := range b.Access.User.Roles {
		result.Roles = append(result.Roles, Role{Name: r.Name})
	}

	for _, entry := range b.Access.ServiceCatalog {
		ce := CatalogEntry{Name: entry.Name, Type: entry.Type}
		for _, ep := range entry.Endpoints {
			for _, urlIface := range []struct {
				url   string
				iface string
			}{
				{ep.PublicURL, "public"},
				{ep.InternalURL, "internal"},
				{ep.AdminURL, "admin"},
			} {
				if urlIface.url == "" {
					continue
				}
				ce.Endpoints = append(ce.Endpoints, Endpoint{
					Region:    ep.Region,
					Interface: urlIface.iface,
					URL:       urlIface.url,
				})
			}
		}
		result.Catalog.Entries = append(result.Catalog.Entries, ce)
	}

	return &result
}
