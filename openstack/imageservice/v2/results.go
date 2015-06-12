package v2

import (
	"errors"
	"fmt"
	
	"github.com/rackspace/gophercloud"
	//"github.com/rackspace/gophercloud/pagination"

	"github.com/mitchellh/mapstructure"
)

// does not include the literal image data; just metadata.
// returned by listing images, and by fetching a specific image.
type Image struct {
	Id string
	
	Name string
	
	Status ImageStatus
	
	Tags []string
	
	ContainerFormat string
	DiskFormat string
	
	MinDiskGigabytes int `mapstructure:"min_disk"`
	MinRamMegabytes int `mapstructure:"min_ram"`
	
	Owner string
	
	Protected bool
	Visibility ImageVisibility

	Checksum *string `mapstructure:"checksum"`
	SizeBytes *int `mapstructure:"size"`
	
	Metadata map[string]string `mapstructure:"metadata"`
	Properties map[string]string `mapstructure:"properties"`
}

// implements pagination.Page<Image>, pagination.MarkerPage
// DOESN'T implement Page. Why? How does openstack/compute/v2
// type ImagePage struct {
// 	pagination.MarkerPageBase  // pagination.MarkerPageBase<Image>
// }

type CreateResult struct {
	gophercloud.ErrResult
}

func asBool(any interface{}) (bool, error) {
	if b, ok := any.(bool); ok {
		return b, nil
	} else {
		return false, errors.New(fmt.Sprintf("expected bool value, but found: %#v", any))
	}
}

func asInt(any interface{}) (int, error) {
	// FIXME integers decoded as float64s
	if f, ok := any.(float64); ok {
		i := int(f)
		return i, nil
	} else {
		return 0, errors.New(fmt.Sprintf("expected int value, but found: %#v", any))
	}
}

func asString(any interface{}) (string, error) {
	if str, ok := any.(string); ok {
		return str, nil
	} else {
		return "", errors.New(fmt.Sprintf("expected string value, but found: %#v", any))
	}
}

func asNoneableString(any interface{}) (*string, error) {
	if str, ok := any.(string); ok {
		if str == "None" {
			return nil, nil
		} else {
			return &str, nil
		}
	} else {
		return nil, errors.New(fmt.Sprintf("expected string value, but found: %#v", any))
	}
}

func asNoneableInteger(any interface{}) (*int, error) {
	// FIXME problem here is that provider_client.go uses: json.NewDecoder(resp.Body).Decode(options.JSONResponse)
	// which apparently converts integers in JSON to float64 values
	if f, ok := any.(float64); ok {
		i := int(f)
		return &i, nil
	} else if s, ok := any.(string); ok {
		if s == "None" {
			return nil, nil
		} else {
			return nil, errors.New(fmt.Sprintf("expected \"None\" or integer value, but found unexpected string: \"%s\"", s))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("expected \"None\" or integer value, but found: %T(%#v)", any, any))
	}
}

func asMapStringString(any interface{}) (map[string]string, error) {
	if mss, ok := any.(map[string]string); ok {
		return mss, nil
	} else {
		return nil, errors.New(fmt.Sprintf("expected map[string]string, but found: %#v", any))
	}
}

func extractBoolAtKey(m map[string]interface{}, k string) (bool, error) {
	if any, ok := m[k]; ok {
		return asBool(any)
	} else {
		return false, errors.New(fmt.Sprintf("expected key \"%s\" in map, but this key is not present", k))
	}
}

func extractIntAtKey(m map[string]interface{}, k string) (int, error) {
	if any, ok := m[k]; ok {
		return asInt(any)
	} else {
		return 0, errors.New(fmt.Sprintf("expected key \"%s\" in map, but this key is not present", k))
	}
}

func extractStringAtKey(m map[string]interface{}, k string) (string, error) {
	if any, ok := m[k]; ok {
		return asString(any)
	} else {
		return "", errors.New(fmt.Sprintf("expected key \"%s\" in map, but this key is not present", k))
	}
}

func extractNoneableStringAtKey(m map[string]interface{}, k string) (*string, error) {
	if any, ok := m[k]; ok {
		return asNoneableString(any)
	} else {
		return nil, errors.New(fmt.Sprintf("expected key \"%s\" in map, but this key is not present", k))
	}
}

func extractNoneableIntegerAtKey(m map[string]interface{}, k string) (*int, error) {
	if any, ok := m[k]; ok {
		return asNoneableInteger(any)
	} else {
		return nil, errors.New(fmt.Sprintf("expected key \"%s\" in map, but this key is not present", k))
	}
}

