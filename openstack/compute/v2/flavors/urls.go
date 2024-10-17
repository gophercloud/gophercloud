package flavors

import (
	"github.com/gophercloud/gophercloud/v2"
)

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("flavors", id)
}

func listURL(client gophercloud.Client) string {
	return client.ServiceURL("flavors", "detail")
}

func createURL(client gophercloud.Client) string {
	return client.ServiceURL("flavors")
}

func updateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("flavors", id)
}

func deleteURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("flavors", id)
}

func accessURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("flavors", id, "os-flavor-access")
}

func accessActionURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("flavors", id, "action")
}

func extraSpecsListURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs")
}

func extraSpecsGetURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}

func extraSpecsCreateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs")
}

func extraSpecUpdateURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}

func extraSpecDeleteURL(client gophercloud.Client, id, key string) string {
	return client.ServiceURL("flavors", id, "os-extra_specs", key)
}
