package images

import (
	"strings"

	"github.com/rackspace/gophercloud"
)

// `listURL` is a pure function. `listURL(c)` is a URL for which a GET
// request will respond with a list of images in the service `c`.
func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("images")
}

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("images")
}

// `imageURL(c,i)` is the URL for the image identified by ID `i` in
// the service `c`.
func imageURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL("images", imageID)
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
	return c.ServiceURL("images", imageID, "file")
}

var getDataURL = imageDataURL

var updateDataURL = imageDataURL

func imageTagURL(c *gophercloud.ServiceClient, imageID string, tag string) string {
	return c.ServiceURL("images", imageID, "tags", tag)
}

var createTagURL = imageTagURL

var deleteTagURL = imageTagURL

func imageMembersURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL("images", imageID, "members")
}

func reactivateImageURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL("images", imageID, "actions", "reactivate")
}

func deactivateImageURL(c *gophercloud.ServiceClient, imageID string) string {
	return c.ServiceURL("images", imageID, "actions", "deactivate")
}

// builds next page full url based on current url
func nextPageURL(currentURL string, next string) string {
	base := currentURL[:strings.Index(currentURL, "/images")]
	return base + next
}
