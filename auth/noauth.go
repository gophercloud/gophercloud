package auth

import (
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

// NoAuthOpts builds a synthetic token for a service configured without an
// identity service.
type NoAuthOpts interface {
	ToNoAuthToken() (string, error)
}

// NoAuthV2Opts contains options for Cinder API v2 no-auth compatibility.
type NoAuthV2Opts struct {
	// Username is the synthetic username.
	Username string

	// TenantName is the synthetic tenant name.
	TenantName string
}

// ToNoAuthToken builds the synthetic username:tenant token.
func (opts NoAuthV2Opts) ToNoAuthToken() (string, error) {
	username := opts.Username
	if username == "" {
		username = "admin"
	}
	tenantName := opts.TenantName
	if tenantName == "" {
		tenantName = "admin"
	}

	if strings.Contains(username, ":") {
		return "", gophercloud.ErrInvalidInput{
			ErrMissingInput: gophercloud.ErrMissingInput{Argument: "Username"},
			Value:           username,
		}
	}
	if strings.Contains(tenantName, ":") {
		return "", gophercloud.ErrInvalidInput{
			ErrMissingInput: gophercloud.ErrMissingInput{Argument: "TenantName"},
			Value:           tenantName,
		}
	}

	return username + ":" + tenantName, nil
}

// NoAuthV3Opts contains options for Cinder API v3 no-auth compatibility.
type NoAuthV3Opts struct {
	// UserID is the synthetic user ID.
	UserID string

	// ProjectID is the synthetic project ID. It defaults to UserID.
	ProjectID string
}

// ToNoAuthToken builds the synthetic user ID:project ID token.
func (opts NoAuthV3Opts) ToNoAuthToken() (string, error) {
	if opts.UserID == "" {
		return "", gophercloud.ErrMissingInput{Argument: "UserID"}
	}
	projectID := opts.ProjectID
	if projectID == "" {
		projectID = opts.UserID
	}

	if strings.Contains(opts.UserID, ":") {
		return "", gophercloud.ErrInvalidInput{
			ErrMissingInput: gophercloud.ErrMissingInput{Argument: "UserID"},
			Value:           opts.UserID,
		}
	}
	if strings.Contains(projectID, ":") {
		return "", gophercloud.ErrInvalidInput{
			ErrMissingInput: gophercloud.ErrMissingInput{Argument: "ProjectID"},
			Value:           projectID,
		}
	}

	return opts.UserID + ":" + projectID, nil
}
