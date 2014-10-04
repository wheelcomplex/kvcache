# kvcache

kvcache implements a simple key/value datastore for a particular use case:

- Keys and values are just bytes
- Keys are small (on the order of 10 bytes) and values are large (kilobytes)
- Keys and values only need to be accessed for a fixed time window.
- Lookups are very much more frequent for more recent keys than for older keys

General implementation notes:

- Key/value pairs have a fixed expiration duration
- Key/value pairs are immutable, once written (until expired)
- Key/value pairs are stored on disk in a rotating set of fixed-size append-only logs
- The keys are associated with offsets into the log by an in-memory hashtable (can be reconstructed from the
  logs)
- Recent key/value pairs are duplicated in another in-memory hashtable (with the complete value) for fast
  lookup

## Speedup

- Slow initial loading when DB is large
  - Profile + optimize
  - Parallelize

- General
  - Optimistically serialize values before grabbing mutex (most writes will not be collisions)
  - See the space/latency difference of switching back to FLATE using some real test data
