package objects

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestExtractObjectMetadata(t *testing.T) {
	getResult := &http.Response{}

	expected := map[string]string{}

	actual := ExtractMetadata(getResult)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v\nActual:%+v", expected, actual)
	}
}

func TestExtractContent(t *testing.T) {
	responseBody := "'Twas brillig, and the slithy toves"
	downloadResult := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}
	expected := []byte("'Twas brillig, and the slithy toves")
	actual, err := ExtractContent(downloadResult)
	if err != nil {
		t.Errorf("Error extracting object content: %s", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %+v\nActual:%+v", expected, actual)
	}
}

func TestExtractObjectInfo(t *testing.T) {
	responseBody := `
		[
		    {
				"hash": "451e372e48e0f6b1114fa0724aa79fa1",
		        "last_modified": "2014-01-15T16:41:49.390270",
				"bytes": 14,
				"name": "goodbye",
				"content_type": "application/octet-stream"
			},
			{
		        "hash": "ed076287532e86365e841e92bfc50d8c",
			    "last_modified": "2014-01-15T16:37:43.427570",
				"bytes": 12,
				"name": "helloworld",
				"content_type": "application/octet-stream"
			}
		]	
	`

	listResult := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	var expected []Object
	err := json.Unmarshal([]byte(responseBody), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON: %s", err)
	}

	actual, err := ExtractInfo(listResult)
	if err != nil {
		t.Errorf("Error extracting objects info: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v\nActual: %+v", expected, actual)
	}
}

func TestExtractObjectNames(t *testing.T) {
	responseBody := "goodbye\nhelloworld\n"

	listResult := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	expected := []string{"goodbye", "helloworld"}

	actual, err := ExtractNames(listResult)
	if err != nil {
		t.Errorf("Error extracting object names: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v\nActual:%+v", expected, actual)
	}
}
