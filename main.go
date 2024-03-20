package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/takumi616/go-restapi-hands-on/config"
)

func run(ctx context.Context) error {
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

	//Get routing info
	//cleanup is used to close *sql.DB
	mux, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		return err
	}

	//Close *sql.DB
	defer cleanup()

	//Get Server config
	s := NewServer(l, mux)

	//Start http server
	return s.Run(ctx)
}

func main() {
	err := run(context.Background())
	if err != nil {
		log.Printf("Failed to terminate server: %v", err)
	}
}
