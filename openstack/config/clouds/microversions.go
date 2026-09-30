package clouds

import (
	"encoding/json"
	"strings"

	"go.yaml.in/yaml/v3"
)

const defaultMicroversionSuffix = "_default_microversion"

// UnmarshalYAML preserves service defaults without interpreting unrelated keys.
func (c *Cloud) UnmarshalYAML(unmarshal func(any) error) error {
	type plain Cloud
	var value plain
	if err := unmarshal(&value); err != nil {
		return err
	}
	var fields map[string]yaml.Node
	if err := unmarshal(&fields); err != nil {
		return err
	}
	for key, node := range fields {
		if service, ok := strings.CutSuffix(key, defaultMicroversionSuffix); ok && service != "" {
			var version string
			if err := node.Decode(&version); err != nil {
				return err
			}
			if value.DefaultMicroversions == nil {
				value.DefaultMicroversions = make(map[string]string)
			}
			value.DefaultMicroversions[strings.ReplaceAll(service, "_", "-")] = version
		}
	}
	*c = Cloud(value)
	return nil
}

// MarshalYAML keeps service defaults in the clouds.yaml format.
func (c Cloud) MarshalYAML() (any, error) {
	type plain Cloud
	fields := make(map[string]string, len(c.DefaultMicroversions))
	for service, version := range c.DefaultMicroversions {
		fields[strings.ReplaceAll(service, "-", "_")+defaultMicroversionSuffix] = version
	}
	return struct {
		Cloud    plain             `yaml:",inline"`
		Defaults map[string]string `yaml:",inline"`
	}{plain(c), fields}, nil
}

// MarshalJSON preserves flat service keys for configuration merging.
func (c Cloud) MarshalJSON() ([]byte, error) {
	type plain Cloud
	data, err := json.Marshal(plain(c))
	if err != nil {
		return nil, err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	for service, version := range c.DefaultMicroversions {
		data, err := json.Marshal(version)
		if err != nil {
			return nil, err
		}
		fields[strings.ReplaceAll(service, "-", "_")+defaultMicroversionSuffix] = data
	}
	return json.Marshal(fields)
}

// UnmarshalJSON restores service defaults after configuration merging.
func (c *Cloud) UnmarshalJSON(data []byte) error {
	type plain Cloud
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	for key, data := range fields {
		if service, ok := strings.CutSuffix(key, defaultMicroversionSuffix); ok && service != "" {
			var version string
			if err := json.Unmarshal(data, &version); err != nil {
				return err
			}
			if value.DefaultMicroversions == nil {
				value.DefaultMicroversions = make(map[string]string)
			}
			value.DefaultMicroversions[strings.ReplaceAll(service, "_", "-")] = version
		}
	}
	*c = Cloud(value)
	return nil
}
