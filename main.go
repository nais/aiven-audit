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

const (
	continuousSyncInterval = 10 * time.Second
)

func root(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Aiven Audit")
	if err != nil {
		log.Fatal("Could not respond")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go httpd(ctx)
	go syncEvents(ctx)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Infof("Received %s, shutting down...", <-interrupt)
}

func httpd(ctx context.Context) {
	log.Info("Starting HTTP server")

	// TODO

	//http.HandleFunc("/", root)
	//nais.InitNaisHandlers()   // Setup handling of nais paths
	//metrics.SetupPrometheus() // Setup metrics path
}

func syncEvents(ctx context.Context) {
	log.Infof("Starting continuous event sync, initial sync scheduled in %s", continuousSyncInterval)
	cfg := config.FromEnv()

	audit := aivensync.NewAuditLog(cfg.AuditLogAddr, "aiven-audit")
	sync := aivensync.NewAivenSync(&audit, cfg.AivenAPIToken)

	ticker := time.NewTicker(continuousSyncInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := sync.Synchronize()
			if err != nil {
				log.Errorf("synchronize: %s", err)
			}
		}
	}
}
