# Tips

## Implementing default logging and re-authentication attempts

You can implement custom logging and/or limit re-auth attempts by creating a custom HTTP client
like the following and setting it as the provider client's HTTP Client (via the
`gophercloud.ProviderClient.HTTPClient` field):

```go
//...

// LogRoundTripper satisfies the http.RoundTripper interface and is used to
// customize the default Gophercloud RoundTripper to allow for logging.
type LogRoundTripper struct {
	rt                http.RoundTripper
	numReauthAttempts int
}

// newHTTPClient return a custom HTTP client that allows for logging relevant
// information before and after the HTTP request.
func newHTTPClient() http.Client {
	return http.Client{
		Transport: &LogRoundTripper{
			rt: http.DefaultTransport,
		},
	}
}

// RoundTrip performs a round-trip HTTP request and logs relevant information about it.
func (lrt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	glog.Infof("Request URL: %s\n", request.URL)

	response, err := lrt.rt.RoundTrip(request)
	if response == nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		if lrt.numReauthAttempts == 3 {
			return response, fmt.Errorf("Tried to re-authenticate 3 times with no success.")
		}
		lrt.numReauthAttempts++
	}

	glog.Debugf("Response Status: %s\n", response.Status)

	return response, nil
}

endpoint := "https://127.0.0.1/auth"
pc := openstack.NewClient(endpoint)
pc.HTTPClient = newHTTPClient()

//...
```
