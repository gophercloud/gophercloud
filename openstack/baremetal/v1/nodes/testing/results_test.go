package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/inventory"
	invtest "github.com/gophercloud/gophercloud/v2/openstack/baremetal/inventory/testing"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetalintrospection/v1/introspection"
	insptest "github.com/gophercloud/gophercloud/v2/openstack/baremetalintrospection/v1/introspection/testing"
	"github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestStandardPluginData(t *testing.T) {
	var pluginData nodes.PluginData
	err := pluginData.RawMessage.UnmarshalJSON([]byte(invtest.StandardPluginDataSample))
	testhelper.AssertNoErr(t, err)

	parsedData, err := pluginData.AsStandardData()
	testhelper.AssertNoErr(t, err)
	testhelper.CheckDeepEquals(t, invtest.StandardPluginData, parsedData)

	irData, inspData, err := pluginData.GuessFormat()
	testhelper.AssertNoErr(t, err)
	testhelper.CheckDeepEquals(t, invtest.StandardPluginData, *irData)
	testhelper.CheckEquals(t, (*introspection.Data)(nil), inspData)
}

func TestInspectorPluginData(t *testing.T) {
	var pluginData nodes.PluginData
	err := pluginData.RawMessage.UnmarshalJSON([]byte(insptest.IntrospectionDataJSONSample))
	testhelper.AssertNoErr(t, err)

	parsedData, err := pluginData.AsInspectorData()
	testhelper.AssertNoErr(t, err)
	testhelper.CheckDeepEquals(t, insptest.IntrospectionDataRes, parsedData)

	irData, inspData, err := pluginData.GuessFormat()
	testhelper.AssertNoErr(t, err)
	testhelper.CheckEquals(t, (*inventory.StandardPluginData)(nil), irData)
	testhelper.CheckDeepEquals(t, insptest.IntrospectionDataRes, *inspData)
}

func TestGuessFormatUnknownDefaultsToIronic(t *testing.T) {
	var pluginData nodes.PluginData
	err := pluginData.RawMessage.UnmarshalJSON([]byte("{}"))
	testhelper.AssertNoErr(t, err)

	irData, inspData, err := pluginData.GuessFormat()
	testhelper.CheckDeepEquals(t, inventory.StandardPluginData{}, *irData)
	testhelper.CheckEquals(t, (*introspection.Data)(nil), inspData)
	testhelper.AssertNoErr(t, err)
}

func TestGuessFormatErrors(t *testing.T) {
	var pluginData nodes.PluginData
	err := pluginData.RawMessage.UnmarshalJSON([]byte("\"banana\""))
	testhelper.AssertNoErr(t, err)

	irData, inspData, err := pluginData.GuessFormat()
	testhelper.CheckEquals(t, (*inventory.StandardPluginData)(nil), irData)
	testhelper.CheckEquals(t, (*introspection.Data)(nil), inspData)
	testhelper.AssertErr(t, err)

	failsInspectorConversion := `{
	    "interfaces": "banana"
	}`
	err = pluginData.RawMessage.UnmarshalJSON([]byte(failsInspectorConversion))
	testhelper.AssertNoErr(t, err)

	irData, inspData, err = pluginData.GuessFormat()
	testhelper.CheckEquals(t, (*inventory.StandardPluginData)(nil), irData)
	testhelper.CheckEquals(t, (*introspection.Data)(nil), inspData)
	testhelper.AssertErr(t, err)
}
