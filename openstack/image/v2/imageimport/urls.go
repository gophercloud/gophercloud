package imageimport

import "github.com/gophercloud/gophercloud/v2"

const (
	rootPath     = "images"
	infoPath     = "info"
	resourcePath = "import"
)

func infoURL(c gophercloud.Client) string {
	return c.ServiceURL(infoPath, resourcePath)
}

func importURL(c gophercloud.Client, imageID string) string {
	return c.ServiceURL(rootPath, imageID, resourcePath)
}
