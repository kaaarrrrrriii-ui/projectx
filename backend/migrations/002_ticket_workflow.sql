BEGIN;

-- API timestamps are UTC instants. Existing timestamp values are interpreted
-- as UTC while converting the columns to timezone-aware PostgreSQL values.
ALTER TABLE tickets
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC',
    ALTER COLUMN closed_at TYPE TIMESTAMPTZ USING closed_at AT TIME ZONE 'UTC';

ALTER TABLE messages
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';

ALTER TABLE attachments
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';

ALTER TABLE tickets_workers
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';

ALTER TABLE crisis_contact
    ALTER COLUMN created_at TYPE TIMESTAMPTZ USING created_at AT TIME ZONE 'UTC';

ALTER TABLE tickets
    ADD CONSTRAINT chk_tickets_status
        CHECK (status BETWEEN 1 AND 9);

ALTER TABLE messages
    ADD CONSTRAINT chk_messages_type
        CHECK (type BETWEEN 1 AND 4);

COMMIT;
