package aivensync

import (
	"fmt"
	"log/slog"

	"github.com/nais/aiven-audit/pkg/metrics"
)

var log = slog.Default().With("subsystem", "app")
var auditLog = slog.Default().With("subsystem", "audit")

type AivenSync struct {
	lastAckedEvent map[string]*AivenEvent
	client         AivenClient
	metrics        *metrics.Metrics
}

func NewAivenSync(aivenToken string, m *metrics.Metrics) AivenSync {
	return AivenSync{
		lastAckedEvent: make(map[string]*AivenEvent),
		client:         NewAivenClient(aivenToken),
		metrics:        m,
	}
}

func (as *AivenSync) Synchronize() error {
	log.Info("syncing")
	projects, err := as.client.GetProjects()
	if err != nil {
		return fmt.Errorf("get projects: %w", err)
	}

	for _, project := range projects {
		log.Info("fetching events", "project", project.ProjectName)
		events, err := as.client.GetProjectEvents(project.ProjectName)
		if err != nil {
			return fmt.Errorf("get project events: %w", err)
		}

		for i := FindStartIndex(events, as.lastAckedEvent[project.ProjectName]); i >= 0; i-- {
			auditLog.Info(events[i].EventDesc,
				"actor", events[i].Actor,
				"eventType", events[i].EventType,
				"projectName", project.ProjectName,
				"serviceName", events[i].ServiceName,
				"time", events[i].Time,
			)
			as.lastAckedEvent[project.ProjectName] = events[i]
			as.metrics.EventLogsSyncCounter.Inc()
		}
	}

	return nil
}

func FindStartIndex(events []*AivenEvent, lastAckedEvent *AivenEvent) int {
	if lastAckedEvent == nil {
		return len(events) - 1
	}

	for i, event := range events {
		if event.Equals(lastAckedEvent) {
			return i - 1
		}
	}

	return len(events) - 1
}
