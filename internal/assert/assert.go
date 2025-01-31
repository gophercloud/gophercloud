package assert

import (
	"fmt"
	"strconv"
)

const errf = "unknown type for %s: %T (value: %v)"

func Int(v any, name string) (int, error) {
	switch t := v.(type) {
	case nil:
		return 0, nil
	case int:
		return t, nil
	case int64:
		return int(t), nil
	case float32:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		return strconv.Atoi(t)
	}

	return 0, fmt.Errorf(errf, name, v, v)
}

func Int64(v any, name string) (int64, error) {
	switch t := v.(type) {
	case nil:
		return 0, nil
	case int:
		return int64(t), nil
	case int64:
		return t, nil
	case float32:
		return int64(t), nil
	case float64:
		return int64(t), nil
	case string:
		return strconv.ParseInt(t, 10, 64)
	}

	return 0, fmt.Errorf(errf, name, v, v)
}

func String(v any, name string) (string, error) {
	switch t := v.(type) {
	case nil:
		return "", nil
	case string:
		return t, nil
	}

	return "", fmt.Errorf(errf, name, v, v)
}

func Bool(v any, name string) (bool, error) {
	switch t := v.(type) {
	case nil:
		return false, nil
	case bool:
		return t, nil
	case string:
		return strconv.ParseBool(t)
	}

	return false, fmt.Errorf(errf, name, v, v)
}
