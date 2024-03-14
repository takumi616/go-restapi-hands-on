package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func AssertJSON(t *testing.T, want, got []byte) {
	//Helper marks the calling function as a test helper function.
	//When printing file and line information, this function will be skipped.
	t.Helper()

	//Convert json data into golang's data type
	var jw, jg any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("Failed to unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jg); err != nil {
		t.Fatalf("Failed to unmarshal got %q: %v", got, err)
	}

	//Compare diff
	if diff := cmp.Diff(jg, jw); diff != "" {
		t.Errorf("got differs: (-got +want)\n%s", diff)
	}
}

func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	//Helper marks the calling function as a test helper function.
	//When printing file and line information, this function will be skipped.
	t.Helper()

	//Read response body
	t.Cleanup(func() { _ = got.Body.Close() })
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}

	//Check http status code
	if got.StatusCode != status {
		t.Fatalf("want status %d, but got %d, body: %q", status, got.StatusCode, gb)
	}

	if len(gb) == 0 && len(body) == 0 {
		//Not need to call AssertJSON()
		//Because the length of want body and got body are empty
		return
	}

	//Check response body
	AssertJSON(t, body, gb)
}

func LoadFile(t *testing.T, path string) []byte {
	//Helper marks the calling function as a test helper function.
	//When printing file and line information, this function will be skipped.
	t.Helper()

	//Read test data file
	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return bt
}
