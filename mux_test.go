package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMux(t *testing.T) {
	t.SkipNow()

	//ResponseRecorder is an implementation of http.ResponseWriter
	//that records its mutations for later inspection in tests
	w := httptest.NewRecorder()
	//r := httptest.NewRequest(http.MethodGet, "/health", nil)

	//Get routing
	//sut := NewMux()
	//Send http request
	//sut.ServeHTTP(w, r)
	//Get http response
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	//Check http status code
	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}

	//Read http response body
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	//Compare received http response body to expected one
	want := `{"status": "ok"}`
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
}
