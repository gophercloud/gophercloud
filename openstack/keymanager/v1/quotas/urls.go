package quotas

import "github.com/gophercloud/gophercloud"

const resourcePathProject = "project-quotas"
const resourcePath = "quotas"

func resourceURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePath)
}

func getURL(c *gophercloud.ServiceClient) string {
	return resourceURL(c)
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL(resourcePathProject)
}

func resourceProjectURL(c *gophercloud.ServiceClient, projectID string) string {
	return c.ServiceURL(resourcePathProject, projectID)
}

func getProjectURL(c *gophercloud.ServiceClient, projectID string) string {
	return resourceProjectURL(c, projectID)
}

func updateProjectURL(c *gophercloud.ServiceClient, projectID string) string {
	return resourceProjectURL(c, projectID)
}

func deleteProjectURL(c *gophercloud.ServiceClient, projectID string) string {
	return resourceProjectURL(c, projectID)
}
