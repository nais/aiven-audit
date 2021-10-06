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
	err := run()
	if err != nil {
		log.Fatalf("fatal: %s", err)
		os.Exit(1)
	}
}

func run() error {
	programContext, cancel := context.WithCancel(context.Background())

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		log.Infof("Received %s, shutting down", <-interrupt)
		cancel()
	}()

	m := metrics.SetupMetrics()
	cfg, err := config.New()
	if err != nil {
		return err
	}

	audit := aivensync.NewAuditLog(cfg.AuditLogAddress, "aiven-audit")
	aivenSync := aivensync.NewAivenSync(&audit, cfg.AivenAPIToken, m)

	go syncEvents(programContext, aivenSync)

	return httpd(programContext)
}

func syncEvents(ctx context.Context, aivenSync aivensync.AivenSync) {
	log.Infof("Starting continuous event sync, initial sync scheduled in %s", continuousSyncInterval)

	ticker := time.NewTicker(continuousSyncInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := aivenSync.Synchronize()
			if err != nil {
				log.Errorf("Continuous sync: %s", err)
			}
		}
	}
}

func httpd(ctx context.Context) error {
	log.Info("Starting HTTP server")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(writer, "Aiven audit log sync")
	})

	nais.Handlers(mux)
	metrics.Handlers(mux)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			log.Errorf("Shutdown HTTP server: %v", err)
		}
	}()

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Serve HTTP: %v", err)
	} else {
		log.Info("HTTP server closed")
	}

	return nil
}
