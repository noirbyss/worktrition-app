CREATE TABLE IF NOT EXISTS ai_generation_requests (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    plan_type VARCHAR(32) NOT NULL,
    status VARCHAR(16) NOT NULL,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ai_prompt_versions (
    id BIGSERIAL PRIMARY KEY,
    generation_id VARCHAR(36) NOT NULL REFERENCES ai_generation_requests(id),
    system_prompt TEXT NOT NULL,
    user_prompt TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS ai_generation_results (
    id BIGSERIAL PRIMARY KEY,
    generation_id VARCHAR(36) NOT NULL REFERENCES ai_generation_requests(id),
    raw_response TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS generated_plan_versions (
    id BIGSERIAL PRIMARY KEY,
    generation_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    plan_type VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ai_generation_requests_user_id ON ai_generation_requests(user_id);
