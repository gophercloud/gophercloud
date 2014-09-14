package containers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)


func TestExtractContainerMetadata(t *testing.T) {
	getResult := &http.Response{}

	expected := map[string]string{}

	actual := ExtractMetadata(getResult)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v\nActual:%+v", expected, actual)
	}
}

func TestExtractContainerInfo(t *testing.T) {
	responseBody := `
		[
			{
				"count": 3,
				"bytes": 2000,
				"name": "artemis"
			},
			{
				"count": 1,
				"bytes": 450,
				"name": "diana"
			}
		]
	`

	listResult := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	var expected []Container
	err := json.Unmarshal([]byte(responseBody), &expected)
	if err != nil {
		t.Errorf("Error unmarshaling JSON: %s", err)
	}

	actual, err := ExtractInfo(listResult)
	if err != nil {
		t.Errorf("Error extracting containers info: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nExpected: %+v\nActual:   %+v", expected, actual)
	}
}

func TestExtractConatinerNames(t *testing.T) {
	responseBody := "artemis\ndiana\n"

	listResult := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}

	expected := []string{"artemis", "diana"}

	actual, err := ExtractNames(listResult)
	if err != nil {
		t.Errorf("Error extracting container names: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected: %+v\nActual:%+v", expected, actual)
	}
}
