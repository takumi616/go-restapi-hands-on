package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/takumi616/go-restapi-hands-on/config"
	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	//Create ctx with stop signal
	//Server sends response that connection is closed
	//before finishing process by these signals
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	//Get environment variables
	cfg, err := config.New()
	if err != nil {
		return err
	}

	//Create http listener with port received from config file
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen port %d: %v", cfg.Port, err)
	}

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("Start with: %v", url)

	//Create http server config
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
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
	err := run(context.Background())
	if err != nil {
		log.Printf("Failed to terminate server: %v", err)
	}
}
