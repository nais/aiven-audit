# Aiven Audit (Go) 📝🕵️
Transfers project event logs from Aiven API to ArcSight

## TODO
- [x] Impl connect to Aiven API and get logs
- [x] Impl connect to db and upsert logs
- [x] Impl connect to ArcSight and sync logs

## Configuration

| environment variable  | description |
| ------------- | ------------- |
| AIVEN_AUDIT_PAT  | Access token for Aiven API  |
| AIVEN_API_URL  | Aiven API URL  |
| AIVEN_PROJECTS  | Aiven projects to fetch audit logs for  |
| NAIS_DATABASE_AIVENAUDIT_EVENTHASHDB_URL  | Database URL for storing event hashes  |
| AUDIT_LOG_ADDR  | Syslog address to send CEF logs to  |


## How Aiven Audit works
![Sequence diagram](doc/aiven-audit.png)

## Sync loop
0. hash batch
1. Hash message
2. upsert with hash as prim key, and with batch hash as column
3. Fetch rows with batch in question
4. Publish to Arcsight