package gophercloud

import (
	"errors"
	"strings"
	"time"
)

// WaitFor polls a predicate function once per second up to secs times to wait
// for a certain state to arrive.
func WaitFor(timeout int, predicate func() (bool, error)) error {
	start := time.Now().Second()
	for {
		// Force a 1s sleep
		time.Sleep(1 * time.Second)

		// If a timeout is set, and that's been exceeded, shut it down
		if timeout >= 0 && time.Now().Second()-start >= timeout {
			return errors.New("A timeout occurred")
		}

		// Execute the function
		satisfied, err := predicate()
		if err != nil {
			return err
		}
		if satisfied {
			return nil
		}
	}
}

// NormalizeURL ensures that each endpoint URL has a closing `/`, as expected by ServiceClient.
func NormalizeURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	}
	return url
}

// BuildQuery constructs the query section of a URI from a map.
func BuildQuery(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	query := "?"
	for k, v := range params {
		query += k + "=" + v + "&"
	}
	query = query[:len(query)-1]
	return query
}
