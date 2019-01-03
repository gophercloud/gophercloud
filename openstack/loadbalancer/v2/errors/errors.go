package errors

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
)

type NeutronError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type errorDetails struct {
	NeutronError NeutronError `json:"NeutronError"`
}

// error types from provider_client.go
func ExtractErrorInto(rawError error, neutronError *NeutronError) (err error) {
	var eDetails errorDetails

	switch e := rawError.(type) {
	case gophercloud.ErrDefault400:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	case gophercloud.ErrDefault401:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	case gophercloud.ErrDefault403:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	case gophercloud.ErrDefault404:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	case gophercloud.ErrDefault405:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	case gophercloud.ErrDefault408:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	case gophercloud.ErrDefault429:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	case gophercloud.ErrDefault500:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	case gophercloud.ErrDefault503:
		log.Printf("[DEBUG] lberror: %s", e.ErrUnexpectedResponseCode.Body)
		err = json.Unmarshal(e.ErrUnexpectedResponseCode.Body, &eDetails)
	default:
		err = fmt.Errorf("Unable to extract detailed error message")
	}

	*neutronError = eDetails.NeutronError

	return err
}
