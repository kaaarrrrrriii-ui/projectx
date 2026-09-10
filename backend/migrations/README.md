# Database initialization

`001_init.sql` is the executable PostgreSQL schema based on the supplied UML.
Docker Compose mounts it into `/docker-entrypoint-initdb.d`, so PostgreSQL runs
it automatically when the `pgdata` volume is created for the first time.

Two UML inconsistencies are normalized so the foreign keys are valid:

- `tickets.id` is the internal integer primary key used by every `ticket_id`;
  `tickets.track_id` remains a unique `VARCHAR(64)` identifier.
- `tickets_workers.ticker_id` is treated as the intended `ticket_id`.

The UML leaves the attachment metadata as a placeholder. The `attachments`
table therefore includes the storage key, safe name, MIME type, byte sizes,
image dimensions, and creation time required by the existing private image
storage code.

The script does not rewrite an already initialized database. Apply future
schema changes as new numbered migrations instead of editing an applied file.

`002_ticket_workflow.sql` assigns the integer status/message-type ranges used
by the ticket API and converts stored timestamps to timezone-aware values. For
an existing Docker volume, apply it with the project's normal migration
process; PostgreSQL's initialization directory only runs for a new volume.

`003_operator_workflow.sql` fixes the two allowed priority values, the three
applicant types, adds operator/return message types, and creates the event log
used for assignment history and timing analytics.

`005_operator_api.sql` expands the existing priority constraint for the third
product priority, `low`. It does not add tables or columns.

For an already initialized Docker volume, apply the new migration explicitly:

```sh
docker compose exec -T postgres psql -U appuser -d appdb \
  -f /docker-entrypoint-initdb.d/003_operator_workflow.sql
```

Apply `005_operator_api.sql` through the same migration mechanism before using
the `low` priority.
