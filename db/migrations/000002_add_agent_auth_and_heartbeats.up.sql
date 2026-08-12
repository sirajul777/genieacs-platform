ALTER TABLE agents
    ADD COLUMN token_hash BYTEA,
    ADD COLUMN last_seen TIMESTAMPTZ;

CREATE TABLE heartbeats (
    id BIGSERIAL PRIMARY KEY,
    agent_id UUID NOT NULL REFERENCES agents (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX heartbeats_agent_id_created_at_idx ON heartbeats (agent_id, created_at DESC);
