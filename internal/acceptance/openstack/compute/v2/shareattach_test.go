//go:build acceptance || compute || shareattach

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	sfs "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/sharedfilesystems/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestShareAttach(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	shareClient, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	server, err := CreateMicroversionServer(t, client)
	th.AssertNoErr(t, err)

	share, err := sfs.CreateShare(t, shareClient)
	th.AssertNoErr(t, err)
	defer sfs.DeleteShare(t, shareClient, share)

	defer DeleteServer(t, client, server)

	client.Microversion = "2.97"

	err = servers.Stop(context.TODO(), client, server.ID).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForComputeStatus(client, server, "SHUTOFF")
	th.AssertNoErr(t, err)

	shareAttachment, err := CreateShareAttachment(t, client, server, share)
	th.AssertNoErr(t, err)
	defer DeleteShareAttachment(t, client, server, shareAttachment)

	tools.PrintResource(t, shareAttachment)

	th.AssertEquals(t, share.ID, shareAttachment.ShareID)
}
