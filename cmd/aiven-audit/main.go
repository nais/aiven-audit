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
	"github.com/nais/aiven-audit/pkg/metrics"
	"github.com/nais/aiven-audit/pkg/nais"
)

const (
	continuousSyncInterval = 10 * time.Second
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go syncEvents(ctx)
	go httpd()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Infof("Received %s, shutting down", <-interrupt)
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
				log.Errorf("Continuous sync: %s", err)
			}
		}
	}
}

func httpd() {
	log.Info("Starting HTTP server")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(writer, "Aiven audit log sync")
	})

	nais.Handlers(mux)
	metrics.Handlers(mux)

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	defer shutdownHttpd(srv)

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Errorf("Serve HTTP: %v", err)
	} else {
		log.Info("HTTP server closed")
	}
}

func shutdownHttpd(srv http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := srv.Shutdown(ctx)
	if err != nil {
		log.Errorf("Shutdown HTTP server: %v", err)
	}

	err = ctx.Err()
	if err != nil {
		log.Errorf("Shutdown HTTP server context error: %v", err)
	}
}
