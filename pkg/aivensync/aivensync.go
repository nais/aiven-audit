package aivensync

import (
	"fmt"
	"log/slog"

	"github.com/nais/aiven-audit/pkg/metrics"
)

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
	slog.Info("syncing")
	projects, err := as.client.GetProjects()
	if err != nil {
		return fmt.Errorf("get projects: %w", err)
	}

	for _, project := range projects {
		slog.Info("fetching events", "project", project.ProjectName)
		events, err := as.client.GetProjectEvents(project.ProjectName)
		if err != nil {
			return fmt.Errorf("get project events: %w", err)
		}

		for i := FindStartIndex(events, as.lastAckedEvent[project.ProjectName]); i >= 0; i-- {
			slog.Info(events[i].EventDesc,
				"AivenAudit_Actor", events[i].Actor,
				"AivenAudit_EventType", events[i].EventType,
				"AivenAudit_ProjectName", project.ProjectName,
				"AivenAudit_ServiceName", events[i].ServiceName,
				"AivenAudit_Time", events[i].Time,
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
