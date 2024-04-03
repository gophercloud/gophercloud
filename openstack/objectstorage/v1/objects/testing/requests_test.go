package testing

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	v1 "github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1"
	accountTesting "github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/accounts/testing"
	containerTesting "github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers/testing"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestContainerNames(t *testing.T) {
	for _, tc := range [...]struct {
		name          string
		containerName string
		expectedError error
	}{
		{
			"rejects_a_slash",
			"one/two",
			v1.ErrInvalidContainerName{},
		},
		{
			"rejects_an_empty_string",
			"",
			v1.ErrEmptyContainerName{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("list", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleListObjectsInfoSuccessfully(t, WithPath("/"))

				_, err := objects.List(fake.ServiceClient(), tc.containerName, nil).AllPages(context.TODO())
				th.CheckErr(t, err, &tc.expectedError)
			})
			t.Run("download", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleDownloadObjectSuccessfully(t, WithPath("/"))

				_, err := objects.Download(context.TODO(), fake.ServiceClient(), tc.containerName, "testObject", nil).Extract()
				th.CheckErr(t, err, &tc.expectedError)
			})
			t.Run("create", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				content := "Ceci n'est pas une pipe"
				HandleCreateTextObjectSuccessfully(t, content, WithPath("/"))

				res := objects.Create(context.TODO(), fake.ServiceClient(), tc.containerName, "testObject", &objects.CreateOpts{
					ContentType: "text/plain",
					Content:     strings.NewReader(content),
				})
				th.CheckErr(t, res.Err, &tc.expectedError)
			})
			t.Run("delete", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleDeleteObjectSuccessfully(t, WithPath("/"))

				res := objects.Delete(context.TODO(), fake.ServiceClient(), tc.containerName, "testObject", nil)
				th.CheckErr(t, res.Err, &tc.expectedError)
			})
			t.Run("get", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleGetObjectSuccessfully(t, WithPath("/"))

				_, err := objects.Get(context.TODO(), fake.ServiceClient(), tc.containerName, "testObject", nil).ExtractMetadata()
				th.CheckErr(t, err, &tc.expectedError)
			})
			t.Run("update", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleUpdateObjectSuccessfully(t)

				res := objects.Update(context.TODO(), fake.ServiceClient(), tc.containerName, "testObject", &objects.UpdateOpts{
					Metadata: map[string]string{"Gophercloud-Test": "objects"},
				})
				th.CheckErr(t, res.Err, &tc.expectedError)
			})
			t.Run("createTempURL", func(t *testing.T) {
				port := 33200
				th.SetupHTTP()
				th.SetupPersistentPortHTTP(t, port)
				defer th.TeardownHTTP()

				// Handle fetching of secret key inside of CreateTempURL
				containerTesting.HandleGetContainerSuccessfully(t)
				accountTesting.HandleGetAccountSuccessfully(t)
				client := fake.ServiceClient()

				// Append v1/ to client endpoint URL to be compliant with tempURL generator
				client.Endpoint = client.Endpoint + "v1/"
				_, err := objects.CreateTempURL(context.TODO(), client, tc.containerName, "testObject/testFile.txt", objects.CreateTempURLOpts{
					Method:    http.MethodGet,
					TTL:       60,
					Timestamp: time.Date(2020, 07, 01, 01, 12, 00, 00, time.UTC),
				})

				th.CheckErr(t, err, &tc.expectedError)
			})
			t.Run("bulk-delete", func(t *testing.T) {
				th.SetupHTTP()
				defer th.TeardownHTTP()
				HandleBulkDeleteSuccessfully(t)

				res := objects.BulkDelete(context.TODO(), fake.ServiceClient(), tc.containerName, []string{"testObject"})
				th.CheckErr(t, res.Err, &tc.expectedError)
			})
		})
	}
}

func TestDownloadReader(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDownloadObjectSuccessfully(t)

	response := objects.Download(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", nil)
	defer response.Body.Close()

	// Check reader
	buf := bytes.NewBuffer(make([]byte, 0))
	_, err := io.CopyN(buf, response.Body, 10)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "Successful", buf.String())
}

func TestDownloadExtraction(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDownloadObjectSuccessfully(t)

	response := objects.Download(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", nil)

	// Check []byte extraction
	bytes, err := response.ExtractContent()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "Successful download with Gophercloud", string(bytes))

	expected := &objects.DownloadHeader{
		ContentLength:     36,
		ContentType:       "text/plain; charset=utf-8",
		Date:              time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		StaticLargeObject: true,
		LastModified:      time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
	}
	actual, err := response.Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)
}

