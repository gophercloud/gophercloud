package federation

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type UserType string

const (
	UserTypeEphemeral UserType = "ephemeral"
	UserTypeLocal     UserType = "local"
)

// Mapping a set of rules to map federation protocol attributes to
// Identity API objects.
type Mapping struct {
	// The Federation Mapping unique ID
	ID string `json:"id"`

	// Links contains referencing links to the limit.
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
func (c MappingsPage) NextPageURL() (string, error) {
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
