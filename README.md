# Aiven Audit (Go) 📝🕵️
Transfers project event logs from Aiven API to ArcSight

## TODO
- [ ] Impl connect to Aiven API and get logs
- [ ] Impl connect to db and upsert logs
- [ ] Impl connect to ArcSight and sync logs


## Sync loop
0. Get latest row fro db and compare etag to HTTP GETed document.

   if new etag

1. Hash message
2. upsert with hash as prim key, and with etag as column
3. Fetch rows with etag in question
4. Publish to Arcsight