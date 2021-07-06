package main

import (
	"context"
	"fmt"
	api_client "github.com/nais/aiven-audit/pkg/api-client"
	"github.com/nais/aiven-audit/pkg/nais"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func root(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Aiven Audit")
	if err != nil {
		log.Fatal("Could not respond")
	}
}

func main() {
	http.HandleFunc("/", root)
	log.Println("starting aiven-audit...")
	nais.InitNaisHandlers() // Setup handling of nais paths

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Handle common signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
			case <-signalChan:
				log.Printf("Got SIGINT/SIGTERM, exiting.")
				cancel()
				os.Exit(1)
			case <-ctx.Done():
				log.Printf("Done.")
				os.Exit(1)
		}
	}()

	// Always cancel child go routines
	defer func() {
		cancel()
	}()

	_ = api_client.CreateAivenClient()

	// Run log event syncer
	go func() {
		err := run(ctx)
		if err != nil {
			log.Fatal("Runner error: ")
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func run(ctx context.Context) error {
	tick, _ := time.ParseDuration("5s")
	for  {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(tick):
			log.Println("Ticking along...")
		}
	}
}



