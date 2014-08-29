package tokens

import (
	"testing"
	"time"
)

func TestTokenID(t *testing.T) {
	result := TokenCreateResult{tokenID: "1234"}

	token, _ := result.TokenID()
	if token != "1234" {
		t.Errorf("Expected tokenID of 1234, got %s", token)
	}
}

func TestExpiresAt(t *testing.T) {
	resp := map[string]interface{}{
		"token": map[string]string{
			"expires_at": "2013-02-02T18:30:59.000000Z",
		},
	}

	result := TokenCreateResult{
		tokenID:  "1234",
		response: resp,
	}

	expected, _ := time.Parse(time.UnixDate, "Sat Feb 2 18:30:59 UTC 2013")
	actual, err := result.ExpiresAt()
	if err != nil {
		t.Errorf("Error extraction expiration time: %v", err)
	}
	if actual != expected {
		t.Errorf("Expected expiration time %s, but was %s", expected.Format(time.UnixDate), actual.Format(time.UnixDate))
	}
}
