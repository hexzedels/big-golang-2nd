-- +goose Up
-- +goose StatementBegin
CREATE TABLE executions (
    id VARCHAR(255) PRIMARY KEY,
    job_id VARCHAR(255) NOT NULl,
    worker_id VARCHAR(255) NULL,
    status VARCHAR(50) NOT NULL,
    queued_at BIGINT NOT NULL,
    started_at BIGINT,
    finished_at BIGINT,
    error JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS executions;
-- +goose StatementEnd
    