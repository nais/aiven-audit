package aiven_sync

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/aiven/aiven-go-client"
	"log"
)

type AivenSync struct {
	api      *aiven.ProjectsHandler
	ses      *SyncedEventsStore
	audit    *AuditLog
	projects []string
}

func NewAivenSync(api *aiven.ProjectsHandler, ses *SyncedEventsStore, audit *AuditLog, projects []string) AivenSync {
	return AivenSync{
		api,
		ses,
		audit,
		projects,
	}
}

type EventProject interface {
	GetEventLog(project string) ([]*aiven.ProjectEvent, error)
}

func (as *AivenSync) Synchronize(project EventProject) {
	for _, aivenProject := range as.projects {

		log.Printf("processing logs for project = %v", aivenProject)
		eventLog, err := project.GetEventLog(aivenProject)
		if err != nil {
			log.Fatal("Error getting event logs %w", err)
		}

		// create hashes of batch and individual events
		eventLogBatchHash := batchHash(eventLog)
		hashToEvent := hashMapOf(eventLog)

		// upsert event hashes
		var nrRecordsChanged int64
		for hash, _ := range hashToEvent {

			affected, err := as.ses.UpsertLogEvent(hash, eventLogBatchHash)
			if err != nil {
				log.Fatalf("error upserting logs: %v", err)
			}
			nrRecordsChanged += affected
		}

		if nrRecordsChanged == 0 {
			log.Printf("Records changed after upserting batch = %v is 0, no new events, aborting sync for batch", eventLogBatchHash)
			return
		}

		// get upnserted events with current batch hash
		newEventHashes, err := as.ses.LogEventHashesForBatch(eventLogBatchHash)
		if err != nil {
			log.Fatal("Error querying events for batch: ", err)
		}

		// get the new, unique events
		newEvents := filterNewEvents(newEventHashes, hashToEvent)

		// log events to ArchSight
		log.Printf("events to be logged to Archsight: %v\n", len(newEvents))
		_, err = as.audit.Log(newEvents)
		if err != nil {
			log.Printf("Error while logging to naudit, rolling back upserts: %v", err)
			as.ses.Rollback()
		} else {
			as.ses.Commit()
		}
	}
}

func hashMapOf(events []*aiven.ProjectEvent) map[string]*aiven.ProjectEvent {
	hashToEvent := make(map[string]*aiven.ProjectEvent) // make mapping for filtering
	for _, event := range events {
		hash := eventHash(event)
		hashToEvent[string(hash)] = event // store for filtering later
	}
	return hashToEvent
}

func filterNewEvents(newEventHashes []string, hashToEvent map[string]*aiven.ProjectEvent) []*aiven.ProjectEvent {
	newEvents := make([]*aiven.ProjectEvent, 0)
	for _, hash := range newEventHashes {
		event := hashToEvent[hash]
		if event != nil {
			newEvents = append(newEvents, event)
		}
	}
	return newEvents
}

func eventHash(event *aiven.ProjectEvent) string {
	hasher := sha256.New()
	hasher.Write([]byte(event.EventType))
	hasher.Write([]byte(event.EventDesc))
	hasher.Write([]byte(event.Time))
	hasher.Write([]byte(event.Actor))
	hasher.Write([]byte(event.ServiceName))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func batchHash(events []*aiven.ProjectEvent) string {
	hasher := sha256.New()
	for _, event := range events {
		hasher.Write([]byte(eventHash(event)))
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
