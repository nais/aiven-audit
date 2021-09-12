package aiven_sync

import (
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	log "github.com/sirupsen/logrus"
)

type SyncedEventsStore struct {
	db *sql.DB
}

func NewSyncedEventsStore(dbUrl string) SyncedEventsStore {
	db, err := sql.Open("pgx",
		dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	return SyncedEventsStore{
		db: db,
	}
}

func (ses SyncedEventsStore) Init() {
	_, err := ses.db.Exec(DbSchema)
	if err != nil {
		log.Fatalf("Error while initializing db: %v", err)
	}
}

func (ses SyncedEventsStore) Close() {
	err := ses.db.Close()
	if err != nil {
		log.Errorf("Error while closing db: %v", err)
	}
}

func (ses SyncedEventsStore) UpsertLogEvent(eventHash, batchHash string) (int64, error) {
	query := "INSERT INTO consumed_aiven_events (event_hash, batch_hash) " +
		"VALUES($1, $2) " +
		"ON CONFLICT (event_hash) " +
		"DO NOTHING;"

	result, err := ses.db.Exec(query, eventHash, batchHash)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (ses SyncedEventsStore) LogEventHashesForBatch(batchHash string) ([]string, error) {
	rows, err := ses.db.Query(
		"SELECT event_hash FROM consumed_aiven_events "+
			"WHERE batch_hash = $1;", batchHash)
	if err != nil {
		return nil, err
	}

	var eventHashesForBatch []string

	for rows.Next() {
		var eventHash string
		if err := rows.Scan(&eventHash); err != nil {
			return eventHashesForBatch, err
		}
		eventHashesForBatch = append(eventHashesForBatch, eventHash)
	}
	return eventHashesForBatch, nil
}
