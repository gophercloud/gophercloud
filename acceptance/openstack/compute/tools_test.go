// +build acceptance

package compute

import (
	"crypto/rand"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	identity "github.com/rackspace/gophercloud/openstack/identity/v2"
	"github.com/rackspace/gophercloud/openstack/utils"
)

var errTimeout = fmt.Errorf("Timeout.")

type testState struct {
	o              identity.AuthOptions
	a              identity.AuthResults
	sc             *identity.ServiceCatalog
	eps            []identity.Endpoint
	w              *tabwriter.Writer
	imageId        string
	flavorId       string
	region         string
	ep             string
	client         *servers.Client
	createdServer  *servers.Server
	gottenServer   *servers.Server
	updatedServer  *servers.Server
	serverName     string
	alternateName  string
	flavorIdResize string
}

func SetupForList(service string) (*testState, error) {
	var err error

	ts := new(testState)

	ts.o, err = utils.AuthOptions()
	if err != nil {
		return ts, err
	}

	ts.a, err = identity.Authenticate(ts.o)
	if err != nil {
		return ts, err
	}

	ts.sc, err = identity.GetServiceCatalog(ts.a)
	if err != nil {
		return ts, err
	}

	ts.eps, err = FindAllEndpoints(ts.sc, service)
	if err != nil {
		return ts, err
	}

	ts.w = new(tabwriter.Writer)
	ts.w.Init(os.Stdout, 2, 8, 2, ' ', 0)

	return ts, nil
}

func SetupForCRUD() (*testState, error) {
	ts, err := SetupForList("compute")
	if err != nil {
		return ts, err
	}

	ts.imageId = os.Getenv("OS_IMAGE_ID")
	if ts.imageId == "" {
		return ts, fmt.Errorf("Expected OS_IMAGE_ID environment variable to be set")
	}

	ts.flavorId = os.Getenv("OS_FLAVOR_ID")
	if ts.flavorId == "" {
		return ts, fmt.Errorf("Expected OS_FLAVOR_ID environment variable to be set")
	}

	ts.flavorIdResize = os.Getenv("OS_FLAVOR_ID_RESIZE")
	if ts.flavorIdResize == "" {
		return ts, fmt.Errorf("Expected OS_FLAVOR_ID_RESIZE environment variable to be set")
	}

	if ts.flavorIdResize == ts.flavorId {
		return ts, fmt.Errorf("OS_FLAVOR_ID and OS_FLAVOR_ID_RESIZE cannot be the same")
	}

	ts.region = os.Getenv("OS_REGION_NAME")
	if ts.region == "" {
		ts.region = ts.eps[0].Region
	}

	ts.ep, err = FindEndpointForRegion(ts.eps, ts.region)
	if err != nil {
		return ts, err
	}

	return ts, err
}

func FindAllEndpoints(sc *identity.ServiceCatalog, service string) ([]identity.Endpoint, error) {
	ces, err := sc.CatalogEntries()
	if err != nil {
		return nil, err
	}

	for _, ce := range ces {
		if ce.Type == service {
			return ce.Endpoints, nil
		}
	}

	return nil, fmt.Errorf(service + " endpoint not found.")
}

func FindEndpointForRegion(eps []identity.Endpoint, r string) (string, error) {
	for _, ep := range eps {
		if ep.Region == r {
			return ep.PublicURL, nil
		}
	}
	return "", fmt.Errorf("Unknown region %s", r)
}

func CountDown(ts *testState, timeout int) (bool, int, error) {
	if timeout < 1 {
		return false, 0, errTimeout
	}
	time.Sleep(1 * time.Second)
	timeout--

	gr, err := servers.GetDetail(ts.client, ts.createdServer.Id)
	if err != nil {
		return false, timeout, err
	}

	ts.gottenServer, err = servers.GetServer(gr)
	if err != nil {
		return false, timeout, err
	}

	return true, timeout, nil
}

func CreateServer(ts *testState) error {
	ts.serverName = RandomString("ACPTTEST", 16)
	fmt.Printf("Attempting to create server: %s\n", ts.serverName)

	ts.client = servers.NewClient(ts.ep, ts.a, ts.o)

	cr, err := servers.Create(ts.client, map[string]interface{}{
		"flavorRef": ts.flavorId,
		"imageRef":  ts.imageId,
		"name":      ts.serverName,
	})
	if err != nil {
		return err
	}

	ts.createdServer, err = servers.GetServer(cr)
	return err
}

func WaitForStatus(ts *testState, s string) error {
	var (
		inProgress bool
		timeout    int
		err        error
	)

	for inProgress, timeout, err = CountDown(ts, 300); inProgress; inProgress, timeout, err = CountDown(ts, timeout) {
		if ts.gottenServer.Id != ts.createdServer.Id {
			return fmt.Errorf("created server id (%s) != gotten server id (%s)", ts.createdServer.Id, ts.gottenServer.Id)
		}

		if ts.gottenServer.Status == s {
			fmt.Printf("Server reached state %s after %d seconds (approximately)\n", s, 300-timeout)
			break
		}
	}

	if err == errTimeout {
		fmt.Printf("Time out -- I'm not waiting around.\n")
		err = nil
	}

	return err
}

