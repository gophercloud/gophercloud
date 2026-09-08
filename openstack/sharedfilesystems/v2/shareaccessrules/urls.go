package shareaccessrules

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
)

const shareAccessRulesEndpoint = "share-access-rules"

func getURL(c *gophercloud.ServiceClient, accessID string) string {
	return c.ServiceURL(shareAccessRulesEndpoint, accessID)
}

func listURL(c *gophercloud.ServiceClient, shareID string) string {
	return fmt.Sprintf("%s?share_id=%s", c.ServiceURL(shareAccessRulesEndpoint), shareID)
}

func metadataURL(c *gophercloud.ServiceClient, accessID string) string {
	return c.ServiceURL(shareAccessRulesEndpoint, accessID, "metadata")
}

func metadatumURL(c *gophercloud.ServiceClient, accessID, key string) string {
	return c.ServiceURL(shareAccessRulesEndpoint, accessID, "metadata", key)
}