func TestDownloadWithLastModified(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDownloadObjectSuccessfully(t)

	options1 := &objects.DownloadOpts{
		IfUnmodifiedSince: time.Date(2009, time.November, 10, 22, 59, 59, 0, time.UTC),
	}
	response1 := objects.Download(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options1)
	_, err1 := response1.Extract()
	th.AssertErr(t, err1)

	options2 := &objects.DownloadOpts{
		IfModifiedSince: time.Date(2009, time.November, 10, 23, 0, 1, 0, time.UTC),
	}
	response2 := objects.Download(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options2)
	content, err2 := response2.ExtractContent()
	th.AssertNoErr(t, err2)
	th.AssertEquals(t, 0, len(content))
}

func TestListObjectInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListObjectsInfoSuccessfully(t)

	count := 0
	options := &objects.ListOpts{}
	err := objects.List(fake.ServiceClient(), "testContainer", options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := objects.ExtractInfo(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedListInfo, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListObjectSubdir(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSubdirSuccessfully(t)

	count := 0
	options := &objects.ListOpts{Prefix: "", Delimiter: "/"}
	err := objects.List(fake.ServiceClient(), "testContainer", options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := objects.ExtractInfo(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedListSubdir, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListObjectNames(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListObjectsInfoSuccessfully(t)

	// Check without delimiter.
	count := 0
	options := &objects.ListOpts{}
	err := objects.List(fake.ServiceClient(), "testContainer", options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := objects.ExtractNames(page)
		if err != nil {
			t.Errorf("Failed to extract container names: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, ExpectedListNames, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)

	// Check with delimiter.
	count = 0
	options = &objects.ListOpts{Delimiter: "/"}
	err = objects.List(fake.ServiceClient(), "testContainer", options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := objects.ExtractNames(page)
		if err != nil {
			t.Errorf("Failed to extract container names: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, ExpectedListNames, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestListZeroObjectNames204(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListZeroObjectNames204(t)

	count := 0
	options := &objects.ListOpts{}
	err := objects.List(fake.ServiceClient(), "testContainer", options).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := objects.ExtractNames(page)
		if err != nil {
			t.Errorf("Failed to extract container names: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, []string{}, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 0, count)
}

func TestCreateObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	content := "Did gyre and gimble in the wabe"

	HandleCreateTextObjectSuccessfully(t, content)

	options := &objects.CreateOpts{ContentType: "text/plain", Content: strings.NewReader(content)}
	res := objects.Create(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options)
	th.AssertNoErr(t, res.Err)
}

func TestCreateObjectWithCacheControl(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	content := "All mimsy were the borogoves"

	HandleCreateTextWithCacheControlSuccessfully(t, content)

	options := &objects.CreateOpts{
		CacheControl: `max-age="3600", public`,
		Content:      strings.NewReader(content),
	}
	res := objects.Create(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options)
	th.AssertNoErr(t, res.Err)
}

func TestCreateObjectWithoutContentType(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	content := "The sky was the color of television, tuned to a dead channel."

	HandleCreateTypelessObjectSuccessfully(t, content)

	res := objects.Create(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", &objects.CreateOpts{Content: strings.NewReader(content)})
	th.AssertNoErr(t, res.Err)
}

func TestCopyObject(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()
		HandleCopyObjectSuccessfully(t, "/newTestContainer/newTestObject")

		options := &objects.CopyOpts{Destination: "/newTestContainer/newTestObject"}
		res := objects.Copy(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options)
		th.AssertNoErr(t, res.Err)
	})
	t.Run("slash", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()
		HandleCopyObjectSuccessfully(t, "/newTestContainer/path%2Fto%2FnewTestObject")

		options := &objects.CopyOpts{Destination: "/newTestContainer/path/to/newTestObject"}
		res := objects.Copy(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options)
		th.AssertNoErr(t, res.Err)
	})
	t.Run("emojis", func(t *testing.T) {
		th.SetupHTTP()
		defer th.TeardownHTTP()
		HandleCopyObjectSuccessfully(t, "/newTestContainer/new%F0%9F%98%8ATest%2C%3B%22O%28bject%21_%E7%AF%84")

		options := &objects.CopyOpts{Destination: "/newTestContainer/new😊Test,;\"O(bject!_範"}
		res := objects.Copy(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options)
		th.AssertNoErr(t, res.Err)
	})
}

func TestCopyObjectVersion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCopyObjectVersionSuccessfully(t)

	options := &objects.CopyOpts{Destination: "/newTestContainer/newTestObject", ObjectVersionID: "123456788"}
	res, err := objects.Copy(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "123456789", res.ObjectVersionID)
}

func TestDeleteObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteObjectSuccessfully(t)

	res := objects.Delete(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", nil)
	th.AssertNoErr(t, res.Err)
}

func TestBulkDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleBulkDeleteSuccessfully(t)

	expected := objects.BulkDeleteResponse{
		ResponseStatus: "foo",
		ResponseBody:   "bar",
		NumberDeleted:  2,
		Errors:         [][]string{},
	}

	resp, err := objects.BulkDelete(context.TODO(), fake.ServiceClient(), "testContainer", []string{"testObject1", "testObject2"}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, *resp)
}

func TestUpateObjectMetadata(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateObjectSuccessfully(t)

	s := new(string)
	i := new(int64)
	options := &objects.UpdateOpts{
		Metadata:           map[string]string{"Gophercloud-Test": "objects"},
		RemoveMetadata:     []string{"Gophercloud-Test-Remove"},
		ContentDisposition: s,
		ContentEncoding:    s,
		ContentType:        s,
		DeleteAt:           i,
		DetectContentType:  new(bool),
	}
	res := objects.Update(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", options)
	th.AssertNoErr(t, res.Err)
}

func TestGetObject(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetObjectSuccessfully(t)

	expected := map[string]string{"Gophercloud-Test": "objects"}
	actual, err := objects.Get(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", nil).ExtractMetadata()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expected, actual)

	getOpts := objects.GetOpts{
		Newest: true,
	}
	actualHeaders, err := objects.Get(context.TODO(), fake.ServiceClient(), "testContainer", "testObject", getOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, actualHeaders.StaticLargeObject)
}

func TestETag(t *testing.T) {
	content := "some example object"
	createOpts := objects.CreateOpts{
		Content: strings.NewReader(content),
		NoETag:  true,
	}

	_, headers, _, err := createOpts.ToObjectCreateParams()
	th.AssertNoErr(t, err)
	_, ok := headers["ETag"]
	th.AssertEquals(t, false, ok)

	hash := md5.New()
	_, err = io.WriteString(hash, content)
	th.AssertNoErr(t, err)
	localChecksum := fmt.Sprintf("%x", hash.Sum(nil))

	createOpts = objects.CreateOpts{
		Content: strings.NewReader(content),
		ETag:    localChecksum,
	}

	_, headers, _, err = createOpts.ToObjectCreateParams()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, localChecksum, headers["ETag"])
}

func TestObjectCreateParamsWithoutSeek(t *testing.T) {
	content := "I do not implement Seek()"
	buf := strings.NewReader(content)

	createOpts := objects.CreateOpts{Content: buf}
	reader, headers, _, err := createOpts.ToObjectCreateParams()

	th.AssertNoErr(t, err)

	_, ok := reader.(io.ReadSeeker)
	th.AssertEquals(t, true, ok)

	c, err := io.ReadAll(reader)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, content, string(c))

	_, ok = headers["ETag"]
	th.AssertEquals(t, true, ok)
}

func TestObjectCreateParamsWithSeek(t *testing.T) {
	content := "I implement Seek()"
	createOpts := objects.CreateOpts{Content: strings.NewReader(content)}
	reader, headers, _, err := createOpts.ToObjectCreateParams()

	th.AssertNoErr(t, err)

	_, ok := reader.(io.ReadSeeker)
	th.AssertEquals(t, ok, true)

	c, err := io.ReadAll(reader)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, content, string(c))

	_, ok = headers["ETag"]
	th.AssertEquals(t, true, ok)
}

func TestCreateTempURL(t *testing.T) {
	port := 33200
	th.SetupHTTP()
	th.SetupPersistentPortHTTP(t, port)
	defer th.TeardownHTTP()

	// Handle fetching of secret key inside of CreateTempURL
	containerTesting.HandleGetContainerSuccessfully(t)
	accountTesting.HandleGetAccountSuccessfully(t)
	client := fake.ServiceClient()

	// Append v1/ to client endpoint URL to be compliant with tempURL generator
	client.Endpoint = client.Endpoint + "v1/"
	tempURL, err := objects.CreateTempURL(context.TODO(), client, "testContainer", "testObject/testFile.txt", objects.CreateTempURLOpts{
		Method:    http.MethodGet,
		TTL:       60,
		Timestamp: time.Date(2020, 07, 01, 01, 12, 00, 00, time.UTC),
	})

	sig := "89be454a9c7e2e9f3f50a8441815e0b5801cba5b"
	expiry := "1593565980"
	expectedURL := fmt.Sprintf("http://127.0.0.1:%v/v1/testContainer/testObject%%2FtestFile.txt?temp_url_sig=%v&temp_url_expires=%v", port, sig, expiry)

	th.AssertNoErr(t, err)
	th.AssertEquals(t, expectedURL, tempURL)

	// Test TTL=0, but different timestamp
	tempURL, err = objects.CreateTempURL(context.TODO(), client, "testContainer", "testObject/testFile.txt", objects.CreateTempURLOpts{
		Method:    http.MethodGet,
		Timestamp: time.Date(2020, 07, 01, 01, 13, 00, 00, time.UTC),
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, expectedURL, tempURL)
}
