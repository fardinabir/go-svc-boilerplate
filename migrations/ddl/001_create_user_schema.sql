-- Users table DDL hardening.
-- AutoMigrate creates the table and uniqueIndex on email; this file adds
-- constraints and indexes that GORM does not generate automatically.
-- All statements are idempotent — safe to re-run on existing databases.

-- Index on created_at for the default ORDER BY on list queries.
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users (created_at DESC);

-- Enforce minimum name length at the DB level (app-layer validation is not
-- the last line of defence against direct DB writes or migration scripts).
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_users_name_length' AND conrelid = 'users'::regclass
    ) THEN
        ALTER TABLE users
            ADD CONSTRAINT chk_users_name_length
            CHECK (char_length(trim(name)) BETWEEN 2 AND 50);
    END IF;
END$$;

-- Enforce lowercase email storage.
-- Prevents duplicate registrations that differ only in case.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_users_email_lowercase' AND conrelid = 'users'::regclass
    ) THEN
        ALTER TABLE users
            ADD CONSTRAINT chk_users_email_lowercase
            CHECK (email = lower(email));
    END IF;
END$$;

COMMENT ON TABLE  users            IS 'Application users — identity domain aggregate.';
COMMENT ON COLUMN users.id         IS 'Surrogate primary key.';
COMMENT ON COLUMN users.name       IS 'Display name. 2–50 chars, letters/spaces/hyphens/apostrophes.';
COMMENT ON COLUMN users.email      IS 'Unique contact email. Stored lowercase.';
COMMENT ON COLUMN users.created_at IS 'Row creation timestamp (UTC).';
COMMENT ON COLUMN users.updated_at IS 'Last update timestamp (UTC).';
