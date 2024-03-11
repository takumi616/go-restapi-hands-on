package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context, l net.Listener) error {
	//Create http server config
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	//Run http server in another groutine
	//to be able to stop it from external action
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("Failed to close: %+v", err)
			return err
		}
		return nil
	})

	//Wait until ctx is canceled
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("Failed to shutdown: %+v", err)
	}

	return eg.Wait()
}

func main() {
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}

	//Create http listener with port received from os argument
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("Failed to listen port %s: %v", p, err)
	}

	err = run(context.Background(), l)
	if err != nil {
		log.Printf("Failed to terminate server: %v", err)
	}
}
