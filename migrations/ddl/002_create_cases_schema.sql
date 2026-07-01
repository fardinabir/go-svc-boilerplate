-- Cases table DDL hardening.
-- AutoMigrate creates the table and uniqueIndex on file_number; this file adds
-- FK constraints, check constraints, and indexes that GORM does not generate.
-- All statements are idempotent — safe to re-run on existing databases.

-- FK: assignee_id → users(id).
-- RESTRICT on delete: a user with active cases cannot be deleted without first
-- reassigning or closing those cases. RESTRICT is intentional — silent SET NULL
-- would leave cases with no owner, breaking downstream milestone workflows.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'fk_cases_assignee' AND conrelid = 'cases'::regclass
    ) THEN
        ALTER TABLE cases
            ADD CONSTRAINT fk_cases_assignee
            FOREIGN KEY (assignee_id)
            REFERENCES users (id)
            ON UPDATE CASCADE
            ON DELETE RESTRICT;
    END IF;
END$$;

-- Index on assignee_id for FK lookup and "cases assigned to user X" queries.
CREATE INDEX IF NOT EXISTS idx_cases_assignee_id ON cases (assignee_id);

-- Index on status for filtering by lifecycle stage.
CREATE INDEX IF NOT EXISTS idx_cases_status ON cases (status);

-- Index on servicer_id for filtering cases by servicer.
CREATE INDEX IF NOT EXISTS idx_cases_servicer_id ON cases (servicer_id);

-- Index on created_at for the default ORDER BY on list queries.
CREATE INDEX IF NOT EXISTS idx_cases_created_at ON cases (created_at DESC);

-- Constrain status to known foreclosure lifecycle values.
-- Extend this list as new workflow stages are introduced; keep it additive.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_cases_status' AND conrelid = 'cases'::regclass
    ) THEN
        ALTER TABLE cases
            ADD CONSTRAINT chk_cases_status
            CHECK (status IN (
                'open',
                'pending_title',
                'title_review',
                'pending_judgment',
                'judgment_entered',
                'sale_scheduled',
                'closed',
                'dismissed'
            ));
    END IF;
END$$;

-- Enforce non-empty file number (GORM not-null only blocks NULL, not '').
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'chk_cases_file_number_nonempty' AND conrelid = 'cases'::regclass
    ) THEN
        ALTER TABLE cases
            ADD CONSTRAINT chk_cases_file_number_nonempty
            CHECK (char_length(trim(file_number)) > 0);
    END IF;
END$$;

COMMENT ON TABLE  cases              IS 'Foreclosure cases — cases domain aggregate.';
COMMENT ON COLUMN cases.id           IS 'Surrogate primary key.';
COMMENT ON COLUMN cases.file_number  IS 'Unique case file reference (e.g. FC-2024-001).';
COMMENT ON COLUMN cases.status       IS 'Lifecycle stage. Constrained to known foreclosure workflow values.';
COMMENT ON COLUMN cases.servicer_id  IS 'Reference to the loan servicer (reference domain, not FK yet).';
COMMENT ON COLUMN cases.assignee_id  IS 'User responsible for this case. FK → users(id).';
COMMENT ON COLUMN cases.created_at   IS 'Row creation timestamp (UTC).';
