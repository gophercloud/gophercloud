package gophercloud

// CommonResult acts as a base struct that other results can embed. It contains
// the deserialized JSON structure returned from the server (Resp), and any
// errors that might have occurred during transport or deserialization.
type CommonResult struct {
	Resp map[string]interface{}
	Err  error
}

// RFC3339Milli describes a time format used by API responses.
const RFC3339Milli = "2006-01-02T15:04:05.999999Z"
