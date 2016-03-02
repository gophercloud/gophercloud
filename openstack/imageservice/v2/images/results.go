package images

import (
	"fmt"
	"io"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Image model
// Does not include the literal image data; just metadata.
// returned by listing images, and by fetching a specific image.
type Image struct {
	ID string

	Name string

	Status ImageStatus

	Tags []string

	ContainerFormat *string
	DiskFormat      *string

	MinDiskGigabytes *int
	MinRAMMegabytes  *int

	Owner *string

	Protected  bool
	Visibility ImageVisibility

	Checksum  *string
	SizeBytes *int

	Metadata   map[string]string
	Properties map[string]string
}

// CreateResult model
type CreateResult struct {
	gophercloud.ErrResult
}

func asBool(any interface{}) (bool, error) {
	if b, ok := any.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("expected bool value, but found: %#v", any)
}

func asInt(any interface{}) (int, error) {
	// FIXME integers decoded as float64s
	if f, ok := any.(float64); ok {
		i := int(f)
		return i, nil
	}
	return 0, fmt.Errorf("expected int value, but found: %#v", any)
}

func asString(any interface{}, key string) (string, error) {
	if str, ok := any.(string); ok {
		return str, nil
	}
	return "", fmt.Errorf("expected string value for key '%s', but found: %#v", key, any)
}

func asNoneableString(any interface{}, key string) (*string, error) {
	// JSON null values could be also returned according to behaviour https://bugs.launchpad.net/glance/+bug/1481512
	if any == nil {
		return nil, nil
	}
	if str, ok := any.(string); ok {
		if str == "None" || &str == nil {
			return nil, nil
		}
		return &str, nil
	}
	return nil, fmt.Errorf("expected string value for key '%s', but found: %#v", key, any)
}

func asNoneableInteger(any interface{}, key string) (*int, error) {
	// FIXME problem here is that provider_client.go uses: json.NewDecoder(resp.Body).Decode(options.JSONResponse)
	// which apparently converts integers in JSON to float64 values
	// JSON null values could be also returned according to behaviour https://bugs.launchpad.net/glance/+bug/1481512
	if any == nil {
		return nil, nil
	}
	if f, ok := any.(float64); ok {
		i := int(f)
		return &i, nil
	} else if s, ok := any.(string); ok {
		if s == "None" {
			return nil, nil
		}
		return nil, fmt.Errorf("expected \"None\" or integer value for key '%s', but found unexpected string: \"%s\"", key, s)
	}
	return nil, fmt.Errorf("expected \"None\" or integer value for key '%s', but found: %T(%#v)", key, any, any)
}

func asMapStringString(any interface{}) (map[string]string, error) {
	if mss, ok := any.(map[string]string); ok {
		return mss, nil
	}
	return nil, fmt.Errorf("expected map[string]string, but found: %#v", any)
}

func extractBoolAtKey(m map[string]interface{}, k string) (bool, error) {
	if any, ok := m[k]; ok {
		return asBool(any)
	}
	return false, fmt.Errorf("expected key \"%s\" in map, but this key is not present", k)
}

func extractIntAtKey(m map[string]interface{}, k string) (int, error) {
	if any, ok := m[k]; ok {
		return asInt(any)
	}
	return 0, fmt.Errorf("expected key \"%s\" in map, but this key is not present", k)
}

func extractStringAtKey(m map[string]interface{}, k string) (string, error) {
	if any, ok := m[k]; ok {
		return asString(any, k)
	}
	return "", fmt.Errorf("expected key \"%s\" in map, but this key is not present", k)
}

func extractNoneableStringAtKey(m map[string]interface{}, k string) (*string, error) {
	if any, ok := m[k]; ok {
		return asNoneableString(any, k)
	}
	return nil, fmt.Errorf("expected key \"%s\" in map, but this key is not present", k)
}

func extractNoneableIntegerAtKey(m map[string]interface{}, k string) (*int, error) {
	if any, ok := m[k]; ok {
		return asNoneableInteger(any, k)
	}
	return nil, fmt.Errorf("expected key \"%s\" in map, but this key is not present", k)
}

func extractStringSliceAtKey(m map[string]interface{}, key string) ([]string, error) {
	if any, ok := m[key]; ok {
		if slice, ok := any.([]interface{}); ok {
			res := make([]string, len(slice))
			for k, v := range slice {
				var err error
				if res[k], err = asString(v, key); err != nil {
					return nil, err
				}
			}
			return res, nil
		}
		return nil, fmt.Errorf("expected slice as \"%s\" value, but found: %#v", key, any)
	}
	return nil, fmt.Errorf("expected key \"%s\" in map, but this key is not present", key)
}

func stringToImageStatus(s string) (ImageStatus, error) {
	if s == "queued" {
		return ImageStatusQueued, nil
	} else if s == "active" {
		return ImageStatusActive, nil
	} else {
		return "", fmt.Errorf("expected \"queued\" or \"active\" as image status, but found: \"%s\"", s)
	}
}

func extractImageStatusAtKey(m map[string]interface{}, k string) (ImageStatus, error) {
	if any, ok := m[k]; ok {
		if str, ok := any.(string); ok {
			return stringToImageStatus(str)
		}
		return "", fmt.Errorf("expected string as \"%s\" value, but found: %#v", k, any)
	}
	return "", fmt.Errorf("expected key \"%s\" in map, but this key is not present", k)

}

func stringToImageVisibility(s string) (ImageVisibility, error) {
	if s == "public" {
		return ImageVisibilityPublic, nil
	} else if s == "private" {
		return ImageVisibilityPrivate, nil
	} else {
		return "", fmt.Errorf("expected \"public\" or \"private\" as image status, but found: \"%s\"", s)
	}
}

func extractImageVisibilityAtKey(m map[string]interface{}, k string) (ImageVisibility, error) {
	if any, ok := m[k]; ok {
		if str, ok := any.(string); ok {
			return stringToImageVisibility(str)
		}
		return "", fmt.Errorf("expected string as \"%s\" value, but found: %#v", k, any)
	}
	return "", fmt.Errorf("expected key \"%s\" in map, but this key is not present", k)
}

func extractBoolAtKeyOptional(m map[string]interface{}, k string, ifMissing bool) (bool, error) {
	if any, ok := m[k]; ok {
		return asBool(any)
	}
	return ifMissing, nil

}

func extractMapStringStringAtKeyOptional(m map[string]interface{}, k string, ifMissing map[string]string) (map[string]string, error) {
	if any, ok := m[k]; ok {
		return asMapStringString(any)
	}
	return ifMissing, nil
}

func extractImage(res gophercloud.ErrResult) (*Image, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	body, ok := res.Body.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map as result body, but found: %#v", res.Body)
	}

	var image Image

	var err error

	if image.ID, err = extractStringAtKey(body, "id"); err != nil {
		return nil, err
	}

	if image.Name, err = extractStringAtKey(body, "name"); err != nil {
		return nil, err
	}

	if image.Status, err = extractImageStatusAtKey(body, "status"); err != nil {
		return nil, err
	}

	if image.Tags, err = extractStringSliceAtKey(body, "tags"); err != nil {
		return nil, err
	}

	if image.ContainerFormat, err = extractNoneableStringAtKey(body, "container_format"); err != nil {
		return nil, err
	}

	if image.DiskFormat, err = extractNoneableStringAtKey(body, "disk_format"); err != nil {
		return nil, err
	}

	if image.MinDiskGigabytes, err = extractNoneableIntegerAtKey(body, "min_disk"); err != nil {
		return nil, err
	}

	if image.MinRAMMegabytes, err = extractNoneableIntegerAtKey(body, "min_ram"); err != nil {
		return nil, err
	}

	if image.Owner, err = extractNoneableStringAtKey(body, "owner"); err != nil {
		return nil, err
	}

	// FIXME should this key actually be optional? Is a missing key equivalent to "protected": false ?
	if image.Protected, err = extractBoolAtKeyOptional(body, "protected", false); err != nil {
		return nil, err
	}

	if image.Visibility, err = extractImageVisibilityAtKey(body, "visibility"); err != nil {
		return nil, err
	}

	if image.Checksum, err = extractNoneableStringAtKey(body, "checksum"); err != nil {
		return nil, err
	}

	if image.SizeBytes, err = extractNoneableIntegerAtKey(body, "size"); err != nil {
		return nil, err
	}

	if image.Metadata, err = extractMapStringStringAtKeyOptional(body, "metadata", make(map[string]string)); err != nil {
		return nil, err
	}

	if image.Properties, err = extractMapStringStringAtKeyOptional(body, "properties", make(map[string]string)); err != nil {
		return nil, err
	}

	return &image, nil
}

