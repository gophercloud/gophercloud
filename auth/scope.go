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
		if opts.ProjectID != "" || opts.ProjectName != "" || opts.DomainID != "" || opts.DomainName != "" || opts.TrustID != "" {
			return nil, gophercloud.ErrSystemScopeAndOtherScope{}
		}

		return map[string]any{
			"system": map[string]any{
				"all": true,
			},
		}, nil
	}

	if opts.TrustID != "" {
		if opts.ProjectID != "" || opts.ProjectName != "" || opts.DomainID != "" || opts.DomainName != "" || opts.System {
			return nil, gophercloud.ErrTrustScopeAndOtherScope{}
		}

		return map[string]any{
			"OS-TRUST:trust": map[string]string{
				"id": opts.TrustID,
			},
		}, nil
	}

	if opts.ProjectID != "" || opts.ProjectName != "" {
		// Domain scoped options cannot be set with project scoped options
		if opts.DomainID != "" || opts.DomainName != "" {
			var projectOption string
			if opts.ProjectID != "" {
				projectOption = "ProjectID"
			} else {
				projectOption = "ProjectName"
			}

			var domainOption string
			if opts.DomainID != "" {
				domainOption = "DomainID"
			} else {
				domainOption = "DomainName"
			}

			return nil, gophercloud.ErrProjectScopeAndDomainScope{
				DomainOption:  domainOption,
				ProjectOption: projectOption,
			}
		}

		// Prioritize ProjectID over ProjectName
		if opts.ProjectID != "" {
			return map[string]any{
				"project": map[string]any{
					"id": &opts.ProjectID,
				},
			}, nil
		}

		if opts.ProjectDomainID == "" && opts.ProjectDomainName == "" {
			return nil, gophercloud.ErrScopeProjectDomainIDOrProjectDomainName{}
		}

		// Prioritize ProjectDomainID over ProjectDomainName
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

	} else if opts.DomainID != "" || opts.DomainName != "" {
		if opts.ProjectDomainID != "" || opts.ProjectDomainName != "" {
			var projectOption string
			if opts.ProjectDomainID != "" {
				projectOption = "ProjectDomainID"
			} else {
				projectOption = "ProjectDomainName"
			}

			var domainOption string
			if opts.DomainID != "" {
				domainOption = "DomainID"
			} else {
				domainOption = "DomainName"
			}

			return nil, gophercloud.ErrProjectScopeAndDomainScope{
				DomainOption:  domainOption,
				ProjectOption: projectOption,
			}
		}

		// Prioritize DomainID over DomainName
		if opts.DomainID != "" {
			return map[string]any{
				"domain": map[string]any{
					"id": &opts.DomainID,
				},
			}, nil
		}

		// Fallback to DomainName
		return map[string]any{
			"domain": map[string]any{
				"name": &opts.DomainName,
			},
		}, nil

	}

	return nil, nil
}
