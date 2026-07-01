-- Development seed data.
-- Inserts a minimal set of users and cases for local development and manual
-- testing. All inserts use ON CONFLICT DO NOTHING — idempotent on re-run.
-- Do NOT add production data here; use environment-specific seed scripts instead.

-- Seed users (passwords / auth not modelled in this boilerplate layer).
INSERT INTO users (name, email, created_at, updated_at)
VALUES
    ('Alice Mercer',   'alice@example.com',   now(), now()),
    ('Bob Harrington', 'bob@example.com',     now(), now()),
    ('Carol Dean',     'carol@example.com',   now(), now())
ON CONFLICT (email) DO NOTHING;

-- Seed cases assigned to the seed users above.
-- Covers a spread of lifecycle statuses for local UI / API testing.
INSERT INTO cases (file_number, status, servicer_id, assignee_id, created_at)
SELECT 'FC-2024-001', 'open',             1, id, now() FROM users WHERE email = 'alice@example.com'
ON CONFLICT (file_number) DO NOTHING;

INSERT INTO cases (file_number, status, servicer_id, assignee_id, created_at)
SELECT 'FC-2024-002', 'title_review',     1, id, now() FROM users WHERE email = 'bob@example.com'
ON CONFLICT (file_number) DO NOTHING;

INSERT INTO cases (file_number, status, servicer_id, assignee_id, created_at)
SELECT 'FC-2024-003', 'pending_judgment', 2, id, now() FROM users WHERE email = 'carol@example.com'
ON CONFLICT (file_number) DO NOTHING;

INSERT INTO cases (file_number, status, servicer_id, assignee_id, created_at)
SELECT 'FC-2024-004', 'closed',           2, id, now() FROM users WHERE email = 'alice@example.com'
ON CONFLICT (file_number) DO NOTHING;
