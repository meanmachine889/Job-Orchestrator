CREATE TABLE job_attempts (
  id UUID PRIMARY KEY,
  job_id UUID REFERENCES jobs(id),
  attempt_number INT NOT NULL,
  status job_status NOT NULL,
  error TEXT,
  started_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ
);