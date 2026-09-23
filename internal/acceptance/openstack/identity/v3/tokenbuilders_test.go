//go:build acceptance || identity

package v3

type tokenCreateOpts struct {
	Body    map[string]any
	Headers map[string]string
}

func (opts tokenCreateOpts) ToTokenCreateMap() (map[string]any, error) {
	return opts.Body, nil
}

func (opts tokenCreateOpts) ToTokenCreateHeaders(method, url string) (map[string]string, error) {
	return opts.Headers, nil
}

func (opts tokenCreateOpts) ToS3TokenValidateMap() (map[string]any, error) {
	return opts.Body, nil
}
