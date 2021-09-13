package aiven_sync

const DbSchema = "" +
	"CREATE TABLE IF NOT EXISTS consumed_aiven_events (" +
	"id serial, " +
	"event_hash varchar(64) NOT NULL," +
	"batch_hash varchar(64) NOT NULL," +
	"PRIMARY KEY (event_hash)" +
	")"
