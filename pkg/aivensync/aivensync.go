package aivensync

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type AivenSync struct {
	audit          *AuditLog
	lastAckedEvent map[string]*AivenEvent
	client         AivenClient
}

func NewAivenSync(audit *AuditLog, aivenToken string) AivenSync {
	return AivenSync{
		audit:          audit,
		lastAckedEvent: make(map[string]*AivenEvent),
		client:         NewAivenClient(aivenToken),
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

		// Events are ordered by
		i := FindLastAckedEvent(events, as.lastAckedEvent[project.ProjectName])
		for ; i >= 0; i-- {
			log.Infof("(%s): %+v", project, events[i])
/*			err := as.audit.Log(event)
			if err != nil {
				log.Errorf("Failed to log event: %v, err: %s", event, err)
				break
			}
*/
			as.lastAckedEvent[project.ProjectName] = events[i]
		}
	}

	return nil
}

func FindLastAckedEvent(events []*AivenEvent, lastAckedEvent *AivenEvent) int {
	if lastAckedEvent == nil {
		return len(events)-1
	}

	for i, event := range events {
		if event.Equals(lastAckedEvent) {
			return i-1
		}
	}

	return len(events)-1
}
