CREATE TABLE IF NOT EXISTS jobs (
    id VARCHAR(36) PRIMARY KEY,
	once VARCHAR(255) NULL,
	interval VARCHAR(50) NULL,
	status VARCHAR(20) NOT NULL,
	createdAt BIGINT NOT NULL,
	lastFinishedAt BIGINT NOT NULL,
	payload JSONB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(createdAt);