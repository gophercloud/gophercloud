package testing

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/imagedata"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestUpload(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandlePutImageDataSuccessfully(t, fakeServer)

	err := imagedata.Upload(
		context.TODO(),
		client.ServiceClient(fakeServer),
		"da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		readSeekerOfBytes([]byte{5, 3, 7, 24})).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestStage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleStageImageDataSuccessfully(t, fakeServer)

	err := imagedata.Stage(
		context.TODO(),
		client.ServiceClient(fakeServer),
		"da3b75d9-3f4a-40e7-8a2c-bfab23927dea",
		readSeekerOfBytes([]byte{5, 3, 7, 24})).ExtractErr()

	th.AssertNoErr(t, err)
}

func readSeekerOfBytes(bs []byte) io.ReadSeeker {
	return &RS{bs: bs}
}

// implements io.ReadSeeker
type RS struct {
	bs     []byte
	offset int
}

func (rs *RS) Read(p []byte) (int, error) {
	leftToRead := len(rs.bs) - rs.offset

	if 0 < leftToRead {
		bytesToWrite := min(leftToRead, len(p))
		for i := 0; i < bytesToWrite; i++ {
			p[i] = rs.bs[rs.offset]
			rs.offset++
		}
		return bytesToWrite, nil
	}
	return 0, io.EOF
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func (rs *RS) Seek(offset int64, whence int) (int64, error) {
	var offsetInt = int(offset)
	switch whence {
	case 0:
		rs.offset = offsetInt
	case 1:
		rs.offset = rs.offset + offsetInt
	case 2:
		rs.offset = len(rs.bs) - offsetInt
	default:
		return 0, fmt.Errorf("for parameter `whence`, expected value in {0,1,2} but got: %#v", whence)
	}

	return int64(rs.offset), nil
}

func TestDownload(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetImageDataSuccessfully(t, fakeServer)

	rdr, err := imagedata.Download(context.TODO(), client.ServiceClient(fakeServer), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea").Extract()
	th.AssertNoErr(t, err)

	defer rdr.Close()

	bs, err := io.ReadAll(rdr)
	th.AssertNoErr(t, err)

	th.AssertByteArrayEquals(t, []byte{34, 87, 0, 23, 23, 23, 56, 255, 254, 0}, bs)
}
