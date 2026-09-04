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

	if opts.ProjectName != "" {
		// ProjectName provided: exactly one of ProjectDomainID, ProjectDomainName, DomainID
		// or DomainName must also be supplied.
		if opts.ProjectDomainID == "" && opts.ProjectDomainName == "" {
			if opts.DomainID == "" && opts.DomainName == "" {
				return nil, gophercloud.ErrScopeDomainIDOrDomainName{}
			}
		}

		// ProjectID may not be supplied.
		if opts.ProjectID != "" {
			return nil, gophercloud.ErrScopeProjectIDOrProjectName{}
		}

		// Prioritize ProjectDomainID and ProjectDomainName before falling back on
		// DomainID or DomainName
		if opts.ProjectDomainID != "" {
			// ProjectName + ProjectDomainID
			return map[string]any{
				"project": map[string]any{
					"name":   &opts.ProjectName,
					"domain": map[string]any{"id": &opts.ProjectDomainID},
				},
			}, nil
		}

		if opts.ProjectDomainName != "" {
			// ProjectName + ProjectDomainName
			return map[string]any{
				"project": map[string]any{
					"name":   &opts.ProjectName,
					"domain": map[string]any{"name": &opts.ProjectDomainName},
				},
			}, nil
		}

		if opts.DomainID != "" {
			// ProjectName + DomainID
			return map[string]any{
				"project": map[string]any{
					"name":   &opts.ProjectName,
					"domain": map[string]any{"id": &opts.DomainID},
				},
			}, nil
		}

		if opts.DomainName != "" {
			// ProjectName + DomainName
			return map[string]any{
				"project": map[string]any{
					"name":   &opts.ProjectName,
					"domain": map[string]any{"name": &opts.DomainName},
				},
			}, nil
		}
	} else if opts.ProjectID != "" {
		// ProjectID provided. ProjectName, DomainID, DomainName,
		// ProjectDomainID, and ProjectDomainName may not be provided.
		if opts.DomainID != "" {
			return nil, gophercloud.ErrScopeProjectIDAlone{Reason: "DomainID"}
		}
		if opts.DomainName != "" {
			return nil, gophercloud.ErrScopeProjectIDAlone{Reason: "DomainName"}
		}
		if opts.ProjectDomainID != "" {
			return nil, gophercloud.ErrScopeProjectIDAlone{Reason: "ProjectDomainID"}
		}
		if opts.ProjectDomainName != "" {
			return nil, gophercloud.ErrScopeProjectIDAlone{Reason: "ProjectDomainName"}
		}

		// ProjectID
		return map[string]any{
			"project": map[string]any{
				"id": &opts.ProjectID,
			},
		}, nil
	} else if opts.DomainID != "" {
		// DomainID provided. ProjectID, ProjectName, and DomainName may not be provided.
		if opts.DomainName != "" {
			return nil, gophercloud.ErrScopeDomainIDOrDomainName{}
		}

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
