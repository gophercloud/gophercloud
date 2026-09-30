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
		if service, ok := strings.CutSuffix(strings.ReplaceAll(key, "-", "_"), defaultMicroversionSuffix); ok && service != "" {
			service = strings.ReplaceAll(service, "_", "-")
			if node.Tag == "!!null" {
				if value.nullDefaultMicroversions == nil {
					value.nullDefaultMicroversions = make(map[string]bool)
				}
				value.nullDefaultMicroversions[service] = true
				continue
			}
			var version string
			if err := node.Decode(&version); err != nil {
				return err
			}
			if value.DefaultMicroversions == nil {
				value.DefaultMicroversions = make(map[string]string)
			}
			value.DefaultMicroversions[service] = version
		}
	}
	*c = Cloud(value)
	return nil
}

// MarshalYAML keeps service defaults in the clouds.yaml format.
func (c Cloud) MarshalYAML() (any, error) {
	type plain Cloud
	fields := make(map[string]any, len(c.DefaultMicroversions)+len(c.nullDefaultMicroversions))
	for service := range c.nullDefaultMicroversions {
		fields[strings.ReplaceAll(service, "-", "_")+defaultMicroversionSuffix] = nil
	}
	for service, version := range c.DefaultMicroversions {
		fields[strings.ReplaceAll(service, "-", "_")+defaultMicroversionSuffix] = version
	}
	return struct {
		Cloud    plain          `yaml:",inline"`
		Defaults map[string]any `yaml:",inline"`
	}{plain(c), fields}, nil
}

// MarshalJSON uses the same flat keys and scalar types as the YAML representation.
func (c Cloud) MarshalJSON() ([]byte, error) {
	var node yaml.Node
	if err := node.Encode(c); err != nil {
		return nil, err
	}
	var fields map[string]any
	if err := node.Decode(&fields); err != nil {
		return nil, err
	}
	return json.Marshal(fields)
}

// UnmarshalJSON shares YAML's handling of service keys and nullable defaults.
func (c *Cloud) UnmarshalJSON(data []byte) error {
	return yaml.Unmarshal(data, c)
}
