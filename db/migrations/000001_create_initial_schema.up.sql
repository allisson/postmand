-- webhooks table

CREATE TABLE IF NOT EXISTS webhooks(
   id uuid PRIMARY KEY,
   name VARCHAR NOT NULL,
   url VARCHAR NOT NULL,
   content_type VARCHAR NOT NULL,
   valid_status_codes SMALLINT[] NOT NULL,
   secret_token VARCHAR NOT NULL,
   active BOOLEAN NOT NULL,
   max_delivery_attempts SMALLINT NOT NULL,
   delivery_attempt_timeout SMALLINT NOT NULL,
   retry_min_backoff SMALLINT NOT NULL,
   retry_max_backoff SMALLINT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS webhooks_name_idx ON webhooks (name);
CREATE INDEX IF NOT EXISTS webhooks_active_idx ON webhooks (active);
CREATE INDEX IF NOT EXISTS webhooks_created_at_idx ON webhooks USING BRIN(created_at);
