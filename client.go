package gophercloud

import (
	"context"
	"net/http"
)

type Client interface {
	ServiceURL(...string) string
	Get(ctx context.Context, url string, JSONResponse any, opts *RequestOpts) (*http.Response, error)
	Post(ctx context.Context, url string, JSONBody any, JSONResponse any, opts *RequestOpts) (*http.Response, error)
	Put(ctx context.Context, url string, JSONBody any, JSONResponse any, opts *RequestOpts) (*http.Response, error)
	Patch(ctx context.Context, url string, JSONBody any, JSONResponse any, opts *RequestOpts) (*http.Response, error)
	Delete(ctx context.Context, url string, opts *RequestOpts) (*http.Response, error)
	Head(ctx context.Context, url string, opts *RequestOpts) (*http.Response, error)
	Request(ctx context.Context, method, url string, options *RequestOpts) (*http.Response, error)
}
