package auth

import "github.com/gophercloud/gophercloud/v2"

type Scope struct {
	// Scope to a Domain ID
	DomainID string
	// Scope to a Domain Name
	DomainName string
	// Scope to a Project Domain ID
	ProjectDomainID string
	// Scope to a Project Domain Name
	ProjectDomainName string
	// Scope to a Project ID
	ProjectID string
	// Scope to a Project Name
	ProjectName string
	// Scope for system operations
	System bool
	// ID of the trust to use as a trustee user
	TrustID string
}

func (opts *Scope) ToScopeMap() (map[string]any, error) {
	if opts == nil {
		return nil, nil
	}

	if opts.System {
		return map[string]any{
			"system": map[string]any{
				"all": true,
			},
		}, nil
	}

	if opts.TrustID != "" {
		return map[string]any{
			"OS-TRUST:trust": map[string]string{
				"id": opts.TrustID,
			},
		}, nil
	}

	// Prioritize project ID over project name
	if opts.ProjectID != "" {
		// Domain scoped options cannot be set with project scoped options
		if opts.DomainID != "" || opts.DomainName != "" {
			var domainOption string
			if opts.DomainID != "" {
				domainOption = "DomainID"
			} else {
				domainOption = "DomainName"
			}
			return nil, gophercloud.ErrProjectScopeAndDomainScope{
				DomainOption:  domainOption,
				ProjectOption: "ProjectID",
			}
		}

		// ProjectID
		return map[string]any{
			"project": map[string]any{
				"id": &opts.ProjectID,
			},
		}, nil

	} else if opts.ProjectName != "" {
		// Domain scoped options cannot be set with project scoped options
		if opts.DomainID != "" || opts.DomainName != "" {
			var domainOption string
			if opts.DomainID != "" {
				domainOption = "DomainID"
			} else {
				domainOption = "DomainName"
			}
			return nil, gophercloud.ErrProjectScopeAndDomainScope{
				DomainOption:  domainOption,
				ProjectOption: "ProjectName",
			}
		}

		if opts.ProjectDomainID == "" && opts.ProjectDomainName == "" {
			return nil, gophercloud.ErrScopeProjectDomainIDOrProjectDomainName{}
		}

		// Prioritize ProjectDomainID
		var k, v string
		if opts.ProjectDomainID != "" {
			k = "id"
			v = opts.ProjectDomainID
		} else if opts.ProjectDomainName != "" {
			k = "name"
			v = opts.ProjectDomainName
		}

		return map[string]any{
			"project": map[string]any{
				"name":   &opts.ProjectName,
				"domain": map[string]any{k: &v},
			},
		}, nil

	} else if opts.DomainID != "" {
		// DomainID
		return map[string]any{
			"domain": map[string]any{
				"id": &opts.DomainID,
			},
		}, nil

	} else if opts.DomainName != "" {
		// DomainName
		return map[string]any{
			"domain": map[string]any{
				"name": &opts.DomainName,
			},
		}, nil
	}

	return nil, nil
}
