package aivensync

import (
	"fmt"

	log "github.com/sirupsen/logrus"

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

		for i := FindStartIndex(events, as.lastAckedEvent[project.ProjectName]); i >= 0; i-- {
			log.WithFields(log.Fields{
				"AivenAudit_Actor":       events[i].Actor,
				"AivenAudit_EventType":   events[i].EventType,
				"AivenAudit_ServiceName": events[i].ServiceName,
				"AivenAudit_Time":        events[i].Time,
			}).Info(events[i].EventDesc)
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
