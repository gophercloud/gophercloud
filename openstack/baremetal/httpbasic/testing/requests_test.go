package testing

import (
	"encoding/base64"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/httpbasic"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestHttpBasic(t *testing.T) {
	httpClient, err := httpbasic.NewBareMetalHTTPBasic(httpbasic.EndpointOpts{
		IronicEndpoint:     "http://ironic:6385/v1",
		IronicUser:         "myUser",
		IronicUserPassword: "myPasswd",
	})
	th.AssertNoErr(t, err)
	encToken := base64.StdEncoding.EncodeToString([]byte("myUser:myPasswd"))
	headerValue := "Basic " + encToken

	th.AssertEquals(t, headerValue, httpClient.MoreHeaders["Authorization"])

	errTest1, err := httpbasic.NewBareMetalHTTPBasic(httpbasic.EndpointOpts{
		IronicEndpoint: "http://ironic:6385/v1",
	})
	_ = errTest1
	th.AssertEquals(t, "User and Password are required", err.Error())

	errTest2, err := httpbasic.NewBareMetalHTTPBasic(httpbasic.EndpointOpts{
		IronicUser:         "myUser",
		IronicUserPassword: "myPasswd",
	})
	_ = errTest2
	th.AssertEquals(t, "IronicEndpoint is required", err.Error())

}
