package aivensync

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/nais/aiven-audit/pkg/metrics"
)

type AivenSync struct {
	audit          *AuditLog
	lastAckedEvent map[string]*AivenEvent
	client         AivenClient
	metrics        *metrics.Metrics
}

func NewAivenSync(audit *AuditLog, aivenToken string, m *metrics.Metrics) AivenSync {
	return AivenSync{
		audit:          audit,
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
		log.Infof("fetching events for: %v", project.ProjectName)
		events, err := as.client.GetProjectEvents(project.ProjectName)
		if err != nil {
			return fmt.Errorf("get project events: %w", err)
		}

		for i := FindStartIndex(events, as.lastAckedEvent[project.ProjectName]) ; i >= 0; i-- {
			log.Infof("(%s): %+v", project, events[i])
			err := as.audit.Log(events[i])

			if err != nil {
				as.metrics.EventLogsFailedSyncCounter.Inc()
				log.Errorf("Failed to log event: %v, err: %s", events[i], err)
				break
			}

			as.metrics.EventLogsSyncCounter.Inc()
			as.lastAckedEvent[project.ProjectName] = events[i]
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
