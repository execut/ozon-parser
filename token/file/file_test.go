package file

import (
	"execut/ozon_parser/internal/testutils"
	"testing"
)

func TestReadTokenSuccess(t *testing.T) {
	testTokenPath, err := testutils.GetTestDataFilePath("token.txt")
	if err != nil {
		t.Fatalf("Failed to load test data file")
	}

	token, _ := File(testTokenPath)

	wantValue := "test-token-value"
	if token != wantValue {
		t.Errorf("Want token value %q; got %q", wantValue, token)
	}
}
