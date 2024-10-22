package attributestags

import "github.com/gophercloud/gophercloud/v2"

const (
	tagsPath = "tags"
)

func replaceURL(c gophercloud.Client, r_type string, id string) string {
	return c.ServiceURL(r_type, id, tagsPath)
}

func listURL(c gophercloud.Client, r_type string, id string) string {
	return c.ServiceURL(r_type, id, tagsPath)
}

func deleteAllURL(c gophercloud.Client, r_type string, id string) string {
	return c.ServiceURL(r_type, id, tagsPath)
}

func addURL(c gophercloud.Client, r_type string, id string, tag string) string {
	return c.ServiceURL(r_type, id, tagsPath, tag)
}

func deleteURL(c gophercloud.Client, r_type string, id string, tag string) string {
	return c.ServiceURL(r_type, id, tagsPath, tag)
}

func confirmURL(c gophercloud.Client, r_type string, id string, tag string) string {
	return c.ServiceURL(r_type, id, tagsPath, tag)
}