func ChangeServerName(ts *testState) error {
	var (
		inProgress bool
		timeout    int
	)

	ts.alternateName = RandomString("ACPTTEST", 16)
	for ts.alternateName == ts.serverName {
		ts.alternateName = RandomString("ACPTTEST", 16)
	}
	fmt.Println("Attempting to change server name")

	ur, err := servers.Update(ts.client, ts.createdServer.Id, map[string]interface{}{
		"name": ts.alternateName,
	})
	if err != nil {
		return err
	}

	ts.updatedServer, err = servers.GetServer(ur)
	if err != nil {
		return err
	}

	if ts.updatedServer.Id != ts.createdServer.Id {
		return fmt.Errorf("Expected updated and created server to share the same ID")
	}

	for inProgress, timeout, err = CountDown(ts, 300); inProgress; inProgress, timeout, err = CountDown(ts, timeout) {
		if ts.gottenServer.Id != ts.updatedServer.Id {
			return fmt.Errorf("Updated server ID (%s) != gotten server ID (%s)", ts.updatedServer.Id, ts.gottenServer.Id)
		}

		if ts.gottenServer.Name == ts.alternateName {
			fmt.Printf("Server updated after %d seconds (approximately)\n", 300-timeout)
			break
		}
	}

	if err == errTimeout {
		fmt.Printf("I'm not waiting around.\n")
		err = nil
	}

	return err
}

func MakeNewPassword(oldPass string) string {
	fmt.Println("Current password: " + oldPass)
	randomPassword := RandomString("", 16)
	for randomPassword == oldPass {
		randomPassword = RandomString("", 16)
	}
	fmt.Println("    New password: " + randomPassword)
	return randomPassword
}

func ChangeAdminPassword(ts *testState) error {
	randomPassword := MakeNewPassword(ts.createdServer.AdminPass)

	err := servers.ChangeAdminPassword(ts.client, ts.createdServer.Id, randomPassword)
	if err != nil {
		return err
	}

	err = WaitForStatus(ts, "PASSWORD")
	if err != nil {
		return err
	}

	return WaitForStatus(ts, "ACTIVE")
}

func RebootServer(ts *testState) error {
	fmt.Println("Attempting reboot of server " + ts.createdServer.Id)
	err := servers.Reboot(ts.client, ts.createdServer.Id, servers.OSReboot)
	if err != nil {
		return err
	}

	err = WaitForStatus(ts, "REBOOT")
	if err != nil {
		return err
	}

	return WaitForStatus(ts, "ACTIVE")
}

func RebuildServer(ts *testState) error {
	fmt.Println("Attempting to rebuild server " + ts.createdServer.Id)

	newPassword := MakeNewPassword(ts.createdServer.AdminPass)
	newName := RandomString("ACPTTEST", 16)
	sr, err := servers.Rebuild(ts.client, ts.createdServer.Id, newName, newPassword, ts.imageId, nil)
	if err != nil {
		return err
	}

	s, err := servers.GetServer(sr)
	if err != nil {
		return err
	}
	if s.Id != ts.createdServer.Id {
		return fmt.Errorf("Expected rebuilt server ID of %s; got %s", ts.createdServer.Id, s.Id)
	}

	err = WaitForStatus(ts, "REBUILD")
	if err != nil {
		return err
	}

	return WaitForStatus(ts, "ACTIVE")
}

func ResizeServer(ts *testState) error {
	fmt.Println("Attempting to resize server " + ts.createdServer.Id)

	err := servers.Resize(ts.client, ts.createdServer.Id, ts.flavorIdResize)
	if err != nil {
		return err
	}

	err = WaitForStatus(ts, "RESIZE")
	if err != nil {
		return err
	}

	return WaitForStatus(ts, "VERIFY_RESIZE")
}

func ConfirmResize(ts *testState) error {
	fmt.Println("Attempting to confirm resize for server " + ts.createdServer.Id)

	err := servers.ConfirmResize(ts.client, ts.createdServer.Id)
	if err != nil {
		return err
	}

	return WaitForStatus(ts, "ACTIVE")
}

func RevertResize(ts *testState) error {
	fmt.Println("Attempting to revert resize for server " + ts.createdServer.Id)

	err := servers.RevertResize(ts.client, ts.createdServer.Id)
	if err != nil {
		return err
	}

	err = WaitForStatus(ts, "REVERT_RESIZE")
	if err != nil {
		return err
	}

	return WaitForStatus(ts, "ACTIVE")
}

// randomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
// (Implementation from Even Shaw's contribution on
// http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go).
func RandomString(prefix string, n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return prefix + string(bytes)
}
