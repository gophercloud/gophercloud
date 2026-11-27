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
	// Name is the server's display name.
	Name string
	// FlavorRef is the flavor ID.
	FlavorRef string
	// ImageRef is the image ID.
	ImageRef string
}
```

to `servers.Create`, simply implement the `servers.CreateOptsBuilder` interface:

```go
func (o MyCreateServerOpts) ToServerCreateMap() (map[string]any, error) {
	return map[string]any{
		"server": map[string]any{
			"name": o.Name,
			"flavorRef": o.FlavorRef,
			"imageRef": o.ImageRef,
		},
	}, nil
}
```

create an instance of your custom options object, and pass it to `servers.Create`:

```go
// ...
myOpts := MyCreateServerOpts{
	Name:      "s1",
	FlavorRef: "flavor-id",
	ImageRef:  "image-id",
}
server, err := servers.Create(context.TODO(), computeClient, myOpts, nil).Extract()
// ...
```

### Custom response objects

Define a custom response object when a deployment returns additional fields
that are not represented by Gophercloud's standard response types. For example,
a deployment might add a `custom_field` attribute to a volume:

```go
// ...
type MyVolume struct {
	// ID is the volume ID.
	ID string `json:"id"`
	// Name is the volume's display name.
	Name string `json:"name"`
	// CustomField is an additional deployment-specific attribute.
	CustomField string `json:"custom_field"`
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
	// ID is the volume ID.
	ID string `json:"id"`
	// TimeCreated is decoded by UnmarshalJSON.
	TimeCreated time.Time `json:"-"`
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
	*r = MyVolume(s.tmp)

	r.TimeCreated = time.Time(s.TimeCreated)

	return err
}
// ...
```
