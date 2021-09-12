package aiven_sync

const DbSchema = "" +
	"CREATE TABLE IF NOT EXISTS consumed_aiven_events (" +
	"id  serial," +
	"event_hash integer NOT NULL," +
	"PRIMARY KEY (event_hash))"
