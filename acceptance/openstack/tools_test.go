// +build acceptance

package openstack

import (
	"crypto/rand"
	"fmt"
	"github.com/rackspace/gophercloud/openstack/compute/servers"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
	"os"
	"text/tabwriter"
	"time"
)

var errTimeout = fmt.Errorf("Timeout.")

type testState struct {
	o             identity.AuthOptions
	a             identity.AuthResults
	sc            *identity.ServiceCatalog
	eps           []identity.Endpoint
	w             *tabwriter.Writer
	imageId       string
	flavorId      string
	region        string
	ep            string
	client        *servers.Client
	createdServer *servers.Server
	gottenServer  *servers.Server
	updatedServer *servers.Server
	serverName    string
	alternateName string
	flavorIdResize string
}

func setupForList(service string) (*testState, error) {
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

	ts.eps, err = findAllEndpoints(ts.sc, service)
	if err != nil {
		return ts, err
	}

	ts.w = new(tabwriter.Writer)
	ts.w.Init(os.Stdout, 2, 8, 2, ' ', 0)

	return ts, nil
}

func setupForCRUD() (*testState, error) {
	ts, err := setupForList("compute")
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

	ts.ep, err = findEndpointForRegion(ts.eps, ts.region)
	if err != nil {
		return ts, err
	}

	return ts, err
}

func findAllEndpoints(sc *identity.ServiceCatalog, service string) ([]identity.Endpoint, error) {
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

func findEndpointForRegion(eps []identity.Endpoint, r string) (string, error) {
	for _, ep := range eps {
		if ep.Region == r {
			return ep.PublicURL, nil
		}
	}
	return "", fmt.Errorf("Unknown region %s", r)
}

func countDown(ts *testState, timeout int) (bool, int, error) {
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

func createServer(ts *testState) error {
	ts.serverName = randomString("ACPTTEST", 16)
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

func waitForStatus(ts *testState, s string) error {
	var (
		inProgress bool
		timeout    int
		err        error
	)

	for inProgress, timeout, err = countDown(ts, 300); inProgress; inProgress, timeout, err = countDown(ts, timeout) {
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

func changeServerName(ts *testState) error {
	var (
		inProgress bool
		timeout    int
	)

	ts.alternateName = randomString("ACPTTEST", 16)
	for ts.alternateName == ts.serverName {
		ts.alternateName = randomString("ACPTTEST", 16)
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

	for inProgress, timeout, err = countDown(ts, 300); inProgress; inProgress, timeout, err = countDown(ts, timeout) {
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

func makeNewPassword(oldPass string) string {
	fmt.Println("Current password: "+oldPass)
	randomPassword := randomString("", 16)
	for randomPassword == oldPass {
		randomPassword = randomString("", 16)
	}
	fmt.Println("    New password: "+randomPassword)
	return randomPassword
}

func changeAdminPassword(ts *testState) error {
	randomPassword := makeNewPassword(ts.createdServer.AdminPass)
	
	err := servers.ChangeAdminPassword(ts.client, ts.createdServer.Id, randomPassword)
	if err != nil {
		return err
	}
	
	err = waitForStatus(ts, "PASSWORD")
	if err != nil {
		return err
	}
	
	return waitForStatus(ts, "ACTIVE")
}

func rebootServer(ts *testState) error {
	fmt.Println("Attempting reboot of server "+ts.createdServer.Id)
	err := servers.Reboot(ts.client, ts.createdServer.Id, servers.OSReboot)
	if err != nil {
		return err
	}
	
	err = waitForStatus(ts, "REBOOT")
	if err != nil {
		return err
	}
	
	return waitForStatus(ts, "ACTIVE")
}

func rebuildServer(ts *testState) error {
	fmt.Println("Attempting to rebuild server "+ts.createdServer.Id)

	newPassword := makeNewPassword(ts.createdServer.AdminPass)
	newName := randomString("ACPTTEST", 16)
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

	err = waitForStatus(ts, "REBUILD")
	if err != nil {
		return err
	}
	
	return waitForStatus(ts, "ACTIVE")
}

func resizeServer(ts *testState) error {
	fmt.Println("Attempting to resize server "+ts.createdServer.Id)

	err := servers.Resize(ts.client, ts.createdServer.Id, ts.flavorIdResize)
	if err != nil {
		return err
	}

	err = waitForStatus(ts, "RESIZE")
	if err != nil {
		return err
	}

	return waitForStatus(ts, "VERIFY_RESIZE")
}

func confirmResize(ts *testState) error {
	fmt.Println("Attempting to confirm resize for server "+ts.createdServer.Id)
	
	err := servers.ConfirmResize(ts.client, ts.createdServer.Id)
	if err != nil {
		return err
	}
	
	return waitForStatus(ts, "ACTIVE")
}

func revertResize(ts *testState) error {
	fmt.Println("Attempting to revert resize for server "+ts.createdServer.Id)
	
	err := servers.RevertResize(ts.client, ts.createdServer.Id)
	if err != nil {
		return err
	}

	err = waitForStatus(ts, "REVERT_RESIZE")
	if err != nil {
		return err
	}

	return waitForStatus(ts, "ACTIVE")
}

// randomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
// (Implementation from Even Shaw's contribution on
// http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go).
func randomString(prefix string, n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return prefix + string(bytes)
}
