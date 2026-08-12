DROP TABLE IF EXISTS heartbeats;

ALTER TABLE agents
    DROP COLUMN IF EXISTS last_seen,
    DROP COLUMN IF EXISTS token_hash;
