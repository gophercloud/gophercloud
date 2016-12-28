package gophercloud

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// WaitFor polls a predicate function, once per second, up to a timeout limit.
// It usually does this to wait for a resource to transition to a certain state.
// To handle situations when the predicate response might hang indefinitely,
// the predicate will be retried every "retry" seconds.
// Resource packages will wrap this in a more convenient function that's
// specific to a certain resource, but it can also be useful on its own.
func WaitFor(timeout, retry int, predicate func() (bool, error)) error {
	type WaitForResult struct {
		Success bool
		Error   error
	}

	start := time.Now().Unix()

	for {
		// If a timeout is set, and that's been exceeded, shut it down
		if timeout >= 0 && time.Now().Unix()-start >= int64(timeout) {
			return errors.New("A timeout occurred")
		}

		time.Sleep(1 * time.Second)

		var result WaitForResult
		ch := make(chan WaitForResult, 1)
		go func() {
			defer close(ch)
			satisfied, err := predicate()
			result.Success = satisfied
			result.Error = err
		}()

		if retry > 0 {
			select {
			case <-ch:
				fmt.Println("tick")
				if result.Error != nil {
					return result.Error
				}
				if result.Success {
					return nil
				}
			case <-time.After(time.Duration(retry) * time.Second):
				fmt.Println("tock")
				continue
			}
		} else {
			<-ch
			fmt.Println("tick")
			if result.Error != nil {
				return result.Error
			}
			if result.Success {
				return nil
			}
		}
	}
}

// NormalizeURL is an internal function to be used by provider clients.
//
// It ensures that each endpoint URL has a closing `/`, as expected by
// ServiceClient's methods.
func NormalizeURL(url string) string {
	if !strings.HasSuffix(url, "/") {
		return url + "/"
	}
	return url
}

// NormalizePathURL is used to convert rawPath to a fqdn, using basePath as
// a reference in the filesystem, if necessary. basePath is assumed to contain
// either '.' when first used, or the file:// type fqdn of the parent resource.
// e.g. myFavScript.yaml => file://opt/lib/myFavScript.yaml
func NormalizePathURL(basePath, rawPath string) (string, error) {
	u, err := url.Parse(rawPath)
	if err != nil {
		return "", err
	}
	// if a scheme is defined, it must be a fqdn already
	if u.Scheme != "" {
		return u.String(), nil
	}
	// if basePath is a url, then child resources are assumed to be relative to it
	bu, err := url.Parse(basePath)
	if err != nil {
		return "", err
	}
	var basePathSys, absPathSys string
	if bu.Scheme != "" {
		basePathSys = filepath.FromSlash(bu.Path)
		absPathSys = filepath.Join(basePathSys, rawPath)
		bu.Path = filepath.ToSlash(absPathSys)
		return bu.String(), nil
	}

	absPathSys = filepath.Join(basePath, rawPath)
	u.Path = filepath.ToSlash(absPathSys)
	if err != nil {
		return "", err
	}
	u.Scheme = "file"
	return u.String(), nil

}
