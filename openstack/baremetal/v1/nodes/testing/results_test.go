package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/inventory"
	invtest "github.com/gophercloud/gophercloud/v2/openstack/baremetal/inventory/testing"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetalintrospection/v1/introspection"
	insptest "github.com/gophercloud/gophercloud/v2/openstack/baremetalintrospection/v1/introspection/testing"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestStandardPluginData(t *testing.T) {
	var pluginData nodes.PluginData
	err := pluginData.RawMessage.UnmarshalJSON([]byte(invtest.StandardPluginDataSample))
	th.AssertNoErr(t, err)

	parsedData, err := pluginData.AsStandardData()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, invtest.StandardPluginData, parsedData)

	irData, inspData, err := pluginData.GuessFormat()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, invtest.StandardPluginData, *irData)
	th.CheckEquals(t, (*introspection.Data)(nil), inspData)
}

func TestInspectorPluginData(t *testing.T) {
	var pluginData nodes.PluginData
	err := pluginData.RawMessage.UnmarshalJSON([]byte(insptest.IntrospectionDataJSONSample))
	th.AssertNoErr(t, err)

	parsedData, err := pluginData.AsInspectorData()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, insptest.IntrospectionDataRes, parsedData)

	irData, inspData, err := pluginData.GuessFormat()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, (*inventory.StandardPluginData)(nil), irData)
	th.CheckDeepEquals(t, insptest.IntrospectionDataRes, *inspData)
}

func TestGuessFormatUnknownDefaultsToIronic(t *testing.T) {
	var pluginData nodes.PluginData
	err := pluginData.RawMessage.UnmarshalJSON([]byte("{}"))
	th.AssertNoErr(t, err)

	irData, inspData, err := pluginData.GuessFormat()
	th.CheckDeepEquals(t, inventory.StandardPluginData{}, *irData)
	th.CheckEquals(t, (*introspection.Data)(nil), inspData)
	th.AssertNoErr(t, err)
}

func TestGuessFormatErrors(t *testing.T) {
	var pluginData nodes.PluginData
	err := pluginData.RawMessage.UnmarshalJSON([]byte("\"banana\""))
	th.AssertNoErr(t, err)

	irData, inspData, err := pluginData.GuessFormat()
	th.CheckEquals(t, (*inventory.StandardPluginData)(nil), irData)
	th.CheckEquals(t, (*introspection.Data)(nil), inspData)
	th.AssertErr(t, err)

	failsInspectorConversion := `{
	    "interfaces": "banana"
	}`
	err = pluginData.RawMessage.UnmarshalJSON([]byte(failsInspectorConversion))
	th.AssertNoErr(t, err)

	irData, inspData, err = pluginData.GuessFormat()
	th.CheckEquals(t, (*inventory.StandardPluginData)(nil), irData)
	th.CheckEquals(t, (*introspection.Data)(nil), inspData)
	th.AssertErr(t, err)
}
