// +build acceptance

package tools

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
	O             identity.AuthOptions
	A             identity.AuthResults
	SC            *identity.ServiceCatalog
	EPs           []identity.Endpoint
	W             *tabwriter.Writer
	ImageId       string
	FlavorId      string
	Region        string
	EP            string
	Client        *servers.Client
	CreatedServer *servers.Server
	GottenServer  *servers.Server
	UpdatedServer *servers.Server
	ServerName    string
	AlternateName string
	FlavorIdResize string
}

func SetupForList(service string) (*testState, error) {
	var err error

	ts := new(testState)

	ts.O, err = utils.AuthOptions()
	if err != nil {
		return ts, err
	}

	ts.A, err = identity.Authenticate(ts.O)
	if err != nil {
		return ts, err
	}

	ts.SC, err = identity.GetServiceCatalog(ts.A)
	if err != nil {
		return ts, err
	}

	ts.EPs, err = FindAllEndpoints(ts.SC, service)
	if err != nil {
		return ts, err
	}

	ts.W = new(tabwriter.Writer)
	ts.W.Init(os.Stdout, 2, 8, 2, ' ', 0)

	return ts, nil
}

func SetupForCRUD() (*testState, error) {
	ts, err := SetupForList("compute")
	if err != nil {
		return ts, err
	}

	ts.ImageId = os.Getenv("OS_IMAGE_ID")
	if ts.ImageId == "" {
		return ts, fmt.Errorf("Expected OS_IMAGE_ID environment variable to be set")
	}

	ts.FlavorId = os.Getenv("OS_FLAVOR_ID")
	if ts.FlavorId == "" {
		return ts, fmt.Errorf("Expected OS_FLAVOR_ID environment variable to be set")
	}

	ts.FlavorIdResize = os.Getenv("OS_FLAVOR_ID_RESIZE")
	if ts.FlavorIdResize == "" {
		return ts, fmt.Errorf("Expected OS_FLAVOR_ID_RESIZE environment variable to be set")
	}

	if ts.FlavorIdResize == ts.FlavorId {
		return ts, fmt.Errorf("OS_FLAVOR_ID and OS_FLAVOR_ID_RESIZE cannot be the same")
	}

	ts.Region = os.Getenv("OS_REGION_NAME")
	if ts.Region == "" {
		ts.Region = ts.EPs[0].Region
	}

	ts.EP, err = FindEndpointForRegion(ts.EPs, ts.Region)
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

	gr, err := servers.GetDetail(ts.Client, ts.CreatedServer.Id)
	if err != nil {
		return false, timeout, err
	}

	ts.GottenServer, err = servers.GetServer(gr)
	if err != nil {
		return false, timeout, err
	}

	return true, timeout, nil
}

func CreateServer(ts *testState) error {
	ts.ServerName = RandomString("ACPTTEST", 16)
	fmt.Printf("Attempting to create server: %s\n", ts.ServerName)

	ts.Client = servers.NewClient(ts.EP, ts.A, ts.O)

	cr, err := servers.Create(ts.Client, map[string]interface{}{
		"flavorRef": ts.FlavorId,
		"imageRef":  ts.ImageId,
		"name":      ts.ServerName,
	})
	if err != nil {
		return err
	}

	ts.CreatedServer, err = servers.GetServer(cr)
	return err
}

func WaitForStatus(ts *testState, s string) error {
	var (
		inProgress bool
		timeout    int
		err        error
	)

	for inProgress, timeout, err = CountDown(ts, 300); inProgress; inProgress, timeout, err = CountDown(ts, timeout) {
		if ts.GottenServer.Id != ts.CreatedServer.Id {
			return fmt.Errorf("created server id (%s) != gotten server id (%s)", ts.CreatedServer.Id, ts.GottenServer.Id)
		}

		if ts.GottenServer.Status == s {
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

	ts.AlternateName = RandomString("ACPTTEST", 16)
	for ts.AlternateName == ts.ServerName {
		ts.AlternateName = RandomString("ACPTTEST", 16)
	}
	fmt.Println("Attempting to change server name")

	ur, err := servers.Update(ts.Client, ts.CreatedServer.Id, map[string]interface{}{
		"name": ts.AlternateName,
	})
	if err != nil {
		return err
	}

	ts.UpdatedServer, err = servers.GetServer(ur)
	if err != nil {
		return err
	}

	if ts.UpdatedServer.Id != ts.CreatedServer.Id {
		return fmt.Errorf("Expected updated and created server to share the same ID")
	}

	for inProgress, timeout, err = CountDown(ts, 300); inProgress; inProgress, timeout, err = CountDown(ts, timeout) {
		if ts.GottenServer.Id != ts.UpdatedServer.Id {
			return fmt.Errorf("Updated server ID (%s) != gotten server ID (%s)", ts.UpdatedServer.Id, ts.GottenServer.Id)
		}

		if ts.GottenServer.Name == ts.AlternateName {
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
	fmt.Println("Current password: "+oldPass)
	randomPassword := RandomString("", 16)
	for randomPassword == oldPass {
		randomPassword = RandomString("", 16)
	}
	fmt.Println("    New password: "+randomPassword)
	return randomPassword
}

func ChangeAdminPassword(ts *testState) error {
	randomPassword := MakeNewPassword(ts.CreatedServer.AdminPass)
	
	err := servers.ChangeAdminPassword(ts.Client, ts.CreatedServer.Id, randomPassword)
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
	fmt.Println("Attempting reboot of server "+ts.CreatedServer.Id)
	err := servers.Reboot(ts.Client, ts.CreatedServer.Id, servers.OSReboot)
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
	fmt.Println("Attempting to rebuild server "+ts.CreatedServer.Id)

	newPassword := MakeNewPassword(ts.CreatedServer.AdminPass)
	newName := RandomString("ACPTTEST", 16)
	sr, err := servers.Rebuild(ts.Client, ts.CreatedServer.Id, newName, newPassword, ts.ImageId, nil)
	if err != nil {
		return err
	}
	
	s, err := servers.GetServer(sr)
	if err != nil {
		return err
	}
	if s.Id != ts.CreatedServer.Id {
		return fmt.Errorf("Expected rebuilt server ID of %s; got %s", ts.CreatedServer.Id, s.Id)
	}

	err = WaitForStatus(ts, "REBUILD")
	if err != nil {
		return err
	}
	
	return WaitForStatus(ts, "ACTIVE")
}

func ResizeServer(ts *testState) error {
	fmt.Println("Attempting to resize server "+ts.CreatedServer.Id)

	err := servers.Resize(ts.Client, ts.CreatedServer.Id, ts.FlavorIdResize)
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
	fmt.Println("Attempting to confirm resize for server "+ts.CreatedServer.Id)
	
	err := servers.ConfirmResize(ts.Client, ts.CreatedServer.Id)
	if err != nil {
		return err
	}
	
	return WaitForStatus(ts, "ACTIVE")
}

func RevertResize(ts *testState) error {
	fmt.Println("Attempting to revert resize for server "+ts.CreatedServer.Id)
	
	err := servers.RevertResize(ts.Client, ts.CreatedServer.Id)
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
