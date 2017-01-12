package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

// CreateProject will create a project with a random name.
// It takes an optional createOpts parameter since creating a project
// has so many options. An error will be returned if the project was
// unable to be created.
func CreateProject(t *testing.T, client *gophercloud.ServiceClient, c *projects.CreateOpts) (*projects.Project, error) {
	name := tools.RandomString("ACPTTEST", 8)
	t.Logf("Attempting to create project: %s", name)

	var createOpts projects.CreateOpts
	if c != nil {
		createOpts = *c
	} else {
		createOpts = projects.CreateOpts{}
	}

	createOpts.Name = name

	project, err := projects.Create(client, createOpts).Extract()
	if err != nil {
		return project, err
	}

	t.Logf("Successfully created project %s with ID %s", name, project.ID)

	return project, nil
}

// DeleteProject will delete a project by ID. A fatal error will occur if
// the project ID failed to be deleted. This works best when using it as
// a deferred function.
func DeleteProject(t *testing.T, client *gophercloud.ServiceClient, projectID string) {
	err := projects.Delete(client, projectID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete project %s: %v", projectID, err)
	}

	t.Logf("Deleted project: %s", projectID)
}

// PrintEndpoint will print an endpoint and all of its attributes.
func PrintEndpoint(t *testing.T, endpoint *endpoints.Endpoint) {
	t.Logf("ID: %s", endpoint.ID)
	t.Logf("Availability: %s", endpoint.Availability)
	t.Logf("Name: %s", endpoint.Name)
	t.Logf("Region: %s", endpoint.Region)
	t.Logf("ServiceID: %s", endpoint.ServiceID)
	t.Logf("URL: %s", endpoint.URL)
}

// PrintProject will print a project and all of its attributes.
func PrintProject(t *testing.T, project *projects.Project) {
	t.Logf("ID: %s", project.ID)
	t.Logf("IsDomain: %t", project.IsDomain)
	t.Logf("Description: %s", project.Description)
	t.Logf("DomainID: %s", project.DomainID)
	t.Logf("Enabled: %t", project.Enabled)
	t.Logf("Name: %s", project.Name)
	t.Logf("ParentID: %s", project.ParentID)
}

// PrintService will print a service and all of its attributes.
func PrintService(t *testing.T, service *services.Service) {
	t.Logf("ID: %s", service.ID)
	t.Logf("Description: %s", service.Description)
	t.Logf("Name: %s", service.Name)
	t.Logf("Type: %s", service.Type)
}

// PrintToken will print a token and all of its attributes.
func PrintToken(t *testing.T, token *tokens.Token) {
	t.Logf("ID: %s", token.ID)
	t.Logf("ExpiresAt: %v", token.ExpiresAt)
}
