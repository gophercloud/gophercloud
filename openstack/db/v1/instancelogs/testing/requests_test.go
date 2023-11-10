package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/db/v1/instancelogs"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleList(t)

	pages := 0
	err := instancelogs.List(fake.ServiceClient(), instanceID).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := instancelogs.ExtractLogs(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, []instancelogs.Log{expectedLog}, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestShow(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleShow(t)

	log, err := instancelogs.Show(fake.ServiceClient(), instanceID, logName).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedShowLog, log)
}

func TestEnable(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleEnable(t)

	log, err := instancelogs.Enable(fake.ServiceClient(), instanceID, logName).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedEnableLog, log)
}

func TestDisable(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDisable(t)

	log, err := instancelogs.Disable(fake.ServiceClient(), instanceID, logName).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedDisableLog, log)
}

func TestPublish(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePublish(t)

	log, err := instancelogs.Publish(fake.ServiceClient(), instanceID, logName).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedPublishLog, log)
}

func TestDiscard(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDiscard(t)

	log, err := instancelogs.Discard(fake.ServiceClient(), instanceID, logName).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedDiscardLog, log)
}
