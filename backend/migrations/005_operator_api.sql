BEGIN;

-- The product specification defines standard, urgent and low priorities.
ALTER TABLE tickets
    DROP CONSTRAINT IF EXISTS chk_tickets_priority,
    ADD CONSTRAINT chk_tickets_priority
        CHECK (priority BETWEEN 1 AND 3);

COMMIT;
