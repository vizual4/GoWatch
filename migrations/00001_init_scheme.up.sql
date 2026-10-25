-- +goose up

CREATE TABLE IF NOT EXISTS monitors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    interval_sec INT NOT NULL DEFAULT 60,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    next_check_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_next_check ON monitors(next_check_at);

CREATE TABLE IF NOT EXISTS checks (
    id SERIAL PRIMARY KEY,
    monitor_id INT NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
    status_code INT NOT NULL,
    response_time_ms INT NOT NULL,
    error_msg TEXT,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_checks_monitor_id ON checks(monitor_id);

-- +goose Down
DROP TABLE IF EXISTS checks;
DROP TABLE IF EXISTS monitors;