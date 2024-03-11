package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	//Create http server config
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	//Run http server in another groutine
	//to be able to stop it from external action
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil &&
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
	err := run(context.Background())
	if err != nil {
		fmt.Printf("Failed to terminate server: %v", err)
		os.Exit(1)
	}
}
