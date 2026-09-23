package auth

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

// newRequestV2 builds request options for identity v2 authentication.
func newRequestV2(opts AuthOptionsBuilderV2) (*gophercloud.RequestOpts, error) {
	authData, err := opts.ToAuthBody()
	if err != nil {
		return nil, err
	}

	if len(authData) != 1 {
		return nil, gophercloud.ErrInvalidInput{
			ErrMissingInput: gophercloud.ErrMissingInput{Argument: "AuthMethods"},
			Value:           len(authData),
		}
	}

	var authBody map[string]any
	for _, body := range authData {
		authBody = body
	}

	return &gophercloud.RequestOpts{
		JSONBody:    map[string]any{"auth": authBody},
		OkCodes:     []int{200, 203},
		OmitHeaders: []string{"X-Auth-Token", "Authorization"},
	}, nil
}

// RequestOption configures NewRequestV3
type requestOption func(*requestOptions)

type requestOptions struct {
	tokenURL string
}

func withTokenURL(tokenURL string) requestOption {
	return func(o *requestOptions) { o.tokenURL = tokenURL }
}

// newRequestV3 builds request options for identity v3 authentication.
func newRequestV3(opts AuthOptionsBuilderV3, options ...requestOption) (*gophercloud.RequestOpts, error) {
	authData, err := opts.ToAuthBody()
	if err != nil {
		return nil, err
	}

	methods := slices.Collect(maps.Keys(authData))
	slices.Sort(methods)
	identity := map[string]any{"methods": methods}
	for method, body := range authData {
		identity[method] = body
	}

	authBody := map[string]any{"identity": identity}
	scope, err := opts.ToAuthScope()
	if err != nil {
		return nil, err
	}
	if scope != nil {
		authBody["scope"] = scope
	}

	headers, err := opts.ToAuthHeaders(options...)
	if err != nil {
		return nil, err
	}
	moreHeaders := make(map[string]string, len(headers))
	omitHeaders := []string{"X-Auth-Token", "Authorization"}
	for name, value := range headers {
		moreHeaders[name] = fmt.Sprint(value)
		if strings.EqualFold(name, "Authorization") {
			omitHeaders = omitHeaders[:1]
		}
	}

	return &gophercloud.RequestOpts{
		JSONBody:    map[string]any{"auth": authBody},
		MoreHeaders: moreHeaders,
		OmitHeaders: omitHeaders,
	}, nil
}

// newRequestEC2 builds request options from EC2 credentials, omitting the
// credential elements which are unused by the target endpoint.
func newRequestEC2(opts AuthOptionsBuilderEC2, omit ...string) (*gophercloud.RequestOpts, error) {
	authData, err := opts.ToAuthBody()
	if err != nil {
		return nil, err
	}

	if len(authData) != 1 {
		return nil, gophercloud.ErrInvalidInput{
			ErrMissingInput: gophercloud.ErrMissingInput{Argument: "Credentials"},
			Value:           len(authData),
		}
	}

	authBody := make(map[string]any, len(authData))
	for name, credentials := range authData {
		for _, k := range omit {
			delete(credentials, k)
		}
		authBody[name] = credentials
	}

	return &gophercloud.RequestOpts{
		JSONBody: authBody,
		OkCodes:  []int{200},
	}, nil
}
