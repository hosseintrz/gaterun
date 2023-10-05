CREATE TABLE IF NOT EXISTS consumers(
    id UUID PRIMARY KEY,
    username TEXT UNIQUE,
    custom_id TEXT UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP(0) AT TIME ZONE 'UTC')
);

CREATE INDEX IF NOT EXISTS "consumers_username_idx" ON "consumers" (LOWER("username"));