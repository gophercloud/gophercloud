//go:build acceptance || compute || remoteconsoles

package v2

import (
	"net/url"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestRemoteConsoleCreate(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	client.Microversion = "2.6"

	server, err := CreateServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	remoteConsole, err := CreateRemoteConsole(t, client, server.ID)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, remoteConsole)
}

func TestConsoleGet(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	client.Microversion = "2.96"

	server, err := CreateMicroversionServer(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, server)

	remoteConsole, err := CreateRemoteConsole(t, client, server.ID)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, remoteConsole)

	u, err := url.Parse(remoteConsole.URL)
	th.AssertNoErr(t, err)

	encodedPath := u.Query().Get("path")
	if encodedPath == "" {
		t.Fatalf("console URL missing path param: %s", remoteConsole.URL)
	}

	decodedPath, err := url.QueryUnescape(encodedPath)
	th.AssertNoErr(t, err)

	inner, err := url.Parse(decodedPath)
	th.AssertNoErr(t, err)

	tokenID := inner.Query().Get("token")
	if tokenID == "" {
		t.Fatalf("failed to extract token from console URL: %s", remoteConsole.URL)
	}

	console, err := GetConsole(t, client, tokenID)
	if err != nil {
		t.Fatalf("failed to get console: %v", err)
	}
	tools.PrintResource(t, console)
}
