package gophercloud

import (
	"testing"
	"time"

	th "github.com/rackspace/gophercloud/testhelper"
)

func TestWaitFor(t *testing.T) {
	err := WaitFor(0, func() (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	})
	if err == nil {
		t.Errorf("Expected error: 'Time out in WaitFor'")
	}

	err = WaitFor(5, func() (bool, error) {
		return true, nil
	})
	th.CheckNoErr(t, err)
}
