package gophercloud

import (
	"fmt"
	"strings"
	"time"
)

// WaitFor polls a predicate function once per second up to secs times to wait for a certain state to arrive.
func WaitFor(secs int, predicate func() (bool, error)) error {
	for i := 0; i < secs; i++ {
		time.Sleep(1 * time.Second)

		satisfied, err := predicate()
		if err != nil {
			return err
		}
		if satisfied {
			return nil
		}
	}
	return fmt.Errorf("Time out in WaitFor.")
}

// NormalizeURL ensures that each endpoint URL has a closing `/`, as expected by ServiceClient.
func NormalizeURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	}
	return url
}
