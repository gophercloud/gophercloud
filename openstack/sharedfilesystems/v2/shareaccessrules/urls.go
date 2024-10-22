package shareaccessrules

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
)

const shareAccessRulesEndpoint = "share-access-rules"

func getURL(c gophercloud.Client, accessID string) string {
	return c.ServiceURL(shareAccessRulesEndpoint, accessID)
}

func listURL(c gophercloud.Client, shareID string) string {
	return fmt.Sprintf("%s?share_id=%s", c.ServiceURL(shareAccessRulesEndpoint), shareID)
}
