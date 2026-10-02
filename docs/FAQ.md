# Tips

## Handling Microversions

Please see our dedicated document [here](MICROVERSIONS.md).

## Customizing HTTP transport and reauthentication

Use `config.WithHTTPClient` to supply an HTTP client for authentication and
subsequent service requests, including automatic reauthentication. Whether
reauthentication is permitted is determined by the authenticator: set
`AllowReauth: true` on reusable manual credentials, or use the credential
loaders' defaults. One-time passcodes cannot be reused for reauthentication.

A logging transport can wrap the standard transport without keeping mutable
per-request state. Do not log tokens, passwords, or request bodies.

```go
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
)

// LogRoundTripper logs request methods and response status codes.
type LogRoundTripper struct {
	// Transport handles the actual HTTP requests.
	Transport http.RoundTripper
}

func (lrt *LogRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := lrt.Transport.RoundTrip(request)
	if err != nil {
		return response, err
	}
	log.Printf("%s: %s", request.Method, response.Status)
	return response, nil
}

func main() {
	opts, err := auth.AuthOptionsFromEnv()
	if err != nil {
		panic(err)
	}
	httpClient := http.Client{
		Transport: &LogRoundTripper{Transport: http.DefaultTransport},
	}
	provider, err := config.NewProviderClient(
		context.Background(), opts, config.WithHTTPClient(httpClient),
	)
	if err != nil {
		panic(err)
	}
	_ = provider // Use the provider to construct service clients.
}
```

## Implementing custom objects

OpenStack request/response objects may differ among variable names or types.

### Custom request objects

To pass custom options to a request, implement the desired `<ACTION>OptsBuilder` interface. For
example, to pass in

```go
type MyCreateServerOpts struct {
	Name string
	Size int
}
```

to `servers.Create`, simply implement the `servers.CreateOptsBuilder` interface:

```go
func (o MyCreateServeropts) ToServerCreateMap() (map[string]interface{}, error) {
	return map[string]interface{}{
		"name": o.Name,
		"size": o.Size,
	}, nil
}
```

create an instance of your custom options object, and pass it to `servers.Create`:

```go
// ...
myOpts := MyCreateServerOpts{
	Name: "s1",
	Size: "100",
}
server, err := servers.Create(context.TODO(), computeClient, myOpts).Extract()
// ...
```

### Custom response objects

Some OpenStack services have extensions. Extensions that are supported in Gophercloud can be
combined to create a custom object:

```go
// ...
type MyVolume struct {
  volumes.Volume
  tenantattr.VolumeExt
}

var v struct {
  MyVolume `json:"volume"`
}

err := volumes.Get(context.TODO(), client, volID).ExtractInto(&v)
// ...
```

## Overriding default `UnmarshalJSON` method

For some response objects, a field may be a custom type or may be allowed to take on
different types. In these cases, overriding the default `UnmarshalJSON` method may be
necessary. To do this, declare the JSON `struct` field tag as "-" and create an `UnmarshalJSON`
method on the type:

```go
// ...
type MyVolume struct {
	ID string `json: "id"`
	TimeCreated time.Time `json: "-"`
}

func (r *MyVolume) UnmarshalJSON(b []byte) error {
	type tmp MyVolume
	var s struct {
		tmp
		TimeCreated gophercloud.JSONRFC3339MilliNoZ `json:"created_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Volume(s.tmp)

	r.TimeCreated = time.Time(s.CreatedAt)

	return err
}
// ...
```
