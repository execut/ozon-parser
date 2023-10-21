package chromeCookie

import (
	"testing"
)

func TestReadTokenSuccess(t *testing.T) {
	token, _ := ReadToken()

	wantValue := "test-token-value"
	if token != wantValue {
		t.Errorf("Want token value %q; got %q", wantValue, token)
	}
}
