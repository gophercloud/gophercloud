package v2

import (
	"encoding/json"
	"testing"
)

func TestAccessToken(t *testing.T) {
	authResults := make(map[string]interface{})
	err := json.Unmarshal([]byte(authResultsOK), &authResults)
	if err != nil {
		t.Error(err)
		return
	}

	tok, err := GetToken(authResults)
	if err != nil {
		t.Error(err)
		return
	}
	if tok.ID != "ab48a9efdfedb23ty3494" {
		t.Errorf("Expected token \"ab48a9efdfedb23ty3494\"; got \"%s\" instead", tok.ID)
		return
	}
}
