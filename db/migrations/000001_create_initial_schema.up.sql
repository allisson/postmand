-- webhooks table

CREATE TABLE IF NOT EXISTS webhooks(
   id UUID PRIMARY KEY,
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

-- deliveries table

CREATE TABLE IF NOT EXISTS deliveries(
   id UUID PRIMARY KEY,
   webhook_id UUID NOT NULL,
   payload TEXT NOT NULL,
   scheduled_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   delivery_attempts SMALLINT NOT NULL,
   status VARCHAR NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   FOREIGN KEY (webhook_id) REFERENCES webhooks (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS deliveries_webhook_id_idx ON deliveries (webhook_id);
CREATE INDEX IF NOT EXISTS deliveries_status_idx ON deliveries (status);
CREATE INDEX IF NOT EXISTS deliveries_scheduled_at_idx ON deliveries USING BRIN(scheduled_at);
CREATE INDEX IF NOT EXISTS deliveries_created_at_idx ON deliveries USING BRIN(created_at);
