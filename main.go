package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nais/aiven-audit/pkg/aivensync"
	"github.com/nais/aiven-audit/pkg/config"
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
	//http.HandleFunc("/", root)
	log.Println("starting aiven-audit...")
	//nais.InitNaisHandlers()   // Setup handling of nais paths
	//metrics.SetupPrometheus() // Setup metrics path

	programContext, cancel := context.WithCancel(context.Background())

	// Handle common signals
	handleCommonSignals(programContext, cancel)

	// Always cancel child go routines
	defer cancel()

	cfg := config.FromEnv()

	audit := aivensync.NewAuditLog(cfg.AuditLogAddr, "aiven-audit")
	aivenSync := aivensync.NewAivenSync(&audit, cfg.AivenAPIToken)
	synchronizeContinuously(programContext, aivenSync)

	return nil
	//return http.ListenAndServe(":8080", nil)
}

func synchronizeContinuously(ctx context.Context, s aivensync.AivenSync) {
	log.Infof("Starting synchronizer, first sync scheduled in 10s")
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Info("Program done, shutting down synchronizer")
			return
		case <-ticker.C:
			err := s.Synchronize()
			if err != nil {
				log.Errorf("synchronize: %s", err)
			}
		}
	}
}

func handleCommonSignals(programContext context.Context, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-signalChan:
			log.Printf("aiven-audit got SIGINT/SIGTERM, exiting")
			cancel()
		case <-programContext.Done():
			log.Printf("aiven-audit done")
		}
	}()
}
