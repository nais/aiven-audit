package aiven_sync

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/aiven/aiven-go-client"
	"log"
)

type AivenSync struct {
	api   *aiven.ProjectsHandler
	ses   *SyncedEventsStore
	audit *AuditLog
}

func NewAivenSync(api *aiven.ProjectsHandler, ses *SyncedEventsStore, audit *AuditLog) AivenSync {
	return AivenSync{
		api,
		ses,
		audit,
	}
}

type EventProject interface {
	GetEventLog(project string) ([]*aiven.ProjectEvent, error)
}

func (as *AivenSync) Synchronize(project EventProject) {
	eventLog, err := project.GetEventLog("nav-dev") // TODO make configurable
	if err != nil {
		log.Fatal("Error getting event logs %w", err)
	}

	// upsert received events
	eventLogBatchHash := batchHash(eventLog)
	hashToEvent := make(map[string]*aiven.ProjectEvent) // make mapping for filtering

	var nrRecordsChanged int64
	for _, event := range eventLog {
		hash := eventHash(event)
		hashToEvent[string(hash)] = event // store for filtering later

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

	// fetch inserted events with batch hash
	newEventHashes, err := as.ses.LogEventHashesForBatch(eventLogBatchHash)
	if err != nil {
		log.Fatal("Error querying events for batch: ", err)
	}

	newEvents := make([]*aiven.ProjectEvent, 0)
	for _, hash := range newEventHashes {
		event := hashToEvent[hash]
		if event != nil {
			newEvents = append(newEvents, event)
		}
	}

	// push events to ArchSight
	log.Printf("events to be logged to Archsight: %v\n", len(newEvents))
	as.audit.Log(newEvents) // TODO handle error when logging to naudit
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