func extractStringSliceAtKey(m map[string]interface{}, k string) ([]string, error) {
	if any, ok := m[k]; ok {
		if slice, ok := any.([]interface{}); ok {
			res := make([]string, len(slice))
			for k, v := range slice {
				var err error
				if res[k], err = asString(v); err != nil {
					return nil, err
				}
			}
			return res, nil
		} else {
			return nil, errors.New(fmt.Sprintf("expected slice as \"%s\" value, but found: %#v", k, any))
		}
	} else {
		return nil, errors.New(fmt.Sprintf("expected key \"%s\" in map, but this key is not present", k))
	}
}

func stringToImageStatus(s string) (ImageStatus, error) {
	if s == "queued" {
		return ImageStatusQueued, nil
	} else if s == "active" {
		return ImageStatusActive, nil
	} else {
		return "", errors.New(fmt.Sprintf("expected \"queued\" or \"active\" as image status, but found: \"%s\"", s))
	}
}

func extractImageStatusAtKey(m map[string]interface{}, k string) (ImageStatus, error) {
	if any, ok := m[k]; ok {
		if str, ok := any.(string); ok {
			return stringToImageStatus(str)
		} else {
			return "", errors.New(fmt.Sprintf("expected string as \"%s\" value, but found: %#v", k, any))
		}
	} else {
		return "", errors.New(fmt.Sprintf("expected key \"%s\" in map, but this key is not present", k))
	}
}

func stringToImageVisibility(s string) (ImageVisibility, error) {
	if s == "public" {
		return ImageVisibilityPublic, nil
	} else if s == "private" {
		return ImageVisibilityPrivate, nil
	} else {
		return "", errors.New(fmt.Sprintf("expected \"public\" or \"private\" as image status, but found: \"%s\"", s))
	}
}

func extractImageVisibilityAtKey(m map[string]interface{}, k string) (ImageVisibility, error) {
	if any, ok := m[k]; ok {
		if str, ok := any.(string); ok {
			return stringToImageVisibility(str)
		} else {
			return "", errors.New(fmt.Sprintf("expected string as \"%s\" value, but found: %#v", k, any))
		}
	} else {
		return "", errors.New(fmt.Sprintf("expected key \"%s\" in map, but this key is not present", k))
	}
}

func extractMapStringStringAtKeyOptional(m map[string]interface{}, k string, ifMissing map[string]string) (map[string]string, error) {
	if any, ok := m[k]; ok {
		return asMapStringString(any)
	} else {
		return ifMissing, nil
	}
}

func extractImage(res gophercloud.ErrResult) (*Image, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	body, ok := res.Body.(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("expected map as result body, but found: %#v", res.Body))
	}
	
	var image Image

	var err error

	if image.Id, err = extractStringAtKey(body, "id"); err != nil {
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

	if image.ContainerFormat, err = extractStringAtKey(body, "container_format"); err != nil {
		return nil, err
	}
	
	if image.DiskFormat, err = extractStringAtKey(body, "disk_format"); err != nil {
		return nil, err
	}

	if image.MinDiskGigabytes, err = extractIntAtKey(body, "min_disk"); err != nil {
		return nil, err
	}

	if image.MinRamMegabytes, err = extractIntAtKey(body, "min_ram"); err != nil {
		return nil, err
	}
	
	if image.Owner, err = extractStringAtKey(body, "owner"); err != nil {
		return nil, err
	}

	if image.Protected, err = extractBoolAtKey(body, "protected"); err != nil {
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
	
	// TODO Metadata map[string]string `mapstructure:"metadata"`
	// TODO Properties map[string]string `mapstructure:"properties"`

	return &image, nil
}

// The response to `POST /images` follows the same schema as `GET /images/:id`.
func extractImageOld(res gophercloud.ErrResult) (*Image, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	var image Image

	err := mapstructure.Decode(res.Body, &image)
	
	return &image, err
}

func (c CreateResult) Extract() (*Image, error) {
	return extractImage(c.ErrResult)
}

type DeleteResult struct {
	gophercloud.ErrResult
}

type GetResult struct {
	gophercloud.ErrResult
}

func (c GetResult) Extract() (*Image, error) {
	return extractImage(c.ErrResult)
}

type UpdateResult struct {
	gophercloud.ErrResult
}
