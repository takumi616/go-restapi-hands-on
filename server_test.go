package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	//Create http listener
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen port %v", err)
	}

	//Create handler
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	//Create ctx with cancel func to test if http server will be stopped
	//by external action intentionally
	ctx, cancel := context.WithCancel(context.Background())

	//Run http server in another groutine
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		s := NewServer(l, mux)
		return s.Run(ctx)
	})

	//Test if http server works correctly
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	//Check randomly selected port number
	t.Logf("Request URL: %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("Failed to get: %+v", err)
	}

	//Compare response body to expected one
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
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
