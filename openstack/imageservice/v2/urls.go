package v2

import "github.com/rackspace/gophercloud"

// `listURL` is a pure function. `listURL(c)` is a URL for which a GET
// request will respond with a list of images in the service `c`.
func listURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}

func createURL(c *gophercloud.ServiceClient) string {
	return c.Endpoint
}

// `imageURL(c,i)` is the URL for the image identified by ID `i` in
// the service `c`.
func imageURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL(imageID)
}

// `getURL(c,i)` is a URL for which a GET request will respond with
// information about the image identified by ID `i` in the service
// `c`.
var getURL = imageURL

var updateURL = imageURL

var deleteURL = imageURL

// `imageDataURL(c,i)` is the URL for the binary image data for the
// image identified by ID `i` in the service `c`.
func imageDataURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL(imageID, "file")
}

var getDataURL = imageDataURL

var updateDataURL = imageDataURL

func imageTagURL(c *gophercloud.ServiceClient, imageID string, tag string) string {
	return c.ServiceURL(imageID, "tags", tag)
}

var createTagURL = imageTagURL

var deleteTagURL = imageTagURL

func imageMembersURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL(imageID, "members")
}

var listMembersURL = imageMembersURL

var createMemberURL = imageMembersURL

func imageMemberURL(c *gophercloud.ServiceClient, imageID string, memberID string) string {
	return c.ServiceURL(imageID, "members", memberID)
}

var getMemberURL = imageMemberURL

var updateMemberURL = imageMemberURL

var deleteMemberURL = imageMemberURL

func reactivateImageURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL(imageID, "actions", "reactivate")
}

func deactivateImageURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL(imageID, "actions", "deactivate")
}
