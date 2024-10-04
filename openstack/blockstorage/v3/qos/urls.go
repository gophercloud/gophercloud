package qos

import "github.com/gophercloud/gophercloud/v2"

func getURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("qos-specs", id)
}

func createURL(c gophercloud.Client) string {
	return c.ServiceURL("qos-specs")
}

func listURL(c gophercloud.Client) string {
	return c.ServiceURL("qos-specs")
}

func deleteURL(c gophercloud.Client, id string) string {
	return c.ServiceURL("qos-specs", id)
}

func updateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("qos-specs", id)
}

func deleteKeysURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("qos-specs", id, "delete_keys")
}

func associateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("qos-specs", id, "associate")
}

func disassociateURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("qos-specs", id, "disassociate")
}

func disassociateAllURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("qos-specs", id, "disassociate_all")
}

func listAssociationsURL(client gophercloud.Client, id string) string {
	return client.ServiceURL("qos-specs", id, "associations")
}
