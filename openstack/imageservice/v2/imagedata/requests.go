package imagedata

import (
	"context"
	"io"

	"github.com/gophercloud/gophercloud/v2"
)

// Upload uploads an image file.
func Upload(ctx context.Context, client *gophercloud.ServiceClient, id string, data io.Reader) (r UploadResult) {
	resp, err := client.PutWithContext(ctx, uploadURL(client, id), data, nil, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/octet-stream"},
		OkCodes:     []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Stage performs PUT call on the existing image object in the Imageservice with
// the provided file.
// Existing image object must be in the "queued" status.
func Stage(ctx context.Context, client *gophercloud.ServiceClient, id string, data io.Reader) (r StageResult) {
	resp, err := client.PutWithContext(ctx, stageURL(client, id), data, nil, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/octet-stream"},
		OkCodes:     []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Download retrieves an image.
func Download(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DownloadResult) {
	resp, err := client.GetWithContext(ctx, downloadURL(client, id), nil, &gophercloud.RequestOpts{
		KeepResponseBody: true,
	})
	r.Body, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