// Extract build CreateResults from imput Image
func (c CreateResult) Extract() (*Image, error) {
	return extractImage(c.ErrResult)
}

//DeleteResult model
type DeleteResult struct {
	gophercloud.ErrResult
}

// GetResult model
type GetResult struct {
	gophercloud.ErrResult
}

// Extract builds GetResult
func (c GetResult) Extract() (*Image, error) {
	return extractImage(c.ErrResult)
}

// UpdateResult model
type UpdateResult struct {
	gophercloud.ErrResult
}

// Extract builds UpdateResult
func (u UpdateResult) Extract() (*Image, error) {
	return extractImage(u.ErrResult)
}

// PutImageDataResult is model put image respose
type PutImageDataResult struct {
	gophercloud.ErrResult
}

// GetImageDataResult model for image response
type GetImageDataResult struct {
	gophercloud.ErrResult
}

// Extract builds images model from io.Reader
func (g GetImageDataResult) Extract() (io.Reader, error) {
	if r, ok := g.Body.(io.Reader); ok {
		return r, nil
	}
	return nil, fmt.Errorf("Expected io.Reader but got: %T(%#v)", g.Body, g.Body)
}

// ImagePage represents page
type ImagePage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Images results.
func (page ImagePage) IsEmpty() (bool, error) {
	images, err := ExtractImages(page)
	if err != nil {
		return true, err
	}
	return len(images) == 0, nil
}

// NextPageURL uses the response's embedded link reference to navigate to the next page of results.
func (page ImagePage) NextPageURL() (string, error) {
	type resp struct {
		Next string `mapstructure:"next"`
	}

	var r resp
	err := mapstructure.Decode(page.Body, &r)
	if err != nil {
		return "", err
	}

	return nextPageURL(page.URL.String(), r.Next), nil
}

func toMapFromString(from reflect.Kind, to reflect.Kind, data interface{}) (interface{}, error) {
	if (from == reflect.String) && (to == reflect.Map) {
		return map[string]interface{}{}, nil
	}
	return data, nil
}

// ExtractImages interprets the results of a single page from a List() call, producing a slice of Image entities.
func ExtractImages(page pagination.Page) ([]Image, error) {
	casted := page.(ImagePage).Body

	var response struct {
		Images []Image `mapstructure:"images"`
	}

	config := &mapstructure.DecoderConfig{
		DecodeHook: toMapFromString,
		Result:     &response,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(casted)
	if err != nil {
		fmt.Printf("Error happened %v \n", err)
	}

	return response.Images, err
}
