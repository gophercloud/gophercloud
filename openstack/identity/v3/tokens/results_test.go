package tokens

import "testing"

func TestTokenID(t *testing.T) {
	result := TokenCreateResult{tokenID: "1234"}

	token, _ := result.TokenID()
	if token != "1234" {
		t.Errorf("Expected tokenID of 1234, got %s", token)
	}
}
