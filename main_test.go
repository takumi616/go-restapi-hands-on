package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	//Create ctx with cancel func to test if http server will be stopped
	//by external action intentionally
	ctx, cancel := context.WithCancel(context.Background())

	//Run http server in another groutine
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})

	//Test if http server works correctly
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}

	//Compare response body to expected one
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	//Send cancel signal to groutine which http server runs in
	cancel()

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
