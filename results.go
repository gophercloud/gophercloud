package gophercloud

// CommonResult acts as a base struct that other results can embed. It contains
// the deserialized JSON structure returned from the server (Resp), and any
// errors that might have occurred during transport or deserialization.
type CommonResult struct {
	Resp map[string]interface{}
	Err  error
}
