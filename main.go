package main

import (
	"context"
	"fmt"
	"github.com/aiven/aiven-go-client"
	"github.com/joho/godotenv"
	eventsync "github.com/nais/aiven-audit/pkg/aiven_sync"
	config2 "github.com/nais/aiven-audit/pkg/config"
	"github.com/nais/aiven-audit/pkg/metrics"
	"github.com/nais/aiven-audit/pkg/nais"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func root(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Aiven Audit")
	if err != nil {
		log.Fatal("Could not respond")
	}
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	http.HandleFunc("/", root)
	log.Println("starting aiven-audit...")
	nais.InitNaisHandlers()   // Setup handling of nais paths
	metrics.SetupPrometheus() // Setup metrics path

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Handle common signals
	handleCommonSignals(ctx, cancel)

	// Always cancel child go routines
	defer func() {
		cancel()
	}()

	config := config2.ConfigFromEnv()
	apiClient, err := aiven.NewTokenClient(config.AivenAPIToken, "")
	if err != nil {
		log.Fatalf("Could not create Aiven Client - err: %s", err)
	}

	ses := eventsync.NewSyncedEventsStore(config.DbHost)
	ses.Init()

	audit := eventsync.NewAuditLog(config.AuditLogAddr, "aiven-audit")
	aivenSync := eventsync.NewAivenSync(apiClient.Projects, &ses, &audit, config.Projects)

	// Run log event syncer
	go func() {
		err := syncLogEvents(ctx, apiClient, &aivenSync)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return http.ListenAndServe(":8080", nil)
}

func syncLogEvents(ctx context.Context, client *aiven.Client, aivernSyncer *eventsync.AivenSync) error {
	tick, _ := time.ParseDuration("5s")
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(tick):
			//log.Println("Ticking along...")
			aivernSyncer.Synchronize(client.Projects)
			metrics.EventLogsSyncCounter.Inc()
		}
	}
}

func handleCommonSignals(ctx context.Context, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-signalChan:
			log.Printf("aiven-audit got SIGINT/SIGTERM, exiting")
			cancel()
			os.Exit(1)
		case <-ctx.Done():
			log.Printf("aiven-audit done")
			os.Exit(1)
		}
	}()
}
