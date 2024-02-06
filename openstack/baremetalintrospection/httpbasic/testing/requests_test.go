package testing

import (
	"encoding/base64"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetalintrospection/httpbasic"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestNoAuth(t *testing.T) {
	httpClient, err := httpbasic.NewBareMetalIntrospectionHTTPBasic(httpbasic.EndpointOpts{
		IronicInspectorEndpoint:     "http://ironic:5050/v1",
		IronicInspectorUser:         "myUser",
		IronicInspectorUserPassword: "myPasswd",
	})
	th.AssertNoErr(t, err)
	encToken := base64.StdEncoding.EncodeToString([]byte("myUser:myPasswd"))
	headerValue := "Basic " + encToken
	th.AssertEquals(t, headerValue, httpClient.MoreHeaders["Authorization"])

	errTest1, err := httpbasic.NewBareMetalIntrospectionHTTPBasic(httpbasic.EndpointOpts{
		IronicInspectorEndpoint: "http://ironic:5050/v1",
	})
	_ = errTest1
	th.AssertEquals(t, "User and Password are required", err.Error())

	errTest2, err := httpbasic.NewBareMetalIntrospectionHTTPBasic(httpbasic.EndpointOpts{
		IronicInspectorUser:         "myUser",
		IronicInspectorUserPassword: "myPasswd",
	})
	_ = errTest2
	th.AssertEquals(t, "IronicInspectorEndpoint is required", err.Error())
}
