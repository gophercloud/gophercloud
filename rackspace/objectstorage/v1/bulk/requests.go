package bulk

import (
  "archive/tar"
  "compress/gzip"
  "compress/bzip2"
  "errors"
  "io"
	"net/url"
  "os"
  "path/filepath"
	"strings"

	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// DeleteOptsBuilder allows extensions to add additional parameters to the
// Delete request.
type DeleteOptsBuilder interface {
	ToBulkDeleteBody() (string, error)
}

// DeleteOpts is a structure that holds parameters for deleting an object.
type DeleteOpts []string

// ToBulkDeleteBody formats a DeleteOpts into a request body.
func (opts DeleteOpts) ToBulkDeleteBody() (string, error) {
	return url.QueryEscape(strings.Join(opts, "\n")), nil
}

// Delete will delete objects or containers in bulk.
func Delete(c *gophercloud.ServiceClient, opts DeleteOptsBuilder) DeleteResult {
	var res DeleteResult

	if opts == nil {
		return res
	}

	reqString, err := opts.ToBulkDeleteBody()
	if err != nil {
		res.Err = err
		return res
	}

  reqBody := strings.NewReader(reqString)

	resp, err := perigee.Request("DELETE", deleteURL(c), perigee.Options{
    ContentType: "text/plain",
		MoreHeaders: c.Provider.AuthenticatedHeaders(),
		OkCodes:     []int{200},
		ReqBody:     reqBody,
		Results:     &res.Body,
	})
	res.Header = resp.HttpResponse.Header
	res.Err = err
	return res
}

// Extract will extract the files in `file` and create objects in object storage
// from them.
func Extract(c *gophercloud.ServiceClient, file string) ExtractResult {
  var res ExtractResult

  if file == ""{
    res.Err = errors.New("Missing required field 'f'.")
    return res
  }

  var ext string
  var reqBody io.Reader
  f, err := os.Open(file)
  if err != nil {
    res.Err = errors.New("Error opening file.")
    return res
  }
  defer f.Close()

  switch filepath.Ext(file) {
    case "tar":
      ext = "tar"
      reqBody = tar.NewReader(f)
    case "gz":
      ext = "tar.gz"
      reqBody, err = gzip.NewReader(f)
      if err != nil {
        res.Err = err
        return res
      }
    case "bz2":
      ext = "tar.bz2"
      reqBody = bzip2.NewReader(f)
    default:
      res.Err = errors.New("Unsupported extension type.")
      return res
  }

  resp, err := perigee.Request("PUT", extractURL(c, ext), perigee.Options{
    MoreHeaders: c.Provider.AuthenticatedHeaders(),
    OkCodes:     []int{200},
    ReqBody:     reqBody,
    Results:     &res.Body,
  })
  res.Header = resp.HttpResponse.Header
  res.Err = err
  return res
}
