package runbooks

import "github.com/gophercloud/gophercloud/v2"

func rootURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("runbooks")
}

func resourceURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("runbooks", id)
}

func traitsURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("runbooks", id, "traits")
}

func traitURL(client *gophercloud.ServiceClient, id string, trait string) string {
	return client.ServiceURL("runbooks", id, "traits", trait)
}
